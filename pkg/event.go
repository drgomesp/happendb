package happendb

type EventType string

type Event struct {
	AggregateID   string `json:"aggregate_id"`
	AggregateType string `json:"aggregate_type"`
	ID            string `json:"id"`
	Type          string `json:"type"`
	Version       *int   `json:"version,string"`
}

func NewEvent(t EventType, id string, version int, aggregateID, aggregateType string) *Event {
	return &Event{
		AggregateID:   aggregateID,
		AggregateType: aggregateType,
		ID:            id,
		Type:          string(t),
		Version:       &version,
	}
}