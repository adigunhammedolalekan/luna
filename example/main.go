package main

import (
	"fmt"
	"github.com/adigunhammedolalekan/luna"
	"github.com/gin-gonic/gin"
)

func main() {

	l := luna.New()
	g := gin.Default()

	g.LoadHTMLFiles("example/files/index.html")
	g.Static("/assets", "./example/assets")

	g.GET("/ws/connect", func(context *gin.Context) {
		l.HandleHttpRequest(context.Writer, context.Request)
	})

	g.GET("/home", func(context *gin.Context) {
		context.HTML(200, "index.html", nil)
	})

	l.Handle("/rooms/{id}/message", func(c *luna.Context) {
		//Message has been sent to appropriate Websocket clients.
		//do other stuffs here, like saving message into a persistence layer?

		//Db.Save(c.Data)
		vars := c.Vars //Grab path parameters
		fmt.Println("Id => ", vars["id"] . (string))
		fmt.Println("Got message from path => " +  c.Path)
		fmt.Println("Data => ", c.Data)
	})

	g.Run("0.0.0.0:8009")
}


