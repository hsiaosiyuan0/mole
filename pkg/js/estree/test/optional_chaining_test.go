package estree_test

import (
	"testing"

	"github.com/hsiaosiyuan0/mole/pkg/assert"
)

func TestOptionalChain1(t *testing.T) {
	ast, err := compile("obj?.foo")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
        "type": "ChainExpression",
        "start": 0,
        "end": 8,
        "expression": {
          "type": "MemberExpression",
          "start": 0,
          "end": 8,
          "object": {
            "type": "Identifier",
            "start": 0,
            "end": 3,
            "name": "obj"
          },
          "property": {
            "type": "Identifier",
            "start": 5,
            "end": 8,
            "name": "foo"
          },
          "computed": false,
          "optional": true
        }
      }
    }
  ]
}
`, ast)
}

func TestOptionalChain2(t *testing.T) {
	ast, err := compile("obj?.[foo]")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
        "type": "ChainExpression",
        "start": 0,
        "end": 10,
        "expression": {
          "type": "MemberExpression",
          "start": 0,
          "end": 10,
          "object": {
            "type": "Identifier",
            "start": 0,
            "end": 3,
            "name": "obj"
          },
          "property": {
            "type": "Identifier",
            "start": 6,
            "end": 9,
            "name": "foo"
          },
          "computed": true,
          "optional": true
        }
      }
    }
  ]
}
`, ast)
}

func TestOptionalChain3(t *testing.T) {
	ast, err := compile("obj?.()")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 7,
  "body": [
    {
      "type": "ExpressionStatement",
      "start": 0,
      "end": 7,
      "expression": {
        "type": "ChainExpression",
        "start": 0,
        "end": 7,
        "expression": {
          "type": "CallExpression",
          "start": 0,
          "end": 7,
          "callee": {
            "type": "Identifier",
            "start": 0,
            "end": 3,
            "name": "obj"
          },
          "arguments": [],
          "optional": true
        }
      }
    }
  ]
}
`, ast)
}

func TestOptionalChain4(t *testing.T) {
	ast, err := compile("obj ?. foo")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
        "type": "ChainExpression",
        "start": 0,
        "end": 10,
        "expression": {
          "type": "MemberExpression",
          "start": 0,
          "end": 10,
          "object": {
            "type": "Identifier",
            "start": 0,
            "end": 3,
            "name": "obj"
          },
          "property": {
            "type": "Identifier",
            "start": 7,
            "end": 10,
            "name": "foo"
          },
          "computed": false,
          "optional": true
        }
      }
    }
  ]
}
`, ast)
}

func TestOptionalChain5(t *testing.T) {
	ast, err := compile("obj ?. [foo]")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
        "type": "ChainExpression",
        "start": 0,
        "end": 12,
        "expression": {
          "type": "MemberExpression",
          "start": 0,
          "end": 12,
          "object": {
            "type": "Identifier",
            "start": 0,
            "end": 3,
            "name": "obj"
          },
          "property": {
            "type": "Identifier",
            "start": 8,
            "end": 11,
            "name": "foo"
          },
          "computed": true,
          "optional": true
        }
      }
    }
  ]
}
`, ast)
}

func TestOptionalChain6(t *testing.T) {
	ast, err := compile("obj ?. ()")
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
        "type": "ChainExpression",
        "start": 0,
        "end": 9,
        "expression": {
          "type": "CallExpression",
          "start": 0,
          "end": 9,
          "callee": {
            "type": "Identifier",
            "start": 0,
            "end": 3,
            "name": "obj"
          },
          "arguments": [],
          "optional": true
        }
      }
    }
  ]
}
`, ast)
}

