package response

import (
	"log/slog"
)

// Helper function to log pointers to ResponseObject.
func (r *ResponseObject) LogValue() slog.Value {
	if r == nil {
		return slog.StringValue("<nil>")
	}
	headerAttrs := make([]slog.Attr, 0, len(r.Headers))
	for key, value := range r.Headers {
		headerAttrs = append(headerAttrs, slog.String(key, value))
	}
	cookieAttrs := make([]slog.Attr, 0, len(r.Cookies))
	for _, cookie := range r.Cookies {
		cookieAttrs = append(cookieAttrs, slog.String(cookie.Name, cookie.Value))
	}
	return slog.GroupValue(
		slog.String("requestHash", r.RequestHash),
		slog.String("templateHash", r.TemplateHash),
		slog.Int("statusCode", r.StatusCode),
		slog.Int64("size", r.Size),
		slog.String("body", r.Body),
		slog.Attr{
			Key:   "headers",
			Value: slog.GroupValue(headerAttrs...),
		},
		slog.Attr{
			Key:   "cookies",
			Value: slog.GroupValue(cookieAttrs...),
		},
		slog.Attr{
			Key:   "times",
			Value: slog.GroupValue(r.Timing.LogValue().Group()...),
		},
	)
}

// Helper function to log pointers to ResponseTimes.
func (r *ResponseTimes) LogValue() slog.Value {
	if r == nil {
		return slog.StringValue("<nil>")
	}
	return slog.GroupValue(
		slog.Time("dnsLookup", r.DNSStart),
		slog.Time("dnsDone", r.DNSDone),
		slog.Time("connectStart", r.ConnectStart),
		slog.Time("connectDone", r.ConnectDone),
		slog.Time("tlsHandshakeStart", r.TLSHandshakeStart),
		slog.Time("tlsHandshakeDone", r.TLSHandshakeDone),
		slog.Time("gotFirstResponseByte", r.GotFirstResponseByte),
		slog.Duration("dnsLookup", r.DNSLookup),
		slog.Duration("tcpConnect", r.TCPConnect),
		slog.Duration("tlsHandshake", r.TLSHandshake),
		slog.Duration("firstByte", r.FirstByte),
		slog.Duration("total", r.Total),
	)
}
