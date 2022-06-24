package pack

import (
	"strings"

	"github.com/hsiaosiyuan0/mole/span"
)

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

	Name() string
	Version() string

	Lang() string

	setFile(string)
	File() string

	Size() int

	Entry() bool
	setAsEntry()

	Outside() bool

	setUmbrella(uint64)
	Umbrella() uint64

	IsUmbrella() bool

	Scanned() bool

	Inlets() []*Relation
	Outlets() []*Relation
}

type JsModule struct {
	id   uint64
	lang string

	name    string
	version string

	file    string
	size    int
	scanned bool

	entry    bool
	outside  bool
	umbrella uint64

	inlets  []*Relation
	outlets []*Relation
}

func (m *JsModule) setId(id uint64) {
	m.id = id
}

func (m *JsModule) Id() uint64 {
	return m.id
}

func (m *JsModule) Name() string {
	return m.name
}

func (m *JsModule) Version() string {
	return m.version
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

func (m *JsModule) Scanned() bool {
	return m.scanned
}

func (m *JsModule) Outside() bool {
	return m.outside
}

func (m *JsModule) setUmbrella(id uint64) {
	m.umbrella = id
}

func (m *JsModule) IsUmbrella() bool {
	return m.id == m.umbrella
}

func (m *JsModule) IsJson() bool {
	return strings.HasSuffix(m.file, ".json")
}

func (m *JsModule) Umbrella() uint64 {
	return m.umbrella
}

func (m *JsModule) Inlets() []*Relation {
	return m.inlets
}

func (m *JsModule) Outlets() []*Relation {
	return m.outlets
}
