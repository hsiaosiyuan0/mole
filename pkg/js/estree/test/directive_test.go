package estree_test

import (
	"testing"

	"github.com/hsiaosiyuan0/mole/pkg/assert"
	"github.com/hsiaosiyuan0/mole/pkg/js/parser"
)

// No directives
func TestDirective1(t *testing.T) {
	ast, err := compile("foo")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 3,
  "body": [
    {
      "type": "ExpressionStatement",
      "start": 0,
      "end": 3,
      "expression": {
        "type": "Identifier",
        "start": 0,
        "end": 3,
        "name": "foo"
      },
      "directive": null
    }
  ]
}
`, ast)
}

func TestDirective2(t *testing.T) {
	ast, err := compile("function wrap() { foo }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 23,
  "body": [
    {
      "type": "FunctionDeclaration",
      "start": 0,
      "end": 23,
      "id": {
        "type": "Identifier",
        "start": 9,
        "end": 13,
        "name": "wrap"
      },
      "generator": false,
      "async": false,
      "params": [],
      "body": {
        "type": "BlockStatement",
        "start": 16,
        "end": 23,
        "body": [
          {
            "type": "ExpressionStatement",
            "start": 18,
            "end": 21,
            "expression": {
              "type": "Identifier",
              "start": 18,
              "end": 21,
              "name": "foo"
            },
            "directive": null
          }
        ]
      }
    }
  ]
}
`, ast)
}

func TestDirective3(t *testing.T) {
	ast, err := compile("!function wrap() { foo }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 24,
  "body": [
    {
      "type": "ExpressionStatement",
      "start": 0,
      "end": 24,
      "expression": {
        "type": "UnaryExpression",
        "start": 0,
        "end": 24,
        "operator": "!",
        "prefix": true,
        "argument": {
          "type": "FunctionExpression",
          "start": 1,
          "end": 24,
          "id": {
            "type": "Identifier",
            "start": 10,
            "end": 14,
            "name": "wrap"
          },
          "expression": false,
          "generator": false,
          "async": false,
          "params": [],
          "body": {
            "type": "BlockStatement",
            "start": 17,
            "end": 24,
            "body": [
              {
                "type": "ExpressionStatement",
                "start": 19,
                "end": 22,
                "expression": {
                  "type": "Identifier",
                  "start": 19,
                  "end": 22,
                  "name": "foo"
                }
              }
            ]
          }
        }
      },
      "directive": null
    }
  ]
}
`, ast)
}

func TestDirective4(t *testing.T) {
	ast, err := compile("() => { foo }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
        "expression": false,
        "generator": false,
        "async": false,
        "params": [],
        "body": {
          "type": "BlockStatement",
          "start": 6,
          "end": 13,
          "body": [
            {
              "type": "ExpressionStatement",
              "start": 8,
              "end": 11,
              "expression": {
                "type": "Identifier",
                "start": 8,
                "end": 11,
                "name": "foo"
              }
            }
          ]
        }
      },
      "directive": null
    }
  ]
}
`, ast)
}

func TestDirective5(t *testing.T) {
	ast, err := compile("100")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 3,
  "body": [
    {
      "type": "ExpressionStatement",
      "start": 0,
      "end": 3,
      "expression": {
        "type": "Literal",
        "start": 0,
        "end": 3,
        "value": 100,
        "raw": "100"
      },
      "directive": null
    }
  ]
}
`, ast)
}

func TestDirective6(t *testing.T) {
	ast, err := compile("\"use strict\" + 1")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
        "type": "BinaryExpression",
        "start": 0,
        "end": 16,
        "left": {
          "type": "Literal",
          "start": 0,
          "end": 12,
          "value": "use strict",
          "raw": "\"use strict\""
        },
        "operator": "+",
        "right": {
          "type": "Literal",
          "start": 15,
          "end": 16,
          "value": 1,
          "raw": "1"
        }
      },
      "directive": null
    }
  ]
}
`, ast)
}

