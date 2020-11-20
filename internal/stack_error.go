package internal

import (
	"fmt"
	"github.com/pkg/errors"
	"strings"
)

type stackTracer interface {
	StackTrace() errors.StackTrace
}

// BuildStack will build and return the stack of an error that implementing stackTracer
func BuildStack(err error, skip int) []string {
	traces := make([]string, 0)
	if err, ok := err.(stackTracer); ok {
		frames := err.StackTrace()
		l := len(frames)
		for i := skip; i <= skip+4; i++ {
			if i >= l {
				break
			}
			f := frames[i]
			traces = append(traces, strings.ReplaceAll(
				strings.TrimSpace(fmt.Sprintf("%+s:%d", f, f)),
				"\n\t",
				" ",
			))
		}
	}
	return traces
}
