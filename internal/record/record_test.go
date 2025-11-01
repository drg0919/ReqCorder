package record

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"reqcorder/internal/request"
	"reqcorder/internal/response"
	"reqcorder/pkg/utils"
	"testing"
	"time"

	"log/slog"
)

func TestSuccessfulGetSortedResponses(t *testing.T) {
	root := t.TempDir()
	recordOne := &RecordStore{
		RecordStorePath: root,
		RequestHash:     "requestOne",
		ResponseYaml:    []byte("key: value"),
	}
	recordTwo := &RecordStore{
		RecordStorePath: root,
		RequestHash:     "requestTwo",
		ResponseYaml:    []byte("key: value"),
	}
	err := recordOne.recordResponse()
	if err != nil {
		t.Fatalf("Expected no error, received %v", err)
	}
	err = recordTwo.recordResponse()
	if err != nil {
		t.Fatalf("Expected no error, received %v", err)
	}
	files, err := recordTwo.GetSortedResponses()
	if err != nil {
		t.Fatalf("Expected no error, received %v", err)
	}
	n := len(files)
	for i := range n - 1 {
		if files[i+1].ModTime.After(files[i].ModTime) {
			t.Fatalf("Incorrect sorting order obtained")
		}
	}
}

func TestFailedGetSortedResponses_StoreFailureRoot(t *testing.T) {
	root := t.TempDir()
	recordOne := &RecordStore{
		RecordStorePath: root,
		RequestHash:     "requestOne",
		ResponseYaml:    []byte("key: value"),
	}
	_, err := recordOne.GetSortedResponses()
	if err == nil {
		t.Fatalf("Expected error, received nil")
	}
	expectedErr := ErrorFailedToReadDirectory
	if !errors.Is(err, expectedErr) {
		t.Fatalf("Expected error %v, received %v", expectedErr, err)
	}
}

func TestFailedGetSortedResponses_StoreFailureSub(t *testing.T) {
	root := t.TempDir()
	recordOne := &RecordStore{
		RecordStorePath: root,
		RequestHash:     "requestOne",
		ResponseYaml:    []byte("key: value"),
	}
	recordOne.recordResponse()
	os.Chmod(filepath.Join(root, "responses", recordOne.RequestHash), 0o311)
	t.Cleanup(func() {
		_ = os.Chmod(filepath.Join(root, "responses", recordOne.RequestHash), 0o755)
	})
	_, err := recordOne.GetSortedResponses()
	if err == nil {
		t.Fatalf("Expected error, received nil")
	}
	expectedErr := ErrorFailedToReadDirectory
	if !errors.Is(err, expectedErr) {
		t.Fatalf("Expected error %v, received %v", expectedErr, err)
	}
}

func TestFailedGetSortedResponses_StatFailure(t *testing.T) {
	root := t.TempDir()
	recordOne := &RecordStore{
		RecordStorePath: root,
		RequestHash:     "requestOne",
		ResponseYaml:    []byte("key: value"),
	}
	err := recordOne.recordResponse()
	if err != nil {
		t.Fatalf("Expected no error, received %v", err)
	}
	responseDir := filepath.Join(root, "responses", recordOne.RequestHash)
	os.Chmod(responseDir, 0644)
	t.Cleanup(func() {
		_ = os.Chmod(responseDir, 0o755)
	})
	_, err = recordOne.GetSortedResponses()
	if err == nil {
		t.Fatalf("Expected error, received nil")
	}
	expectedErr := ErrorFailedToStatPath
	if !errors.Is(err, expectedErr) {
		t.Fatalf("Expected error %v, received %v", expectedErr, err)
	}
}

func TestSuccessfulGetSortedRequests(t *testing.T) {
	root := t.TempDir()
	recordOne := &RecordStore{
		RecordStorePath: root,
		RequestHash:     "requestOne",
		RequestYaml:     []byte("key: value"),
		TemplateHash:    "hashOne",
	}
	recordTwo := &RecordStore{
		RecordStorePath: root,
		RequestHash:     "requestTwo",
		RequestYaml:     []byte("key: value"),
		TemplateHash:    "hashTwo",
	}
	err := recordOne.recordRequest()
	if err != nil {
		t.Fatalf("Expected no error, received %v", err)
	}
	err = recordTwo.recordRequest()
	if err != nil {
		t.Fatalf("Expected no error, received %v", err)
	}
	files, err := recordTwo.GetSortedRequests()
	if err != nil {
		t.Fatalf("Expected no error, received %v", err)
	}
	n := len(files)
	for i := range n - 1 {
		if files[i+1].ModTime.After(files[i].ModTime) {
			t.Fatalf("Incorrect sorting order obtained")
		}
	}
}

func TestFailedGetSortedRequests_StoreFailureRoot(t *testing.T) {
	root := t.TempDir()
	recordOne := &RecordStore{
		RecordStorePath: root,
		RequestHash:     "requestOne",
		RequestYaml:     []byte("key: value"),
		TemplateHash:    "hashOne",
	}
	_, err := recordOne.GetSortedRequests()
	if err == nil {
		t.Fatalf("Expected error, received nil")
	}
	expectedErr := ErrorFailedToReadDirectory
	if !errors.Is(err, expectedErr) {
		t.Fatalf("Expected error %v, received %v", expectedErr, err)
	}
}

func TestFailedGetSortedRequests_StoreFailureSub(t *testing.T) {
	root := t.TempDir()
	recordOne := &RecordStore{
		RecordStorePath: root,
		RequestHash:     "requestOne",
		RequestYaml:     []byte("key: value"),
		TemplateHash:    "hashOne",
	}
	recordOne.recordRequest()
	os.Chmod(filepath.Join(root, "requests", recordOne.TemplateHash), 0o311)
	t.Cleanup(func() {
		_ = os.Chmod(filepath.Join(root, "requests", recordOne.TemplateHash), 0o755)
	})
	_, err := recordOne.GetSortedRequests()
	if err == nil {
		t.Fatalf("Expected error, received nil")
	}
	expectedErr := ErrorFailedToReadDirectory
	if !errors.Is(err, expectedErr) {
		t.Fatalf("Expected error %v, received %v", expectedErr, err)
	}
}
func TestFailedGetSortedRequests_StatFailure(t *testing.T) {
	root := t.TempDir()
	recordOne := &RecordStore{
		RecordStorePath: root,
		TemplateHash:    "hashOne",
		TemplateYaml:    []byte("key: value"),
	}
	requestsDir := filepath.Join(root, "requests", recordOne.TemplateHash)
	err := recordOne.recordRequest()
	if err != nil {
		t.Fatalf("Expected no error, received %v", err)
	}
	os.Chmod(requestsDir, 0644)
	t.Cleanup(func() {
		_ = os.Chmod(requestsDir, 0o755)
	})
	_, err = recordOne.GetSortedRequests()
	if err == nil {
		t.Fatalf("Expected error, received nil")
	}
	expectedErr := ErrorFailedToStatPath
	if !errors.Is(err, expectedErr) {
		t.Fatalf("Expected error %v, received %v", expectedErr, err)
	}
}

