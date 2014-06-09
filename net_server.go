package main

import (
	"fmt"
	"log"
	"net"
	"strings"
	"crypto/tls"
	"crypto/rand"
	"crypto/x509"
)

func (M *Main) NET_Server() {

	DebugLn("NET_Server:Init")

	config := tls.Config{Certificates: []tls.Certificate{M.Certs.Cert}, ClientAuth: tls.RequireAnyClientCert}
	config.Rand = rand.Reader
	config.BuildNameToCertificate()
	listener, err := tls.Listen("tcp", M.Startup_Config.Server.Host+":"+string(M.Startup_Config.Shared.Port), &config)
	if err != nil {
		log.Fatalf("NET_Server:tls.Listen:%s",err)
	}
	
	Debug("NET_Server:Listening on %s %i\n", M.Startup_Config.Server.Host, M.Startup_Config.Shared.Port)

	DebugLn(listener)

//	w := make(chan Multiplex,100)

	/*
	 * IMPORTANT.
	 * Need an initial policies_config, otherwise this is pointless
	 */
	 policies_config,err := M.MG_Lookup_Policies_Config_Raw(M.Startup_Config.Server.Policies)
	 if err != nil {
		 log.Fatal("NET_Server:Unable to obtain policies_config",err)
	 }

	 M.Policies_Config = policies_config

	/*
	 * Special MPLX channel
	 */
	M.W = make(chan *MPLX,100)

	go M.MPLX_Init()
	go M.WATCH_Init()
	go M.HUB_Init()
	go M.HTTP_Init()

	for {
		conn, err := listener.Accept()
		if err != nil {
			break
		}
		go M.NET_Server_Handle_Client(conn)
	}

	return
}

func (M *Main) NET_Server_Handle_Client(Conn net.Conn) {
	Debug("NET_Server_Handle_Client:%q\n", Conn)

	authenticated := false
	defer_command := MPLX_REQ_NOOP
	var account *Account

	mplx := MPLX_Create()

	defer func() {
		Conn.Close()
		mplx.Op = defer_command
		mplx.Data = account
		M.W <- mplx
		_ = mplx.W
		return
	}()

	tlsConn := Conn.(*tls.Conn)
	err := tlsConn.Handshake()
	/*
	 * Handshake needed in order to verify client pubkey against db.accounts
	 */
	if err != nil {
		return
	}

	state := tlsConn.ConnectionState()
	for _, v := range state.PeerCertificates {
		_, err := x509.MarshalPKIXPublicKey(v.PublicKey)
		if err != nil {
			continue
		}

		hash := SHA_Cert_Raw([]byte(v.Raw))
		mplx.Op = MPLX_REQ_LOOKUP_ACCOUNT
		mplx.Arg = hash
		M.W<-mplx
		message := <-mplx.W

		if message.Op == MPLX_REP_OK {
			if acc,ok := message.Data.(*Account); ok {
				authenticated = true
				account = acc
				break
			}
		}
	}

	if authenticated == false {
		return
	}

	DebugLn("NET_Server_Handle_Client:Authenticated!")


	/* Get rid of ip:port */
	ip_arr := strings.Split(fmt.Sprintf("%s", Conn.RemoteAddr()), ":")
	if len(ip_arr) != 2 {
		authenticated = false
		return
	}

	account.Ip = ip_arr[0]

	Debug("NET_Server_Handle_Client:IP=%v\n",account.Ip)

	mplx.Op = MPLX_REQ_NEW_CONNECTION
	mplx.Arg = account.Host
	mplx.Data = account
	M.W<-mplx
	message := <-mplx.W

	if message.Op == MPLX_REP_DIE {
		return
	}

	defer_command = MPLX_REQ_END_CONNECTION

	/*
	 * Request policies_config from the MPLX worker thread
	 */
	mplx.Op = MPLX_REQ_POLICIES_CONFIG
	mplx.Data = account
	M.W<-mplx
	message = <-mplx.W

	/*
	 * We need to make sure we send over a working config
	 */
	if message.Op == MPLX_REP_DATA {
		if data,ok := message.Data.(*Policies_Config); ok {

			Debug("SENDING POLICIES_CONFIG TO CLIENT: %q\n", data)

			wop := WOP_Gen_Rep_Policies_Config(data)
			wop_bytes, err := WOP_To_Bytes(wop)
			if err != nil {
				return
			}
			Conn.Write(wop_bytes)
		} else {
			/* What? */
			return
		}

		/* Force dynamic configuration of hostname, eventually make this optional */
		wop := WOP_Gen_Rep_Hostname_Config(account.Host)
		wop_bytes, _ := WOP_To_Bytes(wop)
		Conn.Write(wop_bytes)

	} else {
		return
	}


	buf := make([]byte, 10024)

	for {
		sz, err := Conn.Read(buf)
		if err != nil {
			Debug("NET_Server_Handle_Client:Read:Err:%q\n",err)
			return
		}

		Debug("NET_Server_Handle_Client:Sz=%d\n",sz)

		/*
		 * Handle operations from the client, these will be:
		 *
		 * tasks - alerts
		 * state - state_report's
		 */

		 /*
		mplx.Op = MPLX_REQ_MON
		mplx.Data = buf
		M.W<-mplx
		message = <-mplx.W
		*/

		wop,err := WOP_Parse(buf)
		if err != nil {
			continue
		}

		switch wop.Code {
			case WOP_REQ_TASK: {
				mplx.Op = MPLX_REQ_TASK
			}
			case WOP_REQ_STATE: {
				mplx.Op = MPLX_REQ_STATE
			}
		}

		mplx.Data = wop.Data
		M.W<-mplx
		message = <-mplx.W

		switch message.Op {
			case MPLX_REP_DIE: {
				return
			}
			case MPLX_REP_OK: {
				continue
			}
			case MPLX_REP_DATA: {
				if data,ok := message.Data.([]byte); ok {
					Conn.Write(data[0:len(data)])
				}
			}
		}
	}

}
