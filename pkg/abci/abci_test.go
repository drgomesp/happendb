package abci

import (
	"reflect"
	"testing"

	"github.com/tendermint/tendermint/abci/types"

	happendb "github.com/drgomesp/happendb/pkg"
	"github.com/drgomesp/happendb/pkg/store"
)

func TestNewApplication(t *testing.T) {
	memStore := store.NewMemory()
	expectedApp := &Application{
		BaseApplication: types.NewBaseApplication(),
		store:           memStore,
	}

	if got := NewApplication(memStore); !reflect.DeepEqual(got, expectedApp) {
		t.Errorf("NewApplication() = %v, want %v", got, expectedApp)
	}
}

func TestApplication_CheckTx(t *testing.T) {
	type fields struct {
		BaseApplication *types.BaseApplication
		store           happendb.EventStore
	}
	type args struct {
		req types.RequestCheckTx
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   types.ResponseCheckTx
	}{
		{},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Application{
				BaseApplication: tt.fields.BaseApplication,
				store:           tt.fields.store,
			}
			if got := a.CheckTx(tt.args.req); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CheckTx() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApplication_DeliverTx(t *testing.T) {
	type fields struct {
		BaseApplication *types.BaseApplication
		store           happendb.EventStore
	}
	type args struct {
		req types.RequestDeliverTx
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   types.ResponseDeliverTx
	}{
		{},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Application{
				BaseApplication: tt.fields.BaseApplication,
				store:           tt.fields.store,
			}
			if got := a.DeliverTx(tt.args.req); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeliverTx() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApplication_Query(t *testing.T) {
	type fields struct {
		BaseApplication *types.BaseApplication
		store           happendb.EventStore
	}
	type args struct {
		req types.RequestQuery
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   types.ResponseQuery
	}{
		{},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Application{
				BaseApplication: tt.fields.BaseApplication,
				store:           tt.fields.store,
			}
			if got := a.Query(tt.args.req); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Query() = %v, want %v", got, tt.want)
			}
		})
	}
}