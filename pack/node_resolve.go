package pack

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/hsiaosiyuan0/mole/util"
)

type NodeResolver struct {
	pkgLoader *PkginfoLoader

	// both `LOAD_PACKAGE_IMPORTS(X, DIR)` and `LOAD_PACKAGE_SELF(X, DIR)` has a nerd strategy
	// described as `Find the closest package scope SCOPE to DIR`, which means each part of
	// the `DIR` should be taken in account to found out if there is a `package.json` in it.
	//
	// consider below example:
	//
	// ```js
	// // current file `a-package/b/c/d.mjs`
	// import "a-package/b/c/d.mjs"
	// ```
	//
	// above import semantic will be handled by the `LOAD_PACKAGE_SELF(X, DIR)` function, then
	// `a-package/b/c` and `a-package/b` are needed to be checked in `order` to find out if there
	// is a `package.json` under one of them, obviously it will lost the performance
	//
	// a trade-off is made in mole after considering the resolving-performance and spec-compatibility,
	// its key is the `pki` field, the caller specify the `approximately` closest scope via this field
	// instead of bubbling up the directory parts
	//
	// let's still base on above example, the value of `pki` should be the instance of the `package.json`
	// under the directory of the module `a-package`, expressed like `a-package/package.json`, in other words
	// the processes to of `a-package/b/c/package.json` and `a-package/b/package.json` are skipped, so if
	// there is really either a-package/b/c/package.json` or `a-package/b/package.json` then the result is
	// a mistake, that's why the `a-package/package.json` is called a `approximately` closest scope
	pki *Pkginfo

	target  string
	cw      string
	imports [][]string
	exports [][]string
	exts    map[string]bool
	builtin map[string]bool

	tried []string
}

func NewNodeResolver(pki *Pkginfo, target string, cw string, exports [][]string, imports [][]string,
	exts map[string]bool, builtin map[string]bool, pkginfoLoader *PkginfoLoader) *NodeResolver {

	r := &NodeResolver{
		pkgLoader: pkginfoLoader,

		pki:     pki,
		target:  target,
		cw:      cw,
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
func (r *NodeResolver) Resolve() (string, error) {
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

var jsDefaultExtensions = map[string]bool{
	".js":   true,
	".json": true,
	".node": true,
}

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

	for ext := range r.exts {
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
	pkg := filepath.Join(filepath.Join(target...), "package.json")
	pki, err := r.pkgLoader.Load(pkg)
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
	if r.pki.imports == nil {
		return "", newNoModErr(r)
	}

	ok, m := r.pki.imports.Match(path.Join(target...), r.exports)
	if !ok {
		return "", nil
	}

	file := append(r.pki.dir, pathSplit(m)...)
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
	if r.pki.exports == nil {
		return "", nil
	}

	name, subpath := subpathOf(target)
	if r.pki.Name != path.Join(name...) {
		return "", nil
	}

	sp := "."
	if len(subpath) > 0 {
		sp = "./" + path.Join(subpath...)
	}

	// load as file
	if sp == "." || filepath.Ext(sp) != "" {
		ok, m := r.pki.exports.Match(sp, r.exports)
		if !ok {
			return "", nil
		}

		file := append(r.pki.dir, pathSplit(m)...)
		if f := r.loadAsFile(file); f != "" {
			return f, nil
		}
	} else {
		for ext := range r.exts {
			ok, m := r.pki.exports.Match(sp+ext, r.exports)
			if !ok {
				return "", nil
			}

			file := append(r.pki.dir, pathSplit(m)...)
			if f := r.loadAsFile(file); f != "" {
				return f, nil
			}
		}
	}

	// load as dir
	ok, m := r.pki.exports.Match(sp, r.exports)
	if !ok {
		return "", nil
	}
	f, err := r.loadAsDir(append(r.pki.dir, pathSplit(m)...), false)
	if err != nil {
		return "", err
	}
	return f, nil
}

func (r *NodeResolver) loadPkgExports(target []string, dir []string) (string, error) {
	name, subpath := subpathOf(target)
	scope := append(dir, name...)

	pkg := filepath.Join(append(scope, "package.json")...)
	pki, err := r.pkgLoader.Load(pkg)
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
		for ext := range r.exts {
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
	// path => info
	cache *util.LruCache[string, *Pkginfo]
}

func NewPkginfoLoader(cap int, clear int) *PkginfoLoader {
	return &PkginfoLoader{
		cache: util.NewLruCache[string, *Pkginfo](cap, clear),
	}
}

func (lo *PkginfoLoader) Load(file string) (*Pkginfo, error) {
	if lo.cache.HasKey(file) {
		return lo.cache.Get(file), nil
	}

	code, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	pi := &Pkginfo{
		dir: osPathSplit(filepath.Dir(file)),
	}
	if err := json.Unmarshal(code, pi); err != nil {
		return nil, err
	}

	if err := pi.compile(); err != nil {
		return nil, err
	}

	lo.cache.Set(file, pi)
	return pi, nil
}

type NoModErr struct {
	Target string
	Cw     string
	Exts   []string
	Tried  []string
}

func (m *NoModErr) Error() string {
	return fmt.Sprintf("failed to load `%s` in `%s` with exts `%v`, tried these paths:\n %s", m.Target, m.Cw, m.Exts, strings.Join(m.Tried, "\n"))
}

func newNoModErr(r *NodeResolver) *NoModErr {
	keys := []string{}
	for key := range r.exts {
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
