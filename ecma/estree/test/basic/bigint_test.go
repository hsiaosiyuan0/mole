package estree_test

import (
	"testing"

	. "github.com/hsiaosiyuan0/mole/ecma/estree/test"
	. "github.com/hsiaosiyuan0/mole/util"
)

func TestBigint1(t *testing.T) {
	ast, err := Compile("let i = 0n")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 10,
  "body": [
    {
      "type": "VariableDeclaration",
      "start": 0,
      "end": 10,
      "declarations": [
        {
          "type": "VariableDeclarator",
          "start": 4,
          "end": 10,
          "id": {
            "type": "Identifier",
            "start": 4,
            "end": 5,
            "name": "i"
          },
          "init": {
            "type": "Literal",
            "start": 8,
            "end": 10,
            "value": 0,
            "raw": "0n"
          }
        }
      ],
      "kind": "let"
    }
  ]
}
`, ast)
}

func TestBigint2(t *testing.T) {
	ast, err := Compile("i = 0n")
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
        "type": "AssignmentExpression",
        "start": 0,
        "end": 6,
        "operator": "=",
        "left": {
          "type": "Identifier",
          "start": 0,
          "end": 1,
          "name": "i"
        },
        "right": {
          "type": "Literal",
          "start": 4,
          "end": 6,
          "value": 0,
          "raw": "0n"
        }
      }
    }
  ]
}
`, ast)
}

func TestBigint3(t *testing.T) {
	ast, err := Compile("((i = 0n) => {})")
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
        "start": 1,
        "end": 15,
        "id": null,
        "expression": false,
        "generator": false,
        "async": false,
        "params": [
          {
            "type": "AssignmentPattern",
            "start": 2,
            "end": 8,
            "left": {
              "type": "Identifier",
              "start": 2,
              "end": 3,
              "name": "i"
            },
            "right": {
              "type": "Literal",
              "start": 6,
              "end": 8,
              "value": 0,
              "raw": "0n"
            }
          }
        ],
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
`, ast)
}

func TestBigint4(t *testing.T) {
	ast, err := Compile("for (let i = 0n; i < 0n;++i) {}")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 31,
  "body": [
    {
      "type": "ForStatement",
      "start": 0,
      "end": 31,
      "init": {
        "type": "VariableDeclaration",
        "start": 5,
        "end": 15,
        "declarations": [
          {
            "type": "VariableDeclarator",
            "start": 9,
            "end": 15,
            "id": {
              "type": "Identifier",
              "start": 9,
              "end": 10,
              "name": "i"
            },
            "init": {
              "type": "Literal",
              "start": 13,
              "end": 15,
              "value": 0,
              "raw": "0n"
            }
          }
        ],
        "kind": "let"
      },
      "test": {
        "type": "BinaryExpression",
        "start": 17,
        "end": 23,
        "left": {
          "type": "Identifier",
          "start": 17,
          "end": 18,
          "name": "i"
        },
        "operator": "<",
        "right": {
          "type": "Literal",
          "start": 21,
          "end": 23,
          "value": 0,
          "raw": "0n"
        }
      },
      "update": {
        "type": "UpdateExpression",
        "start": 24,
        "end": 27,
        "operator": "++",
        "prefix": true,
        "argument": {
          "type": "Identifier",
          "start": 26,
          "end": 27,
          "name": "i"
        }
      },
      "body": {
        "type": "BlockStatement",
        "start": 29,
        "end": 31,
        "body": []
      }
    }
  ]
}
`, ast)
}

func TestBigint5(t *testing.T) {
	ast, err := Compile("i + 0n")
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
          "name": "i"
        },
        "operator": "+",
        "right": {
          "type": "Literal",
          "start": 4,
          "end": 6,
          "value": 0,
          "raw": "0n"
        }
      }
    }
  ]
}
`, ast)
}

func TestBigint6(t *testing.T) {
	ast, err := Compile("let i = 2n")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 10,
  "body": [
    {
      "type": "VariableDeclaration",
      "start": 0,
      "end": 10,
      "declarations": [
        {
          "type": "VariableDeclarator",
          "start": 4,
          "end": 10,
          "id": {
            "type": "Identifier",
            "start": 4,
            "end": 5,
            "name": "i"
          },
          "init": {
            "type": "Literal",
            "start": 8,
            "end": 10,
            "value": 2,
            "raw": "2n"
          }
        }
      ],
      "kind": "let"
    }
  ]
}
`, ast)
}

