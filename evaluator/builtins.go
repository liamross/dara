package evaluator

var builtins = map[string]*Builtin{
	"len": {
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError("invalid operation: too many arguments for len (expected %d, found %d)",
					1, len(args))
			}
			switch arg := args[0].(type) {
			case *String:
				return &Number{Value: float64(len(arg.Value))}
			case *Array:
				return &Number{Value: float64(len(arg.Elements))}
			default:
				return newError("invalid argument: %s (%s) for len",
					arg.Inspect(), arg.Type())
			}
		},
	},
}
