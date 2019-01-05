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

type sqltExecer interface {
	PrepareNamedContext(context.Context, string) (*sqlx.NamedStmt, error)
	Maker
}

func query(ctx context.Context, ext sqltExecer, id string, data interface{}, h ExtractFunc) error {
	sql := MustSql(ext, id, data)
	stmt, err := ext.PrepareNamedContext(ctx, sql)
	if err != nil {
		return err
	}
	defer stmt.Close()
	rows, err := stmt.QueryxContext(ctx, data)
	if err != nil {
		return err
	}
	defer rows.Close()
	return h(rows)
}

func exec(ctx context.Context, ext sqltExecer, id string, data interface{}) (r sql.Result, e error) {
	sql := MustSql(ext, id, data)
	stmt, err := ext.PrepareNamedContext(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	r, e = stmt.ExecContext(ctx, data)
	return
}

type TExecer interface {
	TQuery(context.Context, string, interface{}, ExtractFunc) error
	TExec(context.Context, string, interface{}) (sql.Result, error)
	TExecRtn(context.Context, string, interface{}, ExtractFunc) error
}

func Query(execer TExecer, ctx context.Context, id string, param interface{}, h ExtractFunc) (err error) {
	err = execer.TQuery(ctx, id, param, h)
	return
}

func Exec(execer TExecer, ctx context.Context, id string, param interface{}) (r sql.Result, err error) {
	r, err = execer.TExec(ctx, id, param)
	return
}

func ExecRtn(execer TExecer, ctx context.Context, id string, param interface{}, h ExtractFunc) (err error) {
	err = execer.TExecRtn(ctx, id, param, h)
	return
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