func TestBigint7(t *testing.T) {
	ast, err := Compile("i = 2n")
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
        "type": "AssignmentExpression",
        "start": 0,
        "end": 6,
        "operator": "=",
        "left": {
          "type": "Identifier",
          "start": 0,
          "end": 1,
          "name": "i"
        },
        "right": {
          "type": "Literal",
          "start": 4,
          "end": 6,
          "value": 2,
          "raw": "2n"
        }
      }
    }
  ]
}
`, ast)
}

func TestBigint8(t *testing.T) {
	ast, err := Compile("((i = 2n) => {})")
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
        "start": 1,
        "end": 15,
        "id": null,
        "expression": false,
        "generator": false,
        "async": false,
        "params": [
          {
            "type": "AssignmentPattern",
            "start": 2,
            "end": 8,
            "left": {
              "type": "Identifier",
              "start": 2,
              "end": 3,
              "name": "i"
            },
            "right": {
              "type": "Literal",
              "start": 6,
              "end": 8,
              "value": 2,
              "raw": "2n"
            }
          }
        ],
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
`, ast)
}

func TestBigint9(t *testing.T) {
	ast, err := Compile("for (let i = 0n; i < 2n;++i) {}")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 31,
  "body": [
    {
      "type": "ForStatement",
      "start": 0,
      "end": 31,
      "init": {
        "type": "VariableDeclaration",
        "start": 5,
        "end": 15,
        "declarations": [
          {
            "type": "VariableDeclarator",
            "start": 9,
            "end": 15,
            "id": {
              "type": "Identifier",
              "start": 9,
              "end": 10,
              "name": "i"
            },
            "init": {
              "type": "Literal",
              "start": 13,
              "end": 15,
              "value": 0,
              "raw": "0n"
            }
          }
        ],
        "kind": "let"
      },
      "test": {
        "type": "BinaryExpression",
        "start": 17,
        "end": 23,
        "left": {
          "type": "Identifier",
          "start": 17,
          "end": 18,
          "name": "i"
        },
        "operator": "<",
        "right": {
          "type": "Literal",
          "start": 21,
          "end": 23,
          "value": 2,
          "raw": "2n"
        }
      },
      "update": {
        "type": "UpdateExpression",
        "start": 24,
        "end": 27,
        "operator": "++",
        "prefix": true,
        "argument": {
          "type": "Identifier",
          "start": 26,
          "end": 27,
          "name": "i"
        }
      },
      "body": {
        "type": "BlockStatement",
        "start": 29,
        "end": 31,
        "body": []
      }
    }
  ]
}
`, ast)
}

func TestBigint10(t *testing.T) {
	ast, err := Compile("i + 2n")
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
          "name": "i"
        },
        "operator": "+",
        "right": {
          "type": "Literal",
          "start": 4,
          "end": 6,
          "value": 2,
          "raw": "2n"
        }
      }
    }
  ]
}
`, ast)
}

func TestBigint11(t *testing.T) {
	ast, err := Compile("let i = 0x2n")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 12,
  "body": [
    {
      "type": "VariableDeclaration",
      "start": 0,
      "end": 12,
      "declarations": [
        {
          "type": "VariableDeclarator",
          "start": 4,
          "end": 12,
          "id": {
            "type": "Identifier",
            "start": 4,
            "end": 5,
            "name": "i"
          },
          "init": {
            "type": "Literal",
            "start": 8,
            "end": 12,
            "value": 2,
            "raw": "0x2n"
          }
        }
      ],
      "kind": "let"
    }
  ]
}
`, ast)
}

func TestBigint12(t *testing.T) {
	ast, err := Compile("i = 0x2n")
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
        "operator": "=",
        "left": {
          "type": "Identifier",
          "start": 0,
          "end": 1,
          "name": "i"
        },
        "right": {
          "type": "Literal",
          "start": 4,
          "end": 8,
          "value": 2,
          "raw": "0x2n",
          "bigint": "2"
        }
      }
    }
  ]
}
`, ast)
}

