# Luna
Handle websocket connections and events like rest api endpoints

## Why
Websocket servers are little bit tricky and not very easy to build, we normally fall back to easy alternative BaaS platform like Firebase to handle realtime data processing in our apps.
`Luna wants to make creating Websocket servers in Golang something anyone can do.`


### How To
It only takes some few steps to use Luna.

Add luna to your project
`go get github.com/adigunhammedolalekan/luna`

1. create Luna instance

```Go 

l := luna.New()

```

2.

Allow luna to handle websocket connection requests
```Go
//ws://localhost:port/ws/connect

//Using Go default http stack
http.Handle("/ws/connect", func(w http.ResponseWriter, r *http.Request) {
	//This is a websocket connection, allow luna to handle it
	l.HandleHttpRequest(w, r)
})

//Using Gin
g := gin.Default()
g.GET("/ws/connect", func(context *gin.Context) {
	l.HandleHttpRequest(context.Writer, context.Request)
})
```

3.

Create json data payload from clients(Js frontend, Android, iOS etc) and use any Websocket client lib to send it
```
//Example payload
{
	"action" : "subscribe",
	"path" : "/rooms/22/message"
}
/*This payload simply subscribe the WS client that sends it to path/channel 
`/room/22/message`, after subscription, all messages received by `/rooms/22/message`
 would be sent to this client and every other clients that subscribed
*/

//Example message payload
payload = {
	"action" : "message",
	"path" : "/rooms/22/message",
	"data" : {
		"sent_by" : "L3kan",
		"text" : "Hey, there. Luna is awesome",
		"time" : "12:30"
	}
}

//This payload would send payload.data to every clients that has subscribed to `/rooms/22/message` channel/path. Amazing!
```

As you can see, you can set up a websocket server with 2 lines of code using Luna.
```Go

func main() {

	g := gin.Default()
	l := luna.New() //create luna

	g.GET("/ws/connect", func(c *gin.Context) {
		l.HandleHttpRequest(c.Writer, c.Request)
	})

	//OPTIONAL: Handle a message (payload.data) sent to a particular channel/path
	l.Handle("/rooms/{id}/message", func(c *luna.Context) {
		//Message has been sent to appropriate Websocket clients.
		//do other stuffs here, like saving message into a persistence layer?

		//Db.Save(c.Data)
		vars := c.Vars //Grab path parameters
		fmt.Println("Id => ", vars["id"] . (string))
		fmt.Println("Got message from path => " +  c.Path)
		fmt.Println("Data => ", c.Data)
	})
}
```

A full example with a Javascript client can be found in example folder

![Console](https://github.com/adigunhammedolalekan/luna/blob/master/art/console.PNG)

![Web](https://github.com/adigunhammedolalekan/luna/blob/master/art/web.PNG)

## TODO

* Add support for FCM - Firebase Cloud Messaging
* Write more tests

Thanks to [@olahol](https://github.com/olahol) for creating [melody](https://github.com/olahol/melody).

Follow me on Twitter [@L3kanAdigun](https://twitter.com/L3kanAdigun), my DM is open, come ask me questions!
