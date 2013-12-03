package main

import (
	"log"
	"io/ioutil"
	"encoding/json"
)

type Shared_Config struct {
	Port string
}

type Server_Config struct {
	Host string
	Policies string
	Daemonize bool
}

type Client_Config struct {
	Host string
	Daemonize bool
}

type Mongo_Config struct {
	Addr string
	User string
	Pass string
	Db string
}

type Http_Config struct {
	Host string
	Port string
}


type Startup_Config struct {
	Shared Shared_Config
	Server Server_Config
	Client Client_Config
	Mongo Mongo_Config
	Http Http_Config
}

func (M *Main) SC_Init() {

	Debug("startup_config:Init:%i\n", M.Type)

	SC := new(Startup_Config)

	file, err := ioutil.ReadFile("conf.json")
	if err != nil {
		log.Fatal("startup_config:Init",err)
	}

	err = json.Unmarshal(file, SC)
	if err != nil {
		log.Fatal("startup_config:Json:Unmarshal", err)
	}

	M.Startup_Config = SC
	Debug("startup_config:%q\n", M.Startup_Config)
	return
}
