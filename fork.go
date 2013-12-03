package main

/*
 * This is the "fork in the road" between becoming a client or a server
 */

func (M *Main) Fork() {

	GOB_Init()

	switch M.Type {
		case SERVER: {
			M.Fork_Server()
		}
		case CLIENT: {
			M.Fork_Client()
		}
	}
}

func (M *Main) Fork_Server() {

	if M.Startup_Config.Server.Daemonize == true {
		Daemon()
	}

	M.MG_Init()
	M.NET_Server()
}

func (M *Main) Fork_Client() {

	if M.Startup_Config.Client.Daemonize == true {
		Daemon()
	}

	M.NET_Client()
}
