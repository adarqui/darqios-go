package main

import (
	"bytes"
	"encoding/gob"
)

const (
	WOP_REQ_NOWOP = 0
	WOP_REQ_PING = 1
	WOP_REQ_PONG = 2
	WOP_REQ_POLICIES_CONFIG = 3
	WOP_REQ_TASK = 4
	WOP_REQ_STATE = 5
)

const (
	WOP_REP_NOWOP = 0
	WOP_REP_OK = 1
	WOP_REP_NOTOK = 2
	WOP_REP_DIE = 3
	WOP_REP_POLICIES_CONFIG = 4
)

/*
 * WOPCODE
 */
type WOP struct {
	Code int
//	Arg string
	Data interface{}
}

func Bytes_To_WOP(Bytes []byte) (*WOP, error) {
	wop := new(WOP)
	buf := bytes.NewBuffer(Bytes)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&wop)
	if err != nil {
		Debug("Bytes_To_WOP:Err:%q\n",err)
		return nil, err
	}

	return wop, nil
}


func WOP_To_Bytes(wop *WOP) ([]byte,error) {
	var Bytes bytes.Buffer
	enc := gob.NewEncoder(&Bytes)
	err := enc.Encode(wop)
	if err != nil {
		Debug("WOP_To_Bytes:Err:%q\n", err)
		return nil, err
	}
	
	raw := Bytes.Bytes()
	return raw, nil
}


func WOP_Parse(Bytes []byte) (*WOP,error) {
	wop, err := Bytes_To_WOP(Bytes)
	if err != nil {
		return nil, err
	}

	return wop, nil
}

/*
func WOP_Gen_Ping() (*WOP) {
	wop := new(WOP)
	wop.Code = WOP_REQ_PING
}
*/


func WOP_Gen_Rep_Policies_Config(Policies_Config *Policies_Config)(*WOP){
	wop := new(WOP)
	wop.Code = WOP_REP_POLICIES_CONFIG
	wop.Data = Policies_Config
	return wop
}


func WOP_Gen_Task_Notification(Task interface{}) (*WOP){
	wop := new(WOP)
	wop.Code = WOP_REQ_TASK
	wop.Data = Task
	return wop
}


func WOP_Gen_State_Report(State interface{}) (*WOP) {
	wop := new(WOP)
	wop.Code = WOP_REQ_STATE
	wop.Data = State
	return wop
}
