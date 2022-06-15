package linter

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"
	"path"
	"path/filepath"
	"plugin"
	"regexp"
	"strings"

	"github.com/go-git/go-git/v5/plumbing/format/gitignore"
	"github.com/hsiaosiyuan0/mole/ecma/parser"
	plg "github.com/hsiaosiyuan0/mole/plugin"
	"github.com/hsiaosiyuan0/mole/util"
)

var builtinRuleFacts = map[string]RuleFact{}

func registerBuiltin(r RuleFact) {
	builtinRuleFacts[r.Name()] = r
}

func init() {
	rules := []RuleFact{
		&NoAlert{},
		&NoUnreachable{},
		&GetterReturn{},
	}

	for _, rule := range rules {
		registerBuiltin(rule)
	}
}

type Config struct {
	Plugins        []string               `json:"plugins"`
	Rules          map[string]interface{} `json:"rules"`
	IgnorePatterns []string               `json:"ignorePatterns"`
	ParserOptions  map[string]interface{} `json:"parserOptions"`

	cwd     string
	plugins map[string]*plugin.Plugin

	ruleFacts     map[string]RuleFact
	rulesCfg      map[string]*RuleCfg
	ruleFactsLang map[string]map[string]RuleFact // lang => rules

	igPatterns []gitignore.Pattern
	matcher    gitignore.Matcher

	outer *Config
}

func NewConfig(cf string, outer *Config) (*Config, error) {
	ext := filepath.Ext(cf)

	var cfg *Config
	var err error

	if ext == ".js" {
		cfg, err = readJsCfg(cf)
		if err != nil {
			return nil, err
		}
	} else {
		cfg, err = readJsonCfg(cf, outer)
		if err != nil {
			return nil, err
		}
	}

	cfg.cwd = path.Dir(cf)
	cfg.plugins = map[string]*plugin.Plugin{}

	cfg.ruleFacts = map[string]RuleFact{}
	cfg.rulesCfg = map[string]*RuleCfg{}
	cfg.ruleFactsLang = map[string]map[string]RuleFact{}

	cfg.IgnorePatterns = append(cfg.IgnorePatterns, "node_modules/")
	cfg.rulesCfg = map[string]*RuleCfg{}

	// inherits plugins and ruleFacts from outer config
	if outer != nil {
		for _, rf := range outer.ruleFacts {
			cfg.AddRuleFact(rf)
		}
		cfg.outer = outer
	}

	if cfg.ParserOptions == nil && cfg.outer != nil {
		cfg.ParserOptions = cfg.outer.ParserOptions
	}

	return cfg, nil
}

func LoadCfgInDir(dir string, outer *Config) (*Config, error) {
	cfg := selCfgFile(dir)
	if cfg == "" {
		return nil, nil
	}
	return NewConfig(cfg, outer)
}

