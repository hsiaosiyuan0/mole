package pack

import (
	"container/list"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	"github.com/hsiaosiyuan0/mole/ecma/astutil"
	"github.com/hsiaosiyuan0/mole/ecma/parser"
	"github.com/hsiaosiyuan0/mole/ecma/walk"
	"github.com/hsiaosiyuan0/mole/span"
	"github.com/hsiaosiyuan0/mole/util"
)

type DepUnitFact interface {
	Lang() []string
	New(*DepScanner, *DepFileReq) DepUnit
	NewModule(string) Module
}

type DepUnit interface {
	Load() error
}

const (
	TGT_WEB  string = "web"
	TGT_NODE        = "node"
	TGT_RN          = "react-native"
)

type DepScannerOpts struct {
	target string

	Dir        string
	Entries    []string
	Extensions []string

	Exports [][]string
	Imports [][]string

	Vars map[string]interface{}

	Builtin map[string]bool

	concurrent int
	unitFacts  map[string]DepUnitFact

	TsConfig   *TsConfig
	Ts         bool
	ParserOpts map[string]interface{}
}

func NewDepScannerOpts() *DepScannerOpts {
	opts := &DepScannerOpts{
		Entries:    []string{},
		concurrent: 16,
		unitFacts:  map[string]DepUnitFact{},
		Ts:         true,
	}

	opts.regUnitFact(&JsUnitFact{}).regUnitFact(&JsonUnitFact{})
	return opts
}

func (s *DepScannerOpts) ResolveBuiltin() {
	if s.target == TGT_NODE {
		s.Builtin = nodeBuiltin
	} else if s.target == TGT_RN {
		s.Builtin = rnBuiltin
	}
}

func (s *DepScannerOpts) SerVars(vars map[string]interface{}) {
	s.Vars = vars
}

func (s *DepScannerOpts) SetTsconfig(dir string, file string, ts bool) error {
	var err error
	s.TsConfig, err = NewTsConfig(dir, file)
	if err != nil {
		return err
	}

	_, err = s.TsConfig.PathMaps()
	if err != nil {
		return err
	}

	s.Ts = ts
	return nil
}

func (s *DepScannerOpts) pathMaps() *PathMaps {
	if s.TsConfig == nil {
		return nil
	}
	return s.TsConfig.pathMaps
}

func (s *DepScannerOpts) regUnitFact(u DepUnitFact) *DepScannerOpts {
	for _, lang := range u.Lang() {
		s.unitFacts[lang] = u
	}
	return s
}

type DepScanner struct {
	opts *DepScannerOpts

	fileLoader *FileLoader
	pkgLoader  *PkginfoLoader

	mId         int64
	allModules  map[int64]Module
	fileModules map[string]Module
	modulesLock sync.Mutex

	umbrellas     map[string]Module
	umbrellasLock sync.Mutex

	fileReqList     *list.List
	fileReqListLock sync.Mutex

	newJob chan bool
	wg     sync.WaitGroup
	wgFin  chan bool

	stop  chan error
	fin   chan bool
	fatal error

	minor  chan error
	minors []error
}

func NewDepScanner(opts *DepScannerOpts) *DepScanner {
	fileLoader := NewFileLoader(1024, 10)

	s := &DepScanner{
		opts: opts,

		fileLoader: fileLoader,
		pkgLoader:  NewPkginfoLoader(fileLoader),

		allModules:  map[int64]Module{},
		fileModules: map[string]Module{},
		modulesLock: sync.Mutex{},

		umbrellas:     map[string]Module{},
		umbrellasLock: sync.Mutex{},

		fileReqList:     list.New(),
		fileReqListLock: sync.Mutex{},

		newJob: make(chan bool),
		wg:     sync.WaitGroup{},
		wgFin:  make(chan bool),

		stop: make(chan error),
		fin:  make(chan bool),

		minor:  make(chan error),
		minors: []error{},
	}

	return s.initWorkers()
}

func (s *DepScanner) Modules() map[int64]Module {
	return s.allModules
}

func (s *DepScanner) Umbrellas() map[string]Module {
	return s.umbrellas
}

type FileReqTimeoutErr struct {
	Target string
	Cw     string
}

func (e *FileReqTimeoutErr) Error() string {
	return fmt.Sprintf("file request timeout, file: %s cw: %s", e.Target, e.Cw)
}

