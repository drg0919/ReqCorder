package diff

import "errors"

var (
	ErrorFailedToRenderDiff = errors.New("failed to render diff")
	ErrorInvalidDiffType    = errors.New("invalid diff type")
)
