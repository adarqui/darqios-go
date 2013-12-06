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

func (TD *Task_Data) MG8_Hash_Exists(Policy_Name string, Policy_Idx string, Policy_Actual string) (bool) {
	/* Check to see if a task was hashed => [name][idx][actual] */
	return false
}

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

	base_dir := fmt.Sprintf("%s/mitigate", TD.State.M.Startup_Config.Client.Base)

	for _, script := range TD.Policy.Mitigate {
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
			T2.Subject = fmt.Sprintf("DARQIOS: %s:%s:%s Failed to mitigate %s:%s", T2.Name, T2.Idx, T2.Type, T2.Actual, T2.Threshold)
			T2.Body = fmt.Sprintf("ERROR:%q\n",err)


			mon := new(MON)
			mon.Op = MON_REQ_TASK
			mon.Data = &T2
			TD.State.M.M <- mon
		}
	}

	return false
}

func (TD *Task_Data) MG8_Handle_Clear(T *Task) (bool) {
	return false
}
