package my_orm

import (
	"context"
	"database/sql"
	"fmt"
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
func PingTest(c string) string {

	db, err := sql.Open("mysql", c)
	//尝试与数据库进行连接
	err = db.Ping()
	if err != nil {
		fmt.Println("数据库连接失败", err)
		return "01*"
	}
	return "02"
}
func MyOrmRemoteTest(c string) string {

	db, _ := MyMysql(c, DBWithDialect(SQLite3))

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

	NewInserter[MyTable](db).Values(
		&MyTable{
			Id:   1,
			Name: "Deng",
		}).Columns("Id", "Name").Exec(context.Background())
	NewInserter[MyTable](db).Values(
		&MyTable{
			Id:   2,
			Name: "Ming",
		}).Columns("Id", "Name").Exec(context.Background())
	if err != nil {
		return "3*"
	}

	ress, err := NewSelector[MyTable](db).GetMyMulti(context.Background())
	if err != nil {
		return "4*"
	}
	return ress[0].Name + " " + ress[1].Name
}
