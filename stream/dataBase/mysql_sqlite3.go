package dataBase

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"

	"github.com/ddkwork/golibrary/mylog"
)

type (
	Interface interface {
		Init(driverName, dataSourceName string)
		CreatTables(DDL string)
		Query(query string)
		QueryResult() any
		Update(query string, args ...any)
		Insert(query string, args ...any)
	}
	object struct {
		db          *sql.DB
		stmt        *sql.Stmt
		rows        *sql.Rows
		result      sql.Result
		queryResult any
	}
)

func New() Interface { return &object{} }

func (o *object) Init(driverName, dataSourceName string) {
	db := mylog.Check2(sql.Open(driverName, dataSourceName))
	o.db = db
	o.db.SetMaxOpenConns(1000)
	o.db.SetMaxIdleConns(30000)
	mylog.Check(o.db.Ping())
}

func (o *object) CreatTables(DDL string) { mylog.Check2(o.db.Exec(DDL)) }
func (o *object) QueryResult() any       { return o.queryResult }
func (o *object) Query(query string) {
	rows := mylog.Check2(o.db.Query(query))
	o.rows = rows
	mylog.CheckNil(o.rows)
	defer func() {
		mylog.Check(o.rows.Close())
	}()
	for o.rows.Next() {
		mylog.Check(o.rows.Scan(&o.queryResult))
	}
}

func (o *object) Update(query string, args ...any) { o.stmtExec(query, args) }
func (o *object) Insert(query string, args ...any) { o.stmtExec(query, args) }
func (o *object) stmtExec(query string, args ...any) {
	stmt := mylog.Check2(o.db.Prepare(query))
	o.stmt = stmt
	defer func() {
		mylog.Check(o.stmt == nil)
		mylog.Check(o.stmt.Close())
	}()
	o.result = mylog.Check2(o.stmt.Exec(args))
}
