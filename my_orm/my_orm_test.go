package my_orm

import (
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	_ "paj/my_orm/internal/errs"
	"paj/my_orm/internal/valuer"
	"testing"
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
func sendRemote() {

}
func Test_My_Orm(t *testing.T) {
	db, _ := MyDb("../my_test.db", DBWithDialect(SQLite3))

	db.Db.Exec("DROP TABLE IF EXISTS  `my_test_model`")

	_, err := db.Db.Exec(MyTestModel{}.CreateSQL())
	if err != nil {
		t.Fatal(err)
	}

	res, err := db.Db.Exec("INSERT INTO `my_test_model`(`id0`,`first_name`,`age`,`last_name`)"+
		"VALUES (?,?,?,?)", 14, "1", 18, "Ming")
	if err != nil {
		t.Fatal(err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		t.Fatal(err)
	}
	if affected == 0 {
		t.Fatal()
	}
	t.Run("unsafe", func(t *testing.T) {
		db.ValCreator = valuer.NewUnsafeValue
		_, err = NewSelector[MyTestModel](db).Get(context.Background())
		if err != nil {
			t.Fatal(err)
		}
	})
}
func Test_My_Orm_Insert_Build(t *testing.T) {
	db := memoryDB(t)
	testCases := []struct {
		name      string
		q         QueryBuilder
		wantQuery *Query
		wantErr   error
	}{
		{
			name: "single values",
			q: NewInserter[TestModel](db).Values(
				&TestModel{
					Id:        1,
					FirstName: "Deng",
					Age:       18,
					LastName:  &sql.NullString{String: "Ming", Valid: true},
				}),
			wantQuery: &Query{
				SQL:  "INSERT INTO `test_model`(`id`,`first_name`,`age`,`last_name`) VALUES(?,?,?,?);",
				Args: []any{int64(1), "Deng", int8(18), &sql.NullString{String: "Ming", Valid: true}},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			query, err := tc.q.Build()
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantQuery, query)
		})
	}
}
func Test_My_Orm_Insert(t *testing.T) {
	db := memoryDB(t)
	testCases := []struct {
		name      string
		q         QueryBuilder
		wantQuery *Query
		wantErr   error
	}{
		{
			name: "single values",
			q: NewInserter[TestModel](db).Values(
				&TestModel{
					Id:        1,
					FirstName: "Deng",
					Age:       18,
					LastName:  &sql.NullString{String: "Ming", Valid: true},
				}),
			wantQuery: &Query{
				SQL:  "INSERT INTO `test_model`(`id`,`first_name`,`age`,`last_name`) VALUES(?,?,?,?);",
				Args: []any{int64(1), "Deng", int8(18), &sql.NullString{String: "Ming", Valid: true}},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			query, err := tc.q.Build()
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantQuery, query)
		})
	}
}
func Test_My_Orm_Get(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = mockDB.Close() }()
	db, err := OpenDB(mockDB)
	if err != nil {
		t.Fatal(err)
	}

	//q := NewInserter[TestModel](db).Values(
	//	&TestModel{
	//		Id:        1,
	//		FirstName: "Deng",
	//		Age:       18,
	//		LastName:  &sql.NullString{String: "Ming", Valid: true},
	//	})
	//query, err := q.Build()

	if err != nil {
		return
	}

	testCases := []struct {
		name     string
		query    string
		mockErr  error
		mockRows *sqlmock.Rows
		wantErr  error
		wantVal  *TestModel
	}{
		{
			name:  "get data",
			query: "SELECT .*",
			mockRows: func() *sqlmock.Rows {
				res := sqlmock.NewRows([]string{"id", "first_name", "age", "last_name"})
				res.AddRow([]byte("1"), []byte("Da"), []byte("18"), []byte("Ming"))
				return res
			}(),
			wantVal: &TestModel{
				Id:        1,
				FirstName: "Da",
				Age:       18,
				LastName:  &sql.NullString{String: "Ming", Valid: true},
			},
		},
	}
	for _, tc := range testCases {
		exp := mock.ExpectQuery(tc.query)
		if tc.mockErr != nil {
			exp.WillReturnError(tc.mockErr)
		} else {
			exp.WillReturnRows(tc.mockRows)
		}
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := NewSelector[TestModel](db).Get(context.Background())
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantVal, res)
		})
	}
}
