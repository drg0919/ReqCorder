package request

import (
	"net/http/cookiejar"
	"time"
)

// RequestObject represents an HTTP request with all its configuration options and metadata.
type RequestObject struct {
	TemplateHash   string            `yaml:"template_hash"`
	URL            string            `yaml:"url"`
	Method         string            `yaml:"method"`
	Headers        map[string]string `yaml:"headers,omitempty"`
	Cookies        map[string]string `yaml:"cookies,omitempty"`
	CookieJar      *cookiejar.Jar    `yaml:"-"`
	Auth           string            `yaml:"auth,omitempty"`
	AuthType       string            `yaml:"auth_type,omitempty"`
	AuthHeaderName string            `yaml:"auth_header_name,omitempty"`
	UserAgent      string            `yaml:"user_agent,omitempty"`
	Body           string            `yaml:"body,omitempty"`
	Timeout        time.Duration     `yaml:"-"`
	TimeoutSeconds float64           `yaml:"timeout,omitempty"`
	BodyVars       map[string]string `yaml:"body_vars,omitempty"`
	SSLVerify      *bool             `yaml:"ssl_verify,omitempty"`
	CACertPath     string            `yaml:"ca_cert_path,omitempty"`
}
