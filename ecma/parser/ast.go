package parser

import (
	span "github.com/hsiaosiyuan0/mole/span"
)

// AST nodes are referred to [ESTree](https://github.com/estree/estree/blob/master/es5.md) with some variants:
// - flattened struct is used instead of inheritance
// - fields are not fully described as they are claimed in ESTree, eg. the field names in this file are shorter then
//   their equivalent of ESTree. for the requirement what the ESTree compatible output is needed, use the `estree` package
//   to do the transformation
type Node interface {
	Type() NodeType
	Range() span.Range
}

// #[visitor(Body)]
type Prog struct {
	typ NodeType
	rng span.Range

	stmts []Node
}

func (n *Prog) Type() NodeType {
	return n.typ
}

func (n *Prog) Range() span.Range {
	return n.rng
}

func (n *Prog) Body() []Node {
	return n.stmts
}

// #[visitor(Expr)]
type ExprStmt struct {
	typ  NodeType
	rng  span.Range
	expr Node
	dir  bool
}

func (n *ExprStmt) Type() NodeType {
	return n.typ
}

func (n *ExprStmt) Range() span.Range {
	return n.rng
}

func (n *ExprStmt) Dir() bool {
	return n.dir
}

func (n *ExprStmt) Expr() Node {
	return n.expr
}

type EmptyStmt struct {
	typ NodeType
	rng span.Range
}

func (n *EmptyStmt) Type() NodeType {
	return n.typ
}

func (n *EmptyStmt) Range() span.Range {
	return n.rng
}

type InParenNode interface {
	OuterParen() span.Range
	SetOuterParen(span.Range)
}

func IsNodeInParen(node Node) bool {
	n, ok := node.(InParenNode)
	return ok && !n.OuterParen().Empty()
}

type NullLit struct {
	typ NodeType
	rng span.Range
	opa span.Range
	ti  *TypInfo
}

func (n *NullLit) Type() NodeType {
	return n.typ
}

func (n *NullLit) Range() span.Range {
	return n.rng
}

func (n *NullLit) SetOuterParen(rng span.Range) {
	n.opa = rng
}

func (n *NullLit) TypInfo() *TypInfo {
	return n.ti
}

func (n *NullLit) SetTypInfo(ti *TypInfo) {
	n.ti = ti
}

type BoolLit struct {
	typ NodeType
	rng span.Range
	val bool
	opa span.Range
	ti  *TypInfo
}

func (n *BoolLit) Type() NodeType {
	return n.typ
}

func (n *BoolLit) Range() span.Range {
	return n.rng
}

func (n *BoolLit) Val() bool {
	return n.val
}

func (n *BoolLit) OuterParen() span.Range {
	return n.opa
}

func (n *BoolLit) SetOuterParen(rng span.Range) {
	n.opa = rng
}

func (n *BoolLit) TypInfo() *TypInfo {
	return n.ti
}

func (n *BoolLit) SetTypInfo(ti *TypInfo) {
	n.ti = ti
}

type NumLit struct {
	typ NodeType
	rng span.Range
	opa span.Range
}

func (n *NumLit) Type() NodeType {
	return n.typ
}

func (n *NumLit) Range() span.Range {
	return n.rng
}

func (n *NumLit) OuterParen() span.Range {
	return n.opa
}

func (n *NumLit) SetOuterParen(rng span.Range) {
	n.opa = rng
}

type StrLit struct {
	typ   NodeType
	rng   span.Range
	val   string
	loSeq bool
	opa   span.Range
	ti    *TypInfo
}

func (n *StrLit) Type() NodeType {
	return n.typ
}

func (n *StrLit) Range() span.Range {
	return n.rng
}

func (n *StrLit) Val() string {
	return n.val
}

func (n *StrLit) OuterParen() span.Range {
	return n.opa
}

func (n *StrLit) SetOuterParen(rng span.Range) {
	n.opa = rng
}

func (n *StrLit) TypInfo() *TypInfo {
	return n.ti
}

func (n *StrLit) SetTypInfo(ti *TypInfo) {
	n.ti = ti
}

type RegLit struct {
	typ     NodeType
	rng     span.Range
	val     string
	pattern string
	flags   string
	opa     span.Range
	ti      *TypInfo
}

func (n *RegLit) Type() NodeType {
	return n.typ
}

func (n *RegLit) Range() span.Range {
	return n.rng
}

func (n *RegLit) Pattern() string {
	return n.pattern
}

func (n *RegLit) Flags() string {
	return n.flags
}

func (n *RegLit) OuterParen() span.Range {
	return n.opa
}

func (n *RegLit) SetOuterParen(rng span.Range) {
	n.opa = rng
}

func (n *RegLit) TypInfo() *TypInfo {
	return n.ti
}

func (n *RegLit) SetTypInfo(ti *TypInfo) {
	n.ti = ti
}

// #[visitor(Elems)]
type ArrLit struct {
	typ   NodeType
	rng   span.Range
	elems []Node
	opa   span.Range
	ti    *TypInfo
}

func (n *ArrLit) Type() NodeType {
	return n.typ
}

func (n *ArrLit) Range() span.Range {
	return n.rng
}

func (n *ArrLit) Elems() []Node {
	return n.elems
}

func (n *ArrLit) OuterParen() span.Range {
	return n.opa
}

func (n *ArrLit) SetOuterParen(rng span.Range) {
	n.opa = rng
}

func (n *ArrLit) TypInfo() *TypInfo {
	return n.ti
}

