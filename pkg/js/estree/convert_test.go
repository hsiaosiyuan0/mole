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

	assert.EqualJson(t, `
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
            }
          },
          "name": "Object"
        },
        "arguments": null
      }
    }
  ]
}
  `, ast)
}

func Test2(t *testing.T) {
	ast, err := compile("this\n")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, trim(`
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
          }
        }
      }
    }
  ]
}
  `), ast)
}

func Test3(t *testing.T) {
	ast, err := compile("null\n")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
          }
        },
        "value": null
      }
    }
  ]
}
  `, ast)
}

func Test4(t *testing.T) {
	ast, err := compile("\n    42\n\n")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
          }
        },
        "value": 42
      }
    }
  ]
}
  `, ast)
}

func Test5(t *testing.T) {
	ast, err := compile("/foobar/")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
  `, ast)
}

func Test6(t *testing.T) {
	ast, err := compile("/[a-z]/g")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
  `, ast)
}

func Test7(t *testing.T) {
	ast, err := compile("(1 + 2 ) * 3")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
            }
          },
          "value": 3
        }
      }
    }
  ]
}
  `, ast)
}

func Test8(t *testing.T) {
	ast, err := compile("(x = 23)")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
            }
          },
          "value": 23
        }
      }
    }
  ]
}
  `, ast)
}

func Test9(t *testing.T) {
	ast, err := compile("x = []")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
      "column": 6
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
          "column": 6
        }
      },
      "expression": {
        "type": "AssignmentExpression",
        "loc": {
          "source": "",
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 6
          }
        },
        "operator": "=",
        "left": {
          "type": "Identifier",
          "loc": {
            "source": "",
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 1
            }
          },
          "name": "x"
        },
        "right": {
          "type": "ArrayExpression",
          "loc": {
            "source": "",
            "start": {
              "line": 1,
              "column": 4
            },
            "end": {
              "line": 1,
              "column": 6
            }
          },
          "elements": []
        }
      }
    }
  ]
}
  `, ast)
}

func Test10(t *testing.T) {
	ast, err := compile("x = [ 42 ]")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
        }
      },
      "expression": {
        "type": "AssignmentExpression",
        "loc": {
          "source": "",
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 10
          }
        },
        "operator": "=",
        "left": {
          "type": "Identifier",
          "loc": {
            "source": "",
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 1
            }
          },
          "name": "x"
        },
        "right": {
          "type": "ArrayExpression",
          "loc": {
            "source": "",
            "start": {
              "line": 1,
              "column": 4
            },
            "end": {
              "line": 1,
              "column": 10
            }
          },
          "elements": [
            {
              "type": "Literal",
              "loc": {
                "source": "",
                "start": {
                  "line": 1,
                  "column": 6
                },
                "end": {
                  "line": 1,
                  "column": 8
                }
              },
              "value": 42
            }
          ]
        }
      }
    }
  ]
}
  `, ast)
}

func Test11(t *testing.T) {
	ast, err := compile("x = [ 42, ]")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
      "column": 11
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
          "column": 11
        }
      },
      "expression": {
        "type": "AssignmentExpression",
        "loc": {
          "source": "",
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 11
          }
        },
        "operator": "=",
        "left": {
          "type": "Identifier",
          "loc": {
            "source": "",
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 1
            }
          },
          "name": "x"
        },
        "right": {
          "type": "ArrayExpression",
          "loc": {
            "source": "",
            "start": {
              "line": 1,
              "column": 4
            },
            "end": {
              "line": 1,
              "column": 11
            }
          },
          "elements": [
            {
              "type": "Literal",
              "loc": {
                "source": "",
                "start": {
                  "line": 1,
                  "column": 6
                },
                "end": {
                  "line": 1,
                  "column": 8
                }
              },
              "value": 42
            }
          ]
        }
      }
    }
  ]
}
  `, ast)
}

func Test12(t *testing.T) {
	ast, err := compile("x = [ ,, 42 ]")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
      "column": 13
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
          "column": 13
        }
      },
      "expression": {
        "type": "AssignmentExpression",
        "loc": {
          "source": "",
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 13
          }
        },
        "operator": "=",
        "left": {
          "type": "Identifier",
          "loc": {
            "source": "",
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 1
            }
          },
          "name": "x"
        },
        "right": {
          "type": "ArrayExpression",
          "loc": {
            "source": "",
            "start": {
              "line": 1,
              "column": 4
            },
            "end": {
              "line": 1,
              "column": 13
            }
          },
          "elements": [
            null,
            null,
            {
              "type": "Literal",
              "loc": {
                "source": "",
                "start": {
                  "line": 1,
                  "column": 9
                },
                "end": {
                  "line": 1,
                  "column": 11
                }
              },
              "value": 42
            }
          ]
        }
      }
    }
  ]
}
  `, ast)
}

