package main

import (
//	"fmt"
	"time"
	"crypto/tls"
//	"crypto/x509"
)

func (M *Main) NET_Client() {
	DebugLn("NET_Client:Init")

	M.M = make(chan *MON, 100)
	M.Net = new(Net)

	go M.NET_Client_MON()
	go M.NET_Client_Messenger()

	for {
		DebugLn("NET_Client:Initializing a new client connection")
		M.NET_Client_Loop()
		time.Sleep(1*time.Second)
	}

	return
}

func (M *Main) NET_Client_Loop() (bool) {
	config := tls.Config{Certificates: []tls.Certificate{M.Certs.Cert}, InsecureSkipVerify: true, ClientAuth: tls.RequireAnyClientCert}

	config.BuildNameToCertificate()

	conn, err := tls.Dial("tcp", M.Startup_Config.Client.Host+":"+M.Startup_Config.Shared.Port, &config)
	if err != nil {
		Debug("NET_Client_Loop:Error:%q\n",err)
		return false
	}

	Debug("NET_Client_Loop:Connected to:%v\n",conn.RemoteAddr())

	M.Net.Conn = conn
	M.NET_Client_Handle_Server()

	return true
}

func (M *Main) NET_Client_Handle_Server() {
	defer func() {
		M.Net.Conn.Close()
		M.Net.Conn = nil
		M.Net.Count += 1
//		M.Policies_Config = nil
	} ()

	reply := make([]byte, 10024)

	for {
		sz, err := M.Net.Conn.Read(reply)
		if err != nil {
			Debug("NET_Client_Handle_Serve):Read:%q\n",err)
			break
		}

		Debug("NET_Client_Handle_Server:%d:%q\n",sz,reply)

		wop, err := WOP_Parse(reply)
		if err != nil {
			Debug("NET_Client_Handle_Server:WOP_Parse:Err:%q\n", err)
			continue
		}
		
		M.NET_Client_Handle_WOP(wop)
	}
}



func (M *Main) NET_Client_Handle_WOP(wop *WOP) {
	Debug("WOP! CODE=%v %q\n", wop.Code, wop)

	switch wop.Code {
		case WOP_REP_POLICIES_CONFIG: M.NET_Client_Handle_WOP_Policies_Config(wop)
	}
}


func (M *Main) NET_Client_Handle_WOP_Policies_Config(wop *WOP) {
	if policies_config,ok := wop.Data.(Policies_Config); ok {
		DebugLn("NET_Client_Handle_WOP_Policies_Config:Success")
		M.Policies_Config = &policies_config
	}
	return
}
