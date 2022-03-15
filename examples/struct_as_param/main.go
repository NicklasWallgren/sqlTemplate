package main

import (
	"fmt"

	sqlTemplate "github.com/NicklasWallgren/sqlTemplate/pkg"
)

type searchCriteria struct {
	Name    string
	Surname string
	Order   string
}

func main() {
	sqlT := sqlTemplate.NewQueryTemplateEngine()
	if err := sqlT.Register("users", "queries/users", ".tsql"); err != nil {
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
