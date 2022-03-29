package estree_test

import (
	"testing"

	. "github.com/hsiaosiyuan0/mole/ecma/estree/test"
	"github.com/hsiaosiyuan0/mole/ecma/parser"
	. "github.com/hsiaosiyuan0/mole/util"
)

// below are tests referred from:
// https://github.com/acornjs/acorn/blob/164bf8fc88e02bf5905be7788a9167c34176b50c/test/tests-asyncawait.js

func TestAsyncAwait1(t *testing.T) {
	ast, err := Compile("function foo() { }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 18,
  "body": [
    {
      "type": "FunctionDeclaration",
      "start": 0,
      "end": 18,
      "id": {
        "type": "Identifier",
        "start": 9,
        "end": 12,
        "name": "foo"
      },
      "generator": false,
      "async": false,
      "params": [],
      "body": {
        "type": "BlockStatement",
        "start": 15,
        "end": 18,
        "body": []
      }
    }
  ]
}
`, ast)
}

func TestAsyncAwait2(t *testing.T) {
	ast, err := Compile("async function foo() { }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 24,
  "body": [
    {
      "type": "FunctionDeclaration",
      "start": 0,
      "end": 24,
      "id": {
        "type": "Identifier",
        "start": 15,
        "end": 18,
        "name": "foo"
      },
      "generator": false,
      "async": true,
      "params": [],
      "body": {
        "type": "BlockStatement",
        "start": 21,
        "end": 24,
        "body": []
      }
    }
  ]
}
`, ast)
}

func TestAsyncAwait3(t *testing.T) {
	ast, err := Compile("async\nfunction foo() { }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 24,
  "body": [
    {
      "type": "ExpressionStatement",
      "start": 0,
      "end": 5,
      "expression": {
        "type": "Identifier",
        "start": 0,
        "end": 5,
        "name": "async"
      }
    },
    {
      "type": "FunctionDeclaration",
      "start": 6,
      "end": 24,
      "id": {
        "type": "Identifier",
        "start": 15,
        "end": 18,
        "name": "foo"
      },
      "generator": false,
      "async": false,
      "params": [],
      "body": {
        "type": "BlockStatement",
        "start": 21,
        "end": 24,
        "body": []
      }
    }
  ]
}
`, ast)
}

func TestAsyncAwait4(t *testing.T) {
	ast, err := Compile("export async function foo() { }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 31,
  "body": [
    {
      "type": "ExportNamedDeclaration",
      "start": 0,
      "end": 31,
      "declaration": {
        "type": "FunctionDeclaration",
        "start": 7,
        "end": 31,
        "id": {
          "type": "Identifier",
          "start": 22,
          "end": 25,
          "name": "foo"
        },
        "generator": false,
        "async": true,
        "params": [],
        "body": {
          "type": "BlockStatement",
          "start": 28,
          "end": 31,
          "body": []
        }
      },
      "specifiers": [],
      "source": null
    }
  ]
}
`, ast)
}

func TestAsyncAwait5(t *testing.T) {
	ast, err := Compile("export default async function() { }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 35,
  "body": [
    {
      "type": "ExportDefaultDeclaration",
      "start": 0,
      "end": 35,
      "declaration": {
        "type": "FunctionDeclaration",
        "start": 15,
        "end": 35,
        "id": null,
        "generator": false,
        "async": true,
        "params": [],
        "body": {
          "type": "BlockStatement",
          "start": 32,
          "end": 35,
          "body": []
        }
      }
    }
  ]
}
`, ast)
}

func TestAsyncAwait6(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_STRICT).Off(parser.FEAT_GLOBAL_ASYNC)
	ast, err := CompileWithOpts("async function await() { }", opts)
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 26,
  "body": [
    {
      "type": "FunctionDeclaration",
      "start": 0,
      "end": 26,
      "id": {
        "type": "Identifier",
        "start": 15,
        "end": 20,
        "name": "await"
      },
      "generator": false,
      "async": true,
      "params": [],
      "body": {
        "type": "BlockStatement",
        "start": 23,
        "end": 26,
        "body": []
      }
    }
  ]
}
`, ast)
}

func TestAsyncAwait7(t *testing.T) {
	ast, err := Compile("(function foo() { })")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 20,
  "body": [
    {
      "type": "ExpressionStatement",
      "start": 0,
      "end": 20,
      "expression": {
        "type": "FunctionExpression",
        "start": 1,
        "end": 19,
        "id": {
          "type": "Identifier",
          "start": 10,
          "end": 13,
          "name": "foo"
        },
        "generator": false,
        "async": false,
        "params": [],
        "body": {
          "type": "BlockStatement",
          "start": 16,
          "end": 19,
          "body": []
        }
      }
    }
  ]
}
`, ast)
}

func TestAsyncAwait8(t *testing.T) {
	ast, err := Compile("(async function foo() { })")
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
        "type": "FunctionExpression",
        "start": 1,
        "end": 25,
        "id": {
          "type": "Identifier",
          "start": 16,
          "end": 19,
          "name": "foo"
        },
        "generator": false,
        "async": true,
        "params": [],
        "body": {
          "type": "BlockStatement",
          "start": 22,
          "end": 25,
          "body": []
        }
      }
    }
  ]
}
`, ast)
}

func TestAsyncAwait9(t *testing.T) {
	ast, err := Compile("export default (async function() { })")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 37,
  "body": [
    {
      "type": "ExportDefaultDeclaration",
      "start": 0,
      "end": 37,
      "declaration": {
        "type": "FunctionExpression",
        "start": 16,
        "end": 36,
        "id": null,
        "generator": false,
        "async": true,
        "params": [],
        "body": {
          "type": "BlockStatement",
          "start": 33,
          "end": 36,
          "body": []
        }
      }
    }
  ]
}
`, ast)
}

func TestAsyncAwait10(t *testing.T) {
	ast, err := Compile("a => a")
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
        "type": "ArrowFunctionExpression",
        "start": 0,
        "end": 6,
        "id": null,
        "expression": true,
        "generator": false,
        "async": false,
        "params": [
          {
            "type": "Identifier",
            "start": 0,
            "end": 1,
            "name": "a"
          }
        ],
        "body": {
          "type": "Identifier",
          "start": 5,
          "end": 6,
          "name": "a"
        }
      }
    }
  ]
}
`, ast)
}

func TestAsyncAwait11(t *testing.T) {
	ast, err := Compile("(a) => a")
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
        "type": "ArrowFunctionExpression",
        "start": 0,
        "end": 8,
        "id": null,
        "expression": true,
        "generator": false,
        "async": false,
        "params": [
          {
            "type": "Identifier",
            "start": 1,
            "end": 2,
            "name": "a"
          }
        ],
        "body": {
          "type": "Identifier",
          "start": 7,
          "end": 8,
          "name": "a"
        }
      }
    }
  ]
}
`, ast)
}

func TestAsyncAwait12(t *testing.T) {
	ast, err := Compile("async a => a")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 12,
  "body": [
    {
      "type": "ExpressionStatement",
      "start": 0,
      "end": 12,
      "expression": {
        "type": "ArrowFunctionExpression",
        "start": 0,
        "end": 12,
        "id": null,
        "expression": true,
        "generator": false,
        "async": true,
        "params": [
          {
            "type": "Identifier",
            "start": 6,
            "end": 7,
            "name": "a"
          }
        ],
        "body": {
          "type": "Identifier",
          "start": 11,
          "end": 12,
          "name": "a"
        }
      }
    }
  ]
}
`, ast)
}

func TestAsyncAwait13(t *testing.T) {
	ast, err := Compile("async () => a")
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
        "type": "ArrowFunctionExpression",
        "start": 0,
        "end": 13,
        "id": null,
        "expression": true,
        "generator": false,
        "async": true,
        "params": [],
        "body": {
          "type": "Identifier",
          "start": 12,
          "end": 13,
          "name": "a"
        }
      }
    }
  ]
}
`, ast)
}

func TestAsyncAwait14(t *testing.T) {
	ast, err := Compile("async (a, b) => a")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 17,
  "body": [
    {
      "type": "ExpressionStatement",
      "start": 0,
      "end": 17,
      "expression": {
        "type": "ArrowFunctionExpression",
        "start": 0,
        "end": 17,
        "id": null,
        "expression": true,
        "generator": false,
        "async": true,
        "params": [
          {
            "type": "Identifier",
            "start": 7,
            "end": 8,
            "name": "a"
          },
          {
            "type": "Identifier",
            "start": 10,
            "end": 11,
            "name": "b"
          }
        ],
        "body": {
          "type": "Identifier",
          "start": 16,
          "end": 17,
          "name": "a"
        }
      }
    }
  ]
}
`, ast)
}

