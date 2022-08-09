package store

import (
	"context"

	happendb "github.com/drgomesp/happendb/pkg"
)

var _ happendb.Store = Memory{}

type Memory struct {
	store map[string][]*happendb.Event
}

func NewMemory() *Memory {
	return &Memory{
		store: make(map[string][]*happendb.Event, 0),
	}
}

func (m Memory) Save(ctx context.Context, events []*happendb.Event, fromVersion int) error {
	if len(events) == 0 {
		return happendb.ErrStoreMissingEvents
	}

	head := events[0]
	id := string(head.Type)
	empty := len(m.store[id]) == 0

	if !empty && fromVersion == 0 {
		return happendb.ErrStoreInvalidVersion
	}

	pending := make([]*happendb.Event, len(events))

	for i := 0; i < len(events); i++ {
		e := events[i]

		if len(m.store[id]) > 0 {
			if fromVersion != 0 && *e.Version != fromVersion+i {
				return happendb.ErrStoreInvalidVersion
			}
		}

		pending[i] = e
	}

	m.store[id] = append(m.store[id], pending...)

	return nil
}

func (m Memory) Load(ctx context.Context, t happendb.EventType) ([]*happendb.Event, error) {
	events, ok := m.store[string(t)]
	if ok && len(events) > 0 {
		return events, nil
	}

	return nil, nil
}