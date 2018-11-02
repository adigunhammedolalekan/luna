package main

import (
	"fmt"
	"github.com/adigunhammedolalekan/luna"
	"github.com/gin-gonic/gin"
)

func main() {

	l := luna.New()
	g := gin.Default()
	g.GET("/ws/connect", func(context *gin.Context) {
		l.HandleHttpRequest(context.Writer, context.Request)
	})

	g.GET("/home", func(context *gin.Context) {
		context.JSON(200, gin.H{"message" : "Hello"})
	})

	l.Handle("/rooms/{id}/message", func(c *luna.Context) {
		//Message has been sent to appropriate Websocket clients.
		//do other stuffs here, like saving message into a persistence layer?

		//Db.Save(c.Data)
		vars := c.Vars //Grab path parameters
		fmt.Println(vars["id"] . (string))
		fmt.Println("Got message from path => " +  c.Path)
	})

	g.Run("0.0.0.0:8009")
}


