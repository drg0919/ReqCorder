package record

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"reqcorder/internal/request"
	"reqcorder/internal/response"
	"reqcorder/pkg/utils"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var globalExecutionCounter atomic.Uint64

// Record request-response cycle.
func (r *RecordStore) Record() error {
	slog.Debug("Starting to record request-response cycle", slog.Any("recordStore", r))
	r.TemplateHash = utils.CalculateMD5Hash(r.TemplateYaml)
	r.Request.TemplateHash = r.TemplateHash
	slog.Debug("Calculated template hash", slog.String("templateHash", r.TemplateHash))
	requestYaml, err := utils.ConvertToYAML(r.Request)
	if err != nil {
		err = errors.Join(ErrorFailedToConvertRequest, err)
		slog.Error("Failed to convert request to YAML", "error", err)
		return err
	}
	r.RequestYaml = requestYaml
	r.RequestHash = utils.CalculateMD5Hash(r.RequestYaml)
	slog.Debug("Calculated request hash", slog.String("requestHash", r.RequestHash))
	r.Response.TemplateHash = r.TemplateHash
	r.Response.RequestHash = r.RequestHash
	slog.Debug("Converting response to YAML")
	r.ResponseYaml, err = utils.ConvertToYAML(r.Response)
	if err != nil {
		err = errors.Join(ErrorFailedToConvertResponse, err)
		slog.Error("Failed to convert response to YAML", "error", err)
		return err
	}
	slog.Debug("Successfully converted response to YAML")
	var wg sync.WaitGroup
	errChan := make(chan error, 3)
	wg.Go(func() {
		if err := r.recordTemplate(); err != nil {
			errChan <- err
		}
	})
	wg.Go(func() {
		if err := r.recordRequest(); err != nil {
			errChan <- err
		}
	})
	wg.Go(func() {
		if err := r.recordResponse(); err != nil {
			errChan <- err
		}
	})
	wg.Wait()
	close(errChan)
	for err := range errChan {
		slog.Error("Failed to record artifact", "error", err)
		return err
	}
	slog.Debug("Successfully recorded all artifacts")
	return nil
}

// Retrieve response from YAML file.
func (r *RecordStore) GetResponse() error {
	slog.Debug("Starting to retrieve response", slog.String("requestHash", r.RequestHash), slog.String("responseId", r.ResponseID))
	responsesDir := filepath.Join(r.RecordStorePath, "responses", r.RequestHash)
	slog.Debug("Responses directory", slog.String("responsesDir", responsesDir))
	responsePath := filepath.Join(responsesDir, r.ResponseID+".yaml")
	slog.Debug("Response file path", slog.String("responsePath", responsePath))
	var res response.ResponseObject
	slog.Debug("Reading response YAML file", slog.String("responsePath", responsePath))
	err := utils.ReadYAMLFile(responsePath, &res)
	if err != nil {
		err = errors.Join(ErrorFailedToGetResponse, err)
		slog.Error("Failed to read response YAML file", "error", err)
		return fmt.Errorf("failed to get response from path %q: %w", responsePath, err)
	}
	r.Response = &res
	slog.Debug("Successfully retrieved response", slog.Any("responseObject", &res))
	return nil
}

