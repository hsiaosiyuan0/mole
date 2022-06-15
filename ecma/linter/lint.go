package linter

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"regexp"
	"strings"
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

	file string
	lang string

	parser *parser.Parser
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

	u.parser = p
	u.ana = analysis.NewAnalysis(ast, u.parser.Symtab())

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
	return u.parser.Ast()
}

func (u *JsUnit) SymTab() *parser.SymTab {
	return u.parser.Symtab()
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
			id := fmt.Sprintf("_%s_%s_%d", lang, rf.Name(), nt)
			rule.Listeners[nt] = &walk.Listener{Id: id, Handle: fn}
			u.skipped[rf.Name()] = rule
		}
	}
	return u
}

// represents the command triggered via the comments:
// - eslint-disable
// - eslint-enable
// - eslint-disable no-alert, no-console
// - eslint-enable no-alert, no-console
// - eslint-disable-next-line
// - eslint-disable-next-line no-alert, no-console
type CmtCmd struct {
	enableRule      []string
	disableRule     []string
	disableNextLine []string
}

func (u *JsUnit) parseCmtPrefix(c string) (string, []string) {
	s := strings.Index(c, " ")
	if s == -1 {
		return "", nil
	}

	p := c[:s]
	rules := strings.Split(c[s:], ",")
	for i, r := range rules {
		rules[i] = strings.Trim(r, " ")
	}
	return p, rules
}

func (u *JsUnit) parseCmt(cs *span.Range) *CmtCmd {
	cmt := cs.Text()
	if strings.HasPrefix(cmt, "/*") {
		cmt = regexp.MustCompile(`^/\*\s*|\s*\*/`).ReplaceAllString(cmt, "")
	} else {
		cmt = regexp.MustCompile(`^//\s*|\s*$`).ReplaceAllString(cmt, "")
	}

	switch cmt {
	case "eslint-disable":
		return &CmtCmd{disableRule: []string{}}
	case "eslint-enable":
		return &CmtCmd{enableRule: []string{}}
	case "eslint-disable-next-line":
		return &CmtCmd{disableNextLine: []string{}}
	}

	p, rules := u.parseCmtPrefix(cmt)
	if rules == nil {
		return nil
	}

	switch p {
	case "eslint-disable":
		return &CmtCmd{disableRule: rules}
	case "eslint-enable":
		return &CmtCmd{enableRule: rules}
	case "eslint-disable-next-line":
		return &CmtCmd{disableNextLine: rules}
	}

	return nil
}

func (u *JsUnit) handCmt(c *span.Range) {
	cmd := u.parseCmt(c)
	if cmd == nil {
		return
	}

	if cmd.disableRule != nil {
		if len(cmd.disableRule) == 0 {
			u.disableAllRules(false)
		} else {
			for _, rule := range cmd.disableRule {
				u.disableRule(rule, false)
			}
		}
		return
	}

	if cmd.enableRule != nil {
		if len(cmd.enableRule) == 0 {
			u.enableAllRules(false)
		} else {
			for _, rule := range cmd.enableRule {
				u.enableRule(rule, false)
			}
		}
		return
	}

	if cmd.disableNextLine != nil {
		if len(cmd.enableRule) == 0 {
			u.disableAllRules(true)
		} else {
			for _, rule := range cmd.disableNextLine {
				u.disableRule(rule, true)
			}
		}
		return
	}
}

func (u *JsUnit) enableRule(name string, eph bool) {
	var rule *Rule
	if !eph {
		rule = u.skipped[name]
	} else {
		rule = u.ephSkipped[name]
	}

	if rule == nil {
		return
	}

	for nt, lis := range rule.Listeners {
		walk.AddListener(&u.ana.WalkCtx.Listeners, nt, lis)
	}

	if !eph {
		delete(u.skipped, name)
	} else {
		delete(u.ephSkipped, name)
	}
	u.rules[name] = rule
}

func (u *JsUnit) enableAllRules(eph bool) {
	for name := range u.skipped {
		u.enableRule(name, eph)
	}
}

func (u *JsUnit) disableRule(name string, eph bool) {
	rule := u.rules[name]
	if rule == nil {
		return
	}

	for nt, lis := range rule.Listeners {
		walk.RemoveListener(&u.ana.WalkCtx.Listeners, nt, lis)
	}

	delete(u.rules, name)
	if !eph {
		u.skipped[name] = rule
	} else {
		u.ephSkipped[name] = rule
	}
}

func (u *JsUnit) disableAllRules(eph bool) {
	for name := range u.rules {
		u.disableRule(name, eph)
	}
}

func (u *JsUnit) setupCmtHandles() *JsUnit {
	for nt := range walk.StmtNodeTypes {
		walk.AddNodeBeforeListener(&u.ana.WalkCtx.Listeners, nt, &walk.Listener{
			Id: "_linter_cmt_BeforeStmtListener",
			Handle: func(node parser.Node, key string, ctx *walk.VisitorCtx) {
				if cs := u.parser.PrevStmtCmts(node); cs != nil {
					for _, c := range cs {
						u.handCmt(c)
					}
				}
			},
		})

		walk.AddNodeAfterListener(&u.ana.WalkCtx.Listeners, nt, &walk.Listener{
			Id: "_linter_cmt_AfterStmtListener",
			Handle: func(node parser.Node, key string, ctx *walk.VisitorCtx) {
				for len(u.ephSkipped) > 0 {
					for name := range u.ephSkipped {
						u.enableRule(name, true)
						break
					}
				}
			},
		})
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

	r *Reports
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
		r:         newReports(),
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

func (l *Linter) Config() *Config {
	return l.cfg
}

func (l *Linter) Process() *Reports {
	w := util.NewDirWalker(l.dir, 0, func(f string, dir bool, dw *util.DirWalker) {
		defer func() {
			if r := recover(); r != nil {
				msg := fmt.Sprintf("unexpected error in uint routine, error: %v, src: %v", r, f)
				err := errors.New(msg)
				l.r.addAbnormals(err)
			}
		}()

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

		cfg := l.outerCfg(f)
		if cfg.IsIgnored(f) {
			return
		}

		if isJsFile(f) {
			code, err := ioutil.ReadFile(f)
			if err != nil {
				l.r.addAbnormals(err)
				return
			}

			u, err := NewJsUnit(f, string(code), cfg)
			if err != nil {
				l.r.addAbnormals(err)
				return
			}

			u.linter = l
			u.initRules().setupCmtHandles().enableAllRules(false)
			u.ana.Analyze()
		}
	})

	w.Walk()

	if w.Err() != nil {
		l.r.addAbnormals(w.Err())
	}

	return l.mrkReports()
}

func (l *Linter) mrkReports() *Reports {
	for _, line := range l.diags {
		for _, dig := range line {
			l.r.Diagnoses = append(l.r.Diagnoses, dig...)
		}
	}
	return l.r
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
	// the errors occur in the internal, they maybe caused by some bugs in the internal
	Abnormals []error      `json:"abnormals"`
	Diagnoses []*Diagnosis `json:"diagnoses"`

	abnormalsLock sync.Mutex
}

func newReports() *Reports {
	return &Reports{
		Abnormals: []error{},
		Diagnoses: []*Diagnosis{},

		abnormalsLock: sync.Mutex{},
	}
}

func (r *Reports) addAbnormals(err error) {
	r.abnormalsLock.Lock()
	defer r.abnormalsLock.Unlock()

	r.Abnormals = append(r.Abnormals, err)
}
