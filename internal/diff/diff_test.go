package diff

import (
	"bytes"
	"errors"
	"reqcorder/internal/record"
	"reqcorder/internal/request"
	"reqcorder/internal/response"
	"strings"
	"testing"
)

func TestSuccessfulDefaultDiff(t *testing.T) {
	root := t.TempDir()
	templateYaml := []byte("key: value")
	requestYaml1 := []byte("url: https://example1.com\nmethod: GET")
	requestYaml2 := []byte("url: https://example2.com\nmethod: POST")
	responseYaml1 := []byte("status_code: 200\nbody: response1")
	responseYaml2 := []byte("status_code: 404\nbody: response2")
	recordStore1 := &record.RecordStore{
		RecordStorePath: root,
		TemplateYaml:    templateYaml,
		RequestYaml:     requestYaml1,
		ResponseYaml:    responseYaml1,
		Request: &request.RequestObject{
			URL:    "https://example1.com",
			Method: "GET",
		},
		Response: &response.ResponseObject{
			StatusCode: 200,
		},
	}
	err := recordStore1.Record()
	if err != nil {
		t.Fatalf("Expected no error, received %v\n", err)
	}
	recordStore2 := &record.RecordStore{
		RecordStorePath: root,
		TemplateYaml:    templateYaml,
		RequestYaml:     requestYaml2,
		ResponseYaml:    responseYaml2,
		Request: &request.RequestObject{
			URL:    "https://example2.com",
			Method: "POST",
		},
		Response: &response.ResponseObject{
			StatusCode: 404,
		},
	}
	err = recordStore2.Record()
	if err != nil {
		t.Fatalf("Expected no error, received %v\n", err)
	}
	var buf bytes.Buffer
	err = DefaultDiff(&buf, root, recordStore1.ResponseID, recordStore2.ResponseID, "response")
	if err != nil {
		t.Fatalf("Expected no error, received %v\n", err)
	}
	output := buf.String()
	if !strings.Contains(output, "---") || !strings.Contains(output, "+++") {
		t.Fatal("Expected diff output to contain git-style diff headers")
	}
	buf.Reset()
	err = DefaultDiff(&buf, root, recordStore1.RequestHash, recordStore2.RequestHash, "request")
	if err != nil {
		t.Fatalf("Expected no error, received %v\n", err)
	}
	output = buf.String()
	if !strings.Contains(output, "---") || !strings.Contains(output, "+++") {
		t.Fatal("Expected diff output to contain git-style diff headers")
	}
	buf.Reset()
	err = DefaultDiff(&buf, root, recordStore1.TemplateHash, recordStore2.TemplateHash, "template")
	if err != nil {
		t.Fatalf("Expected no error, received %v\n", err)
	}
	output = buf.String()
	if !strings.Contains(output, "---") || !strings.Contains(output, "+++") {
		t.Fatal("Expected diff output to contain git-style diff headers")
	}
}

func TestFailedDefaultDiff_InvalidResourceType(t *testing.T) {
	root := t.TempDir()
	var buf bytes.Buffer
	err := DefaultDiff(&buf, root, "source", "target", "invalid")
	if err == nil {
		t.Fatal("Expected error, received nil")
	}
	expectedErr := ErrorInvalidDiffType
	if !errors.Is(err, expectedErr) {
		t.Fatalf("Expected error %v, received %v\n", expectedErr, err)
	}
}

func TestFailedDefaultDiff_RecordFailure(t *testing.T) {
	root := t.TempDir()
	var buf bytes.Buffer
	err := DefaultDiff(&buf, root, "nonexistent", "nonexistent", "response")
	if err == nil {
		t.Fatal("Expected error, received nil")
	}
	expectedErr := record.ErrorFailedToStatPath
	if !errors.Is(err, expectedErr) {
		t.Fatalf("Expected error %v, received %v\n", expectedErr, err)
	}
}

