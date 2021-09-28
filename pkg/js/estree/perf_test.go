package estree

import (
	"io/ioutil"
	"path"
	"runtime"
	"testing"
)

func BenchmarkParsingToESTree(t *testing.B) {
	libs := []struct {
		name string
		code string
	}{
		{"angular.js", ""},
		{"backbone.js", ""},
		{"ember.js", ""},
		{"jquery.js", ""},
		{"react-dom.js", ""},
		{"react.js", ""},
	}

	_, fileName, _, _ := runtime.Caller(0)
	for _, lib := range libs {
		b, err := ioutil.ReadFile(path.Join(path.Dir(fileName), "assets", lib.name))
		if err != nil {
			t.Fatal(err)
		}
		lib.code = string(b)
	}

	for _, lib := range libs {
		t.Run(lib.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := compile(lib.code)
				if err != nil {
					t.Fatal(err)
				}
			}
		})
	}
}