func TestBigint13(t *testing.T) {
	ast, err := Compile("((i = 0x2n) => {})")
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
        "type": "ArrowFunctionExpression",
        "start": 1,
        "end": 17,
        "id": null,
        "expression": false,
        "generator": false,
        "async": false,
        "params": [
          {
            "type": "AssignmentPattern",
            "start": 2,
            "end": 10,
            "left": {
              "type": "Identifier",
              "start": 2,
              "end": 3,
              "name": "i"
            },
            "right": {
              "type": "Literal",
              "start": 6,
              "end": 10,
              "value": 2,
              "raw": "0x2n",
              "bigint": "2"
            }
          }
        ],
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
`, ast)
}

func TestBigint14(t *testing.T) {
	ast, err := Compile("for (let i = 0n; i < 0x2n;++i) {}")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 33,
  "body": [
    {
      "type": "ForStatement",
      "start": 0,
      "end": 33,
      "init": {
        "type": "VariableDeclaration",
        "start": 5,
        "end": 15,
        "declarations": [
          {
            "type": "VariableDeclarator",
            "start": 9,
            "end": 15,
            "id": {
              "type": "Identifier",
              "start": 9,
              "end": 10,
              "name": "i"
            },
            "init": {
              "type": "Literal",
              "start": 13,
              "end": 15,
              "value": 0,
              "raw": "0n",
              "bigint": "0"
            }
          }
        ],
        "kind": "let"
      },
      "test": {
        "type": "BinaryExpression",
        "start": 17,
        "end": 25,
        "left": {
          "type": "Identifier",
          "start": 17,
          "end": 18,
          "name": "i"
        },
        "operator": "<",
        "right": {
          "type": "Literal",
          "start": 21,
          "end": 25,
          "value": 2,
          "raw": "0x2n",
          "bigint": "2"
        }
      },
      "update": {
        "type": "UpdateExpression",
        "start": 26,
        "end": 29,
        "operator": "++",
        "prefix": true,
        "argument": {
          "type": "Identifier",
          "start": 28,
          "end": 29,
          "name": "i"
        }
      },
      "body": {
        "type": "BlockStatement",
        "start": 31,
        "end": 33,
        "body": []
      }
    }
  ]
}
`, ast)
}

func TestBigint15(t *testing.T) {
	ast, err := Compile("i + 0x2n")
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
          "type": "Identifier",
          "start": 0,
          "end": 1,
          "name": "i"
        },
        "operator": "+",
        "right": {
          "type": "Literal",
          "start": 4,
          "end": 8,
          "value": 2,
          "raw": "0x2n",
          "bigint": "2"
        }
      }
    }
  ]
}
`, ast)
}

func TestBigint16(t *testing.T) {
	ast, err := Compile("let i = 0o2n")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 12,
  "body": [
    {
      "type": "VariableDeclaration",
      "start": 0,
      "end": 12,
      "declarations": [
        {
          "type": "VariableDeclarator",
          "start": 4,
          "end": 12,
          "id": {
            "type": "Identifier",
            "start": 4,
            "end": 5,
            "name": "i"
          },
          "init": {
            "type": "Literal",
            "start": 8,
            "end": 12,
            "value": 2,
            "raw": "0o2n"
          }
        }
      ],
      "kind": "let"
    }
  ]
}
`, ast)
}

