package main

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

func Task_Pipe(M *Main, S *State, TD *Task_Data) (bool) {
	Debug("Task_Pipe:Entered\n")

	switch TD.Policy.Idx {
		case "exists" : return Task_Pipe_Exists(M, S, TD)
		case "!exists" : return Task_Pipe_Not_Exists(M, S, TD)
		default: return false
	}
}

func Task_Pipe_Generic(M *Main, S *State, TD *Task_Data, Exists bool) (bool) {

	truth := false
	alert_level := Task_Get_Alert_Level(TD.Policy)

	params := TD.Policy.Params
	
	if len(params) != 2 {
		return false
	}

	command_argv := params[0]
	match_string := params[1]
	argv := strings.Split(command_argv, " ")

	cmd := exec.Command("/bin/ascp", "-A")
	cmd.Path = argv[0]
	cmd.Args = argv

	command_output, err := cmd.CombinedOutput()
	if err != nil {
		/* Warn if error executing the argv */
		mon := S.MON_Gen_Task(alert_level, command_argv, TD.Policy, fmt.Sprintf("Error executing %s: %s", command_argv, err), "None")
		M.M <- mon
		return false
	}

	re, err := regexp.CompilePOSIX(match_string)
	if err != nil {
		return false
	}
	matched := re.MatchString(string(command_output))

	if Exists == true {
		if matched == true {
			/* Warn me if file exists */
			mon := S.MON_Gen_Task(alert_level, command_argv, TD.Policy, fmt.Sprintf("Executing %s produced %s", command_argv, match_string), "None")
			M.M <- mon
			truth = true
		} else {
			S.STATE_Hash_All_Clear(TD, command_argv)
		}
	} else if Exists == false {
		if matched == false {
			/* Warn me if file does not exists */
			mon := S.MON_Gen_Task(alert_level, command_argv, TD.Policy, fmt.Sprintf("Executing %s did not produce %s", command_argv, match_string), "None")
			M.M <- mon
			truth = true
		} else {
			S.STATE_Hash_All_Clear(TD, command_argv)
		}
	}

	return truth
}


func Task_Pipe_Not_Exists(M *Main, S *State, TD *Task_Data) (bool) {
	return Task_Pipe_Generic(M, S, TD, false)
}


func Task_Pipe_Exists(M *Main, S *State, TD *Task_Data) (bool) {
	return Task_Pipe_Generic(M, S, TD, true)
}