func TestSuccessfulGetSortedTemplates(t *testing.T) {
	root := t.TempDir()
	recordOne := &RecordStore{
		RecordStorePath: root,
		TemplateHash:    "hashOne",
		TemplateYaml:    []byte("key: value"),
	}
	recordTwo := &RecordStore{
		RecordStorePath: root,
		TemplateHash:    "hashTwo",
		TemplateYaml:    []byte("key: value"),
	}
	err := recordOne.recordTemplate()
	if err != nil {
		t.Fatalf("Expected no error, received %v", err)
	}
	err = recordTwo.recordTemplate()
	if err != nil {
		t.Fatalf("Expected no error, received %v", err)
	}
	files, err := recordTwo.GetSortedTemplates()
	if err != nil {
		t.Fatalf("Expected no error, received %v", err)
	}
	n := len(files)
	for i := range n - 1 {
		if files[i+1].ModTime.After(files[i].ModTime) {
			t.Fatalf("Incorrect sorting order obtained")
		}
	}
}

func TestFailedGetSortedTemplates_StoreFailure(t *testing.T) {
	root := t.TempDir()
	recordOne := &RecordStore{
		RecordStorePath: root,
		TemplateHash:    "hashOne",
		TemplateYaml:    []byte("key: value"),
	}
	_, err := recordOne.GetSortedTemplates()
	if err == nil {
		t.Fatalf("Expected error, received nil")
	}
	expectedErr := ErrorFailedToReadDirectory
	if !errors.Is(err, expectedErr) {
		t.Fatalf("Expected error %v, received %v", expectedErr, err)
	}
}

func TestFailedGetSortedTemplates_StatFailure(t *testing.T) {
	root := t.TempDir()
	templatesDir := filepath.Join(root, "templates")
	recordOne := &RecordStore{
		RecordStorePath: root,
		TemplateHash:    "hashOne",
		TemplateYaml:    []byte("key: value"),
	}
	err := recordOne.recordTemplate()
	if err != nil {
		t.Fatalf("Expected no error, received %v", err)
	}
	os.Chmod(templatesDir, 0644)
	t.Cleanup(func() {
		_ = os.Chmod(templatesDir, 0o755)
	})
	_, err = recordOne.GetSortedTemplates()
	if err == nil {
		t.Fatalf("Expected error, received nil")
	}
	expectedErr := ErrorFailedToStatPath
	if !errors.Is(err, expectedErr) {
		t.Fatalf("Expected error %v, received %v", expectedErr, err)
	}
}

func TestSuccessfulRecordTemplate(t *testing.T) {
	root := t.TempDir()
	recordStore := &RecordStore{
		RecordStorePath: root,
		TemplateHash:    "hash",
		TemplateYaml:    []byte("key: value"),
	}
	err := recordStore.recordTemplate()
	if err != nil {
		t.Fatalf("Expected no error, received %v", err)
	}
	templatePath := filepath.Join(recordStore.RecordStorePath, "templates", recordStore.TemplateHash+".yaml")
	info, _ := os.Stat(templatePath)
	if info.IsDir() {
		t.Fatalf("Expected file, obtained directory %q", templatePath)
	}
	content, _ := os.ReadFile(templatePath)
	if string(content) != string(recordStore.TemplateYaml) {
		t.Fatalf("Expected %s, received %s", content, recordStore.RequestYaml)
	}
}

func TestFailedRecordTemplate_StoreFailure(t *testing.T) {
	root := t.TempDir()
	templatesDir := filepath.Join(root, "templates")
	recordStore := &RecordStore{
		RecordStorePath: root,
		RequestHash:     "hash",
		RequestYaml:     []byte("key: value"),
	}
	os.Chmod(root, 0o555)
	os.Chmod(templatesDir, 0o555)
	t.Cleanup(func() {
		_ = os.Chmod(root, 0o755)
		_ = os.Chmod(templatesDir, 0o755)
	})
	expectedError := ErrorFailedToRecord
	err := recordStore.recordTemplate()
	if err == nil {
		t.Fatal("Expected error, received nil")
	}
	if !errors.Is(err, expectedError) {
		t.Fatalf("Expected error %v, received %v", expectedError, err)
	}
}

func TestFailedRecordTemplate_WriteFailure(t *testing.T) {
	root := t.TempDir()
	templatesDir := filepath.Join(root, "templates")
	recordStore := &RecordStore{
		RecordStorePath: root,
		RequestHash:     "hash",
		RequestYaml:     []byte("key: value"),
	}
	os.Mkdir(templatesDir, 0o755)
	os.Chmod(templatesDir, 0o555)
	t.Cleanup(func() {
		_ = os.Chmod(root, 0o755)
		_ = os.Chmod(templatesDir, 0o755)
	})
	expectedError := ErrorFailedToRecord
	err := recordStore.recordTemplate()
	if err == nil {
		t.Fatal("Expected error, received nil")
	}
	if !errors.Is(err, expectedError) {
		t.Fatalf("Expected error %v, received %v", expectedError, err)
	}
}

func TestSuccessfulRecordRequest(t *testing.T) {
	root := t.TempDir()
	recordStore := &RecordStore{
		RecordStorePath: root,
		RequestHash:     "hash",
		RequestYaml:     []byte("key: value"),
	}
	err := recordStore.recordRequest()
	if err != nil {
		t.Fatalf("Expected no error, received %v", err)
	}
	requestPath := filepath.Join(recordStore.RecordStorePath, "requests", recordStore.RequestHash+".yaml")
	info, _ := os.Stat(requestPath)
	if info.IsDir() {
		t.Fatalf("Expected file, obtained directory %q", requestPath)
	}
	content, _ := os.ReadFile(requestPath)
	if string(content) != string(recordStore.RequestYaml) {
		t.Fatalf("Expected %s, received %s", content, recordStore.RequestYaml)
	}
}
func TestFailedRecordRequest_StoreFailure(t *testing.T) {
	root := t.TempDir()
	requestsDir := filepath.Join(root, "requests")
	recordStore := &RecordStore{
		RecordStorePath: root,
		RequestHash:     "hash",
		RequestYaml:     []byte("key: value"),
	}
	os.Chmod(root, 0o555)
	os.Chmod(requestsDir, 0o555)
	t.Cleanup(func() {
		_ = os.Chmod(root, 0o755)
		_ = os.Chmod(requestsDir, 0o755)
	})
	expectedError := ErrorFailedToRecord
	err := recordStore.recordRequest()
	if err == nil {
		t.Fatal("Expected error, received nil")
	}
	if !errors.Is(err, expectedError) {
		t.Fatalf("Expected error %v, received %v", expectedError, err)
	}
}

func TestFailedRecordRequest_WriteFailure(t *testing.T) {
	root := t.TempDir()
	requestsDir := filepath.Join(root, "requests")
	recordStore := &RecordStore{
		RecordStorePath: root,
		RequestHash:     "hash",
		RequestYaml:     []byte("key: value"),
	}
	os.Mkdir(requestsDir, 0o755)
	os.Chmod(requestsDir, 0o555)
	t.Cleanup(func() {
		_ = os.Chmod(root, 0o755)
		_ = os.Chmod(requestsDir, 0o755)
	})
	expectedError := ErrorFailedToRecord
	err := recordStore.recordRequest()
	if err == nil {
		t.Fatal("Expected error, received nil")
	}
	if !errors.Is(err, expectedError) {
		t.Fatalf("Expected error %v, received %v", expectedError, err)
	}
}

