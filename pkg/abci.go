package happendb

import (
	"context"

	"github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/json"
)

type TxEvents struct {
	Events []*Event `json:"events"`
}

const (
	_ = iota

	ErrCodeBadRequest
	ErrCodeSaveFailed
	ErrCodeLoadFailed
)

type Application struct {
	*types.BaseApplication

	store Store
}

func NewApplication(store Store) *Application {
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

func (a *Application) parseTxEvents(data []byte) ([]*Event, error) {
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
	previous, err := a.store.Load(ctx, head.AggregateID)
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
	data := string(req.GetData())
	_ = data

	return types.ResponseQuery{Code: types.CodeTypeOK}
}
