package main

import (
	"github.com/stretchr/goweb"
	"github.com/stretchr/goweb/context"
)

func (M *Main) HTTP_Help_Init() {
	M.HTTP_Help_Routes()
}

func (M *Main) HTTP_Help_Routes() {

	goweb.Map("GET", "/help", func(c context.Context) (error) {
		DebugLn("PING")
		HTTP_Help(M,c)
		return nil
	})

}

func HTTP_Help(M *Main, c context.Context) {
	goweb.Respond.With(c,200,[]byte("pong"))
}
