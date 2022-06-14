package linter

import (
	"encoding/json"
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
	Lang() string
	Config() *Config
	Report(*Diagnosis)
}

type RuleCtx struct {
	unit     Unit
	ruleFact RuleFact
}

func (u *RuleCtx) Config() *Config {
	return u.unit.Config()
}

func (u *RuleCtx) Report(node parser.Node, msg string, level DiagLevel) {
	lang := u.unit.Lang()
	rule := u.ruleFact.Name()
	lvl := u.Config().LevelOfRule(lang, rule)
	if lvl == DL_NONE {
		lvl = level
	}

	dig := &Diagnosis{
		Loc:   node.Loc().Clone(),
		Lang:  lang,
		Rule:  rule,
		Msg:   msg,
		Level: lvl,
	}
	u.unit.Report(dig)
}

type JsUnit struct {
	linter *Linter
	cfg    *Config

	file   string
	lang   string
	ast    parser.Node
	symTab *parser.SymTab
	ana    *analysis.Analysis

	rules      map[string]*Rule
	skipped    map[string]*Rule
	ephSkipped map[string]*Rule
}

func NewJsUnit(file string, code string, cfg *Config) (*JsUnit, error) {
	u := &JsUnit{
		cfg:        cfg,
		file:       file,
		lang:       filepath.Ext(file),
		rules:      map[string]*Rule{},
		skipped:    map[string]*Rule{},
		ephSkipped: map[string]*Rule{},
	}

	s := span.NewSource(file, string(code))

	p := parser.NewParser(s, cfg.ParserOpts())
	ast, err := p.Prog()
	if err != nil {
		return nil, err
	}

	u.ast = ast
	u.symTab = p.Symtab()
	u.ana = analysis.NewAnalysis(ast, u.symTab)

	return u, nil
}

func (u *JsUnit) Config() *Config {
	return u.cfg
}

func (u *JsUnit) Lang() string {
	return u.lang
}

func (u *JsUnit) Report(dig *Diagnosis) {
	u.linter.report(dig)
}

func isJsFile(f string) bool {
	ext := filepath.Ext(f)
	return ext == ".js" || ext == ".jsx"
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
		ctx := &RuleCtx{u, rf}
		rule := &Rule{rf, map[parser.NodeType]*walk.Listener{}}
		for nt, fn := range rf.Create(ctx) {
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
			walk.AddListener(&u.ana.WalkCtx.Listeners, nt, listener)
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

	// file => line => diagnoses
	diags     map[string]map[uint32][]*Diagnosis
	diagsLock sync.Mutex
}

func NewLinter(dir string, cfg *Config, skipBuiltin bool) (*Linter, error) {
	var err error

	if cfg == nil {
		cfg, err = LoadCfgInDir(dir, nil)
		if err != nil {
			return nil, err
		}
		if cfg == nil {
			return nil, &LoadConfigErr{"no config file detected", nil}
		}
		if err = cfg.Init(); err != nil {
			return nil, err
		}
	}

	if !skipBuiltin {
		// inherits ruleFacts from builtin
		for key, roleFacts := range builtinRuleFacts {
			cfg.ruleFacts[key] = map[string]RuleFact{}
			for name, roleFact := range roleFacts {
				cfg.ruleFacts[key][name] = roleFact
			}
		}
	}

	l := &Linter{
		dir:       dir,
		cfg:       cfg,
		cfgMap:    map[string]*Config{},
		diags:     map[string]map[uint32][]*Diagnosis{},
		diagsLock: sync.Mutex{},
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

func (l *Linter) report(dig *Diagnosis) {
	l.diagsLock.Lock()
	defer l.diagsLock.Unlock()

	file := dig.Loc.Source()
	line := dig.Loc.Begin().Line()

	list := l.diags[file]
	if list == nil {
		list = map[uint32][]*Diagnosis{}
		l.diags[file] = list
	}

	diags := list[line]
	if diags == nil {
		diags = []*Diagnosis{}
		list[line] = diags
	}

	list[line] = append(diags, dig)
}

func (l *Linter) Process() *Reports {
	w := util.NewDirWalker(l.dir, 0, func(f string, dir bool, dw *util.DirWalker) {
		if dir {
			if _, ok := l.cfgMap[f]; ok {
				return
			}

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

			code, err := ioutil.ReadFile(f)
			if err != nil {
				dw.Stop(err)
			}

			u, err := NewJsUnit(f, string(code), cfg)
			if err != nil {
				dw.Stop(err)
			}

			u.linter = l
			u.initRules().enableAllRules()
			u.ana.Analyze()
		}
	})

	w.Walk()

	return l.genReports(w.Err())
}

func (l *Linter) genReports(internalErr error) *Reports {
	r := &Reports{
		InternalError: internalErr,
		Diagnoses:     []*Diagnosis{},
	}

	for _, line := range l.diags {
		for _, dig := range line {
			r.Diagnoses = append(r.Diagnoses, dig...)
		}
	}
	return r
}

type DiagLevel uint16

const (
	DL_NONE DiagLevel = iota
	DL_WARNING
	DL_ERROR
)

type Diagnosis struct {
	Lang string
	Rule string

	Loc   *parser.Loc
	Level DiagLevel
	Msg   string
}

func (d *Diagnosis) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Line  uint32 `json:"line"`
		Col   uint32 `json:"col"`
		Lang  string `json:"lang"`
		Rule  string `json:"rule"`
		Level uint16 `json:"level"`
		Msg   string `json:"Msg"`
	}{
		Line:  d.Loc.Begin().Line(),
		Col:   d.Loc.Begin().Column(),
		Lang:  d.Lang,
		Rule:  d.Rule,
		Level: uint16(d.Level),
		Msg:   d.Msg,
	})
}

type Reports struct {
	InternalError error        `json:"internalError"`
	Diagnoses     []*Diagnosis `json:"diagnoses"`
}
