package sqlt

import (
	"context"
	"database/sql"
	"text/template"

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

type (
	Dbop struct {
		Maker
		*sqlx.DB
	}
)

func New(driverName, dataSourceName, pattern string) *Dbop {
	db := sqlx.MustConnect(driverName, dataSourceName)
	return NewWithDB(db, pattern)
}

func NewWithDB(db *sqlx.DB, pattern string) *Dbop {
	return &Dbop{
		DB:    db,
		Maker: NewSqlTemplate(pattern, make(template.FuncMap)),
	}
}

func (c *Dbop) Exec(ctx context.Context, id string, data interface{}) (r sql.Result, e error) {
	r, e = exec(ctx, c, id, data)
	return
}

func (c *Dbop) ExecRtn(ctx context.Context, id string, data interface{}, h Handler) error {
	return query(ctx, c, id, data, h)
}

func (c *Dbop) BeginTrans(ctx context.Context, opt *sql.TxOptions) (*Txop, error) {
	tx, err := c.BeginTxx(ctx, opt)
	if err != nil {
		return nil, err
	}

	return &Txop{
		Tx:    tx,
		Maker: c.Maker,
	}, nil
}

type Txop struct {
	Maker
	*sqlx.Tx
}

func (t *Txop) ExecContext(ctx context.Context, id string, data interface{}) (r sql.Result, e error) {
	r, e = exec(ctx, t, id, data)
	return
}

func (t *Txop) ExecRtnContext(ctx context.Context, id string, data interface{}, h Handler) error {
	return query(ctx, t, id, data, h)
}

func (t *Txop) CommitTrans() error {
	return t.Commit()
}

func (t *Txop) RollbackTrans() error {
	return t.Rollback()
}
