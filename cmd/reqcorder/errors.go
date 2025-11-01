package main

import "errors"

var (
	ErrorInvalidUsage              = errors.New("invalid usage")
	ErrorFailedToReadHomeDirectory = errors.New("failed to read home directory for current user")
	ErrorInvalidShowType           = errors.New("invalid usage, invalid show type")
	ErrorInvalidListType           = errors.New("invalid usage, invalid list type")
	ErrorFailedToOpenLogFile          = errors.New("failed to open log file")
)
