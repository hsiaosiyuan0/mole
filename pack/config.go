package pack

import (
	"encoding/json"
)

type Config struct {
	// web, node, react-native
	Target string `json:"target"`

	Entries                []string               `json:"entries"`
	Extensions             []string               `json:"extensions"`
	DefinedVars            map[string]interface{} `json:"definedVars"`
	Tsconfig               string                 `json:"tsconfig"`
	Ts                     bool                   `json:"ts"`
	ParserOptions          map[string]interface{} `json:"parserOptions"`
	SideEffectsFreeModules []string               `json:"sideEffectsFreeModules"`

	dir string
}

func NewConfig(dir string, raw []byte) (*Config, error) {
	cfg := &Config{dir: dir, Ts: true, ParserOptions: map[string]interface{}{}}
	if err := json.Unmarshal(raw, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) NewDepScannerOpts() *DepScannerOpts {
	o := NewDepScannerOpts()
	o.target = c.Target
	o.Dir = c.dir
	o.Entries = c.Entries
	o.Extensions = c.Extensions

	o.ResolveBuiltin()

	if c.Tsconfig != "" {
		err := o.SetTsconfig(o.Dir, c.Tsconfig, c.Ts)
		if err != nil {
			panic(err)
		}
	}

	o.SerVars(c.DefinedVars)
	o.ParserOpts = c.ParserOptions
	o.FillDefault()

	if c.SideEffectsFreeModules != nil {
		dict := map[string]bool{}
		for _, name := range c.SideEffectsFreeModules {
			dict[name] = true
		}
		o.SideEffectsFreeModules = dict
	}

	return o
}
