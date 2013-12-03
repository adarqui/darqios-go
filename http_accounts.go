package main

import (
	"time"
	"encoding/json"
	"github.com/stretchr/goweb"
	"github.com/stretchr/goweb/context"
)

func (M *Main) HTTP_Accounts_Init() {
	M.HTTP_Accounts_Routes()
}

func (M *Main) HTTP_Accounts_Routes() {

	goweb.Map("GET", "/accounts:list", func(c context.Context) (error) {
		DebugLn("ACCOUNTS")
		HTTP_Accounts_Index(M, c)
		return nil
	})

	accounts_get_routes := []string{"/accounts/get","/accounts/get/","/accounts/get/{accounts}"}

	for _, route := range accounts_get_routes {
		goweb.Map("GET", route, func(c context.Context) (error) {
			DebugLn("ACCOUNTS GET")
			HTTP_Accounts_Get(M, c)
			return nil
		})
	}


	ignore_get_routes := []string{"/accounts/ignore","/accounts/ignore/"}
	for _, route := range ignore_get_routes {
		goweb.Map("GET", route, func(c context.Context) (error) {
			DebugLn("ACCOUNTS IGNORE LIST")
			HTTP_Accounts_Ignore_List(M,c)
			return nil
		})
	}

	goweb.Map("GET", "/accounts/ignore/{accounts}", func(c context.Context) (error) {
		DebugLn("ACCOUNTS IGNORE")
		HTTP_Accounts_Ignore(M, c)
		return nil
	})

	goweb.Map("GET", "/accounts/unignore/{accounts}", func(c context.Context) (error) {
		DebugLn("ACCOUNTS UNIGNORE")
		HTTP_Accounts_Unignore(M, c)
		return nil
	})

	enable_get_routes := []string{"/accounts/enable","/accounts/enable/"}
	for _, route := range enable_get_routes {
		goweb.Map("GET",route,func(c context.Context) (error) {
			HTTP_Accounts_Enable_List(M,c)
			return nil
		})
	}

	disable_get_routes := []string{"/accounts/disable","/accounts/disable/"}
	for _,route := range disable_get_routes {
		goweb.Map("GET",route,func(c context.Context) (error) {
			HTTP_Accounts_Disable_List(M,c)
			return nil
		})
	}

	goweb.Map("GET", "/accounts/enable/{accounts}", func(c context.Context) (error) {
		HTTP_Accounts_Enable(M,c)
		return nil
	})

	goweb.Map("GET", "/accounts/disable/{accounts}", func(c context.Context) (error) {
		HTTP_Accounts_Disable(M,c)
		return nil
	})


	missing_get_routes := []string{"/accounts/missing","/accounts/missing/"}
	for _,route := range missing_get_routes {
		goweb.Map("GET", route, func(c context.Context) (error) {
			HTTP_Accounts_Missing(M,c)
			return nil
		})
	}

	goweb.Map("GET", "/accounts/add/{hash}/{host}/{groups}/{status}", func (c context.Context) (error) {
		HTTP_Accounts_Add(M,c)
		return nil
	})
}

type Accounts_List struct {
	Accounts []string
}

func HTTP_Accounts_Index(M *Main, c context.Context) {
	accounts,err := M.MG_Accounts(nil, nil)
	if err != nil {
		return
	}

	accounts_list := Accounts_List{}
	accounts_list.Accounts = make([]string, len(accounts))

	i:=0
	for _,acc := range accounts {
		accounts_list.Accounts[i] = acc.Host
		i++
	}

	jsn, err := json.Marshal(&accounts_list)
	if err != nil {
		return
	}

	goweb.Respond.With(c, 200, []byte(jsn))
}


func HTTP_Accounts_Get(M *Main, c context.Context) {

	list,_ := HTTP_Filter(c, "accounts")

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


func HTTP_Accounts_Ignore_List(M *Main, c context.Context) {
	accounts, err := M.MG_Accounts_Find_Ignore()
	if err != nil {
		return
	}

	jsn, err := json.Marshal(accounts)
	goweb.Respond.With(c, 200, []byte(jsn))
}

func HTTP_Accounts_Ignore(M *Main, c context.Context) {
	HTTP_Accounts_IgUnig_Nore(M, c, true)
}

func HTTP_Accounts_Unignore(M *Main, c context.Context) {
	HTTP_Accounts_IgUnig_Nore(M, c, false)
}

/* HAH */
func HTTP_Accounts_IgUnig_Nore(M *Main, c context.Context, Truth bool) {

	list,_ := HTTP_Filter(c, "accounts")

	_, err := M.MG_Accounts_Ignore(list, nil, Truth)
	if err != nil {
		return
	}

	goweb.Respond.With(c, 200, []byte("+OK"))
}





func HTTP_Accounts_Enable_List(M *Main, c context.Context) {
	accounts, err := M.MG_Accounts_Find_Enable()
	if err != nil {
		return
	}

	jsn, err := json.Marshal(accounts)
	goweb.Respond.With(c, 200, []byte(jsn))
}

func HTTP_Accounts_Disable_List(M *Main, c context.Context) {
	accounts, err := M.MG_Accounts_Find_Disable()
	if err != nil {
		return
	}
	jsn, err := json.Marshal(accounts)
	goweb.Respond.With(c, 200, []byte(jsn))
}

func HTTP_Accounts_Enable(M *Main, c context.Context) {
	HTTP_Accounts_AcInac_Tive(M, c, true)
}

func HTTP_Accounts_Disable(M *Main, c context.Context) {
	HTTP_Accounts_AcInac_Tive(M, c, false)
}

/* HAH */
func HTTP_Accounts_AcInac_Tive(M *Main, c context.Context, Truth bool) {

	list,_ := HTTP_Filter(c, "accounts")

	_, err := M.MG_Accounts_Enable(list, nil, Truth)
	if err != nil {
		return
	}

	goweb.Respond.With(c, 200, []byte("+OK"))
}


func HTTP_Accounts_Missing(M *Main, c context.Context) {
	t := time.Now()
	accounts, err := M.MG_Accounts_Find_Missing(t)
	if err != nil {
		return
	}

	missing_hash := make(map[string]Account)

	for _,acc := range accounts {
		if _,ok := M.MS.Clients[acc.Host]; !ok {
			/* Add if we can't find a session */
			missing_hash[acc.Host] = acc
		}
	}

	jsn,err := json.Marshal(missing_hash)
	if err != nil {
		return
	}

	goweb.Respond.With(c, 200, []byte(jsn))
}


func HTTP_Accounts_Add(M *Main, c context.Context) {
//	"/accounts/add/{hash}/{host}/{status}"
	
	hash := c.PathValue("hash")
	host := c.PathValue("host")
	status := c.PathValue("status")

	status_bool := false

	groups,err := HTTP_Filter(c, "groups")
	if err != nil {
		return
	}

	switch status {
		case "true" : status_bool = true
		case "false" : status_bool = false
		default: return
	}

	account,err := M.MG_Insert_Account(hash, host, groups, status_bool)
	if err != nil {
		return
	}

	jsn, err := json.Marshal(account)
	if err != nil {
		return
	}

	goweb.Respond.With(c, 200, []byte(jsn))
	
}
