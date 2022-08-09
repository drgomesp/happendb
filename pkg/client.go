package happendb

import (
	"context"
	"fmt"

	hexbytes "github.com/tendermint/tendermint/libs/bytes"
	"github.com/tendermint/tendermint/libs/json"
	"github.com/tendermint/tendermint/rpc/client"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
)

var _ Store = Client{}

type Client struct {
	abci  client.ABCIClient
	store Store
}

func NewClient(remote string) (*Client, error) {
	abci, err := rpchttp.New(remote)
	if err != nil {
		return nil, err
	}

	return &Client{
		abci: abci,
	}, nil
}

func (c Client) Save(ctx context.Context, events []*Event, version int) error {
	type EventsTx struct {
		Events []*Event `json:"events"`
	}

	data, err := json.Marshal(EventsTx{Events: events})
	if err != nil {
		return err
	}

	res, err := c.abci.BroadcastTxCommit(context.Background(), data)
	if err != nil {
		return err
	}

	if res.CheckTx.IsErr() || res.DeliverTx.IsErr() {
		return err
	}

	return nil
}

func (c Client) Load(ctx context.Context, t EventType) ([]*Event, error) {
	// Now try to fetch the value for the key
	res, err := c.abci.ABCIQuery(ctx, string(t), hexbytes.HexBytes{})
	if err != nil {
		return nil, err
	}
	if res.Response.IsErr() {
		return nil, err
	}

	fmt.Println("Got:", string(res.Response.Value))

	return nil, nil
}