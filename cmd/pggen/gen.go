package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var db *sqlx.DB

const tableName = "req_order"

type Origin struct {
	DataBaseName   string `json:"data_base_name"`
	DataSourceName string `json:"data_source_name"`
	TableName      string `json:"table_name"`
}

//"dbname=testdb sslmode=disable")

var FileName string = "origin.json"

func init() {
	flag.StringVar(&FileName, "origin", "origin.json", "param")
}

func main() {

	flag.Parse()

	bs, err := ioutil.ReadFile(FileName)
	if err != nil {
		panic(err)
	}

	origin := new(Origin)
	origin.DataSourceName = "dbname=testdb sslmode=disable"

	if err := json.Unmarshal(bs, origin); err != nil {
		panic(err)
	}

	db = sqlx.MustConnect("postgres", origin.DataSourceName)

	r, err := db.Queryx(sql(origin.TableName))
	if err != nil {
		panic(err)
	}

	table := new(Table)
	table.TableName = tableName

	for r.Next() {
		f := new(Field)
		r.StructScan(f)
		table.Fields = append(table.Fields, f)
	}

	/*
		for _, f := range table.Fields {
			fmt.Println(f)
		}
	*/
	genTabelName(table)
	genConst(table)
	genStruct(table)

	genInsert(table)
	genSelect(table)
	genUpdate(table)
}
