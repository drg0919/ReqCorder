package history

import (
	"errors"
	"reflect"
	"reqcorder/internal/record"
	"reqcorder/internal/request"
	"reqcorder/internal/response"
	"testing"
	"time"

	"log/slog"
)

var (
	illegalHash = "illegal_hash"
)

func TestSuccessfulGetTemplateByHash(t *testing.T) {
	root := t.TempDir()
	recordStore := &record.RecordStore{
		RecordStorePath: root,
		RequestYaml:     []byte("key: value"),
		ResponseYaml:    []byte("key: value"),
		TemplateYaml:    []byte("key: value"),
		Request:         &request.RequestObject{},
		Response:        &response.ResponseObject{},
	}
	err := recordStore.Record()
	if err != nil {
		t.Fatalf("Expected no error, received %v\n", err)
	}
	historyStore := &HistoryStore{
		RecordStorePath: root,
	}
	_, err = historyStore.GetTemplateByHash(recordStore.TemplateHash)
	if err != nil {
		t.Fatalf("Expected no error, received %v\n", err)
	}
	// 	expectedOutput := `
	// Template

	// key: value`
	//
	//	if res != expectedOutput {
	//		t.Fatal("Received incorrect output string")
	//	}
}

func TestFailedGetTemplateByHash_StoreFailure(t *testing.T) {
	root := t.TempDir()
	historyStore := &HistoryStore{
		RecordStorePath: root,
	}
	_, err := historyStore.GetTemplateByHash(illegalHash)
	expectedErr := record.ErrorFailedToStatPath
	if err == nil {
		t.Fatal("Expected error, received nil")
	}
	if !errors.Is(err, expectedErr) {
		t.Fatalf("Expected %v, received %v\n", expectedErr, err)
	}
}

func TestSuccessfulGetRequestByHash(t *testing.T) {
	root := t.TempDir()
	recordStore := &record.RecordStore{
		RecordStorePath: root,
		RequestYaml:     []byte("key: value"),
		ResponseYaml:    []byte("key: value"),
		TemplateYaml:    []byte("key: value"),
		Request:         &request.RequestObject{},
		Response:        &response.ResponseObject{},
	}
	err := recordStore.Record()
	if err != nil {
		t.Fatalf("Expected no error, received %v\n", err)
	}
	historyStore := &HistoryStore{
		RecordStorePath: root,
	}
	res, err := historyStore.GetRequestByHash(recordStore.RequestHash)
	if err != nil {
		t.Fatalf("Expected no error, received %v\n", err)
	}
	expectedOutput := "\nTemplate Hash: " + recordStore.TemplateHash + "\nRequest:\n\ntemplate_hash: " + recordStore.TemplateHash + "\nurl: \"\"\nmethod: \"\"\n"
	if res != expectedOutput {
		t.Fatal("Received incorrect output string\n")
	}
}

func TestFailedGetRequestByHash_StoreFailure(t *testing.T) {
	root := t.TempDir()
	historyStore := &HistoryStore{
		RecordStorePath: root,
	}
	_, err := historyStore.GetRequestByHash(illegalHash)
	expectedErr := record.ErrorFailedToStatPath
	if err == nil {
		t.Fatal("Expected error, received nil")
	}
	if !errors.Is(err, expectedErr) {
		t.Fatalf("Expected %v, received %v\n", expectedErr, err)
	}
}

func TestSuccessfulGetResponseByID(t *testing.T) {
	root := t.TempDir()
	recordStore := &record.RecordStore{
		RecordStorePath: root,
		RequestYaml:     []byte("key: value"),
		ResponseYaml:    []byte("key: value"),
		TemplateYaml:    []byte("key: value"),
		Request:         &request.RequestObject{},
		Response:        &response.ResponseObject{},
	}
	err := recordStore.Record()
	if err != nil {
		t.Fatalf("Expected no error, received %v\n", err)
	}
	historyStore := &HistoryStore{
		RecordStorePath: root,
	}
	res, err := historyStore.GetResponseByID(recordStore.ResponseID)
	if err != nil {
		t.Fatalf("Expected no error, received %v\n", err)
	}
	expectedOutput := "\nTemplate Hash: " + recordStore.TemplateHash + "\nRequest Hash: " + recordStore.RequestHash + "\nResponse:\n\nrequest_hash: " + recordStore.RequestHash + "\ntemplate_hash: " + recordStore.TemplateHash + "\nstatus_code: 0\nheaders: {}\nbody: \"\"\nsize_bytes: 0\ntiming:\n  dns_start: 0001-01-01T00:00:00Z\n  dns_end: 0001-01-01T00:00:00Z\n  connect_start: 0001-01-01T00:00:00Z\n  connect_done: 0001-01-01T00:00:00Z\n  tls_handshake_start: 0001-01-01T00:00:00Z\n  tls_handshake_done: 0001-01-01T00:00:00Z\n  got_first_response_byte: 0001-01-01T00:00:00Z\n  dns_lookup: 0s\n  tcp_connect: 0s\n  tls_handshake: 0s\n  time_to_first_byte: 0s\n  total_duration: 0s\ncookies: []\n"
	if res != expectedOutput {
		t.Fatal("Received incorrect output string\n")
	}
}

