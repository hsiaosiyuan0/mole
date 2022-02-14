package parser

import "github.com/hsiaosiyuan0/mole/fuzz"

// #[visitor(TsTyp)]
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
	typ        NodeType
	loc        *Loc
	ques       *Loc
	outerParen *Loc
}

func (n *TsPredef) OuterParen() *Loc {
	return n.outerParen
}

func (n *TsPredef) SetOuterParen(loc *Loc) {
	n.outerParen = loc
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

// #[visitor(Lit)]
type TsLit struct {
	typ        NodeType
	loc        *Loc
	lit        Node
	outerParen *Loc
}

func (n *TsLit) OuterParen() *Loc {
	return n.outerParen
}

func (n *TsLit) SetOuterParen(loc *Loc) {
	n.outerParen = loc
}

func (n *TsLit) Lit() Node {
	return n.lit
}

func (n *TsLit) Type() NodeType {
	return n.typ
}

func (n *TsLit) Loc() *Loc {
	return n.loc
}

type TsThis struct {
	typ        NodeType
	loc        *Loc
	outerParen *Loc
}

func (n *TsThis) OuterParen() *Loc {
	return n.outerParen
}

func (n *TsThis) SetOuterParen(loc *Loc) {
	n.outerParen = loc
}

func (n *TsThis) Type() NodeType {
	return n.typ
}

func (n *TsThis) Loc() *Loc {
	return n.loc
}

// #[visitor(Lhs,Rhs)]
type TsNsName struct {
	typ        NodeType
	loc        *Loc
	lhs        Node
	dot        *Loc
	rhs        Node
	outerParen *Loc
}

func (n *TsNsName) OuterParen() *Loc {
	return n.outerParen
}

func (n *TsNsName) SetOuterParen(loc *Loc) {
	n.outerParen = loc
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

// #[visitor(Name,Args)]
type TsRef struct {
	typ        NodeType
	loc        *Loc
	name       Node
	lt         *Loc
	args       Node
	outerParen *Loc
}

func (n *TsRef) OuterParen() *Loc {
	return n.outerParen
}

func (n *TsRef) SetOuterParen(loc *Loc) {
	n.outerParen = loc
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

// #[visitor(Arg)]
type TsTypQuery struct {
	typ        NodeType
	loc        *Loc
	arg        Node // any ts Typ
	outerParen *Loc
}

func (n *TsTypQuery) OuterParen() *Loc {
	return n.outerParen
}

func (n *TsTypQuery) SetOuterParen(loc *Loc) {
	n.outerParen = loc
}

func (n *TsTypQuery) Arg() Node {
	return n.arg
}

func (n *TsTypQuery) Type() NodeType {
	return n.typ
}

func (n *TsTypQuery) Loc() *Loc {
	return n.loc
}

// #[visitor(Arg)]
type TsParen struct {
	typ        NodeType
	loc        *Loc
	arg        Node // name or nsName
	outerParen *Loc
}

func (n *TsParen) OuterParen() *Loc {
	return n.outerParen
}

func (n *TsParen) SetOuterParen(loc *Loc) {
	n.outerParen = loc
}

func (n *TsParen) Arg() Node {
	return n.arg
}

func (n *TsParen) Type() NodeType {
	return n.typ
}

func (n *TsParen) Loc() *Loc {
	return n.loc
}

// #[visitor(Arg)]
type TsArr struct {
	typ        NodeType
	loc        *Loc
	bracket    *Loc
	arg        Node // name or nsName
	outerParen *Loc
}

func (n *TsArr) OuterParen() *Loc {
	return n.outerParen
}

func (n *TsArr) SetOuterParen(loc *Loc) {
	n.outerParen = loc
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

// #[visitor(Obj,Idx)]
type TsIdxAccess struct {
	typ        NodeType
	loc        *Loc
	obj        Node
	idx        Node
	outerParen *Loc
}

func (n *TsIdxAccess) OuterParen() *Loc {
	return n.outerParen
}

func (n *TsIdxAccess) SetOuterParen(loc *Loc) {
	n.outerParen = loc
}

func (n *TsIdxAccess) Obj() Node {
	return n.obj
}

func (n *TsIdxAccess) Idx() Node {
	return n.idx
}

func (n *TsIdxAccess) Type() NodeType {
	return n.typ
}

func (n *TsIdxAccess) Loc() *Loc {
	return n.loc
}

// #[visitor(Args)]
type TsTuple struct {
	typ        NodeType
	loc        *Loc
	args       []Node
	outerParen *Loc
}

func (n *TsTuple) OuterParen() *Loc {
	return n.outerParen
}

func (n *TsTuple) SetOuterParen(loc *Loc) {
	n.outerParen = loc
}

func (n *TsTuple) Args() []Node {
	return n.args
}

func (n *TsTuple) Type() NodeType {
	return n.typ
}

func (n *TsTuple) Loc() *Loc {
	return n.loc
}

// #[visitor(Arg)]
type TsRest struct {
	typ NodeType
	loc *Loc
	arg Node
}

func (n *TsRest) Arg() Node {
	return n.arg
}

func (n *TsRest) Type() NodeType {
	return n.typ
}

func (n *TsRest) Loc() *Loc {
	return n.loc
}

// #[visitor(Label,Val)]
type TsTupleNamedMember struct {
	typ   NodeType
	loc   *Loc
	label Node
	opt   bool
	val   Node
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

func (n *TsTupleNamedMember) Type() NodeType {
	return n.typ
}

func (n *TsTupleNamedMember) Loc() *Loc {
	return n.loc
}

// #[visitor(Props)]
type TsObj struct {
	typ        NodeType
	loc        *Loc
	props      []Node
	outerParen *Loc
}

func (n *TsObj) OuterParen() *Loc {
	return n.outerParen
}

func (n *TsObj) SetOuterParen(loc *Loc) {
	n.outerParen = loc
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

// #[visitor(Key,Val)]
type TsProp struct {
	typ        NodeType
	loc        *Loc
	key        Node
	val        Node
	ques       *Loc
	kind       PropKind
	computeLoc *Loc
	readonly   bool
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
	return n.ques != nil
}

func (n *TsProp) Readonly() bool {
	if wt, ok := n.key.(NodeWithTypInfo); ok {
		return wt.TypInfo().Readonly()
	}
	return false
}

func (n *TsProp) Computed() bool {
	return n.computeLoc != nil
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

func (n *TsProp) Type() NodeType {
	return n.typ
}

func (n *TsProp) Loc() *Loc {
	return n.loc
}

// #[visitor(TypParams,Params,RetTyp)]
type TsCallSig struct {
	typ       NodeType
	loc       *Loc
	typParams Node
	params    []Node
	retTyp    Node
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

func (n *TsCallSig) Type() NodeType {
	return n.typ
}

func (n *TsCallSig) Loc() *Loc {
	return n.loc
}

// #[visitor(TypParams,Params,RetTyp)]
type TsNewSig struct {
	typ       NodeType
	loc       *Loc
	typParams Node
	params    []Node
	retTyp    Node
	abstract  bool
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

func (n *TsNewSig) Type() NodeType {
	return n.typ
}

func (n *TsNewSig) Loc() *Loc {
	return n.loc
}

// #[visitor(Key,KeyType,Val)]
type TsIdxSig struct {
	typ  NodeType
	loc  *Loc
	key  Node
	val  Node
	ques *Loc
}

func (n *TsIdxSig) Key() Node {
	return n.key
}

func (n *TsIdxSig) KeyType() Node {
	return n.key.(NodeWithTypInfo).TypInfo().typAnnot
}

func (n *TsIdxSig) Optional() bool {
	return n.ques != nil
}

func (n *TsIdxSig) Val() Node {
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

// #[visitor(Params)]
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

// #[visitor(Params)]
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

// #[visitor(Name,Cons,Default)]
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

// #[visitor(TypParams,Params,RetTyp)]
type TsFnTyp struct {
	typ        NodeType
	loc        *Loc
	typParams  Node
	params     []Node
	retTyp     Node
	outerParen *Loc
}

func (n *TsFnTyp) OuterParen() *Loc {
	return n.outerParen
}

func (n *TsFnTyp) SetOuterParen(loc *Loc) {
	n.outerParen = loc
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

func (n *TsFnTyp) Type() NodeType {
	return n.typ
}

func (n *TsFnTyp) Loc() *Loc {
	return n.loc
}

// #[visitor(Elems)]
type TsUnionTyp struct {
	typ        NodeType
	loc        *Loc
	op         *Loc
	elems      []Node
	outerParen *Loc
}

func (n *TsUnionTyp) OuterParen() *Loc {
	return n.outerParen
}

func (n *TsUnionTyp) SetOuterParen(loc *Loc) {
	n.outerParen = loc
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

// #[visitor(Elems)]
type TsIntersecTyp struct {
	typ        NodeType
	loc        *Loc
	op         *Loc
	elems      []Node
	outerParen *Loc
}

func (n *TsIntersecTyp) OuterParen() *Loc {
	return n.outerParen
}

func (n *TsIntersecTyp) SetOuterParen(loc *Loc) {
	n.outerParen = loc
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

// #[visitor(Typ,Expr)]
type TsTypAssert struct {
	typ        NodeType
	loc        *Loc
	des        Node
	arg        Node
	outerParen *Loc
}

func (n *TsTypAssert) OuterParen() *Loc {
	return n.outerParen
}

func (n *TsTypAssert) SetOuterParen(loc *Loc) {
	n.outerParen = loc
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

// #[visitor(Id,TypParams)]
type TsTypDec struct {
	typ  NodeType
	loc  *Loc
	name Node
	ti   *TypInfo
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

func (n *TsTypDec) Type() NodeType {
	return n.typ
}

func (n *TsTypDec) Loc() *Loc {
	return n.loc
}

// #[visitor(Id,TypParams,Supers,Body)]
type TsInferface struct {
	typ    NodeType
	loc    *Loc
	name   Node
	params Node
	supers []Node
	body   Node
}

func (n *TsInferface) Id() Node {
	return n.name
}

func (n *TsInferface) TypParams() Node {
	return n.params
}

func (n *TsInferface) Supers() []Node {
	return n.supers
}

func (n *TsInferface) Body() Node {
	return n.body
}

func (n *TsInferface) Type() NodeType {
	return n.typ
}

func (n *TsInferface) Loc() *Loc {
	return n.loc
}

// #[visitor(Body)]
type TsInferfaceBody struct {
	typ  NodeType
	loc  *Loc
	body []Node
}

func (n *TsInferfaceBody) Body() []Node {
	return n.body
}

func (n *TsInferfaceBody) Type() NodeType {
	return n.typ
}

func (n *TsInferfaceBody) Loc() *Loc {
	return n.loc
}

// #[visitor(Id,Members)]
type TsEnum struct {
	typ   NodeType
	loc   *Loc
	name  Node
	items []Node
	cons  bool
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

func (n *TsEnum) Type() NodeType {
	return n.typ
}

func (n *TsEnum) Loc() *Loc {
	return n.loc
}

// #[visitor(Key,Val)]
type TsEnumMember struct {
	typ NodeType
	loc *Loc
	key Node
	val Node
}

func (n *TsEnumMember) Key() Node {
	return n.key
}

func (n *TsEnumMember) Val() Node {
	return n.val
}

func (n *TsEnumMember) Type() NodeType {
	return n.typ
}

func (n *TsEnumMember) Loc() *Loc {
	return n.loc
}

// #[visitor(Name,Val)]
type TsImportAlias struct {
	typ    NodeType
	loc    *Loc
	name   Node
	val    Node
	export bool
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

func (n *TsImportAlias) Type() NodeType {
	return n.typ
}

func (n *TsImportAlias) Loc() *Loc {
	return n.loc
}

// #[visitor(Id,Body)]
type TsNS struct {
	typ   NodeType
	loc   *Loc
	name  Node
	body  Node
	alias bool
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

func (n *TsNS) Type() NodeType {
	return n.typ
}

func (n *TsNS) Loc() *Loc {
	return n.loc
}

// #[visitor(Name,Expr)]
type TsImportRequire struct {
	typ        NodeType
	loc        *Loc
	name       Node
	expr       Node
	outerParen *Loc
}

func (n *TsImportRequire) OuterParen() *Loc {
	return n.outerParen
}

func (n *TsImportRequire) SetOuterParen(loc *Loc) {
	n.outerParen = loc
}

func (n *TsImportRequire) Name() Node {
	return n.name
}

func (n *TsImportRequire) Expr() Node {
	return n.expr
}

func (n *TsImportRequire) Type() NodeType {
	return n.typ
}

func (n *TsImportRequire) Loc() *Loc {
	return n.loc
}

// #[visitor(Expr)]
type TsExportAssign struct {
	typ  NodeType
	loc  *Loc
	expr Node
}

func (n *TsExportAssign) Expr() Node {
	return n.expr
}

func (n *TsExportAssign) Type() NodeType {
	return n.typ
}

func (n *TsExportAssign) Loc() *Loc {
	return n.loc
}

// #[visitor(Name,Inner)]
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
//
// #[visitor(Name,Typ)]
type TsTypPredicate struct {
	typ        NodeType
	loc        *Loc
	name       Node
	des        Node
	assert     bool
	outerParen *Loc
}

func (n *TsTypPredicate) OuterParen() *Loc {
	return n.outerParen
}

func (n *TsTypPredicate) SetOuterParen(loc *Loc) {
	n.outerParen = loc
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

// #[visitor(Arg)]
type TsNoNull struct {
	typ        NodeType
	loc        *Loc
	arg        Node
	outerParen *Loc
}

func (n *TsNoNull) OuterParen() *Loc {
	return n.outerParen
}

func (n *TsNoNull) SetOuterParen(loc *Loc) {
	n.outerParen = loc
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

// #[visitor(Arg,Qualifier,TypArg)]
type TsImportType struct {
	typ        NodeType
	loc        *Loc
	arg        Node
	qualifier  Node
	typArgs    Node
	outerParen *Loc
}

func (n *TsImportType) OuterParen() *Loc {
	return n.outerParen
}

func (n *TsImportType) SetOuterParen(loc *Loc) {
	n.outerParen = loc
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

func (n *TsImportType) Type() NodeType {
	return n.typ
}

func (n *TsImportType) Loc() *Loc {
	return n.loc
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

	return st.SuperTypArgs()
}

func (n *ClassDec) TypParams() Node {
	if fuzz.IsNilPtr(n.ti) {
		return nil
	}
	return n.ti.typParams
}

// #[visitor(CheckTyp,ExtTyp,TrueTyp,FalseTyp)]
type TsCondType struct {
	typ        NodeType
	loc        *Loc
	check      Node
	ext        Node
	trueTyp    Node
	falseTyp   Node
	outerParen *Loc
}

func (n *TsCondType) OuterParen() *Loc {
	return n.outerParen
}

func (n *TsCondType) SetOuterParen(loc *Loc) {
	n.outerParen = loc
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

func (n *TsCondType) Type() NodeType {
	return n.typ
}

func (n *TsCondType) Loc() *Loc {
	return n.loc
}

// #[visitor(Arg)]
type TsTypInfer struct {
	typ        NodeType
	loc        *Loc
	arg        Node
	outerParen *Loc
}

func (n *TsTypInfer) OuterParen() *Loc {
	return n.outerParen
}

func (n *TsTypInfer) SetOuterParen(loc *Loc) {
	n.outerParen = loc
}

func (n *TsTypInfer) Arg() Node {
	return n.arg
}

func (n *TsTypInfer) Type() NodeType {
	return n.typ
}

func (n *TsTypInfer) Loc() *Loc {
	return n.loc
}

// #[visitor(Name,Key,Val)]
type TsMapped struct {
	typ        NodeType
	loc        *Loc
	readonly   int // 0: not set, 1: set, 2: positive, 3: negative
	optional   int // 0: not set, 1: set, 2: positive, 3: negative
	key        Node
	name       Node
	val        Node
	outerParen *Loc
}

func (n *TsMapped) OuterParen() *Loc {
	return n.outerParen
}

func (n *TsMapped) SetOuterParen(loc *Loc) {
	n.outerParen = loc
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

func (n *TsMapped) Type() NodeType {
	return n.typ
}

func (n *TsMapped) Loc() *Loc {
	return n.loc
}

// #[visitor(Arg)]
type TsTypOp struct {
	typ        NodeType
	loc        *Loc
	op         string
	arg        Node
	outerParen *Loc
}

func (n *TsTypOp) OuterParen() *Loc {
	return n.outerParen
}

func (n *TsTypOp) SetOuterParen(loc *Loc) {
	n.outerParen = loc
}

func (n *TsTypOp) Op() string {
	return n.op
}

func (n *TsTypOp) Arg() Node {
	return n.arg
}

func (n *TsTypOp) Type() NodeType {
	return n.typ
}

func (n *TsTypOp) Loc() *Loc {
	return n.loc
}

// #[visitor(Arg)]
type TsOpt struct {
	typ NodeType
	loc *Loc
	arg Node
}

func (n *TsOpt) Arg() Node {
	return n.arg
}

func (n *TsOpt) Type() NodeType {
	return n.typ
}

func (n *TsOpt) Loc() *Loc {
	return n.loc
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
