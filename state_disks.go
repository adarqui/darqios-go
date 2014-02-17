package main

import (
//	"log"
	"time"
	"strings"
	"strconv"
	"os/exec"
//	"github.com/c9s/goprocinfo/linux"
)

func (S *State) STATE_Init_Disks() {
	S.History_Disks = make(map[string]Disk)
}

func (S * State) STATE_Disks_Bandwidth(D *Disk) {
	if disk,ok := S.History_Disks[D.Mount]; ok {
		time_diff := D.TS_Last.Sub(disk.TS_Last)
		quot := uint64(time_diff/time.Second)
		if quot != 0 {
			if D.Used >= disk.Used {
				D.Bandwidth = int64((D.Used - disk.Used) / quot)
			} else {
				D.Bandwidth = int64(((disk.Used - D.Used) / quot))*(-1)
			}
		}

	}

	S.History_Disks[D.Mount] = *D
}

func (S *State) STATE_Get_Disks() (*Disks) {

	/* name size used avail mount */
	var err error
	D := new(Disks)

	D.Map = make(map[string]Disk)

	df := exec.Command("df", "-a")
	output, _ := df.Output()

	t := time.Now()

	for k,v := range strings.Split(string(output), "\n") {
		if k == 0 || len(v) == 0 {
			continue
		}
		fields := strings.Fields(v)
		if len(fields) < 5 {
			continue
		}

//		log.Printf("FIELDS: %q %v %v\n", fields, fields[0], fields[1])
		disk := new(Disk)
		disk.Name = fields[0]
		disk.Size, err = strconv.ParseUint(fields[1],10,64)
		if err != nil {
			continue
		}

		if disk.Size == 0 {
			continue
		}

		disk.Used, err = strconv.ParseUint(fields[2],10,64)
		if err != nil {
			continue
		}
		disk.Avail, err = strconv.ParseUint(fields[3],10,64)
		if err != nil {
			continue
		}

		if disk.Size != 0 {
			disk.AvailP = float64(disk.Avail) / float64(disk.Size) * float64(100)
			disk.UsedP = float64(disk.Used) / float64(disk.Size) * float64(100)
		}

		disk.Mount = fields[5]
		disk.TS_Last = t

		S.STATE_Disks_Bandwidth(disk)

		D.Map[disk.Mount] = *disk
	}

	return D
}
