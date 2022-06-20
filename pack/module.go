package pack

import "github.com/hsiaosiyuan0/mole/span"

type Relation struct {
	Lhs   uint64
	Rhs   uint64
	Cross bool
	Loc   *span.Range
}

// methods on `Module` should be thread-safe
type Module interface {
	Id() uint64
	setId(id uint64)

	Lang() string

	setFile(string)
	File() string

	Size() int

	Entry() bool
	setAsEntry()

	Ext() bool
	ExtMain() uint64

	Parsed() bool

	Inlets() []*Relation
	Outlets() []*Relation
}

type JsModule struct {
	id   uint64
	lang string

	file   string
	size   int
	parsed bool

	entry   bool
	ext     bool
	extMain uint64

	inlets  []*Relation
	outlets []*Relation
}

func (m *JsModule) setId(id uint64) {
	m.id = id
}

func (m *JsModule) Id() uint64 {
	return m.id
}

func (m *JsModule) Lang() string {
	return m.lang
}

func (m *JsModule) setFile(file string) {
	m.file = file
}

func (m *JsModule) File() string {
	return m.file
}

func (m *JsModule) Size() int {
	return m.size
}

func (m *JsModule) setAsEntry() {
	m.entry = true
}

func (m *JsModule) Entry() bool {
	return m.entry
}

func (m *JsModule) Parsed() bool {
	return m.parsed
}

func (m *JsModule) Ext() bool {
	return m.ext
}

func (m *JsModule) ExtMain() uint64 {
	return m.extMain
}

func (m *JsModule) Inlets() []*Relation {
	return m.inlets
}

func (m *JsModule) Outlets() []*Relation {
	return m.outlets
}
