package main

import (
	"fmt"
	"time"
)

func (M *Main) WATCH_Init() {
	DebugLn("WATCH_Init():Initialized")

	mplx := MPLX_Create()

	mplx.Op = MPLX_REQ_LIST_SESSIONS

	for {
		/*
		DebugLn("WATCH_Init:Looping")
		M.W<-mplx
		message := <-mplx.W

		Debug("WATCH_Init():Message:%q\n", message)

		* FIXME - if no sessions, won't get reports *
		if sessions,ok := message.Data.(*MPLX_Sessions); ok {
			Debug("WATCH_Init():Sessions:%q\n", sessions)

			now := time.Now()
			accounts, err := M.MG_Find_Missing(now)

			Debug("ACCOUNTS:%q\n", accounts)

			if err != nil {
				Debug("WATCH_Init:Err:%q\n",err)
				continue
			}

			for _, acc := range accounts {
				Debug("WATCH_Init:Notifying for:%q\n",acc)
				time_diff := now.Sub(acc.Last)
				go M.WATCH_Notify(mplx, &acc, time_diff)
			}
		}
		*/
		time.Sleep(5*time.Second)

		now := time.Now()
		accounts, err := M.MG_Accounts_Find_Missing(now)

		if err != nil {
			continue
		}

		for k, acc := range accounts {
			if acc.Host == "" {
				/* Empty host */
				continue
			}
			Debug("WATCH_Init:Notifying for:%q\n", acc)
			time_diff := now.Sub(acc.Last)
			acc_ptr := accounts[k]
			mplx_new := MPLX_Create()
			if _, ok := M.MS.Clients[acc_ptr.Host]; !ok {
				/* Only notify if the client is not in session */
				go M.WATCH_Notify(mplx_new, &acc_ptr, time_diff)
			}
		}
	}
}

func (M *Main) WATCH_Notify(Msg *MPLX, A *Account, Time_Diff time.Duration) {

	P := Policy{}
	P.Name = "Lost"

	task := MON_Gen_Task_Raw("high", &P, fmt.Sprintf("Lost - Level (%s) : %s hasn't checked in for the last %s seconds", "high", A.Host, Time_Diff), "Bleh")
//	wop := WOP_Gen_Task_Notification(mon)

/*
	Msg.Op = MPLX_REQ_MON
	Msg.Arg = A.Host
	Msg.Data = mon
	Msg.Supp = A
	M.W<-Msg
	message := <-Msg.W
	*/

	Msg.Op = MPLX_REQ_MISSING
	Msg.Arg = A.Host
	Msg.Data = task
	Msg.Account = A
	M.W<-Msg
	message := <-Msg.W

	Debug("WATCH_Notify:%q\n", message)
}
