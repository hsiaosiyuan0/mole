package pack

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"syscall"

	"github.com/hsiaosiyuan0/mole/util"
)

// the main algorithm refers: https://nodejs.org/api/modules.html#modules_all_together
type NodeResolver struct {
	pkgLoader *PkginfoLoader

	self *Pkginfo

	target string
	cw     string

	imports [][]string
	exports [][]string
	exts    []string
	builtin map[string]bool

	ts       bool
	baseUrl  string
	pathMaps *PathMaps

	tried []string
}

func NewNodeResolver(exports [][]string, imports [][]string, exts []string,
	builtin map[string]bool, pkginfoLoader *PkginfoLoader, ts bool, pathMaps *PathMaps) *NodeResolver {

	r := &NodeResolver{
		pkgLoader: pkginfoLoader,

		exports: exports,
		imports: imports,
		exts:    exts,
		builtin: builtin,

		ts:       ts,
		pathMaps: pathMaps,

		tried: []string{},
	}

	if len(r.exports) == 0 {
		r.exports = [][]string{{"node", "require"}}
	}
	if len(r.imports) == 0 {
		r.imports = [][]string{{"node", "require"}}
	}
	if len(r.builtin) == 0 {
		r.builtin = jsBuiltin
	}

	r.exports = append(r.exports, []string{"default"})
	r.imports = append(r.imports, []string{"default"})

	if r.exts == nil {
		if r.ts {
			r.exts = tsDefaultExtensions
		} else {
			r.exts = jsDefaultExtensions
		}
	}
	return r
}

// `target` and `cw` cannot be empty string.
//
// because of the resolution strategy in typescript also grab the `.d.ts` declared in
// the `types` field of the `package.json`, so the first return value is `[]string`
// instead of a bare `string`
//
// `*Pkginfo` is parent module of the resolved files
func (r *NodeResolver) Resolve(target string, cw string) ([]string, *Pkginfo, error) {
	r.target = target
	r.cw = cw

	if r.isBuiltin(r.target) {
		return nil, nil, nil
	}

	parts := pathSplit(r.target)
	cwp := osPathSplit(r.cw)

	var err error
	if r.self, err = r.pkgLoader.closest(r.cw); err != nil {
		return nil, nil, err
	}

	c := r.target[0]
	if c == '.' || c == '/' {
		return r.loadRelative(parts, cwp, r.self)
	}

	if r.pathMaps != nil {
		if f := r.pathMaps.Match(r.target, r); f != nil {
			return f, r.self, nil
		}
	}

	if c == '#' {
		return r.loadPkgImports(parts, r.self)
	}

	f, pi, err := r.loadPkgSelf(parts, cwp, r.self)
	if err != nil {
		return nil, nil, err
	}
	if len(f) != 0 {
		return f, pi, nil
	}

	return r.loadModule(parts, cwp)
}

func (r *NodeResolver) ResolveRoots(target string, cws []string) ([]string, *Pkginfo, error) {
	ure := &NoModUnderRootsErr{[]*NoModErr{}}
	for _, cw := range cws {
		f, pi, err := r.Resolve(target, cw)
		if err != nil {
			if e, ok := err.(*NoModErr); ok {
				ure.errs = append(ure.errs, e)
				continue
			}
			return nil, nil, err
		} else {
			return f, pi, nil
		}
	}
	return nil, nil, ure
}

func (r *NodeResolver) isBuiltin(target string) bool {
	if r.builtin == nil {
		return false
	}
	_, ok := r.builtin[target]
	return ok
}

func (r *NodeResolver) try(target string) {
	r.tried = append(r.tried, target)
}

var jsDefaultExtensions = []string{".js", ".jsx", ".mjs", ".json", ".node"}
var tsDefaultExtensions = append([]string{".ts", ".tsx", ".d.ts"}, jsDefaultExtensions...)

