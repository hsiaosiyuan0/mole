package resolver

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
	"sync"
	"syscall"

	"github.com/hsiaosiyuan0/mole/util"
)

// only stores the info for module resolution
type PkgJson struct {
	Name        string                 `json:"name"`
	Version     string                 `json:"version"`
	Private     bool                   `json:"private"`
	Main        string                 `json:"main"`
	Browser     interface{}            `json:"browser"`
	Module      string                 `json:"module"`
	Type        string                 `json:"type"`
	Types       string                 `json:"types"`
	RawExports  interface{}            `json:"exports"`
	RawImports  map[string]interface{} `json:"imports"`
	SideEffects interface{}            `json:"sideEffects"`

	// the filesystem location of the pjson
	file string
	dir  string

	// the compiled subpath of the exports and imports declared in pjson
	exports *SubpathGrp
	imports *SubpathGrp

	onlyMain bool
	main     []string
}

func (pj *PkgJson) IsSideEffectsFree() bool {
	switch v := pj.SideEffects.(type) {
	case bool:
		return v == false
	case map[string]interface{}:
		return v != nil
	}
	return false
}

func (pj *PkgJson) Dir() string {
	return pj.dir
}

func (pj *PkgJson) File() string {
	return pj.file
}

func (pj *PkgJson) compile(raw []byte, browser bool) error {
	var err error

	if err := json.Unmarshal(raw, pj); err != nil {
		return &InvalidPkgCfgErr{pj.file, err.Error()}
	}

	if pj.RawImports != nil {
		pj.imports, err = NewSubpathGrp(pj.RawImports)
		if err != nil {
			return err
		}
	}

	// spec says that if `browser` is a string then it will take the place of the `main`
	main := pj.Main
	if browser && pj.Browser != nil {
		if s, ok := pj.Browser.(string); ok {
			main = s
		}
	}

	pj.main = []string{}
	if pj.Module != "" {
		pj.main = append(pj.main, pj.Module)
	}

	if main != "" {
		pj.main = append(pj.main, main)
	}

	if pj.RawExports == nil && main != "" && (pj.Browser == nil || pj.Browser == main) {
		pj.onlyMain = true
		return nil
	}

	// set an empty value for normalizing and merge the `main`
	var rawExports map[string]interface{}
	if pj.RawExports != nil {
		rawExports, err = NormalizeSubpath(pj.RawExports)
		if err != nil {
			return err
		}
	} else {
		rawExports = map[string]interface{}{}
	}

	if main != "" {
		mainCond, err := NormalizeSubpath(main)
		if err != nil {
			return err
		}

		util.MergeMap(rawExports, mainCond)
	}

	// do the partial replacement if `browser` is a map
	if browser && pj.Browser != nil {
		if bro, ok := pj.Browser.(map[string]interface{}); ok {
			browserFirst := func(path []string, val interface{}, key interface{}, parent interface{}, arr bool) bool {
				if s, ok := val.(string); ok {
					if s[0] != '.' {
						s = "./" + s
					}
					if bro[s] != nil {
						if arr {
							parent.([]interface{})[key.(int)] = bro[s]
						} else {
							parent.(map[string]interface{})[key.(string)] = bro[s]
						}
					}
				}
				return true
			}
			util.WalkObj(rawExports, make([]string, 0), browserFirst, nil, nil, false)

			// merge the ignored settings
			for k, v := range bro {
				switch v.(type) {
				case nil, bool:
					rawExports[k] = v
				}
			}
		}
	}

	pj.exports, err = NewSubpathGrp(rawExports)
	if err != nil {
		return err
	}

	return nil
}

type PjsonLoader struct {
	loader *FileLoader

	store map[string]*PkgJson
	lock  sync.RWMutex

	// cache the paths which does not have `package.json` in them
	notFound     map[string]bool
	notFoundLock sync.RWMutex

	// use the `browser` instead of the `main` in `package.json`
	browser bool
}

func NewPjsonLoader(fl *FileLoader) *PjsonLoader {
	return &PjsonLoader{
		loader: fl,

		store: map[string]*PkgJson{},
		lock:  sync.RWMutex{},

		notFound:     map[string]bool{},
		notFoundLock: sync.RWMutex{},
	}
}

// directly get info from cache
func (pl *PjsonLoader) Get(file string) *PkgJson {
	pl.lock.RLock()
	defer pl.lock.RUnlock()

	return pl.store[file]
}

func (pl *PjsonLoader) SetBrowser(browser bool) {
	pl.browser = browser
}

func (pl *PjsonLoader) SetNotFound(file string, err error) {
	pl.notFoundLock.Lock()
	defer pl.notFoundLock.Unlock()

	switch ev := err.(type) {
	case *fs.PathError:
		if ev.Unwrap() == syscall.ENOENT {
			pl.notFound[file] = true
		}
	}
}

func (pl *PjsonLoader) IsNotFound(file string) bool {
	pl.notFoundLock.RLock()
	defer pl.notFoundLock.RUnlock()

	_, ok := pl.notFound[file]
	return ok
}

func (pl *PjsonLoader) Load(file string) (*PkgJson, error) {
	if c := pl.Get(file); c != nil {
		return c, nil
	}

	fr := <-pl.loader.Load(file) // wait
	if fr.Err != nil {
		pl.SetNotFound(file, fr.Err)
		return nil, fr.Err
	}

	pj, err := pl.compile(file, fr.Raw)
	if err != nil {
		return nil, err
	}

	pl.lock.Lock()
	pl.store[file] = pj
	pl.lock.Unlock()
	return pj, nil
}

func (pl *PjsonLoader) compile(file string, raw []byte) (*PkgJson, error) {
	pj := &PkgJson{
		file: file,
		dir:  filepath.Dir(file),
	}

	if err := pj.compile(raw, pl.browser); err != nil {
		return nil, err
	}
	return pj, nil
}

func (pl *PjsonLoader) LookupPkgScope(start string) *PkgJson {
	for {
		if start == "/" || strings.HasSuffix(start, "/node_modules") {
			break
		}

		file := filepath.Join(start, "package.json")

		if pl.IsNotFound(file) {
			break
		}

		pi, err := pl.Load(file)
		// some modules use internal package.json to redirect user's imports, if that internal
		// package.json was tagged as `private` then use its outer scope instead
		if err == nil && !(pi.Private && strings.Index(file, "/node_modules") != -1) {
			return pi
		}

		start = filepath.Dir(start)
	}

	return nil
}

type InvalidPkgCfgErr struct {
	File string
	Msg  string
}

func (e *InvalidPkgCfgErr) Error() string {
	return fmt.Sprintf("Invalid package configuration in %s, reason %s", e.File, e.Msg)
}
