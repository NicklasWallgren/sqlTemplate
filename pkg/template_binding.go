package pkg

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

// NewBindingEngine build a binding engine based on the default binding engine.
func NewBindingEngine() bindingEngine {
	return &DefaultBindingEngine{values: []any{}, index: 0, placeholderFunc: defaultPlaceholderFunc}
}

// bind is a dummy function which is never used while executing a template.
func bind(_ interface{}) string {
	panic("dummy function which should never be used.")
}
