package request

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"reflect"
	"testing"
	"time"

	"log/slog"
)

func compareStructFields(t *testing.T, got, want any, fieldsToCheck []string) {
	gotVal := reflect.ValueOf(got)
	wantVal := reflect.ValueOf(want)
	gotType := reflect.TypeOf(got)
	if gotVal.Kind() == reflect.Pointer {
		if gotVal.IsNil() {
			t.Fatalf("Got is nil")
		}
		gotVal = gotVal.Elem()
		gotType = gotType.Elem()
	}
	if wantVal.Kind() == reflect.Ptr {
		if wantVal.IsNil() {
			t.Fatalf("Want is nil")
		}
		wantVal = wantVal.Elem()
	}
	if gotVal.Kind() != reflect.Struct {
		t.Fatalf("Expected got struct, received %v", gotVal.Kind())
	}
	if wantVal.Kind() != reflect.Struct {
		t.Fatalf("Expected want struct, received %v", wantVal.Kind())
	}
	for i := 0; i < gotVal.NumField(); i++ {
		fieldName := gotType.Field(i).Name
		if len(fieldsToCheck) > 0 {
			found := false
			for _, checkField := range fieldsToCheck {
				if checkField == fieldName {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}
		gotField := gotVal.Field(i)
		wantField := wantVal.Field(i)
		if wantField.IsZero() {
			continue
		}

		if !reflect.DeepEqual(gotField.Interface(), wantField.Interface()) {
			t.Errorf("Field %s: expected %v, received %v",
				fieldName, wantField.Interface(), gotField.Interface())
		}
	}
}

func TestProcessBasics(t *testing.T) {
	floatTimeout := 66.6
	tests := []struct {
		name           string
		input          *RequestObject
		expectedOutput *RequestObject
		expectedError  error
	}{
		{
			name: "Invalid method",
			input: &RequestObject{
				Method: "WOW",
				URL:    "https://example.com",
			},
			expectedError: ErrorInvalidMethod,
		},
		{
			name: "Invalid URL",
			input: &RequestObject{
				Method: "GET",
				URL:    "invalid-url",
			},
			expectedError: ErrorInvalidURL,
		},
		{
			name: "Valid GET",
			input: &RequestObject{
				Method: "GET",
				URL:    "https://example.com",
			},
			expectedOutput: &RequestObject{
				Method:         "GET",
				URL:            "https://example.com",
				TimeoutSeconds: 30,
				Timeout:        30 * time.Second,
				UserAgent:      "ReqCorder",
			},
		},
		{
			name: "Valid POST",
			input: &RequestObject{
				Method: "POST",
				URL:    "https://example.com",
			},
			expectedOutput: &RequestObject{
				Method:         "POST",
				URL:            "https://example.com",
				TimeoutSeconds: 30,
				Timeout:        30 * time.Second,
				UserAgent:      "ReqCorder",
			},
		},
		{
			name: "Valid DELETE",
			input: &RequestObject{
				Method: "DELETE",
				URL:    "https://example.com",
			},
			expectedOutput: &RequestObject{
				Method:         "DELETE",
				URL:            "https://example.com",
				TimeoutSeconds: 30,
				Timeout:        30 * time.Second,
				UserAgent:      "ReqCorder",
			},
		},
		{
			name: "Valid OPTIONS",
			input: &RequestObject{
				Method: "OPTIONS",
				URL:    "https://example.com",
			},
			expectedOutput: &RequestObject{
				Method:         "OPTIONS",
				URL:            "https://example.com",
				TimeoutSeconds: 30,
				Timeout:        30 * time.Second,
				UserAgent:      "ReqCorder",
			},
		},
		{
			name: "Valid PATCH",
			input: &RequestObject{
				Method: "PATCH",
				URL:    "https://example.com",
			},
			expectedOutput: &RequestObject{
				Method:         "PATCH",
				URL:            "https://example.com",
				TimeoutSeconds: 30,
				Timeout:        30 * time.Second,
				UserAgent:      "ReqCorder",
			},
		},
		{
			name: "Valid PUT",
			input: &RequestObject{
				Method: "PUT",
				URL:    "https://example.com",
			},
			expectedOutput: &RequestObject{
				Method:         "PUT",
				URL:            "https://example.com",
				TimeoutSeconds: 30,
				Timeout:        30 * time.Second,
				UserAgent:      "ReqCorder",
			},
		},
		{
			name: "Valid HEAD",
			input: &RequestObject{
				Method: "HEAD",
				URL:    "https://example.com",
			},
			expectedOutput: &RequestObject{
				Method:         "HEAD",
				URL:            "https://example.com",
				TimeoutSeconds: 30,
				Timeout:        30 * time.Second,
				UserAgent:      "ReqCorder",
			},
		},
		{
			name: "Valid float timeout",
			input: &RequestObject{
				Method:         "PATCH",
				URL:            "https://example.com",
				TimeoutSeconds: floatTimeout,
			},
			expectedOutput: &RequestObject{
				Method:         "PATCH",
				URL:            "https://example.com",
				TimeoutSeconds: floatTimeout,
				Timeout:        time.Duration(floatTimeout) * time.Second,
				UserAgent:      "ReqCorder",
			},
		},
		{
			name: "Valid int timeout",
			input: &RequestObject{
				Method:         "PATCH",
				URL:            "https://example.com",
				TimeoutSeconds: 60,
			},
			expectedOutput: &RequestObject{
				Method:         "PATCH",
				URL:            "https://example.com",
				TimeoutSeconds: 60,
				Timeout:        time.Duration(60) * time.Second,
				UserAgent:      "ReqCorder",
			},
		},
		{
			name: "Negative timeout",
			input: &RequestObject{
				Method:         "PATCH",
				URL:            "https://example.com",
				TimeoutSeconds: -10,
			},
			expectedOutput: &RequestObject{
				Method:         "PATCH",
				URL:            "https://example.com",
				TimeoutSeconds: 30,
				Timeout:        30 * time.Second,
				UserAgent:      "ReqCorder",
			},
		},
		{
			name: "Zero timeout",
			input: &RequestObject{
				Method:         "PATCH",
				URL:            "https://example.com",
				TimeoutSeconds: 0,
			},
			expectedOutput: &RequestObject{
				Method:         "PATCH",
				URL:            "https://example.com",
				TimeoutSeconds: 30,
				Timeout:        30 * time.Second,
				UserAgent:      "ReqCorder",
			},
		},
		{
			name: "Custom User Agent",
			input: &RequestObject{
				Method:         "PATCH",
				URL:            "https://example.com",
				TimeoutSeconds: 30,
				UserAgent:      "MyBrowser",
			},
			expectedOutput: &RequestObject{
				Method:         "PATCH",
				URL:            "https://example.com",
				TimeoutSeconds: 30,
				Timeout:        30 * time.Second,
				UserAgent:      "MyBrowser",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.input.processBasics()
			if tt.expectedError != nil && err == nil {
				t.Fatalf("Expected error %v, received nil", tt.expectedError)
			}
			if !errors.Is(err, tt.expectedError) {
				t.Fatalf("Expected error %v, received %v", tt.expectedError, err)

			}
			if tt.expectedOutput != nil {
				compareStructFields(t, tt.input, tt.expectedOutput, []string{"Method", "URL", "TimeoutSeconds", "Timeout", "UserAgent"})
			}
		})
	}
}

func TestProcessBodyVars(t *testing.T) {
	tests := []struct {
		name           string
		input          *RequestObject
		expectedOutput *RequestObject
		expectedError  error
	}{
		{
			name: "Valid body and body vars",
			input: &RequestObject{
				Body: "{\"name\": \"{{name}}\"}",
				BodyVars: map[string]string{
					"name": "Something",
				},
			},
			expectedOutput: &RequestObject{
				Body: "{\"name\": \"Something\"}",
			},
		},
		{
			name: "Empty body and valid body vars",
			input: &RequestObject{
				Body: "",
				BodyVars: map[string]string{
					"name": "Something",
				},
			},
			expectedOutput: &RequestObject{
				Body: "",
			},
		},
		{
			name: "Mutliple references to body var",
			input: &RequestObject{
				Body: "{\"name\": \"{{name}}\", \"firstName\": \"{{name}}\"}",
				BodyVars: map[string]string{
					"name": "Something",
				},
			},
			expectedOutput: &RequestObject{
				Body: "{\"name\": \"Something\", \"firstName\": \"Something\"}",
			},
		},
		{
			name: "Mutliple body vars",
			input: &RequestObject{
				Body: "{\"name\": \"{{name}}\", \"firstName\": \"{{firstName}}\"}",
				BodyVars: map[string]string{
					"name":      "Something",
					"firstName": "Wrong",
				},
			},
			expectedOutput: &RequestObject{
				Body: "{\"name\": \"Something\", \"firstName\": \"Wrong\"}",
			},
		},
		{
			name: "Missing body var",
			input: &RequestObject{
				Body:     "{\"name\": \"{{name}}\", \"firstName\": \"{{firstName}}\"}",
				BodyVars: map[string]string{},
			},
			expectedOutput: &RequestObject{
				Body: "{\"name\": \"{{name}}\", \"firstName\": \"{{firstName}}\"}",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.input.processBodyVars()
			compareStructFields(t, tt.input, tt.expectedOutput, []string{"Body"})
		})
	}
}

func TestProcessEnvVars(t *testing.T) {
	path := os.Getenv("PATH")
	home := os.Getenv("HOME")
	tests := []struct {
		name           string
		input          *RequestObject
		expectedOutput *RequestObject
		expectedError  error
	}{
		{
			name: "Valid environment variable",
			input: &RequestObject{
				Body: "{\"name\": \"{{env:PATH}}\"}",
			},
			expectedOutput: &RequestObject{
				Body: fmt.Sprintf("{\"name\": \"%s\"}", path),
			},
		},
		{
			name: "Non existent environment variable",
			input: &RequestObject{
				Body: "{\"name\": \"{{env:DAMN}}\"}",
			},
			expectedOutput: &RequestObject{
				Body: "{\"name\": \"\"}",
			},
		},
		{
			name: "Multiple valid environment variable",
			input: &RequestObject{
				Body: "{\"name\": \"{{env:PATH}}\", \"firstName\": \"{{env:HOME}}\"}",
			},
			expectedOutput: &RequestObject{
				Body: fmt.Sprintf("{\"name\": \"%s\", \"firstName\": \"%s\"}", path, home),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.input.processEnvVars()
			compareStructFields(t, tt.input, tt.expectedOutput, []string{"Body"})
		})
	}
}

func TestProcessAuth(t *testing.T) {
	tests := []struct {
		name           string
		input          *RequestObject
		expectedOutput *RequestObject
		expectedError  error
	}{
		{
			name: "Valid bearer",
			input: &RequestObject{
				AuthType: "Bearer",
				Auth:     "token",
			},
			expectedOutput: &RequestObject{
				Auth: "Bearer token",
			},
		},
		{
			name: "Valid basic",
			input: &RequestObject{
				AuthType: "Basic",
				Auth:     "token",
			},
			expectedOutput: &RequestObject{
				Auth: "Basic token",
			},
		},
		{
			name: "Valid bearer mixed case",
			input: &RequestObject{
				AuthType: "bEAreR",
				Auth:     "token",
			},
			expectedOutput: &RequestObject{
				Auth: "Bearer token",
			},
		},
		{
			name: "Valid basic mixed case",
			input: &RequestObject{
				AuthType: "BaSiC",
				Auth:     "token",
			},
			expectedOutput: &RequestObject{
				Auth: "Basic token",
			},
		},
		{
			name: "Bearer already present",
			input: &RequestObject{
				AuthType: "Bearer",
				Auth:     "bearer token",
			},
			expectedOutput: &RequestObject{
				Auth: "bearer token",
			},
		},
		{
			name: "Basic already present",
			input: &RequestObject{
				AuthType: "Basic",
				Auth:     "basic token",
			},
			expectedOutput: &RequestObject{
				Auth: "basic token",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.input.processAuth()
			compareStructFields(t, tt.input, tt.expectedOutput, []string{"Auth"})
		})
	}
}

func TestProcessCookies(t *testing.T) {
	tests := []struct {
		name           string
		input          *RequestObject
		expectedOutput *RequestObject
		expectedError  error
	}{
		{
			name: "Valid cookies",
			input: &RequestObject{
				URL: "http://example.com",
				Cookies: map[string]string{
					"cookie1": "cook",
				},
			},
			expectedOutput: &RequestObject{
				URL: "http://example.com",
				Cookies: map[string]string{
					"cookie1": "cook",
				},
			},
		},
		{
			name: "Invalid URL cookies",
			input: &RequestObject{
				URL: "wwww",
				Cookies: map[string]string{
					"cookie1": "cook",
				},
			},
			expectedError: ErrorInvalidURL,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsedURL, _ := url.ParseRequestURI(tt.input.URL)
			err := tt.input.processCookies()
			if !errors.Is(err, tt.expectedError) {
				t.Fatalf("Expected error %v, received %v", tt.expectedError, err)
			}
			if err != nil {
				return
			}
			compareStructFields(t, tt.input, tt.expectedOutput, []string{"URL", "Cookies"})
			if len(tt.input.Cookies) > 0 {
				if tt.input.CookieJar == nil {
					t.Fatal("Expected CookieJar to be created")
				}
				actualCookies := make(map[string]string)
				for _, cookie := range tt.input.CookieJar.Cookies(parsedURL) {
					actualCookies[cookie.Name] = cookie.Value
				}
				for expectedName, expectedValue := range tt.input.Cookies {
					if actualValue, exists := actualCookies[expectedName]; !exists {
						t.Errorf("Expected cookie %q not found in jar", expectedName)
					} else if actualValue != expectedValue {
						t.Errorf("Cookie %q expected %q, received %q", expectedName, expectedValue, actualValue)
					}
				}
			}
		})
	}
}

// func TestProcessCookiesError(t *testing.T) {
// 	originalNewCookieJar := newCookieJar
// 	defer func() { newCookieJar = originalNewCookieJar }()
// 	newCookieJar = func(*cookiejar.Options) (*cookiejar.Jar, error) {
// 		return nil, ErrorFailedToCreateCookieJar
// 	}
// 	input := &RequestObject{
// 		URL:       "https://example.com",
// 		CookieJar: nil,
// 	}
// 	err := input.processCookies()
// 	if err == nil {
// 		t.Fatal("Expected error, got nil")
// 	}
// 	if !errors.Is(err, ErrorFailedToCreateCookieJar) {
// 		t.Fatalf("Expected error containing %v, received %v", ErrorFailedToCreateCookieJar, err)
// 	}
// 	if input.CookieJar != nil {
// 		t.Error("Expected CookieJar to remain nil on error")
// 	}
// }

func TestValidate(t *testing.T) {
	tests := []struct {
		name           string
		input          *RequestObject
		expectedOutput *RequestObject
		expectedError  error
	}{
		{
			name: "Valid request with all fields",
			input: &RequestObject{
				URL:            "https://example.com",
				Method:         "GET",
				Headers:        map[string]string{"Content-Type": "application/json"},
				Cookies:        map[string]string{"session": "abc123"},
				Auth:           "token",
				AuthType:       "Bearer",
				Body:           "{\"name\": \"{{name}}\"}",
				BodyVars:       map[string]string{"name": "John"},
				TimeoutSeconds: 60,
				UserAgent:      "TestAgent",
			},
			expectedOutput: &RequestObject{
				URL:            "https://example.com",
				Method:         "GET",
				Headers:        map[string]string{"Content-Type": "application/json"},
				Cookies:        map[string]string{"session": "abc123"},
				Auth:           "Bearer token",
				AuthType:       "Bearer",
				Body:           "{\"name\": \"John\"}",
				BodyVars:       map[string]string{"name": "John"},
				TimeoutSeconds: 60,
				Timeout:        60 * time.Second,
				UserAgent:      "TestAgent",
			},
		},
		{
			name: "Invalid URL",
			input: &RequestObject{
				URL:    "invalid-url",
				Method: "GET",
			},
			expectedError: ErrorInvalidURL,
		},
		{
			name: "Invalid method",
			input: &RequestObject{
				URL:    "https://example.com",
				Method: "WOW",
			},
			expectedError: ErrorInvalidMethod,
		},
		{
			name: "Valid request with environment variables",
			input: &RequestObject{
				URL:    "https://example.com",
				Method: "GET",
				Body:   "{\"path\": \"{{env:PATH}}\", \"home\": \"{{env:HOME}}\"}",
			},
			expectedOutput: &RequestObject{
				URL:            "https://example.com",
				Method:         "GET",
				Body:           fmt.Sprintf("{\"path\": \"%s\", \"home\": \"%s\"}", os.Getenv("PATH"), os.Getenv("HOME")),
				TimeoutSeconds: 30,
				Timeout:        30 * time.Second,
				UserAgent:      "ReqCorder",
			},
		},
		{
			name: "Valid request with basic auth",
			input: &RequestObject{
				URL:      "https://example.com",
				Method:   "GET",
				Auth:     "token",
				AuthType: "Basic",
			},
			expectedOutput: &RequestObject{
				URL:            "https://example.com",
				Method:         "GET",
				Auth:           "Basic token",
				AuthType:       "Basic",
				TimeoutSeconds: 30,
				Timeout:        30 * time.Second,
				UserAgent:      "ReqCorder",
			},
		},
		{
			name: "Request with cookies",
			input: &RequestObject{
				URL:     "https://example.com",
				Method:  "GET",
				Cookies: map[string]string{"session": "abc123"},
			},
			expectedOutput: &RequestObject{
				URL:            "https://example.com",
				Method:         "GET",
				Cookies:        map[string]string{"session": "abc123"},
				TimeoutSeconds: 30,
				Timeout:        30 * time.Second,
				UserAgent:      "ReqCorder",
			},
		},
		// {
		// 	name: "Request with failed cookie jar creation",
		// 	input: &RequestObject{
		// 		URL:     "https://example.com",
		// 		Method:  "GET",
		// 		Cookies: map[string]string{"session": "abc123"},
		// 	},
		// 	expectedError: ErrorFailedToCreateCookieJar,
		// },
		{
			name: "Request with default values",
			input: &RequestObject{
				URL:    "https://example.com",
				Method: "GET",
			},
			expectedOutput: &RequestObject{
				URL:            "https://example.com",
				Method:         "GET",
				TimeoutSeconds: 30,
				Timeout:        30 * time.Second,
				UserAgent:      "ReqCorder",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// if tt.expectedError == ErrorFailedToCreateCookieJar {
			// 	originalNewCookieJar := newCookieJar
			// 	defer func() { newCookieJar = originalNewCookieJar }()
			// 	newCookieJar = func(*cookiejar.Options) (*cookiejar.Jar, error) {
			// 		return nil, ErrorFailedToCreateCookieJar
			// 	}
			// }

			err := tt.input.Validate()
			if tt.expectedError != nil {
				if err == nil {
					t.Fatalf("Expected error %v, received nil", tt.expectedError)
				}
				if !errors.Is(err, tt.expectedError) {
					t.Fatalf("Expected error %v, received %v", tt.expectedError, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if tt.expectedOutput != nil {
				compareStructFields(t, tt.input, tt.expectedOutput, []string{
					"URL", "Method", "Headers", "Cookies", "Auth", "AuthType",
					"Body", "BodyVars", "TimeoutSeconds", "Timeout", "UserAgent",
				})
				if len(tt.input.Cookies) > 0 && tt.input.CookieJar != nil {
					parsedURL, _ := url.ParseRequestURI(tt.input.URL)
					actualCookies := make(map[string]string)
					for _, cookie := range tt.input.CookieJar.Cookies(parsedURL) {
						actualCookies[cookie.Name] = cookie.Value
					}
					for expectedName, expectedValue := range tt.input.Cookies {
						if actualValue, exists := actualCookies[expectedName]; !exists {
							t.Errorf("Expected cookie %q not found in jar", expectedName)
						} else if actualValue != expectedValue {
							t.Errorf("Cookie %q expected %q, received %q", expectedName, expectedValue, actualValue)
						}
					}
				}
			}
		})
	}
}

func TestLogValue(t *testing.T) {
	tests := []struct {
		name     string
		input    *RequestObject
		expected slog.Value
	}{
		{
			name:     "Empty request",
			input:    nil,
			expected: slog.StringValue("<nil>"),
		},
		{
			name: "Valid request with all fields",
			input: &RequestObject{
				TemplateHash:   "hash123",
				URL:            "https://example.com",
				Method:         "GET",
				UserAgent:      "TestAgent",
				Timeout:        30 * time.Second,
				TimeoutSeconds: 30,
				Body:           "{\"name\": \"John\"}",
				AuthType:       "Bearer",
				Auth:           "token",
				AuthHeaderName: "Authorization",
				SSLVerify:      new(bool),
				CACertPath:     "/path/to/cert",
				Headers: map[string]string{
					"Content-Type": "application/json",
					"Accept":       "application/json",
				},
				Cookies: map[string]string{
					"session": "abc123",
					"theme":   "dark",
				},
				BodyVars: map[string]string{
					"name": "John",
				},
			},
			expected: slog.GroupValue(
				slog.String("templateHash", "hash123"),
				slog.String("url", "https://example.com"),
				slog.String("method", "GET"),
				slog.String("userAgent", "TestAgent"),
				slog.Duration("timeout", 30*time.Second),
				slog.String("body", "{\"name\": \"John\"}"),
				slog.String("authType", "Bearer"),
				slog.String("auth", "token"),
				slog.String("authHeaderName", "Authorization"),
				slog.Float64("timeoutSeconds", 30),
				slog.Bool("sslVerify", false),
				slog.String("caCertPath", "/path/to/cert"),
				slog.Attr{
					Key:   "headers",
					Value: slog.GroupValue(slog.String("Content-Type", "application/json"), slog.String("Accept", "application/json")),
				},
				slog.Attr{
					Key:   "cookies",
					Value: slog.GroupValue(slog.String("session", "abc123"), slog.String("theme", "dark")),
				},
				slog.Attr{
					Key:   "bodyVars",
					Value: slog.GroupValue(slog.String("name", "John")),
				},
			),
		},
		{
			name: "Request with empty maps",
			input: &RequestObject{
				TemplateHash:   "hash123",
				URL:            "https://example.com",
				Method:         "GET",
				UserAgent:      "TestAgent",
				Timeout:        30 * time.Second,
				TimeoutSeconds: 30,
				Body:           "",
				AuthType:       "",
				Auth:           "",
				AuthHeaderName: "",
				SSLVerify:      new(bool),
				CACertPath:     "",
				Headers:        map[string]string{},
				Cookies:        map[string]string{},
				BodyVars:       map[string]string{},
			},
			expected: slog.GroupValue(
				slog.String("templateHash", "hash123"),
				slog.String("url", "https://example.com"),
				slog.String("method", "GET"),
				slog.String("userAgent", "TestAgent"),
				slog.Duration("timeout", 30*time.Second),
				slog.String("body", ""),
				slog.String("authType", ""),
				slog.String("auth", ""),
				slog.String("authHeaderName", ""),
				slog.Float64("timeoutSeconds", 30),
				slog.Bool("sslVerify", false),
				slog.String("caCertPath", ""),
				slog.Attr{
					Key:   "headers",
					Value: slog.GroupValue(),
				},
				slog.Attr{
					Key:   "cookies",
					Value: slog.GroupValue(),
				},
				slog.Attr{
					Key:   "bodyVars",
					Value: slog.GroupValue(),
				},
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.LogValue()
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("LogValue() = %v, want %v", result, tt.expected)
			}
		})
	}
}
