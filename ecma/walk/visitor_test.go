package walk

import (
	"testing"

	"github.com/hsiaosiyuan0/mole/ecma/parser"
	"github.com/hsiaosiyuan0/mole/span"
	. "github.com/hsiaosiyuan0/mole/util"
)

func newParser(code string, opts *parser.ParserOpts) *parser.Parser {
	if opts == nil {
		opts = parser.NewParserOpts()
	}
	s := span.NewSource("", code)
	return parser.NewParser(s, opts)
}

func compile(code string, opts *parser.ParserOpts) (*parser.Parser, parser.Node, *parser.SymTab, error) {
	p := newParser(code, opts)
	ast, err := p.Prog()
	if err != nil {
		return nil, nil, nil, err
	}
	return p, ast, p.Symtab(), nil
}

func TestVisitor(t *testing.T) {
	_, ast, symtab, err := compile("a + b - c", nil)
	AssertEqual(t, nil, err, "should pass")

	names := make([]string, 0)

	ctx := NewWalkCtx(ast, symtab)
	SetVisitor(&ctx.Visitors, N_NAME, func(node parser.Node, key string, ctx *VisitorCtx) {
		n := node.(*parser.Ident)
		names = append(names, n.Val())
	})

	VisitNode(ast, "", ctx.VisitorCtx())
	AssertEqual(t, 3, len(names), "should pass")
}

func TestVisitorFnScope(t *testing.T) {
	_, ast, symtab, err := compile(`
function a() {
}
  `, nil)
	AssertEqual(t, nil, err, "should pass")

	ctx := NewWalkCtx(ast, symtab)

	fnScopeId := 0
	SetVisitor(&ctx.Visitors, N_STMT_BLOCK_BEFORE, func(node parser.Node, key string, ctx *VisitorCtx) {
		fnScopeId = ctx.ScopeId()
	})

	VisitNode(ast, "", ctx.VisitorCtx())
	AssertEqual(t, 1, fnScopeId, "should pass")
}

func TestVisitorFnNestedScope(t *testing.T) {
	_, ast, symtab, err := compile(`
function a() {
  function b() {
  }
}
  `, nil)
	AssertEqual(t, nil, err, "should pass")

	ctx := NewWalkCtx(ast, symtab)

	fnScopeId := 0
	SetVisitor(&ctx.Visitors, N_STMT_BLOCK_AFTER, func(node parser.Node, key string, ctx *VisitorCtx) {
		if ctx.ScopeId() > fnScopeId {
			fnScopeId = ctx.ScopeId()
		}
	})

	VisitNode(ast, "", ctx.VisitorCtx())
	AssertEqual(t, 2, fnScopeId, "should pass")
}

func TestVisitorScopeBlock(t *testing.T) {
	_, ast, symtab, err := compile(`
{
  {
  }
}
  `, nil)
	AssertEqual(t, nil, err, "should pass")

	ctx := NewWalkCtx(ast, symtab)

	fnScopeId := 0
	SetVisitor(&ctx.Visitors, N_STMT_BLOCK_BEFORE, func(node parser.Node, key string, ctx *VisitorCtx) {
		if ctx.ScopeId() > fnScopeId {
			fnScopeId = ctx.ScopeId()
		}
	})

	VisitNode(ast, "", ctx.VisitorCtx())
	AssertEqual(t, 2, fnScopeId, "should pass")
}

func TestVisitorScopeWhile(t *testing.T) {
	_, ast, symtab, err := compile(`
while(1) {
  let a = 1
}
  `, nil)
	AssertEqual(t, nil, err, "should pass")

	ctx := NewWalkCtx(ast, symtab)

	scopeId := 0
	SetVisitor(&ctx.Visitors, N_STMT_BLOCK_BEFORE, func(node parser.Node, key string, ctx *VisitorCtx) {
		scopeId = ctx.ScopeId()
	})

	VisitNode(ast, "", ctx.VisitorCtx())
	AssertEqual(t, 2, scopeId, "should pass")

	AssertEqual(t, true, symtab.Scopes[2].HasName("a"), "should pass")
}

func TestVisitorScopeArrowFn(t *testing.T) {
	_, ast, symtab, err := compile(`
let a = () => {
  let b = 1
}
  `, nil)
	AssertEqual(t, nil, err, "should pass")

	ctx := NewWalkCtx(ast, symtab)

	fnScopeId := 0
	SetVisitor(&ctx.Visitors, N_STMT_BLOCK_BEFORE, func(node parser.Node, key string, ctx *VisitorCtx) {
		if ctx.ScopeId() > fnScopeId {
			fnScopeId = ctx.ScopeId()
		}
	})

	VisitNode(ast, "", ctx.VisitorCtx())
	AssertEqual(t, 1, fnScopeId, "should pass")

	AssertEqual(t, true, symtab.Scopes[1].HasName("b"), "should pass")
}