func TestDirective7(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_STRICT)
	ast, err := compileWithOpts("; 'use strict'; with ({}) {}", opts)
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 28,
  "body": [
    {
      "type": "EmptyStatement",
      "start": 0,
      "end": 1
    },
    {
      "type": "ExpressionStatement",
      "start": 2,
      "end": 15,
      "expression": {
        "type": "Literal",
        "start": 2,
        "end": 14,
        "value": "use strict",
        "raw": "'use strict'"
      },
      "directive": null
    },
    {
      "type": "WithStatement",
      "start": 16,
      "end": 28,
      "object": {
        "type": "ObjectExpression",
        "start": 22,
        "end": 24,
        "properties": []
      },
      "body": {
        "type": "BlockStatement",
        "start": 26,
        "end": 28,
        "body": []
      }
    }
  ]
}
`, ast)
}

// One directive
func TestDirective8(t *testing.T) {
	ast, err := compile("\"use strict\"\n foo")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 17,
  "body": [
    {
      "type": "ExpressionStatement",
      "start": 0,
      "end": 12,
      "expression": {
        "type": "Literal",
        "start": 0,
        "end": 12,
        "value": "use strict",
        "raw": "\"use strict\""
      },
      "directive": "use strict"
    },
    {
      "type": "ExpressionStatement",
      "start": 14,
      "end": 17,
      "expression": {
        "type": "Identifier",
        "start": 14,
        "end": 17,
        "name": "foo"
      }
    }
  ]
}
`, ast)
}

func TestDirective9(t *testing.T) {
	ast, err := compile("'use strict'; foo")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 17,
  "body": [
    {
      "type": "ExpressionStatement",
      "start": 0,
      "end": 13,
      "expression": {
        "type": "Literal",
        "start": 0,
        "end": 12,
        "value": "use strict",
        "raw": "'use strict'"
      },
      "directive": "use strict"
    },
    {
      "type": "ExpressionStatement",
      "start": 14,
      "end": 17,
      "expression": {
        "type": "Identifier",
        "start": 14,
        "end": 17,
        "name": "foo"
      }
    }
  ]
}
`, ast)
}

func TestDirective10(t *testing.T) {
	ast, err := compile("function wrap() { \"use strict\"\n foo }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 37,
  "body": [
    {
      "type": "FunctionDeclaration",
      "start": 0,
      "end": 37,
      "id": {
        "type": "Identifier",
        "start": 9,
        "end": 13,
        "name": "wrap"
      },
      "generator": false,
      "async": false,
      "params": [],
      "body": {
        "type": "BlockStatement",
        "start": 16,
        "end": 37,
        "body": [
          {
            "type": "ExpressionStatement",
            "start": 18,
            "end": 30,
            "expression": {
              "type": "Literal",
              "start": 18,
              "end": 30,
              "value": "use strict",
              "raw": "\"use strict\""
            },
            "directive": "use strict"
          },
          {
            "type": "ExpressionStatement",
            "start": 32,
            "end": 35,
            "expression": {
              "type": "Identifier",
              "start": 32,
              "end": 35,
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

func TestDirective11(t *testing.T) {
	ast, err := compile("!function wrap() { \"use strict\"\n foo }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 38,
  "body": [
    {
      "type": "ExpressionStatement",
      "start": 0,
      "end": 38,
      "expression": {
        "type": "UnaryExpression",
        "start": 0,
        "end": 38,
        "operator": "!",
        "prefix": true,
        "argument": {
          "type": "FunctionExpression",
          "start": 1,
          "end": 38,
          "id": {
            "type": "Identifier",
            "start": 10,
            "end": 14,
            "name": "wrap"
          },
          "expression": false,
          "generator": false,
          "async": false,
          "params": [],
          "body": {
            "type": "BlockStatement",
            "start": 17,
            "end": 38,
            "body": [
              {
                "type": "ExpressionStatement",
                "start": 19,
                "end": 31,
                "expression": {
                  "type": "Literal",
                  "start": 19,
                  "end": 31,
                  "value": "use strict",
                  "raw": "\"use strict\""
                },
                "directive": "use strict"
              },
              {
                "type": "ExpressionStatement",
                "start": 33,
                "end": 36,
                "expression": {
                  "type": "Identifier",
                  "start": 33,
                  "end": 36,
                  "name": "foo"
                }
              }
            ]
          }
        }
      }
    }
  ]
}
`, ast)
}

