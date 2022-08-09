package happendb

import (
	"context"

	"github.com/tendermint/tendermint/rpc/client"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
)

type Client struct {
	client.ABCIClient
}

func NewClient(remote string) (*Client, error) {
	abci, err := rpchttp.New(remote)
	if err != nil {
		return nil, err
	}

	events := `{
  "events": [
    {
      "id": "222a1f3a-47ec-4acb-949e-915d5a7ee889",
      "type": "RepositoryInitialized",
      "version": "1",
      "aggregate_id": "7df974e0-1282-42b4-8924-78712a6568e0",
      "aggregate_type": "repository"
    }
  ]
}`

	// Broadcast the transaction and wait for it to commit (rather use
	// c.BroadcastTxSync though in production).
	bres, err := abci.BroadcastTxCommit(context.Background(), []byte(events))
	if err != nil {
		return nil, err
	}

	if bres.CheckTx.IsErr() || bres.DeliverTx.IsErr() {
		return nil, err
	}

	return &Client{
		ABCIClient: abci,
	}, nil
}
