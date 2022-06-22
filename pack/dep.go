package pack

import (
	"container/list"
	"errors"
	"path/filepath"
	"sync"
	"sync/atomic"

	"github.com/hsiaosiyuan0/mole/ecma/astutil"
	"github.com/hsiaosiyuan0/mole/ecma/parser"
	"github.com/hsiaosiyuan0/mole/ecma/walk"
	"github.com/hsiaosiyuan0/mole/span"
)

type DepUnitFact interface {
	Lang() []string
	New(*DepScanner, *DepFileReq) DepUnit
	NewModule(string) Module
}

type DepUnit interface {
	Load() error
}

type DepScannerOpts struct {
	dir        string
	entries    []string
	extensions []string

	exports [][]string
	imports [][]string

	builtin map[string]bool

	concurrent int
	unitFacts  map[string]DepUnitFact

	tsConfig *TsConfig
	ts       bool
}

func NewDepScannerOpts() *DepScannerOpts {
	opts := &DepScannerOpts{
		entries:    []string{},
		concurrent: 16,
		unitFacts:  map[string]DepUnitFact{},
	}

	opts.regUnitFact(&JsUnitFact{})
	return opts
}

func (s *DepScannerOpts) SetTsconfig(dir string, file string, ts bool) error {
	var err error
	s.tsConfig, err = NewTsConfig(dir, file)
	if err != nil {
		return err
	}

	_, err = s.tsConfig.PathMaps()
	if err != nil {
		return err
	}

	s.ts = ts
	return nil
}

func (s *DepScannerOpts) regUnitFact(u DepUnitFact) {
	for _, lang := range u.Lang() {
		s.unitFacts[lang] = u
	}
}

type DepScanner struct {
	opts *DepScannerOpts

	fileLoader *FileLoader
	pkgLoader  *PkginfoLoader

	mId         uint64
	modules     map[string]Module
	modulesLock sync.Mutex

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

		modules:     map[string]Module{},
		modulesLock: sync.Mutex{},

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

			if err := unit.Load(); err != nil {
				s.Minor(err)
			}
			s.wg.Done()

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
	m.setId(atomic.AddUint64(&s.mId, 1))
	return m
}

func (s *DepScanner) getOrNewModule(file string) Module {
	s.modulesLock.Lock()
	defer s.modulesLock.Unlock()

	if m, ok := s.modules[file]; ok {
		return m
	}

	m := s.newModule(file)
	s.modules[file] = m
	return m
}