func TestSuccessfulRecordResponse(t *testing.T) {
	root := t.TempDir()
	recordStore := &RecordStore{
		RecordStorePath: root,
		ResponseID:      "id",
		ResponseYaml:    []byte("key: value"),
	}
	err := recordStore.recordResponse()
	if err != nil {
		t.Fatalf("Expected no error, received %v", err)
	}
	responsePath := filepath.Join(recordStore.RecordStorePath, "responses", recordStore.ResponseID+".yaml")
	info, _ := os.Stat(responsePath)
	if info.IsDir() {
		t.Fatalf("Expected file, obtained directory %q", responsePath)
	}
	content, _ := os.ReadFile(responsePath)
	if string(content) != string(recordStore.ResponseYaml) {
		t.Fatalf("Expected %s, received %s", content, recordStore.ResponseYaml)
	}
}

func TestFailedRecordResponse_StoreFailure(t *testing.T) {
	root := t.TempDir()
	responsesDir := filepath.Join(root, "responses")
	recordStore := &RecordStore{
		RecordStorePath: root,
		ResponseID:      "id",
		ResponseYaml:    []byte("key: value"),
	}
	os.Chmod(root, 0o555)
	os.Chmod(responsesDir, 0o555)
	t.Cleanup(func() {
		_ = os.Chmod(root, 0o755)
		_ = os.Chmod(responsesDir, 0o755)
	})
	expectedError := ErrorFailedToRecord
	err := recordStore.recordResponse()
	if err == nil {
		t.Fatal("Expected error, received nil")
	}
	if !errors.Is(err, expectedError) {
		t.Fatalf("Expected error %v, received %v", expectedError, err)
	}
}

func TestFailedRecordResponse_WriteFailure(t *testing.T) {
	root := t.TempDir()
	responsesDir := filepath.Join(root, "responses")
	recordStore := &RecordStore{
		RecordStorePath: root,
		ResponseID:      "id",
		ResponseYaml:    []byte("key: value"),
	}
	os.Mkdir(responsesDir, 0o755)
	os.Chmod(responsesDir, 0o555)
	t.Cleanup(func() {
		_ = os.Chmod(root, 0o755)
		_ = os.Chmod(responsesDir, 0o755)
	})
	expectedError := ErrorFailedToRecord
	err := recordStore.recordResponse()
	if err == nil {
		t.Fatal("Expected error, received nil")
	}
	if !errors.Is(err, expectedError) {
		t.Fatalf("Expected error %v, received %v", expectedError, err)
	}
}

func TestSuccessfulGetTemplateByHash(t *testing.T) {
	root := t.TempDir()
	templateHash := "hashOne"
	templateYaml := []byte("key: value")
	recordStore := &RecordStore{
		RecordStorePath: root,
		TemplateHash:    templateHash,
		TemplateYaml:    templateYaml,
	}
	err := recordStore.recordTemplate()
	if err != nil {
		t.Fatalf("Expected no error, received %v", err)
	}
	getRecord := &RecordStore{
		RecordStorePath: root,
		TemplateHash:    templateHash,
	}
	err = getRecord.GetTemplateByHash()
	if err != nil {
		t.Fatalf("Expected no error, received %v", err)
	}
	if string(getRecord.TemplateYaml) != string(templateYaml) {
		t.Fatalf("Expected template content %s, received %s", string(templateYaml), string(getRecord.TemplateYaml))
	}
}

func TestFailedGetTemplateByHash_DirectoryNotFound(t *testing.T) {
	root := t.TempDir()
	getRecord := &RecordStore{
		RecordStorePath: root,
		TemplateHash:    "nonexistent",
	}
	err := getRecord.GetTemplateByHash()
	if err == nil {
		t.Fatalf("Expected error, received nil")
	}
	expectedErr := ErrorFailedToStatPath
	if !errors.Is(err, expectedErr) {
		t.Fatalf("Expected error %v, received %v", expectedErr, err)
	}
}

