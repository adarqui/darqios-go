package main

import (
	"fmt"
	"os/exec"
)

/*
type Task struct {
	⇥   Type string
	⇥   Name string
	⇥   Idx string
	⇥   Subject string
	⇥   Body string
	⇥   Time time.Time
	⇥   Actual string
	⇥   Operator string
	⇥   Threshold string
}
*/


func (TD * Task_Data) MG8_Launch(T *Task) (bool) {

	switch T.Type {
		case "low","med","high" : return TD.MG8_Handle_Alert(T)
		case "clear" : return TD.MG8_Handle_Clear(T)
		default : {
			return false
		}
	}

	return false
}


func (TD *Task_Data) MG8_Handle_Alert(T *Task) (bool) {

	if TD.State.M.Startup_Config.Client.Base == "" {
		return false
	}

	base_dir := fmt.Sprintf("%s/mitigate", TD.State.M.Startup_Config.Client.Base)

	script := ""

	l := len(TD.Policy.Mitigate)
	switch l {
		case 0: return false
		case 1: script = TD.Policy.Mitigate[0]
		default: {
			alert_index,err := Task_Alert_2_Index(T.Type)
			if err != nil {
				return false
			}
			if l < alert_index {
				return false
			}
			script = TD.Policy.Mitigate[alert_index]
		}
	}


	script_path := fmt.Sprintf("%s/%s", base_dir, script)
	Debug("MG8_Handle_Alert: %q %q %q\n", script, base_dir, script_path)
	cmd := exec.Command(
		script_path,
		T.Type,
		T.Name,
		T.Idx,
		T.Actual,
		T.Threshold,
	)
	err := cmd.Run()
	if err != nil {
		Debug("MG8_Handle_Alert:cmd.Run:Err:%q\n",err)
		T2 := *T
		T2.Subject = fmt.Sprintf("%s:%s:%s Failed to mitigate %s:%s", T2.Name, T2.Idx, T2.Type, T2.Actual, T2.Threshold)
		T2.Body = fmt.Sprintf("ERROR:%q\n",err)


		mon := new(MON)
		mon.Op = MON_REQ_TASK
		mon.Data = &T2
		TD.State.M.M <- mon
	} else {
		mon := new(MON)
		mon.Op = MON_REQ_TASK
		T2 := *T
		T2.Subject = fmt.Sprintf("%s:%s:%s - %s = Success", T2.Name, T2.Idx, T2.Type, script_path)
		T2.Body = "None"
		mon.Data = &T2
		TD.State.M.M <- mon
	}

	return false
}

func (TD *Task_Data) MG8_Handle_Clear(T *Task) (bool) {
	return false
}
