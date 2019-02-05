package main

import "fmt"

func genStruct(t *Table) {
	fmt.Printf("type %s struct {\n", t.SnakeName())

	for _, f := range t.Fields {
		fmt.Printf("%s %s `sqlt:\"%s\" json:\"%s\"` // \n", f.SnakeName(), f.TypeString(), f.Name, f.SnakeName())
	}

	fmt.Printf("}\n")
}
