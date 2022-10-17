package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/abci/types"
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/libs/service"
	nm "github.com/tendermint/tendermint/node"
	p2p2 "github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/privval"
	"github.com/tendermint/tendermint/proxy"

	happendb "github.com/drgomesp/happendb/pkg/abci"
	"github.com/drgomesp/happendb/pkg/store"
)

var configFile string
var socketAddr string

func init() {
	flag.StringVar(&socketAddr, "socket-addr", "tcp://127.0.0.1:26658", "Unix domain socket address")
	flag.StringVar(&configFile, "config", "$HOME/.tendermint/config/config.toml", "Path to config.toml")
}

func main() {
	flag.Parse()

	app := happendb.NewApplication(store.NewMemory())

	//server := abciserver.NewSocketServer(socketAddr, app)
	//server.SetLogger(logger)

	node, err := newTendermint(app, configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(2)
	}

	node.Start()
	defer func() {
		node.Stop()
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	os.Exit(0)

}
func newTendermint(app types.Application, configFile string) (service.Service, error) {
	// read config
	config := cfg.DefaultConfig()
	config.RootDir = filepath.Dir(filepath.Dir(configFile))
	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("viper failed to read config file: %w", err)
	}
	if err := viper.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("viper failed to unmarshal config: %w", err)
	}
	if err := config.ValidateBasic(); err != nil {
		return nil, fmt.Errorf("config is invalid: %w", err)
	}

	// read private validator
	pv := privval.LoadFilePV(
		config.PrivValidatorKeyFile(),
		config.PrivValidatorStateFile(),
	)

	// read node key
	nodeKey, err := p2p2.LoadNodeKey(config.NodeKeyFile())
	if err != nil {
		return nil, fmt.Errorf("failed to load node's key: %w", err)
	}

	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
	// create node
	node, err := nm.NewNode(
		config,
		pv,
		nodeKey,
		proxy.NewLocalClientCreator(app),
		nm.DefaultGenesisDocProviderFunc(config),
		nm.DefaultDBProvider,
		nm.DefaultMetricsProvider(config.Instrumentation),
		logger,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create new Tendermint node: %w", err)
	}

	return node, nil
}
