package main

import (
	"time"
	"fmt"
	"sort"
)

const (
	MON_REQ_TASK = 1
	MON_REQ_STATE = 2
)

type MON struct {
	Op int
	Data interface{}
}


/*
 * This generates a "MON" Object which is sent up to the monitor channel for relay to the server
 */
func (S *State) MON_Gen_Task(Type string, Actual string, Policy *Policy, Subject string, Body string) (*MON) {
	mon := new(MON)

	mon.Op = MON_REQ_TASK
	task := MON_Gen_Task_Raw(Type, Actual, Policy, Subject, Body)
	mon.Data = task

	Debug("MON_Gen_Task:MON_REQ_TASK:%q\n", mon)

	/* Where to put this? */

	S.STATE_Hash_Add(Policy.Name, Policy.Idx, Actual)

	/*
	 * Mitigation entry point
	 */
	TD := new(Task_Data)
	TD.State = S
	TD.Policy = Policy
	go TD.MG8_Launch(task)

	return mon
}

func MON_Gen_Task_Raw(Type string, Actual string, Policy *Policy, Subject string, Body string) (*Task) {

	task := new(Task)
	task.Type = Type
	task.Name = Policy.Name
	task.Idx = Policy.Idx
	task.Actual = Actual
	task.Subject = fmt.Sprintf("DARQIOS: %s:%s - (%s) - %s", Policy.Name, Policy.Idx, Type, Subject)
	task.Body = Body
	task.Time = time.Now()

	return task
}


func MON_Gen_State(State *State) (*MON) {

	mon := new(MON)

	mon.Op = MON_REQ_STATE
	mon.Data = MON_Gen_State_Raw(State)

	Debug("MON_Gen_State:MON_REQ_STATE:%q\n",mon)

	return mon
}


func MON_Gen_State_Raw(State *State) (*State_Report) {

	state_report := new(State_Report)

	state_report.LoadAvg = State.LoadAvg
	state_report.Proc.Total = len(State.Processes.Map)

	for _,user := range State.Users.Map {
		state_report.Users.Total += user.Count
	}

//	state_report.Memory = State.Memory
	mem_free := State.Memory["MemFree"]

	state_report.Memory.Total = float64(State.Memory["MemTotal"])
	state_report.Memory.Free = float64(mem_free) / float64(state_report.Memory.Total) * float64(100)

	state_report.Network = *State.Network
	Debug("%q\n",state_report.Network)
//	log.Fatal(state_report.Network)
	state_report.Interfaces = STATE_XInterfaces_2MAP(State.Interfaces)
	state_report.Disks = State.Disks

//	state_report.Network.Connections = State.Network.Connections
//	state_report.Network.Map,_ = bson.Marshal(State.Network.Map)
//	state_report.Interfaces = State.Interfaces

	state_report.Ts = time.Now()


	/*
	 * Generate two map's of processes sorted by : Highest cpu, Highest mem
	 */

	state_report.Proc.ByCpu = make(map[string]float64)
	state_report.Proc.ByMem = make(map[string]float64)

	j := 0
	sort.Sort(Proc_ByCpu(State.Processes.Array))
	for _,elm := range State.Processes.Array {
		if val,ok := state_report.Proc.ByCpu[elm.Comm]; ok {
			val = val + elm.Pcpu
			state_report.Proc.ByCpu[elm.Comm] = val
		} else {
			j += 1
			state_report.Proc.ByCpu[elm.Comm] = elm.Pcpu
		}

		if j > 5 {
			break
		}
	}


	j = 0
	sort.Sort(Proc_ByMem(State.Processes.Array))
	for _, elm := range State.Processes.Array {
		if val,ok := state_report.Proc.ByMem[elm.Comm]; ok {
			val = val + elm.Pmem
			state_report.Proc.ByMem[elm.Comm] = val
		} else {
			j += 1
			state_report.Proc.ByMem[elm.Comm] = elm.Pmem
		}

		if j > 5 {
			break
		}
	}


	return state_report
}



func MON_Gen_Truthy(P *Policy) ([]bool, error) {

	truthy := make([]bool, len(P.Params))

	return truthy, nil
}
