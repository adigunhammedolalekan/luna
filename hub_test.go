package luna

import (
	"github.com/olahol/melody"
	"testing"
)

func TestChannel_Subscribe(t *testing.T) {

	sess := &melody.Session{}
	id := "/rooms/22/message"

	h := &Hub{
		Channels: make([]*Channel, 0),
	}

	h.Subscribe(id, sess)

	if h.Count() != 1 {
		t.Error("Expected subscribe count of channel to equals 1")
	}
}
