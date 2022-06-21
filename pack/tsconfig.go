package pack

import (
	"encoding/json"
	"io/ioutil"
	"path"
	"path/filepath"

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
	if path.IsAbs(thePath) {
		return filepath.Join(pathSplit(thePath)...)
	}
	parts := pathSplit(thePath)
	return filepath.Join(append([]string{dir}, parts...)...)
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
		file = path.Join(dir, file)
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

	for cc["extends"] != nil {
		cf := cc["extends"].(string)
		delete(cc, "extends")

		if !filepath.IsAbs(cf) {
			cf = path.Join(cd, cf)
			cd = filepath.Dir(cf)
		}
		raw, err := ioutil.ReadFile(cf)
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

	if !path.IsAbs(c.CompilerOptions.BaseUrl) {
		u := pathSplit(c.CompilerOptions.BaseUrl)
		c.CompilerOptions.BaseUrl = filepath.Join(dir, filepath.Join(u...))
	}

	return c, nil
}

func (c *TsConfig) PathMaps() (*PathMaps, error) {
	if c.pathMaps != nil {
		return c.pathMaps, nil
	}
	var err error
	c.pathMaps, err = NewPathMaps(c.CompilerOptions.BaseUrl, c.CompilerOptions.Paths)
	return c.pathMaps, err
}
