/*
 * This file needs to be refactored for sure...
 */

package main

import (
	"errors"
	"strings"
	"encoding/base64"
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

	if strings.HasPrefix(filter, "b64=") == true {
		decoded, err := base64.StdEncoding.DecodeString(filter[4:len(filter)])
		if err != nil {
			return
		}
		filter = string(decoded)
	}

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
		case "Users": return HTTP_Query_Filter_Users(SP,paths)
		case "Interfaces" : return HTTP_Query_Filter_Interfaces(SP,paths)
		case "Disks" : return HTTP_Query_Filter_Disk(SP,paths)
	}
	
	return nil, nil
}



func HTTP_Query_Filter_Disk(SP []State_Report, Paths []string) (*Query_Generic, error) {
	/*
	 * Disks:<name>:<field>
	 * Disks:/:Avail
	 * Disks:/:Bandwidth
	 */

	 if len(Paths) < 3 {
		 return nil, errors.New("-EPARAM")
	 }

	 qg := new(Query_Generic)
	 qg.Arr = make([]Query_Datum, 0)

	 disk_name := Paths[1]
	 disk_field := Paths[2]

	 for _, sp := range SP {
		 if disk, ok := sp.Disks.Map[disk_name]; ok {
			 qd := Query_Datum{}
			 qd.Host = sp.Host

			 ret, err := disk.PARSE_Disk_Req_Into_Datum(disk_field)
			 if err != nil {
				 continue
			 }

			qd.Data = ret
			qg.Arr = append(qg.Arr, qd)
		 }
	 }

	 return qg, nil
}


func (DISK *Disk) PARSE_Disk_Req_Into_Datum(DISK_FIELD string) (interface{}, error) {

	switch DISK_FIELD {
		case "Size" : {
			return DISK.Size, nil
		}
		case "Used" : {
			return DISK.Used, nil
		}
		case "Avail" : {
			return DISK.Avail, nil
		}
		case "Bandwidth" : {
			return DISK.Bandwidth, nil
		}
		case "AvailP" : {
			return DISK.AvailP, nil
		}
		case "UsedP" : {
			return DISK.UsedP, nil
		}
		default : {
			return nil, errors.New("-EPARAM")
		}
	}
}



func HTTP_Query_Filter_Interfaces(SP []State_Report, Paths []string) (*Query_Generic, error) {

	/*
	 * Interfaces:<if>:{TX,RX,BOTH}:<field>
	 * Interfaces:eth0:TX:Bandwidth
	 * Interfaces:eth0:BOTH:Bandwidth
	 */

	if len(Paths) < 4 {
		return nil, errors.New("-EPARAM")
	}

	qg := new(Query_Generic)
	qg.Arr = make([]Query_Datum,0)

	xif_name := Paths[1]
	xif_txrx := Paths[2]
	xif_field := Paths[3]

	for _, sp := range SP {

		if xif, ok := sp.Interfaces.Map[xif_name]; ok {

			qd := Query_Datum{}
			qd.Host = sp.Host

			ret, err := xif.PARSE_Interface_Req_Into_Datum(xif_txrx, xif_field)
			if err != nil {
				continue
			}

			qd.Data = ret
			qg.Arr = append(qg.Arr, qd)
		}
	}

	return qg, nil
}


func (XIF *XInterface) PARSE_Interface_Req_Into_Datum(XIF_TXRX string, XIF_FIELD string) (uint64, error) {
	switch(XIF_FIELD) {
		case "Bytes" : {
			return XIF.PARSE_Interface_Bytes(XIF_TXRX)
		}
		case  "Packets" : {
			return XIF.PARSE_Interface_Packets(XIF_TXRX)
		}
		case "Errors" : {
			return XIF.PARSE_Interface_Errors(XIF_TXRX)
		}
		case "Drops" : {
			return XIF.PARSE_Interface_Drops(XIF_TXRX)
		}
		case "Bandwidth" : {
			return XIF.PARSE_Interface_Bandwidth(XIF_TXRX)
		}
		default: {
			return 0, errors.New("-EPARAM")
		}
	}

}

func (XIF *XInterface) PARSE_Interface_Bytes(XIF_TXRX string) (uint64, error) {
	switch XIF_TXRX {
		case "TX": {
			return XIF.Tx.Bytes, nil
		}
		case "RX": {
			return XIF.Rx.Bytes, nil
		}
		case "BOTH" : {
			return XIF.Tx.Bytes + XIF.Rx.Bytes, nil
		}
		default: {
			return 0, errors.New("-EPARAM")
		}
	}
}

func (XIF *XInterface) PARSE_Interface_Packets(XIF_TXRX string) (uint64, error) {
	switch XIF_TXRX {
		case "TX" : {
			return XIF.Tx.Packets, nil
		}
		case "RX" : {
			return XIF.Rx.Packets, nil
		}
		case "BOTH" : {
			return XIF.Tx.Packets + XIF.Rx.Packets, nil
		}
		default : {
			return 0, errors.New("-EPARAM")
		}
	}
}


func (XIF *XInterface) PARSE_Interface_Errors(XIF_TXRX string) (uint64, error) {
	switch XIF_TXRX {
		case "TX" : {
			return XIF.Tx.Errors, nil
		}
		case "RX" : {
			return XIF.Rx.Errors, nil
		}
		case "BOTH" : {
			return XIF.Tx.Errors + XIF.Rx.Errors, nil
		}
		default : {
			return 0, errors.New("-EPARAM")
		}
	}
}


func (XIF *XInterface) PARSE_Interface_Drops(XIF_TXRX string) (uint64, error) {
	switch XIF_TXRX {
		case "TX" : {
			return XIF.Tx.Drops, nil
		}
		case "RX" : {
			return XIF.Rx.Drops, nil
		}
		case "BOTH" : {
			return XIF.Tx.Drops + XIF.Rx.Drops, nil
		}
		default : {
			return 0, errors.New("-EPARAM")
		}
	}
}



func (XIF *XInterface) PARSE_Interface_Bandwidth(XIF_TXRX string) (uint64, error) {
	switch XIF_TXRX {
		case "TX" : {
			return XIF.Tx.Bandwidth, nil
		}
		case "RX" : {
			return XIF.Rx.Bandwidth, nil
		}
		case "BOTH" : {
			return XIF.Tx.Bandwidth + XIF.Rx.Bandwidth, nil
		}
		default : {
			return 0, errors.New("-EPARAM")
		}
	}
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


func HTTP_Query_Filter_Users(SP []State_Report, Paths []string) (*Query_Generic, error) {
	if len(Paths) < 2 {
		return nil, errors.New("-EPARAM")
	}

	qg := new(Query_Generic)
	qg.Arr = make([]Query_Datum,0)

	for _, sp := range SP {
		qd := Query_Datum{}
		qd.Host = sp.Host
		switch Paths[1] {
			case "Total": qd.Data = sp.Users.Total
			default: continue
		}
		qg.Arr = append(qg.Arr, qd)
	}
	return qg, nil
}