func (n *ArrLit) SetTypInfo(ti *TypInfo) {
	n.ti = ti
}

// #[visitor(Arg)]
type Spread struct {
	typ   NodeType
	rng   span.Range
	arg   Node
	tcLoc span.Range // range of the trailing comma
	opa   span.Range
	ti    *TypInfo
}

func (n *Spread) Type() NodeType {
	return n.typ
}

func (n *Spread) Range() span.Range {
	return n.rng
}

func (n *Spread) Arg() Node {
	return n.arg
}

func (n *Spread) OuterParen() span.Range {
	return n.opa
}

func (n *Spread) SetOuterParen(rng span.Range) {
	n.opa = rng
}

func (n *Spread) TypInfo() *TypInfo {
	return n.ti
}

func (n *Spread) SetTypInfo(ti *TypInfo) {
	n.ti = ti
}

// #[visitor(Props)]
type ObjLit struct {
	typ   NodeType
	rng   span.Range
	props []Node
	opa   span.Range
	ti    *TypInfo
}

func (n *ObjLit) Type() NodeType {
	return n.typ
}

func (n *ObjLit) Range() span.Range {
	return n.rng
}

func (n *ObjLit) Props() []Node {
	return n.props
}

func (n *ObjLit) OuterParen() span.Range {
	return n.opa
}

func (n *ObjLit) SetOuterParen(rng span.Range) {
	n.opa = rng
}

func (n *ObjLit) TypInfo() *TypInfo {
	return n.ti
}

func (n *ObjLit) SetTypInfo(ti *TypInfo) {
	n.ti = ti
}

type ACC_MOD uint8

const (
	ACC_MOD_NONE ACC_MOD = iota
	ACC_MOD_PUB
	ACC_MOD_PRI
	ACC_MOD_PRO
)

func (a ACC_MOD) String() string {
	switch a {
	case ACC_MOD_PUB:
		return "public"
	case ACC_MOD_PRI:
		return "private"
	case ACC_MOD_PRO:
		return "protected"
	}
	return ""
}

type NodeWithTypInfo interface {
	TypInfo() *TypInfo
	SetTypInfo(*TypInfo)
}

type Ident struct {
	typ            NodeType
	rng            span.Range
	val            string
	pvt            bool
	containsEscape bool
	opa            span.Range

	// consider below statements:
	// `export { if } from "a"` is legal
	// `export { if } ` is illegal
	// for reporting `if` is a keyword, firstly produce a
	// Ident with content `if` and flag it's a keyword by
	// setting this field to true, later report the `unexpected token`
	// error if the coming token is not `from`
	kw bool

	ti *TypInfo
}

func (n *Ident) Type() NodeType {
	return n.typ
}

func (n *Ident) Range() span.Range {
	return n.rng
}

func (n *Ident) Val() string {
	return n.val
}

func (n *Ident) IsPrivate() bool {
	return n.pvt
}

func (n *Ident) ContainsEscape() bool {
	return n.containsEscape
}

func (n *Ident) OuterParen() span.Range {
	return n.opa
}

func (n *Ident) SetOuterParen(rng span.Range) {
	n.opa = rng
}

func (n *Ident) TypInfo() *TypInfo {
	return n.ti
}

func (n *Ident) SetTypInfo(ti *TypInfo) {
	n.ti = ti
}

// #[visitor(Callee,Args)]
type NewExpr struct {
	typ    NodeType
	rng    span.Range
	callee Node
	args   []Node
	opa    span.Range
	ti     *TypInfo
}

func (n *NewExpr) Type() NodeType {
	return n.typ
}

func (n *NewExpr) Range() span.Range {
	return n.rng
}

func (n *NewExpr) Callee() Node {
	return n.callee
}

func (n *NewExpr) Args() []Node {
	return n.args
}

func (n *NewExpr) OuterParen() span.Range {
	return n.opa
}

func (n *NewExpr) SetOuterParen(rng span.Range) {
	n.opa = rng
}

func (n *NewExpr) TypInfo() *TypInfo {
	return n.ti
}

func (n *NewExpr) SetTypInfo(ti *TypInfo) {
	n.ti = ti
}

// #[visitor(Obj,Prop)]
type MemberExpr struct {
	typ      NodeType
	rng      span.Range
	obj      Node
	prop     Node
	compute  bool
	optional bool
	opa      span.Range
}

func (n *MemberExpr) Type() NodeType {
	return n.typ
}

func (n *MemberExpr) Range() span.Range {
	return n.rng
}

func (n *MemberExpr) Obj() Node {
	return n.obj
}

func (n *MemberExpr) Prop() Node {
	return n.prop
}

func (n *MemberExpr) Compute() bool {
	return n.compute
}

func (n *MemberExpr) Optional() bool {
	return n.optional
}

func (n *MemberExpr) OuterParen() span.Range {
	return n.opa
}

func (n *MemberExpr) SetOuterParen(rng span.Range) {
	n.opa = rng
}

// #[visitor(Callee,Args)]
type CallExpr struct {
	typ      NodeType
	rng      span.Range
	callee   Node
	args     []Node
	optional bool
	opa      span.Range
	ti       *TypInfo
}

func (n *CallExpr) Type() NodeType {
	return n.typ
}

func (n *CallExpr) Range() span.Range {
	return n.rng
}

func (n *CallExpr) Callee() Node {
	return n.callee
}

func (n *CallExpr) Args() []Node {
	return n.args
}

func (n *CallExpr) Optional() bool {
	return n.optional
}

func (n *CallExpr) OuterParen() span.Range {
	return n.opa
}

