package main

import (
	"os"
	"fmt"
)

func main() {
	dir := os.Args[1]

	certs := CERTS_Init()
	certs.Load_Certs(dir)
	hash := certs.SHA_Cert()
	fmt.Printf("THE_KEY:%s\n", hash)
	return
}
