package parser

import (
	span "github.com/hsiaosiyuan0/mole/span"
	"github.com/hsiaosiyuan0/mole/util"
)

// #[visitor(TsTyp)]
type TsTypAnnot struct {
	typ   NodeType
	rng   span.Range
	tsTyp Node
}

func NewTsTypAnnot(typ Node) *TsTypAnnot {
	return &TsTypAnnot{N_TS_TYP_ANNOT, typ.Range(), typ}
}

func (n *TsTypAnnot) Type() NodeType {
	return n.typ
}

func (n *TsTypAnnot) Range() span.Range {
	return n.rng
}

func (n *TsTypAnnot) TsTyp() Node {
	return n.tsTyp
}

type TsPredef struct {
	typ  NodeType
	rng  span.Range
	ques span.Range
	opa  span.Range
}

func (n *TsPredef) Type() NodeType {
	return n.typ
}

func (n *TsPredef) Range() span.Range {
	return n.rng
}

func (n *TsPredef) OuterParen() span.Range {
	return n.opa
}

func (n *TsPredef) SetOuterParen(rng span.Range) {
	n.opa = rng
}

// #[visitor(Lit)]
type TsLit struct {
	typ NodeType
	rng span.Range
	lit Node
	opa span.Range
}

func (n *TsLit) Type() NodeType {
	return n.typ
}

func (n *TsLit) Range() span.Range {
	return n.rng
}

func (n *TsLit) OuterParen() span.Range {
	return n.opa
}

func (n *TsLit) SetOuterParen(rng span.Range) {
	n.opa = rng
}

func (n *TsLit) Lit() Node {
	return n.lit
}

type TsThis struct {
	typ NodeType
	rng span.Range
	opa span.Range
}

func (n *TsThis) Type() NodeType {
	return n.typ
}

func (n *TsThis) Range() span.Range {
	return n.rng
}

func (n *TsThis) OuterParen() span.Range {
	return n.opa
}

func (n *TsThis) SetOuterParen(rng span.Range) {
	n.opa = rng
}

// #[visitor(Lhs,Rhs)]
type TsNsName struct {
	typ NodeType
	rng span.Range
	lhs Node
	dot span.Range
	rhs Node
	opa span.Range
}

func (n *TsNsName) Type() NodeType {
	return n.typ
}

func (n *TsNsName) Range() span.Range {
	return n.rng
}

func (n *TsNsName) OuterParen() span.Range {
	return n.opa
}

func (n *TsNsName) SetOuterParen(rng span.Range) {
	n.opa = rng
}

func (n *TsNsName) Lhs() Node {
	return n.lhs
}

func (n *TsNsName) Rhs() Node {
	return n.rhs
}

// #[visitor(Name,Args)]
type TsRef struct {
	typ  NodeType
	rng  span.Range
	name Node
	lt   span.Range
	args Node
	opa  span.Range
}

func (n *TsRef) Type() NodeType {
	return n.typ
}

func (n *TsRef) Range() span.Range {
	return n.rng
}

func (n *TsRef) OuterParen() span.Range {
	return n.opa
}

