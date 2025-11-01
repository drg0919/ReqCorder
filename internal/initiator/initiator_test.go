package initiator

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/http/httptrace"
	"os"
	"path/filepath"
	"reqcorder/internal/request"
	"reqcorder/internal/response"
	"testing"
)

var (
	illegalMethod  = "//"
	externalDomain = "https://www.github.com"
	illegalDomain  = "https://this-does-not-exist-at-all.ai"
)

func TestSuccessfulFlattenHeaders(t *testing.T) {
	req, _ := http.NewRequest("GET", "localhost", nil)
	req.Header.Set("head1", "one")
	req.Header.Set("head2", "two")
	headers := flattenHeaders(req.Header)
	expectedLen := 2
	if len(headers) != expectedLen {
		t.Errorf("Expected %d headers, received %d\n", expectedLen, len(headers))
	}
	for key, val := range headers {
		if key == "head1" && val != "one" {
			t.Errorf("Expected header %s to have value %s, received %s\n", "head1", "one", val)
		}
		if key == "head2" && val != "two" {
			t.Errorf("Expected header %s to have value %s, received %s\n", "head2", "two", val)
		}
	}
}

func TestFlattenHeaders_Empty(t *testing.T) {
	req, _ := http.NewRequest("GET", "localhost", nil)
	headers := flattenHeaders(req.Header)
	expectedLen := 0
	if len(headers) != expectedLen {
		t.Errorf("Expected %d headers, received %d\n", expectedLen, len(headers))
	}
}

func TestSuccessfulConstructHTTPRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"ping": "pong"}`))
	}))
	defer server.Close()
	r := &request.RequestObject{
		Method:    "GET",
		URL:       server.URL + "/api/ping",
		Body:      "hello",
		UserAgent: "ReqCorder/Test",
		Headers:   map[string]string{"one": "1", "two": "2"},
		Auth:      "top-secret",
	}
	_, err := constructHTTPRequest(r)
	if err != nil {
		t.Fatalf("Expected no error, received %v\n", err)
	}
}

func TestConstructHTTPRequest_IllegalMethod(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"ping": "pong"}`))
	}))
	defer server.Close()
	r := &request.RequestObject{
		Method:    illegalMethod,
		URL:       server.URL + "/api/ping",
		Body:      "hello",
		UserAgent: "ReqCorder/Test",
		Headers:   map[string]string{"one": "1", "two": "2"},
		Auth:      "top-secret",
	}
	_, err := constructHTTPRequest(r)
	if err == nil {
		t.Fatalf("Expected error, received nil\n")
	}
	if !errors.Is(err, ErrorFailedToBuildRequest) {
		t.Fatalf("Expected error %v, received %v\n", ErrorFailedToBuildRequest, err)
	}
}

