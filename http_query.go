package main

import (
	"github.com/stretchr/goweb"
	"github.com/stretchr/goweb/context"
)

func (M *Main) HTTP_Query_Init() {
	M.HTTP_Query_Routes()
}

func (M *Main) HTTP_Query_Routes() {

	goweb.Map("GET", "/query/{accounts}/{groups}/{ts_start}/{ts_end}/{limit}/{filter}", func(c context.Context) (error) {
		DebugLn("QUERY")
		HTTP_Query(M,c)
		return nil
	})

}

func HTTP_Query(M *Main, c context.Context) {
	goweb.Respond.With(c,200,[]byte("pong"))
}
