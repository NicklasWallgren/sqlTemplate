package pkg

import (
	"fmt"
	"io/fs"
	"text/template"
)

// QueryTemplateEngine is the interface implemented by types that can parse sql templates.
type QueryTemplateEngine interface {
	// Parse parses a sql template and returns the 'QueryTemplate'
	Parse(namespace string, templateName string) (QueryTemplate, error)
	// ParseWithValuesFromMap parses a sql template with values from a map and returns the 'QueryTemplate'
	ParseWithValuesFromMap(namespace string, templateName string, parameters map[string]interface{}) (QueryTemplate, error)
	// ParseWithValuesFromStruct parses a sql template with values from a struct and returns the 'QueryTemplate'
	ParseWithValuesFromStruct(namespace string, templateName string, parameters interface{}) (QueryTemplate, error)
	// Register registers a new namespace by template filesystem and extension
	Register(namespace string, filesystem fs.FS, extensions string) error
}

// QueryTemplate is the interface implemented by types that holds the parsed template sql query context.
type QueryTemplate interface {
	// GetQuery returns the query containing named values.
	GetQuery() string
	// GetParams returns the values in order.
	GetParams() []interface{}
}

// Option definition.
type Option func(*queryTemplateEngine)

type queryTemplateEngine struct {
	repository    *repository
	bindingEngine bindingEngine
}

type queryTemplate struct {
	template string
	params   []interface{}
}

func (t queryTemplate) GetQuery() string {
	return t.template
}

func (t queryTemplate) GetParams() []interface{} {
	return t.params
}

// WithTemplateFunctions creates an Option func to set template functions.
// nolint:deadcode
func WithTemplateFunctions(funcMap template.FuncMap) Option {
	return func(queryTypeEngine *queryTemplateEngine) {
		queryTypeEngine.repository.addFunctions(funcMap)
	}
}

// WithBindingEngine creates an Option func to set custom binding engine.
// nolint:deadcode
func WithBindingEngine(bEngine bindingEngine) Option {
	return func(queryTypeEngine *queryTemplateEngine) {
		queryTypeEngine.bindingEngine = bEngine
	}
}

// WithPlaceholderFunc creates an Option func to set custom placeholder function.
// nolint:deadcode
func WithPlaceholderFunc(placeholderfunc placeholderFunc) Option {
	return func(queryTypeEngine *queryTemplateEngine) {
		if queryTypeEngine.bindingEngine == nil {
			queryTypeEngine.bindingEngine = NewBindingEngine()
		}

		queryTypeEngine.bindingEngine.SetPlaceholderFunc(placeholderfunc)
	}
}

// NewQueryTemplateEngine returns a new instance of 'QueryTemplateEngine'.
func NewQueryTemplateEngine(options ...Option) QueryTemplateEngine {
	templateEngine := &queryTemplateEngine{repository: newRepository(), bindingEngine: nil}

	// Apply options if there are any, can overwrite default
	for _, option := range options {
		option(templateEngine)
	}

	return templateEngine
}

func (q queryTemplateEngine) Register(namespace string, filesystem fs.FS, ext string) error {
	err := q.repository.add(namespace, filesystem, ext)
	if err != nil {
		return fmt.Errorf("could not register the namespace %s %w", namespace, err)
	}

	return nil
}

func (q queryTemplateEngine) Parse(namespace string, templateName string) (QueryTemplate, error) {
	sqlQuery, bindings, err := q.repository.parse(namespace, templateName, nil, q.bindingEngine)
	if err != nil {
		return nil, fmt.Errorf("unable to parse %s for namespace %s %w", templateName, namespace, err)
	}

	return &queryTemplate{sqlQuery, bindings}, nil
}

func (q queryTemplateEngine) ParseWithValuesFromMap(namespace string, templateName string, parameters map[string]interface{}) (QueryTemplate, error) {
	sqlQuery, bindings, err := q.repository.parse(namespace, templateName, parameters, q.bindingEngine)
	if err != nil {
		return nil, fmt.Errorf("unable to parse %s for namespace %s %w", templateName, namespace, err)
	}

	return &queryTemplate{sqlQuery, bindings}, nil
}

func (q queryTemplateEngine) ParseWithValuesFromStruct(namespace string, templateName string, parameters interface{}) (QueryTemplate, error) {
	sqlQuery, bindings, err := q.repository.parse(namespace, templateName, parameters, q.bindingEngine)
	if err != nil {
		return nil, fmt.Errorf("unable to parse %s for namespace %s %w", templateName, namespace, err)
	}

	return &queryTemplate{sqlQuery, bindings}, nil
}
