package parser

import "github.com/hsiaosiyuan0/mole/fuzz"

type TsTypAnnot struct {
	typ   NodeType
	loc   *Loc
	tsTyp Node
}

func NewTsTypAnnot(typ Node) *TsTypAnnot {
	return &TsTypAnnot{
		N_TS_TYP_ANNOT,
		typ.Loc().Clone(),
		typ,
	}
}

func (n *TsTypAnnot) Type() NodeType {
	return n.typ
}

func (n *TsTypAnnot) Loc() *Loc {
	return n.loc
}

func (n *TsTypAnnot) TsTyp() Node {
	return n.tsTyp
}

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

type TsLit struct {
	typ NodeType
	loc *Loc
	lit Node
}

func (n *TsLit) Type() NodeType {
	return n.typ
}

func (n *TsLit) Loc() *Loc {
	return n.loc
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

func (n *TsNsName) Lhs() Node {
	return n.lhs
}

func (n *TsNsName) Rhs() Node {
	return n.rhs
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
	args Node
}

func (n *TsRef) Type() NodeType {
	return n.typ
}

func (n *TsRef) Loc() *Loc {
	return n.loc
}

func (n *TsRef) Name() Node {
	return n.name
}

func (n *TsRef) Args() []Node {
	if n.HasArgs() {
		return n.args.(*TsParamsInst).params
	}
	return []Node{}
}

func (n *TsRef) ParamsInst() Node {
	return n.args
}

func (n *TsRef) HasArgs() bool {
	return n.args != nil && len(n.args.(*TsParamsInst).params) > 0
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

func (n *TsArr) Arg() Node {
	return n.arg
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

func (n *TsObj) Props() []Node {
	return n.props
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

func (n *TsProp) Key() Node {
	return n.key
}

func (n *TsProp) Val() Node {
	return n.val
}

func (n *TsProp) Optional() bool {
	return n.ques != nil
}

func (n *TsProp) Computed() bool {
	return n.computeLoc != nil
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
	typParams Node
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
	typParams Node
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

func (n *TsIdxSig) Key() Node {
	return n.id
}

func (n *TsIdxSig) KeyType() Node {
	return n.key
}

func (n *TsIdxSig) Value() Node {
	return n.val
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

type TsParamsInst struct {
	typ    NodeType
	loc    *Loc
	params []Node
}

func (n *TsParamsInst) Type() NodeType {
	return n.typ
}

func (n *TsParamsInst) Loc() *Loc {
	return n.loc
}

func (n *TsParamsInst) Params() []Node {
	return n.params
}

type TsParamsDec struct {
	typ    NodeType
	loc    *Loc
	params []Node
}

func (n *TsParamsDec) Type() NodeType {
	return n.typ
}

func (n *TsParamsDec) Loc() *Loc {
	return n.loc
}

func (n *TsParamsDec) Params() []Node {
	return n.params
}

type TsParam struct {
	typ  NodeType
	loc  *Loc
	name Node
	cons Node // the constraint
	val  Node // the default
}

func (n *TsParam) Type() NodeType {
	return n.typ
}

func (n *TsParam) Loc() *Loc {
	return n.loc
}

func (n *TsParam) Name() Node {
	return n.name
}

func (n *TsParam) Cons() Node {
	return n.cons
}

func (n *TsParam) Default() Node {
	return n.val
}

type TsFnTyp struct {
	typ       NodeType
	loc       *Loc
	typParams Node
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
	typ   NodeType
	loc   *Loc
	op    *Loc
	elems []Node
}

func (n *TsUnionTyp) Type() NodeType {
	return n.typ
}

func (n *TsUnionTyp) Loc() *Loc {
	return n.loc
}

func (n *TsUnionTyp) Elems() []Node {
	return n.elems
}

type TsIntersecTyp struct {
	typ   NodeType
	loc   *Loc
	op    *Loc
	elems []Node
}

func (n *TsIntersecTyp) Type() NodeType {
	return n.typ
}

func (n *TsIntersecTyp) Loc() *Loc {
	return n.loc
}

func (n *TsIntersecTyp) Elems() []Node {
	return n.elems
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

func (n *TsTypAssert) Typ() Node {
	return n.des
}

func (n *TsTypAssert) Expr() Node {
	return n.arg
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
	params Node
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
	typ   NodeType
	loc   *Loc
	name  Node
	items []Node
	cons  bool
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

type TsNS struct {
	typ   NodeType
	loc   *Loc
	name  Node
	stmts []Node
}

func (n *TsNS) Type() NodeType {
	return n.typ
}

func (n *TsNS) Loc() *Loc {
	return n.loc
}

type TsImportRequire struct {
	typ  NodeType
	loc  *Loc
	name Node
	args []Node
}

func (n *TsImportRequire) Type() NodeType {
	return n.typ
}

func (n *TsImportRequire) Loc() *Loc {
	return n.loc
}

type TsExportAssign struct {
	typ  NodeType
	loc  *Loc
	name Node
}

func (n *TsExportAssign) Type() NodeType {
	return n.typ
}

func (n *TsExportAssign) Loc() *Loc {
	return n.loc
}

type TsDec struct {
	typ   NodeType
	loc   *Loc
	name  Node
	inner Node
}

func (n *TsDec) Type() NodeType {
	return n.typ
}

func (n *TsDec) Loc() *Loc {
	return n.loc
}

func (n *TsDec) Name() Node {
	return n.name
}

func (n *TsDec) Inner() Node {
	return n.inner
}

// [assertion](https://www.typescriptlang.org/docs/handbook/release-notes/typescript-3-7.html#assertion-functions) and
// [type-predicates](https://www.typescriptlang.org/docs/handbook/2/narrowing.html#using-type-predicates) have almost
// the same syntax so they'll share the same definition by using `assert` to distinguish them
type TsTypPredicate struct {
	typ    NodeType
	loc    *Loc
	name   Node
	des    Node
	assert bool
}

func (n *TsTypPredicate) Type() NodeType {
	return n.typ
}

func (n *TsTypPredicate) Loc() *Loc {
	return n.loc
}

func (n *TsTypPredicate) Name() Node {
	return n.name
}

func (n *TsTypPredicate) Typ() Node {
	return n.des
}

func (n *TsTypPredicate) Asserts() bool {
	return n.assert
}

type TsNoNull struct {
	typ NodeType
	loc *Loc
	arg Node
}

func (n *TsNoNull) Type() NodeType {
	return n.typ
}

func (n *TsNoNull) Loc() *Loc {
	return n.loc
}

func (n *TsNoNull) Arg() Node {
	return n.arg
}

// ClassDec
func (n *ClassDec) Implements() []Node {
	if fuzz.IsNilPtr(n.ti) {
		return nil
	}
	return n.ti.Implements()
}

func (n *ClassDec) Abstract() bool {
	if fuzz.IsNilPtr(n.ti) {
		return false
	}
	return n.ti.Abstract()
}

func (n *ClassDec) SuperTypArgs() Node {
	super := n.super
	if super == nil {
		return nil
	}

	wt, ok := super.(NodeWithTypInfo)
	if !ok {
		return nil
	}

	st := wt.TypInfo()
	if st == nil {
		return nil
	}

	switch super.Type() {
	case N_EXPR_CALL:
		return st.SuperTypArgs()
	}
	return st.TypArgs()
}

func (n *ClassDec) TypParams() Node {
	if fuzz.IsNilPtr(n.ti) {
		return nil
	}
	return n.ti.typParams
}

// Field
func (n *Field) IsTsSig() bool {
	return n.ti != nil && n.val == nil && n.computed
}
