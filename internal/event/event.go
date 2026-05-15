package event

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

type Event struct {
	Time     time.Time
	PlayerID int
	EventID  int
	Extra    string
}

func ParseEvent(line string) (Event, error) {
	parts := strings.SplitN(line, "]", 2)

	if len(parts) != 2 {
		return Event{}, errors.New("invalid event format")
	}

	timePart := strings.TrimPrefix(parts[0], "[")
	t, err := time.Parse("15:04:05", timePart)
	if err != nil {
		return Event{}, err
	}
	fields := strings.Fields(parts[1])

	if len(fields) < 2 {
		return Event{}, errors.New("invalid event format")
	}

	playerID, err := strconv.Atoi(fields[0])
	if err != nil {
		return Event{}, err
	}

	eventID, err := strconv.Atoi(fields[1])
	if err != nil {
		return Event{}, err
	}

	event := Event{
		Time:     t,
		PlayerID: playerID,
		EventID:  eventID,
	}

	if len(fields) > 2 {
		event.Extra = strings.Join(fields[2:], " ")
	}

	return event, nil

}
