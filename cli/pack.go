package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/hsiaosiyuan0/mole/ecma/parser"
	"github.com/hsiaosiyuan0/mole/pack"
	"github.com/hsiaosiyuan0/mole/pack/resolver"
	"github.com/hsiaosiyuan0/mole/util"
)

type PkgAnalysis struct {
	dir     string
	entries []string
	out     string
}

type DupVersion struct {
	Id      int64  `json:"id"`
	Version string `json:"version"`
	Size    int64  `json:"size"`
}

type DupItem struct {
	Name     string        `json:"name"`
	Size     int64         `json:"size"`
	Versions []*DupVersion `json:"versions"`
}

type DupItems []*DupItem

func (d DupItems) Len() int           { return len(d) }
func (d DupItems) Less(i, j int) bool { return d[i].Size > d[j].Size }
func (d DupItems) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }

func (d *DupItem) addVersion(m *pack.Module) {
	for _, v := range d.Versions {
		if v.Id == m.Id() {
			return
		}
	}
	v := &DupVersion{m.Id(), m.Version(), m.Size()}
	d.Versions = append(d.Versions, v)
	d.Size += v.Size
}

type ImportPoint struct {
	File   string   `json:"file"`
	Reason []string `json:"reason"`
}

type ImportInfo struct {
	IncludePath  [][]int64      `json:"includePath"`
	ImportPoints []*ImportPoint `json:"importPoints"`
}

type Result struct {
	Elapsed int64 `json:"elapsed"`

	// name => versions
	DupModules []*DupItem `json:"dupModules"`

	// name+moduleId => []*ImportPoint
	ImportInfo map[string]*ImportInfo `json:"importInfo"`

	Modules map[int64]*pack.Module `json:"modules"`

	ParserErrors  []error `json:"parserErrors"`
	ResolveErrors []error `json:"resolveErrors"`
	TimeoutErrors []error `json:"timeoutErrors"`
}

func (a *PkgAnalysis) Process(opts *Options) bool {
	if !opts.packAna {
		return false
	}

	a.dir = opts.dir
	a.out = opts.out

	cfg := opts.cfg
	if cfg == "" {
		cfg = "./mole.json"
	}

	cfgFile := filepath.Join(opts.dir, cfg)
	cfgRaw, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		panic(err)
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
		DupModules:    []*DupItem{},
		ImportInfo:    map[string]*ImportInfo{},
		ParserErrors:  []error{},
		ResolveErrors: []error{},
		TimeoutErrors: []error{},
	}

	begin := time.Now()
	err = s.ResolveDeps()
	if err != nil {
		panic(err)
	}

	s.DCE()

	res.Elapsed = time.Since(begin).Milliseconds()

	// find the dup umbrellas
	umbrellas := s.Umbrellas()
	modules := s.Modules()

	moduleIds := map[string][]int64{} // module name => ids
	dups := []string{}
	for _, m := range umbrellas {
		id := m.Id()
		name := m.Name()

		ids := moduleIds[name]
		if !util.Includes(ids, id) {
			moduleIds[name] = append(moduleIds[name], id)
			cnt := len(moduleIds[name])
			if cnt > 1 {
				if cnt == 2 {
					dups = append(dups, modules[moduleIds[name][0]].File())
				}
				dups = append(dups, m.File())
			}
		}
	}

	dupItemsMap := map[string]*DupItem{}
	for _, mf := range dups {
		m := umbrellas[mf]
		n := m.Name()

		for _, id := range moduleIds[n] {
			sm := modules[id]

			dupItem := dupItemsMap[m.Name()]
			if dupItem == nil {
				dupItem = &DupItem{m.Name(), 0, []*DupVersion{}}
				dupItemsMap[m.Name()] = dupItem
			}

			dupItem.addVersion(sm)
			res.ImportInfo[sm.Name()+"@"+strconv.Itoa(int(sm.Id()))] = resolveImportInfo(sm, modules)
		}
	}

	dupItems := []*DupItem{}
	for _, item := range dupItemsMap {
		dupItems = append(dupItems, item)
	}
	sort.Sort(DupItems(dupItems))

	res.DupModules = dupItems
	res.Modules = modules

	errs := s.Minors()
	for _, err := range errs {
		switch e := err.(type) {
		case *parser.ParserError, *parser.LexerError:
			res.ParserErrors = append(res.ParserErrors, e)
		case *resolver.NoModErr, *resolver.FileExtErr,
			*resolver.InvalidPkgCfgErr, *resolver.InvalidSpecifierErr:
			res.ResolveErrors = append(res.ResolveErrors, e)
		case *pack.FileReqTimeoutErr:
			res.TimeoutErrors = append(res.TimeoutErrors, e)
		}
	}

	outData, err := json.Marshal(res)
	if err != nil {
		panic(err)
	}

	out := a.out
	if out == "" {
		out = filepath.Join(a.dir, fmt.Sprintf("mole-pkg-analysis-%d.json", time.Now().Unix()))
	}

	if util.FileExist(out) {
		panic(fmt.Sprintf("Output file `%s` already exists, abort to overwrite it", out))
	}

	err = ioutil.WriteFile(out, outData, 0644)
	if err != nil {
		panic(err)
	}

	return true
}

func umbrellasOfFrames(c int64, frames []*pack.ImportFrame, modules map[int64]*pack.Module) ([]int64, string) {
	ret := []int64{}
	key := []string{}
	for _, frame := range frames {
		um := modules[modules[frame.Mid].Umbrella()]
		uid := um.Id()
		cnt := len(ret)
		if cnt == 0 || ret[cnt-1] != uid {
			ret = append(ret, uid)
			key = append(key, strconv.Itoa(int(uid)))
		}
	}
	if ret[len(ret)-1] != c {
		ret = append(ret, c)
		key = append(key, strconv.Itoa(int(c)))
	}
	return ret, strings.Join(key, "-")
}

// find out the import points which cause the umbrella being introduced
func resolveImportInfo(main *pack.Module, modules map[int64]*pack.Module) *ImportInfo {
	subs := []*pack.Module{}
	for _, m := range modules {
		if m == nil {
			continue
		}
		if !m.IsUmbrella() && m.Umbrella() == main.Id() {
			subs = append(subs, m)
		}
	}

	points := []*ImportPoint{}
	unique := map[string]bool{}
	includePaths := [][]int64{}
	for _, c := range subs {
		frames := c.ImportStk()
		ums, key := umbrellasOfFrames(c.Umbrella(), frames, modules)
		if !unique[key] {
			unique[key] = true
			includePaths = append(includePaths, ums)
		}

		stk := []string{}
		for _, frame := range frames {
			fm := modules[frame.Mid]
			loc := frame.S.OfstLineCol(frame.Rng.Lo)
			stk = append(stk, fmt.Sprintf("%s(%d:%d)", fm.File(), loc.Line, loc.Col))
		}
		points = append(points, &ImportPoint{c.File(), stk})
	}

	return &ImportInfo{includePaths, points}
}
