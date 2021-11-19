package estree

import (
	"testing"

	"github.com/hsiaosiyuan0/mole/pkg/assert"
)

// below are tests referred from:
// https://github.com/acornjs/acorn/blob/134ede4084a6611f2e0d60e676983443d2426405/test/tests-async-iteration.js

// for-await-of

func TestAsyncIteration1(t *testing.T) {
	ast, err := compile("async function f() { for await (x of xs); }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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

func TestAsyncIteration2(t *testing.T) {
	ast, err := compile("async function f() { for await (var x of xs); }")
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
      "end": 47,
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
        "end": 47,
        "body": [
          {
            "type": "ForOfStatement",
            "start": 21,
            "end": 45,
            "await": true,
            "left": {
              "type": "VariableDeclaration",
              "start": 32,
              "end": 37,
              "declarations": [
                {
                  "type": "VariableDeclarator",
                  "start": 36,
                  "end": 37,
                  "id": {
                    "type": "Identifier",
                    "start": 36,
                    "end": 37,
                    "name": "x"
                  },
                  "init": null
                }
              ],
              "kind": "var"
            },
            "right": {
              "type": "Identifier",
              "start": 41,
              "end": 43,
              "name": "xs"
            },
            "body": {
              "type": "EmptyStatement",
              "start": 44,
              "end": 45
            }
          }
        ]
      }
    }
  ]
}
    `, ast)
}

func TestAsyncIteration3(t *testing.T) {
	ast, err := compile("async function f() { for await (let x of xs); }")
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
      "end": 47,
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
        "end": 47,
        "body": [
          {
            "type": "ForOfStatement",
            "start": 21,
            "end": 45,
            "await": true,
            "left": {
              "type": "VariableDeclaration",
              "start": 32,
              "end": 37,
              "declarations": [
                {
                  "type": "VariableDeclarator",
                  "start": 36,
                  "end": 37,
                  "id": {
                    "type": "Identifier",
                    "start": 36,
                    "end": 37,
                    "name": "x"
                  },
                  "init": null
                }
              ],
              "kind": "let"
            },
            "right": {
              "type": "Identifier",
              "start": 41,
              "end": 43,
              "name": "xs"
            },
            "body": {
              "type": "EmptyStatement",
              "start": 44,
              "end": 45
            }
          }
        ]
      }
    }
  ]
}
    `, ast)
}

