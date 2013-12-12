package main

import (
	"errors"
	"strings"
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



type Query_Generic struct {
//	Map map[string]*Query_Datum
	Arr []Query_Datum
}

type Query_Datum struct {
	Host string
	Data interface{}
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

	query,err := M.MG_Query(collection, hosts, groups, ts_start, ts_end, limit, filter)
	if err != nil {
		return
	}

	/*
	 * FIXME clean this up, dirty - tired
	 */
	if filter != "" {
		var q []State_Report

		switch val := query.(type) {
			case []State_Report: {
				q = val
			}
		}

		query_filter, err := HTTP_Query_Filter(q,filter)
		if err != nil {
			return
		}

		jsn, err := json.Marshal(query_filter)
		if err != nil {
			return
		}
		goweb.Respond.With(c,200,[]byte(jsn))
	} else {
		jsn, err := json.Marshal(query)
		if err != nil {
			return
		}

		goweb.Respond.With(c,200,[]byte(jsn))
		}
	return
}



// FIXME - these functions are all pretty similar... need to generic it

func HTTP_Query_Filter(SP []State_Report, Filter string) (*Query_Generic, error) {
	paths := strings.Split(Filter, ":")


	switch paths[0] {
		case "LoadAvg": return HTTP_Query_Filter_LoadAvg(SP,paths)
		case "Memory": return HTTP_Query_Filter_Memory(SP,paths)
		case "Process": return HTTP_Query_Filter_Process(SP,paths)
		case "Network": return HTTP_Query_Filter_Network(SP,paths)
	}
	return nil, nil
}

func HTTP_Query_Filter_LoadAvg(SP []State_Report, Paths []string) (*Query_Generic, error) {

	if len(Paths) < 2 {
		return nil, errors.New("-EPARAM")
	}

	qg := new(Query_Generic)
	qg.Arr = make([]Query_Datum,0)

	for _, sp := range SP {
		qd := Query_Datum{}
		qd.Host = sp.Host
		switch Paths[1] {
			case "last1min": qd.Data = sp.LoadAvg.Last1Min
			case "last5min": qd.Data = sp.LoadAvg.Last5Min
			case "last15min": qd.Data = sp.LoadAvg.Last15Min
			default: continue
		}
		qg.Arr = append(qg.Arr, qd)
	}

	return qg, nil
}


func HTTP_Query_Filter_Memory(SP []State_Report, Paths []string) (*Query_Generic, error) {

	if len(Paths) < 2 {
		return nil, errors.New("-EPARAM")
	}

	qg := new(Query_Generic)
	qg.Arr = make([]Query_Datum,0)

	for _, sp := range SP {
		qd := Query_Datum{}
		qd.Host = sp.Host
		switch Paths[1] {
			case "Free": qd.Data = sp.Memory.Free
			default: continue
		}
		qg.Arr = append(qg.Arr, qd)
	}

	return qg, nil
}


func HTTP_Query_Filter_Process(SP []State_Report, Paths []string) (*Query_Generic, error) {
	if len(Paths) < 2 {
		return nil, errors.New("-EPARAM")
	}

	qg := new(Query_Generic)
	qg.Arr = make([]Query_Datum,0)

	for _, sp := range SP {
		qd := Query_Datum{}
		qd.Host = sp.Host
		switch Paths[1] {
			case "Total" : qd.Data = sp.Proc.Total
			default: continue
		}
		qg.Arr = append(qg.Arr, qd)
	}

	return qg, nil
}


func HTTP_Query_Filter_Network(SP []State_Report, Paths []string) (*Query_Generic, error) {
	if len(Paths) < 2 {
		return nil, errors.New("-EPARAM")
	}

	qg := new(Query_Generic)
	qg.Arr = make([]Query_Datum,0)

	for _, sp := range SP {
		qd := Query_Datum{}
		qd.Host = sp.Host
		switch Paths[1] {
			case "Connections": qd.Data = sp.Network.Connections
			default: continue
		}
		qg.Arr = append(qg.Arr, qd)
	}

	return qg, nil
}
