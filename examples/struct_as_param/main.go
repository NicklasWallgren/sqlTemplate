package main

import (
	"fmt"
	"os"

	sqlTemplate "github.com/NicklasWallgren/sqlTemplate/pkg"
)

type searchCriteria struct {
	Name    string
	Surname string
	Order   string
}

func main() {
	wd, _ := os.Getwd()
	fs := os.DirFS(wd + "/examples/struct_as_param/queries/users")

	sqlT := sqlTemplate.NewTemplateEngine()
	if err := sqlT.Register("users", fs, ".tsql"); err != nil {
		panic(err)
	}

	criteria := searchCriteria{Name: "Bill", Surname: "Gates", Order: "id"}

	tmpl, err := sqlT.ParseWithValuesFromStruct("users", "findByName", criteria)
	if err != nil {
		panic(err)
	}

	// nolint:forbidigo
	fmt.Printf("query %v\n", tmpl.GetQuery())
	// nolint:forbidigo
	fmt.Printf("query parameters %v\n", tmpl.GetParams())
}