func readJsCfg(cf string) (*Config, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("node", "-e", fmt.Sprintf("var cfg = require('%s'); console.log(JSON.stringify(cfg));", cf))

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	cfg := &Config{}
	if err := json.Unmarshal(stdout.Bytes(), &cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func readJsonCfg(cf string, outer *Config) (*Config, error) {
	raw, err := ioutil.ReadFile(cf)
	if err != nil {
		return nil, err
	}

	cfg := &Config{}
	if err := json.Unmarshal(raw, &cfg); err != nil {
		return nil, err
	}

	cfg.outer = outer
	return cfg, nil
}

func (c *Config) Init() error {
	if err := c.InitIgPatterns(); err != nil {
		return err
	}
	if err := c.InitPlugins(); err != nil {
		return err
	}
	if err := c.InitRulesCfg(); err != nil {
		return err
	}
	return nil
}

func (c *Config) AddRuleFact(rf RuleFact) error {
	name := rf.Name()
	if _, dup := c.ruleFacts[name]; dup {
		return errors.New(fmt.Sprintf("duplicated rule `%s`", name))
	}

	c.ruleFacts[name] = rf
	for _, la := range rf.Meta().Lang {
		if c.ruleFactsLang[la] == nil {
			c.ruleFactsLang[la] = map[string]RuleFact{}
		}
		c.ruleFactsLang[la][name] = rf
	}
	return nil
}

func (c *Config) AddRuleFacts(rfs []RuleFact) error {
	for _, rf := range rfs {
		if err := c.AddRuleFact(rf); err != nil {
			return err
		}
	}
	return nil
}

// the ruleFacts aggregated by their effect langs
func (c *Config) RuleFactsLang() map[string]map[string]RuleFact {
	return c.ruleFactsLang
}

func (c *Config) InitPlugins() error {
	for _, pn := range c.Plugins {
		pg, err := plg.Resolve(c.cwd, pn)
		if err != nil {
			return &LoadConfigErr{fmt.Sprintf("unable to resolve plugin %s", pn), err}
		}
		c.plugins[pn] = pg
		rfs, err := register(pg, pn)
		if err != nil {
			return err
		}
		if err = c.AddRuleFacts(rfs); err != nil {
			return err
		}
	}
	return nil
}

func (c *Config) InitIgPatterns() error {
	ps := c.IgnorePatterns
	psf := path.Join(c.cwd, ".eslintignore")
	if util.FileExist(psf) {
		raw, err := ioutil.ReadFile(psf)
		if err != nil {
			return err
		}
		str := strings.Trim(string(raw), "\r\n")
		ps = append(ps, regexp.MustCompile("\r?\n").Split(str, -1)...)
	}

	domain := strings.Split(c.cwd, string(filepath.Separator))
	c.igPatterns = []gitignore.Pattern{}
	for _, p := range ps {
		pattern := gitignore.ParsePattern(p, domain)
		c.igPatterns = append(c.igPatterns, pattern)
	}

	c.matcher = gitignore.NewMatcher(c.igPatterns)
	return nil
}

type RuleCfg struct {
	Name  string
	Level DiagLevel
	Opts  []interface{}
}

func parseDiagLevel(lvl interface{}) DiagLevel {
	switch lvl.(type) {
	case int:
		switch lvl.(int) {
		case 0:
			return DL_OFF
		case 1:
			return DL_WARN
		case 2:
			return DL_ERROR
		}
	case string:
		switch lvl.(string) {
		case "off":
			return DL_OFF
		case "warn":
			return DL_WARN
		case "error":
			return DL_ERROR
		}
	}
	return DL_ERROR
}

func newRuleCfg(name string, raw interface{}, rf RuleFact) (*RuleCfg, error) {
	c := &RuleCfg{Name: name}

	switch v := raw.(type) {
	case int, string:
		c.Level = parseDiagLevel(raw)

	case []interface{}:
		if len(v) >= 1 {
			c.Level = parseDiagLevel(v[0])

			if len(v) > 1 {
				os := rf.Options()
				if os != nil {
					js, err := json.Marshal(v[1:])
					if err != nil {
						return nil, err
					}
					c.Opts, err = os.ParseOpts(string(js), rf.Validate(), rf.Validates())
					if err != nil {
						return nil, err
					}
				}
			}
		}
	}
	return c, nil
}

func (c *Config) InitRulesCfg() error {
	for name, raw := range c.Rules {
		rf := c.ruleFacts[name]
		if rf == nil {
			continue
		}
		cfg, err := newRuleCfg(name, raw, rf)
		if err != nil {
			return err
		}
		c.rulesCfg[name] = cfg
	}
	return nil
}

func (c *Config) IsIgnored(f string) bool {
	if c.matcher == nil {
		return false
	}
	isDir, _ := util.IsDir(f)
	return c.matcher.Match(strings.Split(f, string(filepath.Separator)), isDir)
}

func (c *Config) ParserOpts() *parser.ParserOpts {
	opts := parser.NewParserOpts()
	if c.ParserOptions != nil {
		opts.MergeJson(c.ParserOptions)
	}
	return opts
}

func (c *Config) CfgOfRule(name string) *RuleCfg {
	return c.rulesCfg[name]
}

func (c *Config) LevelOfRule(name string) DiagLevel {
	rc := c.rulesCfg[name]
	if rc == nil {
		return DL_ERROR
	}
	return rc.Level
}

func selCfgFile(dir string) string {
	cfs := []string{
		".eslintrc",
		".eslintrc.js",
		".eslintrc.json",
	}

	for _, cf := range cfs {
		cf = path.Join(dir, cf)
		if util.FileExist(cf) {
			return cf
		}
	}
	return ""
}

// all the plugins should impl this function
type Register = func() []RuleFact

func register(p *plugin.Plugin, name string) ([]RuleFact, error) {
	e := &LoadConfigErr{fmt.Sprintf("deformed `Register` func in plugin: %s", name), nil}
	fn, err := p.Lookup("Register")
	if err != nil {
		e.cause = err
		return nil, e
	}
	r, ok := fn.(Register)
	if !ok {
		return nil, e
	}
	return r(), nil
}

type LoadConfigErr struct {
	reason string
	cause  error
}

func (e *LoadConfigErr) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(", cause: %s", e.cause.Error())
	}
	return fmt.Sprintf("Failed to load config, reason: %s %s", e.reason, cause)
}
