package response

import (
	"net/http"
	"time"
)

// ResponseObject represents the complete response from an HTTP request, including metadata, headers, body, and timing information.
type ResponseObject struct {
	RequestHash  string            `yaml:"request_hash"`
	TemplateHash string            `yaml:"template_hash"`
	StatusCode   int               `yaml:"status_code"`
	Headers      map[string]string `yaml:"headers"`
	Body         string            `yaml:"body"`
	Size         int64             `yaml:"size_bytes"`
	Timing       ResponseTimes     `yaml:"timing"`
	Cookies      []*http.Cookie    `yaml:"cookies"`
}

// ResponseTimes contains timing information for various stages of an HTTP request.
type ResponseTimes struct {
	DNSStart             time.Time     `yaml:"dns_start"`
	DNSDone              time.Time     `yaml:"dns_end"`
	ConnectStart         time.Time     `yaml:"connect_start"`
	ConnectDone          time.Time     `yaml:"connect_done"`
	TLSHandshakeStart    time.Time     `yaml:"tls_handshake_start"`
	TLSHandshakeDone     time.Time     `yaml:"tls_handshake_done"`
	GotFirstResponseByte time.Time     `yaml:"got_first_response_byte"`
	DNSLookup            time.Duration `yaml:"dns_lookup"`
	TCPConnect           time.Duration `yaml:"tcp_connect"`
	TLSHandshake         time.Duration `yaml:"tls_handshake"`
	FirstByte            time.Duration `yaml:"time_to_first_byte"`
	Total                time.Duration `yaml:"total_duration"`
}
