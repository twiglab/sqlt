package sqlt

import (
	"context"
	"database/sql"
	"text/template"

	"github.com/jmoiron/sqlx"
)

type ExtractFunc func(*sqlx.Rows) error

type Dbop struct {
	Maker
	*sqlx.DB
}

func New(db *sqlx.DB, pattern string) *Dbop {
	return &Dbop{
		DB:    db,
		Maker: NewSqlTemplate(pattern, make(template.FuncMap)),
	}
}

func (c *Dbop) TQuery(ctx context.Context, id string, param interface{}, h ExtractFunc) error {
	return query(ctx, c, id, param, h)
}

func (c *Dbop) TExec(ctx context.Context, id string, param interface{}) (r sql.Result, e error) {
	r, e = exec(ctx, c, id, param)
	return
}

func (c *Dbop) TExecRtn(ctx context.Context, id string, param interface{}, h ExtractFunc) error {
	return query(ctx, c, id, param, h)
}

func (c *Dbop) TBegin(ctx context.Context, opt *sql.TxOptions) (*Txop, error) {
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

func (t *Txop) TQuery(ctx context.Context, id string, param interface{}, h ExtractFunc) error {
	return query(ctx, t, id, param, h)
}

func (t *Txop) TExec(ctx context.Context, id string, param interface{}) (r sql.Result, e error) {
	r, e = exec(ctx, t, id, param)
	return
}

func (t *Txop) TExecRtn(ctx context.Context, id string, param interface{}, h ExtractFunc) error {
	return query(ctx, t, id, param, h)
}

func (t *Txop) TCommit() error {
	return t.Commit()
}

func (t *Txop) TRollback() error {
	return t.Rollback()
}