func TestFailedGetResponseByID_StoreFailure(t *testing.T) {
	root := t.TempDir()
	historyStore := &HistoryStore{
		RecordStorePath: root,
	}
	_, err := historyStore.GetResponseByID(illegalHash)
	expectedErr := record.ErrorFailedToStatPath
	if err == nil {
		t.Fatal("Expected error, received nil")
	}
	if !errors.Is(err, expectedErr) {
		t.Fatalf("Expected %v, received %v\n", expectedErr, err)
	}
}

func TestSuccessfulGetAllRequestsSorted(t *testing.T) {
	root := t.TempDir()
	recordStore := &record.RecordStore{
		RecordStorePath: root,
		RequestYaml:     []byte("key: value"),
		ResponseYaml:    []byte("key: value"),
		TemplateYaml:    []byte("key: value"),
		Request:         &request.RequestObject{},
		Response:        &response.ResponseObject{},
	}
	err := recordStore.Record()
	if err != nil {
		t.Fatalf("Expected no error, received %v\n", err)
	}

	recordStore2 := &record.RecordStore{
		RecordStorePath: root,
		RequestYaml:     []byte("key2: value2"),
		ResponseYaml:    []byte("key2: value2"),
		TemplateYaml:    []byte("key2: value2"),
		Request:         &request.RequestObject{},
		Response:        &response.ResponseObject{},
	}
	err = recordStore2.Record()
	if err != nil {
		t.Fatalf("Expected no error, received %v\n", err)
	}

	historyStore := &HistoryStore{
		RecordStorePath: root,
	}
	requests, err := historyStore.GetAllRequestsSorted(0)
	if err != nil {
		t.Fatalf("Expected no error, received %v\n", err)
	}

	if len(requests) != 2 {
		t.Fatalf("Expected 2 requests, received %d\n", len(requests))
	}

	for _, request := range requests {
		if len(request) != 3 {
			t.Fatalf("Expected 3 fields in request data, got %d\n", len(request))
		}
		if request[0] == "" {
			t.Fatal("Request hash cannot be empty")
		}
		if request[1] == "" {
			t.Fatal("Template hash cannot be empty")
		}
		if request[2] == "" {
			t.Fatal("Timestamp cannot be empty")
		}
	}
}

func TestSuccessfulGetAllRequestsSorted_Limit(t *testing.T) {
	root := t.TempDir()
	recordStore := &record.RecordStore{
		RecordStorePath: root,
		RequestYaml:     []byte("key: value"),
		ResponseYaml:    []byte("key: value"),
		TemplateYaml:    []byte("key: value"),
		Request:         &request.RequestObject{},
		Response:        &response.ResponseObject{},
	}
	err := recordStore.Record()
	if err != nil {
		t.Fatalf("Expected no error, received %v\n", err)
	}

	recordStore2 := &record.RecordStore{
		RecordStorePath: root,
		RequestYaml:     []byte("key2: value2"),
		ResponseYaml:    []byte("key2: value2"),
		TemplateYaml:    []byte("key2: value2"),
		Request:         &request.RequestObject{},
		Response:        &response.ResponseObject{},
	}
	err = recordStore2.Record()
	if err != nil {
		t.Fatalf("Expected no error, received %v\n", err)
	}

	historyStore := &HistoryStore{
		RecordStorePath: root,
	}
	requests, err := historyStore.GetAllRequestsSorted(1)
	if err != nil {
		t.Fatalf("Expected no error, received %v\n", err)
	}

	if len(requests) != 1 {
		t.Fatalf("Expected 1 request due to limit, got %d\n", len(requests))
	}
}

