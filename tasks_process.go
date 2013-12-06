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
		case "!running_single": return Task_Process_Not_Running_Single(M, S, TD)
		default : return false
	}
}


func Task_Process_Running(M *Main, S * State, TD *Task_Data) (bool) {

	found := false

	alert_level := Task_Get_Alert_Level(TD.Policy)

	for _, check_process := range TD.Policy.Params {
		found = false
		for _, actual_process := range S.Processes.Map {
			if check_process == actual_process.Comm {
				found = true
				break
			}
		}
		if found == true {
			mon := S.MON_Gen_Task(alert_level, check_process, TD.Policy, fmt.Sprintf("%s is running", check_process),  "None.")
			M.M<-mon
		}
	}

	return found
}


func Task_Process_Not_Running(M * Main, S *State, TD *Task_Data) (bool) {
	return Task_Process_Not_Running_Generic(M, S, TD, false)
}


func Task_Process_Not_Running_Single(M * Main, S *State, TD *Task_Data) (bool) {
	return Task_Process_Not_Running_Generic(M, S, TD, true)
}


func Task_Process_Not_Running_Generic(M *Main, S *State, TD *Task_Data, Single bool) (bool) {

	found := false

	alert_level := Task_Get_Alert_Level(TD.Policy)

	for _, check_process := range TD.Policy.Params {
		found = false
		for _, actual_process := range S.Processes.Map {
			if check_process == actual_process.Comm {
				found = true
				break
			}
		}
		if found == false {
			mon := S.MON_Gen_Task(alert_level, check_process, TD.Policy, fmt.Sprintf("%s is not running", check_process), "None")
			M.M<-mon

			if Single == true {
				/* When set, only fire off one Gen_Task.. for example, when multiple processes being down can trigger one and only one restart script. */
				return found
			}
		} else {
			/*
			 * Here's where it gets tricky... We need to notify users of a cleared alert
			 */
			truth := S.STATE_Hash_Exists(TD.Policy.Name, TD.Policy.Idx, check_process)
			if truth == true {
				/* This means we've had an alert, that is now clear, so notify */
				S.STATE_Hash_Clear(TD.Policy.Name, TD.Policy.Idx, check_process, TD.Policy)
			}

		}
	}

	return found
}



func Task_Process_Check(M *Main, S *State, TD *Task_Data) (bool) {
	return false
}
