#Luna
Handle websocket connections and channel like rest api endpoints

##Why
Websocket servers are little bit tricky and not very easy to build, we love to fall back to easy alternative BaaS platform like Firebase to handle realtime processing in our apps.

Luna aims to simplify how you interact with websockets by letting you define channels as endpoints So it feels like you can actually send a message to a websocket channel by calling an endpoint.

###Using Luna
You can get up and running with few lines of code.

```
func main() {

	g := gin.Defaults()
	l := luna.New()

	g.GET("/ws/connect", func(c *gin.Context) {
		l.HandleHttpRequest(c.Writer, c.Request)
	})

	//OPTIONAL: Handle a message sent to a particular endpoint/channel
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

A full example with a Javascript client is in example folder

Thanks
Thanks to [@olahol](https://github.com/olahol/melody) for creating melody.

Follow me on Twitter [@L3kanAdigun](https://twitter.com/L3kanAdigun)
