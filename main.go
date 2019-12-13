package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"

	"simple-http-server/config"
	"simple-http-server/pkg/server"
)

var version = "1.0.0"

var (
	h bool
	v bool
	f bool

	c string
)

func init() {
	flag.BoolVar(&h, "h", false, "this help")
	flag.BoolVar(&v, "v", false, "show version and exit")
	flag.StringVar(&c, "c", "./config/config.json", "set configuration `file`")
}

func main() {
	flag.Parse()

	if h {
		flag.Usage()
		return
	}

	if v {
		fmt.Printf("version: %s\n", version)
		return
	}

	err := config.LoadConfig(c)
	if err != nil {
		fmt.Println("LoadConfig error", err)
		return
	}

	config, err := config.GetConfig()
	if err != nil {
		fmt.Println("GetConfig error", err)
		return
	}

	q := make(chan error, 1)

	addr := ":" + strconv.Itoa(config.Port)
	// Run our server in a goroutine so that it doesn't block.
	server.Serve(addr, q)

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	select {
	case <-c:
		fmt.Println("srv got SIGINT!!!")
		server.Stop()
		fmt.Println("shutting down")
		os.Exit(0)
	case err = <-q:
		fmt.Println("srv exit error", err)
		os.Exit(-1)
	}
}
