package request

import "errors"

var (
	ErrorInvalidURL              = errors.New("invalid URL passed")
	ErrorInvalidMethod           = errors.New("invalid method passed")
	ErrorUnsupportedType         = errors.New("unsupported type")
	ErrorFailedToConvertBodyVar  = errors.New("failed to convert bodyVar to string")
	ErrorFailedToCreateCookieJar = errors.New("failed to create cookie jar")
)
