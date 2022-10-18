package abci

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/abci/example/code"
	"github.com/tendermint/tendermint/abci/types"

	"github.com/drgomesp/happendb/pkg/store"
)

func TestNewApplication(t *testing.T) {
	memStore := store.NewMemory()
	expectedApp := &Application{
		BaseApplication: types.NewBaseApplication(),
		store:           memStore,
	}

	if got := NewApplication(memStore); !reflect.DeepEqual(got, expectedApp) {
		t.Errorf("NewApplication() = %v, expectedResponse %v", got, expectedApp)
	}
}

func TestApplication_CheckTx(t *testing.T) {
	memStore := store.NewMemory()
	app := &Application{
		BaseApplication: types.NewBaseApplication(),
		store:           memStore,
	}

	type args struct {
		req types.RequestCheckTx
	}
	tests := []struct {
		name             string
		args             args
		expectedResponse types.ResponseCheckTx
		expectedError    error
	}{
		{
			name: "test send events",
			args: args{
				req: types.RequestCheckTx{
					Tx: []byte(`{
					  "events": [
						{
						  "id": "222a1f3a-47ec-4acb-949e-915d5a7ee889",
						  "type": "RepositoryCreated",
						  "version": "3",
						  "aggregate_id": "7df974e0-1282-42b4-8924-78712a6568e0",
						  "aggregate_type": "repository"
						},
						{
						  "id": "68af701f-0578-4132-ada1-6efa5a189666",
						  "type": "RepositoryUpdated",
						  "version": "4",
						  "aggregate_id": "8e31acbf-f554-4bca-ad42-1e03c64c09dd",
						  "aggregate_type": "repository"
						}
					  ]
					}`),
					Type: types.CheckTxType_New,
				},
			},
			expectedResponse: types.ResponseCheckTx{
				Code: types.CodeTypeOK,
			},
		},
		{
			name: "test invalid tx",
			args: args{
				req: types.RequestCheckTx{
					Tx:   []byte(`[notgood}`),
					Type: types.CheckTxType_New,
				},
			},
			expectedError: errors.New("oops"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := app.CheckTx(tt.args.req)

			if tt.expectedError != nil {
				assert.NotEqual(t, types.CodeTypeOK, got.Code)
			} else {
				assert.Equal(t, tt.expectedResponse, got)
			}
		})
	}
}

func TestApplication_DeliverTx(t *testing.T) {
	type args struct {
		req types.RequestDeliverTx
	}
	tests := []struct {
		name             string
		args             args
		expectedResponse types.ResponseDeliverTx
		expectedError    error
	}{
		{
			name: "test send events",
			args: args{
				req: types.RequestDeliverTx{
					Tx: []byte(`{
					  "events": [
						{
						  "id": "222a1f3a-47ec-4acb-949e-915d5a7ee889",
						  "type": "RepositoryCreated",
						  "version": "3",
						  "aggregate_id": "7df974e0-1282-42b4-8924-78712a6568e0",
						  "aggregate_type": "repository"
						},
						{
						  "id": "68af701f-0578-4132-ada1-6efa5a189666",
						  "type": "RepositoryUpdated",
						  "version": "4",
						  "aggregate_id": "8e31acbf-f554-4bca-ad42-1e03c64c09dd",
						  "aggregate_type": "repository"
						}
					  ]
					}`),
				},
			},
			expectedResponse: types.ResponseDeliverTx{
				Code: types.CodeTypeOK,
			},
		},
	}

	memStore := store.NewMemory()
	app := &Application{
		BaseApplication: types.NewBaseApplication(),
		store:           memStore,
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := app.DeliverTx(tt.args.req)

			if tt.expectedError != nil {
				assert.NotEqual(t, types.CodeTypeOK, got.Code)
			} else {
				assert.Equal(t, tt.expectedResponse, got)
			}
		})
	}
}

func TestApplication_Query(t *testing.T) {
	type args struct {
		req types.RequestQuery
	}
	tests := []struct {
		name             string
		args             args
		expectedResponse types.ResponseQuery
		expectedError    error
	}{
		{
			name: "test query events",
			args: args{
				req: types.RequestQuery{
					Data: []byte("54e260be-26ce-451a-815d-b2a16e4f3cd0"),
				},
			},
		},
	}

	memStore := store.NewMemory()
	app := &Application{
		BaseApplication: types.NewBaseApplication(),
		store:           memStore,
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := app.Query(tt.args.req)

			if tt.expectedError != nil {
				assert.NotEqual(t, types.CodeTypeOK, got.Code)
			} else {
				assert.Equal(t, tt.expectedResponse.Code, code.CodeTypeOK)
			}
		})
	}
}
