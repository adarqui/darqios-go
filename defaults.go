package main

func (M *Main) Defaults() {
	DEBUG = false
	switch M.Type {
		case SERVER: {
			M.Prefix = "./certs/server"
		}
		case CLIENT: {
			M.Prefix = "./certs/client"
		}
	}
}
