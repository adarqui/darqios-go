package main

import (
	"time"
	"fmt"
	"os/exec"
)

type MPLX struct {
	Op int
	Arg string
	Args []string
	W chan *MPLX
	Data interface{}
	Account *Account
}

type Session struct {
	Account *Account
	TS_Last time.Time
	Count int
	Active bool
}

type MPLX_Sessions struct {
	Clients map[string]*Session
}

type MPLX_Fn struct {
	Name string
	Fn *func ()bool
}

/*
 * REQUESTS
 */
const (
	MPLX_REQ_NOOP = 0
	MPLX_REQ_LIST_SESSIONS = 1
	MPLX_REQ_NEW_CONNECTION = 2
	MPLX_REQ_END_CONNECTION = 3
	MPLX_REQ_POLICIES_CONFIG = 4
	MPLX_REQ_LOOKUP_ACCOUNT = 5
//	MPLX_REQ_MON = 6
	MPLX_REQ_TASK = 6
	MPLX_REQ_STATE = 7
	MPLX_REQ_MISSING = 8
)

/*
 * REPLIES
 */
const (
	MPLX_REP_OK = 1
	MPLX_REP_NOTOK = 2
	MPLX_REP_DATA = 3
	MPLX_REP_DIE = 4
)

/*
 * MISc
 */
const (
	MPLX_OPT_CONNECTION_NEW = 0
	MPLX_OPT_CONNECTION_END = 1
)

func (M *Main) MPLX_Init() {
	/*
	 * Initialized the main MPLX 'circuit'
	 */
	DebugLn("MPLX_Init():Initialized")

	M.MS = new(MPLX_Sessions)
	M.MS.Clients = make(map[string]*Session)

	defer func() {
		DebugLn("MPLX_Init:Exiting")
	}()

	fns := map[int]func(*Main,*MPLX)bool {
		MPLX_REQ_NOOP : MPLX_REQ_NoOp,
		MPLX_REQ_LIST_SESSIONS : MPLX_REQ_List_Sessions,
		MPLX_REQ_NEW_CONNECTION : MPLX_REQ_New_Connection,
		MPLX_REQ_END_CONNECTION : MPLX_REQ_End_Connection,
		MPLX_REQ_POLICIES_CONFIG : MPLX_REQ_Policies_Config,
		MPLX_REQ_LOOKUP_ACCOUNT : MPLX_REQ_Lookup_Account,
//		MPLX_REQ_MON : MPLX_REQ_Mon,
		MPLX_REQ_TASK : MPLX_REQ_Task,
		MPLX_REQ_STATE: MPLX_REQ_State,
		MPLX_REQ_MISSING: MPLX_REQ_Task,
	}

	for {
		message := <-M.W
//		Debug("MPLX_Init():Message:%q\n",message)

		Debug("MPLX_Init():%v\n", message.Op)

		if idx,ok := fns[message.Op]; ok {
			/*
			 * Process an MPLX operation
			 */
			go func() {
				/*
				 * Return op is used for MPLX_Broadcast etc.. tells us what we 'really are'
				 */
				message_copy := *message
				go M.MPLX_Broadcast(&message_copy)
				truth := idx(M,message)
				Debug("MPLX_Init():Truth=%q\n",truth)
				message.W<-message

				/*
				 * Broadcast an MPLX operation to portential listeners/subscribers on a websocket
				 */
				//go M.MPLX_Broadcast(&message_copy)
			}()
		}
	}
}

func MPLX_Create() (*MPLX) {
	/*
	 * Creates an MPLX 'circuit'
	 */
	mplx := new(MPLX)
	mplx.W = make(chan *MPLX)
	return mplx
}


func MPLX_REQ_NoOp(M *Main, Msg *MPLX) (bool) {
	/*
	Msg.Op = MPLX_REP_OK
	Msg.W<-Msg
	*/
	return true
}

func MPLX_REQ_List_Sessions(M *Main, Msg *MPLX) (bool) {
	/*
	Msg.W<-Msg
	*/
	Msg.Op = MPLX_REP_OK
	Msg.Data = M.MS
	return false
}

func MPLX_REQ_New_Connection(M *Main, Msg *MPLX) (bool) {
	return MPLX_REQ_Generic_Connection(M,Msg,MPLX_OPT_CONNECTION_NEW)
}

func MPLX_REQ_End_Connection(M *Main, Msg *MPLX) (bool) {
	return MPLX_REQ_Generic_Connection(M,Msg,MPLX_OPT_CONNECTION_END)
}

