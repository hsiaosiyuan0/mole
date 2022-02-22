package walk

import (
	"testing"

	"github.com/hsiaosiyuan0/mole/ecma/parser"
	. "github.com/hsiaosiyuan0/mole/fuzz"
	"github.com/hsiaosiyuan0/mole/span"
)

func newParser(code string, opts *parser.ParserOpts) *parser.Parser {
	if opts == nil {
		opts = parser.NewParserOpts()
	}
	s := span.NewSource("", code)
	return parser.NewParser(s, opts)
}

func compile(code string, opts *parser.ParserOpts) (parser.Node, *parser.SymTab, error) {
	p := newParser(code, opts)
	ast, err := p.Prog()
	if err != nil {
		return nil, nil, err
	}
	return ast, p.Symtab(), nil
}

func TestVisitor(t *testing.T) {
	ast, symtab, err := compile("a + b - c", nil)
	AssertEqual(t, nil, err, "should pass")

	names := make([]string, 0)

	ctx := NewWalkCtx(ast, symtab)
	SetVisitor(&ctx.Visitors, VK_NAME, func(node parser.Node, key string, ctx *WalkCtx) {
		n := node.(*parser.Ident)
		names = append(names, n.Text())
	})

	VisitNode(ast, "", ctx)
	AssertEqual(t, 3, len(names), "should pass")
}

func TestVisitorFnScope(t *testing.T) {
	ast, symtab, err := compile(`
function a() {
}
  `, nil)
	AssertEqual(t, nil, err, "should pass")

	ctx := NewWalkCtx(ast, symtab)

	fnScopeId := 0
	SetVisitor(&ctx.Visitors, VK_STMT_BLOCK_BEFORE, func(node parser.Node, key string, ctx *WalkCtx) {
		fnScopeId = ctx.ScopeId()
	})

	VisitNode(ast, "", ctx)
	AssertEqual(t, 1, fnScopeId, "should pass")
}

func TestVisitorFnNestedScope(t *testing.T) {
	ast, symtab, err := compile(`
function a() {
  function b() {
  }
}
  `, nil)
	AssertEqual(t, nil, err, "should pass")

	ctx := NewWalkCtx(ast, symtab)

	fnScopeId := 0
	SetVisitor(&ctx.Visitors, VK_STMT_BLOCK_AFTER, func(node parser.Node, key string, ctx *WalkCtx) {
		if ctx.ScopeId() > fnScopeId {
			fnScopeId = ctx.ScopeId()
		}
	})

	VisitNode(ast, "", ctx)
	AssertEqual(t, 2, fnScopeId, "should pass")
}

func TestVisitorScopeBlock(t *testing.T) {
	ast, symtab, err := compile(`
{
  {
  }
}
  `, nil)
	AssertEqual(t, nil, err, "should pass")

	ctx := NewWalkCtx(ast, symtab)

	fnScopeId := 0
	SetVisitor(&ctx.Visitors, VK_STMT_BLOCK_BEFORE, func(node parser.Node, key string, ctx *WalkCtx) {
		if ctx.ScopeId() > fnScopeId {
			fnScopeId = ctx.ScopeId()
		}
	})

	VisitNode(ast, "", ctx)
	AssertEqual(t, 2, fnScopeId, "should pass")
}

func TestVisitorScopeWhile(t *testing.T) {
	ast, symtab, err := compile(`
while(1) {
  let a = 1
}
  `, nil)
	AssertEqual(t, nil, err, "should pass")

	ctx := NewWalkCtx(ast, symtab)

	fnScopeId := 0
	SetVisitor(&ctx.Visitors, VK_STMT_BLOCK_BEFORE, func(node parser.Node, key string, ctx *WalkCtx) {
		if ctx.ScopeId() > fnScopeId {
			fnScopeId = ctx.ScopeId()
		}
	})

	VisitNode(ast, "", ctx)
	AssertEqual(t, 1, fnScopeId, "should pass")

	AssertEqual(t, true, symtab.Scopes[1].HasName("a"), "should pass")
}

