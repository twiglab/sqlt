package sqlt

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type sqlExecer interface {
	PrepareNamedContext(context.Context, string) (*sqlx.NamedStmt, error)
	Maker
}

func query(ctx context.Context, ext sqlExecer, id string, data interface{}, h Handler) error {
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
	return h.HandleRows(rows)
}

func exec(ctx context.Context, ext sqlExecer, id string, data interface{}) (r sql.Result, e error) {
	sql := MustSql(ext, id, data)
	stmt, err := ext.PrepareNamedContext(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	r, e = stmt.ExecContext(ctx, data)
	return
}

type Execer interface {
	Exec(context.Context, string, interface{}) (sql.Result, error)
	ExecRtn(context.Context, string, interface{}, Handler) error
}

func Exec(execer Execer, ctx context.Context, id string, param interface{}) (r sql.Result, err error) {
	r, err = execer.Exec(ctx, id, param)
	return
}

func ExecRtn(execer Execer, ctx context.Context, id string, param interface{}, h Handler) (err error) {
	err = execer.ExecRtn(ctx, id, param, h)
	return
}

type TxBegin interface {
	BeginTrans(context.Context, *sql.TxOptions) (*Txop, error)
}

func Begin(b TxBegin, ctx context.Context, opt *sql.TxOptions) (*Txop, error) {
	return b.BeginTrans(ctx, opt)
}

type TxEnd interface {
	CommitTrans() error
	Rollback() error
}

func Commit(t TxEnd) error {
	return t.CommitTrans()
}

func Rollback(t TxEnd) error {
	return t.CommitTrans()
}
