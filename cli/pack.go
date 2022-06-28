package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/hsiaosiyuan0/mole/ecma/parser"
	"github.com/hsiaosiyuan0/mole/pack"
	"github.com/hsiaosiyuan0/mole/util"
)

type PkgAnalysis struct {
	dir     string
	entries []string
}

type Result struct {
	Elapsed int64 `json:"elapsed"`

	// name => versions
	DupModules map[int64][]string `json:"dupModules"`

	// name+version => [][]Loc
	ImportPoints map[string][][]string `json:"importPoints"`

	Modules map[int64]pack.Module `json:"modules"`

	ParserErrors  []error `json:"parserErrors"`
	ResolveErrors []error `json:"resolveErrors"`
	TimeoutErrors []error `json:"timeoutErrors"`
}

func (a *PkgAnalysis) Process(opts *Options) bool {
	if !opts.packAna {
		return false
	}

	a.dir = opts.dir

	cfg := opts.cfg
	if cfg == "" {
		cfg = "./mole.json"
	}

	cfgFile := filepath.Join(opts.dir, cfg)
	cfgRaw, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		panic(cfgRaw)
	}

	c, err := pack.NewConfig(opts.dir, cfgRaw)
	if err != nil {
		panic(err)
	}

	tsCfg := c.Tsconfig
	if tsCfg == "" {
		if util.FileExist(filepath.Join(a.dir, "tsconfig.json")) {
			c.Tsconfig = "tsconfig.json"
			c.Ts = true
		} else if util.FileExist(filepath.Join(a.dir, "jsconfig.json")) {
			c.Tsconfig = "jsconfig.json"
			c.Ts = false
		}
	}

	packOpts := c.NewDepScannerOpts()
	s := pack.NewDepScanner(packOpts)

	res := &Result{
		DupModules:    map[int64][]string{},
		ImportPoints:  map[string][][]string{},
		ParserErrors:  []error{},
		ResolveErrors: []error{},
	}

	begin := time.Now()
	err = s.Run()
	if err != nil {
		panic(err)
	}

	res.Elapsed = time.Since(begin).Milliseconds()

	// find the dup umbrellas
	umbrellas := s.Umbrellas()
	mvsMap := map[string][]string{} // module name => version names
	dupVs := map[string][]int64{}   // module name => version ids
	dups := []string{}
	for _, m := range umbrellas {
		vs := mvsMap[m.Name()]
		if mvsMap[m.Name()] == nil {
			mvsMap[m.Name()] = []string{}
			vs = mvsMap[m.Name()]

			dupVs[m.Name()] = []int64{}
		}
		ds := dupVs[m.Name()]

		if !util.Includes(vs, m.Version()) {
			mvsMap[m.Name()] = append(vs, m.Version())
			dupVs[m.Name()] = append(ds, m.Id())
			if len(mvsMap[m.Name()]) > 1 {
				dups = append(dups, m.File())
			}
		}
	}

	modules := s.Modules()
	for _, mf := range dups {
		m := umbrellas[mf]
		n := m.Name()

		for _, v := range dupVs[n] {
			sm := modules[v]

			if res.DupModules[m.Id()] == nil {
				res.DupModules[m.Id()] = []string{}
			}

			res.DupModules[m.Id()] = append(res.DupModules[m.Id()], sm.Version())
			res.ImportPoints[sm.Name()+"@"+sm.Version()] = findImportPoints(sm, modules)
		}
	}

	res.Modules = modules

	errs := s.Minors()
	for _, err := range errs {
		switch e := err.(type) {
		case *parser.ParserError, *parser.LexerError:
			res.ParserErrors = append(res.ParserErrors, e)
		case *pack.NoModErr:
			res.ResolveErrors = append(res.ResolveErrors, e)
		case *pack.FileReqTimeoutErr:
			res.TimeoutErrors = append(res.TimeoutErrors, e)
		}
	}

	outData, err := json.Marshal(res)
	if err != nil {
		panic(err)
	}

	out := filepath.Join(a.dir, "mole-pkg-analysis.json")
	err = ioutil.WriteFile(out, outData, 0644)
	if err != nil {
		panic(err)
	}

	return true
}

// find the import points where cause the umbrella being introduced
func findImportPoints(main pack.Module, modules map[int64]pack.Module) [][]string {
	subs := []pack.Module{}
	for _, m := range modules {
		if m == nil {
			continue
		}
		if !m.IsUmbrella() && m.Umbrella() == main.Id() {
			subs = append(subs, m)
		}
	}

	ret := [][]string{}
	for _, c := range subs {
		stk := []string{}
		for _, frame := range c.ImportStk() {
			fm := modules[frame.Mid]
			stk = append(stk, fmt.Sprintf("%s(%d:%d)", fm.File(), frame.Line, frame.Col))
		}
		ret = append(ret, stk)
	}

	return ret
}
