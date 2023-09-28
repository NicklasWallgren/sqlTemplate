package pkg_test

import (
	"embed"
	"fmt"
	"testing"
	"text/template"

	"github.com/NicklasWallgren/sqlTemplate/pkg"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed testdata/*.tsql
var fs embed.FS // nolint: varnamelen

func TestNewQueryTemplateEngine(_ *testing.T) {
	pkg.NewQueryTemplateEngine()
}

func TestNewQueryTemplateEngine_WithTemplateFunctions(t *testing.T) {
	if sqlT := pkg.NewQueryTemplateEngine(pkg.WithTemplateFunctions(template.FuncMap{})); sqlT == nil {
		t.Fatal()
	}
}

func TestQueryTemplateEngine_Register(t *testing.T) {
	sqlT := pkg.NewQueryTemplateEngine(pkg.WithTemplateFunctions(template.FuncMap{}))

	if err := sqlT.Register("users", fs, ".tsql"); err != nil {
		t.Fatal()
	}
}

func TestQueryTemplateEngine_Parse(t *testing.T) {
	sqlT := pkg.NewQueryTemplateEngine(pkg.WithTemplateFunctions(template.FuncMap{}))

	if err := sqlT.Register("users", fs, ".tsql"); err != nil {
		t.Fatal()
	}

	template, err := sqlT.Parse("users", "findUsers")
	assert.Nil(t, err)
	assert.Equal(t, "\n    SELECT *\n    FROM users\n", template.GetQuery())
	assert.Equal(t, []any{}, template.GetParams())
}

func TestQueryTemplateEngine_ParseWithValuesFromMap(t *testing.T) {
	sqlT := pkg.NewQueryTemplateEngine(pkg.WithTemplateFunctions(template.FuncMap{}))

	if err := sqlT.Register("users", fs, ".tsql"); err != nil {
		t.Fatal()
	}

	criteria := map[string]interface{}{"ID": 1, "Order": "id"}

	template, err := sqlT.ParseWithValuesFromMap("users", "findById", criteria)
	assert.Nil(t, err)
	assert.Equal(t, "\n    SELECT *\n    FROM users\n    WHERE id=?\n    ORDER BY id\n", template.GetQuery())
	assert.Equal(t, []any{1}, template.GetParams())
}

func TestQueryTemplateEngine_ParseWithValuesFromStruct(t *testing.T) {
	sqlT := pkg.NewQueryTemplateEngine(pkg.WithTemplateFunctions(template.FuncMap{}))

	if err := sqlT.Register("users", fs, ".tsql"); err != nil {
		t.Fatal()
	}

	type searchCriteria struct { // nolint:govet
		ID    int
		Order string
	}

	criteria := searchCriteria{ID: 1, Order: "id"}

	template, err := sqlT.ParseWithValuesFromStruct("users", "findById", criteria)
	require.Nil(t, err)
	assert.Equal(t, "\n    SELECT *\n    FROM users\n    WHERE id=?\n    ORDER BY id\n", template.GetQuery())
	assert.Equal(t, []any{1}, template.GetParams())
}

func getExtendedCriteria() map[string]any {
	return map[string]interface{}{"a": "a", "b": "b", "c": "c"}
}

func getExpectedExtendedQuery(isPgBind bool) string {
	expected := `
SELECT *
FROM users
WHERE TRUE
`
	if isPgBind {
		expected += `  AND a=$1
  AND b=$2
  AND c=$3
`
	} else {
		expected += `  AND a=?
  AND b=?
  AND c=?
`
	}

	return expected
}

func testExtendedParams(t *testing.T, params []any) {
	t.Helper()
	assert.Equal(t, 3, len(params), "wrong parameters length")
	assert.Equal(t, "a", params[0], "parameter does not match")
	assert.Equal(t, "b", params[1], "parameter does not match")
	assert.Equal(t, "c", params[2], "parameter does not match")
}

func TestQueryTemplateEngine_ParseWithValuesFromMap2(t *testing.T) {
	sqlT := pkg.NewQueryTemplateEngine(pkg.WithTemplateFunctions(template.FuncMap{}))

	if err := sqlT.Register("users", fs, ".tsql"); err != nil {
		t.Fatal()
	}

	template, err := sqlT.ParseWithValuesFromMap("users", "multipleBinds", getExtendedCriteria())
	assert.Nil(t, err)

	assert.Equal(t, getExpectedExtendedQuery(false), template.GetQuery())

	testExtendedParams(t, template.GetParams())
}

func TestQueryTemplateEngine_ParseWithValuesFromMapCustomPlaceholder(t *testing.T) {
	pgPlaceholder := func(_ any, index int) string { return fmt.Sprintf("$%d", index+1) }

	sqlT := pkg.NewQueryTemplateEngine(pkg.WithPlaceholderFunc(pgPlaceholder))

	if err := sqlT.Register("users", fs, ".tsql"); err != nil {
		t.Fatal()
	}

	template, err := sqlT.ParseWithValuesFromMap("users", "multipleBinds", getExtendedCriteria())
	require.Nil(t, err)

	assert.Equal(t, getExpectedExtendedQuery(true), template.GetQuery(), "expected query failed")
	testExtendedParams(t, template.GetParams())
}

func TestQueryTemplateEngine_ParseWithValuesFromMapCustomBindingEngine(t *testing.T) {
	pgPlaceholder := func(_ any, index int) string { return fmt.Sprintf("$%d", index+1) }
	bindingEninge := pkg.NewBindingEngine()
	bindingEninge.SetPlaceholderFunc(pgPlaceholder)
	sqlT := pkg.NewQueryTemplateEngine(pkg.WithBindingEngine(bindingEninge))

	if err := sqlT.Register("users", fs, ".tsql"); err != nil {
		t.Fatal()
	}

	template, err := sqlT.ParseWithValuesFromMap("users", "multipleBinds", getExtendedCriteria())
	require.Nil(t, err)

	assert.Equal(t, getExpectedExtendedQuery(true), template.GetQuery(), "expected query failed")
	testExtendedParams(t, template.GetParams())
}