func TestBigint17(t *testing.T) {
	ast, err := Compile("i = 0o2n")
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
        "operator": "=",
        "left": {
          "type": "Identifier",
          "start": 0,
          "end": 1,
          "name": "i"
        },
        "right": {
          "type": "Literal",
          "start": 4,
          "end": 8,
          "value": 2,
          "raw": "0o2n"
        }
      }
    }
  ]
}
`, ast)
}

func TestBigint18(t *testing.T) {
	ast, err := Compile("((i = 0o2n) => {})")
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
        "type": "ArrowFunctionExpression",
        "start": 1,
        "end": 17,
        "id": null,
        "expression": false,
        "generator": false,
        "async": false,
        "params": [
          {
            "type": "AssignmentPattern",
            "start": 2,
            "end": 10,
            "left": {
              "type": "Identifier",
              "start": 2,
              "end": 3,
              "name": "i"
            },
            "right": {
              "type": "Literal",
              "start": 6,
              "end": 10,
              "value": 2,
              "raw": "0o2n"
            }
          }
        ],
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
`, ast)
}

func TestBigint19(t *testing.T) {
	ast, err := Compile("for (let i = 0n; i < 0o2n;++i) {}")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 33,
  "body": [
    {
      "type": "ForStatement",
      "start": 0,
      "end": 33,
      "init": {
        "type": "VariableDeclaration",
        "start": 5,
        "end": 15,
        "declarations": [
          {
            "type": "VariableDeclarator",
            "start": 9,
            "end": 15,
            "id": {
              "type": "Identifier",
              "start": 9,
              "end": 10,
              "name": "i"
            },
            "init": {
              "type": "Literal",
              "start": 13,
              "end": 15,
              "value": 0,
              "raw": "0n"
            }
          }
        ],
        "kind": "let"
      },
      "test": {
        "type": "BinaryExpression",
        "start": 17,
        "end": 25,
        "left": {
          "type": "Identifier",
          "start": 17,
          "end": 18,
          "name": "i"
        },
        "operator": "<",
        "right": {
          "type": "Literal",
          "start": 21,
          "end": 25,
          "value": 2,
          "raw": "0o2n"
        }
      },
      "update": {
        "type": "UpdateExpression",
        "start": 26,
        "end": 29,
        "operator": "++",
        "prefix": true,
        "argument": {
          "type": "Identifier",
          "start": 28,
          "end": 29,
          "name": "i"
        }
      },
      "body": {
        "type": "BlockStatement",
        "start": 31,
        "end": 33,
        "body": []
      }
    }
  ]
}
`, ast)
}

func TestBigint20(t *testing.T) {
	ast, err := Compile("i + 0o2n")
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
          "type": "Identifier",
          "start": 0,
          "end": 1,
          "name": "i"
        },
        "operator": "+",
        "right": {
          "type": "Literal",
          "start": 4,
          "end": 8,
          "value": 2,
          "raw": "0o2n"
        }
      }
    }
  ]
}
`, ast)
}

func TestBigint21(t *testing.T) {
	ast, err := Compile("let i = 0b10n")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 13,
  "body": [
    {
      "type": "VariableDeclaration",
      "start": 0,
      "end": 13,
      "declarations": [
        {
          "type": "VariableDeclarator",
          "start": 4,
          "end": 13,
          "id": {
            "type": "Identifier",
            "start": 4,
            "end": 5,
            "name": "i"
          },
          "init": {
            "type": "Literal",
            "start": 8,
            "end": 13,
            "value": 2,
            "raw": "0b10n"
          }
        }
      ],
      "kind": "let"
    }
  ]
}
`, ast)
}

func TestBigint22(t *testing.T) {
	ast, err := Compile("i = 0b10n")
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
        "type": "AssignmentExpression",
        "start": 0,
        "end": 9,
        "operator": "=",
        "left": {
          "type": "Identifier",
          "start": 0,
          "end": 1,
          "name": "i"
        },
        "right": {
          "type": "Literal",
          "start": 4,
          "end": 9,
          "value": 2,
          "raw": "0b10n"
        }
      }
    }
  ]
}
`, ast)
}

