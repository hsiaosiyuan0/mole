package estree_test

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/hsiaosiyuan0/mole/pkg/assert"
)

type Fixture struct {
	name   string
	input  string
	output string
}

func scanFixtures(name string) (map[string]*Fixture, error) {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	fxs := make(map[string]*Fixture)
	err := filepath.Walk(path.Join(basepath, "fixtures", name), func(path string, info os.FileInfo, err error) error {
		if err == nil {
			if info.IsDir() {
				return nil
			}
			dir := filepath.Dir(path)
			name := strings.Trim(strings.Replace(dir, basepath, "", 1), string(os.PathSeparator))
			fx := fxs[name]
			if fx == nil {
				fx = &Fixture{name, "", ""}
				fxs[name] = fx
			}
			if strings.HasPrefix(info.Name(), "input") {
				fx.input = path
			} else {
				fx.output = path
			}
		}
		return err
	})
	if err != nil {
		return nil, err
	}
	return fxs, nil
}

func runFixtures(t *testing.T, name string) {
	fxs, err := scanFixtures(name)
	if err != nil {
		t.Fatalf("failed to run fixture [%s] %v", name, err)
	}
	for _, fx := range fxs {
		t.Run(fx.name, func(t *testing.T) {
			code, err := ioutil.ReadFile(fx.input)
			if err != nil {
				t.Fatalf("failed to read fixture code at: %s\nerror: %v", fx.input, err)
			}
			output, err := ioutil.ReadFile(fx.output)
			if err != nil {
				t.Fatalf("failed to read fixture output at: %s\nerror: %v", fx.output, err)
			}
			jsonObj := make(map[string]interface{})
			if err = json.Unmarshal(output, &jsonObj); err != nil {
				t.Fatalf("failed to decode fixture output at: %s\nerror: %v", fx.output, err)
			}

			ast, err := compile(string(code))
			if jsonObj["throws"] != nil {
				if err == nil {
					t.Fatalf("should not pass code:\n%s\nast:\n%v", code, ast)
				}
				assert.Equal(t, jsonObj["throws"].(string), err.Error(), "")
			} else {
				if err != nil {
					t.Fatalf("failed to parse fixture at: %s\nerror: %v", fx.input, err)
				}
				assert.EqualJson(t, string(output), ast)
			}
		})
	}
}

func TestFixtures_es2015(t *testing.T) {
	runFixtures(t, "es2015")
}
