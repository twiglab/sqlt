package main

import "fmt"

func genSelect(t *Table) {
	fmt.Printf("{{ define \"%s.select\"}}\n", t.SnakeName())

	fmt.Printf("select\n")

	for _, f := range t.Fields {
		fmt.Printf("\t%s,\n", f.Name)
	}

	fmt.Printf("from\n\t%s\nwhere\n", t.SnakeName())
	for _, f := range t.Fields {
		fmt.Printf("\t{{if .%s}} %s = :%s and {{end}}\n", f.SnakeName(), f.Name, f.SnakeName())
	}

	fmt.Printf("{{end}}\n")
}

func genInsert(t *Table) {
	fmt.Printf("{{ define \"%s.insert\"}}\n", t.SnakeName())
	fmt.Printf("insert into %s (\n", t.SnakeName())
	for _, f := range t.Fields {
		fmt.Printf("\t{{if .%s}} ,%s {{end}}\n", f.SnakeName(), f.Name)
	}
	fmt.Printf(") values (\n")
	for _, f := range t.Fields {
		fmt.Printf("\t{{if .%s}} ,:%s {{end}}\n", f.SnakeName(), f.SnakeName())
	}
	fmt.Printf(")\n")
	fmt.Printf("{{end}}\n")
}

func genUpdate(t *Table) {
	fmt.Printf("{{ define \"%s.update\"}}\n", t.SnakeName())
	fmt.Printf("update\n\t%s\nset\n", t.SnakeName())
	for _, f := range t.Fields {
		fmt.Printf("\t{{if .%s}} %s = :%s, {{end}}\n", f.SnakeName(), f.Name, f.SnakeName())
	}
	fmt.Printf("nwhere\n", t.SnakeName())
	for _, f := range t.Fields {
		fmt.Printf("\t{{if .%s}} %s = :%s and {{end}}\n", f.SnakeName(), f.Name, f.SnakeName())
	}

	fmt.Printf("{{end}}\n")
}
