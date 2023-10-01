package pkg

import (
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	txtTemplate "text/template"
)

type templateEntry struct {
	template *txtTemplate.Template
	fullPath string
}

// repository stores SQL templates.
type repository struct {
	templates map[string]*templateEntry
	functions txtTemplate.FuncMap
}

// newRepository creates a new repository.
func newRepository() *repository {
	return &repository{
		templates: make(map[string]*templateEntry),
		functions: txtTemplate.FuncMap{"bind": bind},
	}
}

func (r *repository) addFunctions(functions txtTemplate.FuncMap) {
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

	rootTemplate := txtTemplate.New(namespace).Funcs(r.functions)

	for _, filename := range filesInFilesystem {
		parsedTemplate, err := rootTemplate.ParseFS(filesystem, filename)
		if err != nil {
			return fmt.Errorf("unable to parse file %s %w", filename, err)
		}

		r.templates[namespace] = &templateEntry{parsedTemplate, filename}
	}

	return nil
}

// getTemplate returns the template by his namespace.
func (r *repository) getTemplate(namespace string) (*txtTemplate.Template, error) {
	entry, ok := r.templates[namespace]
	if !ok {
		return nil, errors.New("unable to locate namespace " + namespace)
	}

	// We clone the template to prevent simultaneous mutation of the template.FuncMap
	// otherwise the binds functions might be replaced during execution of a template
	clonedTmpl, err := entry.template.Clone()
	if err != nil {
		return nil, fmt.Errorf("unable to parse template %w", err)
	}

	return clonedTmpl, nil
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
