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

// bindings stores the values of any placeholder parameter in the query.
type bindings struct {
	values []interface{}
}

// bind stores the given `value` and returns a placeholder parameter.
func (b *bindings) bind(value interface{}) string {
	b.values = append(b.values, value)

	return "?"
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
func (r *repository) add(namespace string, filesystem fs.FS, extension string) (err error) {
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
func (r *repository) parse(namespace string, name string, data interface{}) (string, []interface{}, error) {
	entry, ok := r.templates[namespace]
	if !ok {
		return "", nil, errors.New("unable to locate namespace " + namespace)
	}

	// We clone the template to prevent simultaneous mutation of the template.FuncMap
	// otherwise the bind function might be replaced during execution of a template
	clonedTmpl, err := entry.template.Clone()
	if err != nil {
		return "", nil, fmt.Errorf("unable to parse template %w", err)
	}

	// Apply the bind function which stores the values for any placeholder parameters
	values := &bindings{values: []interface{}{}}

	clonedTmpl.Funcs(template.FuncMap{"bind": values.bind})

	var b bytes.Buffer
	if err := clonedTmpl.ExecuteTemplate(&b, name, data); err != nil {
		return "", nil, fmt.Errorf("unable to execute template %w", err)
	}

	return b.String(), values.values, nil
}

// bind is a dummy function which is never used while executing a template.
func bind(param interface{}) string {
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
