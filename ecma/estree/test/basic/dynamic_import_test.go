package estree_test

import (
	"testing"

	. "github.com/hsiaosiyuan0/mole/ecma/estree/test"
	. "github.com/hsiaosiyuan0/mole/util"
)

func TestDynamicImport1(t *testing.T) {
	ast, err := Compile("import('dynamicImport.js')")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 26,
  "body": [
    {
      "type": "ExpressionStatement",
      "start": 0,
      "end": 26,
      "expression": {
        "type": "ImportExpression",
        "start": 0,
        "end": 26,
        "source": {
          "type": "Literal",
          "start": 7,
          "end": 25,
          "value": "dynamicImport.js",
          "raw": "'dynamicImport.js'"
        }
      }
    }
  ]
}
`, ast)
}

func TestDynamicImport2(t *testing.T) {
	ast, err := Compile("import(a = 'dynamicImport.js')")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 30,
  "body": [
    {
      "type": "ExpressionStatement",
      "start": 0,
      "end": 30,
      "expression": {
        "type": "ImportExpression",
        "start": 0,
        "end": 30,
        "source": {
          "type": "AssignmentExpression",
          "start": 7,
          "end": 29,
          "operator": "=",
          "left": {
            "type": "Identifier",
            "start": 7,
            "end": 8,
            "name": "a"
          },
          "right": {
            "type": "Literal",
            "start": 11,
            "end": 29,
            "value": "dynamicImport.js",
            "raw": "'dynamicImport.js'"
          }
        }
      }
    }
  ]
}
`, ast)
}

func TestDynamicImport3(t *testing.T) {
	ast, err := Compile("new (import(s))")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 15,
  "body": [
    {
      "type": "ExpressionStatement",
      "start": 0,
      "end": 15,
      "expression": {
        "type": "NewExpression",
        "start": 0,
        "end": 15,
        "callee": {
          "type": "ImportExpression",
          "start": 5,
          "end": 14,
          "source": {
            "type": "Identifier",
            "start": 12,
            "end": 13,
            "name": "s"
          }
        },
        "arguments": []
      }
    }
  ]
}
`, ast)
}

func TestDynamicImport4(t *testing.T) {
	ast, err := Compile("import((s,t))")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 13,
  "body": [
    {
      "type": "ExpressionStatement",
      "start": 0,
      "end": 13,
      "expression": {
        "type": "ImportExpression",
        "start": 0,
        "end": 13,
        "source": {
          "type": "SequenceExpression",
          "start": 8,
          "end": 11,
          "expressions": [
            {
              "type": "Identifier",
              "start": 8,
              "end": 9,
              "name": "s"
            },
            {
              "type": "Identifier",
              "start": 10,
              "end": 11,
              "name": "t"
            }
          ]
        }
      }
    }
  ]
}
`, ast)
}
