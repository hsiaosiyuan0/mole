package parser

import (
	"reflect"

	"github.com/hsiaosiyuan0/mole/fuzz"
)

type TypInfo struct {
	ques      *Loc
	typAnnot  *TsTypAnnot
	typParams Node
	typArgs   Node
	clsTyp    *ClsTypInfo
}

func NewTypInfo() *TypInfo {
	return &TypInfo{}
}

func (ti *TypInfo) Ques() *Loc {
	return ti.ques
}

func (ti *TypInfo) SetQues(loc *Loc) {
	ti.ques = loc
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

type ClsTypInfo struct {
	accMod       ACC_MOD
	superTypArgs Node
	implements   []Node
	abstractLoc  *Loc
	readonlyLoc  *Loc
	overrideLoc  *Loc
	declareLoc   *Loc
}

func (ti *TypInfo) intiClsTyp() {
	if ti.clsTyp != nil {
		return
	}
	ti.clsTyp = &ClsTypInfo{}
}

func (ti *TypInfo) AccMod() ACC_MOD {
	if fuzz.IsNilPtr(ti.clsTyp) {
		return ACC_MOD_NONE
	}
	return ti.clsTyp.accMod
}

func (ti *TypInfo) SetAccMod(accMod ACC_MOD) {
	ti.intiClsTyp()
	ti.clsTyp.accMod = accMod
}

func (ti *TypInfo) SuperTypArgs() Node {
	if fuzz.IsNilPtr(ti.clsTyp) {
		return nil
	}
	return ti.clsTyp.superTypArgs
}

func (ti *TypInfo) SetSuperTypArgs(node Node) {
	ti.intiClsTyp()
	ti.clsTyp.superTypArgs = node
}

func (ti *TypInfo) AbstractLoc() *Loc {
	if fuzz.IsNilPtr(ti.clsTyp) {
		return nil
	}
	return ti.clsTyp.abstractLoc
}

func (ti *TypInfo) SetAbstractLoc(loc *Loc) {
	ti.intiClsTyp()
	ti.clsTyp.abstractLoc = loc
}

func (ti *TypInfo) ReadonlyLoc() *Loc {
	if fuzz.IsNilPtr(ti.clsTyp) {
		return nil
	}
	return ti.clsTyp.readonlyLoc
}

func (ti *TypInfo) SetReadonlyLoc(loc *Loc) {
	ti.intiClsTyp()
	ti.clsTyp.abstractLoc = loc
}

func (ti *TypInfo) OverrideLoc() *Loc {
	if fuzz.IsNilPtr(ti.clsTyp) {
		return nil
	}
	return ti.clsTyp.overrideLoc
}

func (ti *TypInfo) SetOverrideLoc(loc *Loc) {
	ti.intiClsTyp()
	ti.clsTyp.overrideLoc = loc
}

func (ti *TypInfo) DeclareLoc() *Loc {
	if fuzz.IsNilPtr(ti.clsTyp) {
		return nil
	}
	return ti.clsTyp.declareLoc
}

func (ti *TypInfo) SetDeclareLoc(loc *Loc) {
	ti.intiClsTyp()
	ti.clsTyp.declareLoc = loc
}

func (ti *TypInfo) Implements() []Node {
	if fuzz.IsNilPtr(ti.clsTyp) {
		return nil
	}
	return ti.clsTyp.implements
}

func (ti *TypInfo) SetImplements(nodes []Node) {
	ti.intiClsTyp()
	ti.clsTyp.implements = nodes
}

func (ti *TypInfo) Abstract() bool {
	if fuzz.IsNilPtr(ti.clsTyp) {
		return false
	}
	return ti.clsTyp.abstractLoc != nil
}

func (ti *TypInfo) Readonly() bool {
	if fuzz.IsNilPtr(ti.clsTyp) {
		return false
	}
	return ti.clsTyp.readonlyLoc != nil
}

func (ti *TypInfo) Override() bool {
	if fuzz.IsNilPtr(ti.clsTyp) {
		return false
	}
	return ti.clsTyp.overrideLoc != nil
}

func (ti *TypInfo) Declare() bool {
	if fuzz.IsNilPtr(ti.clsTyp) {
		return false
	}
	return ti.clsTyp.declareLoc != nil
}
