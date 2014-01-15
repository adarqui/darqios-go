package main

import (
	"fmt"
	"os"
)

func Task_File(M *Main, S *State, TD *Task_Data) (bool) {
	Debug("Task_File:Entered\n")

	switch TD.Policy.Idx {
		case "exists" : return Task_File_Exists(M, S, TD)
		case "!exists" : return Task_File_Not_Exists(M, S, TD)
		default: return false
	}
}

func Task_File_Generic(M *Main, S *State, TD *Task_Data, Exists bool) (bool) {

	truth := false
	alert_level := Task_Get_Alert_Level(TD.Policy)

	for _, param := range TD.Policy.Params {

		_, err := os.Stat(param)

		if Exists == true {
			if err == nil {
				/* Warn me if file exists */
				mon := S.MON_Gen_Task(alert_level, param, TD.Policy, fmt.Sprintf("%s exists", param), "None")
				M.M <- mon
				truth = true
			} else {
				S.STATE_Hash_All_Clear(TD, param)
			}
		} else if Exists == false {
			if err != nil {
				/* Warn me if file does not exists */
				mon := S.MON_Gen_Task(alert_level, param, TD.Policy, fmt.Sprintf("%s does not exist", param), "None")
				M.M <- mon
				truth = true
			} else {
				S.STATE_Hash_All_Clear(TD, param)
			}
		}

	}

	return truth
}


func Task_File_Not_Exists(M *Main, S *State, TD *Task_Data) (bool) {
	return Task_File_Generic(M, S, TD, false)
}


func Task_File_Exists(M *Main, S *State, TD *Task_Data) (bool) {
	return Task_File_Generic(M, S, TD, true)
}
