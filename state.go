package main

import (
//	"log"
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
	Processes *Processes
	LoadAvg *linux.LoadAvg
	Disks *Disks
	Memory linux.MemInfo
	Uptime *linux.Uptime
	Users *Users
	Interval time.Duration
	Interval_Counter int
	Network *Network
	Interfaces *Interfaces
	History_Interfaces map[string]XInterface
	History_Disks map[string]Disk
}

type State_Report_Memory struct {
	Total float64
	Free float64
}

type State_Report_Proc struct {
	Total int
	Running int
	Stopped int
	Zombie int
}

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


func STATE_Init() (*State) {
	S := new(State)
	S.Interval = 1

	S.STATE_Init_Disks()
	S.STATE_Init_Network()

	return S
}

func (S *State) STATE_Sleep() {
	time.Sleep(S.Interval * time.Second)
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
	S.Network = STATE_Get_Network()
	S.Interfaces = S.STATE_Get_Interfaces()
}



func STATE_Get_Processes() (*Processes) {

	var err error
	P := new(Processes)

	P.Map = make(map[int]Process)

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
//		log.Printf("%q %q\n", fields, k, v)
//		log.Printf("%q\n", proc)
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
