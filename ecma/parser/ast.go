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

func (r *Range) SetStart(n int) {
	r.start = n
}

func (r *Range) SetEnd(n int) {
	r.end = n
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

type InParenNode interface {
	OuterParen() *Loc
	SetOuterParen(*Loc)
}

type NullLit struct {
	typ        NodeType
	loc        *Loc
	outerParen *Loc
	ti         *TypInfo
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

func (n *NullLit) OuterParen() *Loc {
	return n.outerParen
}

func (n *NullLit) SetOuterParen(loc *Loc) {
	n.outerParen = loc
}

func (n *NullLit) TypInfo() *TypInfo {
	return n.ti
}

func (n *NullLit) SetTypInfo(ti *TypInfo) {
	n.ti = ti
}

type BoolLit struct {
	typ        NodeType
	loc        *Loc
	val        bool
	outerParen *Loc
	ti         *TypInfo
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

func (n *BoolLit) OuterParen() *Loc {
	return n.outerParen
}

func (n *BoolLit) SetOuterParen(loc *Loc) {
	n.outerParen = loc
}

func (n *BoolLit) TypInfo() *TypInfo {
	return n.ti
}

func (n *BoolLit) SetTypInfo(ti *TypInfo) {
	n.ti = ti
}

type NumLit struct {
	typ        NodeType
	loc        *Loc
	outerParen *Loc
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

// unsafe method, use `IsBigint` before this method
// return the numerical value if it's safe to be represented in JSON as number
// otherwise zero is returned
func (n *NumLit) Float() float64 {
	if n.IsBigint() {
		i := n.ToBigint()
		f := big.NewFloat(0)
		f.SetInt(i)
		max := big.NewFloat(0)
		max.SetUint64(math.MaxInt)
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

func (n *NumLit) OuterParen() *Loc {
	return n.outerParen
}

func (n *NumLit) SetOuterParen(loc *Loc) {
	n.outerParen = loc
}

type StrLit struct {
	typ                  NodeType
	loc                  *Loc
	val                  string
	legacyOctalEscapeSeq bool
	outerParen           *Loc
	ti                   *TypInfo
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

func (n *StrLit) OuterParen() *Loc {
	return n.outerParen
}

func (n *StrLit) SetOuterParen(loc *Loc) {
	n.outerParen = loc
}

func (n *StrLit) TypInfo() *TypInfo {
	return n.ti
}

func (n *StrLit) SetTypInfo(ti *TypInfo) {
	n.ti = ti
}

type RegLit struct {
	typ        NodeType
	loc        *Loc
	val        string
	pattern    string
	flags      string
	outerParen *Loc
	ti         *TypInfo
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

func (n *RegLit) OuterParen() *Loc {
	return n.outerParen
}

func (n *RegLit) SetOuterParen(loc *Loc) {
	n.outerParen = loc
}

func (n *RegLit) TypInfo() *TypInfo {
	return n.ti
}

func (n *RegLit) SetTypInfo(ti *TypInfo) {
	n.ti = ti
}

type ArrLit struct {
	typ        NodeType
	loc        *Loc
	elems      []Node
	outerParen *Loc
	ti         *TypInfo
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

func (n *ArrLit) OuterParen() *Loc {
	return n.outerParen
}

func (n *ArrLit) SetOuterParen(loc *Loc) {
	n.outerParen = loc
}

func (n *ArrLit) TypInfo() *TypInfo {
	return n.ti
}

func (n *ArrLit) SetTypInfo(ti *TypInfo) {
	n.ti = ti
}

type Spread struct {
	typ              NodeType
	loc              *Loc
	arg              Node
	trailingCommaLoc *Loc
	outerParen       *Loc
	ti               *TypInfo
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

func (n *Spread) OuterParen() *Loc {
	return n.outerParen
}

func (n *Spread) SetOuterParen(loc *Loc) {
	n.outerParen = loc
}

func (n *Spread) TypInfo() *TypInfo {
	return n.ti
}

func (n *Spread) SetTypInfo(ti *TypInfo) {
	n.ti = ti
}

type ObjLit struct {
	typ        NodeType
	loc        *Loc
	props      []Node
	outerParen *Loc
	ti         *TypInfo
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

func (n *ObjLit) OuterParen() *Loc {
	return n.outerParen
}

func (n *ObjLit) SetOuterParen(loc *Loc) {
	n.outerParen = loc
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

type TypInfo struct {
	accMod    ACC_MOD
	ques      *Loc
	typAnnot  Node
	typParams Node
	typArgs   Node
}

func (ti *TypInfo) Ques() *Loc {
	return ti.ques
}

func (ti *TypInfo) Optional() bool {
	return ti.ques != nil
}

func (ti *TypInfo) TypAnnot() Node {
	return ti.typAnnot
}

func (ti *TypInfo) TypParams() Node {
	return ti.typParams
}

func (ti *TypInfo) TypArgs() Node {
	return ti.typArgs
}

type NodeWithTypInfo interface {
	TypInfo() *TypInfo
	SetTypInfo(*TypInfo)
}

func locOfNode(node Node) *Loc {
	if node == nil {
		return nil
	}
	return node.Loc()
}

func startOf(locs ...*Loc) *Loc {
	start := 0
	line := math.MaxInt
	col := math.MaxInt
	for i, loc := range locs {
		if loc == nil {
			continue
		}
		if loc.begin.line < line || (loc.begin.line == line && loc.begin.col < col) {
			line = loc.begin.line
			col = loc.begin.col
			start = i
		}
	}
	return locs[start]
}

func endOf(locs ...*Loc) *Loc {
	end := 0
	line := -1
	col := -1
	for i, loc := range locs {
		if loc == nil {
			continue
		}
		if loc.end.line > line || (loc.end.line == line && loc.end.col > col) {
			line = loc.end.line
			col = loc.end.col
			end = i
		}
	}
	return locs[end]
}

func LocWithTypeInfo(node Node) *Loc {
	nw, ok := node.(NodeWithTypInfo)
	if !ok {
		return node.Loc()
	}

	ti := nw.TypInfo()
	loc := node.Loc().Clone()

	start := startOf(locOfNode(ti.TypParams()), locOfNode(node))
	loc.begin.line = start.begin.line
	loc.begin.col = start.begin.col
	loc.rng.start = start.rng.start

	end := endOf(locOfNode(ti.TypAnnot()), ti.Ques(), locOfNode(node))
	loc.end.line = end.end.line
	loc.end.col = end.end.col
	loc.rng.end = end.rng.end

	return loc
}

type Ident struct {
	typ            NodeType
	loc            *Loc
	val            string
	pvt            bool
	containsEscape bool
	outerParen     *Loc

	// consider below statements:
	// `export { if } from "a"` is legal
	// `export { if } ` is illegal
	// for reporting `if` is a keyword, firstly produce a
	// Ident with conent `if` and flag it's a keyword by
	// setting this field to true, later report the `unexpected token`
	// error if the coming token is not `from`
	kw bool

	ti *TypInfo
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

func (n *Ident) OuterParen() *Loc {
	return n.outerParen
}

func (n *Ident) SetOuterParen(loc *Loc) {
	n.outerParen = loc
}

func (n *Ident) TypInfo() *TypInfo {
	return n.ti
}

func (n *Ident) SetTypInfo(ti *TypInfo) {
	n.ti = ti
}

type NewExpr struct {
	typ        NodeType
	loc        *Loc
	callee     Node
	args       []Node
	outerParen *Loc
	ti         *TypInfo
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

func (n *NewExpr) OuterParen() *Loc {
	return n.outerParen
}

func (n *NewExpr) SetOuterParen(loc *Loc) {
	n.outerParen = loc
}

func (n *NewExpr) TypInfo() *TypInfo {
	return n.ti
}

func (n *NewExpr) SetTypInfo(ti *TypInfo) {
	n.ti = ti
}

type MemberExpr struct {
	typ        NodeType
	loc        *Loc
	obj        Node
	prop       Node
	compute    bool
	optional   bool
	outerParen *Loc
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

func (n *MemberExpr) OuterParen() *Loc {
	return n.outerParen
}

func (n *MemberExpr) SetOuterParen(loc *Loc) {
	n.outerParen = loc
}

type CallExpr struct {
	typ        NodeType
	loc        *Loc
	callee     Node
	args       []Node
	optional   bool
	outerParen *Loc
	ti         *TypInfo
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

func (n *CallExpr) OuterParen() *Loc {
	return n.outerParen
}

func (n *CallExpr) SetOuterParen(loc *Loc) {
	n.outerParen = loc
}

func (n *CallExpr) TypInfo() *TypInfo {
	return n.ti
}

func (n *CallExpr) SetTypInfo(ti *TypInfo) {
	n.ti = ti
}

type BinExpr struct {
	typ        NodeType
	loc        *Loc
	op         TokenValue
	opLoc      *Loc
	lhs        Node
	rhs        Node
	outerParen *Loc
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

func (n *BinExpr) OuterParen() *Loc {
	return n.outerParen
}

func (n *BinExpr) SetOuterParen(loc *Loc) {
	n.outerParen = loc
}

type UnaryExpr struct {
	typ        NodeType
	loc        *Loc
	op         TokenValue
	arg        Node
	outerParen *Loc
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

func (n *UnaryExpr) OuterParen() *Loc {
	return n.outerParen
}

func (n *UnaryExpr) SetOuterParen(loc *Loc) {
	n.outerParen = loc
}

type UpdateExpr struct {
	typ        NodeType
	loc        *Loc
	op         TokenValue
	prefix     bool
	arg        Node
	outerParen *Loc
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

func (n *UpdateExpr) OuterParen() *Loc {
	return n.outerParen
}

func (n *UpdateExpr) SetOuterParen(loc *Loc) {
	n.outerParen = loc
}

type CondExpr struct {
	typ        NodeType
	loc        *Loc
	test       Node
	cons       Node
	alt        Node
	outerParen *Loc
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

func (n *CondExpr) OuterParen() *Loc {
	return n.outerParen
}

func (n *CondExpr) SetOuterParen(loc *Loc) {
	n.outerParen = loc
}

type AssignExpr struct {
	typ        NodeType
	loc        *Loc
	op         TokenValue
	opLoc      *Loc
	lhs        Node
	rhs        Node
	outerParen *Loc
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

func (n *AssignExpr) OuterParen() *Loc {
	return n.outerParen
}

func (n *AssignExpr) SetOuterParen(loc *Loc) {
	n.outerParen = loc
}

type ThisExpr struct {
	typ        NodeType
	loc        *Loc
	outerParen *Loc
	ti         *TypInfo
}

func (n *ThisExpr) Type() NodeType {
	return n.typ
}

func (n *ThisExpr) Loc() *Loc {
	return n.loc
}

func (n *ThisExpr) OuterParen() *Loc {
	return n.outerParen
}

func (n *ThisExpr) SetOuterParen(loc *Loc) {
	n.outerParen = loc
}

func (n *ThisExpr) TypInfo() *TypInfo {
	return n.ti
}

func (n *ThisExpr) SetTypInfo(ti *TypInfo) {
	n.ti = ti
}

type SeqExpr struct {
	typ        NodeType
	loc        *Loc
	elems      []Node
	outerParen *Loc
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

func (n *SeqExpr) OuterParen() *Loc {
	return n.outerParen
}

func (n *SeqExpr) SetOuterParen(loc *Loc) {
	n.outerParen = loc
}

type ParenExpr struct {
	typ        NodeType
	loc        *Loc
	expr       Node
	outerParen *Loc
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

func (n *ParenExpr) OuterParen() *Loc {
	return n.outerParen
}

func (n *ParenExpr) SetOuterParen(loc *Loc) {
	n.outerParen = loc
}

// there is no information kept to describe the program order of the quasis and expresions
// according to below link descries how the quasis and expresion are being walk over:
// https://opensource.apple.com/source/WebInspectorUI/WebInspectorUI-7602.2.14.0.5/UserInterface/Workers/Formatter/ESTreeWalker.js.auto.html
// some meaningless output should be taken into its estree result, such as put first quasis as
// a emptry string if the first element in `elems` is a expression
type TplExpr struct {
	typ        NodeType
	loc        *Loc
	tag        Node
	elems      []Node
	outerParen *Loc
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

func (n *TplExpr) OuterParen() *Loc {
	return n.outerParen
}

func (n *TplExpr) SetOuterParen(loc *Loc) {
	n.outerParen = loc
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
	typ        NodeType
	loc        *Loc
	outerParen *Loc
	ti         *TypInfo
}

func (n *Super) Type() NodeType {
	return n.typ
}

func (n *Super) Loc() *Loc {
	return n.loc
}

func (n *Super) OuterParen() *Loc {
	return n.outerParen
}

func (n *Super) SetOuterParen(loc *Loc) {
	n.outerParen = loc
}

func (n *Super) TypInfo() *TypInfo {
	return n.ti
}

func (n *Super) SetTypInfo(ti *TypInfo) {
	n.ti = ti
}

type ImportCall struct {
	typ        NodeType
	loc        *Loc
	src        Node
	outerParen *Loc
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

func (n *ImportCall) OuterParen() *Loc {
	return n.outerParen
}

func (n *ImportCall) SetOuterParen(loc *Loc) {
	n.outerParen = loc
}

type YieldExpr struct {
	typ        NodeType
	loc        *Loc
	delegate   bool
	arg        Node
	outerParen *Loc
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

func (n *YieldExpr) OuterParen() *Loc {
	return n.outerParen
}

func (n *YieldExpr) SetOuterParen(loc *Loc) {
	n.outerParen = loc
}

type ArrPat struct {
	typ        NodeType
	loc        *Loc
	elems      []Node
	outerParen *Loc
	ti         *TypInfo
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

func (n *ArrPat) OuterParen() *Loc {
	return n.outerParen
}

func (n *ArrPat) SetOuterParen(loc *Loc) {
	n.outerParen = loc
}

func (n *ArrPat) TypInfo() *TypInfo {
	return n.ti
}

func (n *ArrPat) SetTypInfo(ti *TypInfo) {
	n.ti = ti
}

type AssignPat struct {
	typ        NodeType
	loc        *Loc
	lhs        Node
	rhs        Node
	outerParen *Loc
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

func (n *AssignPat) OuterParen() *Loc {
	return n.outerParen
}

type RestPat struct {
	typ        NodeType
	loc        *Loc
	arg        Node
	outerParen *Loc
	ti         *TypInfo
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

func (n *RestPat) OuterParen() *Loc {
	return n.outerParen
}

func (n *RestPat) SetOuterParen(loc *Loc) {
	n.outerParen = loc
}

func (n *RestPat) TypInfo() *TypInfo {
	return n.ti
}

func (n *RestPat) Optional() bool {
	return n.ti.ques != nil
}

func (n *RestPat) hoistTypInfo() {
	if wt, ok := n.arg.(NodeWithTypInfo); ok {
		n.ti = wt.TypInfo()
		wt.SetTypInfo(nil)
	}
}

type ObjPat struct {
	typ        NodeType
	loc        *Loc
	props      []Node
	outerParen *Loc
	ti         *TypInfo
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

func (n *ObjPat) OuterParen() *Loc {
	return n.outerParen
}

func (n *ObjPat) SetOuterParen(loc *Loc) {
	n.outerParen = loc
}

func (n *ObjPat) TypInfo() *TypInfo {
	return n.ti
}

func (n *ObjPat) SetTypInfo(ti *TypInfo) {
	n.ti = ti
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

	kind    PropKind
	accMode ACC_MOD
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

type FnDec struct {
	typ        NodeType
	loc        *Loc
	id         Node
	generator  bool
	async      bool
	params     []Node
	body       Node
	outerParen *Loc
	ti         *TypInfo
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

func (n *FnDec) IsSig() bool {
	return n.body == nil
}

func (n *FnDec) OuterParen() *Loc {
	return n.outerParen
}

func (n *FnDec) SetOuterParen(loc *Loc) {
	n.outerParen = loc
}

func (n *FnDec) TypInfo() *TypInfo {
	return n.ti
}

func (n *FnDec) SetTypInfo(ti *TypInfo) {
	n.ti = ti
}

type ArrowFn struct {
	typ        NodeType
	loc        *Loc
	arrowLoc   *Loc
	async      bool
	params     []Node
	body       Node
	expr       bool
	outerParen *Loc
	ti         *TypInfo
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

func (n *ArrowFn) OuterParen() *Loc {
	return n.outerParen
}

func (n *ArrowFn) SetOuterParen(loc *Loc) {
	n.outerParen = loc
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
	kind     PropKind
	value    Node
	accMode  ACC_MOD
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

type Field struct {
	typ      NodeType
	loc      *Loc
	key      Node
	static   bool
	computed bool
	value    Node
	accMode  ACC_MOD
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

func (n *StaticBlock) Body() []Node {
	return n.body
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
