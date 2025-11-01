package history

import (
	"fmt"
	"log/slog"
	"reqcorder/internal/record"
	"reqcorder/pkg/utils"
	"strconv"
	"time"
)

// Retrieve sorted responses for a specific template hash with optional limit.
func (h *HistoryStore) GetSortedResponsesByTemplateHash(templateHash string, limit uint64) ([][]string, error) {
	slog.Debug("Getting sorted responses by template hash", "templateHash", templateHash, "limit", limit)
	recordStore := &record.RecordStore{
		RecordStorePath: h.RecordStorePath,
		TemplateHash:    templateHash,
	}
	responses, err := recordStore.GetSortedResponsesByTemplateHash()
	if err != nil {
		slog.Error("Failed to get sorted responses by template hash", "error", err)
		return nil, err
	}
	slog.Debug("Retrieved responses count", "count", len(responses), "templateHash", templateHash)
	if limit > 0 {
		responses = responses[:min(len(responses), int(limit))]
		slog.Debug("Applied limit to responses", "originalCount", len(responses), "limitedCount", limit)
	}
	var data [][]string
	for i, response := range responses {
		slog.Debug("Processing response", "index", i, "responseID", response.ResponseID)
		recordStore.RequestHash = response.RequestHash
		recordStore.ResponseID = response.ResponseID
		err := recordStore.GetResponse()
		if err != nil {
			slog.Error("Failed to get response", "error", err)
			return nil, err
		}
		statusIndicator := " ✅"
		if recordStore.Response.StatusCode >= 400 {
			statusIndicator = " ❌"
		}
		timestampStr := recordStore.ResponseID[:19]
		parsedTimestamp, err := time.Parse("20060102_150405_000", timestampStr)
		if err != nil {
			slog.Error("Failed to parse timestamp", "error", err)
			return nil, fmt.Errorf("%w %q: %v", ErrorFailedToParseTimestamp, timestampStr, err)
		}
		data = append(data, []string{recordStore.ResponseID, strconv.Itoa(recordStore.Response.StatusCode) + statusIndicator, recordStore.Response.Timing.Total.String(), parsedTimestamp.String()})
	}
	slog.Debug("Successfully retrieved sorted responses by template hash", "templateHash", templateHash, "dataCount", len(data))
	return data, nil
}

// Retrieve sorted responses for a specific request hash with optional limit.
func (h *HistoryStore) GetSortedResponsesByRequestHash(requestHash string, limit uint64) ([][]string, error) {
	slog.Debug("Getting sorted responses by request hash", "requestHash", requestHash, "limit", limit)
	recordStore := &record.RecordStore{
		RecordStorePath: h.RecordStorePath,
		RequestHash:     requestHash,
	}
	responses, err := recordStore.GetSortedResponsesByRequestHash()
	if err != nil {
		slog.Error("Failed to get sorted responses by request hash", "error", err)
		return nil, err
	}
	slog.Debug("Retrieved responses count", "count", len(responses), "requestHash", requestHash)
	if limit > 0 {
		responses = responses[:min(len(responses), int(limit))]
		slog.Debug("Applied limit to responses", "originalCount", len(responses), "limitedCount", limit)
	}
	var data [][]string
	for i, response := range responses {
		slog.Debug("Processing response", "index", i, "responseID", response.ResponseID)
		recordStore.RequestHash = response.RequestHash
		recordStore.ResponseID = response.ResponseID
		err := recordStore.GetResponse()
		if err != nil {
			slog.Error("Failed to get response", "error", err)
			return nil, err
		}
		statusIndicator := " ✅"
		if recordStore.Response.StatusCode >= 400 {
			statusIndicator = " ❌"
		}
		timestampStr := recordStore.ResponseID[:19]
		parsedTimestamp, err := time.Parse("20060102_150405_000", timestampStr)
		if err != nil {
			slog.Error("Failed to parse timestamp", "error", err)
			return nil, fmt.Errorf("%w %q: %v", ErrorFailedToParseTimestamp, timestampStr, err)
		}
		data = append(data, []string{recordStore.ResponseID, strconv.Itoa(recordStore.Response.StatusCode) + statusIndicator, recordStore.Response.Timing.Total.String(), parsedTimestamp.String()})
	}
	slog.Debug("Successfully retrieved sorted responses by request hash", "requestHash", requestHash, "dataCount", len(data))
	return data, nil
}