func TestFailedGetAllRequestsSorted_StoreFailure(t *testing.T) {
	root := t.TempDir()
	historyStore := &HistoryStore{
		RecordStorePath: root,
	}
	requests, err := historyStore.GetAllRequestsSorted(0)
	if err != nil {
		t.Fatalf("Expected no error, received %v\n", err)
	}
	if len(requests) != 0 {
		t.Fatalf("Expected empty slice, got %d requests\n", len(requests))
	}
}

func TestSuccessfulGetAllResponsesSorted(t *testing.T) {
	root := t.TempDir()
	recordStore := &record.RecordStore{
		RecordStorePath: root,
		RequestYaml:     []byte("key: value"),
		ResponseYaml:    []byte("key: value"),
		TemplateYaml:    []byte("key: value"),
		Request:         &request.RequestObject{},
		Response:        &response.ResponseObject{},
	}
	err := recordStore.Record()
	if err != nil {
		t.Fatalf("Expected no error, received %v\n", err)
	}

	recordStore2 := &record.RecordStore{
		RecordStorePath: root,
		RequestYaml:     []byte("key2: value2"),
		ResponseYaml:    []byte("key2: value2"),
		TemplateYaml:    []byte("key2: value2"),
		Request:         &request.RequestObject{},
		Response:        &response.ResponseObject{},
	}
	err = recordStore2.Record()
	if err != nil {
		t.Fatalf("Expected no error, received %v\n", err)
	}

	historyStore := &HistoryStore{
		RecordStorePath: root,
	}
	responses, err := historyStore.GetAllResponsesSorted(0)
	if err != nil {
		t.Fatalf("Expected no error, received %v\n", err)
	}

	if len(responses) != 2 {
		t.Fatalf("Expected 2 responses, received %d\n", len(responses))
	}

	for _, response := range responses {
		if len(response) != 4 {
			t.Fatalf("Expected 4 fields in response data, got %d\n", len(response))
		}
		if response[0] == "" {
			t.Fatal("Response ID cannot be empty")
		}
		if response[1] == "" {
			t.Fatal("Status code with indicator cannot be empty")
		}
		if response[2] == "" {
			t.Fatal("Timing total cannot be empty")
		}
		if response[3] == "" {
			t.Fatal("Timestamp cannot be empty")
		}
	}

	firstResponseID := responses[0][0]
	secondResponseID := responses[1][0]

	firstTimestamp := firstResponseID[:19]
	secondTimestamp := secondResponseID[:19]

	if firstTimestamp < secondTimestamp {
		t.Fatalf("Expected responses to be in descending order by timestamp, but received %s before %s", firstTimestamp, secondTimestamp)
	}
}

func TestSuccessfulGetAllResponsesSorted_StatusCode400(t *testing.T) {
	root := t.TempDir()
	recordStore := &record.RecordStore{
		RecordStorePath: root,
		RequestYaml:     []byte("key: value"),
		ResponseYaml:    []byte("key: value"),
		TemplateYaml:    []byte("key: value"),
		Request:         &request.RequestObject{},
		Response: &response.ResponseObject{
			StatusCode: 400,
		},
	}
	err := recordStore.Record()
	if err != nil {
		t.Fatalf("Expected no error, received %v\n", err)
	}

	historyStore := &HistoryStore{
		RecordStorePath: root,
	}
	responses, err := historyStore.GetAllResponsesSorted(0)
	if err != nil {
		t.Fatalf("Expected no error, received %v\n", err)
	}

	if len(responses) != 1 {
		t.Fatalf("Expected 1 response, received %d\n", len(responses))
	}

	statusWithIndicator := responses[0][1]
	if statusWithIndicator != "400 ❌" {
		t.Fatalf("Expected '400 ❌', got '%s'\n", statusWithIndicator)
	}
}

