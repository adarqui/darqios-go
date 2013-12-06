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
		time.Sleep(M.Startup_Config.Server.Watcher*time.Second)

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
	P.Name = "Missing"

	P.Idx = "lost"
	task := MON_Gen_Task_Raw("high", A.Host, &P, fmt.Sprintf("%s hasn't checked in for the last %s seconds", A.Host, Time_Diff), "Bleh")

	Msg.Op = MPLX_REQ_MISSING
	Msg.Arg = A.Host
	Msg.Data = task
	Msg.Account = A
	M.W<-Msg
	message := <-Msg.W

	Debug("WATCH_Notify:%q\n", message)
}