func TestSuccessfulInlineDiff(t *testing.T) {
	root := t.TempDir()
	templateYaml := []byte("key: value")
	requestYaml1 := []byte("url: https://example1.com\nmethod: GET")
	requestYaml2 := []byte("url: https://example2.com\nmethod: POST")
	responseYaml1 := []byte("status_code: 200\nbody: response1")
	responseYaml2 := []byte("status_code: 404\nbody: response2")
	recordStore1 := &record.RecordStore{
		RecordStorePath: root,
		TemplateYaml:    templateYaml,
		RequestYaml:     requestYaml1,
		ResponseYaml:    responseYaml1,
		Request: &request.RequestObject{
			URL:    "https://example1.com",
			Method: "GET",
		},
		Response: &response.ResponseObject{
			StatusCode: 200,
		},
	}
	err := recordStore1.Record()
	if err != nil {
		t.Fatalf("Expected no error, received %v\n", err)
	}
	recordStore2 := &record.RecordStore{
		RecordStorePath: root,
		TemplateYaml:    templateYaml,
		RequestYaml:     requestYaml2,
		ResponseYaml:    responseYaml2,
		Request: &request.RequestObject{
			URL:    "https://example2.com",
			Method: "POST",
		},
		Response: &response.ResponseObject{
			StatusCode: 404,
		},
	}
	err = recordStore2.Record()
	if err != nil {
		t.Fatalf("Expected no error, received %v\n", err)
	}
	var buf bytes.Buffer
	err = InlineDiff(&buf, root, recordStore1.ResponseID, recordStore2.ResponseID, "response")
	if err != nil {
		t.Fatalf("Expected no error, received %v\n", err)
	}
	output := buf.String()
	if !strings.Contains(output, ColorGreen) || !strings.Contains(output, ColorRed) {
		t.Fatal("Expected diff output to contain color codes")
	}
	buf.Reset()
	err = InlineDiff(&buf, root, recordStore1.RequestHash, recordStore2.RequestHash, "request")
	if err != nil {
		t.Fatalf("Expected no error, received %v\n", err)
	}
	output = buf.String()
	if !strings.Contains(output, ColorGreen) || !strings.Contains(output, ColorRed) {
		t.Fatal("Expected diff output to contain color codes")
	}
	buf.Reset()
	err = InlineDiff(&buf, root, recordStore1.TemplateHash, recordStore2.TemplateHash, "template")
	if err != nil {
		t.Fatalf("Expected no error, received %v\n", err)
	}
	output = buf.String()
	if !strings.Contains(output, ColorGreen) || !strings.Contains(output, ColorRed) {
		t.Fatal("Expected diff output to contain color codes")
	}
}

func TestFailedInlineDiff_InvalidResourceType(t *testing.T) {
	root := t.TempDir()
	var buf bytes.Buffer
	err := InlineDiff(&buf, root, "source", "target", "invalid")
	if err == nil {
		t.Fatal("Expected error, received nil")
	}
	expectedErr := ErrorInvalidDiffType
	if !errors.Is(err, expectedErr) {
		t.Fatalf("Expected error %v, received %v\n", expectedErr, err)
	}
}

func TestFailedInlineDiff_RecordFailure(t *testing.T) {
	root := t.TempDir()
	var buf bytes.Buffer
	err := InlineDiff(&buf, root, "nonexistent", "nonexistent", "response")
	if err == nil {
		t.Fatal("Expected error, received nil")
	}
	expectedErr := record.ErrorFailedToStatPath
	if !errors.Is(err, expectedErr) {
		t.Fatalf("Expected error %v, received %v\n", expectedErr, err)
	}
}

