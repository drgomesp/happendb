package main

import (
	"context"
	"flag"
	"fmt"
	stdlog "log"
	"os"
	"os/signal"
	"syscall"

	abciserver "github.com/tendermint/tendermint/abci/server"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/drgomesp/happendb/pkg/abci"
	"github.com/drgomesp/happendb/pkg/store"
)

var socketAddr string

func init() {
	flag.StringVar(&socketAddr, "socket-addr", "tcp://happendb:26658", "Unix domain socket address")
}

func main() {
	flag.Parse()

	ctx := context.Background()
	sig := make(chan os.Signal, 1)

	go handleSignals(ctx, sig, func(ctx context.Context) error {
		stdlog.Println("shutting down...")
		return nil
	})

	app := abci.NewApplication(store.NewMemory())
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

func handleSignals(ctx context.Context, sig chan os.Signal, shutdownFunc func(ctx context.Context) error) {
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sig)

	<-sig
	stdlog.Println("attempting graceful shutdown...")
	if err := shutdownFunc(ctx); err != nil {
		stdlog.Println(err)
		os.Exit(-1)
	}

	os.Exit(0)
}
