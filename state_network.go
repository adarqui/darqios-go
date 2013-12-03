package main

import (
//	"log"
	"time"
	"io/ioutil"
//	"time"
	"strings"
	"strconv"
	"os/exec"
)


func (S *State) STATE_Init_Network() {
	S.History_Interfaces = make(map[string]XInterface)
}

func (S *State) STATE_XInterface_Bandwidth(XI *XInterface) {
	if xif,ok := S.History_Interfaces[XI.Name]; ok {
		time_diff := XI.TS_Last.Sub(xif.TS_Last)
		quot := uint64(time_diff/time.Second)
		if quot != 0 {

			XI.Tx.Bandwidth = (XI.Tx.Bytes - xif.Tx.Bytes) / quot
			XI.Rx.Bandwidth = (XI.Rx.Bytes - xif.Rx.Bytes) / quot
		}
		/*
		S.History_Interfaces[XI.Name] = *XI
		Debug("%q:::%q\n", xif, XI)
	} else {
		*/
	}

		S.History_Interfaces[XI.Name] = *XI
}


func STATE_XInterfaces_2MAP(Interfaces *Interfaces) (*State_Report_Interfaces) {
	SRI := new(State_Report_Interfaces)
	SRI.Map = make(map[string]XInterface)
	for _,xif := range Interfaces.Interfaces {
		SRI.Map[xif.Name] = *xif
	}
	return SRI
}

func STATE_Get_Network() (*Network) {

	N := new(Network)
	N.Listeners = make([]*Network_Listener,0)

	netstat_connections := exec.Command("netstat","-4","-6","-n")
	output_connections, _ := netstat_connections.Output()

	connections := strings.Split(string(output_connections),"\n")
	N.Connections = len(connections)

	netstat_map := exec.Command("netstat","-nlp","-4","-6")
	output_map, _ := netstat_map.Output()

	for k,v := range strings.Split(string(output_map), "\n" ) {
		if k == 0 || len(v) == 0 {
			continue
		}

		fields := strings.Fields(v)

		if len(fields) < 6 {
			continue
		}

		proto := fields[0]
		listen := fields[3]

		idx := 0
		switch len(fields) {
			case 6: idx = 5
			case 7: idx = 6
			default: continue
		}
		pid_proc := fields[idx]

		fields_pid_proc := strings.Split(pid_proc,"/")
		if len(fields_pid_proc) != 2 {
			continue
		}

		pid := fields_pid_proc[0]
		proc := fields_pid_proc[1]

		NL := new(Network_Listener)
		NL.Pid = pid
		NL.Proto = proto
		NL.Process = proc
		NL.Listen = listen

		N.Listeners = append(N.Listeners,NL)

		Debug("Get_Network:NL:%q\n", NL)
	}

	Debug("---------------------------------------------------\n%q----------------------------------\n", N)
	return N
}


func (S *State) STATE_Get_Interfaces() (*Interfaces) {

	I := new(Interfaces)

	interfaces := S.STATE_Interfaces()

	I.Interfaces = interfaces

	Debug("STATE_Get_Interfaces:Interfaces:%q\n", I.Interfaces)

	return I
}


func (S *State) STATE_Interfaces() ([]*XInterface) {

	file,err := ioutil.ReadFile("/proc/net/dev")
	if err != nil {
		return nil
	}

	t:=time.Now()

	xif_array := make([]*XInterface,0)

	ifs := strings.Split(string(file), "\n")

	for k,v := range ifs {

		if k < 2 {
			continue
		}

		fields := strings.Fields(v)

		if len(fields) < 10 {
			continue
		}

//		Debug("FIELDS:%q\n", fields)

		xif := new(XInterface)
		xif.Name = fields[0]

		xif.Name = strings.Replace(xif.Name,":","",-1)

		xif.Rx.Bytes,_ = strconv.ParseUint(fields[1],10,64)
		xif.Rx.Packets,_ = strconv.ParseUint(fields[2],10,64)
		xif.Rx.Errors,_ = strconv.ParseUint(fields[3],10,64)
		xif.Rx.Drops,_ = strconv.ParseUint(fields[4],10,64)

		xif.Tx.Bytes,_ = strconv.ParseUint(fields[9],10,64)
		xif.Tx.Packets,_ = strconv.ParseUint(fields[10],10,64)
		xif.Tx.Errors,_ = strconv.ParseUint(fields[11],10,64)
		xif.Tx.Drops,_ = strconv.ParseUint(fields[12],10,64)

		xif.TS_Last = t

		xif_array = append(xif_array,xif)

		S.STATE_XInterface_Bandwidth(xif)

/*
interface bytes packets errs drop fifo frame compressed multicast bytes packets errors drop fifo colls carrier compressed

interface = 0
rx.bytes = 1
rx.packets = 2
rx.errors = 3
rx.drops = 4
tx.bytes = 9
tx.packets = 10
tx.errors = 11
tx.drops = 12

type XInterface_Stat struct {
	⇥   Bytes uint64
	⇥   Packets uint64
	⇥   Errors uint64
	⇥   Drops uint64
}*/
	}

	return xif_array
}