func TestAsyncIteration4(t *testing.T) {
	ast, err := compile("async function f() { for\nawait (x of xs); }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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

func TestAsyncIteration5(t *testing.T) {
	ast, err := compile("f = async function() { for await (x of xs); }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 45,
  "body": [
    {
      "type": "ExpressionStatement",
      "start": 0,
      "end": 45,
      "expression": {
        "type": "AssignmentExpression",
        "start": 0,
        "end": 45,
        "operator": "=",
        "left": {
          "type": "Identifier",
          "start": 0,
          "end": 1,
          "name": "f"
        },
        "right": {
          "type": "FunctionExpression",
          "start": 4,
          "end": 45,
          "id": null,
          "expression": false,
          "generator": false,
          "async": true,
          "params": [],
          "body": {
            "type": "BlockStatement",
            "start": 21,
            "end": 45,
            "body": [
              {
                "type": "ForOfStatement",
                "start": 23,
                "end": 43,
                "await": true,
                "left": {
                  "type": "Identifier",
                  "start": 34,
                  "end": 35,
                  "name": "x"
                },
                "right": {
                  "type": "Identifier",
                  "start": 39,
                  "end": 41,
                  "name": "xs"
                },
                "body": {
                  "type": "EmptyStatement",
                  "start": 42,
                  "end": 43
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

func TestAsyncIteration6(t *testing.T) {
	ast, err := compile("f = async() => { for await (x of xs); }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 39,
  "body": [
    {
      "type": "ExpressionStatement",
      "start": 0,
      "end": 39,
      "expression": {
        "type": "AssignmentExpression",
        "start": 0,
        "end": 39,
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
          "end": 39,
          "id": null,
          "expression": false,
          "generator": false,
          "async": true,
          "params": [],
          "body": {
            "type": "BlockStatement",
            "start": 15,
            "end": 39,
            "body": [
              {
                "type": "ForOfStatement",
                "start": 17,
                "end": 37,
                "await": true,
                "left": {
                  "type": "Identifier",
                  "start": 28,
                  "end": 29,
                  "name": "x"
                },
                "right": {
                  "type": "Identifier",
                  "start": 33,
                  "end": 35,
                  "name": "xs"
                },
                "body": {
                  "type": "EmptyStatement",
                  "start": 36,
                  "end": 37
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

func TestAsyncIteration7(t *testing.T) {
	ast, err := compile("obj = { async f() { for await (x of xs); } }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 44,
  "body": [
    {
      "type": "ExpressionStatement",
      "start": 0,
      "end": 44,
      "expression": {
        "type": "AssignmentExpression",
        "start": 0,
        "end": 44,
        "operator": "=",
        "left": {
          "type": "Identifier",
          "start": 0,
          "end": 3,
          "name": "obj"
        },
        "right": {
          "type": "ObjectExpression",
          "start": 6,
          "end": 44,
          "properties": [
            {
              "type": "Property",
              "start": 8,
              "end": 42,
              "method": true,
              "shorthand": false,
              "computed": false,
              "key": {
                "type": "Identifier",
                "start": 14,
                "end": 15,
                "name": "f"
              },
              "kind": "init",
              "value": {
                "type": "FunctionExpression",
                "start": 15,
                "end": 42,
                "id": null,
                "expression": false,
                "generator": false,
                "async": true,
                "params": [],
                "body": {
                  "type": "BlockStatement",
                  "start": 18,
                  "end": 42,
                  "body": [
                    {
                      "type": "ForOfStatement",
                      "start": 20,
                      "end": 40,
                      "await": true,
                      "left": {
                        "type": "Identifier",
                        "start": 31,
                        "end": 32,
                        "name": "x"
                      },
                      "right": {
                        "type": "Identifier",
                        "start": 36,
                        "end": 38,
                        "name": "xs"
                      },
                      "body": {
                        "type": "EmptyStatement",
                        "start": 39,
                        "end": 40
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

func TestAsyncIteration8(t *testing.T) {
	ast, err := compile("class A { async f() { for await (x of xs); } }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 46,
  "body": [
    {
      "type": "ClassDeclaration",
      "start": 0,
      "end": 46,
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
        "end": 46,
        "body": [
          {
            "type": "MethodDefinition",
            "start": 10,
            "end": 44,
            "static": false,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 16,
              "end": 17,
              "name": "f"
            },
            "kind": "method",
            "value": {
              "type": "FunctionExpression",
              "start": 17,
              "end": 44,
              "id": null,
              "expression": false,
              "generator": false,
              "async": true,
              "params": [],
              "body": {
                "type": "BlockStatement",
                "start": 20,
                "end": 44,
                "body": [
                  {
                    "type": "ForOfStatement",
                    "start": 22,
                    "end": 42,
                    "await": true,
                    "left": {
                      "type": "Identifier",
                      "start": 33,
                      "end": 34,
                      "name": "x"
                    },
                    "right": {
                      "type": "Identifier",
                      "start": 38,
                      "end": 40,
                      "name": "xs"
                    },
                    "body": {
                      "type": "EmptyStatement",
                      "start": 41,
                      "end": 42
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

func TestAsyncIteration9(t *testing.T) {
	ast, err := compile("for (x of xs);")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 14,
  "body": [
    {
      "type": "ForOfStatement",
      "start": 0,
      "end": 14,
      "await": false,
      "left": {
        "type": "Identifier",
        "start": 5,
        "end": 6,
        "name": "x"
      },
      "right": {
        "type": "Identifier",
        "start": 10,
        "end": 12,
        "name": "xs"
      },
      "body": {
        "type": "EmptyStatement",
        "start": 13,
        "end": 14
      }
    }
  ]
}
`, ast)
}

func TestAsyncIteration10(t *testing.T) {
	ast, err := compile("for (x in xs);")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 14,
  "body": [
    {
      "type": "ForInStatement",
      "start": 0,
      "end": 14,
      "left": {
        "type": "Identifier",
        "start": 5,
        "end": 6,
        "name": "x"
      },
      "right": {
        "type": "Identifier",
        "start": 10,
        "end": 12,
        "name": "xs"
      },
      "body": {
        "type": "EmptyStatement",
        "start": 13,
        "end": 14
      }
    }
  ]
}
`, ast)
}

func TestAsyncIteration11(t *testing.T) {
	ast, err := compile("for (x of xs);")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 14,
  "body": [
    {
      "type": "ForOfStatement",
      "start": 0,
      "end": 14,
      "await": false,
      "left": {
        "type": "Identifier",
        "start": 5,
        "end": 6,
        "name": "x"
      },
      "right": {
        "type": "Identifier",
        "start": 10,
        "end": 12,
        "name": "xs"
      },
      "body": {
        "type": "EmptyStatement",
        "start": 13,
        "end": 14
      }
    }
  ]
}
`, ast)
}

// FunctionDeclaration#await

func TestAsyncIteration12(t *testing.T) {
	ast, err := compile("async function* f() { await a; yield b; }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 41,
  "body": [
    {
      "type": "FunctionDeclaration",
      "start": 0,
      "end": 41,
      "id": {
        "type": "Identifier",
        "start": 16,
        "end": 17,
        "name": "f"
      },
      "generator": true,
      "async": true,
      "params": [],
      "body": {
        "type": "BlockStatement",
        "start": 20,
        "end": 41,
        "body": [
          {
            "type": "ExpressionStatement",
            "start": 22,
            "end": 30,
            "expression": {
              "type": "AwaitExpression",
              "start": 22,
              "end": 29,
              "argument": {
                "type": "Identifier",
                "start": 28,
                "end": 29,
                "name": "a"
              }
            }
          },
          {
            "type": "ExpressionStatement",
            "start": 31,
            "end": 39,
            "expression": {
              "type": "YieldExpression",
              "start": 31,
              "end": 38,
              "delegate": false,
              "argument": {
                "type": "Identifier",
                "start": 37,
                "end": 38,
                "name": "b"
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

// FunctionExpression#await
func TestAsyncIteration13(t *testing.T) {
	ast, err := compile("f = async function*() { await a; yield b; }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 43,
  "body": [
    {
      "type": "ExpressionStatement",
      "start": 0,
      "end": 43,
      "expression": {
        "type": "AssignmentExpression",
        "start": 0,
        "end": 43,
        "operator": "=",
        "left": {
          "type": "Identifier",
          "start": 0,
          "end": 1,
          "name": "f"
        },
        "right": {
          "type": "FunctionExpression",
          "start": 4,
          "end": 43,
          "id": null,
          "expression": false,
          "generator": true,
          "async": true,
          "params": [],
          "body": {
            "type": "BlockStatement",
            "start": 22,
            "end": 43,
            "body": [
              {
                "type": "ExpressionStatement",
                "start": 24,
                "end": 32,
                "expression": {
                  "type": "AwaitExpression",
                  "start": 24,
                  "end": 31,
                  "argument": {
                    "type": "Identifier",
                    "start": 30,
                    "end": 31,
                    "name": "a"
                  }
                }
              },
              {
                "type": "ExpressionStatement",
                "start": 33,
                "end": 41,
                "expression": {
                  "type": "YieldExpression",
                  "start": 33,
                  "end": 40,
                  "delegate": false,
                  "argument": {
                    "type": "Identifier",
                    "start": 39,
                    "end": 40,
                    "name": "b"
                  }
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

func TestAsyncIteration14(t *testing.T) {
	ast, err := compile("obj = { async* f() { await a; yield b; } }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 42,
  "body": [
    {
      "type": "ExpressionStatement",
      "start": 0,
      "end": 42,
      "expression": {
        "type": "AssignmentExpression",
        "start": 0,
        "end": 42,
        "operator": "=",
        "left": {
          "type": "Identifier",
          "start": 0,
          "end": 3,
          "name": "obj"
        },
        "right": {
          "type": "ObjectExpression",
          "start": 6,
          "end": 42,
          "properties": [
            {
              "type": "Property",
              "start": 8,
              "end": 40,
              "method": true,
              "shorthand": false,
              "computed": false,
              "key": {
                "type": "Identifier",
                "start": 15,
                "end": 16,
                "name": "f"
              },
              "kind": "init",
              "value": {
                "type": "FunctionExpression",
                "start": 16,
                "end": 40,
                "id": null,
                "expression": false,
                "generator": true,
                "async": true,
                "params": [],
                "body": {
                  "type": "BlockStatement",
                  "start": 19,
                  "end": 40,
                  "body": [
                    {
                      "type": "ExpressionStatement",
                      "start": 21,
                      "end": 29,
                      "expression": {
                        "type": "AwaitExpression",
                        "start": 21,
                        "end": 28,
                        "argument": {
                          "type": "Identifier",
                          "start": 27,
                          "end": 28,
                          "name": "a"
                        }
                      }
                    },
                    {
                      "type": "ExpressionStatement",
                      "start": 30,
                      "end": 38,
                      "expression": {
                        "type": "YieldExpression",
                        "start": 30,
                        "end": 37,
                        "delegate": false,
                        "argument": {
                          "type": "Identifier",
                          "start": 36,
                          "end": 37,
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
  ]
}
`, ast)
}

func TestAsyncIteration15(t *testing.T) {
	ast, err := compile("class A { async* f() { await a; yield b; } }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 44,
  "body": [
    {
      "type": "ClassDeclaration",
      "start": 0,
      "end": 44,
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
        "end": 44,
        "body": [
          {
            "type": "MethodDefinition",
            "start": 10,
            "end": 42,
            "static": false,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 17,
              "end": 18,
              "name": "f"
            },
            "kind": "method",
            "value": {
              "type": "FunctionExpression",
              "start": 18,
              "end": 42,
              "id": null,
              "expression": false,
              "generator": true,
              "async": true,
              "params": [],
              "body": {
                "type": "BlockStatement",
                "start": 21,
                "end": 42,
                "body": [
                  {
                    "type": "ExpressionStatement",
                    "start": 23,
                    "end": 31,
                    "expression": {
                      "type": "AwaitExpression",
                      "start": 23,
                      "end": 30,
                      "argument": {
                        "type": "Identifier",
                        "start": 29,
                        "end": 30,
                        "name": "a"
                      }
                    }
                  },
                  {
                    "type": "ExpressionStatement",
                    "start": 32,
                    "end": 40,
                    "expression": {
                      "type": "YieldExpression",
                      "start": 32,
                      "end": 39,
                      "delegate": false,
                      "argument": {
                        "type": "Identifier",
                        "start": 38,
                        "end": 39,
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
  ]
}
`, ast)
}

func TestAsyncIteration16(t *testing.T) {
	ast, err := compile("class A { static async* f() { await a; yield b; } }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 51,
  "body": [
    {
      "type": "ClassDeclaration",
      "start": 0,
      "end": 51,
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
        "end": 51,
        "body": [
          {
            "type": "MethodDefinition",
            "start": 10,
            "end": 49,
            "static": true,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 24,
              "end": 25,
              "name": "f"
            },
            "kind": "method",
            "value": {
              "type": "FunctionExpression",
              "start": 25,
              "end": 49,
              "id": null,
              "expression": false,
              "generator": true,
              "async": true,
              "params": [],
              "body": {
                "type": "BlockStatement",
                "start": 28,
                "end": 49,
                "body": [
                  {
                    "type": "ExpressionStatement",
                    "start": 30,
                    "end": 38,
                    "expression": {
                      "type": "AwaitExpression",
                      "start": 30,
                      "end": 37,
                      "argument": {
                        "type": "Identifier",
                        "start": 36,
                        "end": 37,
                        "name": "a"
                      }
                    }
                  },
                  {
                    "type": "ExpressionStatement",
                    "start": 39,
                    "end": 47,
                    "expression": {
                      "type": "YieldExpression",
                      "start": 39,
                      "end": 46,
                      "delegate": false,
                      "argument": {
                        "type": "Identifier",
                        "start": 45,
                        "end": 46,
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
  ]
}
`, ast)
}

func TestAsyncIteration17(t *testing.T) {
	ast, err := compile("async function* x() {}")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 22,
  "body": [
    {
      "type": "FunctionDeclaration",
      "start": 0,
      "end": 22,
      "id": {
        "type": "Identifier",
        "start": 16,
        "end": 17,
        "name": "x"
      },
      "generator": true,
      "async": true,
      "params": [],
      "body": {
        "type": "BlockStatement",
        "start": 20,
        "end": 22,
        "body": []
      }
    }
  ]
}
`, ast)
}

func TestAsyncIteration18(t *testing.T) {
	ast, err := compile("ref = async function*() {}")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
        "type": "AssignmentExpression",
        "start": 0,
        "end": 26,
        "operator": "=",
        "left": {
          "type": "Identifier",
          "start": 0,
          "end": 3,
          "name": "ref"
        },
        "right": {
          "type": "FunctionExpression",
          "start": 6,
          "end": 26,
          "id": null,
          "expression": false,
          "generator": true,
          "async": true,
          "params": [],
          "body": {
            "type": "BlockStatement",
            "start": 24,
            "end": 26,
            "body": []
          }
        }
      }
    }
  ]
}
`, ast)
}

func TestAsyncIteration19(t *testing.T) {
	ast, err := compile("(async function*() {})")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
        "type": "FunctionExpression",
        "start": 1,
        "end": 21,
        "id": null,
        "expression": false,
        "generator": true,
        "async": true,
        "params": [],
        "body": {
          "type": "BlockStatement",
          "start": 19,
          "end": 21,
          "body": []
        }
      }
    }
  ]
}
`, ast)
}

func TestAsyncIteration20(t *testing.T) {
	ast, err := compile("var gen = { async *method() {} }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 32,
  "body": [
    {
      "type": "VariableDeclaration",
      "start": 0,
      "end": 32,
      "declarations": [
        {
          "type": "VariableDeclarator",
          "start": 4,
          "end": 32,
          "id": {
            "type": "Identifier",
            "start": 4,
            "end": 7,
            "name": "gen"
          },
          "init": {
            "type": "ObjectExpression",
            "start": 10,
            "end": 32,
            "properties": [
              {
                "type": "Property",
                "start": 12,
                "end": 30,
                "method": true,
                "shorthand": false,
                "computed": false,
                "key": {
                  "type": "Identifier",
                  "start": 19,
                  "end": 25,
                  "name": "method"
                },
                "kind": "init",
                "value": {
                  "type": "FunctionExpression",
                  "start": 25,
                  "end": 30,
                  "id": null,
                  "expression": false,
                  "generator": true,
                  "async": true,
                  "params": [],
                  "body": {
                    "type": "BlockStatement",
                    "start": 28,
                    "end": 30,
                    "body": []
                  }
                }
              }
            ]
          }
        }
      ],
      "kind": "var"
    }
  ]
}
`, ast)
}

func TestAsyncIteration21(t *testing.T) {
	ast, err := compile("export default async function*() {}")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
        "generator": true,
        "async": true,
        "params": [],
        "body": {
          "type": "BlockStatement",
          "start": 33,
          "end": 35,
          "body": []
        }
      }
    }
  ]
}
`, ast)
}

func TestAsyncIteration22(t *testing.T) {
	ast, err := compile("var C = class { async *method() {} }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 36,
  "body": [
    {
      "type": "VariableDeclaration",
      "start": 0,
      "end": 36,
      "declarations": [
        {
          "type": "VariableDeclarator",
          "start": 4,
          "end": 36,
          "id": {
            "type": "Identifier",
            "start": 4,
            "end": 5,
            "name": "C"
          },
          "init": {
            "type": "ClassExpression",
            "start": 8,
            "end": 36,
            "id": null,
            "superClass": null,
            "body": {
              "type": "ClassBody",
              "start": 14,
              "end": 36,
              "body": [
                {
                  "type": "MethodDefinition",
                  "start": 16,
                  "end": 34,
                  "static": false,
                  "computed": false,
                  "key": {
                    "type": "Identifier",
                    "start": 23,
                    "end": 29,
                    "name": "method"
                  },
                  "kind": "method",
                  "value": {
                    "type": "FunctionExpression",
                    "start": 29,
                    "end": 34,
                    "id": null,
                    "expression": false,
                    "generator": true,
                    "async": true,
                    "params": [],
                    "body": {
                      "type": "BlockStatement",
                      "start": 32,
                      "end": 34,
                      "body": []
                    }
                  }
                }
              ]
            }
          }
        }
      ],
      "kind": "var"
    }
  ]
}
`, ast)
}

func TestAsyncIteration23(t *testing.T) {
	ast, err := compile("({ async *method(){} })")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
              "start": 10,
              "end": 16,
              "name": "method"
            },
            "kind": "init",
            "value": {
              "type": "FunctionExpression",
              "start": 16,
              "end": 20,
              "id": null,
              "expression": false,
              "generator": true,
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

func TestAsyncIteration24(t *testing.T) {
	ast, err := compile("({async() { }})")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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

func TestAsyncIteration25(t *testing.T) {
	ast, err := compile("({async = 0} = {})")
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

func TestAsyncIteration26(t *testing.T) {
	ast, err := compile("({async, foo})")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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

func TestAsyncIteration27(t *testing.T) {
	ast, err := compile("({async})")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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

func TestAsyncIteration28(t *testing.T) {
	ast, err := compile("({async: true})")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
            "method": false,
            "shorthand": false,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 2,
              "end": 7,
              "name": "async"
            },
            "value": {
              "type": "Literal",
              "start": 9,
              "end": 13,
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

func TestAsyncIteration29(t *testing.T) {
	ast, err := compile("async () => { for await (async of []); }")
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
        "type": "ArrowFunctionExpression",
        "start": 0,
        "end": 40,
        "id": null,
        "expression": false,
        "generator": false,
        "async": true,
        "params": [],
        "body": {
          "type": "BlockStatement",
          "start": 12,
          "end": 40,
          "body": [
            {
              "type": "ForOfStatement",
              "start": 14,
              "end": 38,
              "await": true,
              "left": {
                "type": "Identifier",
                "start": 25,
                "end": 30,
                "name": "async"
              },
              "right": {
                "type": "ArrayExpression",
                "start": 34,
                "end": 36,
                "elements": []
              },
              "body": {
                "type": "EmptyStatement",
                "start": 37,
                "end": 38
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

func TestAsyncIteration30(t *testing.T) {
	ast, err := compile("let a = async a => {}")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 21,
  "body": [
    {
      "type": "VariableDeclaration",
      "start": 0,
      "end": 21,
      "declarations": [
        {
          "type": "VariableDeclarator",
          "start": 4,
          "end": 21,
          "id": {
            "type": "Identifier",
            "start": 4,
            "end": 5,
            "name": "a"
          },
          "init": {
            "type": "ArrowFunctionExpression",
            "start": 8,
            "end": 21,
            "id": null,
            "expression": false,
            "generator": false,
            "async": true,
            "params": [
              {
                "type": "Identifier",
                "start": 14,
                "end": 15,
                "name": "a"
              }
            ],
            "body": {
              "type": "BlockStatement",
              "start": 19,
              "end": 21,
              "body": []
            }
          }
        }
      ],
      "kind": "let"
    }
  ]
}
`, ast)
}
