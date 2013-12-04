package main

import (
	"fmt"
	"encoding/json"
)

func Task_Disk(M *Main, S *State, TD *Task_Data) (bool) {
	Debug("Task_Disk:Entered\n")

	switch TD.Policy.Idx {
		case "free": {
			return Task_Disk_Free(M, S, TD)
		}
		case "lost": {
			return Task_Disk_Lost(M, S, TD)
		}
		default: {
			return false
		}
	}

}


func Task_Disk_Free(M *Main, S *State, TD *Task_Data) (bool) {

	Debug("Task_Disk_Free:Entered\n")

	found := false

	for _, disk_name := range TD.Policy.Params {
		if disk, ok := S.Disks.Map[disk_name]; ok {

			percent := float64(disk.Used) / float64(disk.Size) * float64(100)

			result, result_from_threshold := Task_Compare_Numbers(TD.Policy.Thresholds, percent, '>')

			if result == "" {
				continue
			}

			jsn, err := json.Marshal(disk)
			if err != nil {
				continue
			}

			mon := MON_Gen_Task(result, TD.Policy, fmt.Sprintf("%s (%f%%) exceeds %s%% capacity", disk.Mount, percent, result_from_threshold), string(jsn))

			M.M<-mon

			found = true
		}
	}
	return found
}



func Task_Disk_Lost(M *Main, S *State, TD *Task_Data) (bool) {
	return false
}
