package sqlt

import "github.com/jmoiron/sqlx"

func MustConnect(dbname, dsname string) *sqlx.DB {
	return sqlx.MustConnect(dbname, dsname)
}
