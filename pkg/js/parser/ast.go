package parser

import (
	"math"
	"math/big"
	"strconv"
	"strings"
)

// below AST nodes are described as: https://github.com/estree/estree/blob/master/es5.md
// however the flatterned struct is used instead of inheritance

type Node interface {
	Type() NodeType
	Loc() *Loc
	Extra() interface{}
	setExtra(interface{})
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

func (l *Loc) Text() string {
	return l.src.code[l.rng.start:l.rng.end]
}

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

func (n *Prog) Extra() interface{} {
	return nil
}

func (n *Prog) setExtra(_ interface{}) {
}

type ExprStmt struct {
	typ  NodeType
	loc  *Loc
	expr Node
	dir  bool
}

func (n *ExprStmt) Type() NodeType {
	return n.typ
}

func (n *ExprStmt) Loc() *Loc {
	return n.loc
}

func (n *ExprStmt) Dir() bool {
	return n.dir
}

func (n *ExprStmt) DirStr() string {
	s := n.expr.(*StrLit)
	raw := s.Raw()
	return s.Raw()[1 : len(raw)-1]
}

func (n *ExprStmt) Expr() Node {
	return n.expr
}

func (n *ExprStmt) Extra() interface{} {
	return nil
}

func (n *ExprStmt) setExtra(_ interface{}) {
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

func (n *EmptyStmt) Extra() interface{} {
	return nil
}

func (n *EmptyStmt) setExtra(_ interface{}) {
}

type ExprExtra struct {
	OuterParen *Loc
}

type NullLit struct {
	typ   NodeType
	loc   *Loc
	extra *ExprExtra
}

func (n *NullLit) Type() NodeType {
	return n.typ
}

func (n *NullLit) Loc() *Loc {
	return n.loc
}

func (n *NullLit) Text() string {
	return "null"
}

func (n *NullLit) Extra() interface{} {
	return nil
}

func (n *NullLit) setExtra(ext interface{}) {
	n.extra = ext.(*ExprExtra)
}

type BoolLit struct {
	typ   NodeType
	loc   *Loc
	val   bool
	extra *ExprExtra
}

func (n *BoolLit) Value() bool {
	return n.val
}

func (n *BoolLit) Type() NodeType {
	return n.typ
}

func (n *BoolLit) Loc() *Loc {
	return n.loc
}

func (n *BoolLit) Text() string {
	if n.val {
		return "true"
	}
	return "false"
}

func (n *BoolLit) Extra() interface{} {
	return nil
}

func (n *BoolLit) setExtra(ext interface{}) {
	n.extra = ext.(*ExprExtra)
}

type NumLit struct {
	typ   NodeType
	loc   *Loc
	extra *ExprExtra
}

func (n *NumLit) Type() NodeType {
	return n.typ
}

func (n *NumLit) Loc() *Loc {
	return n.loc
}

func (n *NumLit) Text() string {
	return n.loc.Text()
}

func (n *NumLit) IsBigint() bool {
	t := n.loc.Text()
	return t[len(t)-1] == 'n'
}

// unsafe method, use `IsBigint` before this method
func (n *NumLit) ToBigint() *big.Int {
	t := n.loc.Text()
	t = strings.ReplaceAll(t[:len(t)-1], "_", "")

	if strings.HasPrefix(t, "0x") || strings.HasPrefix(t, "0X") {
		i := big.NewInt(0)
		i.SetString(t[2:], 16)
		return i
	}
	if strings.HasPrefix(t, "0o") || strings.HasPrefix(t, "0O") {
		i := big.NewInt(0)
		i.SetString(t[2:], 8)
		return i
	}
	if strings.HasPrefix(t, "0b") || strings.HasPrefix(t, "0B") {
		i := big.NewInt(0)
		i.SetString(t[2:], 2)
		return i
	}
	if strings.HasPrefix(t, "0") && len(t) > 1 {
		t = strings.TrimLeft(t, "0")
		i := big.NewInt(0)
		i.SetString(t, 8)
		return i
	}
	i := big.NewInt(0)
	i.SetString(t, 10)
	return i
}

// unsafe method, use `IsBigint` before this method
func (n *NumLit) ToFloat() float64 {
	t := n.loc.Text()
	t = strings.ReplaceAll(t, "_", "")

	if strings.HasPrefix(t, "0x") || strings.HasPrefix(t, "0X") {
		f, _ := strconv.ParseUint(t[2:], 16, 32)
		return float64(f)
	}
	if strings.HasPrefix(t, "0o") || strings.HasPrefix(t, "0O") {
		f, _ := strconv.ParseUint(t[2:], 8, 32)
		return float64(f)
	}
	if strings.HasPrefix(t, "0b") || strings.HasPrefix(t, "0B") {
		f, _ := strconv.ParseUint(t[2:], 2, 32)
		return float64(f)
	}
	if strings.HasPrefix(t, "0") && len(t) > 1 {
		t = strings.TrimLeft(t, "0")
		f, _ := strconv.ParseUint(t, 8, 32)
		return float64(f)
	}
	f, _ := strconv.ParseFloat(t, 64)
	return f
}

const maxSafeInt = (uint64(1) << 53) - 1

// unsafe method, use `IsBigint` before this method
// return the numerical value if it's safe to be represented in JSON as number
// otherwise zero is returned
func (n *NumLit) Float() float64 {
	if n.IsBigint() {
		i := n.ToBigint()
		f := big.NewFloat(0)
		f.SetInt(i)
		max := big.NewFloat(0)
		max.SetUint64(maxSafeInt)
		c := f.Cmp(max)
		if c == -1 || c == 0 {
			ff, _ := f.Float64()
			return ff
		}
		return 0
	}
	f := n.ToFloat()
	if f > math.MaxFloat64 {
		return 0
	}
	return f
}

func (n *NumLit) Extra() interface{} {
	return nil
}

func (n *NumLit) setExtra(ext interface{}) {
	n.extra = ext.(*ExprExtra)
}

type StrLit struct {
	typ                  NodeType
	loc                  *Loc
	val                  string
	legacyOctalEscapeSeq bool
	extra                *ExprExtra
}

func (n *StrLit) Type() NodeType {
	return n.typ
}

func (n *StrLit) Loc() *Loc {
	return n.loc
}

func (n *StrLit) Text() string {
	return n.val
}

func (n *StrLit) Raw() string {
	return n.loc.Text()
}

func (n *StrLit) Extra() interface{} {
	return nil
}

func (n *StrLit) setExtra(ext interface{}) {
	n.extra = ext.(*ExprExtra)
}

type RegLit struct {
	typ     NodeType
	loc     *Loc
	val     string
	pattern string
	flags   string
	extra   *ExprExtra
}

func (n *RegLit) Type() NodeType {
	return n.typ
}

func (n *RegLit) Loc() *Loc {
	return n.loc
}

func (n *RegLit) Pattern() string {
	return n.pattern
}

func (n *RegLit) Flags() string {
	return n.flags
}

func (n *RegLit) Extra() interface{} {
	return nil
}

func (n *RegLit) setExtra(ext interface{}) {
	n.extra = ext.(*ExprExtra)
}

type ArrLit struct {
	typ      NodeType
	loc      *Loc
	elems    []Node
	extra    *ExprExtra
	typAnnot Node
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

func (n *ArrLit) Extra() interface{} {
	return nil
}

func (n *ArrLit) setExtra(ext interface{}) {
	n.extra = ext.(*ExprExtra)
}

type Spread struct {
	typ              NodeType
	loc              *Loc
	arg              Node
	trailingCommaLoc *Loc
	extra            *ExprExtra
	typAnnot         Node
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

func (n *Spread) Extra() interface{} {
	return n.extra
}

func (n *Spread) setExtra(ext interface{}) {
	n.extra = ext.(*ExprExtra)
}

type ObjLit struct {
	typ      NodeType
	loc      *Loc
	props    []Node
	extra    *ExprExtra
	typAnnot Node
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

func (n *ObjLit) Extra() interface{} {
	return n.extra
}

func (n *ObjLit) setExtra(ext interface{}) {
	n.extra = ext.(*ExprExtra)
}

type Ident struct {
	typ            NodeType
	loc            *Loc
	val            string
	pvt            bool
	containsEscape bool
	extra          *ExprExtra
	// consider below statements:
	// `export { if } from "a"` is legal
	// `export { if } ` is illegal
	// for reporting `if` is a keyword, firstly produce a
	// Ident with conent `if` and flag it's a keyword by
	// setting this field to true, later report the `unexpected token`
	// error if the coming token is not `from`
	kw       bool
	typAnnot Node
}

func (n *Ident) Type() NodeType {
	return n.typ
}

func (n *Ident) Loc() *Loc {
	return n.loc
}

func (n *Ident) Text() string {
	return n.val
}

func (n *Ident) IsPrivate() bool {
	return n.pvt
}

func (n *Ident) ContainsEscape() bool {
	return n.containsEscape
}

func (n *Ident) Extra() interface{} {
	return n.extra
}

func (n *Ident) setExtra(ext interface{}) {
	n.extra = ext.(*ExprExtra)
}

type NewExpr struct {
	typ    NodeType
	loc    *Loc
	callee Node
	args   []Node
	extra  *ExprExtra
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

func (n *NewExpr) Extra() interface{} {
	return n.extra
}

func (n *NewExpr) setExtra(ext interface{}) {
	n.extra = ext.(*ExprExtra)
}

type MemberExpr struct {
	typ      NodeType
	loc      *Loc
	obj      Node
	prop     Node
	compute  bool
	optional bool
	extra    *ExprExtra
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

func (n *MemberExpr) Extra() interface{} {
	return n.extra
}

func (n *MemberExpr) setExtra(ext interface{}) {
	n.extra = ext.(*ExprExtra)
}

type CallExpr struct {
	typ      NodeType
	loc      *Loc
	callee   Node
	args     []Node
	optional bool
	extra    *ExprExtra
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

func (n *CallExpr) Type() NodeType {
	return n.typ
}

func (n *CallExpr) Loc() *Loc {
	return n.loc
}

func (n *CallExpr) Extra() interface{} {
	return n.extra
}

func (n *CallExpr) setExtra(ext interface{}) {
	n.extra = ext.(*ExprExtra)
}

type BinExpr struct {
	typ   NodeType
	loc   *Loc
	op    TokenValue
	lhs   Node
	rhs   Node
	extra *ExprExtra
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

func (n *BinExpr) Type() NodeType {
	return n.typ
}

func (n *BinExpr) Loc() *Loc {
	return n.loc
}

func (n *BinExpr) Extra() interface{} {
	return n.extra
}

func (n *BinExpr) setExtra(ext interface{}) {
	n.extra = ext.(*ExprExtra)
}

type UnaryExpr struct {
	typ   NodeType
	loc   *Loc
	op    TokenValue
	arg   Node
	extra *ExprExtra
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

func (n *UnaryExpr) Type() NodeType {
	return n.typ
}

func (n *UnaryExpr) Loc() *Loc {
	return n.loc
}

func (n *UnaryExpr) Extra() interface{} {
	return n.extra
}

func (n *UnaryExpr) setExtra(ext interface{}) {
	n.extra = ext.(*ExprExtra)
}

type UpdateExpr struct {
	typ    NodeType
	loc    *Loc
	op     TokenValue
	prefix bool
	arg    Node
	extra  *ExprExtra
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

func (n *UpdateExpr) Type() NodeType {
	return n.typ
}

func (n *UpdateExpr) Loc() *Loc {
	return n.loc
}

func (n *UpdateExpr) Extra() interface{} {
	return n.extra
}

func (n *UpdateExpr) setExtra(ext interface{}) {
	n.extra = ext.(*ExprExtra)
}

type CondExpr struct {
	typ   NodeType
	loc   *Loc
	test  Node
	cons  Node
	alt   Node
	extra *ExprExtra
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

func (n *CondExpr) Extra() interface{} {
	return n.extra
}

func (n *CondExpr) setExtra(ext interface{}) {
	n.extra = ext.(*ExprExtra)
}

type AssignExpr struct {
	typ   NodeType
	loc   *Loc
	op    TokenValue
	opLoc *Loc
	lhs   Node
	rhs   Node
	extra *ExprExtra
}

func (n *AssignExpr) Op() TokenValue {
	return n.op
}

func (n *AssignExpr) OpText() string {
	return TokenKinds[n.op].Name
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

func (n *AssignExpr) Extra() interface{} {
	return n.extra
}

func (n *AssignExpr) setExtra(ext interface{}) {
	n.extra = ext.(*ExprExtra)
}

type ThisExpr struct {
	typ   NodeType
	loc   *Loc
	extra *ExprExtra
}

func (n *ThisExpr) Type() NodeType {
	return n.typ
}

func (n *ThisExpr) Loc() *Loc {
	return n.loc
}

func (n *ThisExpr) Extra() interface{} {
	return n.extra
}

func (n *ThisExpr) setExtra(ext interface{}) {
	n.extra = ext.(*ExprExtra)
}

type SeqExpr struct {
	typ   NodeType
	loc   *Loc
	elems []Node
	extra *ExprExtra
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

func (n *SeqExpr) Extra() interface{} {
	return n.extra
}

func (n *SeqExpr) setExtra(ext interface{}) {
	n.extra = ext.(*ExprExtra)
}

type ParenExpr struct {
	typ   NodeType
	loc   *Loc
	expr  Node
	extra *ExprExtra
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

func (n *ParenExpr) Extra() interface{} {
	return n.extra
}

func (n *ParenExpr) setExtra(ext interface{}) {
	n.extra = ext.(*ExprExtra)
}

// there is no information kept to describe the program order of the quasis and expresions
// according to below link descries how the quasis and expresion are being walk over:
// https://opensource.apple.com/source/WebInspectorUI/WebInspectorUI-7602.2.14.0.5/UserInterface/Workers/Formatter/ESTreeWalker.js.auto.html
// some meaningless output should be taken into its estree result, such as put first quasis as
// a emptry string if the first element in `elems` is a expression
type TplExpr struct {
	typ   NodeType
	loc   *Loc
	tag   Node
	elems []Node
	extra *ExprExtra
}

func (n *TplExpr) Tag() Node {
	return n.tag
}

func (n *TplExpr) Elems() []Node {
	return n.elems
}

func (n *TplExpr) Type() NodeType {
	return n.typ
}

// loc without the `tag`
func (n *TplExpr) Loc() *Loc {
	return n.loc
}

func (n *TplExpr) Extra() interface{} {
	return n.extra
}

func (n *TplExpr) setExtra(ext interface{}) {
	n.extra = ext.(*ExprExtra)
}

func (n *TplExpr) LocWithTag() *Loc {
	loc := n.loc.Clone()
	if n.tag != nil {
		tl := n.tag.Loc()
		loc.begin = tl.begin.Clone()
		loc.rng.start = tl.rng.start
	}
	return loc
}

type Super struct {
	typ   NodeType
	loc   *Loc
	extra *ExprExtra
}

func (n *Super) Type() NodeType {
	return n.typ
}

func (n *Super) Loc() *Loc {
	return n.loc
}

func (n *Super) Extra() interface{} {
	return n.extra
}

func (n *Super) setExtra(ext interface{}) {
	n.extra = ext.(*ExprExtra)
}

type ImportCall struct {
	typ   NodeType
	loc   *Loc
	src   Node
	extra *ExprExtra
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

func (n *ImportCall) Extra() interface{} {
	return n.extra
}

func (n *ImportCall) setExtra(ext interface{}) {
	n.extra = ext.(*ExprExtra)
}

type YieldExpr struct {
	typ      NodeType
	loc      *Loc
	delegate bool
	arg      Node
	extra    *ExprExtra
}

func (n *YieldExpr) Delegate() bool {
	return n.delegate
}

func (n *YieldExpr) Arg() Node {
	return n.arg
}

func (n *YieldExpr) Type() NodeType {
	return n.typ
}

func (n *YieldExpr) Loc() *Loc {
	return n.loc
}

func (n *YieldExpr) Extra() interface{} {
	return n.extra
}

func (n *YieldExpr) setExtra(ext interface{}) {
	n.extra = ext.(*ExprExtra)
}

type ArrPat struct {
	typ      NodeType
	loc      *Loc
	elems    []Node
	extra    *ExprExtra
	typAnnot Node
}

func (n *ArrPat) Type() NodeType {
	return n.typ
}

func (n *ArrPat) Loc() *Loc {
	return n.loc
}

func (n *ArrPat) Elems() []Node {
	return n.elems
}

func (n *ArrPat) Extra() interface{} {
	return n.extra
}

func (n *ArrPat) setExtra(ext interface{}) {
	n.extra = ext.(*ExprExtra)
}

type AssignPat struct {
	typ      NodeType
	loc      *Loc
	lhs      Node
	rhs      Node
	extra    *ExprExtra
	typAnnot Node
}

func (n *AssignPat) Left() Node {
	return n.lhs
}

func (n *AssignPat) Right() Node {
	return n.rhs
}

func (n *AssignPat) Type() NodeType {
	return n.typ
}

func (n *AssignPat) Loc() *Loc {
	return n.loc
}

func (n *AssignPat) Extra() interface{} {
	return n.extra
}

func (n *AssignPat) setExtra(ext interface{}) {
	n.extra = ext.(*ExprExtra)
}

type RestPat struct {
	typ      NodeType
	loc      *Loc
	arg      Node
	extra    *ExprExtra
	typAnnot Node
}

func (n *RestPat) Arg() Node {
	return n.arg
}

func (n *RestPat) Type() NodeType {
	return n.typ
}

func (n *RestPat) Loc() *Loc {
	return n.loc
}

func (n *RestPat) Extra() interface{} {
	return n.extra
}

func (n *RestPat) setExtra(ext interface{}) {
	n.extra = ext.(*ExprExtra)
}

type ObjPat struct {
	typ      NodeType
	loc      *Loc
	props    []Node
	extra    *ExprExtra
	typAnnot Node
}

func (n *ObjPat) Props() []Node {
	return n.props
}

func (n *ObjPat) Type() NodeType {
	return n.typ
}

func (n *ObjPat) Loc() *Loc {
	return n.loc
}

func (n *ObjPat) Extra() interface{} {
	return n.extra
}

func (n *ObjPat) setExtra(ext interface{}) {
	n.extra = ext.(*ExprExtra)
}

type PropKind uint8

const (
	PK_INIT PropKind = iota
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

type Prop struct {
	typ      NodeType
	loc      *Loc
	key      Node
	opLoc    *Loc
	value    Node
	computed bool
	method   bool

	shorthand bool
	// it's `true` if the prop value is in assign pattern
	assign bool

	kind PropKind
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

func (n *Prop) Value() Node {
	return n.value
}

func (n *Prop) Computed() bool {
	return n.computed
}

func (n *Prop) Shorthand() bool {
	return n.shorthand
}

func (n *Prop) Type() NodeType {
	return n.typ
}

func (n *Prop) Loc() *Loc {
	return n.loc
}

func (n *Prop) Extra() interface{} {
	return nil
}

func (n *Prop) setExtra(_ interface{}) {
}

type FnDec struct {
	typ       NodeType
	loc       *Loc
	id        Node
	generator bool
	async     bool
	params    []Node
	body      Node
	extra     *ExprExtra
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

func (n *FnDec) Extra() interface{} {
	return n.extra
}

func (n *FnDec) setExtra(ext interface{}) {
	n.extra = ext.(*ExprExtra)
}

type ArrowFn struct {
	typ      NodeType
	loc      *Loc
	arrowLoc *Loc
	async    bool
	params   []Node
	body     Node
	expr     bool
	extra    *ExprExtra
}

func (n *ArrowFn) Async() bool {
	return n.async
}

func (n *ArrowFn) Params() []Node {
	return n.params
}

func (n *ArrowFn) Body() Node {
	return n.body
}

func (n *ArrowFn) Type() NodeType {
	return n.typ
}

func (n *ArrowFn) Loc() *Loc {
	return n.loc
}

func (n *ArrowFn) Extra() interface{} {
	return n.extra
}

func (n *ArrowFn) setExtra(ext interface{}) {
	n.extra = ext.(*ExprExtra)
}

func (n *ArrowFn) Expr() bool {
	return n.expr
}

type VarDecStmt struct {
	typ     NodeType
	loc     *Loc
	kind    TokenValue
	decList []Node
	names   []Node
}

func (n *VarDecStmt) Kind() string {
	return TokenKinds[n.kind].Name
}

func (n *VarDecStmt) DecList() []Node {
	return n.decList
}

func (n *VarDecStmt) Type() NodeType {
	return n.typ
}

func (n *VarDecStmt) Loc() *Loc {
	return n.loc
}

func (n *VarDecStmt) Extra() interface{} {
	return nil
}

func (n *VarDecStmt) setExtra(_ interface{}) {
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

func (n *VarDec) Extra() interface{} {
	return nil
}

func (n *VarDec) setExtra(_ interface{}) {
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

func (n *BlockStmt) Extra() interface{} {
	return nil
}

func (n *BlockStmt) setExtra(_ interface{}) {
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

func (n *DoWhileStmt) Extra() interface{} {
	return nil
}

func (n *DoWhileStmt) setExtra(_ interface{}) {
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

func (n *WhileStmt) Extra() interface{} {
	return nil
}

func (n *WhileStmt) setExtra(_ interface{}) {
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

func (n *ForStmt) Extra() interface{} {
	return nil
}

func (n *ForStmt) setExtra(_ interface{}) {
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

func (n *ForInOfStmt) Extra() interface{} {
	return nil
}

func (n *ForInOfStmt) setExtra(_ interface{}) {
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

func (n *IfStmt) Extra() interface{} {
	return nil
}

func (n *IfStmt) setExtra(_ interface{}) {
}

type SwitchStmt struct {
	typ   NodeType
	loc   *Loc
	test  Node
	cases []Node
}

func (n *SwitchStmt) Cases() []Node {
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

func (n *SwitchStmt) Extra() interface{} {
	return nil
}

func (n *SwitchStmt) setExtra(_ interface{}) {
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

func (n *SwitchCase) Extra() interface{} {
	return nil
}

func (n *SwitchCase) setExtra(_ interface{}) {
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

func (n *BrkStmt) Extra() interface{} {
	return nil
}

func (n *BrkStmt) setExtra(_ interface{}) {
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

func (n *ContStmt) Extra() interface{} {
	return nil
}

func (n *ContStmt) setExtra(_ interface{}) {
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

func (n *LabelStmt) Extra() interface{} {
	return nil
}

func (n *LabelStmt) setExtra(_ interface{}) {
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

func (n *RetStmt) Extra() interface{} {
	return nil
}

func (n *RetStmt) setExtra(_ interface{}) {
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

func (n *ThrowStmt) Extra() interface{} {
	return nil
}

func (n *ThrowStmt) setExtra(_ interface{}) {
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

func (n *Catch) Extra() interface{} {
	return nil
}

func (n *Catch) setExtra(_ interface{}) {
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

func (n *TryStmt) Extra() interface{} {
	return nil
}

func (n *TryStmt) setExtra(_ interface{}) {
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

func (n *DebugStmt) Extra() interface{} {
	return nil
}

func (n *DebugStmt) setExtra(_ interface{}) {
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

func (n *WithStmt) Extra() interface{} {
	return nil
}

func (n *WithStmt) setExtra(_ interface{}) {
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

func (n *ClassDec) Extra() interface{} {
	return nil
}

func (n *ClassDec) setExtra(_ interface{}) {
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

func (n *ClassBody) Extra() interface{} {
	return nil
}

func (n *ClassBody) setExtra(_ interface{}) {
}

type Method struct {
	typ      NodeType
	loc      *Loc
	key      Node
	static   bool
	computed bool
	kind     PropKind
	value    Node
}

func (n *Method) Kind() string {
	return n.kind.ToString()
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

func (n *Method) Extra() interface{} {
	return nil
}

func (n *Method) setExtra(_ interface{}) {
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
	return n.value
}

func (n *Field) Static() bool {
	return n.static
}

func (n *Field) Computed() bool {
	return n.computed
}

func (n *Field) Type() NodeType {
	return n.typ
}

func (n *Field) Loc() *Loc {
	return n.loc
}

func (n *Field) Extra() interface{} {
	return nil
}

func (n *Field) setExtra(_ interface{}) {
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

func (n *StaticBlock) Extra() interface{} {
	return nil
}

func (n *StaticBlock) Body() []Node {
	return n.body
}

func (n *StaticBlock) setExtra(_ interface{}) {
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

func (n *MetaProp) Extra() interface{} {
	return nil
}

func (n *MetaProp) setExtra(_ interface{}) {
}

type ImportDec struct {
	typ   NodeType
	loc   *Loc
	specs []Node
	src   Node
}

func (n *ImportDec) Specs() []Node {
	return n.specs
}

func (n *ImportDec) Src() Node {
	return n.src
}

func (n *ImportDec) Type() NodeType {
	return n.typ
}

func (n *ImportDec) Loc() *Loc {
	return n.loc
}

func (n *ImportDec) Extra() interface{} {
	return nil
}

func (n *ImportDec) setExtra(_ interface{}) {
}

type ImportSpec struct {
	typ   NodeType
	loc   *Loc
	def   bool
	ns    bool
	local Node
	id    Node
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

func (n *ImportSpec) Type() NodeType {
	return n.typ
}

func (n *ImportSpec) Loc() *Loc {
	return n.loc
}

func (n *ImportSpec) Extra() interface{} {
	return nil
}

func (n *ImportSpec) setExtra(_ interface{}) {
}

type ExportDec struct {
	typ   NodeType
	loc   *Loc
	all   bool
	def   *Loc
	dec   Node
	specs []Node
	src   Node
}

func (n *ExportDec) All() bool {
	return n.all
}

func (n *ExportDec) Default() bool {
	return n.def != nil
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

func (n *ExportDec) Type() NodeType {
	return n.typ
}

func (n *ExportDec) Loc() *Loc {
	return n.loc
}

func (n *ExportDec) Extra() interface{} {
	return nil
}

func (n *ExportDec) setExtra(_ interface{}) {
}

type ExportSpec struct {
	typ   NodeType
	loc   *Loc
	ns    bool
	local Node
	id    Node // exported
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

func (n *ExportSpec) Type() NodeType {
	return n.typ
}

func (n *ExportSpec) Loc() *Loc {
	return n.loc
}

func (n *ExportSpec) Extra() interface{} {
	return nil
}

func (n *ExportSpec) setExtra(_ interface{}) {
}

type ChainExpr struct {
	typ  NodeType
	loc  *Loc
	expr Node
}

func (n *ChainExpr) Type() NodeType {
	return n.typ
}

func (n *ChainExpr) Loc() *Loc {
	return n.loc
}

func (n *ChainExpr) Expr() Node {
	return n.expr
}

func (n *ChainExpr) Extra() interface{} {
	return nil
}

func (n *ChainExpr) setExtra(ext interface{}) {

}
