package pack

import (
	"encoding/json"
	"io/ioutil"
	"path"
	"path/filepath"

	"github.com/hsiaosiyuan0/mole/util"
)

type TscOptions struct {
	ModuleResolution string              `json:"moduleResolution"`
	BaseUrl          string              `json:"baseUrl"`
	Paths            map[string][]string `json:"paths"`

	Module string `json:"module"`
	Target string `json:"target"`
}

// represents both `tsconfig` and `jsconfig`, prefix `Ts` is used instead of `Js` since
// this structure is developed by microsoft and firstly used in vscode
type TsConfig struct {
	Extends string `json:"extends"`

	CompilerOptions TscOptions `json:"compilerOptions"`

	Include []string `json:"include"`
	Exclude []string `json:"exclude"`
	Files   []string `json:"files"`
}

func NewTsConfig(dir string, file string) (*TsConfig, error) {
	if !path.IsAbs(file) {
		file = path.Join(dir, file)
	}

	rawCfg := map[string]interface{}{}

	raw, err := ioutil.ReadFile(file)
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

		if !path.IsAbs(cf) {
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
		supers = append(supers, sc)

		cc = sc
	}

	util.MergeMaps(rawCfg, supers...)
	js, err := json.Marshal(rawCfg)
	if err != nil {
		return nil, err
	}

	ret := &TsConfig{}
	if err := json.Unmarshal(js, &ret); err != nil {
		return nil, err
	}

	return ret, nil
}
