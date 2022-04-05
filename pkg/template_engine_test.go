package pkg

import (
	"testing"
	"text/template"
)

func TestNewQueryTemplateEngine(t *testing.T) {
	NewQueryTemplateEngine()
}

func TestNewQueryTemplateEngine_WithTemplateFunctions(t *testing.T) {
	_ = NewQueryTemplateEngine(WithTemplateFunctions(template.FuncMap{}))
}

func TestQueryTemplateEngine_Register(t *testing.T) {
}

func TestQueryTemplateEngine_Parse(t *testing.T) {
}

func TestQueryTemplateEngine_ParseWithValuesFromMap(t *testing.T) {
}

func TestQueryTemplateEngine_ParseWithValuesFromStruct(t *testing.T) {
}
