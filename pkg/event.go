package happendb

type EventType string

type Event struct {
	AggregateRef

	ID      string `json:"id"`
	Type    string `json:"type"`
	Version *int   `json:"version,string"`
}

func NewEvent(t EventType, id string, version int, agg AggregateRef) *Event {
	return &Event{
		AggregateRef: agg,
		ID:           id,
		Type:         string(t),
		Version:      &version,
	}
}
