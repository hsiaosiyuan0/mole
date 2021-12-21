package parser

type TsPredef struct {
	typ  NodeType
	loc  *Loc
	ques *Loc
}

func (n *TsPredef) Type() NodeType {
	return n.typ
}

func (n *TsPredef) Loc() *Loc {
	return n.loc
}

func (n *TsPredef) Text() string {
	return n.loc.Text()
}

type TsThis struct {
	typ NodeType
	loc *Loc
}

func (n *TsThis) Type() NodeType {
	return n.typ
}

func (n *TsThis) Loc() *Loc {
	return n.loc
}

type TsNsName struct {
	typ NodeType
	loc *Loc
	lhs Node
	dot *Loc
	rhs Node
}

func (n *TsNsName) Type() NodeType {
	return n.typ
}

func (n *TsNsName) Loc() *Loc {
	return n.loc
}

type TsRef struct {
	typ  NodeType
	loc  *Loc
	name Node
	lt   *Loc
	args []Node
}

func (n *TsRef) Type() NodeType {
	return n.typ
}

func (n *TsRef) Loc() *Loc {
	return n.loc
}

func (n *TsRef) HasArgs() bool {
	return n.args != nil && len(n.args) > 0
}

type TsQuery struct {
	typ NodeType
	loc *Loc
	arg Node // name or nsName
}

func (n *TsQuery) Type() NodeType {
	return n.typ
}

func (n *TsQuery) Loc() *Loc {
	return n.loc
}

type TsParen struct {
	typ NodeType
	loc *Loc
	arg Node // name or nsName
}

func (n *TsParen) Type() NodeType {
	return n.typ
}

func (n *TsParen) Loc() *Loc {
	return n.loc
}

type TsArr struct {
	typ     NodeType
	loc     *Loc
	bracket *Loc
	arg     Node // name or nsName
}

func (n *TsArr) Type() NodeType {
	return n.typ
}

func (n *TsArr) Loc() *Loc {
	return n.loc
}

type TsTuple struct {
	typ  NodeType
	loc  *Loc
	args []Node
}

func (n *TsTuple) Type() NodeType {
	return n.typ
}

func (n *TsTuple) Loc() *Loc {
	return n.loc
}

type TsObj struct {
	typ   NodeType
	loc   *Loc
	props []Node
}

func (n *TsObj) Type() NodeType {
	return n.typ
}

func (n *TsObj) Loc() *Loc {
	return n.loc
}

type TsProp struct {
	typ        NodeType
	loc        *Loc
	key        Node
	val        Node
	ques       *Loc
	computeLoc *Loc
}

func (n *TsProp) Type() NodeType {
	return n.typ
}

func (n *TsProp) Loc() *Loc {
	return n.loc
}

type TsCallSig struct {
	typ       NodeType
	loc       *Loc
	typParams []Node
	params    []Node
	retTyp    Node
}

func (n *TsCallSig) Type() NodeType {
	return n.typ
}

func (n *TsCallSig) Loc() *Loc {
	return n.loc
}

type TsNewSig struct {
	typ       NodeType
	loc       *Loc
	typParams []Node
	params    []Node
	retTyp    Node
}

func (n *TsNewSig) Type() NodeType {
	return n.typ
}

func (n *TsNewSig) Loc() *Loc {
	return n.loc
}

type TsIdxSig struct {
	typ NodeType
	loc *Loc
	id  Node
	key Node
	val Node
}

func (n *TsIdxSig) Type() NodeType {
	return n.typ
}

func (n *TsIdxSig) Loc() *Loc {
	return n.loc
}

type TsRoughParam struct {
	typ   NodeType
	loc   *Loc
	name  Node
	colon *Loc
	ti    *TypInfo
}

func (n *TsRoughParam) Type() NodeType {
	return n.typ
}

func (n *TsRoughParam) Loc() *Loc {
	return n.loc
}

type TsParam struct {
	typ  NodeType
	loc  *Loc
	name Node
	cons Node
}

func (n *TsParam) Type() NodeType {
	return n.typ
}

func (n *TsParam) Loc() *Loc {
	return n.loc
}

type TsFnTyp struct {
	typ       NodeType
	loc       *Loc
	typParams []Node
	params    []Node
	retTyp    Node
}

func (n *TsFnTyp) Type() NodeType {
	return n.typ
}

func (n *TsFnTyp) Loc() *Loc {
	return n.loc
}

type TsUnionTyp struct {
	typ NodeType
	loc *Loc
	lhs Node
	op  *Loc
	rhs Node
}

func (n *TsUnionTyp) Type() NodeType {
	return n.typ
}

func (n *TsUnionTyp) Loc() *Loc {
	return n.loc
}

type TsIntersecTyp struct {
	typ NodeType
	loc *Loc
	lhs Node
	op  *Loc
	rhs Node
}

func (n *TsIntersecTyp) Type() NodeType {
	return n.typ
}

func (n *TsIntersecTyp) Loc() *Loc {
	return n.loc
}

type TsTypAssert struct {
	typ NodeType
	loc *Loc
	des Node
	arg Node
}

func (n *TsTypAssert) Type() NodeType {
	return n.typ
}

func (n *TsTypAssert) Loc() *Loc {
	return n.loc
}

type TsTypDec struct {
	typ  NodeType
	loc  *Loc
	name Node
	ti   *TypInfo
}

func (n *TsTypDec) Type() NodeType {
	return n.typ
}

func (n *TsTypDec) Loc() *Loc {
	return n.loc
}

type TsInferface struct {
	typ    NodeType
	loc    *Loc
	name   Node
	params []Node
	supers []Node
	body   Node
}

func (n *TsInferface) Type() NodeType {
	return n.typ
}

func (n *TsInferface) Loc() *Loc {
	return n.loc
}

type TsEnum struct {
	typ  NodeType
	loc  *Loc
	name Node
	mems []Node
	cons bool
}

func (n *TsEnum) Type() NodeType {
	return n.typ
}

func (n *TsEnum) Loc() *Loc {
	return n.loc
}

type TsEnumMember struct {
	typ NodeType
	loc *Loc
	key Node
	val Node
}

func (n *TsEnumMember) Type() NodeType {
	return n.typ
}

func (n *TsEnumMember) Loc() *Loc {
	return n.loc
}

type TsImportAlias struct {
	typ    NodeType
	loc    *Loc
	name   Node
	val    Node
	export bool
}

func (n *TsImportAlias) Type() NodeType {
	return n.typ
}

func (n *TsImportAlias) Loc() *Loc {
	return n.loc
}
