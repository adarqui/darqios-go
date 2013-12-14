package main

import (
//"log"
	"fmt"
	"os/exec"
	"strings"
)

func Task_Process(M *Main, S *State, TD *Task_Data) (bool) {
	Debug("Task_Process:Entered\n")

	switch TD.Policy.Idx {
		case "check" : return Task_Process_Check(M, S, TD)
		case "running" : return Task_Process_Running(M, S, TD)
		case "!running" : return Task_Process_Not_Running(M, S, TD)
		case "!running_single": return Task_Process_Not_Running_Single(M, S, TD)
		case "!running_args" : return Task_Process_Not_Running_Args(M, S, TD)
		case "running_args" : return Task_Process_Running_Args(M, S, TD)
		case "!running_pgrep" : return Task_Process_Not_Running_Pgrep(M, S, TD)
		case "running_pgrep" : return Task_Process_Running_Pgrep(M, S, TD)
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
		} else {
			/*
			truth := S.STATE_Hash_Exists(TD.Policy.Name, TD.Policy.Idx, check_process)
			if truth == true {
				S.STATE_Hash_Clear(TD.Policy.Name, TD.Policy.Idx, check_process, TD.Policy)
			}
			*/

			S.STATE_Hash_All_Clear(TD, check_process)
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
			S.STATE_Hash_All_Clear(TD, check_process)

		}
	}

	return found
}


func Task_Process_Args_Generic(M *Main, S *State, TD *Task_Data, Running bool) (bool) {
	truth := false

	alert_level := Task_Get_Alert_Level(TD.Policy)
	/* process:count */
	Map := make(map[string]int)

	ps := exec.Command("ps", "-e", "-o", "args")
	output, _ := ps.Output()

	for k, v := range strings.Split(string(output), "\n") {
		if k == 0 || len(v) == 0 {
			continue
		}

		if _,ok := Map[v]; ok {
			Map[v] += 1
		} else {
			Map[v] = 1
		}
	}

	for _, param := range TD.Policy.Params {
		if val, ok := Map[param]; ok {

			if Running == true {
				/* process found and running true means, notify when this process is running */
				mon := S.MON_Gen_Task(alert_level, param, TD.Policy, fmt.Sprintf("%s is running (count=%d)", param, val), "None")
				M.M<-mon
				truth = true
			} else {
				S.STATE_Hash_All_Clear(TD, param)
			}
		} else {
			if Running == false {
				/* process not found, and running false means, notify when this process is not running */
				mon := S.MON_Gen_Task(alert_level, param, TD.Policy, fmt.Sprintf("%s is not running", param), "None")
				M.M<-mon
				truth = true
			} else {
				S.STATE_Hash_All_Clear(TD, param)
			}
		}
	}

	return truth
}

func Task_Process_Running_Args(M *Main, S *State, TD *Task_Data) (bool) {
	return Task_Process_Args_Generic(M,S,TD,true)
}

func Task_Process_Not_Running_Args(M *Main, S *State, TD *Task_Data) (bool) {
	return Task_Process_Args_Generic(M,S,TD,false)
}



func Task_Process_Pgrep_Generic(M *Main, S *State, TD *Task_Data, Running bool) (bool) {
	truth := false

	alert_level := Task_Get_Alert_Level(TD.Policy)
	/* process:count */
	for _, param := range TD.Policy.Params {

		ps := exec.Command("pgrep", "-f", param)
		output, _ := ps.Output()

		if len(output) > 0 {
			if Running == true {
				/* processes found and running = true, bad */
				mon := S.MON_Gen_Task(alert_level, param, TD.Policy, fmt.Sprintf("%s is running", param), "None")
				M.M<-mon
				truth = true
			} else {
				S.STATE_Hash_All_Clear(TD, param)
			}
		} else {
			if Running == false {
				/* process not found and running = false, bad */
				mon := S.MON_Gen_Task(alert_level, param, TD.Policy, fmt.Sprintf("%s is running", param), "None")
				M.M<-mon
				truth = true
			} else {
				S.STATE_Hash_All_Clear(TD, param)
			}
		}

	}

	return truth
}

func Task_Process_Running_Pgrep(M *Main, S *State, TD *Task_Data) (bool) {
	return Task_Process_Pgrep_Generic(M,S,TD,true)
}

func Task_Process_Not_Running_Pgrep(M *Main, S *State, TD *Task_Data) (bool) {
	return Task_Process_Pgrep_Generic(M,S,TD,false)
}




func Task_Process_Check(M *Main, S *State, TD *Task_Data) (bool) {
	return false
}
