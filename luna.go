package luna

import (
	"encoding/json"
	"fmt"
	"github.com/olahol/melody"
	"net/http"
)

var (
	COMMAND_SUBSCRIBE   = "subscribe"
	COMMAND_MESSAGE     = "message"
	COMMAND_UNSUBSCRIBE = "unsubscribe"
)

type Luna struct {
	melody *melody.Melody
	hub    *Hub
	routes []*Route
}

type WsMessage struct {
	Path   string      `json:"path"`
	Action string      `json:"action"`
	Data   interface{} `json:"data"`
}

// Config holds luna configuration data
type Config struct {

	BufferSize int
	MaxMessageSize int64
}

var DefaultConfig = &Config{BufferSize: 512 * 10, MaxMessageSize: 512 * 10}

//New creates a new Luna instance
func New(config *Config) *Luna {

	m := melody.New()
	h := &Hub{
		Channels: make([]*Channel, 0),
	}

	if config == nil {
		config = DefaultConfig
	}

	m.Config.MessageBufferSize = config.BufferSize
	m.Config.MaxMessageSize = config.MaxMessageSize

	go h.EnsureClean()

	luna := &Luna{
		melody: m,
		hub:    h,
		routes: make([]*Route, 0),
	}

	go luna.handleMessages()
	return luna
}

//Handle registers a new Route
func (l *Luna) Handle(path string, f OnMessageHandler) {

	route := &Route{}
	route.Path = path
	route.OnNewMessage = f
	l.routes = append(l.routes, route)
}

func (l *Luna) HandleHttpRequest(wr http.ResponseWriter, req *http.Request) error {

	keys := make(map[string] interface{})
	keys["session_token"] = req.Header.Get("Authorization")
	return l.melody.HandleRequestWithKeys(wr, req, keys)
}

//handleMessages starts to listen for new websocket events on a seperate goroutine
func (l *Luna) handleMessages() {

	l.melody.HandleMessage(func(session *melody.Session, bytes []byte) {

		message := &WsMessage{}
		err := json.Unmarshal(bytes, message)
		if err != nil {
			fmt.Println("Error while creating object from json payload ", string(bytes), err)
			return
		}

		switch message.Action {

		case COMMAND_SUBSCRIBE:
			l.hub.Subscribe(message.Path, session)
		case COMMAND_UNSUBSCRIBE:
			l.hub.UnSubscribe(message.Path, session)
		case COMMAND_MESSAGE:
			l.hub.Send(message.Path, message.Data)
			for _, route := range l.routes {

				fmt.Println("Route => ", route.Path)
				fmt.Println("Message Route => ", message.Path)
				if MatchRoute(route.Path, message.Path) {

					fmt.Println("Match!")
					if route.OnNewMessage != nil {
						ctx := &Context{}
						ctx.Path = message.Path
						ctx.Vars, _ = ExtractParams(route.Path, message.Path)
						ctx.Data = message.Data

						fmt.Println("Called")
						route.OnNewMessage(ctx)
					}
				}
			}
		}
	})
}

// Publish sends @param data to channel @param channel
func (l *Luna) Publish(channel string, data interface{}) error {
	return l.hub.Send(channel, data)
}

//
type OnMessageHandler func(context *Context)

type Route struct {
	Path         string
	OnNewMessage OnMessageHandler
}