func TestDirective12(t *testing.T) {
	ast, err := compile("() => { \"use strict\"\n foo }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
        "expression": false,
        "generator": false,
        "async": false,
        "params": [],
        "body": {
          "type": "BlockStatement",
          "start": 6,
          "end": 27,
          "body": [
            {
              "type": "ExpressionStatement",
              "start": 8,
              "end": 20,
              "expression": {
                "type": "Literal",
                "start": 8,
                "end": 20,
                "value": "use strict",
                "raw": "\"use strict\""
              },
              "directive": "use strict"
            },
            {
              "type": "ExpressionStatement",
              "start": 22,
              "end": 25,
              "expression": {
                "type": "Identifier",
                "start": 22,
                "end": 25,
                "name": "foo"
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

func TestDirective13(t *testing.T) {
	ast, err := compile("() => \"use strict\"")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
        "type": "ArrowFunctionExpression",
        "start": 0,
        "end": 18,
        "id": null,
        "expression": true,
        "generator": false,
        "async": false,
        "params": [],
        "body": {
          "type": "Literal",
          "start": 6,
          "end": 18,
          "value": "use strict",
          "raw": "\"use strict\""
        }
      }
    }
  ]
}
`, ast)
}

func TestDirective14(t *testing.T) {
	ast, err := compile("({ wrap() { \"use strict\"; foo } })")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 34,
  "body": [
    {
      "type": "ExpressionStatement",
      "start": 0,
      "end": 34,
      "expression": {
        "type": "ObjectExpression",
        "start": 1,
        "end": 33,
        "properties": [
          {
            "type": "Property",
            "start": 3,
            "end": 31,
            "method": true,
            "shorthand": false,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 3,
              "end": 7,
              "name": "wrap"
            },
            "kind": "init",
            "value": {
              "type": "FunctionExpression",
              "start": 7,
              "end": 31,
              "id": null,
              "expression": false,
              "generator": false,
              "async": false,
              "params": [],
              "body": {
                "type": "BlockStatement",
                "start": 10,
                "end": 31,
                "body": [
                  {
                    "type": "ExpressionStatement",
                    "start": 12,
                    "end": 25,
                    "expression": {
                      "type": "Literal",
                      "start": 12,
                      "end": 24,
                      "value": "use strict",
                      "raw": "\"use strict\""
                    },
                    "directive": "use strict"
                  },
                  {
                    "type": "ExpressionStatement",
                    "start": 26,
                    "end": 29,
                    "expression": {
                      "type": "Identifier",
                      "start": 26,
                      "end": 29,
                      "name": "foo"
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

func TestDirective15(t *testing.T) {
	ast, err := compile("(class { wrap() { \"use strict\"; foo } })")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 40,
  "body": [
    {
      "type": "ExpressionStatement",
      "start": 0,
      "end": 40,
      "expression": {
        "type": "ClassExpression",
        "start": 1,
        "end": 39,
        "id": null,
        "superClass": null,
        "body": {
          "type": "ClassBody",
          "start": 7,
          "end": 39,
          "body": [
            {
              "type": "MethodDefinition",
              "start": 9,
              "end": 37,
              "static": false,
              "computed": false,
              "key": {
                "type": "Identifier",
                "start": 9,
                "end": 13,
                "name": "wrap"
              },
              "kind": "method",
              "value": {
                "type": "FunctionExpression",
                "start": 13,
                "end": 37,
                "id": null,
                "expression": false,
                "generator": false,
                "async": false,
                "params": [],
                "body": {
                  "type": "BlockStatement",
                  "start": 16,
                  "end": 37,
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
                        "raw": "\"use strict\""
                      },
                      "directive": "use strict"
                    },
                    {
                      "type": "ExpressionStatement",
                      "start": 32,
                      "end": 35,
                      "expression": {
                        "type": "Identifier",
                        "start": 32,
                        "end": 35,
                        "name": "foo"
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
  ]
}
`, ast)
}

// Should not decode escape sequence.
func TestDirective16(t *testing.T) {
	ast, err := compile("\"\\u0075se strict\"")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
        "type": "Literal",
        "start": 0,
        "end": 17,
        "value": "use strict",
        "raw": "\"\\u0075se strict\""
      },
      "directive": "\\u0075se strict"
    }
  ]
}
`, ast)
}

// Two or more directives.
func TestDirective17(t *testing.T) {
	ast, err := compile("\"use asm\"; \"use strict\"; foo")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 28,
  "body": [
    {
      "type": "ExpressionStatement",
      "start": 0,
      "end": 10,
      "expression": {
        "type": "Literal",
        "start": 0,
        "end": 9,
        "value": "use asm",
        "raw": "\"use asm\""
      },
      "directive": "use asm"
    },
    {
      "type": "ExpressionStatement",
      "start": 11,
      "end": 24,
      "expression": {
        "type": "Literal",
        "start": 11,
        "end": 23,
        "value": "use strict",
        "raw": "\"use strict\""
      },
      "directive": "use strict"
    },
    {
      "type": "ExpressionStatement",
      "start": 25,
      "end": 28,
      "expression": {
        "type": "Identifier",
        "start": 25,
        "end": 28,
        "name": "foo"
      }
    }
  ]
}
`, ast)
}

