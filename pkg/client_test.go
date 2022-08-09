package happendb

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_Save(t *testing.T) {
	tests := []struct {
		name        string
		events      []*Event
		fromVersion int
	}{
		{
			name: "save",
			events: []*Event{
				NewEvent(
					"repository.RepositoryInitialized",
					"54e260be-26ce-451a-815d-b2a16e4f3cd0",
					1,
					"my-app",
				),
			},
			fromVersion: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient("http://localhost:26657")
			assert.NoError(t, err)
			assert.NotNil(t, client)

			ctx := context.Background()
			err = client.Save(ctx, tt.events, tt.fromVersion)

			assert.NoError(t, err)
		})
	}
}

func TestClient_Load(t *testing.T) {
	tests := []struct {
		name           string
		t              EventType
		expectedEvents []*Event
	}{
		{
			name:           "load",
			t:              EventType("repository.RepositoryInitialized"),
			expectedEvents: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient("http://localhost:26657")
			assert.NoError(t, err)
			assert.NotNil(t, client)

			ctx := context.Background()
			events, err := client.Load(ctx, tt.t)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedEvents, events)
		})
	}
}