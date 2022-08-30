package pack

import (
	"container/list"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/hsiaosiyuan0/mole/ecma/astutil"
	"github.com/hsiaosiyuan0/mole/ecma/parser"
	"github.com/hsiaosiyuan0/mole/ecma/walk"
	"github.com/hsiaosiyuan0/mole/pack/resolver"
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

	TsConfig   *resolver.TsConfig
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
		s.Builtin = resolver.NodeBuiltin
	} else if s.target == TGT_RN {
		s.Builtin = resolver.RnBuiltin
	}
}

func (s *DepScannerOpts) SerVars(vars map[string]interface{}) {
	s.Vars = vars
}

func (s *DepScannerOpts) SetTsconfig(dir string, file string, ts bool) error {
	var err error

	s.Ts = ts

	s.TsConfig, err = resolver.NewTsConfig(dir, file)
	if err != nil {
		return err
	}

	return nil
}

func (s *DepScannerOpts) FillDefault() {
	if len(s.Exports) == 0 {
		s.Exports = [][]string{{"browser", "require"}}
	}
	if len(s.Imports) == 0 {
		s.Imports = [][]string{{"browser", "require"}}
	}

	s.Exports = append(s.Exports, []string{"default"})
	s.Imports = append(s.Imports, []string{"default"})

	if s.Extensions == nil {
		if s.Ts {
			s.Extensions = resolver.DefaultTsExts
		} else {
			s.Extensions = resolver.DefaultJsExts
		}
	}
}

func (s *DepScannerOpts) regUnitFact(u DepUnitFact) *DepScannerOpts {
	for _, lang := range u.Lang() {
		s.unitFacts[lang] = u
	}
	return s
}

