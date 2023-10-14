package main

import "github.com/gin-gonic/gin"

func main() {
	println("hello paj!")

	r := gin.Default()
	r.GET("/test", func(c *gin.Context) {
		c.String(200, "Hello, Geektutu")
	})
	//r.Run()
	r.Run("0.0.0.0:8004") // listen and serve on 0.0.0.0:8080
}
