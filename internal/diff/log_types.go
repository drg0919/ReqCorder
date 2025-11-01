package diff

import "log/slog"

// Helper function to log pointers to DiffStore.
func (d *DiffStore) LogValue() slog.Value {
	if d == nil {
		return slog.StringValue("<nil>")
	}
	return slog.GroupValue(
		slog.String("recordStorePath", d.RecordStorePath),
	)
}