func (n *CallExpr) SetOuterParen(rng span.Range) {
	n.opa = rng
}

func (n *CallExpr) TypInfo() *TypInfo {
	return n.ti
}

func (n *CallExpr) SetTypInfo(ti *TypInfo) {
	n.ti = ti
}

// #[visitor(Lhs,Rhs)]
type BinExpr struct {
	typ   NodeType
	rng   span.Range
	op    TokenValue
	opLoc span.Range
	lhs   Node
	rhs   Node
	opa   span.Range
}

func (n *BinExpr) Type() NodeType {
	return n.typ
}

func (n *BinExpr) Range() span.Range {
	return n.rng
}

func (n *BinExpr) Op() TokenValue {
	return n.op
}

func (n *BinExpr) OpText() string {
	return TokenKinds[n.op].Name
}

func (n *BinExpr) Lhs() Node {
	return n.lhs
}

func (n *BinExpr) Rhs() Node {
	return n.rhs
}

func (n *BinExpr) OuterParen() span.Range {
	return n.opa
}

func (n *BinExpr) SetOuterParen(rng span.Range) {
	n.opa = rng
}

// #[visitor(Arg)]
type UnaryExpr struct {
	typ NodeType
	rng span.Range
	op  TokenValue
	arg Node
	opa span.Range
}

func (n *UnaryExpr) Type() NodeType {
	return n.typ
}

func (n *UnaryExpr) Range() span.Range {
	return n.rng
}

func (n *UnaryExpr) Arg() Node {
	return n.arg
}

func (n *UnaryExpr) Op() TokenValue {
	return n.op
}

func (n *UnaryExpr) OpText() string {
	return TokenKinds[n.op].Name
}

func (n *UnaryExpr) OpName() string {
	return TokenKinds[n.op].Name
}

func (n *UnaryExpr) OuterParen() span.Range {
	return n.opa
}

func (n *UnaryExpr) SetOuterParen(rng span.Range) {
	n.opa = rng
}

// #[visitor(Arg)]
type UpdateExpr struct {
	typ    NodeType
	rng    span.Range
	op     TokenValue
	prefix bool
	arg    Node
	opa    span.Range
}

func (n *UpdateExpr) Type() NodeType {
	return n.typ
}

func (n *UpdateExpr) Range() span.Range {
	return n.rng
}

func (n *UpdateExpr) Arg() Node {
	return n.arg
}

func (n *UpdateExpr) Prefix() bool {
	return n.prefix
}

func (n *UpdateExpr) Op() TokenValue {
	return n.op
}

func (n *UpdateExpr) OpText() string {
	return TokenKinds[n.op].Name
}

func (n *UpdateExpr) OuterParen() span.Range {
	return n.opa
}

func (n *UpdateExpr) SetOuterParen(rng span.Range) {
	n.opa = rng
}

// #[visitor(Test,Cons,Alt)]
type CondExpr struct {
	typ  NodeType
	rng  span.Range
	test Node
	cons Node
	alt  Node
	opa  span.Range
}

func (n *CondExpr) Type() NodeType {
	return n.typ
}

func (n *CondExpr) Range() span.Range {
	return n.rng
}

func (n *CondExpr) Test() Node {
	return n.test
}

func (n *CondExpr) Cons() Node {
	return n.cons
}

func (n *CondExpr) Alt() Node {
	return n.alt
}

func (n *CondExpr) OuterParen() span.Range {
	return n.opa
}

func (n *CondExpr) SetOuterParen(rng span.Range) {
	n.opa = rng
}

// #[visitor(Lhs,Rhs)]
type AssignExpr struct {
	typ   NodeType
	rng   span.Range
	op    TokenValue
	opLoc span.Range
	lhs   Node
	rhs   Node
	opa   span.Range
	ti    *TypInfo
}

func (n *AssignExpr) Type() NodeType {
	return n.typ
}

func (n *AssignExpr) Range() span.Range {
	return n.rng
}

func (n *AssignExpr) Op() TokenValue {
	return n.op
}

func (n *AssignExpr) OpName() string {
	return TokenKinds[n.op].Name
}

func (n *AssignExpr) Lhs() Node {
	return n.lhs
}

func (n *AssignExpr) Rhs() Node {
	return n.rhs
}

func (n *AssignExpr) OuterParen() span.Range {
	return n.opa
}

func (n *AssignExpr) SetOuterParen(rng span.Range) {
	n.opa = rng
}

func (n *AssignExpr) TypInfo() *TypInfo {
	return n.ti
}

func (n *AssignExpr) SetTypInfo(ti *TypInfo) {
	n.ti = ti
}

type ThisExpr struct {
	typ NodeType
	rng span.Range
	opa span.Range
	ti  *TypInfo
}

func (n *ThisExpr) Type() NodeType {
	return n.typ
}

func (n *ThisExpr) Range() span.Range {
	return n.rng
}

func (n *ThisExpr) OuterParen() span.Range {
	return n.opa
}

func (n *ThisExpr) SetOuterParen(rng span.Range) {
	n.opa = rng
}

func (n *ThisExpr) TypInfo() *TypInfo {
	return n.ti
}

func (n *ThisExpr) SetTypInfo(ti *TypInfo) {
	n.ti = ti
}

// #[visitor(Elems)]
type SeqExpr struct {
	typ   NodeType
	rng   span.Range
	elems []Node
	opa   span.Range
}

func (n *SeqExpr) Type() NodeType {
	return n.typ
}

