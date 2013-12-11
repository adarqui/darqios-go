package main

import (
	"fmt"
	"time"
	"strings"
	"strconv"
	"os/exec"
	"github.com/c9s/goprocinfo/linux"
)

type State_Report struct {
	Host string
	LoadAvg *linux.LoadAvg
	Disks *Disks
	Memory State_Report_Memory
	Proc State_Report_Proc
	Interfaces *State_Report_Interfaces
	Network Network
//	Network State_Report_Network
	Users State_Report_Users
	Ts time.Time
}

type State struct {
	M * Main
	Processes *Processes
	LoadAvg *linux.LoadAvg
	Disks *Disks
	Memory linux.MemInfo
	Uptime *linux.Uptime
	Users *Users
	Interval time.Duration
	Interval_Counter int
	Network *Network
	Network_Map *Network_Map
	Interfaces *Interfaces
	History_Interfaces map[string]XInterface
	History_Disks map[string]Disk
	/* [Policy.Name][Policy.Idx]Threshold */
	Alerts_Hash map[string]bool
}

type State_Report_Memory struct {
	Total float64
	Free float64
}

type State_Report_Proc struct {
	Total int
	ByMem map[string]float64
	ByCpu map[string]float64
}

/*
 * For sorting
 */
type Proc_ByCpu []Process
type Proc_ByMem []Process



type XInterface struct {
	Name string
	Tx XInterface_Stat
	Rx XInterface_Stat
	TS_Last time.Time
}

type XInterface_Stat struct {
	Bytes uint64
	Packets uint64
	Errors uint64
	Drops uint64
	Bandwidth uint64
}

type Interfaces struct {
//	Connections int
	Interfaces []*XInterface
}

type State_Report_Interfaces struct {
	Map map[string]XInterface
}

/*
type State_Report_Network struct {
	Map map[string]Network_Listener
}
*/

type State_Report_Users struct {
	Total int
}


type Network_Listener struct {
	Pid string
	Proto string
	Process string
	Listen string
}

/*
type State_Report_Network struct {
	Connections int
	Map []byte
}
*/

type Network struct {
	Connections int
	Listeners []*Network_Listener
}

type Network_Map struct {
	Sockets map[int]map[string]bool
}

/*
type Interface struct {
	Name string
	Bandwidth int
	Errors int
	Tx uint64
	Rx uint64
}
*/

type Process struct {
	Pid int
	Comm string
	Ppid int
	Pcpu float64
	Pmem float64
}

type Processes struct {
	Map map[int]Process
	Array []Process
}


type Disk struct {
	Name string
	Size uint64
	Used uint64
	Avail uint64
	Mount string
	Bandwidth uint64
	TS_Last time.Time
}

type Disks struct {
	Map map[string]Disk
}


type User struct {
	Name string
	Count int
}

type Users struct {
	Map map[string]User
}


func (M *Main) STATE_Init() (*State) {
	S := new(State)

	/* dirty, tired of passing around M.. sad */
	S.M = M
//	S.Interval = 1

/*	sleep := M.STATE_Init_Interval()
	S.Interval = time.Duration(sleep)
	*/
	S.Interval = 1*time.Second
	S.STATE_Init_Disks()
	S.STATE_Init_Network()

	S.STATE_Hash_Init()

	return S
}

func (M *Main) STATE_Init_Interval() (int) {
	for _, policy := range M.Policies_Config.Policies {
		if policy.Name == "State" && policy.Active == true {
			return policy.Interval
		}
	}
	return 1
}

func (S *State) STATE_Sleep() {
	time.Sleep(S.Interval)
	S.Interval_Counter++
}

func (S *State) STATE_Should_Run(P *Policy) (bool) {
	if P.Active != true {
		return false
	}

	if (S.Interval_Counter % P.Interval) != 0 {
		return false
	}

	return true
}

func (S *State) STATE_Get() {
	DebugLn("STATE_Get:Initialize")
	S.Processes = STATE_Get_Processes()
	S.Disks = S.STATE_Get_Disks()
	S.Users = STATE_Get_Users()
	S.LoadAvg, _ = STATE_Get_LoadAvg()
	S.Memory, _ = STATE_Get_Memory()
	S.Network = S.STATE_Get_Network()
	S.Interfaces = S.STATE_Get_Interfaces()
}



func (p Proc_ByCpu) Len() (int) {
	return len(p)
}

func (p Proc_ByCpu) Swap(a, b int) {
	p[a], p[b] = p[b], p[a]
}

