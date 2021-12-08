package parser

type BuiltinType struct {
	typ NodeType
	loc *Loc
}

func (n *BuiltinType) Type() NodeType {
	return n.typ
}

func (n *BuiltinType) Loc() *Loc {
	return n.loc
}

func (n *BuiltinType) Extra() interface{} {
	return nil
}

func (n *BuiltinType) setExtra(ext interface{}) {

}
