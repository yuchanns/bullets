package common

import (
	"github.com/yuchanns/bullet/internal"
)

var Logger internal.ILogger

func init() {
	Logger = internal.NewBuiltinLogger()
}
