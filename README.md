# SQL Template library

A simple template engine for writing dynamic SQL queries.

[![Build Status](https://github.com/NicklasWallgren/sqlTemplate/workflows/Test/badge.svg)](https://github.com/NicklasWallgren/sqlTemplate/actions?query=workflow%3ATest)
[![Reviewdog](https://github.com/NicklasWallgren/sqlTemplate/workflows/reviewdog/badge.svg)](https://github.com/NicklasWallgren/sqlTemplate/actions?query=workflow%3Areviewdog)
[![Go Report Card](https://goreportcard.com/badge/github.com/NicklasWallgren/sqlTemplate)](https://goreportcard.com/report/github.com/NicklasWallgren/sqlTemplate)
[![GoDoc](https://godoc.org/github.com/NicklasWallgren/sqlTemplate?status.svg)](https://godoc.org/github.com/NicklasWallgren/sqlTemplate)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/cabd5fbbcde543ec959fb4a3581600ed)](https://app.codacy.com/gh/NicklasWallgren/sqlTemplate?utm_source=github.com&utm_medium=referral&utm_content=NicklasWallgren/sqlTemplate&utm_campaign=Badge_Grade)

Sometimes it can be hard to write comprehensible SQL queries with tools like SQL builders ([squirrel](https://github.com/Masterminds/squirrel)
or [dbr](https://github.com/gocraft/dbr)), specially dynamic queries with optional statements and joins.
It's hard to see the overall cohesive structure of the queries, and the primary goal.

The main motivation of this library is to separate the SQL queries from the Go code, and to improve the readability of complex dynamic queries.

Check out the API Documentation http://godoc.org/github.com/NicklasWallgren/sqlTemplate

# Installation
The library can be installed through `go get`
```bash
go get github.com/NicklasWallgren/sqlTemplate
```

# Supported versions
We support the two major Go versions, which are 1.17 and 1.18 at the moment.

# Features and benefits
- Separates SQL och Go code.
- Keeps the templated query as close as possible to the actual SQL query.
- Extensible template language with support for https://github.com/Masterminds/sprig
- No third party dependencies

# SDK
```go
// Parse parses a sql template and returns the 'QueryTemplate'
Parse(namespace string, templateName string) (QueryTemplate, error)

// ParseWithValuesFromMap parses a sql template with values from a map and returns the 'QueryTemplate'
ParseWithValuesFromMap(namespace string, templateName string, parameters map[string]interface{}) (QueryTemplate, error)

// ParseWithValuesFromStruct parses a sql template with values from a struct and returns the 'QueryTemplate'
ParseWithValuesFromStruct(namespace string, templateName string, parameters interface{}) (QueryTemplate, error)

// Register registers a new namespace by template root and extension
Register(namespace string, templateRoot string, extensions string) error
```

# Examples 

## Register a namespace and parse a template
```go
sqlt := sqlTemplate.NewQueryTemplateEngine()
sqlt.Register("users", "queries/users", ".tsql");

criteria := map[string]interface{}{"Id": "1", "Order": "id"}
tmpl, _ := sqlt.ParseWithValuesFromMap("users", "findById", criteria)

fmt.Printf("query %v\n", tmpl.GetQuery())
fmt.Printf("query parameters %v\n", tmpl.GetParams())
```

```
-- File ./queries/users/users.tsql
{{define "findById"}}
    SELECT *
    FROM users
    WHERE id={{bind .Id}}
    {{if .Order}}ORDER BY {{.Order}}{{end}}
{{end}}
```

## Unit tests
```bash
go test -v -race $(go list ./... | grep -v vendor)
```

### Code Guide

We use GitHub Actions to make sure the codebase is consistent (`golangci-lint run`) and continuously tested (`go test -v -race $(go list ./... | grep -v vendor)`). We try to keep comments at a maximum of 120 characters of length and code at 120.

## Contributing

If you find any problems or have suggestions about this library, please submit an issue. Moreover, any pull request, code review and feedback are welcome.

## Contributors
- [Nicklas Wallgren](https://github.com/NicklasWallgren)
- [All Contributors][link-contributors]

[link-contributors]: ../../contributors

## License

[MIT](./LICENSE)