package evaluator

import "monkey/object"

var builtins = map[string]*object.Builtin{
	"push": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. expected=2 got=%d", len(args))
			}

			switch arg := args[0].(type) {
			case *object.Array:
				return &object.Array{Elements: append(arg.Elements, args[1])}
			default:
				return newError("argument to `push` not supported, got %s", args[0].Type())
			}
		},
	},
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. expected=1 got=%d", len(args))
			}

			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			default:
				return newError("argument to `len` not supported, got %s", args[0].Type())
			}
		},
	},
	"first": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. expected=1 got=%d", len(args))
			}

			switch arg := args[0].(type) {
			case *object.Array:
				if len(arg.Elements) == 0 {
					return NULL
				}
				return arg.Elements[0]
			default:
				return newError("argument to `first` not supported, got %s", args[0].Type())
			}
		},
	},
	"last": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. expected=1 got=%d", len(args))
			}

			switch arg := args[0].(type) {
			case *object.Array:
				if len(arg.Elements) == 0 {
					return NULL
				}
				return arg.Elements[len(arg.Elements)-1]
			default:
				return newError("argument to `last` not supported, got %s", args[0].Type())
			}
		},
	},
	"rest": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. expected=1 got=%d", len(args))
			}

			switch arg := args[0].(type) {
			case *object.Array:
				if len(arg.Elements) < 2 {
					return NULL
				}
				return &object.Array{Elements: arg.Elements[1:]}
			default:
				return newError("argument to `rest` not supported, got %s", args[0].Type())
			}
		},
	},
}
