package app

import (
	"errors"
	"fmt"
	"testing"
)

func TestError(t *testing.T) {
	err := NewRootError(100, 233, "系统错误").
		WithError(errors.New("mysql connect failed")).
		WithError(errors.New("connection timeout")).
		WithError(errors.New("error test")).
		WithErrorDesc("hello world")
	fmt.Printf("\n\n\n%s\n\n\n", err)
}
