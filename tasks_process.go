package main

import (
	"fmt"
)

func Task_Process(M *Main, S *State, TD *Task_Data) (bool) {
	Debug("Task_Process:Entered\n")

	switch TD.Policy.Idx {
		case "check" : return Task_Process_Check(M, S, TD)
		case "running" : return Task_Process_Running(M, S, TD)
		case "!running" : return Task_Process_Not_Running(M, S, TD)
		default : return false
	}
}


func Task_Process_Running(M *Main, S * State, TD *Task_Data) (bool) {

	found := false

	alert_level := Tasks_Alert_Level_From_Idx(TD.Policy)

	for _, check_process := range TD.Policy.Params {
		found = false
		for _, actual_process := range S.Processes.Map {
			if check_process == actual_process.Comm {
				found = true
				break
			}
		}
		if found == true {
			mon := MON_Gen_Task(alert_level, TD.Policy, fmt.Sprintf("%s is running", check_process),  "None.")
			M.M<-mon
		}
	}

	return found
}




func Task_Process_Not_Running(M *Main, S *State, TD *Task_Data) (bool) {

	found := false

	alert_level := Tasks_Alert_Level_From_Idx(TD.Policy)

	for _, check_process := range TD.Policy.Params {
		found = false
		for _, actual_process := range S.Processes.Map {
			if check_process == actual_process.Comm {
				found = true
				break
			}
		}
		if found == false {
			mon := MON_Gen_Task(alert_level, TD.Policy, fmt.Sprintf("%s is not running", check_process), "None")
			M.M<-mon
		}
	}

	return found
}



func Task_Process_Check(M *Main, S *State, TD *Task_Data) (bool) {
	return false
}