func TestVisitorScopeFnSymbol(t *testing.T) {
	_, ast, symtab, err := compile(`
function a() {
  let c = 1
  function b() {
    let d = 1
  }
}
  `, nil)
	AssertEqual(t, nil, err, "should pass")

	ctx := NewWalkCtx(ast, symtab)

	fnScopeId := 0
	SetVisitor(&ctx.Visitors, N_STMT_BLOCK_BEFORE, func(node parser.Node, key string, ctx *VisitorCtx) {
		if ctx.ScopeId() > fnScopeId {
			fnScopeId = ctx.ScopeId()
		}
	})

	VisitNode(ast, "", ctx.VisitorCtx())
	AssertEqual(t, 2, fnScopeId, "should pass")

	AssertEqual(t, true, symtab.Scopes[0].HasName("a"), "should pass")
	AssertEqual(t, true, symtab.Scopes[1].HasName("c"), "should pass")
	AssertEqual(t, true, symtab.Scopes[1].HasName("b"), "should pass")
	AssertEqual(t, true, symtab.Scopes[2].HasName("d"), "should pass")
}

func TestVisitorScopeFor(t *testing.T) {
	_, ast, symtab, err := compile(`
for(let a = 1;;){}
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ctx := NewWalkCtx(ast, symtab)

	fnScopeId := 0
	SetVisitor(&ctx.Visitors, N_STMT_BLOCK, func(node parser.Node, key string, ctx *VisitorCtx) {
		if ctx.ScopeId() > fnScopeId {
			fnScopeId = ctx.ScopeId()
		}
	})

	VisitNode(ast, "", ctx.VisitorCtx())
	AssertEqual(t, 1, fnScopeId, "should pass")

	AssertEqual(t, true, symtab.Scopes[1].HasName("a"), "should pass")
}

func TestVisitorScopeSwitch(t *testing.T) {
	_, ast, symtab, err := compile(`
switch(1) {
  case 1: let a = 1
  case 2: console.log(a)
}
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ctx := NewWalkCtx(ast, symtab)

	fnScopeId := 0
	SetVisitor(&ctx.Visitors, N_STMT_SWITCH_BEFORE, func(node parser.Node, key string, ctx *VisitorCtx) {
		if ctx.ScopeId() > fnScopeId {
			fnScopeId = ctx.ScopeId()
		}
	})

	VisitNode(ast, "", ctx.VisitorCtx())
	AssertEqual(t, 1, fnScopeId, "should pass")

	AssertEqual(t, true, symtab.Scopes[1].HasName("a"), "should pass")
}

func TestVisitorScopeSwitchCase(t *testing.T) {
	_, ast, symtab, err := compile(`
switch(1) {
  case 1: { let a = 1 }
  case 2: console.log(a)
}
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ctx := NewWalkCtx(ast, symtab)

	fnScopeId := 0
	SetVisitor(&ctx.Visitors, N_STMT_BLOCK_BEFORE, func(node parser.Node, key string, ctx *VisitorCtx) {
		if ctx.ScopeId() > fnScopeId {
			fnScopeId = ctx.ScopeId()
		}
	})

	VisitNode(ast, "", ctx.VisitorCtx())
	AssertEqual(t, 2, fnScopeId, "should pass")

	AssertEqual(t, true, symtab.Scopes[2].HasName("a"), "should pass")
	AssertEqual(t, false, symtab.Scopes[1].HasName("a"), "should pass")
}

func TestVisitorDoWhile(t *testing.T) {
	_, ast, symtab, err := compile(`
do {
  let a = 1
} while(1)
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ctx := NewWalkCtx(ast, symtab)

	scopeId := 0
	SetVisitor(&ctx.Visitors, N_STMT_BLOCK_BEFORE, func(node parser.Node, key string, ctx *VisitorCtx) {
		if ctx.ScopeId() > scopeId {
			scopeId = ctx.ScopeId()
		}
	})

	VisitNode(ast, "", ctx.VisitorCtx())
	AssertEqual(t, 2, scopeId, "should pass")

	AssertEqual(t, true, symtab.Scopes[2].HasName("a"), "should pass")
}

func TestVisitorScopeTry(t *testing.T) {
	_, ast, symtab, err := compile(`
try {
  let a = 1
} catch(e) {
  let b = 2
} finally {
  let c = 3
}
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ctx := NewWalkCtx(ast, symtab)

	fnScopeId := 0
	SetVisitor(&ctx.Visitors, N_STMT_BLOCK_BEFORE, func(node parser.Node, key string, ctx *VisitorCtx) {
		if ctx.ScopeId() > fnScopeId {
			fnScopeId = ctx.ScopeId()
		}
	})

	VisitNode(ast, "", ctx.VisitorCtx())
	AssertEqual(t, 3, fnScopeId, "should pass")

	AssertEqual(t, true, symtab.Scopes[1].HasName("a"), "should pass")
	AssertEqual(t, true, symtab.Scopes[2].HasName("b"), "should pass")
	AssertEqual(t, true, symtab.Scopes[3].HasName("c"), "should pass")
}

func TestVisitorScopeClass(t *testing.T) {
	_, ast, symtab, err := compile(`
class A {}
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ctx := NewWalkCtx(ast, symtab)

	fnScopeId := 0
	SetVisitor(&ctx.Visitors, N_CLASS_BODY_BEFORE, func(node parser.Node, key string, ctx *VisitorCtx) {
		if ctx.ScopeId() > fnScopeId {
			fnScopeId = ctx.ScopeId()
		}
	})

	VisitNode(ast, "", ctx.VisitorCtx())
	AssertEqual(t, 1, fnScopeId, "should pass")
}

func TestVisitorScopeMethod(t *testing.T) {
	_, ast, symtab, err := compile(`
class A {
  f(a){
    let b
  }
}
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ctx := NewWalkCtx(ast, symtab)

	fnScopeId := 0
	SetVisitor(&ctx.Visitors, N_STMT_BLOCK_BEFORE, func(node parser.Node, key string, ctx *VisitorCtx) {
		if ctx.ScopeId() > fnScopeId {
			fnScopeId = ctx.ScopeId()
		}
	})

	VisitNode(ast, "", ctx.VisitorCtx())
	AssertEqual(t, 2, fnScopeId, "should pass")
	AssertEqual(t, true, symtab.Scopes[2].HasName("a"), "should pass")
	AssertEqual(t, true, symtab.Scopes[2].HasName("b"), "should pass")
}

func TestListenerDoWhile(t *testing.T) {
	_, ast, symtab, err := compile(`
do {
  let a = 1
} while(1)
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ctx := NewWalkCtx(ast, symtab)

	scopeId := 0
	// also test `AddBeforeListener`
	AddBeforeListener(&ctx.Listeners, &Listener{
		Id: "BeforeListener",
		Handle: func(node parser.Node, key string, ctx *VisitorCtx) {
			if ctx.ScopeId() > scopeId {
				scopeId = ctx.ScopeId()
			}
		},
	})

	idName := ""
	AddListener(&ctx.Listeners, N_NAME_AFTER, &Listener{
		Id: "N_NAME_AFTER",
		Handle: func(node parser.Node, key string, ctx *VisitorCtx) {
			idName = node.(*parser.Ident).Val()
		},
	})

	VisitNode(ast, "", ctx.VisitorCtx())
	AssertEqual(t, 2, scopeId, "should pass")

	AssertEqual(t, true, symtab.Scopes[2].HasName("a"), "should pass")
	AssertEqual(t, "a", idName, "should pass")
}

func TestVisitComments(t *testing.T) {
	p, ast, symtab, err := compile(`
/* 1 */
let a, b;

/* 2 */
console.log(a == b); /* 3 */

/* 4 */ /* 5 */

console.log(a == b); /* 6 */
/* 7 */
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ctx := NewWalkCtx(ast, symtab)

	cmts := []span.Range{}
	for nt := range StmtNodeTypes {
		AddNodeBeforeListener(&ctx.Listeners, nt, &Listener{
			Id: "BeforeListener",
			Handle: func(node parser.Node, key string, ctx *VisitorCtx) {
				if cs := p.PrevCmts(node); cs != nil {
					cmts = append(cmts, cs...)
				}
				if cs := p.PostCmts(node); cs != nil {
					cmts = append(cmts, cs...)
				}
			},
		})
	}

	VisitNode(ast, "", ctx.VisitorCtx())

	AssertEqual(t, "/* 1 */", p.RngText(cmts[0]), "should pass")
	AssertEqual(t, "/* 2 */", p.RngText(cmts[1]), "should pass")
	AssertEqual(t, "/* 3 */", p.RngText(cmts[2]), "should pass")
	AssertEqual(t, "/* 4 */", p.RngText(cmts[3]), "should pass")
	AssertEqual(t, "/* 5 */", p.RngText(cmts[4]), "should pass")
	AssertEqual(t, "/* 6 */", p.RngText(cmts[5]), "should pass")
}
