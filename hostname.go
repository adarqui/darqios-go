package main

import (
	"syscall"
	"io/ioutil"
)

func setHostname(name string) {
	syscall.Sethostname([]byte(name))
	ioutil.WriteFile("/etc/hostname", []byte(name), 0755)
}