func TestSuccessfulGetAllResponsesSorted_StatusCode200(t *testing.T) {
	root := t.TempDir()
	recordStore := &record.RecordStore{
		RecordStorePath: root,
		RequestYaml:     []byte("key: value"),
		ResponseYaml:    []byte("key: value"),
		TemplateYaml:    []byte("key: value"),
		Request:         &request.RequestObject{},
		Response: &response.ResponseObject{
			StatusCode: 200,
		},
	}
	err := recordStore.Record()
	if err != nil {
		t.Fatalf("Expected no error, received %v\n", err)
	}

	historyStore := &HistoryStore{
		RecordStorePath: root,
	}
	responses, err := historyStore.GetAllResponsesSorted(0)
	if err != nil {
		t.Fatalf("Expected no error, received %v\n", err)
	}

	if len(responses) != 1 {
		t.Fatalf("Expected 1 response, received %d\n", len(responses))
	}

	statusWithIndicator := responses[0][1]
	if statusWithIndicator != "200 ✅" {
		t.Fatalf("Expected '200 ✅', got '%s'\n", statusWithIndicator)
	}
}

func TestSuccessfulGetAllResponsesSorted_Limit(t *testing.T) {
	root := t.TempDir()
	recordStore := &record.RecordStore{
		RecordStorePath: root,
		RequestYaml:     []byte("key: value"),
		ResponseYaml:    []byte("key: value"),
		TemplateYaml:    []byte("key: value"),
		Request:         &request.RequestObject{},
		Response:        &response.ResponseObject{},
	}
	err := recordStore.Record()
	if err != nil {
		t.Fatalf("Expected no error, received %v\n", err)
	}

	recordStore2 := &record.RecordStore{
		RecordStorePath: root,
		RequestYaml:     []byte("key2: value2"),
		ResponseYaml:    []byte("key2: value2"),
		TemplateYaml:    []byte("key2: value2"),
		Request:         &request.RequestObject{},
		Response:        &response.ResponseObject{},
	}
	err = recordStore2.Record()
	if err != nil {
		t.Fatalf("Expected no error, received %v\n", err)
	}

	historyStore := &HistoryStore{
		RecordStorePath: root,
	}
	responses, err := historyStore.GetAllResponsesSorted(1)
	if err != nil {
		t.Fatalf("Expected no error, received %v\n", err)
	}

	if len(responses) != 1 {
		t.Fatalf("Expected 1 response due to limit, got %d\n", len(responses))
	}
}

func TestFailedGetAllResponsesSorted_StoreFailure(t *testing.T) {
	root := t.TempDir()
	historyStore := &HistoryStore{
		RecordStorePath: root,
	}
	_, err := historyStore.GetAllResponsesSorted(0)
	if err == nil {
		t.Fatal("Expected error, received nil")
	}
	expectedErr := record.ErrorFailedToReadDirectory
	if !errors.Is(err, expectedErr) {
		t.Fatalf("Expected %v, received %v\n", expectedErr, err)
	}
}

func TestSuccessfulGetSortedResponsesByRequestHash(t *testing.T) {
	root := t.TempDir()
	recordStore := &record.RecordStore{
		RecordStorePath: root,
		RequestYaml:     []byte("key: value"),
		ResponseYaml:    []byte("key: value"),
		TemplateYaml:    []byte("key: value"),
		Request:         &request.RequestObject{},
		Response:        &response.ResponseObject{},
	}
	err := recordStore.Record()
	if err != nil {
		t.Fatalf("Expected no error, received %v\n", err)
	}

	recordStore2 := &record.RecordStore{
		RecordStorePath: root,
		RequestYaml:     []byte("key: value"),
		ResponseYaml:    []byte("key2: value2"),
		TemplateYaml:    []byte("key: value"),
		Request:         &request.RequestObject{},
		Response:        &response.ResponseObject{},
	}
	err = recordStore2.Record()
	if err != nil {
		t.Fatalf("Expected no error, received %v\n", err)
	}

	historyStore := &HistoryStore{
		RecordStorePath: root,
	}
	responses, err := historyStore.GetSortedResponsesByRequestHash(recordStore.RequestHash, 0)
	if err != nil {
		t.Fatalf("Expected no error, received %v\n", err)
	}

	if len(responses) != 2 {
		t.Fatalf("Expected 2 responses, got %d\n", len(responses))
	}

	for _, response := range responses {
		if len(response) != 4 {
			t.Fatalf("Expected 4 fields in response data, got %d\n", len(response))
		}
		if response[0] == "" {
			t.Fatal("Response ID cannot be empty")
		}
		if response[1] == "" {
			t.Fatal("Status code with indicator cannot be empty")
		}
		if response[2] == "" {
			t.Fatal("Timing total cannot be empty")
		}
		if response[3] == "" {
			t.Fatal("Timestamp cannot be empty")
		}
	}

	firstResponseID := responses[0][0]
	secondResponseID := responses[1][0]

	firstTimestamp := firstResponseID[:19]
	secondTimestamp := secondResponseID[:19]

	if firstTimestamp < secondTimestamp {
		t.Fatalf("Expected responses to be in descending order by timestamp, but received %s before %s", firstTimestamp, secondTimestamp)
	}
}