func MPLX_REQ_Generic_Connection(M *Main, Msg *MPLX, Type int) (bool) {
	/*
	 * Type:
	 * 1 = new
	 * 0 = end
	 */

	var account *Account

	/*
	defer_command := MPLX_REP_DIE
	defer func() {
		Msg.Op = defer_command
		Msg.W<-Msg
	}()
	*/

	if acc,ok := Msg.Data.(*Account); ok {
		account = acc
		if _, ok := M.MS.Clients[account.Host]; ok {
			if Type == MPLX_OPT_CONNECTION_NEW {
				/* Already exists */
				return false
			}
		} else /* !ok */ {
			if Type == MPLX_OPT_CONNECTION_END {
				/* Doesn't exist */
				return false
			} else {
				/* Initialize it */
				M.MPLX_Session_Init(acc.Host)
			}
		}
	} else {
		return false
	}

	/*
	 * Ok a bunch of logic
	 */

	 switch Type {
		case MPLX_OPT_CONNECTION_NEW: {
			Ses := M.MS.Clients[account.Host]
			Ses.Active = true
			Ses.Account = account
			Ses.TS_Last = time.Now()
//			defer_command = MPLX_REP_OK
			Msg.Op = MPLX_REP_OK
		}
		case MPLX_OPT_CONNECTION_END: {
			delete(M.MS.Clients,account.Host)
//			defer_command = MPLX_REP_OK
			Msg.Op = MPLX_REP_OK
		}
	 }

	 go func() {
		 _,_ = M.MG_Update_Account_Last_Seen(account)
	 }()

	 return true
}



func MPLX_REQ_Policies_Config(M *Main, Msg *MPLX) (bool) {

	Msg.Op = MPLX_REP_NOTOK

	/*
	defer func() {
		Msg.W<-Msg
	}()
	*/

	if acc,ok := Msg.Data.(*Account); ok {
		policies_config,err := M.MG_Lookup_Policies_Config(M.Startup_Config.Server.Policies,acc)
		if err != nil {
			Debug("MPLX_REQ_Policies_Config:MG_Lookup_Policies_Config:Err:%q\n",err)
			return false
		}
		Msg.Op = MPLX_REP_DATA
		Msg.Data = policies_config
		return true
	}

	return false
}

func MPLX_REQ_Lookup_Account(M *Main, Msg *MPLX) (bool) {

	if len (Msg.Args) > 0 {
		return MPLX_REQ_Lookup_Accounts(M,Msg)
	}

	Msg.Op = MPLX_REP_NOTOK

	/*
	defer func() {
		Msg.W<-Msg
	}()
	*/

	switch len(Msg.Arg) {
		case 128: {
			account,_ := M.MG_Lookup_Account_By_Hash(Msg.Arg)
			if account != nil && account.Status != false {
				Msg.Op = MPLX_REP_OK
				Msg.Data = account
				return true
			}
		}
		default : {
		}
	}
	return false
}

func MPLX_REQ_Lookup_Accounts(M *Main, Msg *MPLX) (bool) {
	return false
}


func MPLX_REQ_Task(M *Main, Msg *MPLX) (bool) {

	var Ses *Session
	var Acc *Account
	var T *Task

	Debug("DATA:%q\n", Msg.Data)
	switch task := Msg.Data.(type) {
		case Task: {
			T = &task
		}
		case *Task: {
			T = task
		}
		default: {
			return false
		}
	}

	Debug("TASK: Policy=%v\n", T.Name)

	if Msg.Account != nil {
		Acc = Msg.Account
	} else if ses,ok := M.MS.Clients[Msg.Arg]; ok {
		Ses = ses
		Acc = Ses.Account
	} else {
		return false
	}

	return MPLX_REQ_Task_Raw(M, Acc, T)
}


func MPLX_REQ_State(M *Main, Msg *MPLX) (bool) {

	var Ses *Session
	var Acc *Account
	var S *State_Report

	Debug("STATE: %q\n", Msg.Data)

	switch state := Msg.Data.(type) {
		case State_Report: {
			S = &state
		}
		case *State_Report: {
			S = state
		}
		default: {
			return false
		}
	}

	/*
	if state,ok := Msg.Data.(*State_Report); ok {
		S = state
	} else {
		return false
	}
	*/

	if Msg.Account != nil {
		Acc = Msg.Account
	} else if ses,ok := M.MS.Clients[Msg.Arg]; ok {
		Ses = ses
		Acc = Ses.Account
	} else {
		return false
	}

	return MPLX_REQ_State_Raw(M, Acc, S)
}

