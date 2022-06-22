package exec

import (
	"container/list"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"time"

	"github.com/hsiaosiyuan0/mole/ecma/parser"
	"github.com/hsiaosiyuan0/mole/ecma/walk"
)

type ExprEvaluator struct {
	id      string
	walkCtx *walk.WalkCtx
	vars    map[string]interface{}

	this interface{}
	stk  *list.List
	err  error

	listeners map[parser.NodeType]*walk.Listener
}

var builtinProto = map[reflect.Type]map[string]NativeFn{
	reflect.TypeOf([]interface{}{}): {
		"includes": func(this interface{}, args *Args) interface{} {
			arr := this.([]interface{})
			target := args.Get(0)
			for _, arg := range arr {
				if reflect.DeepEqual(arg, target) {
					return true
				}
			}
			return false
		},
	},
}

func NewExprEvaluator(wc *walk.WalkCtx) *ExprEvaluator {
	ee := &ExprEvaluator{
		id:        fmt.Sprintf("expr_evaluator_%d", time.Now().Nanosecond()),
		walkCtx:   wc,
		vars:      map[string]interface{}{},
		stk:       list.New(),
		listeners: map[parser.NodeType]*walk.Listener{},
	}
	ee.init()
	return ee
}

func (ee *ExprEvaluator) addListener(nt parser.NodeType, impl walk.ListenFn) {
	fn := &walk.Listener{
		Id:     ee.id + nt.String(),
		Handle: impl,
	}
	walk.AddListener(&ee.walkCtx.Listeners, nt, fn)
	ee.listeners[nt] = fn
}

func (ee *ExprEvaluator) push(v interface{}) {
	ee.stk.PushBack(v)
}

func (ee *ExprEvaluator) pop() interface{} {
	b := ee.stk.Back()
	if b != nil {
		ee.stk.Remove(b)
	}
	if b != nil {
		return b.Value
	}
	return nil
}

func (ee *ExprEvaluator) init() {
	ee.addListener(walk.NodeAfterEvent(parser.N_LIT_BOOL),
		func(node parser.Node, key string, ctx *walk.VisitorCtx) {
			if ee.err != nil {
				return
			}

			ee.push(node.(*parser.BoolLit).Val())
		})

	ee.addListener(walk.NodeAfterEvent(parser.N_LIT_NULL),
		func(node parser.Node, key string, ctx *walk.VisitorCtx) {
			if ee.err != nil {
				return
			}

			ee.push(nil)
		})

	ee.addListener(walk.NodeAfterEvent(parser.N_NAME),
		func(node parser.Node, key string, ctx *walk.VisitorCtx) {
			if ee.err != nil {
				return
			}

			name := node.(*parser.Ident).Text()
			if ctx.ParentNode().Type() == parser.N_EXPR_MEMBER && key == "Prop" {
				ee.push(name)
			} else if name == "undefined" {
				ee.push(nil)
			} else {
				ee.push(ee.vars[name])
			}
		})

	ee.addListener(walk.NodeAfterEvent(parser.N_LIT_ARR),
		func(node parser.Node, key string, ctx *walk.VisitorCtx) {
			if ee.err != nil {
				return
			}

			n := node.(*parser.ArrLit)
			arrLen := len(n.Elems())
			arr := make([]interface{}, arrLen)
			for i := len(n.Elems()) - 1; i >= 0; i-- {
				arr[i] = ee.pop()
			}
			ee.push(arr)
		})

	ee.addListener(walk.NodeAfterEvent(parser.N_LIT_NUM),
		func(node parser.Node, key string, ctx *walk.VisitorCtx) {
			if ee.err != nil {
				return
			}

			n := node.(*parser.NumLit)
			i, err := strconv.ParseFloat(n.Text(), 64)
			if err != nil {
				ee.push(math.NaN())
			} else {
				ee.push(i)
			}
		})

	ee.addListener(walk.NodeAfterEvent(parser.N_LIT_STR),
		func(node parser.Node, key string, ctx *walk.VisitorCtx) {
			if ee.err != nil {
				return
			}

			ee.push(node.(*parser.StrLit).Text())
		})

	ee.addListener(walk.NodeAfterEvent(parser.N_EXPR_MEMBER),
		func(node parser.Node, key string, ctx *walk.VisitorCtx) {
			if ee.err != nil {
				return
			}

			prop := ee.pop()
			obj := ee.pop()

			v := GetProp(obj, prop)
			if v == nil {
				v = GetBuiltinProto(reflect.TypeOf(obj), prop)
				if IsNativeFn(v) {
					v = NewBuiltinFn(obj, v.(NativeFn))
				}
			}

			ee.push(v)
		})

	ee.addListener(walk.NodeAfterEvent(parser.N_EXPR_BIN),
		func(node parser.Node, key string, ctx *walk.VisitorCtx) {
			if ee.err != nil {
				return
			}

			n := node.(*parser.BinExpr)
			rhs := ee.pop()
			lhs := ee.pop()

			switch n.Op() {
			case parser.T_EQ, parser.T_EQ_S:
				ee.push(reflect.DeepEqual(lhs, rhs))
			case parser.T_ADD:
				ee.push(Add(lhs, rhs))
			case parser.T_SUB:
				ee.push(ToNum(lhs) - ToNum(rhs))
			case parser.T_MUL:
				ee.push(ToNum(lhs) * ToNum(rhs))
			case parser.T_DIV:
				r := ToNum(rhs)
				if r == 0 {
					ee.push(math.NaN())
				} else {
					ee.push(ToNum(lhs) - r)
				}
			case parser.T_AND:
				ee.push(ToBool(lhs) && ToBool(rhs))
			case parser.T_OR:
				ee.push(ToBool(lhs) || ToBool(rhs))
			}
		})

	ee.addListener(walk.NodeAfterEvent(parser.N_EXPR_UNARY),
		func(node parser.Node, key string, ctx *walk.VisitorCtx) {
			if ee.err != nil {
				return
			}

			n := node.(*parser.UnaryExpr)
			arg := ee.pop()

			switch n.Op() {
			case parser.T_NOT:
				ee.push(!ToBool(arg))
			default:
				ee.push(false)
			}
		})

	ee.addListener(walk.NodeAfterEvent(parser.N_EXPR_CALL),
		func(node parser.Node, key string, ctx *walk.VisitorCtx) {
			if ee.err != nil {
				return
			}

			n := node.(*parser.CallExpr)
			argsLen := len(n.Args())
			args := make([]interface{}, argsLen)

			for i := argsLen - 1; i >= 0; i-- {
				args[i] = ee.pop()
			}

			fn := ee.pop()
			if fo, ok := fn.(*BuiltinFn); ok {
				ee.push(fo.Call(NewArgs(args)))
			} else {
				ee.push(nil)
			}
		})
}

