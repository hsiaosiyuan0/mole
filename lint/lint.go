package lint

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sync"

	"github.com/hsiaosiyuan0/mole/util"
)

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
	}

	if !skipBuiltin {
		// inherits ruleFacts from builtin
		for _, rf := range builtinRuleFacts {
			cfg.AddRuleFact(rf)
		}
	}

	if err = cfg.Init(); err != nil {
		return nil, err
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
	line := dig.Loc.Begin().Line

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
