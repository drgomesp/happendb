package happendb

import (
	"context"
	"errors"
)

var (
	ErrStoreMissingEvents  = errors.New("no events")
	ErrStoreInvalidVersion = errors.New("event version invalid")
)

type Store interface {
	Save(ctx context.Context, events []*Event, version int) error
	Load(ctx context.Context, t EventType) ([]*Event, error)
}