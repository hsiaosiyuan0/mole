package lint

import (
	"fmt"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/hsiaosiyuan0/mole/ecma/analysis"
	"github.com/hsiaosiyuan0/mole/ecma/parser"
	"github.com/hsiaosiyuan0/mole/ecma/walk"
	"github.com/hsiaosiyuan0/mole/span"
)

type Unit interface {
	Lang() string
	Config() *Config
	HasCommentInSpan(span.Range) bool
	Report(*Diagnosis)
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
	u.ana = analysis.NewAnalysis(ast, u.parser.Symtab(), p.Source())

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

var hasCmtReg = regexp.MustCompile("(?s)//.*?\n|/\\*.*?\\*/")

func (u *JsUnit) HasCommentInSpan(rng span.Range) bool {
	str := u.parser.RngText(rng)
	return hasCmtReg.FindIndex([]byte(str)) != nil
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
	for _, rf := range u.cfg.ruleFactsLang[lang] {
		ctx := &RuleCtx{u.parser.Source(), u, rf}
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

var startCmtReg = regexp.MustCompile(`^/\*\s*|\s*\*/`)
var cmtReg = regexp.MustCompile(`^//\s*|\s*$`)

func (u *JsUnit) parseCmt(cs span.Range) *CmtCmd {
	cmt := u.parser.RngText(cs)
	if strings.HasPrefix(cmt, "/*") {
		cmt = startCmtReg.ReplaceAllString(cmt, "")
	} else {
		cmt = cmtReg.ReplaceAllString(cmt, "")
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

func (u *JsUnit) handleCmt(c span.Range) {
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
				if cs := u.parser.PrevCmts(node); cs != nil {
					for _, c := range cs {
						u.handleCmt(c)
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
