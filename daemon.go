package main

import (
	"github.com/VividCortex/godaemon"
)

func Daemon() {
	godaemon.MakeDaemon(&godaemon.DaemonAttr{})
}
