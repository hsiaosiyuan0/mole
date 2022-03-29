package parser

import (
	"reflect"

	"github.com/hsiaosiyuan0/mole/util"
)

func DecoratorsOf(node Node) []Node {
	if node == nil {
		return nil
	}
	if wt, ok := node.(NodeWithTypInfo); ok {
		ti := wt.TypInfo()
		if ti != nil {
			return ti.decorators
		}
	}
	return nil
}

type TypInfo struct {
	ques       *Loc
	not        *Loc
	typAnnot   *TsTypAnnot
	typParams  Node
	typArgs    Node
	decorators []Node
	clsTyp     *ClsTypInfo
}

func NewTypInfo() *TypInfo {
	return &TypInfo{}
}

func (ti *TypInfo) Clone() *TypInfo {
	return &TypInfo{
		ques:     ti.ques,
		not:      ti.not,
		typAnnot: ti.typAnnot,
		typArgs:  ti.typArgs,
		clsTyp:   ti.clsTyp,
	}
}

func (n *TypInfo) Decorators() []Node {
	return n.decorators
}

func (ti *TypInfo) Ques() *Loc {
	return ti.ques
}

func (ti *TypInfo) SetQues(loc *Loc) {
	ti.ques = loc
}

func (ti *TypInfo) Not(loc *Loc) {
	ti.ques = loc
}

func (ti *TypInfo) SetNot(loc *Loc) {
	ti.not = loc
}

func (ti *TypInfo) TypAnnot() *TsTypAnnot {
	return ti.typAnnot
}

func (ti *TypInfo) SetTypAnnot(node Node) {
	if node == nil || reflect.ValueOf(node).IsNil() {
		ti.typAnnot = nil
		return
	}
	if node.Type() != N_TS_TYP_ANNOT {
		node = NewTsTypAnnot(node)
	}
	ti.typAnnot = node.(*TsTypAnnot)
}

func (ti *TypInfo) TypParams() Node {
	return ti.typParams
}

func (ti *TypInfo) SetTypParams(node Node) {
	ti.typParams = node
}

func (ti *TypInfo) TypArgs() Node {
	return ti.typArgs
}

func (ti *TypInfo) SetTypArgs(node Node) {
	ti.typArgs = node
}

func (ti *TypInfo) Optional() bool {
	return ti.ques != nil
}

func (ti *TypInfo) Definite() bool {
	return ti.not != nil
}

type ClsTypInfo struct {
	accMod       ACC_MOD
	superTypArgs Node
	implements   []Node
	beginLoc     *Loc
	abstract     bool
	readonly     bool
	override     bool
	declare      bool
}

func (ti *TypInfo) intiClsTyp() {
	if ti.clsTyp != nil {
		return
	}
	ti.clsTyp = &ClsTypInfo{}
}

func (ti *TypInfo) AccMod() ACC_MOD {
	if util.IsNilPtr(ti.clsTyp) {
		return ACC_MOD_NONE
	}
	return ti.clsTyp.accMod
}

func (ti *TypInfo) SetAccMod(accMod ACC_MOD) {
	ti.intiClsTyp()
	ti.clsTyp.accMod = accMod
}

func (ti *TypInfo) SuperTypArgs() Node {
	if util.IsNilPtr(ti.clsTyp) {
		return nil
	}
	return ti.clsTyp.superTypArgs
}

func (ti *TypInfo) SetSuperTypArgs(node Node) {
	ti.intiClsTyp()
	ti.clsTyp.superTypArgs = node
}

func (ti *TypInfo) BeginLoc() *Loc {
	if util.IsNilPtr(ti.clsTyp) {
		return nil
	}
	return ti.clsTyp.beginLoc
}

func (ti *TypInfo) SetBeginLoc(loc *Loc) {
	ti.intiClsTyp()
	ti.clsTyp.beginLoc = loc
}

func (ti *TypInfo) Abstract() bool {
	if util.IsNilPtr(ti.clsTyp) {
		return false
	}
	return ti.clsTyp.abstract
}

func (ti *TypInfo) SetAbstract(flag bool) {
	ti.intiClsTyp()
	ti.clsTyp.abstract = flag
}

func (ti *TypInfo) Readonly() bool {
	if util.IsNilPtr(ti.clsTyp) {
		return false
	}
	return ti.clsTyp.readonly
}

func (ti *TypInfo) SetReadonly(flag bool) {
	ti.intiClsTyp()
	ti.clsTyp.readonly = flag
}

func (ti *TypInfo) Override() bool {
	if util.IsNilPtr(ti.clsTyp) {
		return false
	}
	return ti.clsTyp.override
}

func (ti *TypInfo) SetOverride(flag bool) {
	ti.intiClsTyp()
	ti.clsTyp.override = flag
}

func (ti *TypInfo) Declare() bool {
	if util.IsNilPtr(ti.clsTyp) {
		return false
	}
	return ti.clsTyp.declare
}

func (ti *TypInfo) SetDeclare(flag bool) {
	ti.intiClsTyp()
	ti.clsTyp.declare = flag
}

func (ti *TypInfo) Implements() []Node {
	if util.IsNilPtr(ti.clsTyp) {
		return nil
	}
	return ti.clsTyp.implements
}

func (ti *TypInfo) SetImplements(nodes []Node) {
	ti.intiClsTyp()
	ti.clsTyp.implements = nodes
}