func TestAsyncAwait15(t *testing.T) {
	ast, err := Compile("async ({a = b}) => a")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 20,
  "body": [
    {
      "type": "ExpressionStatement",
      "start": 0,
      "end": 20,
      "expression": {
        "type": "ArrowFunctionExpression",
        "start": 0,
        "end": 20,
        "id": null,
        "expression": true,
        "generator": false,
        "async": true,
        "params": [
          {
            "type": "ObjectPattern",
            "start": 7,
            "end": 14,
            "properties": [
              {
                "type": "Property",
                "start": 8,
                "end": 13,
                "method": false,
                "shorthand": true,
                "computed": false,
                "key": {
                  "type": "Identifier",
                  "start": 8,
                  "end": 9,
                  "name": "a"
                },
                "kind": "init",
                "value": {
                  "type": "AssignmentPattern",
                  "start": 8,
                  "end": 13,
                  "left": {
                    "type": "Identifier",
                    "start": 8,
                    "end": 9,
                    "name": "a"
                  },
                  "right": {
                    "type": "Identifier",
                    "start": 12,
                    "end": 13,
                    "name": "b"
                  }
                }
              }
            ]
          }
        ],
        "body": {
          "type": "Identifier",
          "start": 19,
          "end": 20,
          "name": "a"
        }
      }
    }
  ]
}
`, ast)
}

func TestAsyncAwait16(t *testing.T) {
	ast, err := Compile("async ({a: b = c}) => a")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 23,
  "body": [
    {
      "type": "ExpressionStatement",
      "start": 0,
      "end": 23,
      "expression": {
        "type": "ArrowFunctionExpression",
        "start": 0,
        "end": 23,
        "id": null,
        "expression": true,
        "generator": false,
        "async": true,
        "params": [
          {
            "type": "ObjectPattern",
            "start": 7,
            "end": 17,
            "properties": [
              {
                "type": "Property",
                "start": 8,
                "end": 16,
                "method": false,
                "shorthand": false,
                "computed": false,
                "key": {
                  "type": "Identifier",
                  "start": 8,
                  "end": 9,
                  "name": "a"
                },
                "value": {
                  "type": "AssignmentPattern",
                  "start": 11,
                  "end": 16,
                  "left": {
                    "type": "Identifier",
                    "start": 11,
                    "end": 12,
                    "name": "b"
                  },
                  "right": {
                    "type": "Identifier",
                    "start": 15,
                    "end": 16,
                    "name": "c"
                  }
                },
                "kind": "init"
              }
            ]
          }
        ],
        "body": {
          "type": "Identifier",
          "start": 22,
          "end": 23,
          "name": "a"
        }
      }
    }
  ]
}
`, ast)
}

func TestAsyncAwait17(t *testing.T) {
	ast, err := Compile("async ({a: b = c})")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 18,
  "body": [
    {
      "type": "ExpressionStatement",
      "start": 0,
      "end": 18,
      "expression": {
        "type": "CallExpression",
        "start": 0,
        "end": 18,
        "callee": {
          "type": "Identifier",
          "start": 0,
          "end": 5,
          "name": "async"
        },
        "arguments": [
          {
            "type": "ObjectExpression",
            "start": 7,
            "end": 17,
            "properties": [
              {
                "type": "Property",
                "start": 8,
                "end": 16,
                "method": false,
                "shorthand": false,
                "computed": false,
                "key": {
                  "type": "Identifier",
                  "start": 8,
                  "end": 9,
                  "name": "a"
                },
                "value": {
                  "type": "AssignmentExpression",
                  "start": 11,
                  "end": 16,
                  "operator": "=",
                  "left": {
                    "type": "Identifier",
                    "start": 11,
                    "end": 12,
                    "name": "b"
                  },
                  "right": {
                    "type": "Identifier",
                    "start": 15,
                    "end": 16,
                    "name": "c"
                  }
                },
                "kind": "init"
              }
            ]
          }
        ],
        "optional": false
      }
    }
  ]
}
`, ast)
}

func TestAsyncAwait18(t *testing.T) {
	ast, err := Compile("async\na => a")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 12,
  "body": [
    {
      "type": "ExpressionStatement",
      "start": 0,
      "end": 5,
      "expression": {
        "type": "Identifier",
        "start": 0,
        "end": 5,
        "name": "async"
      }
    },
    {
      "type": "ExpressionStatement",
      "start": 6,
      "end": 12,
      "expression": {
        "type": "ArrowFunctionExpression",
        "start": 6,
        "end": 12,
        "id": null,
        "expression": true,
        "generator": false,
        "async": false,
        "params": [
          {
            "type": "Identifier",
            "start": 6,
            "end": 7,
            "name": "a"
          }
        ],
        "body": {
          "type": "Identifier",
          "start": 11,
          "end": 12,
          "name": "a"
        }
      }
    }
  ]
}
`, ast)
}

func TestAsyncAwait19(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_MODULE).Off(parser.FEAT_STRICT).Off(parser.FEAT_GLOBAL_ASYNC)
	ast, err := CompileWithOpts("async (await)", opts)
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
        "type": "CallExpression",
        "start": 0,
        "end": 13,
        "callee": {
          "type": "Identifier",
          "start": 0,
          "end": 5,
          "name": "async"
        },
        "arguments": [
          {
            "type": "Identifier",
            "start": 7,
            "end": 12,
            "name": "await"
          }
        ]
      }
    }
  ]
}
`, ast)
}

func TestAsyncAwait20(t *testing.T) {
	ast, err := Compile("async ({await: a}) => 1")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 23,
  "body": [
    {
      "type": "ExpressionStatement",
      "start": 0,
      "end": 23,
      "expression": {
        "type": "ArrowFunctionExpression",
        "start": 0,
        "end": 23,
        "id": null,
        "expression": true,
        "generator": false,
        "async": true,
        "params": [
          {
            "type": "ObjectPattern",
            "start": 7,
            "end": 17,
            "properties": [
              {
                "type": "Property",
                "start": 8,
                "end": 16,
                "method": false,
                "shorthand": false,
                "computed": false,
                "key": {
                  "type": "Identifier",
                  "start": 8,
                  "end": 13,
                  "name": "await"
                },
                "value": {
                  "type": "Identifier",
                  "start": 15,
                  "end": 16,
                  "name": "a"
                },
                "kind": "init"
              }
            ]
          }
        ],
        "body": {
          "type": "Literal",
          "start": 22,
          "end": 23,
          "value": 1,
          "raw": "1"
        }
      }
    }
  ]
}
`, ast)
}

func TestAsyncAwait21(t *testing.T) {
	ast, err := Compile("async (b = {await: a}) => 1")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 27,
  "body": [
    {
      "type": "ExpressionStatement",
      "start": 0,
      "end": 27,
      "expression": {
        "type": "ArrowFunctionExpression",
        "start": 0,
        "end": 27,
        "id": null,
        "expression": true,
        "generator": false,
        "async": true,
        "params": [
          {
            "type": "AssignmentPattern",
            "start": 7,
            "end": 21,
            "left": {
              "type": "Identifier",
              "start": 7,
              "end": 8,
              "name": "b"
            },
            "right": {
              "type": "ObjectExpression",
              "start": 11,
              "end": 21,
              "properties": [
                {
                  "type": "Property",
                  "start": 12,
                  "end": 20,
                  "method": false,
                  "shorthand": false,
                  "computed": false,
                  "key": {
                    "type": "Identifier",
                    "start": 12,
                    "end": 17,
                    "name": "await"
                  },
                  "value": {
                    "type": "Identifier",
                    "start": 19,
                    "end": 20,
                    "name": "a"
                  },
                  "kind": "init"
                }
              ]
            }
          }
        ],
        "body": {
          "type": "Literal",
          "start": 26,
          "end": 27,
          "value": 1,
          "raw": "1"
        }
      }
    }
  ]
}
`, ast)
}

func TestAsyncAwait22(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_STRICT).Off(parser.FEAT_GLOBAL_ASYNC)
	ast, err := CompileWithOpts("async (b = function* await() {}) => 1", opts)
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 37,
  "body": [
    {
      "type": "ExpressionStatement",
      "start": 0,
      "end": 37,
      "expression": {
        "type": "ArrowFunctionExpression",
        "start": 0,
        "end": 37,
        "id": null,
        "expression": true,
        "generator": false,
        "async": true,
        "params": [
          {
            "type": "AssignmentPattern",
            "start": 7,
            "end": 31,
            "left": {
              "type": "Identifier",
              "start": 7,
              "end": 8,
              "name": "b"
            },
            "right": {
              "type": "FunctionExpression",
              "start": 11,
              "end": 31,
              "id": {
                "type": "Identifier",
                "start": 21,
                "end": 26,
                "name": "await"
              },
              "expression": false,
              "generator": true,
              "async": false,
              "params": [],
              "body": {
                "type": "BlockStatement",
                "start": 29,
                "end": 31,
                "body": []
              }
            }
          }
        ],
        "body": {
          "type": "Literal",
          "start": 36,
          "end": 37,
          "value": 1,
          "raw": "1"
        }
      }
    }
  ]
}
`, ast)
}

