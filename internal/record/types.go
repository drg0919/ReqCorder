package record

import (
	"reqcorder/internal/request"
	"reqcorder/internal/response"
	"time"
)

// RecordStore holds all data for a recorded request-response cycle.
type RecordStore struct {
	RecordStorePath string
	TemplateYaml    []byte
	RequestYaml     []byte
	ResponseYaml    []byte
	Request         *request.RequestObject
	Response        *response.ResponseObject
	TemplateHash    string
	RequestHash     string
	ResponseID      string
}

// FileInfo contains metadata about a recorded file.
type FileInfo struct {
	ResponseID   string
	RequestHash  string
	TemplateHash string
	FilePath     string
	ModTime      time.Time
}