// Get sorted responses for current request hash.
func (r *RecordStore) GetSortedResponsesByRequestHash() ([]FileInfo, error) {
	slog.Debug("Starting to get sorted responses by request hash", slog.String("requestHash", r.RequestHash))
	responsesDir := filepath.Join(r.RecordStorePath, "responses", r.RequestHash)
	slog.Debug("Responses directory", slog.String("responsesDir", responsesDir))
	info, err := os.Stat(responsesDir)
	if err != nil {
		slog.Error("Failed to read directory", "error", err)
		return nil, fmt.Errorf("%w %q: %v", ErrorFailedToStatPath, responsesDir, err)
	}
	if !info.IsDir() {
		err = fmt.Errorf("%w %q", ErrorPathIsNotDirectory, responsesDir)
		slog.Error("Path was not a directory", "error", err)
		return nil, err
	}
	files, err := os.ReadDir(responsesDir)
	if err != nil {
		return nil, fmt.Errorf("%w %q: %v", ErrorFailedToReadDirectory, responsesDir, err)
	}
	var allFiles []FileInfo
	for _, response := range files {
		if !response.IsDir() && strings.HasSuffix(response.Name(), ".yaml") {
			responseID := strings.TrimSuffix(response.Name(), ".yaml")
			filePath := filepath.Join(responsesDir, response.Name())
			fileInfo, err := os.Stat(filePath)
			if err != nil {
				return nil, fmt.Errorf("%w %q: %v", ErrorFailedToStatPath, filePath, err)
			}
			allFiles = append(allFiles, FileInfo{
				RequestHash: r.RequestHash,
				ResponseID:  responseID,
				FilePath:    filePath,
				ModTime:     fileInfo.ModTime(),
			})
		}
	}
	sortFilesByTimeInPlace(allFiles)
	return allFiles, nil
}

// Get sorted responses for current template hash.
func (r *RecordStore) GetSortedResponsesByTemplateHash() ([]FileInfo, error) {
	slog.Debug("Starting to get sorted responses by template hash", slog.String("templateHash", r.TemplateHash))
	requestsDir := filepath.Join(r.RecordStorePath, "requests", r.TemplateHash)
	slog.Debug("Requests directory for template hash", slog.String("requestsDir", requestsDir))
	info, err := os.Stat(requestsDir)
	if err != nil {
		slog.Error("Failed to stat requests directory", "error", err)
		return nil, fmt.Errorf("%w %q: %v", ErrorFailedToStatPath, requestsDir, err)
	}
	if !info.IsDir() {
		slog.Error("Requests path is not a directory", "requestsDir", requestsDir)
		return nil, fmt.Errorf("%w %q: %v", ErrorPathIsNotDirectory, requestsDir, err)
	}
	requests, err := os.ReadDir(requestsDir)
	if err != nil {
		slog.Error("Failed to read requests directory", "error", err)
		return nil, fmt.Errorf("%w %q: %v", ErrorFailedToReadDirectory, requestsDir, err)
	}
	slog.Debug("Found requests", slog.Int("count", len(requests)))
	var allFiles []FileInfo
	for _, request := range requests {
		if !request.IsDir() && strings.HasSuffix(request.Name(), ".yaml") {
			requestHash := strings.TrimSuffix(request.Name(), ".yaml")
			slog.Debug("Processing request hash", slog.String("requestHash", requestHash))
			responsesDir := filepath.Join(r.RecordStorePath, "responses", requestHash)
			slog.Debug("Responses directory for request hash", slog.String("responsesDir", responsesDir))
			info, err = os.Stat(responsesDir)
			if err != nil {
				slog.Error("Failed to stat responses directory", "error", err)
				return nil, fmt.Errorf("%w %q: %v", ErrorFailedToStatPath, responsesDir, err)
			}
			if !info.IsDir() {
				slog.Error("Responses path is not a directory", "responsesDir", responsesDir)
				return nil, fmt.Errorf("%w %q: %v", ErrorPathIsNotDirectory, responsesDir, err)
			}
			files, err := os.ReadDir(responsesDir)
			if err != nil {
				slog.Error("Failed to read responses directory", "error", err)
				return nil, fmt.Errorf("%w %q: %v", ErrorFailedToReadDirectory, responsesDir, err)
			}
			slog.Debug("Found response files", slog.Int("count", len(files)), slog.String("requestHash", requestHash))
			for _, file := range files {
				if !file.IsDir() && strings.HasSuffix(file.Name(), ".yaml") {
					responseID := strings.TrimSuffix(file.Name(), ".yaml")
					filePath := filepath.Join(responsesDir, file.Name())
					fileInfo, err := os.Stat(filePath)
					if err != nil {
						slog.Error("Failed to stat response file", "error", err)
						return nil, fmt.Errorf("%w %q: %v", ErrorFailedToStatPath, filePath, err)
					}
					allFiles = append(allFiles, FileInfo{
						RequestHash:  requestHash,
						TemplateHash: r.TemplateHash,
						ResponseID:   responseID,
						FilePath:     filePath,
						ModTime:      fileInfo.ModTime(),
					})
				}
			}
		}
	}
	slog.Debug("Collected all response files", slog.Int("totalCount", len(allFiles)))
	sortFilesByTimeInPlace(allFiles)
	return allFiles, nil
}