func (n *SeqExpr) Range() span.Range {
	return n.rng
}

func (n *SeqExpr) Elems() []Node {
	return n.elems
}

func (n *SeqExpr) OuterParen() span.Range {
	return n.opa
}

func (n *SeqExpr) SetOuterParen(rng span.Range) {
	n.opa = rng
}

// #[visitor(Expr)]
type ParenExpr struct {
	typ  NodeType
	rng  span.Range
	expr Node
	opa  span.Range
}

func (n *ParenExpr) Type() NodeType {
	return n.typ
}

func (n *ParenExpr) Range() span.Range {
	return n.rng
}

func (n *ParenExpr) Expr() Node {
	expr := n.expr
	for {
		if expr.Type() == N_EXPR_PAREN {
			expr = expr.(*ParenExpr).expr
		} else {
			break
		}
	}
	return expr
}

func (n *ParenExpr) OuterParen() span.Range {
	return n.opa
}

func (n *ParenExpr) SetOuterParen(rng span.Range) {
	n.opa = rng
}

// there is no information kept to describe the program order of the quasis and expressions
// according to below link descries how the quasis and expression are being walk over:
// https://opensource.apple.com/source/WebInspectorUI/WebInspectorUI-7602.2.14.0.5/UserInterface/Workers/Formatter/ESTreeWalker.js.auto.html
// some meaningless output should be taken into its estree result, such as put first quasis as
// a empty string if the first element in `elems` is a expression
//
// #[visitor(Tag,Elems)]
type TplExpr struct {
	typ   NodeType
	rng   span.Range
	tag   Node
	elems []Node
	opa   span.Range
}

func (n *TplExpr) Tag() Node {
	return n.tag
}

func (n *TplExpr) Range() span.Range {
	return n.rng
}

func (n *TplExpr) Elems() []Node {
	return n.elems
}

func (n *TplExpr) Type() NodeType {
	return n.typ
}

func (n *TplExpr) OuterParen() span.Range {
	return n.opa
}

func (n *TplExpr) SetOuterParen(rng span.Range) {
	n.opa = rng
}

type Super struct {
	typ NodeType
	rng span.Range
	opa span.Range
	ti  *TypInfo
}

func (n *Super) Type() NodeType {
	return n.typ
}

func (n *Super) Range() span.Range {
	return n.rng
}

func (n *Super) OuterParen() span.Range {
	return n.opa
}

func (n *Super) SetOuterParen(rng span.Range) {
	n.opa = rng
}

func (n *Super) TypInfo() *TypInfo {
	return n.ti
}

func (n *Super) SetTypInfo(ti *TypInfo) {
	n.ti = ti
}

// #[visitor(Src)]
type ImportCall struct {
	typ NodeType
	rng span.Range
	src Node
	opa span.Range
}

func (n *ImportCall) Type() NodeType {
	return n.typ
}

func (n *ImportCall) Range() span.Range {
	return n.rng
}

func (n *ImportCall) Src() Node {
	return n.src
}

func (n *ImportCall) OuterParen() span.Range {
	return n.opa
}

func (n *ImportCall) SetOuterParen(rng span.Range) {
	n.opa = rng
}

// #[visitor(Arg)]
type YieldExpr struct {
	typ      NodeType
	rng      span.Range
	delegate bool
	arg      Node
	opa      span.Range
}

func (n *YieldExpr) Type() NodeType {
	return n.typ
}

func (n *YieldExpr) Range() span.Range {
	return n.rng
}

func (n *YieldExpr) Delegate() bool {
	return n.delegate
}

func (n *YieldExpr) Arg() Node {
	return n.arg
}

func (n *YieldExpr) OuterParen() span.Range {
	return n.opa
}

func (n *YieldExpr) SetOuterParen(rng span.Range) {
	n.opa = rng
}

// #[visitor(Elems)]
type ArrPat struct {
	typ   NodeType
	rng   span.Range
	elems []Node
	opa   span.Range
	ti    *TypInfo
}

func (n *ArrPat) Type() NodeType {
	return n.typ
}

func (n *ArrPat) Range() span.Range {
	return n.rng
}

func (n *ArrPat) Elems() []Node {
	return n.elems
}

func (n *ArrPat) OuterParen() span.Range {
	return n.opa
}

func (n *ArrPat) SetOuterParen(rng span.Range) {
	n.opa = rng
}

func (n *ArrPat) TypInfo() *TypInfo {
	return n.ti
}

func (n *ArrPat) SetTypInfo(ti *TypInfo) {
	n.ti = ti
}

// #[visitor(Lhs,Rhs)]
type AssignPat struct {
	typ NodeType
	rng span.Range
	lhs Node
	rhs Node
	opa span.Range
	ti  *TypInfo
}

func (n *AssignPat) Type() NodeType {
	return n.typ
}

func (n *AssignPat) Range() span.Range {
	return n.rng
}

func (n *AssignPat) Lhs() Node {
	return n.lhs
}

func (n *AssignPat) Rhs() Node {
	return n.rhs
}

func (n *AssignPat) OuterParen() span.Range {
	return n.opa
}

func (n *AssignPat) TypInfo() *TypInfo {
	return n.ti
}

func (n *AssignPat) SetTypInfo(ti *TypInfo) {
	n.ti = ti
}

func (n *AssignPat) hoistTypInfo() {
	if wt, ok := n.lhs.(NodeWithTypInfo); ok {
		n.ti = wt.TypInfo()
	}
}

