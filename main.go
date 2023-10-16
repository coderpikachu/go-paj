package main

import (
	"context"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"paj/my_orm"
	"paj/my_web"
)

type MyTestModel struct {
	Id0       int64
	FirstName string `my_orm:"type=varchar(128)"`
	Age       int8
	LastName  *sql.NullString
}

func (MyTestModel) CreateSQL() string {
	return `
CREATE TABLE IF NOT EXISTS my_test_model(
    id0 INTEGER PRIMARY KEY,
    first_name TEXT NOT NULL,
    age INTEGER,
    last_name TEXT NOT NULL
)
`
}
func myOrmRemoteTest() string {
	db, _ := my_orm.MyDb("../my_test.db", my_orm.DBWithDialect(my_orm.SQLite3))
	_, err := db.Db.Exec("DROP TABLE IF EXISTS `my_test_model`")
	if err != nil {
		return "1*"
	}

	_, err = db.Db.Exec(MyTestModel{}.CreateSQL())
	if err != nil {
		return "2*"
	}
	_, err = db.Db.Exec("INSERT INTO `my_test_model`(`id0`,`first_name`,`age`,`last_name`)"+
		"VALUES (?,?,?,?)", 14, "Deng", 18, "Ming")
	if err != nil {
		return "3*"
	}
	res, err := my_orm.NewSelector[MyTestModel](db).Get(context.Background())
	if err != nil {
		return "4*"
	}
	return res.FirstName
}
func main() {
	println("hello paj!")

	//r := gin.Default()
	//r.GET("/test", func(c *gin.Context) {
	//	c.String(200, "Hello, Geektutu")
	//})
	////r.Run()
	//r.Run("0.0.0.0:8004") // listen and serve on 0.0.0.0:8080

	s := my_web.NewHTTPServer()

	r1 := myOrmRemoteTest()

	s.Get("/test", func(ctx *my_web.Context) {
		ctx.Resp.Write([]byte("hello, my world\n"))
		ctx.Resp.Write([]byte(r1))
	})
	s.Start(":8004")
}
