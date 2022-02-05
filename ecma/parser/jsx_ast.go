package parser

type JsxIdent struct {
	typ        NodeType
	loc        *Loc
	val        string
	outerParen *Loc
	ti         *TypInfo
}

func (n *JsxIdent) Type() NodeType {
	return n.typ
}

func (n *JsxIdent) Loc() *Loc {
	return n.loc
}

func (n *JsxIdent) OuterParen() *Loc {
	return n.outerParen
}

func (n *JsxIdent) SetOuterParen(loc *Loc) {
	n.outerParen = loc
}

func (n *JsxIdent) Text() string {
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
	loc  *Loc
	ns   Node
	name Node
}

func (n *JsxNsName) Type() NodeType {
	return n.typ
}

func (n *JsxNsName) Loc() *Loc {
	return n.loc
}

func (n *JsxNsName) NS() string {
	return n.ns.(*JsxIdent).Text()
}

func (n *JsxNsName) Name() string {
	return n.name.(*JsxIdent).Text()
}

type JsxMemberExpr struct {
	typ  NodeType
	loc  *Loc
	obj  Node
	prop Node
	ti   *TypInfo
}

func (n *JsxMemberExpr) Type() NodeType {
	return n.typ
}

func (n *JsxMemberExpr) Loc() *Loc {
	return n.loc
}

func (n *JsxMemberExpr) Obj() Node {
	return n.obj
}

func (n *JsxMemberExpr) Prop() Node {
	return n.prop
}

func (n *JsxMemberExpr) TypInfo() *TypInfo {
	return n.ti
}

func (n *JsxMemberExpr) SetTypInfo(ti *TypInfo) {
	n.ti = ti
}

type JsxOpen struct {
	typ NodeType
	loc *Loc
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

func (n *JsxOpen) Loc() *Loc {
	return n.loc
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

type JsxClose struct {
	typ     NodeType
	loc     *Loc
	name    Node
	nameStr string
}

func (n *JsxClose) Type() NodeType {
	return n.typ
}

func (n *JsxClose) Loc() *Loc {
	return n.loc
}

func (n *JsxClose) Name() Node {
	return n.name
}

type JsxText struct {
	typ NodeType
	loc *Loc
	val string
}

func (n *JsxText) Type() NodeType {
	return n.typ
}

func (n *JsxText) Loc() *Loc {
	return n.loc
}

func (n *JsxText) Value() string {
	return n.val
}

func (n *JsxText) Raw() string {
	return n.loc.Text()
}

type JsxAttr struct {
	typ     NodeType
	loc     *Loc
	name    Node
	nameStr string
	val     Node
}

func (n *JsxAttr) Type() NodeType {
	return n.typ
}

func (n *JsxAttr) Loc() *Loc {
	return n.loc
}

func (n *JsxAttr) Name() Node {
	return n.name
}

func (n *JsxAttr) Value() Node {
	return n.val
}

type JsxSpreadAttr struct {
	typ NodeType
	loc *Loc
	arg Node
}

func (n *JsxSpreadAttr) Type() NodeType {
	return n.typ
}

func (n *JsxSpreadAttr) Loc() *Loc {
	return n.loc
}

func (n *JsxSpreadAttr) Arg() Node {
	return n.arg
}

type JsxSpreadChild struct {
	typ  NodeType
	loc  *Loc
	expr Node
}

func (n *JsxSpreadChild) Type() NodeType {
	return n.typ
}

func (n *JsxSpreadChild) Loc() *Loc {
	return n.loc
}

func (n *JsxSpreadChild) Expr() Node {
	return n.expr
}

type JsxEmpty struct {
	typ NodeType
	loc *Loc
}

func (n *JsxEmpty) Type() NodeType {
	return n.typ
}

func (n *JsxEmpty) Loc() *Loc {
	return n.loc
}

// https://github.com/facebook/jsx/blob/main/AST.md#jsx-element
type JsxElem struct {
	typ      NodeType
	loc      *Loc
	open     Node
	close    Node
	children []Node // [ JsxText | JsxExpressionContainer | JsxSpreadChild | JsxElement | JsxFragment ]
}

func (n *JsxElem) Type() NodeType {
	return n.typ
}

func (n *JsxElem) Loc() *Loc {
	return n.loc
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

type JsxExprSpan struct {
	typ  NodeType
	loc  *Loc
	expr Node
}

func (n *JsxExprSpan) Type() NodeType {
	return n.typ
}

func (n *JsxExprSpan) Loc() *Loc {
	return n.loc
}

func (n *JsxExprSpan) Expr() Node {
	return n.expr
}
