package parser

type TsPredefType struct {
	typ NodeType
	loc *Loc
}

func (n *TsPredefType) Type() NodeType {
	return n.typ
}

func (n *TsPredefType) Loc() *Loc {
	return n.loc
}

func (n *TsPredefType) Extra() interface{} {
	return nil
}

func (n *TsPredefType) setExtra(ext interface{}) {
}

type TsThisType struct {
	typ NodeType
	loc *Loc
}

func (n *TsThisType) Type() NodeType {
	return n.typ
}

func (n *TsThisType) Loc() *Loc {
	return n.loc
}

func (n *TsThisType) Extra() interface{} {
	return nil
}

func (n *TsThisType) setExtra(ext interface{}) {
}

type TsNsName struct {
	typ NodeType
	loc *Loc
	lhs Node
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

type TsTypeRef struct {
	typ  NodeType
	loc  *Loc
	name Node
	args []Node
}

func (n *TsTypeRef) Type() NodeType {
	return n.typ
}

func (n *TsTypeRef) Loc() *Loc {
	return n.loc
}

func (n *TsTypeRef) Extra() interface{} {
	return nil
}

func (n *TsTypeRef) setExtra(ext interface{}) {
}

type TsTypeQuery struct {
	typ NodeType
	loc *Loc
	arg Node // name or nsName
}

func (n *TsTypeQuery) Type() NodeType {
	return n.typ
}

func (n *TsTypeQuery) Loc() *Loc {
	return n.loc
}

func (n *TsTypeQuery) Extra() interface{} {
	return nil
}

func (n *TsTypeQuery) setExtra(ext interface{}) {
}

type TsArrType struct {
	typ NodeType
	loc *Loc
	arg Node // name or nsName
}

func (n *TsArrType) Type() NodeType {
	return n.typ
}

func (n *TsArrType) Loc() *Loc {
	return n.loc
}

func (n *TsArrType) Extra() interface{} {
	return nil
}

func (n *TsArrType) setExtra(ext interface{}) {
}

type TsTupleType struct {
	typ  NodeType
	loc  *Loc
	args []Node
}

func (n *TsTupleType) Type() NodeType {
	return n.typ
}

func (n *TsTupleType) Loc() *Loc {
	return n.loc
}

func (n *TsTupleType) Extra() interface{} {
	return nil
}

func (n *TsTupleType) setExtra(ext interface{}) {
}
