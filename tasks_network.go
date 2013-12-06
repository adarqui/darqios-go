package main

import (
//	"log"
	"fmt"
	"strconv"
)

func Task_Network(M *Main, S *State, TD *Task_Data) (bool) {

	Debug("Task_Network:Entered\n")

	if TD.Policy.Idx == "tx" || TD.Policy.Idx == "rx" || TD.Policy.Idx == "bandwidth_any" {
		/* Bandwidth check */
		return Task_Network_Bandwidth_Check(M,S,TD)
	} else {
		/* Protocol/Port check */
		return Task_Network_Port_Check(M,S,TD)
	}
}


func Task_Network_Bandwidth_Check(M *Main, S *State, TD *Task_Data) (bool) {
	/*
	 * Check to see if bytes per second is being exceeded for rx, tx, or any
	 */
	return false
}

func Task_Network_Port_Check(M *Main, S *State, TD *Task_Data) (bool) {
	/*
	 * Check to see if a port is 'listening'
	 *
	 * Idx = protocol { tcp,tcp6,udp,udp6,port_any }
	 * Params = ports
	 */

	alert_level := Task_Get_Alert_Level(TD.Policy)

	var found bool
	for _, port_str := range TD.Policy.Params {
		found = false

		port, err := strconv.Atoi(port_str)
		if err != nil {
			continue
		}
		if TD.Policy.Idx == "port_any" {
			if _, ok := S.Network_Map.Sockets[port]; !ok {
				found = true
			}
		} else {
			if _, ok := S.Network_Map.Sockets[port][TD.Policy.Idx]; !ok {
				found = true
			}
		}
		if found == true {
			Task_Network_Port_Check_Alert(M,S,TD,alert_level, port)
		}
	}
	return false
}



func Task_Network_Port_Check_Alert(M *Main, S *State, TD *Task_Data, Alert_Level string, Port int) {
	mon := S.MON_Gen_Task(Alert_Level, string(Port), TD.Policy, fmt.Sprintf("%s:%d is not listening", TD.Policy.Idx, Port), "None.")
	M.M<-mon
}
