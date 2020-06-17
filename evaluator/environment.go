package evaluator

type Environment struct {
	store map[string]Object
	outer *Environment
}

func NewScopedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s}
}

func (e *Environment) Get(name string) (obj Object, ok bool) {
	obj, ok = e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return
}

func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}