// Retrieve all responses sorted by timestamp with optional limit.
func (h *HistoryStore) GetAllResponsesSorted(limit uint64) ([][]string, error) {
	slog.Debug("Getting all responses sorted by timestamp", "limit", limit)
	recordStore := &record.RecordStore{
		RecordStorePath: h.RecordStorePath,
	}
	allFiles, err := recordStore.GetSortedResponses()
	if err != nil {
		slog.Error("Failed to get sorted responses", "error", err)
		return nil, err
	}
	slog.Debug("Retrieved responses count", "count", len(allFiles))
	if limit > 0 {
		allFiles = allFiles[:min(len(allFiles), int(limit))]
		slog.Debug("Applied limit to responses", "originalCount", len(allFiles), "limitedCount", limit)
	}
	var data [][]string
	for i, fileInfo := range allFiles {
		slog.Debug("Processing response file", "index", i, "responseID", fileInfo.ResponseID)
		recordStore.RequestHash = fileInfo.RequestHash
		recordStore.ResponseID = fileInfo.ResponseID
		err := recordStore.GetResponse()
		if err != nil {
			slog.Error("Failed to get response", "error", err)
			return nil, err
		}
		statusIndicator := " ✅"
		if recordStore.Response.StatusCode >= 400 {
			statusIndicator = " ❌"
		}
		timestampStr := recordStore.ResponseID[:19]
		parsedTimestamp, err := time.Parse("20060102_150405_000", timestampStr)
		if err != nil {
			slog.Error("Failed to parse timestamp", "error", err)
			return nil, fmt.Errorf("%w %q: %v", ErrorFailedToParseTimestamp, timestampStr, err)
		}
		data = append(data, []string{
			recordStore.ResponseID,
			strconv.Itoa(recordStore.Response.StatusCode) + statusIndicator,
			recordStore.Response.Timing.Total.String(),
			parsedTimestamp.Format("2006-01-02 15:04:05 +0000 UTC"),
		})
	}
	slog.Debug("Successfully retrieved all responses in sorted order", "dataCount", len(data))
	return data, nil
}

// Retrieve all requests sorted by modification time with optional limit.
func (h *HistoryStore) GetAllRequestsSorted(limit uint64) ([][]string, error) {
	slog.Debug("Getting all requests sorted by modification time", "limit", limit)
	recordStore := &record.RecordStore{
		RecordStorePath: h.RecordStorePath,
	}
	allFiles, err := recordStore.GetSortedRequests()
	if err != nil {
		slog.Warn("Failed to get sorted requests (may not be an error if directory is empty)", "error", err)
	}
	slog.Debug("Retrieved requests count", "count", len(allFiles))
	if limit > 0 {
		allFiles = allFiles[:min(len(allFiles), int(limit))]
		slog.Debug("Applied limit to requests", "originalCount", len(allFiles), "limitedCount", limit)
	}
	var data [][]string
	for i, fileInfo := range allFiles {
		slog.Debug("Processing request file", "index", i, "requestHash", fileInfo.RequestHash)
		recordStore.RequestHash = fileInfo.RequestHash
		err := recordStore.GetRequestByHash()
		if err != nil {
			slog.Error("Failed to get request by hash", "error", err, "requestHash", fileInfo.RequestHash)
			return nil, err
		}
		data = append(data, []string{
			recordStore.RequestHash,
			fileInfo.TemplateHash,
			fileInfo.ModTime.UTC().String(),
		})
	}
	slog.Debug("Successfully retrieved all requests sorted", "dataCount", len(data))
	return data, nil
}