func TestAsyncAwait23(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_STRICT)
	ast, err := CompileWithOpts("async yield => 1", opts)
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 16,
  "body": [
    {
      "type": "ExpressionStatement",
      "start": 0,
      "end": 16,
      "expression": {
        "type": "ArrowFunctionExpression",
        "start": 0,
        "end": 16,
        "id": null,
        "expression": true,
        "generator": false,
        "async": true,
        "params": [
          {
            "type": "Identifier",
            "start": 6,
            "end": 11,
            "name": "yield"
          }
        ],
        "body": {
          "type": "Literal",
          "start": 15,
          "end": 16,
          "value": 1,
          "raw": "1"
        }
      }
    }
  ]
}
`, ast)
}

func TestAsyncAwait24(t *testing.T) {
	ast, err := Compile("({foo() { }})")
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
        "type": "ObjectExpression",
        "start": 1,
        "end": 12,
        "properties": [
          {
            "type": "Property",
            "start": 2,
            "end": 11,
            "method": true,
            "shorthand": false,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 2,
              "end": 5,
              "name": "foo"
            },
            "kind": "init",
            "value": {
              "type": "FunctionExpression",
              "start": 5,
              "end": 11,
              "id": null,
              "expression": false,
              "generator": false,
              "async": false,
              "params": [],
              "body": {
                "type": "BlockStatement",
                "start": 8,
                "end": 11,
                "body": []
              }
            }
          }
        ]
      }
    }
  ]
}
`, ast)
}

func TestAsyncAwait25(t *testing.T) {
	ast, err := Compile("({async foo() { }})")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 19,
  "body": [
    {
      "type": "ExpressionStatement",
      "start": 0,
      "end": 19,
      "expression": {
        "type": "ObjectExpression",
        "start": 1,
        "end": 18,
        "properties": [
          {
            "type": "Property",
            "start": 2,
            "end": 17,
            "method": true,
            "shorthand": false,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 8,
              "end": 11,
              "name": "foo"
            },
            "kind": "init",
            "value": {
              "type": "FunctionExpression",
              "start": 11,
              "end": 17,
              "id": null,
              "expression": false,
              "generator": false,
              "async": true,
              "params": [],
              "body": {
                "type": "BlockStatement",
                "start": 14,
                "end": 17,
                "body": []
              }
            }
          }
        ]
      }
    }
  ]
}
`, ast)
}

func TestAsyncAwait26(t *testing.T) {
	ast, err := Compile("({async() { }})")
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
        "type": "ObjectExpression",
        "start": 1,
        "end": 14,
        "properties": [
          {
            "type": "Property",
            "start": 2,
            "end": 13,
            "method": true,
            "shorthand": false,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 2,
              "end": 7,
              "name": "async"
            },
            "kind": "init",
            "value": {
              "type": "FunctionExpression",
              "start": 7,
              "end": 13,
              "id": null,
              "expression": false,
              "generator": false,
              "async": false,
              "params": [],
              "body": {
                "type": "BlockStatement",
                "start": 10,
                "end": 13,
                "body": []
              }
            }
          }
        ]
      }
    }
  ]
}
`, ast)
}

func TestAsyncAwait27(t *testing.T) {
	ast, err := Compile("({async await() { }})")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 21,
  "body": [
    {
      "type": "ExpressionStatement",
      "start": 0,
      "end": 21,
      "expression": {
        "type": "ObjectExpression",
        "start": 1,
        "end": 20,
        "properties": [
          {
            "type": "Property",
            "start": 2,
            "end": 19,
            "method": true,
            "shorthand": false,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 8,
              "end": 13,
              "name": "await"
            },
            "kind": "init",
            "value": {
              "type": "FunctionExpression",
              "start": 13,
              "end": 19,
              "id": null,
              "expression": false,
              "generator": false,
              "async": true,
              "params": [],
              "body": {
                "type": "BlockStatement",
                "start": 16,
                "end": 19,
                "body": []
              }
            }
          }
        ]
      }
    }
  ]
}
`, ast)
}

func TestAsyncAwait28(t *testing.T) {
	ast, err := Compile("async function wrap() {\n({async await() { }})\n}")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 47,
  "body": [
    {
      "type": "FunctionDeclaration",
      "start": 0,
      "end": 47,
      "id": {
        "type": "Identifier",
        "start": 15,
        "end": 19,
        "name": "wrap"
      },
      "generator": false,
      "async": true,
      "params": [],
      "body": {
        "type": "BlockStatement",
        "start": 22,
        "end": 47,
        "body": [
          {
            "type": "ExpressionStatement",
            "start": 24,
            "end": 45,
            "expression": {
              "type": "ObjectExpression",
              "start": 25,
              "end": 44,
              "properties": [
                {
                  "type": "Property",
                  "start": 26,
                  "end": 43,
                  "method": true,
                  "shorthand": false,
                  "computed": false,
                  "key": {
                    "type": "Identifier",
                    "start": 32,
                    "end": 37,
                    "name": "await"
                  },
                  "kind": "init",
                  "value": {
                    "type": "FunctionExpression",
                    "start": 37,
                    "end": 43,
                    "id": null,
                    "expression": false,
                    "generator": false,
                    "async": true,
                    "params": [],
                    "body": {
                      "type": "BlockStatement",
                      "start": 40,
                      "end": 43,
                      "body": []
                    }
                  }
                }
              ]
            }
          }
        ]
      }
    }
  ]
}
`, ast)
}

func TestAsyncAwait29(t *testing.T) {
	ast, err := Compile("class A {foo() { }}")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 19,
  "body": [
    {
      "type": "ClassDeclaration",
      "start": 0,
      "end": 19,
      "id": {
        "type": "Identifier",
        "start": 6,
        "end": 7,
        "name": "A"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "start": 8,
        "end": 19,
        "body": [
          {
            "type": "MethodDefinition",
            "start": 9,
            "end": 18,
            "static": false,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 9,
              "end": 12,
              "name": "foo"
            },
            "kind": "method",
            "value": {
              "type": "FunctionExpression",
              "start": 12,
              "end": 18,
              "id": null,
              "expression": false,
              "generator": false,
              "async": false,
              "params": [],
              "body": {
                "type": "BlockStatement",
                "start": 15,
                "end": 18,
                "body": []
              }
            }
          }
        ]
      }
    }
  ]
}
`, ast)
}

func TestAsyncAwait30(t *testing.T) {
	ast, err := Compile("class A {async foo() { }}")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 25,
  "body": [
    {
      "type": "ClassDeclaration",
      "start": 0,
      "end": 25,
      "id": {
        "type": "Identifier",
        "start": 6,
        "end": 7,
        "name": "A"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "start": 8,
        "end": 25,
        "body": [
          {
            "type": "MethodDefinition",
            "start": 9,
            "end": 24,
            "static": false,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 15,
              "end": 18,
              "name": "foo"
            },
            "kind": "method",
            "value": {
              "type": "FunctionExpression",
              "start": 18,
              "end": 24,
              "id": null,
              "expression": false,
              "generator": false,
              "async": true,
              "params": [],
              "body": {
                "type": "BlockStatement",
                "start": 21,
                "end": 24,
                "body": []
              }
            }
          }
        ]
      }
    }
  ]
}
`, ast)
}

func TestAsyncAwait31(t *testing.T) {
	ast, err := Compile("class A {static async foo() { }}")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 32,
  "body": [
    {
      "type": "ClassDeclaration",
      "start": 0,
      "end": 32,
      "id": {
        "type": "Identifier",
        "start": 6,
        "end": 7,
        "name": "A"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "start": 8,
        "end": 32,
        "body": [
          {
            "type": "MethodDefinition",
            "start": 9,
            "end": 31,
            "static": true,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 22,
              "end": 25,
              "name": "foo"
            },
            "kind": "method",
            "value": {
              "type": "FunctionExpression",
              "start": 25,
              "end": 31,
              "id": null,
              "expression": false,
              "generator": false,
              "async": true,
              "params": [],
              "body": {
                "type": "BlockStatement",
                "start": 28,
                "end": 31,
                "body": []
              }
            }
          }
        ]
      }
    }
  ]
}
`, ast)
}

func TestAsyncAwait32(t *testing.T) {
	ast, err := Compile("class A {async() { }}")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 21,
  "body": [
    {
      "type": "ClassDeclaration",
      "start": 0,
      "end": 21,
      "id": {
        "type": "Identifier",
        "start": 6,
        "end": 7,
        "name": "A"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "start": 8,
        "end": 21,
        "body": [
          {
            "type": "MethodDefinition",
            "start": 9,
            "end": 20,
            "static": false,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 9,
              "end": 14,
              "name": "async"
            },
            "kind": "method",
            "value": {
              "type": "FunctionExpression",
              "start": 14,
              "end": 20,
              "id": null,
              "expression": false,
              "generator": false,
              "async": false,
              "params": [],
              "body": {
                "type": "BlockStatement",
                "start": 17,
                "end": 20,
                "body": []
              }
            }
          }
        ]
      }
    }
  ]
}
`, ast)
}

