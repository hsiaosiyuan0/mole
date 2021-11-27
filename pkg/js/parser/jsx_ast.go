package parser

type JSXOpen struct {
	typ NodeType
	loc *Loc
	// JSXIdentifier | JSXMemberExpression | JSXNamespacedName
	// `JSXNamespacedName` is a part of the JSX spec though it's
	// not used in the React implementation: https://github.com/facebook/jsx/issues/13
	name  Node
	attrs []Node
}

func (n *JSXOpen) Type() NodeType {
	return n.typ
}

func (n *JSXOpen) Loc() *Loc {
	return n.loc
}

func (n *JSXOpen) Extra() interface{} {
	return nil
}

func (n *JSXOpen) setExtra(ext interface{}) {
}

func (n *JSXOpen) Name() Node {
	return n.name
}

func (n *JSXOpen) Attrs() []Node {
	return n.attrs
}

type JSXClose struct {
	typ  NodeType
	loc  *Loc
	name Node
}

func (n *JSXClose) Type() NodeType {
	return n.typ
}

func (n *JSXClose) Loc() *Loc {
	return n.loc
}

func (n *JSXClose) Extra() interface{} {
	return nil
}

func (n *JSXClose) setExtra(ext interface{}) {
}

func (n *JSXClose) Name() Node {
	return n.name
}

type JSXText struct {
	typ NodeType
	loc *Loc
	val string
}

func (n *JSXText) Type() NodeType {
	return n.typ
}

func (n *JSXText) Loc() *Loc {
	return n.loc
}

func (n *JSXText) Extra() interface{} {
	return nil
}

func (n *JSXText) setExtra(ext interface{}) {
}

func (n *JSXText) Value() string {
	return n.val
}

type JSXAttr struct {
	typ  NodeType
	loc  *Loc
	name Node
	val  Node
}

func (n *JSXAttr) Type() NodeType {
	return n.typ
}

func (n *JSXAttr) Loc() *Loc {
	return n.loc
}

func (n *JSXAttr) Extra() interface{} {
	return nil
}

func (n *JSXAttr) setExtra(ext interface{}) {
}

func (n *JSXAttr) Name() Node {
	return n.name
}

func (n *JSXAttr) Value() Node {
	return n.val
}

type JSXSpreadAttr struct {
	typ NodeType
	loc *Loc
	arg Node
}

func (n *JSXSpreadAttr) Type() NodeType {
	return n.typ
}

func (n *JSXSpreadAttr) Loc() *Loc {
	return n.loc
}

func (n *JSXSpreadAttr) Extra() interface{} {
	return nil
}

func (n *JSXSpreadAttr) setExtra(ext interface{}) {
}

func (n *JSXSpreadAttr) Arg() Node {
	return n.arg
}

type JSXSpreadChild struct {
	typ  NodeType
	loc  *Loc
	expr Node
}

func (n *JSXSpreadChild) Type() NodeType {
	return n.typ
}

func (n *JSXSpreadChild) Loc() *Loc {
	return n.loc
}

func (n *JSXSpreadChild) Extra() interface{} {
	return nil
}

func (n *JSXSpreadChild) setExtra(ext interface{}) {
}

func (n *JSXSpreadChild) Expr() Node {
	return n.expr
}

// https://github.com/facebook/jsx/blob/main/AST.md#jsx-element
type JSXElem struct {
	typ      NodeType
	loc      *Loc
	open     Node
	close    Node
	children []Node // [ JSXText | JSXExpressionContainer | JSXSpreadChild | JSXElement | JSXFragment ]
}

func (n *JSXElem) Type() NodeType {
	return n.typ
}

func (n *JSXElem) Loc() *Loc {
	return n.loc
}

func (n *JSXElem) Extra() interface{} {
	return nil
}

func (n *JSXElem) setExtra(_ interface{}) {
}

func (n *JSXElem) Open() Node {
	return n.open
}

func (n *JSXElem) Close() Node {
	return n.close
}

func (n *JSXElem) Children() []Node {
	return n.children
}
