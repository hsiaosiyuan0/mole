package estree

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"

	"github.com/hsiaosiyuan0/mole/pkg/assert"
	"github.com/hsiaosiyuan0/mole/pkg/js/parser"
)

func compile(code string) (parser.Node, error) {
	s := parser.NewSource("", code)
	p := parser.NewParser(s, make([]string, 0))
	return p.Prog()
}

func TestWorks(t *testing.T) {
	ast, err := compile(`
const {unlink} = require('fs');

unlink('/tmp/hello', (err, data) => {
  if (err) throw err;
  console.log('successfully deleted /tmp/hello');
});
`)
	assert.Equal(t, nil, err, "should be prog ok")

	b, err := json.Marshal(NewProgram(ast.(*parser.Prog)))
	if err != nil {
		t.Fail()
	}
	var out bytes.Buffer
	json.Indent(&out, b, "", "  ")
	out.WriteTo(os.Stdout)
}