func TestOptionalChain7(t *testing.T) {
	ast, err := compile("obj?.0:.1")
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
        "type": "ConditionalExpression",
        "start": 0,
        "end": 9,
        "test": {
          "type": "Identifier",
          "start": 0,
          "end": 3,
          "name": "obj"
        },
        "consequent": {
          "type": "Literal",
          "start": 4,
          "end": 6,
          "value": 0,
          "raw": ".0"
        },
        "alternate": {
          "type": "Literal",
          "start": 7,
          "end": 9,
          "value": 0.1,
          "raw": ".1"
        }
      }
    }
  ]
}
`, ast)
}

func TestOptionalChain8(t *testing.T) {
	ast, err := compile("obj?.aaa?.bbb")
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
        "type": "ChainExpression",
        "start": 0,
        "end": 13,
        "expression": {
          "type": "MemberExpression",
          "start": 0,
          "end": 13,
          "object": {
            "type": "MemberExpression",
            "start": 0,
            "end": 8,
            "object": {
              "type": "Identifier",
              "start": 0,
              "end": 3,
              "name": "obj"
            },
            "property": {
              "type": "Identifier",
              "start": 5,
              "end": 8,
              "name": "aaa"
            },
            "computed": false,
            "optional": true
          },
          "property": {
            "type": "Identifier",
            "start": 10,
            "end": 13,
            "name": "bbb"
          },
          "computed": false,
          "optional": true
        }
      }
    }
  ]
}
`, ast)
}

func TestOptionalChain9(t *testing.T) {
	ast, err := compile("obj?.aaa.bbb")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
        "type": "ChainExpression",
        "start": 0,
        "end": 12,
        "expression": {
          "type": "MemberExpression",
          "start": 0,
          "end": 12,
          "object": {
            "type": "MemberExpression",
            "start": 0,
            "end": 8,
            "object": {
              "type": "Identifier",
              "start": 0,
              "end": 3,
              "name": "obj"
            },
            "property": {
              "type": "Identifier",
              "start": 5,
              "end": 8,
              "name": "aaa"
            },
            "computed": false,
            "optional": true
          },
          "property": {
            "type": "Identifier",
            "start": 9,
            "end": 12,
            "name": "bbb"
          },
          "computed": false,
          "optional": false
        }
      }
    }
  ]
}
`, ast)
}

func TestOptionalChain10(t *testing.T) {
	ast, err := compile("(obj?.aaa)?.bbb")
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
        "type": "ChainExpression",
        "start": 0,
        "end": 15,
        "expression": {
          "type": "MemberExpression",
          "start": 0,
          "end": 15,
          "object": {
            "type": "ChainExpression",
            "start": 1,
            "end": 9,
            "expression": {
              "type": "MemberExpression",
              "start": 1,
              "end": 9,
              "object": {
                "type": "Identifier",
                "start": 1,
                "end": 4,
                "name": "obj"
              },
              "property": {
                "type": "Identifier",
                "start": 6,
                "end": 9,
                "name": "aaa"
              },
              "computed": false,
              "optional": true
            }
          },
          "property": {
            "type": "Identifier",
            "start": 12,
            "end": 15,
            "name": "bbb"
          },
          "computed": false,
          "optional": true
        }
      }
    }
  ]
}
`, ast)
}