func TestAsyncAwait33(t *testing.T) {
	ast, err := Compile("class A {static async() { }}")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 28,
  "body": [
    {
      "type": "ClassDeclaration",
      "start": 0,
      "end": 28,
      "id": {
        "type": "Identifier",
        "start": 6,
        "end": 7,
        "name": "A"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "start": 8,
        "end": 28,
        "body": [
          {
            "type": "MethodDefinition",
            "start": 9,
            "end": 27,
            "static": true,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 16,
              "end": 21,
              "name": "async"
            },
            "kind": "method",
            "value": {
              "type": "FunctionExpression",
              "start": 21,
              "end": 27,
              "id": null,
              "expression": false,
              "generator": false,
              "async": false,
              "params": [],
              "body": {
                "type": "BlockStatement",
                "start": 24,
                "end": 27,
                "body": []
              }
            }
          }
        ]
      }
    }
  ]
}
`, ast)
}

func TestAsyncAwait34(t *testing.T) {
	ast, err := Compile("class A {*async() { }}")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 22,
  "body": [
    {
      "type": "ClassDeclaration",
      "start": 0,
      "end": 22,
      "id": {
        "type": "Identifier",
        "start": 6,
        "end": 7,
        "name": "A"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "start": 8,
        "end": 22,
        "body": [
          {
            "type": "MethodDefinition",
            "start": 9,
            "end": 21,
            "static": false,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 10,
              "end": 15,
              "name": "async"
            },
            "kind": "method",
            "value": {
              "type": "FunctionExpression",
              "start": 15,
              "end": 21,
              "id": null,
              "expression": false,
              "generator": true,
              "async": false,
              "params": [],
              "body": {
                "type": "BlockStatement",
                "start": 18,
                "end": 21,
                "body": []
              }
            }
          }
        ]
      }
    }
  ]
}
`, ast)
}

func TestAsyncAwait35(t *testing.T) {
	ast, err := Compile("class A {static* async() { }}")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 29,
  "body": [
    {
      "type": "ClassDeclaration",
      "start": 0,
      "end": 29,
      "id": {
        "type": "Identifier",
        "start": 6,
        "end": 7,
        "name": "A"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "start": 8,
        "end": 29,
        "body": [
          {
            "type": "MethodDefinition",
            "start": 9,
            "end": 28,
            "static": true,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 17,
              "end": 22,
              "name": "async"
            },
            "kind": "method",
            "value": {
              "type": "FunctionExpression",
              "start": 22,
              "end": 28,
              "id": null,
              "expression": false,
              "generator": true,
              "async": false,
              "params": [],
              "body": {
                "type": "BlockStatement",
                "start": 25,
                "end": 28,
                "body": []
              }
            }
          }
        ]
      }
    }
  ]
}
`, ast)
}

func TestAsyncAwait36(t *testing.T) {
	ast, err := Compile("class A {async await() { }}")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 27,
  "body": [
    {
      "type": "ClassDeclaration",
      "start": 0,
      "end": 27,
      "id": {
        "type": "Identifier",
        "start": 6,
        "end": 7,
        "name": "A"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "start": 8,
        "end": 27,
        "body": [
          {
            "type": "MethodDefinition",
            "start": 9,
            "end": 26,
            "static": false,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 15,
              "end": 20,
              "name": "await"
            },
            "kind": "method",
            "value": {
              "type": "FunctionExpression",
              "start": 20,
              "end": 26,
              "id": null,
              "expression": false,
              "generator": false,
              "async": true,
              "params": [],
              "body": {
                "type": "BlockStatement",
                "start": 23,
                "end": 26,
                "body": []
              }
            }
          }
        ]
      }
    }
  ]
}
`, ast)
}

func TestAsyncAwait37(t *testing.T) {
	ast, err := Compile("class A {static async await() { }}")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 34,
  "body": [
    {
      "type": "ClassDeclaration",
      "start": 0,
      "end": 34,
      "id": {
        "type": "Identifier",
        "start": 6,
        "end": 7,
        "name": "A"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "start": 8,
        "end": 34,
        "body": [
          {
            "type": "MethodDefinition",
            "start": 9,
            "end": 33,
            "static": true,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 22,
              "end": 27,
              "name": "await"
            },
            "kind": "method",
            "value": {
              "type": "FunctionExpression",
              "start": 27,
              "end": 33,
              "id": null,
              "expression": false,
              "generator": false,
              "async": true,
              "params": [],
              "body": {
                "type": "BlockStatement",
                "start": 30,
                "end": 33,
                "body": []
              }
            }
          }
        ]
      }
    }
  ]
}
`, ast)
}

func TestAsyncAwait38(t *testing.T) {
	ast, err := Compile("async function wrap() {\nclass A {async await() { }}\n}")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 53,
  "body": [
    {
      "type": "FunctionDeclaration",
      "start": 0,
      "end": 53,
      "id": {
        "type": "Identifier",
        "start": 15,
        "end": 19,
        "name": "wrap"
      },
      "generator": false,
      "async": true,
      "params": [],
      "body": {
        "type": "BlockStatement",
        "start": 22,
        "end": 53,
        "body": [
          {
            "type": "ClassDeclaration",
            "start": 24,
            "end": 51,
            "id": {
              "type": "Identifier",
              "start": 30,
              "end": 31,
              "name": "A"
            },
            "superClass": null,
            "body": {
              "type": "ClassBody",
              "start": 32,
              "end": 51,
              "body": [
                {
                  "type": "MethodDefinition",
                  "start": 33,
                  "end": 50,
                  "static": false,
                  "computed": false,
                  "key": {
                    "type": "Identifier",
                    "start": 39,
                    "end": 44,
                    "name": "await"
                  },
                  "kind": "method",
                  "value": {
                    "type": "FunctionExpression",
                    "start": 44,
                    "end": 50,
                    "id": null,
                    "expression": false,
                    "generator": false,
                    "async": true,
                    "params": [],
                    "body": {
                      "type": "BlockStatement",
                      "start": 47,
                      "end": 50,
                      "body": []
                    }
                  }
                }
              ]
            }
          }
        ]
      }
    }
  ]
}
`, ast)
}

func TestAsyncAwait39(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_MODULE).Off(parser.FEAT_GLOBAL_ASYNC)
	ast, err := CompileWithOpts("await", opts)
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 5,
  "body": [
    {
      "type": "ExpressionStatement",
      "start": 0,
      "end": 5,
      "expression": {
        "type": "Identifier",
        "start": 0,
        "end": 5,
        "name": "await"
      }
    }
  ]
}
`, ast)
}

func TestAsyncAwait40(t *testing.T) {
	ast, err := Compile("async function foo(a, b) { await a }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 36,
  "body": [
    {
      "type": "FunctionDeclaration",
      "start": 0,
      "end": 36,
      "id": {
        "type": "Identifier",
        "start": 15,
        "end": 18,
        "name": "foo"
      },
      "generator": false,
      "async": true,
      "params": [
        {
          "type": "Identifier",
          "start": 19,
          "end": 20,
          "name": "a"
        },
        {
          "type": "Identifier",
          "start": 22,
          "end": 23,
          "name": "b"
        }
      ],
      "body": {
        "type": "BlockStatement",
        "start": 25,
        "end": 36,
        "body": [
          {
            "type": "ExpressionStatement",
            "start": 27,
            "end": 34,
            "expression": {
              "type": "AwaitExpression",
              "start": 27,
              "end": 34,
              "argument": {
                "type": "Identifier",
                "start": 33,
                "end": 34,
                "name": "a"
              }
            }
          }
        ]
      }
    }
  ]
}
`, ast)
}