// #[visitor(Arg)]
type RestPat struct {
	typ NodeType
	rng span.Range
	arg Node
	opa span.Range
	ti  *TypInfo
}

func (n *RestPat) Type() NodeType {
	return n.typ
}

func (n *RestPat) Range() span.Range {
	return n.rng
}

func (n *RestPat) Arg() Node {
	return n.arg
}

func (n *RestPat) OuterParen() span.Range {
	return n.opa
}

func (n *RestPat) SetOuterParen(rng span.Range) {
	n.opa = rng
}

func (n *RestPat) TypInfo() *TypInfo {
	return n.ti
}

func (n *RestPat) Optional() bool {
	return !n.ti.Ques().Empty()
}

func (n *RestPat) hoistTypInfo() {
	if wt, ok := n.arg.(NodeWithTypInfo); ok {
		n.ti = wt.TypInfo()
		wt.SetTypInfo(nil)
	}
}

// #[visitor(Props)]
type ObjPat struct {
	typ   NodeType
	rng   span.Range
	props []Node
	opa   span.Range
	ti    *TypInfo
}

func (n *ObjPat) Type() NodeType {
	return n.typ
}

func (n *ObjPat) Range() span.Range {
	return n.rng
}

func (n *ObjPat) Props() []Node {
	return n.props
}

func (n *ObjPat) OuterParen() span.Range {
	return n.opa
}

func (n *ObjPat) SetOuterParen(rng span.Range) {
	n.opa = rng
}

func (n *ObjPat) TypInfo() *TypInfo {
	return n.ti
}

func (n *ObjPat) SetTypInfo(ti *TypInfo) {
	n.ti = ti
}

type PropKind uint8

const (
	PK_NONE PropKind = iota
	PK_INIT
	PK_GETTER
	PK_SETTER
	PK_CTOR
	PK_METHOD
)

func (pk PropKind) ToString() string {
	switch pk {
	case PK_INIT:
		return "init"
	case PK_GETTER:
		return "get"
	case PK_SETTER:
		return "set"
	case PK_CTOR:
		return "constructor"
	case PK_METHOD:
		return "method"
	default:
		return ""
	}
}

// #[visitor(Key,Val)]
type Prop struct {
	typ      NodeType
	rng      span.Range
	key      Node
	opLoc    span.Range
	value    Node
	computed bool
	method   bool

	shorthand bool
	assign    bool // it's `true` if the prop value is in assign pattern

	kind    PropKind
	accMode ACC_MOD
}

func (n *Prop) Type() NodeType {
	return n.typ
}

func (n *Prop) Range() span.Range {
	return n.rng
}

func (n *Prop) PropKind() PropKind {
	return n.kind
}

func (n *Prop) Kind() string {
	return n.kind.ToString()
}

func (n *Prop) Method() bool {
	return n.kind == PK_INIT && n.method
}

func (n *Prop) Key() Node {
	return n.key
}

func (n *Prop) Val() Node {
	return n.value
}

func (n *Prop) Computed() bool {
	return n.computed
}

func (n *Prop) Shorthand() bool {
	return n.shorthand
}

// #[visitor(Id,PUSH_SCOPE,Params,Body)]
type FnDec struct {
	typ       NodeType
	rng       span.Range
	id        Node
	generator bool
	async     bool
	params    []Node
	body      Node
	rets      []Node
	opa       span.Range
	ti        *TypInfo
}

func (n *FnDec) Type() NodeType {
	return n.typ
}

func (n *FnDec) Range() span.Range {
	return n.rng
}

func (n *FnDec) Id() Node {
	return n.id
}

func (n *FnDec) Generator() bool {
	return n.generator
}

func (n *FnDec) Async() bool {
	return n.async
}

func (n *FnDec) Params() []Node {
	return n.params
}

func (n *FnDec) Body() Node {
	return n.body
}

func (n *FnDec) Rets() []Node {
	return n.rets
}

func (n *FnDec) ExpRet() bool {
	return len(n.rets) > 0
}

func (n *FnDec) IsSig() bool {
	return n.body == nil
}

func (n *FnDec) OuterParen() span.Range {
	return n.opa
}

func (n *FnDec) SetOuterParen(rng span.Range) {
	n.opa = rng
}

func (n *FnDec) TypInfo() *TypInfo {
	return n.ti
}

func (n *FnDec) SetTypInfo(ti *TypInfo) {
	n.ti = ti
}

// #[visitor(PUSH_SCOPE,Params,Body)]
type ArrowFn struct {
	typ      NodeType
	rng      span.Range
	arrowLoc span.Range
	async    bool
	params   []Node
	body     Node
	expr     bool
	rets     []Node
	opa      span.Range
	ti       *TypInfo
}

func (n *ArrowFn) Type() NodeType {
	return n.typ
}

func (n *ArrowFn) Range() span.Range {
	return n.rng
}

func (n *ArrowFn) Async() bool {
	return n.async
}

func (n *ArrowFn) Params() []Node {
	return n.params
}

func (n *ArrowFn) Rets() []Node {
	return n.rets
}

func (n *ArrowFn) ExpRet() bool {
	return len(n.rets) > 0
}

func (n *ArrowFn) Body() Node {
	return n.body
}

func (n *ArrowFn) OuterParen() span.Range {
	return n.opa
}

func (n *ArrowFn) SetOuterParen(rng span.Range) {
	n.opa = rng
}

func (n *ArrowFn) Expr() bool {
	return n.expr
}

func (n *ArrowFn) TypInfo() *TypInfo {
	return n.ti
}

