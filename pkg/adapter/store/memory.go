package store

import (
	"context"

	"github.com/drgomesp/happendb/pkg/core"
)

var _ core.EventStore = Memory{}

type Memory struct {
	Events map[string][]*core.Event
}

func NewMemory() *Memory {
	return &Memory{
		Events: make(map[string][]*core.Event, 0),
	}
}

func (m Memory) Save(ctx context.Context, events []*core.Event, fromVersion int) error {
	if len(events) == 0 {
		return core.ErrStoreMissingEvents
	}

	head := events[0]
	id := string(head.ID)
	empty := len(m.Events[id]) == 0

	if !empty && fromVersion == 0 {
		return core.ErrStoreInvalidVersion
	}

	pending := make([]*core.Event, len(events))

	for i := 0; i < len(events); i++ {
		e := events[i]

		if len(m.Events[id]) > 0 {
			if fromVersion != 0 && *e.Version != fromVersion+i {
				return core.ErrStoreInvalidVersion
			}
		}

		pending[i] = e
	}

	m.Events[id] = append(m.Events[id], pending...)

	return nil
}

func (m Memory) Load(ctx context.Context, id string) ([]*core.Event, error) {
	events, ok := m.Events[id]
	if ok && len(events) > 0 {
		return events, nil
	}

	return nil, nil
}
