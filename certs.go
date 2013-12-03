package main

import (
	"log"
	"fmt"
	"io"
	"crypto/tls"
	"crypto/sha512"
)

const (
	CERT_NAME_PUB = "cert.pem"
	CERT_NAME_KEY = "cert.key"
)

type Certs struct {
	Cert tls.Certificate
}

func CERTS_Init() (*Certs) {
	C := new(Certs)

	DebugLn("certs:Init")

	return C
}

func (M *Main)CERTS_Init() {
	M.Certs = CERTS_Init()
	M.Certs.Load_Certs(M.Prefix)
}

func (C *Certs) Load_Certs(Prefix string) {
	cert, err := tls.LoadX509KeyPair(Prefix+"/"+CERT_NAME_PUB, Prefix+"/"+CERT_NAME_KEY)
	if err != nil {
		log.Fatal("dgo_certs:LoadX509KeyPair",err)
	}

	C.Cert = cert
}

func (C *Certs) SHA_Cert() (string) {
	sha := sha512.New()
	io.WriteString(sha, string(C.Cert.Certificate[0]))
	hash := fmt.Sprintf("%x",sha.Sum(nil))
	return hash
}

func SHA_Cert_Raw(Bytes []byte) (string) {
	sha := sha512.New()
	io.WriteString(sha, string(Bytes))
	hash := fmt.Sprintf("%x", sha.Sum(nil))
	return hash
}
