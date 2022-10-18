package store_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/drgomesp/happendb/pkg/adapter/store"
	"github.com/drgomesp/happendb/pkg/core"
)

var (
	e1 = core.NewEvent(
		"repository.RepositoryInitialized",
		"54e260be-26ce-451a-815d-b2a16e4f3cd0",
		1,
		"my-app",
	)

	e2 = core.NewEvent(
		"repository.RepositoryUpdated",
		"54e260be-26ce-451a-815d-b2a16e4f3cd0",
		2,
		"my-app",
	)

	e3 = core.NewEvent(
		"repository.RepositoryUpdated",
		"54e260be-26ce-451a-815d-b2a16e4f3cd0",
		3,
		"my-app",
	)

	e4 = core.NewEvent(
		"repository.RepositoryUpdated",
		"54e260be-26ce-451a-815d-b2a16e4f3cd0",
		4,
		"my-app",
	)
)

func TestMemoryStore(t *testing.T) {
	tests := []struct {
		name                       string
		events, existing, expected []*core.Event
		fromVersion                int
		expectedError              error
	}{
		{
			name:        "test save on empty Events",
			existing:    []*core.Event{},
			events:      []*core.Event{e1, e2},
			expected:    []*core.Event{e1, e2},
			fromVersion: 0,
		},
		{
			name:        "test save on empty Events from non-zero version",
			existing:    []*core.Event{},
			events:      []*core.Event{e1, e2},
			expected:    []*core.Event{e1, e2},
			fromVersion: 360,
		},
		{
			name:        "test save on non-empty Events",
			existing:    []*core.Event{e1, e2},
			events:      []*core.Event{e3, e4},
			expected:    []*core.Event{e1, e2, e3, e4},
			fromVersion: 3,
		}, {
			name:          "test save on non-empty Events from version zero",
			existing:      []*core.Event{e1, e2},
			events:        []*core.Event{e3},
			expectedError: core.ErrStoreInvalidVersion,
			fromVersion:   0,
		},
		{
			name:          "test save on non-empty Events from wrong version",
			existing:      []*core.Event{e1, e2},
			events:        []*core.Event{e3},
			expectedError: core.ErrStoreInvalidVersion,
			fromVersion:   1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			memoryStore := store.NewMemory()
			memoryStore.Events = map[string][]*core.Event{
				"repository": tt.existing,
			}

			if len(tt.existing) > 0 {
				assert.NoError(t, memoryStore.Save(ctx, tt.existing, 1))
			}

			err := memoryStore.Save(ctx, tt.events, tt.fromVersion)

			if tt.expectedError == nil {
				loaded, loadErr := memoryStore.Load(ctx, tt.events[0].ID)
				assert.NoError(t, loadErr)
				assert.Equal(t, tt.expected, loaded)
			} else {
				assert.Equal(t, tt.expectedError, err)
			}
		})
	}
}
