package main

import (
	"log"
	"fmt"
	"os"
)

/*
const (
	DEBUG = false
)
*/
var DEBUG bool

var xstd = log.New(os.Stderr, "", log.LstdFlags)

func Debug(format string, v ... interface{}) {
	if DEBUG == true {
		xstd.Output(2, fmt.Sprintf(format, v...))
	}
}

func DebugFatal(Status int, format string, v ... interface{}) {
	if DEBUG == true {
		xstd.Output(2, fmt.Sprintf(format, v...))
		os.Exit(Status)
	}
}

func DebugLn(v ... interface{}) {
	if DEBUG == true {
		xstd.Output(2, fmt.Sprintln(v...))
	}
}
