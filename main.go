package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"paj/my_orm"
	"paj/my_web"
)

type MyTable struct {
	Id   int64
	Name string `my_orm:"type=varchar(45)"`
}

func (MyTable) CreateSQL() string {
	return `
CREATE TABLE IF NOT EXISTS my_table(
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL
)
`
}
func pingTest(c string) string {

	db, err := sql.Open("mysql", c)
	//尝试与数据库进行连接
	err = db.Ping()
	if err != nil {
		fmt.Println("数据库连接失败", err)
		return "01*"
	}
	return "02"
}
func myOrmRemoteTest(c string) string {

	db, _ := my_orm.MyMysql(c, my_orm.DBWithDialect(my_orm.SQLite3))

	_, err := db.Db.Exec("DROP TABLE IF EXISTS `my_table`")
	if err != nil {
		return "1*"
	}
	_, err = db.Db.Exec(MyTable{}.CreateSQL())
	if err != nil {
		return "2*"
	}
	//_, err = db.Db.Exec("INSERT INTO `my_table`(`id`,`name`)"+
	//	"VALUES (?,?)", 11, "2")

	my_orm.NewInserter[MyTable](db).Values(
		&MyTable{
			Id:   1,
			Name: "Deng",
		}).Columns("Id", "Name").Exec(context.Background())
	my_orm.NewInserter[MyTable](db).Values(
		&MyTable{
			Id:   2,
			Name: "Ming",
		}).Columns("Id", "Name").Exec(context.Background())
	if err != nil {
		return "3*"
	}

	ress, err := my_orm.NewSelector[MyTable](db).GetMyMulti(context.Background())
	if err != nil {
		return "4*"
	}
	return ress[0].Name + " " + ress[1].Name
}
func main() {
	println("hello paj!")

	s := my_web.NewHTTPServer()
	c := "root:SWB1436001@tcp(101.43.168.151:3306)/my_test_db"
	r0 := pingTest(c)
	r1 := myOrmRemoteTest(c)

	s.Get("/test", func(ctx *my_web.Context) {
		ctx.Resp.Write([]byte("test mysql \n"))
		ctx.Resp.Write([]byte(r0))
		ctx.Resp.Write([]byte("\n"))
		ctx.Resp.Write([]byte(r1))
		ctx.Resp.Write([]byte("\n"))
		ctx.Resp.Write([]byte(c))
	})
	s.Start(":8004")
}
