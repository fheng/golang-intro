package main

import (
	"fmt"

	"github.com/pkg/errors"
)

func main() {
	err := errors.New("my error with stack")
	fmt.Printf("%+v", err) //the +v indicates to print the trace
}
