package object

// Create a new environment
func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s}
}

// Create an enclosed environment
func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

// Environment
type Environment struct {
	store map[string]Object
	outer *Environment
}

// Get a value from the environment by name
func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

// Set a value in the environment by name
func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}
