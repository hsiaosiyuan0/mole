package pack

import (
	"encoding/json"
)

type Config struct {
	Entries       []string               `json:"entries"`
	DefinedVars   map[string]interface{} `json:"definedVars"`
	Tsconfig      string                 `json:"tsconfig"`
	Ts            bool                   `json:"ts"`
	ParserOptions map[string]interface{} `json:"parserOptions"`

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
	o.Dir = c.dir
	o.Entries = c.Entries

	err := o.SetTsconfig(o.Dir, c.Tsconfig, c.Ts)
	if err != nil {
		panic(err)
	}

	o.SerVars(c.DefinedVars)
	o.ParserOpts = c.ParserOptions

	return o
}
