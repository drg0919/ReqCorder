package initiator

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptrace"
	"os"
	"strings"
	"time"

	"reqcorder/internal/request"
	"reqcorder/internal/response"
)

// Perform HTTP request based on the request configuration.
func InitiateRequest(r *request.RequestObject) (*response.ResponseObject, error) {
	slog.Debug("Initiating HTTP request", "request", r)
	var client *http.Client
	if (*r.SSLVerify) && r.CACertPath != "" {
		slog.Debug("Configuring TLS with custom CA certificate", "caCertPath", r.CACertPath)
		pool, _ := x509.SystemCertPool()
		if pool == nil {
			slog.Debug("System certificate pool not available, creating new pool")
			pool = x509.NewCertPool()
		}

		pem, err := os.ReadFile(r.CACertPath)
		if err != nil {
			slog.Error("Failed to read CA certificate file", "caCertPath", r.CACertPath, "error", err)
			return nil, ErrorFailedToReadCert
		}
		if ok := pool.AppendCertsFromPEM(pem); !ok {
			slog.Error("Failed to append CA certificate to pool", "caCertPath", r.CACertPath)
			return nil, fmt.Errorf("%w %q", ErrorFailedToReadCert, r.CACertPath)
		}

		tlsConfig := &tls.Config{
			RootCAs:    pool,
			MinVersion: tls.VersionTLS12,
		}

		transport := &http.Transport{TLSClientConfig: tlsConfig}
		client = &http.Client{
			Timeout:   r.Timeout,
			Jar:       r.CookieJar,
			Transport: transport,
		}
		slog.Debug("HTTP client configured with custom TLS settings", "timeout", r.Timeout)
	} else {
		slog.Debug("Configuring HTTP client with default TLS settings", "sslVerify", *r.SSLVerify)
		client = &http.Client{
			Timeout: r.Timeout,
			Jar:     r.CookieJar,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: !(*r.SSLVerify),
				},
			},
		}
		slog.Debug("HTTP client configured", "timeout", r.Timeout, "insecureSkipVerify", !(*r.SSLVerify))
	}

	slog.Debug("Constructing HTTP request", "method", r.Method, "url", r.URL)
	req, err := constructHTTPRequest(r)
	if err != nil {
		slog.Error("Failed to construct HTTP request", "error", err)
		return nil, err
	}
	slog.Debug("HTTP request constructed successfully", "method", r.Method, "url", r.URL)

	var timing response.ResponseTimes
	trace := createTrace(&timing)
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))

	slog.Info("Executing HTTP request", "method", r.Method, "url", r.URL, "timeout", r.Timeout)
	start := time.Now().UTC()
	res, err := client.Do(req)
	timing.Total = time.Since(start)

	if err != nil {
		slog.Error("HTTP request failed to execute", "error", err)

		return &response.ResponseObject{
			StatusCode: 1000,
			Body:       "This request failed: " + err.Error() + "\n Refer to the logs for more details",
			Timing:     timing,
		}, fmt.Errorf("%w: %v", ErrorRequestFailed, err)
	}

	slog.Info("HTTP request completed", "method", r.Method, "url", r.URL, "statusCode", res.StatusCode, "duration", timing.Total)

	slog.Debug("Reading response body", "contentLength", res.Header.Get("Content-Length"), "contentType", res.Header.Get("Content-Type"))
	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()

	if err != nil {
		slog.Error("Failed to read response body", "error", err, "statusCode", res.StatusCode)
		return nil, fmt.Errorf("%w: %v", ErrorFailedToReadResponse, err)
	}

	slog.Debug("Response body read successfully", "bodySize", len(body), "statusCode", res.StatusCode)

	return &response.ResponseObject{
		StatusCode: res.StatusCode,
		Headers:    flattenHeaders(res.Header),
		Body:       string(body),
		Size:       int64(len(body)),
		Timing:     timing,
		Cookies:    res.Cookies(),
	}, nil
}

// Construct request related items before initiating it.
func constructHTTPRequest(r *request.RequestObject) (*http.Request, error) {
	var bodyReader io.Reader
	if r.Body != "" {
		slog.Debug("Reading request body")
		bodyReader = strings.NewReader(r.Body)
	}
	req, err := http.NewRequest(r.Method, r.URL, bodyReader)
	if err != nil {
		slog.Error("Failed to build HTTP request", "error", err)
		return nil, fmt.Errorf("%w: %v", ErrorFailedToBuildRequest, err)
	}

	for key, value := range r.Headers {
		slog.Debug("Processing header", "header", key, "headerValue", value)
		req.Header.Set(key, value)
	}

	if r.Auth != "" {
		headerName := r.AuthHeaderName
		if headerName == "" {
			slog.Debug("Defaulting to Authorization header")
			headerName = "Authorization"
		}
		req.Header.Set(headerName, r.Auth)
	}

	if r.UserAgent != "" {
		slog.Debug("Setting user agent", "userAgent", r.UserAgent)
		req.Header.Set("User-Agent", r.UserAgent)
	}

	return req, nil
}

// Capture timing stats for response.
func createTrace(timing *response.ResponseTimes) *httptrace.ClientTrace {
	return &httptrace.ClientTrace{
		DNSStart: func(_ httptrace.DNSStartInfo) {
			timing.DNSStart = time.Now()
		},
		DNSDone: func(_ httptrace.DNSDoneInfo) {
			timing.DNSDone = time.Now()
			timing.DNSLookup = timing.DNSDone.Sub(timing.DNSStart)
		},
		ConnectStart: func(_, _ string) {
			timing.ConnectStart = time.Now()
		},
		ConnectDone: func(_, _ string, _ error) {
			timing.ConnectDone = time.Now()
			timing.TCPConnect = timing.ConnectDone.Sub(timing.ConnectStart)
		},
		TLSHandshakeStart: func() {
			timing.TLSHandshakeStart = time.Now()
		},
		TLSHandshakeDone: func(_ tls.ConnectionState, _ error) {
			timing.TLSHandshakeDone = time.Now()
			timing.TLSHandshake = timing.TLSHandshakeDone.Sub(timing.TLSHandshakeStart)
		},
		GotFirstResponseByte: func() {
			timing.GotFirstResponseByte = time.Now()
			timing.FirstByte = timing.GotFirstResponseByte.Sub(timing.TLSHandshakeDone)
		},
	}
}

// Convert headers to map.
func flattenHeaders(headers http.Header) map[string]string {
	flat := make(map[string]string)
	for key, values := range headers {
		flat[key] = strings.Join(values, ", ")
	}
	return flat
}