func (p Proc_ByCpu) Less(a, b int) (bool) {
	return p[a].Pcpu > p[b].Pcpu
}


func (p Proc_ByMem) Len() (int) {
	return len(p)
}

func (p Proc_ByMem) Swap(a, b int) {
	p[a], p[b] = p[b], p[a]
}

func (p Proc_ByMem) Less(a, b int) (bool) {
	return p[a].Pmem > p[b].Pmem
}


func STATE_Get_Processes() (*Processes) {

	var err error
	P := new(Processes)

	P.Map = make(map[int]Process)
	/* Not sure what to do here, we need an array for sorting */
	P.Array = make([]Process, 0)

	ps := exec.Command("ps", "-e", "-opid,comm,ppid,pcpu,pmem")
	output, _ := ps.Output()
	for k, v := range strings.Split(string(output), "\n") {
		if k == 0 || len(v) == 0 {
			continue
		}
		fields := strings.Fields(v)
		proc := new(Process)
		proc.Pid, err = strconv.Atoi(fields[0])
		if err != nil {
			continue
		}
		proc.Comm = fields[1]
		proc.Ppid, err = strconv.Atoi(fields[2])
		if err != nil {
			continue
		}

		proc.Pcpu, err = strconv.ParseFloat(fields[3],64)
		if err != nil {
			continue
		}
		proc.Pmem, err = strconv.ParseFloat(fields[4],64)
		if err != nil {
			continue
		}
		P.Map[proc.Pid] = *proc


		P.Array = append(P.Array, *proc)
	}

	return P
}



func STATE_Get_LoadAvg() (*linux.LoadAvg,error) {
	load, err := linux.ReadLoadAvg("/proc/loadavg")
	return load, err
}

func STATE_Get_Memory() (linux.MemInfo,error) {
	mem, err := linux.ReadMemInfo("/proc/meminfo")
	return mem, err
}


func STATE_Get_Users() (*Users) {
	U := new(Users)
	U.Map = make(map[string]User)

	users := exec.Command("users")
	output, _ := users.Output()

	for k,v := range strings.Split(string(output), " ") {
		if k == 0 || len(v) == 0 {
			continue
		}

//		log.Printf("K=%q V=%q\n", k, v)

		if _,ok := U.Map[v]; !ok {
			user := new(User)
			user.Name = v
			user.Count = 1
			U.Map[v] = *user
		} else {
			val := U.Map[v]
			val.Count += 1
			U.Map[v] = val
		}
		
	}

	return U
}


func (S *State) STATE_Hash_Init() {
	S.Alerts_Hash = make(map[string]bool)
}

/* GENERIC FUNCTION */
func (S *State) STATE_Hash_Exists(Policy_Name string, Policy_Idx string, Policy_Actual string) (bool) {

	/* Check to see if a task was hashed => [name][idx][actual] */
	hash_index := fmt.Sprintf("%s:%s:%s", Policy_Name, Policy_Idx, Policy_Actual)
	if _, ok := S.Alerts_Hash[hash_index]; ok {
		return true
	}
	return false
}

func (S *State) STATE_Hash_Clear(Policy_Name string, Policy_Idx string, Policy_Actual string, P *Policy) (bool) {
	hash_index := fmt.Sprintf("%s:%s:%s", Policy_Name, Policy_Idx, Policy_Actual)
	if _,ok := S.Alerts_Hash[hash_index]; ok {
		delete(S.Alerts_Hash, hash_index)

		mon := new(MON)
		mon.Op = MON_REQ_TASK
		task := MON_Gen_Task_Raw("clear", Policy_Actual, P, fmt.Sprintf("%s is no longer an issue", Policy_Actual), "None")
		mon.Data = task
		S.M.M<-mon
	}
	return false
}


func (S *State) STATE_Hash_Add(Policy_Name string, Policy_Idx string, Policy_Actual string) (bool) {

	hash_index := fmt.Sprintf("%s:%s:%s", Policy_Name, Policy_Idx, Policy_Actual)
	if _, ok := S.Alerts_Hash[hash_index]; !ok {
		S.Alerts_Hash[hash_index] = true
	}

	return false
}


func (S *State) STATE_Hash_All_Clear(TD *Task_Data, Actual string) {
	truth := S.STATE_Hash_Exists(TD.Policy.Name, TD.Policy.Idx, Actual)
	if truth == true {
		S.STATE_Hash_Clear(TD.Policy.Name, TD.Policy.Idx, Actual, TD.Policy)
	}
}