// Locate and retrieve response by ID.
func (r *RecordStore) GetResponseByID() error {
	slog.Debug("Starting to locate and retrieve response by ID", slog.String("responseId", r.ResponseID))
	responsesRootDir := filepath.Join(r.RecordStorePath, "responses")

	info, err := os.Stat(responsesRootDir)
	if err != nil {
		return fmt.Errorf("%w %q: %v", ErrorFailedToStatPath, responsesRootDir, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("%w %q: %v", ErrorPathIsNotDirectory, responsesRootDir, err)
	}

	requestDirs, err := os.ReadDir(responsesRootDir)
	if err != nil {
		slog.Error("Failed to read responses root directory", "error", err)
		return fmt.Errorf("%w %q: %v", ErrorFailedToReadDirectory, responsesRootDir, err)
	}

	slog.Debug("Searching for response ID across all request directories")
	found := false
	for _, requestDir := range requestDirs {
		if requestDir.IsDir() {
			responsePath := filepath.Join(responsesRootDir, requestDir.Name(), r.ResponseID+".yaml")
			if _, err := os.Stat(responsePath); err == nil {
				slog.Debug("Found response file", slog.String("responsePath", responsePath))
				var res response.ResponseObject
				err := utils.ReadYAMLFile(responsePath, &res)
				if err != nil {
					err = errors.Join(ErrorFailedToGetResponse, err)
					slog.Error("Failed to read response YAML file", "error", err)
					return fmt.Errorf("failed to get response from path %q: %w", responsePath, err)
				}
				r.Response = &res
				r.RequestHash = r.Response.RequestHash
				r.TemplateHash = r.Response.TemplateHash
				r.ResponseYaml, _ = utils.ConvertToYAML(r.Response)
				slog.Debug("Successfully retrieved response by ID", slog.Any("responseObject", &res))
				found = true
				return nil
			}
		}
	}

	if !found {
		slog.Debug("Response ID not found", slog.String("responseId", r.ResponseID))
	}
	return fmt.Errorf("%w: %q not found", ErrorFailedToGetResponse, r.ResponseID)
}

// Retrieve request by hash.
func (r *RecordStore) GetRequestByHash() error {
	slog.Debug("Starting to retrieve request by hash", slog.String("requestHash", r.RequestHash))
	requestsRootDir := filepath.Join(r.RecordStorePath, "requests")

	info, err := os.Stat(requestsRootDir)
	if err != nil {
		return fmt.Errorf("%w %q: %v", ErrorFailedToStatPath, requestsRootDir, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("%w %q: %v", ErrorPathIsNotDirectory, requestsRootDir, err)
	}

	templateDirs, err := os.ReadDir(requestsRootDir)
	if err != nil {
		slog.Error("Failed to read requests root directory", "error", err)
		return fmt.Errorf("%w %q: %v", ErrorFailedToReadDirectory, requestsRootDir, err)
	}

	slog.Debug("Searching for request hash across all template directories")
	found := false
	for _, templateDir := range templateDirs {
		if templateDir.IsDir() {
			requestPath := filepath.Join(requestsRootDir, templateDir.Name(), r.RequestHash+".yaml")

			if _, err := os.Stat(requestPath); err == nil {
				slog.Debug("Found request file", slog.String("requestPath", requestPath))
				var req request.RequestObject
				err := utils.ReadYAMLFile(requestPath, &req)
				if err != nil {
					err = errors.Join(ErrorFailedToGetRequest, err)
					slog.Error("Failed to read request YAML file", "error", err)
					return fmt.Errorf("failed to get request from path %q: %w", requestPath, err)
				}
				r.Request = &req
				r.TemplateHash = r.Request.TemplateHash
				r.RequestYaml, _ = utils.ConvertToYAML(r.Request)
				slog.Debug("Successfully retrieved request by hash", slog.Any("requestObject", &req))
				found = true
				return nil
			}
		}
	}

	if !found {
		slog.Debug("Request hash not found", slog.String("requestHash", r.RequestHash))
	}
	return fmt.Errorf("%w: %q not found", ErrorFailedToGetRequest, r.RequestHash)
}

// Retrieve template by hash.
func (r *RecordStore) GetTemplateByHash() error {
	slog.Debug("Starting to retrieve template by hash", slog.String("templateHash", r.TemplateHash))
	templatesRootDir := filepath.Join(r.RecordStorePath, "templates")
	slog.Debug("Templates root directory", slog.String("templatesRootDir", templatesRootDir))
	info, err := os.Stat(templatesRootDir)
	if err != nil {
		slog.Error("Failed to stat templates root directory", "error", err)
		return fmt.Errorf("%w %q: %v", ErrorFailedToStatPath, templatesRootDir, err)
	}
	if !info.IsDir() {
		slog.Error("Templates root path is not a directory", "templatesRootDir", templatesRootDir)
		return fmt.Errorf("%w %q: %v", ErrorPathIsNotDirectory, templatesRootDir, err)
	}
	templatePath := filepath.Join(templatesRootDir, r.TemplateHash+".yaml")
	slog.Debug("Looking for template file", slog.String("templatePath", templatePath))
	if _, err := os.Stat(templatePath); err == nil {
		slog.Debug("Template file found, reading content")
		templateYaml, err := utils.ReadFile(templatePath)
		if err != nil {
			err = errors.Join(ErrorFailedToGetTemplate, err)
			slog.Error("Failed to read template file", "error", err)
			return fmt.Errorf("failed to get template from path %q: %w", templatePath, err)
		}
		r.TemplateYaml = templateYaml
		slog.Debug("Successfully retrieved template", slog.String("templateHash", r.TemplateHash))
		return nil
	}
	slog.Debug("Template file not found", slog.String("templateHash", r.TemplateHash))
	return fmt.Errorf("%w: %q not found", ErrorFailedToGetTemplate, r.TemplateHash)
}

// Get all responses sorted in descending order of modification.
func (r *RecordStore) GetSortedResponses() ([]FileInfo, error) {
	slog.Debug("Starting to get all sorted responses")
	responsesRootDir := filepath.Join(r.RecordStorePath, "responses")
	slog.Debug("Responses root directory", slog.String("responsesRootDir", responsesRootDir))
	var allFiles []FileInfo
	requestDirs, err := os.ReadDir(responsesRootDir)
	if err != nil {
		slog.Error("Failed to read responses root directory", "error", err)
		return nil, fmt.Errorf("%w %q: %v", ErrorFailedToReadDirectory, responsesRootDir, err)
	}
	slog.Debug("Found request directories", slog.Int("count", len(requestDirs)))
	for _, requestDir := range requestDirs {
		if requestDir.IsDir() {
			requestDirPath := filepath.Join(responsesRootDir, requestDir.Name())
			slog.Debug("Processing request directory", slog.String("requestDirPath", requestDirPath))
			responseFiles, err := os.ReadDir(requestDirPath)
			if err != nil {
				slog.Error("Failed to read request directory", "error", err)
				return nil, fmt.Errorf("%w %q: %v", ErrorFailedToReadDirectory, requestDirPath, err)
			}

			slog.Debug("Found response files", slog.Int("count", len(responseFiles)))
			for _, responseFile := range responseFiles {
				if !responseFile.IsDir() && strings.HasSuffix(responseFile.Name(), ".yaml") {
					responseID := strings.TrimSuffix(responseFile.Name(), ".yaml")
					slog.Debug("Processing response file", slog.String("responseID", responseID))
					filePath := filepath.Join(requestDirPath, responseFile.Name())

					fileInfo, err := os.Stat(filePath)
					if err != nil {
						return nil, fmt.Errorf("%w %q: %v", ErrorFailedToStatPath, filePath, err)
					}
					allFiles = append(allFiles, FileInfo{
						ResponseID:  responseID,
						RequestHash: requestDir.Name(),
						FilePath:    filePath,
						ModTime:     fileInfo.ModTime(),
					})
				}
			}
		}
	}
	sortFilesByTimeInPlace(allFiles)
	return allFiles, nil
}

// Get all requests sorted in descending order of modification.
func (r *RecordStore) GetSortedRequests() ([]FileInfo, error) {
	slog.Debug("Starting to get all sorted requests")
	requestsRootDir := filepath.Join(r.RecordStorePath, "requests")
	slog.Debug("Requests root directory", slog.String("requestsRootDir", requestsRootDir))
	var allFiles []FileInfo
	templateDirs, err := os.ReadDir(requestsRootDir)
	if err != nil {
		slog.Error("Failed to read requests root directory", "error", err)
		return nil, fmt.Errorf("%w %q: %v", ErrorFailedToReadDirectory, requestsRootDir, err)
	}
	slog.Debug("Found template directories", slog.Int("count", len(templateDirs)))
	for _, templateDir := range templateDirs {
		if templateDir.IsDir() {
			requestDirPath := filepath.Join(requestsRootDir, templateDir.Name())
			slog.Debug("Processing template directory", slog.String("templateDir", templateDir.Name()))
			responseFiles, err := os.ReadDir(requestDirPath)
			if err != nil {
				slog.Error("Failed to read request directory", "error", err)
				return nil, fmt.Errorf("%w %q: %v", ErrorFailedToReadDirectory, requestDirPath, err)
			}
			slog.Debug("Found request files", slog.Int("count", len(responseFiles)), slog.String("templateDir", templateDir.Name()))
			for _, responseFile := range responseFiles {
				if !responseFile.IsDir() && strings.HasSuffix(responseFile.Name(), ".yaml") {
					requestHash := strings.TrimSuffix(responseFile.Name(), ".yaml")
					slog.Debug("Processing request file", slog.String("requestHash", requestHash))
					filePath := filepath.Join(requestDirPath, responseFile.Name())
					fileInfo, err := os.Stat(filePath)
					if err != nil {
						return nil, fmt.Errorf("%w %q: %v", ErrorFailedToStatPath, filePath, err)
					}
					allFiles = append(allFiles, FileInfo{
						RequestHash:  requestHash,
						TemplateHash: templateDir.Name(),
						FilePath:     filePath,
						ModTime:      fileInfo.ModTime(),
					})
				}
			}
		}
	}
	sortFilesByTimeInPlace(allFiles)
	return allFiles, nil
}

// Get all templates sorted in descending order of modification.
func (r *RecordStore) GetSortedTemplates() ([]FileInfo, error) {
	slog.Debug("Starting to get all sorted templates")
	templatesRootDir := filepath.Join(r.RecordStorePath, "templates")
	slog.Debug("Templates root directory", slog.String("templatesRootDir", templatesRootDir))
	var allFiles []FileInfo
	templateItems, err := os.ReadDir(templatesRootDir)
	if err != nil {
		slog.Error("Failed to read templates root directory", "error", err)
		return nil, fmt.Errorf("%w %q: %v", ErrorFailedToReadDirectory, templatesRootDir, err)
	}
	slog.Debug("Found template files", slog.Int("count", len(templateItems)))
	for _, templateItem := range templateItems {
		if !templateItem.IsDir() && strings.HasSuffix(templateItem.Name(), ".yaml") {
			templateHash := strings.TrimSuffix(templateItem.Name(), ".yaml")
			slog.Debug("Processing template file", slog.String("templateHash", templateHash))
			filePath := filepath.Join(templatesRootDir, templateItem.Name())
			fileInfo, err := os.Stat(filePath)
			if err != nil {
				return nil, fmt.Errorf("%w %q: %v", ErrorFailedToStatPath, filePath, err)
			}
			allFiles = append(allFiles, FileInfo{
				TemplateHash: templateHash,
				FilePath:     filePath,
				ModTime:      fileInfo.ModTime(),
			})
		}
	}
	sortFilesByTimeInPlace(allFiles)
	return allFiles, nil
}

// Generate unique response ID.
func (r *RecordStore) generateResponseID() string {
	now := time.Now().UTC()
	counter := globalExecutionCounter.Add(1)
	return fmt.Sprintf("%s_%04d", now.Format("20060102_150405_000"), counter%10000)
}

// Sort files by modification time.
func sortFilesByTimeInPlace(files []FileInfo) {
	sort.Slice(files, func(i, j int) bool {
		return files[i].ModTime.After(files[j].ModTime)
	})
}

// Write template YAML to templates directory.
func (r *RecordStore) recordTemplate() error {
	slog.Debug("Starting to record template", slog.String("templateHash", r.TemplateHash))
	templateDir := filepath.Join(r.RecordStorePath, "templates")
	slog.Debug("Templates directory", slog.String("templateDir", templateDir))
	if err := utils.EnsureDir(templateDir); err != nil {
		err = errors.Join(ErrorFailedToRecord, err)
		slog.Error("Failed to ensure templates directory", "error", err)
		return err
	}
	templatePath := filepath.Join(templateDir, r.TemplateHash+".yaml")
	slog.Debug("Writing template file", slog.String("templatePath", templatePath))
	if err := os.WriteFile(templatePath, r.TemplateYaml, 0644); err != nil {
		slog.Error("Failed to write template file", "error", err)
		return fmt.Errorf("%w: %v", ErrorFailedToRecord, err)
	}
	slog.Debug("Successfully recorded template", slog.String("templateHash", r.TemplateHash))
	return nil
}

// Write request YAML to requests directory.
func (r *RecordStore) recordRequest() error {
	slog.Debug("Starting to record request", slog.String("requestHash", r.RequestHash), slog.String("templateHash", r.TemplateHash))
	requestDir := filepath.Join(r.RecordStorePath, "requests", r.TemplateHash)
	slog.Debug("Request directory", slog.String("requestDir", requestDir))
	if err := utils.EnsureDir(requestDir); err != nil {
		err = errors.Join(ErrorFailedToRecord, err)
		slog.Error("Failed to ensure request directory", "error", err)
		return err
	}
	requestPath := filepath.Join(requestDir, r.RequestHash+".yaml")
	slog.Debug("Writing request file", slog.String("requestPath", requestPath))
	if err := os.WriteFile(requestPath, r.RequestYaml, 0644); err != nil {
		slog.Error("Failed to write request file", "error", err)
		return fmt.Errorf("%w: %v", ErrorFailedToRecord, err)
	}
	slog.Debug("Successfully recorded request", slog.String("requestHash", r.RequestHash))
	return nil
}

// Write response YAML to responses directory.
func (r *RecordStore) recordResponse() error {
	slog.Debug("Starting to record response", slog.String("requestHash", r.RequestHash))
	responseDir := filepath.Join(r.RecordStorePath, "responses", r.RequestHash)
	slog.Debug("Response directory", slog.String("responseDir", responseDir))
	if err := utils.EnsureDir(responseDir); err != nil {
		err = errors.Join(ErrorFailedToRecord, err)
		slog.Error("Failed to ensure response directory", "error", err)
		return err
	}
	responseID := r.generateResponseID()
	r.ResponseID = responseID
	slog.Debug("Generated response ID", slog.String("responseId", responseID))
	responsePath := filepath.Join(responseDir, responseID+".yaml")
	slog.Debug("Writing response file", slog.String("responsePath", responsePath))
	if err := os.WriteFile(responsePath, r.ResponseYaml, 0644); err != nil {
		slog.Error("Failed to write response file", "error", err)
		return fmt.Errorf("%w: %v", ErrorFailedToRecord, err)
	}
	slog.Debug("Successfully recorded response", slog.String("responseId", responseID))
	return nil
}
