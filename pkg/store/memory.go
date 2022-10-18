package store

import (
	"context"

	happendb "github.com/drgomesp/happendb/pkg"
)

var _ happendb.EventStore = Memory{}

type Memory struct {
	Events map[string][]*happendb.Event
}

func NewMemory() *Memory {
	return &Memory{
		Events: make(map[string][]*happendb.Event, 0),
	}
}

func (m Memory) Save(ctx context.Context, events []*happendb.Event, fromVersion int) error {
	if len(events) == 0 {
		return happendb.ErrStoreMissingEvents
	}

	head := events[0]
	id := string(head.ID)
	empty := len(m.Events[id]) == 0

	if !empty && fromVersion == 0 {
		return happendb.ErrStoreInvalidVersion
	}

	pending := make([]*happendb.Event, len(events))

	for i := 0; i < len(events); i++ {
		e := events[i]

		if len(m.Events[id]) > 0 {
			if fromVersion != 0 && *e.Version != fromVersion+i {
				return happendb.ErrStoreInvalidVersion
			}
		}

		pending[i] = e
	}

	m.Events[id] = append(m.Events[id], pending...)

	return nil
}

func (m Memory) Load(ctx context.Context, id string) ([]*happendb.Event, error) {
	events, ok := m.Events[id]
	if ok && len(events) > 0 {
		return events, nil
	}

	return nil, nil
}
