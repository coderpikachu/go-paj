package my_orm

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"log"
	"paj/my_orm/internal/errs"
	"paj/my_orm/internal/valuer"
	"paj/my_orm/model"
	"time"
)

type DBOption func(*DB)

type DB struct {
	core
	Db *sql.DB
}

// Wait 会等待数据库连接
// 注意只能用于测试
func (db *DB) Wait() error {
	err := db.Db.Ping()
	for err == driver.ErrBadConn {
		log.Printf("等待数据库启动...")
		err = db.Db.Ping()
		time.Sleep(time.Second)
	}
	return err
}

// Open 创建一个 DB 实例。
// 默认情况下，该 DB 将使用 MySQL 作为方言
// 如果你使用了其它数据库，可以使用 DBWithDialect 指定
func MyDb(dbStr string, opts ...DBOption) (*DB, error) {
	db, err := sql.Open("sqlite3", dbStr)
	if err != nil {
		return nil, err
	}
	return OpenDB(db, opts...)
}
func Open(driver string, dsn string, opts ...DBOption) (*DB, error) {
	//db, err := sql.Open(driver, dsn)
	db, err := sql.Open("sqlite3", "./my_test.db")
	if err != nil {
		return nil, err
	}
	return OpenDB(db, opts...)
}

func OpenDB(db *sql.DB, opts ...DBOption) (*DB, error) {
	res := &DB{
		core: core{
			dialect:    MySQL,
			r:          model.NewRegistry(),
			ValCreator: valuer.NewUnsafeValue,
		},
		Db: db,
	}
	for _, opt := range opts {
		opt(res)
	}
	return res, nil
}

func DBWithDialect(dialect Dialect) DBOption {
	return func(db *DB) {
		db.dialect = dialect
	}
}

func DBWithRegistry(r model.Registry) DBOption {
	return func(db *DB) {
		db.r = r
	}
}

func DBUseReflectValuer() DBOption {
	return func(db *DB) {
		db.ValCreator = valuer.NewReflectValue
	}
}

func DBWithMiddleware(ms ...Middleware) DBOption {
	return func(db *DB) {
		db.ms = ms
	}
}

// MustNewDB 创建一个 DB，如果失败则会 panic
// 我个人不太喜欢这种
func MustNewDB(driver string, dsn string, opts ...DBOption) *DB {
	db, err := Open(driver, dsn, opts...)
	if err != nil {
		panic(err)
	}
	return db
}

// BeginTx 开启事务
func (db *DB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	tx, err := db.Db.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}
	return &Tx{tx: tx, db: db}, nil
}

// type txKey struct {
//
// }

// BeginTxV2 事务扩散
// 个人不太喜欢
// func (db *DB) BeginTxV2(ctx context.Context,
// 	opts *sql.TxOptions) (context.Context, *Tx, error) {
// 	val := ctx.Value(txKey{})
// 	if val != nil {
// 		tx := val.(*Tx)
// 		if !tx.done {
// 			return ctx, tx, nil
// 		}
// 	}
// 	tx, err := db.BeginTx(ctx, opts)
// 	if err != nil {
// 		return ctx, nil, err
// 	}
// 	ctx = context.WithValue(ctx, txKey{}, tx)
// 	return ctx, tx, nil
// }

// DoTx 将会开启事务执行 fn。如果 fn 返回错误或者发生 panic，事务将会回滚，
// 否则提交事务
func (db *DB) DoTx(ctx context.Context,
	fn func(ctx context.Context, tx *Tx) error,
	opts *sql.TxOptions) (err error) {
	var tx *Tx
	tx, err = db.BeginTx(ctx, opts)
	if err != nil {
		return err
	}

	panicked := true
	defer func() {
		if panicked || err != nil {
			e := tx.Rollback()
			if e != nil {
				err = errs.NewErrFailToRollbackTx(err, e, panicked)
			}
		} else {
			err = tx.Commit()
		}
	}()

	err = fn(ctx, tx)
	panicked = false
	return err
}

func (db *DB) Close() error {
	return db.Db.Close()
}

func (db *DB) getCore() core {
	return db.core
}

func (db *DB) queryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return db.Db.QueryContext(ctx, query, args...)
}

func (db *DB) execContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return db.Db.ExecContext(ctx, query, args...)
}
