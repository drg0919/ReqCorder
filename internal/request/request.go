package request

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
)

// Run all processing steps to prepare request for execution.
func (r *RequestObject) Validate() error {
	slog.Debug("Raw request object", slog.Any("requestObject", r))
	err := r.processBasics()
	if err != nil {
		slog.Error("Error processing request object", "error", err)
		return err
	}
	slog.Debug("Processing body variables")
	r.processBodyVars()
	slog.Debug("Processing environment variables")
	r.processEnvVars()
	slog.Debug("Processing authentication")
	r.processAuth()
	slog.Debug("Processing cookies")
	r.processCookies()
	slog.Debug("Request object post processing", slog.Any("requestObject", r))
	return nil
}

// Validate and set default values for basic request fields.
func (r *RequestObject) processBasics() error {
	slog.Debug("Validating basic request fields")
	validMethods := map[string]bool{
		"GET":     true,
		"POST":    true,
		"PUT":     true,
		"DELETE":  true,
		"PATCH":   true,
		"HEAD":    true,
		"OPTIONS": true,
	}
	slog.Debug("Validating URL", "url", r.URL)
	_, err := url.ParseRequestURI(r.URL)
	if err != nil {
		return fmt.Errorf("%w %q: %w", ErrorInvalidURL, r.URL, err)
	}
	cleanedMethod := strings.TrimSpace(strings.ToUpper(r.Method))
	if _, exists := validMethods[cleanedMethod]; !exists {
		return fmt.Errorf("%w %q", ErrorInvalidMethod, r.Method)
	}
	r.Method = cleanedMethod
	if r.TimeoutSeconds > 0 {
		r.Timeout = time.Duration(r.TimeoutSeconds) * time.Second
	} else {
		r.TimeoutSeconds = 30
		r.Timeout = 30 * time.Second
	}
	if r.UserAgent == "" {
		r.UserAgent = "ReqCorder"
	}
	t := true
	if r.SSLVerify == nil {
		r.SSLVerify = &t
	}
	slog.Debug("Basic validation completed", "method", r.Method, "timeout", r.Timeout, "userAgent", r.UserAgent)
	return nil
}

// Replace placeholders in request body with BodyVars values.
func (r *RequestObject) processBodyVars() {
	if r.Body == "" {
		slog.Debug("No body content, skipping body variable processing")
		return
	}
	slog.Debug("Processing body variables", "bodyVarsCount", len(r.BodyVars), "bodyLength", len(r.Body))
	body := r.Body
	for key, value := range r.BodyVars {
		placeholder := "{{" + key + "}}"
		body = strings.ReplaceAll(body, placeholder, value)
		slog.Debug("Replaced body variable", "key", key, "placeholder", placeholder)
	}
	r.Body = body
	slog.Debug("Body variable processing completed", "newBodyLength", len(body))
}

// Substitute {{env:VAR}} placeholders with environment variable values.
func (r *RequestObject) processEnvVars() {
	slog.Debug("Processing environment variables in request body")
	envRegex := regexp.MustCompile(`\{\{env:([A-Z_][A-Z0-9_]+)\}\}`)
	matches := envRegex.FindAllStringSubmatch(r.Body, -1)
	if len(matches) == 0 {
		slog.Debug("No environment variable placeholders found in request body")
		return
	}
	for _, match := range matches {
		if len(match) > 1 {
			varName := match[1]
			slog.Debug("Found environment variable placeholder", "varName", varName)
		}
	}
	r.Body = envRegex.ReplaceAllStringFunc(r.Body, func(match string) string {
		varName := envRegex.FindStringSubmatch(match)[1]
		value := os.Getenv(varName)
		slog.Debug("Replacing environment variable", "varName", varName, "found", value != "")
		return value
	})
	slog.Debug("Environment variable processing completed", "bodyLength", len(r.Body))
}

// Ensure Authorization header is correctly prefixed.
func (r *RequestObject) processAuth() {
	slog.Debug("Processing authentication", "authType", r.AuthType, "hasAuth", r.Auth != "")
	if r.Auth == "" {
		slog.Debug("No authentication configured, skipping auth processing")
		return
	}
	if strings.ToLower(r.AuthType) == "bearer" && !strings.HasPrefix(r.Auth, "bearer") && !strings.HasPrefix(r.Auth, "Bearer") {
		r.Auth = "Bearer " + r.Auth
		slog.Debug("Added Bearer prefix to authentication token")
	}
	if strings.ToLower(r.AuthType) == "basic" && !strings.HasPrefix(r.Auth, "basic") && !strings.HasPrefix(r.Auth, "Basic") {
		r.Auth = "Basic " + r.Auth
		slog.Debug("Added Basic prefix to authentication token")
	}
	slog.Debug("Authentication processing completed", "finalAuth", r.Auth)
}

// Initialize cookie jar and add defined cookies for request URL.
func (r *RequestObject) processCookies() error {
	slog.Debug("Processing cookies", "cookieCount", len(r.Cookies))
	if r.CookieJar == nil {
		slog.Debug("Cookie jar not initialized, creating new jar")
		jar, _ := cookiejar.New(nil)
		r.CookieJar = jar
	}
	if len(r.Cookies) == 0 {
		slog.Debug("No cookies defined, skipping cookie processing")
		return nil
	}
	parsedURL, err := url.ParseRequestURI(r.URL)
	if err != nil {
		return fmt.Errorf("%w %q: %v", ErrorInvalidURL, r.URL, err)
	}
	cookies := []*http.Cookie{}
	for name, value := range r.Cookies {
		cookie := &http.Cookie{
			Name:   name,
			Value:  value,
			Path:   "/",
			Domain: parsedURL.Host,
		}
		cookies = append(cookies, cookie)
		slog.Debug("Added cookie to request", "name", name, "domain", parsedURL.Host)
	}
	r.CookieJar.SetCookies(parsedURL, cookies)
	slog.Debug("Cookie processing completed", "url", r.URL)
	return nil
}