func (n *ArrowFn) SetTypInfo(ti *TypInfo) {
	n.ti = ti
}

// #[visitor(DecList)]
type VarDecStmt struct {
	typ     NodeType
	rng     span.Range
	kind    TokenValue
	decList []Node
	names   []Node
}

func (n *VarDecStmt) Type() NodeType {
	return n.typ
}

func (n *VarDecStmt) Range() span.Range {
	return n.rng
}

func (n *VarDecStmt) Names() []Node {
	return n.names
}

func (n *VarDecStmt) Kind() string {
	return TokenKinds[n.kind].Name
}

func (n *VarDecStmt) DecList() []Node {
	return n.decList
}

// #[visitor(Id,Init)]
type VarDec struct {
	typ  NodeType
	rng  span.Range
	id   Node
	init Node
}

func (n *VarDec) Type() NodeType {
	return n.typ
}

func (n *VarDec) Range() span.Range {
	return n.rng
}

func (n *VarDec) Id() Node {
	return n.id
}

func (n *VarDec) Init() Node {
	return n.init
}

// #[visitor(PUSH_SCOPE,Body)]
type BlockStmt struct {
	typ      NodeType
	rng      span.Range
	body     []Node
	newScope bool // whether or not to introduce a new scope
}

func (n *BlockStmt) Type() NodeType {
	return n.typ
}
func (n *BlockStmt) Range() span.Range {
	return n.rng
}

func (n *BlockStmt) Body() []Node {
	return n.body
}

func (n *BlockStmt) NewScope() bool {
	return n.newScope
}

// #[visitor(PUSH_SCOPE,Body,Test)]
type DoWhileStmt struct {
	typ  NodeType
	rng  span.Range
	test Node
	body Node
}

func (n *DoWhileStmt) Type() NodeType {
	return n.typ
}

func (n *DoWhileStmt) Range() span.Range {
	return n.rng
}

func (n *DoWhileStmt) Test() Node {
	return n.test
}

func (n *DoWhileStmt) Body() Node {
	return n.body
}

// #[visitor(PUSH_SCOPE,Test,Body)]
type WhileStmt struct {
	typ  NodeType
	rng  span.Range
	test Node
	body Node
}

func (n *WhileStmt) Type() NodeType {
	return n.typ
}

func (n *WhileStmt) Range() span.Range {
	return n.rng
}

func (n *WhileStmt) Test() Node {
	return n.test
}

func (n *WhileStmt) Body() Node {
	return n.body
}

// #[visitor(PUSH_SCOPE,Init,Test,Update,Body)]
type ForStmt struct {
	typ    NodeType
	rng    span.Range
	init   Node
	test   Node
	update Node
	body   Node
}

func (n *ForStmt) Type() NodeType {
	return n.typ
}

func (n *ForStmt) Range() span.Range {
	return n.rng
}

func (n *ForStmt) Init() Node {
	return n.init
}

func (n *ForStmt) Test() Node {
	return n.test
}

func (n *ForStmt) Update() Node {
	return n.update
}

func (n *ForStmt) Body() Node {
	return n.body
}

// #[visitor(Left,Right,Body)]
type ForInOfStmt struct {
	typ   NodeType
	rng   span.Range
	in    bool
	await bool
	left  Node
	right Node
	body  Node
}

func (n *ForInOfStmt) Type() NodeType {
	return n.typ
}

func (n *ForInOfStmt) Range() span.Range {
	return n.rng
}

func (n *ForInOfStmt) In() bool {
	return n.in
}

func (n *ForInOfStmt) Await() bool {
	return n.await
}

func (n *ForInOfStmt) Left() Node {
	return n.left
}

func (n *ForInOfStmt) Right() Node {
	return n.right
}

func (n *ForInOfStmt) Body() Node {
	return n.body
}

// #[visitor(Test,Cons,Alt)]
type IfStmt struct {
	typ  NodeType
	rng  span.Range
	test Node
	cons Node
	alt  Node
}

func (n *IfStmt) Type() NodeType {
	return n.typ
}

func (n *IfStmt) Range() span.Range {
	return n.rng
}

func (n *IfStmt) Test() Node {
	return n.test
}

func (n *IfStmt) Cons() Node {
	return n.cons
}

func (n *IfStmt) Alt() Node {
	return n.alt
}

// #[visitor(PUSH_SCOPE,Test,Cases)]
type SwitchStmt struct {
	typ   NodeType
	rng   span.Range
	test  Node
	cases []Node
}

func (n *SwitchStmt) Type() NodeType {
	return n.typ
}

func (n *SwitchStmt) Range() span.Range {
	return n.rng
}

func (n *SwitchStmt) Cases() []Node {
	return n.cases
}

func (n *SwitchStmt) Test() Node {
	return n.test
}

// #[visitor(Test,Cons)]
type SwitchCase struct {
	typ  NodeType
	rng  span.Range
	test Node // nil in default clause
	cons []Node
}

func (n *SwitchCase) Type() NodeType {
	return n.typ
}

func (n *SwitchCase) Range() span.Range {
	return n.rng
}

func (n *SwitchCase) Test() Node {
	return n.test
}

func (n *SwitchCase) Cons() []Node {
	return n.cons
}

// #[visitor(Label)]
type BrkStmt struct {
	typ    NodeType
	rng    span.Range
	label  Node
	target Node
}

func (n *BrkStmt) Type() NodeType {
	return n.typ
}

func (n *BrkStmt) Range() span.Range {
	return n.rng
}

func (n *BrkStmt) Label() Node {
	return n.label
}

