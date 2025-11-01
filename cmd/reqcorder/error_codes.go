package main

import (
	"reqcorder/internal/diff"
	"reqcorder/internal/history"
	"reqcorder/internal/initiator"
	"reqcorder/internal/record"
	"reqcorder/internal/request"
	"reqcorder/pkg/utils"
)

var errorCodes = map[error]int{
	// File IO errors
	ErrorFailedToReadHomeDirectory:     1,
	ErrorFailedToOpenLogFile:           1,
	utils.ErrorFailedToCreateDirectory: 1,
	record.ErrorFailedToStatPath:       1,
	record.ErrorFailedToReadDirectory:  1,
	record.ErrorPathIsNotDirectory:     1,
	// Usage errors
	ErrorInvalidUsage:          2,
	diff.ErrorInvalidDiffType:  2,
	ErrorInvalidShowType:       2,
	ErrorInvalidListType:       2,
	request.ErrorInvalidURL:    2,
	request.ErrorInvalidMethod: 2,
	// Processing data errors
	utils.ErrorFailedToMarshalJSON:       3,
	history.ErrorFailedToParseTimestamp:  3,
	utils.ErrorFailedToUnmarshalYAML:     3,
	request.ErrorFailedToConvertBodyVar:  3,
	request.ErrorFailedToCreateCookieJar: 3,
	initiator.ErrorFailedToReadCert:      3,
	initiator.ErrorFailedToBuildRequest:  3,
	initiator.ErrorRequestFailed:         3,
	initiator.ErrorFailedToReadResponse:  3,
	// Broad fetching errors
	record.ErrorFailedToGetRequest:  4,
	record.ErrorFailedToGetTemplate: 4,
	record.ErrorFailedToGetResponse: 4,
	// Rendering errors
	diff.ErrorFailedToRenderDiff: 5,
}
