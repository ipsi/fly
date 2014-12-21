package eventstream

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/concourse/atc"
	"github.com/concourse/atc/event"
	"github.com/vito/go-sse/sse"
)

//go:generate counterfeiter . EventStream
type EventStream interface {
	NextEvent() (atc.Event, error)
}

type SSEEventStream struct {
	sseReader *sse.EventSource
}

func NewSSEEventStream(reader *sse.EventSource) *SSEEventStream {
	return &SSEEventStream{sseReader: reader}
}

func (s *SSEEventStream) NextEvent() (atc.Event, error) {
	se, err := s.sseReader.Next()
	if err != nil {
		return nil, err
	}

	switch se.Name {
	case "event":
		var message event.Message
		err = json.Unmarshal(se.Data, &message)
		if err != nil {
			return nil, err
		}

		return message.Event, nil

	case "end":
		return nil, io.EOF

	default:
		return nil, fmt.Errorf("unknown event name: %s", se.Name)
	}
}