func TestFailedGetTemplateByHash_PathIsNotDirectory(t *testing.T) {
	root := t.TempDir()
	templatesDir := filepath.Join(root, "templates")
	err := os.WriteFile(templatesDir, []byte("not a directory"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	getRecord := &RecordStore{
		RecordStorePath: root,
		TemplateHash:    "any",
	}
	err = getRecord.GetTemplateByHash()
	if err == nil {
		t.Fatalf("Expected error, received nil")
	}
	expectedErr := ErrorPathIsNotDirectory
	if !errors.Is(err, expectedErr) {
		t.Fatalf("Expected error %v, received %v", expectedErr, err)
	}
}

func TestFailedGetTemplateByHash_TemplateNotFound(t *testing.T) {
	root := t.TempDir()
	templatesDir := filepath.Join(root, "templates")
	err := os.Mkdir(templatesDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create templates directory: %v", err)
	}
	getRecord := &RecordStore{
		RecordStorePath: root,
		TemplateHash:    "nonexistent",
	}
	err = getRecord.GetTemplateByHash()
	if err == nil {
		t.Fatalf("Expected error, received nil")
	}
	expectedErr := ErrorFailedToGetTemplate
	if !errors.Is(err, expectedErr) {
		t.Fatalf("Expected error %v, received %v", expectedErr, err)
	}
}

func TestFailedGetTemplateByHash_ReadFailure(t *testing.T) {
	root := t.TempDir()
	templateHash := "hashOne"
	templateYaml := []byte("key: value")
	recordStore := &RecordStore{
		RecordStorePath: root,
		TemplateHash:    templateHash,
		TemplateYaml:    templateYaml,
	}
	err := recordStore.recordTemplate()
	if err != nil {
		t.Fatalf("Expected no error, received %v", err)
	}
	templatePath := filepath.Join(root, "templates", templateHash+".yaml")
	os.Chmod(templatePath, 0o222)
	t.Cleanup(func() {
		_ = os.Chmod(templatePath, 0o644)
	})
	getRecord := &RecordStore{
		RecordStorePath: root,
		TemplateHash:    templateHash,
	}
	err = getRecord.GetTemplateByHash()
	if err == nil {
		t.Fatalf("Expected error, received nil")
	}
	expectedErr := ErrorFailedToGetTemplate
	if !errors.Is(err, expectedErr) {
		t.Fatalf("Expected error %v, received %v", expectedErr, err)
	}
}

func TestSuccessfulGetRequestByHash(t *testing.T) {
	root := t.TempDir()
	requestHash := "hashOne"
	recordStore := &RecordStore{
		RecordStorePath: root,
		RequestHash:     requestHash,
		TemplateHash:    "templateHash",
		Request: &request.RequestObject{
			TemplateHash: "templateHash",
		},
	}
	recordStore.RequestYaml, _ = utils.ConvertToYAML(recordStore.Request)
	err := recordStore.recordRequest()
	if err != nil {
		t.Fatalf("Expected no error, received %v", err)
	}
	getRecord := &RecordStore{
		RecordStorePath: root,
		RequestHash:     requestHash,
	}
	err = getRecord.GetRequestByHash()
	if err != nil {
		t.Fatalf("Expected no error, received %v", err)
	}
	if getRecord.TemplateHash != "templateHash" {
		t.Fatalf("Expected template hash %s, received %s", "templateHash", getRecord.TemplateHash)
	}
}

func TestFailedGetRequestByHash_RequestNotFound(t *testing.T) {
	root := t.TempDir()
	requestsDir := filepath.Join(root, "requests")
	err := os.Mkdir(requestsDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create requests directory: %v", err)
	}
	getRecord := &RecordStore{
		RecordStorePath: root,
		RequestHash:     "nonexistent",
	}
	err = getRecord.GetRequestByHash()
	if err == nil {
		t.Fatalf("Expected error, received nil")
	}
	expectedErr := ErrorFailedToGetRequest
	if !errors.Is(err, expectedErr) {
		t.Fatalf("Expected error %v, received %v", expectedErr, err)
	}
}

func TestFailedGetRequestByHash_ReadFailure(t *testing.T) {
	root := t.TempDir()
	requestHash := "hashOne"
	requestYaml := []byte("key: value")
	recordStore := &RecordStore{
		RecordStorePath: root,
		RequestHash:     requestHash,
		RequestYaml:     requestYaml,
		TemplateHash:    "templateHash",
	}
	err := recordStore.recordRequest()
	if err != nil {
		t.Fatalf("Expected no error, received %v", err)
	}
	requestPath := filepath.Join(root, "requests", "templateHash", requestHash+".yaml")
	os.Chmod(requestPath, 0o222)
	t.Cleanup(func() {
		_ = os.Chmod(requestPath, 0o644)
	})
	getRecord := &RecordStore{
		RecordStorePath: root,
		RequestHash:     requestHash,
	}
	err = getRecord.GetRequestByHash()
	if err == nil {
		t.Fatalf("Expected error, received nil")
	}
	expectedErr := ErrorFailedToGetRequest
	if !errors.Is(err, expectedErr) {
		t.Fatalf("Expected error %v, received %v", expectedErr, err)
	}
}

func TestFailedGetRequestByHash_StatFailure(t *testing.T) {
	root := t.TempDir()
	getRecord := &RecordStore{
		RecordStorePath: root,
		RequestHash:     "any",
	}
	err := getRecord.GetRequestByHash()
	if err == nil {
		t.Fatalf("Expected error, received nil")
	}
	expectedErr := ErrorFailedToStatPath
	if !errors.Is(err, expectedErr) {
		t.Fatalf("Expected error %v, received %v", expectedErr, err)
	}
}

func TestFailedGetRequestByHash_PathIsNotDirectory(t *testing.T) {
	root := t.TempDir()
	requestsRootDir := filepath.Join(root, "requests")
	err := os.WriteFile(requestsRootDir, []byte("not a directory"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	getRecord := &RecordStore{
		RecordStorePath: root,
		RequestHash:     "any",
	}
	err = getRecord.GetRequestByHash()
	if err == nil {
		t.Fatalf("Expected error, received nil")
	}
	expectedErr := ErrorPathIsNotDirectory
	if !errors.Is(err, expectedErr) {
		t.Fatalf("Expected error %v, received %v", expectedErr, err)
	}
}

func TestFailedGetRequestByHash_ReadDirFailure(t *testing.T) {
	root := t.TempDir()
	requestsRootDir := filepath.Join(root, "requests")
	err := os.Mkdir(requestsRootDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create requests directory: %v", err)
	}
	os.Chmod(requestsRootDir, 0o222)
	t.Cleanup(func() {
		_ = os.Chmod(requestsRootDir, 0o755)
	})
	getRecord := &RecordStore{
		RecordStorePath: root,
		RequestHash:     "any",
	}
	err = getRecord.GetRequestByHash()
	if err == nil {
		t.Fatalf("Expected error, received nil")
	}
	expectedErr := ErrorFailedToReadDirectory
	if !errors.Is(err, expectedErr) {
		t.Fatalf("Expected error %v, received %v", expectedErr, err)
	}
}

func TestSuccessfulGetResponseByID(t *testing.T) {
	root := t.TempDir()
	recordStore := &RecordStore{
		RecordStorePath: root,
		RequestHash:     "requestHash",
		TemplateHash:    "templateHash",
		Response: &response.ResponseObject{
			TemplateHash: "templateHash",
			RequestHash:  "requestHash",
		},
	}
	recordStore.ResponseYaml, _ = utils.ConvertToYAML(recordStore.Response)
	err := recordStore.recordResponse()
	if err != nil {
		t.Fatalf("Expected no error, received %v", err)
	}
	getRecord := &RecordStore{
		RecordStorePath: root,
		ResponseID:      recordStore.ResponseID,
	}
	err = getRecord.GetResponseByID()
	if err != nil {
		t.Fatalf("Expected no error, received %v", err)
	}
	if getRecord.RequestHash != "requestHash" {
		t.Fatalf("Expected request hash %s, received %s", "requestHash", getRecord.RequestHash)
	}
	if getRecord.TemplateHash != "templateHash" {
		t.Fatalf("Expected template hash %s, received %s", "templateHash", getRecord.TemplateHash)
	}
}

func TestFailedGetResponseByID_ResponseNotFound(t *testing.T) {
	root := t.TempDir()
	responsesDir := filepath.Join(root, "responses")
	err := os.Mkdir(responsesDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create responses directory: %v", err)
	}
	getRecord := &RecordStore{
		RecordStorePath: root,
		ResponseID:      "nonexistent",
	}
	err = getRecord.GetResponseByID()
	if err == nil {
		t.Fatalf("Expected error, received nil")
	}
	expectedErr := ErrorFailedToGetResponse
	if !errors.Is(err, expectedErr) {
		t.Fatalf("Expected error %v, received %v", expectedErr, err)
	}
}

func TestFailedGetResponseByID_ReadFailure(t *testing.T) {
	root := t.TempDir()
	recordStore := &RecordStore{
		RecordStorePath: root,
		RequestHash:     "requestHash",
		TemplateHash:    "templateHash",
		Response: &response.ResponseObject{
			TemplateHash: "templateHash",
			RequestHash:  "requestHash",
		},
	}
	recordStore.ResponseYaml, _ = utils.ConvertToYAML(recordStore.Response)
	err := recordStore.recordResponse()
	if err != nil {
		t.Fatalf("Expected no error, received %v", err)
	}
	responsePath := filepath.Join(root, "responses", "requestHash", recordStore.ResponseID+".yaml")
	if _, err := os.Stat(responsePath); err != nil {
		t.Fatalf("Response file not found at path: %s", responsePath)
	}
	os.Chmod(responsePath, 0o222)
	t.Cleanup(func() {
		_ = os.Chmod(responsePath, 0o644)
	})
	getRecord := &RecordStore{
		RecordStorePath: root,
		ResponseID:      recordStore.ResponseID,
	}
	err = getRecord.GetResponseByID()
	if err == nil {
		t.Fatalf("Expected error, received nil")
	}
	expectedErr := ErrorFailedToGetResponse
	if !errors.Is(err, expectedErr) {
		t.Fatalf("Expected error %v, received %v", expectedErr, err)
	}
}

func TestFailedGetResponseByID_StatFailure(t *testing.T) {
	root := t.TempDir()
	getRecord := &RecordStore{
		RecordStorePath: root,
		ResponseID:      "any",
	}
	err := getRecord.GetResponseByID()
	if err == nil {
		t.Fatalf("Expected error, received nil")
	}
	expectedErr := ErrorFailedToStatPath
	if !errors.Is(err, expectedErr) {
		t.Fatalf("Expected error %v, received %v", expectedErr, err)
	}
}

func TestFailedGetResponseByID_PathIsNotDirectory(t *testing.T) {
	root := t.TempDir()
	responsesRootDir := filepath.Join(root, "responses")
	err := os.WriteFile(responsesRootDir, []byte("not a directory"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	getRecord := &RecordStore{
		RecordStorePath: root,
		ResponseID:      "any",
	}
	err = getRecord.GetResponseByID()
	if err == nil {
		t.Fatalf("Expected error, received nil")
	}
	expectedErr := ErrorPathIsNotDirectory
	if !errors.Is(err, expectedErr) {
		t.Fatalf("Expected error %v, received %v", expectedErr, err)
	}
}

func TestFailedGetResponseByID_ReadDirFailure(t *testing.T) {
	root := t.TempDir()
	responsesRootDir := filepath.Join(root, "responses")
	err := os.Mkdir(responsesRootDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create responses directory: %v", err)
	}
	os.Chmod(responsesRootDir, 0o222)
	t.Cleanup(func() {
		_ = os.Chmod(responsesRootDir, 0o755)
	})
	getRecord := &RecordStore{
		RecordStorePath: root,
		ResponseID:      "any",
	}
	err = getRecord.GetResponseByID()
	if err == nil {
		t.Fatalf("Expected error, received nil")
	}
	expectedErr := ErrorFailedToReadDirectory
	if !errors.Is(err, expectedErr) {
		t.Fatalf("Expected error %v, received %v", expectedErr, err)
	}
}

func TestSuccessfulGetSortedResponsesByRequestHash(t *testing.T) {
	root := t.TempDir()
	recordOne := &RecordStore{
		RecordStorePath: root,
		RequestHash:     "requestOne",
		ResponseYaml:    []byte("key: value"),
	}
	recordTwo := &RecordStore{
		RecordStorePath: root,
		RequestHash:     "requestOne",
		ResponseYaml:    []byte("key: value"),
	}
	err := recordOne.recordResponse()
	if err != nil {
		t.Fatalf("Expected no error, received %v", err)
	}
	err = recordTwo.recordResponse()
	if err != nil {
		t.Fatalf("Expected no error, received %v", err)
	}
	getRecord := &RecordStore{
		RecordStorePath: root,
		RequestHash:     "requestOne",
	}
	files, err := getRecord.GetSortedResponsesByRequestHash()
	if err != nil {
		t.Fatalf("Expected no error, received %v", err)
	}
	n := len(files)
	for i := range n - 1 {
		if files[i+1].ModTime.After(files[i].ModTime) {
			t.Fatalf("Incorrect sorting order obtained")
		}
	}
}

func TestFailedGetSortedResponsesByRequestHash_StoreFailure(t *testing.T) {
	root := t.TempDir()
	recordOne := &RecordStore{
		RecordStorePath: root,
		RequestHash:     "requestOne",
		ResponseYaml:    []byte("key: value"),
	}
	_, err := recordOne.GetSortedResponsesByRequestHash()
	if err == nil {
		t.Fatalf("Expected error, received nil")
	}
	expectedErr := ErrorFailedToStatPath
	if !errors.Is(err, expectedErr) {
		t.Fatalf("Expected error %v, received %v", expectedErr, err)
	}
}

func TestFailedGetSortedResponsesByRequestHash_PathIsNotDirectory(t *testing.T) {
	root := t.TempDir()
	// Create the responses directory first
	responsesRootDir := filepath.Join(root, "responses")
	err := os.Mkdir(responsesRootDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create responses directory: %v", err)
	}
	// Create the request-specific directory and make it a file instead
	requestResponsesDir := filepath.Join(responsesRootDir, "requestOne")
	err = os.WriteFile(requestResponsesDir, []byte("not a directory"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	recordOne := &RecordStore{
		RecordStorePath: root,
		RequestHash:     "requestOne",
	}
	_, err = recordOne.GetSortedResponsesByRequestHash()
	if err == nil {
		t.Fatalf("Expected error, received nil")
	}
	expectedErr := ErrorPathIsNotDirectory
	if !errors.Is(err, expectedErr) {
		t.Fatalf("Expected error %v, received %v", expectedErr, err)
	}
}

func TestFailedGetSortedResponsesByRequestHash_ReadFailure(t *testing.T) {
	root := t.TempDir()
	recordOne := &RecordStore{
		RecordStorePath: root,
		RequestHash:     "requestOne",
		ResponseYaml:    []byte("key: value"),
	}
	err := recordOne.recordResponse()
	if err != nil {
		t.Fatalf("Expected no error, received %v", err)
	}
	responsesDir := filepath.Join(root, "responses", "requestOne")
	os.Chmod(responsesDir, 0o222)
	t.Cleanup(func() {
		_ = os.Chmod(responsesDir, 0o755)
	})
	_, err = recordOne.GetSortedResponsesByRequestHash()
	if err == nil {
		t.Fatalf("Expected error, received nil")
	}
	expectedErr := ErrorFailedToReadDirectory
	if !errors.Is(err, expectedErr) {
		t.Fatalf("Expected error %v, received %v", expectedErr, err)
	}
}

func TestSuccessfulGetSortedResponsesByTemplateHash(t *testing.T) {
	root := t.TempDir()
	recordOne := &RecordStore{
		RecordStorePath: root,
		RequestHash:     "requestOne",
		TemplateHash:    "templateOne",
		ResponseYaml:    []byte("key: value"),
	}
	recordTwo := &RecordStore{
		RecordStorePath: root,
		RequestHash:     "requestTwo",
		TemplateHash:    "templateOne",
		ResponseYaml:    []byte("key: value"),
	}
	err := recordOne.recordRequest()
	if err != nil {
		t.Fatalf("Expected no error, received %v", err)
	}
	err = recordTwo.recordRequest()
	if err != nil {
		t.Fatalf("Expected no error, received %v", err)
	}
	err = recordOne.recordResponse()
	if err != nil {
		t.Fatalf("Expected no error, received %v", err)
	}
	err = recordTwo.recordResponse()
	if err != nil {
		t.Fatalf("Expected no error, received %v", err)
	}
	getRecord := &RecordStore{
		RecordStorePath: root,
		TemplateHash:    "templateOne",
	}
	files, err := getRecord.GetSortedResponsesByTemplateHash()
	if err != nil {
		t.Fatalf("Expected no error, received %v", err)
	}
	n := len(files)
	for i := range n - 1 {
		if files[i+1].ModTime.After(files[i].ModTime) {
			t.Fatalf("Incorrect sorting order obtained")
		}
	}
}

func TestFailedGetSortedResponsesByTemplateHash_StoreFailure(t *testing.T) {
	root := t.TempDir()
	recordOne := &RecordStore{
		RecordStorePath: root,
		TemplateHash:    "templateOne",
	}
	_, err := recordOne.GetSortedResponsesByTemplateHash()
	if err == nil {
		t.Fatalf("Expected error, received nil")
	}
	expectedErr := ErrorFailedToStatPath
	if !errors.Is(err, expectedErr) {
		t.Fatalf("Expected error %v, received %v", expectedErr, err)
	}
}

func TestFailedGetSortedResponsesByTemplateHash_PathIsNotDirectory(t *testing.T) {
	root := t.TempDir()
	requestsRootDir := filepath.Join(root, "requests")
	err := os.Mkdir(requestsRootDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create requests directory: %v", err)
	}
	templateRequestsDir := filepath.Join(requestsRootDir, "templateOne")
	err = os.WriteFile(templateRequestsDir, []byte("not a directory"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	recordOne := &RecordStore{
		RecordStorePath: root,
		TemplateHash:    "templateOne",
	}
	_, err = recordOne.GetSortedResponsesByTemplateHash()
	if err == nil {
		t.Fatalf("Expected error, received nil")
	}
	expectedErr := ErrorPathIsNotDirectory
	if !errors.Is(err, expectedErr) {
		t.Fatalf("Expected error %v, received %v", expectedErr, err)
	}
}

func TestFailedGetSortedResponsesByTemplateHash_ReadFailure(t *testing.T) {
	root := t.TempDir()
	recordOne := &RecordStore{
		RecordStorePath: root,
		RequestHash:     "requestOne",
		TemplateHash:    "templateOne",
		ResponseYaml:    []byte("key: value"),
	}
	err := recordOne.recordRequest()
	if err != nil {
		t.Fatalf("Expected no error, received %v", err)
	}
	requestsDir := filepath.Join(root, "requests", "templateOne")
	os.Chmod(requestsDir, 0o222)
	t.Cleanup(func() {
		_ = os.Chmod(requestsDir, 0o755)
	})
	_, err = recordOne.GetSortedResponsesByTemplateHash()
	if err == nil {
		t.Fatalf("Expected error, received nil")
	}
	expectedErr := ErrorFailedToReadDirectory
	if !errors.Is(err, expectedErr) {
		t.Fatalf("Expected error %v, received %v", expectedErr, err)
	}
}

// func TestFailedGetSortedResponsesByTemplateHash_ResponseDirStoreFailure(t *testing.T) {
// 	root := t.TempDir()
// 	// Create the responses directory first
// 	responsesRootDir := filepath.Join(root, "responses")
// 	err := os.Mkdir(responsesRootDir, 0755)
// 	if err != nil {
// 		t.Fatalf("Failed to create responses directory: %v", err)
// 	}
// 	// Create the request-specific directory
// 	requestResponsesDir := filepath.Join(responsesRootDir, "requestOne")
// 	err = os.Mkdir(requestResponsesDir, 0755)
// 	if err != nil {
// 		t.Fatalf("Failed to create request responses directory: %v", err)
// 	}
// 	// Create a response file in the request-specific directory
// 	responseFile := filepath.Join(requestResponsesDir, "response1.yaml")
// 	err = os.WriteFile(responseFile, []byte("key: value"), 0644)
// 	if err != nil {
// 		t.Fatalf("Failed to create response file: %v", err)
// 	}
// 	recordOne := &RecordStore{
// 		RecordStorePath: root,
// 		RequestHash:     "requestOne",
// 		TemplateHash:    "templateOne",
// 		ResponseYaml:    []byte("key: value"),
// 	}
// 	err = recordOne.recordRequest()
// 	if err != nil {
// 		t.Fatalf("Expected no error, received %v", err)
// 	}
// 	// Make the response directory inaccessible
// 	os.Chmod(requestResponsesDir, 0o000)
// 	t.Cleanup(func() {
// 		_ = os.Chmod(requestResponsesDir, 0o755)
// 	})
// 	_, err = recordOne.GetSortedResponsesByTemplateHash()
// 	fmt.Println(err)
// 	if err == nil {
// 		t.Fatalf("Expected error, received nil")
// 	}
// 	expectedErr := ErrorFailedToStatPath
// 	if !errors.Is(err, expectedErr) {
// 		t.Fatalf("Expected error %v, received %v", expectedErr, err)
// 	}
// }

func TestFailedGetSortedResponsesByTemplateHash_ResponseDirPathIsNotDirectory(t *testing.T) {
	root := t.TempDir()
	// Create the responses directory first
	responsesRootDir := filepath.Join(root, "responses")
	err := os.Mkdir(responsesRootDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create responses directory: %v", err)
	}
	// Create the request-specific directory as a file instead
	requestResponsesDir := filepath.Join(responsesRootDir, "requestOne")
	err = os.WriteFile(requestResponsesDir, []byte("not a directory"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	recordOne := &RecordStore{
		RecordStorePath: root,
		RequestHash:     "requestOne",
		TemplateHash:    "templateOne",
		ResponseYaml:    []byte("key: value"),
	}
	err = recordOne.recordRequest()
	if err != nil {
		t.Fatalf("Expected no error, received %v", err)
	}
	_, err = recordOne.GetSortedResponsesByTemplateHash()
	if err == nil {
		t.Fatalf("Expected error, received nil")
	}
	expectedErr := ErrorPathIsNotDirectory
	if !errors.Is(err, expectedErr) {
		t.Fatalf("Expected error %v, received %v", expectedErr, err)
	}
}

func TestFailedGetSortedResponsesByTemplateHash_ResponseDirReadFailure(t *testing.T) {
	root := t.TempDir()
	recordOne := &RecordStore{
		RecordStorePath: root,
		RequestHash:     "requestOne",
		TemplateHash:    "templateOne",
		ResponseYaml:    []byte("key: value"),
	}
	err := recordOne.recordRequest()
	if err != nil {
		t.Fatalf("Expected no error, received %v", err)
	}
	err = recordOne.recordResponse()
	if err != nil {
		t.Fatalf("Expected no error, received %v", err)
	}
	responsesDir := filepath.Join(root, "responses", "requestOne")
	os.Chmod(responsesDir, 0o222)
	t.Cleanup(func() {
		_ = os.Chmod(responsesDir, 0o755)
	})
	_, err = recordOne.GetSortedResponsesByTemplateHash()
	if err == nil {
		t.Fatalf("Expected error, received nil")
	}
	expectedErr := ErrorFailedToReadDirectory
	if !errors.Is(err, expectedErr) {
		t.Fatalf("Expected error %v, received %v", expectedErr, err)
	}
}

func TestSuccessfulGetResponse(t *testing.T) {
	root := t.TempDir()
	recordStore := &RecordStore{
		RecordStorePath: root,
		ResponseID:      "responseID",
		Response:        &response.ResponseObject{StatusCode: 200, RequestHash: "requestHash"},
		Request:         &request.RequestObject{URL: "https://example.com", Method: "GET", TemplateHash: "templateHash"},
	}
	err := recordStore.Record()
	if err != nil {
		t.Fatalf("Expected no error, received %v", err)
	}
	getRecord := &RecordStore{
		RecordStorePath: root,
		RequestHash:     recordStore.RequestHash,
		ResponseID:      recordStore.ResponseID,
	}
	err = getRecord.GetResponse()
	if err != nil {
		t.Fatalf("Expected no error, received %v", err)
	}
	getRecord.ResponseYaml, _ = utils.ConvertToYAML(getRecord.Response)
	if string(getRecord.ResponseYaml) != string(recordStore.ResponseYaml) {
		t.Fatalf("Expected response YAML %s, received %s", string(recordStore.ResponseYaml), string(getRecord.ResponseYaml))
	}
}

func TestFailedGetResponse_ResponseNotFound(t *testing.T) {
	root := t.TempDir()
	responsesDir := filepath.Join(root, "responses", "requestHash")
	err := os.MkdirAll(responsesDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create responses directory: %v", err)
	}
	getRecord := &RecordStore{
		RecordStorePath: root,
		RequestHash:     "requestHash",
		ResponseID:      "nonexistent",
	}
	err = getRecord.GetResponse()
	if err == nil {
		t.Fatalf("Expected error, received nil")
	}
	expectedErr := ErrorFailedToGetResponse
	if !errors.Is(err, expectedErr) {
		t.Fatalf("Expected error %v, received %v", expectedErr, err)
	}
}

func TestFailedGetResponse_ReadFailure(t *testing.T) {
	root := t.TempDir()
	recordStore := &RecordStore{
		RecordStorePath: root,
		RequestHash:     "requestHash",
		ResponseID:      "responseID",
		ResponseYaml:    []byte("key: value"),
	}
	err := recordStore.recordResponse()
	if err != nil {
		t.Fatalf("Expected no error, received %v", err)
	}
	responsePath := filepath.Join(root, "responses", "requestHash", "responseID.yaml")
	os.Chmod(responsePath, 0o222)
	t.Cleanup(func() {
		_ = os.Chmod(responsePath, 0o755)
	})
	getRecord := &RecordStore{
		RecordStorePath: root,
		RequestHash:     "requestHash",
		ResponseID:      "responseID",
	}
	err = getRecord.GetResponse()
	if err == nil {
		t.Fatalf("Expected error, received nil")
	}
	expectedErr := ErrorFailedToGetResponse
	if !errors.Is(err, expectedErr) {
		t.Fatalf("Expected error %v, received %v", expectedErr, err)
	}
}

func TestSuccessfulRecord(t *testing.T) {
	root := t.TempDir()
	recordStore := &RecordStore{
		RecordStorePath: root,
		TemplateYaml:    []byte("template: value"),
		Request: &request.RequestObject{
			TemplateHash: "templateHash",
			URL:          "https://example.com",
			Method:       "GET",
		},
		Response: &response.ResponseObject{
			TemplateHash: "templateHash",
			RequestHash:  "requestHash",
			StatusCode:   200,
		},
	}
	// Manually set the hashes since we're not using the full Record() flow
	recordStore.TemplateHash = utils.CalculateMD5Hash(recordStore.TemplateYaml)
	recordStore.RequestHash = "requestHash"

	// Record each component individually
	err := recordStore.recordTemplate()
	if err != nil {
		t.Fatalf("Expected no error, received %v", err)
	}
	err = recordStore.recordRequest()
	if err != nil {
		t.Fatalf("Expected no error, received %v", err)
	}
	err = recordStore.recordResponse()
	if err != nil {
		t.Fatalf("Expected no error, received %v", err)
	}

	// Verify template was recorded
	templatePath := filepath.Join(root, "templates", recordStore.TemplateHash+".yaml")
	if _, err := os.Stat(templatePath); err != nil {
		t.Fatalf("Expected template file to exist at %q", templatePath)
	}

	// Verify request was recorded
	requestPath := filepath.Join(root, "requests", recordStore.TemplateHash, recordStore.RequestHash+".yaml")
	if _, err := os.Stat(requestPath); err != nil {
		t.Fatalf("Expected request file to exist at %q", requestPath)
	}

	// Verify response was recorded
	responseDir := filepath.Join(root, "responses", recordStore.RequestHash)
	responseFiles, err := os.ReadDir(responseDir)
	if err != nil {
		t.Fatalf("Expected to read response directory, received %v", err)
	}
	if len(responseFiles) != 1 {
		t.Fatalf("Expected 1 response file, received %d", len(responseFiles))
	}
}

// func TestFailedRecord_ConvertRequestFailure(t *testing.T) {
// 	root := t.TempDir()
// 	recordStore := &RecordStore{
// 		RecordStorePath: root,
// 		TemplateYaml:    []byte("template: value"),
// 		Request:         &request.RequestObject{}, // Invalid request that will cause conversion to fail
// 		Response: &response.ResponseObject{
// 			TemplateHash: "templateHash",
// 			RequestHash:  "requestHash",
// 		},
// 	}
// 	// Manually set the template hash
// 	recordStore.TemplateHash = utils.CalculateMD5Hash(recordStore.TemplateYaml)
// 	// Manually convert response to YAML
// 	var err error
// 	recordStore.ResponseYaml, err = utils.ConvertToYAML(recordStore.Response)
// 	if err != nil {
// 		t.Fatalf("Failed to convert response to YAML: %v", err)
// 	}
// 	// Manually set request hash (empty request will have empty hash)
// 	recordStore.RequestHash = ""

// 	// Now try to record - this should fail during request conversion
// 	err = recordStore.Record()
// 	if err == nil {
// 		t.Fatalf("Expected error, received nil")
// 	}
// 	expectedErr := ErrorFailedToConvertRequest
// 	if !errors.Is(err, expectedErr) {
// 		t.Fatalf("Expected error %v, received %v", expectedErr, err)
// 	}
// }

// func TestFailedRecord_ConvertResponseFailure(t *testing.T) {
// 	root := t.TempDir()
// 	recordStore := &RecordStore{
// 		RecordStorePath: root,
// 		TemplateYaml:    []byte("template: value"),
// 		Request: &request.RequestObject{
// 			TemplateHash: "templateHash",
// 			URL:          "https://example.com",
// 			Method:       "GET",
// 		},
// 		Response: &response.ResponseObject{}, // Invalid response that will cause conversion to fail
// 	}
// 	// Manually set the template hash
// 	recordStore.TemplateHash = utils.CalculateMD5Hash(recordStore.TemplateYaml)
// 	// Manually convert request to YAML
// 	var err error
// 	recordStore.RequestYaml, err = utils.ConvertToYAML(recordStore.Request)
// 	if err != nil {
// 		t.Fatalf("Failed to convert request to YAML: %v", err)
// 	}
// 	// Manually set request hash
// 	recordStore.RequestHash = utils.CalculateMD5Hash(recordStore.RequestYaml)

// 	// Now try to record - this should fail during response conversion
// 	err = recordStore.Record()
// 	if err == nil {
// 		t.Fatalf("Expected error, received nil")
// 	}
// 	expectedErr := ErrorFailedToConvertResponse
// 	if !errors.Is(err, expectedErr) {
// 		t.Fatalf("Expected error %v, received %v", expectedErr, err)
// 	}
// }

func TestLogValue(t *testing.T) {
	tests := []struct {
		name     string
		input    *RecordStore
		expected slog.Value
	}{
		{
			name:     "Nil pointer",
			input:    nil,
			expected: slog.StringValue("<nil>"),
		},
		{
			name: "Valid record store with all fields",
			input: &RecordStore{
				RecordStorePath: "/test/path",
				TemplateHash:    "template123",
				RequestHash:     "request456",
				ResponseID:      "response789",
			},
			expected: slog.GroupValue(
				slog.String("recordStorePath", "/test/path"),
				slog.String("templateHash", "template123"),
				slog.String("requestHash", "request456"),
				slog.String("responseId", "response789"),
			),
		},
		{
			name: "Valid record store with empty fields",
			input: &RecordStore{
				RecordStorePath: "",
				TemplateHash:    "",
				RequestHash:     "",
				ResponseID:      "",
			},
			expected: slog.GroupValue(
				slog.String("recordStorePath", ""),
				slog.String("templateHash", ""),
				slog.String("requestHash", ""),
				slog.String("responseId", ""),
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.LogValue()
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("LogValue() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestFileInfo_LogValue(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name     string
		input    FileInfo
		expected slog.Value
	}{
		{
			name: "Valid file info with all fields",
			input: FileInfo{
				FilePath:     "/test/file.yaml",
				TemplateHash: "template123",
				RequestHash:  "request456",
				ResponseID:   "response789",
				ModTime:      now,
			},
			expected: slog.GroupValue(
				slog.String("filePath", "/test/file.yaml"),
				slog.String("templateHash", "template123"),
				slog.String("requestHash", "request456"),
				slog.String("responseId", "response789"),
				slog.Time("modTime", now),
			),
		},
		{
			name: "Valid file info with empty fields",
			input: FileInfo{
				FilePath:     "",
				TemplateHash: "",
				RequestHash:  "",
				ResponseID:   "",
				ModTime:      time.Time{},
			},
			expected: slog.GroupValue(
				slog.String("filePath", ""),
				slog.String("templateHash", ""),
				slog.String("requestHash", ""),
				slog.String("responseId", ""),
				slog.Time("modTime", time.Time{}),
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.LogValue()
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("LogValue() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestFailedRecord_TemplateWriteFailure(t *testing.T) {
	root := t.TempDir()
	templatesDir := filepath.Join(root, "templates")
	recordStore := &RecordStore{
		RecordStorePath: root,
		TemplateYaml:    []byte("template: value"),
		Request: &request.RequestObject{
			TemplateHash: "templateHash",
		},
		Response: &response.ResponseObject{
			TemplateHash: "templateHash",
			RequestHash:  "requestHash",
		},
	}
	recordStore.RequestYaml, _ = utils.ConvertToYAML(recordStore.Request)
	recordStore.ResponseYaml, _ = utils.ConvertToYAML(recordStore.Response)
	os.Mkdir(templatesDir, 0755)
	os.Chmod(templatesDir, 0o555)
	t.Cleanup(func() {
		_ = os.Chmod(templatesDir, 0o755)
	})
	err := recordStore.Record()
	if err == nil {
		t.Fatalf("Expected error, received nil")
	}
	expectedErr := ErrorFailedToRecord
	if !errors.Is(err, expectedErr) {
		t.Fatalf("Expected error %v, received %v", expectedErr, err)
	}
}

// func TestFailedRecord_RequestWriteFailure(t *testing.T) {
// 	root := t.TempDir()
// 	requestsDir := filepath.Join(root, "requests", "templateHash")
// 	recordStore := &RecordStore{
// 		RecordStorePath: root,
// 		TemplateYaml:    []byte("template: value"),
// 		Request: &request.RequestObject{
// 			TemplateHash: "templateHash",
// 			URL:          "https://example.com",
// 			Method:       "GET",
// 		},
// 		Response: &response.ResponseObject{
// 			TemplateHash: "templateHash",
// 			RequestHash:  "requestHash",
// 			StatusCode:   200,
// 		},
// 	}
// 	// Manually set the template hash
// 	recordStore.TemplateHash = utils.CalculateMD5Hash(recordStore.TemplateYaml)
// 	// Manually convert request to YAML
// 	var err error
// 	recordStore.RequestYaml, err = utils.ConvertToYAML(recordStore.Request)
// 	if err != nil {
// 		t.Fatalf("Failed to convert request to YAML: %v", err)
// 	}
// 	// Manually set request hash
// 	recordStore.RequestHash = utils.CalculateMD5Hash(recordStore.RequestYaml)
// 	// Manually convert response to YAML
// 	recordStore.ResponseYaml, err = utils.ConvertToYAML(recordStore.Response)
// 	if err != nil {
// 		t.Fatalf("Failed to convert response to YAML: %v", err)
// 	}

// 	// Create and make the requests directory unwritable
// 	err = os.MkdirAll(requestsDir, 0755)
// 	if err != nil {
// 		t.Fatalf("Failed to create requests directory: %v", err)
// 	}
// 	err = os.Chmod(requestsDir, 0o555)
// 	if err != nil {
// 		t.Fatalf("Failed to change directory permissions: %v", err)
// 	}
// 	t.Cleanup(func() {
// 		_ = os.Chmod(requestsDir, 0o755)
// 	})

// 	// Now try to record - this should fail during request write
// 	err = recordStore.Record()
// 	if err == nil {
// 		t.Fatalf("Expected error, received nil")
// 	}
// 	expectedErr := ErrorFailedToRecord
// 	if !errors.Is(err, expectedErr) {
// 		t.Fatalf("Expected error %v, received %v", expectedErr, err)
// 	}
// }

// func TestFailedRecord_ResponseWriteFailure(t *testing.T) {
// 	root := t.TempDir()
// 	responsesDir := filepath.Join(root, "responses", "requestHash")
// 	recordStore := &RecordStore{
// 		RecordStorePath: root,
// 		TemplateYaml:    []byte("template: value"),
// 		Request: &request.RequestObject{
// 			TemplateHash: "templateHash",
// 			URL:          "https://example.com",
// 			Method:       "GET",
// 		},
// 		Response: &response.ResponseObject{
// 			TemplateHash: "templateHash",
// 			RequestHash:  "requestHash",
// 			StatusCode:   200,
// 		},
// 	}
// 	// Manually set the template hash
// 	recordStore.TemplateHash = utils.CalculateMD5Hash(recordStore.TemplateYaml)
// 	// Manually convert request to YAML
// 	var err error
// 	recordStore.RequestYaml, err = utils.ConvertToYAML(recordStore.Request)
// 	if err != nil {
// 		t.Fatalf("Failed to convert request to YAML: %v", err)
// 	}
// 	// Manually set request hash
// 	recordStore.RequestHash = utils.CalculateMD5Hash(recordStore.RequestYaml)
// 	// Manually convert response to YAML
// 	recordStore.ResponseYaml, err = utils.ConvertToYAML(recordStore.Response)
// 	if err != nil {
// 		t.Fatalf("Failed to convert response to YAML: %v", err)
// 	}

// 	// Create and make the responses directory unwritable
// 	err = os.MkdirAll(responsesDir, 0755)
// 	if err != nil {
// 		t.Fatalf("Failed to create responses directory: %v", err)
// 	}
// 	err = os.Chmod(responsesDir, 0o555)
// 	if err != nil {
// 		t.Fatalf("Failed to change directory permissions: %v", err)
// 	}
// 	t.Cleanup(func() {
// 		_ = os.Chmod(responsesDir, 0o755)
// 	})

// 	// Now try to record - this should fail during response write
// 	err = recordStore.Record()
// 	if err == nil {
// 		t.Fatalf("Expected error, received nil")
// 	}
// 	expectedErr := ErrorFailedToRecord
// 	if !errors.Is(err, expectedErr) {
// 		t.Fatalf("Expected error %v, received %v", expectedErr, err)
// 	}
// }
