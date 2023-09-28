package pkg

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"text/template"
)

type templateEntry struct {
	template *template.Template
	fullPath string
}

type placeholderFunc func(value any, index int) string

func defaultPlaceholderFunc(_ any, _ int) string { return "?" }

type bindingEngine interface {
	// StoreValue stores the given `value` and return the index order of
	// the stored value.
	storeValue(any) int
	// GetValues get the stored values.
	getValues() []any
	// new returns an new instance.
	new() bindingEngine
	// SetPlaceholderFunc allows to set custom placeholderFunc.
	SetPlaceholderFunc(placeholderFunc)
	// getPlaceholderFunc returns the placeholder function.
	getPlaceholderFunc() placeholderFunc
}

// DefaultBindingEngine is a base engine to handle bindings (placeholder "?" for MySQL-like dbe).
// It can be oveload to provide other type if bindings (placeholder "$i" for PostgreSQL-like dbe).
type DefaultBindingEngine struct {
	values          []any
	index           int
	placeholderFunc placeholderFunc
}

func (b *DefaultBindingEngine) new() bindingEngine {
	newBE := NewBindingEngine()
	newBE.SetPlaceholderFunc(b.placeholderFunc)

	return newBE
}

func (b *DefaultBindingEngine) getPlaceholderFunc() placeholderFunc {
	if b.placeholderFunc == nil {
		return defaultPlaceholderFunc
	}

	return b.placeholderFunc
}

// SetPlaceholderFunc allows to set custom placeholder function.
func (b *DefaultBindingEngine) SetPlaceholderFunc(placeholderfunc placeholderFunc) {
	b.placeholderFunc = placeholderfunc
}

func (b *DefaultBindingEngine) storeValue(value any) int {
	b.values = append(b.values, value)
	b.index++

	return b.index - 1
}

func (b *DefaultBindingEngine) getValues() []any {
	return b.values
}

func NewBindingEngine() bindingEngine {
	return &DefaultBindingEngine{values: []any{}, index: 0, placeholderFunc: defaultPlaceholderFunc}
}

// repository stores SQL templates.
type repository struct {
	templates map[string]*templateEntry
	functions template.FuncMap
}

// newRepository creates a new repository.
func newRepository() *repository {
	return &repository{
		templates: make(map[string]*templateEntry),
		functions: template.FuncMap{"bind": bind},
	}
}

func (r *repository) addFunctions(functions template.FuncMap) {
	for name, function := range functions {
		r.functions[name] = function
	}
}

// add walks a filesystem and parses the corresponding templates.
func (r *repository) add(namespace string, filesystem fs.FS, extension string) error {
	filesInFilesystem, err := getFilesInFilesystem(filesystem, extension)
	if err != nil {
		return fmt.Errorf("unable to retrieve files in directory %w", err)
	}

	rootTemplate := template.New(namespace).Funcs(r.functions)

	for _, filename := range filesInFilesystem {
		parsedTemplate, err := rootTemplate.ParseFS(filesystem, filename)
		if err != nil {
			return fmt.Errorf("unable to parse file %s %w", filename, err)
		}

		r.templates[namespace] = &templateEntry{parsedTemplate, filename}
	}

	return nil
}

// parse executes the template and returns the resulting SQL or an error.
func (r *repository) parse(namespace string, name string, data interface{}, bEngine bindingEngine) (string, []interface{}, error) {
	entry, ok := r.templates[namespace]
	if !ok {
		return "", nil, errors.New("unable to locate namespace " + namespace)
	}

	// We clone the template to prevent simultaneous mutation of the template.FuncMap
	// otherwise the binds functions might be replaced during execution of a template
	clonedTmpl, err := entry.template.Clone()
	if err != nil {
		return "", nil, fmt.Errorf("unable to parse template %w", err)
	}

	// Apply the bind function which stores the values for any placeholder parameters
	if bEngine == nil {
		bEngine = &DefaultBindingEngine{values: []any{}, index: 0, placeholderFunc: defaultPlaceholderFunc}
	} else {
		bEngine = bEngine.new()
	}

	clonedTmpl.Funcs(template.FuncMap{"bind": func(value any) string {
		index := bEngine.storeValue(value)

		return bEngine.getPlaceholderFunc()(value, index)
	}})

	var b bytes.Buffer
	if err := clonedTmpl.ExecuteTemplate(&b, name, data); err != nil {
		return "", nil, fmt.Errorf("unable to execute template %w", err)
	}

	return b.String(), bEngine.getValues(), nil
}

// bind is a dummy function which is never used while executing a template.
func bind(_ interface{}) string {
	return "?"
}

// getFilesInFilesystem walks the directory tree and returns a slice of files with the given extension.
func getFilesInFilesystem(filesystem fs.FS, extension string) ([]string, error) {
	var files []string

	err := fs.WalkDir(filesystem, ".", func(path string, d fs.DirEntry, err error) error {
		if filepath.Ext(path) == extension {
			files = append(files, path)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve files in directory %w", err)
	}

	return files, nil
}