func TestBigint23(t *testing.T) {
	ast, err := Compile("((i = 0b10n) => {})")
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
        "type": "ArrowFunctionExpression",
        "start": 1,
        "end": 18,
        "id": null,
        "expression": false,
        "generator": false,
        "async": false,
        "params": [
          {
            "type": "AssignmentPattern",
            "start": 2,
            "end": 11,
            "left": {
              "type": "Identifier",
              "start": 2,
              "end": 3,
              "name": "i"
            },
            "right": {
              "type": "Literal",
              "start": 6,
              "end": 11,
              "value": 2,
              "raw": "0b10n"
            }
          }
        ],
        "body": {
          "type": "BlockStatement",
          "start": 16,
          "end": 18,
          "body": []
        }
      }
    }
  ]
}
`, ast)
}

func TestBigint24(t *testing.T) {
	ast, err := Compile("for (let i = 0n; i < 0b10n;++i) {}")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 34,
  "body": [
    {
      "type": "ForStatement",
      "start": 0,
      "end": 34,
      "init": {
        "type": "VariableDeclaration",
        "start": 5,
        "end": 15,
        "declarations": [
          {
            "type": "VariableDeclarator",
            "start": 9,
            "end": 15,
            "id": {
              "type": "Identifier",
              "start": 9,
              "end": 10,
              "name": "i"
            },
            "init": {
              "type": "Literal",
              "start": 13,
              "end": 15,
              "value": 0,
              "raw": "0n"
            }
          }
        ],
        "kind": "let"
      },
      "test": {
        "type": "BinaryExpression",
        "start": 17,
        "end": 26,
        "left": {
          "type": "Identifier",
          "start": 17,
          "end": 18,
          "name": "i"
        },
        "operator": "<",
        "right": {
          "type": "Literal",
          "start": 21,
          "end": 26,
          "value": 2,
          "raw": "0b10n"
        }
      },
      "update": {
        "type": "UpdateExpression",
        "start": 27,
        "end": 30,
        "operator": "++",
        "prefix": true,
        "argument": {
          "type": "Identifier",
          "start": 29,
          "end": 30,
          "name": "i"
        }
      },
      "body": {
        "type": "BlockStatement",
        "start": 32,
        "end": 34,
        "body": []
      }
    }
  ]
}
`, ast)
}

func TestBigint25(t *testing.T) {
	ast, err := Compile("i + 0b10n")
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
          "type": "Identifier",
          "start": 0,
          "end": 1,
          "name": "i"
        },
        "operator": "+",
        "right": {
          "type": "Literal",
          "start": 4,
          "end": 9,
          "value": 2,
          "raw": "0b10n"
        }
      }
    }
  ]
}
`, ast)
}

func TestBigint26(t *testing.T) {
	ast, err := Compile("let i = -0xbf2ed51ff75d380fd3be813ec6185780n")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 44,
  "body": [
    {
      "type": "VariableDeclaration",
      "start": 0,
      "end": 44,
      "declarations": [
        {
          "type": "VariableDeclarator",
          "start": 4,
          "end": 44,
          "id": {
            "type": "Identifier",
            "start": 4,
            "end": 5,
            "name": "i"
          },
          "init": {
            "type": "UnaryExpression",
            "start": 8,
            "end": 44,
            "operator": "-",
            "prefix": true,
            "argument": {
              "type": "Literal",
              "start": 9,
              "end": 44,
              "value": 0,
              "raw": "0xbf2ed51ff75d380fd3be813ec6185780n",
              "bigint": "254125715536285641815112686497309415296"
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

func TestBigint27(t *testing.T) {
	ast, err := Compile("i = -0xbf2ed51ff75d380fd3be813ec6185780n")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
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
        "type": "AssignmentExpression",
        "start": 0,
        "end": 40,
        "operator": "=",
        "left": {
          "type": "Identifier",
          "start": 0,
          "end": 1,
          "name": "i"
        },
        "right": {
          "type": "UnaryExpression",
          "start": 4,
          "end": 40,
          "operator": "-",
          "prefix": true,
          "argument": {
            "type": "Literal",
            "start": 5,
            "end": 40,
            "value": 0,
            "raw": "0xbf2ed51ff75d380fd3be813ec6185780n",
            "bigint": "254125715536285641815112686497309415296"
          }
        }
      }
    }
  ]
}
`, ast)
}

