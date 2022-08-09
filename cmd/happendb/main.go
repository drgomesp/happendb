package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	abciserver "github.com/tendermint/tendermint/abci/server"
	"github.com/tendermint/tendermint/libs/log"

	happendb "github.com/drgomesp/happendb/pkg/abci"
	"github.com/drgomesp/happendb/pkg/store"
)

var socketAddr string

func init() {
	flag.StringVar(&socketAddr, "socket-addr", "tcp://127.0.0.1:26658", "Unix domain socket address")
}

func main() {
	flag.Parse()

	app := happendb.NewApplication(store.NewMemory())
	logger := log.MustNewDefaultLogger(log.LogFormatPlain, log.LogLevelDebug, false)
	server := abciserver.NewSocketServer(socketAddr, app)
	server.SetLogger(logger)

	if err := server.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "error starting socket server: %v", err)
		os.Exit(1)
	}

	defer server.Stop()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	os.Exit(0)
}