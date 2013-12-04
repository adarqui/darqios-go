package main

import (
	"fmt"
	"encoding/json"
)

func Task_Load(M *Main, S *State, TD *Task_Data) (bool) {
	Debug("Task_Load:Entered\n")

	result, result_from_threshold := Task_Compare_Numbers(TD.Policy.Thresholds, S.LoadAvg.Last1Min, '>')

	fmt.Printf("result=%q result_from_threshold=%q\n", result, result_from_threshold)
	if result == "" {
		return false
	}

	jsn, err := json.Marshal(S.LoadAvg)
	if err != nil {
		return false
	}

	mon := MON_Gen_Task(result, TD.Policy, fmt.Sprintf("%f > %s", S.LoadAvg.Last1Min, result_from_threshold), string(jsn))

	M.M<-mon

	return true
}
