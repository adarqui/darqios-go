package main

import (
	"encoding/json"
)

func (M *Main) MPLX_Broadcast(Msg *MPLX) {
	Debug("MPLX_Broadcast:Entered:%v\n", Msg)


	bmsg := BMSG{}

	bmsg_data := BMSG_DATA{}
	bmsg_data.Channel = MPLX_REQ_2_STRING(Msg.Op)
	bmsg_data.Host = Msg.Arg
	bmsg_data.Data = Msg.Data

	jsn, err := json.Marshal(bmsg_data)
	if err != nil {
		Debug("MPLX_Broadcast:json.Marshal:Err:%q\n", err)
		return
	}

	bmsg.Channel = bmsg_data.Channel
	bmsg.Data = jsn

	HUB.broadcast <- bmsg
}
