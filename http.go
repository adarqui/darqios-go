package main

import (
	"fmt"
	"log"
	"time"
	"strings"
	"net/http"
	"github.com/stretchr/goweb"
	"github.com/stretchr/goweb/context"
)

type HTTP struct {
	Ses *http.Server
	Addr string
	ReadTimeout time.Duration
	WriteTimeout time.Duration
	MaxHeaderBytes int
	W chan *MPLX
}

func (M *Main) HTTP_Init() {
	DebugLn("HTTP_Init():Initialized")

	M.H = new(HTTP)
	M.HTTP_Defaults()
	M.HTTP_Setup()
	M.HTTP_Routes()

	go HUB.WS_Run()
	go M.HTTP_Serve()
}


func (M *Main) HTTP_Defaults() {
/* FIXME - only allow localhost */
// H.Addr = "127.0.0.1:8080
	M.H.Addr = fmt.Sprintf("%s:%s",M.Startup_Config.Http.Host, M.Startup_Config.Http.Port) //"0.0.0.0:911"
	M.H.ReadTimeout = time.Duration(2*time.Second)
	M.H.WriteTimeout = time.Duration(2*time.Second)
	M.H.MaxHeaderBytes = 1 << 20
	return
}


func (M *Main) HTTP_Setup() {

	M.H.Ses = &http.Server{
		Addr : M.H.Addr,
		ReadTimeout : M.H.ReadTimeout,
		WriteTimeout : M.H.WriteTimeout,
		MaxHeaderBytes : M.H.MaxHeaderBytes,
		Handler : goweb.DefaultHttpHandler(),
	}

	goweb.MapStaticFile("/", "pub/indeH.html")
	goweb.MapStatic("/pub", "pub")

	goweb.MapBefore(func(c context.Context) (error) {
		req := c.HttpRequest()
		Debug("LOGGER: %v %v\n", req.Method, req.URL)
		return nil
	})

	goweb.MapAfter(func(c context.Context) (error) {
		req := c.HttpRequest()
		Debug("LOGGER: %v %v\n", req.Method, req.URL)
		return nil
	})

 return
}



func (M *Main) HTTP_Routes() {

	M.HTTP_WS_Init()
	M.HTTP_Ping_Init()
	M.HTTP_Accounts_Init()
	M.HTTP_Sessions_Init()
	M.HTTP_State_Init()
	M.HTTP_Help_Init()

}


func (M *Main) HTTP_Serve() {
	err := M.H.Ses.ListenAndServe()
	if err != nil {
		log.Fatal("HTTP_Serve:ListenAndServe()", err)
	}
	return
}


func HTTP_Filter(c context.Context, path string) ([]string, error) {
	list := c.PathValue(path)

	if list == "" {
		return nil, nil
	}
	list_array := strings.Split(string(list),":")

	return list_array, nil
}
