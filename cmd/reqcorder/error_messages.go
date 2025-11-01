package main

import (
	"reqcorder/internal/diff"
	"reqcorder/internal/history"
	"reqcorder/internal/initiator"
	"reqcorder/internal/record"
	"reqcorder/internal/request"
	"reqcorder/pkg/utils"
)

var errorMessages = map[error]string{
	utils.ErrorFailedToCreateDirectory:   "failed to create directory",
	utils.ErrorFailedToMarshalJSON:       "failed to format output",
	utils.ErrorFailedToUnmarshalYAML:     "failed to read file contents",
	record.ErrorFailedToStatPath:         "failed to read required path",
	record.ErrorFailedToReadDirectory:    "failed to read required directory",
	record.ErrorPathIsNotDirectory:       "expected directory",
	record.ErrorFailedToGetRequest:       "failed to get request details",
	record.ErrorFailedToGetResponse:      "failed to get response details",
	record.ErrorFailedToGetTemplate:      "failed to get template details",
	diff.ErrorInvalidDiffType:            "invalid usage, invalid diff type",
	diff.ErrorFailedToRenderDiff:         "failed to render diff",
	history.ErrorFailedToParseTimestamp:  "failed to format timestamp",
	request.ErrorInvalidURL:              "invalid URL passed",
	request.ErrorInvalidMethod:           "invalid HTTP method passed",
	request.ErrorFailedToConvertBodyVar:  "failed to process body var(s)",
	request.ErrorFailedToCreateCookieJar: "failed to process request cookie(s)",
	initiator.ErrorFailedToReadCert:      "failed to read certificate path",
	initiator.ErrorFailedToBuildRequest:  "failed to process request",
	initiator.ErrorRequestFailed:         "failed to process request",
	initiator.ErrorFailedToReadResponse:  "failed to read response",
}
