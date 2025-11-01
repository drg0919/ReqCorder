package initiator

import "errors"

var (
	ErrorFailedToBuildRequest = errors.New("failed to build request")
	ErrorFailedToReadCert     = errors.New("failed to read cert")
	ErrorRequestFailed        = errors.New("request failed")
	ErrorFailedToReadResponse = errors.New("failed to read response")
)
