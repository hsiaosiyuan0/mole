package estree_test

import (
	"testing"

	. "github.com/hsiaosiyuan0/mole/ecma/estree/test"
	"github.com/hsiaosiyuan0/mole/ecma/parser"
	. "github.com/hsiaosiyuan0/mole/internal"
)

func TestEs7th1(t *testing.T) {
	ast, err := Compile("x **= 42")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 8,
  "body": [
    {
      "type": "ExpressionStatement",
      "start": 0,
      "end": 8,
      "expression": {
        "type": "AssignmentExpression",
        "start": 0,
        "end": 8,
        "operator": "**=",
        "left": {
          "type": "Identifier",
          "start": 0,
          "end": 1,
          "name": "x"
        },
        "right": {
          "type": "Literal",
          "start": 6,
          "end": 8,
          "value": 42,
          "raw": "42"
        }
      }
    }
  ]
}
`, ast)
}

func TestEs7th2(t *testing.T) {
	ast, err := Compile("x ** y")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 6,
  "body": [
    {
      "type": "ExpressionStatement",
      "start": 0,
      "end": 6,
      "expression": {
        "type": "BinaryExpression",
        "start": 0,
        "end": 6,
        "left": {
          "type": "Identifier",
          "start": 0,
          "end": 1,
          "name": "x"
        },
        "operator": "**",
        "right": {
          "type": "Identifier",
          "start": 5,
          "end": 6,
          "name": "y"
        }
      }
    }
  ]
}
`, ast)
}

func TestEs7th3(t *testing.T) {
	ast, err := Compile("3 ** 5 * 1")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 10,
  "body": [
    {
      "type": "ExpressionStatement",
      "start": 0,
      "end": 10,
      "expression": {
        "type": "BinaryExpression",
        "start": 0,
        "end": 10,
        "left": {
          "type": "BinaryExpression",
          "start": 0,
          "end": 6,
          "left": {
            "type": "Literal",
            "start": 0,
            "end": 1,
            "value": 3,
            "raw": "3"
          },
          "operator": "**",
          "right": {
            "type": "Literal",
            "start": 5,
            "end": 6,
            "value": 5,
            "raw": "5"
          }
        },
        "operator": "*",
        "right": {
          "type": "Literal",
          "start": 9,
          "end": 10,
          "value": 1,
          "raw": "1"
        }
      }
    }
  ]
}
`, ast)
}

func TestEs7th4(t *testing.T) {
	ast, err := Compile("3 % 5 ** 1")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 10,
  "body": [
    {
      "type": "ExpressionStatement",
      "start": 0,
      "end": 10,
      "expression": {
        "type": "BinaryExpression",
        "start": 0,
        "end": 10,
        "left": {
          "type": "Literal",
          "start": 0,
          "end": 1,
          "value": 3,
          "raw": "3"
        },
        "operator": "%",
        "right": {
          "type": "BinaryExpression",
          "start": 4,
          "end": 10,
          "left": {
            "type": "Literal",
            "start": 4,
            "end": 5,
            "value": 5,
            "raw": "5"
          },
          "operator": "**",
          "right": {
            "type": "Literal",
            "start": 9,
            "end": 10,
            "value": 1,
            "raw": "1"
          }
        }
      }
    }
  ]
}
`, ast)
}

func TestEs7th5(t *testing.T) {
	ast, err := Compile("-a * 5")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 6,
  "body": [
    {
      "type": "ExpressionStatement",
      "start": 0,
      "end": 6,
      "expression": {
        "type": "BinaryExpression",
        "start": 0,
        "end": 6,
        "left": {
          "type": "UnaryExpression",
          "start": 0,
          "end": 2,
          "operator": "-",
          "prefix": true,
          "argument": {
            "type": "Identifier",
            "start": 1,
            "end": 2,
            "name": "a"
          }
        },
        "operator": "*",
        "right": {
          "type": "Literal",
          "start": 5,
          "end": 6,
          "value": 5,
          "raw": "5"
        }
      }
    }
  ]
}
`, ast)
}

func TestEs7th6(t *testing.T) {
	ast, err := Compile("(-5) ** y")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 9,
  "body": [
    {
      "type": "ExpressionStatement",
      "start": 0,
      "end": 9,
      "expression": {
        "type": "BinaryExpression",
        "start": 0,
        "end": 9,
        "left": {
          "type": "UnaryExpression",
          "start": 1,
          "end": 3,
          "operator": "-",
          "prefix": true,
          "argument": {
            "type": "Literal",
            "start": 2,
            "end": 3,
            "value": 5,
            "raw": "5"
          }
        },
        "operator": "**",
        "right": {
          "type": "Identifier",
          "start": 8,
          "end": 9,
          "name": "y"
        }
      }
    }
  ]
}
`, ast)
}

