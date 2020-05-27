package main

import (
	"fmt"
	"github.com/pkg/errors"
)

func main() {
	err := errors.Wrap(errors.Wrap(errors.New("test"), "Hello world"), "This is a error package test")
	fmt.Println(err)
	fmt.Println(errors.Cause(err))
	fmt.Printf("%+v", err)

	err = errors.Wrapf(err, "this is %d level", 5)
	fmt.Println(err)

	err = errors.WithMessage(err, "You are stupid")
	fmt.Println(err)
}