func (n *BrkStmt) Target() Node {
	return n.target
}

// #[visitor(Label)]
type ContStmt struct {
	typ    NodeType
	rng    span.Range
	label  Node
	target Node
}

func (n *ContStmt) Type() NodeType {
	return n.typ
}

func (n *ContStmt) Range() span.Range {
	return n.rng
}

func (n *ContStmt) Label() Node {
	return n.label
}

func (n *ContStmt) Target() Node {
	return n.target
}

// #[visitor(Label,Body)]
type LabelStmt struct {
	typ   NodeType
	rng   span.Range
	label Node
	body  Node
	used  bool
}

func (n *LabelStmt) Type() NodeType {
	return n.typ
}

func (n *LabelStmt) Range() span.Range {
	return n.rng
}

func (n *LabelStmt) Label() Node {
	return n.label
}

func (n *LabelStmt) Body() Node {
	return n.body
}

func (n *LabelStmt) Used() bool {
	return n.used
}

// #[visitor(Arg)]
type RetStmt struct {
	typ NodeType
	rng span.Range
	arg Node
}

func (n *RetStmt) Type() NodeType {
	return n.typ
}

func (n *RetStmt) Range() span.Range {
	return n.rng
}

func (n *RetStmt) Arg() Node {
	return n.arg
}

// #[visitor(Arg)]
type ThrowStmt struct {
	typ    NodeType
	rng    span.Range
	arg    Node
	target Node
}

func (n *ThrowStmt) Type() NodeType {
	return n.typ
}

func (n *ThrowStmt) Range() span.Range {
	return n.rng
}

func (n *ThrowStmt) Arg() Node {
	return n.arg
}

func (n *ThrowStmt) Target() Node {
	return n.target
}

// #[visitor(PUSH_SCOPE,Param,Body)]
type Catch struct {
	typ   NodeType
	rng   span.Range
	param Node
	body  Node
}

func (n *Catch) Type() NodeType {
	return n.typ
}

func (n *Catch) Range() span.Range {
	return n.rng
}

func (n *Catch) Param() Node {
	return n.param
}

func (n *Catch) Body() Node {
	return n.body
}

// #[visitor(Try,Catch,Fin)]
type TryStmt struct {
	typ   NodeType
	rng   span.Range
	try   Node
	catch Node
	fin   Node
}

func (n *TryStmt) Type() NodeType {
	return n.typ
}

func (n *TryStmt) Range() span.Range {
	return n.rng
}

func (n *TryStmt) Try() Node {
	return n.try
}

func (n *TryStmt) Catch() Node {
	return n.catch
}

func (n *TryStmt) Fin() Node {
	return n.fin
}

type DebugStmt struct {
	typ NodeType
	rng span.Range
}

func (n *DebugStmt) Type() NodeType {
	return n.typ
}

func (n *DebugStmt) Range() span.Range {
	return n.rng
}

// #[visitor(Expr,Body)]
type WithStmt struct {
	typ  NodeType
	rng  span.Range
	expr Node
	body Node
}

func (n *WithStmt) Type() NodeType {
	return n.typ
}

func (n *WithStmt) Range() span.Range {
	return n.rng
}

func (n *WithStmt) Expr() Node {
	return n.expr
}

func (n *WithStmt) Body() Node {
	return n.body
}

// #[visitor(Id,Super,Body)]
type ClassDec struct {
	typ     NodeType
	rng     span.Range
	id      Node
	super   Node
	body    Node
	declare bool
	ti      *TypInfo
}

func (n *ClassDec) Type() NodeType {
	return n.typ
}

func (n *ClassDec) Range() span.Range {
	return n.rng
}

func (n *ClassDec) Declare() bool {
	return n.declare
}

func (n *ClassDec) Id() Node {
	return n.id
}

func (n *ClassDec) Super() Node {
	return n.super
}

func (n *ClassDec) Body() Node {
	return n.body
}

func (n *ClassDec) TypInfo() *TypInfo {
	return n.ti
}

func (n *ClassDec) SetTypInfo(ti *TypInfo) {
	n.ti = ti
}

// #[visitor(PUSH_SCOPE,Elems)]
type ClassBody struct {
	typ   NodeType
	rng   span.Range
	elems []Node
}

func (n *ClassBody) Type() NodeType {
	return n.typ
}

func (n *ClassBody) Range() span.Range {
	return n.rng
}

func (n *ClassBody) Elems() []Node {
	return n.elems
}

// #[visitor(Key,Val)]
type Method struct {
	typ      NodeType
	rng      span.Range
	key      Node
	static   bool
	computed bool
	kind     PropKind
	val      Node
	ti       *TypInfo
}

func (n *Method) Type() NodeType {
	return n.typ
}

func (n *Method) Range() span.Range {
	return n.rng
}

func (n *Method) PropKind() PropKind {
	return n.kind
}

func (n *Method) Kind() string {
	return n.kind.ToString()
}

func (n *Method) Key() Node {
	return n.key
}

func (n *Method) Val() Node {
	return n.val
}

func (n *Method) HasBody() bool {
	return n.val != nil && n.val.(*FnDec).body != nil
}

func (n *Method) Computed() bool {
	return n.computed
}

func (n *Method) Static() bool {
	return n.static
}

func (n *Method) Declare() bool {
	fn := n.val.(*FnDec)
	return fn.body == nil
}

func (n *Method) TypInfo() *TypInfo {
	return n.ti
}

func (n *Method) SetTypInfo(ti *TypInfo) {
	n.ti = ti
}