var jsBuiltin = map[string]bool{
	"assert":         true,
	"buffer":         true,
	"child_process":  true,
	"cluster":        true,
	"crypto":         true,
	"dgram":          true,
	"dns":            true,
	"domain":         true,
	"events":         true,
	"fs":             true,
	"http":           true,
	"https":          true,
	"net":            true,
	"os":             true,
	"path":           true,
	"querystring":    true,
	"readline":       true,
	"stream":         true,
	"string_decoder": true,
	"timers":         true,
	"tls":            true,
	"tty":            true,
	"url":            true,
	"util":           true,
	"v8":             true,
	"vm":             true,
	"zlib":           true,
}

// target can be either directory or normal file
func (r *NodeResolver) loadAsFile(target []string) (string, *Pkginfo) {
	file := filepath.Join(target...)
	r.try(file)

	f := ""
	if filepath.Ext(file) != "" {
		if util.IsFile(file) {
			f = file
		}
	}

	if f == "" {
		for _, ext := range r.exts {
			fe := file + ext
			if util.IsFile(fe) {
				f = fe
			}
		}
	}

	if f != "" {
		pi, _ := r.pkgLoader.closest(filepath.Join(target...))
		return f, pi
	}

	return "", nil
}

func (r *NodeResolver) loadRelative(target []string, cw []string, pi *Pkginfo) ([]string, *Pkginfo, error) {
	var file []string
	if target[0][0] == '.' {
		file = append(cw, target...)
	} else {
		file = target
	}

	if f, pi := r.loadAsFile(file); f != "" {
		return []string{f}, pi, nil
	}

	f, pi, err := r.loadAsDir(file, true, true)
	if err != nil {
		return nil, nil, err
	}
	return f, pi, nil
}

// target must be a directory
func (r *NodeResolver) loadIndex(target []string) (string, *Pkginfo) {
	target = append(target, "index")
	return r.loadAsFile(target)
}

// target must be a directory
func (r *NodeResolver) loadAsDir(target []string, raise bool, skipPkgInfo bool) ([]string, *Pkginfo, error) {
	var pki *Pkginfo
	var err error
	if !skipPkgInfo {
		pki, err = r.pkgLoader.Load(filepath.Join(append(target, "package.json")...))
		if err != nil {
			switch ev := err.(type) {
			case *fs.PathError:
				if ev.Unwrap() != syscall.ENOENT {
					return nil, nil, err
				}
			default:
				return nil, nil, err
			}
		}
	}

	res := []string{}
	if pki != nil {
		if pki.Main != "" {
			file := append(target, pki.Main)
			if f, _ := r.loadAsFile(file); f != "" {
				res = append(res, f)
			} else if f, _, _ := r.loadAsDir(file, false, true); len(f) > 0 {
				res = append(res, f...)
			} else {
				return nil, nil, newNoModErr(r)
			}
		}
		if r.ts && pki.Types != "" {
			file := append(target, pki.Types)
			if f, _ := r.loadAsFile(file); f != "" {
				res = append(res, f)
			} else {
				return nil, nil, newNoModErr(r)
			}
		}
		if len(res) > 0 {
			return res, pki, nil
		}
	}

	f, pki := r.loadIndex(target)
	if f == "" {
		if raise {
			return nil, nil, newNoModErr(r)
		}
		return nil, nil, nil
	}

	res = append(res, f)
	return res, pki, nil
}