func TestSuccessfulGetResponseByID(t *testing.T) {
	root := t.TempDir()
	recordStore := &record.RecordStore{
		RecordStorePath: root,
		TemplateYaml:    []byte("key: value"),
		Request: &request.RequestObject{
			URL:    "https://example.com",
			Method: "GET",
		},
		Response: &response.ResponseObject{
			StatusCode: 200,
		},
	}
	err := recordStore.Record()
	if err != nil {
		t.Fatalf("Expected no error, received %v\n", err)
	}
	diffStore := &DiffStore{
		RecordStorePath: root,
	}
	content, err := diffStore.getResponseByID(recordStore.ResponseID)
	if err != nil {
		t.Fatalf("Expected no error, received %v\n", err)
	}
	if content == "" {
		t.Fatal("Expected non-empty content, received empty string")
	}
}

func TestFailedGetResponseByID_RecordFailure(t *testing.T) {
	root := t.TempDir()
	diffStore := &DiffStore{
		RecordStorePath: root,
	}
	_, err := diffStore.getResponseByID("nonexistent")
	if err == nil {
		t.Fatal("Expected error, received nil")
	}
	expectedErr := record.ErrorFailedToStatPath
	if !errors.Is(err, expectedErr) {
		t.Fatalf("Expected error %v, received %v\n", expectedErr, err)
	}
}

func TestSuccessfulGetRequestByHash(t *testing.T) {
	root := t.TempDir()
	requestYaml := []byte("url: https://example.com\nmethod: GET")
	requestHash := "one"
	recordStore := &record.RecordStore{
		TemplateHash:    requestHash,
		TemplateYaml:    requestYaml,
		RecordStorePath: root,
		RequestYaml:     requestYaml,
		Request: &request.RequestObject{
			URL:          "https://example.com",
			Method:       "GET",
			TemplateHash: requestHash,
		},
		RequestHash: requestHash,
		Response: &response.ResponseObject{
			StatusCode: 200,
		},
	}
	err := recordStore.Record()
	if err != nil {
		t.Fatalf("Expected no error, received %v\n", err)
	}
	diffStore := &DiffStore{
		RecordStorePath: root,
	}
	_, err = diffStore.getRequestByHash(recordStore.Response.RequestHash)
	if err != nil {
		t.Fatalf("Expected no error, received %v\n", err)
	}
}

func TestFailedGetRequestByHash_RecordFailure(t *testing.T) {
	root := t.TempDir()
	diffStore := &DiffStore{
		RecordStorePath: root,
	}
	_, err := diffStore.getRequestByHash("nonexistent")
	if err == nil {
		t.Fatal("Expected error, received nil")
	}
	expectedErr := record.ErrorFailedToStatPath
	if !errors.Is(err, expectedErr) {
		t.Fatalf("Expected error %v, received %v\n", expectedErr, err)
	}
}

func TestSuccessfulGetTemplateByHash(t *testing.T) {
	root := t.TempDir()
	templateYaml := []byte("key: value")
	templateHash := "one"
	recordStore := &record.RecordStore{
		TemplateHash:    templateHash,
		TemplateYaml:    templateYaml,
		RecordStorePath: root,
		Request: &request.RequestObject{
			URL:          "https://example.com",
			Method:       "GET",
			TemplateHash: templateHash,
		},
		Response: &response.ResponseObject{
			StatusCode: 200,
		},
	}
	err := recordStore.Record()
	if err != nil {
		t.Fatalf("Expected no error, received %v\n", err)
	}

	diffStore := &DiffStore{
		RecordStorePath: root,
	}
	content, err := diffStore.getTemplateByHash(recordStore.TemplateHash)
	if err != nil {
		t.Fatalf("Expected no error, received %v\n", err)
	}
	if content == "" {
		t.Fatal("Expected non-empty content, received empty string")
	}
}

func TestFailedGetTemplateByHash_RecordFailure(t *testing.T) {
	root := t.TempDir()
	diffStore := &DiffStore{
		RecordStorePath: root,
	}
	_, err := diffStore.getTemplateByHash("nonexistent")
	if err == nil {
		t.Fatal("Expected error, received nil")
	}
	expectedErr := record.ErrorFailedToStatPath
	if !errors.Is(err, expectedErr) {
		t.Fatalf("Expected error %v, received %v\n", expectedErr, err)
	}
}

