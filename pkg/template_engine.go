// Package pkg is the public lib of SQLTemplate.
// SQLTemplate is a simple template engine for writing
// dynamic SQL queries.
package pkg

import (
	"bytes"
	"fmt"
	"io/fs"
	txtTemplate "text/template"
)

// TemplateEngine is the interface implemented by types that can parse sql templates.
type TemplateEngine interface {
	// Parse parses a sql template and returns the 'Template'.
	Parse(namespace string, templateName string) (Template, error)
	// ParseWithValuesFromMap parses a sql template with values from a map and returns the 'Template'.
	ParseWithValuesFromMap(namespace string, templateName string, parameters map[string]interface{}) (Template, error)
	// ParseWithValuesFromStruct parses a sql template with values from a struct and returns the 'Template'.
	ParseWithValuesFromStruct(namespace string, templateName string, parameters interface{}) (Template, error)
	// Register registers a new namespace by template filesystem and extension.
	Register(namespace string, filesystem fs.FS, extensions string) error
}

// QueryTemplateEngine is for backwards compatibility.
// Deprecated: use TemplateEngine.
type QueryTemplateEngine = TemplateEngine

// Template is the interface implemented by types that holds the parsed template sql context.
type Template interface {
	// GetQuery returns the query containing named values.
	GetQuery() string
	// GetParams returns the values in order.
	GetParams() []interface{}
}

// QueryTemplate is for backwards compatibility.
// Deprecated: use Template.
type QueryTemplate = Template

// Option definition.
type Option func(*templateEngine)

type templateEngine struct {
	repository    *repository
	bindingEngine bindingEngine
}

type template struct {
	template string
	params   []interface{}
}

func (t template) GetQuery() string {
	return t.template
}

func (t template) GetParams() []interface{} {
	return t.params
}

// WithTemplateFunctions creates an Option func to set template functions.
// nolint:deadcode
func WithTemplateFunctions(funcMap txtTemplate.FuncMap) Option {
	return func(templateTypeEngine *templateEngine) {
		templateTypeEngine.repository.addFunctions(funcMap)
	}
}

// WithBindingEngine creates an Option func to set custom binding engine.
// nolint:deadcode
func WithBindingEngine(bEngine bindingEngine) Option {
	return func(templateTypeEngine *templateEngine) {
		templateTypeEngine.bindingEngine = bEngine
	}
}

// WithPlaceholderFunc creates an Option func to set custom placeholder function.
// nolint:deadcode
func WithPlaceholderFunc(placeholderfunc placeholderFunc) Option {
	return func(templateTypeEngine *templateEngine) {
		if templateTypeEngine.bindingEngine == nil {
			templateTypeEngine.bindingEngine = NewBindingEngine()
		}

		templateTypeEngine.bindingEngine.SetPlaceholderFunc(placeholderfunc)
	}
}

// NewTemplateEngine returns a new instance of 'TemplateEngine'.
func NewTemplateEngine(options ...Option) TemplateEngine {
	templateEngine := &templateEngine{repository: newRepository(), bindingEngine: nil}

	// Apply options if there are any, can overwrite default
	for _, option := range options {
		option(templateEngine)
	}

	return templateEngine
}

// Deprecated: use NewTemplateEngine.
func NewQueryTemplateEngine(options ...Option) TemplateEngine {
	return NewTemplateEngine(options...)
}

func (q templateEngine) Register(namespace string, filesystem fs.FS, ext string) error {
	err := q.repository.add(namespace, filesystem, ext)
	if err != nil {
		return fmt.Errorf("could not register the namespace %s %w", namespace, err)
	}

	return nil
}

func (q templateEngine) Parse(namespace string, templateName string) (Template, error) {
	sqlQuery, bindings, err := q.parse(namespace, templateName, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to parse %s for namespace %s %w", templateName, namespace, err)
	}

	return &template{sqlQuery, bindings}, nil
}

func (q templateEngine) ParseWithValuesFromMap(namespace string, templateName string, parameters map[string]interface{}) (Template, error) {
	sqlQuery, bindings, err := q.parse(namespace, templateName, parameters)
	if err != nil {
		return nil, fmt.Errorf("unable to parse %s for namespace %s %w", templateName, namespace, err)
	}

	return &template{sqlQuery, bindings}, nil
}

func (q templateEngine) ParseWithValuesFromStruct(namespace string, templateName string, parameters interface{}) (Template, error) {
	sqlQuery, bindings, err := q.parse(namespace, templateName, parameters)
	if err != nil {
		return nil, fmt.Errorf("unable to parse %s for namespace %s %w", templateName, namespace, err)
	}

	return &template{sqlQuery, bindings}, nil
}

// parse executes the template and returns the resulting SQL or an error.
func (q templateEngine) parse(namespace string, name string, data interface{}) (string, []interface{}, error) {
	tmpl, err := q.repository.getTemplate(namespace)
	if err != nil {
		return "", nil, err
	}

	var bEngine bindingEngine

	// Apply the bind function which stores the values for any placeholder parameters
	if q.bindingEngine == nil {
		bEngine = &DefaultBindingEngine{values: []any{}, index: 0, placeholderFunc: defaultPlaceholderFunc}
	} else {
		bEngine = q.bindingEngine.new()
	}

	tmpl.Funcs(txtTemplate.FuncMap{"bind": func(value any) string {
		index := bEngine.storeValue(value)

		return bEngine.getPlaceholderFunc()(value, index)
	}})

	var b bytes.Buffer
	if err := tmpl.ExecuteTemplate(&b, name, data); err != nil {
		return "", nil, fmt.Errorf("unable to execute template %w", err)
	}

	return b.String(), bEngine.getValues(), nil
}
