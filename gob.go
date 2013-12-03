package main

import (
	"encoding/gob"
)


func GOB_Init() {
	gob.Register(Policies_Config{})
	gob.Register(MON{})
	gob.Register(Task{})
	gob.Register(State_Report{})
}