func MPLX_REQ_Task_Raw(M *Main, Acc *Account, T *Task) (bool) {
	/*
	 * Handle a raw task: This is an alert/notification
	 */

	Debug("MPLX_REQ_Task_Raw:Entered\n")

	if Acc == nil {
		return false
	}

	if T.Type == "none" {
		return false
	}

	M.MG_Insert_Task(Acc, T)

	alert_type := new(Alert)
	switch(T.Type) {
		case "low" : {
			alert_type = M.Policies_Config.Alerts.Low
		}
		case "med" : {
			alert_type = M.Policies_Config.Alerts.Med
		}
		case "high" : {
			alert_type = M.Policies_Config.Alerts.High
		}
		case "clear" : {
			alert_type = M.Policies_Config.Alerts.Clear
		}
		default : {
			return false
		}
	}

	for k, j := range(alert_type.Scripts) {
		cmd := exec.Command(
		fmt.Sprintf("%s/alerts/%s", M.Policies_Config.Base, j),
			Acc.Host,
			Acc.Ip,
			T.Type,
			T.Name,
			T.Idx,
			T.Subject,
			T.Body,
			fmt.Sprintf("%s",T.Time),
			T.Actual,
			T.Operator,
			T.Threshold,
		)
		err := cmd.Run()
		if err != nil {
			Debug("cmd.Run() error %q\n", err, k, j)
			return false
		}
		return true
	}


	return false
}


func MPLX_REQ_State_Raw(M *Main, Acc *Account, S *State_Report) (bool) {
	/*
	 * Handle a raw state report: This is a notification about device state, irregardless of thresholds
	 */

	M.MG_Update_Account_State(Acc, S)

	return false
}


/*
func MPLX_REQ_Mon(M *Main, Msg *MPLX) (bool) {

	*
	 * Expects Msg.Arg to be the "host" db.accounts.find({host:...})
	 *

	Msg.Op = MPLX_REP_NOTOK

	var Ses *Session
	var Acc *Account

	if ses,ok := M.MS.Clients[Msg.Arg]; ok {
		Ses = ses
		Acc = Ses.Account
	} else {
		if Msg.Supp != nil {
			if acc,ok := Msg.Supp.(*Account); ok {
				Acc = acc
			} else {
				return false
			}
		} else {
			return false
		}
	}

//	Acc := Ses.Account

	Debug("MPLX_REQ_Lookup_Accounts:Acc:%q\n", Acc)

	var mon *MON

	switch data := Msg.Data.(type) {
		case *MON: {
			DebugLn("MON")
			mon = data
		}
		case []byte: {
			DebugLn("BYTE")
			wop, err := WOP_Parse(data)
			if err != nil {
				return false
			}

			Debug("WOP:%q\n", wop)

			switch wop_data := wop.Data.(type) {
				case MON: {
					DebugLn("MON!!!!!!")
					mon = &wop_data
				}
				default: {
					DebugLn("NEITHER!!!!!")
					return false
				}
			}
		}
		default : {
			return false
		}
	}

	if mon != nil {
		DebugLn("MON != nil")
		switch mon_data := mon.Data.(type) {
			case *Task: {
				task := mon_data
				return MPLX_REQ_Task_Raw(M, Acc, task)
			}
			case Task: {
				task := mon_data
				return MPLX_REQ_Task_Raw(M, Acc, &task)
			}
			case State_Report: {
				state_report := mon_data
				return MPLX_REQ_State_Raw(M, Acc, &state_report)
			}
			default : {
				return false
			}
		}
	}


	return false
}
*/




func (M *Main)MPLX_Session_Init(Host string) {
	Ses := new(Session)
	M.MS.Clients[Host] = Ses
}


func MPLX_REQ_2_STRING(Op int) (string) {
	switch Op {
		case MPLX_REQ_NOOP: return "noop"
		case MPLX_REQ_NEW_CONNECTION: return "new:connection"
		case MPLX_REQ_END_CONNECTION: return "end:connection"
		case MPLX_REQ_TASK: return "task"
		case MPLX_REQ_STATE: return "state"
		case MPLX_REQ_MISSING: return "missing"
		default: return "unknown"
	}
}

/*
⇥   MPLX_REQ_NOOP = 0
⇥   MPLX_REQ_LIST_SESSIONS = 1
⇥   MPLX_REQ_NEW_CONNECTION = 2
⇥   MPLX_REQ_END_CONNECTION = 3
⇥   MPLX_REQ_POLICIES_CONFIG = 4
⇥   MPLX_REQ_LOOKUP_ACCOUNT = 5
⇥   MPLX_REQ_MON = 6
⇥   MPLX_REQ_TASK = 7
⇥   MPLX_REQ_STATE = 8
*/
