package main

import (
	"paj/my_web"
)

func main() {
	println("hello paj!")

	//r := gin.Default()
	//r.GET("/test", func(c *gin.Context) {
	//	c.String(200, "Hello, Geektutu")
	//})
	////r.Run()
	//r.Run("0.0.0.0:8004") // listen and serve on 0.0.0.0:8080

	s := my_web.NewHTTPServer()
	s.Get("/test", func(ctx *my_web.Context) {
		ctx.Resp.Write([]byte("hello, my world"))
	})

	s.Start(":8004")
}
