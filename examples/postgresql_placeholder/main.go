package main

import (
	"fmt"
	"os"

	sqlTemplate "github.com/NicklasWallgren/sqlTemplate/pkg"
)

func main() {
	wd, _ := os.Getwd()
	fs := os.DirFS(wd + "/examples/postgresql_placeholder/queries/users")

	pgPlaceholder := func(_ any, index int) string { return fmt.Sprintf("$%d", index+1) }

	sqlT := sqlTemplate.NewQueryTemplateEngine(sqlTemplate.WithPlaceholderFunc(pgPlaceholder))
	if err := sqlT.Register("users", fs, ".tsql"); err != nil {
		panic(err)
	}

	criteria := map[string]interface{}{"a": "a", "b": "b", "c": "c"}

	tmpl, err := sqlT.ParseWithValuesFromMap("users", "multipleBinds", criteria)
	if err != nil {
		panic(err)
	}

	// nolint:forbidigo
	fmt.Printf("query %v\n", tmpl.GetQuery())
	// nolint:forbidigo
	fmt.Printf("query parameters %v\n", tmpl.GetParams())
}
