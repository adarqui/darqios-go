package main

func (M *Main) Defaults() {
	switch M.Type {
		case SERVER: {
			M.Prefix = "./certs/server"
		}
		case CLIENT: {
			M.Prefix = "./certs/client"
		}
	}
}