func Test13(t *testing.T) {
	ast, err := compile("x = [ 1, 2, 3, ]")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
      "column": 16
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
          "column": 16
        }
      },
      "expression": {
        "type": "AssignmentExpression",
        "loc": {
          "source": "",
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 16
          }
        },
        "operator": "=",
        "left": {
          "type": "Identifier",
          "loc": {
            "source": "",
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 1
            }
          },
          "name": "x"
        },
        "right": {
          "type": "ArrayExpression",
          "loc": {
            "source": "",
            "start": {
              "line": 1,
              "column": 4
            },
            "end": {
              "line": 1,
              "column": 16
            }
          },
          "elements": [
            {
              "type": "Literal",
              "loc": {
                "source": "",
                "start": {
                  "line": 1,
                  "column": 6
                },
                "end": {
                  "line": 1,
                  "column": 7
                }
              },
              "value": 1
            },
            {
              "type": "Literal",
              "loc": {
                "source": "",
                "start": {
                  "line": 1,
                  "column": 9
                },
                "end": {
                  "line": 1,
                  "column": 10
                }
              },
              "value": 2
            },
            {
              "type": "Literal",
              "loc": {
                "source": "",
                "start": {
                  "line": 1,
                  "column": 12
                },
                "end": {
                  "line": 1,
                  "column": 13
                }
              },
              "value": 3
            }
          ]
        }
      }
    }
  ]
}
  `, ast)
}

func Test14(t *testing.T) {
	ast, err := compile("x = [ 1, 2,, 3, ]")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
      "column": 17
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
          "column": 17
        }
      },
      "expression": {
        "type": "AssignmentExpression",
        "loc": {
          "source": "",
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 17
          }
        },
        "operator": "=",
        "left": {
          "type": "Identifier",
          "loc": {
            "source": "",
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 1
            }
          },
          "name": "x"
        },
        "right": {
          "type": "ArrayExpression",
          "loc": {
            "source": "",
            "start": {
              "line": 1,
              "column": 4
            },
            "end": {
              "line": 1,
              "column": 17
            }
          },
          "elements": [
            {
              "type": "Literal",
              "loc": {
                "source": "",
                "start": {
                  "line": 1,
                  "column": 6
                },
                "end": {
                  "line": 1,
                  "column": 7
                }
              },
              "value": 1
            },
            {
              "type": "Literal",
              "loc": {
                "source": "",
                "start": {
                  "line": 1,
                  "column": 9
                },
                "end": {
                  "line": 1,
                  "column": 10
                }
              },
              "value": 2
            },
            null,
            {
              "type": "Literal",
              "loc": {
                "source": "",
                "start": {
                  "line": 1,
                  "column": 13
                },
                "end": {
                  "line": 1,
                  "column": 14
                }
              },
              "value": 3
            }
          ]
        }
      }
    }
  ]
}
  `, ast)
}

func Test15(t *testing.T) {
	ast, err := compile("日本語 = []")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
        }
      },
      "expression": {
        "type": "AssignmentExpression",
        "loc": {
          "source": "",
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 8
          }
        },
        "operator": "=",
        "left": {
          "type": "Identifier",
          "loc": {
            "source": "",
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 3
            }
          },
          "name": "日本語"
        },
        "right": {
          "type": "ArrayExpression",
          "loc": {
            "source": "",
            "start": {
              "line": 1,
              "column": 6
            },
            "end": {
              "line": 1,
              "column": 8
            }
          },
          "elements": []
        }
      }
    }
  ]
}
  `, ast)
}

func Test16(t *testing.T) {
	ast, err := compile("T‿ = []")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
      "column": 7
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
          "column": 7
        }
      },
      "expression": {
        "type": "AssignmentExpression",
        "loc": {
          "source": "",
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 7
          }
        },
        "operator": "=",
        "left": {
          "type": "Identifier",
          "loc": {
            "source": "",
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 2
            }
          },
          "name": "T‿"
        },
        "right": {
          "type": "ArrayExpression",
          "loc": {
            "source": "",
            "start": {
              "line": 1,
              "column": 5
            },
            "end": {
              "line": 1,
              "column": 7
            }
          },
          "elements": []
        }
      }
    }
  ]
}
  `, ast)
}

func Test17(t *testing.T) {
	ast, err := compile("T\u200c = []")
	assert.Equal(t, nil, err, "should be prog ok")

	//lint:ignore ST1018 below `T` contains two codepoints `T\u200c`
	assert.EqualJson(t, `
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
      "column": 7
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
          "column": 7
        }
      },
      "expression": {
        "type": "AssignmentExpression",
        "loc": {
          "source": "",
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 7
          }
        },
        "operator": "=",
        "left": {
          "type": "Identifier",
          "loc": {
            "source": "",
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 2
            }
          },
          "name": "T‌"
        },
        "right": {
          "type": "ArrayExpression",
          "loc": {
            "source": "",
            "start": {
              "line": 1,
              "column": 5
            },
            "end": {
              "line": 1,
              "column": 7
            }
          },
          "elements": []
        }
      }
    }
  ]
}
  `, ast)
}

