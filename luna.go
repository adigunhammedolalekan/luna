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
	config *Config
}

type WsMessage struct {
	Path   string      `json:"path"`
	Action string      `json:"action"`
	Data   interface{} `json:"data"`
}

// ExtractKeyFunc to assign a unique key for each session
// can be used like :-
//
// keyFunc := func(r *http.Request) string {
// 		return r.Header.Get("x-auth-token")
// }
type ExtractKeyFunc func(*http.Request) string

// Config holds luna configurations
type Config struct {
	BufferSize     int
	MaxMessageSize int64
	KeyExtractor   ExtractKeyFunc
}

var DefaultConfig = &Config{BufferSize: 512 * 10, MaxMessageSize: 512 * 10}

// New creates a new Luna instance
func New(config *Config) *Luna {

	m := melody.New()
	h := &Hub{
		Channels: make([]*Channel, 0),
	}

	if config == nil {
		panic("initialize luna with a config")
	}

	if config.KeyExtractor == nil {
		panic("KeyExtractor function is nil")
	}

	if config.MaxMessageSize == 0 {
		config.MaxMessageSize = DefaultConfig.MaxMessageSize
	}

	if config.BufferSize == 0 {
		config.BufferSize = DefaultConfig.BufferSize
	}

	m.Config.MessageBufferSize = config.BufferSize
	m.Config.MaxMessageSize = config.MaxMessageSize

	go h.EnsureClean()

	luna := &Luna{
		melody: m,
		hub:    h,
		routes: make([]*Route, 0),
		config: config,
	}

	go luna.handleMessages()
	return luna
}

// Handle registers a new Route
func (l *Luna) Handle(path string, f OnMessageHandler) {

	route := &Route{}
	route.Path = path
	route.OnNewMessage = f
	l.routes = append(l.routes, route)
}

func (l *Luna) HandleHttpRequest(wr http.ResponseWriter, req *http.Request) error {

	keys := make(map[string]interface{})
	keys["session_token"] = l.config.KeyExtractor(req)
	return l.melody.HandleRequestWithKeys(wr, req, keys)
}

// handleMessages starts to listen for new websocket events on a seperate goroutine
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

				if MatchRoute(route.Path, message.Path) {

					if route.OnNewMessage != nil {
						ctx := &Context{}
						ctx.Path = message.Path
						ctx.Vars, _ = ExtractParams(route.Path, message.Path)

						bytes, _ := json.Marshal(message.Data)
						ctx.Data = bytes
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
