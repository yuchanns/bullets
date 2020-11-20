package common

import (
	"github.com/yuchanns/bullet/internal"
)

var DefaultLogger internal.ILogger

func init() {
	DefaultLogger = internal.NewBuiltinLogger()
}
