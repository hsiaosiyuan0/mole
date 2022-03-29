package estree_test

import (
	"testing"

	. "github.com/hsiaosiyuan0/mole/ecma/estree/test"
	. "github.com/hsiaosiyuan0/mole/util"
)

func TestExportAllAsNS1(t *testing.T) {
	ast, err := Compile("export * as ns from \"source\"")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 28,
  "body": [
    {
      "type": "ExportAllDeclaration",
      "start": 0,
      "end": 28,
      "exported": {
        "type": "Identifier",
        "start": 12,
        "end": 14,
        "name": "ns"
      },
      "source": {
        "type": "Literal",
        "start": 20,
        "end": 28,
        "value": "source",
        "raw": "\"source\""
      }
    }
  ]
}
	`, ast)
}

func TestExportAllAsNS2(t *testing.T) {
	ast, err := Compile("export * as foo from \"bar\"")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 26,
  "body": [
    {
      "type": "ExportAllDeclaration",
      "start": 0,
      "end": 26,
      "exported": {
        "type": "Identifier",
        "start": 12,
        "end": 15,
        "name": "foo"
      },
      "source": {
        "type": "Literal",
        "start": 21,
        "end": 26,
        "value": "bar",
        "raw": "\"bar\""
      }
    }
  ]
}
	`, ast)
}

func TestExportAllAsNS3(t *testing.T) {
	ast, err := Compile("export * from \"source\"")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 22,
  "body": [
    {
      "type": "ExportAllDeclaration",
      "start": 0,
      "end": 22,
      "exported": null,
      "source": {
        "type": "Literal",
        "start": 14,
        "end": 22,
        "value": "source",
        "raw": "\"source\""
      }
    }
  ]
}
	`, ast)
}