func TestOptionalChain11(t *testing.T) {
	ast, err := compile("(obj?.aaa).bbb")
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
        "type": "MemberExpression",
        "start": 0,
        "end": 14,
        "object": {
          "type": "ChainExpression",
          "start": 1,
          "end": 9,
          "expression": {
            "type": "MemberExpression",
            "start": 1,
            "end": 9,
            "object": {
              "type": "Identifier",
              "start": 1,
              "end": 4,
              "name": "obj"
            },
            "property": {
              "type": "Identifier",
              "start": 6,
              "end": 9,
              "name": "aaa"
            },
            "computed": false,
            "optional": true
          }
        },
        "property": {
          "type": "Identifier",
          "start": 11,
          "end": 14,
          "name": "bbb"
        },
        "computed": false,
        "optional": false
      }
    }
  ]
}
`, ast)
}

func TestOptionalChain12(t *testing.T) {
	ast, err := compile("(obj?.aaa.bbb).ccc")
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
        "type": "MemberExpression",
        "start": 0,
        "end": 18,
        "object": {
          "type": "ChainExpression",
          "start": 1,
          "end": 13,
          "expression": {
            "type": "MemberExpression",
            "start": 1,
            "end": 13,
            "object": {
              "type": "MemberExpression",
              "start": 1,
              "end": 9,
              "object": {
                "type": "Identifier",
                "start": 1,
                "end": 4,
                "name": "obj"
              },
              "property": {
                "type": "Identifier",
                "start": 6,
                "end": 9,
                "name": "aaa"
              },
              "computed": false,
              "optional": true
            },
            "property": {
              "type": "Identifier",
              "start": 10,
              "end": 13,
              "name": "bbb"
            },
            "computed": false,
            "optional": false
          }
        },
        "property": {
          "type": "Identifier",
          "start": 15,
          "end": 18,
          "name": "ccc"
        },
        "computed": false,
        "optional": false
      }
    }
  ]
}
`, ast)
}

func TestOptionalChain13(t *testing.T) {
	ast, err := compile("func?.()?.bbb")
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
        "type": "ChainExpression",
        "start": 0,
        "end": 13,
        "expression": {
          "type": "MemberExpression",
          "start": 0,
          "end": 13,
          "object": {
            "type": "CallExpression",
            "start": 0,
            "end": 8,
            "callee": {
              "type": "Identifier",
              "start": 0,
              "end": 4,
              "name": "func"
            },
            "arguments": [],
            "optional": true
          },
          "property": {
            "type": "Identifier",
            "start": 10,
            "end": 13,
            "name": "bbb"
          },
          "computed": false,
          "optional": true
        }
      }
    }
  ]
}
`, ast)
}

func TestOptionalChain14(t *testing.T) {
	ast, err := compile("func?.().bbb")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
        "type": "ChainExpression",
        "start": 0,
        "end": 12,
        "expression": {
          "type": "MemberExpression",
          "start": 0,
          "end": 12,
          "object": {
            "type": "CallExpression",
            "start": 0,
            "end": 8,
            "callee": {
              "type": "Identifier",
              "start": 0,
              "end": 4,
              "name": "func"
            },
            "arguments": [],
            "optional": true
          },
          "property": {
            "type": "Identifier",
            "start": 9,
            "end": 12,
            "name": "bbb"
          },
          "computed": false,
          "optional": false
        }
      }
    }
  ]
}
`, ast)
}

func TestOptionalChain15(t *testing.T) {
	ast, err := compile("(func?.())?.bbb")
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
        "type": "ChainExpression",
        "start": 0,
        "end": 15,
        "expression": {
          "type": "MemberExpression",
          "start": 0,
          "end": 15,
          "object": {
            "type": "ChainExpression",
            "start": 1,
            "end": 9,
            "expression": {
              "type": "CallExpression",
              "start": 1,
              "end": 9,
              "callee": {
                "type": "Identifier",
                "start": 1,
                "end": 5,
                "name": "func"
              },
              "arguments": [],
              "optional": true
            }
          },
          "property": {
            "type": "Identifier",
            "start": 12,
            "end": 15,
            "name": "bbb"
          },
          "computed": false,
          "optional": true
        }
      }
    }
  ]
}
`, ast)
}

func TestOptionalChain16(t *testing.T) {
	ast, err := compile("(func?.()).bbb")
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
        "type": "MemberExpression",
        "start": 0,
        "end": 14,
        "object": {
          "type": "ChainExpression",
          "start": 1,
          "end": 9,
          "expression": {
            "type": "CallExpression",
            "start": 1,
            "end": 9,
            "callee": {
              "type": "Identifier",
              "start": 1,
              "end": 5,
              "name": "func"
            },
            "arguments": [],
            "optional": true
          }
        },
        "property": {
          "type": "Identifier",
          "start": 11,
          "end": 14,
          "name": "bbb"
        },
        "computed": false,
        "optional": false
      }
    }
  ]
}
`, ast)
}