func TestVisitorScopeArrowFn(t *testing.T) {
	ast, symtab, err := compile(`
let a = () => {
  let b = 1
}
  `, nil)
	AssertEqual(t, nil, err, "should pass")

	ctx := NewWalkCtx(ast, symtab)

	fnScopeId := 0
	SetVisitor(&ctx.Visitors, VK_STMT_BLOCK_BEFORE, func(node parser.Node, key string, ctx *WalkCtx) {
		if ctx.ScopeId() > fnScopeId {
			fnScopeId = ctx.ScopeId()
		}
	})

	VisitNode(ast, "", ctx)
	AssertEqual(t, 1, fnScopeId, "should pass")

	AssertEqual(t, true, symtab.Scopes[1].HasName("b"), "should pass")
}

func TestVisitorScopeFnSymbol(t *testing.T) {
	ast, symtab, err := compile(`
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
	SetVisitor(&ctx.Visitors, VK_STMT_BLOCK_BEFORE, func(node parser.Node, key string, ctx *WalkCtx) {
		if ctx.ScopeId() > fnScopeId {
			fnScopeId = ctx.ScopeId()
		}
	})

	VisitNode(ast, "", ctx)
	AssertEqual(t, 2, fnScopeId, "should pass")

	AssertEqual(t, true, symtab.Scopes[0].HasName("a"), "should pass")
	AssertEqual(t, true, symtab.Scopes[1].HasName("c"), "should pass")
	AssertEqual(t, true, symtab.Scopes[1].HasName("b"), "should pass")
	AssertEqual(t, true, symtab.Scopes[2].HasName("d"), "should pass")
}

func TestVisitorScopeFor(t *testing.T) {
	ast, symtab, err := compile(`
for(let a = 1;;){}
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ctx := NewWalkCtx(ast, symtab)

	fnScopeId := 0
	SetVisitor(&ctx.Visitors, VK_STMT_BLOCK, func(node parser.Node, key string, ctx *WalkCtx) {
		if ctx.ScopeId() > fnScopeId {
			fnScopeId = ctx.ScopeId()
		}
	})

	VisitNode(ast, "", ctx)
	AssertEqual(t, 1, fnScopeId, "should pass")

	AssertEqual(t, true, symtab.Scopes[1].HasName("a"), "should pass")
}

func TestVisitorScopeSwitch(t *testing.T) {
	ast, symtab, err := compile(`
switch(1) {
  case 1: let a = 1
  case 2: console.log(a)
}
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ctx := NewWalkCtx(ast, symtab)

	fnScopeId := 0
	SetVisitor(&ctx.Visitors, VK_STMT_SWITCH_BEFORE, func(node parser.Node, key string, ctx *WalkCtx) {
		if ctx.ScopeId() > fnScopeId {
			fnScopeId = ctx.ScopeId()
		}
	})

	VisitNode(ast, "", ctx)
	AssertEqual(t, 1, fnScopeId, "should pass")

	AssertEqual(t, true, symtab.Scopes[1].HasName("a"), "should pass")
}

func TestVisitorScopeSwitchCase(t *testing.T) {
	ast, symtab, err := compile(`
switch(1) {
  case 1: { let a = 1 }
  case 2: console.log(a)
}
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ctx := NewWalkCtx(ast, symtab)

	fnScopeId := 0
	SetVisitor(&ctx.Visitors, VK_STMT_BLOCK_BEFORE, func(node parser.Node, key string, ctx *WalkCtx) {
		if ctx.ScopeId() > fnScopeId {
			fnScopeId = ctx.ScopeId()
		}
	})

	VisitNode(ast, "", ctx)
	AssertEqual(t, 2, fnScopeId, "should pass")

	AssertEqual(t, true, symtab.Scopes[2].HasName("a"), "should pass")
	AssertEqual(t, false, symtab.Scopes[1].HasName("a"), "should pass")
}

func TestVisitorDoWhile(t *testing.T) {
	ast, symtab, err := compile(`
do {
  let a = 1
} while(1)
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ctx := NewWalkCtx(ast, symtab)

	fnScopeId := 0
	SetVisitor(&ctx.Visitors, VK_STMT_BLOCK_BEFORE, func(node parser.Node, key string, ctx *WalkCtx) {
		if ctx.ScopeId() > fnScopeId {
			fnScopeId = ctx.ScopeId()
		}
	})

	VisitNode(ast, "", ctx)
	AssertEqual(t, 1, fnScopeId, "should pass")

	AssertEqual(t, true, symtab.Scopes[1].HasName("a"), "should pass")
}

func TestVisitorScopeTry(t *testing.T) {
	ast, symtab, err := compile(`
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
	SetVisitor(&ctx.Visitors, VK_STMT_BLOCK_BEFORE, func(node parser.Node, key string, ctx *WalkCtx) {
		if ctx.ScopeId() > fnScopeId {
			fnScopeId = ctx.ScopeId()
		}
	})

	VisitNode(ast, "", ctx)
	AssertEqual(t, 3, fnScopeId, "should pass")

	AssertEqual(t, true, symtab.Scopes[1].HasName("a"), "should pass")
	AssertEqual(t, true, symtab.Scopes[2].HasName("b"), "should pass")
	AssertEqual(t, true, symtab.Scopes[3].HasName("c"), "should pass")
}

func TestVisitorScopeClass(t *testing.T) {
	ast, symtab, err := compile(`
class A {}
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ctx := NewWalkCtx(ast, symtab)

	fnScopeId := 0
	SetVisitor(&ctx.Visitors, VK_CLASS_BODY_BEFORE, func(node parser.Node, key string, ctx *WalkCtx) {
		if ctx.ScopeId() > fnScopeId {
			fnScopeId = ctx.ScopeId()
		}
	})

	VisitNode(ast, "", ctx)
	AssertEqual(t, 1, fnScopeId, "should pass")
}

func TestVisitorScopeMethod(t *testing.T) {
	ast, symtab, err := compile(`
class A {
  f(a){
    let b
  }
}
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ctx := NewWalkCtx(ast, symtab)

	fnScopeId := 0
	SetVisitor(&ctx.Visitors, VK_STMT_BLOCK_BEFORE, func(node parser.Node, key string, ctx *WalkCtx) {
		if ctx.ScopeId() > fnScopeId {
			fnScopeId = ctx.ScopeId()
		}
	})

	VisitNode(ast, "", ctx)
	AssertEqual(t, 2, fnScopeId, "should pass")
	AssertEqual(t, true, symtab.Scopes[2].HasName("a"), "should pass")
	AssertEqual(t, true, symtab.Scopes[2].HasName("b"), "should pass")
}

func TestListenerDoWhile(t *testing.T) {
	ast, symtab, err := compile(`
do {
  let a = 1
} while(1)
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ctx := NewWalkCtx(ast, symtab)

	fnScopeId := 0
	AddListener(&ctx.Listeners, LK_STMT_BLOCK_BEFORE, func(node parser.Node, key string, ctx *WalkCtx) {
		if ctx.ScopeId() > fnScopeId {
			fnScopeId = ctx.ScopeId()
		}
	})

	idName := ""
	AddListener(&ctx.Listeners, LK_NAME, func(node parser.Node, key string, ctx *WalkCtx) {
		idName = node.(*parser.Ident).Text()
	})

	VisitNode(ast, "", ctx)
	AssertEqual(t, 1, fnScopeId, "should pass")

	AssertEqual(t, true, symtab.Scopes[1].HasName("a"), "should pass")
	AssertEqual(t, "a", idName, "should pass")
}
