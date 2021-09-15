package estree

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/hsiaosiyuan0/mole/pkg/assert"
	"github.com/hsiaosiyuan0/mole/pkg/js/parser"
)

func compile(code string) (string, error) {
	s := parser.NewSource("", code)
	p := parser.NewParser(s, make([]string, 0))
	ast, err := p.Prog()
	if err != nil {
		return "", err
	}

	b, err := json.Marshal(program(ast.(*parser.Prog)))
	if err != nil {
		return "", err
	}
	var out bytes.Buffer
	json.Indent(&out, b, "", "  ")

	return out.String(), nil
}

func trim(str string) string {
	return strings.Trim(str, "\n ")
}

func TestNewExpr(t *testing.T) {
	ast, err := compile(`new Object`)
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualString(t, trim(`
{
  "type": "Program",
  "loc": {
    "source": "",
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 10
    },
    "range": {
      "start": 0,
      "end": 10
    }
  },
  "sourceType": "",
  "body": [
    {
      "type": "ExpressionStatement",
      "loc": {
        "source": "",
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 10
        },
        "range": {
          "start": 0,
          "end": 10
        }
      },
      "expression": {
        "type": "NewExpression",
        "loc": {
          "source": "",
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 10
          },
          "range": {
            "start": 0,
            "end": 10
          }
        },
        "callee": {
          "type": "Identifier",
          "loc": {
            "source": "",
            "start": {
              "line": 1,
              "column": 4
            },
            "end": {
              "line": 1,
              "column": 10
            },
            "range": {
              "start": 4,
              "end": 10
            }
          },
          "name": "Object"
        },
        "arguments": null
      }
    }
  ]
}
  `), ast, "should pass")
}

func TestThisExpr(t *testing.T) {
	ast, err := compile("this\n")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualString(t, trim(`
{
  "type": "Program",
  "loc": {
    "source": "",
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 4
    },
    "range": {
      "start": 0,
      "end": 4
    }
  },
  "sourceType": "",
  "body": [
    {
      "type": "ExpressionStatement",
      "loc": {
        "source": "",
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 4
        },
        "range": {
          "start": 0,
          "end": 4
        }
      },
      "expression": {
        "type": "ThisExpression",
        "loc": {
          "source": "",
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 4
          },
          "range": {
            "start": 0,
            "end": 4
          }
        }
      }
    }
  ]
}
  `), ast, "should pass")
}

func TestNull(t *testing.T) {
	ast, err := compile("null\n")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualString(t, trim(`
{
  "type": "Program",
  "loc": {
    "source": "",
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 4
    },
    "range": {
      "start": 0,
      "end": 4
    }
  },
  "sourceType": "",
  "body": [
    {
      "type": "ExpressionStatement",
      "loc": {
        "source": "",
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 4
        },
        "range": {
          "start": 0,
          "end": 4
        }
      },
      "expression": {
        "type": "Literal",
        "loc": {
          "source": "",
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 4
          },
          "range": {
            "start": 0,
            "end": 4
          }
        },
        "value": null
      }
    }
  ]
}
  `), ast, "should pass")
}