// #[visitor(Key,Val)]
type Field struct {
	typ      NodeType
	rng      span.Range
	key      Node
	static   bool
	computed bool
	val      Node
	ti       *TypInfo
}

func (n *Field) Type() NodeType {
	return n.typ
}

func (n *Field) Range() span.Range {
	return n.rng
}

func (n *Field) Key() Node {
	return n.key
}

func (n *Field) Val() Node {
	return n.val
}

func (n *Field) Static() bool {
	return n.static
}

func (n *Field) Computed() bool {
	return n.computed
}

func (n *Field) TypInfo() *TypInfo {
	return n.ti
}

func (n *Field) SetTypInfo(ti *TypInfo) {
	n.ti = ti
}

// #[visitor(Body)]
type StaticBlock struct {
	typ  NodeType
	rng  span.Range
	body []Node
	ti   *TypInfo
}

func (n *StaticBlock) Type() NodeType {
	return n.typ
}

func (n *StaticBlock) Range() span.Range {
	return n.rng
}

func (n *StaticBlock) Body() []Node {
	return n.body
}

func (n *StaticBlock) TypInfo() *TypInfo {
	return n.ti
}

func (n *StaticBlock) SetTypInfo(ti *TypInfo) {
	n.ti = ti
}

// #[visitor(Meta,Prop)]
type MetaProp struct {
	typ  NodeType
	rng  span.Range
	meta Node
	prop Node
}

func (n *MetaProp) Type() NodeType {
	return n.typ
}

func (n *MetaProp) Range() span.Range {
	return n.rng
}

func (n *MetaProp) Meta() Node {
	return n.meta
}

func (n *MetaProp) Prop() Node {
	return n.prop
}

// #[visitor(Specs,Src)]
type ImportDec struct {
	typ   NodeType
	rng   span.Range
	specs []Node
	src   Node
	tsTyp bool
}

func (n *ImportDec) Type() NodeType {
	return n.typ
}

func (n *ImportDec) Range() span.Range {
	return n.rng
}

func (n *ImportDec) TsTyp() bool {
	return n.tsTyp
}

func (n *ImportDec) Kind() string {
	if n.tsTyp {
		return "type"
	}
	return "value"
}

func (n *ImportDec) Specs() []Node {
	return n.specs
}

func (n *ImportDec) Src() Node {
	return n.src
}

// #[visitor(Local,Id)]
type ImportSpec struct {
	typ   NodeType
	rng   span.Range
	def   bool
	ns    bool
	local Node
	id    Node
	tsTyp bool // if represents the ts type
}

func (n *ImportSpec) Type() NodeType {
	return n.typ
}

func (n *ImportSpec) Range() span.Range {
	return n.rng
}

func (n *ImportSpec) TsTyp() bool {
	return n.tsTyp
}

func (n *ImportSpec) Kind() string {
	if n.tsTyp {
		return "type"
	}
	return "value"
}

func (n *ImportSpec) Default() bool {
	return n.def
}

func (n *ImportSpec) NameSpace() bool {
	return n.ns
}

func (n *ImportSpec) Local() Node {
	return n.local
}

func (n *ImportSpec) Id() Node {
	return n.id
}

// #[visitor(Dec,Specs,Src)]
type ExportDec struct {
	typ   NodeType
	rng   span.Range
	all   bool
	def   span.Range
	dec   Node
	specs []Node
	src   Node
	tsTyp bool
}

func (n *ExportDec) Type() NodeType {
	return n.typ
}

func (n *ExportDec) Range() span.Range {
	return n.rng
}

func (n *ExportDec) TsTyp() bool {
	return n.tsTyp
}

func (n *ExportDec) Kind() string {
	if n.tsTyp {
		return "type"
	}
	return "value"
}

func (n *ExportDec) All() bool {
	return n.all
}

func (n *ExportDec) Default() bool {
	return !n.def.Empty()
}

func (n *ExportDec) Dec() Node {
	return n.dec
}

func (n *ExportDec) Specs() []Node {
	return n.specs
}

func (n *ExportDec) Src() Node {
	return n.src
}

// #[visitor(Local,Id)]
type ExportSpec struct {
	typ   NodeType
	rng   span.Range
	ns    bool
	local Node
	id    Node // exported
	tsTyp bool // if represents the ts type
}

func (n *ExportSpec) Type() NodeType {
	return n.typ
}

func (n *ExportSpec) Range() span.Range {
	return n.rng
}

func (n *ExportSpec) Kind() string {
	if n.tsTyp {
		return "type"
	}
	return "value"
}

func (n *ExportSpec) TsTyp() bool {
	return n.tsTyp
}

func (n *ExportSpec) NameSpace() bool {
	return n.ns
}

func (n *ExportSpec) Local() Node {
	return n.local
}

func (n *ExportSpec) Id() Node {
	return n.id
}

// #[visitor(Expr)]
type ChainExpr struct {
	typ  NodeType
	rng  span.Range
	expr Node
}

func (n *ChainExpr) Type() NodeType {
	return n.typ
}

func (n *ChainExpr) Range() span.Range {
	return n.rng
}

func (n *ChainExpr) Expr() Node {
	return n.expr
}

// #[visitor(Expr)]
type Decorator struct {
	typ  NodeType
	rng  span.Range
	expr Node
}

func (n *Decorator) Type() NodeType {
	return n.typ
}

func (n *Decorator) Range() span.Range {
	return n.rng
}

func (n *Decorator) Expr() Node {
	return n.expr
}
