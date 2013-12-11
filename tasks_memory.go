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

	mem_free := S.Memory["MemFree"]
	mem_tot := S.Memory["MemTotal"]

	percent := float64(mem_free) / float64(mem_tot) * float64(100)

	result, result_from_threshold := Task_Compare_Numbers(TD.Policy.Thresholds, percent, '<')

	if result == "" {

		S.STATE_Hash_All_Clear(TD, "free")

		return false
	}


	jsn, err := json.Marshal(S.Memory)
	if err != nil {
		return false
	}

	mon := S.MON_Gen_Task(result, "free", TD.Policy, fmt.Sprintf("Free memory (%f%%) is below %s%%", percent, result_from_threshold), string(jsn))

	M.M<-mon

	return true
}