func TestSuccessfulGetTexts(t *testing.T) {
	root := t.TempDir()
	templateYaml := []byte("key: value")
	recordStore1 := &record.RecordStore{
		RecordStorePath: root,
		TemplateYaml:    templateYaml,
		Request: &request.RequestObject{
			URL:    "https://example1.com",
			Method: "GET",
		},
		Response: &response.ResponseObject{
			StatusCode: 200,
			Body:       "response1",
		},
	}
	err := recordStore1.Record()
	if err != nil {
		t.Fatalf("Expected no error, received %v\n", err)
	}
	recordStore2 := &record.RecordStore{
		RecordStorePath: root,
		TemplateYaml:    templateYaml,
		Request: &request.RequestObject{
			URL:    "https://example2.com",
			Method: "POST",
		},
		Response: &response.ResponseObject{
			StatusCode: 404,
			Body:       "response2",
		},
	}
	err = recordStore2.Record()
	if err != nil {
		t.Fatalf("Expected no error, received %v\n", err)
	}
	text1, text2, err := getTexts(root, recordStore1.ResponseID, recordStore2.ResponseID, "response")
	if err != nil {
		t.Fatalf("Expected no error, received %v\n", err)
	}
	if text1 == "" {
		t.Fatal("Expected non-empty text1, received empty string")
	}
	if text2 == "" {
		t.Fatal("Expected non-empty text2, received empty string")
	}
	text1, text2, err = getTexts(root, recordStore1.RequestHash, recordStore2.RequestHash, "request")
	if err != nil {
		t.Fatalf("Expected no error, received %v\n", err)
	}
	if text1 == "" {
		t.Fatal("Expected non-empty text1, received empty string")
	}
	if text2 == "" {
		t.Fatal("Expected non-empty text2, received empty string")
	}
	text1, text2, err = getTexts(root, recordStore1.TemplateHash, recordStore2.TemplateHash, "template")
	if err != nil {
		t.Fatalf("Expected no error, received %v\n", err)
	}
	if text1 == "" {
		t.Fatal("Expected non-empty text1, received empty string")
	}
	if text2 == "" {
		t.Fatal("Expected non-empty text2, received empty string")
	}
}

func TestFailedGetTexts_InvalidResourceType(t *testing.T) {
	root := t.TempDir()
	_, _, err := getTexts(root, "source", "target", "invalid")
	if err == nil {
		t.Fatal("Expected error, received nil")
	}
	expectedErr := ErrorInvalidDiffType
	if !errors.Is(err, expectedErr) {
		t.Fatalf("Expected error %v, received %v\n", expectedErr, err)
	}
}

func TestFailedGetTexts_RecordFailure(t *testing.T) {
	root := t.TempDir()
	_, _, err := getTexts(root, "nonexistent", "nonexistent", "response")
	if err == nil {
		t.Fatal("Expected error, received nil")
	}
	expectedErr := record.ErrorFailedToStatPath
	if !errors.Is(err, expectedErr) {
		t.Fatalf("Expected error %v, received %v\n", expectedErr, err)
	}
}

func TestFailedGetTexts_ResponseFailure(t *testing.T) {
	root := t.TempDir()
	recordStore := &record.RecordStore{
		RecordStorePath: root,
		TemplateYaml:    []byte("key: value"),
		Request: &request.RequestObject{
			URL:    "https://example.com",
			Method: "GET",
		},
		Response: &response.ResponseObject{
			StatusCode: 200,
		},
	}
	err := recordStore.Record()
	if err != nil {
		t.Fatalf("Expected no error, received %v\n", err)
	}
	_, _, err = getTexts(root, recordStore.ResponseID, "nonexistent", "response")
	if err == nil {
		t.Fatal("Expected error, received nil")
	}
	expectedErr := record.ErrorFailedToGetResponse
	if !errors.Is(err, expectedErr) {
		t.Fatalf("Expected error %v, received %v\n", expectedErr, err)
	}
}

