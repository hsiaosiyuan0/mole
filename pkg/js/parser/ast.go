package parser

import (
	"strconv"
	"strings"
)

// AST nodes is described as: https://github.com/estree/estree/blob/master/es5.md
// flatterned struct is used instead of inheritance

type Node interface {
	Type() NodeType
	Loc() *Loc
}

type Range struct {
	start int
	end   int
}

func (r *Range) Start() int {
	return r.start
}

func (r *Range) End() int {
	return r.end
}

func (r *Range) Clone() *Range {
	return &Range{
		start: r.start,
		end:   r.end,
	}
}

type Loc struct {
	src   *Source
	begin *Pos
	end   *Pos
	rng   *Range
}

func NewLoc() *Loc {
	return &Loc{
		src:   nil,
		begin: &Pos{},
		end:   &Pos{},
		rng:   &Range{},
	}
}

func (l *Loc) Source() string {
	return l.src.path
}

func (l *Loc) Begin() *Pos {
	return l.begin
}

func (l *Loc) End() *Pos {
	return l.end
}

func (l *Loc) Range() *Range {
	return l.rng
}

func (l *Loc) Clone() *Loc {
	return &Loc{
		src:   l.src,
		begin: l.begin.Clone(),
		end:   l.end.Clone(),
		rng:   l.rng.Clone(),
	}
}

type NodeType int

const (
	N_ILLEGAL NodeType = iota
	N_PROG

	N_STMT_BEGIN
	N_STMT_EXPR
	N_STMT_EMPTY
	N_STMT_VAR_DEC
	N_STMT_FN
	N_STMT_BLOCK
	N_STMT_DO_WHILE
	N_STMT_WHILE
	N_STMT_FOR
	N_STMT_FOR_IN_OF
	N_STMT_IF
	N_STMT_SWITCH
	N_STMT_BRK
	N_STMT_CONT
	N_STMT_LABEL
	N_STMT_RET
	N_STMT_THROW
	N_STMT_TRY
	N_STMT_DEBUG
	N_STMT_WITH
	N_STMT_CLASS
	N_STMT_END

	N_EXPR_BEGIN
	N_LIT_BEGIN
	N_LIT_NULL
	N_LIT_BOOL
	N_LIT_NUM
	N_LIT_STR
	N_LIT_ARR
	N_LIT_OBJ
	N_LIT_REGEXP
	N_LIT_END

	N_EXPR_NEW
	N_EXPR_MEMBER
	N_EXPR_CALL
	N_EXPR_BIN
	N_EXPR_UNARY
	N_EXPR_UPDATE
	N_EXPR_COND
	N_EXPR_ASSIGN
	N_EXPR_FN
	N_EXPR_THIS
	N_EXPR_PAREN
	N_EXPR_ARROW
	N_EXPR_SEQ
	N_EXPR_CLASS
	N_EXPR_TPL
	N_NAME

	N_VAR_DEC
	N_PATTERN_REST
	N_PATTERN_ARRAY
	N_PATTERN_ASSIGN
	N_PATTERN_OBJ
	N_PROP
	N_SPREAD
	N_SWITCH_CASE
	N_CATCH
	N_ClASS_BODY
	N_STATIC_BLOCK
	N_METHOD
	N_FIELD
	N_SUPER
	N_IMPORT_CALL
	N_META_PROP

	N_EXPR_END
)

type Prog struct {
	typ   NodeType
	loc   *Loc
	stmts []Node
}

func (n *Prog) Body() []Node {
	return n.stmts
}

func NewProg() *Prog {
	return &Prog{N_PROG, &Loc{}, make([]Node, 0)}
}

func (n *Prog) Type() NodeType {
	return n.typ
}

func (n *Prog) Loc() *Loc {
	return n.loc
}

type ExprStmt struct {
	typ  NodeType
	loc  *Loc
	expr Node
}

func (n *ExprStmt) Type() NodeType {
	return n.typ
}

