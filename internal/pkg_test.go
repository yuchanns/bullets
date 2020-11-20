package internal

import (
	"fmt"
	"github.com/pkg/errors"
	"testing"
)

func TestBuildStack(t *testing.T) {
	tables := []error{
		errors.New("a new error"),
		errors.Wrap(errors.New("a new error"), "a wrap error"),
		errors.Wrapf(errors.New("a new error"), "a wrap error with %s", "format"),
	}
	for i := range tables {
		stack := BuildStack(tables[i], 0)
		for i := range stack {
			fmt.Println(tables[i].Error())
			fmt.Println(stack[i])
		}
	}
}
