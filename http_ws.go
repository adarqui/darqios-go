package main

import (
	"github.com/stretchr/goweb"
	"github.com/stretchr/goweb/context"
)

func (M *Main) HTTP_WS_Init() {
	M.HTTP_WS_Routes()
}

func (M *Main) HTTP_WS_Routes() {

	goweb.Map("UPDATE", "/ws", func(c context.Context) (error) {
		DebugLn("WS")
		WS_Serve(M, c.HttpResponseWriter(), c.HttpRequest())
		return nil
	})

	goweb.Map("GET", "/ws", func(c context.Context) (error) {
		DebugLn("WS")
		WS_Serve(M, c.HttpResponseWriter(), c.HttpRequest())
		return nil
	})

}
