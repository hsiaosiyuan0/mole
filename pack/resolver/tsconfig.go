package resolver

import (
	"encoding/json"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"

	"github.com/hsiaosiyuan0/mole/util"
)

type TscOptions struct {
	Module string `json:"module"`
	Target string `json:"target"`

	ModuleResolution string              `json:"moduleResolution"`
	BaseUrl          string              `json:"baseUrl"`
	Paths            map[string][]string `json:"paths"`
	RootDirs         []string            `json:"rootDirs"`
}

// represents both `tsconfig` and `jsconfig`, prefix `Ts` is used instead of `Js` since
// this structure is developed by microsoft and firstly used in vscode
type TsConfig struct {
	Extends string `json:"extends"`

	CompilerOptions TscOptions `json:"compilerOptions"`

	Include []string `json:"include"`
	Exclude []string `json:"exclude"`
	Files   []string `json:"files"`

	pathMaps *PathMaps
}

func fixRelativePath(thePath, dir string) string {
	if filepath.IsAbs(thePath) {
		return thePath
	}
	return filepath.Join(dir, thePath)
}

func fixRelativePaths(thePath interface{}, dir string) interface{} {
	switch v := thePath.(type) {
	case string:
		return fixRelativePath(v, dir)
	case []interface{}:
		ret := make([]interface{}, len(v))
		for i, p := range v {
			ret[i] = fixRelativePath(p.(string), dir)
		}
		return ret
	case []string:
		ret := make([]interface{}, len(v))
		for i, p := range v {
			ret[i] = fixRelativePath(p, dir)
		}
		return ret
	case map[string]interface{}:
		ret := map[string][]interface{}{}
		for k, ps := range v {
			ret[k] = fixRelativePaths(ps, dir).([]interface{})
		}
		return ret
	}
	return nil
}

func fixRelativePathInCfg(c map[string]interface{}, dir string) {
	c["include"] = fixRelativePaths(c["include"], dir)
	c["exclude"] = fixRelativePaths(c["exclude"], dir)
	c["files"] = fixRelativePaths(c["files"], dir)

	if opts := c["compilerOptions"]; opts != nil {
		if c, ok := opts.(map[string]interface{}); ok {
			c["baseUrl"] = fixRelativePaths(c["baseUrl"], dir)
			c["paths"] = fixRelativePaths(c["paths"], dir)
			c["rootDirs"] = fixRelativePaths(c["rootDirs"], dir)
		}
	}
}

func NewTsConfig(dir string, file string) (*TsConfig, error) {
	if !filepath.IsAbs(file) {
		file = filepath.Join(dir, file)
	}

	rawCfg := map[string]interface{}{}

	raw, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	raw, err = util.RemoveJsonComments(string(raw))
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(raw, &rawCfg); err != nil {
		return nil, err
	}

	supers := []map[string]interface{}{}
	cc := rawCfg
	cd := dir

	// load the base config

	fileLoader := NewFileLoader(2048, 128)
	pkgLoader := NewPjsonLoader(fileLoader)
	r := NewModResolver(true, nil, nil, DefaultJsExts, nil, "", nil, pkgLoader)

	for cc["extends"] != nil {
		extFile := cc["extends"].(string)
		delete(cc, "extends")

		if strings.HasPrefix(extFile, ".") || strings.HasPrefix(extFile, "/") {
			extFile = filepath.Join(cd, extFile)
		} else {

			t, err := r.NewTask(extFile, cd, nil, nil)
			if err != nil {
				return nil, err
			}

			rd, err := t.Resolve()
			if err != nil {
				return nil, err
			}
			extFile = rd.File
		}

		cd = filepath.Dir(extFile)
		raw, err := ioutil.ReadFile(extFile)
		if err != nil {
			return nil, err
		}

		sc := map[string]interface{}{}
		if err := json.Unmarshal(raw, &sc); err != nil {
			return nil, err
		}

		// `references` is excluded from inheritance
		delete(sc, "references")
		fixRelativePathInCfg(sc, cd)
		supers = append(supers, sc)

		cc = sc
	}

	util.MergeMaps(rawCfg, supers...)
	js, err := json.Marshal(rawCfg)
	if err != nil {
		return nil, err
	}

	c := &TsConfig{}
	if err := json.Unmarshal(js, &c); err != nil {
		return nil, err
	}

	bu := c.CompilerOptions.BaseUrl
	if bu != "" && !path.IsAbs(bu) {
		c.CompilerOptions.BaseUrl = filepath.Join(dir, bu)
	}

	c.pathMaps, err = NewPathMaps(c.CompilerOptions.BaseUrl, c.CompilerOptions.Paths)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c *TsConfig) PathMaps() *PathMaps {
	return c.pathMaps
}