func TestSuccessfulGetSortedResponsesByRequestHash_Limit(t *testing.T) {
	root := t.TempDir()
	recordStore := &record.RecordStore{
		RecordStorePath: root,
		RequestYaml:     []byte("key: value"),
		ResponseYaml:    []byte("key: value"),
		TemplateYaml:    []byte("key: value"),
		Request:         &request.RequestObject{},
		Response:        &response.ResponseObject{},
	}
	err := recordStore.Record()
	if err != nil {
		t.Fatalf("Expected no error, received %v\n", err)
	}

	recordStore2 := &record.RecordStore{
		RecordStorePath: root,
		RequestYaml:     []byte("key2: value2"),
		ResponseYaml:    []byte("key2: value2"),
		TemplateYaml:    []byte("key2: value2"),
		Request:         &request.RequestObject{},
		Response:        &response.ResponseObject{},
	}
	err = recordStore2.Record()
	if err != nil {
		t.Fatalf("Expected no error, received %v\n", err)
	}

	historyStore := &HistoryStore{
		RecordStorePath: root,
	}
	responses, err := historyStore.GetSortedResponsesByRequestHash(recordStore.RequestHash, 1)
	if err != nil {
		t.Fatalf("Expected no error, received %v\n", err)
	}

	if len(responses) != 1 {
		t.Fatalf("Expected 1 response due to limit, got %d\n", len(responses))
	}
}

func TestFailedGetSortedResponsesByRequestHash_StoreFailure(t *testing.T) {
	root := t.TempDir()
	recordStore := &record.RecordStore{
		RecordStorePath: root,
		RequestYaml:     []byte("key: value"),
		ResponseYaml:    []byte("key: value"),
		TemplateYaml:    []byte("key: value"),
		Request:         &request.RequestObject{},
		Response:        &response.ResponseObject{},
	}
	err := recordStore.Record()
	if err != nil {
		t.Fatalf("Expected no error, received %v\n", err)
	}
	historyStore := &HistoryStore{
		RecordStorePath: root,
	}
	_, err = historyStore.GetSortedResponsesByRequestHash(illegalHash, 0)
	if err == nil {
		t.Fatal("Expected error, received nil")
	}
	expectedErr := record.ErrorFailedToStatPath
	if !errors.Is(err, expectedErr) {
		t.Fatalf("Expected %v, received %v\n", expectedErr, err)
	}
}

func TestSuccessfulGetSortedResponsesByRequestHash_StatusCode400(t *testing.T) {
	root := t.TempDir()
	recordStore := &record.RecordStore{
		RecordStorePath: root,
		RequestYaml:     []byte("key: value"),
		ResponseYaml:    []byte("key: value"),
		TemplateYaml:    []byte("key: value"),
		Request:         &request.RequestObject{},
		Response: &response.ResponseObject{
			StatusCode: 400,
		},
	}
	err := recordStore.Record()
	if err != nil {
		t.Fatalf("Expected no error, received %v\n", err)
	}

	historyStore := &HistoryStore{
		RecordStorePath: root,
	}
	responses, err := historyStore.GetSortedResponsesByRequestHash(recordStore.RequestHash, 0)
	if err != nil {
		t.Fatalf("Expected no error, received %v\n", err)
	}

	if len(responses) != 1 {
		t.Fatalf("Expected 1 response, received %d\n", len(responses))
	}

	statusWithIndicator := responses[0][1]
	if statusWithIndicator != "400 ❌" {
		t.Fatalf("Expected '400 ❌', got '%s'\n", statusWithIndicator)
	}
}