func (s *DepScanner) handleFileReq() {
loop:
	for {
		select {
		case <-s.newJob:
			req := s.shift()
			if req == nil {
				continue
			}

			unit := s.newUnit(req)
			if unit == nil {
				s.wg.Done()
				continue
			}

			done := make(chan bool)

			go func() {
				if err := unit.Load(); err != nil {
					s.Minor(err)
				}
				done <- true
			}()

			select {
			case <-time.After(5 * time.Second):
				s.Minor(&FileReqTimeoutErr{req.target, req.cw})
				s.wg.Done()
			case <-done:
				s.wg.Done()
			}

		case <-s.fin:
			break loop
		}
	}
}

func (s *DepScanner) newUnit(req *DepFileReq) DepUnit {
	uf := s.opts.unitFacts[req.lang]
	if uf == nil {
		return nil
	}
	return uf.New(s, req)
}

func (s *DepScanner) initWorkers() *DepScanner {
	for i := 0; i < s.opts.concurrent; i++ {
		go s.handleFileReq()
	}
	return s
}

func (s *DepScanner) shift() *DepFileReq {
	if s.fileReqList.Len() == 0 {
		return nil
	}
	s.fileReqListLock.Lock()
	defer s.fileReqListLock.Unlock()

	f := s.fileReqList.Front()
	s.fileReqList.Remove(f)
	return f.Value.(*DepFileReq)
}

func (s *DepScanner) push(req *DepFileReq) {
	s.fileReqListLock.Lock()
	defer s.fileReqListLock.Unlock()

	s.fileReqList.PushBack(req)
}

func (s *DepScanner) addNewJob(req *DepFileReq) {
	s.wg.Add(1)
	s.push(req)
	go func() {
		s.newJob <- true
	}()
}

func (s *DepScanner) Run() error {
	err := s.prepareEntries()
	if err != nil {
		return err
	}

	go func() {
		s.wg.Wait()
		s.wgFin <- true
	}()

loop:
	for {
		select {
		case <-s.wgFin:
			break loop
		case err := <-s.stop:
			s.fatal = err
			break loop
		case err := <-s.minor:
			s.minors = append(s.minors, err)
		}
	}

	s.fin <- true
	return s.fatal
}

func (s *DepScanner) newModule(file string) Module {
	lang := filepath.Ext(file)
	uf := s.opts.unitFacts[lang]
	if uf == nil {
		return nil
	}
	m := uf.NewModule(file)
	m.setId(atomic.AddInt64(&s.mId, 1))
	return m
}

func (s *DepScanner) getOrNewUmbrella(pi *Pkginfo) Module {
	s.umbrellasLock.Lock()
	defer s.umbrellasLock.Unlock()

	file := pi.File()

	if m, ok := s.umbrellas[file]; ok {
		return m
	}

	m := &JsModule{
		id:          atomic.AddInt64(&s.mId, 1),
		file:        file,
		name:        pi.Name,
		version:     pi.Version,
		inlets:      []*Relation{},
		inletsLock:  sync.Mutex{},
		outlets:     []*Relation{},
		outletsLock: sync.Mutex{},
	}
	m.setUmbrella(m.Id())
	s.umbrellas[file] = m
	s.addModule(m)
	return m
}

func (s *DepScanner) addModule(m Module) {
	s.modulesLock.Lock()
	defer s.modulesLock.Unlock()

	s.allModules[m.Id()] = m
}

func (s *DepScanner) getOrNewModule(file string) Module {
	s.modulesLock.Lock()
	defer s.modulesLock.Unlock()

	if m, ok := s.fileModules[file]; ok {
		return m
	}

	m := s.newModule(file)
	if m != nil {
		s.fileModules[file] = m
		s.allModules[m.Id()] = m
	}
	return m
}

func (s *DepScanner) prepareEntries() error {
	dir := s.opts.Dir
	for _, entry := range s.opts.Entries {
		file := filepath.Join(dir, entry)
		m := s.newModule(file)
		m.setAsEntry()
		s.fileModules[file] = m
		s.allModules[m.Id()] = m

		req := &DepFileReq{[]*ImportFrame{}, nil, file, dir, filepath.Ext(file)}
		s.addNewJob(req)
	}
	return nil
}

func (s *DepScanner) Minors() []error {
	return s.minors
}

func (s *DepScanner) Minor(err error) {
	s.minor <- err
}

func (s *DepScanner) Stop(err error) {
	s.stop <- err
}

func (s *DepScanner) Fin() chan bool {
	return s.fin
}

