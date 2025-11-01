package request

import (
	"log/slog"
)

// Helper function to log pointers to RequestObject.
func (r *RequestObject) LogValue() slog.Value {
	if r == nil {
		return slog.StringValue("<nil>")
	}
	headerAttrs := make([]slog.Attr, 0, len(r.Headers))
	for key, value := range r.Headers {
		headerAttrs = append(headerAttrs, slog.String(key, value))
	}
	cookieAttrs := make([]slog.Attr, 0, len(r.Cookies))
	for key, value := range r.Cookies {
		cookieAttrs = append(cookieAttrs, slog.String(key, value))
	}
	bodyVarAttrs := make([]slog.Attr, 0, len(r.BodyVars))
	for key, value := range r.BodyVars {
		bodyVarAttrs = append(bodyVarAttrs, slog.String(key, value))
	}
	return slog.GroupValue(
		slog.String("templateHash", r.TemplateHash),
		slog.String("url", r.URL),
		slog.String("method", r.Method),
		slog.String("userAgent", r.UserAgent),
		slog.Duration("timeout", r.Timeout),
		slog.String("body", r.Body),
		slog.String("authType", r.AuthType),
		slog.String("auth", r.Auth),
		slog.String("authHeaderName", r.AuthHeaderName),
		slog.Float64("timeoutSeconds", r.TimeoutSeconds),
		slog.Bool("sslVerify", *r.SSLVerify),
		slog.String("caCertPath", r.CACertPath),
		slog.Attr{
			Key:   "headers",
			Value: slog.GroupValue(headerAttrs...),
		},
		slog.Attr{
			Key:   "cookies",
			Value: slog.GroupValue(cookieAttrs...),
		},
		slog.Attr{
			Key:   "bodyVars",
			Value: slog.GroupValue(bodyVarAttrs...),
		},
	)
}
