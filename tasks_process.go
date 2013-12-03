package main

import (
	"fmt"
)

func Task_Process(M *Main, S *State, TD *Task_Data) (bool) {
	Debug("Task_Process:Entered\n")

	switch TD.Policy.Idx {
		case "check" :{
			return Task_Process_Check(M, S, TD)
		}
		case "running" : {
			return Task_Process_Running(M, S, TD)
		}
		case "!running" : {
			return Task_Process_Not_Running(M, S, TD)
		}
		default : {
			return false
		}
	}
}


func Task_Process_Running(M *Main, S * State, TD *Task_Data) (bool) {

	var alert_level string

	if len(TD.Policy.Thresholds) == 0 {
		alert_level = "high"
	} else {
		alert_level = TD.Policy.Thresholds[0]
	}


	found := false

	for _, v := range S.Processes.Map {

		for _, w := range TD.Policy.Params {
			if v.Comm == w {

				found = true

				mon := MON_Gen_Task(alert_level, TD.Policy, fmt.Sprintf("%s:%s - Level (%s) : %s is running", TD.Policy.Name, TD.Policy.Idx, alert_level, w), "None.")

				M.M<-mon
			}
		}
	}

	return found
}




func Task_Process_Not_Running(M *Main, S *State, TD *Task_Data) (bool) {

	var alert_level string

	found := false

	if len(TD.Policy.Thresholds) == 0 {
		alert_level = "high"
	} else {
		alert_level = TD.Policy.Thresholds[0]
	}

	for _, v := range S.Processes.Map {
		for _, w := range TD.Policy.Params {
			if v.Comm == w {

				found = true

				mon := MON_Gen_Task(alert_level, TD.Policy, fmt.Sprintf("%s:%s - Level (%s) : %s is not running", TD.Policy.Name, TD.Policy.Idx, alert_level, w), "None.")

				M.M<-mon
			}
		}
	}

	return found
}



func Task_Process_Check(M *Main, S *State, TD *Task_Data) (bool) {
	return false
}
