package resolver

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/hsiaosiyuan0/mole/util"
)

// [cjs](https://nodejs.org/api/modules.html#modules_all_together)
// [esm](https://nodejs.org/api/esm.html#resolution-algorithm)
type ModResolver struct {
	browser bool // if use the `browser` spec

	impConds [][]string
	expConds [][]string
	exts     []string
	builtin  map[string]bool

	baseUrl  string
	pathMaps *PathMaps

	pkgLoader *PjsonLoader
}

var DefaultJsExts = []string{".js", ".jsx", ".mjs", ".cjs", ".json", ".node"}
var DefaultTsExts = []string{".ts", ".tsx", ".js", ".jsx", ".mjs", ".d.ts", ".json", ".node"}
var DefaultImpConds = [][]string{{"browser", "require"}}
var DefaultExpConds = [][]string{{"browser", "require"}}

// use this method to ensure the caller wouldn't forget some required options
func NewModResolver(browser bool,
	impConds, expConds [][]string, exts []string, builtin map[string]bool,
	baseUrl string, pathMaps *PathMaps,
	pkgLoader *PjsonLoader) *ModResolver {

	return &ModResolver{
		browser,

		impConds,
		expConds,
		exts,
		builtin,

		baseUrl,
		pathMaps,

		pkgLoader,
	}
}

func (r *ModResolver) NewTask(specifier, cw string, psc, sc *PkgJson) (*ModResolveTask, error) {
	s, err := NewSpecifier(specifier, cw)
	if err != nil {
		return nil, err
	}
	return &ModResolveTask{r, s, cw, psc, sc, "", []string{}}, nil
}

func (t *ModResolver) LookupPkgScope(d string) *PkgJson {
	return t.pkgLoader.LookupPkgScope(d)
}

type ModResolved struct {
	File  string
	Pjson *PkgJson
}

type ModResolveTask struct {
	r   *ModResolver
	sp  *Specifier
	cw  string
	psc *PkgJson
	sc  *PkgJson

	rd    string
	tried []string
}

func (t *ModResolveTask) result() *ModResolved {
	if t.rd == "" {
		return nil
	}
	return &ModResolved{t.rd, t.sc}
}

func (t *ModResolveTask) setResult(r *ModResolved) {
	t.rd = r.File
	t.sc = r.Pjson
}

// return an empty string if specifier represents a builtin module or hits the ignored pattern
func (t *ModResolveTask) Resolve() (*ModResolved, error) {
	if t.isBuiltin() {
		return t.result(), nil
	}

	if t.r.pathMaps != nil {
		if r := t.r.pathMaps.Match(t.sp.raw, t.cw, t.psc, t.r); r != nil {
			return r, nil
		}
	}

	switch t.sp.kind {
	case SPK_FILE:
		err := t.resolveFile()
		if err == nil {
			t.sc = t.psc
		}
		return t.result(), err
	case SPK_BARE:
		err := t.resolvePkg()
		return t.result(), err
	case SPK_IMPT:
		err := t.resolvePkgImports()
		return t.result(), err
	case SPK_NODE, SPK_DATA:
		// `SPK_NODE` is handled by the above `isBuiltin` branch
		// `SPK_DATA` has nothing to do with it
		return t.result(), nil
	default:
		panic("unreachable")
	}
}

func (t *ModResolveTask) newNoModErr() error {
	return newNoModErr(t.cw, t.sp.s, t.tried)
}

func (t *ModResolveTask) resolvePkgImports() error {
	sc := t.lookupPkgScope(t.cw)
	if sc == nil || sc.RawImports == nil {
		return t.newNoModErr()
	}

	pos, neg, m := t.sc.imports.Match(t.sp.s, t.r.expConds)
	if neg {
		return nil
	}
	if pos {
		if t.loadAsFile(filepath.Join(sc.dir, m)) {
			return nil
		}
	}
	return t.newNoModErr()
}

func (t *ModResolveTask) resolvePkg() error {
	if t.resolvePkgSelf(t.sp.s, t.sp.ss) {
		return nil
	}
	if t.loadNodeModules(t.sp.s, t.sp.ss) {
		return nil
	}
	return t.newNoModErr()
}

var sep = string(filepath.Separator)

func (t *ModResolveTask) loadNodeModules(pkgName, subpath string) bool {
	if t.sc != nil && t.sc.exports != nil {
		_, neg, _ := t.sc.exports.Match(pkgName, t.r.expConds)
		if neg {
			return true
		}
	}

	spans := strings.Split(t.cw, "/")
	for i := len(spans) - 1; i >= 0; i-- {
		if spans[i] == "" {
			break
		}

		dir := spans[0 : i+1]
		if spans[i] != "node_modules" {
			dir = append(dir, "node_modules")
		}

		d := filepath.Join(dir...)
		d = filepath.Join(sep, d, pkgName)

		tt, err := t.r.NewTask(subpath, d, t.sc, nil)
		if err != nil {
			continue
		}

		if tt.resolvePkgExports(subpath) {
			t.setResult(tt.result())
			return true
		}

		d = filepath.Join(d, subpath)
		if tt.loadAsFile(d) {
			t.setResult(tt.result())
			return true
		}
		if tt.loadAsDir(d) == nil {
			r := tt.result()
			if r.Pjson.Private {
				r.Pjson = t.lookupPkgScope(r.Pjson.dir)
			}
			t.setResult(r)
			return true
		}
	}
	return false
}

