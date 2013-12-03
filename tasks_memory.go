package main

import (
	"fmt"
	"encoding/json"
)

func Task_Memory(M *Main, S *State, TD *Task_Data) (bool) {
	Debug("Task_Memory:Entered\n")

	return Task_Memory_Free(M,S,TD)
}

func Task_Memory_Free(M *Main, S *State, TD *Task_Data) (bool) {
	fmt.Printf("%v %d %d\n", S.Memory, S.Memory["MemFree"], S.Memory["MemTotal"])

	mem_free := S.Memory["MemFree"]
	mem_tot := S.Memory["MemTotal"]

	fmt.Printf("%d %d %f\n", mem_free, mem_tot, float64(mem_free/mem_tot))
	percent := float64(mem_free) / float64(mem_tot) * float64(100)

	result, result_from_threshold := Task_Compare_Numbers(TD.Policy.Thresholds, percent, '<')

	if result == "" {
		return false
	}


	jsn, err := json.Marshal(S.Memory)
	if err != nil {
		return false
	}

	mon := MON_Gen_Task(result, TD.Policy, fmt.Sprintf("%s:%s - Level (%s) : Free memory is %s", TD.Policy.Name, TD.Policy.Idx, result, result_from_threshold), string(jsn))

	M.M<-mon

	return true
}
