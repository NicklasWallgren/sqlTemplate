package main

import (
	"embed"
	"fmt"

	sqlTemplate "github.com/NicklasWallgren/sqlTemplate/pkg"
)

//go:embed queries/users/*.tsql
var fs embed.FS

func main() {
	sqlT := sqlTemplate.NewTemplateEngine()
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