func (s *DepScanner) prepareEntries() error {
	dir := s.opts.dir
	for _, entry := range s.opts.entries {
		file := filepath.Join(dir, entry)
		m := s.newModule(file)
		m.setAsEntry()
		s.modules[file] = m

		req := &DepFileReq{file, dir, filepath.Ext(file)}
		s.addNewJob(req)
	}
	return nil
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

type DepFileReq struct {
	target string
	cw     string
	lang   string
}

type JsUnitFact struct{}

func (j *JsUnitFact) New(s *DepScanner, req *DepFileReq) DepUnit {
	opts := s.opts
	r := NewNodeResolver(opts.exports, opts.imports, opts.extensions, opts.builtin,
		s.pkgLoader, opts.ts, opts.tsConfig.pathMaps)
	return &JsUnit{s, req, r}
}

func (j *JsUnitFact) Lang() []string {
	return []string{".js", ".jsx", ".ts", ".tsx"}
}

func (j *JsUnitFact) NewModule(file string) Module {
	return &JsModule{
		file:    file,
		inlets:  []*Relation{},
		outlets: []*Relation{},
	}
}

type JsUnit struct {
	s   *DepScanner
	req *DepFileReq
	r   *NodeResolver
}

func getParentCond(ctx *walk.VisitorCtx) parser.Node {
	barrier := []parser.NodeType{parser.N_EXPR_FN, parser.N_STMT_FN, parser.N_EXPR_ARROW}
	pn, ctx := astutil.GetParent(ctx, []parser.NodeType{parser.N_STMT_IF}, barrier)
	if pn != nil {
		return pn
	}

	pn, ctx = astutil.GetParent(ctx, []parser.NodeType{parser.N_EXPR_COND}, barrier)
	if pn != nil {
		pnIf, _ := astutil.GetParent(ctx, []parser.NodeType{parser.N_STMT_IF}, barrier)
		if pnIf != nil {
			return pnIf
		}
		return pn
	}

	pn, ctx = astutil.GetParent(ctx, []parser.NodeType{parser.N_EXPR_BIN}, barrier)
	if pn != nil {
		op := pn.(*parser.BinExpr).Op()
		if op == parser.T_AND || op == parser.T_OR {
			pnIf, _ := astutil.GetParent(ctx, []parser.NodeType{parser.N_STMT_IF}, barrier)
			if pnIf != nil {
				return pnIf
			}
			return pn
		}
	}

	return nil
}

// TODO: select require calls from true branches

func parseDep(file, code string) ([]string, error) {
	s := span.NewSource(file, code)
	p := parser.NewParser(s, parser.NewParserOpts())
	ast, err := p.Prog()
	if err != nil {
		return nil, err
	}

	ctx := walk.NewWalkCtx(ast, p.Symtab())
	derived := []string{}
	walk.SetVisitor(&ctx.Visitors, parser.N_STMT_IMPORT, func(node parser.Node, key string, ctx *walk.VisitorCtx) {
		n := node.(*parser.ImportDec)
		derived = append(derived, n.Src().(*parser.StrLit).Text())
	})

	reqRebound := false
	walk.SetVisitor(&ctx.Visitors, parser.N_EXPR_ASSIGN, func(node parser.Node, key string, ctx *walk.VisitorCtx) {
		n := node.(*parser.AssignExpr)
		if astutil.GetName(n.Lhs()) == "require" {
			s := ctx.WalkCtx.Scope()
			ref := s.BindingOf("require")
			reqRebound = ref == nil
		}
	})

	walk.SetVisitor(&ctx.Visitors, parser.N_EXPR_CALL, func(node parser.Node, key string, ctx *walk.VisitorCtx) {
		n := node.(*parser.CallExpr)
		callee := n.Callee()
		s := ctx.WalkCtx.Scope()
		isRequireCall := !reqRebound && astutil.GetName(callee) == "require" &&
			s.BindingOf("require") == nil && len(n.Args()) == 1 && n.Args()[0].Type() == parser.N_LIT_STR
		if isRequireCall {

			derived = append(derived, n.Args()[0].(*parser.StrLit).Text())
		}
	})

	walk.SetVisitor(&ctx.Visitors, parser.N_IMPORT_CALL, func(node parser.Node, key string, ctx *walk.VisitorCtx) {
		n := node.(*parser.ImportCall)
		derived = append(derived, n.Src().(*parser.StrLit).Text())
	})

	walk.VisitNode(ast, "", ctx.VisitorCtx())
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

func (j *JsUnit) parse(m Module) ([]string, error) {
	f, err := j.load(m)
	if err != nil {
		return nil, err
	}

	jm := m.(*JsModule)
	jm.size = len(f)
	jm.parsed = true
	return parseDep(jm.file, string(f))
}

func (j *JsUnit) Load() error {
	req := j.req
	file, err := j.r.Resolve(req.target, req.cw)
	if err != nil {
		return err
	}

	m := j.s.getOrNewModule(file[0])
	if m == nil {
		return errors.New("unsupported file: " + file[0])
	}

	if !m.Parsed() {
		derived, err := j.parse(m)
		if err != nil {
			return err
		}

		lang := filepath.Ext(file[0])
		cw := filepath.Dir(file[0])
		for _, f := range derived {
			j.s.addNewJob(&DepFileReq{f, cw, lang})
		}
	}

	return nil
}
