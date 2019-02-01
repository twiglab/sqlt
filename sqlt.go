package sqlt

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type Rows interface {
	Next() bool
	NextResultSet() bool
	Err() error

	Scan(...interface{}) error
	MapScan(map[string]interface{}) error
	StructScan(interface{}) error

	ColumnTypes() ([]*sql.ColumnType, error)
	Columns() ([]string, error)
}

type ExtractFunc func(Rows) error

func (e ExtractFunc) Extract(rs Rows) error {
	return e(rs)
}

type RowsExtractor interface {
	Extract(Rows) error
}

type Dbop struct {
	Maker
	*sqlx.DB
}

func Default(dbname, dbsource, pattern string) *Dbop {
	dbx := sqlx.MustConnect(dbname, dbsource)
	maker := NewSqlTemplate(pattern)
	return New(dbx, maker)
}

func New(db *sqlx.DB, maker Maker) *Dbop {
	return &Dbop{
		DB:    db,
		Maker: maker,
	}
}

func (c *Dbop) TQuery(ctx context.Context, id string, param interface{}, h RowsExtractor) error {
	return query(ctx, c, id, param, h)
}

func (c *Dbop) TExec(ctx context.Context, id string, param interface{}) (r sql.Result, e error) {
	r, e = exec(ctx, c, id, param)
	return
}

func (c *Dbop) TExecRtn(ctx context.Context, id string, param interface{}, h RowsExtractor) error {
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

func (t *Txop) TQuery(ctx context.Context, id string, param interface{}, h RowsExtractor) error {
	return query(ctx, t, id, param, h)
}

func (t *Txop) TExec(ctx context.Context, id string, param interface{}) (r sql.Result, e error) {
	r, e = exec(ctx, t, id, param)
	return
}

func (t *Txop) TExecRtn(ctx context.Context, id string, param interface{}, h RowsExtractor) error {
	return query(ctx, t, id, param, h)
}

func (t *Txop) TCommit() error {
	return t.Commit()
}

func (t *Txop) TRollback() error {
	return t.Rollback()
}