func TestSuccessfulCreateTrace(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"ping": "pong"}`))
	}))
	defer server.Close()
	r := &request.RequestObject{
		Method:    "GET",
		URL:       server.URL + "/api/ping",
		Body:      "hello",
		UserAgent: "ReqCorder/Test",
		Headers:   map[string]string{"one": "1", "two": "2"},
		Auth:      "top-secret",
	}
	var timing response.ResponseTimes
	trace := createTrace(&timing)
	req, err := constructHTTPRequest(r)
	if err != nil {
		t.Fatalf("Expected no error, received %v\n", err)
	}
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	_, err = client.Do(req)
	if err != nil {
		t.Fatalf("Expected no error, received %v\n", err)
	}
	if timing.DNSStart.After(timing.DNSDone) {
		t.Fatalf("Incorrect DNS timings")
	}
	if timing.ConnectStart.After(timing.ConnectDone) {
		t.Fatalf("Incorrect connection timings")
	}
	if timing.TLSHandshakeStart.After(timing.TLSHandshakeDone) {
		t.Fatalf("Incorrect handshake timings")
	}
	if timing.GotFirstResponseByte.Before(timing.TLSHandshakeDone) {
		t.Fatalf("Incorrect first byte timing")
	}
}

func TestSuccessfulCreateTrace_ExternalDomain(t *testing.T) {
	r := &request.RequestObject{
		Method:    "GET",
		URL:       externalDomain,
		UserAgent: "ReqCorder/Test",
		Headers:   map[string]string{"one": "1", "two": "2"},
		Auth:      "top-secret",
	}
	var timing response.ResponseTimes
	trace := createTrace(&timing)
	req, err := constructHTTPRequest(r)
	if err != nil {
		t.Fatalf("Expected no error, received %v\n", err)
	}
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	_, err = client.Do(req)
	if err != nil {
		t.Fatalf("Expected no error, received %v\n", err)
	}
	if timing.DNSStart.After(timing.DNSDone) {
		t.Fatalf("Incorrect DNS timings")
	}
	if timing.ConnectStart.After(timing.ConnectDone) {
		t.Fatalf("Incorrect connection timings")
	}
	if timing.TLSHandshakeStart.After(timing.TLSHandshakeDone) {
		t.Fatalf("Incorrect handshake timings")
	}
	if timing.GotFirstResponseByte.Before(timing.TLSHandshakeDone) {
		t.Fatalf("Incorrect first byte timing")
	}
}

func TestSuccessfulInitiateRequest(t *testing.T) {
	resBody := `{"ping": "pong"}`
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(resBody))
	}))
	defer server.Close()
	certFile := writeCertToTempFile(t, server.Certificate())
	defer os.Remove(certFile) // Clean up
	sslVerify := true
	r := &request.RequestObject{
		Method:     "GET",
		URL:        server.URL + "/api/ping",
		UserAgent:  "ReqCorder/Test",
		Headers:    map[string]string{"one": "1", "two": "2"},
		Auth:       "top-secret",
		SSLVerify:  &sslVerify,
		CACertPath: certFile,
	}
	err := r.Validate()
	if err != nil {
		t.Fatalf("Expected no error, received %v\n", err)
	}
	res, _ := InitiateRequest(r)
	if res == nil {
		t.Fatalf("Expected non nil response")
	}
}

func TestInitiateRequest_ConstructFailure(t *testing.T) {
	sslVerify := true
	r := &request.RequestObject{
		Method:    illegalMethod,
		URL:       externalDomain,
		UserAgent: "ReqCorder/Test",
		Headers:   map[string]string{"one": "1", "two": "2"},
		Auth:      "top-secret",
		SSLVerify: &sslVerify,
	}
	_ = r.Validate()
	_, err := InitiateRequest(r)
	expectedErr := ErrorFailedToBuildRequest
	if err == nil {
		t.Fatalf("Expected error %v, received nil\n", expectedErr)
	}
	if !errors.Is(err, expectedErr) {
		t.Fatalf("Expected %v, received %v\n", expectedErr, err)
	}
}

func TestInitiateRequest_IllegalDomain(t *testing.T) {
	sslVerify := true
	r := &request.RequestObject{
		Method:    "GET",
		URL:       illegalDomain,
		UserAgent: "ReqCorder/Test",
		Headers:   map[string]string{"one": "1", "two": "2"},
		Auth:      "top-secret",
		SSLVerify: &sslVerify,
	}
	_ = r.Validate()
	_, err := InitiateRequest(r)
	expectedErr := ErrorRequestFailed
	if err == nil {
		t.Fatalf("Expected error %v, received nil\n", expectedErr)
	}
	if !errors.Is(err, expectedErr) {
		t.Fatalf("Expected %v, received %v\n", expectedErr, err)
	}
}

func TestInitiateRequest_ReadFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", "1000")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("short body"))
	}))
	defer server.Close()
	sslVerify := true
	r := &request.RequestObject{
		Method:    "GET",
		URL:       server.URL,
		UserAgent: "ReqCorder/Test",
		Headers:   map[string]string{"one": "1", "two": "2"},
		Auth:      "top-secret",
		SSLVerify: &sslVerify,
	}
	_ = r.Validate()
	_, err := InitiateRequest(r)
	expectedErr := ErrorFailedToReadResponse
	if err == nil {
		t.Fatalf("Expected error %v, received nil\n", expectedErr)
	}
	if !errors.Is(err, expectedErr) {
		t.Fatalf("Expected %v, received %v\n", expectedErr, err)
	}
}

func TestSuccessfulInitiateRequest_IllegalCertPath(t *testing.T) {
	resBody := `{"ping": "pong"}`
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(resBody))
	}))
	defer server.Close()
	certFile := writeCertToTempFile(t, server.Certificate())
	defer os.Remove(certFile) // Clean up
	sslVerify := true
	r := &request.RequestObject{
		Method:     "GET",
		URL:        server.URL + "/api/ping",
		UserAgent:  "ReqCorder/Test",
		Headers:    map[string]string{"one": "1", "two": "2"},
		Auth:       "top-secret",
		SSLVerify:  &sslVerify,
		CACertPath: certFile + ".com",
	}
	err := r.Validate()
	if err != nil {
		t.Fatalf("Expected no error, received %v\n", err)
	}
	_, err = InitiateRequest(r)
	expectedErr := ErrorFailedToReadCert
	if err == nil {
		t.Fatalf("Expected error %v, received nil\n", expectedErr)
	}
	if !errors.Is(err, expectedErr) {
		t.Fatalf("Expected %v, received %v\n", expectedErr, err)
	}
}

func TestSuccessfulInitiateRequest_IllegalCert(t *testing.T) {
	resBody := `{"ping": "pong"}`
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(resBody))
	}))
	defer server.Close()
	certFile := writeInvalidCertToTempFile(t)
	defer os.Remove(certFile) // Clean up
	sslVerify := true
	r := &request.RequestObject{
		Method:     "GET",
		URL:        server.URL + "/api/ping",
		UserAgent:  "ReqCorder/Test",
		Headers:    map[string]string{"one": "1", "two": "2"},
		Auth:       "top-secret",
		SSLVerify:  &sslVerify,
		CACertPath: certFile,
	}
	err := r.Validate()
	if err != nil {
		t.Fatalf("Expected no error, received %v\n", err)
	}
	_, err = InitiateRequest(r)
	expectedErr := ErrorFailedToReadCert
	if err == nil {
		t.Fatalf("Expected error %v, received nil\n", expectedErr)
	}
	if !errors.Is(err, expectedErr) {
		t.Fatalf("Expected %v, received %v\n", expectedErr, err)
	}
}

// func TestSuccessfulInitiateRequest_CustomCert(t *testing.T) {
// 	// Create TLS test server
// 	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		w.WriteHeader(http.StatusOK)
// 		w.Write([]byte("success"))
// 	}))
// 	defer server.Close()

// 	// Write server certificate to temporary file
// 	certFile := writeCertToTempFile(t, server.Certificate())
// 	defer os.Remove(certFile) // Clean up

// 	// Test your function with the cert file path
// 	result, err := YourFunction(server.URL, certFile)

// 	if err != nil {
// 		t.Fatalf("unexpected error with valid CA cert: %v", err)
// 	}

// 	if result == nil {
// 		t.Error("expected result, got nil")
// 	}
// }

func writeCertToTempFile(t *testing.T, cert *x509.Certificate) string {
	t.Helper()
	tmpDir := t.TempDir()
	certPath := filepath.Join(tmpDir, "ca-cert.pem")
	certPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert.Raw,
	})
	err := os.WriteFile(certPath, certPEM, 0644)
	if err != nil {
		t.Fatalf("failed to write cert file: %v", err)
	}
	return certPath
}

func writeInvalidCertToTempFile(t *testing.T) string {
	t.Helper()
	tmpDir := t.TempDir()
	certPath := filepath.Join(tmpDir, "invalid-cert.pem")
	if err := os.WriteFile(certPath, []byte("not a valid certificate"), 0644); err != nil {
		t.Fatalf("failed to write invalid cert file: %v", err)
	}
	return certPath
}

// func TestFailedGetSortedResponsesByTemplateHash_StoreFailure(t *testing.T) {
// 	root := t.TempDir()
// 	recordOne := &request.RequestObject{
// 		RecordStorePath: root,
// 		TemplateHash:    "templateOne",
// 	}
// 	_, err := recordOne.GetSortedResponsesByTemplateHash()
// 	if err == nil {
// 		t.Fatalf("Expected error, received nil")
// 	}
// 	expectedErr := ErrorFailedToStatPath
// 	if !errors.Is(err, expectedErr) {
// 		t.Fatalf("Expected error %v, received %v", expectedErr, err)
// 	}
// }