func TestOptionalChain17(t *testing.T) {
	ast, err := compile("obj?.aaa?.()")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
        "type": "ChainExpression",
        "start": 0,
        "end": 12,
        "expression": {
          "type": "CallExpression",
          "start": 0,
          "end": 12,
          "callee": {
            "type": "MemberExpression",
            "start": 0,
            "end": 8,
            "object": {
              "type": "Identifier",
              "start": 0,
              "end": 3,
              "name": "obj"
            },
            "property": {
              "type": "Identifier",
              "start": 5,
              "end": 8,
              "name": "aaa"
            },
            "computed": false,
            "optional": true
          },
          "arguments": [],
          "optional": true
        }
      }
    }
  ]
}
`, ast)
}

func TestOptionalChain18(t *testing.T) {
	ast, err := compile("obj?.aaa()")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
        "type": "ChainExpression",
        "start": 0,
        "end": 10,
        "expression": {
          "type": "CallExpression",
          "start": 0,
          "end": 10,
          "callee": {
            "type": "MemberExpression",
            "start": 0,
            "end": 8,
            "object": {
              "type": "Identifier",
              "start": 0,
              "end": 3,
              "name": "obj"
            },
            "property": {
              "type": "Identifier",
              "start": 5,
              "end": 8,
              "name": "aaa"
            },
            "computed": false,
            "optional": true
          },
          "arguments": [],
          "optional": false
        }
      }
    }
  ]
}
`, ast)
}

func TestOptionalChain19(t *testing.T) {
	ast, err := compile("(obj?.aaa)?.()")
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
        "type": "ChainExpression",
        "start": 0,
        "end": 14,
        "expression": {
          "type": "CallExpression",
          "start": 0,
          "end": 14,
          "callee": {
            "type": "ChainExpression",
            "start": 1,
            "end": 9,
            "expression": {
              "type": "MemberExpression",
              "start": 1,
              "end": 9,
              "object": {
                "type": "Identifier",
                "start": 1,
                "end": 4,
                "name": "obj"
              },
              "property": {
                "type": "Identifier",
                "start": 6,
                "end": 9,
                "name": "aaa"
              },
              "computed": false,
              "optional": true
            }
          },
          "arguments": [],
          "optional": true
        }
      }
    }
  ]
}
`, ast)
}

func TestOptionalChain20(t *testing.T) {
	ast, err := compile("(obj?.aaa)()")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
        "type": "CallExpression",
        "start": 0,
        "end": 12,
        "callee": {
          "type": "ChainExpression",
          "start": 1,
          "end": 9,
          "expression": {
            "type": "MemberExpression",
            "start": 1,
            "end": 9,
            "object": {
              "type": "Identifier",
              "start": 1,
              "end": 4,
              "name": "obj"
            },
            "property": {
              "type": "Identifier",
              "start": 6,
              "end": 9,
              "name": "aaa"
            },
            "computed": false,
            "optional": true
          }
        },
        "arguments": [],
        "optional": false
      }
    }
  ]
}
`, ast)
}

func TestOptionalChain21(t *testing.T) {
	ast, err := compile("delete obj?.foo")
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
        "type": "UnaryExpression",
        "start": 0,
        "end": 15,
        "operator": "delete",
        "prefix": true,
        "argument": {
          "type": "ChainExpression",
          "start": 7,
          "end": 15,
          "expression": {
            "type": "MemberExpression",
            "start": 7,
            "end": 15,
            "object": {
              "type": "Identifier",
              "start": 7,
              "end": 10,
              "name": "obj"
            },
            "property": {
              "type": "Identifier",
              "start": 12,
              "end": 15,
              "name": "foo"
            },
            "computed": false,
            "optional": true
          }
        }
      }
    }
  ]
}
`, ast)
}

func TestOptionalChain22(t *testing.T) {
	ast, err := compile("new (obj?.foo)()")
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
        "type": "NewExpression",
        "start": 0,
        "end": 16,
        "callee": {
          "type": "ChainExpression",
          "start": 5,
          "end": 13,
          "expression": {
            "type": "MemberExpression",
            "start": 5,
            "end": 13,
            "object": {
              "type": "Identifier",
              "start": 5,
              "end": 8,
              "name": "obj"
            },
            "property": {
              "type": "Identifier",
              "start": 10,
              "end": 13,
              "name": "foo"
            },
            "computed": false,
            "optional": true
          }
        },
        "arguments": []
      }
    }
  ]
}
`, ast)
}