func (t *ModResolveTask) resolvePkgSelf(pkgName, subpath string) bool {
	sc := t.lookupPkgScope(t.cw)
	if sc == nil {
		return false
	}

	if sc.Name != pkgName {
		return false
	}

	if sc.exports == nil || sc.exports.IsEmpty() {
		return false
	}

	if subpath == "." && sc.onlyMain && len(sc.main) > 0 {
		for _, main := range sc.main {
			f := filepath.Join(sc.dir, main)
			if t.loadAsFile(f) {
				return true
			}
			return false
		}
	}

	return t.resolvePkgExports(subpath)
}

func (t *ModResolveTask) resolvePkgExports(subpath string) bool {
	if t.lookupPkgScope(t.cw) == nil {
		return false
	}

	if t.sc.exports == nil {
		return false
	}

	pos, neg, m := t.sc.exports.Match(subpath, t.r.expConds)
	if neg {
		return true
	}
	if pos {
		return t.loadAsFile(filepath.Join(t.sc.dir, m))
	}

	for _, ext := range t.r.exts {
		pos, neg, m := t.sc.exports.Match(subpath+ext, t.r.expConds)
		if neg {
			return true
		}
		if pos {
			return t.loadAsFile(filepath.Join(t.sc.dir, m))
		}
	}

	return false
}

func (t *ModResolveTask) isBuiltin() bool {
	if t.sp.kind == SPK_NODE {
		return true
	}
	if t.sp.kind != SPK_BARE {
		return false
	}
	if t.r.builtin == nil {
		return false
	}
	return t.r.builtin[t.sp.s]
}

func (t *ModResolveTask) resolveFile() error {
	if t.loadAsFile(t.sp.s) {
		return nil
	}

	return t.loadAsDir(t.sp.s)
}

func (t *ModResolveTask) loadAsFile(s string) bool {
	t.trace(s)

	if util.IsFile(s) {
		t.rd = s
		return true
	}

	// if the cjs compatible mode is off then stop trying the exts
	if filepath.Ext(s) != "" {
		return false
	}

	for _, ext := range t.r.exts {
		f := s + ext
		t.trace(f)
		if util.IsFile(f) {
			t.rd = f
			return true
		}
	}

	return false
}

// try to resolve as directory in cjs compatible mode
func (t *ModResolveTask) loadAsDir(s string) error {
	if util.IsDir(s) {
		sc, err := t.r.pkgLoader.Load(filepath.Join(s, "package.json"))
		if err != nil {
			goto LOAD_AS_INDEX
		}

		if len(sc.main) == 0 {
			goto LOAD_AS_INDEX
		}

		for _, main := range sc.main {
			main = filepath.Join(s, main)
			if t.loadAsFile(main) {
				t.sc = sc
				return nil
			}

			err = t.loadIndex(main)
			if err == nil {
				return nil
			}
		}
	}

LOAD_AS_INDEX:
	return t.loadIndex(s)
}

func (t *ModResolveTask) loadIndex(s string) error {
	s = filepath.Join(s, "index")
	if t.loadAsFile(s) {
		return nil
	}

	return t.newNoModErr()
}

func (t *ModResolveTask) trace(s string) {
	t.tried = append(t.tried, s)
}

func (t *ModResolveTask) lookupPkgScope(d string) *PkgJson {
	if t.sc != nil {
		return t.sc
	}
	t.sc = t.r.pkgLoader.LookupPkgScope(d)
	return t.sc
}

type NoModErr struct {
	Cw        string
	Specifier string
	Tried     []string
}

func newNoModErr(cw, specifier string, tried []string) *NoModErr {
	return &NoModErr{cw, specifier, tried}
}

func (e *NoModErr) Error() string {
	return fmt.Sprintf(`Module Not Found, specifier: "%s", tried paths:\n%s`, e.Specifier, strings.Join(e.Tried, "\n"))
}

type FileExtErr struct {
	Specifier string
	Msg       string
}

func newFileExtErr(specifier, msg string) *FileExtErr {
	return &FileExtErr{specifier, msg}
}

func (e *FileExtErr) Error() string {
	return fmt.Sprintf(`Unsupported File Extension, reason: %s specifier: "%s"`, e.Msg, e.Specifier)
}