func TestFailedGetTexts_RequestFailure(t *testing.T) {
	root := t.TempDir()
	recordStore := &record.RecordStore{
		RecordStorePath: root,
		TemplateYaml:    []byte("key: value"),
		Request: &request.RequestObject{
			URL:    "https://example.com",
			Method: "GET",
		},
		Response: &response.ResponseObject{
			StatusCode: 200,
		},
	}
	err := recordStore.Record()
	if err != nil {
		t.Fatalf("Expected no error, received %v\n", err)
	}
	_, _, err = getTexts(root, recordStore.RequestHash, "nonexistent", "request")
	if err == nil {
		t.Fatal("Expected error, received nil")
	}
	expectedErr := record.ErrorFailedToGetRequest
	if !errors.Is(err, expectedErr) {
		t.Fatalf("Expected error %v, received %v\n", expectedErr, err)
	}
}

func TestFailedGetTexts_TemplateFailure(t *testing.T) {
	root := t.TempDir()
	recordStore := &record.RecordStore{
		RecordStorePath: root,
		TemplateYaml:    []byte("key: value"),
		Request: &request.RequestObject{
			URL:    "https://example.com",
			Method: "GET",
		},
		Response: &response.ResponseObject{
			StatusCode: 200,
		},
	}
	err := recordStore.Record()
	if err != nil {
		t.Fatalf("Expected no error, received %v\n", err)
	}

	_, _, err = getTexts(root, recordStore.TemplateHash, "nonexistent", "template")
	if err == nil {
		t.Fatal("Expected error, received nil")
	}
	expectedErr := record.ErrorFailedToGetTemplate
	if !errors.Is(err, expectedErr) {
		t.Fatalf("Expected error %v, received %v\n", expectedErr, err)
	}
}

func TestSuccessfulGitStyleDiff(t *testing.T) {
	text1 := "line1\nline2\nline3"
	text2 := "line1\nline2 modified\nline3\nline4"
	var buf bytes.Buffer
	err := gitStyleDiff(&buf, text1, text2, "file1", "file2")
	if err != nil {
		t.Fatalf("Expected no error, received %v\n", err)
	}

	output := buf.String()
	if !strings.Contains(output, "--- file1") || !strings.Contains(output, "+++ file2") {
		t.Fatal("Expected diff output to contain git-style diff headers")
	}
	if !strings.Contains(output, ColorRed) || !strings.Contains(output, ColorGreen) {
		t.Fatal("Expected diff output to contain color codes")
	}
}

func TestFailedGitStyleDiff_WriteFailure(t *testing.T) {
	writer := &failingWriter{}
	err := gitStyleDiff(writer, "text1", "text2", "file1", "file2")
	if err == nil {
		t.Fatal("Expected error, received nil")
	}
	expectedErr := ErrorFailedToRenderDiff
	if !errors.Is(err, expectedErr) {
		t.Fatalf("Expected error %v, received %v\n", expectedErr, err)
	}
}

func TestSuccessfulInlineDiffFunction(t *testing.T) {
	text1 := "line1\nline2\nline3"
	text2 := "line1\nline2 modified\nline3\nline4"
	var buf bytes.Buffer
	err := inlineDiff(&buf, text1, text2, "file1", "file2")
	if err != nil {
		t.Fatalf("Expected no error, received %v\n", err)
	}
	output := buf.String()
	if !strings.Contains(output, ColorGreen+"file2"+ColorReset) || !strings.Contains(output, ColorRed+"file1"+ColorReset) {
		t.Fatal("Expected diff output to contain colored file names")
	}
}

