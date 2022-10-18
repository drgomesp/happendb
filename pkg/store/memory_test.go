package store_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	happendb "github.com/drgomesp/happendb/pkg"
	"github.com/drgomesp/happendb/pkg/store"
)

var (
	e1 = happendb.NewEvent(
		"repository.RepositoryInitialized",
		"54e260be-26ce-451a-815d-b2a16e4f3cd0",
		1,
		"my-app",
	)

	e2 = happendb.NewEvent(
		"repository.RepositoryUpdated",
		"54e260be-26ce-451a-815d-b2a16e4f3cd0",
		2,
		"my-app",
	)

	e3 = happendb.NewEvent(
		"repository.RepositoryUpdated",
		"54e260be-26ce-451a-815d-b2a16e4f3cd0",
		3,
		"my-app",
	)

	e4 = happendb.NewEvent(
		"repository.RepositoryUpdated",
		"54e260be-26ce-451a-815d-b2a16e4f3cd0",
		4,
		"my-app",
	)
)

func TestMemoryStore(t *testing.T) {
	tests := []struct {
		name                       string
		events, existing, expected []*happendb.Event
		fromVersion                int
		expectedError              error
	}{
		{
			name:        "test save on empty Events",
			existing:    []*happendb.Event{},
			events:      []*happendb.Event{e1, e2},
			expected:    []*happendb.Event{e1, e2},
			fromVersion: 0,
		},
		{
			name:        "test save on empty Events from non-zero version",
			existing:    []*happendb.Event{},
			events:      []*happendb.Event{e1, e2},
			expected:    []*happendb.Event{e1, e2},
			fromVersion: 360,
		},
		{
			name:        "test save on non-empty Events",
			existing:    []*happendb.Event{e1, e2},
			events:      []*happendb.Event{e3, e4},
			expected:    []*happendb.Event{e1, e2, e3, e4},
			fromVersion: 3,
		}, {
			name:          "test save on non-empty Events from version zero",
			existing:      []*happendb.Event{e1, e2},
			events:        []*happendb.Event{e3},
			expectedError: happendb.ErrStoreInvalidVersion,
			fromVersion:   0,
		},
		{
			name:          "test save on non-empty Events from wrong version",
			existing:      []*happendb.Event{e1, e2},
			events:        []*happendb.Event{e3},
			expectedError: happendb.ErrStoreInvalidVersion,
			fromVersion:   1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			memoryStore := store.NewMemory()
			memoryStore.Events = map[string][]*happendb.Event{
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