func TestOptionalChain23(t *testing.T) {
	ast, err := compile("(obj?.foo)`template`")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
        "type": "TaggedTemplateExpression",
        "start": 0,
        "end": 20,
        "tag": {
          "type": "ChainExpression",
          "start": 1,
          "end": 9,
          "expression": {
            "type": "MemberExpression",
            "start": 1,
            "end": 9,
            "object": {
              "type": "Identifier",
              "start": 1,
              "end": 4,
              "name": "obj"
            },
            "property": {
              "type": "Identifier",
              "start": 6,
              "end": 9,
              "name": "foo"
            },
            "computed": false,
            "optional": true
          }
        },
        "quasi": {
          "type": "TemplateLiteral",
          "start": 10,
          "end": 20,
          "expressions": [],
          "quasis": [
            {
              "type": "TemplateElement",
              "start": 11,
              "end": 19,
              "value": {
                "raw": "template",
                "cooked": "template"
              },
              "tail": true
            }
          ]
        }
      }
    }
  ]
}
`, ast)
}

func TestOptionalChain24(t *testing.T) {
	ast, err := compile("(obj?.foo).bar = 0")
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
        "start": 0,
        "end": 18,
        "operator": "=",
        "left": {
          "type": "MemberExpression",
          "start": 0,
          "end": 14,
          "object": {
            "type": "ChainExpression",
            "start": 1,
            "end": 9,
            "expression": {
              "type": "MemberExpression",
              "start": 1,
              "end": 9,
              "object": {
                "type": "Identifier",
                "start": 1,
                "end": 4,
                "name": "obj"
              },
              "property": {
                "type": "Identifier",
                "start": 6,
                "end": 9,
                "name": "foo"
              },
              "computed": false,
              "optional": true
            }
          },
          "property": {
            "type": "Identifier",
            "start": 11,
            "end": 14,
            "name": "bar"
          },
          "computed": false,
          "optional": false
        },
        "right": {
          "type": "Literal",
          "start": 17,
          "end": 18,
          "value": 0,
          "raw": "0"
        }
      }
    }
  ]
}
`, ast)
}

func TestOptionalChain25(t *testing.T) {
	ast, err := compile("(obj?.foo).bar++")
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
        "type": "UpdateExpression",
        "start": 0,
        "end": 16,
        "operator": "++",
        "prefix": false,
        "argument": {
          "type": "MemberExpression",
          "start": 0,
          "end": 14,
          "object": {
            "type": "ChainExpression",
            "start": 1,
            "end": 9,
            "expression": {
              "type": "MemberExpression",
              "start": 1,
              "end": 9,
              "object": {
                "type": "Identifier",
                "start": 1,
                "end": 4,
                "name": "obj"
              },
              "property": {
                "type": "Identifier",
                "start": 6,
                "end": 9,
                "name": "foo"
              },
              "computed": false,
              "optional": true
            }
          },
          "property": {
            "type": "Identifier",
            "start": 11,
            "end": 14,
            "name": "bar"
          },
          "computed": false,
          "optional": false
        }
      }
    }
  ]
}
`, ast)
}

func TestOptionalChain26(t *testing.T) {
	ast, err := compile("for ((obj?.foo).bar of []);")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 27,
  "body": [
    {
      "type": "ForOfStatement",
      "start": 0,
      "end": 27,
      "await": false,
      "left": {
        "type": "MemberExpression",
        "start": 5,
        "end": 19,
        "object": {
          "type": "ChainExpression",
          "start": 6,
          "end": 14,
          "expression": {
            "type": "MemberExpression",
            "start": 6,
            "end": 14,
            "object": {
              "type": "Identifier",
              "start": 6,
              "end": 9,
              "name": "obj"
            },
            "property": {
              "type": "Identifier",
              "start": 11,
              "end": 14,
              "name": "foo"
            },
            "computed": false,
            "optional": true
          }
        },
        "property": {
          "type": "Identifier",
          "start": 16,
          "end": 19,
          "name": "bar"
        },
        "computed": false,
        "optional": false
      },
      "right": {
        "type": "ArrayExpression",
        "start": 23,
        "end": 25,
        "elements": []
      },
      "body": {
        "type": "EmptyStatement",
        "start": 26,
        "end": 27
      }
    }
  ]
}
`, ast)
}