func TestDirective18(t *testing.T) {
	ast, err := compile("function wrap() { \"use asm\"; \"use strict\"; foo }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 48,
  "body": [
    {
      "type": "FunctionDeclaration",
      "start": 0,
      "end": 48,
      "id": {
        "type": "Identifier",
        "start": 9,
        "end": 13,
        "name": "wrap"
      },
      "generator": false,
      "async": false,
      "params": [],
      "body": {
        "type": "BlockStatement",
        "start": 16,
        "end": 48,
        "body": [
          {
            "type": "ExpressionStatement",
            "start": 18,
            "end": 28,
            "expression": {
              "type": "Literal",
              "start": 18,
              "end": 27,
              "value": "use asm",
              "raw": "\"use asm\""
            },
            "directive": "use asm"
          },
          {
            "type": "ExpressionStatement",
            "start": 29,
            "end": 42,
            "expression": {
              "type": "Literal",
              "start": 29,
              "end": 41,
              "value": "use strict",
              "raw": "\"use strict\""
            },
            "directive": "use strict"
          },
          {
            "type": "ExpressionStatement",
            "start": 43,
            "end": 46,
            "expression": {
              "type": "Identifier",
              "start": 43,
              "end": 46,
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

// One string after other expressions.
func TestDirective19(t *testing.T) {
	ast, err := compile("\"use strict\"; foo; \"use asm\"")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 28,
  "body": [
    {
      "type": "ExpressionStatement",
      "start": 0,
      "end": 13,
      "expression": {
        "type": "Literal",
        "start": 0,
        "end": 12,
        "value": "use strict",
        "raw": "\"use strict\""
      },
      "directive": "use strict"
    },
    {
      "type": "ExpressionStatement",
      "start": 14,
      "end": 18,
      "expression": {
        "type": "Identifier",
        "start": 14,
        "end": 17,
        "name": "foo"
      }
    },
    {
      "type": "ExpressionStatement",
      "start": 19,
      "end": 28,
      "expression": {
        "type": "Literal",
        "start": 19,
        "end": 28,
        "value": "use asm",
        "raw": "\"use asm\""
      }
    }
  ]
}
`, ast)
}

func TestDirective20(t *testing.T) {
	ast, err := compile("function wrap() { \"use asm\"; foo; \"use strict\" }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 48,
  "body": [
    {
      "type": "FunctionDeclaration",
      "start": 0,
      "end": 48,
      "id": {
        "type": "Identifier",
        "start": 9,
        "end": 13,
        "name": "wrap"
      },
      "generator": false,
      "async": false,
      "params": [],
      "body": {
        "type": "BlockStatement",
        "start": 16,
        "end": 48,
        "body": [
          {
            "type": "ExpressionStatement",
            "start": 18,
            "end": 28,
            "expression": {
              "type": "Literal",
              "start": 18,
              "end": 27,
              "value": "use asm",
              "raw": "\"use asm\""
            },
            "directive": "use asm"
          },
          {
            "type": "ExpressionStatement",
            "start": 29,
            "end": 33,
            "expression": {
              "type": "Identifier",
              "start": 29,
              "end": 32,
              "name": "foo"
            }
          },
          {
            "type": "ExpressionStatement",
            "start": 34,
            "end": 46,
            "expression": {
              "type": "Literal",
              "start": 34,
              "end": 46,
              "value": "use strict",
              "raw": "\"use strict\""
            }
          }
        ]
      }
    }
  ]
}
`, ast)
}

// One string in a block.
func TestDirective21(t *testing.T) {
	ast, err := compile("{ \"use strict\"; }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 17,
  "body": [
    {
      "type": "BlockStatement",
      "start": 0,
      "end": 17,
      "body": [
        {
          "type": "ExpressionStatement",
          "start": 2,
          "end": 15,
          "expression": {
            "type": "Literal",
            "start": 2,
            "end": 14,
            "value": "use strict",
            "raw": "\"use strict\""
          }
        }
      ]
    }
  ]
}
`, ast)
}

// One string in a block.
func TestDirective22(t *testing.T) {
	ast, err := compile("function wrap() { { \"use strict\" } foo }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 40,
  "body": [
    {
      "type": "FunctionDeclaration",
      "start": 0,
      "end": 40,
      "id": {
        "type": "Identifier",
        "start": 9,
        "end": 13,
        "name": "wrap"
      },
      "generator": false,
      "async": false,
      "params": [],
      "body": {
        "type": "BlockStatement",
        "start": 16,
        "end": 40,
        "body": [
          {
            "type": "BlockStatement",
            "start": 18,
            "end": 34,
            "body": [
              {
                "type": "ExpressionStatement",
                "start": 20,
                "end": 32,
                "expression": {
                  "type": "Literal",
                  "start": 20,
                  "end": 32,
                  "value": "use strict",
                  "raw": "\"use strict\""
                }
              }
            ]
          },
          {
            "type": "ExpressionStatement",
            "start": 35,
            "end": 38,
            "expression": {
              "type": "Identifier",
              "start": 35,
              "end": 38,
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

// One string with parentheses.
func TestDirective23(t *testing.T) {
	ast, err := compile("(\"use strict\"); foo")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 19,
  "body": [
    {
      "type": "ExpressionStatement",
      "start": 0,
      "end": 15,
      "expression": {
        "type": "Literal",
        "start": 1,
        "end": 13,
        "value": "use strict",
        "raw": "\"use strict\""
      }
    },
    {
      "type": "ExpressionStatement",
      "start": 16,
      "end": 19,
      "expression": {
        "type": "Identifier",
        "start": 16,
        "end": 19,
        "name": "foo"
      }
    }
  ]
}
`, ast)
}

func TestDirective24(t *testing.T) {
	ast, err := compile("function wrap() { (\"use strict\"); foo }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
        "start": 9,
        "end": 13,
        "name": "wrap"
      },
      "generator": false,
      "async": false,
      "params": [],
      "body": {
        "type": "BlockStatement",
        "start": 16,
        "end": 39,
        "body": [
          {
            "type": "ExpressionStatement",
            "start": 18,
            "end": 33,
            "expression": {
              "type": "Literal",
              "start": 19,
              "end": 31,
              "value": "use strict",
              "raw": "\"use strict\""
            }
          },
          {
            "type": "ExpressionStatement",
            "start": 34,
            "end": 37,
            "expression": {
              "type": "Identifier",
              "start": 34,
              "end": 37,
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

// Complex cases such as the function in a default parameter.
func TestDirective25(t *testing.T) {
	ast, err := compile("function a() { \"use strict\" } \"use strict\"; foo")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 47,
  "body": [
    {
      "type": "FunctionDeclaration",
      "start": 0,
      "end": 29,
      "id": {
        "type": "Identifier",
        "start": 9,
        "end": 10,
        "name": "a"
      },
      "generator": false,
      "async": false,
      "params": [],
      "body": {
        "type": "BlockStatement",
        "start": 13,
        "end": 29,
        "body": [
          {
            "type": "ExpressionStatement",
            "start": 15,
            "end": 27,
            "expression": {
              "type": "Literal",
              "start": 15,
              "end": 27,
              "value": "use strict",
              "raw": "\"use strict\""
            },
            "directive": "use strict"
          }
        ]
      }
    },
    {
      "type": "ExpressionStatement",
      "start": 30,
      "end": 43,
      "expression": {
        "type": "Literal",
        "start": 30,
        "end": 42,
        "value": "use strict",
        "raw": "\"use strict\""
      }
    },
    {
      "type": "ExpressionStatement",
      "start": 44,
      "end": 47,
      "expression": {
        "type": "Identifier",
        "start": 44,
        "end": 47,
        "name": "foo"
      }
    }
  ]
}
`, ast)
}

func TestDirective26(t *testing.T) {
	// below cases are skipped, since they cannot pass the rule:
	// `Illegal 'use strict' directive in function with non-simple parameter list`

	// `function a(a = function() { "use strict"; foo }) { "use strict" }`
	// `(a = () => { "use strict"; foo }) => { "use strict" }`
}
