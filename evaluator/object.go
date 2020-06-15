package evaluator

import "fmt"

type ObjectType string

const (
	NUMBER_OBJ       ObjectType = "number"
	BOOLEAN_OBJ      ObjectType = "boolean"
	STRING_OBJ       ObjectType = "string"
	NIL_OBJ          ObjectType = "nil"
	RETURN_VALUE_OBJ ObjectType = "return_value"
	ERROR_OBJ        ObjectType = "error"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Number struct {
	Value float64
}

func (n *Number) Type() ObjectType { return NUMBER_OBJ }
func (n *Number) Inspect() string  { return fmt.Sprintf("%v", n.Value) }

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }

type String struct {
	Value string
}

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string  { return `"` + s.Value + `"` }

type Nil struct{}

func (n *Nil) Type() ObjectType { return NIL_OBJ }
func (n *Nil) Inspect() string  { return "nil" }

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }
func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return RETURN_VALUE_OBJ }
func (e *Error) Inspect() string  { return e.Message }