func TestAsyncAwait41(t *testing.T) {
	ast, err := Compile("(async function foo(a) { await a })")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 35,
  "body": [
    {
      "type": "ExpressionStatement",
      "start": 0,
      "end": 35,
      "expression": {
        "type": "FunctionExpression",
        "start": 1,
        "end": 34,
        "id": {
          "type": "Identifier",
          "start": 16,
          "end": 19,
          "name": "foo"
        },
        "expression": false,
        "generator": false,
        "async": true,
        "params": [
          {
            "type": "Identifier",
            "start": 20,
            "end": 21,
            "name": "a"
          }
        ],
        "body": {
          "type": "BlockStatement",
          "start": 23,
          "end": 34,
          "body": [
            {
              "type": "ExpressionStatement",
              "start": 25,
              "end": 32,
              "expression": {
                "type": "AwaitExpression",
                "start": 25,
                "end": 32,
                "argument": {
                  "type": "Identifier",
                  "start": 31,
                  "end": 32,
                  "name": "a"
                }
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

func TestAsyncAwait42(t *testing.T) {
	ast, err := Compile("(async (a) => await a)")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 22,
  "body": [
    {
      "type": "ExpressionStatement",
      "start": 0,
      "end": 22,
      "expression": {
        "type": "ArrowFunctionExpression",
        "start": 1,
        "end": 21,
        "id": null,
        "expression": true,
        "generator": false,
        "async": true,
        "params": [
          {
            "type": "Identifier",
            "start": 8,
            "end": 9,
            "name": "a"
          }
        ],
        "body": {
          "type": "AwaitExpression",
          "start": 14,
          "end": 21,
          "argument": {
            "type": "Identifier",
            "start": 20,
            "end": 21,
            "name": "a"
          }
        }
      }
    }
  ]
}
`, ast)
}

func TestAsyncAwait43(t *testing.T) {
	ast, err := Compile("({async foo(a) { await a }})")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 28,
  "body": [
    {
      "type": "ExpressionStatement",
      "start": 0,
      "end": 28,
      "expression": {
        "type": "ObjectExpression",
        "start": 1,
        "end": 27,
        "properties": [
          {
            "type": "Property",
            "start": 2,
            "end": 26,
            "method": true,
            "shorthand": false,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 8,
              "end": 11,
              "name": "foo"
            },
            "kind": "init",
            "value": {
              "type": "FunctionExpression",
              "start": 11,
              "end": 26,
              "id": null,
              "expression": false,
              "generator": false,
              "async": true,
              "params": [
                {
                  "type": "Identifier",
                  "start": 12,
                  "end": 13,
                  "name": "a"
                }
              ],
              "body": {
                "type": "BlockStatement",
                "start": 15,
                "end": 26,
                "body": [
                  {
                    "type": "ExpressionStatement",
                    "start": 17,
                    "end": 24,
                    "expression": {
                      "type": "AwaitExpression",
                      "start": 17,
                      "end": 24,
                      "argument": {
                        "type": "Identifier",
                        "start": 23,
                        "end": 24,
                        "name": "a"
                      }
                    }
                  }
                ]
              }
            }
          }
        ]
      }
    }
  ]
}
`, ast)
}

func TestAsyncAwait44(t *testing.T) {
	ast, err := Compile("async function foo(a, b) { await a + await b }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 46,
  "body": [
    {
      "type": "FunctionDeclaration",
      "start": 0,
      "end": 46,
      "id": {
        "type": "Identifier",
        "start": 15,
        "end": 18,
        "name": "foo"
      },
      "generator": false,
      "async": true,
      "params": [
        {
          "type": "Identifier",
          "start": 19,
          "end": 20,
          "name": "a"
        },
        {
          "type": "Identifier",
          "start": 22,
          "end": 23,
          "name": "b"
        }
      ],
      "body": {
        "type": "BlockStatement",
        "start": 25,
        "end": 46,
        "body": [
          {
            "type": "ExpressionStatement",
            "start": 27,
            "end": 44,
            "expression": {
              "type": "BinaryExpression",
              "start": 27,
              "end": 44,
              "left": {
                "type": "AwaitExpression",
                "start": 27,
                "end": 34,
                "argument": {
                  "type": "Identifier",
                  "start": 33,
                  "end": 34,
                  "name": "a"
                }
              },
              "operator": "+",
              "right": {
                "type": "AwaitExpression",
                "start": 37,
                "end": 44,
                "argument": {
                  "type": "Identifier",
                  "start": 43,
                  "end": 44,
                  "name": "b"
                }
              }
            }
          }
        ]
      }
    }
  ]
}
`, ast)
}

func TestAsyncAwait45(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_MODULE)
	ast, err := CompileWithOpts("function foo() { await + 1 }", opts)
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 28,
  "body": [
    {
      "type": "FunctionDeclaration",
      "start": 0,
      "end": 28,
      "id": {
        "type": "Identifier",
        "start": 9,
        "end": 12,
        "name": "foo"
      },
      "generator": false,
      "async": false,
      "params": [],
      "body": {
        "type": "BlockStatement",
        "start": 15,
        "end": 28,
        "body": [
          {
            "type": "ExpressionStatement",
            "start": 17,
            "end": 26,
            "expression": {
              "type": "BinaryExpression",
              "start": 17,
              "end": 26,
              "left": {
                "type": "Identifier",
                "start": 17,
                "end": 22,
                "name": "await"
              },
              "operator": "+",
              "right": {
                "type": "Literal",
                "start": 25,
                "end": 26,
                "value": 1,
                "raw": "1"
              }
            }
          }
        ]
      }
    }
  ]
}
`, ast)
}

func TestAsyncAwait46(t *testing.T) {
	ast, err := Compile("async function foo() { await + 1 }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 34,
  "body": [
    {
      "type": "FunctionDeclaration",
      "start": 0,
      "end": 34,
      "id": {
        "type": "Identifier",
        "start": 15,
        "end": 18,
        "name": "foo"
      },
      "generator": false,
      "async": true,
      "params": [],
      "body": {
        "type": "BlockStatement",
        "start": 21,
        "end": 34,
        "body": [
          {
            "type": "ExpressionStatement",
            "start": 23,
            "end": 32,
            "expression": {
              "type": "AwaitExpression",
              "start": 23,
              "end": 32,
              "argument": {
                "type": "UnaryExpression",
                "start": 29,
                "end": 32,
                "operator": "+",
                "prefix": true,
                "argument": {
                  "type": "Literal",
                  "start": 31,
                  "end": 32,
                  "value": 1,
                  "raw": "1"
                }
              }
            }
          }
        ]
      }
    }
  ]
}
`, ast)
}

func TestAsyncAwait47(t *testing.T) {
	ast, err := Compile("async function foo(a = async function foo() { await b }) {}")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 59,
  "body": [
    {
      "type": "FunctionDeclaration",
      "start": 0,
      "end": 59,
      "id": {
        "type": "Identifier",
        "start": 15,
        "end": 18,
        "name": "foo"
      },
      "generator": false,
      "async": true,
      "params": [
        {
          "type": "AssignmentPattern",
          "start": 19,
          "end": 55,
          "left": {
            "type": "Identifier",
            "start": 19,
            "end": 20,
            "name": "a"
          },
          "right": {
            "type": "FunctionExpression",
            "start": 23,
            "end": 55,
            "id": {
              "type": "Identifier",
              "start": 38,
              "end": 41,
              "name": "foo"
            },
            "expression": false,
            "generator": false,
            "async": true,
            "params": [],
            "body": {
              "type": "BlockStatement",
              "start": 44,
              "end": 55,
              "body": [
                {
                  "type": "ExpressionStatement",
                  "start": 46,
                  "end": 53,
                  "expression": {
                    "type": "AwaitExpression",
                    "start": 46,
                    "end": 53,
                    "argument": {
                      "type": "Identifier",
                      "start": 52,
                      "end": 53,
                      "name": "b"
                    }
                  }
                }
              ]
            }
          }
        }
      ],
      "body": {
        "type": "BlockStatement",
        "start": 57,
        "end": 59,
        "body": []
      }
    }
  ]
}
`, ast)
}

func TestAsyncAwait48(t *testing.T) {
	ast, err := Compile("async function foo(a = async () => await b) {}")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 46,
  "body": [
    {
      "type": "FunctionDeclaration",
      "start": 0,
      "end": 46,
      "id": {
        "type": "Identifier",
        "start": 15,
        "end": 18,
        "name": "foo"
      },
      "generator": false,
      "async": true,
      "params": [
        {
          "type": "AssignmentPattern",
          "start": 19,
          "end": 42,
          "left": {
            "type": "Identifier",
            "start": 19,
            "end": 20,
            "name": "a"
          },
          "right": {
            "type": "ArrowFunctionExpression",
            "start": 23,
            "end": 42,
            "id": null,
            "expression": true,
            "generator": false,
            "async": true,
            "params": [],
            "body": {
              "type": "AwaitExpression",
              "start": 35,
              "end": 42,
              "argument": {
                "type": "Identifier",
                "start": 41,
                "end": 42,
                "name": "b"
              }
            }
          }
        }
      ],
      "body": {
        "type": "BlockStatement",
        "start": 44,
        "end": 46,
        "body": []
      }
    }
  ]
}
`, ast)
}

func TestAsyncAwait49(t *testing.T) {
	ast, err := Compile("async function foo(a = {async bar() { await b }}) {}")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 52,
  "body": [
    {
      "type": "FunctionDeclaration",
      "start": 0,
      "end": 52,
      "id": {
        "type": "Identifier",
        "start": 15,
        "end": 18,
        "name": "foo"
      },
      "generator": false,
      "async": true,
      "params": [
        {
          "type": "AssignmentPattern",
          "start": 19,
          "end": 48,
          "left": {
            "type": "Identifier",
            "start": 19,
            "end": 20,
            "name": "a"
          },
          "right": {
            "type": "ObjectExpression",
            "start": 23,
            "end": 48,
            "properties": [
              {
                "type": "Property",
                "start": 24,
                "end": 47,
                "method": true,
                "shorthand": false,
                "computed": false,
                "key": {
                  "type": "Identifier",
                  "start": 30,
                  "end": 33,
                  "name": "bar"
                },
                "kind": "init",
                "value": {
                  "type": "FunctionExpression",
                  "start": 33,
                  "end": 47,
                  "id": null,
                  "expression": false,
                  "generator": false,
                  "async": true,
                  "params": [],
                  "body": {
                    "type": "BlockStatement",
                    "start": 36,
                    "end": 47,
                    "body": [
                      {
                        "type": "ExpressionStatement",
                        "start": 38,
                        "end": 45,
                        "expression": {
                          "type": "AwaitExpression",
                          "start": 38,
                          "end": 45,
                          "argument": {
                            "type": "Identifier",
                            "start": 44,
                            "end": 45,
                            "name": "b"
                          }
                        }
                      }
                    ]
                  }
                }
              }
            ]
          }
        }
      ],
      "body": {
        "type": "BlockStatement",
        "start": 50,
        "end": 52,
        "body": []
      }
    }
  ]
}
`, ast)
}

func TestAsyncAwait50(t *testing.T) {
	ast, err := Compile("async function foo(a = class {async bar() { await b }}) {}")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 58,
  "body": [
    {
      "type": "FunctionDeclaration",
      "start": 0,
      "end": 58,
      "id": {
        "type": "Identifier",
        "start": 15,
        "end": 18,
        "name": "foo"
      },
      "generator": false,
      "async": true,
      "params": [
        {
          "type": "AssignmentPattern",
          "start": 19,
          "end": 54,
          "left": {
            "type": "Identifier",
            "start": 19,
            "end": 20,
            "name": "a"
          },
          "right": {
            "type": "ClassExpression",
            "start": 23,
            "end": 54,
            "id": null,
            "superClass": null,
            "body": {
              "type": "ClassBody",
              "start": 29,
              "end": 54,
              "body": [
                {
                  "type": "MethodDefinition",
                  "start": 30,
                  "end": 53,
                  "static": false,
                  "computed": false,
                  "key": {
                    "type": "Identifier",
                    "start": 36,
                    "end": 39,
                    "name": "bar"
                  },
                  "kind": "method",
                  "value": {
                    "type": "FunctionExpression",
                    "start": 39,
                    "end": 53,
                    "id": null,
                    "expression": false,
                    "generator": false,
                    "async": true,
                    "params": [],
                    "body": {
                      "type": "BlockStatement",
                      "start": 42,
                      "end": 53,
                      "body": [
                        {
                          "type": "ExpressionStatement",
                          "start": 44,
                          "end": 51,
                          "expression": {
                            "type": "AwaitExpression",
                            "start": 44,
                            "end": 51,
                            "argument": {
                              "type": "Identifier",
                              "start": 50,
                              "end": 51,
                              "name": "b"
                            }
                          }
                        }
                      ]
                    }
                  }
                }
              ]
            }
          }
        }
      ],
      "body": {
        "type": "BlockStatement",
        "start": 56,
        "end": 58,
        "body": []
      }
    }
  ]
}
`, ast)
}

func TestAsyncAwait51(t *testing.T) {
	ast, err := Compile("async function wrap() {\n(a = await b)\n}")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 39,
  "body": [
    {
      "type": "FunctionDeclaration",
      "start": 0,
      "end": 39,
      "id": {
        "type": "Identifier",
        "start": 15,
        "end": 19,
        "name": "wrap"
      },
      "generator": false,
      "async": true,
      "params": [],
      "body": {
        "type": "BlockStatement",
        "start": 22,
        "end": 39,
        "body": [
          {
            "type": "ExpressionStatement",
            "start": 24,
            "end": 37,
            "expression": {
              "type": "AssignmentExpression",
              "start": 25,
              "end": 36,
              "operator": "=",
              "left": {
                "type": "Identifier",
                "start": 25,
                "end": 26,
                "name": "a"
              },
              "right": {
                "type": "AwaitExpression",
                "start": 29,
                "end": 36,
                "argument": {
                  "type": "Identifier",
                  "start": 35,
                  "end": 36,
                  "name": "b"
                }
              }
            }
          }
        ]
      }
    }
  ]
}
`, ast)
}

