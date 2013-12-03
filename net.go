package main

import (
	"net"
	"crypto/tls"
)

type Net struct {
	Conn net.Conn
	TlsConn *tls.Conn
	Count int
}
