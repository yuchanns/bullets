package common

import (
	"github.com/yuchanns/bullets/internal"
)

var Logger internal.ILogger

func init() {
	Logger = internal.NewBuiltinLogger()
}
