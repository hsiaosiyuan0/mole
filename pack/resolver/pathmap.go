package resolver

import (
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/hsiaosiyuan0/mole/util"
)

// covers the path-mapping concept used in typescript:
// https://www.typescriptlang.org/docs/handbook/module-resolution.html#path-mapping
type PathMap struct {
	pattern interface{} // string|*regexp.Regexp
	trials  []string
}

// join the `trails` with the given `baseUrl`, the wildcard still in the trails
// ant will be replaced to the matched parts by the latter `Match` calls
func NewPathMap(pat string, baseUrl string, trials []string) (*PathMap, error) {
	p, err := compileSubpath(pat)
	if err != nil {
		return nil, err
	}
	for i, t := range trials {
		if !filepath.IsAbs(t) {
			trials[i] = filepath.Join(baseUrl, t)
		}
	}
	return &PathMap{p, trials}, nil
}

func (m *PathMap) Match(s, cw string, sc *PkgJson, r *ModResolver) *ModResolved {
	mc := false
	var mcs []string
	switch v := m.pattern.(type) {
	case string:
		mc = s == v
	case *regexp.Regexp:
		mcs = v.FindStringSubmatch(s)
		mc = len(mcs) > 0
	}
	if !mc {
		return nil
	}

	for _, d := range m.trials {
		if len(mcs) > 0 {
			d = strings.Replace(d, "*", mcs[1], -1)
		}
		t, err := r.NewTask(d, cw, sc, nil)
		if err != nil {
			continue
		}
		if err := t.resolveFile(); err == nil {
			t.sc = sc
			return t.result()
		}
	}
	return nil
}

type PathMaps struct {
	maps []*PathMap
}

func NewPathMaps(baseUrl string, c map[string][]string) (*PathMaps, error) {
	maps := []*PathMap{}
	keys := util.Keys(c)

	// sort the pattern as desc order by their string length, give the
	// longest pattern the most weight
	sort.Slice(keys, func(i, j int) bool {
		return len(keys[i]) > len(keys[j])
	})

	for _, p := range keys {
		cond := c[p]
		m, err := NewPathMap(p, baseUrl, cond)
		if err != nil {
			return nil, err
		}
		maps = append(maps, m)
	}
	return &PathMaps{maps}, nil
}

func (p *PathMaps) match(s, cw string, sc *PkgJson, r *ModResolver) *ModResolved {
	for _, m := range p.maps {
		mr := m.Match(s, cw, sc, r)
		if mr != nil {
			return mr
		}
	}
	return nil
}

func (p *PathMaps) Match(s, cw string, sc *PkgJson, r *ModResolver) *ModResolved {
	mc := p.match(s, cw, sc, r)
	if mc != nil {
		return mc
	}

	// for `@models` can hit `@models/* => /some/path`, make `@models` to `@models/index`
	// then use the latter to try the rules again
	if mc == nil && strings.HasPrefix(s, "@") && strings.IndexRune(s, '/') == -1 {
		return p.match(s+"/index", cw, sc, r)
	}

	return nil
}