type ImportFrame struct {
	Mid       int64
	Line, Col uint32
	Import    bool
}

type DepFileReq struct {
	iptStk []*ImportFrame
	parent Module
	target string
	cw     string
	lang   string
}

type JsUnitFact struct{}

func (j *JsUnitFact) New(s *DepScanner, req *DepFileReq) DepUnit {
	opts := s.opts
	r := NewNodeResolver(
		opts.Exports, opts.Imports, opts.Extensions, opts.Builtin, s.pkgLoader, opts.Ts, opts.pathMaps())

	if s.opts.TsConfig != nil {
		r.baseUrl = s.opts.TsConfig.CompilerOptions.BaseUrl
	}

	return &JsUnit{s, req, r}
}

func (j *JsUnitFact) Lang() []string {
	return []string{".js", ".jsx", ".ts", ".tsx"}
}

func (j *JsUnitFact) NewModule(file string) Module {
	return &JsModule{
		file:        file,
		inlets:      []*Relation{},
		inletsLock:  sync.Mutex{},
		outlets:     []*Relation{},
		outletsLock: sync.Mutex{},
	}
}

type JsUnit struct {
	s   *DepScanner
	req *DepFileReq
	r   *NodeResolver
}

type importPoint struct {
	file      string
	line, col uint32
	ipt       bool
}

func parseDep(file, code string, vars map[string]interface{}, parserOpts map[string]interface{}, m *JsModule) ([]*importPoint, error) {
	s := span.NewSource(file, code)
	opts := parser.NewParserOpts()
	if parserOpts != nil {
		opts.MergeJson(parserOpts)
	}

	ext := filepath.Ext(file)
	if ext == ".ts" {
		opts.Feature = opts.Feature.On(parser.FEAT_STRICT)
		opts.Feature = opts.Feature.On(parser.FEAT_TS)
		opts.Feature = opts.Feature.Off(parser.FEAT_JSX)
	} else if ext == ".tsx" {
		opts.Feature = opts.Feature.On(parser.FEAT_STRICT)
		opts.Feature = opts.Feature.On(parser.FEAT_TS)
		opts.Feature = opts.Feature.On(parser.FEAT_JSX)
	} else if ext == ".jsx" {
		opts.Feature = opts.Feature.On(parser.FEAT_JSX)
	} else {
		opts.Feature = opts.Feature.Off(parser.FEAT_STRICT)
	}

	p := parser.NewParser(s, opts)

	ast, err := p.Prog()
	if err != nil {
		return nil, err
	}

	if m != nil {
		m.setStrict(p.Symtab().Root.IsKind(parser.SPK_STRICT))
	}

	ctx := walk.NewWalkCtx(ast, p.Symtab())
	derived := []*importPoint{}

	// collect the import statements
	walk.AddNodeAfterListener(&ctx.Listeners, parser.N_STMT_IMPORT, &walk.Listener{
		Id: "parseDep",
		Handle: func(node parser.Node, key string, ctx *walk.VisitorCtx) {
			n := node.(*parser.ImportDec)
			loc := n.Loc().Begin()
			derived = append(derived, &importPoint{n.Src().(*parser.StrLit).Text(), loc.Line, loc.Col, true})
		},
	})

	// check if the `require` has been rebound to other values
	reqRebound := false
	walk.AddNodeAfterListener(&ctx.Listeners, parser.N_EXPR_ASSIGN, &walk.Listener{
		Id: "parseDep",
		Handle: func(node parser.Node, key string, ctx *walk.VisitorCtx) {
			n := node.(*parser.AssignExpr)
			if astutil.GetName(n.Lhs()) == "require" {
				s := ctx.WalkCtx.Scope()
				ref := s.BindingOf("require")
				reqRebound = ref == nil
			}
		},
	})

	// collect the require calls first, which will be filtered by below condition judgement
	candidates := map[parser.Node]parser.Node{}
	walk.AddNodeAfterListener(&ctx.Listeners, parser.N_EXPR_CALL, &walk.Listener{
		Id: "parseDep",
		Handle: func(node parser.Node, key string, ctx *walk.VisitorCtx) {
			n := node.(*parser.CallExpr)
			s := ctx.WalkCtx.Scope()
			callee := n.Callee()
			args := n.Args()

			isRequire :=
				!reqRebound && astutil.GetName(callee) == "require" && s.BindingOf("require") == nil &&
					len(args) == 1 && args[0].Type() == parser.N_LIT_STR

			if isRequire {
				candidates[node] = node
			}
		},
	})

	// since `import` is keyword instead of variable, collect the import points directly
	walk.AddNodeAfterListener(&ctx.Listeners, parser.N_IMPORT_CALL, &walk.Listener{
		Id: "parseDep",
		Handle: func(node parser.Node, key string, ctx *walk.VisitorCtx) {
			candidates[node] = node
		},
	})

	// do the collecting processes mentioned above
	walk.VisitNode(ast, "", ctx.VisitorCtx())

	// find the all call exprs in the true branches
	interests := []parser.NodeType{parser.N_IMPORT_CALL}
	if !reqRebound {
		interests = append(interests, parser.N_EXPR_CALL)
	}
	nodes := astutil.CollectNodesInTrueBranches(ast, interests, vars)

	// filter out the dead require calls
	for _, node := range nodes {
		if it, ok := candidates[node]; ok {
			switch it.Type() {
			case parser.N_EXPR_CALL:
				n := node.(*parser.CallExpr)
				loc := n.Loc().Begin()
				derived = append(derived, &importPoint{n.Args()[0].(*parser.StrLit).Text(), loc.Line, loc.Col, false})
			case parser.N_IMPORT_CALL:
				n := node.(*parser.ImportCall)
				loc := n.Loc().Begin()
				derived = append(derived, &importPoint{n.Src().(*parser.StrLit).Text(), loc.Line, loc.Col, true})
			}
		}
	}

	return derived, nil
}

