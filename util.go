package sqlt

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

var DefaultTxOptions *sql.TxOptions = NewTxOptions(sql.LevelDefault, false)

func NewTxOptions(level sql.IsolationLevel, readonly bool) *sql.TxOptions {
	return &sql.TxOptions{Isolation: level, ReadOnly: readonly}
}

type Param = map[string]interface{}

var null = make(Param)

func dummy(p interface{}) interface{} {
	if p == nil {
		return null
	}
	return p
}

type sqltExecer interface {
	PrepareNamedContext(context.Context, string) (*sqlx.NamedStmt, error)
	Maker
}

func query(ctx context.Context, ext sqltExecer, id string, data interface{}, h RowsExtractor) error {
	param := dummy(data)

	sql := MustSql(ext, id, param)
	stmt, err := ext.PrepareNamedContext(ctx, sql)
	if err != nil {
		return err
	}
	defer stmt.Close()
	rows, err := stmt.QueryxContext(ctx, param)
	if err != nil {
		return err
	}
	defer rows.Close()
	return h.Extract(rows)
}

func exec(ctx context.Context, ext sqltExecer, id string, data interface{}) (r sql.Result, e error) {
	param := dummy(data)

	sql := MustSql(ext, id, param)
	stmt, err := ext.PrepareNamedContext(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	r, e = stmt.ExecContext(ctx, param)
	return
}

type TExecer interface {
	TQuery(context.Context, string, interface{}, RowsExtractor) error
	TExec(context.Context, string, interface{}) (sql.Result, error)
	TExecRtn(context.Context, string, interface{}, RowsExtractor) error
}

func Query(execer TExecer, ctx context.Context, id string, param interface{}, h RowsExtractor) (err error) {
	err = execer.TQuery(ctx, id, param, h)
	return
}

func MustQuery(execer TExecer, ctx context.Context, id string, param interface{}, h RowsExtractor) {
	if err := Query(execer, ctx, id, param, h); err != nil {
		panic(err)
	}
}

func Exec(execer TExecer, ctx context.Context, id string, param interface{}) (r sql.Result, err error) {
	r, err = execer.TExec(ctx, id, param)
	return
}

func MustExec(execer TExecer, ctx context.Context, id string, param interface{}) (r sql.Result) {
	var err error
	if r, err = Exec(execer, ctx, id, param); err != nil {
		panic(err)
	}
	return
}

func ExecRtn(execer TExecer, ctx context.Context, id string, param interface{}, h RowsExtractor) (err error) {
	err = execer.TExecRtn(ctx, id, param, h)
	return
}

func MustExecRtn(execer TExecer, ctx context.Context, id string, param interface{}, h RowsExtractor) {
	if err := ExecRtn(execer, ctx, id, param, h); err != nil {
		panic(err)
	}
}

type TxBegin interface {
	TBegin(context.Context, *sql.TxOptions) (*Txop, error)
}

func Begin(b TxBegin, ctx context.Context, opt *sql.TxOptions) (*Txop, error) {
	return b.TBegin(ctx, opt)
}

type TxEnd interface {
	TCommit() error
	TRollback() error
}

func Commit(t TxEnd) error {
	return t.TCommit()
}

func Rollback(t TxEnd) error {
	return t.TRollback()
}
