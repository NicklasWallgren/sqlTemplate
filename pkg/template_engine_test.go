package pkg_test

import (
	"embed"
	"fmt"
	"math/rand"
	"testing"
	"text/template"
	"time"

	"github.com/NicklasWallgren/sqlTemplate/pkg"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed testdata/*.tsql
var fs embed.FS // nolint: varnamelen

func TestNewTemplateEngine(_ *testing.T) {
	pkg.NewTemplateEngine()
}

func TestNewTemplateEngine_WithTemplateFunctions(t *testing.T) {
	if sqlT := pkg.NewTemplateEngine(pkg.WithTemplateFunctions(template.FuncMap{})); sqlT == nil {
		t.Fatal()
	}
}

func TestTemplateEngine_Register(t *testing.T) {
	sqlT := pkg.NewTemplateEngine(pkg.WithTemplateFunctions(template.FuncMap{}))

	if err := sqlT.Register("users", fs, ".tsql"); err != nil {
		t.Fatal()
	}
}

func TestTemplateEngine_Parse(t *testing.T) {
	sqlT := pkg.NewTemplateEngine(pkg.WithTemplateFunctions(template.FuncMap{}))

	if err := sqlT.Register("users", fs, ".tsql"); err != nil {
		t.Fatal()
	}

	template, err := sqlT.Parse("users", "findUsers")
	assert.Nil(t, err)
	assert.Equal(t, "\n    SELECT *\n    FROM users\n", template.GetQuery())
	assert.Equal(t, []any{}, template.GetParams())
}

func TestTemplateEngine_ParseWithValuesFromMap(t *testing.T) {
	sqlT := pkg.NewTemplateEngine(pkg.WithTemplateFunctions(template.FuncMap{}))

	if err := sqlT.Register("users", fs, ".tsql"); err != nil {
		t.Fatal()
	}

	criteria := map[string]interface{}{"ID": 1, "Order": "id"}

	template, err := sqlT.ParseWithValuesFromMap("users", "findById", criteria)
	assert.Nil(t, err)
	assert.Equal(t, "\n    SELECT *\n    FROM users\n    WHERE id=?\n    ORDER BY id\n", template.GetQuery())
	assert.Equal(t, []any{1}, template.GetParams())
}

func TestTemplateEngine_ParseWithValuesFromStruct(t *testing.T) {
	sqlT := pkg.NewTemplateEngine(pkg.WithTemplateFunctions(template.FuncMap{}))

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

func getExtendedCriteria(i int64) map[string]any {
	myRand := rand.New(rand.NewSource(i + time.Now().UnixNano()))
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	nbLetters := len(letters)
	aVal := letters[myRand.Intn(nbLetters)]
	bVal := letters[myRand.Intn(nbLetters)]
	cVal := letters[myRand.Intn(nbLetters)]

	return map[string]interface{}{"a": aVal, "b": bVal, "c": cVal}
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

func testExtendedParams(t *testing.T, actual []any, expected map[string]any) {
	t.Helper()
	assert.Equal(t, 3, len(actual), "wrong parameters length")
	assert.Equal(t, expected["a"], actual[0], "parameter does not match")
	assert.Equal(t, expected["b"], actual[1], "parameter does not match")
	assert.Equal(t, expected["c"], actual[2], "parameter does not match")
}

func TestTemplateEngine_ParseWithValuesFromMap2(t *testing.T) {
	sqlT := pkg.NewTemplateEngine(pkg.WithTemplateFunctions(template.FuncMap{}))

	if err := sqlT.Register("users", fs, ".tsql"); err != nil {
		t.Fatal(err)
	}

	query, params, criteria, err := parseWithValuesFromMap(sqlT)
	require.Nil(t, err)
	assert.Equal(t, getExpectedExtendedQuery(false), query)
	testExtendedParams(t, params, criteria)
}

func parseWithValuesFromMap(sqlT pkg.TemplateEngine) (query string, params []any, criteria map[string]any, err error) {
	criteria = getExtendedCriteria(50)
	template, err := sqlT.ParseWithValuesFromMap("users", "multipleBinds", criteria)
	query = template.GetQuery()
	params = template.GetParams()

	return
}

var pQuery string
var pParams []any

func TestTemplateEngine_ParseWithValuesFromMapConcurrent(t *testing.T) {
	sqlT := pkg.NewTemplateEngine(pkg.WithTemplateFunctions(template.FuncMap{}))

	if err := sqlT.Register("users", fs, ".tsql"); err != nil {
		t.Fatal(err)
	}

	start := time.Now()
	for i := int64(0); i < 100; i++ {
		myRand := rand.New(rand.NewSource(i - time.Now().UnixNano()))
		go func(j int64) {
			for time.Since(start) < time.Second {
				time.Sleep(time.Duration(myRand.Intn(40)))

				query, params, criteria, err := parseWithValuesFromMap(sqlT)
				require.Nil(t, err)
				assert.Equal(t, getExpectedExtendedQuery(false), query)
				testExtendedParams(t, params, criteria)
			}
		}(i)
	}

	for time.Since(start) < time.Second {
		query, params, criteria, err := parseWithValuesFromMap(sqlT)
		require.Nil(t, err)
		assert.Equal(t, getExpectedExtendedQuery(false), query)
		testExtendedParams(t, params, criteria)
	}
}

func BenchmarkTemplateEngine_ParseWithValuesFromMap(b *testing.B) {
	sqlT := pkg.NewTemplateEngine(pkg.WithTemplateFunctions(template.FuncMap{}))

	if err := sqlT.Register("users", fs, ".tsql"); err != nil {
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		query, params, _, _ := parseWithValuesFromMap(sqlT)

		// always store the result to a package level variable
		// so the compiler cannot eliminate the Benchmark itself.
		pQuery = query
		pParams = params
	}
}

func TestTemplateEngine_ParseWithValuesFromMapCustomPlaceholder(t *testing.T) {
	pgPlaceholder := func(_ any, index int) string { return fmt.Sprintf("$%d", index+1) }

	sqlT := pkg.NewTemplateEngine(pkg.WithPlaceholderFunc(pgPlaceholder))

	if err := sqlT.Register("users", fs, ".tsql"); err != nil {
		t.Fatal()
	}

	query, params, criteria, err := parseWithValuesFromMap(sqlT)
	require.Nil(t, err)
	assert.Equal(t, getExpectedExtendedQuery(true), query)
	testExtendedParams(t, params, criteria)
}

func TestTemplateEngine_ParseWithValuesFromMapCustomiBndingEngine(t *testing.T) {
	pgPlaceholder := func(_ any, index int) string { return fmt.Sprintf("$%d", index+1) }
	bindingEninge := pkg.NewBindingEngine()
	bindingEninge.SetPlaceholderFunc(pgPlaceholder)
	sqlT := pkg.NewTemplateEngine(pkg.WithBindingEngine(bindingEninge))

	if err := sqlT.Register("users", fs, ".tsql"); err != nil {
		t.Fatal()
	}

	query, params, criteria, err := parseWithValuesFromMap(sqlT)
	require.Nil(t, err)
	assert.Equal(t, getExpectedExtendedQuery(true), query)
	testExtendedParams(t, params, criteria)
}
