package utils

import "errors"

var (
	ErrorFailedToReadFile        = errors.New("failed to read file")
	ErrorFailedToUnmarshalYAML   = errors.New("failed to unmarshal from YAML")
	ErrorFailedToMarshalJSON     = errors.New("failed to marshal to JSON")
	ErrorFailedToMarshalYAML     = errors.New("failed to marshal to YAML")
	ErrorFailedToCreateDirectory = errors.New("failed to create directory")
)
