package main

import (
	"encoding/json"
	"github.com/stretchr/goweb"
	"github.com/stretchr/goweb/context"
)

func (M *Main) HTTP_Sessions_Init() {
	M.HTTP_Sessions_Routes()
}

func (M *Main) HTTP_Sessions_Routes() {

	sessions_route_list := []string{"/sessions","/sessions/"}
	for _, route := range sessions_route_list {
		goweb.Map("GET", route, func(c context.Context) (error) {
			HTTP_Sessions(M,c)
			return nil
		})
	}
}

func HTTP_Sessions(M *Main, c context.Context) {

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
