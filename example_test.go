package faultline_test

import (
	"fmt"

	"github.com/agentine/faultline"
)

func ExampleNew() {
	err := faultline.New("whoops")
	fmt.Println(err)
	// Output: whoops
}

func ExampleErrorf() {
	err := faultline.Errorf("count: %d", 42)
	fmt.Println(err)
	// Output: count: 42
}

func ExampleWrap() {
	cause := faultline.New("original")
	err := faultline.Wrap(cause, "context")
	fmt.Println(err)
	// Output: context: original
}

func ExampleCause() {
	cause := faultline.New("original")
	err := faultline.Wrap(cause, "context")
	fmt.Println(faultline.Cause(err))
	// Output: original
}

func ExampleWithStack() {
	cause := fmt.Errorf("simple error")
	err := faultline.WithStack(cause)
	fmt.Println(err)
	// Output: simple error
}

func ExampleWithMessage() {
	cause := faultline.New("original")
	err := faultline.WithMessage(cause, "added context")
	fmt.Println(err)
	// Output: added context: original
}
