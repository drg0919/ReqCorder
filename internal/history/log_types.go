package history

import "log/slog"

// Helper function to log pointers to HistoryStore.
func (h *HistoryStore) LogValue() slog.Value {
	if h == nil {
		return slog.StringValue("<nil>")
	}
	return slog.GroupValue(
		slog.String("recordStorePath", h.RecordStorePath),
	)
}

// Helper function to log FileInfo.
func (f FileInfo) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("filePath", f.FilePath),
		slog.String("requestHash", f.RequestHash),
		slog.String("responseId", f.ResponseID),
		slog.Time("modTime", f.ModTime),
	)
}