func (r *NodeResolver) loadModule(target []string, start []string) ([]string, *Pkginfo, error) {
	parts := util.Copy(start)

	for len(parts) > 0 {
		dir := parts
		pl := len(parts)
		if parts[pl-1] != "node_modules" {
			dir = append(parts, "node_modules")
		}
		parts = parts[:pl-1]

		f, pi, err := r.loadPkgExports(target, dir)
		if err != nil {
			return nil, nil, err
		}
		if len(f) != 0 {
			return f, pi, nil
		}

		file := append(dir, target...)
		if f, pi := r.loadAsFile(file); f != "" {
			return []string{f}, pi, nil
		}

		f, pi, err = r.loadAsDir(file, false, false)
		if err != nil {
			return nil, nil, err
		}
		if len(f) != 0 {
			return f, pi, nil
		}
	}

	if r.ts && r.baseUrl != "" {
		prefix := osPathSplit(r.baseUrl)
		target := append(prefix, target...)

		f, pi := r.loadAsFile(target)
		if f != "" {
			return []string{f}, pi, nil
		}
		return r.loadAsDir(target, true, false)
	}

	return nil, nil, newNoModErr(r)
}

func (r *NodeResolver) loadPkgImports(target []string, pi *Pkginfo) ([]string, *Pkginfo, error) {
	if r.self.imports == nil {
		return nil, nil, newNoModErr(r)
	}

	ok, m := r.self.imports.Match(path.Join(target...), r.exports)
	if !ok {
		return nil, nil, nil
	}

	file := append(r.self.dir, pathSplit(m)...)
	if f, pi := r.loadAsFile(file); f != "" {
		return []string{f}, pi, nil
	}

	f, pi, err := r.loadAsDir(file, true, false)
	if err != nil {
		return nil, nil, err
	}
	return f, pi, nil
}

func subpathOf(target []string) ([]string, []string) {
	name := target[0:1]
	subpath := target[1:]
	if target[0][0] == '@' && len(target) > 1 {
		name = name[0:2]
		subpath = target[2:]
	}
	return name, subpath
}

func (r *NodeResolver) loadPkgSelf(target []string, dir []string, pi *Pkginfo) ([]string, *Pkginfo, error) {
	if r.self.exports == nil {
		return nil, nil, nil
	}

	name, subpath := subpathOf(target)
	if r.self.Name != path.Join(name...) {
		return nil, nil, nil
	}

	sp := "."
	if len(subpath) > 0 {
		sp = "./" + path.Join(subpath...)
	}

	// load as file
	if sp == "." || filepath.Ext(sp) != "" {
		ok, m := r.self.exports.Match(sp, r.exports)
		if !ok {
			return nil, nil, nil
		}

		file := append(r.self.dir, pathSplit(m)...)
		if f, pi := r.loadAsFile(file); f != "" {
			return []string{f}, pi, nil
		}
	} else {
		for _, ext := range r.exts {
			ok, m := r.self.exports.Match(sp+ext, r.exports)
			if !ok {
				return nil, nil, nil
			}

			file := append(r.self.dir, pathSplit(m)...)
			if f, pi := r.loadAsFile(file); f != "" {
				return []string{f}, pi, nil
			}
		}
	}

	// load as dir
	ok, m := r.self.exports.Match(sp, r.exports)
	if !ok {
		return nil, nil, nil
	}
	f, pi, err := r.loadAsDir(append(r.self.dir, pathSplit(m)...), false, false)
	if err != nil {
		return nil, nil, err
	}
	return f, pi, nil
}

func (r *NodeResolver) loadPkgExports(target []string, dir []string) ([]string, *Pkginfo, error) {
	name, subpath := subpathOf(target)
	scope := append(dir, name...)

	pki, err := r.pkgLoader.Load(filepath.Join(append(scope, "package.json")...))
	if err != nil {
		switch ev := err.(type) {
		case *fs.PathError:
			if ev.Unwrap() != syscall.ENOENT {
				return nil, nil, err
			}
		default:
			return nil, nil, err
		}
	}

	if pki == nil || pki.exports == nil {
		return nil, nil, nil
	}

	sp := "."
	if len(subpath) > 0 {
		sp = "./" + path.Join(subpath...)
	}

	// load as file
	if sp == "." || filepath.Ext(sp) != "" {
		ok, m := pki.exports.Match(sp, r.exports)
		if !ok {
			return nil, nil, nil
		}

		file := append(scope, pathSplit(m)...)
		if f, pki := r.loadAsFile(file); f != "" {
			return []string{f}, pki, nil
		}
	} else {
		for _, ext := range r.exts {
			ok, m := pki.exports.Match(sp+ext, r.exports)
			if !ok {
				return nil, nil, nil
			}

			file := append(scope, pathSplit(m)...)
			if f, pki := r.loadAsFile(file); f != "" {
				return []string{f}, pki, nil
			}
		}
	}

	// load as dir
	ok, m := pki.exports.Match(sp, r.exports)
	if !ok {
		return nil, nil, nil
	}
	f, pki, err := r.loadAsDir(append(scope, pathSplit(m)...), false, false)
	if err != nil {
		return nil, nil, err
	}
	return f, pki, nil
}

