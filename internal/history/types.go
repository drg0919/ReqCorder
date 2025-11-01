package history

import "time"

type HistoryStore struct {
	RecordStorePath string
}

type FileInfo struct {
	ResponseID  string
	RequestHash string
	FilePath    string
	ModTime     time.Time
}
