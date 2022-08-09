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
		"107342db-0a17-4314-aa8a-2120842a7645",
		2,
		"my-app",
	)

	e3 = happendb.NewEvent(
		"repository.RepositoryUpdated",
		"a6271955-3b95-4f16-a00e-8e21e57ac106",
		3,
		"my-app",
	)

	e4 = happendb.NewEvent(
		"repository.RepositoryUpdated",
		"e3d88938-d492-4d3d-b108-f24ea77ce4dd",
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
			name:        "test save on empty store",
			existing:    []*happendb.Event{},
			events:      []*happendb.Event{e1, e2},
			expected:    []*happendb.Event{e1, e2},
			fromVersion: 0,
		},
		{
			name:        "test save on empty store from non-zero version",
			existing:    []*happendb.Event{},
			events:      []*happendb.Event{e1, e2},
			expected:    []*happendb.Event{e1, e2},
			fromVersion: 360,
		},
		{
			name:        "test save on non-empty store",
			existing:    []*happendb.Event{e1, e2},
			events:      []*happendb.Event{e3, e4},
			expected:    []*happendb.Event{e1, e2, e3, e4},
			fromVersion: 3,
		}, {
			name:          "test save on non-empty store from version zero",
			existing:      []*happendb.Event{e1, e2},
			events:        []*happendb.Event{e3},
			expectedError: happendb.ErrStoreInvalidVersion,
			fromVersion:   0,
		},
		{
			name:          "test save on non-empty store from wrong version",
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

			if len(tt.existing) > 0 {
				assert.NoError(t, memoryStore.Save(ctx, tt.existing, 0))
			}

			err := memoryStore.Save(ctx, tt.events, tt.fromVersion)

			if tt.expectedError == nil {
				loaded, loadErr := memoryStore.Load(ctx, tt.events[0].Type)
				assert.NoError(t, loadErr)
				assert.Equal(t, tt.expected, loaded)
			} else {
				assert.Equal(t, tt.expectedError, err)
			}

		})
	}
}