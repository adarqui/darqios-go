package main

func (M *Main) NET_Client_Messenger() {
	DebugLn("NET_Client_Messenger:Initialized")

	for {
		message := <-M.M
		if M.Net.Conn != nil {

			Debug("NET_Client_Messenger:Message:%q\n", message)

			if message.Op == MON_REQ_TASK {

				/* message.Data = IMPORTANT!!!!!! */
				wop := WOP_Gen_Task_Notification(message.Data)
				data, err := WOP_To_Bytes(wop)
				if err != nil {
					continue
				}

				M.Net.Conn.Write(data)

			} else if message.Op == MON_REQ_STATE {

				/* message.Data = IMPORTANT!!!!!! */
				wop := WOP_Gen_State_Report(message.Data)
				data, err := WOP_To_Bytes(wop)
				if err != nil {
					continue
				}

				M.Net.Conn.Write(data)
			}
		}
	}
}
