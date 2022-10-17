package happendb

import (
	"context"

	"github.com/davecgh/go-spew/spew"
	hexbytes "github.com/tendermint/tendermint/libs/bytes"
	"github.com/tendermint/tendermint/libs/json"
	"github.com/tendermint/tendermint/rpc/client"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
)

var _ EventStore = Client{}

type Client struct {
	abci  client.ABCIClient
	store EventStore
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

func (c Client) Save(ctx context.Context, events []*Event, fromVersion int) error {
	type EventsTx struct {
		Events []*Event `json:"events"`
	}

	data, err := json.Marshal(EventsTx{Events: events})
	if err != nil {
		return err
	}

	_ = fromVersion // ignore, since it will be enforced by the event store

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

	var events []*Event
	spew.Dump(string(res.Response.Value))
	if err = json.Unmarshal(res.Response.Value, &events); err != nil {
		return nil, err
	}

	return events, nil
}