func Test18(t *testing.T) {
	ast, err := compile("ⅣⅡ = []")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
      "column": 7
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
          "column": 7
        }
      },
      "expression": {
        "type": "AssignmentExpression",
        "loc": {
          "source": "",
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 7
          }
        },
        "operator": "=",
        "left": {
          "type": "Identifier",
          "loc": {
            "source": "",
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 2
            }
          },
          "name": "ⅣⅡ"
        },
        "right": {
          "type": "ArrayExpression",
          "loc": {
            "source": "",
            "start": {
              "line": 1,
              "column": 5
            },
            "end": {
              "line": 1,
              "column": 7
            }
          },
          "elements": []
        }
      }
    }
  ]
}
  `, ast)
}

func Test19(t *testing.T) {
	// the `u200a` is after `ⅣⅡ`
	ast, err := compile("ⅣⅡ = []")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
      "column": 7
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
          "column": 7
        }
      },
      "expression": {
        "type": "AssignmentExpression",
        "loc": {
          "source": "",
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 7
          }
        },
        "operator": "=",
        "left": {
          "type": "Identifier",
          "loc": {
            "source": "",
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 2
            }
          },
          "name": "ⅣⅡ"
        },
        "right": {
          "type": "ArrayExpression",
          "loc": {
            "source": "",
            "start": {
              "line": 1,
              "column": 5
            },
            "end": {
              "line": 1,
              "column": 7
            }
          },
          "elements": []
        }
      }
    }
  ]
}
  `, ast)
}

func Test20(t *testing.T) {
	ast, err := compile("x = {}")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
      "column": 6
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
          "column": 6
        }
      },
      "expression": {
        "type": "AssignmentExpression",
        "loc": {
          "source": "",
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 6
          }
        },
        "operator": "=",
        "left": {
          "type": "Identifier",
          "loc": {
            "source": "",
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 1
            }
          },
          "name": "x"
        },
        "right": {
          "type": "ObjectExpression",
          "loc": {
            "source": "",
            "start": {
              "line": 1,
              "column": 4
            },
            "end": {
              "line": 1,
              "column": 6
            }
          },
          "properties": []
        }
      }
    }
  ]
}
  `, ast)
}

func Test21(t *testing.T) {
	ast, err := compile("x = { }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
      "column": 7
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
          "column": 7
        }
      },
      "expression": {
        "type": "AssignmentExpression",
        "loc": {
          "source": "",
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 7
          }
        },
        "operator": "=",
        "left": {
          "type": "Identifier",
          "loc": {
            "source": "",
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 1
            }
          },
          "name": "x"
        },
        "right": {
          "type": "ObjectExpression",
          "loc": {
            "source": "",
            "start": {
              "line": 1,
              "column": 4
            },
            "end": {
              "line": 1,
              "column": 7
            }
          },
          "properties": []
        }
      }
    }
  ]
}
  `, ast)
}

func Test22(t *testing.T) {
	ast, err := compile("x = { if: 42 }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
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
          "column": 14
        }
      },
      "expression": {
        "type": "AssignmentExpression",
        "loc": {
          "source": "",
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 14
          }
        },
        "operator": "=",
        "left": {
          "type": "Identifier",
          "loc": {
            "source": "",
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 1
            }
          },
          "name": "x"
        },
        "right": {
          "type": "ObjectExpression",
          "loc": {
            "source": "",
            "start": {
              "line": 1,
              "column": 4
            },
            "end": {
              "line": 1,
              "column": 14
            }
          },
          "properties": [
            {
              "type": "Property",
              "loc": {
                "source": "",
                "start": {
                  "line": 1,
                  "column": 6
                },
                "end": {
                  "line": 1,
                  "column": 12
                }
              },
              "key": {
                "type": "Identifier",
                "loc": {
                  "source": "",
                  "start": {
                    "line": 1,
                    "column": 6
                  },
                  "end": {
                    "line": 1,
                    "column": 8
                  }
                },
                "name": "if"
              },
              "value": {
                "type": "Literal",
                "loc": {
                  "source": "",
                  "start": {
                    "line": 1,
                    "column": 10
                  },
                  "end": {
                    "line": 1,
                    "column": 12
                  }
                },
                "value": 42
              },
              "kind": "init",
              "method": false,
              "shorthand": false,
              "computed": false
            }
          ]
        }
      }
    }
  ]
}
  `, ast)
}
