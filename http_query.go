package main

import (
	"encoding/json"
	"github.com/stretchr/goweb"
	"github.com/stretchr/goweb/context"
)

func (M *Main) HTTP_Query_Init() {
	M.HTTP_Query_Routes()
}

func (M *Main) HTTP_Query_Routes() {

	goweb.Map("GET", "/query/{collection}/{hosts}/{groups}/{ts_start}/{ts_end}/{limit}/{filter}", func(c context.Context) (error) {
		DebugLn("QUERY")
		HTTP_Query(M,c)
		return nil
	})

}

func HTTP_Query(M *Main, c context.Context) {

	collection := c.PathValue("collection")

	hosts, err := HTTP_Filter(c, "hosts")
	if err != nil {
		return
	}

	groups, err := HTTP_Filter(c, "groups")
	if err != nil {
		return
	}

	ts_start := c.PathValue("ts_start")
	ts_end := c.PathValue("ts_end")
	limit := c.PathValue("limit")
	filter := c.PathValue("filter")

	Debug("%q %q %q %q %q %q\n", hosts, groups, ts_start, ts_end, limit, filter)

	query,err := M.MG_Query(collection, hosts, groups, ts_start, ts_end, limit, filter)
	if err != nil {
		return
	}

	Debug("%q\n", query)

	jsn, err := json.Marshal(query)
	if err != nil {
		return
	}

	goweb.Respond.With(c,200,[]byte(jsn))
}
