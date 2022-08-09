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
					EventType("RepositoryInitialized"),
					"54e260be-26ce-451a-815d-b2a16e4f3cd0",
					1,
					"3aa25321-1ca3-4b00-8aee-d73e311383b2",
					"repository",
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