package record

import "errors"

var (
	ErrorFailedToConvertRequest  = errors.New("failed to convert request")
	ErrorFailedToConvertResponse = errors.New("failed to convert response")
	ErrorFailedToRecord          = errors.New("failed to create record store entry")
	ErrorFailedToStatPath        = errors.New("failed to get information of path")
	ErrorPathIsNotDirectory      = errors.New("path is not directory")
	ErrorFailedToReadDirectory   = errors.New("failed to read directory")
	ErrorFailedToGetRequest      = errors.New("failed to get request")
	ErrorFailedToGetResponse     = errors.New("failed to get response")
	ErrorFailedToGetTemplate     = errors.New("failed to get template")
)
