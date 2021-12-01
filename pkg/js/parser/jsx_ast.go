package parser

type JSXIdent struct {
	typ   NodeType
	loc   *Loc
	val   string
	extra *ExprExtra
}

func (n *JSXIdent) Type() NodeType {
	return n.typ
}

func (n *JSXIdent) Loc() *Loc {
	return n.loc
}

func (n *JSXIdent) Extra() interface{} {
	return n.extra
}

func (n *JSXIdent) setExtra(ext interface{}) {
	n.extra = ext.(*ExprExtra)
}

func (n *JSXIdent) Text() string {
	return n.val
}

type JSXNsName struct {
	typ  NodeType
	loc  *Loc
	ns   Node
	name Node
}

func (n *JSXNsName) Type() NodeType {
	return n.typ
}

func (n *JSXNsName) Loc() *Loc {
	return n.loc
}

func (n *JSXNsName) Extra() interface{} {
	return nil
}

func (n *JSXNsName) setExtra(ext interface{}) {
}

func (n *JSXNsName) NS() string {
	return n.ns.(*JSXIdent).Text()
}

func (n *JSXNsName) Name() string {
	return n.name.(*JSXIdent).Text()
}

type JSXMemberExpr struct {
	typ  NodeType
	loc  *Loc
	obj  Node
	prop Node
}

func (n *JSXMemberExpr) Type() NodeType {
	return n.typ
}

func (n *JSXMemberExpr) Loc() *Loc {
	return n.loc
}

func (n *JSXMemberExpr) Extra() interface{} {
	return nil
}

func (n *JSXMemberExpr) setExtra(ext interface{}) {
}

func (n *JSXMemberExpr) Obj() Node {
	return n.obj
}

func (n *JSXMemberExpr) Prop() Node {
	return n.prop
}

type JSXOpen struct {
	typ NodeType
	loc *Loc
	// JSXIdentifier | JSXMemberExpression | JSXNamespacedName
	// `JSXNamespacedName` is a part of the JSX spec though it's
	// not used in the React implementation: https://github.com/facebook/jsx/issues/13
	name    Node
	nameStr string
	attrs   []Node
	closed  bool
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

func (n *JSXOpen) Closed() bool {
	return n.closed
}

type JSXClose struct {
	typ     NodeType
	loc     *Loc
	name    Node
	nameStr string
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

func (n *JSXText) Raw() string {
	return n.loc.Text()
}

type JSXAttr struct {
	typ     NodeType
	loc     *Loc
	name    Node
	nameStr string
	val     Node
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

type JSXEmpty struct {
	typ NodeType
	loc *Loc
}

func (n *JSXEmpty) Type() NodeType {
	return n.typ
}

func (n *JSXEmpty) Loc() *Loc {
	return n.loc
}

func (n *JSXEmpty) Extra() interface{} {
	return nil
}

func (n *JSXEmpty) setExtra(ext interface{}) {
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

func (n *JSXElem) IsFragment() bool {
	return n.open.(*JSXOpen).name == nil
}

type JSXExprSpan struct {
	typ  NodeType
	loc  *Loc
	expr Node
}

func (n *JSXExprSpan) Type() NodeType {
	return n.typ
}

func (n *JSXExprSpan) Loc() *Loc {
	return n.loc
}

func (n *JSXExprSpan) Extra() interface{} {
	return nil
}

func (n *JSXExprSpan) setExtra(ext interface{}) {
}

func (n *JSXExprSpan) Expr() Node {
	return n.expr
}