func (n *ExprStmt) Loc() *Loc {
	return n.loc
}

func (n *ExprStmt) Expr() Node {
	return n.expr
}

type EmptyStmt struct {
	typ NodeType
	loc *Loc
}

func (n *EmptyStmt) Type() NodeType {
	return n.typ
}

func (n *EmptyStmt) Loc() *Loc {
	return n.loc
}

type NullLit struct {
	typ NodeType
	loc *Loc
}

func (n *NullLit) Type() NodeType {
	return n.typ
}

func (n *NullLit) Loc() *Loc {
	return n.loc
}

type BoolLit struct {
	typ NodeType
	loc *Loc
	val *Token
}

func (n *BoolLit) Value() bool {
	return n.val.Text() == "true"
}

func (n *BoolLit) Type() NodeType {
	return n.typ
}

func (n *BoolLit) Loc() *Loc {
	return n.loc
}

type NumLit struct {
	typ NodeType
	loc *Loc
	val *Token
}

func (n *NumLit) Type() NodeType {
	return n.typ
}

func (n *NumLit) Loc() *Loc {
	return n.loc
}

func (n *NumLit) ToFloat() float64 {
	t := n.val.Text()
	if strings.HasPrefix(t, "0x") || strings.HasPrefix(t, "0X") {
		s, _ := strconv.ParseUint(t[2:], 16, 32)
		return float64(s)
	}
	if strings.HasPrefix(t, "0x") || strings.HasPrefix(t, "0X") {
		s, _ := strconv.ParseUint(t[2:], 8, 32)
		return float64(s)
	}
	if strings.HasPrefix(t, "0") && len(t) > 1 {
		t = strings.TrimLeft(t, "0")
		s, _ := strconv.ParseUint(t, 8, 32)
		return float64(s)
	}
	s, _ := strconv.ParseFloat(n.val.Text(), 64)
	return s
}

type StrLit struct {
	typ NodeType
	loc *Loc
	val *Token
}

func (n *StrLit) Type() NodeType {
	return n.typ
}

func (n *StrLit) Loc() *Loc {
	return n.loc
}

func (n *StrLit) Text() string {
	return n.val.Text()
}

type RegexpLit struct {
	typ     NodeType
	loc     *Loc
	val     *Token
	pattern string
	flags   string
}

func (n *RegexpLit) Type() NodeType {
	return n.typ
}

func (n *RegexpLit) Loc() *Loc {
	return n.loc
}

func (n *RegexpLit) Pattern() string {
	return n.pattern
}

func (n *RegexpLit) Flags() string {
	return n.flags
}

type ArrLit struct {
	typ   NodeType
	loc   *Loc
	elems []Node
}

func (n *ArrLit) Type() NodeType {
	return n.typ
}

func (n *ArrLit) Loc() *Loc {
	return n.loc
}

func (n *ArrLit) Elems() []Node {
	return n.elems
}

type Spread struct {
	typ NodeType
	loc *Loc
	arg Node
}

func (n *Spread) Arg() Node {
	return n.arg
}

func (n *Spread) Type() NodeType {
	return n.typ
}

func (n *Spread) Loc() *Loc {
	return n.loc
}

type ObjLit struct {
	typ   NodeType
	loc   *Loc
	props []Node
}

func (n *ObjLit) Props() []Node {
	return n.props
}

func (n *ObjLit) Type() NodeType {
	return n.typ
}

func (n *ObjLit) Loc() *Loc {
	return n.loc
}

type Ident struct {
	typ NodeType
	loc *Loc
	val *Token
	pvt bool
}

func (n *Ident) Type() NodeType {
	return n.typ
}

func (n *Ident) Loc() *Loc {
	return n.loc
}

func (n *Ident) Text() string {
	return n.val.Text()
}

type NewExpr struct {
	typ    NodeType
	loc    *Loc
	callee Node
	args   []Node
}

