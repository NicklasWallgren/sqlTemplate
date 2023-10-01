# SQL Template library

A simple template engine for writing dynamic SQL queries.

[![Build Status](https://github.com/NicklasWallgren/sqlTemplate/workflows/Test/badge.svg)](https://github.com/NicklasWallgren/sqlTemplate/actions?query=workflow%3ATest)
[![Reviewdog](https://github.com/NicklasWallgren/sqlTemplate/workflows/reviewdog/badge.svg)](https://github.com/NicklasWallgren/sqlTemplate/actions?query=workflow%3Areviewdog)
[![Go Report Card](https://goreportcard.com/badge/github.com/NicklasWallgren/sqlTemplate)](https://goreportcard.com/report/github.com/NicklasWallgren/sqlTemplate)
[![GoDoc](https://godoc.org/github.com/NicklasWallgren/sqlTemplate?status.svg)](https://godoc.org/github.com/NicklasWallgren/sqlTemplate)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/cabd5fbbcde543ec959fb4a3581600ed)](https://app.codacy.com/gh/NicklasWallgren/sqlTemplate?utm_source=github.com&utm_medium=referral&utm_content=NicklasWallgren/sqlTemplate&utm_campaign=Badge_Grade)

Sometimes it can be hard to write comprehensible SQL queries with
tools like SQL builders ([squirrel](https://github.com/Masterminds/squirrel)
or [dbr](https://github.com/gocraft/dbr)), specially dynamic queries
with optional statements and joins.
It can be hard to see the overall cohesive structure of the queries,
and the primary goal.

The main motivation of this library is to separate the SQL queries
from the Go code, and to improve the readability of complex dynamic
queries.

Check out the API Documentation http://godoc.org/github.com/NicklasWallgren/sqlTemplate

# Installation
The library can be installed through `go get`
```bash
go get github.com/NicklasWallgren/sqlTemplate
```

# Supported versions
We support the latest major Go version, which are 1.19 at the moment.

# Features and benefits
- Separates SQL and Go code.
- Keeps the templated query as close as possible to the actual SQL query.
- Extensible template language with support for https://github.com/Masterminds/sprig
- No third party dependencies
- Support for embedded filesystem

# API
```go
// Parse parses a sql template and returns the 'QueryTemplate'
Parse(namespace string, templateName string) (QueryTemplate, error)

// ParseWithValuesFromMap parses a sql template with values from a map and returns the 'QueryTemplate'
ParseWithValuesFromMap(namespace string, templateName string, parameters map[string]interface{}) (QueryTemplate, error)

// ParseWithValuesFromStruct parses a sql template with values from a struct and returns the 'QueryTemplate'
ParseWithValuesFromStruct(namespace string, templateName string, parameters interface{}) (QueryTemplate, error)

// Register registers a new namespace by template filesystem and extension
Register(namespace string, filesystem fs.FS, extensions string) error
```

# Examples 

## Register a namespace and parse a template
```go
//go:embed queries/users/*.tsql
var fs embed.FS

sqlt := sqlTemplate.NewQueryTemplateEngine()
sqlt.Register("users", fs, ".tsql");

criteria := map[string]interface{}{"Name": "Bill", "Order": "id"}
tmpl, _ := sqlt.ParseWithValuesFromMap("users", "findByName", criteria)

sql.QueryRowContext(context.Background(), tmpl.GetQuery(), tmpl.GetParams())

fmt.Printf("query %v\n", tmpl.GetQuery())
fmt.Printf("query parameters %v\n", tmpl.GetParams())
```

```sql
-- File ./queries/users/users.tsql
{{define "findByName"}}
    SELECT *
    FROM users
    WHERE name={{bind .Name}}
    {{if .Order}}ORDER BY {{.Order}}{{end}}
{{end}}
```

## Unit tests

```bash
go test -v -race ./pkg
```

For benchmark :

```bash
go test ./pkg -bench=.
```

### Code Guide

We use GitHub Actions to make sure the codebase is consistent
(`golangci-lint run`) and continuously tested (`go test -v -race
./pkg`). We try to keep comments at a maximum of 120 characters of
length and code at 120.

## Contributing

If you find any problems or have suggestions about this library,
please submit an issue. Moreover, any pull request, code review and
feedback are welcome.

## Contributors
- [Nicklas Wallgren](https://github.com/NicklasWallgren)
- [All Contributors][link-contributors]

[link-contributors]: ../../contributors

## License

[MIT](./LICENSE)
