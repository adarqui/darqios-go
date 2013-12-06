package main

import (
	"fmt"
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

func (TD *Task_Data) MG8_Hash_Exists(T *Task) (bool) {
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
	}

	return false
}

func (TD *Task_Data) MG8_Handle_Clear(T *Task) (bool) {
	return false
}
