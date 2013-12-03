package main

import (
	"github.com/stretchr/goweb"
	"github.com/stretchr/goweb/context"
)

func (M *Main) HTTP_Ping_Init() {
	M.HTTP_Ping_Routes()
}

func (M *Main) HTTP_Ping_Routes() {

	goweb.Map("GET", "/ping", func(c context.Context) (error) {
		DebugLn("PING")
		HTTP_Ping(M,c)
		return nil
	})

}

func HTTP_Ping(M *Main, c context.Context) {
	goweb.Respond.With(c,200,[]byte("pong"))
}
