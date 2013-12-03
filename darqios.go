package main

/*
put http params in conf.js (or policies.js)
put disconnect warn option in conf(or policies)
get mitigate done
clean up MON_Gen_Task.. only allow concat'n somehing to subject..
*/



import (
	"log"
	"os"
)

const (
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

func usage() {
	log.Fatal("go run darqios.go <server|client>")
}

func Main_Init() (*Main) {
	m := new(Main)
	return m
}


func (M *Main) Init() {
	if len(os.Args) > 2 {
		M.Prefix = os.Args[2]
	}
}


func main() {

	if len(os.Args) <= 1 {
		usage()
	}

	DebugLn("darqios: Initialized");

	M := Main_Init()

	mode := os.Args[1]
	switch mode {
		case "server" : {
			M.Type = SERVER
		}
		case "client" : {
			M.Type = CLIENT
		}
		default : {
			usage()
		}
	}

	M.Defaults()
	M.Init()
	M.SC_Init()
	M.CERTS_Init()

	Debug("main:M=%q\n",M)

	M.Fork()

	select {}
}
