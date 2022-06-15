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

var builtinRuleFacts = map[string]map[string]RuleFact{}

func registerBuiltin(r RuleFact) {
	for _, la := range r.Meta().Lang {
		if builtinRuleFacts[la] == nil {
			builtinRuleFacts[la] = map[string]RuleFact{}
		}
		builtinRuleFacts[la][r.Name()] = r
	}
}

func init() {
	registerBuiltin(&NoAlert{})
	registerBuiltin(&NoUnreachable{})
}

type Config struct {
	Plugins        []string                 `json:"plugins"`
	Rules          map[string][]interface{} `json:"rules"`
	IgnorePatterns []string                 `json:"ignorePatterns"`
	ParserOptions  map[string]interface{}   `json:"parserOptions"`

	cwd       string
	plugins   map[string]*plugin.Plugin
	ruleFacts map[string]map[string]RuleFact // lang => rules

	igPatterns []gitignore.Pattern
	matcher    gitignore.Matcher

	ruleLevels map[string]map[string]DiagLevel

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
	cfg.ruleFacts = map[string]map[string]RuleFact{}
	cfg.IgnorePatterns = append(cfg.IgnorePatterns, "node_modules/")
	cfg.ruleLevels = map[string]map[string]DiagLevel{}

	// inherits plugins and ruleFacts from outer config
	if outer != nil {
		for key, roleFacts := range outer.ruleFacts {
			if cfg.ruleFacts[key] == nil {
				cfg.ruleFacts[key] = map[string]RuleFact{}
			}
			for name, roleFact := range roleFacts {
				cfg.ruleFacts[key][name] = roleFact
			}
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
	return nil
}

func (c *Config) addRuleFact(rfs []RuleFact) error {
	for _, rf := range rfs {
		for _, la := range rf.Meta().Lang {
			if c.ruleFacts[la] == nil {
				c.ruleFacts[la] = map[string]RuleFact{}
			}
			name := rf.Name()
			if _, ok := c.ruleFacts[la][name]; ok {
				return errors.New(fmt.Sprintf("duplicated rule `%s` for lang `%s`", name, la))
			}
			c.ruleFacts[la][name] = rf
		}
	}
	return nil
}

func (c *Config) RuleFact() map[string]map[string]RuleFact {
	return c.ruleFacts
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
		if err = c.addRuleFact(rfs); err != nil {
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

func (c *Config) LevelOfRule(lang string, name string) DiagLevel {
	rules := c.ruleLevels[lang]
	if rules == nil {
		return DL_NONE
	}
	lvl, ok := rules[name]
	if !ok {
		return DL_NONE
	}
	return lvl
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