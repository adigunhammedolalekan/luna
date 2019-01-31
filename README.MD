# Luna [![Made in Nigeria](https://img.shields.io/badge/made%20in-nigeria-008751.svg?style=flat-square)](https://github.com/acekyd/made-in-nigeria)
Create webSocket servers in a minute

## Why
Websocket servers are little bit tricky and not very easy to build, we normally fall back to easy alternative BaaS platform like Firebase or services like Pusher to handle realtime data processing in our apps.
`Luna aim to make creating WebSocket servers in Golang something anyone can do.`


### How To
It only takes some few steps to use Luna.

Add luna to your project
`go get github.com/adigunhammedolalekan/luna`

1. create Luna instance

```Go 
    g := gin.Default()

    // keyExtractor is a function that should return a
    // unique 'string'. It takes a *http.Request pointer.
    // The intention is to extract a unique Id from request's
    // pointer to set as a unique session key for every
    // webSocket request

	keyExtractor := func(req *http.Request) string {
    		return req.Header.Get("Authorization")
    }

    config := &luna.Config{
        KeyExtractor: keyExtractor,
    }

    l := luna.New(config)

```

2.

Allow luna to handle websocket connection requests
``` Go
// Example URL - ws://localhost:port/ws/connect

// Using Go default http stack
http.Handle("/ws/connect", func(w http.ResponseWriter, r *http.Request) {
	// This is a websocket connection, allow luna to handle it
	l.HandleHttpRequest(w, r)
})

// Using Gin
g := gin.Default()
g.GET("/ws/connect", func(context *gin.Context) {
	l.HandleHttpRequest(context.Writer, context.Request)
})
```

3. Create json data payload from clients(Js frontend, Android, iOS etc) and use any Websocket client lib to send it
```
// Example payload
{
	"action" : "subscribe",
	"path" : "/rooms/22/message"
}
// This payload simply subscribe the sender's websocket client to path/channel
// `/rooms/22/message`, after subscription, all messages received by `/rooms/22/message`
// would be sent to this client and every other client that subscribed

// Example message payload
messagePayload = {
	"action" : "message",
	"path" : "/rooms/22/message",
	"data" : {
		"sent_by" : "L3kan",
		"text" : "Hey, there. Luna is awesome",
		"time" : "12:30"
	}
}

```
    **action - can be any of the three option(subscribe, message, unsubscribe), what you use depends on what you want to do
    **path - path is where you are sending the message to
    **data - is the actual message you are sending

```
// This payload would send `payload.data` to every client that has subscribed to `/rooms/22/message` channel/path. Amazing!

To stop receiving messages on the client side(i.e UnSubscribe from a channel)
{
    "action" : "unsubscribe",
    "path" : "/rooms/22/message"
}

Although, Luna has a feature that automatically removes idle websocket clients.
```

As you can see, you can set up a websocket server with few lines of code using Luna.
```Go

func main() {

	g := gin.Default()

	// setup keyExtractor
	keyExtractor := func(req *http.Request) string {
    		return req.Header.Get("Authorization")
    }

    // create config
    config := &luna.Config{
        KeyExtractor: keyExtractor,
    }

    // start websocket server
	l := luna.New(config)

    // handle websocket connection requests
	g.GET("/ws/connect", func(c *gin.Context) {
		l.HandleHttpRequest(c.Writer, c.Request)
	})

	// OPTIONAL: Handle a message `payload.data` sent to a particular channel/path
	l.Handle("/rooms/{id}/message", func(c *luna.Context) {
		// Message has been sent to appropriate Websocket clients.
		// do other stuffs here, like saving message into a persistence layer?

        m := &Message{}
        err := json.UnMarshall(c.Data, m)
        if err == nil {
		    Db.Save(m)

		    // you can also use the extracted path
		    // parameters
		    vars := c.Vars
            fmt.Println("Id => ", vars["id"] . (string))
            fmt.Println("Got message from path => " +  c.Path)
        }
	})
}
```

A full example with a Javascript client can be found in example folder

![Console](https://github.com/adigunhammedolalekan/luna/blob/master/art/console.PNG)

![Web](https://github.com/adigunhammedolalekan/luna/blob/master/art/web.PNG)

## TODO

* Add support for FCM - Firebase Cloud Messaging
* Write more tests

I built a small wrapper lib for android. [luna-android](https://github.com/adigunhammedolalekan/luna-android)

Thanks to [@olahol](https://github.com/olahol) for creating [melody](https://github.com/olahol/melody).

Follow me on Twitter [@L3kanAdigun](https://twitter.com/L3kanAdigun) if you have any question.