func TestFailedInlineDiffFunction_WriteFailure(t *testing.T) {
	writer := &failingWriter{}
	err := inlineDiff(writer, "text1", "text2", "file1", "file2")
	if err == nil {
		t.Fatal("Expected error, received nil")
	}
	expectedErr := ErrorFailedToRenderDiff
	if !errors.Is(err, expectedErr) {
		t.Fatalf("Expected error %v, received %v\n", expectedErr, err)
	}
}

func TestFailedGitStyleDiff_FirstHeaderWriteFailure(t *testing.T) {
	writer := &failingWriterOnFirstCall{}
	err := gitStyleDiff(writer, "text1", "text2", "file1", "file2")
	if err == nil {
		t.Fatal("Expected error, received nil")
	}
	expectedErr := ErrorFailedToRenderDiff
	if !errors.Is(err, expectedErr) {
		t.Fatalf("Expected error %v, received %v\n", expectedErr, err)
	}
}

func TestFailedGitStyleDiff_SecondHeaderWriteFailure(t *testing.T) {
	writer := &failingWriterOnSecondCall{}
	err := gitStyleDiff(writer, "text1", "text2", "file1", "file2")
	if err == nil {
		t.Fatal("Expected error, received nil")
	}
	expectedErr := ErrorFailedToRenderDiff
	if !errors.Is(err, expectedErr) {
		t.Fatalf("Expected error %v, received %v\n", expectedErr, err)
	}
}

func TestFailedInlineDiffFunction_FirstHeaderWriteFailure(t *testing.T) {
	writer := &failingWriterOnFirstCall{}
	err := inlineDiff(writer, "text1", "text2", "file1", "file2")
	if err == nil {
		t.Fatal("Expected error, received nil")
	}
	expectedErr := ErrorFailedToRenderDiff
	if !errors.Is(err, expectedErr) {
		t.Fatalf("Expected error %v, received %v\n", expectedErr, err)
	}
}

func TestFailedInlineDiffFunction_SecondHeaderWriteFailure(t *testing.T) {
	writer := &failingWriterOnSecondCall{}
	err := inlineDiff(writer, "text1", "text2", "file1", "file2")
	if err == nil {
		t.Fatal("Expected error, received nil")
	}
	expectedErr := ErrorFailedToRenderDiff
	if !errors.Is(err, expectedErr) {
		t.Fatalf("Expected error %v, received %v\n", expectedErr, err)
	}
}

func TestFailedInlineDiffFunction_DiffTextWriteFailure(t *testing.T) {
	writer := &failingWriterOnThirdCall{}
	err := inlineDiff(writer, "text1", "text2", "file1", "file2")
	if err == nil {
		t.Fatal("Expected error, received nil")
	}
	expectedErr := ErrorFailedToRenderDiff
	if !errors.Is(err, expectedErr) {
		t.Fatalf("Expected error %v, received %v\n", expectedErr, err)
	}
}

type failingWriterOnFirstCall struct {
	callCount int
}

func (fw *failingWriterOnFirstCall) Write(p []byte) (n int, err error) {
	fw.callCount++
	if fw.callCount == 1 {
		return 0, errors.New("write failed")
	}
	return len(p), nil
}

type failingWriterOnSecondCall struct {
	callCount int
}

func (fw *failingWriterOnSecondCall) Write(p []byte) (n int, err error) {
	fw.callCount++
	if fw.callCount == 2 {
		return 0, errors.New("write failed")
	}
	return len(p), nil
}

type failingWriterOnThirdCall struct {
	callCount int
}

func (fw *failingWriterOnThirdCall) Write(p []byte) (n int, err error) {
	fw.callCount++
	if fw.callCount == 3 {
		return 0, errors.New("write failed")
	}
	return len(p), nil
}

type failingWriter struct{}

func (fw *failingWriter) Write(p []byte) (n int, err error) {
	return 0, errors.New("write failed")
}