// only stores the info for module resolution
type Pkginfo struct {
	Name       string                 `json:"name"`
	Version    string                 `json:"version"`
	Main       string                 `json:"main"`
	Types      string                 `json:"types"`
	RawExports interface{}            `json:"exports"`
	RawImports map[string]interface{} `json:"imports"`

	file    string
	dir     []string
	exports *SubpathGrp
	imports *SubpathGrp
}

func (pi *Pkginfo) Dir() []string {
	return pi.dir
}

func (pi *Pkginfo) File() string {
	return pi.file
}

func (pi *Pkginfo) compile() error {
	var err error
	if pi.RawExports != nil {
		pi.exports, err = NewSubpathGrp(pi.RawExports)
		if err != nil {
			return err
		}
	}
	if pi.RawImports != nil {
		pi.imports, err = NewSubpathGrp(pi.RawImports)
		if err != nil {
			return err
		}
	}
	return nil
}

type PkginfoLoader struct {
	loader *FileLoader

	store map[string]*Pkginfo
	lock  sync.RWMutex

	// cache the path where `package.json` does not exist
	notFound     map[string]bool
	notFoundLock sync.RWMutex
}

func NewPkginfoLoader(fl *FileLoader) *PkginfoLoader {
	if fl == nil {
		fl = NewFileLoader(1024, 10)
	}

	return &PkginfoLoader{
		loader: fl,

		store: map[string]*Pkginfo{},
		lock:  sync.RWMutex{},

		notFound:     map[string]bool{},
		notFoundLock: sync.RWMutex{},
	}
}

// directly get info from cache
func (lo *PkginfoLoader) Get(file string) *Pkginfo {
	lo.lock.RLock()
	defer lo.lock.RUnlock()

	return lo.store[file]
}

func (lo *PkginfoLoader) setNotFound(file string, err error) {
	lo.notFoundLock.Lock()
	defer lo.notFoundLock.Unlock()

	switch ev := err.(type) {
	case *fs.PathError:
		if ev.Unwrap() == syscall.ENOENT {
			lo.notFound[file] = true
		}
	}
}

func (lo *PkginfoLoader) isNotFound(file string) bool {
	lo.lock.RLock()
	defer lo.lock.RUnlock()

	return lo.notFound[file] == true
}

func (lo *PkginfoLoader) Load(file string) (*Pkginfo, error) {
	if hit := lo.Get(file); hit != nil {
		return hit, nil
	}

	f, err := lo.loader.Load(file)
	if err != nil {
		return nil, err
	}

	var pi *Pkginfo
	switch fv := f.(type) {
	case []byte: // done
		if pi, err = lo.compile(file, fv); err != nil {
			lo.setNotFound(file, err)
			return nil, err
		}
	case chan *FileLoadResult:
		f := <-fv // wait
		if f.err != nil {
			return nil, f.err
		}
		if pi, err = lo.compile(file, f.raw); err != nil {
			lo.setNotFound(file, err)
			return nil, err
		}
	default:
		panic("unreachable")
	}

	lo.lock.Lock()
	lo.store[file] = pi
	lo.lock.Unlock()
	return pi, nil
}

