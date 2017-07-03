package main

import (
	"flag"
	"fmt"
	"net/http"
	"runtime"
	"strings"

	"github.com/swanwish/fileserver/handlers/helper"
	"github.com/swanwish/fileserver/settings"
	"github.com/swanwish/go-common/logs"
	"github.com/swanwish/go-common/utils"
)

const (
	DefaultPort = 8080
)

var (
	port           int64
	rootPath       string
	configFilePath string
)

func parseCmdLineArgs() {
	flag.Int64Var(&port, "port", DefaultPort, "The port to listen")
	flag.StringVar(&rootPath, "path", ".", "The path of the md files")
	flag.StringVar(&configFilePath, "c", "conf/app.ini", "The configuration file path")
	flag.Parse()
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	parseCmdLineArgs()

	if configFilePath != "" {
		settings.ConfigFilePath = configFilePath
	}
	logs.Debugf("The configuration file path is %s", settings.ConfigFilePath)

	settings.LoadAppSetting()

	if port == 0 {
		port = DefaultPort
	}

	localIps, err := utils.GetLocalIPAddrs()
	if err != nil {
		fmt.Println("Failed to get local ip addresses.")
		return
	}

	fmt.Printf("Service listen on port \x1b[31;1m%d\x1b[0m and server ip addresses are \x1b[31;1m%s\x1b[0m\n", port, strings.Join(localIps, ", "))

	http.Handle("/", http.StripPrefix("/", http.FileServer(helper.FSDir(rootPath))))

	httpAddr := fmt.Sprintf(":%d", port)
	if err := http.ListenAndServe(httpAddr, nil); err != nil {
		fmt.Printf("http.ListendAndServer() failed with %s\n", err)
	}

	fmt.Println("Exited")
}