type DepScanner struct {
	opts *DepScannerOpts

	fileLoader  *resolver.FileLoader
	pkgLoader   *resolver.PjsonLoader
	modResolver *resolver.ModResolver

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
	fileLoader := resolver.NewFileLoader(2048, 128)

	s := &DepScanner{
		opts: opts,

		fileLoader: fileLoader,
		pkgLoader:  resolver.NewPjsonLoader(fileLoader),

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

	browser := opts.target != TGT_NODE
	s.pkgLoader.SetBrowser(browser)

	baseUrl := ""
	var pathMaps *resolver.PathMaps
	if opts.TsConfig != nil {
		baseUrl = opts.TsConfig.CompilerOptions.BaseUrl
		pathMaps = opts.TsConfig.PathMaps()
	}

	s.modResolver = resolver.NewModResolver(browser, opts.Imports, opts.Exports, opts.Extensions, opts.Builtin, baseUrl, pathMaps, s.pkgLoader)

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

func (s *DepScanner) getOrNewUmbrella(pi *resolver.PkgJson) Module {
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

	entries := make([]string, 0, len(s.opts.Entries))
	for _, entry := range s.opts.Entries {
		if strings.IndexRune(entry, '*') != -1 {
			p := filepath.Join(dir, entry)
			matches, err := filepath.Glob(p)
			if err != nil {
				fmt.Printf("deformed entry `%s` with error `%v` ", entry, err)
			} else {
				entries = append(entries, matches...)
			}
		} else {
			entries = append(entries, entry)
		}
	}

	sc := s.modResolver.LookupPkgScope(dir)
	if sc == nil {
		return errors.New("no package.json detected under " + dir)
	}

	for _, file := range entries {
		m := s.newModule(file)
		m.setAsEntry()
		s.fileModules[file] = m
		s.allModules[m.Id()] = m

		req := &DepFileReq{true, sc, []*ImportFrame{}, nil, file, dir, filepath.Ext(file)}
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
	S      *span.Source
	Mid    int64
	Rng    span.Range
	Import bool
}

type DepFileReq struct {
	entry  bool
	sc     *resolver.PkgJson
	stk    []*ImportFrame
	parent Module
	target string
	cw     string
	lang   string
}

type JsUnitFact struct{}

func (j *JsUnitFact) New(s *DepScanner, req *DepFileReq) DepUnit {
	return &JsUnit{s, req}
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
}

type importPoint struct {
	s    *span.Source
	file string
	rng  span.Range
	ipt  bool
}

var isFlowReg = regexp.MustCompile(`(?s)/\*.*?\*\s*@flow.*?\*/`)

func isFlow(code string) bool {
	return len(isFlowReg.Find([]byte(code))) > 0
}

func parseDep(file, code string, vars map[string]interface{}, parserOpts map[string]interface{}, m *JsModule, skipFlow bool) ([]*importPoint, int64, int64, error) {
	if skipFlow && isFlow(code) {
		return []*importPoint{}, 0, 0, nil
	}

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
		if strings.HasSuffix(file, ".d.ts") {
			opts.Feature = opts.Feature.On(parser.FEAT_DTS)
		}
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

	start := time.Now()
	ast, err := p.Prog()
	parseTime := time.Since(start)

	if err != nil {
		return nil, 0, 0, err
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
			derived = append(derived, &importPoint{p.Source(), n.Src().(*parser.StrLit).Val(), n.Range(), true})
		},
	})

	// collect the export statements
	walk.AddNodeAfterListener(&ctx.Listeners, parser.N_STMT_EXPORT, &walk.Listener{
		Id: "parseDep",
		Handle: func(node parser.Node, key string, ctx *walk.VisitorCtx) {
			n := node.(*parser.ExportDec)
			if n.Src() != nil {
				derived = append(derived, &importPoint{p.Source(), n.Src().(*parser.StrLit).Val(), n.Range(), true})
			}
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
	start = time.Now()
	walk.VisitNode(ast, "", ctx.VisitorCtx())
	walkTime := time.Since(start)

	// find the all call exprs in the true branches
	interests := []parser.NodeType{parser.N_IMPORT_CALL}
	if !reqRebound {
		interests = append(interests, parser.N_EXPR_CALL)
	}
	nodes := astutil.CollectNodesInTrueBranches(ast, interests, vars, p)

	// filter out the dead require calls
	for _, node := range nodes {
		if it, ok := candidates[node]; ok {
			switch it.Type() {
			case parser.N_EXPR_CALL:
				n := node.(*parser.CallExpr)
				derived = append(derived, &importPoint{p.Source(), n.Args()[0].(*parser.StrLit).Val(), n.Range(), false})
			case parser.N_IMPORT_CALL:
				n := node.(*parser.ImportCall)
				derived = append(derived, &importPoint{p.Source(), n.Src().(*parser.StrLit).Val(), n.Range(), true})
			}
		}
	}

	return derived, parseTime.Nanoseconds(), walkTime.Nanoseconds(), nil
}

func (j *JsUnit) load(m Module) ([]byte, error) {
	c := j.s.fileLoader.Load(m.File())

	f := <-c // wait
	return f.Raw, f.Err
}

func (j *JsUnit) scan(m Module) ([]*importPoint, error) {
	f, err := j.load(m)
	if err != nil {
		return nil, err
	}

	jm := m.(*JsModule)
	jm.size = int64(len(f))
	jm.scanned = true

	stk, parseTime, walkTime, err := parseDep(jm.file, string(f), j.s.opts.Vars, j.s.opts.ParserOpts, jm, j.s.opts.target == TGT_RN)
	if err != nil {
		return nil, err
	}

	jm.setParseTime(parseTime)
	jm.setWalkTime(walkTime)
	return stk, nil
}

func (j *JsUnit) Load() error {
	req := j.req
	t, err := j.s.modResolver.NewTask(req.target, req.cw, req.sc, nil)

	if err != nil {
		return err
	}

	r, err := t.Resolve()
	if err != nil {
		return err
	}

	if r == nil {
		return nil // builtin module or ignored
	}

	m := j.s.getOrNewModule(r.File)
	if m == nil {
		return errors.New("unsupported file: " + r.File)
	}
	m.setImportStk(req.stk)

	if req.parent != nil {
		link(req.parent, m)
	}

	umb := j.s.getOrNewUmbrella(r.Pjson)
	m.setUmbrella(umb.Id())

	if jm, ok := m.(*JsModule); ok && !m.Scanned() {
		if !jm.IsJson() {
			derived, err := j.scan(m)
			if err != nil {
				return err
			}

			curLang := filepath.Ext(r.File)
			cw := filepath.Dir(r.File)
			for _, d := range derived {
				frame := &ImportFrame{d.s, m.Id(), d.rng, d.ipt}
				stk := util.Copy(req.stk)
				lang := filepath.Ext(d.file)
				if lang == "" {
					// if there is no ext in the importing target use
					// the host file ext instead
					lang = curLang
				}
				j.s.addNewJob(&DepFileReq{false, req.sc, append(stk, frame), m, d.file, cw, lang})
			}

			umb.addSize(m.Size())
		} else if !jm.IsUmbrella() {
			jm.scanned = true

			s, err := os.Stat(r.File)
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
	return &JsonUnit{s, req}
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
}

func (j *JsonUnit) Load() error {
	return nil
}
