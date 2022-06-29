package pack

import (
	"encoding/json"
	"strings"
	"sync"
	"sync/atomic"
)

type Relation struct {
	Lhs int64 `json:"lhs"`
	Rhs int64 `json:"rhs"`
}

func FindOutlet(a Module, dst int64, create bool) (*Relation, bool) {
	rs := a.Outlets()

	var edge *Relation
	for _, edge = range rs {
		if dst == edge.Rhs {
			return edge, false
		}
	}

	if edge == nil && create {
		edge = &Relation{}
		edge.Lhs = a.Id()
		edge.Rhs = dst
		return edge, true
	}
	return nil, false
}

func FindInlet(a Module, src int64, create bool) (*Relation, bool) {
	rs := a.Inlets()

	var edge *Relation
	for _, edge = range rs {
		if src == edge.Lhs {
			return edge, false
		}
	}

	if edge == nil && create {
		edge = &Relation{}
		edge.Lhs = src
		edge.Rhs = a.Id()
		return edge, true
	}
	return nil, false
}

func link(a, b Module) {
	a.AddOutlet(b.Id())
	b.AddInlet(a.Id())
}

// methods on `Module` should be thread-safe
type Module interface {
	Id() int64
	setId(id int64)

	Name() string
	Version() string

	Lang() string

	setFile(string)
	File() string

	setStrict(bool)
	Strict() bool

	addSize(int64)
	Size() int64
	setSize(int64)

	Entry() bool
	setAsEntry()

	Outside() bool

	setUmbrella(int64)
	Umbrella() int64

	IsUmbrella() bool

	Scanned() bool

	AddInlet(int64)
	Inlets() []*Relation

	AddOutlet(int64)
	Outlets() []*Relation

	setImportStk([]*ImportFrame)
	ImportStk() []*ImportFrame

	MarshalJSON() ([]byte, error)
}

type JsModule struct {
	id   int64
	lang string

	name    string
	version string

	file    string
	size    int64
	strict  bool
	scanned bool

	entry    bool
	outside  bool
	umbrella int64

	inlets     []*Relation
	inletsLock sync.Mutex

	outlets     []*Relation
	outletsLock sync.Mutex

	iptStk []*ImportFrame
}

func (m *JsModule) setId(id int64) {
	m.id = id
}

func (m *JsModule) Id() int64 {
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

func (m *JsModule) addSize(s int64) {
	atomic.AddInt64(&m.size, s)
}

func (m *JsModule) Size() int64 {
	return m.size
}

func (m *JsModule) setSize(s int64) {
	atomic.StoreInt64(&m.size, s)
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

func (m *JsModule) setUmbrella(id int64) {
	m.umbrella = id
}

func (m *JsModule) IsUmbrella() bool {
	return m.id == m.umbrella
}

func (m *JsModule) setStrict(f bool) {
	m.strict = f
}

func (m *JsModule) Strict() bool {
	return m.strict
}

func (m *JsModule) IsJson() bool {
	return strings.HasSuffix(m.file, ".json")
}

func (m *JsModule) Umbrella() int64 {
	return m.umbrella
}

func (m *JsModule) AddInlet(src int64) {
	m.inletsLock.Lock()
	defer m.inletsLock.Unlock()

	edge, new := FindInlet(m, src, true)
	if new {
		m.inlets = append(m.inlets, edge)
	}
}

func (m *JsModule) Inlets() []*Relation {
	return m.inlets
}

func (m *JsModule) AddOutlet(dst int64) {
	m.outletsLock.Lock()
	defer m.outletsLock.Unlock()

	edge, new := FindOutlet(m, dst, true)
	if new {
		m.outlets = append(m.outlets, edge)
	}
}

func (m *JsModule) Outlets() []*Relation {
	return m.outlets
}

func (m *JsModule) setImportStk(s []*ImportFrame) {
	m.iptStk = s
}

func (m *JsModule) ImportStk() []*ImportFrame {
	return m.iptStk
}

func (m *JsModule) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		ID       int64       `json:"id"`
		Name     string      `json:"name"`
		Version  string      `json:"version"`
		File     string      `json:"file"`
		Size     int64       `json:"size"`
		Strict   bool        `json:"strict"`
		Entry    bool        `json:"entry"`
		Umbrella int64       `json:"umbrella"`
		Inlets   []*Relation `json:"inlets"`
		Outlets  []*Relation `json:"outlets"`
	}{
		ID:       m.id,
		Name:     m.name,
		Version:  m.version,
		File:     m.file,
		Size:     m.size,
		Strict:   m.strict,
		Entry:    m.entry,
		Umbrella: m.umbrella,
		Inlets:   m.inlets,
		Outlets:  m.outlets,
	})
}
