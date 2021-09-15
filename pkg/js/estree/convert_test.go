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

func Test1(t *testing.T) {
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

func Test2(t *testing.T) {
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

func Test3(t *testing.T) {
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

func Test4(t *testing.T) {
	ast, err := compile("\n    42\n\n")
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
      "line": 2,
      "column": 6
    },
    "range": {
      "start": 0,
      "end": 7
    }
  },
  "sourceType": "",
  "body": [
    {
      "type": "ExpressionStatement",
      "loc": {
        "source": "",
        "start": {
          "line": 2,
          "column": 4
        },
        "end": {
          "line": 2,
          "column": 6
        },
        "range": {
          "start": 5,
          "end": 7
        }
      },
      "expression": {
        "type": "Literal",
        "loc": {
          "source": "",
          "start": {
            "line": 2,
            "column": 4
          },
          "end": {
            "line": 2,
            "column": 6
          },
          "range": {
            "start": 5,
            "end": 7
          }
        },
        "value": 42
      }
    }
  ]
}
  `), ast, "should pass")
}

func Test5(t *testing.T) {
	ast, err := compile("/foobar/")
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
      "column": 8
    },
    "range": {
      "start": 0,
      "end": 8
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
          "column": 8
        },
        "range": {
          "start": 0,
          "end": 8
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
            "column": 8
          },
          "range": {
            "start": 0,
            "end": 8
          }
        },
        "value": null,
        "regexp": {
          "pattern": "foobar",
          "flags": ""
        }
      }
    }
  ]
}
  `), ast, "should pass")
}

func Test6(t *testing.T) {
	ast, err := compile("/[a-z]/g")
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
      "column": 8
    },
    "range": {
      "start": 0,
      "end": 8
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
          "column": 8
        },
        "range": {
          "start": 0,
          "end": 8
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
            "column": 8
          },
          "range": {
            "start": 0,
            "end": 8
          }
        },
        "value": null,
        "regexp": {
          "pattern": "[a-z]",
          "flags": "g"
        }
      }
    }
  ]
}
  `), ast, "should pass")
}

func Test7(t *testing.T) {
	ast, err := compile("(1 + 2 ) * 3")
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
      "column": 12
    },
    "range": {
      "start": 0,
      "end": 12
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
          "column": 12
        },
        "range": {
          "start": 0,
          "end": 12
        }
      },
      "expression": {
        "type": "BinaryExpression",
        "loc": {
          "source": "",
          "start": {
            "line": 1,
            "column": 1
          },
          "end": {
            "line": 1,
            "column": 12
          },
          "range": {
            "start": 1,
            "end": 12
          }
        },
        "operator": "*",
        "left": {
          "type": "BinaryExpression",
          "loc": {
            "source": "",
            "start": {
              "line": 1,
              "column": 1
            },
            "end": {
              "line": 1,
              "column": 12
            },
            "range": {
              "start": 1,
              "end": 12
            }
          },
          "operator": "+",
          "left": {
            "type": "Literal",
            "loc": {
              "source": "",
              "start": {
                "line": 1,
                "column": 1
              },
              "end": {
                "line": 1,
                "column": 12
              },
              "range": {
                "start": 1,
                "end": 12
              }
            },
            "value": 1
          },
          "right": {
            "type": "Literal",
            "loc": {
              "source": "",
              "start": {
                "line": 1,
                "column": 5
              },
              "end": {
                "line": 1,
                "column": 6
              },
              "range": {
                "start": 5,
                "end": 6
              }
            },
            "value": 2
          }
        },
        "right": {
          "type": "Literal",
          "loc": {
            "source": "",
            "start": {
              "line": 1,
              "column": 11
            },
            "end": {
              "line": 1,
              "column": 12
            },
            "range": {
              "start": 11,
              "end": 12
            }
          },
          "value": 3
        }
      }
    }
  ]
}
  `), ast, "should pass")
}

func Test8(t *testing.T) {
	ast, err := compile("(x = 23)")
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
      "column": 8
    },
    "range": {
      "start": 0,
      "end": 8
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
          "column": 8
        },
        "range": {
          "start": 0,
          "end": 8
        }
      },
      "expression": {
        "type": "AssignmentExpression",
        "loc": {
          "source": "",
          "start": {
            "line": 1,
            "column": 1
          },
          "end": {
            "line": 1,
            "column": 7
          },
          "range": {
            "start": 1,
            "end": 7
          }
        },
        "operator": "=",
        "left": {
          "type": "Identifier",
          "loc": {
            "source": "",
            "start": {
              "line": 1,
              "column": 1
            },
            "end": {
              "line": 1,
              "column": 2
            },
            "range": {
              "start": 1,
              "end": 2
            }
          },
          "name": "x"
        },
        "right": {
          "type": "Literal",
          "loc": {
            "source": "",
            "start": {
              "line": 1,
              "column": 5
            },
            "end": {
              "line": 1,
              "column": 7
            },
            "range": {
              "start": 5,
              "end": 7
            }
          },
          "value": 23
        }
      }
    }
  ]
}
  `), ast, "should pass")
}
