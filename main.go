package main

import (
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang/mock/mockgen/model"
	_ "github.com/mattn/go-sqlite3"
	"paj/my_cache"
	"paj/my_web"
)

func main() {
	println("hello paj!")

	s := my_web.NewHTTPServer()
	//orm
	//c := "root:SWB1436001@tcp(101.43.168.151:3306)/my_test_db"
	//r0 := my_orm.PingTest(c)
	//r1 := my_orm.MyOrmRemoteTest(c)

	r2 := my_cache.MyCacheTest()

	s.Get("/test", func(ctx *my_web.Context) {
		ctx.Resp.Write([]byte("test mysql \n"))
		ctx.Resp.Write([]byte(r2))
		ctx.Resp.Write([]byte("\n"))
	})
	s.Start(":8004")
}
