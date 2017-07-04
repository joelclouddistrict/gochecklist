package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/jllopis/getconf"
	"github.com/joelclouddistrict/gochecklist/log"
	"github.com/joelclouddistrict/gochecklist/server"
)

var (
	// BuildDate holds the date the binary was built. It is valued at compile time
	BuildDate string
	// Version holds the version number of the build. It is valued at compile time
	Version string
	// Revision holds the git revision of the binary. It is valued at compile time
	Revision string
	// config struct
	config *getconf.GetConf
)

// Config proporciona la configuración del servicio para ser utilizado por getconf
type Config struct {
	GwAddr    string `getconf:"gwaddr, default 127.0.0.1, info http rest gateway address"`
	Port      string `getconf:"port, default 8000, info default port to listen on"`
	Mode      string `getconf:"mode, default dev, info startup mode"`
	StoreHost string `getconf:"store-host, default db.acb.info, info store server address"`
	StorePort string `getconf:"store-port, default 5432, info store server port"`
	StoreName string `getconf:"store-name, default annotationdb, info store database name"`
	StoreUser string `getconf:"store-user, default annotationadm, info store user to connect to db"`
	StorePass string `getconf:"store-pass, default 00000000, info store user password"`
}

func main() {
	defer func() {
		if e := recover(); e != nil {
			stack := make([]byte, 1<<16)
			stackSize := runtime.Stack(stack, true)
			log.Info(string(stack[:stackSize]))
		}
	}()

	// load config (jllopis/getconf)
	config = getconf.New("TODO", &Config{})

	// capture signals
	setupSignals()

	mssrv := server.New(config)
	mssrv.Register()

	log.Err(mssrv.Serve().Error())
}

// setupSignals configura la captura de señales de sistema y actúa basándose en ellas
func setupSignals() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		for sig := range sc {
			log.Info(fmt.Sprintf("Capturada señal: %v", sig))
			os.Exit(1)
		}
	}()
}
