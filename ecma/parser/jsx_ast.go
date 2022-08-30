package parser

import span "github.com/hsiaosiyuan0/mole/span"

// grammar: https://facebook.github.io/jsx/

type JsxIdent struct {
	typ NodeType
	rng span.Range
	val string
	opa span.Range
	ti  *TypInfo
}

func (n *JsxIdent) Type() NodeType {
	return n.typ
}

func (n *JsxIdent) Range() span.Range {
	return n.rng
}

func (n *JsxIdent) OuterParen() span.Range {
	return n.opa
}

func (n *JsxIdent) SetOuterParen(rng span.Range) {
	n.opa = rng
}

func (n *JsxIdent) Val() string {
	return n.val
}

func (n *JsxIdent) TypInfo() *TypInfo {
	return n.ti
}

func (n *JsxIdent) SetTypInfo(ti *TypInfo) {
	n.ti = ti
}

type JsxNsName struct {
	typ  NodeType
	rng  span.Range
	ns   Node
	name Node
}

func (n *JsxNsName) Type() NodeType {
	return n.typ
}

func (n *JsxNsName) Range() span.Range {
	return n.rng
}

func (n *JsxNsName) NS() string {
	return n.ns.(*JsxIdent).Val()
}

func (n *JsxNsName) Name() string {
	return n.name.(*JsxIdent).Val()
}

// #[visitor(Obj,Prop)]
type JsxMember struct {
	typ  NodeType
	rng  span.Range
	obj  Node
	prop Node
	ti   *TypInfo
}

func (n *JsxMember) Type() NodeType {
	return n.typ
}

func (n *JsxMember) Range() span.Range {
	return n.rng
}

func (n *JsxMember) Obj() Node {
	return n.obj
}

func (n *JsxMember) Prop() Node {
	return n.prop
}

func (n *JsxMember) TypInfo() *TypInfo {
	return n.ti
}

func (n *JsxMember) SetTypInfo(ti *TypInfo) {
	n.ti = ti
}

// #[visitor(Name,Attrs)]
type JsxOpen struct {
	typ NodeType
	rng span.Range
	// JsxIdentifier | JsxMemberExpression | JsxNamespacedName
	// `JsxNamespacedName` is a part of the Jsx spec though it's
	// not used in the React implementation: https://github.com/facebook/jsx/issues/13
	name    Node
	nameStr string
	attrs   []Node
	closed  bool
}

func (n *JsxOpen) Type() NodeType {
	return n.typ
}

func (n *JsxOpen) Range() span.Range {
	return n.rng
}

func (n *JsxOpen) Name() Node {
	return n.name
}

func (n *JsxOpen) Attrs() []Node {
	return n.attrs
}

func (n *JsxOpen) Closed() bool {
	return n.closed
}

// #[visitor(Name)]
type JsxClose struct {
	typ     NodeType
	rng     span.Range
	name    Node
	nameStr string
}

func (n *JsxClose) Type() NodeType {
	return n.typ
}

func (n *JsxClose) Range() span.Range {
	return n.rng
}

func (n *JsxClose) Name() Node {
	return n.name
}

type JsxText struct {
	typ NodeType
	rng span.Range
	val string
}

func (n *JsxText) Type() NodeType {
	return n.typ
}

func (n *JsxText) Range() span.Range {
	return n.rng
}

func (n *JsxText) Val() string {
	return n.val
}

// #[visitor(Name,Val)]
type JsxAttr struct {
	typ     NodeType
	rng     span.Range
	name    Node
	nameStr string
	val     Node
}

func (n *JsxAttr) Type() NodeType {
	return n.typ
}

func (n *JsxAttr) Range() span.Range {
	return n.rng
}

func (n *JsxAttr) Name() Node {
	return n.name
}

func (n *JsxAttr) NameStr() string {
	return n.nameStr
}

func (n *JsxAttr) Val() Node {
	return n.val
}

// #[visitor(Arg)]
type JsxSpreadAttr struct {
	typ NodeType
	rng span.Range
	arg Node
}

func (n *JsxSpreadAttr) Type() NodeType {
	return n.typ
}

func (n *JsxSpreadAttr) Range() span.Range {
	return n.rng
}

func (n *JsxSpreadAttr) Arg() Node {
	return n.arg
}

// #[visitor(Expr)]
type JsxSpreadChild struct {
	typ  NodeType
	rng  span.Range
	expr Node
}

func (n *JsxSpreadChild) Type() NodeType {
	return n.typ
}

func (n *JsxSpreadChild) Range() span.Range {
	return n.rng
}

func (n *JsxSpreadChild) Expr() Node {
	return n.expr
}

type JsxEmpty struct {
	typ NodeType
	rng span.Range
}

func (n *JsxEmpty) Type() NodeType {
	return n.typ
}

func (n *JsxEmpty) Range() span.Range {
	return n.rng
}

// https://github.com/facebook/jsx/blob/main/AST.md#jsx-element
//
// #[visitor(Open,Children,Close)]
type JsxElem struct {
	typ      NodeType
	rng      span.Range
	open     Node
	close    Node
	children []Node // [ JsxText | JsxExpressionContainer | JsxSpreadChild | JsxElement | JsxFragment ]
}

func (n *JsxElem) Type() NodeType {
	return n.typ
}

func (n *JsxElem) Range() span.Range {
	return n.rng
}

func (n *JsxElem) Open() Node {
	return n.open
}

func (n *JsxElem) Close() Node {
	return n.close
}

func (n *JsxElem) Children() []Node {
	return n.children
}

func (n *JsxElem) IsFragment() bool {
	return n.open.(*JsxOpen).name == nil
}

// #[visitor(Expr)]
type JsxExprSpan struct {
	typ  NodeType
	rng  span.Range
	expr Node
}

func (n *JsxExprSpan) Type() NodeType {
	return n.typ
}

func (n *JsxExprSpan) Range() span.Range {
	return n.rng
}

func (n *JsxExprSpan) Expr() Node {
	return n.expr
}