func TestSuccessfulGetSortedResponsesByRequestHash_StatusCode200(t *testing.T) {
	root := t.TempDir()
	recordStore := &record.RecordStore{
		RecordStorePath: root,
		RequestYaml:     []byte("key: value"),
		ResponseYaml:    []byte("key: value"),
		TemplateYaml:    []byte("key: value"),
		Request:         &request.RequestObject{},
		Response: &response.ResponseObject{
			StatusCode: 200,
		},
	}
	err := recordStore.Record()
	if err != nil {
		t.Fatalf("Expected no error, received %v\n", err)
	}

	historyStore := &HistoryStore{
		RecordStorePath: root,
	}
	responses, err := historyStore.GetSortedResponsesByRequestHash(recordStore.RequestHash, 0)
	if err != nil {
		t.Fatalf("Expected no error, received %v\n", err)
	}

	if len(responses) != 1 {
		t.Fatalf("Expected 1 response, received %d\n", len(responses))
	}

	statusWithIndicator := responses[0][1]
	if statusWithIndicator != "200 ✅" {
		t.Fatalf("Expected '200 ✅', got '%s'\n", statusWithIndicator)
	}
}

func TestSuccessfulGetSortedResponsesByTemplateHash(t *testing.T) {
	root := t.TempDir()
	recordStore := &record.RecordStore{
		RecordStorePath: root,
		RequestYaml:     []byte("key: value"),
		ResponseYaml:    []byte("key: value"),
		TemplateYaml:    []byte("key: value"),
		Request:         &request.RequestObject{},
		Response:        &response.ResponseObject{},
	}
	err := recordStore.Record()
	if err != nil {
		t.Fatalf("Expected no error, received %v\n", err)
	}

	recordStore2 := &record.RecordStore{
		RecordStorePath: root,
		RequestYaml:     []byte("key2: value2"),
		ResponseYaml:    []byte("key2: value2"),
		TemplateYaml:    []byte("key: value"),
		Request:         &request.RequestObject{},
		Response:        &response.ResponseObject{},
	}
	err = recordStore2.Record()
	if err != nil {
		t.Fatalf("Expected no error, received %v\n", err)
	}

	historyStore := &HistoryStore{
		RecordStorePath: root,
	}
	responses, err := historyStore.GetSortedResponsesByTemplateHash(recordStore.TemplateHash, 0)
	if err != nil {
		t.Fatalf("Expected no error, received %v\n", err)
	}

	if len(responses) != 2 {
		t.Fatalf("Expected 2 responses, got %d\n", len(responses))
	}

	for _, response := range responses {
		if len(response) != 4 {
			t.Fatalf("Expected 4 fields in response data, got %d\n", len(response))
		}
		if response[0] == "" {
			t.Fatal("Response ID cannot be empty")
		}
		if response[1] == "" {
			t.Fatal("Status code with indicator cannot be empty")
		}
		if response[2] == "" {
			t.Fatal("Timing total cannot be empty")
		}
		if response[3] == "" {
			t.Fatal("Timestamp cannot be empty")
		}
	}

	firstResponseID := responses[0][0]
	secondResponseID := responses[1][0]

	firstTimestamp := firstResponseID[:19]
	secondTimestamp := secondResponseID[:19]

	if firstTimestamp < secondTimestamp {
		t.Fatalf("Expected responses to be in descending order by timestamp, but received %s before %s", firstTimestamp, secondTimestamp)
	}
}

func TestSuccessfulGetSortedResponsesByTemplateHash_Limit(t *testing.T) {
	root := t.TempDir()
	recordStore := &record.RecordStore{
		RecordStorePath: root,
		RequestYaml:     []byte("key: value"),
		ResponseYaml:    []byte("key: value"),
		TemplateYaml:    []byte("key: value"),
		Request:         &request.RequestObject{},
		Response:        &response.ResponseObject{},
	}
	err := recordStore.Record()
	if err != nil {
		t.Fatalf("Expected no error, received %v\n", err)
	}

	recordStore2 := &record.RecordStore{
		RecordStorePath: root,
		RequestYaml:     []byte("key2: value2"),
		ResponseYaml:    []byte("key2: value2"),
		TemplateYaml:    []byte("key: value"),
		Request:         &request.RequestObject{},
		Response:        &response.ResponseObject{},
	}
	err = recordStore2.Record()
	if err != nil {
		t.Fatalf("Expected no error, received %v\n", err)
	}

	historyStore := &HistoryStore{
		RecordStorePath: root,
	}
	responses, err := historyStore.GetSortedResponsesByTemplateHash(recordStore.TemplateHash, 1)
	if err != nil {
		t.Fatalf("Expected no error, received %v\n", err)
	}

	if len(responses) != 1 {
		t.Fatalf("Expected 1 response due to limit, got %d\n", len(responses))
	}
}

