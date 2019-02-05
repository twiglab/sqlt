package main

import (
	"fmt"
	"strings"
)

const sqlString = `
SELECT a.attnum,
       a.attname AS field,
       t.typname AS type,
       a.attlen AS length,
       a.atttypmod AS lengthvar,
       a.atttypid as typid
  FROM pg_class c join pg_attribute a on(c.oid = a.attrelid)
       join pg_type t on (a.atttypid = t.oid)
 WHERE c.relname = '%s' and a.attnum > 0
 ORDER BY a.attnum;
 `

func sql(tablename string) (s string) {
	s = fmt.Sprintf(sqlString, tablename)
	return
}

type Field struct {
	Num     int    `db:"attnum"`
	Name    string `db:"field"`
	TypName string `db:"type"`
	Len     int    `db:"length"`
	LenVar  int    `db:"lengthvar"`
	TypID   int    `db:"typid"`
}

func (f *Field) SnakeName() string {
	return snake(f.Name)
}

func (f *Field) TypeString() string {
	return typeName(f)
}

type Table struct {
	TableName string
	Fields    []*Field
}

func (t *Table) SnakeName() string {
	return snake(t.TableName)
}

func snake(name string) string {
	newstr := make([]rune, 0)
	upNextChar := true

	name = strings.ToLower(name)

	for _, chr := range name {
		switch {
		case upNextChar:
			upNextChar = false
			if 'a' <= chr && chr <= 'z' {
				chr -= ('a' - 'A')
			}
		case chr == '_':
			upNextChar = true
			continue
		}

		newstr = append(newstr, chr)
	}

	return string(newstr)
}

func typeName(f *Field) string {
	switch f.TypID {
	case 1042: // bpchar
		return "string"
	case 1114: // timestamp
		return "time.Time"
	case 1043: // varchar
		return "string"
	case 23: // int4
		return "int"
	case 20: // int8
		return "int64"
	case 1700: // numeric
		return "float64"
	}

	panic(fmt.Errorf("%s, %s, %d", f.Name, f.TypName, f.TypID))
}
