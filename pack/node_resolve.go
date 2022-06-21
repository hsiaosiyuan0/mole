package pack

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"syscall"

	"github.com/hsiaosiyuan0/mole/util"
)

type NodeResolver struct {
	pkgLoader *PkginfoLoader

	self *Pkginfo

	target  string
	cw      string
	imports [][]string
	exports [][]string
	exts    []string
	builtin map[string]bool

	tried []string
}

func NewNodeResolver(exports [][]string, imports [][]string,
	exts []string, builtin map[string]bool, pkginfoLoader *PkginfoLoader) *NodeResolver {

	r := &NodeResolver{
		pkgLoader: pkginfoLoader,

		exports: exports,
		imports: imports,
		exts:    exts,
		builtin: builtin,

		tried: []string{},
	}

	if r.exports == nil {
		r.exports = [][]string{{"node", "require"}}
	}
	if r.imports == nil {
		r.imports = [][]string{{"node", "require"}}
	}

	r.exports = append(r.exports, []string{"default"})
	r.imports = append(r.imports, []string{"default"})

	if r.exts == nil {
		r.exts = jsDefaultExtensions
	}
	return r
}

// checkout: https://nodejs.org/api/modules.html#modules_all_together
//
// `target` and `cw` cannot be empty string
func (r *NodeResolver) Resolve(target string, cw string) (string, error) {
	r.target = target
	r.cw = cw

	if r.isBuiltin(r.target) {
		return "", nil
	}

	parts := pathSplit(r.target)
	cwp := osPathSplit(r.cw)

	c := r.target[0]
	if c == '.' || c == '/' {
		var file []string
		if c == '.' {
			file = append(cwp, parts...)
		} else {
			file = parts
		}

		if f := r.loadAsFile(file); f != "" {
			return f, nil
		}

		f, err := r.loadAsDir(file, true)
		if err != nil {
			return "", err
		}
		return f, nil
	}

	var err error
	if r.self, err = r.pkgLoader.closest(r.cw); err != nil {
		return "", err
	}

	if c == '#' {
		return r.loadPkgImports(parts)
	}

	f, err := r.loadPkgSelf(parts, cwp)
	if err != nil {
		return "", err
	}
	if f != "" {
		return f, nil
	}

	return r.loadModule(parts, cwp)
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

var jsDefaultExtensions = []string{".js", ".json", ".node"}

// target can be either directory or normal file
func (r *NodeResolver) loadAsFile(target []string) string {
	file := filepath.Join(target...)
	r.try(file)

	if filepath.Ext(file) != "" {
		if util.IsFile(file) {
			return file
		}
		return ""
	}

	for _, ext := range r.exts {
		fe := file + ext
		if util.IsFile(fe) {
			return fe
		}
	}
	return ""
}

// target must be a directory
func (r *NodeResolver) loadIndex(target []string) string {
	target = append(target, "index")
	return r.loadAsFile(target)
}

// target must be a directory
func (r *NodeResolver) loadAsDir(target []string, raise bool) (string, error) {
	pki, err := r.pkgLoader.Load(filepath.Join(append(target, "package.json")...))
	if err != nil {
		switch ev := err.(type) {
		case *fs.PathError:
			if ev.Unwrap() != syscall.ENOENT {
				return "", err
			}
		default:
			return "", err
		}
	}

	if pki != nil && pki.Main != "" {
		file := append(target, pki.Main)
		if f := r.loadAsFile(file); f != "" {
			return f, nil
		}
		return "", newNoModErr(r)
	}

	f := r.loadIndex(target)
	if f == "" && raise {
		return "", newNoModErr(r)
	}
	return f, nil
}

func (r *NodeResolver) loadModule(target []string, start []string) (string, error) {
	parts := start

	for len(parts) > 0 {
		dir := []string{}
		pl := len(parts)
		if parts[pl-1] != "node_modules" {
			dir = append(parts, "node_modules")
		}
		parts = parts[:pl-1]

		f, err := r.loadPkgExports(target, dir)
		if err != nil {
			return "", nil
		}
		if f != "" {
			return f, nil
		}

		file := append(dir, target...)
		if f := r.loadAsFile(file); f != "" {
			return f, nil
		}

		f, err = r.loadAsDir(file, false)
		if err != nil {
			return "", err
		}
		if f != "" {
			return f, nil
		}
	}

	return "", newNoModErr(r)
}

func (r *NodeResolver) loadPkgImports(target []string) (string, error) {
	if r.self.imports == nil {
		return "", newNoModErr(r)
	}

	ok, m := r.self.imports.Match(path.Join(target...), r.exports)
	if !ok {
		return "", nil
	}

	file := append(r.self.dir, pathSplit(m)...)
	if f := r.loadAsFile(file); f != "" {
		return f, nil
	}

	f, err := r.loadAsDir(file, true)
	if err != nil {
		return "", err
	}
	return f, nil
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

func (r *NodeResolver) loadPkgSelf(target []string, dir []string) (string, error) {
	if r.self.exports == nil {
		return "", nil
	}

	name, subpath := subpathOf(target)
	if r.self.Name != path.Join(name...) {
		return "", nil
	}

	sp := "."
	if len(subpath) > 0 {
		sp = "./" + path.Join(subpath...)
	}

	// load as file
	if sp == "." || filepath.Ext(sp) != "" {
		ok, m := r.self.exports.Match(sp, r.exports)
		if !ok {
			return "", nil
		}

		file := append(r.self.dir, pathSplit(m)...)
		if f := r.loadAsFile(file); f != "" {
			return f, nil
		}
	} else {
		for _, ext := range r.exts {
			ok, m := r.self.exports.Match(sp+ext, r.exports)
			if !ok {
				return "", nil
			}

			file := append(r.self.dir, pathSplit(m)...)
			if f := r.loadAsFile(file); f != "" {
				return f, nil
			}
		}
	}

	// load as dir
	ok, m := r.self.exports.Match(sp, r.exports)
	if !ok {
		return "", nil
	}
	f, err := r.loadAsDir(append(r.self.dir, pathSplit(m)...), false)
	if err != nil {
		return "", err
	}
	return f, nil
}

func (r *NodeResolver) loadPkgExports(target []string, dir []string) (string, error) {
	name, subpath := subpathOf(target)
	scope := append(dir, name...)

	pki, err := r.pkgLoader.Load(filepath.Join(append(scope, "package.json")...))
	if err != nil {
		switch ev := err.(type) {
		case *fs.PathError:
			if ev.Unwrap() != syscall.ENOENT {
				return "", err
			}
		default:
			return "", err
		}
	}

	if pki == nil || pki.exports == nil {
		return "", nil
	}

	sp := "."
	if len(subpath) > 0 {
		sp = "./" + path.Join(subpath...)
	}

	// load as file
	if sp == "." || filepath.Ext(sp) != "" {
		ok, m := pki.exports.Match(sp, r.exports)
		if !ok {
			return "", nil
		}

		file := append(scope, pathSplit(m)...)
		if f := r.loadAsFile(file); f != "" {
			return f, nil
		}
	} else {
		for _, ext := range r.exts {
			ok, m := pki.exports.Match(sp+ext, r.exports)
			if !ok {
				return "", nil
			}

			file := append(scope, pathSplit(m)...)
			if f := r.loadAsFile(file); f != "" {
				return f, nil
			}
		}
	}

	// load as dir
	ok, m := pki.exports.Match(sp, r.exports)
	if !ok {
		return "", nil
	}
	f, err := r.loadAsDir(append(scope, pathSplit(m)...), false)
	if err != nil {
		return "", err
	}
	return f, nil
}

// only stores the info for module resolution
type Pkginfo struct {
	Name       string                 `json:"name"`
	Main       string                 `json:"main"`
	RawExports interface{}            `json:"exports"`
	RawImports map[string]interface{} `json:"imports"`

	dir     []string
	exports *SubpathGrp
	imports *SubpathGrp
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
		dir: osPathSplit(filepath.Dir(file)),
	}
	if err := json.Unmarshal(code, pi); err != nil {
		return nil, err
	}

	if err := pi.compile(); err != nil {
		return nil, err
	}
	return pi, nil
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