func TestAsyncAwait52(t *testing.T) {
	ast, err := Compile("async function wrap() {\n({a = await b} = obj)\n}")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 47,
  "body": [
    {
      "type": "FunctionDeclaration",
      "start": 0,
      "end": 47,
      "id": {
        "type": "Identifier",
        "start": 15,
        "end": 19,
        "name": "wrap"
      },
      "generator": false,
      "async": true,
      "params": [],
      "body": {
        "type": "BlockStatement",
        "start": 22,
        "end": 47,
        "body": [
          {
            "type": "ExpressionStatement",
            "start": 24,
            "end": 45,
            "expression": {
              "type": "AssignmentExpression",
              "start": 25,
              "end": 44,
              "operator": "=",
              "left": {
                "type": "ObjectPattern",
                "start": 25,
                "end": 38,
                "properties": [
                  {
                    "type": "Property",
                    "start": 26,
                    "end": 37,
                    "method": false,
                    "shorthand": true,
                    "computed": false,
                    "key": {
                      "type": "Identifier",
                      "start": 26,
                      "end": 27,
                      "name": "a"
                    },
                    "kind": "init",
                    "value": {
                      "type": "AssignmentPattern",
                      "start": 26,
                      "end": 37,
                      "left": {
                        "type": "Identifier",
                        "start": 26,
                        "end": 27,
                        "name": "a"
                      },
                      "right": {
                        "type": "AwaitExpression",
                        "start": 30,
                        "end": 37,
                        "argument": {
                          "type": "Identifier",
                          "start": 36,
                          "end": 37,
                          "name": "b"
                        }
                      }
                    }
                  }
                ]
              },
              "right": {
                "type": "Identifier",
                "start": 41,
                "end": 44,
                "name": "obj"
              }
            }
          }
        ]
      }
    }
  ]
}
`, ast)
}

func TestAsyncAwait53(t *testing.T) {
	ast, err := Compile("async function f() { for await (x of xs); }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 43,
  "body": [
    {
      "type": "FunctionDeclaration",
      "start": 0,
      "end": 43,
      "id": {
        "type": "Identifier",
        "start": 15,
        "end": 16,
        "name": "f"
      },
      "generator": false,
      "async": true,
      "params": [],
      "body": {
        "type": "BlockStatement",
        "start": 19,
        "end": 43,
        "body": [
          {
            "type": "ForOfStatement",
            "start": 21,
            "end": 41,
            "await": true,
            "left": {
              "type": "Identifier",
              "start": 32,
              "end": 33,
              "name": "x"
            },
            "right": {
              "type": "Identifier",
              "start": 37,
              "end": 39,
              "name": "xs"
            },
            "body": {
              "type": "EmptyStatement",
              "start": 40,
              "end": 41
            }
          }
        ]
      }
    }
  ]
}
`, ast)
}