func TestEs7th7(t *testing.T) {
	ast, err := Compile("++a ** 2")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 8,
  "body": [
    {
      "type": "ExpressionStatement",
      "start": 0,
      "end": 8,
      "expression": {
        "type": "BinaryExpression",
        "start": 0,
        "end": 8,
        "left": {
          "type": "UpdateExpression",
          "start": 0,
          "end": 3,
          "operator": "++",
          "prefix": true,
          "argument": {
            "type": "Identifier",
            "start": 2,
            "end": 3,
            "name": "a"
          }
        },
        "operator": "**",
        "right": {
          "type": "Literal",
          "start": 7,
          "end": 8,
          "value": 2,
          "raw": "2"
        }
      }
    }
  ]
}
`, ast)
}

func TestEs7th8(t *testing.T) {
	ast, err := Compile("a-- ** 2")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 8,
  "body": [
    {
      "type": "ExpressionStatement",
      "start": 0,
      "end": 8,
      "expression": {
        "type": "BinaryExpression",
        "start": 0,
        "end": 8,
        "left": {
          "type": "UpdateExpression",
          "start": 0,
          "end": 3,
          "operator": "--",
          "prefix": false,
          "argument": {
            "type": "Identifier",
            "start": 0,
            "end": 1,
            "name": "a"
          }
        },
        "operator": "**",
        "right": {
          "type": "Literal",
          "start": 7,
          "end": 8,
          "value": 2,
          "raw": "2"
        }
      }
    }
  ]
}
`, ast)
}

func TestEs7th9(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_STRICT)
	ast, err := CompileWithOpts("if (x) function f() {}", opts)
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 22,
  "body": [
    {
      "type": "IfStatement",
      "start": 0,
      "end": 22,
      "test": {
        "type": "Identifier",
        "start": 4,
        "end": 5,
        "name": "x"
      },
      "consequent": {
        "type": "FunctionDeclaration",
        "start": 7,
        "end": 22,
        "id": {
          "type": "Identifier",
          "start": 16,
          "end": 17,
          "name": "f"
        },
        "generator": false,
        "async": false,
        "params": [],
        "body": {
          "type": "BlockStatement",
          "start": 20,
          "end": 22,
          "body": []
        }
      },
      "alternate": null
    }
  ]
}
`, ast)
}

func TestEs7th10(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_STRICT)
	ast, err := CompileWithOpts("if (x) function f() { return 23; } else function f() { return 42; }", opts)
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 67,
  "body": [
    {
      "type": "IfStatement",
      "start": 0,
      "end": 67,
      "test": {
        "type": "Identifier",
        "start": 4,
        "end": 5,
        "name": "x"
      },
      "consequent": {
        "type": "FunctionDeclaration",
        "start": 7,
        "end": 34,
        "id": {
          "type": "Identifier",
          "start": 16,
          "end": 17,
          "name": "f"
        },
        "generator": false,
        "async": false,
        "params": [],
        "body": {
          "type": "BlockStatement",
          "start": 20,
          "end": 34,
          "body": [
            {
              "type": "ReturnStatement",
              "start": 22,
              "end": 32,
              "argument": {
                "type": "Literal",
                "start": 29,
                "end": 31,
                "value": 23,
                "raw": "23"
              }
            }
          ]
        }
      },
      "alternate": {
        "type": "FunctionDeclaration",
        "start": 40,
        "end": 67,
        "id": {
          "type": "Identifier",
          "start": 49,
          "end": 50,
          "name": "f"
        },
        "generator": false,
        "async": false,
        "params": [],
        "body": {
          "type": "BlockStatement",
          "start": 53,
          "end": 67,
          "body": [
            {
              "type": "ReturnStatement",
              "start": 55,
              "end": 65,
              "argument": {
                "type": "Literal",
                "start": 62,
                "end": 64,
                "value": 42,
                "raw": "42"
              }
            }
          ]
        }
      }
    }
  ]
}
`, ast)
}

func TestEs7th11(t *testing.T) {
	ast, err := Compile("function foo(a) { 'use strict'; }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 33,
  "body": [
    {
      "type": "FunctionDeclaration",
      "start": 0,
      "end": 33,
      "id": {
        "type": "Identifier",
        "start": 9,
        "end": 12,
        "name": "foo"
      },
      "generator": false,
      "async": false,
      "params": [
        {
          "type": "Identifier",
          "start": 13,
          "end": 14,
          "name": "a"
        }
      ],
      "body": {
        "type": "BlockStatement",
        "start": 16,
        "end": 33,
        "body": [
          {
            "type": "ExpressionStatement",
            "start": 18,
            "end": 31,
            "expression": {
              "type": "Literal",
              "start": 18,
              "end": 30,
              "value": "use strict",
              "raw": "'use strict'"
            }
          }
        ]
      }
    }
  ]
}
`, ast)
}
