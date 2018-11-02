package luna

import (
	"encoding/json"
	"fmt"
	"github.com/olahol/melody"
	"net/http"
)

var (
	COMMAND_SUBSCRIBE = "subscribe"
	COMMAND_MESSAGE = "message"
)

type Luna struct {
	melody *melody.Melody
	hub *Hub
	routes []*Route
}

type WsMessage struct {
	Path string `json:"path"`
	Action string `json:"action"`
	Data interface{} `json:"data"`
}

func New() *Luna {

	m := melody.New()
	h := &Hub{
		Channels: make([]*Channel, 0),
	}

	go h.EnsureClean()

	luna := &Luna{
		melody: m,
		hub: h,
		routes: make([]*Route, 0),
	}

	go luna.handleMessages()
	return luna
}

func (l *Luna) Handle(path string, f OnMessageHandler)  {

	route := &Route{}
	route.Path = path
	route.OnMessage = f
	l.routes = append(l.routes, route)
}

func (l *Luna) HandleHttpRequest(wr http.ResponseWriter, req *http.Request)  {

	l.melody.HandleRequest(wr, req)
}

func (l *Luna) handleMessages() {

	l.melody.HandleMessage(func(session *melody.Session, bytes []byte) {

		fmt.Println("New Message => ", string(bytes))
		message := &WsMessage{}
		err := json.Unmarshal(bytes, message)
		if err != nil {
			fmt.Println("Error while creating object from json payload ", string(bytes), err)
			return
		}

		switch message.Action {

		case COMMAND_SUBSCRIBE:
			l.hub.Subscribe(message.Path, session)
		case COMMAND_MESSAGE:
			l.hub.Send(message.Path, message.Data)
			for _, route := range l.routes {

				if Match(route.Path, message.Path) {

					if route.OnMessage != nil {
						ctx := &Context{}
						ctx.Path = message.Path
						ctx.Vars, _ = ExtractParams(route.Path, message.Path)
						ctx.Data = message.Data

						route.OnMessage(ctx)
					}
				}
			}
		}
	})
}

type OnMessageHandler func(context *Context)

type Route struct {
	Path string
	OnMessage OnMessageHandler
}