func (ee *ExprEvaluator) Release() {
	for nt, lis := range ee.listeners {
		walk.RemoveListener(&ee.walkCtx.Listeners, nt, lis)
	}
}

func (ee *ExprEvaluator) GetResult() interface{} {
	return ee.pop()
}

func Add(a, b interface{}) interface{} {
	_, sa := a.(string)
	_, sb := b.(string)
	if sa || sb {
		return ToStr(a) + ToStr(b)
	}

	return ToNum(a) + ToNum(b)
}

func ToNum(v interface{}) float64 {
	switch vv := v.(type) {
	case float64:
		return vv
	case string:
		if i, err := strconv.ParseFloat(vv, 64); err == nil {
			return i
		}
	}
	return math.NaN()
}

func ToStr(v interface{}) string {
	switch vv := v.(type) {
	case float64:
		return strconv.FormatFloat(vv, 'f', -1, 64)
	case string:
		return vv
	case bool:
		if vv {
			return "true"
		}
		return "false"
	}
	return fmt.Sprintf("%v", v)
}

func ToBool(v interface{}) bool {
	switch vv := v.(type) {
	case float64:
		return vv != 0
	case string:
		return vv != ""
	case bool:
		return vv
	case nil:
		return false
	}
	return true
}

func GetProp(obj, prop interface{}) interface{} {
	if obj == nil || prop == nil {
		return nil
	}

	p := ToStr(prop)
	switch ov := obj.(type) {
	case []interface{}:
		return GetElem(obj, prop)
	case map[string]interface{}:
		return ov[p]
	default:
		return nil
	}
}

func GetBuiltinProto(typ reflect.Type, prop interface{}) interface{} {
	if prop == nil {
		return nil
	}
	p := ToStr(prop)
	fns := builtinProto[typ]
	if fns == nil {
		return nil
	}
	if f, ok := fns[p]; ok {
		return f
	}
	return nil
}

func GetElem(obj, prop interface{}) interface{} {
	if obj == nil || prop == nil {
		return nil
	}

	p := ToNum(prop)
	if math.IsNaN(p) || p < 0 {
		return nil
	}

	arr, ok := obj.([]interface{})
	if !ok {
		return nil
	}

	i := int(p)
	if i > len(arr)-1 {
		return nil
	}
	return arr[i]
}

type Args struct {
	args []interface{}
}

func NewArgs(args []interface{}) *Args {
	return &Args{args}
}

func (a *Args) Get(i int) interface{} {
	if i < len(a.args) {
		return a.args[i]
	}
	return nil
}

type NativeFn = func(this interface{}, args *Args) interface{}

func IsNativeFn(o interface{}) bool {
	_, ok := o.(NativeFn)
	return ok
}

type BuiltinFn struct {
	this interface{}
	fn   NativeFn
}

func NewBuiltinFn(this interface{}, fn NativeFn) *BuiltinFn {
	return &BuiltinFn{this, fn}
}

func (f *BuiltinFn) Call(args *Args) interface{} {
	return f.fn(f.this, args)
}