func TestAsyncAwait54(t *testing.T) {
	ast, err := Compile("f = ({ w = counter(), x = counter(), y = counter(), z = counter() } = { w: null, x: 0, y: false, z: '' }) => {}")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 111,
  "body": [
    {
      "type": "ExpressionStatement",
      "start": 0,
      "end": 111,
      "expression": {
        "type": "AssignmentExpression",
        "start": 0,
        "end": 111,
        "operator": "=",
        "left": {
          "type": "Identifier",
          "start": 0,
          "end": 1,
          "name": "f"
        },
        "right": {
          "type": "ArrowFunctionExpression",
          "start": 4,
          "end": 111,
          "id": null,
          "expression": false,
          "generator": false,
          "async": false,
          "params": [
            {
              "type": "AssignmentPattern",
              "start": 5,
              "end": 104,
              "left": {
                "type": "ObjectPattern",
                "start": 5,
                "end": 67,
                "properties": [
                  {
                    "type": "Property",
                    "start": 7,
                    "end": 20,
                    "method": false,
                    "shorthand": true,
                    "computed": false,
                    "key": {
                      "type": "Identifier",
                      "start": 7,
                      "end": 8,
                      "name": "w"
                    },
                    "kind": "init",
                    "value": {
                      "type": "AssignmentPattern",
                      "start": 7,
                      "end": 20,
                      "left": {
                        "type": "Identifier",
                        "start": 7,
                        "end": 8,
                        "name": "w"
                      },
                      "right": {
                        "type": "CallExpression",
                        "start": 11,
                        "end": 20,
                        "callee": {
                          "type": "Identifier",
                          "start": 11,
                          "end": 18,
                          "name": "counter"
                        },
                        "arguments": []
                      }
                    }
                  },
                  {
                    "type": "Property",
                    "start": 22,
                    "end": 35,
                    "method": false,
                    "shorthand": true,
                    "computed": false,
                    "key": {
                      "type": "Identifier",
                      "start": 22,
                      "end": 23,
                      "name": "x"
                    },
                    "kind": "init",
                    "value": {
                      "type": "AssignmentPattern",
                      "start": 22,
                      "end": 35,
                      "left": {
                        "type": "Identifier",
                        "start": 22,
                        "end": 23,
                        "name": "x"
                      },
                      "right": {
                        "type": "CallExpression",
                        "start": 26,
                        "end": 35,
                        "callee": {
                          "type": "Identifier",
                          "start": 26,
                          "end": 33,
                          "name": "counter"
                        },
                        "arguments": []
                      }
                    }
                  },
                  {
                    "type": "Property",
                    "start": 37,
                    "end": 50,
                    "method": false,
                    "shorthand": true,
                    "computed": false,
                    "key": {
                      "type": "Identifier",
                      "start": 37,
                      "end": 38,
                      "name": "y"
                    },
                    "kind": "init",
                    "value": {
                      "type": "AssignmentPattern",
                      "start": 37,
                      "end": 50,
                      "left": {
                        "type": "Identifier",
                        "start": 37,
                        "end": 38,
                        "name": "y"
                      },
                      "right": {
                        "type": "CallExpression",
                        "start": 41,
                        "end": 50,
                        "callee": {
                          "type": "Identifier",
                          "start": 41,
                          "end": 48,
                          "name": "counter"
                        },
                        "arguments": []
                      }
                    }
                  },
                  {
                    "type": "Property",
                    "start": 52,
                    "end": 65,
                    "method": false,
                    "shorthand": true,
                    "computed": false,
                    "key": {
                      "type": "Identifier",
                      "start": 52,
                      "end": 53,
                      "name": "z"
                    },
                    "kind": "init",
                    "value": {
                      "type": "AssignmentPattern",
                      "start": 52,
                      "end": 65,
                      "left": {
                        "type": "Identifier",
                        "start": 52,
                        "end": 53,
                        "name": "z"
                      },
                      "right": {
                        "type": "CallExpression",
                        "start": 56,
                        "end": 65,
                        "callee": {
                          "type": "Identifier",
                          "start": 56,
                          "end": 63,
                          "name": "counter"
                        },
                        "arguments": []
                      }
                    }
                  }
                ]
              },
              "right": {
                "type": "ObjectExpression",
                "start": 70,
                "end": 104,
                "properties": [
                  {
                    "type": "Property",
                    "start": 72,
                    "end": 79,
                    "method": false,
                    "shorthand": false,
                    "computed": false,
                    "key": {
                      "type": "Identifier",
                      "start": 72,
                      "end": 73,
                      "name": "w"
                    },
                    "value": {
                      "type": "Literal",
                      "start": 75,
                      "end": 79,
                      "value": null
                    },
                    "kind": "init"
                  },
                  {
                    "type": "Property",
                    "start": 81,
                    "end": 85,
                    "method": false,
                    "shorthand": false,
                    "computed": false,
                    "key": {
                      "type": "Identifier",
                      "start": 81,
                      "end": 82,
                      "name": "x"
                    },
                    "value": {
                      "type": "Literal",
                      "start": 84,
                      "end": 85,
                      "value": 0
                    },
                    "kind": "init"
                  },
                  {
                    "type": "Property",
                    "start": 87,
                    "end": 95,
                    "method": false,
                    "shorthand": false,
                    "computed": false,
                    "key": {
                      "type": "Identifier",
                      "start": 87,
                      "end": 88,
                      "name": "y"
                    },
                    "value": {
                      "type": "Literal",
                      "start": 90,
                      "end": 95,
                      "value": false
                    },
                    "kind": "init"
                  },
                  {
                    "type": "Property",
                    "start": 97,
                    "end": 102,
                    "method": false,
                    "shorthand": false,
                    "computed": false,
                    "key": {
                      "type": "Identifier",
                      "start": 97,
                      "end": 98,
                      "name": "z"
                    },
                    "value": {
                      "type": "Literal",
                      "start": 100,
                      "end": 102,
                      "value": "",
                      "raw": "''"
                    },
                    "kind": "init"
                  }
                ]
              }
            }
          ],
          "body": {
            "type": "BlockStatement",
            "start": 109,
            "end": 111,
            "body": []
          }
        }
      }
    }
  ]
}
`, ast)
}

func TestAsyncAwait55(t *testing.T) {
	ast, err := Compile("({ async: true })")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 17,
  "body": [
    {
      "type": "ExpressionStatement",
      "start": 0,
      "end": 17,
      "expression": {
        "type": "ObjectExpression",
        "start": 1,
        "end": 16,
        "properties": [
          {
            "type": "Property",
            "start": 3,
            "end": 14,
            "method": false,
            "shorthand": false,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 3,
              "end": 8,
              "name": "async"
            },
            "value": {
              "type": "Literal",
              "start": 10,
              "end": 14,
              "value": true
            },
            "kind": "init"
          }
        ]
      }
    }
  ]
}
`, ast)
}

func TestAsyncAwait56(t *testing.T) {
	ast, err := Compile("({async})")
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
        "type": "ObjectExpression",
        "start": 1,
        "end": 8,
        "properties": [
          {
            "type": "Property",
            "start": 2,
            "end": 7,
            "method": false,
            "shorthand": true,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 2,
              "end": 7,
              "name": "async"
            },
            "kind": "init",
            "value": {
              "type": "Identifier",
              "start": 2,
              "end": 7,
              "name": "async"
            }
          }
        ]
      }
    }
  ]
}
`, ast)
}

func TestAsyncAwait57(t *testing.T) {
	ast, err := Compile("({async, foo})")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 14,
  "body": [
    {
      "type": "ExpressionStatement",
      "start": 0,
      "end": 14,
      "expression": {
        "type": "ObjectExpression",
        "start": 1,
        "end": 13,
        "properties": [
          {
            "type": "Property",
            "start": 2,
            "end": 7,
            "method": false,
            "shorthand": true,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 2,
              "end": 7,
              "name": "async"
            },
            "kind": "init",
            "value": {
              "type": "Identifier",
              "start": 2,
              "end": 7,
              "name": "async"
            }
          },
          {
            "type": "Property",
            "start": 9,
            "end": 12,
            "method": false,
            "shorthand": true,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 9,
              "end": 12,
              "name": "foo"
            },
            "kind": "init",
            "value": {
              "type": "Identifier",
              "start": 9,
              "end": 12,
              "name": "foo"
            }
          }
        ]
      }
    }
  ]
}
`, ast)
}

func TestAsyncAwait58(t *testing.T) {
	ast, err := Compile("({async = 0} = {})")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 18,
  "body": [
    {
      "type": "ExpressionStatement",
      "start": 0,
      "end": 18,
      "expression": {
        "type": "AssignmentExpression",
        "start": 1,
        "end": 17,
        "operator": "=",
        "left": {
          "type": "ObjectPattern",
          "start": 1,
          "end": 12,
          "properties": [
            {
              "type": "Property",
              "start": 2,
              "end": 11,
              "method": false,
              "shorthand": true,
              "computed": false,
              "key": {
                "type": "Identifier",
                "start": 2,
                "end": 7,
                "name": "async"
              },
              "kind": "init",
              "value": {
                "type": "AssignmentPattern",
                "start": 2,
                "end": 11,
                "left": {
                  "type": "Identifier",
                  "start": 2,
                  "end": 7,
                  "name": "async"
                },
                "right": {
                  "type": "Literal",
                  "start": 10,
                  "end": 11,
                  "value": 0,
                  "raw": "0"
                }
              }
            }
          ]
        },
        "right": {
          "type": "ObjectExpression",
          "start": 15,
          "end": 17,
          "properties": []
        }
      }
    }
  ]
}
`, ast)
}

func TestAsyncAwait59(t *testing.T) {
	ast, err := Compile("({async \"foo\"(){}})")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 19,
  "body": [
    {
      "type": "ExpressionStatement",
      "start": 0,
      "end": 19,
      "expression": {
        "type": "ObjectExpression",
        "start": 1,
        "end": 18,
        "properties": [
          {
            "type": "Property",
            "start": 2,
            "end": 17,
            "method": true,
            "shorthand": false,
            "computed": false,
            "key": {
              "type": "Literal",
              "start": 8,
              "end": 13,
              "value": "foo",
              "raw": "\"foo\""
            },
            "kind": "init",
            "value": {
              "type": "FunctionExpression",
              "start": 13,
              "end": 17,
              "id": null,
              "expression": false,
              "generator": false,
              "async": true,
              "params": [],
              "body": {
                "type": "BlockStatement",
                "start": 15,
                "end": 17,
                "body": []
              }
            }
          }
        ]
      }
    }
  ]
}
`, ast)
}

