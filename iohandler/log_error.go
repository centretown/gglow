package iohandler

import "log"

func LogError(cause string, err error) {
	log.Printf("%s:\n\t%v\n", cause, err)
}
