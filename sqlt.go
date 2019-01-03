package sqlt

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type Handler interface {
	HandleRows(*sqlx.Rows) error
}

var (
	DefaultTxOptions *sql.TxOptions = NewTxOptions(sql.LevelDefault, false)
)

func NewTxOptions(level sql.IsolationLevel, readonly bool) *sql.TxOptions {
	return &sql.TxOptions{Isolation: level, ReadOnly: readonly}
}

type sqlExecer interface {
	PrepareNamedContext(context.Context, string) (*sqlx.NamedStmt, error)
	MustSql(string, interface{}) string
}

func query(ctx context.Context, ext sqlExecer, id string, data interface{}, h Handler) error {
	sql := ext.MustSql(id, data)
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
	sql := ext.MustSql(id, data)
	stmt, err := ext.PrepareNamedContext(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	r, e = stmt.ExecContext(ctx, data)
	return
}

type (
	Dbop struct {
		*StdSqlAssembler
		*sqlx.DB
	}
)

func New(driverName, dataSourceName, pattern string) *Dbop {
	return &Dbop{
		DB:              sqlx.MustConnect(driverName, dataSourceName),
		StdSqlAssembler: NewStdSqlAssemblerDefault(pattern),
	}
}

func (c *Dbop) QueryContext(ctx context.Context, id string, data interface{}, h Handler) error {
	return query(ctx, c, id, data, h)
}

func (c *Dbop) ExecContext(ctx context.Context, id string, data interface{}) (r sql.Result, e error) {
	r, e = exec(ctx, c, id, data)
	return
}

func (c *Dbop) ExecRtnContext(ctx context.Context, id string, data interface{}, mrh Handler) error {
	return c.QueryContext(ctx, id, data, mrh)
}

func (c *Dbop) BeginTrans(ctx context.Context, opt *sql.TxOptions) (*Txop, error) {
	tx, err := c.BeginTxx(ctx, opt)
	if err != nil {
		return nil, err
	}

	return &Txop{
		Tx:              tx,
		StdSqlAssembler: c.StdSqlAssembler,
	}, nil
}

type Txop struct {
	*StdSqlAssembler
	*sqlx.Tx
}

func (t *Txop) QueryContext(ctx context.Context, id string, data interface{}, h Handler) error {
	return query(ctx, t, id, data, h)
}

func (t *Txop) ExecContext(ctx context.Context, id string, data interface{}) (r sql.Result, e error) {
	r, e = exec(ctx, t, id, data)
	return
}

func (t *Txop) ExecRtnContext(ctx context.Context, id string, data interface{}, h Handler) error {
	return t.QueryContext(ctx, id, data, h)
}

func (t *Txop) CommitTrans() error {
	return t.Commit()
}

func (t *Txop) RollbackTrans() error {
	return t.Rollback()
}