func (n *NewExpr) Type() NodeType {
	return n.typ
}

func (n *NewExpr) Loc() *Loc {
	return n.loc
}

func (n *NewExpr) Callee() Node {
	return n.callee
}

func (n *NewExpr) Args() []Node {
	return n.args
}

type MemberExpr struct {
	typ      NodeType
	loc      *Loc
	obj      Node
	prop     Node
	compute  bool
	optional bool
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

func (n *MemberExpr) Type() NodeType {
	return n.typ
}

func (n *MemberExpr) Loc() *Loc {
	return n.loc
}

type CallExpr struct {
	typ    NodeType
	loc    *Loc
	callee Node
	args   []Node
}

func (n *CallExpr) Callee() Node {
	return n.callee
}

func (n *CallExpr) Args() []Node {
	return n.args
}

func (n *CallExpr) Type() NodeType {
	return n.typ
}

func (n *CallExpr) Loc() *Loc {
	return n.loc
}

type BinExpr struct {
	typ NodeType
	loc *Loc
	op  *Token
	lhs Node
	rhs Node
}

func (n *BinExpr) Op() *Token {
	return n.op
}

func (n *BinExpr) Lhs() Node {
	return n.lhs
}

func (n *BinExpr) Rhs() Node {
	return n.rhs
}

func (n *BinExpr) Type() NodeType {
	return n.typ
}

func (n *BinExpr) Loc() *Loc {
	return n.loc
}

type UnaryExpr struct {
	typ NodeType
	loc *Loc
	op  *Token
	arg Node
}

func (n *UnaryExpr) Arg() Node {
	return n.arg
}

func (n *UnaryExpr) Op() *Token {
	return n.op
}

func (n *UnaryExpr) Type() NodeType {
	return n.typ
}

func (n *UnaryExpr) Loc() *Loc {
	return n.loc
}

type UpdateExpr struct {
	typ    NodeType
	loc    *Loc
	op     *Token
	prefix bool
	arg    Node
}

func (n *UpdateExpr) Arg() Node {
	return n.arg
}

func (n *UpdateExpr) Prefix() bool {
	return n.prefix
}

func (n *UpdateExpr) Op() *Token {
	return n.op
}

func (n *UpdateExpr) Type() NodeType {
	return n.typ
}

func (n *UpdateExpr) Loc() *Loc {
	return n.loc
}

type CondExpr struct {
	typ  NodeType
	loc  *Loc
	test Node
	cons Node
	alt  Node
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

func (n *CondExpr) Type() NodeType {
	return n.typ
}

func (n *CondExpr) Loc() *Loc {
	return n.loc
}

type AssignExpr struct {
	typ NodeType
	loc *Loc
	op  *Token
	lhs Node
	rhs Node
}

func (n *AssignExpr) Op() *Token {
	return n.op
}

func (n *AssignExpr) Lhs() Node {
	return n.lhs
}

func (n *AssignExpr) Rhs() Node {
	return n.rhs
}

func (n *AssignExpr) Type() NodeType {
	return n.typ
}

func (n *AssignExpr) Loc() *Loc {
	return n.loc
}

type VarDecStmt struct {
	typ     NodeType
	loc     *Loc
	kind    TokenValue
	decList []*VarDec
}

func (n *VarDecStmt) Kind() string {
	return TokenKinds[n.kind].Name
}

func (n *VarDecStmt) DecList() []*VarDec {
	return n.decList
}

func (n *VarDecStmt) Type() NodeType {
	return n.typ
}

func (n *VarDecStmt) Loc() *Loc {
	return n.loc
}

type VarDec struct {
	typ  NodeType
	loc  *Loc
	id   Node
	init Node
}

func (n *VarDec) Id() Node {
	return n.id
}

func (n *VarDec) Init() Node {
	return n.init
}

func (n *VarDec) Type() NodeType {
	return n.typ
}

func (n *VarDec) Loc() *Loc {
	return n.loc
}

type ArrayPattern struct {
	typ   NodeType
	loc   *Loc
	elems []Node
}

func (n *ArrayPattern) Type() NodeType {
	return n.typ
}

func (n *ArrayPattern) Loc() *Loc {
	return n.loc
}

type AssignPattern struct {
	typ   NodeType
	loc   *Loc
	left  Node
	right Node
}

func (n *AssignPattern) Type() NodeType {
	return n.typ
}

func (n *AssignPattern) Loc() *Loc {
	return n.loc
}

type RestPattern struct {
	typ NodeType
	loc *Loc
	arg Node
}

func (n *RestPattern) Arg() Node {
	return n.arg
}

func (n *RestPattern) Type() NodeType {
	return n.typ
}

func (n *RestPattern) Loc() *Loc {
	return n.loc
}

type ObjPattern struct {
	typ   NodeType
	loc   *Loc
	props []Node
}

func (n *ObjPattern) Type() NodeType {
	return n.typ
}

func (n *ObjPattern) Loc() *Loc {
	return n.loc
}

type Prop struct {
	typ      NodeType
	loc      *Loc
	key      Node
	value    Node
	computed bool
	kind     *Token
}

func (n *Prop) Kind() string {
	if n.kind != nil {
		return n.kind.Text()
	}
	return "init"
}

func (n *Prop) Method() bool {
	return n.value.Type() == N_EXPR_FN
}

func (n *Prop) Key() Node {
	return n.key
}

func (n *Prop) Value() Node {
	return n.value
}

func (n *Prop) Computed() bool {
	return n.computed
}

func (n *Prop) Type() NodeType {
	return n.typ
}

func (n *Prop) Loc() *Loc {
	return n.loc
}

type FnDec struct {
	typ       NodeType
	loc       *Loc
	id        Node
	generator bool
	async     bool
	params    []Node
	body      Node
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

func (n *FnDec) Type() NodeType {
	return n.typ
}

func (n *FnDec) Loc() *Loc {
	return n.loc
}

type ArrowFn struct {
	typ       NodeType
	loc       *Loc
	generator bool
	async     bool
	params    []Node
	body      Node
}

func (n *ArrowFn) Type() NodeType {
	return n.typ
}

func (n *ArrowFn) Loc() *Loc {
	return n.loc
}

type BlockStmt struct {
	typ  NodeType
	loc  *Loc
	body []Node
}

func (n *BlockStmt) Body() []Node {
	return n.body
}

func (n *BlockStmt) Type() NodeType {
	return n.typ
}

func (n *BlockStmt) Loc() *Loc {
	return n.loc
}

type DoWhileStmt struct {
	typ  NodeType
	loc  *Loc
	test Node
	body Node
}

func (n *DoWhileStmt) Test() Node {
	return n.test
}

func (n *DoWhileStmt) Body() Node {
	return n.body
}

func (n *DoWhileStmt) Type() NodeType {
	return n.typ
}

func (n *DoWhileStmt) Loc() *Loc {
	return n.loc
}

type WhileStmt struct {
	typ  NodeType
	loc  *Loc
	test Node
	body Node
}

func (n *WhileStmt) Test() Node {
	return n.test
}

func (n *WhileStmt) Body() Node {
	return n.body
}

func (n *WhileStmt) Type() NodeType {
	return n.typ
}

func (n *WhileStmt) Loc() *Loc {
	return n.loc
}

type ForStmt struct {
	typ    NodeType
	loc    *Loc
	init   Node
	test   Node
	update Node
	body   Node
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

func (n *ForStmt) Type() NodeType {
	return n.typ
}

func (n *ForStmt) Loc() *Loc {
	return n.loc
}

type ForInOfStmt struct {
	typ   NodeType
	loc   *Loc
	in    bool
	await bool
	left  Node
	right Node
	body  Node
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

func (n *ForInOfStmt) Type() NodeType {
	return n.typ
}

func (n *ForInOfStmt) Loc() *Loc {
	return n.loc
}

type IfStmt struct {
	typ  NodeType
	loc  *Loc
	test Node
	cons Node
	alt  Node
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

func (n *IfStmt) Type() NodeType {
	return n.typ
}

func (n *IfStmt) Loc() *Loc {
	return n.loc
}

type SwitchStmt struct {
	typ   NodeType
	loc   *Loc
	test  Node
	cases []*SwitchCase
}

func (n *SwitchStmt) Cases() []*SwitchCase {
	return n.cases
}

func (n *SwitchStmt) Test() Node {
	return n.test
}

func (n *SwitchStmt) Type() NodeType {
	return n.typ
}

func (n *SwitchStmt) Loc() *Loc {
	return n.loc
}

type SwitchCase struct {
	typ  NodeType
	loc  *Loc
	test Node // nil in default clause
	cons []Node
}

func (n *SwitchCase) Test() Node {
	return n.test
}

func (n *SwitchCase) Cons() []Node {
	return n.cons
}

func (n *SwitchCase) Type() NodeType {
	return n.typ
}

func (n *SwitchCase) Loc() *Loc {
	return n.loc
}

type BrkStmt struct {
	typ   NodeType
	loc   *Loc
	label Node
}

func (n *BrkStmt) Label() Node {
	return n.label
}

func (n *BrkStmt) Type() NodeType {
	return n.typ
}

func (n *BrkStmt) Loc() *Loc {
	return n.loc
}

type ContStmt struct {
	typ   NodeType
	loc   *Loc
	label Node
}

func (n *ContStmt) Label() Node {
	return n.label
}

func (n *ContStmt) Type() NodeType {
	return n.typ
}

func (n *ContStmt) Loc() *Loc {
	return n.loc
}

type LabelStmt struct {
	typ   NodeType
	loc   *Loc
	label Node
	body  Node
}

func (n *LabelStmt) Label() Node {
	return n.label
}

func (n *LabelStmt) Body() Node {
	return n.body
}

func (n *LabelStmt) Type() NodeType {
	return n.typ
}

func (n *LabelStmt) Loc() *Loc {
	return n.loc
}

type RetStmt struct {
	typ NodeType
	loc *Loc
	arg Node
}

func (n *RetStmt) Arg() Node {
	return n.arg
}

func (n *RetStmt) Type() NodeType {
	return n.typ
}

func (n *RetStmt) Loc() *Loc {
	return n.loc
}

type ThrowStmt struct {
	typ NodeType
	loc *Loc
	arg Node
}

func (n *ThrowStmt) Arg() Node {
	return n.arg
}

func (n *ThrowStmt) Type() NodeType {
	return n.typ
}

func (n *ThrowStmt) Loc() *Loc {
	return n.loc
}

type Catch struct {
	typ   NodeType
	loc   *Loc
	param Node
	body  Node
}

func (n *Catch) Param() Node {
	return n.param
}

func (n *Catch) Body() Node {
	return n.body
}

func (n *Catch) Type() NodeType {
	return n.typ
}

func (n *Catch) Loc() *Loc {
	return n.loc
}

type TryStmt struct {
	typ   NodeType
	loc   *Loc
	try   Node
	catch Node
	fin   Node
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

func (n *TryStmt) Type() NodeType {
	return n.typ
}

func (n *TryStmt) Loc() *Loc {
	return n.loc
}

type DebugStmt struct {
	typ NodeType
	loc *Loc
}

func (n *DebugStmt) Type() NodeType {
	return n.typ
}

func (n *DebugStmt) Loc() *Loc {
	return n.loc
}

type WithStmt struct {
	typ  NodeType
	loc  *Loc
	expr Node
	body Node
}

func (n *WithStmt) Expr() Node {
	return n.expr
}

func (n *WithStmt) Body() Node {
	return n.body
}

func (n *WithStmt) Type() NodeType {
	return n.typ
}

func (n *WithStmt) Loc() *Loc {
	return n.loc
}

type ClassDec struct {
	typ   NodeType
	loc   *Loc
	id    Node
	super Node
	body  Node
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

func (n *ClassDec) Type() NodeType {
	return n.typ
}

func (n *ClassDec) Loc() *Loc {
	return n.loc
}

type ClassBody struct {
	typ   NodeType
	loc   *Loc
	elems []Node
}

func (n *ClassBody) Elems() []Node {
	return n.elems
}

func (n *ClassBody) Type() NodeType {
	return n.typ
}

func (n *ClassBody) Loc() *Loc {
	return n.loc
}

type Method struct {
	typ      NodeType
	loc      *Loc
	key      Node
	static   bool
	computed bool
	kind     *Token
	value    Node
}

func (n *Method) Kind() string {
	if n.kind != nil {
		return n.kind.Text()
	}
	return "init"
}

func (n *Method) Key() Node {
	return n.key
}

func (n *Method) Value() Node {
	return n.value
}

func (n *Method) Computed() bool {
	return n.computed
}

func (n *Method) Static() bool {
	return n.static
}

func (n *Method) Type() NodeType {
	return n.typ
}

func (n *Method) Loc() *Loc {
	return n.loc
}

type Field struct {
	typ      NodeType
	loc      *Loc
	key      Node
	static   bool
	computed bool
	value    Node
}

func (n *Field) Key() Node {
	return n.key
}

func (n *Field) Value() Node {
	return n.key
}

func (n *Field) Static() bool {
	return n.static
}

func (n *Field) Compute() bool {
	return n.computed
}

func (n *Field) Type() NodeType {
	return n.typ
}

func (n *Field) Loc() *Loc {
	return n.loc
}

type StaticBlock struct {
	typ  NodeType
	loc  *Loc
	body []Node
}

func (n *StaticBlock) Type() NodeType {
	return n.typ
}

func (n *StaticBlock) Loc() *Loc {
	return n.loc
}

type ThisExpr struct {
	typ NodeType
	loc *Loc
}

func (n *ThisExpr) Type() NodeType {
	return n.typ
}

func (n *ThisExpr) Loc() *Loc {
	return n.loc
}

type SeqExpr struct {
	typ   NodeType
	loc   *Loc
	elems []Node
}

func (n *SeqExpr) Elems() []Node {
	return n.elems
}

func (n *SeqExpr) Type() NodeType {
	return n.typ
}

func (n *SeqExpr) Loc() *Loc {
	return n.loc
}

type ParenExpr struct {
	typ  NodeType
	loc  *Loc
	expr Node
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

func (n *ParenExpr) Type() NodeType {
	return n.typ
}

func (n *ParenExpr) Loc() *Loc {
	return n.loc
}

type TplExpr struct {
	typ   NodeType
	loc   *Loc
	tag   Node
	elems []Node
}

func (n *TplExpr) Type() NodeType {
	return n.typ
}

func (n *TplExpr) Loc() *Loc {
	return n.loc
}

type Super struct {
	typ NodeType
	loc *Loc
}

func (n *Super) Type() NodeType {
	return n.typ
}

func (n *Super) Loc() *Loc {
	return n.loc
}

type ImportCall struct {
	typ NodeType
	loc *Loc
	src Node
}

func (n *ImportCall) Src() Node {
	return n.src
}

func (n *ImportCall) Type() NodeType {
	return n.typ
}

func (n *ImportCall) Loc() *Loc {
	return n.loc
}

type MetaProp struct {
	typ  NodeType
	loc  *Loc
	meta Node
	prop Node
}

func (n *MetaProp) Meta() Node {
	return n.meta
}

func (n *MetaProp) Prop() Node {
	return n.prop
}

func (n *MetaProp) Type() NodeType {
	return n.typ
}

func (n *MetaProp) Loc() *Loc {
	return n.loc
}
