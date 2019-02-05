package main

import "fmt"

func genConst(t *Table) {
	fmt.Println("const (")
	for _, f := range t.Fields {
		fmt.Printf("%s = \"%s\"\n", f.SnakeName(), f.Name)
	}
	fmt.Println(")")
}
