package happendb

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type clientMock struct {
	mock.Mock
}

func (m *clientMock) Save(ctx context.Context, events []*Event, fromVersion int) error {
	args := m.Called(ctx, events, fromVersion)
	return args.Error(0)

}

func (m *clientMock) Load(ctx context.Context, eventType EventType) ([]*Event, error) {
	args := m.Called(ctx, eventType)
	return args.Get(0).([]*Event), args.Error(1)
}

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
			ctx := context.Background()

			client := new(clientMock)
			client.On("Save", ctx, tt.events, tt.fromVersion).Return(nil)

			err := client.Save(ctx, tt.events, tt.fromVersion)

			assert.NoError(t, err)
			assert.NotNil(t, client)
		})
	}
}

func TestClient_Load(t *testing.T) {
	tests := []struct {
		name           string
		eventType      EventType
		expectedEvents []*Event
	}{
		{
			name:      "load",
			eventType: EventType("repository.RepositoryInitialized"),
			expectedEvents: []*Event{
				NewEvent(
					"repository.RepositoryInitialized",
					"54e260be-26ce-451a-815d-b2a16e4f3cd0",
					1,
					"my-app",
				),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			client := new(clientMock)
			client.On("Load", ctx, tt.eventType).Return(tt.expectedEvents, nil)

			events, err := client.Load(ctx, tt.eventType)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedEvents, events)
		})
	}
}
