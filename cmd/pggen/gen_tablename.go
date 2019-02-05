package main

import "fmt"

func genTabelName(t *Table) {
	fmt.Printf("const TableName = \"%s\"\n\n", t.SnakeName())
}
