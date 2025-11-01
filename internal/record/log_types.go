package record

import "log/slog"

// Helper function to log pointers to RecordStore.
func (r *RecordStore) LogValue() slog.Value {
	if r == nil {
		return slog.StringValue("<nil>")
	}
	return slog.GroupValue(
		slog.String("recordStorePath", r.RecordStorePath),
		slog.String("templateHash", r.TemplateHash),
		slog.String("requestHash", r.RequestHash),
		slog.String("responseId", r.ResponseID),
	)
}

// Helper function to log FileInfo.
func (f FileInfo) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("filePath", f.FilePath),
		slog.String("templateHash", f.TemplateHash),
		slog.String("requestHash", f.RequestHash),
		slog.String("responseId", f.ResponseID),
		slog.Time("modTime", f.ModTime),
	)
}