func TestAsyncAwait60(t *testing.T) {
	ast, err := Compile("({async 'foo'(){}})")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 19,
  "body": [
    {
      "type": "ExpressionStatement",
      "start": 0,
      "end": 19,
      "expression": {
        "type": "ObjectExpression",
        "start": 1,
        "end": 18,
        "properties": [
          {
            "type": "Property",
            "start": 2,
            "end": 17,
            "method": true,
            "shorthand": false,
            "computed": false,
            "key": {
              "type": "Literal",
              "start": 8,
              "end": 13,
              "value": "foo",
              "raw": "'foo'"
            },
            "kind": "init",
            "value": {
              "type": "FunctionExpression",
              "start": 13,
              "end": 17,
              "id": null,
              "expression": false,
              "generator": false,
              "async": true,
              "params": [],
              "body": {
                "type": "BlockStatement",
                "start": 15,
                "end": 17,
                "body": []
              }
            }
          }
        ]
      }
    }
  ]
}
`, ast)
}

func TestAsyncAwait61(t *testing.T) {
	ast, err := Compile("({async 100(){}})")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 17,
  "body": [
    {
      "type": "ExpressionStatement",
      "start": 0,
      "end": 17,
      "expression": {
        "type": "ObjectExpression",
        "start": 1,
        "end": 16,
        "properties": [
          {
            "type": "Property",
            "start": 2,
            "end": 15,
            "method": true,
            "shorthand": false,
            "computed": false,
            "key": {
              "type": "Literal",
              "start": 8,
              "end": 11,
              "value": 100,
              "raw": "100"
            },
            "kind": "init",
            "value": {
              "type": "FunctionExpression",
              "start": 11,
              "end": 15,
              "id": null,
              "expression": false,
              "generator": false,
              "async": true,
              "params": [],
              "body": {
                "type": "BlockStatement",
                "start": 13,
                "end": 15,
                "body": []
              }
            }
          }
        ]
      }
    }
  ]
}
`, ast)
}

func TestAsyncAwait62(t *testing.T) {
	ast, err := Compile("({async [foo](){}})")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 19,
  "body": [
    {
      "type": "ExpressionStatement",
      "start": 0,
      "end": 19,
      "expression": {
        "type": "ObjectExpression",
        "start": 1,
        "end": 18,
        "properties": [
          {
            "type": "Property",
            "start": 2,
            "end": 17,
            "method": true,
            "shorthand": false,
            "computed": true,
            "key": {
              "type": "Identifier",
              "start": 9,
              "end": 12,
              "name": "foo"
            },
            "kind": "init",
            "value": {
              "type": "FunctionExpression",
              "start": 13,
              "end": 17,
              "id": null,
              "expression": false,
              "generator": false,
              "async": true,
              "params": [],
              "body": {
                "type": "BlockStatement",
                "start": 15,
                "end": 17,
                "body": []
              }
            }
          }
        ]
      }
    }
  ]
}
`, ast)
}

func TestAsyncAwait63(t *testing.T) {
	ast, err := Compile("({ async delete() {} })")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 23,
  "body": [
    {
      "type": "ExpressionStatement",
      "start": 0,
      "end": 23,
      "expression": {
        "type": "ObjectExpression",
        "start": 1,
        "end": 22,
        "properties": [
          {
            "type": "Property",
            "start": 3,
            "end": 20,
            "method": true,
            "shorthand": false,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 9,
              "end": 15,
              "name": "delete"
            },
            "kind": "init",
            "value": {
              "type": "FunctionExpression",
              "start": 15,
              "end": 20,
              "id": null,
              "expression": false,
              "generator": false,
              "async": true,
              "params": [],
              "body": {
                "type": "BlockStatement",
                "start": 18,
                "end": 20,
                "body": []
              }
            }
          }
        ]
      }
    }
  ]
}
`, ast)
}

func TestAsyncAwait64(t *testing.T) {
	ast, err := Compile("(async() => { await (4 ** 2) })()")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 33,
  "body": [
    {
      "type": "ExpressionStatement",
      "start": 0,
      "end": 33,
      "expression": {
        "type": "CallExpression",
        "start": 0,
        "end": 33,
        "callee": {
          "type": "ArrowFunctionExpression",
          "start": 1,
          "end": 30,
          "id": null,
          "expression": false,
          "generator": false,
          "async": true,
          "params": [],
          "body": {
            "type": "BlockStatement",
            "start": 12,
            "end": 30,
            "body": [
              {
                "type": "ExpressionStatement",
                "start": 14,
                "end": 28,
                "expression": {
                  "type": "AwaitExpression",
                  "start": 14,
                  "end": 28,
                  "argument": {
                    "type": "BinaryExpression",
                    "start": 21,
                    "end": 27,
                    "left": {
                      "type": "Literal",
                      "start": 21,
                      "end": 22,
                      "value": 4,
                      "raw": "4"
                    },
                    "operator": "**",
                    "right": {
                      "type": "Literal",
                      "start": 26,
                      "end": 27,
                      "value": 2,
                      "raw": "2"
                    }
                  }
                }
              }
            ]
          }
        },
        "arguments": []
      }
    }
  ]
}
`, ast)
}

func TestAsyncAwait65(t *testing.T) {
	ast, err := Compile("async() => (await (1 ** 3))")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 27,
  "body": [
    {
      "type": "ExpressionStatement",
      "start": 0,
      "end": 27,
      "expression": {
        "type": "ArrowFunctionExpression",
        "start": 0,
        "end": 27,
        "id": null,
        "expression": true,
        "generator": false,
        "async": true,
        "params": [],
        "body": {
          "type": "AwaitExpression",
          "start": 12,
          "end": 26,
          "argument": {
            "type": "BinaryExpression",
            "start": 19,
            "end": 25,
            "left": {
              "type": "Literal",
              "start": 19,
              "end": 20,
              "value": 1,
              "raw": "1"
            },
            "operator": "**",
            "right": {
              "type": "Literal",
              "start": 24,
              "end": 25,
              "value": 3,
              "raw": "3"
            }
          }
        }
      }
    }
  ]
}
`, ast)
}

func TestAsyncAwait66(t *testing.T) {
	ast, err := Compile("async() => await (5 ** 6)")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 25,
  "body": [
    {
      "type": "ExpressionStatement",
      "start": 0,
      "end": 25,
      "expression": {
        "type": "ArrowFunctionExpression",
        "start": 0,
        "end": 25,
        "id": null,
        "expression": true,
        "generator": false,
        "async": true,
        "params": [],
        "body": {
          "type": "AwaitExpression",
          "start": 11,
          "end": 25,
          "argument": {
            "type": "BinaryExpression",
            "start": 18,
            "end": 24,
            "left": {
              "type": "Literal",
              "start": 18,
              "end": 19,
              "value": 5,
              "raw": "5"
            },
            "operator": "**",
            "right": {
              "type": "Literal",
              "start": 23,
              "end": 24,
              "value": 6,
              "raw": "6"
            }
          }
        }
      }
    }
  ]
}
`, ast)
}

func TestAsyncAwait67(t *testing.T) {
	ast, err := Compile("async* ({a: b = c})")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 19,
  "body": [
    {
      "type": "ExpressionStatement",
      "start": 0,
      "end": 19,
      "expression": {
        "type": "BinaryExpression",
        "start": 0,
        "end": 19,
        "left": {
          "type": "Identifier",
          "start": 0,
          "end": 5,
          "name": "async"
        },
        "operator": "*",
        "right": {
          "type": "ObjectExpression",
          "start": 8,
          "end": 18,
          "properties": [
            {
              "type": "Property",
              "start": 9,
              "end": 17,
              "method": false,
              "shorthand": false,
              "computed": false,
              "key": {
                "type": "Identifier",
                "start": 9,
                "end": 10,
                "name": "a"
              },
              "value": {
                "type": "AssignmentExpression",
                "start": 12,
                "end": 17,
                "operator": "=",
                "left": {
                  "type": "Identifier",
                  "start": 12,
                  "end": 13,
                  "name": "b"
                },
                "right": {
                  "type": "Identifier",
                  "start": 16,
                  "end": 17,
                  "name": "c"
                }
              },
              "kind": "init"
            }
          ]
        }
      }
    }
  ]
}
`, ast)
}

func TestAsyncAwait68(t *testing.T) {
	// 	ast, err := Compile("async function f() { for await (x of xs); }")
	// 	AssertEqual(t, nil, err, "should be prog ok")

	// 	AssertEqualJson(t, `

	// `, ast)
}

func TestAsyncAwait69(t *testing.T) {
	// 	ast, err := Compile("async function f() { for await (x of xs); }")
	// 	AssertEqual(t, nil, err, "should be prog ok")

	// 	AssertEqualJson(t, `

	// `, ast)
}
