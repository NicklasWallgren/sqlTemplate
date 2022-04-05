package main

import (
	"fmt"
	"os"

	sqlTemplate "github.com/NicklasWallgren/sqlTemplate/pkg"
)

func main() {
	wd, _ := os.Getwd()
	fs := os.DirFS(wd + "/examples/map_as_param/queries/users")

	sqlT := sqlTemplate.NewQueryTemplateEngine()
	if err := sqlT.Register("users", fs, ".tsql"); err != nil {
		panic(err)
	}

	criteria := map[string]interface{}{"Id": "1", "Order": "id"}

	tmpl, err := sqlT.ParseWithValuesFromMap("users", "findById", criteria)
	if err != nil {
		panic(err)
	}

	// nolint:forbidigo
	fmt.Printf("query %v\n", tmpl.GetQuery())
	// nolint:forbidigo
	fmt.Printf("query parameters %v\n", tmpl.GetParams())
}