func TestFailedGetSortedResponsesByTemplateHash_StoreFailure(t *testing.T) {
	root := t.TempDir()
	recordStore := &record.RecordStore{
		RecordStorePath: root,
		RequestYaml:     []byte("key: value"),
		ResponseYaml:    []byte("key: value"),
		TemplateYaml:    []byte("key: value"),
		Request:         &request.RequestObject{},
		Response:        &response.ResponseObject{},
	}
	err := recordStore.Record()
	if err != nil {
		t.Fatalf("Expected no error, received %v\n", err)
	}
	historyStore := &HistoryStore{
		RecordStorePath: root,
	}
	_, err = historyStore.GetSortedResponsesByTemplateHash(illegalHash, 0)
	if err == nil {
		t.Fatal("Expected error, received nil")
	}
	expectedErr := record.ErrorFailedToStatPath
	if !errors.Is(err, expectedErr) {
		t.Fatalf("Expected %v, received %v\n", expectedErr, err)
	}
}

func TestSuccessfulGetSortedResponsesByTemplateHash_StatusCode400(t *testing.T) {
	root := t.TempDir()
	recordStore := &record.RecordStore{
		RecordStorePath: root,
		RequestYaml:     []byte("key: value"),
		ResponseYaml:    []byte("key: value"),
		TemplateYaml:    []byte("key: value"),
		Request:         &request.RequestObject{},
		Response: &response.ResponseObject{
			StatusCode: 400,
		},
	}
	err := recordStore.Record()
	if err != nil {
		t.Fatalf("Expected no error, received %v\n", err)
	}

	historyStore := &HistoryStore{
		RecordStorePath: root,
	}
	responses, err := historyStore.GetSortedResponsesByTemplateHash(recordStore.TemplateHash, 0)
	if err != nil {
		t.Fatalf("Expected no error, received %v\n", err)
	}

	if len(responses) != 1 {
		t.Fatalf("Expected 1 response, received %d\n", len(responses))
	}

	statusWithIndicator := responses[0][1]
	if statusWithIndicator != "400 ❌" {
		t.Fatalf("Expected '400 ❌', got '%s'\n", statusWithIndicator)
	}
}

func TestSuccessfulGetSortedResponsesByTemplateHash_StatusCode200(t *testing.T) {
	root := t.TempDir()
	recordStore := &record.RecordStore{
		RecordStorePath: root,
		RequestYaml:     []byte("key: value"),
		ResponseYaml:    []byte("key: value"),
		TemplateYaml:    []byte("key: value"),
		Request:         &request.RequestObject{},
		Response: &response.ResponseObject{
			StatusCode: 200,
		},
	}
	err := recordStore.Record()
	if err != nil {
		t.Fatalf("Expected no error, received %v\n", err)
	}

	historyStore := &HistoryStore{
		RecordStorePath: root,
	}
	responses, err := historyStore.GetSortedResponsesByTemplateHash(recordStore.TemplateHash, 0)
	if err != nil {
		t.Fatalf("Expected no error, received %v\n", err)
	}

	if len(responses) != 1 {
		t.Fatalf("Expected 1 response, received %d\n", len(responses))
	}

	statusWithIndicator := responses[0][1]
	if statusWithIndicator != "200 ✅" {
		t.Fatalf("Expected '200 ✅', got '%s'\n", statusWithIndicator)
	}
}

func TestLogValue_NilHistoryStore(t *testing.T) {
	var h *HistoryStore
	result := h.LogValue()
	expected := slog.StringValue("<nil>")
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("LogValue() = %v, want %v", result, expected)
	}
}

func TestLogValue_ValidHistoryStore(t *testing.T) {
	h := &HistoryStore{
		RecordStorePath: "/path/to/record/store",
	}
	result := h.LogValue()
	expected := slog.GroupValue(
		slog.String("recordStorePath", "/path/to/record/store"),
	)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("LogValue() = %v, want %v", result, expected)
	}
}

