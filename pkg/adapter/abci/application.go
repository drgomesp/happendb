package abci

import (
	"context"

	"github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/json"

	"github.com/drgomesp/happendb/pkg/core"
)

type TxEvents struct {
	Events []*core.Event `json:"events"`
}

const (
	_ = iota

	ErrCodeBadRequest
	ErrCodeSaveFailed
	ErrCodeLoadFailed
)

type Application struct {
	*types.BaseApplication

	store core.EventStore
}

func NewApplication(store core.EventStore) *Application {
	return &Application{
		BaseApplication: types.NewBaseApplication(),
		store:           store,
	}
}

func (a *Application) CheckTx(req types.RequestCheckTx) types.ResponseCheckTx {
	if _, err := a.parseTxEvents(req.GetTx()); err != nil {
		return types.ResponseCheckTx{
			Code: ErrCodeBadRequest,
			Log:  err.Error(),
		}
	}

	return types.ResponseCheckTx{
		Code: types.CodeTypeOK,
	}
}

func (a *Application) parseTxEvents(data []byte) ([]*core.Event, error) {
	var tx TxEvents
	err := json.Unmarshal(data, &tx)

	if err != nil && len(tx.Events) == 0 {
		return nil, err
	}

	return tx.Events, nil
}

func (a *Application) DeliverTx(req types.RequestDeliverTx) types.ResponseDeliverTx {
	ctx := context.Background()

	events, err := a.parseTxEvents(req.GetTx())
	if err != nil {
		return types.ResponseDeliverTx{
			Code: ErrCodeBadRequest,
			Log:  err.Error(),
		}
	}

	head := events[0]
	previous, err := a.store.Load(ctx, head.ID)
	if err != nil {
		return types.ResponseDeliverTx{
			Code: ErrCodeLoadFailed,
			Log:  err.Error(),
		}
	}

	fromVersion := 0
	if len(previous) > 0 {
		tail := previous[len(previous)-1]
		fromVersion = *tail.Version
	}

	err = a.store.Save(ctx, events, fromVersion)
	if err != nil {
		return types.ResponseDeliverTx{
			Code: ErrCodeSaveFailed,
			Log:  err.Error(),
		}
	}

	return types.ResponseDeliverTx{Code: types.CodeTypeOK}
}

func (a *Application) Query(req types.RequestQuery) types.ResponseQuery {
	id := req.GetData()

	ctx := context.Background()
	events, err := a.store.Load(ctx, string(id))

	if err != nil {
		return types.ResponseQuery{Code: ErrCodeLoadFailed}
	}

	data, err := json.Marshal(events)
	if err != nil {
		return types.ResponseQuery{Code: ErrCodeLoadFailed}
	}

	return types.ResponseQuery{
		Code:  types.CodeTypeOK,
		Value: data,
	}
}
