package main

const (
	NIL_STRING = ""
	SERVER = 0
	CLIENT = 1
	VERSION = 1
)

type Main struct {
	Type int
	Prefix string
	Startup_Config *Startup_Config
	Policies_Config *Policies_Config
	Certs *Certs
	/*
	 * SERVER only
	 * Mongo db connection
	 */
	Mongo *Mongo
	Net *Net
	/* SERVER only
	 * MPLX channel. Server communicates with clients, watcher, websockets, etc via this channel
	 */
	W chan *MPLX
	/*
	 * CLIENT only
	 * Monitoring channel. Client passes messages here to communicate with server
	 */
	M chan *MON
//	M chan interface{}
	
	/*
	 * SERVER only
	 * MS.Clients: stores clients in a hash
	 */
	MS *MPLX_Sessions
	/*
	 * SERVER only
	 * http & websocket handlers
	 */
	H *HTTP
}
