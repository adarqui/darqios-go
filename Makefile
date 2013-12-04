OBJS=debug.go defaults.go startup_config.go fork.go certs.go mongo.go schema.go net.go net_server.go net_client.go mplx.go watch.go hub.go http.go mongo_policies.go policy.go wire_ops.go gob.go net_client_mon.go net_client_messenger.go mon.go state.go tasks.go tasks_ping.go tasks_load.go tasks_process.go tasks_memory.go tasks_disk.go tasks_state.go tasks_scheduler.go ws.go routes.go mplx_broadcast.go http_ws.go http_accounts.go http_ping.go http_sessions.go http_state.go http_policies.go http_help.go state_network.go state_disks.go daemon.go inc.go

all:
	make clean
	go build darqios.go $(OBJS)
	go build sha_cert.go $(OBJS)


deps:
	go get labix.org/v2/mgo
	go get github.com/c9s/goprocinfo/linux
	go get github.com/gorilla/websocket
	go get github.com/stretchr/goweb
	go get github.com/VividCortex/godaemon

run_server:
	./darqios server

run_client:
	./darqios client testing/clients/certs/client

clean:
	rm -f darqios
