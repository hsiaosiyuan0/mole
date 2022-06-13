package linter

import (
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"sync"

	"github.com/hsiaosiyuan0/mole/ecma/analysis"
	"github.com/hsiaosiyuan0/mole/ecma/parser"
	"github.com/hsiaosiyuan0/mole/ecma/walk"
	"github.com/hsiaosiyuan0/mole/span"
	"github.com/hsiaosiyuan0/mole/util"
)

type Unit interface {
	Config() *Config
}

type JsUnit struct {
	cfg *Config

	file    string
	ast     parser.Node
	symTab  *parser.SymTab
	walkCtx *walk.WalkCtx
	ana     *analysis.Analysis

	rules      map[string]*Rule
	skipped    map[string]*Rule
	ephSkipped map[string]*Rule
}

func NewJsUnit(file string, cfg *Config) (*JsUnit, error) {
	u := &JsUnit{
		cfg:        cfg,
		file:       file,
		rules:      map[string]*Rule{},
		skipped:    map[string]*Rule{},
		ephSkipped: map[string]*Rule{},
	}

	code, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	s := span.NewSource(file, string(code))

	p := parser.NewParser(s, cfg.ParserOpts())
	ast, err := p.Prog()
	if err != nil {
		return u, nil
	}

	u.ast = ast
	u.symTab = p.Symtab()
	u.walkCtx = walk.NewWalkCtx(ast, u.symTab)

	return u, nil
}

func isJsFile(f string) bool {
	ext := filepath.Ext(f)
	return ext == ".js" || ext == ".jsx"
}

func (u *JsUnit) Config() *Config {
	return u.cfg
}

func (u *JsUnit) File() string {
	return u.file
}

func (u *JsUnit) Ast() parser.Node {
	return u.ast
}

func (u *JsUnit) SymTab() *parser.SymTab {
	return u.symTab
}

func (u *JsUnit) Analysis() *analysis.Analysis {
	return u.ana
}

func (u *JsUnit) initRules() *JsUnit {
	lang := path.Ext(u.file)
	for _, rf := range u.cfg.ruleFacts[lang] {
		for nt, fn := range rf.Create(u) {
			rule := &Rule{rf, map[parser.NodeType]*walk.Listener{}}
			id := fmt.Sprintf("%s_%s_%d", lang, rf.Name(), nt)
			rule.Listeners[nt] = &walk.Listener{Id: id, Handle: fn}
			u.skipped[rf.Name()] = rule
		}
	}
	return u
}

func (u *JsUnit) enableAllRules() *JsUnit {
	for rn, rule := range u.skipped {
		for nt, listener := range rule.Listeners {
			walk.AddListener(&u.walkCtx.Listeners, nt, listener)
		}
		u.rules[rn] = rule
		delete(u.skipped, rn)
	}
	return u
}

type Linter struct {
	dir string
	cfg *Config

	// dir => cfg
	cfgMap       map[string]*Config
	cfgMapRwLock sync.RWMutex

	// lang => line => diagnoses
	results map[string]map[int][]*Diagnosis
}

func NewLinter(dir string) (*Linter, error) {
	cfg, err := LoadCfgInDir(dir, nil)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, &LoadConfigErr{"no config file detected", nil}
	}
	if err = cfg.Init(); err != nil {
		return nil, err
	}

	l := &Linter{
		dir:     dir,
		cfg:     cfg,
		cfgMap:  map[string]*Config{},
		results: map[string]map[int][]*Diagnosis{},
	}
	l.setCfgMap(dir, cfg)
	return l, nil
}

func (l *Linter) setCfgMap(dir string, cfg *Config) {
	l.cfgMapRwLock.Lock()
	defer l.cfgMapRwLock.Unlock()

	l.cfgMap[dir] = cfg
}

func (l *Linter) cfgOfDir(dir string) *Config {
	l.cfgMapRwLock.RLock()
	defer l.cfgMapRwLock.RUnlock()

	return l.cfgMap[dir]
}

func (l *Linter) outerCfg(file string) *Config {
	l.cfgMapRwLock.RLock()
	defer l.cfgMapRwLock.RUnlock()

	file = filepath.Dir(file)
	if l.cfgMap[file] != nil {
		return l.cfgMap[file]
	}
	return l.cfg
}

func (l *Linter) Process() {
	w := util.NewDirWalker(l.dir, 0, func(f string, dir bool, dw *util.DirWalker) {
		if dir {
			pc := l.outerCfg(f)
			cfg, err := LoadCfgInDir(f, pc)
			if err != nil {
				dw.Stop(err)
				return
			}
			if cfg != nil {
				if err = cfg.Init(); err != nil {
					dw.Stop(err)
					return
				}

				l.setCfgMap(f, cfg)
			}
			return
		}

		if isJsFile(f) {
			cfg := l.outerCfg(f)
			u, err := NewJsUnit(f, cfg)
			if err != nil {
				dw.Stop(err)
			}

			u.initRules().enableAllRules()
			walk.VisitNode(u.ast, "", u.walkCtx.VisitorCtx())

			fmt.Println(f)
		}

	})

	w.Walk()
}

type Severity uint16

const (
	STY_UNKNOWN Severity = iota
	STY_WARNING
	STY_ERROR

	STY_INTERNAL_ERROR
)

type Diagnosis struct {
	Severity Severity
	Msg      string
}