func TestSuccessfulGetAllTemplatesSorted(t *testing.T) {
	root := t.TempDir()
	recordStore := &record.RecordStore{
		RecordStorePath: root,
		RequestYaml:     []byte("key: value"),
		ResponseYaml:    []byte("key: value"),
		TemplateYaml:    []byte("key: value"),
		Request:         &request.RequestObject{},
		Response:        &response.ResponseObject{},
	}
	err := recordStore.Record()
	if err != nil {
		t.Fatalf("Expected no error, received %v\n", err)
	}

	recordStore2 := &record.RecordStore{
		RecordStorePath: root,
		RequestYaml:     []byte("key2: value2"),
		ResponseYaml:    []byte("key2: value2"),
		TemplateYaml:    []byte("key2: value2"),
		Request:         &request.RequestObject{},
		Response:        &response.ResponseObject{},
	}
	err = recordStore2.Record()
	if err != nil {
		t.Fatalf("Expected no error, received %v\n", err)
	}

	historyStore := &HistoryStore{
		RecordStorePath: root,
	}
	templates, err := historyStore.GetAllTemplatesSorted(0)
	if err != nil {
		t.Fatalf("Expected no error, received %v\n", err)
	}

	if len(templates) != 2 {
		t.Fatalf("Expected 2 templates, received %d\n", len(templates))
	}

	for _, template := range templates {
		if len(template) != 2 {
			t.Fatalf("Expected 2 fields in template data, received %d\n", len(template))
		}
		if template[0] == "" {
			t.Fatal("Template hash cannot be empty")
		}
		if template[1] == "" {
			t.Fatal("Timestamp cannot be empty")
		}
	}

	firstTemplateHash := templates[0][0]
	secondTemplateHash := templates[1][0]
	if firstTemplateHash == "" || secondTemplateHash == "" {
		t.Fatal("Received empty template hashe(s)")
	}
}

func TestSuccessfulGetAllTemplatesSorted_Limit(t *testing.T) {
	root := t.TempDir()
	recordStore := &record.RecordStore{
		RecordStorePath: root,
		RequestYaml:     []byte("key: value"),
		ResponseYaml:    []byte("key: value"),
		TemplateYaml:    []byte("key: value"),
		Request:         &request.RequestObject{},
		Response:        &response.ResponseObject{},
	}
	err := recordStore.Record()
	if err != nil {
		t.Fatalf("Expected no error, received %v\n", err)
	}

	recordStore2 := &record.RecordStore{
		RecordStorePath: root,
		RequestYaml:     []byte("key2: value2"),
		ResponseYaml:    []byte("key2: value2"),
		TemplateYaml:    []byte("key2: value2"),
		Request:         &request.RequestObject{},
		Response:        &response.ResponseObject{},
	}
	err = recordStore2.Record()
	if err != nil {
		t.Fatalf("Expected no error, received %v\n", err)
	}

	historyStore := &HistoryStore{
		RecordStorePath: root,
	}
	templates, err := historyStore.GetAllTemplatesSorted(1)
	if err != nil {
		t.Fatalf("Expected no error, received %v\n", err)
	}

	if len(templates) != 1 {
		t.Fatalf("Expected 1 template due to limit, received %d\n", len(templates))
	}
}

func TestFailedGetAllTemplatesSorted_StoreFailure(t *testing.T) {
	root := t.TempDir()
	historyStore := &HistoryStore{
		RecordStorePath: root,
	}
	templates, err := historyStore.GetAllTemplatesSorted(0)
	if err != nil {
		t.Fatalf("Expected no error, received %v\n", err)
	}
	if len(templates) != 0 {
		t.Fatalf("Expected empty slice, received %d templates\n", len(templates))
	}
}

func TestLogValue_FileInfo(t *testing.T) {
	now := time.Now()
	fileInfo := FileInfo{
		FilePath:    "/path/to/file",
		RequestHash: "request123",
		ResponseID:  "response456",
		ModTime:     now,
	}
	result := fileInfo.LogValue()
	expected := slog.GroupValue(
		slog.String("filePath", "/path/to/file"),
		slog.String("requestHash", "request123"),
		slog.String("responseId", "response456"),
		slog.Time("modTime", now),
	)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("LogValue() = %v, want %v", result, expected)
	}
}
