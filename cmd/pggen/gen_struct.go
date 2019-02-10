package main

import "fmt"

func genStruct(t *Table) {
	fmt.Printf("type %s struct {\n", t.SnakeName())

	for _, f := range t.Fields {
		fmt.Printf("%s %s `db:\"%s\"` // \n", f.SnakeName(), f.TypeString(), f.Name)
	}

	fmt.Printf("}\n")
}
