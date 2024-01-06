package iohandler

import "fmt"

func LogError(cause string, err error) {
	fmt.Printf("%s:\n\t%v\n", cause, err)
}
