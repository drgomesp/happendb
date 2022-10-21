package core

type EventType string

type Event struct {
	// ID is a Version 4 UUID
	ID string `json:"id"`

	// Type is a general event type
	Type EventType `json:"type"`

	// Version starting from 0 (zero).
	Version *int `json:"version,string"`

	// AggregateID is a Version 4 UUID of the aggregate
	AggregateID string `json:"aggregate_id"`

	// AggregateType is the aggregate type
	AggregateType string `json:"aggregate_type"`
}

func NewEvent(t EventType, id string, version int, aggregateID string, aggregateType string) *Event {
	return &Event{
		ID:            id,
		Version:       &version,
		Type:          t,
		AggregateID:   aggregateID,
		AggregateType: aggregateType,
	}
}
