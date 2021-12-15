package parser

type TsPredef struct {
	typ NodeType
	loc *Loc
}

func (n *TsPredef) Type() NodeType {
	return n.typ
}

func (n *TsPredef) Loc() *Loc {
	return n.loc
}

func (n *TsPredef) Extra() interface{} {
	return nil
}

func (n *TsPredef) setExtra(ext interface{}) {
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

func (n *TsThis) Extra() interface{} {
	return nil
}

func (n *TsThis) setExtra(ext interface{}) {
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

func (n *TsNsName) Extra() interface{} {
	return nil
}

func (n *TsNsName) setExtra(ext interface{}) {
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

func (n *TsRef) Extra() interface{} {
	return nil
}

func (n *TsRef) setExtra(ext interface{}) {
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

func (n *TsQuery) Extra() interface{} {
	return nil
}

func (n *TsQuery) setExtra(ext interface{}) {
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

func (n *TsParen) Extra() interface{} {
	return nil
}

func (n *TsParen) setExtra(ext interface{}) {
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

func (n *TsArr) Extra() interface{} {
	return nil
}

func (n *TsArr) setExtra(ext interface{}) {
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

func (n *TsTuple) Extra() interface{} {
	return nil
}

func (n *TsTuple) setExtra(ext interface{}) {
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

func (n *TsObj) Extra() interface{} {
	return nil
}

func (n *TsObj) setExtra(ext interface{}) {
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

func (n *TsProp) Extra() interface{} {
	return nil
}

func (n *TsProp) setExtra(ext interface{}) {
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

func (n *TsCallSig) Extra() interface{} {
	return nil
}

func (n *TsCallSig) setExtra(ext interface{}) {
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

func (n *TsNewSig) Extra() interface{} {
	return nil
}

func (n *TsNewSig) setExtra(ext interface{}) {
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

func (n *TsIdxSig) Extra() interface{} {
	return nil
}

func (n *TsIdxSig) setExtra(ext interface{}) {
}

type TsRoughParam struct {
	typ      NodeType
	loc      *Loc
	name     Node
	colon    *Loc
	typAnnot Node
}

func (n *TsRoughParam) Type() NodeType {
	return n.typ
}

func (n *TsRoughParam) Loc() *Loc {
	return n.loc
}

func (n *TsRoughParam) Extra() interface{} {
	return nil
}

func (n *TsRoughParam) setExtra(ext interface{}) {
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

func (n *TsParam) Extra() interface{} {
	return nil
}

func (n *TsParam) setExtra(ext interface{}) {
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

func (n *TsFnTyp) Extra() interface{} {
	return nil
}

func (n *TsFnTyp) setExtra(ext interface{}) {
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

func (n *TsUnionTyp) Extra() interface{} {
	return nil
}

func (n *TsUnionTyp) setExtra(ext interface{}) {
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

func (n *TsIntersecTyp) Extra() interface{} {
	return nil
}

func (n *TsIntersecTyp) setExtra(ext interface{}) {
}
