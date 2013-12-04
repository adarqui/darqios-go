package main

import (
	"encoding/json"
	"github.com/stretchr/goweb"
	"github.com/stretchr/goweb/context"
)

func (M *Main) HTTP_Policies_Init() {
	M.HTTP_Policies_Routes()
}

func (M *Main) HTTP_Policies_Routes() {

	get_routes := []string{"/policies","/policies/","/policies/get","/policies/get/"}
	for _, route := range get_routes {
		goweb.Map("GET", route, func(c context.Context) (error) {
			HTTP_Policies_Index(M,c)
			return nil
		})
	}

	goweb.Map("GET", "/policies/get/{name}", func(c context.Context) (error) {
		HTTP_Policies_Get(M,c)
		return nil
	})


}

func HTTP_Policies_Get(M *Main, c context.Context) {

	name := c.PathValue("name")

	policies_config, err := M.MG_Lookup_Policies_Config_Raw(name)
	if err != nil {
		return
	}

	jsn,err := json.Marshal(policies_config)
	if err != nil {
		return
	}
	goweb.Respond.With(c,200,[]byte(jsn))
}


func HTTP_Policies_Index(M *Main, c context.Context) {

	names := make([]string,0)

	policies_config, err := M.MG_Lookup_Policies_Configs()
	if err != nil {
		return
	}

	for _,name := range policies_config {
		names = append(names, name.Name)
	}

	jsn, err := json.Marshal(names)
	if err != nil {
		return
	}
	goweb.Respond.With(c, 200, []byte(jsn))
}
