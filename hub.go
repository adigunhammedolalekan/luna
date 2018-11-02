package luna

import (
	"encoding/json"
	"fmt"
	"github.com/olahol/melody"
	"sync"
	"time"
)

//Hub holds all channels in a slice. Clients can subscribe and send message to channels through the Hub
type Hub struct {

	Channels []*Channel //Slice of channels in this Hub
}

//Subscribe a client to a channel. Create channel if it does not exists before
func (h *Hub) Subscribe(id string, session *melody.Session) *Channel {

	ch := h.GetChannel(id)
	if ch == nil { //create channel if this is the first time or if channel does not not exists
		ch := &Channel{}
		ch.Id = id
		ch.Mtx = &sync.Mutex{}

		h.Channels = append(h.Channels, ch)
	}

	ch.Subscribe(session)
	return ch
}

//returns a Channel identified by @param id
func (h *Hub) GetChannel(id string) *Channel {

	for _, ch := range h.Channels {

		if ch != nil && ch.Id == id {
			return ch
		}
	}

	return nil
}

//Send data payload to a channel. This broadcast @params data to all connected clients
func (h *Hub) Send(channel string, data interface{})  {

	ch := h.GetChannel(channel)
	if ch != nil {
		ch.Send(data)
	}
}

//Keep clients slice clean. Remove all clients that has been idle for more than 10minutes
func (h *Hub) EnsureClean()  {

	for {
		time.Sleep(10 * time.Minute) //Run every 10min

		for _, ch := range h.Channels {

			for i, v := range ch.Clients {

				if t := time.Now().Sub(v.LastSeen); t.Minutes() >= 10 && v.Session.IsClosed() {

					ch.Clients[i] = nil
				}
			}
		}
	}
}

//returns number of connect channels
func (h *Hub) Count() int {
	return len(h.Channels)
}

type Client struct {

	Session *melody.Session
	LastSeen time.Time
}

//A message channel. Multiple clients can subscribe to a channel, message will be broadcast to all clients
//anytime Channel.Send(data) is called
type Channel struct {

	Id string
	Mtx *sync.Mutex
	Clients []*Client
}

func (ch *Channel) Lock()  {
	ch.Mtx.Lock()
}

func (ch *Channel) UnLock()  {
	ch.Mtx.Unlock()
}

//Subscribe a session to a client.
func (ch *Channel) Subscribe(session *melody.Session)  {

	ch.Lock()
	defer ch.UnLock()

	client := &Client{}
	client.Session = session
	client.LastSeen = time.Now()
	ch.Clients = append(ch.Clients, client)
}

//Broadcast message to all connected client and update the last activity time
func (ch *Channel) Send(data interface{}) {

	value, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error while mashalling JSON ", err)
		return
	}

	ch.Lock()
	defer ch.UnLock()
	for _, v := range ch.Clients {

		if v != nil && !v.Session.IsClosed() {

			v.Session.Write(value)
			v.LastSeen = time.Now()
		}
	}
}