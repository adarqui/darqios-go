package main

import (
	"encoding/json"
	"github.com/stretchr/goweb"
	"github.com/stretchr/goweb/context"
)

func (M *Main) HTTP_State_Init() {
	M.HTTP_State_Routes()
}

func (M *Main) HTTP_State_Routes() {

	goweb.Map("GET", "/state/{accounts}/{path}", func(c context.Context) (error) {
		HTTP_State(M,c)
		return nil
	})
}

func HTTP_State(M *Main, c context.Context) {

	list := make([]string,0)
	for _, ses := range M.MS.Clients {
		if ses.Account != nil {
			list = append(list, ses.Account.Host)
		}
	}

	accounts, err := M.MG_Accounts(list, nil)
	if err != nil {
		return
	}

	jsn, err := json.Marshal(accounts)
	if err != nil {
		return
	}

	goweb.Respond.With(c, 200, []byte(jsn))
}