func (j *JsUnit) load(m Module) ([]byte, error) {
	f, err := j.s.fileLoader.Load(m.File())
	if err != nil {
		return nil, err
	}

	switch fv := f.(type) {
	case []byte: // done
		return fv, nil
	case chan *FileLoadResult:
		f := <-fv // wait
		if f.err != nil {
			return nil, err
		}
		return f.raw, nil
	}

	panic("unreachable")
}

func (j *JsUnit) scan(m Module) ([]*importPoint, error) {
	f, err := j.load(m)
	if err != nil {
		return nil, err
	}

	jm := m.(*JsModule)
	jm.size = int64(len(f))
	jm.scanned = true

	return parseDep(jm.file, string(f), j.s.opts.Vars, j.s.opts.ParserOpts, jm)
}

func (j *JsUnit) Load() error {
	req := j.req
	file, pi, err := j.r.Resolve(req.target, req.cw)
	if err != nil {
		return err
	}

	if len(file) == 0 {
		return nil // builtin module
	}

	m := j.s.getOrNewModule(file[0])
	if m == nil {
		return errors.New("unsupported file: " + file[0])
	}
	m.setImportStk(req.iptStk)

	if req.parent != nil {
		link(req.parent, m)
	}

	umb := j.s.getOrNewUmbrella(pi)
	m.setUmbrella(umb.Id())

	if jm, ok := m.(*JsModule); ok && !m.Scanned() {
		if !jm.IsJson() {
			derived, err := j.scan(m)
			if err != nil {
				return err
			}

			lang := filepath.Ext(file[0])
			cw := filepath.Dir(file[0])
			for _, d := range derived {
				frame := &ImportFrame{m.Id(), d.line, d.col, d.ipt}
				stk := util.Copy(req.iptStk)
				j.s.addNewJob(&DepFileReq{append(stk, frame), m, d.file, cw, lang})
			}

			umb.addSize(m.Size())
		} else if !jm.IsUmbrella() {
			jm.scanned = true

			s, err := os.Stat(file[0])
			if err != nil {
				return err
			}
			jm.size = s.Size()
			umb.addSize(jm.size)
		}
	}

	return nil
}

type JsonUnitFact struct{}

func (j *JsonUnitFact) New(s *DepScanner, req *DepFileReq) DepUnit {
	opts := s.opts
	r := NewNodeResolver(
		opts.Exports, opts.Imports, opts.Extensions, opts.Builtin, s.pkgLoader, opts.Ts, opts.pathMaps())
	return &JsonUnit{s, req, r}
}

func (j *JsonUnitFact) Lang() []string {
	return []string{".json"}
}

func (j *JsonUnitFact) NewModule(file string) Module {
	return &JsModule{
		file:    file,
		inlets:  []*Relation{},
		outlets: []*Relation{},
	}
}

type JsonUnit struct {
	s   *DepScanner
	req *DepFileReq
	r   *NodeResolver
}

func (j *JsonUnit) Load() error {
	return nil
}
