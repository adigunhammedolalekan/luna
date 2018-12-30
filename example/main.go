package main

import (
	"fmt"
	"github.com/adigunhammedolalekan/luna"
	"github.com/gin-gonic/gin"
)

func main() {

	l := luna.New(nil)
	g := gin.Default()

	g.LoadHTMLFiles("example/files/index.html")
	g.Static("/assets", "./example/assets")

	g.GET("/ws/connect", func(context *gin.Context) {
		l.HandleHttpRequest(context.Writer, context.Request)
	})

	g.GET("/home", func(context *gin.Context) {
		context.HTML(200, "index.html", nil)
	})

	l.Handle("/rooms/{id}/message", func(context *luna.Context) {

		fmt.Println("Path Data => ", context.Vars)
		fmt.Println("Message => ", context.Data)

		data := context.Data . (map[string] interface{})
		fmt.Println(data["text"])
	})

	g.Run("0.0.0.0:8009")
}