func TestBigint28(t *testing.T) {
	ast, err := Compile("((i = -0xbf2ed51ff75d380fd3be813ec6185780n) => {})")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 50,
  "body": [
    {
      "type": "ExpressionStatement",
      "start": 0,
      "end": 50,
      "expression": {
        "type": "ArrowFunctionExpression",
        "start": 1,
        "end": 49,
        "id": null,
        "expression": false,
        "generator": false,
        "async": false,
        "params": [
          {
            "type": "AssignmentPattern",
            "start": 2,
            "end": 42,
            "left": {
              "type": "Identifier",
              "start": 2,
              "end": 3,
              "name": "i"
            },
            "right": {
              "type": "UnaryExpression",
              "start": 6,
              "end": 42,
              "operator": "-",
              "prefix": true,
              "argument": {
                "type": "Literal",
                "start": 7,
                "end": 42,
                "value": 0,
                "raw": "0xbf2ed51ff75d380fd3be813ec6185780n",
                "bigint": "254125715536285641815112686497309415296"
              }
            }
          }
        ],
        "body": {
          "type": "BlockStatement",
          "start": 47,
          "end": 49,
          "body": []
        }
      }
    }
  ]
}
`, ast)
}

func TestBigint29(t *testing.T) {
	ast, err := Compile("for (let i = 0n; i < -0xbf2ed51ff75d380fd3be813ec6185780n;++i) {}")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 65,
  "body": [
    {
      "type": "ForStatement",
      "start": 0,
      "end": 65,
      "init": {
        "type": "VariableDeclaration",
        "start": 5,
        "end": 15,
        "declarations": [
          {
            "type": "VariableDeclarator",
            "start": 9,
            "end": 15,
            "id": {
              "type": "Identifier",
              "start": 9,
              "end": 10,
              "name": "i"
            },
            "init": {
              "type": "Literal",
              "start": 13,
              "end": 15,
              "value": 0,
              "raw": "0n",
              "bigint": "0"
            }
          }
        ],
        "kind": "let"
      },
      "test": {
        "type": "BinaryExpression",
        "start": 17,
        "end": 57,
        "left": {
          "type": "Identifier",
          "start": 17,
          "end": 18,
          "name": "i"
        },
        "operator": "<",
        "right": {
          "type": "UnaryExpression",
          "start": 21,
          "end": 57,
          "operator": "-",
          "prefix": true,
          "argument": {
            "type": "Literal",
            "start": 22,
            "end": 57,
            "value": 0,
            "raw": "0xbf2ed51ff75d380fd3be813ec6185780n",
            "bigint": "254125715536285641815112686497309415296"
          }
        }
      },
      "update": {
        "type": "UpdateExpression",
        "start": 58,
        "end": 61,
        "operator": "++",
        "prefix": true,
        "argument": {
          "type": "Identifier",
          "start": 60,
          "end": 61,
          "name": "i"
        }
      },
      "body": {
        "type": "BlockStatement",
        "start": 63,
        "end": 65,
        "body": []
      }
    }
  ]
}
`, ast)
}

func TestBigint30(t *testing.T) {
	ast, err := Compile("i + -0xbf2ed51ff75d380fd3be813ec6185780n")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
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
        "type": "BinaryExpression",
        "start": 0,
        "end": 40,
        "left": {
          "type": "Identifier",
          "start": 0,
          "end": 1,
          "name": "i"
        },
        "operator": "+",
        "right": {
          "type": "UnaryExpression",
          "start": 4,
          "end": 40,
          "operator": "-",
          "prefix": true,
          "argument": {
            "type": "Literal",
            "start": 5,
            "end": 40,
            "value": 0,
            "raw": "0xbf2ed51ff75d380fd3be813ec6185780n",
            "bigint": "254125715536285641815112686497309415296"
          }
        }
      }
    }
  ]
}
`, ast)
}