func (n *TsRef) SetOuterParen(rng span.Range) {
	n.opa = rng
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

// #[visitor(Arg)]
type TsTypQuery struct {
	typ NodeType
	rng span.Range
	arg Node // any ts Typ
	opa span.Range
}

func (n *TsTypQuery) Type() NodeType {
	return n.typ
}

func (n *TsTypQuery) Range() span.Range {
	return n.rng
}

func (n *TsTypQuery) OuterParen() span.Range {
	return n.opa
}

func (n *TsTypQuery) SetOuterParen(rng span.Range) {
	n.opa = rng
}

func (n *TsTypQuery) Arg() Node {
	return n.arg
}

// #[visitor(Arg)]
type TsParen struct {
	typ NodeType
	rng span.Range
	arg Node // name or nsName
	opa span.Range
}

func (n *TsParen) Type() NodeType {
	return n.typ
}

func (n *TsParen) Range() span.Range {
	return n.rng
}

func (n *TsParen) OuterParen() span.Range {
	return n.opa
}

func (n *TsParen) SetOuterParen(rng span.Range) {
	n.opa = rng
}

func (n *TsParen) Arg() Node {
	return n.arg
}

// #[visitor(Arg)]
type TsArr struct {
	typ     NodeType
	rng     span.Range
	bracket span.Range
	arg     Node // name or nsName
	opa     span.Range
}

func (n *TsArr) Type() NodeType {
	return n.typ
}

func (n *TsArr) Range() span.Range {
	return n.rng
}

func (n *TsArr) OuterParen() span.Range {
	return n.opa
}

func (n *TsArr) SetOuterParen(rng span.Range) {
	n.opa = rng
}

func (n *TsArr) Arg() Node {
	return n.arg
}

// #[visitor(Obj,Idx)]
type TsIdxAccess struct {
	typ NodeType
	rng span.Range
	obj Node
	idx Node
	opa span.Range
}

func (n *TsIdxAccess) Type() NodeType {
	return n.typ
}

func (n *TsIdxAccess) Range() span.Range {
	return n.rng
}

func (n *TsIdxAccess) OuterParen() span.Range {
	return n.opa
}

func (n *TsIdxAccess) SetOuterParen(rng span.Range) {
	n.opa = rng
}

func (n *TsIdxAccess) Obj() Node {
	return n.obj
}

func (n *TsIdxAccess) Idx() Node {
	return n.idx
}

// #[visitor(Args)]
type TsTuple struct {
	typ  NodeType
	rng  span.Range
	args []Node
	opa  span.Range
}

func (n *TsTuple) Type() NodeType {
	return n.typ
}

func (n *TsTuple) Range() span.Range {
	return n.rng
}

func (n *TsTuple) OuterParen() span.Range {
	return n.opa
}

func (n *TsTuple) SetOuterParen(rng span.Range) {
	n.opa = rng
}

func (n *TsTuple) Args() []Node {
	return n.args
}

// #[visitor(Arg)]
type TsRest struct {
	typ NodeType
	rng span.Range
	arg Node
}

func (n *TsRest) Type() NodeType {
	return n.typ
}

func (n *TsRest) Range() span.Range {
	return n.rng
}

func (n *TsRest) Arg() Node {
	return n.arg
}

// #[visitor(Label,Val)]
type TsTupleNamedMember struct {
	typ   NodeType
	rng   span.Range
	label Node
	opt   bool
	val   Node
}

func (n *TsTupleNamedMember) Type() NodeType {
	return n.typ
}

func (n *TsTupleNamedMember) Range() span.Range {
	return n.rng
}

func (n *TsTupleNamedMember) Label() Node {
	return n.label
}

func (n *TsTupleNamedMember) Opt() bool {
	return n.opt
}

func (n *TsTupleNamedMember) Val() Node {
	return n.val
}

// #[visitor(Props)]
type TsObj struct {
	typ   NodeType
	rng   span.Range
	props []Node
	opa   span.Range
}

func (n *TsObj) Type() NodeType {
	return n.typ
}

func (n *TsObj) Range() span.Range {
	return n.rng
}

func (n *TsObj) OuterParen() span.Range {
	return n.opa
}

func (n *TsObj) SetOuterParen(rng span.Range) {
	n.opa = rng
}

func (n *TsObj) Props() []Node {
	return n.props
}

// #[visitor(Key,Val)]
type TsProp struct {
	typ      NodeType
	rng      span.Range
	key      Node
	val      Node
	ques     span.Range
	kind     PropKind
	compute  span.Range
	readonly bool
}

func (n *TsProp) Type() NodeType {
	return n.typ
}

func (n *TsProp) Range() span.Range {
	return n.rng
}

func (n *TsProp) Kind() PropKind {
	return n.kind
}

func (n *TsProp) Key() Node {
	return n.key
}

func (n *TsProp) Val() Node {
	return n.val
}

func (n *TsProp) Optional() bool {
	return !n.ques.Empty()
}

func (n *TsProp) Readonly() bool {
	if wt, ok := n.key.(NodeWithTypInfo); ok {
		return wt.TypInfo().Readonly()
	}
	return false
}

func (n *TsProp) Computed() bool {
	return !n.compute.Empty()
}

func (n *TsProp) Method() *TsCallSig {
	if n.val == nil {
		return nil
	}
	return n.val.(*TsCallSig)
}

func (n *TsProp) IsMethod() bool {
	if n.val == nil {
		return false
	}
	return n.val.Type() == N_TS_CALL_SIG
}

// #[visitor(TypParams,Params,RetTyp)]
type TsCallSig struct {
	typ       NodeType
	rng       span.Range
	typParams Node
	params    []Node
	retTyp    Node
}

func (n *TsCallSig) Type() NodeType {
	return n.typ
}

func (n *TsCallSig) Range() span.Range {
	return n.rng
}

func (n *TsCallSig) TypParams() Node {
	return n.typParams
}

func (n *TsCallSig) Params() []Node {
	return n.params
}

func (n *TsCallSig) RetTyp() Node {
	return n.retTyp
}

// #[visitor(TypParams,Params,RetTyp)]
type TsNewSig struct {
	typ       NodeType
	rng       span.Range
	typParams Node
	params    []Node
	retTyp    Node
	abstract  bool
}

func (n *TsNewSig) Type() NodeType {
	return n.typ
}

func (n *TsNewSig) Range() span.Range {
	return n.rng
}

func (n *TsNewSig) Abstract() bool {
	return n.abstract
}

func (n *TsNewSig) TypParams() Node {
	return n.typParams
}

func (n *TsNewSig) Params() []Node {
	return n.params
}

func (n *TsNewSig) RetTyp() Node {
	return n.retTyp
}

// #[visitor(Key,KeyType,Val)]
type TsIdxSig struct {
	typ  NodeType
	rng  span.Range
	key  Node
	val  Node
	ques span.Range
}

func (n *TsIdxSig) Type() NodeType {
	return n.typ
}

func (n *TsIdxSig) Range() span.Range {
	return n.rng
}

func (n *TsIdxSig) Key() Node {
	return n.key
}

func (n *TsIdxSig) KeyType() Node {
	return n.key.(NodeWithTypInfo).TypInfo().typAnnot
}

func (n *TsIdxSig) Optional() bool {
	return !n.ques.Empty()
}

func (n *TsIdxSig) Val() Node {
	return n.val
}

type TsRoughParam struct {
	typ   NodeType
	rng   span.Range
	name  Node
	colon span.Range
	ti    *TypInfo
}

func (n *TsRoughParam) Type() NodeType {
	return n.typ
}

func (n *TsRoughParam) Range() span.Range {
	return n.rng
}

// #[visitor(Params)]
type TsParamsInst struct {
	typ    NodeType
	rng    span.Range
	params []Node
}

func (n *TsParamsInst) Type() NodeType {
	return n.typ
}

func (n *TsParamsInst) Range() span.Range {
	return n.rng
}

func (n *TsParamsInst) Params() []Node {
	return n.params
}

// #[visitor(Params)]
type TsParamsDec struct {
	typ    NodeType
	rng    span.Range
	params []Node
}

func (n *TsParamsDec) Type() NodeType {
	return n.typ
}

func (n *TsParamsDec) Range() span.Range {
	return n.rng
}

func (n *TsParamsDec) Params() []Node {
	return n.params
}

// #[visitor(Name,Cons,Default)]
type TsParam struct {
	typ  NodeType
	rng  span.Range
	name Node
	cons Node // the constraint
	val  Node // the default
}

func (n *TsParam) Type() NodeType {
	return n.typ
}

func (n *TsParam) Range() span.Range {
	return n.rng
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

// #[visitor(TypParams,Params,RetTyp)]
type TsFnTyp struct {
	typ       NodeType
	rng       span.Range
	typParams Node
	params    []Node
	retTyp    Node
	opa       span.Range
}

func (n *TsFnTyp) Type() NodeType {
	return n.typ
}

func (n *TsFnTyp) Range() span.Range {
	return n.rng
}

func (n *TsFnTyp) OuterParen() span.Range {
	return n.opa
}

func (n *TsFnTyp) SetOuterParen(rng span.Range) {
	n.opa = rng
}

func (n *TsFnTyp) TypParams() Node {
	return n.typParams
}

func (n *TsFnTyp) Params() []Node {
	return n.params
}

func (n *TsFnTyp) RetTyp() Node {
	return n.retTyp
}

// #[visitor(Elems)]
type TsUnionTyp struct {
	typ   NodeType
	rng   span.Range
	op    span.Range
	elems []Node
	opa   span.Range
}

func (n *TsUnionTyp) Type() NodeType {
	return n.typ
}

func (n *TsUnionTyp) Range() span.Range {
	return n.rng
}

func (n *TsUnionTyp) OuterParen() span.Range {
	return n.opa
}

func (n *TsUnionTyp) SetOuterParen(rng span.Range) {
	n.opa = rng
}

func (n *TsUnionTyp) Elems() []Node {
	return n.elems
}

// #[visitor(Elems)]
type TsIntersectTyp struct {
	typ   NodeType
	rng   span.Range
	op    span.Range
	elems []Node
	opa   span.Range
}

func (n *TsIntersectTyp) Type() NodeType {
	return n.typ
}

func (n *TsIntersectTyp) Range() span.Range {
	return n.rng
}

func (n *TsIntersectTyp) OuterParen() span.Range {
	return n.opa
}

func (n *TsIntersectTyp) SetOuterParen(rng span.Range) {
	n.opa = rng
}

func (n *TsIntersectTyp) Elems() []Node {
	return n.elems
}

// #[visitor(Typ,Expr)]
type TsTypAssert struct {
	typ NodeType
	rng span.Range
	des Node
	arg Node
	opa span.Range
}

func (n *TsTypAssert) Type() NodeType {
	return n.typ
}

func (n *TsTypAssert) Range() span.Range {
	return n.rng
}

func (n *TsTypAssert) OuterParen() span.Range {
	return n.opa
}

func (n *TsTypAssert) SetOuterParen(rng span.Range) {
	n.opa = rng
}

func (n *TsTypAssert) Typ() Node {
	return n.des
}

func (n *TsTypAssert) Expr() Node {
	return n.arg
}

// #[visitor(Id,TypParams)]
type TsTypDec struct {
	typ  NodeType
	rng  span.Range
	name Node
	ti   *TypInfo
}

func (n *TsTypDec) Type() NodeType {
	return n.typ
}

func (n *TsTypDec) Range() span.Range {
	return n.rng
}

func (n *TsTypDec) Id() Node {
	return n.name
}

func (n *TsTypDec) TypParams() Node {
	return n.ti.typParams
}

func (n *TsTypDec) TypInfo() *TypInfo {
	return n.ti
}

// #[visitor(Id,TypParams,Supers,Body)]
type TsInterface struct {
	typ    NodeType
	rng    span.Range
	name   Node
	params Node
	supers []Node
	body   Node
}

func (n *TsInterface) Type() NodeType {
	return n.typ
}

func (n *TsInterface) Range() span.Range {
	return n.rng
}

func (n *TsInterface) Id() Node {
	return n.name
}

func (n *TsInterface) TypParams() Node {
	return n.params
}

func (n *TsInterface) Supers() []Node {
	return n.supers
}

func (n *TsInterface) Body() Node {
	return n.body
}

// #[visitor(Body)]
type TsInterfaceBody struct {
	typ  NodeType
	rng  span.Range
	body []Node
}

func (n *TsInterfaceBody) Type() NodeType {
	return n.typ
}

func (n *TsInterfaceBody) Range() span.Range {
	return n.rng
}

func (n *TsInterfaceBody) Body() []Node {
	return n.body
}

// #[visitor(Id,Members)]
type TsEnum struct {
	typ   NodeType
	rng   span.Range
	name  Node
	items []Node
	cons  bool
}

func (n *TsEnum) Type() NodeType {
	return n.typ
}

func (n *TsEnum) Range() span.Range {
	return n.rng
}

func (n *TsEnum) Const() bool {
	return n.cons
}

func (n *TsEnum) Id() Node {
	return n.name
}

func (n *TsEnum) Members() []Node {
	return n.items
}

// #[visitor(Key,Val)]
type TsEnumMember struct {
	typ NodeType
	rng span.Range
	key Node
	val Node
}

func (n *TsEnumMember) Type() NodeType {
	return n.typ
}

func (n *TsEnumMember) Range() span.Range {
	return n.rng
}

func (n *TsEnumMember) Key() Node {
	return n.key
}

func (n *TsEnumMember) Val() Node {
	return n.val
}

// #[visitor(Name,Val)]
type TsImportAlias struct {
	typ    NodeType
	rng    span.Range
	name   Node
	val    Node
	export bool
}

func (n *TsImportAlias) Type() NodeType {
	return n.typ
}

func (n *TsImportAlias) Range() span.Range {
	return n.rng
}

func (n *TsImportAlias) Name() Node {
	return n.name
}

func (n *TsImportAlias) Val() Node {
	return n.val
}

func (n *TsImportAlias) Export() bool {
	return n.export
}

// #[visitor(Id,Body)]
type TsNS struct {
	typ   NodeType
	rng   span.Range
	name  Node
	body  Node
	alias bool
}

func (n *TsNS) Type() NodeType {
	return n.typ
}

func (n *TsNS) Range() span.Range {
	return n.rng
}

func (n *TsNS) Id() Node {
	return n.name
}

func (n *TsNS) Body() Node {
	return n.body
}

func (n *TsNS) Alias() bool {
	return n.alias
}

// #[visitor(Name,Expr)]
type TsImportRequire struct {
	typ  NodeType
	rng  span.Range
	name Node
	expr Node
	opa  span.Range
}

func (n *TsImportRequire) Type() NodeType {
	return n.typ
}

func (n *TsImportRequire) Range() span.Range {
	return n.rng
}

func (n *TsImportRequire) OuterParen() span.Range {
	return n.opa
}

func (n *TsImportRequire) SetOuterParen(rng span.Range) {
	n.opa = rng
}

func (n *TsImportRequire) Name() Node {
	return n.name
}

func (n *TsImportRequire) Expr() Node {
	return n.expr
}

// #[visitor(Expr)]
type TsExportAssign struct {
	typ  NodeType
	rng  span.Range
	expr Node
}

func (n *TsExportAssign) Type() NodeType {
	return n.typ
}

func (n *TsExportAssign) Range() span.Range {
	return n.rng
}

func (n *TsExportAssign) Expr() Node {
	return n.expr
}

// #[visitor(Name,Inner)]
type TsDec struct {
	typ   NodeType
	rng   span.Range
	name  Node
	inner Node
}

func (n *TsDec) Type() NodeType {
	return n.typ
}

func (n *TsDec) Range() span.Range {
	return n.rng
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
//
// #[visitor(Name,Typ)]
type TsTypPredicate struct {
	typ    NodeType
	rng    span.Range
	name   Node
	des    Node
	assert bool
	opa    span.Range
}

func (n *TsTypPredicate) Type() NodeType {
	return n.typ
}

func (n *TsTypPredicate) Range() span.Range {
	return n.rng
}

func (n *TsTypPredicate) OuterParen() span.Range {
	return n.opa
}

func (n *TsTypPredicate) SetOuterParen(rng span.Range) {
	n.opa = rng
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

// #[visitor(Arg)]
type TsNoNull struct {
	typ NodeType
	rng span.Range
	arg Node
	opa span.Range
}

func (n *TsNoNull) Type() NodeType {
	return n.typ
}

func (n *TsNoNull) Range() span.Range {
	return n.rng
}

func (n *TsNoNull) OuterParen() span.Range {
	return n.opa
}

func (n *TsNoNull) SetOuterParen(rng span.Range) {
	n.opa = rng
}

func (n *TsNoNull) Arg() Node {
	return n.arg
}

// #[visitor(Arg,Qualifier,TypArg)]
type TsImportType struct {
	typ       NodeType
	rng       span.Range
	arg       Node
	qualifier Node
	typArgs   Node
	opa       span.Range
}

func (n *TsImportType) Type() NodeType {
	return n.typ
}

func (n *TsImportType) Range() span.Range {
	return n.rng
}

func (n *TsImportType) OuterParen() span.Range {
	return n.opa
}

func (n *TsImportType) SetOuterParen(rng span.Range) {
	n.opa = rng
}

func (n *TsImportType) Arg() Node {
	return n.arg
}

func (n *TsImportType) Qualifier() Node {
	return n.qualifier
}

func (n *TsImportType) TypArg() Node {
	return n.typArgs
}

// ClassDec
func (n *ClassDec) Implements() []Node {
	if util.IsNilPtr(n.ti) {
		return nil
	}
	return n.ti.Implements()
}

func (n *ClassDec) Abstract() bool {
	if util.IsNilPtr(n.ti) {
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

	return st.SuperTypArgs()
}

func (n *ClassDec) TypParams() Node {
	if util.IsNilPtr(n.ti) {
		return nil
	}
	return n.ti.typParams
}

// #[visitor(CheckTyp,ExtTyp,TrueTyp,FalseTyp)]
type TsCondType struct {
	typ      NodeType
	rng      span.Range
	check    Node
	ext      Node
	trueTyp  Node
	falseTyp Node
	opa      span.Range
}

func (n *TsCondType) Type() NodeType {
	return n.typ
}

func (n *TsCondType) Range() span.Range {
	return n.rng
}

func (n *TsCondType) OuterParen() span.Range {
	return n.opa
}

func (n *TsCondType) SetOuterParen(rng span.Range) {
	n.opa = rng
}

func (n *TsCondType) CheckTyp() Node {
	return n.check
}

func (n *TsCondType) ExtTyp() Node {
	return n.ext
}

func (n *TsCondType) TrueTyp() Node {
	return n.trueTyp
}

func (n *TsCondType) FalseTyp() Node {
	return n.falseTyp
}

// #[visitor(Arg)]
type TsTypInfer struct {
	typ NodeType
	rng span.Range
	arg Node
	opa span.Range
}

func (n *TsTypInfer) Type() NodeType {
	return n.typ
}

func (n *TsTypInfer) Range() span.Range {
	return n.rng
}

func (n *TsTypInfer) OuterParen() span.Range {
	return n.opa
}

func (n *TsTypInfer) SetOuterParen(rng span.Range) {
	n.opa = rng
}

func (n *TsTypInfer) Arg() Node {
	return n.arg
}

// #[visitor(Name,Key,Val)]
type TsMapped struct {
	typ      NodeType
	rng      span.Range
	readonly int // 0: not set, 1: set, 2: positive, 3: negative
	optional int // 0: not set, 1: set, 2: positive, 3: negative
	key      Node
	name     Node
	val      Node
	opa      span.Range
}

func (n *TsMapped) Type() NodeType {
	return n.typ
}

func (n *TsMapped) Range() span.Range {
	return n.rng
}

func (n *TsMapped) OuterParen() span.Range {
	return n.opa
}

func (n *TsMapped) SetOuterParen(rng span.Range) {
	n.opa = rng
}

func (n *TsMapped) ReadonlyFmt() interface{} {
	switch n.readonly {
	case 1:
		return true
	case 2:
		return "+"
	case 3:
		return "-"
	}
	return ""
}

func (n *TsMapped) OptionalFmt() interface{} {
	switch n.optional {
	case 1:
		return true
	case 2:
		return "+"
	case 3:
		return "-"
	}
	return ""
}

func (n *TsMapped) Readonly() int {
	return n.readonly
}

func (n *TsMapped) Optional() int {
	return n.optional
}

func (n *TsMapped) Name() Node {
	return n.name
}

func (n *TsMapped) Key() Node {
	return n.key
}

func (n *TsMapped) Val() Node {
	return n.val
}

// #[visitor(Arg)]
type TsTypOp struct {
	typ NodeType
	rng span.Range
	op  string
	arg Node
	opa span.Range
}

func (n *TsTypOp) Type() NodeType {
	return n.typ
}

func (n *TsTypOp) Range() span.Range {
	return n.rng
}

func (n *TsTypOp) OuterParen() span.Range {
	return n.opa
}

func (n *TsTypOp) SetOuterParen(rng span.Range) {
	n.opa = rng
}

func (n *TsTypOp) Op() string {
	return n.op
}

func (n *TsTypOp) Arg() Node {
	return n.arg
}

// #[visitor(Arg)]
type TsOpt struct {
	typ NodeType
	rng span.Range
	arg Node
}

func (n *TsOpt) Type() NodeType {
	return n.typ
}

func (n *TsOpt) Range() span.Range {
	return n.rng
}

func (n *TsOpt) Arg() Node {
	return n.arg
}

// Field
func (n *Field) IsTsSig() bool {
	if n.ti == nil || n.val != nil || !n.computed {
		return false
	}
	if wt, ok := n.key.(NodeWithTypInfo); ok {
		ti := wt.TypInfo()
		return ti != nil && ti.typAnnot != nil
	}
	return false
}