func (lo *PkginfoLoader) closest(start string) (*Pkginfo, error) {
	for {
		if strings.HasSuffix(start, "node_modules") {
			break
		}
		file := filepath.Join(start, "package.json")
		if lo.isNotFound(file) {
			break
		}
		if pi, err := lo.Load(file); err == nil {
			return pi, nil
		}
		start = filepath.Dir(start)
	}
	return nil, errors.New("failed to find the closest scope from " + start)
}

func (lo *PkginfoLoader) compile(file string, code []byte) (*Pkginfo, error) {
	pi := &Pkginfo{
		file: file,
		dir:  osPathSplit(filepath.Dir(file)),
	}
	if err := json.Unmarshal(code, pi); err != nil {
		return nil, err
	}

	if err := pi.compile(); err != nil {
		return nil, err
	}
	return pi, nil
}

type NoModUnderRootsErr struct {
	errs []*NoModErr
}

func (m *NoModUnderRootsErr) Error() string {
	sb := strings.Builder{}
	for _, e := range m.errs {
		sb.WriteString(e.Error() + "\n")
	}
	return sb.String()
}

type NoModErr struct {
	Target string
	Cw     string
	Exts   []string
	Tried  []string
}

func (m *NoModErr) Error() string {
	return fmt.Sprintf("failed to load `%s` under `%s` with exts `%v`, tried paths:\n%s", m.Target, m.Cw, m.Exts, strings.Join(m.Tried, "\n"))
}

func newNoModErr(r *NodeResolver) *NoModErr {
	keys := []string{}
	for _, key := range r.exts {
		keys = append(keys, key)
	}
	return &NoModErr{r.target, r.cw, keys, r.tried}
}

var osSep = string(filepath.Separator)

func pathSplit(f string) []string {
	parts := strings.Split(f, "/")
	if f[0] == '/' {
		parts[0] = osSep
	}
	return parts
}

func osPathSplit(f string) []string {
	parts := strings.Split(f, osSep)
	if f[0] == '/' {
		parts[0] = osSep
	}
	return parts
}

// covers the path-mapping concept used in typescript:
// https://www.typescriptlang.org/docs/handbook/module-resolution.html#path-mapping
type PathMap struct {
	pat  interface{} // string|*regexp.Regexp
	cond []string
}

func NewPathMap(pat string, baseUrl string, cond []string) (*PathMap, error) {
	p, err := compileSubpath(pat)
	if err != nil {
		return nil, err
	}
	for i, c := range cond {
		if !path.IsAbs(c) {
			c = path.Join(baseUrl, c)
		}
		cond[i] = filepath.Join(pathSplit(c)...)
	}
	return &PathMap{p, cond}, nil
}

func (m *PathMap) Match(nom string, r *NodeResolver) []string {
	mc := false
	var mcs []string
	switch v := m.pat.(type) {
	case string:
		mc = nom == v
	case *regexp.Regexp:
		mcs = v.FindStringSubmatch(nom)
		mc = len(mcs) > 0
	}
	if !mc {
		return nil
	}

	for _, d := range m.cond {
		d = strings.Replace(d, "*", mcs[1], -1)
		if f, _, _ := r.loadRelative(osPathSplit(d), nil, nil); len(f) != 0 {
			return f
		}
	}
	return nil
}

type PathMaps struct {
	maps []*PathMap
}

func NewPathMaps(baseUrl string, c map[string][]string) (*PathMaps, error) {
	maps := []*PathMap{}
	for p, cond := range c {
		m, err := NewPathMap(p, baseUrl, cond)
		if err != nil {
			return nil, err
		}
		maps = append(maps, m)
	}
	return &PathMaps{maps}, nil
}

func (p *PathMaps) Match(file string, r *NodeResolver) []string {
	for _, m := range p.maps {
		mr := m.Match(file, r)
		if mr != nil {
			return mr
		}
	}
	return nil
}