func TestOptionalChain27(t *testing.T) {
	ast, err := compile("(obj?.aaa).bbb")
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
        "type": "MemberExpression",
        "start": 0,
        "end": 14,
        "object": {
          "type": "ChainExpression",
          "start": 1,
          "end": 9,
          "expression": {
            "type": "MemberExpression",
            "start": 1,
            "end": 9,
            "object": {
              "type": "Identifier",
              "start": 1,
              "end": 4,
              "name": "obj"
            },
            "property": {
              "type": "Identifier",
              "start": 6,
              "end": 9,
              "name": "aaa"
            },
            "computed": false,
            "optional": true
          }
        },
        "property": {
          "type": "Identifier",
          "start": 11,
          "end": 14,
          "name": "bbb"
        },
        "computed": false,
        "optional": false
      }
    }
  ]
}
`, ast)
}

func TestOptionalChain28(t *testing.T) {
	ast, err := compile("(obj.foo.bar)?.buz")
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
        "type": "ChainExpression",
        "start": 0,
        "end": 18,
        "expression": {
          "type": "MemberExpression",
          "start": 0,
          "end": 18,
          "object": {
            "type": "MemberExpression",
            "start": 1,
            "end": 12,
            "object": {
              "type": "MemberExpression",
              "start": 1,
              "end": 8,
              "object": {
                "type": "Identifier",
                "start": 1,
                "end": 4,
                "name": "obj"
              },
              "property": {
                "type": "Identifier",
                "start": 5,
                "end": 8,
                "name": "foo"
              },
              "computed": false,
              "optional": false
            },
            "property": {
              "type": "Identifier",
              "start": 9,
              "end": 12,
              "name": "bar"
            },
            "computed": false,
            "optional": false
          },
          "property": {
            "type": "Identifier",
            "start": 15,
            "end": 18,
            "name": "buz"
          },
          "computed": false,
          "optional": true
        }
      }
    }
  ]
}
`, ast)
}

func TestOptionalChain29(t *testing.T) {
	ast, err := compile("(obj.foo())?.buz")
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
        "type": "ChainExpression",
        "start": 0,
        "end": 16,
        "expression": {
          "type": "MemberExpression",
          "start": 0,
          "end": 16,
          "object": {
            "type": "CallExpression",
            "start": 1,
            "end": 10,
            "callee": {
              "type": "MemberExpression",
              "start": 1,
              "end": 8,
              "object": {
                "type": "Identifier",
                "start": 1,
                "end": 4,
                "name": "obj"
              },
              "property": {
                "type": "Identifier",
                "start": 5,
                "end": 8,
                "name": "foo"
              },
              "computed": false,
              "optional": false
            },
            "arguments": [],
            "optional": false
          },
          "property": {
            "type": "Identifier",
            "start": 13,
            "end": 16,
            "name": "buz"
          },
          "computed": false,
          "optional": true
        }
      }
    }
  ]
}
`, ast)
}