// Retrieve a specific response by its ID.
func (h *HistoryStore) GetResponseByID(responseID string) (string, error) {
	slog.Debug("Getting response by ID", "responseID", responseID)
	recordStore := &record.RecordStore{
		RecordStorePath: h.RecordStorePath,
		ResponseID:      responseID,
	}
	err := recordStore.GetResponseByID()
	if err != nil {
		slog.Error("Failed to get response by ID", "error", err)
		return "", err
	}
	slog.Debug("Successfully retrieved response data", "responseID", responseID, "templateHash", recordStore.TemplateHash, "requestHash", recordStore.RequestHash)
	response, err := utils.Prettify(string(recordStore.ResponseYaml))
	if err != nil {
		slog.Error("Failed to prettify response YAML", "error", err)
		return "", err
	}
	result := fmt.Sprintf("\nTemplate Hash: %s\n", recordStore.TemplateHash)
	result += fmt.Sprintf("Request Hash: %s\n", recordStore.RequestHash)
	result += "Response:\n\n"
	result += response
	slog.Debug("Successfully formatted response", "responseID", responseID)
	return result, nil
}

// Retrieve a specific request by its hash.
func (h *HistoryStore) GetRequestByHash(requestHash string) (string, error) {
	slog.Debug("Getting request by hash", "requestHash", requestHash)
	recordStore := &record.RecordStore{
		RecordStorePath: h.RecordStorePath,
		RequestHash:     requestHash,
	}
	err := recordStore.GetRequestByHash()
	if err != nil {
		slog.Error("Failed to get request by hash", "error", err)
		return "", err
	}
	slog.Debug("Successfully retrieved request data", "requestHash", requestHash, "templateHash", recordStore.TemplateHash)
	request, err := utils.Prettify(string(recordStore.RequestYaml))
	if err != nil {
		slog.Error("Failed to prettify request YAML", "error", err)
		return "", err
	}
	result := fmt.Sprintf("\nTemplate Hash: %s\n", recordStore.TemplateHash)
	result += "Request:\n\n"
	result += request
	slog.Debug("Successfully formatted request", "requestHash", requestHash)
	return result, nil
}

// Retrieve all templates sorted by descending order of modification time with optional limit.
func (h *HistoryStore) GetAllTemplatesSorted(limit uint64) ([][]string, error) {
	slog.Debug("Getting all templates sorted by modification time", "limit", limit)
	recordStore := &record.RecordStore{
		RecordStorePath: h.RecordStorePath,
	}
	allFiles, err := recordStore.GetSortedTemplates()
	if err != nil {
		slog.Warn("Failed to get sorted templates (may not be an error if directory is empty)", "error", err)
	}
	slog.Debug("Retrieved templates count", "count", len(allFiles))
	if limit > 0 {
		allFiles = allFiles[:min(len(allFiles), int(limit))]
		slog.Debug("Applied limit to templates", "originalCount", len(allFiles), "limitedCount", limit)
	}
	var data [][]string
	for i, fileInfo := range allFiles {
		slog.Debug("Processing template file", "index", i, "templateHash", fileInfo.TemplateHash)
		data = append(data, []string{
			fileInfo.TemplateHash,
			fileInfo.ModTime.UTC().String(),
		})
	}
	slog.Debug("Successfully retrieved all templates sorted", "dataCount", len(data))
	return data, nil
}

// Retrieve a specific template by its hash.
func (h *HistoryStore) GetTemplateByHash(templateHash string) (string, error) {
	slog.Debug("Getting template by hash", "templateHash", templateHash)
	recordStore := &record.RecordStore{
		RecordStorePath: h.RecordStorePath,
		TemplateHash:    templateHash,
	}
	err := recordStore.GetTemplateByHash()
	if err != nil {
		slog.Error("Failed to get template by hash", "error", err, "templateHash", templateHash)
		return "", err
	}
	slog.Debug("Successfully retrieved template data", "templateHash", templateHash)
	template, err := utils.Prettify(string(recordStore.TemplateYaml))
	if err != nil {
		slog.Error("Failed to prettify template YAML", "error", err, "templateHash", templateHash)
		return "", err
	}
	result := "\nTemplate:\n\n"
	result += template
	slog.Debug("Successfully formatted template", "templateHash", templateHash)
	return result, nil
}
