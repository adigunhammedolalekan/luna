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

	sess2 := &melody.Session{}
	h.Subscribe(id, sess2)

	if h.Count() != 1 {
		t.Error("Expected subscribe count of channel to equals 1")
	}

	if h.ClientsCount(id) != 2 {
		t.Errorf("Exptected clients count to be 2, %d returned", h.ClientsCount(id))
	}
}

func TestChannel_UnSubscribe(t *testing.T) {

	sess := &melody.Session{}

	id := "/rooms/22/message"

	h := &Hub{
		Channels: make([]*Channel, 0),
	}

	h.Subscribe(id, sess)
	sess2 := &melody.Session{}
	h.Subscribe(id, sess2)

	if h.Count() != 1 {
		t.Error("Expected count of channel to equals 1")
	}

	if h.ClientsCount(id) != 2 {
		t.Errorf("Exptected clients count to be 2, %d returned", h.ClientsCount(id))
	}

	h.UnSubscribe(id, sess)

	if h.ClientsCount(id) != 1 {
		t.Errorf("Expected clients count to be 0, %d found", h.ClientsCount(id))
	}

	h.UnSubscribe(id, sess2)

	if h.ClientsCount(id) != 0 {
		t.Errorf("Expected clients count to be 0, %d found", h.ClientsCount(id))
	}
}

func TestHub_GetChannel(t *testing.T) {

	sess := &melody.Session{}
	id := "/rooms/22/message"

	h := &Hub{
		Channels: make([]*Channel, 0),
	}

	h.Subscribe(id, sess)

	ch := h.GetChannel(id)
	if ch == nil {
		t.Errorf("Expected ch to hold a value, nil is returned instead")
	}
}
