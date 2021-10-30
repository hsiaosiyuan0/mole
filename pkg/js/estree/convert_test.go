package estree

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/hsiaosiyuan0/mole/pkg/assert"
	"github.com/hsiaosiyuan0/mole/pkg/js/parser"
)

func compile(code string) (string, error) {
	s := parser.NewSource("", code)
	p := parser.NewParser(s, parser.NewParserOpts())
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

func Test1(t *testing.T) {
	ast, err := compile(`new Object`)
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 10,
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 10
    }
  },
  "body": [
    {
      "type": "ExpressionStatement",
      "start": 0,
      "end": 10,
      "loc": {
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
        "start": 0,
        "end": 10,
        "loc": {
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
          "start": 4,
          "end": 10,
          "loc": {
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
        "arguments": []
      }
    }
  ]
}
  `, ast)
}

func Test2(t *testing.T) {
	ast, err := compile("this\n")
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
      "column": 0
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
  `, ast)
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
      "line": 2,
      "column": 0
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
      "line": 4,
      "column": 0
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
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "BinaryExpression",
        "left": {
          "type": "BinaryExpression",
          "left": {
            "type": "Literal",
            "value": 1,
            "loc": {
              "start": {
                "line": 1,
                "column": 1
              },
              "end": {
                "line": 1,
                "column": 2
              }
            }
          },
          "operator": "+",
          "right": {
            "type": "Literal",
            "value": 2,
            "loc": {
              "start": {
                "line": 1,
                "column": 5
              },
              "end": {
                "line": 1,
                "column": 6
              }
            }
          },
          "loc": {
            "start": {
              "line": 1,
              "column": 1
            },
            "end": {
              "line": 1,
              "column": 6
            }
          }
        },
        "operator": "*",
        "right": {
          "type": "Literal",
          "value": 3,
          "loc": {
            "start": {
              "line": 1,
              "column": 11
            },
            "end": {
              "line": 1,
              "column": 12
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 12
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 12
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 12
    }
  }
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

func Test23(t *testing.T) {
	ast, err := compile("x = { true: 42 }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "AssignmentExpression",
        "operator": "=",
        "left": {
          "type": "Identifier",
          "name": "x",
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 1
            }
          }
        },
        "right": {
          "type": "ObjectExpression",
          "properties": [
            {
              "type": "Property",
              "key": {
                "type": "Identifier",
                "name": "true",
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 6
                  },
                  "end": {
                    "line": 1,
                    "column": 10
                  }
                }
              },
              "value": {
                "type": "Literal",
                "value": 42,
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 12
                  },
                  "end": {
                    "line": 1,
                    "column": 14
                  }
                }
              },
              "kind": "init"
            }
          ],
          "loc": {
            "start": {
              "line": 1,
              "column": 4
            },
            "end": {
              "line": 1,
              "column": 16
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 16
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 16
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 16
    }
  }
}
  `, ast)
}

func Test24(t *testing.T) {
	ast, err := compile("x = { false: 42 }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "AssignmentExpression",
        "operator": "=",
        "left": {
          "type": "Identifier",
          "name": "x",
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 1
            }
          }
        },
        "right": {
          "type": "ObjectExpression",
          "properties": [
            {
              "type": "Property",
              "key": {
                "type": "Identifier",
                "name": "false",
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 6
                  },
                  "end": {
                    "line": 1,
                    "column": 11
                  }
                }
              },
              "value": {
                "type": "Literal",
                "value": 42,
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 13
                  },
                  "end": {
                    "line": 1,
                    "column": 15
                  }
                }
              },
              "kind": "init"
            }
          ],
          "loc": {
            "start": {
              "line": 1,
              "column": 4
            },
            "end": {
              "line": 1,
              "column": 17
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 17
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 17
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 17
    }
  }
}
  `, ast)
}

func Test25(t *testing.T) {
	ast, err := compile("x = { null: 42 }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "AssignmentExpression",
        "operator": "=",
        "left": {
          "type": "Identifier",
          "name": "x",
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 1
            }
          }
        },
        "right": {
          "type": "ObjectExpression",
          "properties": [
            {
              "type": "Property",
              "key": {
                "type": "Identifier",
                "name": "null",
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 6
                  },
                  "end": {
                    "line": 1,
                    "column": 10
                  }
                }
              },
              "value": {
                "type": "Literal",
                "value": 42,
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 12
                  },
                  "end": {
                    "line": 1,
                    "column": 14
                  }
                }
              },
              "kind": "init"
            }
          ],
          "loc": {
            "start": {
              "line": 1,
              "column": 4
            },
            "end": {
              "line": 1,
              "column": 16
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 16
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 16
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 16
    }
  }
}
  `, ast)
}

func Test26(t *testing.T) {
	ast, err := compile("x = { \"answer\": 42 }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "AssignmentExpression",
        "operator": "=",
        "left": {
          "type": "Identifier",
          "name": "x",
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 1
            }
          }
        },
        "right": {
          "type": "ObjectExpression",
          "properties": [
            {
              "type": "Property",
              "key": {
                "type": "Literal",
                "value": "answer",
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 6
                  },
                  "end": {
                    "line": 1,
                    "column": 14
                  }
                }
              },
              "value": {
                "type": "Literal",
                "value": 42,
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 16
                  },
                  "end": {
                    "line": 1,
                    "column": 18
                  }
                }
              },
              "kind": "init"
            }
          ],
          "loc": {
            "start": {
              "line": 1,
              "column": 4
            },
            "end": {
              "line": 1,
              "column": 20
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 20
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 20
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 20
    }
  }
}
  `, ast)
}

func Test27(t *testing.T) {
	ast, err := compile("x = { x: 1, x: 2 }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "AssignmentExpression",
        "operator": "=",
        "left": {
          "type": "Identifier",
          "name": "x",
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 1
            }
          }
        },
        "right": {
          "type": "ObjectExpression",
          "properties": [
            {
              "type": "Property",
              "key": {
                "type": "Identifier",
                "name": "x",
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 6
                  },
                  "end": {
                    "line": 1,
                    "column": 7
                  }
                }
              },
              "value": {
                "type": "Literal",
                "value": 1,
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 9
                  },
                  "end": {
                    "line": 1,
                    "column": 10
                  }
                }
              },
              "kind": "init"
            },
            {
              "type": "Property",
              "key": {
                "type": "Identifier",
                "name": "x",
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 12
                  },
                  "end": {
                    "line": 1,
                    "column": 13
                  }
                }
              },
              "value": {
                "type": "Literal",
                "value": 2,
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 15
                  },
                  "end": {
                    "line": 1,
                    "column": 16
                  }
                }
              },
              "kind": "init"
            }
          ],
          "loc": {
            "start": {
              "line": 1,
              "column": 4
            },
            "end": {
              "line": 1,
              "column": 18
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 18
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 18
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 18
    }
  }
}
  `, ast)
}

func Test28(t *testing.T) {
	ast, err := compile("x = { get width() { return m_width } }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "AssignmentExpression",
        "operator": "=",
        "left": {
          "type": "Identifier",
          "name": "x",
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 1
            }
          }
        },
        "right": {
          "type": "ObjectExpression",
          "properties": [
            {
              "type": "Property",
              "key": {
                "type": "Identifier",
                "name": "width",
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 10
                  },
                  "end": {
                    "line": 1,
                    "column": 15
                  }
                }
              },
              "kind": "get",
              "value": {
                "type": "FunctionExpression",
                "id": null,
                "params": [],
                "body": {
                  "type": "BlockStatement",
                  "body": [
                    {
                      "type": "ReturnStatement",
                      "argument": {
                        "type": "Identifier",
                        "name": "m_width",
                        "loc": {
                          "start": {
                            "line": 1,
                            "column": 27
                          },
                          "end": {
                            "line": 1,
                            "column": 34
                          }
                        }
                      },
                      "loc": {
                        "start": {
                          "line": 1,
                          "column": 20
                        },
                        "end": {
                          "line": 1,
                          "column": 34
                        }
                      }
                    }
                  ],
                  "loc": {
                    "start": {
                      "line": 1,
                      "column": 18
                    },
                    "end": {
                      "line": 1,
                      "column": 36
                    }
                  }
                },
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 15
                  },
                  "end": {
                    "line": 1,
                    "column": 36
                  }
                }
              }
            }
          ],
          "loc": {
            "start": {
              "line": 1,
              "column": 4
            },
            "end": {
              "line": 1,
              "column": 38
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 38
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 38
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 38
    }
  }
}
  `, ast)
}

func Test29(t *testing.T) {
	ast, err := compile("x = { get undef() {} }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "AssignmentExpression",
        "operator": "=",
        "left": {
          "type": "Identifier",
          "name": "x",
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 1
            }
          }
        },
        "right": {
          "type": "ObjectExpression",
          "properties": [
            {
              "type": "Property",
              "key": {
                "type": "Identifier",
                "name": "undef",
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 10
                  },
                  "end": {
                    "line": 1,
                    "column": 15
                  }
                }
              },
              "kind": "get",
              "value": {
                "type": "FunctionExpression",
                "id": null,
                "params": [],
                "body": {
                  "type": "BlockStatement",
                  "body": [],
                  "loc": {
                    "start": {
                      "line": 1,
                      "column": 18
                    },
                    "end": {
                      "line": 1,
                      "column": 20
                    }
                  }
                },
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 15
                  },
                  "end": {
                    "line": 1,
                    "column": 20
                  }
                }
              }
            }
          ],
          "loc": {
            "start": {
              "line": 1,
              "column": 4
            },
            "end": {
              "line": 1,
              "column": 22
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 22
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 22
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 22
    }
  }
}
  `, ast)
}

func Test30(t *testing.T) {
	ast, err := compile("x = { get if() {} }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "AssignmentExpression",
        "operator": "=",
        "left": {
          "type": "Identifier",
          "name": "x",
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 1
            }
          }
        },
        "right": {
          "type": "ObjectExpression",
          "properties": [
            {
              "type": "Property",
              "key": {
                "type": "Identifier",
                "name": "if",
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 10
                  },
                  "end": {
                    "line": 1,
                    "column": 12
                  }
                }
              },
              "kind": "get",
              "value": {
                "type": "FunctionExpression",
                "id": null,
                "params": [],
                "body": {
                  "type": "BlockStatement",
                  "body": [],
                  "loc": {
                    "start": {
                      "line": 1,
                      "column": 15
                    },
                    "end": {
                      "line": 1,
                      "column": 17
                    }
                  }
                },
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 12
                  },
                  "end": {
                    "line": 1,
                    "column": 17
                  }
                }
              }
            }
          ],
          "loc": {
            "start": {
              "line": 1,
              "column": 4
            },
            "end": {
              "line": 1,
              "column": 19
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 19
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 19
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 19
    }
  }
}
  `, ast)
}

func Test31(t *testing.T) {
	ast, err := compile("x = { get true() {} }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "AssignmentExpression",
        "operator": "=",
        "left": {
          "type": "Identifier",
          "name": "x",
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 1
            }
          }
        },
        "right": {
          "type": "ObjectExpression",
          "properties": [
            {
              "type": "Property",
              "key": {
                "type": "Identifier",
                "name": "true",
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 10
                  },
                  "end": {
                    "line": 1,
                    "column": 14
                  }
                }
              },
              "kind": "get",
              "value": {
                "type": "FunctionExpression",
                "id": null,
                "params": [],
                "body": {
                  "type": "BlockStatement",
                  "body": [],
                  "loc": {
                    "start": {
                      "line": 1,
                      "column": 17
                    },
                    "end": {
                      "line": 1,
                      "column": 19
                    }
                  }
                },
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 14
                  },
                  "end": {
                    "line": 1,
                    "column": 19
                  }
                }
              }
            }
          ],
          "loc": {
            "start": {
              "line": 1,
              "column": 4
            },
            "end": {
              "line": 1,
              "column": 21
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 21
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 21
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 21
    }
  }
}
  `, ast)
}

func Test32(t *testing.T) {
	ast, err := compile("x = { get false() {} }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "AssignmentExpression",
        "operator": "=",
        "left": {
          "type": "Identifier",
          "name": "x",
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 1
            }
          }
        },
        "right": {
          "type": "ObjectExpression",
          "properties": [
            {
              "type": "Property",
              "key": {
                "type": "Identifier",
                "name": "false",
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 10
                  },
                  "end": {
                    "line": 1,
                    "column": 15
                  }
                }
              },
              "kind": "get",
              "value": {
                "type": "FunctionExpression",
                "id": null,
                "params": [],
                "body": {
                  "type": "BlockStatement",
                  "body": [],
                  "loc": {
                    "start": {
                      "line": 1,
                      "column": 18
                    },
                    "end": {
                      "line": 1,
                      "column": 20
                    }
                  }
                },
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 15
                  },
                  "end": {
                    "line": 1,
                    "column": 20
                  }
                }
              }
            }
          ],
          "loc": {
            "start": {
              "line": 1,
              "column": 4
            },
            "end": {
              "line": 1,
              "column": 22
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 22
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 22
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 22
    }
  }
}
  `, ast)
}

func Test33(t *testing.T) {
	ast, err := compile("x = { get null() {} }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "AssignmentExpression",
        "operator": "=",
        "left": {
          "type": "Identifier",
          "name": "x",
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 1
            }
          }
        },
        "right": {
          "type": "ObjectExpression",
          "properties": [
            {
              "type": "Property",
              "key": {
                "type": "Identifier",
                "name": "null",
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 10
                  },
                  "end": {
                    "line": 1,
                    "column": 14
                  }
                }
              },
              "kind": "get",
              "value": {
                "type": "FunctionExpression",
                "id": null,
                "params": [],
                "body": {
                  "type": "BlockStatement",
                  "body": [],
                  "loc": {
                    "start": {
                      "line": 1,
                      "column": 17
                    },
                    "end": {
                      "line": 1,
                      "column": 19
                    }
                  }
                },
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 14
                  },
                  "end": {
                    "line": 1,
                    "column": 19
                  }
                }
              }
            }
          ],
          "loc": {
            "start": {
              "line": 1,
              "column": 4
            },
            "end": {
              "line": 1,
              "column": 21
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 21
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 21
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 21
    }
  }
}
  `, ast)
}

func Test34(t *testing.T) {
	ast, err := compile("x = { get \"undef\"() {} }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "AssignmentExpression",
        "operator": "=",
        "left": {
          "type": "Identifier",
          "name": "x",
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 1
            }
          }
        },
        "right": {
          "type": "ObjectExpression",
          "properties": [
            {
              "type": "Property",
              "key": {
                "type": "Literal",
                "value": "undef",
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 10
                  },
                  "end": {
                    "line": 1,
                    "column": 17
                  }
                }
              },
              "kind": "get",
              "value": {
                "type": "FunctionExpression",
                "id": null,
                "params": [],
                "body": {
                  "type": "BlockStatement",
                  "body": [],
                  "loc": {
                    "start": {
                      "line": 1,
                      "column": 20
                    },
                    "end": {
                      "line": 1,
                      "column": 22
                    }
                  }
                },
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 17
                  },
                  "end": {
                    "line": 1,
                    "column": 22
                  }
                }
              }
            }
          ],
          "loc": {
            "start": {
              "line": 1,
              "column": 4
            },
            "end": {
              "line": 1,
              "column": 24
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 24
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 24
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 24
    }
  }
}
  `, ast)
}

func Test35(t *testing.T) {
	ast, err := compile("x = { get 10() {} }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "AssignmentExpression",
        "operator": "=",
        "left": {
          "type": "Identifier",
          "name": "x",
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 1
            }
          }
        },
        "right": {
          "type": "ObjectExpression",
          "properties": [
            {
              "type": "Property",
              "key": {
                "type": "Literal",
                "value": 10,
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 10
                  },
                  "end": {
                    "line": 1,
                    "column": 12
                  }
                }
              },
              "kind": "get",
              "value": {
                "type": "FunctionExpression",
                "id": null,
                "params": [],
                "body": {
                  "type": "BlockStatement",
                  "body": [],
                  "loc": {
                    "start": {
                      "line": 1,
                      "column": 15
                    },
                    "end": {
                      "line": 1,
                      "column": 17
                    }
                  }
                },
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 12
                  },
                  "end": {
                    "line": 1,
                    "column": 17
                  }
                }
              }
            }
          ],
          "loc": {
            "start": {
              "line": 1,
              "column": 4
            },
            "end": {
              "line": 1,
              "column": 19
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 19
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 19
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 19
    }
  }
}
  `, ast)
}

func Test36(t *testing.T) {
	ast, err := compile("x = { get 10() {} }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "AssignmentExpression",
        "operator": "=",
        "left": {
          "type": "Identifier",
          "name": "x",
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 1
            }
          }
        },
        "right": {
          "type": "ObjectExpression",
          "properties": [
            {
              "type": "Property",
              "key": {
                "type": "Literal",
                "value": 10,
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 10
                  },
                  "end": {
                    "line": 1,
                    "column": 12
                  }
                }
              },
              "kind": "get",
              "value": {
                "type": "FunctionExpression",
                "id": null,
                "params": [],
                "body": {
                  "type": "BlockStatement",
                  "body": [],
                  "loc": {
                    "start": {
                      "line": 1,
                      "column": 15
                    },
                    "end": {
                      "line": 1,
                      "column": 17
                    }
                  }
                },
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 12
                  },
                  "end": {
                    "line": 1,
                    "column": 17
                  }
                }
              }
            }
          ],
          "loc": {
            "start": {
              "line": 1,
              "column": 4
            },
            "end": {
              "line": 1,
              "column": 19
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 19
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 19
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 19
    }
  }
}
  `, ast)
}

func Test37(t *testing.T) {
	ast, err := compile("x = { set width(w) { m_width = w } }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "AssignmentExpression",
        "operator": "=",
        "left": {
          "type": "Identifier",
          "name": "x",
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 1
            }
          }
        },
        "right": {
          "type": "ObjectExpression",
          "properties": [
            {
              "type": "Property",
              "key": {
                "type": "Identifier",
                "name": "width",
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 10
                  },
                  "end": {
                    "line": 1,
                    "column": 15
                  }
                }
              },
              "kind": "set",
              "value": {
                "type": "FunctionExpression",
                "id": null,
                "params": [
                  {
                    "type": "Identifier",
                    "name": "w",
                    "loc": {
                      "start": {
                        "line": 1,
                        "column": 16
                      },
                      "end": {
                        "line": 1,
                        "column": 17
                      }
                    }
                  }
                ],
                "body": {
                  "type": "BlockStatement",
                  "body": [
                    {
                      "type": "ExpressionStatement",
                      "expression": {
                        "type": "AssignmentExpression",
                        "operator": "=",
                        "left": {
                          "type": "Identifier",
                          "name": "m_width",
                          "loc": {
                            "start": {
                              "line": 1,
                              "column": 21
                            },
                            "end": {
                              "line": 1,
                              "column": 28
                            }
                          }
                        },
                        "right": {
                          "type": "Identifier",
                          "name": "w",
                          "loc": {
                            "start": {
                              "line": 1,
                              "column": 31
                            },
                            "end": {
                              "line": 1,
                              "column": 32
                            }
                          }
                        },
                        "loc": {
                          "start": {
                            "line": 1,
                            "column": 21
                          },
                          "end": {
                            "line": 1,
                            "column": 32
                          }
                        }
                      },
                      "loc": {
                        "start": {
                          "line": 1,
                          "column": 21
                        },
                        "end": {
                          "line": 1,
                          "column": 32
                        }
                      }
                    }
                  ],
                  "loc": {
                    "start": {
                      "line": 1,
                      "column": 19
                    },
                    "end": {
                      "line": 1,
                      "column": 34
                    }
                  }
                },
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 15
                  },
                  "end": {
                    "line": 1,
                    "column": 34
                  }
                }
              }
            }
          ],
          "loc": {
            "start": {
              "line": 1,
              "column": 4
            },
            "end": {
              "line": 1,
              "column": 36
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 36
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 36
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 36
    }
  }
}
  `, ast)
}

func Test38(t *testing.T) {
	ast, err := compile("x = { set if(w) { m_if = w } }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "AssignmentExpression",
        "operator": "=",
        "left": {
          "type": "Identifier",
          "name": "x",
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 1
            }
          }
        },
        "right": {
          "type": "ObjectExpression",
          "properties": [
            {
              "type": "Property",
              "key": {
                "type": "Identifier",
                "name": "if",
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 10
                  },
                  "end": {
                    "line": 1,
                    "column": 12
                  }
                }
              },
              "kind": "set",
              "value": {
                "type": "FunctionExpression",
                "id": null,
                "params": [
                  {
                    "type": "Identifier",
                    "name": "w",
                    "loc": {
                      "start": {
                        "line": 1,
                        "column": 13
                      },
                      "end": {
                        "line": 1,
                        "column": 14
                      }
                    }
                  }
                ],
                "body": {
                  "type": "BlockStatement",
                  "body": [
                    {
                      "type": "ExpressionStatement",
                      "expression": {
                        "type": "AssignmentExpression",
                        "operator": "=",
                        "left": {
                          "type": "Identifier",
                          "name": "m_if",
                          "loc": {
                            "start": {
                              "line": 1,
                              "column": 18
                            },
                            "end": {
                              "line": 1,
                              "column": 22
                            }
                          }
                        },
                        "right": {
                          "type": "Identifier",
                          "name": "w",
                          "loc": {
                            "start": {
                              "line": 1,
                              "column": 25
                            },
                            "end": {
                              "line": 1,
                              "column": 26
                            }
                          }
                        },
                        "loc": {
                          "start": {
                            "line": 1,
                            "column": 18
                          },
                          "end": {
                            "line": 1,
                            "column": 26
                          }
                        }
                      },
                      "loc": {
                        "start": {
                          "line": 1,
                          "column": 18
                        },
                        "end": {
                          "line": 1,
                          "column": 26
                        }
                      }
                    }
                  ],
                  "loc": {
                    "start": {
                      "line": 1,
                      "column": 16
                    },
                    "end": {
                      "line": 1,
                      "column": 28
                    }
                  }
                },
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 12
                  },
                  "end": {
                    "line": 1,
                    "column": 28
                  }
                }
              }
            }
          ],
          "loc": {
            "start": {
              "line": 1,
              "column": 4
            },
            "end": {
              "line": 1,
              "column": 30
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 30
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 30
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 30
    }
  }
}
  `, ast)
}

func Test39(t *testing.T) {
	ast, err := compile("x = { set true(w) { m_true = w } }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "AssignmentExpression",
        "operator": "=",
        "left": {
          "type": "Identifier",
          "name": "x",
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 1
            }
          }
        },
        "right": {
          "type": "ObjectExpression",
          "properties": [
            {
              "type": "Property",
              "key": {
                "type": "Identifier",
                "name": "true",
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 10
                  },
                  "end": {
                    "line": 1,
                    "column": 14
                  }
                }
              },
              "kind": "set",
              "value": {
                "type": "FunctionExpression",
                "id": null,
                "params": [
                  {
                    "type": "Identifier",
                    "name": "w",
                    "loc": {
                      "start": {
                        "line": 1,
                        "column": 15
                      },
                      "end": {
                        "line": 1,
                        "column": 16
                      }
                    }
                  }
                ],
                "body": {
                  "type": "BlockStatement",
                  "body": [
                    {
                      "type": "ExpressionStatement",
                      "expression": {
                        "type": "AssignmentExpression",
                        "operator": "=",
                        "left": {
                          "type": "Identifier",
                          "name": "m_true",
                          "loc": {
                            "start": {
                              "line": 1,
                              "column": 20
                            },
                            "end": {
                              "line": 1,
                              "column": 26
                            }
                          }
                        },
                        "right": {
                          "type": "Identifier",
                          "name": "w",
                          "loc": {
                            "start": {
                              "line": 1,
                              "column": 29
                            },
                            "end": {
                              "line": 1,
                              "column": 30
                            }
                          }
                        },
                        "loc": {
                          "start": {
                            "line": 1,
                            "column": 20
                          },
                          "end": {
                            "line": 1,
                            "column": 30
                          }
                        }
                      },
                      "loc": {
                        "start": {
                          "line": 1,
                          "column": 20
                        },
                        "end": {
                          "line": 1,
                          "column": 30
                        }
                      }
                    }
                  ],
                  "loc": {
                    "start": {
                      "line": 1,
                      "column": 18
                    },
                    "end": {
                      "line": 1,
                      "column": 32
                    }
                  }
                },
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 14
                  },
                  "end": {
                    "line": 1,
                    "column": 32
                  }
                }
              }
            }
          ],
          "loc": {
            "start": {
              "line": 1,
              "column": 4
            },
            "end": {
              "line": 1,
              "column": 34
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 34
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 34
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 34
    }
  }
}
  `, ast)
}

func Test40(t *testing.T) {
	ast, err := compile("x = { set false(w) { m_false = w } }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "AssignmentExpression",
        "operator": "=",
        "left": {
          "type": "Identifier",
          "name": "x",
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 1
            }
          }
        },
        "right": {
          "type": "ObjectExpression",
          "properties": [
            {
              "type": "Property",
              "key": {
                "type": "Identifier",
                "name": "false",
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 10
                  },
                  "end": {
                    "line": 1,
                    "column": 15
                  }
                }
              },
              "kind": "set",
              "value": {
                "type": "FunctionExpression",
                "id": null,
                "params": [
                  {
                    "type": "Identifier",
                    "name": "w",
                    "loc": {
                      "start": {
                        "line": 1,
                        "column": 16
                      },
                      "end": {
                        "line": 1,
                        "column": 17
                      }
                    }
                  }
                ],
                "body": {
                  "type": "BlockStatement",
                  "body": [
                    {
                      "type": "ExpressionStatement",
                      "expression": {
                        "type": "AssignmentExpression",
                        "operator": "=",
                        "left": {
                          "type": "Identifier",
                          "name": "m_false",
                          "loc": {
                            "start": {
                              "line": 1,
                              "column": 21
                            },
                            "end": {
                              "line": 1,
                              "column": 28
                            }
                          }
                        },
                        "right": {
                          "type": "Identifier",
                          "name": "w",
                          "loc": {
                            "start": {
                              "line": 1,
                              "column": 31
                            },
                            "end": {
                              "line": 1,
                              "column": 32
                            }
                          }
                        },
                        "loc": {
                          "start": {
                            "line": 1,
                            "column": 21
                          },
                          "end": {
                            "line": 1,
                            "column": 32
                          }
                        }
                      },
                      "loc": {
                        "start": {
                          "line": 1,
                          "column": 21
                        },
                        "end": {
                          "line": 1,
                          "column": 32
                        }
                      }
                    }
                  ],
                  "loc": {
                    "start": {
                      "line": 1,
                      "column": 19
                    },
                    "end": {
                      "line": 1,
                      "column": 34
                    }
                  }
                },
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 15
                  },
                  "end": {
                    "line": 1,
                    "column": 34
                  }
                }
              }
            }
          ],
          "loc": {
            "start": {
              "line": 1,
              "column": 4
            },
            "end": {
              "line": 1,
              "column": 36
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 36
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 36
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 36
    }
  }
}
  `, ast)
}

func Test41(t *testing.T) {
	ast, err := compile("x = { set null(w) { m_null = w } }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "AssignmentExpression",
        "operator": "=",
        "left": {
          "type": "Identifier",
          "name": "x",
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 1
            }
          }
        },
        "right": {
          "type": "ObjectExpression",
          "properties": [
            {
              "type": "Property",
              "key": {
                "type": "Identifier",
                "name": "null",
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 10
                  },
                  "end": {
                    "line": 1,
                    "column": 14
                  }
                }
              },
              "kind": "set",
              "value": {
                "type": "FunctionExpression",
                "id": null,
                "params": [
                  {
                    "type": "Identifier",
                    "name": "w",
                    "loc": {
                      "start": {
                        "line": 1,
                        "column": 15
                      },
                      "end": {
                        "line": 1,
                        "column": 16
                      }
                    }
                  }
                ],
                "body": {
                  "type": "BlockStatement",
                  "body": [
                    {
                      "type": "ExpressionStatement",
                      "expression": {
                        "type": "AssignmentExpression",
                        "operator": "=",
                        "left": {
                          "type": "Identifier",
                          "name": "m_null",
                          "loc": {
                            "start": {
                              "line": 1,
                              "column": 20
                            },
                            "end": {
                              "line": 1,
                              "column": 26
                            }
                          }
                        },
                        "right": {
                          "type": "Identifier",
                          "name": "w",
                          "loc": {
                            "start": {
                              "line": 1,
                              "column": 29
                            },
                            "end": {
                              "line": 1,
                              "column": 30
                            }
                          }
                        },
                        "loc": {
                          "start": {
                            "line": 1,
                            "column": 20
                          },
                          "end": {
                            "line": 1,
                            "column": 30
                          }
                        }
                      },
                      "loc": {
                        "start": {
                          "line": 1,
                          "column": 20
                        },
                        "end": {
                          "line": 1,
                          "column": 30
                        }
                      }
                    }
                  ],
                  "loc": {
                    "start": {
                      "line": 1,
                      "column": 18
                    },
                    "end": {
                      "line": 1,
                      "column": 32
                    }
                  }
                },
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 14
                  },
                  "end": {
                    "line": 1,
                    "column": 32
                  }
                }
              }
            }
          ],
          "loc": {
            "start": {
              "line": 1,
              "column": 4
            },
            "end": {
              "line": 1,
              "column": 34
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 34
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 34
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 34
    }
  }
}
  `, ast)
}

func Test42(t *testing.T) {
	ast, err := compile("x = { set \"null\"(w) { m_null = w } }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "AssignmentExpression",
        "operator": "=",
        "left": {
          "type": "Identifier",
          "name": "x",
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 1
            }
          }
        },
        "right": {
          "type": "ObjectExpression",
          "properties": [
            {
              "type": "Property",
              "key": {
                "type": "Literal",
                "value": "null",
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 10
                  },
                  "end": {
                    "line": 1,
                    "column": 16
                  }
                }
              },
              "kind": "set",
              "value": {
                "type": "FunctionExpression",
                "id": null,
                "params": [
                  {
                    "type": "Identifier",
                    "name": "w",
                    "loc": {
                      "start": {
                        "line": 1,
                        "column": 17
                      },
                      "end": {
                        "line": 1,
                        "column": 18
                      }
                    }
                  }
                ],
                "body": {
                  "type": "BlockStatement",
                  "body": [
                    {
                      "type": "ExpressionStatement",
                      "expression": {
                        "type": "AssignmentExpression",
                        "operator": "=",
                        "left": {
                          "type": "Identifier",
                          "name": "m_null",
                          "loc": {
                            "start": {
                              "line": 1,
                              "column": 22
                            },
                            "end": {
                              "line": 1,
                              "column": 28
                            }
                          }
                        },
                        "right": {
                          "type": "Identifier",
                          "name": "w",
                          "loc": {
                            "start": {
                              "line": 1,
                              "column": 31
                            },
                            "end": {
                              "line": 1,
                              "column": 32
                            }
                          }
                        },
                        "loc": {
                          "start": {
                            "line": 1,
                            "column": 22
                          },
                          "end": {
                            "line": 1,
                            "column": 32
                          }
                        }
                      },
                      "loc": {
                        "start": {
                          "line": 1,
                          "column": 22
                        },
                        "end": {
                          "line": 1,
                          "column": 32
                        }
                      }
                    }
                  ],
                  "loc": {
                    "start": {
                      "line": 1,
                      "column": 20
                    },
                    "end": {
                      "line": 1,
                      "column": 34
                    }
                  }
                },
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 16
                  },
                  "end": {
                    "line": 1,
                    "column": 34
                  }
                }
              }
            }
          ],
          "loc": {
            "start": {
              "line": 1,
              "column": 4
            },
            "end": {
              "line": 1,
              "column": 36
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 36
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 36
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 36
    }
  }
}
  `, ast)
}

func Test43(t *testing.T) {
	ast, err := compile("x = { set 10(w) { m_null = w } }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "AssignmentExpression",
        "operator": "=",
        "left": {
          "type": "Identifier",
          "name": "x",
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 1
            }
          }
        },
        "right": {
          "type": "ObjectExpression",
          "properties": [
            {
              "type": "Property",
              "key": {
                "type": "Literal",
                "value": 10,
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 10
                  },
                  "end": {
                    "line": 1,
                    "column": 12
                  }
                }
              },
              "kind": "set",
              "value": {
                "type": "FunctionExpression",
                "id": null,
                "params": [
                  {
                    "type": "Identifier",
                    "name": "w",
                    "loc": {
                      "start": {
                        "line": 1,
                        "column": 13
                      },
                      "end": {
                        "line": 1,
                        "column": 14
                      }
                    }
                  }
                ],
                "body": {
                  "type": "BlockStatement",
                  "body": [
                    {
                      "type": "ExpressionStatement",
                      "expression": {
                        "type": "AssignmentExpression",
                        "operator": "=",
                        "left": {
                          "type": "Identifier",
                          "name": "m_null",
                          "loc": {
                            "start": {
                              "line": 1,
                              "column": 18
                            },
                            "end": {
                              "line": 1,
                              "column": 24
                            }
                          }
                        },
                        "right": {
                          "type": "Identifier",
                          "name": "w",
                          "loc": {
                            "start": {
                              "line": 1,
                              "column": 27
                            },
                            "end": {
                              "line": 1,
                              "column": 28
                            }
                          }
                        },
                        "loc": {
                          "start": {
                            "line": 1,
                            "column": 18
                          },
                          "end": {
                            "line": 1,
                            "column": 28
                          }
                        }
                      },
                      "loc": {
                        "start": {
                          "line": 1,
                          "column": 18
                        },
                        "end": {
                          "line": 1,
                          "column": 28
                        }
                      }
                    }
                  ],
                  "loc": {
                    "start": {
                      "line": 1,
                      "column": 16
                    },
                    "end": {
                      "line": 1,
                      "column": 30
                    }
                  }
                },
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 12
                  },
                  "end": {
                    "line": 1,
                    "column": 30
                  }
                }
              }
            }
          ],
          "loc": {
            "start": {
              "line": 1,
              "column": 4
            },
            "end": {
              "line": 1,
              "column": 32
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 32
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 32
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 32
    }
  }
}
  `, ast)
}

func Test44(t *testing.T) {
	ast, err := compile("x = { get: 42 }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "AssignmentExpression",
        "operator": "=",
        "left": {
          "type": "Identifier",
          "name": "x",
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 1
            }
          }
        },
        "right": {
          "type": "ObjectExpression",
          "properties": [
            {
              "type": "Property",
              "key": {
                "type": "Identifier",
                "name": "get",
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 6
                  },
                  "end": {
                    "line": 1,
                    "column": 9
                  }
                }
              },
              "value": {
                "type": "Literal",
                "value": 42,
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 11
                  },
                  "end": {
                    "line": 1,
                    "column": 13
                  }
                }
              },
              "kind": "init"
            }
          ],
          "loc": {
            "start": {
              "line": 1,
              "column": 4
            },
            "end": {
              "line": 1,
              "column": 15
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 15
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 15
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 15
    }
  }
}
  `, ast)
}

func Test45(t *testing.T) {
	ast, err := compile("x = { set: 43 }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "AssignmentExpression",
        "operator": "=",
        "left": {
          "type": "Identifier",
          "name": "x",
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 1
            }
          }
        },
        "right": {
          "type": "ObjectExpression",
          "properties": [
            {
              "type": "Property",
              "key": {
                "type": "Identifier",
                "name": "set",
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 6
                  },
                  "end": {
                    "line": 1,
                    "column": 9
                  }
                }
              },
              "value": {
                "type": "Literal",
                "value": 43,
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 11
                  },
                  "end": {
                    "line": 1,
                    "column": 13
                  }
                }
              },
              "kind": "init"
            }
          ],
          "loc": {
            "start": {
              "line": 1,
              "column": 4
            },
            "end": {
              "line": 1,
              "column": 15
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 15
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 15
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 15
    }
  }
}
  `, ast)
}

func Test46(t *testing.T) {
	ast, err := compile("/* block comment */ 42")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "Literal",
        "value": 42,
        "loc": {
          "start": {
            "line": 1,
            "column": 20
          },
          "end": {
            "line": 1,
            "column": 22
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 20
        },
        "end": {
          "line": 1,
          "column": 22
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 22
    }
  }
}
  `, ast)
}

func Test47(t *testing.T) {
	ast, err := compile("42 /*The*/ /*Answer*/")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "Literal",
        "value": 42,
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 2
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 2
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 21
    }
  }
}
  `, ast)
}

func Test48(t *testing.T) {
	ast, err := compile("/* multiline\ncomment\nshould\nbe\nignored */ 42")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "Literal",
        "value": 42,
        "loc": {
          "start": {
            "line": 5,
            "column": 11
          },
          "end": {
            "line": 5,
            "column": 13
          }
        }
      },
      "loc": {
        "start": {
          "line": 5,
          "column": 11
        },
        "end": {
          "line": 5,
          "column": 13
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 5,
      "column": 13
    }
  }
}
  `, ast)
}

func Test49(t *testing.T) {
	ast, err := compile("/*a\r\nb*/ 42")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "Literal",
        "value": 42,
        "loc": {
          "start": {
            "line": 2,
            "column": 4
          },
          "end": {
            "line": 2,
            "column": 6
          }
        }
      },
      "loc": {
        "start": {
          "line": 2,
          "column": 4
        },
        "end": {
          "line": 2,
          "column": 6
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 2,
      "column": 6
    }
  }
}
  `, ast)
}

func Test50(t *testing.T) {
	ast, err := compile("/*a\rb*/ 42")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "Literal",
        "value": 42,
        "loc": {
          "start": {
            "line": 2,
            "column": 4
          },
          "end": {
            "line": 2,
            "column": 6
          }
        }
      },
      "loc": {
        "start": {
          "line": 2,
          "column": 4
        },
        "end": {
          "line": 2,
          "column": 6
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 2,
      "column": 6
    }
  }
}
  `, ast)
}

func Test51(t *testing.T) {
	ast, err := compile("/*a\nb*/ 42")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "Literal",
        "value": 42,
        "loc": {
          "start": {
            "line": 2,
            "column": 4
          },
          "end": {
            "line": 2,
            "column": 6
          }
        }
      },
      "loc": {
        "start": {
          "line": 2,
          "column": 4
        },
        "end": {
          "line": 2,
          "column": 6
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 2,
      "column": 6
    }
  }
}
  `, ast)
}

func Test52(t *testing.T) {
	ast, err := compile("/*a\nc*/ 42")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "Literal",
        "value": 42,
        "loc": {
          "start": {
            "line": 2,
            "column": 4
          },
          "end": {
            "line": 2,
            "column": 6
          }
        }
      },
      "loc": {
        "start": {
          "line": 2,
          "column": 4
        },
        "end": {
          "line": 2,
          "column": 6
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 2,
      "column": 6
    }
  }
}
  `, ast)
}

func Test53(t *testing.T) {
	ast, err := compile("// line comment\n42")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "Literal",
        "value": 42,
        "loc": {
          "start": {
            "line": 2,
            "column": 0
          },
          "end": {
            "line": 2,
            "column": 2
          }
        }
      },
      "loc": {
        "start": {
          "line": 2,
          "column": 0
        },
        "end": {
          "line": 2,
          "column": 2
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 2,
      "column": 2
    }
  }
}
  `, ast)
}

func Test54(t *testing.T) {
	ast, err := compile("42 // line comment")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "Literal",
        "value": 42,
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 2
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 2
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 18
    }
  }
}
  `, ast)
}

func Test56(t *testing.T) {
	ast, err := compile("// Hello, world!\n42")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "Literal",
        "value": 42,
        "loc": {
          "start": {
            "line": 2,
            "column": 0
          },
          "end": {
            "line": 2,
            "column": 2
          }
        }
      },
      "loc": {
        "start": {
          "line": 2,
          "column": 0
        },
        "end": {
          "line": 2,
          "column": 2
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 2,
      "column": 2
    }
  }
}
  `, ast)
}

func Test57(t *testing.T) {
	ast, err := compile("// Hello, world!\n")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 2,
      "column": 0
    }
  }
}
  `, ast)
}

func Test58(t *testing.T) {
	ast, err := compile("//\n42")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "Literal",
        "value": 42,
        "loc": {
          "start": {
            "line": 2,
            "column": 0
          },
          "end": {
            "line": 2,
            "column": 2
          }
        }
      },
      "loc": {
        "start": {
          "line": 2,
          "column": 0
        },
        "end": {
          "line": 2,
          "column": 2
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 2,
      "column": 2
    }
  }
}
  `, ast)
}

func Test59(t *testing.T) {
	ast, err := compile("//")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 2
    }
  }
}
  `, ast)
}

func Test60(t *testing.T) {
	ast, err := compile("// ")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 3
    }
  }
}
  `, ast)
}

func Test61(t *testing.T) {
	ast, err := compile("/**/42")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "Literal",
        "value": 42,
        "loc": {
          "start": {
            "line": 1,
            "column": 4
          },
          "end": {
            "line": 1,
            "column": 6
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 4
        },
        "end": {
          "line": 1,
          "column": 6
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 6
    }
  }
}
  `, ast)
}

func Test62(t *testing.T) {
	ast, err := compile("// Hello, world!\n\n//   Another hello\n42")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "Literal",
        "value": 42,
        "loc": {
          "start": {
            "line": 4,
            "column": 0
          },
          "end": {
            "line": 4,
            "column": 2
          }
        }
      },
      "loc": {
        "start": {
          "line": 4,
          "column": 0
        },
        "end": {
          "line": 4,
          "column": 2
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 4,
      "column": 2
    }
  }
}
  `, ast)
}

func Test63(t *testing.T) {
	ast, err := compile("if (x) { // Some comment\ndoThat(); }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "IfStatement",
      "test": {
        "type": "Identifier",
        "name": "x",
        "loc": {
          "start": {
            "line": 1,
            "column": 4
          },
          "end": {
            "line": 1,
            "column": 5
          }
        }
      },
      "consequent": {
        "type": "BlockStatement",
        "body": [
          {
            "type": "ExpressionStatement",
            "expression": {
              "type": "CallExpression",
              "callee": {
                "type": "Identifier",
                "name": "doThat",
                "loc": {
                  "start": {
                    "line": 2,
                    "column": 0
                  },
                  "end": {
                    "line": 2,
                    "column": 6
                  }
                }
              },
              "arguments": [],
              "loc": {
                "start": {
                  "line": 2,
                  "column": 0
                },
                "end": {
                  "line": 2,
                  "column": 8
                }
              }
            },
            "loc": {
              "start": {
                "line": 2,
                "column": 0
              },
              "end": {
                "line": 2,
                "column": 9
              }
            }
          }
        ],
        "loc": {
          "start": {
            "line": 1,
            "column": 7
          },
          "end": {
            "line": 2,
            "column": 11
          }
        }
      },
      "alternate": null,
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 2,
          "column": 11
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 2,
      "column": 11
    }
  }
}
  `, ast)
}

func Test64(t *testing.T) {
	ast, err := compile("switch (answer) { case 42: /* perfect */ bingo() }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "SwitchStatement",
      "discriminant": {
        "type": "Identifier",
        "name": "answer",
        "loc": {
          "start": {
            "line": 1,
            "column": 8
          },
          "end": {
            "line": 1,
            "column": 14
          }
        }
      },
      "cases": [
        {
          "type": "SwitchCase",
          "consequent": [
            {
              "type": "ExpressionStatement",
              "expression": {
                "type": "CallExpression",
                "callee": {
                  "type": "Identifier",
                  "name": "bingo",
                  "loc": {
                    "start": {
                      "line": 1,
                      "column": 41
                    },
                    "end": {
                      "line": 1,
                      "column": 46
                    }
                  }
                },
                "arguments": [],
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 41
                  },
                  "end": {
                    "line": 1,
                    "column": 48
                  }
                }
              },
              "loc": {
                "start": {
                  "line": 1,
                  "column": 41
                },
                "end": {
                  "line": 1,
                  "column": 48
                }
              }
            }
          ],
          "test": {
            "type": "Literal",
            "value": 42,
            "loc": {
              "start": {
                "line": 1,
                "column": 23
              },
              "end": {
                "line": 1,
                "column": 25
              }
            }
          },
          "loc": {
            "start": {
              "line": 1,
              "column": 18
            },
            "end": {
              "line": 1,
              "column": 48
            }
          }
        }
      ],
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 50
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 50
    }
  }
}
  `, ast)
}

func Test65(t *testing.T) {
	ast, err := compile("0")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "Literal",
        "value": 0,
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 1
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 1
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 1
    }
  }
}
  `, ast)
}

func Test66(t *testing.T) {
	ast, err := compile("3")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "Literal",
        "value": 3,
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 1
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 1
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 1
    }
  }
}
  `, ast)
}

func Test67(t *testing.T) {
	ast, err := compile("5")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "Literal",
        "value": 5,
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 1
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 1
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 1
    }
  }
}
  `, ast)
}

func Test68(t *testing.T) {
	ast, err := compile("42")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "Literal",
        "value": 42,
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 2
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 2
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 2
    }
  }
}
	`, ast)
}

func Test69(t *testing.T) {
	ast, err := compile(".14")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "Literal",
        "value": 0.14,
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 3
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 3
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 3
    }
  }
}
	`, ast)
}

func Test70(t *testing.T) {
	ast, err := compile("3.14159")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "Literal",
        "value": 3.14159,
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 7
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 7
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 7
    }
  }
}
	`, ast)
}

func Test71(t *testing.T) {
	ast, err := compile("6.02214179e+23")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "Literal",
        "value": 6.02214179e+23,
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 14
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 14
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 14
    }
  }
}
	`, ast)
}

func Test72(t *testing.T) {
	ast, err := compile("1.492417830e-10")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "Literal",
        "value": 1.49241783e-10,
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 15
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 15
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 15
    }
  }
}
	`, ast)
}

func Test73(t *testing.T) {
	ast, err := compile("0x0")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "Literal",
        "value": 0,
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 3
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 3
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 3
    }
  }
}
	`, ast)
}

func Test74(t *testing.T) {
	ast, err := compile("0e+100")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "Literal",
        "value": 0,
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 6
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 6
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 6
    }
  }
}
	`, ast)
}

func Test75(t *testing.T) {
	ast, err := compile("0xabc")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "Literal",
        "value": 2748,
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 5
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 5
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 5
    }
  }
}
	`, ast)
}

func Test76(t *testing.T) {
	ast, err := compile("0xdef")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "Literal",
        "value": 3567,
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 5
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 5
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 5
    }
  }
}
	`, ast)
}

func Test77(t *testing.T) {
	ast, err := compile("0X1A")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "Literal",
        "value": 26,
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 4
          }
        }
      },
      "loc": {
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
  ],
  "loc": {
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
	`, ast)
}

func Test78(t *testing.T) {
	ast, err := compile("0x10")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "Literal",
        "value": 16,
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 4
          }
        }
      },
      "loc": {
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
  ],
  "loc": {
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
	`, ast)
}

func Test79(t *testing.T) {
	ast, err := compile("0x100")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "Literal",
        "value": 256,
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 5
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 5
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 5
    }
  }
}
	`, ast)
}

func Test80(t *testing.T) {
	ast, err := compile("0X04")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
  {
    "type": "Program",
    "body": [
      {
        "type": "ExpressionStatement",
        "expression": {
          "type": "Literal",
          "value": 4,
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 4
            }
          }
        },
        "loc": {
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
    ],
    "loc": {
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
	`, ast)
}

func Test81(t *testing.T) {
	ast, err := compile("02")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "Literal",
        "value": 2,
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 2
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 2
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 2
    }
  }
}
	`, ast)
}

func Test82(t *testing.T) {
	ast, err := compile("012")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "Literal",
        "value": 10,
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 3
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 3
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 3
    }
  }
}
	`, ast)
}

func Test83(t *testing.T) {
	ast, err := compile("0012")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "Literal",
        "value": 10,
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 4
          }
        }
      },
      "loc": {
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
  ],
  "loc": {
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
	`, ast)
}

func Test84(t *testing.T) {
	ast, err := compile("\"Hello\"")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "Literal",
        "value": "Hello",
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 7
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 7
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 7
    }
  }
}
	`, ast)
}

func Test85(t *testing.T) {
	ast, err := compile("\"\\n\\r\\t\\v\\b\\f\\\\\\'\\\"\\0\"")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "Literal",
        "value": "\n\r\t\u000b\b\f\\'\"\u0000",
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 22
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 22
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 22
    }
  }
}
	`, ast)
}

func Test86(t *testing.T) {
	ast, err := compile("\"\\u0061\"")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "Literal",
        "value": "a",
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 8
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 8
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 8
    }
  }
}
	`, ast)
}

func Test87(t *testing.T) {
	ast, err := compile("\"\\x61\"")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "Literal",
        "value": "a",
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 6
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 6
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 6
    }
  }
}
	`, ast)
}

func Test88(t *testing.T) {
	ast, err := compile("\"Hello\\nworld\"")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "Literal",
        "value": "Hello\nworld",
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 14
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 14
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 14
    }
  }
}
	`, ast)
}

func Test89(t *testing.T) {
	ast, err := compile("\"Hello\\\nworld\"")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "Literal",
        "value": "Helloworld",
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 2,
            "column": 6
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 2,
          "column": 6
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 2,
      "column": 6
    }
  }
}
	`, ast)
}

func Test90(t *testing.T) {
	ast, err := compile("\"Hello\\\u2028world\"")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "Literal",
        "value": "Helloworld",
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 14
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 14
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 14
    }
  }
}
	`, ast)
}

func Test91(t *testing.T) {
	ast, err := compile("\"Hello\\\u2029world\"")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "Literal",
        "value": "Helloworld",
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 14
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 14
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 14
    }
  }
}
	`, ast)
}

func Test92(t *testing.T) {
	ast, err := compile("\"Hello\\02World\"")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "Literal",
        "value": "Hello\u0002World",
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 15
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 15
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 15
    }
  }
}
	`, ast)
}

func Test93(t *testing.T) {
	ast, err := compile("\"Hello\\012World\"")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "Literal",
        "value": "Hello\nWorld",
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 16
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 16
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 16
    }
  }
}
	`, ast)
}

func Test94(t *testing.T) {
	ast, err := compile("\"Hello\\122World\"")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "Literal",
        "value": "HelloRWorld",
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 16
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 16
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 16
    }
  }
}
	`, ast)
}

func Test95(t *testing.T) {
	ast, err := compile("\"Hello\\0122World\"")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "Literal",
        "value": "Hello\n2World",
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 17
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 17
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 17
    }
  }
}
	`, ast)
}

func Test96(t *testing.T) {
	ast, err := compile("\"Hello\\312World\"")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "Literal",
        "value": "HelloÊWorld",
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 16
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 16
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 16
    }
  }
}
	`, ast)
}

func Test97(t *testing.T) {
	ast, err := compile("\"Hello\\412World\"")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "Literal",
        "value": "Hello!2World",
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 16
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 16
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 16
    }
  }
}
	`, ast)
}

func Test98(t *testing.T) {
	ast, err := compile("\"Hello\\812World\"")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "Literal",
        "value": "Hello812World",
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 16
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 16
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 16
    }
  }
}
	`, ast)
}

func Test99(t *testing.T) {
	ast, err := compile("\"Hello\\712World\"")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "Literal",
        "value": "Hello92World",
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 16
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 16
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 16
    }
  }
}
	`, ast)
}

func Test100(t *testing.T) {
	ast, err := compile("\"Hello\\0World\"")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "Literal",
        "value": "Hello\u0000World",
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 14
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 14
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 14
    }
  }
}
	`, ast)
}

func Test101(t *testing.T) {
	ast, err := compile("\"Hello\\\r\nworld\"")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "Literal",
        "value": "Helloworld",
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 2,
            "column": 6
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 2,
          "column": 6
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 2,
      "column": 6
    }
  }
}
	`, ast)
}

func Test102(t *testing.T) {
	ast, err := compile("\"Hello\\1World\"")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "Literal",
        "value": "Hello\u0001World",
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 14
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 14
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 14
    }
  }
}
	`, ast)
}

func Test103(t *testing.T) {
	ast, err := compile("var x = /[a-z]/i")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "VariableDeclaration",
      "declarations": [
        {
          "type": "VariableDeclarator",
          "id": {
            "type": "Identifier",
            "name": "x",
            "loc": {
              "start": {
                "line": 1,
                "column": 4
              },
              "end": {
                "line": 1,
                "column": 5
              }
            }
          },
          "init": {
            "type": "Literal",
            "value": null,
            "regexp": {
              "pattern": "[a-z]",
              "flags": "i"
            },
            "loc": {
              "start": {
                "line": 1,
                "column": 8
              },
              "end": {
                "line": 1,
                "column": 16
              }
            }
          },
          "loc": {
            "start": {
              "line": 1,
              "column": 4
            },
            "end": {
              "line": 1,
              "column": 16
            }
          }
        }
      ],
      "kind": "var",
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 16
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 16
    }
  }
}
	`, ast)
}

func Test104(t *testing.T) {
	ast, err := compile("var x = /[x-z]/i")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "VariableDeclaration",
      "declarations": [
        {
          "type": "VariableDeclarator",
          "id": {
            "type": "Identifier",
            "name": "x",
            "loc": {
              "start": {
                "line": 1,
                "column": 4
              },
              "end": {
                "line": 1,
                "column": 5
              }
            }
          },
          "init": {
            "type": "Literal",
            "value": null,
            "regexp": {
              "pattern": "[x-z]",
              "flags": "i"
            },
            "loc": {
              "start": {
                "line": 1,
                "column": 8
              },
              "end": {
                "line": 1,
                "column": 16
              }
            }
          },
          "loc": {
            "start": {
              "line": 1,
              "column": 4
            },
            "end": {
              "line": 1,
              "column": 16
            }
          }
        }
      ],
      "kind": "var",
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 16
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 16
    }
  }
}
	`, ast)
}

func Test105(t *testing.T) {
	ast, err := compile("var x = /[a-c]/i")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "VariableDeclaration",
      "declarations": [
        {
          "type": "VariableDeclarator",
          "id": {
            "type": "Identifier",
            "name": "x",
            "loc": {
              "start": {
                "line": 1,
                "column": 4
              },
              "end": {
                "line": 1,
                "column": 5
              }
            }
          },
          "init": {
            "type": "Literal",
            "value": null,
            "regexp": {
              "pattern": "[a-c]",
              "flags": "i"
            },
            "loc": {
              "start": {
                "line": 1,
                "column": 8
              },
              "end": {
                "line": 1,
                "column": 16
              }
            }
          },
          "loc": {
            "start": {
              "line": 1,
              "column": 4
            },
            "end": {
              "line": 1,
              "column": 16
            }
          }
        }
      ],
      "kind": "var",
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 16
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 16
    }
  }
}
	`, ast)
}

func Test106(t *testing.T) {
	ast, err := compile("var x = /[P QR]/i")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "VariableDeclaration",
      "declarations": [
        {
          "type": "VariableDeclarator",
          "id": {
            "type": "Identifier",
            "name": "x",
            "loc": {
              "start": {
                "line": 1,
                "column": 4
              },
              "end": {
                "line": 1,
                "column": 5
              }
            }
          },
          "init": {
            "type": "Literal",
            "value": null,
            "regexp": {
              "pattern": "[P QR]",
              "flags": "i"
            },
            "loc": {
              "start": {
                "line": 1,
                "column": 8
              },
              "end": {
                "line": 1,
                "column": 17
              }
            }
          },
          "loc": {
            "start": {
              "line": 1,
              "column": 4
            },
            "end": {
              "line": 1,
              "column": 17
            }
          }
        }
      ],
      "kind": "var",
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 17
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 17
    }
  }
}
	`, ast)
}

func Test107(t *testing.T) {
	ast, err := compile("var x = /foo\\/bar/")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "VariableDeclaration",
      "declarations": [
        {
          "type": "VariableDeclarator",
          "id": {
            "type": "Identifier",
            "name": "x",
            "loc": {
              "start": {
                "line": 1,
                "column": 4
              },
              "end": {
                "line": 1,
                "column": 5
              }
            }
          },
          "init": {
            "type": "Literal",
            "value": null,
            "regexp": {
              "pattern": "foo\\/bar",
              "flags": ""
            },
            "loc": {
              "start": {
                "line": 1,
                "column": 8
              },
              "end": {
                "line": 1,
                "column": 18
              }
            }
          },
          "loc": {
            "start": {
              "line": 1,
              "column": 4
            },
            "end": {
              "line": 1,
              "column": 18
            }
          }
        }
      ],
      "kind": "var",
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 18
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 18
    }
  }
}
	`, ast)
}

func Test108(t *testing.T) {
	ast, err := compile("var x = /=([^=\\s])+/g")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "VariableDeclaration",
      "declarations": [
        {
          "type": "VariableDeclarator",
          "id": {
            "type": "Identifier",
            "name": "x",
            "loc": {
              "start": {
                "line": 1,
                "column": 4
              },
              "end": {
                "line": 1,
                "column": 5
              }
            }
          },
          "init": {
            "type": "Literal",
            "value": null,
            "regexp": {
              "pattern": "=([^=\\s])+",
              "flags": "g"
            },
            "loc": {
              "start": {
                "line": 1,
                "column": 8
              },
              "end": {
                "line": 1,
                "column": 21
              }
            }
          },
          "loc": {
            "start": {
              "line": 1,
              "column": 4
            },
            "end": {
              "line": 1,
              "column": 21
            }
          }
        }
      ],
      "kind": "var",
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 21
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 21
    }
  }
}
	`, ast)
}

func Test109(t *testing.T) {
	ast, err := compile("new Button")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "NewExpression",
        "callee": {
          "type": "Identifier",
          "name": "Button",
          "loc": {
            "start": {
              "line": 1,
              "column": 4
            },
            "end": {
              "line": 1,
              "column": 10
            }
          }
        },
        "arguments": [],
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 10
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 10
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 10
    }
  }
}
	`, ast)
}

func Test110(t *testing.T) {
	ast, err := compile("new Button()")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "NewExpression",
        "callee": {
          "type": "Identifier",
          "name": "Button",
          "loc": {
            "start": {
              "line": 1,
              "column": 4
            },
            "end": {
              "line": 1,
              "column": 10
            }
          }
        },
        "arguments": [],
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 12
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 12
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 12
    }
  }
}
	`, ast)
}

func Test111(t *testing.T) {
	ast, err := compile("new new foo")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "NewExpression",
        "callee": {
          "type": "NewExpression",
          "callee": {
            "type": "Identifier",
            "name": "foo",
            "loc": {
              "start": {
                "line": 1,
                "column": 8
              },
              "end": {
                "line": 1,
                "column": 11
              }
            }
          },
          "arguments": [],
          "loc": {
            "start": {
              "line": 1,
              "column": 4
            },
            "end": {
              "line": 1,
              "column": 11
            }
          }
        },
        "arguments": [],
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 11
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 11
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 11
    }
  }
}
	`, ast)
}

func Test112(t *testing.T) {
	ast, err := compile("new new foo()")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "NewExpression",
        "callee": {
          "type": "NewExpression",
          "callee": {
            "type": "Identifier",
            "name": "foo",
            "loc": {
              "start": {
                "line": 1,
                "column": 8
              },
              "end": {
                "line": 1,
                "column": 11
              }
            }
          },
          "arguments": [],
          "loc": {
            "start": {
              "line": 1,
              "column": 4
            },
            "end": {
              "line": 1,
              "column": 13
            }
          }
        },
        "arguments": [],
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 13
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 13
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 13
    }
  }
}
	`, ast)
}

func Test113(t *testing.T) {
	ast, err := compile("new foo().bar()")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "CallExpression",
        "callee": {
          "type": "MemberExpression",
          "object": {
            "type": "NewExpression",
            "callee": {
              "type": "Identifier",
              "name": "foo",
              "loc": {
                "start": {
                  "line": 1,
                  "column": 4
                },
                "end": {
                  "line": 1,
                  "column": 7
                }
              }
            },
            "arguments": [],
            "loc": {
              "start": {
                "line": 1,
                "column": 0
              },
              "end": {
                "line": 1,
                "column": 9
              }
            }
          },
          "property": {
            "type": "Identifier",
            "name": "bar",
            "loc": {
              "start": {
                "line": 1,
                "column": 10
              },
              "end": {
                "line": 1,
                "column": 13
              }
            }
          },
          "computed": false,
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 13
            }
          }
        },
        "arguments": [],
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 15
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 15
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 15
    }
  }
}
	`, ast)
}

func Test114(t *testing.T) {
	ast, err := compile("new foo[bar]")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "NewExpression",
        "callee": {
          "type": "MemberExpression",
          "object": {
            "type": "Identifier",
            "name": "foo",
            "loc": {
              "start": {
                "line": 1,
                "column": 4
              },
              "end": {
                "line": 1,
                "column": 7
              }
            }
          },
          "property": {
            "type": "Identifier",
            "name": "bar",
            "loc": {
              "start": {
                "line": 1,
                "column": 8
              },
              "end": {
                "line": 1,
                "column": 11
              }
            }
          },
          "computed": true,
          "loc": {
            "start": {
              "line": 1,
              "column": 4
            },
            "end": {
              "line": 1,
              "column": 12
            }
          }
        },
        "arguments": [],
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 12
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 12
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 12
    }
  }
}
	`, ast)
}

func Test115(t *testing.T) {
	ast, err := compile("new foo.bar()")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "NewExpression",
        "callee": {
          "type": "MemberExpression",
          "object": {
            "type": "Identifier",
            "name": "foo",
            "loc": {
              "start": {
                "line": 1,
                "column": 4
              },
              "end": {
                "line": 1,
                "column": 7
              }
            }
          },
          "property": {
            "type": "Identifier",
            "name": "bar",
            "loc": {
              "start": {
                "line": 1,
                "column": 8
              },
              "end": {
                "line": 1,
                "column": 11
              }
            }
          },
          "computed": false,
          "loc": {
            "start": {
              "line": 1,
              "column": 4
            },
            "end": {
              "line": 1,
              "column": 11
            }
          }
        },
        "arguments": [],
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 13
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 13
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 13
    }
  }
}
	`, ast)
}

func Test116(t *testing.T) {
	ast, err := compile("( new foo).bar()")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "CallExpression",
        "callee": {
          "type": "MemberExpression",
          "object": {
            "type": "NewExpression",
            "callee": {
              "type": "Identifier",
              "name": "foo",
              "loc": {
                "start": {
                  "line": 1,
                  "column": 6
                },
                "end": {
                  "line": 1,
                  "column": 9
                }
              }
            },
            "arguments": [],
            "loc": {
              "start": {
                "line": 1,
                "column": 2
              },
              "end": {
                "line": 1,
                "column": 9
              }
            }
          },
          "property": {
            "type": "Identifier",
            "name": "bar",
            "loc": {
              "start": {
                "line": 1,
                "column": 11
              },
              "end": {
                "line": 1,
                "column": 14
              }
            }
          },
          "computed": false,
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 14
            }
          }
        },
        "arguments": [],
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 16
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 16
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 16
    }
  }
}
	`, ast)
}

func Test117(t *testing.T) {
	ast, err := compile("foo(bar, baz)")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "CallExpression",
        "callee": {
          "type": "Identifier",
          "name": "foo",
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 3
            }
          }
        },
        "arguments": [
          {
            "type": "Identifier",
            "name": "bar",
            "loc": {
              "start": {
                "line": 1,
                "column": 4
              },
              "end": {
                "line": 1,
                "column": 7
              }
            }
          },
          {
            "type": "Identifier",
            "name": "baz",
            "loc": {
              "start": {
                "line": 1,
                "column": 9
              },
              "end": {
                "line": 1,
                "column": 12
              }
            }
          }
        ],
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 13
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 13
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 13
    }
  }
}
	`, ast)
}

func Test118(t *testing.T) {
	ast, err := compile("foo(bar, baz)")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "CallExpression",
        "callee": {
          "type": "Identifier",
          "name": "foo",
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 3
            }
          }
        },
        "arguments": [
          {
            "type": "Identifier",
            "name": "bar",
            "loc": {
              "start": {
                "line": 1,
                "column": 4
              },
              "end": {
                "line": 1,
                "column": 7
              }
            }
          },
          {
            "type": "Identifier",
            "name": "baz",
            "loc": {
              "start": {
                "line": 1,
                "column": 9
              },
              "end": {
                "line": 1,
                "column": 12
              }
            }
          }
        ],
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 13
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 13
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 13
    }
  }
}
	`, ast)
}

func Test119(t *testing.T) {
	ast, err := compile("(    foo  )()")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "CallExpression",
        "callee": {
          "type": "Identifier",
          "name": "foo",
          "loc": {
            "start": {
              "line": 1,
              "column": 5
            },
            "end": {
              "line": 1,
              "column": 8
            }
          }
        },
        "arguments": [],
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 13
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 13
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 13
    }
  }
}
	`, ast)
}

func Test120(t *testing.T) {
	ast, err := compile("universe.milkyway")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "MemberExpression",
        "object": {
          "type": "Identifier",
          "name": "universe",
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 8
            }
          }
        },
        "property": {
          "type": "Identifier",
          "name": "milkyway",
          "loc": {
            "start": {
              "line": 1,
              "column": 9
            },
            "end": {
              "line": 1,
              "column": 17
            }
          }
        },
        "computed": false,
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 17
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 17
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 17
    }
  }
}
	`, ast)
}

func Test121(t *testing.T) {
	ast, err := compile("universe.milkyway.solarsystem")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "MemberExpression",
        "object": {
          "type": "MemberExpression",
          "object": {
            "type": "Identifier",
            "name": "universe",
            "loc": {
              "start": {
                "line": 1,
                "column": 0
              },
              "end": {
                "line": 1,
                "column": 8
              }
            }
          },
          "property": {
            "type": "Identifier",
            "name": "milkyway",
            "loc": {
              "start": {
                "line": 1,
                "column": 9
              },
              "end": {
                "line": 1,
                "column": 17
              }
            }
          },
          "computed": false,
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 17
            }
          }
        },
        "property": {
          "type": "Identifier",
          "name": "solarsystem",
          "loc": {
            "start": {
              "line": 1,
              "column": 18
            },
            "end": {
              "line": 1,
              "column": 29
            }
          }
        },
        "computed": false,
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 29
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 29
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 29
    }
  }
}
	`, ast)
}

func Test122(t *testing.T) {
	ast, err := compile("universe.milkyway.solarsystem.Earth")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "MemberExpression",
        "object": {
          "type": "MemberExpression",
          "object": {
            "type": "MemberExpression",
            "object": {
              "type": "Identifier",
              "name": "universe",
              "loc": {
                "start": {
                  "line": 1,
                  "column": 0
                },
                "end": {
                  "line": 1,
                  "column": 8
                }
              }
            },
            "property": {
              "type": "Identifier",
              "name": "milkyway",
              "loc": {
                "start": {
                  "line": 1,
                  "column": 9
                },
                "end": {
                  "line": 1,
                  "column": 17
                }
              }
            },
            "computed": false,
            "loc": {
              "start": {
                "line": 1,
                "column": 0
              },
              "end": {
                "line": 1,
                "column": 17
              }
            }
          },
          "property": {
            "type": "Identifier",
            "name": "solarsystem",
            "loc": {
              "start": {
                "line": 1,
                "column": 18
              },
              "end": {
                "line": 1,
                "column": 29
              }
            }
          },
          "computed": false,
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 29
            }
          }
        },
        "property": {
          "type": "Identifier",
          "name": "Earth",
          "loc": {
            "start": {
              "line": 1,
              "column": 30
            },
            "end": {
              "line": 1,
              "column": 35
            }
          }
        },
        "computed": false,
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 35
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 35
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 35
    }
  }
}
	`, ast)
}

func Test123(t *testing.T) {
	ast, err := compile("universe[galaxyName, otherUselessName]")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "MemberExpression",
        "object": {
          "type": "Identifier",
          "name": "universe",
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 8
            }
          }
        },
        "property": {
          "type": "SequenceExpression",
          "expressions": [
            {
              "type": "Identifier",
              "name": "galaxyName",
              "loc": {
                "start": {
                  "line": 1,
                  "column": 9
                },
                "end": {
                  "line": 1,
                  "column": 19
                }
              }
            },
            {
              "type": "Identifier",
              "name": "otherUselessName",
              "loc": {
                "start": {
                  "line": 1,
                  "column": 21
                },
                "end": {
                  "line": 1,
                  "column": 37
                }
              }
            }
          ],
          "loc": {
            "start": {
              "line": 1,
              "column": 9
            },
            "end": {
              "line": 1,
              "column": 37
            }
          }
        },
        "computed": true,
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 38
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 38
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 38
    }
  }
}
	`, ast)
}

func Test124(t *testing.T) {
	ast, err := compile("universe[galaxyName]")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "MemberExpression",
        "object": {
          "type": "Identifier",
          "name": "universe",
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 8
            }
          }
        },
        "property": {
          "type": "Identifier",
          "name": "galaxyName",
          "loc": {
            "start": {
              "line": 1,
              "column": 9
            },
            "end": {
              "line": 1,
              "column": 19
            }
          }
        },
        "computed": true,
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 20
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 20
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 20
    }
  }
}
	`, ast)
}

func Test125(t *testing.T) {
	ast, err := compile("universe[42].galaxies")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "MemberExpression",
        "object": {
          "type": "MemberExpression",
          "object": {
            "type": "Identifier",
            "name": "universe",
            "loc": {
              "start": {
                "line": 1,
                "column": 0
              },
              "end": {
                "line": 1,
                "column": 8
              }
            }
          },
          "property": {
            "type": "Literal",
            "value": 42,
            "loc": {
              "start": {
                "line": 1,
                "column": 9
              },
              "end": {
                "line": 1,
                "column": 11
              }
            }
          },
          "computed": true,
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 12
            }
          }
        },
        "property": {
          "type": "Identifier",
          "name": "galaxies",
          "loc": {
            "start": {
              "line": 1,
              "column": 13
            },
            "end": {
              "line": 1,
              "column": 21
            }
          }
        },
        "computed": false,
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 21
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 21
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 21
    }
  }
}
	`, ast)
}

func Test126(t *testing.T) {
	ast, err := compile("universe(42).galaxies")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "MemberExpression",
        "object": {
          "type": "CallExpression",
          "callee": {
            "type": "Identifier",
            "name": "universe",
            "loc": {
              "start": {
                "line": 1,
                "column": 0
              },
              "end": {
                "line": 1,
                "column": 8
              }
            }
          },
          "arguments": [
            {
              "type": "Literal",
              "value": 42,
              "loc": {
                "start": {
                  "line": 1,
                  "column": 9
                },
                "end": {
                  "line": 1,
                  "column": 11
                }
              }
            }
          ],
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 12
            }
          }
        },
        "property": {
          "type": "Identifier",
          "name": "galaxies",
          "loc": {
            "start": {
              "line": 1,
              "column": 13
            },
            "end": {
              "line": 1,
              "column": 21
            }
          }
        },
        "computed": false,
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 21
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 21
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 21
    }
  }
}
	`, ast)
}

func Test127(t *testing.T) {
	ast, err := compile("universe(42).galaxies(14, 3, 77).milkyway")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "MemberExpression",
        "object": {
          "type": "CallExpression",
          "callee": {
            "type": "MemberExpression",
            "object": {
              "type": "CallExpression",
              "callee": {
                "type": "Identifier",
                "name": "universe",
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 0
                  },
                  "end": {
                    "line": 1,
                    "column": 8
                  }
                }
              },
              "arguments": [
                {
                  "type": "Literal",
                  "value": 42,
                  "loc": {
                    "start": {
                      "line": 1,
                      "column": 9
                    },
                    "end": {
                      "line": 1,
                      "column": 11
                    }
                  }
                }
              ],
              "loc": {
                "start": {
                  "line": 1,
                  "column": 0
                },
                "end": {
                  "line": 1,
                  "column": 12
                }
              }
            },
            "property": {
              "type": "Identifier",
              "name": "galaxies",
              "loc": {
                "start": {
                  "line": 1,
                  "column": 13
                },
                "end": {
                  "line": 1,
                  "column": 21
                }
              }
            },
            "computed": false,
            "loc": {
              "start": {
                "line": 1,
                "column": 0
              },
              "end": {
                "line": 1,
                "column": 21
              }
            }
          },
          "arguments": [
            {
              "type": "Literal",
              "value": 14,
              "loc": {
                "start": {
                  "line": 1,
                  "column": 22
                },
                "end": {
                  "line": 1,
                  "column": 24
                }
              }
            },
            {
              "type": "Literal",
              "value": 3,
              "loc": {
                "start": {
                  "line": 1,
                  "column": 26
                },
                "end": {
                  "line": 1,
                  "column": 27
                }
              }
            },
            {
              "type": "Literal",
              "value": 77,
              "loc": {
                "start": {
                  "line": 1,
                  "column": 29
                },
                "end": {
                  "line": 1,
                  "column": 31
                }
              }
            }
          ],
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 32
            }
          }
        },
        "property": {
          "type": "Identifier",
          "name": "milkyway",
          "loc": {
            "start": {
              "line": 1,
              "column": 33
            },
            "end": {
              "line": 1,
              "column": 41
            }
          }
        },
        "computed": false,
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 41
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 41
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 41
    }
  }
}
	`, ast)
}

func Test128(t *testing.T) {
	ast, err := compile("earth.asia.Indonesia.prepareForElection(2014)")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "CallExpression",
        "callee": {
          "type": "MemberExpression",
          "object": {
            "type": "MemberExpression",
            "object": {
              "type": "MemberExpression",
              "object": {
                "type": "Identifier",
                "name": "earth",
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 0
                  },
                  "end": {
                    "line": 1,
                    "column": 5
                  }
                }
              },
              "property": {
                "type": "Identifier",
                "name": "asia",
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 6
                  },
                  "end": {
                    "line": 1,
                    "column": 10
                  }
                }
              },
              "computed": false,
              "loc": {
                "start": {
                  "line": 1,
                  "column": 0
                },
                "end": {
                  "line": 1,
                  "column": 10
                }
              }
            },
            "property": {
              "type": "Identifier",
              "name": "Indonesia",
              "loc": {
                "start": {
                  "line": 1,
                  "column": 11
                },
                "end": {
                  "line": 1,
                  "column": 20
                }
              }
            },
            "computed": false,
            "loc": {
              "start": {
                "line": 1,
                "column": 0
              },
              "end": {
                "line": 1,
                "column": 20
              }
            }
          },
          "property": {
            "type": "Identifier",
            "name": "prepareForElection",
            "loc": {
              "start": {
                "line": 1,
                "column": 21
              },
              "end": {
                "line": 1,
                "column": 39
              }
            }
          },
          "computed": false,
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 39
            }
          }
        },
        "arguments": [
          {
            "type": "Literal",
            "value": 2014,
            "loc": {
              "start": {
                "line": 1,
                "column": 40
              },
              "end": {
                "line": 1,
                "column": 44
              }
            }
          }
        ],
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 45
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 45
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 45
    }
  }
}
	`, ast)
}

func Test129(t *testing.T) {
	ast, err := compile("universe.if")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "MemberExpression",
        "object": {
          "type": "Identifier",
          "name": "universe",
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 8
            }
          }
        },
        "property": {
          "type": "Identifier",
          "name": "if",
          "loc": {
            "start": {
              "line": 1,
              "column": 9
            },
            "end": {
              "line": 1,
              "column": 11
            }
          }
        },
        "computed": false,
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 11
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 11
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 11
    }
  }
}
	`, ast)
}

func Test130(t *testing.T) {
	ast, err := compile("universe.true")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "MemberExpression",
        "object": {
          "type": "Identifier",
          "name": "universe",
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 8
            }
          }
        },
        "property": {
          "type": "Identifier",
          "name": "true",
          "loc": {
            "start": {
              "line": 1,
              "column": 9
            },
            "end": {
              "line": 1,
              "column": 13
            }
          }
        },
        "computed": false,
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 13
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 13
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 13
    }
  }
}
	`, ast)
}

func Test131(t *testing.T) {
	ast, err := compile("universe.false")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "MemberExpression",
        "object": {
          "type": "Identifier",
          "name": "universe",
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 8
            }
          }
        },
        "property": {
          "type": "Identifier",
          "name": "false",
          "loc": {
            "start": {
              "line": 1,
              "column": 9
            },
            "end": {
              "line": 1,
              "column": 14
            }
          }
        },
        "computed": false,
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 14
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 14
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 14
    }
  }
}
	`, ast)
}

func Test132(t *testing.T) {
	ast, err := compile("universe.null")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "MemberExpression",
        "object": {
          "type": "Identifier",
          "name": "universe",
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 8
            }
          }
        },
        "property": {
          "type": "Identifier",
          "name": "null",
          "loc": {
            "start": {
              "line": 1,
              "column": 9
            },
            "end": {
              "line": 1,
              "column": 13
            }
          }
        },
        "computed": false,
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 13
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 13
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 13
    }
  }
}
	`, ast)
}

func Test133(t *testing.T) {
	ast, err := compile("x++")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "UpdateExpression",
        "operator": "++",
        "prefix": false,
        "argument": {
          "type": "Identifier",
          "name": "x",
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 1
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 3
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 3
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 3
    }
  }
}
	`, ast)
}

func Test134(t *testing.T) {
	ast, err := compile("x--")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "UpdateExpression",
        "operator": "--",
        "prefix": false,
        "argument": {
          "type": "Identifier",
          "name": "x",
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 1
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 3
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 3
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 3
    }
  }
}
	`, ast)
}

func Test135(t *testing.T) {
	ast, err := compile("eval++")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "UpdateExpression",
        "operator": "++",
        "prefix": false,
        "argument": {
          "type": "Identifier",
          "name": "eval",
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 4
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 6
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 6
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 6
    }
  }
}
	`, ast)
}

func Test136(t *testing.T) {
	ast, err := compile("eval--")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "UpdateExpression",
        "operator": "--",
        "prefix": false,
        "argument": {
          "type": "Identifier",
          "name": "eval",
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 4
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 6
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 6
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 6
    }
  }
}
	`, ast)
}

func Test137(t *testing.T) {
	ast, err := compile("arguments++")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "UpdateExpression",
        "operator": "++",
        "prefix": false,
        "argument": {
          "type": "Identifier",
          "name": "arguments",
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 9
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 11
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 11
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 11
    }
  }
}
	`, ast)
}

func Test138(t *testing.T) {
	ast, err := compile("arguments--")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "UpdateExpression",
        "operator": "--",
        "prefix": false,
        "argument": {
          "type": "Identifier",
          "name": "arguments",
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 9
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 11
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 11
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 11
    }
  }
}
	`, ast)
}

func Test139(t *testing.T) {
	ast, err := compile("++x")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "UpdateExpression",
        "operator": "++",
        "prefix": true,
        "argument": {
          "type": "Identifier",
          "name": "x",
          "loc": {
            "start": {
              "line": 1,
              "column": 2
            },
            "end": {
              "line": 1,
              "column": 3
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 3
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 3
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 3
    }
  }
}
	`, ast)
}

func Test140(t *testing.T) {
	ast, err := compile("--x")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "UpdateExpression",
        "operator": "--",
        "prefix": true,
        "argument": {
          "type": "Identifier",
          "name": "x",
          "loc": {
            "start": {
              "line": 1,
              "column": 2
            },
            "end": {
              "line": 1,
              "column": 3
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 3
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 3
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 3
    }
  }
}
	`, ast)
}

func Test141(t *testing.T) {
	ast, err := compile("++eval")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "UpdateExpression",
        "operator": "++",
        "prefix": true,
        "argument": {
          "type": "Identifier",
          "name": "eval",
          "loc": {
            "start": {
              "line": 1,
              "column": 2
            },
            "end": {
              "line": 1,
              "column": 6
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 6
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 6
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 6
    }
  }
}
	`, ast)
}

func Test142(t *testing.T) {
	ast, err := compile("--eval")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "UpdateExpression",
        "operator": "--",
        "prefix": true,
        "argument": {
          "type": "Identifier",
          "name": "eval",
          "loc": {
            "start": {
              "line": 1,
              "column": 2
            },
            "end": {
              "line": 1,
              "column": 6
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 6
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 6
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 6
    }
  }
}
	`, ast)
}

func Test143(t *testing.T) {
	ast, err := compile("++arguments")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "UpdateExpression",
        "operator": "++",
        "prefix": true,
        "argument": {
          "type": "Identifier",
          "name": "arguments",
          "loc": {
            "start": {
              "line": 1,
              "column": 2
            },
            "end": {
              "line": 1,
              "column": 11
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 11
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 11
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 11
    }
  }
}
	`, ast)
}

func Test144(t *testing.T) {
	ast, err := compile("--arguments")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "UpdateExpression",
        "operator": "--",
        "prefix": true,
        "argument": {
          "type": "Identifier",
          "name": "arguments",
          "loc": {
            "start": {
              "line": 1,
              "column": 2
            },
            "end": {
              "line": 1,
              "column": 11
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 11
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 11
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 11
    }
  }
}
	`, ast)
}

func Test145(t *testing.T) {
	ast, err := compile("+x")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "UnaryExpression",
        "operator": "+",
        "prefix": true,
        "argument": {
          "type": "Identifier",
          "name": "x",
          "loc": {
            "start": {
              "line": 1,
              "column": 1
            },
            "end": {
              "line": 1,
              "column": 2
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 2
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 2
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 2
    }
  }
}
	`, ast)
}

func Test146(t *testing.T) {
	ast, err := compile("-x")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "UnaryExpression",
        "operator": "-",
        "prefix": true,
        "argument": {
          "type": "Identifier",
          "name": "x",
          "loc": {
            "start": {
              "line": 1,
              "column": 1
            },
            "end": {
              "line": 1,
              "column": 2
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 2
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 2
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 2
    }
  }
}
	`, ast)
}

func Test147(t *testing.T) {
	ast, err := compile("~x")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "UnaryExpression",
        "operator": "~",
        "prefix": true,
        "argument": {
          "type": "Identifier",
          "name": "x",
          "loc": {
            "start": {
              "line": 1,
              "column": 1
            },
            "end": {
              "line": 1,
              "column": 2
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 2
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 2
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 2
    }
  }
}
	`, ast)
}

func Test148(t *testing.T) {
	ast, err := compile("!x")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "UnaryExpression",
        "operator": "!",
        "prefix": true,
        "argument": {
          "type": "Identifier",
          "name": "x",
          "loc": {
            "start": {
              "line": 1,
              "column": 1
            },
            "end": {
              "line": 1,
              "column": 2
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 2
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 2
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 2
    }
  }
}
	`, ast)
}

func Test149(t *testing.T) {
	ast, err := compile("void x")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "UnaryExpression",
        "operator": "void",
        "prefix": true,
        "argument": {
          "type": "Identifier",
          "name": "x",
          "loc": {
            "start": {
              "line": 1,
              "column": 5
            },
            "end": {
              "line": 1,
              "column": 6
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 6
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 6
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 6
    }
  }
}
	`, ast)
}

func Test150(t *testing.T) {
	ast, err := compile("delete x")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "UnaryExpression",
        "operator": "delete",
        "prefix": true,
        "argument": {
          "type": "Identifier",
          "name": "x",
          "loc": {
            "start": {
              "line": 1,
              "column": 7
            },
            "end": {
              "line": 1,
              "column": 8
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 8
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 8
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 8
    }
  }
}
	`, ast)
}

func Test151(t *testing.T) {
	ast, err := compile("typeof x")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "UnaryExpression",
        "operator": "typeof",
        "prefix": true,
        "argument": {
          "type": "Identifier",
          "name": "x",
          "loc": {
            "start": {
              "line": 1,
              "column": 7
            },
            "end": {
              "line": 1,
              "column": 8
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 8
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 8
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 8
    }
  }
}
	`, ast)
}

func Test152(t *testing.T) {
	ast, err := compile("x * y")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "BinaryExpression",
        "left": {
          "type": "Identifier",
          "name": "x",
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 1
            }
          }
        },
        "operator": "*",
        "right": {
          "type": "Identifier",
          "name": "y",
          "loc": {
            "start": {
              "line": 1,
              "column": 4
            },
            "end": {
              "line": 1,
              "column": 5
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 5
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 5
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 5
    }
  }
}
	`, ast)
}

func Test153(t *testing.T) {
	ast, err := compile("x / y")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "BinaryExpression",
        "left": {
          "type": "Identifier",
          "name": "x",
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 1
            }
          }
        },
        "operator": "/",
        "right": {
          "type": "Identifier",
          "name": "y",
          "loc": {
            "start": {
              "line": 1,
              "column": 4
            },
            "end": {
              "line": 1,
              "column": 5
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 5
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 5
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 5
    }
  }
}
	`, ast)
}

func Test154(t *testing.T) {
	ast, err := compile("x % y")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "BinaryExpression",
        "left": {
          "type": "Identifier",
          "name": "x",
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 1
            }
          }
        },
        "operator": "%",
        "right": {
          "type": "Identifier",
          "name": "y",
          "loc": {
            "start": {
              "line": 1,
              "column": 4
            },
            "end": {
              "line": 1,
              "column": 5
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 5
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 5
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 5
    }
  }
}
	`, ast)
}

func Test155(t *testing.T) {
	ast, err := compile("x + y")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "BinaryExpression",
        "left": {
          "type": "Identifier",
          "name": "x",
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 1
            }
          }
        },
        "operator": "+",
        "right": {
          "type": "Identifier",
          "name": "y",
          "loc": {
            "start": {
              "line": 1,
              "column": 4
            },
            "end": {
              "line": 1,
              "column": 5
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 5
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 5
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 5
    }
  }
}
	`, ast)
}

func Test156(t *testing.T) {
	ast, err := compile("x - y")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "BinaryExpression",
        "left": {
          "type": "Identifier",
          "name": "x",
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 1
            }
          }
        },
        "operator": "-",
        "right": {
          "type": "Identifier",
          "name": "y",
          "loc": {
            "start": {
              "line": 1,
              "column": 4
            },
            "end": {
              "line": 1,
              "column": 5
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 5
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 5
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 5
    }
  }
}
	`, ast)
}

func Test157(t *testing.T) {
	ast, err := compile("x << y")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "BinaryExpression",
        "left": {
          "type": "Identifier",
          "name": "x",
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 1
            }
          }
        },
        "operator": "<<",
        "right": {
          "type": "Identifier",
          "name": "y",
          "loc": {
            "start": {
              "line": 1,
              "column": 5
            },
            "end": {
              "line": 1,
              "column": 6
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 6
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 6
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 6
    }
  }
}
	`, ast)
}

func Test158(t *testing.T) {
	ast, err := compile("x >> y")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "BinaryExpression",
        "left": {
          "type": "Identifier",
          "name": "x",
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 1
            }
          }
        },
        "operator": ">>",
        "right": {
          "type": "Identifier",
          "name": "y",
          "loc": {
            "start": {
              "line": 1,
              "column": 5
            },
            "end": {
              "line": 1,
              "column": 6
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 6
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 6
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 6
    }
  }
}
	`, ast)
}

func Test159(t *testing.T) {
	ast, err := compile("x >>> y")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "BinaryExpression",
        "left": {
          "type": "Identifier",
          "name": "x",
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 1
            }
          }
        },
        "operator": ">>>",
        "right": {
          "type": "Identifier",
          "name": "y",
          "loc": {
            "start": {
              "line": 1,
              "column": 6
            },
            "end": {
              "line": 1,
              "column": 7
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 7
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 7
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 7
    }
  }
}
	`, ast)
}

func Test160(t *testing.T) {
	ast, err := compile("x < y")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "BinaryExpression",
        "left": {
          "type": "Identifier",
          "name": "x",
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 1
            }
          }
        },
        "operator": "<",
        "right": {
          "type": "Identifier",
          "name": "y",
          "loc": {
            "start": {
              "line": 1,
              "column": 4
            },
            "end": {
              "line": 1,
              "column": 5
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 5
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 5
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 5
    }
  }
}
	`, ast)
}

func Test161(t *testing.T) {
	ast, err := compile("x > y")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "BinaryExpression",
        "left": {
          "type": "Identifier",
          "name": "x",
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 1
            }
          }
        },
        "operator": ">",
        "right": {
          "type": "Identifier",
          "name": "y",
          "loc": {
            "start": {
              "line": 1,
              "column": 4
            },
            "end": {
              "line": 1,
              "column": 5
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 5
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 5
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 5
    }
  }
}
	`, ast)
}

func Test162(t *testing.T) {
	ast, err := compile("x <= y")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "BinaryExpression",
        "left": {
          "type": "Identifier",
          "name": "x",
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 1
            }
          }
        },
        "operator": "<=",
        "right": {
          "type": "Identifier",
          "name": "y",
          "loc": {
            "start": {
              "line": 1,
              "column": 5
            },
            "end": {
              "line": 1,
              "column": 6
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 6
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 6
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 6
    }
  }
}
	`, ast)
}

func Test163(t *testing.T) {
	ast, err := compile("x >= y")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "BinaryExpression",
        "left": {
          "type": "Identifier",
          "name": "x",
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 1
            }
          }
        },
        "operator": ">=",
        "right": {
          "type": "Identifier",
          "name": "y",
          "loc": {
            "start": {
              "line": 1,
              "column": 5
            },
            "end": {
              "line": 1,
              "column": 6
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 6
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 6
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 6
    }
  }
}
	`, ast)
}

func Test164(t *testing.T) {
	ast, err := compile("x in y")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "BinaryExpression",
        "left": {
          "type": "Identifier",
          "name": "x",
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 1
            }
          }
        },
        "operator": "in",
        "right": {
          "type": "Identifier",
          "name": "y",
          "loc": {
            "start": {
              "line": 1,
              "column": 5
            },
            "end": {
              "line": 1,
              "column": 6
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 6
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 6
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 6
    }
  }
}
	`, ast)
}

func Test165(t *testing.T) {
	ast, err := compile("x instanceof y")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "BinaryExpression",
        "left": {
          "type": "Identifier",
          "name": "x",
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 1
            }
          }
        },
        "operator": "instanceof",
        "right": {
          "type": "Identifier",
          "name": "y",
          "loc": {
            "start": {
              "line": 1,
              "column": 13
            },
            "end": {
              "line": 1,
              "column": 14
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 14
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 14
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 14
    }
  }
}
	`, ast)
}

func Test166(t *testing.T) {
	ast, err := compile("x < y < z")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "BinaryExpression",
        "left": {
          "type": "BinaryExpression",
          "left": {
            "type": "Identifier",
            "name": "x",
            "loc": {
              "start": {
                "line": 1,
                "column": 0
              },
              "end": {
                "line": 1,
                "column": 1
              }
            }
          },
          "operator": "<",
          "right": {
            "type": "Identifier",
            "name": "y",
            "loc": {
              "start": {
                "line": 1,
                "column": 4
              },
              "end": {
                "line": 1,
                "column": 5
              }
            }
          },
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 5
            }
          }
        },
        "operator": "<",
        "right": {
          "type": "Identifier",
          "name": "z",
          "loc": {
            "start": {
              "line": 1,
              "column": 8
            },
            "end": {
              "line": 1,
              "column": 9
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 9
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 9
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 9
    }
  }
}
	`, ast)
}

func Test167(t *testing.T) {
	ast, err := compile("x == y")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "BinaryExpression",
        "left": {
          "type": "Identifier",
          "name": "x",
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 1
            }
          }
        },
        "operator": "==",
        "right": {
          "type": "Identifier",
          "name": "y",
          "loc": {
            "start": {
              "line": 1,
              "column": 5
            },
            "end": {
              "line": 1,
              "column": 6
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 6
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 6
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 6
    }
  }
}
	`, ast)
}

func Test168(t *testing.T) {
	ast, err := compile("x != y")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "BinaryExpression",
        "left": {
          "type": "Identifier",
          "name": "x",
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 1
            }
          }
        },
        "operator": "!=",
        "right": {
          "type": "Identifier",
          "name": "y",
          "loc": {
            "start": {
              "line": 1,
              "column": 5
            },
            "end": {
              "line": 1,
              "column": 6
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 6
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 6
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 6
    }
  }
}
	`, ast)
}

func Test169(t *testing.T) {
	ast, err := compile("x === y")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "BinaryExpression",
        "left": {
          "type": "Identifier",
          "name": "x",
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 1
            }
          }
        },
        "operator": "===",
        "right": {
          "type": "Identifier",
          "name": "y",
          "loc": {
            "start": {
              "line": 1,
              "column": 6
            },
            "end": {
              "line": 1,
              "column": 7
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 7
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 7
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 7
    }
  }
}
	`, ast)
}

func Test170(t *testing.T) {
	ast, err := compile("x !== y")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "BinaryExpression",
        "left": {
          "type": "Identifier",
          "name": "x",
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 1
            }
          }
        },
        "operator": "!==",
        "right": {
          "type": "Identifier",
          "name": "y",
          "loc": {
            "start": {
              "line": 1,
              "column": 6
            },
            "end": {
              "line": 1,
              "column": 7
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 7
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 7
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 7
    }
  }
}
	`, ast)
}

func Test171(t *testing.T) {
	ast, err := compile("x & y")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "BinaryExpression",
        "left": {
          "type": "Identifier",
          "name": "x",
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 1
            }
          }
        },
        "operator": "&",
        "right": {
          "type": "Identifier",
          "name": "y",
          "loc": {
            "start": {
              "line": 1,
              "column": 4
            },
            "end": {
              "line": 1,
              "column": 5
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 5
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 5
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 5
    }
  }
}
	`, ast)
}

func Test172(t *testing.T) {
	ast, err := compile("x ^ y")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "BinaryExpression",
        "left": {
          "type": "Identifier",
          "name": "x",
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 1
            }
          }
        },
        "operator": "^",
        "right": {
          "type": "Identifier",
          "name": "y",
          "loc": {
            "start": {
              "line": 1,
              "column": 4
            },
            "end": {
              "line": 1,
              "column": 5
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 5
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 5
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 5
    }
  }
}
	`, ast)
}

func Test173(t *testing.T) {
	ast, err := compile("x | y")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "BinaryExpression",
        "left": {
          "type": "Identifier",
          "name": "x",
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 1
            }
          }
        },
        "operator": "|",
        "right": {
          "type": "Identifier",
          "name": "y",
          "loc": {
            "start": {
              "line": 1,
              "column": 4
            },
            "end": {
              "line": 1,
              "column": 5
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 5
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 5
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 5
    }
  }
}
	`, ast)
}

func Test174(t *testing.T) {
	ast, err := compile("x + y + z")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "BinaryExpression",
        "left": {
          "type": "BinaryExpression",
          "left": {
            "type": "Identifier",
            "name": "x",
            "loc": {
              "start": {
                "line": 1,
                "column": 0
              },
              "end": {
                "line": 1,
                "column": 1
              }
            }
          },
          "operator": "+",
          "right": {
            "type": "Identifier",
            "name": "y",
            "loc": {
              "start": {
                "line": 1,
                "column": 4
              },
              "end": {
                "line": 1,
                "column": 5
              }
            }
          },
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 5
            }
          }
        },
        "operator": "+",
        "right": {
          "type": "Identifier",
          "name": "z",
          "loc": {
            "start": {
              "line": 1,
              "column": 8
            },
            "end": {
              "line": 1,
              "column": 9
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 9
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 9
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 9
    }
  }
}
	`, ast)
}

func Test175(t *testing.T) {
	ast, err := compile("x - y + z")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "BinaryExpression",
        "left": {
          "type": "BinaryExpression",
          "left": {
            "type": "Identifier",
            "name": "x",
            "loc": {
              "start": {
                "line": 1,
                "column": 0
              },
              "end": {
                "line": 1,
                "column": 1
              }
            }
          },
          "operator": "-",
          "right": {
            "type": "Identifier",
            "name": "y",
            "loc": {
              "start": {
                "line": 1,
                "column": 4
              },
              "end": {
                "line": 1,
                "column": 5
              }
            }
          },
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 5
            }
          }
        },
        "operator": "+",
        "right": {
          "type": "Identifier",
          "name": "z",
          "loc": {
            "start": {
              "line": 1,
              "column": 8
            },
            "end": {
              "line": 1,
              "column": 9
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 9
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 9
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 9
    }
  }
}
	`, ast)
}

func Test176(t *testing.T) {
	ast, err := compile("x + y - z")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "BinaryExpression",
        "left": {
          "type": "BinaryExpression",
          "left": {
            "type": "Identifier",
            "name": "x",
            "loc": {
              "start": {
                "line": 1,
                "column": 0
              },
              "end": {
                "line": 1,
                "column": 1
              }
            }
          },
          "operator": "+",
          "right": {
            "type": "Identifier",
            "name": "y",
            "loc": {
              "start": {
                "line": 1,
                "column": 4
              },
              "end": {
                "line": 1,
                "column": 5
              }
            }
          },
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 5
            }
          }
        },
        "operator": "-",
        "right": {
          "type": "Identifier",
          "name": "z",
          "loc": {
            "start": {
              "line": 1,
              "column": 8
            },
            "end": {
              "line": 1,
              "column": 9
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 9
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 9
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 9
    }
  }
}
	`, ast)
}

func Test177(t *testing.T) {
	ast, err := compile("x - y - z")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "BinaryExpression",
        "left": {
          "type": "BinaryExpression",
          "left": {
            "type": "Identifier",
            "name": "x",
            "loc": {
              "start": {
                "line": 1,
                "column": 0
              },
              "end": {
                "line": 1,
                "column": 1
              }
            }
          },
          "operator": "-",
          "right": {
            "type": "Identifier",
            "name": "y",
            "loc": {
              "start": {
                "line": 1,
                "column": 4
              },
              "end": {
                "line": 1,
                "column": 5
              }
            }
          },
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 5
            }
          }
        },
        "operator": "-",
        "right": {
          "type": "Identifier",
          "name": "z",
          "loc": {
            "start": {
              "line": 1,
              "column": 8
            },
            "end": {
              "line": 1,
              "column": 9
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 9
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 9
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 9
    }
  }
}
	`, ast)
}

func Test178(t *testing.T) {
	ast, err := compile("x + y * z")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "BinaryExpression",
        "left": {
          "type": "Identifier",
          "name": "x",
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 1
            }
          }
        },
        "operator": "+",
        "right": {
          "type": "BinaryExpression",
          "left": {
            "type": "Identifier",
            "name": "y",
            "loc": {
              "start": {
                "line": 1,
                "column": 4
              },
              "end": {
                "line": 1,
                "column": 5
              }
            }
          },
          "operator": "*",
          "right": {
            "type": "Identifier",
            "name": "z",
            "loc": {
              "start": {
                "line": 1,
                "column": 8
              },
              "end": {
                "line": 1,
                "column": 9
              }
            }
          },
          "loc": {
            "start": {
              "line": 1,
              "column": 4
            },
            "end": {
              "line": 1,
              "column": 9
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 9
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 9
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 9
    }
  }
}
	`, ast)
}

func Test179(t *testing.T) {
	ast, err := compile("x + y / z")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "BinaryExpression",
        "left": {
          "type": "Identifier",
          "name": "x",
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 1
            }
          }
        },
        "operator": "+",
        "right": {
          "type": "BinaryExpression",
          "left": {
            "type": "Identifier",
            "name": "y",
            "loc": {
              "start": {
                "line": 1,
                "column": 4
              },
              "end": {
                "line": 1,
                "column": 5
              }
            }
          },
          "operator": "/",
          "right": {
            "type": "Identifier",
            "name": "z",
            "loc": {
              "start": {
                "line": 1,
                "column": 8
              },
              "end": {
                "line": 1,
                "column": 9
              }
            }
          },
          "loc": {
            "start": {
              "line": 1,
              "column": 4
            },
            "end": {
              "line": 1,
              "column": 9
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 9
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 9
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 9
    }
  }
}
	`, ast)
}

func Test180(t *testing.T) {
	ast, err := compile("x - y % z")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "BinaryExpression",
        "left": {
          "type": "Identifier",
          "name": "x",
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 1
            }
          }
        },
        "operator": "-",
        "right": {
          "type": "BinaryExpression",
          "left": {
            "type": "Identifier",
            "name": "y",
            "loc": {
              "start": {
                "line": 1,
                "column": 4
              },
              "end": {
                "line": 1,
                "column": 5
              }
            }
          },
          "operator": "%",
          "right": {
            "type": "Identifier",
            "name": "z",
            "loc": {
              "start": {
                "line": 1,
                "column": 8
              },
              "end": {
                "line": 1,
                "column": 9
              }
            }
          },
          "loc": {
            "start": {
              "line": 1,
              "column": 4
            },
            "end": {
              "line": 1,
              "column": 9
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 9
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 9
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 9
    }
  }
}
	`, ast)
}

func Test181(t *testing.T) {
	ast, err := compile("x * y * z")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "BinaryExpression",
        "left": {
          "type": "BinaryExpression",
          "left": {
            "type": "Identifier",
            "name": "x",
            "loc": {
              "start": {
                "line": 1,
                "column": 0
              },
              "end": {
                "line": 1,
                "column": 1
              }
            }
          },
          "operator": "*",
          "right": {
            "type": "Identifier",
            "name": "y",
            "loc": {
              "start": {
                "line": 1,
                "column": 4
              },
              "end": {
                "line": 1,
                "column": 5
              }
            }
          },
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 5
            }
          }
        },
        "operator": "*",
        "right": {
          "type": "Identifier",
          "name": "z",
          "loc": {
            "start": {
              "line": 1,
              "column": 8
            },
            "end": {
              "line": 1,
              "column": 9
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 9
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 9
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 9
    }
  }
}
	`, ast)
}

func Test182(t *testing.T) {
	ast, err := compile("x * y / z")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "BinaryExpression",
        "left": {
          "type": "BinaryExpression",
          "left": {
            "type": "Identifier",
            "name": "x",
            "loc": {
              "start": {
                "line": 1,
                "column": 0
              },
              "end": {
                "line": 1,
                "column": 1
              }
            }
          },
          "operator": "*",
          "right": {
            "type": "Identifier",
            "name": "y",
            "loc": {
              "start": {
                "line": 1,
                "column": 4
              },
              "end": {
                "line": 1,
                "column": 5
              }
            }
          },
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 5
            }
          }
        },
        "operator": "/",
        "right": {
          "type": "Identifier",
          "name": "z",
          "loc": {
            "start": {
              "line": 1,
              "column": 8
            },
            "end": {
              "line": 1,
              "column": 9
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 9
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 9
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 9
    }
  }
}
	`, ast)
}

func Test183(t *testing.T) {
	ast, err := compile("x * y % z")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "BinaryExpression",
        "left": {
          "type": "BinaryExpression",
          "left": {
            "type": "Identifier",
            "name": "x",
            "loc": {
              "start": {
                "line": 1,
                "column": 0
              },
              "end": {
                "line": 1,
                "column": 1
              }
            }
          },
          "operator": "*",
          "right": {
            "type": "Identifier",
            "name": "y",
            "loc": {
              "start": {
                "line": 1,
                "column": 4
              },
              "end": {
                "line": 1,
                "column": 5
              }
            }
          },
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 5
            }
          }
        },
        "operator": "%",
        "right": {
          "type": "Identifier",
          "name": "z",
          "loc": {
            "start": {
              "line": 1,
              "column": 8
            },
            "end": {
              "line": 1,
              "column": 9
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 9
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 9
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 9
    }
  }
}
	`, ast)
}

func Test184(t *testing.T) {
	ast, err := compile("x % y * z")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "BinaryExpression",
        "left": {
          "type": "BinaryExpression",
          "left": {
            "type": "Identifier",
            "name": "x",
            "loc": {
              "start": {
                "line": 1,
                "column": 0
              },
              "end": {
                "line": 1,
                "column": 1
              }
            }
          },
          "operator": "%",
          "right": {
            "type": "Identifier",
            "name": "y",
            "loc": {
              "start": {
                "line": 1,
                "column": 4
              },
              "end": {
                "line": 1,
                "column": 5
              }
            }
          },
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 5
            }
          }
        },
        "operator": "*",
        "right": {
          "type": "Identifier",
          "name": "z",
          "loc": {
            "start": {
              "line": 1,
              "column": 8
            },
            "end": {
              "line": 1,
              "column": 9
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 9
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 9
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 9
    }
  }
}
	`, ast)
}

func Test185(t *testing.T) {
	ast, err := compile("x << y << z")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "BinaryExpression",
        "left": {
          "type": "BinaryExpression",
          "left": {
            "type": "Identifier",
            "name": "x",
            "loc": {
              "start": {
                "line": 1,
                "column": 0
              },
              "end": {
                "line": 1,
                "column": 1
              }
            }
          },
          "operator": "<<",
          "right": {
            "type": "Identifier",
            "name": "y",
            "loc": {
              "start": {
                "line": 1,
                "column": 5
              },
              "end": {
                "line": 1,
                "column": 6
              }
            }
          },
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 6
            }
          }
        },
        "operator": "<<",
        "right": {
          "type": "Identifier",
          "name": "z",
          "loc": {
            "start": {
              "line": 1,
              "column": 10
            },
            "end": {
              "line": 1,
              "column": 11
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 11
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 11
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 11
    }
  }
}
	`, ast)
}

func Test186(t *testing.T) {
	ast, err := compile("x | y | z")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "BinaryExpression",
        "left": {
          "type": "BinaryExpression",
          "left": {
            "type": "Identifier",
            "name": "x",
            "loc": {
              "start": {
                "line": 1,
                "column": 0
              },
              "end": {
                "line": 1,
                "column": 1
              }
            }
          },
          "operator": "|",
          "right": {
            "type": "Identifier",
            "name": "y",
            "loc": {
              "start": {
                "line": 1,
                "column": 4
              },
              "end": {
                "line": 1,
                "column": 5
              }
            }
          },
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 5
            }
          }
        },
        "operator": "|",
        "right": {
          "type": "Identifier",
          "name": "z",
          "loc": {
            "start": {
              "line": 1,
              "column": 8
            },
            "end": {
              "line": 1,
              "column": 9
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 9
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 9
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 9
    }
  }
}
	`, ast)
}

func Test187(t *testing.T) {
	ast, err := compile("x & y & z")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "BinaryExpression",
        "left": {
          "type": "BinaryExpression",
          "left": {
            "type": "Identifier",
            "name": "x",
            "loc": {
              "start": {
                "line": 1,
                "column": 0
              },
              "end": {
                "line": 1,
                "column": 1
              }
            }
          },
          "operator": "&",
          "right": {
            "type": "Identifier",
            "name": "y",
            "loc": {
              "start": {
                "line": 1,
                "column": 4
              },
              "end": {
                "line": 1,
                "column": 5
              }
            }
          },
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 5
            }
          }
        },
        "operator": "&",
        "right": {
          "type": "Identifier",
          "name": "z",
          "loc": {
            "start": {
              "line": 1,
              "column": 8
            },
            "end": {
              "line": 1,
              "column": 9
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 9
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 9
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 9
    }
  }
}
	`, ast)
}

func Test188(t *testing.T) {
	ast, err := compile("x ^ y ^ z")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "BinaryExpression",
        "left": {
          "type": "BinaryExpression",
          "left": {
            "type": "Identifier",
            "name": "x",
            "loc": {
              "start": {
                "line": 1,
                "column": 0
              },
              "end": {
                "line": 1,
                "column": 1
              }
            }
          },
          "operator": "^",
          "right": {
            "type": "Identifier",
            "name": "y",
            "loc": {
              "start": {
                "line": 1,
                "column": 4
              },
              "end": {
                "line": 1,
                "column": 5
              }
            }
          },
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 5
            }
          }
        },
        "operator": "^",
        "right": {
          "type": "Identifier",
          "name": "z",
          "loc": {
            "start": {
              "line": 1,
              "column": 8
            },
            "end": {
              "line": 1,
              "column": 9
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 9
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 9
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 9
    }
  }
}
	`, ast)
}

func Test189(t *testing.T) {
	ast, err := compile("x & y | z")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "BinaryExpression",
        "left": {
          "type": "BinaryExpression",
          "left": {
            "type": "Identifier",
            "name": "x",
            "loc": {
              "start": {
                "line": 1,
                "column": 0
              },
              "end": {
                "line": 1,
                "column": 1
              }
            }
          },
          "operator": "&",
          "right": {
            "type": "Identifier",
            "name": "y",
            "loc": {
              "start": {
                "line": 1,
                "column": 4
              },
              "end": {
                "line": 1,
                "column": 5
              }
            }
          },
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 5
            }
          }
        },
        "operator": "|",
        "right": {
          "type": "Identifier",
          "name": "z",
          "loc": {
            "start": {
              "line": 1,
              "column": 8
            },
            "end": {
              "line": 1,
              "column": 9
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 9
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 9
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 9
    }
  }
}
	`, ast)
}

func Test190(t *testing.T) {
	ast, err := compile("x | y ^ z")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "BinaryExpression",
        "left": {
          "type": "Identifier",
          "name": "x",
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 1
            }
          }
        },
        "operator": "|",
        "right": {
          "type": "BinaryExpression",
          "left": {
            "type": "Identifier",
            "name": "y",
            "loc": {
              "start": {
                "line": 1,
                "column": 4
              },
              "end": {
                "line": 1,
                "column": 5
              }
            }
          },
          "operator": "^",
          "right": {
            "type": "Identifier",
            "name": "z",
            "loc": {
              "start": {
                "line": 1,
                "column": 8
              },
              "end": {
                "line": 1,
                "column": 9
              }
            }
          },
          "loc": {
            "start": {
              "line": 1,
              "column": 4
            },
            "end": {
              "line": 1,
              "column": 9
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 9
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 9
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 9
    }
  }
}
	`, ast)
}

func Test191(t *testing.T) {
	ast, err := compile("x | y & z")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "BinaryExpression",
        "left": {
          "type": "Identifier",
          "name": "x",
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 1
            }
          }
        },
        "operator": "|",
        "right": {
          "type": "BinaryExpression",
          "left": {
            "type": "Identifier",
            "name": "y",
            "loc": {
              "start": {
                "line": 1,
                "column": 4
              },
              "end": {
                "line": 1,
                "column": 5
              }
            }
          },
          "operator": "&",
          "right": {
            "type": "Identifier",
            "name": "z",
            "loc": {
              "start": {
                "line": 1,
                "column": 8
              },
              "end": {
                "line": 1,
                "column": 9
              }
            }
          },
          "loc": {
            "start": {
              "line": 1,
              "column": 4
            },
            "end": {
              "line": 1,
              "column": 9
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 9
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 9
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 9
    }
  }
}
	`, ast)
}

func Test192(t *testing.T) {
	ast, err := compile("x || y")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "LogicalExpression",
        "left": {
          "type": "Identifier",
          "name": "x",
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 1
            }
          }
        },
        "operator": "||",
        "right": {
          "type": "Identifier",
          "name": "y",
          "loc": {
            "start": {
              "line": 1,
              "column": 5
            },
            "end": {
              "line": 1,
              "column": 6
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 6
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 6
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 6
    }
  }
}
	`, ast)
}

func Test193(t *testing.T) {
	ast, err := compile("x && y")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "LogicalExpression",
        "left": {
          "type": "Identifier",
          "name": "x",
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 1
            }
          }
        },
        "operator": "&&",
        "right": {
          "type": "Identifier",
          "name": "y",
          "loc": {
            "start": {
              "line": 1,
              "column": 5
            },
            "end": {
              "line": 1,
              "column": 6
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 6
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 6
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 6
    }
  }
}
	`, ast)
}

func Test194(t *testing.T) {
	ast, err := compile("x || y || z")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "LogicalExpression",
        "left": {
          "type": "LogicalExpression",
          "left": {
            "type": "Identifier",
            "name": "x",
            "loc": {
              "start": {
                "line": 1,
                "column": 0
              },
              "end": {
                "line": 1,
                "column": 1
              }
            }
          },
          "operator": "||",
          "right": {
            "type": "Identifier",
            "name": "y",
            "loc": {
              "start": {
                "line": 1,
                "column": 5
              },
              "end": {
                "line": 1,
                "column": 6
              }
            }
          },
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 6
            }
          }
        },
        "operator": "||",
        "right": {
          "type": "Identifier",
          "name": "z",
          "loc": {
            "start": {
              "line": 1,
              "column": 10
            },
            "end": {
              "line": 1,
              "column": 11
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 11
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 11
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 11
    }
  }
}
	`, ast)
}

func Test195(t *testing.T) {
	ast, err := compile("x && y && z")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "LogicalExpression",
        "left": {
          "type": "LogicalExpression",
          "left": {
            "type": "Identifier",
            "name": "x",
            "loc": {
              "start": {
                "line": 1,
                "column": 0
              },
              "end": {
                "line": 1,
                "column": 1
              }
            }
          },
          "operator": "&&",
          "right": {
            "type": "Identifier",
            "name": "y",
            "loc": {
              "start": {
                "line": 1,
                "column": 5
              },
              "end": {
                "line": 1,
                "column": 6
              }
            }
          },
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 6
            }
          }
        },
        "operator": "&&",
        "right": {
          "type": "Identifier",
          "name": "z",
          "loc": {
            "start": {
              "line": 1,
              "column": 10
            },
            "end": {
              "line": 1,
              "column": 11
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 11
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 11
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 11
    }
  }
}
	`, ast)
}

func Test196(t *testing.T) {
	ast, err := compile("x || y && z")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "LogicalExpression",
        "left": {
          "type": "Identifier",
          "name": "x",
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 1
            }
          }
        },
        "operator": "||",
        "right": {
          "type": "LogicalExpression",
          "left": {
            "type": "Identifier",
            "name": "y",
            "loc": {
              "start": {
                "line": 1,
                "column": 5
              },
              "end": {
                "line": 1,
                "column": 6
              }
            }
          },
          "operator": "&&",
          "right": {
            "type": "Identifier",
            "name": "z",
            "loc": {
              "start": {
                "line": 1,
                "column": 10
              },
              "end": {
                "line": 1,
                "column": 11
              }
            }
          },
          "loc": {
            "start": {
              "line": 1,
              "column": 5
            },
            "end": {
              "line": 1,
              "column": 11
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 11
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 11
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 11
    }
  }
}
	`, ast)
}

func Test197(t *testing.T) {
	ast, err := compile("x || y ^ z")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "LogicalExpression",
        "left": {
          "type": "Identifier",
          "name": "x",
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 1
            }
          }
        },
        "operator": "||",
        "right": {
          "type": "BinaryExpression",
          "left": {
            "type": "Identifier",
            "name": "y",
            "loc": {
              "start": {
                "line": 1,
                "column": 5
              },
              "end": {
                "line": 1,
                "column": 6
              }
            }
          },
          "operator": "^",
          "right": {
            "type": "Identifier",
            "name": "z",
            "loc": {
              "start": {
                "line": 1,
                "column": 9
              },
              "end": {
                "line": 1,
                "column": 10
              }
            }
          },
          "loc": {
            "start": {
              "line": 1,
              "column": 5
            },
            "end": {
              "line": 1,
              "column": 10
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 10
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 10
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 10
    }
  }
}
	`, ast)
}

func Test198(t *testing.T) {
	ast, err := compile("y ? 1 : 2")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "ConditionalExpression",
        "test": {
          "type": "Identifier",
          "name": "y",
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 1
            }
          }
        },
        "consequent": {
          "type": "Literal",
          "value": 1,
          "loc": {
            "start": {
              "line": 1,
              "column": 4
            },
            "end": {
              "line": 1,
              "column": 5
            }
          }
        },
        "alternate": {
          "type": "Literal",
          "value": 2,
          "loc": {
            "start": {
              "line": 1,
              "column": 8
            },
            "end": {
              "line": 1,
              "column": 9
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 9
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 9
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 9
    }
  }
}
	`, ast)
}

func Test199(t *testing.T) {
	ast, err := compile("x && y ? 1 : 2")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "ConditionalExpression",
        "test": {
          "type": "LogicalExpression",
          "left": {
            "type": "Identifier",
            "name": "x",
            "loc": {
              "start": {
                "line": 1,
                "column": 0
              },
              "end": {
                "line": 1,
                "column": 1
              }
            }
          },
          "operator": "&&",
          "right": {
            "type": "Identifier",
            "name": "y",
            "loc": {
              "start": {
                "line": 1,
                "column": 5
              },
              "end": {
                "line": 1,
                "column": 6
              }
            }
          },
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 6
            }
          }
        },
        "consequent": {
          "type": "Literal",
          "value": 1,
          "loc": {
            "start": {
              "line": 1,
              "column": 9
            },
            "end": {
              "line": 1,
              "column": 10
            }
          }
        },
        "alternate": {
          "type": "Literal",
          "value": 2,
          "loc": {
            "start": {
              "line": 1,
              "column": 13
            },
            "end": {
              "line": 1,
              "column": 14
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 14
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 14
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 14
    }
  }
}
	`, ast)
}

func Test200(t *testing.T) {
	ast, err := compile("x = 42")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "AssignmentExpression",
        "operator": "=",
        "left": {
          "type": "Identifier",
          "name": "x",
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 1
            }
          }
        },
        "right": {
          "type": "Literal",
          "value": 42,
          "loc": {
            "start": {
              "line": 1,
              "column": 4
            },
            "end": {
              "line": 1,
              "column": 6
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 6
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 6
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 6
    }
  }
}
	`, ast)
}

func Test201(t *testing.T) {
	ast, err := compile("eval = 42")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "AssignmentExpression",
        "operator": "=",
        "left": {
          "type": "Identifier",
          "name": "eval",
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 4
            }
          }
        },
        "right": {
          "type": "Literal",
          "value": 42,
          "loc": {
            "start": {
              "line": 1,
              "column": 7
            },
            "end": {
              "line": 1,
              "column": 9
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 9
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 9
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 9
    }
  }
}
	`, ast)
}

func Test202(t *testing.T) {
	ast, err := compile("arguments = 42")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "AssignmentExpression",
        "operator": "=",
        "left": {
          "type": "Identifier",
          "name": "arguments",
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 9
            }
          }
        },
        "right": {
          "type": "Literal",
          "value": 42,
          "loc": {
            "start": {
              "line": 1,
              "column": 12
            },
            "end": {
              "line": 1,
              "column": 14
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 14
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 14
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 14
    }
  }
}
	`, ast)
}

func Test203(t *testing.T) {
	ast, err := compile("x *= 42")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "AssignmentExpression",
        "operator": "*=",
        "left": {
          "type": "Identifier",
          "name": "x",
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 1
            }
          }
        },
        "right": {
          "type": "Literal",
          "value": 42,
          "loc": {
            "start": {
              "line": 1,
              "column": 5
            },
            "end": {
              "line": 1,
              "column": 7
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 7
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 7
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 7
    }
  }
}
	`, ast)
}

func Test204(t *testing.T) {
	ast, err := compile("x /= 42")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "AssignmentExpression",
        "operator": "/=",
        "left": {
          "type": "Identifier",
          "name": "x",
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 1
            }
          }
        },
        "right": {
          "type": "Literal",
          "value": 42,
          "loc": {
            "start": {
              "line": 1,
              "column": 5
            },
            "end": {
              "line": 1,
              "column": 7
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 7
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 7
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 7
    }
  }
}
	`, ast)
}

func Test205(t *testing.T) {
	ast, err := compile("x %= 42")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "AssignmentExpression",
        "operator": "%=",
        "left": {
          "type": "Identifier",
          "name": "x",
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 1
            }
          }
        },
        "right": {
          "type": "Literal",
          "value": 42,
          "loc": {
            "start": {
              "line": 1,
              "column": 5
            },
            "end": {
              "line": 1,
              "column": 7
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 7
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 7
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 7
    }
  }
}
	`, ast)
}

func Test206(t *testing.T) {
	ast, err := compile("x += 42")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "AssignmentExpression",
        "operator": "+=",
        "left": {
          "type": "Identifier",
          "name": "x",
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 1
            }
          }
        },
        "right": {
          "type": "Literal",
          "value": 42,
          "loc": {
            "start": {
              "line": 1,
              "column": 5
            },
            "end": {
              "line": 1,
              "column": 7
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 7
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 7
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 7
    }
  }
}
	`, ast)
}

func Test207(t *testing.T) {
	ast, err := compile("x -= 42")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "AssignmentExpression",
        "operator": "-=",
        "left": {
          "type": "Identifier",
          "name": "x",
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 1
            }
          }
        },
        "right": {
          "type": "Literal",
          "value": 42,
          "loc": {
            "start": {
              "line": 1,
              "column": 5
            },
            "end": {
              "line": 1,
              "column": 7
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 7
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 7
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 7
    }
  }
}
	`, ast)
}

func Test208(t *testing.T) {
	ast, err := compile("x <<= 42")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "AssignmentExpression",
        "operator": "<<=",
        "left": {
          "type": "Identifier",
          "name": "x",
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 1
            }
          }
        },
        "right": {
          "type": "Literal",
          "value": 42,
          "loc": {
            "start": {
              "line": 1,
              "column": 6
            },
            "end": {
              "line": 1,
              "column": 8
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 8
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 8
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 8
    }
  }
}
	`, ast)
}

func Test209(t *testing.T) {
	ast, err := compile("x >>= 42")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "AssignmentExpression",
        "operator": ">>=",
        "left": {
          "type": "Identifier",
          "name": "x",
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 1
            }
          }
        },
        "right": {
          "type": "Literal",
          "value": 42,
          "loc": {
            "start": {
              "line": 1,
              "column": 6
            },
            "end": {
              "line": 1,
              "column": 8
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 8
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 8
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 8
    }
  }
}
	`, ast)
}

func Test210(t *testing.T) {
	ast, err := compile("x >>>= 42")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "AssignmentExpression",
        "operator": ">>>=",
        "left": {
          "type": "Identifier",
          "name": "x",
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 1
            }
          }
        },
        "right": {
          "type": "Literal",
          "value": 42,
          "loc": {
            "start": {
              "line": 1,
              "column": 7
            },
            "end": {
              "line": 1,
              "column": 9
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 9
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 9
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 9
    }
  }
}
	`, ast)
}

func Test211(t *testing.T) {
	ast, err := compile("x &= 42")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "AssignmentExpression",
        "operator": "&=",
        "left": {
          "type": "Identifier",
          "name": "x",
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 1
            }
          }
        },
        "right": {
          "type": "Literal",
          "value": 42,
          "loc": {
            "start": {
              "line": 1,
              "column": 5
            },
            "end": {
              "line": 1,
              "column": 7
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 7
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 7
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 7
    }
  }
}
	`, ast)
}

func Test212(t *testing.T) {
	ast, err := compile("x ^= 42")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "AssignmentExpression",
        "operator": "^=",
        "left": {
          "type": "Identifier",
          "name": "x",
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 1
            }
          }
        },
        "right": {
          "type": "Literal",
          "value": 42,
          "loc": {
            "start": {
              "line": 1,
              "column": 5
            },
            "end": {
              "line": 1,
              "column": 7
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 7
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 7
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 7
    }
  }
}
	`, ast)
}

func Test213(t *testing.T) {
	ast, err := compile("x |= 42")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "AssignmentExpression",
        "operator": "|=",
        "left": {
          "type": "Identifier",
          "name": "x",
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 1
            }
          }
        },
        "right": {
          "type": "Literal",
          "value": 42,
          "loc": {
            "start": {
              "line": 1,
              "column": 5
            },
            "end": {
              "line": 1,
              "column": 7
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 7
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 7
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 7
    }
  }
}
	`, ast)
}

func Test214(t *testing.T) {
	ast, err := compile("{ foo }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "BlockStatement",
      "body": [
        {
          "type": "ExpressionStatement",
          "expression": {
            "type": "Identifier",
            "name": "foo",
            "loc": {
              "start": {
                "line": 1,
                "column": 2
              },
              "end": {
                "line": 1,
                "column": 5
              }
            }
          },
          "loc": {
            "start": {
              "line": 1,
              "column": 2
            },
            "end": {
              "line": 1,
              "column": 5
            }
          }
        }
      ],
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 7
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 7
    }
  }
}
	`, ast)
}

func Test215(t *testing.T) {
	ast, err := compile("{ doThis(); doThat(); }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "BlockStatement",
      "body": [
        {
          "type": "ExpressionStatement",
          "expression": {
            "type": "CallExpression",
            "callee": {
              "type": "Identifier",
              "name": "doThis",
              "loc": {
                "start": {
                  "line": 1,
                  "column": 2
                },
                "end": {
                  "line": 1,
                  "column": 8
                }
              }
            },
            "arguments": [],
            "loc": {
              "start": {
                "line": 1,
                "column": 2
              },
              "end": {
                "line": 1,
                "column": 10
              }
            }
          },
          "loc": {
            "start": {
              "line": 1,
              "column": 2
            },
            "end": {
              "line": 1,
              "column": 11
            }
          }
        },
        {
          "type": "ExpressionStatement",
          "expression": {
            "type": "CallExpression",
            "callee": {
              "type": "Identifier",
              "name": "doThat",
              "loc": {
                "start": {
                  "line": 1,
                  "column": 12
                },
                "end": {
                  "line": 1,
                  "column": 18
                }
              }
            },
            "arguments": [],
            "loc": {
              "start": {
                "line": 1,
                "column": 12
              },
              "end": {
                "line": 1,
                "column": 20
              }
            }
          },
          "loc": {
            "start": {
              "line": 1,
              "column": 12
            },
            "end": {
              "line": 1,
              "column": 21
            }
          }
        }
      ],
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 23
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 23
    }
  }
}
	`, ast)
}

func Test216(t *testing.T) {
	ast, err := compile("{}")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "BlockStatement",
      "body": [],
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 2
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 2
    }
  }
}
	`, ast)
}

func Test217(t *testing.T) {
	ast, err := compile("var x")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "VariableDeclaration",
      "declarations": [
        {
          "type": "VariableDeclarator",
          "id": {
            "type": "Identifier",
            "name": "x",
            "loc": {
              "start": {
                "line": 1,
                "column": 4
              },
              "end": {
                "line": 1,
                "column": 5
              }
            }
          },
          "init": null,
          "loc": {
            "start": {
              "line": 1,
              "column": 4
            },
            "end": {
              "line": 1,
              "column": 5
            }
          }
        }
      ],
      "kind": "var",
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 5
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 5
    }
  }
}
	`, ast)
}

func Test218(t *testing.T) {
	ast, err := compile("var await")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "VariableDeclaration",
      "declarations": [
        {
          "type": "VariableDeclarator",
          "id": {
            "type": "Identifier",
            "name": "await",
            "loc": {
              "start": {
                "line": 1,
                "column": 4
              },
              "end": {
                "line": 1,
                "column": 9
              }
            }
          },
          "init": null,
          "loc": {
            "start": {
              "line": 1,
              "column": 4
            },
            "end": {
              "line": 1,
              "column": 9
            }
          }
        }
      ],
      "kind": "var",
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 9
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 9
    }
  }
}
	`, ast)
}

func Test219(t *testing.T) {
	ast, err := compile("var x, y;")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "VariableDeclaration",
      "declarations": [
        {
          "type": "VariableDeclarator",
          "id": {
            "type": "Identifier",
            "name": "x",
            "loc": {
              "start": {
                "line": 1,
                "column": 4
              },
              "end": {
                "line": 1,
                "column": 5
              }
            }
          },
          "init": null,
          "loc": {
            "start": {
              "line": 1,
              "column": 4
            },
            "end": {
              "line": 1,
              "column": 5
            }
          }
        },
        {
          "type": "VariableDeclarator",
          "id": {
            "type": "Identifier",
            "name": "y",
            "loc": {
              "start": {
                "line": 1,
                "column": 7
              },
              "end": {
                "line": 1,
                "column": 8
              }
            }
          },
          "init": null,
          "loc": {
            "start": {
              "line": 1,
              "column": 7
            },
            "end": {
              "line": 1,
              "column": 8
            }
          }
        }
      ],
      "kind": "var",
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 9
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 9
    }
  }
}
	`, ast)
}

func Test220(t *testing.T) {
	ast, err := compile("var x = 42")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "VariableDeclaration",
      "declarations": [
        {
          "type": "VariableDeclarator",
          "id": {
            "type": "Identifier",
            "name": "x",
            "loc": {
              "start": {
                "line": 1,
                "column": 4
              },
              "end": {
                "line": 1,
                "column": 5
              }
            }
          },
          "init": {
            "type": "Literal",
            "value": 42,
            "loc": {
              "start": {
                "line": 1,
                "column": 8
              },
              "end": {
                "line": 1,
                "column": 10
              }
            }
          },
          "loc": {
            "start": {
              "line": 1,
              "column": 4
            },
            "end": {
              "line": 1,
              "column": 10
            }
          }
        }
      ],
      "kind": "var",
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 10
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 10
    }
  }
}
	`, ast)
}

func Test221(t *testing.T) {
	ast, err := compile("var eval = 42, arguments = 42")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "VariableDeclaration",
      "declarations": [
        {
          "type": "VariableDeclarator",
          "id": {
            "type": "Identifier",
            "name": "eval",
            "loc": {
              "start": {
                "line": 1,
                "column": 4
              },
              "end": {
                "line": 1,
                "column": 8
              }
            }
          },
          "init": {
            "type": "Literal",
            "value": 42,
            "loc": {
              "start": {
                "line": 1,
                "column": 11
              },
              "end": {
                "line": 1,
                "column": 13
              }
            }
          },
          "loc": {
            "start": {
              "line": 1,
              "column": 4
            },
            "end": {
              "line": 1,
              "column": 13
            }
          }
        },
        {
          "type": "VariableDeclarator",
          "id": {
            "type": "Identifier",
            "name": "arguments",
            "loc": {
              "start": {
                "line": 1,
                "column": 15
              },
              "end": {
                "line": 1,
                "column": 24
              }
            }
          },
          "init": {
            "type": "Literal",
            "value": 42,
            "loc": {
              "start": {
                "line": 1,
                "column": 27
              },
              "end": {
                "line": 1,
                "column": 29
              }
            }
          },
          "loc": {
            "start": {
              "line": 1,
              "column": 15
            },
            "end": {
              "line": 1,
              "column": 29
            }
          }
        }
      ],
      "kind": "var",
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 29
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 29
    }
  }
}
	`, ast)
}

func Test222(t *testing.T) {
	ast, err := compile("var x = 14, y = 3, z = 1977")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "VariableDeclaration",
      "declarations": [
        {
          "type": "VariableDeclarator",
          "id": {
            "type": "Identifier",
            "name": "x",
            "loc": {
              "start": {
                "line": 1,
                "column": 4
              },
              "end": {
                "line": 1,
                "column": 5
              }
            }
          },
          "init": {
            "type": "Literal",
            "value": 14,
            "loc": {
              "start": {
                "line": 1,
                "column": 8
              },
              "end": {
                "line": 1,
                "column": 10
              }
            }
          },
          "loc": {
            "start": {
              "line": 1,
              "column": 4
            },
            "end": {
              "line": 1,
              "column": 10
            }
          }
        },
        {
          "type": "VariableDeclarator",
          "id": {
            "type": "Identifier",
            "name": "y",
            "loc": {
              "start": {
                "line": 1,
                "column": 12
              },
              "end": {
                "line": 1,
                "column": 13
              }
            }
          },
          "init": {
            "type": "Literal",
            "value": 3,
            "loc": {
              "start": {
                "line": 1,
                "column": 16
              },
              "end": {
                "line": 1,
                "column": 17
              }
            }
          },
          "loc": {
            "start": {
              "line": 1,
              "column": 12
            },
            "end": {
              "line": 1,
              "column": 17
            }
          }
        },
        {
          "type": "VariableDeclarator",
          "id": {
            "type": "Identifier",
            "name": "z",
            "loc": {
              "start": {
                "line": 1,
                "column": 19
              },
              "end": {
                "line": 1,
                "column": 20
              }
            }
          },
          "init": {
            "type": "Literal",
            "value": 1977,
            "loc": {
              "start": {
                "line": 1,
                "column": 23
              },
              "end": {
                "line": 1,
                "column": 27
              }
            }
          },
          "loc": {
            "start": {
              "line": 1,
              "column": 19
            },
            "end": {
              "line": 1,
              "column": 27
            }
          }
        }
      ],
      "kind": "var",
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 27
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 27
    }
  }
}
	`, ast)
}

func Test223(t *testing.T) {
	ast, err := compile("var implements, interface, package")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "VariableDeclaration",
      "declarations": [
        {
          "type": "VariableDeclarator",
          "id": {
            "type": "Identifier",
            "name": "implements",
            "loc": {
              "start": {
                "line": 1,
                "column": 4
              },
              "end": {
                "line": 1,
                "column": 14
              }
            }
          },
          "init": null,
          "loc": {
            "start": {
              "line": 1,
              "column": 4
            },
            "end": {
              "line": 1,
              "column": 14
            }
          }
        },
        {
          "type": "VariableDeclarator",
          "id": {
            "type": "Identifier",
            "name": "interface",
            "loc": {
              "start": {
                "line": 1,
                "column": 16
              },
              "end": {
                "line": 1,
                "column": 25
              }
            }
          },
          "init": null,
          "loc": {
            "start": {
              "line": 1,
              "column": 16
            },
            "end": {
              "line": 1,
              "column": 25
            }
          }
        },
        {
          "type": "VariableDeclarator",
          "id": {
            "type": "Identifier",
            "name": "package",
            "loc": {
              "start": {
                "line": 1,
                "column": 27
              },
              "end": {
                "line": 1,
                "column": 34
              }
            }
          },
          "init": null,
          "loc": {
            "start": {
              "line": 1,
              "column": 27
            },
            "end": {
              "line": 1,
              "column": 34
            }
          }
        }
      ],
      "kind": "var",
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 34
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 34
    }
  }
}
	`, ast)
}

func Test224(t *testing.T) {
	ast, err := compile("var private, protected, public, static")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "VariableDeclaration",
      "declarations": [
        {
          "type": "VariableDeclarator",
          "id": {
            "type": "Identifier",
            "name": "private",
            "loc": {
              "start": {
                "line": 1,
                "column": 4
              },
              "end": {
                "line": 1,
                "column": 11
              }
            }
          },
          "init": null,
          "loc": {
            "start": {
              "line": 1,
              "column": 4
            },
            "end": {
              "line": 1,
              "column": 11
            }
          }
        },
        {
          "type": "VariableDeclarator",
          "id": {
            "type": "Identifier",
            "name": "protected",
            "loc": {
              "start": {
                "line": 1,
                "column": 13
              },
              "end": {
                "line": 1,
                "column": 22
              }
            }
          },
          "init": null,
          "loc": {
            "start": {
              "line": 1,
              "column": 13
            },
            "end": {
              "line": 1,
              "column": 22
            }
          }
        },
        {
          "type": "VariableDeclarator",
          "id": {
            "type": "Identifier",
            "name": "public",
            "loc": {
              "start": {
                "line": 1,
                "column": 24
              },
              "end": {
                "line": 1,
                "column": 30
              }
            }
          },
          "init": null,
          "loc": {
            "start": {
              "line": 1,
              "column": 24
            },
            "end": {
              "line": 1,
              "column": 30
            }
          }
        },
        {
          "type": "VariableDeclarator",
          "id": {
            "type": "Identifier",
            "name": "static",
            "loc": {
              "start": {
                "line": 1,
                "column": 32
              },
              "end": {
                "line": 1,
                "column": 38
              }
            }
          },
          "init": null,
          "loc": {
            "start": {
              "line": 1,
              "column": 32
            },
            "end": {
              "line": 1,
              "column": 38
            }
          }
        }
      ],
      "kind": "var",
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 38
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 38
    }
  }
}
	`, ast)
}

func Test225(t *testing.T) {
	ast, err := compile(";")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "EmptyStatement",
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 1
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 1
    }
  }
}
	`, ast)
}

func Test226(t *testing.T) {
	ast, err := compile("x")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "Identifier",
        "name": "x",
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 1
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 1
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 1
    }
  }
}
	`, ast)
}

func Test227(t *testing.T) {
	ast, err := compile("x, y")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "SequenceExpression",
        "expressions": [
          {
            "type": "Identifier",
            "name": "x",
            "loc": {
              "start": {
                "line": 1,
                "column": 0
              },
              "end": {
                "line": 1,
                "column": 1
              }
            }
          },
          {
            "type": "Identifier",
            "name": "y",
            "loc": {
              "start": {
                "line": 1,
                "column": 3
              },
              "end": {
                "line": 1,
                "column": 4
              }
            }
          }
        ],
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 4
          }
        }
      },
      "loc": {
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
  ],
  "loc": {
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
	`, ast)
}

func Test228(t *testing.T) {
	ast, err := compile("\\u0061")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "Identifier",
        "name": "a",
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 6
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 6
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 6
    }
  }
}
	`, ast)
}

func Test229(t *testing.T) {
	ast, err := compile("a\\u0061")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "Identifier",
        "name": "aa",
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 7
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 7
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 7
    }
  }
}
	`, ast)
}

func Test230(t *testing.T) {
	ast, err := compile("if (morning) goodMorning()")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "IfStatement",
      "test": {
        "type": "Identifier",
        "name": "morning",
        "loc": {
          "start": {
            "line": 1,
            "column": 4
          },
          "end": {
            "line": 1,
            "column": 11
          }
        }
      },
      "consequent": {
        "type": "ExpressionStatement",
        "expression": {
          "type": "CallExpression",
          "callee": {
            "type": "Identifier",
            "name": "goodMorning",
            "loc": {
              "start": {
                "line": 1,
                "column": 13
              },
              "end": {
                "line": 1,
                "column": 24
              }
            }
          },
          "arguments": [],
          "loc": {
            "start": {
              "line": 1,
              "column": 13
            },
            "end": {
              "line": 1,
              "column": 26
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 13
          },
          "end": {
            "line": 1,
            "column": 26
          }
        }
      },
      "alternate": null,
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 26
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 26
    }
  }
}
	`, ast)
}

func Test231(t *testing.T) {
	ast, err := compile("if (morning) (function(){})")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "IfStatement",
      "test": {
        "type": "Identifier",
        "name": "morning",
        "loc": {
          "start": {
            "line": 1,
            "column": 4
          },
          "end": {
            "line": 1,
            "column": 11
          }
        }
      },
      "consequent": {
        "type": "ExpressionStatement",
        "expression": {
          "type": "FunctionExpression",
          "id": null,
          "params": [],
          "body": {
            "type": "BlockStatement",
            "body": [],
            "loc": {
              "start": {
                "line": 1,
                "column": 24
              },
              "end": {
                "line": 1,
                "column": 26
              }
            }
          },
          "loc": {
            "start": {
              "line": 1,
              "column": 14
            },
            "end": {
              "line": 1,
              "column": 26
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 13
          },
          "end": {
            "line": 1,
            "column": 27
          }
        }
      },
      "alternate": null,
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 27
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 27
    }
  }
}
	`, ast)
}

func Test232(t *testing.T) {
	ast, err := compile("if (morning) var x = 0;")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "IfStatement",
      "test": {
        "type": "Identifier",
        "name": "morning",
        "loc": {
          "start": {
            "line": 1,
            "column": 4
          },
          "end": {
            "line": 1,
            "column": 11
          }
        }
      },
      "consequent": {
        "type": "VariableDeclaration",
        "declarations": [
          {
            "type": "VariableDeclarator",
            "id": {
              "type": "Identifier",
              "name": "x",
              "loc": {
                "start": {
                  "line": 1,
                  "column": 17
                },
                "end": {
                  "line": 1,
                  "column": 18
                }
              }
            },
            "init": {
              "type": "Literal",
              "value": 0,
              "loc": {
                "start": {
                  "line": 1,
                  "column": 21
                },
                "end": {
                  "line": 1,
                  "column": 22
                }
              }
            },
            "loc": {
              "start": {
                "line": 1,
                "column": 17
              },
              "end": {
                "line": 1,
                "column": 22
              }
            }
          }
        ],
        "kind": "var",
        "loc": {
          "start": {
            "line": 1,
            "column": 13
          },
          "end": {
            "line": 1,
            "column": 23
          }
        }
      },
      "alternate": null,
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 23
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 23
    }
  }
}
	`, ast)
}

func Test233(t *testing.T) {
	ast, err := compile("if (morning) function a(){}")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "IfStatement",
      "test": {
        "type": "Identifier",
        "name": "morning",
        "loc": {
          "start": {
            "line": 1,
            "column": 4
          },
          "end": {
            "line": 1,
            "column": 11
          }
        }
      },
      "consequent": {
        "type": "FunctionDeclaration",
        "id": {
          "type": "Identifier",
          "name": "a",
          "loc": {
            "start": {
              "line": 1,
              "column": 22
            },
            "end": {
              "line": 1,
              "column": 23
            }
          }
        },
        "params": [],
        "body": {
          "type": "BlockStatement",
          "body": [],
          "loc": {
            "start": {
              "line": 1,
              "column": 25
            },
            "end": {
              "line": 1,
              "column": 27
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 13
          },
          "end": {
            "line": 1,
            "column": 27
          }
        }
      },
      "alternate": null,
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 27
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 27
    }
  }
}
	`, ast)
}

func Test234(t *testing.T) {
	ast, err := compile("if (morning) goodMorning(); else goodDay()")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "IfStatement",
      "test": {
        "type": "Identifier",
        "name": "morning",
        "loc": {
          "start": {
            "line": 1,
            "column": 4
          },
          "end": {
            "line": 1,
            "column": 11
          }
        }
      },
      "consequent": {
        "type": "ExpressionStatement",
        "expression": {
          "type": "CallExpression",
          "callee": {
            "type": "Identifier",
            "name": "goodMorning",
            "loc": {
              "start": {
                "line": 1,
                "column": 13
              },
              "end": {
                "line": 1,
                "column": 24
              }
            }
          },
          "arguments": [],
          "loc": {
            "start": {
              "line": 1,
              "column": 13
            },
            "end": {
              "line": 1,
              "column": 26
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 13
          },
          "end": {
            "line": 1,
            "column": 27
          }
        }
      },
      "alternate": {
        "type": "ExpressionStatement",
        "expression": {
          "type": "CallExpression",
          "callee": {
            "type": "Identifier",
            "name": "goodDay",
            "loc": {
              "start": {
                "line": 1,
                "column": 33
              },
              "end": {
                "line": 1,
                "column": 40
              }
            }
          },
          "arguments": [],
          "loc": {
            "start": {
              "line": 1,
              "column": 33
            },
            "end": {
              "line": 1,
              "column": 42
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 33
          },
          "end": {
            "line": 1,
            "column": 42
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 42
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 42
    }
  }
}
	`, ast)
}

func Test235(t *testing.T) {
	ast, err := compile("do keep(); while (true);")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "DoWhileStatement",
      "body": {
        "type": "ExpressionStatement",
        "expression": {
          "type": "CallExpression",
          "callee": {
            "type": "Identifier",
            "name": "keep",
            "loc": {
              "start": {
                "line": 1,
                "column": 3
              },
              "end": {
                "line": 1,
                "column": 7
              }
            }
          },
          "arguments": [],
          "loc": {
            "start": {
              "line": 1,
              "column": 3
            },
            "end": {
              "line": 1,
              "column": 9
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 3
          },
          "end": {
            "line": 1,
            "column": 10
          }
        }
      },
      "test": {
        "type": "Literal",
        "value": true,
        "loc": {
          "start": {
            "line": 1,
            "column": 18
          },
          "end": {
            "line": 1,
            "column": 22
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 24
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 24
    }
  }
}
	`, ast)
}

func Test236(t *testing.T) {
	ast, err := compile("{ do { } while (false);false }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "BlockStatement",
      "body": [
        {
          "type": "DoWhileStatement",
          "body": {
            "type": "BlockStatement",
            "body": [],
            "loc": {
              "start": {
                "line": 1,
                "column": 5
              },
              "end": {
                "line": 1,
                "column": 8
              }
            }
          },
          "test": {
            "type": "Literal",
            "value": false,
            "loc": {
              "start": {
                "line": 1,
                "column": 16
              },
              "end": {
                "line": 1,
                "column": 21
              }
            }
          },
          "loc": {
            "start": {
              "line": 1,
              "column": 2
            },
            "end": {
              "line": 1,
              "column": 23
            }
          }
        },
        {
          "type": "ExpressionStatement",
          "expression": {
            "type": "Literal",
            "value": false,
            "loc": {
              "start": {
                "line": 1,
                "column": 23
              },
              "end": {
                "line": 1,
                "column": 28
              }
            }
          },
          "loc": {
            "start": {
              "line": 1,
              "column": 23
            },
            "end": {
              "line": 1,
              "column": 28
            }
          }
        }
      ],
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 30
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 30
    }
  }
}
	`, ast)
}

func Test237(t *testing.T) {
	ast, err := compile("while (true) doSomething()")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "WhileStatement",
      "test": {
        "type": "Literal",
        "value": true,
        "loc": {
          "start": {
            "line": 1,
            "column": 7
          },
          "end": {
            "line": 1,
            "column": 11
          }
        }
      },
      "body": {
        "type": "ExpressionStatement",
        "expression": {
          "type": "CallExpression",
          "callee": {
            "type": "Identifier",
            "name": "doSomething",
            "loc": {
              "start": {
                "line": 1,
                "column": 13
              },
              "end": {
                "line": 1,
                "column": 24
              }
            }
          },
          "arguments": [],
          "loc": {
            "start": {
              "line": 1,
              "column": 13
            },
            "end": {
              "line": 1,
              "column": 26
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 13
          },
          "end": {
            "line": 1,
            "column": 26
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 26
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 26
    }
  }
}
	`, ast)
}

func Test238(t *testing.T) {
	ast, err := compile("while (x < 10) { x++; y--; }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "WhileStatement",
      "test": {
        "type": "BinaryExpression",
        "left": {
          "type": "Identifier",
          "name": "x",
          "loc": {
            "start": {
              "line": 1,
              "column": 7
            },
            "end": {
              "line": 1,
              "column": 8
            }
          }
        },
        "operator": "<",
        "right": {
          "type": "Literal",
          "value": 10,
          "loc": {
            "start": {
              "line": 1,
              "column": 11
            },
            "end": {
              "line": 1,
              "column": 13
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 7
          },
          "end": {
            "line": 1,
            "column": 13
          }
        }
      },
      "body": {
        "type": "BlockStatement",
        "body": [
          {
            "type": "ExpressionStatement",
            "expression": {
              "type": "UpdateExpression",
              "operator": "++",
              "prefix": false,
              "argument": {
                "type": "Identifier",
                "name": "x",
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 17
                  },
                  "end": {
                    "line": 1,
                    "column": 18
                  }
                }
              },
              "loc": {
                "start": {
                  "line": 1,
                  "column": 17
                },
                "end": {
                  "line": 1,
                  "column": 20
                }
              }
            },
            "loc": {
              "start": {
                "line": 1,
                "column": 17
              },
              "end": {
                "line": 1,
                "column": 21
              }
            }
          },
          {
            "type": "ExpressionStatement",
            "expression": {
              "type": "UpdateExpression",
              "operator": "--",
              "prefix": false,
              "argument": {
                "type": "Identifier",
                "name": "y",
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 22
                  },
                  "end": {
                    "line": 1,
                    "column": 23
                  }
                }
              },
              "loc": {
                "start": {
                  "line": 1,
                  "column": 22
                },
                "end": {
                  "line": 1,
                  "column": 25
                }
              }
            },
            "loc": {
              "start": {
                "line": 1,
                "column": 22
              },
              "end": {
                "line": 1,
                "column": 26
              }
            }
          }
        ],
        "loc": {
          "start": {
            "line": 1,
            "column": 15
          },
          "end": {
            "line": 1,
            "column": 28
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 28
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 28
    }
  }
}
	`, ast)
}

func Test239(t *testing.T) {
	ast, err := compile("for(;;);")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ForStatement",
      "init": null,
      "test": null,
      "update": null,
      "body": {
        "type": "EmptyStatement",
        "loc": {
          "start": {
            "line": 1,
            "column": 7
          },
          "end": {
            "line": 1,
            "column": 8
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 8
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 8
    }
  }
}
	`, ast)
}

func Test240(t *testing.T) {
	ast, err := compile("for(;;){}")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ForStatement",
      "init": null,
      "test": null,
      "update": null,
      "body": {
        "type": "BlockStatement",
        "body": [],
        "loc": {
          "start": {
            "line": 1,
            "column": 7
          },
          "end": {
            "line": 1,
            "column": 9
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 9
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 9
    }
  }
}
	`, ast)
}

func Test241(t *testing.T) {
	ast, err := compile("for(x = 0;;);")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ForStatement",
      "init": {
        "type": "AssignmentExpression",
        "operator": "=",
        "left": {
          "type": "Identifier",
          "name": "x",
          "loc": {
            "start": {
              "line": 1,
              "column": 4
            },
            "end": {
              "line": 1,
              "column": 5
            }
          }
        },
        "right": {
          "type": "Literal",
          "value": 0,
          "loc": {
            "start": {
              "line": 1,
              "column": 8
            },
            "end": {
              "line": 1,
              "column": 9
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 4
          },
          "end": {
            "line": 1,
            "column": 9
          }
        }
      },
      "test": null,
      "update": null,
      "body": {
        "type": "EmptyStatement",
        "loc": {
          "start": {
            "line": 1,
            "column": 12
          },
          "end": {
            "line": 1,
            "column": 13
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 13
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 13
    }
  }
}
	`, ast)
}

func Test242(t *testing.T) {
	ast, err := compile("for(var x = 0;;);")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ForStatement",
      "init": {
        "type": "VariableDeclaration",
        "declarations": [
          {
            "type": "VariableDeclarator",
            "id": {
              "type": "Identifier",
              "name": "x",
              "loc": {
                "start": {
                  "line": 1,
                  "column": 8
                },
                "end": {
                  "line": 1,
                  "column": 9
                }
              }
            },
            "init": {
              "type": "Literal",
              "value": 0,
              "loc": {
                "start": {
                  "line": 1,
                  "column": 12
                },
                "end": {
                  "line": 1,
                  "column": 13
                }
              }
            },
            "loc": {
              "start": {
                "line": 1,
                "column": 8
              },
              "end": {
                "line": 1,
                "column": 13
              }
            }
          }
        ],
        "kind": "var",
        "loc": {
          "start": {
            "line": 1,
            "column": 4
          },
          "end": {
            "line": 1,
            "column": 13
          }
        }
      },
      "test": null,
      "update": null,
      "body": {
        "type": "EmptyStatement",
        "loc": {
          "start": {
            "line": 1,
            "column": 16
          },
          "end": {
            "line": 1,
            "column": 17
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 17
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 17
    }
  }
}
	`, ast)
}

func Test243(t *testing.T) {
	ast, err := compile("for(var x = 0, y = 1;;);")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ForStatement",
      "init": {
        "type": "VariableDeclaration",
        "declarations": [
          {
            "type": "VariableDeclarator",
            "id": {
              "type": "Identifier",
              "name": "x",
              "loc": {
                "start": {
                  "line": 1,
                  "column": 8
                },
                "end": {
                  "line": 1,
                  "column": 9
                }
              }
            },
            "init": {
              "type": "Literal",
              "value": 0,
              "loc": {
                "start": {
                  "line": 1,
                  "column": 12
                },
                "end": {
                  "line": 1,
                  "column": 13
                }
              }
            },
            "loc": {
              "start": {
                "line": 1,
                "column": 8
              },
              "end": {
                "line": 1,
                "column": 13
              }
            }
          },
          {
            "type": "VariableDeclarator",
            "id": {
              "type": "Identifier",
              "name": "y",
              "loc": {
                "start": {
                  "line": 1,
                  "column": 15
                },
                "end": {
                  "line": 1,
                  "column": 16
                }
              }
            },
            "init": {
              "type": "Literal",
              "value": 1,
              "loc": {
                "start": {
                  "line": 1,
                  "column": 19
                },
                "end": {
                  "line": 1,
                  "column": 20
                }
              }
            },
            "loc": {
              "start": {
                "line": 1,
                "column": 15
              },
              "end": {
                "line": 1,
                "column": 20
              }
            }
          }
        ],
        "kind": "var",
        "loc": {
          "start": {
            "line": 1,
            "column": 4
          },
          "end": {
            "line": 1,
            "column": 20
          }
        }
      },
      "test": null,
      "update": null,
      "body": {
        "type": "EmptyStatement",
        "loc": {
          "start": {
            "line": 1,
            "column": 23
          },
          "end": {
            "line": 1,
            "column": 24
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 24
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 24
    }
  }
}
	`, ast)
}

func Test244(t *testing.T) {
	ast, err := compile("for(x = 0; x < 42;);")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ForStatement",
      "init": {
        "type": "AssignmentExpression",
        "operator": "=",
        "left": {
          "type": "Identifier",
          "name": "x",
          "loc": {
            "start": {
              "line": 1,
              "column": 4
            },
            "end": {
              "line": 1,
              "column": 5
            }
          }
        },
        "right": {
          "type": "Literal",
          "value": 0,
          "loc": {
            "start": {
              "line": 1,
              "column": 8
            },
            "end": {
              "line": 1,
              "column": 9
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 4
          },
          "end": {
            "line": 1,
            "column": 9
          }
        }
      },
      "test": {
        "type": "BinaryExpression",
        "left": {
          "type": "Identifier",
          "name": "x",
          "loc": {
            "start": {
              "line": 1,
              "column": 11
            },
            "end": {
              "line": 1,
              "column": 12
            }
          }
        },
        "operator": "<",
        "right": {
          "type": "Literal",
          "value": 42,
          "loc": {
            "start": {
              "line": 1,
              "column": 15
            },
            "end": {
              "line": 1,
              "column": 17
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 11
          },
          "end": {
            "line": 1,
            "column": 17
          }
        }
      },
      "update": null,
      "body": {
        "type": "EmptyStatement",
        "loc": {
          "start": {
            "line": 1,
            "column": 19
          },
          "end": {
            "line": 1,
            "column": 20
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 20
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 20
    }
  }
}
	`, ast)
}

func Test245(t *testing.T) {
	ast, err := compile("for(x = 0; x < 42; x++);")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ForStatement",
      "init": {
        "type": "AssignmentExpression",
        "operator": "=",
        "left": {
          "type": "Identifier",
          "name": "x",
          "loc": {
            "start": {
              "line": 1,
              "column": 4
            },
            "end": {
              "line": 1,
              "column": 5
            }
          }
        },
        "right": {
          "type": "Literal",
          "value": 0,
          "loc": {
            "start": {
              "line": 1,
              "column": 8
            },
            "end": {
              "line": 1,
              "column": 9
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 4
          },
          "end": {
            "line": 1,
            "column": 9
          }
        }
      },
      "test": {
        "type": "BinaryExpression",
        "left": {
          "type": "Identifier",
          "name": "x",
          "loc": {
            "start": {
              "line": 1,
              "column": 11
            },
            "end": {
              "line": 1,
              "column": 12
            }
          }
        },
        "operator": "<",
        "right": {
          "type": "Literal",
          "value": 42,
          "loc": {
            "start": {
              "line": 1,
              "column": 15
            },
            "end": {
              "line": 1,
              "column": 17
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 11
          },
          "end": {
            "line": 1,
            "column": 17
          }
        }
      },
      "update": {
        "type": "UpdateExpression",
        "operator": "++",
        "prefix": false,
        "argument": {
          "type": "Identifier",
          "name": "x",
          "loc": {
            "start": {
              "line": 1,
              "column": 19
            },
            "end": {
              "line": 1,
              "column": 20
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 19
          },
          "end": {
            "line": 1,
            "column": 22
          }
        }
      },
      "body": {
        "type": "EmptyStatement",
        "loc": {
          "start": {
            "line": 1,
            "column": 23
          },
          "end": {
            "line": 1,
            "column": 24
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 24
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 24
    }
  }
}
	`, ast)
}

func Test246(t *testing.T) {
	ast, err := compile("for(x = 0; x < 42; x++) process(x);")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ForStatement",
      "init": {
        "type": "AssignmentExpression",
        "operator": "=",
        "left": {
          "type": "Identifier",
          "name": "x",
          "loc": {
            "start": {
              "line": 1,
              "column": 4
            },
            "end": {
              "line": 1,
              "column": 5
            }
          }
        },
        "right": {
          "type": "Literal",
          "value": 0,
          "loc": {
            "start": {
              "line": 1,
              "column": 8
            },
            "end": {
              "line": 1,
              "column": 9
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 4
          },
          "end": {
            "line": 1,
            "column": 9
          }
        }
      },
      "test": {
        "type": "BinaryExpression",
        "left": {
          "type": "Identifier",
          "name": "x",
          "loc": {
            "start": {
              "line": 1,
              "column": 11
            },
            "end": {
              "line": 1,
              "column": 12
            }
          }
        },
        "operator": "<",
        "right": {
          "type": "Literal",
          "value": 42,
          "loc": {
            "start": {
              "line": 1,
              "column": 15
            },
            "end": {
              "line": 1,
              "column": 17
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 11
          },
          "end": {
            "line": 1,
            "column": 17
          }
        }
      },
      "update": {
        "type": "UpdateExpression",
        "operator": "++",
        "prefix": false,
        "argument": {
          "type": "Identifier",
          "name": "x",
          "loc": {
            "start": {
              "line": 1,
              "column": 19
            },
            "end": {
              "line": 1,
              "column": 20
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 19
          },
          "end": {
            "line": 1,
            "column": 22
          }
        }
      },
      "body": {
        "type": "ExpressionStatement",
        "expression": {
          "type": "CallExpression",
          "callee": {
            "type": "Identifier",
            "name": "process",
            "loc": {
              "start": {
                "line": 1,
                "column": 24
              },
              "end": {
                "line": 1,
                "column": 31
              }
            }
          },
          "arguments": [
            {
              "type": "Identifier",
              "name": "x",
              "loc": {
                "start": {
                  "line": 1,
                  "column": 32
                },
                "end": {
                  "line": 1,
                  "column": 33
                }
              }
            }
          ],
          "loc": {
            "start": {
              "line": 1,
              "column": 24
            },
            "end": {
              "line": 1,
              "column": 34
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 24
          },
          "end": {
            "line": 1,
            "column": 35
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 35
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 35
    }
  }
}
	`, ast)
}

func Test247(t *testing.T) {
	ast, err := compile("for(x in list) process(x);")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ForInStatement",
      "left": {
        "type": "Identifier",
        "name": "x",
        "loc": {
          "start": {
            "line": 1,
            "column": 4
          },
          "end": {
            "line": 1,
            "column": 5
          }
        }
      },
      "right": {
        "type": "Identifier",
        "name": "list",
        "loc": {
          "start": {
            "line": 1,
            "column": 9
          },
          "end": {
            "line": 1,
            "column": 13
          }
        }
      },
      "body": {
        "type": "ExpressionStatement",
        "expression": {
          "type": "CallExpression",
          "callee": {
            "type": "Identifier",
            "name": "process",
            "loc": {
              "start": {
                "line": 1,
                "column": 15
              },
              "end": {
                "line": 1,
                "column": 22
              }
            }
          },
          "arguments": [
            {
              "type": "Identifier",
              "name": "x",
              "loc": {
                "start": {
                  "line": 1,
                  "column": 23
                },
                "end": {
                  "line": 1,
                  "column": 24
                }
              }
            }
          ],
          "loc": {
            "start": {
              "line": 1,
              "column": 15
            },
            "end": {
              "line": 1,
              "column": 25
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 15
          },
          "end": {
            "line": 1,
            "column": 26
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 26
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 26
    }
  }
}
	`, ast)
}

func Test248(t *testing.T) {
	ast, err := compile("for (var x in list) process(x);")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ForInStatement",
      "left": {
        "type": "VariableDeclaration",
        "declarations": [
          {
            "type": "VariableDeclarator",
            "id": {
              "type": "Identifier",
              "name": "x",
              "loc": {
                "start": {
                  "line": 1,
                  "column": 9
                },
                "end": {
                  "line": 1,
                  "column": 10
                }
              }
            },
            "init": null,
            "loc": {
              "start": {
                "line": 1,
                "column": 9
              },
              "end": {
                "line": 1,
                "column": 10
              }
            }
          }
        ],
        "kind": "var",
        "loc": {
          "start": {
            "line": 1,
            "column": 5
          },
          "end": {
            "line": 1,
            "column": 10
          }
        }
      },
      "right": {
        "type": "Identifier",
        "name": "list",
        "loc": {
          "start": {
            "line": 1,
            "column": 14
          },
          "end": {
            "line": 1,
            "column": 18
          }
        }
      },
      "body": {
        "type": "ExpressionStatement",
        "expression": {
          "type": "CallExpression",
          "callee": {
            "type": "Identifier",
            "name": "process",
            "loc": {
              "start": {
                "line": 1,
                "column": 20
              },
              "end": {
                "line": 1,
                "column": 27
              }
            }
          },
          "arguments": [
            {
              "type": "Identifier",
              "name": "x",
              "loc": {
                "start": {
                  "line": 1,
                  "column": 28
                },
                "end": {
                  "line": 1,
                  "column": 29
                }
              }
            }
          ],
          "loc": {
            "start": {
              "line": 1,
              "column": 20
            },
            "end": {
              "line": 1,
              "column": 30
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 20
          },
          "end": {
            "line": 1,
            "column": 31
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 31
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 31
    }
  }
}
	`, ast)
}

func Test249(t *testing.T) {
	ast, err := compile("for (var x = 42 in list) process(x);")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ForInStatement",
      "left": {
        "type": "VariableDeclaration",
        "declarations": [
          {
            "type": "VariableDeclarator",
            "id": {
              "type": "Identifier",
              "name": "x",
              "loc": {
                "start": {
                  "line": 1,
                  "column": 9
                },
                "end": {
                  "line": 1,
                  "column": 10
                }
              }
            },
            "init": {
              "type": "Literal",
              "value": 42,
              "loc": {
                "start": {
                  "line": 1,
                  "column": 13
                },
                "end": {
                  "line": 1,
                  "column": 15
                }
              }
            },
            "loc": {
              "start": {
                "line": 1,
                "column": 9
              },
              "end": {
                "line": 1,
                "column": 15
              }
            }
          }
        ],
        "kind": "var",
        "loc": {
          "start": {
            "line": 1,
            "column": 5
          },
          "end": {
            "line": 1,
            "column": 15
          }
        }
      },
      "right": {
        "type": "Identifier",
        "name": "list",
        "loc": {
          "start": {
            "line": 1,
            "column": 19
          },
          "end": {
            "line": 1,
            "column": 23
          }
        }
      },
      "body": {
        "type": "ExpressionStatement",
        "expression": {
          "type": "CallExpression",
          "callee": {
            "type": "Identifier",
            "name": "process",
            "loc": {
              "start": {
                "line": 1,
                "column": 25
              },
              "end": {
                "line": 1,
                "column": 32
              }
            }
          },
          "arguments": [
            {
              "type": "Identifier",
              "name": "x",
              "loc": {
                "start": {
                  "line": 1,
                  "column": 33
                },
                "end": {
                  "line": 1,
                  "column": 34
                }
              }
            }
          ],
          "loc": {
            "start": {
              "line": 1,
              "column": 25
            },
            "end": {
              "line": 1,
              "column": 35
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 25
          },
          "end": {
            "line": 1,
            "column": 36
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 36
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 36
    }
  }
}
	`, ast)
}

func Test250(t *testing.T) {
	ast, err := compile("for (var i = function() { return 10 in [] } in list) process(x);")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ForInStatement",
      "left": {
        "type": "VariableDeclaration",
        "declarations": [
          {
            "type": "VariableDeclarator",
            "id": {
              "type": "Identifier",
              "name": "i",
              "loc": {
                "start": {
                  "line": 1,
                  "column": 9
                },
                "end": {
                  "line": 1,
                  "column": 10
                }
              }
            },
            "init": {
              "type": "FunctionExpression",
              "id": null,
              "params": [],
              "body": {
                "type": "BlockStatement",
                "body": [
                  {
                    "type": "ReturnStatement",
                    "argument": {
                      "type": "BinaryExpression",
                      "left": {
                        "type": "Literal",
                        "value": 10,
                        "loc": {
                          "start": {
                            "line": 1,
                            "column": 33
                          },
                          "end": {
                            "line": 1,
                            "column": 35
                          }
                        }
                      },
                      "operator": "in",
                      "right": {
                        "type": "ArrayExpression",
                        "elements": [],
                        "loc": {
                          "start": {
                            "line": 1,
                            "column": 39
                          },
                          "end": {
                            "line": 1,
                            "column": 41
                          }
                        }
                      },
                      "loc": {
                        "start": {
                          "line": 1,
                          "column": 33
                        },
                        "end": {
                          "line": 1,
                          "column": 41
                        }
                      }
                    },
                    "loc": {
                      "start": {
                        "line": 1,
                        "column": 26
                      },
                      "end": {
                        "line": 1,
                        "column": 41
                      }
                    }
                  }
                ],
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 24
                  },
                  "end": {
                    "line": 1,
                    "column": 43
                  }
                }
              },
              "loc": {
                "start": {
                  "line": 1,
                  "column": 13
                },
                "end": {
                  "line": 1,
                  "column": 43
                }
              }
            },
            "loc": {
              "start": {
                "line": 1,
                "column": 9
              },
              "end": {
                "line": 1,
                "column": 43
              }
            }
          }
        ],
        "kind": "var",
        "loc": {
          "start": {
            "line": 1,
            "column": 5
          },
          "end": {
            "line": 1,
            "column": 43
          }
        }
      },
      "right": {
        "type": "Identifier",
        "name": "list",
        "loc": {
          "start": {
            "line": 1,
            "column": 47
          },
          "end": {
            "line": 1,
            "column": 51
          }
        }
      },
      "body": {
        "type": "ExpressionStatement",
        "expression": {
          "type": "CallExpression",
          "callee": {
            "type": "Identifier",
            "name": "process",
            "loc": {
              "start": {
                "line": 1,
                "column": 53
              },
              "end": {
                "line": 1,
                "column": 60
              }
            }
          },
          "arguments": [
            {
              "type": "Identifier",
              "name": "x",
              "loc": {
                "start": {
                  "line": 1,
                  "column": 61
                },
                "end": {
                  "line": 1,
                  "column": 62
                }
              }
            }
          ],
          "loc": {
            "start": {
              "line": 1,
              "column": 53
            },
            "end": {
              "line": 1,
              "column": 63
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 53
          },
          "end": {
            "line": 1,
            "column": 64
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 64
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 64
    }
  }
}
	`, ast)
}

func Test251(t *testing.T) {
	ast, err := compile("while (true) { continue; }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "WhileStatement",
      "test": {
        "type": "Literal",
        "value": true,
        "loc": {
          "start": {
            "line": 1,
            "column": 7
          },
          "end": {
            "line": 1,
            "column": 11
          }
        }
      },
      "body": {
        "type": "BlockStatement",
        "body": [
          {
            "type": "ContinueStatement",
            "label": null,
            "loc": {
              "start": {
                "line": 1,
                "column": 15
              },
              "end": {
                "line": 1,
                "column": 24
              }
            }
          }
        ],
        "loc": {
          "start": {
            "line": 1,
            "column": 13
          },
          "end": {
            "line": 1,
            "column": 26
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 26
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 26
    }
  }
}
	`, ast)
}

func Test252(t *testing.T) {
	ast, err := compile("while (true) { continue }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "WhileStatement",
      "test": {
        "type": "Literal",
        "value": true,
        "loc": {
          "start": {
            "line": 1,
            "column": 7
          },
          "end": {
            "line": 1,
            "column": 11
          }
        }
      },
      "body": {
        "type": "BlockStatement",
        "body": [
          {
            "type": "ContinueStatement",
            "label": null,
            "loc": {
              "start": {
                "line": 1,
                "column": 15
              },
              "end": {
                "line": 1,
                "column": 23
              }
            }
          }
        ],
        "loc": {
          "start": {
            "line": 1,
            "column": 13
          },
          "end": {
            "line": 1,
            "column": 25
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 25
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 25
    }
  }
}
	`, ast)
}

func Test253(t *testing.T) {
	ast, err := compile("done: while (true) { continue done }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "LabeledStatement",
      "body": {
        "type": "WhileStatement",
        "test": {
          "type": "Literal",
          "value": true,
          "loc": {
            "start": {
              "line": 1,
              "column": 13
            },
            "end": {
              "line": 1,
              "column": 17
            }
          }
        },
        "body": {
          "type": "BlockStatement",
          "body": [
            {
              "type": "ContinueStatement",
              "label": {
                "type": "Identifier",
                "name": "done",
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 30
                  },
                  "end": {
                    "line": 1,
                    "column": 34
                  }
                }
              },
              "loc": {
                "start": {
                  "line": 1,
                  "column": 21
                },
                "end": {
                  "line": 1,
                  "column": 34
                }
              }
            }
          ],
          "loc": {
            "start": {
              "line": 1,
              "column": 19
            },
            "end": {
              "line": 1,
              "column": 36
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 6
          },
          "end": {
            "line": 1,
            "column": 36
          }
        }
      },
      "label": {
        "type": "Identifier",
        "name": "done",
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 4
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 36
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 36
    }
  }
}
	`, ast)
}

func Test254(t *testing.T) {
	ast, err := compile("done: while (true) { continue done; }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "LabeledStatement",
      "body": {
        "type": "WhileStatement",
        "test": {
          "type": "Literal",
          "value": true,
          "loc": {
            "start": {
              "line": 1,
              "column": 13
            },
            "end": {
              "line": 1,
              "column": 17
            }
          }
        },
        "body": {
          "type": "BlockStatement",
          "body": [
            {
              "type": "ContinueStatement",
              "label": {
                "type": "Identifier",
                "name": "done",
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 30
                  },
                  "end": {
                    "line": 1,
                    "column": 34
                  }
                }
              },
              "loc": {
                "start": {
                  "line": 1,
                  "column": 21
                },
                "end": {
                  "line": 1,
                  "column": 35
                }
              }
            }
          ],
          "loc": {
            "start": {
              "line": 1,
              "column": 19
            },
            "end": {
              "line": 1,
              "column": 37
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 6
          },
          "end": {
            "line": 1,
            "column": 37
          }
        }
      },
      "label": {
        "type": "Identifier",
        "name": "done",
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 4
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 37
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 37
    }
  }
}
	`, ast)
}

func Test255(t *testing.T) {
	ast, err := compile("while (true) { break }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "WhileStatement",
      "test": {
        "type": "Literal",
        "value": true,
        "loc": {
          "start": {
            "line": 1,
            "column": 7
          },
          "end": {
            "line": 1,
            "column": 11
          }
        }
      },
      "body": {
        "type": "BlockStatement",
        "body": [
          {
            "type": "BreakStatement",
            "label": null,
            "loc": {
              "start": {
                "line": 1,
                "column": 15
              },
              "end": {
                "line": 1,
                "column": 20
              }
            }
          }
        ],
        "loc": {
          "start": {
            "line": 1,
            "column": 13
          },
          "end": {
            "line": 1,
            "column": 22
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 22
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 22
    }
  }
}
	`, ast)
}

func Test256(t *testing.T) {
	ast, err := compile("done: while (true) { break done }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "LabeledStatement",
      "body": {
        "type": "WhileStatement",
        "test": {
          "type": "Literal",
          "value": true,
          "loc": {
            "start": {
              "line": 1,
              "column": 13
            },
            "end": {
              "line": 1,
              "column": 17
            }
          }
        },
        "body": {
          "type": "BlockStatement",
          "body": [
            {
              "type": "BreakStatement",
              "label": {
                "type": "Identifier",
                "name": "done",
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 27
                  },
                  "end": {
                    "line": 1,
                    "column": 31
                  }
                }
              },
              "loc": {
                "start": {
                  "line": 1,
                  "column": 21
                },
                "end": {
                  "line": 1,
                  "column": 31
                }
              }
            }
          ],
          "loc": {
            "start": {
              "line": 1,
              "column": 19
            },
            "end": {
              "line": 1,
              "column": 33
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 6
          },
          "end": {
            "line": 1,
            "column": 33
          }
        }
      },
      "label": {
        "type": "Identifier",
        "name": "done",
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 4
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 33
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 33
    }
  }
}
	`, ast)
}

func Test257(t *testing.T) {
	ast, err := compile("done: while (true) { break done; }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "LabeledStatement",
      "body": {
        "type": "WhileStatement",
        "test": {
          "type": "Literal",
          "value": true,
          "loc": {
            "start": {
              "line": 1,
              "column": 13
            },
            "end": {
              "line": 1,
              "column": 17
            }
          }
        },
        "body": {
          "type": "BlockStatement",
          "body": [
            {
              "type": "BreakStatement",
              "label": {
                "type": "Identifier",
                "name": "done",
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 27
                  },
                  "end": {
                    "line": 1,
                    "column": 31
                  }
                }
              },
              "loc": {
                "start": {
                  "line": 1,
                  "column": 21
                },
                "end": {
                  "line": 1,
                  "column": 32
                }
              }
            }
          ],
          "loc": {
            "start": {
              "line": 1,
              "column": 19
            },
            "end": {
              "line": 1,
              "column": 34
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 6
          },
          "end": {
            "line": 1,
            "column": 34
          }
        }
      },
      "label": {
        "type": "Identifier",
        "name": "done",
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 4
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 34
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 34
    }
  }
}
	`, ast)
}

func Test258(t *testing.T) {
	ast, err := compile("done: switch (a) { default: break done }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "LabeledStatement",
      "body": {
        "type": "SwitchStatement",
        "discriminant": {
          "type": "Identifier",
          "loc": {
            "start": {
              "line": 1,
              "column": 14
            }
          }
        },
        "cases": [
          {
            "type": "SwitchCase",
            "consequent": [
              {
                "type": "BreakStatement",
                "label": {
                  "type": "Identifier",
                  "name": "done"
                }
              }
            ],
            "test": null
          }
        ]
      },
      "label": {
        "type": "Identifier",
        "name": "done"
      }
    }
  ]
}
	`, ast)
}

func Test259(t *testing.T) {
	ast, err := compile("(function(){ return })")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "FunctionExpression",
        "id": null,
        "params": [],
        "body": {
          "type": "BlockStatement",
          "body": [
            {
              "type": "ReturnStatement",
              "argument": null,
              "loc": {
                "start": {
                  "line": 1,
                  "column": 13
                },
                "end": {
                  "line": 1,
                  "column": 19
                }
              }
            }
          ],
          "loc": {
            "start": {
              "line": 1,
              "column": 11
            },
            "end": {
              "line": 1,
              "column": 21
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 1
          },
          "end": {
            "line": 1,
            "column": 21
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 22
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 22
    }
  }
}
	`, ast)
}

func Test260(t *testing.T) {
	ast, err := compile("(function(){ return; })")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "FunctionExpression",
        "id": null,
        "params": [],
        "body": {
          "type": "BlockStatement",
          "body": [
            {
              "type": "ReturnStatement",
              "argument": null,
              "loc": {
                "start": {
                  "line": 1,
                  "column": 13
                },
                "end": {
                  "line": 1,
                  "column": 20
                }
              }
            }
          ],
          "loc": {
            "start": {
              "line": 1,
              "column": 11
            },
            "end": {
              "line": 1,
              "column": 22
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 1
          },
          "end": {
            "line": 1,
            "column": 22
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 23
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 23
    }
  }
}
	`, ast)
}

func Test261(t *testing.T) {
	ast, err := compile("(function(){ return x; })")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "FunctionExpression",
        "id": null,
        "params": [],
        "body": {
          "type": "BlockStatement",
          "body": [
            {
              "type": "ReturnStatement",
              "argument": {
                "type": "Identifier",
                "name": "x",
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 20
                  },
                  "end": {
                    "line": 1,
                    "column": 21
                  }
                }
              },
              "loc": {
                "start": {
                  "line": 1,
                  "column": 13
                },
                "end": {
                  "line": 1,
                  "column": 22
                }
              }
            }
          ],
          "loc": {
            "start": {
              "line": 1,
              "column": 11
            },
            "end": {
              "line": 1,
              "column": 24
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 1
          },
          "end": {
            "line": 1,
            "column": 24
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 25
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 25
    }
  }
}
	`, ast)
}

func Test262(t *testing.T) {
	ast, err := compile("(function(){ return x * y })")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "FunctionExpression",
        "id": null,
        "params": [],
        "body": {
          "type": "BlockStatement",
          "body": [
            {
              "type": "ReturnStatement",
              "argument": {
                "type": "BinaryExpression",
                "left": {
                  "type": "Identifier",
                  "name": "x",
                  "loc": {
                    "start": {
                      "line": 1,
                      "column": 20
                    },
                    "end": {
                      "line": 1,
                      "column": 21
                    }
                  }
                },
                "operator": "*",
                "right": {
                  "type": "Identifier",
                  "name": "y",
                  "loc": {
                    "start": {
                      "line": 1,
                      "column": 24
                    },
                    "end": {
                      "line": 1,
                      "column": 25
                    }
                  }
                },
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 20
                  },
                  "end": {
                    "line": 1,
                    "column": 25
                  }
                }
              },
              "loc": {
                "start": {
                  "line": 1,
                  "column": 13
                },
                "end": {
                  "line": 1,
                  "column": 25
                }
              }
            }
          ],
          "loc": {
            "start": {
              "line": 1,
              "column": 11
            },
            "end": {
              "line": 1,
              "column": 27
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 1
          },
          "end": {
            "line": 1,
            "column": 27
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 28
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 28
    }
  }
}
	`, ast)
}

func Test263(t *testing.T) {
	ast, err := compile("with (x) foo = bar")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "WithStatement",
      "object": {
        "type": "Identifier",
        "name": "x",
        "loc": {
          "start": {
            "line": 1,
            "column": 6
          },
          "end": {
            "line": 1,
            "column": 7
          }
        }
      },
      "body": {
        "type": "ExpressionStatement",
        "expression": {
          "type": "AssignmentExpression",
          "operator": "=",
          "left": {
            "type": "Identifier",
            "name": "foo",
            "loc": {
              "start": {
                "line": 1,
                "column": 9
              },
              "end": {
                "line": 1,
                "column": 12
              }
            }
          },
          "right": {
            "type": "Identifier",
            "name": "bar",
            "loc": {
              "start": {
                "line": 1,
                "column": 15
              },
              "end": {
                "line": 1,
                "column": 18
              }
            }
          },
          "loc": {
            "start": {
              "line": 1,
              "column": 9
            },
            "end": {
              "line": 1,
              "column": 18
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 9
          },
          "end": {
            "line": 1,
            "column": 18
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 18
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 18
    }
  }
}
	`, ast)
}

func Test264(t *testing.T) {
	ast, err := compile("with (x) foo = bar;")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "WithStatement",
      "object": {
        "type": "Identifier",
        "name": "x",
        "loc": {
          "start": {
            "line": 1,
            "column": 6
          },
          "end": {
            "line": 1,
            "column": 7
          }
        }
      },
      "body": {
        "type": "ExpressionStatement",
        "expression": {
          "type": "AssignmentExpression",
          "operator": "=",
          "left": {
            "type": "Identifier",
            "name": "foo",
            "loc": {
              "start": {
                "line": 1,
                "column": 9
              },
              "end": {
                "line": 1,
                "column": 12
              }
            }
          },
          "right": {
            "type": "Identifier",
            "name": "bar",
            "loc": {
              "start": {
                "line": 1,
                "column": 15
              },
              "end": {
                "line": 1,
                "column": 18
              }
            }
          },
          "loc": {
            "start": {
              "line": 1,
              "column": 9
            },
            "end": {
              "line": 1,
              "column": 18
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 9
          },
          "end": {
            "line": 1,
            "column": 19
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 19
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 19
    }
  }
}
	`, ast)
}

func Test265(t *testing.T) {
	ast, err := compile("with (x) { foo = bar }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "WithStatement",
      "object": {
        "type": "Identifier",
        "name": "x",
        "loc": {
          "start": {
            "line": 1,
            "column": 6
          },
          "end": {
            "line": 1,
            "column": 7
          }
        }
      },
      "body": {
        "type": "BlockStatement",
        "body": [
          {
            "type": "ExpressionStatement",
            "expression": {
              "type": "AssignmentExpression",
              "operator": "=",
              "left": {
                "type": "Identifier",
                "name": "foo",
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 11
                  },
                  "end": {
                    "line": 1,
                    "column": 14
                  }
                }
              },
              "right": {
                "type": "Identifier",
                "name": "bar",
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 17
                  },
                  "end": {
                    "line": 1,
                    "column": 20
                  }
                }
              },
              "loc": {
                "start": {
                  "line": 1,
                  "column": 11
                },
                "end": {
                  "line": 1,
                  "column": 20
                }
              }
            },
            "loc": {
              "start": {
                "line": 1,
                "column": 11
              },
              "end": {
                "line": 1,
                "column": 20
              }
            }
          }
        ],
        "loc": {
          "start": {
            "line": 1,
            "column": 9
          },
          "end": {
            "line": 1,
            "column": 22
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 22
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 22
    }
  }
}
	`, ast)
}

func Test266(t *testing.T) {
	ast, err := compile("switch (x) {}")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "SwitchStatement",
      "discriminant": {
        "type": "Identifier",
        "name": "x",
        "loc": {
          "start": {
            "line": 1,
            "column": 8
          },
          "end": {
            "line": 1,
            "column": 9
          }
        }
      },
      "cases": [],
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 13
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 13
    }
  }
}
	`, ast)
}

func Test267(t *testing.T) {
	ast, err := compile("switch (answer) { case 42: hi(); break; }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "SwitchStatement",
      "discriminant": {
        "type": "Identifier",
        "name": "answer",
        "loc": {
          "start": {
            "line": 1,
            "column": 8
          },
          "end": {
            "line": 1,
            "column": 14
          }
        }
      },
      "cases": [
        {
          "type": "SwitchCase",
          "consequent": [
            {
              "type": "ExpressionStatement",
              "expression": {
                "type": "CallExpression",
                "callee": {
                  "type": "Identifier",
                  "name": "hi",
                  "loc": {
                    "start": {
                      "line": 1,
                      "column": 27
                    },
                    "end": {
                      "line": 1,
                      "column": 29
                    }
                  }
                },
                "arguments": [],
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 27
                  },
                  "end": {
                    "line": 1,
                    "column": 31
                  }
                }
              },
              "loc": {
                "start": {
                  "line": 1,
                  "column": 27
                },
                "end": {
                  "line": 1,
                  "column": 32
                }
              }
            },
            {
              "type": "BreakStatement",
              "label": null,
              "loc": {
                "start": {
                  "line": 1,
                  "column": 33
                },
                "end": {
                  "line": 1,
                  "column": 39
                }
              }
            }
          ],
          "test": {
            "type": "Literal",
            "value": 42,
            "loc": {
              "start": {
                "line": 1,
                "column": 23
              },
              "end": {
                "line": 1,
                "column": 25
              }
            }
          },
          "loc": {
            "start": {
              "line": 1,
              "column": 18
            },
            "end": {
              "line": 1,
              "column": 39
            }
          }
        }
      ],
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 41
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 41
    }
  }
}
	`, ast)
}

func Test268(t *testing.T) {
	ast, err := compile("switch (answer) { case 42: hi(); break; default: break }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "SwitchStatement",
      "discriminant": {
        "type": "Identifier",
        "name": "answer",
        "loc": {
          "start": {
            "line": 1,
            "column": 8
          },
          "end": {
            "line": 1,
            "column": 14
          }
        }
      },
      "cases": [
        {
          "type": "SwitchCase",
          "consequent": [
            {
              "type": "ExpressionStatement",
              "expression": {
                "type": "CallExpression",
                "callee": {
                  "type": "Identifier",
                  "name": "hi",
                  "loc": {
                    "start": {
                      "line": 1,
                      "column": 27
                    },
                    "end": {
                      "line": 1,
                      "column": 29
                    }
                  }
                },
                "arguments": [],
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 27
                  },
                  "end": {
                    "line": 1,
                    "column": 31
                  }
                }
              },
              "loc": {
                "start": {
                  "line": 1,
                  "column": 27
                },
                "end": {
                  "line": 1,
                  "column": 32
                }
              }
            },
            {
              "type": "BreakStatement",
              "label": null,
              "loc": {
                "start": {
                  "line": 1,
                  "column": 33
                },
                "end": {
                  "line": 1,
                  "column": 39
                }
              }
            }
          ],
          "test": {
            "type": "Literal",
            "value": 42,
            "loc": {
              "start": {
                "line": 1,
                "column": 23
              },
              "end": {
                "line": 1,
                "column": 25
              }
            }
          },
          "loc": {
            "start": {
              "line": 1,
              "column": 18
            },
            "end": {
              "line": 1,
              "column": 39
            }
          }
        },
        {
          "type": "SwitchCase",
          "consequent": [
            {
              "type": "BreakStatement",
              "label": null,
              "loc": {
                "start": {
                  "line": 1,
                  "column": 49
                },
                "end": {
                  "line": 1,
                  "column": 54
                }
              }
            }
          ],
          "test": null,
          "loc": {
            "start": {
              "line": 1,
              "column": 40
            },
            "end": {
              "line": 1,
              "column": 54
            }
          }
        }
      ],
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 56
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 56
    }
  }
}
	`, ast)
}

func Test269(t *testing.T) {
	ast, err := compile("start: for (;;) break start")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "LabeledStatement",
      "body": {
        "type": "ForStatement",
        "init": null,
        "test": null,
        "update": null,
        "body": {
          "type": "BreakStatement",
          "label": {
            "type": "Identifier",
            "name": "start",
            "loc": {
              "start": {
                "line": 1,
                "column": 22
              },
              "end": {
                "line": 1,
                "column": 27
              }
            }
          },
          "loc": {
            "start": {
              "line": 1,
              "column": 16
            },
            "end": {
              "line": 1,
              "column": 27
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 7
          },
          "end": {
            "line": 1,
            "column": 27
          }
        }
      },
      "label": {
        "type": "Identifier",
        "name": "start",
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 5
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 27
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 27
    }
  }
}
	`, ast)
}

func Test270(t *testing.T) {
	ast, err := compile("start: while (true) break start")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "LabeledStatement",
      "body": {
        "type": "WhileStatement",
        "test": {
          "type": "Literal",
          "value": true,
          "loc": {
            "start": {
              "line": 1,
              "column": 14
            },
            "end": {
              "line": 1,
              "column": 18
            }
          }
        },
        "body": {
          "type": "BreakStatement",
          "label": {
            "type": "Identifier",
            "name": "start",
            "loc": {
              "start": {
                "line": 1,
                "column": 26
              },
              "end": {
                "line": 1,
                "column": 31
              }
            }
          },
          "loc": {
            "start": {
              "line": 1,
              "column": 20
            },
            "end": {
              "line": 1,
              "column": 31
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 7
          },
          "end": {
            "line": 1,
            "column": 31
          }
        }
      },
      "label": {
        "type": "Identifier",
        "name": "start",
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 5
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 31
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 31
    }
  }
}
	`, ast)
}

func Test271(t *testing.T) {
	ast, err := compile("throw x;")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ThrowStatement",
      "argument": {
        "type": "Identifier",
        "name": "x",
        "loc": {
          "start": {
            "line": 1,
            "column": 6
          },
          "end": {
            "line": 1,
            "column": 7
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 8
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 8
    }
  }
}
	`, ast)
}

func Test272(t *testing.T) {
	ast, err := compile("throw x * y")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ThrowStatement",
      "argument": {
        "type": "BinaryExpression",
        "left": {
          "type": "Identifier",
          "name": "x",
          "loc": {
            "start": {
              "line": 1,
              "column": 6
            },
            "end": {
              "line": 1,
              "column": 7
            }
          }
        },
        "operator": "*",
        "right": {
          "type": "Identifier",
          "name": "y",
          "loc": {
            "start": {
              "line": 1,
              "column": 10
            },
            "end": {
              "line": 1,
              "column": 11
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 6
          },
          "end": {
            "line": 1,
            "column": 11
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 11
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 11
    }
  }
}
	`, ast)
}

func Test273(t *testing.T) {
	ast, err := compile("throw { message: \"Error\" }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ThrowStatement",
      "argument": {
        "type": "ObjectExpression",
        "properties": [
          {
            "type": "Property",
            "key": {
              "type": "Identifier",
              "name": "message",
              "loc": {
                "start": {
                  "line": 1,
                  "column": 8
                },
                "end": {
                  "line": 1,
                  "column": 15
                }
              }
            },
            "value": {
              "type": "Literal",
              "value": "Error",
              "loc": {
                "start": {
                  "line": 1,
                  "column": 17
                },
                "end": {
                  "line": 1,
                  "column": 24
                }
              }
            },
            "kind": "init"
          }
        ],
        "loc": {
          "start": {
            "line": 1,
            "column": 6
          },
          "end": {
            "line": 1,
            "column": 26
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 26
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 26
    }
  }
}
	`, ast)
}

func Test274(t *testing.T) {
	ast, err := compile("try { } catch (e) { }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "TryStatement",
      "block": {
        "type": "BlockStatement",
        "body": [],
        "loc": {
          "start": {
            "line": 1,
            "column": 4
          },
          "end": {
            "line": 1,
            "column": 7
          }
        }
      },
      "handler": {
        "type": "CatchClause",
        "param": {
          "type": "Identifier",
          "name": "e",
          "loc": {
            "start": {
              "line": 1,
              "column": 15
            },
            "end": {
              "line": 1,
              "column": 16
            }
          }
        },
        "body": {
          "type": "BlockStatement",
          "body": [],
          "loc": {
            "start": {
              "line": 1,
              "column": 18
            },
            "end": {
              "line": 1,
              "column": 21
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 8
          },
          "end": {
            "line": 1,
            "column": 21
          }
        }
      },
      "finalizer": null,
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 21
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 21
    }
  }
}
	`, ast)
}

func Test275(t *testing.T) {
	ast, err := compile("try { } catch (eval) { }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "TryStatement",
      "block": {
        "type": "BlockStatement",
        "body": [],
        "loc": {
          "start": {
            "line": 1,
            "column": 4
          },
          "end": {
            "line": 1,
            "column": 7
          }
        }
      },
      "handler": {
        "type": "CatchClause",
        "param": {
          "type": "Identifier",
          "name": "eval",
          "loc": {
            "start": {
              "line": 1,
              "column": 15
            },
            "end": {
              "line": 1,
              "column": 19
            }
          }
        },
        "body": {
          "type": "BlockStatement",
          "body": [],
          "loc": {
            "start": {
              "line": 1,
              "column": 21
            },
            "end": {
              "line": 1,
              "column": 24
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 8
          },
          "end": {
            "line": 1,
            "column": 24
          }
        }
      },
      "finalizer": null,
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 24
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 24
    }
  }
}
	`, ast)
}

func Test276(t *testing.T) {
	ast, err := compile("try { } catch (arguments) { }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "TryStatement",
      "block": {
        "type": "BlockStatement",
        "body": [],
        "loc": {
          "start": {
            "line": 1,
            "column": 4
          },
          "end": {
            "line": 1,
            "column": 7
          }
        }
      },
      "handler": {
        "type": "CatchClause",
        "param": {
          "type": "Identifier",
          "name": "arguments",
          "loc": {
            "start": {
              "line": 1,
              "column": 15
            },
            "end": {
              "line": 1,
              "column": 24
            }
          }
        },
        "body": {
          "type": "BlockStatement",
          "body": [],
          "loc": {
            "start": {
              "line": 1,
              "column": 26
            },
            "end": {
              "line": 1,
              "column": 29
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 8
          },
          "end": {
            "line": 1,
            "column": 29
          }
        }
      },
      "finalizer": null,
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 29
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 29
    }
  }
}
	`, ast)
}

func Test277(t *testing.T) {
	ast, err := compile("try { } catch (e) { say(e) }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "TryStatement",
      "block": {
        "type": "BlockStatement",
        "body": [],
        "loc": {
          "start": {
            "line": 1,
            "column": 4
          },
          "end": {
            "line": 1,
            "column": 7
          }
        }
      },
      "handler": {
        "type": "CatchClause",
        "param": {
          "type": "Identifier",
          "name": "e",
          "loc": {
            "start": {
              "line": 1,
              "column": 15
            },
            "end": {
              "line": 1,
              "column": 16
            }
          }
        },
        "body": {
          "type": "BlockStatement",
          "body": [
            {
              "type": "ExpressionStatement",
              "expression": {
                "type": "CallExpression",
                "callee": {
                  "type": "Identifier",
                  "name": "say",
                  "loc": {
                    "start": {
                      "line": 1,
                      "column": 20
                    },
                    "end": {
                      "line": 1,
                      "column": 23
                    }
                  }
                },
                "arguments": [
                  {
                    "type": "Identifier",
                    "name": "e",
                    "loc": {
                      "start": {
                        "line": 1,
                        "column": 24
                      },
                      "end": {
                        "line": 1,
                        "column": 25
                      }
                    }
                  }
                ],
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 20
                  },
                  "end": {
                    "line": 1,
                    "column": 26
                  }
                }
              },
              "loc": {
                "start": {
                  "line": 1,
                  "column": 20
                },
                "end": {
                  "line": 1,
                  "column": 26
                }
              }
            }
          ],
          "loc": {
            "start": {
              "line": 1,
              "column": 18
            },
            "end": {
              "line": 1,
              "column": 28
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 8
          },
          "end": {
            "line": 1,
            "column": 28
          }
        }
      },
      "finalizer": null,
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 28
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 28
    }
  }
}
	`, ast)
}

func Test278(t *testing.T) {
	ast, err := compile("try { } finally { cleanup(stuff) }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "TryStatement",
      "block": {
        "type": "BlockStatement",
        "body": [],
        "loc": {
          "start": {
            "line": 1,
            "column": 4
          },
          "end": {
            "line": 1,
            "column": 7
          }
        }
      },
      "handler": null,
      "finalizer": {
        "type": "BlockStatement",
        "body": [
          {
            "type": "ExpressionStatement",
            "expression": {
              "type": "CallExpression",
              "callee": {
                "type": "Identifier",
                "name": "cleanup",
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 18
                  },
                  "end": {
                    "line": 1,
                    "column": 25
                  }
                }
              },
              "arguments": [
                {
                  "type": "Identifier",
                  "name": "stuff",
                  "loc": {
                    "start": {
                      "line": 1,
                      "column": 26
                    },
                    "end": {
                      "line": 1,
                      "column": 31
                    }
                  }
                }
              ],
              "loc": {
                "start": {
                  "line": 1,
                  "column": 18
                },
                "end": {
                  "line": 1,
                  "column": 32
                }
              }
            },
            "loc": {
              "start": {
                "line": 1,
                "column": 18
              },
              "end": {
                "line": 1,
                "column": 32
              }
            }
          }
        ],
        "loc": {
          "start": {
            "line": 1,
            "column": 16
          },
          "end": {
            "line": 1,
            "column": 34
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 34
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 34
    }
  }
}
	`, ast)
}

func Test279(t *testing.T) {
	ast, err := compile("try { doThat(); } catch (e) { say(e) }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "TryStatement",
      "block": {
        "type": "BlockStatement",
        "body": [
          {
            "type": "ExpressionStatement",
            "expression": {
              "type": "CallExpression",
              "callee": {
                "type": "Identifier",
                "name": "doThat",
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 6
                  },
                  "end": {
                    "line": 1,
                    "column": 12
                  }
                }
              },
              "arguments": [],
              "loc": {
                "start": {
                  "line": 1,
                  "column": 6
                },
                "end": {
                  "line": 1,
                  "column": 14
                }
              }
            },
            "loc": {
              "start": {
                "line": 1,
                "column": 6
              },
              "end": {
                "line": 1,
                "column": 15
              }
            }
          }
        ],
        "loc": {
          "start": {
            "line": 1,
            "column": 4
          },
          "end": {
            "line": 1,
            "column": 17
          }
        }
      },
      "handler": {
        "type": "CatchClause",
        "param": {
          "type": "Identifier",
          "name": "e",
          "loc": {
            "start": {
              "line": 1,
              "column": 25
            },
            "end": {
              "line": 1,
              "column": 26
            }
          }
        },
        "body": {
          "type": "BlockStatement",
          "body": [
            {
              "type": "ExpressionStatement",
              "expression": {
                "type": "CallExpression",
                "callee": {
                  "type": "Identifier",
                  "name": "say",
                  "loc": {
                    "start": {
                      "line": 1,
                      "column": 30
                    },
                    "end": {
                      "line": 1,
                      "column": 33
                    }
                  }
                },
                "arguments": [
                  {
                    "type": "Identifier",
                    "name": "e",
                    "loc": {
                      "start": {
                        "line": 1,
                        "column": 34
                      },
                      "end": {
                        "line": 1,
                        "column": 35
                      }
                    }
                  }
                ],
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 30
                  },
                  "end": {
                    "line": 1,
                    "column": 36
                  }
                }
              },
              "loc": {
                "start": {
                  "line": 1,
                  "column": 30
                },
                "end": {
                  "line": 1,
                  "column": 36
                }
              }
            }
          ],
          "loc": {
            "start": {
              "line": 1,
              "column": 28
            },
            "end": {
              "line": 1,
              "column": 38
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 18
          },
          "end": {
            "line": 1,
            "column": 38
          }
        }
      },
      "finalizer": null,
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 38
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 38
    }
  }
}
	`, ast)
}

func Test280(t *testing.T) {
	ast, err := compile("try { doThat(); } catch (e) { say(e) } finally { cleanup(stuff) }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "TryStatement",
      "block": {
        "type": "BlockStatement",
        "body": [
          {
            "type": "ExpressionStatement",
            "expression": {
              "type": "CallExpression",
              "callee": {
                "type": "Identifier",
                "name": "doThat",
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 6
                  },
                  "end": {
                    "line": 1,
                    "column": 12
                  }
                }
              },
              "arguments": [],
              "loc": {
                "start": {
                  "line": 1,
                  "column": 6
                },
                "end": {
                  "line": 1,
                  "column": 14
                }
              }
            },
            "loc": {
              "start": {
                "line": 1,
                "column": 6
              },
              "end": {
                "line": 1,
                "column": 15
              }
            }
          }
        ],
        "loc": {
          "start": {
            "line": 1,
            "column": 4
          },
          "end": {
            "line": 1,
            "column": 17
          }
        }
      },
      "handler": {
        "type": "CatchClause",
        "param": {
          "type": "Identifier",
          "name": "e",
          "loc": {
            "start": {
              "line": 1,
              "column": 25
            },
            "end": {
              "line": 1,
              "column": 26
            }
          }
        },
        "body": {
          "type": "BlockStatement",
          "body": [
            {
              "type": "ExpressionStatement",
              "expression": {
                "type": "CallExpression",
                "callee": {
                  "type": "Identifier",
                  "name": "say",
                  "loc": {
                    "start": {
                      "line": 1,
                      "column": 30
                    },
                    "end": {
                      "line": 1,
                      "column": 33
                    }
                  }
                },
                "arguments": [
                  {
                    "type": "Identifier",
                    "name": "e",
                    "loc": {
                      "start": {
                        "line": 1,
                        "column": 34
                      },
                      "end": {
                        "line": 1,
                        "column": 35
                      }
                    }
                  }
                ],
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 30
                  },
                  "end": {
                    "line": 1,
                    "column": 36
                  }
                }
              },
              "loc": {
                "start": {
                  "line": 1,
                  "column": 30
                },
                "end": {
                  "line": 1,
                  "column": 36
                }
              }
            }
          ],
          "loc": {
            "start": {
              "line": 1,
              "column": 28
            },
            "end": {
              "line": 1,
              "column": 38
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 18
          },
          "end": {
            "line": 1,
            "column": 38
          }
        }
      },
      "finalizer": {
        "type": "BlockStatement",
        "body": [
          {
            "type": "ExpressionStatement",
            "expression": {
              "type": "CallExpression",
              "callee": {
                "type": "Identifier",
                "name": "cleanup",
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 49
                  },
                  "end": {
                    "line": 1,
                    "column": 56
                  }
                }
              },
              "arguments": [
                {
                  "type": "Identifier",
                  "name": "stuff",
                  "loc": {
                    "start": {
                      "line": 1,
                      "column": 57
                    },
                    "end": {
                      "line": 1,
                      "column": 62
                    }
                  }
                }
              ],
              "loc": {
                "start": {
                  "line": 1,
                  "column": 49
                },
                "end": {
                  "line": 1,
                  "column": 63
                }
              }
            },
            "loc": {
              "start": {
                "line": 1,
                "column": 49
              },
              "end": {
                "line": 1,
                "column": 63
              }
            }
          }
        ],
        "loc": {
          "start": {
            "line": 1,
            "column": 47
          },
          "end": {
            "line": 1,
            "column": 65
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 65
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 65
    }
  }
}
	`, ast)
}

func Test281(t *testing.T) {
	ast, err := compile("debugger;")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "DebuggerStatement",
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 9
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 9
    }
  }
}
	`, ast)
}

func Test282(t *testing.T) {
	ast, err := compile("function hello() { sayHi(); }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "FunctionDeclaration",
      "id": {
        "type": "Identifier",
        "name": "hello",
        "loc": {
          "start": {
            "line": 1,
            "column": 9
          },
          "end": {
            "line": 1,
            "column": 14
          }
        }
      },
      "params": [],
      "body": {
        "type": "BlockStatement",
        "body": [
          {
            "type": "ExpressionStatement",
            "expression": {
              "type": "CallExpression",
              "callee": {
                "type": "Identifier",
                "name": "sayHi",
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 19
                  },
                  "end": {
                    "line": 1,
                    "column": 24
                  }
                }
              },
              "arguments": [],
              "loc": {
                "start": {
                  "line": 1,
                  "column": 19
                },
                "end": {
                  "line": 1,
                  "column": 26
                }
              }
            },
            "loc": {
              "start": {
                "line": 1,
                "column": 19
              },
              "end": {
                "line": 1,
                "column": 27
              }
            }
          }
        ],
        "loc": {
          "start": {
            "line": 1,
            "column": 17
          },
          "end": {
            "line": 1,
            "column": 29
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 29
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 29
    }
  }
}
	`, ast)
}

func Test283(t *testing.T) {
	ast, err := compile("function eval() { }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "FunctionDeclaration",
      "id": {
        "type": "Identifier",
        "name": "eval",
        "loc": {
          "start": {
            "line": 1,
            "column": 9
          },
          "end": {
            "line": 1,
            "column": 13
          }
        }
      },
      "params": [],
      "body": {
        "type": "BlockStatement",
        "body": [],
        "loc": {
          "start": {
            "line": 1,
            "column": 16
          },
          "end": {
            "line": 1,
            "column": 19
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 19
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 19
    }
  }
}
	`, ast)
}

func Test284(t *testing.T) {
	ast, err := compile("function arguments() { }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "FunctionDeclaration",
      "id": {
        "type": "Identifier",
        "name": "arguments",
        "loc": {
          "start": {
            "line": 1,
            "column": 9
          },
          "end": {
            "line": 1,
            "column": 18
          }
        }
      },
      "params": [],
      "body": {
        "type": "BlockStatement",
        "body": [],
        "loc": {
          "start": {
            "line": 1,
            "column": 21
          },
          "end": {
            "line": 1,
            "column": 24
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 24
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 24
    }
  }
}
	`, ast)
}

func Test285(t *testing.T) {
	ast, err := compile("function test(t, t) { }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "FunctionDeclaration",
      "id": {
        "type": "Identifier",
        "name": "test",
        "loc": {
          "start": {
            "line": 1,
            "column": 9
          },
          "end": {
            "line": 1,
            "column": 13
          }
        }
      },
      "params": [
        {
          "type": "Identifier",
          "name": "t",
          "loc": {
            "start": {
              "line": 1,
              "column": 14
            },
            "end": {
              "line": 1,
              "column": 15
            }
          }
        },
        {
          "type": "Identifier",
          "name": "t",
          "loc": {
            "start": {
              "line": 1,
              "column": 17
            },
            "end": {
              "line": 1,
              "column": 18
            }
          }
        }
      ],
      "body": {
        "type": "BlockStatement",
        "body": [],
        "loc": {
          "start": {
            "line": 1,
            "column": 20
          },
          "end": {
            "line": 1,
            "column": 23
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 23
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 23
    }
  }
}
	`, ast)
}

func Test286(t *testing.T) {
	ast, err := compile("(function test(t, t) { })")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "FunctionExpression",
        "id": {
          "type": "Identifier",
          "name": "test",
          "loc": {
            "start": {
              "line": 1,
              "column": 10
            },
            "end": {
              "line": 1,
              "column": 14
            }
          }
        },
        "params": [
          {
            "type": "Identifier",
            "name": "t",
            "loc": {
              "start": {
                "line": 1,
                "column": 15
              },
              "end": {
                "line": 1,
                "column": 16
              }
            }
          },
          {
            "type": "Identifier",
            "name": "t",
            "loc": {
              "start": {
                "line": 1,
                "column": 18
              },
              "end": {
                "line": 1,
                "column": 19
              }
            }
          }
        ],
        "body": {
          "type": "BlockStatement",
          "body": [],
          "loc": {
            "start": {
              "line": 1,
              "column": 21
            },
            "end": {
              "line": 1,
              "column": 24
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 1
          },
          "end": {
            "line": 1,
            "column": 24
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 25
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 25
    }
  }
}
	`, ast)
}

func Test287(t *testing.T) {
	ast, err := compile("function eval() { function inner() { \"use strict\" } }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "FunctionDeclaration",
      "id": {
        "type": "Identifier",
        "name": "eval",
        "loc": {
          "start": {
            "line": 1,
            "column": 9
          },
          "end": {
            "line": 1,
            "column": 13
          }
        }
      },
      "params": [],
      "body": {
        "type": "BlockStatement",
        "body": [
          {
            "type": "FunctionDeclaration",
            "id": {
              "type": "Identifier",
              "name": "inner",
              "loc": {
                "start": {
                  "line": 1,
                  "column": 27
                },
                "end": {
                  "line": 1,
                  "column": 32
                }
              }
            },
            "params": [],
            "body": {
              "type": "BlockStatement",
              "body": [
                {
                  "type": "ExpressionStatement",
                  "expression": {
                    "type": "Literal",
                    "value": "use strict",
                    "loc": {
                      "start": {
                        "line": 1,
                        "column": 37
                      },
                      "end": {
                        "line": 1,
                        "column": 49
                      }
                    }
                  },
                  "loc": {
                    "start": {
                      "line": 1,
                      "column": 37
                    },
                    "end": {
                      "line": 1,
                      "column": 49
                    }
                  }
                }
              ],
              "loc": {
                "start": {
                  "line": 1,
                  "column": 35
                },
                "end": {
                  "line": 1,
                  "column": 51
                }
              }
            },
            "loc": {
              "start": {
                "line": 1,
                "column": 18
              },
              "end": {
                "line": 1,
                "column": 51
              }
            }
          }
        ],
        "loc": {
          "start": {
            "line": 1,
            "column": 16
          },
          "end": {
            "line": 1,
            "column": 53
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 53
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 53
    }
  }
}
	`, ast)
}

func Test288(t *testing.T) {
	ast, err := compile("function hello(a) { sayHi(); }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "FunctionDeclaration",
      "id": {
        "type": "Identifier",
        "name": "hello",
        "loc": {
          "start": {
            "line": 1,
            "column": 9
          },
          "end": {
            "line": 1,
            "column": 14
          }
        }
      },
      "params": [
        {
          "type": "Identifier",
          "name": "a",
          "loc": {
            "start": {
              "line": 1,
              "column": 15
            },
            "end": {
              "line": 1,
              "column": 16
            }
          }
        }
      ],
      "body": {
        "type": "BlockStatement",
        "body": [
          {
            "type": "ExpressionStatement",
            "expression": {
              "type": "CallExpression",
              "callee": {
                "type": "Identifier",
                "name": "sayHi",
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 20
                  },
                  "end": {
                    "line": 1,
                    "column": 25
                  }
                }
              },
              "arguments": [],
              "loc": {
                "start": {
                  "line": 1,
                  "column": 20
                },
                "end": {
                  "line": 1,
                  "column": 27
                }
              }
            },
            "loc": {
              "start": {
                "line": 1,
                "column": 20
              },
              "end": {
                "line": 1,
                "column": 28
              }
            }
          }
        ],
        "loc": {
          "start": {
            "line": 1,
            "column": 18
          },
          "end": {
            "line": 1,
            "column": 30
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 30
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 30
    }
  }
}
	`, ast)
}

func Test289(t *testing.T) {
	ast, err := compile("function hello(a, b) { sayHi(); }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "FunctionDeclaration",
      "id": {
        "type": "Identifier",
        "name": "hello",
        "loc": {
          "start": {
            "line": 1,
            "column": 9
          },
          "end": {
            "line": 1,
            "column": 14
          }
        }
      },
      "params": [
        {
          "type": "Identifier",
          "name": "a",
          "loc": {
            "start": {
              "line": 1,
              "column": 15
            },
            "end": {
              "line": 1,
              "column": 16
            }
          }
        },
        {
          "type": "Identifier",
          "name": "b",
          "loc": {
            "start": {
              "line": 1,
              "column": 18
            },
            "end": {
              "line": 1,
              "column": 19
            }
          }
        }
      ],
      "body": {
        "type": "BlockStatement",
        "body": [
          {
            "type": "ExpressionStatement",
            "expression": {
              "type": "CallExpression",
              "callee": {
                "type": "Identifier",
                "name": "sayHi",
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 23
                  },
                  "end": {
                    "line": 1,
                    "column": 28
                  }
                }
              },
              "arguments": [],
              "loc": {
                "start": {
                  "line": 1,
                  "column": 23
                },
                "end": {
                  "line": 1,
                  "column": 30
                }
              }
            },
            "loc": {
              "start": {
                "line": 1,
                "column": 23
              },
              "end": {
                "line": 1,
                "column": 31
              }
            }
          }
        ],
        "loc": {
          "start": {
            "line": 1,
            "column": 21
          },
          "end": {
            "line": 1,
            "column": 33
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 33
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 33
    }
  }
}
	`, ast)
}

func Test290(t *testing.T) {
	ast, err := compile("function hello(...rest) { }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "FunctionDeclaration",
      "id": {
        "type": "Identifier",
        "name": "hello",
        "loc": {
          "start": {
            "line": 1,
            "column": 9
          },
          "end": {
            "line": 1,
            "column": 14
          }
        }
      },
      "params": [
        {
          "type": "RestElement",
          "argument": {
            "type": "Identifier",
            "name": "rest",
            "loc": {
              "start": {
                "line": 1,
                "column": 18
              },
              "end": {
                "line": 1,
                "column": 22
              }
            }
          }
        }
      ],
      "body": {
        "type": "BlockStatement",
        "body": [],
        "loc": {
          "start": {
            "line": 1,
            "column": 24
          },
          "end": {
            "line": 1,
            "column": 27
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 27
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 27
    }
  }
}
	`, ast)
}

func Test291(t *testing.T) {
	ast, err := compile("function hello(a, ...rest) { }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "FunctionDeclaration",
      "id": {
        "type": "Identifier",
        "name": "hello",
        "loc": {
          "start": {
            "line": 1,
            "column": 9
          },
          "end": {
            "line": 1,
            "column": 14
          }
        }
      },
      "params": [
        {
          "type": "Identifier",
          "name": "a",
          "loc": {
            "start": {
              "line": 1,
              "column": 15
            },
            "end": {
              "line": 1,
              "column": 16
            }
          }
        },
        {
          "type": "RestElement",
          "argument": {
            "type": "Identifier",
            "name": "rest",
            "loc": {
              "start": {
                "line": 1,
                "column": 21
              },
              "end": {
                "line": 1,
                "column": 25
              }
            }
          }
        }
      ],
      "body": {
        "type": "BlockStatement",
        "body": [],
        "loc": {
          "start": {
            "line": 1,
            "column": 27
          },
          "end": {
            "line": 1,
            "column": 30
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 30
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 30
    }
  }
}
	`, ast)
}

func Test292(t *testing.T) {
	ast, err := compile("var hi = function() { sayHi() };")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "VariableDeclaration",
      "declarations": [
        {
          "type": "VariableDeclarator",
          "id": {
            "type": "Identifier",
            "name": "hi",
            "loc": {
              "start": {
                "line": 1,
                "column": 4
              },
              "end": {
                "line": 1,
                "column": 6
              }
            }
          },
          "init": {
            "type": "FunctionExpression",
            "id": null,
            "params": [],
            "body": {
              "type": "BlockStatement",
              "body": [
                {
                  "type": "ExpressionStatement",
                  "expression": {
                    "type": "CallExpression",
                    "callee": {
                      "type": "Identifier",
                      "name": "sayHi",
                      "loc": {
                        "start": {
                          "line": 1,
                          "column": 22
                        },
                        "end": {
                          "line": 1,
                          "column": 27
                        }
                      }
                    },
                    "arguments": [],
                    "loc": {
                      "start": {
                        "line": 1,
                        "column": 22
                      },
                      "end": {
                        "line": 1,
                        "column": 29
                      }
                    }
                  },
                  "loc": {
                    "start": {
                      "line": 1,
                      "column": 22
                    },
                    "end": {
                      "line": 1,
                      "column": 29
                    }
                  }
                }
              ],
              "loc": {
                "start": {
                  "line": 1,
                  "column": 20
                },
                "end": {
                  "line": 1,
                  "column": 31
                }
              }
            },
            "loc": {
              "start": {
                "line": 1,
                "column": 9
              },
              "end": {
                "line": 1,
                "column": 31
              }
            }
          },
          "loc": {
            "start": {
              "line": 1,
              "column": 4
            },
            "end": {
              "line": 1,
              "column": 31
            }
          }
        }
      ],
      "kind": "var",
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 32
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 32
    }
  }
}
	`, ast)
}

func Test293(t *testing.T) {
	ast, err := compile("var hi = function (...r) { sayHi() };")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "VariableDeclaration",
      "declarations": [
        {
          "type": "VariableDeclarator",
          "id": {
            "type": "Identifier",
            "name": "hi",
            "loc": {
              "start": {
                "line": 1,
                "column": 4
              },
              "end": {
                "line": 1,
                "column": 6
              }
            }
          },
          "init": {
            "type": "FunctionExpression",
            "id": null,
            "params": [
              {
                "type": "RestElement",
                "argument": {
                  "type": "Identifier",
                  "name": "r",
                  "loc": {
                    "start": {
                      "line": 1,
                      "column": 22
                    },
                    "end": {
                      "line": 1,
                      "column": 23
                    }
                  }
                }
              }
            ],
            "body": {
              "type": "BlockStatement",
              "body": [
                {
                  "type": "ExpressionStatement",
                  "expression": {
                    "type": "CallExpression",
                    "callee": {
                      "type": "Identifier",
                      "name": "sayHi",
                      "loc": {
                        "start": {
                          "line": 1,
                          "column": 27
                        },
                        "end": {
                          "line": 1,
                          "column": 32
                        }
                      }
                    },
                    "arguments": [],
                    "loc": {
                      "start": {
                        "line": 1,
                        "column": 27
                      },
                      "end": {
                        "line": 1,
                        "column": 34
                      }
                    }
                  },
                  "loc": {
                    "start": {
                      "line": 1,
                      "column": 27
                    },
                    "end": {
                      "line": 1,
                      "column": 34
                    }
                  }
                }
              ],
              "loc": {
                "start": {
                  "line": 1,
                  "column": 25
                },
                "end": {
                  "line": 1,
                  "column": 36
                }
              }
            },
            "loc": {
              "start": {
                "line": 1,
                "column": 9
              },
              "end": {
                "line": 1,
                "column": 36
              }
            }
          },
          "loc": {
            "start": {
              "line": 1,
              "column": 4
            },
            "end": {
              "line": 1,
              "column": 36
            }
          }
        }
      ],
      "kind": "var",
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 37
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 37
    }
  }
}
	`, ast)
}

func Test294(t *testing.T) {
	ast, err := compile("var hi = function eval() { };")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "VariableDeclaration",
      "declarations": [
        {
          "type": "VariableDeclarator",
          "id": {
            "type": "Identifier",
            "name": "hi",
            "loc": {
              "start": {
                "line": 1,
                "column": 4
              },
              "end": {
                "line": 1,
                "column": 6
              }
            }
          },
          "init": {
            "type": "FunctionExpression",
            "id": {
              "type": "Identifier",
              "name": "eval",
              "loc": {
                "start": {
                  "line": 1,
                  "column": 18
                },
                "end": {
                  "line": 1,
                  "column": 22
                }
              }
            },
            "params": [],
            "body": {
              "type": "BlockStatement",
              "body": [],
              "loc": {
                "start": {
                  "line": 1,
                  "column": 25
                },
                "end": {
                  "line": 1,
                  "column": 28
                }
              }
            },
            "loc": {
              "start": {
                "line": 1,
                "column": 9
              },
              "end": {
                "line": 1,
                "column": 28
              }
            }
          },
          "loc": {
            "start": {
              "line": 1,
              "column": 4
            },
            "end": {
              "line": 1,
              "column": 28
            }
          }
        }
      ],
      "kind": "var",
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 29
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 29
    }
  }
}
	`, ast)
}

func Test295(t *testing.T) {
	ast, err := compile("var hi = function arguments() { };")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "VariableDeclaration",
      "declarations": [
        {
          "type": "VariableDeclarator",
          "id": {
            "type": "Identifier",
            "name": "hi",
            "loc": {
              "start": {
                "line": 1,
                "column": 4
              },
              "end": {
                "line": 1,
                "column": 6
              }
            }
          },
          "init": {
            "type": "FunctionExpression",
            "id": {
              "type": "Identifier",
              "name": "arguments",
              "loc": {
                "start": {
                  "line": 1,
                  "column": 18
                },
                "end": {
                  "line": 1,
                  "column": 27
                }
              }
            },
            "params": [],
            "body": {
              "type": "BlockStatement",
              "body": [],
              "loc": {
                "start": {
                  "line": 1,
                  "column": 30
                },
                "end": {
                  "line": 1,
                  "column": 33
                }
              }
            },
            "loc": {
              "start": {
                "line": 1,
                "column": 9
              },
              "end": {
                "line": 1,
                "column": 33
              }
            }
          },
          "loc": {
            "start": {
              "line": 1,
              "column": 4
            },
            "end": {
              "line": 1,
              "column": 33
            }
          }
        }
      ],
      "kind": "var",
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 34
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 34
    }
  }
}
	`, ast)
}

func Test296(t *testing.T) {
	ast, err := compile("var hello = function hi() { sayHi() };")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "VariableDeclaration",
      "declarations": [
        {
          "type": "VariableDeclarator",
          "id": {
            "type": "Identifier",
            "name": "hello",
            "loc": {
              "start": {
                "line": 1,
                "column": 4
              },
              "end": {
                "line": 1,
                "column": 9
              }
            }
          },
          "init": {
            "type": "FunctionExpression",
            "id": {
              "type": "Identifier",
              "name": "hi",
              "loc": {
                "start": {
                  "line": 1,
                  "column": 21
                },
                "end": {
                  "line": 1,
                  "column": 23
                }
              }
            },
            "params": [],
            "body": {
              "type": "BlockStatement",
              "body": [
                {
                  "type": "ExpressionStatement",
                  "expression": {
                    "type": "CallExpression",
                    "callee": {
                      "type": "Identifier",
                      "name": "sayHi",
                      "loc": {
                        "start": {
                          "line": 1,
                          "column": 28
                        },
                        "end": {
                          "line": 1,
                          "column": 33
                        }
                      }
                    },
                    "arguments": [],
                    "loc": {
                      "start": {
                        "line": 1,
                        "column": 28
                      },
                      "end": {
                        "line": 1,
                        "column": 35
                      }
                    }
                  },
                  "loc": {
                    "start": {
                      "line": 1,
                      "column": 28
                    },
                    "end": {
                      "line": 1,
                      "column": 35
                    }
                  }
                }
              ],
              "loc": {
                "start": {
                  "line": 1,
                  "column": 26
                },
                "end": {
                  "line": 1,
                  "column": 37
                }
              }
            },
            "loc": {
              "start": {
                "line": 1,
                "column": 12
              },
              "end": {
                "line": 1,
                "column": 37
              }
            }
          },
          "loc": {
            "start": {
              "line": 1,
              "column": 4
            },
            "end": {
              "line": 1,
              "column": 37
            }
          }
        }
      ],
      "kind": "var",
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 38
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 38
    }
  }
}
	`, ast)
}

func Test297(t *testing.T) {
	ast, err := compile("(function(){})")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "FunctionExpression",
        "id": null,
        "params": [],
        "body": {
          "type": "BlockStatement",
          "body": [],
          "loc": {
            "start": {
              "line": 1,
              "column": 11
            },
            "end": {
              "line": 1,
              "column": 13
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 1
          },
          "end": {
            "line": 1,
            "column": 13
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 14
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 14
    }
  }
}
	`, ast)
}

func Test298(t *testing.T) {
	ast, err := compile("{ x\n++y }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "BlockStatement",
      "body": [
        {
          "type": "ExpressionStatement",
          "expression": {
            "type": "Identifier",
            "name": "x",
            "loc": {
              "start": {
                "line": 1,
                "column": 2
              },
              "end": {
                "line": 1,
                "column": 3
              }
            }
          },
          "loc": {
            "start": {
              "line": 1,
              "column": 2
            },
            "end": {
              "line": 1,
              "column": 3
            }
          }
        },
        {
          "type": "ExpressionStatement",
          "expression": {
            "type": "UpdateExpression",
            "operator": "++",
            "prefix": true,
            "argument": {
              "type": "Identifier",
              "name": "y",
              "loc": {
                "start": {
                  "line": 2,
                  "column": 2
                },
                "end": {
                  "line": 2,
                  "column": 3
                }
              }
            },
            "loc": {
              "start": {
                "line": 2,
                "column": 0
              },
              "end": {
                "line": 2,
                "column": 3
              }
            }
          },
          "loc": {
            "start": {
              "line": 2,
              "column": 0
            },
            "end": {
              "line": 2,
              "column": 3
            }
          }
        }
      ],
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 2,
          "column": 5
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 2,
      "column": 5
    }
  }
}
	`, ast)
}

func Test299(t *testing.T) {
	ast, err := compile("{ x\n--y }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "BlockStatement",
      "body": [
        {
          "type": "ExpressionStatement",
          "expression": {
            "type": "Identifier",
            "name": "x",
            "loc": {
              "start": {
                "line": 1,
                "column": 2
              },
              "end": {
                "line": 1,
                "column": 3
              }
            }
          },
          "loc": {
            "start": {
              "line": 1,
              "column": 2
            },
            "end": {
              "line": 1,
              "column": 3
            }
          }
        },
        {
          "type": "ExpressionStatement",
          "expression": {
            "type": "UpdateExpression",
            "operator": "--",
            "prefix": true,
            "argument": {
              "type": "Identifier",
              "name": "y",
              "loc": {
                "start": {
                  "line": 2,
                  "column": 2
                },
                "end": {
                  "line": 2,
                  "column": 3
                }
              }
            },
            "loc": {
              "start": {
                "line": 2,
                "column": 0
              },
              "end": {
                "line": 2,
                "column": 3
              }
            }
          },
          "loc": {
            "start": {
              "line": 2,
              "column": 0
            },
            "end": {
              "line": 2,
              "column": 3
            }
          }
        }
      ],
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 2,
          "column": 5
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 2,
      "column": 5
    }
  }
}
	`, ast)
}

func Test300(t *testing.T) {
	ast, err := compile("var x /* comment */;")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "VariableDeclaration",
      "declarations": [
        {
          "type": "VariableDeclarator",
          "id": {
            "type": "Identifier",
            "name": "x",
            "loc": {
              "start": {
                "line": 1,
                "column": 4
              },
              "end": {
                "line": 1,
                "column": 5
              }
            }
          },
          "init": null,
          "loc": {
            "start": {
              "line": 1,
              "column": 4
            },
            "end": {
              "line": 1,
              "column": 5
            }
          }
        }
      ],
      "kind": "var",
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 20
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 20
    }
  }
}
	`, ast)
}

func Test301(t *testing.T) {
	ast, err := compile("{ var x = 14, y = 3\nz; }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "BlockStatement",
      "body": [
        {
          "type": "VariableDeclaration",
          "declarations": [
            {
              "type": "VariableDeclarator",
              "id": {
                "type": "Identifier",
                "name": "x",
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 6
                  },
                  "end": {
                    "line": 1,
                    "column": 7
                  }
                }
              },
              "init": {
                "type": "Literal",
                "value": 14,
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 10
                  },
                  "end": {
                    "line": 1,
                    "column": 12
                  }
                }
              },
              "loc": {
                "start": {
                  "line": 1,
                  "column": 6
                },
                "end": {
                  "line": 1,
                  "column": 12
                }
              }
            },
            {
              "type": "VariableDeclarator",
              "id": {
                "type": "Identifier",
                "name": "y",
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 14
                  },
                  "end": {
                    "line": 1,
                    "column": 15
                  }
                }
              },
              "init": {
                "type": "Literal",
                "value": 3,
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 18
                  },
                  "end": {
                    "line": 1,
                    "column": 19
                  }
                }
              },
              "loc": {
                "start": {
                  "line": 1,
                  "column": 14
                },
                "end": {
                  "line": 1,
                  "column": 19
                }
              }
            }
          ],
          "kind": "var",
          "loc": {
            "start": {
              "line": 1,
              "column": 2
            },
            "end": {
              "line": 1,
              "column": 19
            }
          }
        },
        {
          "type": "ExpressionStatement",
          "expression": {
            "type": "Identifier",
            "name": "z",
            "loc": {
              "start": {
                "line": 2,
                "column": 0
              },
              "end": {
                "line": 2,
                "column": 1
              }
            }
          },
          "loc": {
            "start": {
              "line": 2,
              "column": 0
            },
            "end": {
              "line": 2,
              "column": 2
            }
          }
        }
      ],
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 2,
          "column": 4
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 2,
      "column": 4
    }
  }
}
	`, ast)
}

func Test302(t *testing.T) {
	ast, err := compile("while (true) { continue\nthere; }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "WhileStatement",
      "test": {
        "type": "Literal",
        "value": true,
        "loc": {
          "start": {
            "line": 1,
            "column": 7
          },
          "end": {
            "line": 1,
            "column": 11
          }
        }
      },
      "body": {
        "type": "BlockStatement",
        "body": [
          {
            "type": "ContinueStatement",
            "label": null,
            "loc": {
              "start": {
                "line": 1,
                "column": 15
              },
              "end": {
                "line": 1,
                "column": 23
              }
            }
          },
          {
            "type": "ExpressionStatement",
            "expression": {
              "type": "Identifier",
              "name": "there",
              "loc": {
                "start": {
                  "line": 2,
                  "column": 0
                },
                "end": {
                  "line": 2,
                  "column": 5
                }
              }
            },
            "loc": {
              "start": {
                "line": 2,
                "column": 0
              },
              "end": {
                "line": 2,
                "column": 6
              }
            }
          }
        ],
        "loc": {
          "start": {
            "line": 1,
            "column": 13
          },
          "end": {
            "line": 2,
            "column": 8
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 2,
          "column": 8
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 2,
      "column": 8
    }
  }
}
	`, ast)
}

func Test303(t *testing.T) {
	ast, err := compile("while (true) { continue // Comment\nthere; }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "WhileStatement",
      "test": {
        "type": "Literal",
        "value": true,
        "loc": {
          "start": {
            "line": 1,
            "column": 7
          },
          "end": {
            "line": 1,
            "column": 11
          }
        }
      },
      "body": {
        "type": "BlockStatement",
        "body": [
          {
            "type": "ContinueStatement",
            "label": null,
            "loc": {
              "start": {
                "line": 1,
                "column": 15
              },
              "end": {
                "line": 1,
                "column": 23
              }
            }
          },
          {
            "type": "ExpressionStatement",
            "expression": {
              "type": "Identifier",
              "name": "there",
              "loc": {
                "start": {
                  "line": 2,
                  "column": 0
                },
                "end": {
                  "line": 2,
                  "column": 5
                }
              }
            },
            "loc": {
              "start": {
                "line": 2,
                "column": 0
              },
              "end": {
                "line": 2,
                "column": 6
              }
            }
          }
        ],
        "loc": {
          "start": {
            "line": 1,
            "column": 13
          },
          "end": {
            "line": 2,
            "column": 8
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 2,
          "column": 8
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 2,
      "column": 8
    }
  }
}
	`, ast)
}

func Test304(t *testing.T) {
	ast, err := compile("while (true) { continue /* Multiline\nComment */there; }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "WhileStatement",
      "test": {
        "type": "Literal",
        "value": true,
        "loc": {
          "start": {
            "line": 1,
            "column": 7
          },
          "end": {
            "line": 1,
            "column": 11
          }
        }
      },
      "body": {
        "type": "BlockStatement",
        "body": [
          {
            "type": "ContinueStatement",
            "label": null,
            "loc": {
              "start": {
                "line": 1,
                "column": 15
              },
              "end": {
                "line": 1,
                "column": 23
              }
            }
          },
          {
            "type": "ExpressionStatement",
            "expression": {
              "type": "Identifier",
              "name": "there",
              "loc": {
                "start": {
                  "line": 2,
                  "column": 10
                },
                "end": {
                  "line": 2,
                  "column": 15
                }
              }
            },
            "loc": {
              "start": {
                "line": 2,
                "column": 10
              },
              "end": {
                "line": 2,
                "column": 16
              }
            }
          }
        ],
        "loc": {
          "start": {
            "line": 1,
            "column": 13
          },
          "end": {
            "line": 2,
            "column": 18
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 2,
          "column": 18
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 2,
      "column": 18
    }
  }
}
	`, ast)
}

func Test305(t *testing.T) {
	ast, err := compile("while (true) { break\nthere; }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "WhileStatement",
      "test": {
        "type": "Literal",
        "value": true,
        "loc": {
          "start": {
            "line": 1,
            "column": 7
          },
          "end": {
            "line": 1,
            "column": 11
          }
        }
      },
      "body": {
        "type": "BlockStatement",
        "body": [
          {
            "type": "BreakStatement",
            "label": null,
            "loc": {
              "start": {
                "line": 1,
                "column": 15
              },
              "end": {
                "line": 1,
                "column": 20
              }
            }
          },
          {
            "type": "ExpressionStatement",
            "expression": {
              "type": "Identifier",
              "name": "there",
              "loc": {
                "start": {
                  "line": 2,
                  "column": 0
                },
                "end": {
                  "line": 2,
                  "column": 5
                }
              }
            },
            "loc": {
              "start": {
                "line": 2,
                "column": 0
              },
              "end": {
                "line": 2,
                "column": 6
              }
            }
          }
        ],
        "loc": {
          "start": {
            "line": 1,
            "column": 13
          },
          "end": {
            "line": 2,
            "column": 8
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 2,
          "column": 8
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 2,
      "column": 8
    }
  }
}
	`, ast)
}

func Test306(t *testing.T) {
	ast, err := compile("while (true) { break // Comment\nthere; }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "WhileStatement",
      "test": {
        "type": "Literal",
        "value": true,
        "loc": {
          "start": {
            "line": 1,
            "column": 7
          },
          "end": {
            "line": 1,
            "column": 11
          }
        }
      },
      "body": {
        "type": "BlockStatement",
        "body": [
          {
            "type": "BreakStatement",
            "label": null,
            "loc": {
              "start": {
                "line": 1,
                "column": 15
              },
              "end": {
                "line": 1,
                "column": 20
              }
            }
          },
          {
            "type": "ExpressionStatement",
            "expression": {
              "type": "Identifier",
              "name": "there",
              "loc": {
                "start": {
                  "line": 2,
                  "column": 0
                },
                "end": {
                  "line": 2,
                  "column": 5
                }
              }
            },
            "loc": {
              "start": {
                "line": 2,
                "column": 0
              },
              "end": {
                "line": 2,
                "column": 6
              }
            }
          }
        ],
        "loc": {
          "start": {
            "line": 1,
            "column": 13
          },
          "end": {
            "line": 2,
            "column": 8
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 2,
          "column": 8
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 2,
      "column": 8
    }
  }
}
	`, ast)
}

func Test307(t *testing.T) {
	ast, err := compile("while (true) { break /* Multiline\nComment */there; }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "WhileStatement",
      "test": {
        "type": "Literal",
        "value": true,
        "loc": {
          "start": {
            "line": 1,
            "column": 7
          },
          "end": {
            "line": 1,
            "column": 11
          }
        }
      },
      "body": {
        "type": "BlockStatement",
        "body": [
          {
            "type": "BreakStatement",
            "label": null,
            "loc": {
              "start": {
                "line": 1,
                "column": 15
              },
              "end": {
                "line": 1,
                "column": 20
              }
            }
          },
          {
            "type": "ExpressionStatement",
            "expression": {
              "type": "Identifier",
              "name": "there",
              "loc": {
                "start": {
                  "line": 2,
                  "column": 10
                },
                "end": {
                  "line": 2,
                  "column": 15
                }
              }
            },
            "loc": {
              "start": {
                "line": 2,
                "column": 10
              },
              "end": {
                "line": 2,
                "column": 16
              }
            }
          }
        ],
        "loc": {
          "start": {
            "line": 1,
            "column": 13
          },
          "end": {
            "line": 2,
            "column": 18
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 2,
          "column": 18
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 2,
      "column": 18
    }
  }
}
	`, ast)
}

func Test308(t *testing.T) {
	ast, err := compile("(function(){ return\nx; })")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "FunctionExpression",
        "id": null,
        "params": [],
        "body": {
          "type": "BlockStatement",
          "body": [
            {
              "type": "ReturnStatement",
              "argument": null,
              "loc": {
                "start": {
                  "line": 1,
                  "column": 13
                },
                "end": {
                  "line": 1,
                  "column": 19
                }
              }
            },
            {
              "type": "ExpressionStatement",
              "expression": {
                "type": "Identifier",
                "name": "x",
                "loc": {
                  "start": {
                    "line": 2,
                    "column": 0
                  },
                  "end": {
                    "line": 2,
                    "column": 1
                  }
                }
              },
              "loc": {
                "start": {
                  "line": 2,
                  "column": 0
                },
                "end": {
                  "line": 2,
                  "column": 2
                }
              }
            }
          ],
          "loc": {
            "start": {
              "line": 1,
              "column": 11
            },
            "end": {
              "line": 2,
              "column": 4
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 1
          },
          "end": {
            "line": 2,
            "column": 4
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 2,
          "column": 5
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 2,
      "column": 5
    }
  }
}
	`, ast)
}

func Test309(t *testing.T) {
	ast, err := compile("(function(){ return // Comment\nx; })")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "FunctionExpression",
        "id": null,
        "params": [],
        "body": {
          "type": "BlockStatement",
          "body": [
            {
              "type": "ReturnStatement",
              "argument": null,
              "loc": {
                "start": {
                  "line": 1,
                  "column": 13
                },
                "end": {
                  "line": 1,
                  "column": 19
                }
              }
            },
            {
              "type": "ExpressionStatement",
              "expression": {
                "type": "Identifier",
                "name": "x",
                "loc": {
                  "start": {
                    "line": 2,
                    "column": 0
                  },
                  "end": {
                    "line": 2,
                    "column": 1
                  }
                }
              },
              "loc": {
                "start": {
                  "line": 2,
                  "column": 0
                },
                "end": {
                  "line": 2,
                  "column": 2
                }
              }
            }
          ],
          "loc": {
            "start": {
              "line": 1,
              "column": 11
            },
            "end": {
              "line": 2,
              "column": 4
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 1
          },
          "end": {
            "line": 2,
            "column": 4
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 2,
          "column": 5
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 2,
      "column": 5
    }
  }
}
	`, ast)
}

func Test310(t *testing.T) {
	ast, err := compile("(function(){ return/* Multiline\nComment */x; })")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "FunctionExpression",
        "id": null,
        "params": [],
        "body": {
          "type": "BlockStatement",
          "body": [
            {
              "type": "ReturnStatement",
              "argument": null,
              "loc": {
                "start": {
                  "line": 1,
                  "column": 13
                },
                "end": {
                  "line": 1,
                  "column": 19
                }
              }
            },
            {
              "type": "ExpressionStatement",
              "expression": {
                "type": "Identifier",
                "name": "x",
                "loc": {
                  "start": {
                    "line": 2,
                    "column": 10
                  },
                  "end": {
                    "line": 2,
                    "column": 11
                  }
                }
              },
              "loc": {
                "start": {
                  "line": 2,
                  "column": 10
                },
                "end": {
                  "line": 2,
                  "column": 12
                }
              }
            }
          ],
          "loc": {
            "start": {
              "line": 1,
              "column": 11
            },
            "end": {
              "line": 2,
              "column": 14
            }
          }
        },
        "loc": {
          "start": {
            "line": 1,
            "column": 1
          },
          "end": {
            "line": 2,
            "column": 14
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 2,
          "column": 15
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 2,
      "column": 15
    }
  }
}
	`, ast)
}

func Test311(t *testing.T) {
	ast, err := compile("{ throw error\nerror; }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "BlockStatement",
      "body": [
        {
          "type": "ThrowStatement",
          "argument": {
            "type": "Identifier",
            "name": "error",
            "loc": {
              "start": {
                "line": 1,
                "column": 8
              },
              "end": {
                "line": 1,
                "column": 13
              }
            }
          },
          "loc": {
            "start": {
              "line": 1,
              "column": 2
            },
            "end": {
              "line": 1,
              "column": 13
            }
          }
        },
        {
          "type": "ExpressionStatement",
          "expression": {
            "type": "Identifier",
            "name": "error",
            "loc": {
              "start": {
                "line": 2,
                "column": 0
              },
              "end": {
                "line": 2,
                "column": 5
              }
            }
          },
          "loc": {
            "start": {
              "line": 2,
              "column": 0
            },
            "end": {
              "line": 2,
              "column": 6
            }
          }
        }
      ],
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 2,
          "column": 8
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 2,
      "column": 8
    }
  }
}
	`, ast)
}

func Test312(t *testing.T) {
	ast, err := compile("{ throw error// Comment\nerror; }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "BlockStatement",
      "body": [
        {
          "type": "ThrowStatement",
          "argument": {
            "type": "Identifier",
            "name": "error",
            "loc": {
              "start": {
                "line": 1,
                "column": 8
              },
              "end": {
                "line": 1,
                "column": 13
              }
            }
          },
          "loc": {
            "start": {
              "line": 1,
              "column": 2
            },
            "end": {
              "line": 1,
              "column": 13
            }
          }
        },
        {
          "type": "ExpressionStatement",
          "expression": {
            "type": "Identifier",
            "name": "error",
            "loc": {
              "start": {
                "line": 2,
                "column": 0
              },
              "end": {
                "line": 2,
                "column": 5
              }
            }
          },
          "loc": {
            "start": {
              "line": 2,
              "column": 0
            },
            "end": {
              "line": 2,
              "column": 6
            }
          }
        }
      ],
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 2,
          "column": 8
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 2,
      "column": 8
    }
  }
}
	`, ast)
}

func Test313(t *testing.T) {
	ast, err := compile("{ throw error/* Multiline\nComment */error; }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "BlockStatement",
      "body": [
        {
          "type": "ThrowStatement",
          "argument": {
            "type": "Identifier",
            "name": "error",
            "loc": {
              "start": {
                "line": 1,
                "column": 8
              },
              "end": {
                "line": 1,
                "column": 13
              }
            }
          },
          "loc": {
            "start": {
              "line": 1,
              "column": 2
            },
            "end": {
              "line": 1,
              "column": 13
            }
          }
        },
        {
          "type": "ExpressionStatement",
          "expression": {
            "type": "Identifier",
            "name": "error",
            "loc": {
              "start": {
                "line": 2,
                "column": 10
              },
              "end": {
                "line": 2,
                "column": 15
              }
            }
          },
          "loc": {
            "start": {
              "line": 2,
              "column": 10
            },
            "end": {
              "line": 2,
              "column": 16
            }
          }
        }
      ],
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 2,
          "column": 18
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 2,
      "column": 18
    }
  }
}
	`, ast)
}

func Test314(t *testing.T) {
	ast, err := compile("")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 0
    }
  }
}
	`, ast)
}

func Test315(t *testing.T) {
	ast, err := compile("foo: if (true) break foo;")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 25
    }
  },
  "body": [
    {
      "type": "LabeledStatement",
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 25
        }
      },
      "body": {
        "type": "IfStatement",
        "loc": {
          "start": {
            "line": 1,
            "column": 5
          },
          "end": {
            "line": 1,
            "column": 25
          }
        },
        "test": {
          "type": "Literal",
          "loc": {
            "start": {
              "line": 1,
              "column": 9
            },
            "end": {
              "line": 1,
              "column": 13
            }
          },
          "value": true
        },
        "consequent": {
          "type": "BreakStatement",
          "loc": {
            "start": {
              "line": 1,
              "column": 15
            },
            "end": {
              "line": 1,
              "column": 25
            }
          },
          "label": {
            "type": "Identifier",
            "loc": {
              "start": {
                "line": 1,
                "column": 21
              },
              "end": {
                "line": 1,
                "column": 24
              }
            },
            "name": "foo"
          }
        },
        "alternate": null
      },
      "label": {
        "type": "Identifier",
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 3
          }
        },
        "name": "foo"
      }
    }
  ]
}
	`, ast)
}

func Test316(t *testing.T) {
	ast, err := compile("a: { b: switch(x) {} }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "LabeledStatement",
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 22
        }
      },
      "body": {
        "type": "BlockStatement",
        "loc": {
          "start": {
            "line": 1,
            "column": 3
          },
          "end": {
            "line": 1,
            "column": 22
          }
        },
        "body": [
          {
            "type": "LabeledStatement",
            "loc": {
              "start": {
                "line": 1,
                "column": 5
              },
              "end": {
                "line": 1,
                "column": 20
              }
            },
            "body": {
              "type": "SwitchStatement",
              "discriminant": {
                "type": "Identifier",
                "name": "x",
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 15
                  },
                  "end": {
                    "line": 1,
                    "column": 16
                  }
                }
              },
              "cases": []
            },
            "label": {
              "type": "Identifier",
              "name": "b"
            }
          }
        ]
      },
      "label": {
        "type": "Identifier",
        "name": "a"
      }
    }
  ]
}
	`, ast)
}

func Test317(t *testing.T) {
	ast, err := compile("(function () {\n 'use strict';\n '\u0000';\n}())")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 4,
      "column": 4
    }
  },
  "body": [
    {
      "type": "ExpressionStatement",
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 4,
          "column": 4
        }
      },
      "expression": {
        "type": "CallExpression",
        "loc": {
          "start": {
            "line": 1,
            "column": 1
          },
          "end": {
            "line": 4,
            "column": 3
          }
        },
        "callee": {
          "type": "FunctionExpression",
          "loc": {
            "start": {
              "line": 1,
              "column": 1
            },
            "end": {
              "line": 4,
              "column": 1
            }
          },
          "id": null,
          "params": [],
          "body": {
            "type": "BlockStatement",
            "loc": {
              "start": {
                "line": 1,
                "column": 13
              },
              "end": {
                "line": 4,
                "column": 1
              }
            },
            "body": [
              {
                "type": "ExpressionStatement",
                "loc": {
                  "start": {
                    "line": 2,
                    "column": 1
                  },
                  "end": {
                    "line": 2,
                    "column": 14
                  }
                },
                "expression": {
                  "type": "Literal",
                  "loc": {
                    "start": {
                      "line": 2,
                      "column": 1
                    },
                    "end": {
                      "line": 2,
                      "column": 13
                    }
                  },
                  "value": "use strict"
                }
              },
              {
                "type": "ExpressionStatement",
                "loc": {
                  "start": {
                    "line": 3,
                    "column": 1
                  },
                  "end": {
                    "line": 3,
                    "column": 5
                  }
                },
                "expression": {
                  "type": "Literal",
                  "loc": {
                    "start": {
                      "line": 3,
                      "column": 1
                    },
                    "end": {
                      "line": 3,
                      "column": 4
                    }
                  },
                  "value": "\u0000"
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

func Test318(t *testing.T) {
	ast, err := compile("123..toString(10)")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "CallExpression",
        "callee": {
          "type": "MemberExpression",
          "object": {
            "type": "Literal",
            "value": 123
          },
          "property": {
            "type": "Identifier",
            "name": "toString"
          },
          "computed": false
        },
        "arguments": [
          {
            "type": "Literal",
            "value": 10
          }
        ]
      }
    }
  ]
}
	`, ast)
}

func Test319(t *testing.T) {
	ast, err := compile("123.+2")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "BinaryExpression",
        "left": {
          "type": "Literal",
          "value": 123
        },
        "operator": "+",
        "right": {
          "type": "Literal",
          "value": 2
        }
      }
    }
  ]
}
	`, ast)
}

func Test320(t *testing.T) {
	ast, err := compile("a\u2028b")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "Identifier",
        "name": "a"
      }
    },
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "Identifier",
        "name": "b"
      }
    }
  ]
}
	`, ast)
}

func Test321(t *testing.T) {
	ast, err := compile("'a\\u0026b'")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "Literal",
        "value": "a&b"
      }
    }
  ]
}
	`, ast)
}

func Test322(t *testing.T) {
	ast, err := compile("foo: 10; foo: 20;")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "LabeledStatement",
      "body": {
        "type": "ExpressionStatement",
        "expression": {
          "type": "Literal",
          "value": 10
        }
      },
      "label": {
        "type": "Identifier",
        "name": "foo"
      }
    },
    {
      "type": "LabeledStatement",
      "body": {
        "type": "ExpressionStatement",
        "expression": {
          "type": "Literal",
          "value": 20
        }
      },
      "label": {
        "type": "Identifier",
        "name": "foo"
      }
    }
  ]
}
	`, ast)
}

func Test323(t *testing.T) {
	ast, err := compile("if(1)/  foo/")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "IfStatement",
      "test": {
        "type": "Literal",
        "value": 1
      },
      "consequent": {
        "type": "ExpressionStatement",
        "expression": {
          "type": "Literal",
          "regexp": {
            "pattern": "  foo"
          }
        }
      },
      "alternate": null
    }
  ]
}
	`, ast)
}

func Test324(t *testing.T) {
	ast, err := compile("price_9̶9̶_89")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "Identifier",
        "name": "price_9̶9̶_89"
      }
    }
  ]
}
	`, ast)
}

func Test325(t *testing.T) {
	ast, err := compile("a.in / b")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "BinaryExpression",
        "left": {
          "type": "MemberExpression",
          "object": {
            "type": "Identifier",
            "name": "a"
          },
          "property": {
            "type": "Identifier",
            "name": "in"
          },
          "computed": false
        },
        "operator": "/",
        "right": {
          "type": "Identifier",
          "name": "b"
        }
      }
    }
  ]
}
	`, ast)
}

func Test326(t *testing.T) {
	ast, err := compile("{}/=/")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "BlockStatement",
      "body": []
    },
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "Literal",
        "regexp": {
          "pattern": "="
        }
      }
    }
  ]
}
	`, ast)
}

func Test327(t *testing.T) {
	ast, err := compile("'use strict';\nobject.static();")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "Literal",
        "value": "use strict"
      }
    },
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "CallExpression",
        "callee": {
          "type": "MemberExpression",
          "object": {
            "type": "Identifier",
            "name": "object"
          },
          "property": {
            "type": "Identifier",
            "name": "static"
          },
          "computed": false
        },
        "arguments": []
      }
    }
  ]
}
	`, ast)
}

func Test328(t *testing.T) {
	ast, err := compile("foo <!--bar\n+baz")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "BinaryExpression",
        "left": {
          "type": "Identifier",
          "name": "foo"
        },
        "operator": "<",
        "right": {
          "type": "BinaryExpression",
          "loc": {
            "start": {
              "line": 1,
              "column": 5
            },
            "end": {
              "line": 2,
              "column": 4
            }
          },
          "left": {
            "type": "UnaryExpression",
            "loc": {
              "start": {
                "line": 1,
                "column": 5
              },
              "end": {
                "line": 1,
                "column": 11
              }
            },
            "operator": "!",
            "prefix": true,
            "argument": {
              "type": "UpdateExpression",
              "loc": {
                "start": {
                  "line": 1,
                  "column": 6
                },
                "end": {
                  "line": 1,
                  "column": 11
                }
              },
              "operator": "--",
              "prefix": true,
              "argument": {
                "type": "Identifier",
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 8
                  },
                  "end": {
                    "line": 1,
                    "column": 11
                  }
                },
                "name": "bar"
              }
            }
          },
          "operator": "+",
          "right": {
            "type": "Identifier",
            "loc": {
              "start": {
                "line": 2,
                "column": 1
              },
              "end": {
                "line": 2,
                "column": 4
              }
            },
            "name": "baz"
          }
        }
      }
    }
  ]
}
	`, ast)
}

func Test329(t *testing.T) {
	ast, err := compile("a = import(a)")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "AssignmentExpression",
        "operator": "=",
        "left": {
          "type": "Identifier",
          "name": "a"
        },
        "right": {
          "type": "ImportExpression",
          "source": {
            "type": "Identifier",
            "name": "a"
          }
        }
      }
    }
  ]
}
	`, ast)
}

func Test330(t *testing.T) {
	ast, err := compile("a = import.meta")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "AssignmentExpression",
        "operator": "=",
        "left": {
          "type": "Identifier",
          "name": "a"
        },
        "right": {
          "type": "MetaProperty",
          "meta": {
            "type": "Identifier",
            "name": "import"
          },
          "property": {
            "type": "Identifier",
            "name": "meta"
          }
        }
      }
    }
  ]
}
	`, ast)
}

func Test331(t *testing.T) {
	ast, err := compile("class A extends B { constructor() { super() } }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ClassDeclaration",
      "id": {
        "type": "Identifier",
        "name": "A"
      },
      "superClass": {
        "type": "Identifier",
        "name": "B"
      },
      "body": {
        "type": "ClassBody",
        "body": [
          {
            "type": "MethodDefinition",
            "static": false,
            "computed": false,
            "key": {
              "type": "Identifier",
              "name": "constructor"
            },
            "kind": "constructor",
            "value": {
              "type": "FunctionExpression",
              "id": null,
              "expression": false,
              "generator": false,
              "async": false,
              "params": [],
              "body": {
                "type": "BlockStatement",
                "body": [
                  {
                    "type": "ExpressionStatement",
                    "expression": {
                      "type": "CallExpression",
                      "callee": {
                        "type": "Super"
                      },
                      "arguments": [],
                      "optional": false
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

func Test332(t *testing.T) {
	ast, err := compile("`\\n${a}b`")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 9
    }
  },
  "body": [
    {
      "type": "ExpressionStatement",
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 9
        }
      },
      "expression": {
        "type": "TemplateLiteral",
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 9
          }
        },
        "expressions": [
          {
            "type": "Identifier",
            "loc": {
              "start": {
                "line": 1,
                "column": 5
              },
              "end": {
                "line": 1,
                "column": 6
              }
            },
            "name": "a"
          }
        ],
        "quasis": [
          {
            "type": "TemplateElement",
            "loc": {
              "start": {
                "line": 1,
                "column": 1
              },
              "end": {
                "line": 1,
                "column": 3
              }
            },
            "value": {
              "raw": "\\n"
            },
            "tail": false
          },
          {
            "type": "TemplateElement",
            "loc": {
              "start": {
                "line": 1,
                "column": 7
              },
              "end": {
                "line": 1,
                "column": 8
              }
            },
            "value": {
              "raw": "b"
            },
            "tail": true
          }
        ]
      }
    }
  ]
}
	`, ast)
}

func Test333(t *testing.T) {
	ast, err := compile("`${a}b${c}`")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 11
    }
  },
  "body": [
    {
      "type": "ExpressionStatement",
      "loc": {
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
        "type": "TemplateLiteral",
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 11
          }
        },
        "expressions": [
          {
            "type": "Identifier",
            "loc": {
              "start": {
                "line": 1,
                "column": 3
              },
              "end": {
                "line": 1,
                "column": 4
              }
            },
            "name": "a"
          },
          {
            "type": "Identifier",
            "loc": {
              "start": {
                "line": 1,
                "column": 8
              },
              "end": {
                "line": 1,
                "column": 9
              }
            },
            "name": "c"
          }
        ],
        "quasis": [
          {
            "type": "TemplateElement",
            "loc": {
              "start": {
                "line": 1,
                "column": 1
              },
              "end": {
                "line": 1,
                "column": 1
              }
            },
            "value": {
              "raw": "",
              "cooked": ""
            },
            "tail": false
          },
          {
            "type": "TemplateElement",
            "loc": {
              "start": {
                "line": 1,
                "column": 5
              },
              "end": {
                "line": 1,
                "column": 6
              }
            },
            "value": {
              "raw": "b",
              "cooked": "b"
            },
            "tail": false
          },
          {
            "type": "TemplateElement",
            "loc": {
              "start": {
                "line": 1,
                "column": 10
              },
              "end": {
                "line": 1,
                "column": 10
              }
            },
            "value": {
              "raw": "",
              "cooked": ""
            },
            "tail": true
          }
        ]
      }
    }
  ]
}
	`, ast)
}

func Test334(t *testing.T) {
	ast, err := compile(" `${a}\nb\n${`c${d}e`}`")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 3,
      "column": 12
    }
  },
  "body": [
    {
      "type": "ExpressionStatement",
      "loc": {
        "start": {
          "line": 1,
          "column": 1
        },
        "end": {
          "line": 3,
          "column": 12
        }
      },
      "expression": {
        "type": "TemplateLiteral",
        "loc": {
          "start": {
            "line": 1,
            "column": 1
          },
          "end": {
            "line": 3,
            "column": 12
          }
        },
        "expressions": [
          {
            "type": "Identifier",
            "loc": {
              "start": {
                "line": 1,
                "column": 4
              },
              "end": {
                "line": 1,
                "column": 5
              }
            },
            "name": "a"
          },
          {
            "type": "TemplateLiteral",
            "loc": {
              "start": {
                "line": 3,
                "column": 2
              },
              "end": {
                "line": 3,
                "column": 10
              }
            },
            "expressions": [
              {
                "type": "Identifier",
                "loc": {
                  "start": {
                    "line": 3,
                    "column": 6
                  },
                  "end": {
                    "line": 3,
                    "column": 7
                  }
                },
                "name": "d"
              }
            ],
            "quasis": [
              {
                "type": "TemplateElement",
                "loc": {
                  "start": {
                    "line": 3,
                    "column": 3
                  },
                  "end": {
                    "line": 3,
                    "column": 4
                  }
                },
                "value": {
                  "raw": "c",
                  "cooked": "c"
                },
                "tail": false
              },
              {
                "type": "TemplateElement",
                "loc": {
                  "start": {
                    "line": 3,
                    "column": 8
                  },
                  "end": {
                    "line": 3,
                    "column": 9
                  }
                },
                "value": {
                  "raw": "e",
                  "cooked": "e"
                },
                "tail": true
              }
            ]
          }
        ],
        "quasis": [
          {
            "type": "TemplateElement",
            "loc": {
              "start": {
                "line": 1,
                "column": 2
              },
              "end": {
                "line": 1,
                "column": 2
              }
            },
            "value": {
              "raw": "",
              "cooked": ""
            },
            "tail": false
          },
          {
            "type": "TemplateElement",
            "loc": {
              "start": {
                "line": 1,
                "column": 6
              },
              "end": {
                "line": 3,
                "column": 0
              }
            },
            "value": {
              "raw": "\nb\n",
              "cooked": "\nb\n"
            },
            "tail": false
          },
          {
            "type": "TemplateElement",
            "loc": {
              "start": {
                "line": 3,
                "column": 11
              },
              "end": {
                "line": 3,
                "column": 11
              }
            },
            "value": {
              "raw": "",
              "cooked": ""
            },
            "tail": true
          }
        ]
      }
    }
  ]
}
	`, ast)
}

func Test335(t *testing.T) {
	ast, err := compile("tag`a`")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 6,
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 6
    }
  },
  "body": [
    {
      "type": "ExpressionStatement",
      "start": 0,
      "end": 6,
      "loc": {
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
        "type": "TaggedTemplateExpression",
        "start": 0,
        "end": 6,
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 6
          }
        },
        "tag": {
          "type": "Identifier",
          "start": 0,
          "end": 3,
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 3
            }
          },
          "name": "tag"
        },
        "quasi": {
          "type": "TemplateLiteral",
          "start": 3,
          "end": 6,
          "loc": {
            "start": {
              "line": 1,
              "column": 3
            },
            "end": {
              "line": 1,
              "column": 6
            }
          },
          "expressions": [],
          "quasis": [
            {
              "type": "TemplateElement",
              "start": 4,
              "end": 5,
              "loc": {
                "start": {
                  "line": 1,
                  "column": 4
                },
                "end": {
                  "line": 1,
                  "column": 5
                }
              },
              "value": {
                "raw": "a",
                "cooked": "a"
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

func Test336(t *testing.T) {
	ast, err := compile("tag.a().b[c]`a`")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 15,
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 15
    }
  },
  "body": [
    {
      "type": "ExpressionStatement",
      "start": 0,
      "end": 15,
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 15
        }
      },
      "expression": {
        "type": "TaggedTemplateExpression",
        "start": 0,
        "end": 15,
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 15
          }
        },
        "tag": {
          "type": "MemberExpression",
          "start": 0,
          "end": 12,
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 12
            }
          },
          "object": {
            "type": "MemberExpression",
            "start": 0,
            "end": 9,
            "loc": {
              "start": {
                "line": 1,
                "column": 0
              },
              "end": {
                "line": 1,
                "column": 9
              }
            },
            "object": {
              "type": "CallExpression",
              "start": 0,
              "end": 7,
              "loc": {
                "start": {
                  "line": 1,
                  "column": 0
                },
                "end": {
                  "line": 1,
                  "column": 7
                }
              },
              "callee": {
                "type": "MemberExpression",
                "start": 0,
                "end": 5,
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 0
                  },
                  "end": {
                    "line": 1,
                    "column": 5
                  }
                },
                "object": {
                  "type": "Identifier",
                  "start": 0,
                  "end": 3,
                  "loc": {
                    "start": {
                      "line": 1,
                      "column": 0
                    },
                    "end": {
                      "line": 1,
                      "column": 3
                    }
                  },
                  "name": "tag"
                },
                "computed": false,
                "property": {
                  "type": "Identifier",
                  "start": 4,
                  "end": 5,
                  "loc": {
                    "start": {
                      "line": 1,
                      "column": 4
                    },
                    "end": {
                      "line": 1,
                      "column": 5
                    }
                  },
                  "name": "a"
                }
              },
              "arguments": []
            },
            "computed": false,
            "property": {
              "type": "Identifier",
              "start": 8,
              "end": 9,
              "loc": {
                "start": {
                  "line": 1,
                  "column": 8
                },
                "end": {
                  "line": 1,
                  "column": 9
                }
              },
              "name": "b"
            }
          },
          "computed": true,
          "property": {
            "type": "Identifier",
            "start": 10,
            "end": 11,
            "loc": {
              "start": {
                "line": 1,
                "column": 10
              },
              "end": {
                "line": 1,
                "column": 11
              }
            },
            "name": "c"
          }
        },
        "quasi": {
          "type": "TemplateLiteral",
          "start": 12,
          "end": 15,
          "loc": {
            "start": {
              "line": 1,
              "column": 12
            },
            "end": {
              "line": 1,
              "column": 15
            }
          },
          "expressions": [],
          "quasis": [
            {
              "type": "TemplateElement",
              "start": 13,
              "end": 14,
              "loc": {
                "start": {
                  "line": 1,
                  "column": 13
                },
                "end": {
                  "line": 1,
                  "column": 14
                }
              },
              "value": {
                "raw": "a",
                "cooked": "a"
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

func Test337(t *testing.T) {
	ast, err := compile("class A { #a; get false() {} b; }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 33,
  "body": [
    {
      "type": "ClassDeclaration",
      "start": 0,
      "end": 33,
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
        "end": 33,
        "body": [
          {
            "type": "PropertyDefinition",
            "start": 10,
            "end": 13,
            "static": false,
            "computed": false,
            "key": {
              "type": "PrivateIdentifier",
              "start": 10,
              "end": 12,
              "name": "a"
            },
            "value": null
          },
          {
            "type": "MethodDefinition",
            "start": 14,
            "end": 28,
            "static": false,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 18,
              "end": 23,
              "name": "false"
            },
            "kind": "get",
            "value": {
              "type": "FunctionExpression",
              "start": 23,
              "end": 28,
              "id": null,
              "expression": false,
              "generator": false,
              "async": false,
              "params": [],
              "body": {
                "type": "BlockStatement",
                "start": 26,
                "end": 28,
                "body": []
              }
            }
          },
          {
            "type": "PropertyDefinition",
            "start": 29,
            "end": 31,
            "static": false,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 29,
              "end": 30,
              "name": "b"
            },
            "value": null
          }
        ]
      }
    }
  ]
}
	`, ast)
}

func Test338(t *testing.T) {
	ast, err := compile("import {b} from \"a\"")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 19
    }
  },
  "body": [
    {
      "type": "ImportDeclaration",
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 19
        }
      },
      "specifiers": [
        {
          "type": "ImportSpecifier",
          "loc": {
            "start": {
              "line": 1,
              "column": 8
            },
            "end": {
              "line": 1,
              "column": 9
            }
          },
          "imported": {
            "type": "Identifier",
            "loc": {
              "start": {
                "line": 1,
                "column": 8
              },
              "end": {
                "line": 1,
                "column": 9
              }
            },
            "name": "b"
          },
          "local": {
            "type": "Identifier",
            "loc": {
              "start": {
                "line": 1,
                "column": 8
              },
              "end": {
                "line": 1,
                "column": 9
              }
            },
            "name": "b"
          }
        }
      ],
      "source": {
        "type": "Literal",
        "loc": {
          "start": {
            "line": 1,
            "column": 16
          },
          "end": {
            "line": 1,
            "column": 19
          }
        },
        "value": "a"
      }
    }
  ]
}
	`, ast)
}

func Test339(t *testing.T) {
	ast, err := compile("import {from as as} from \"a\"")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 28
    }
  },
  "body": [
    {
      "type": "ImportDeclaration",
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 28
        }
      },
      "specifiers": [
        {
          "type": "ImportSpecifier",
          "loc": {
            "start": {
              "line": 1,
              "column": 8
            },
            "end": {
              "line": 1,
              "column": 18
            }
          },
          "imported": {
            "type": "Identifier",
            "loc": {
              "start": {
                "line": 1,
                "column": 8
              },
              "end": {
                "line": 1,
                "column": 12
              }
            },
            "name": "from"
          },
          "local": {
            "type": "Identifier",
            "loc": {
              "start": {
                "line": 1,
                "column": 16
              },
              "end": {
                "line": 1,
                "column": 18
              }
            },
            "name": "as"
          }
        }
      ],
      "source": {
        "type": "Literal",
        "loc": {
          "start": {
            "line": 1,
            "column": 25
          },
          "end": {
            "line": 1,
            "column": 28
          }
        },
        "value": "a"
      }
    }
  ]
}
	`, ast)
}

func Test340(t *testing.T) {
	ast, err := compile("import b, * as c from \"a\"")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 25
    }
  },
  "body": [
    {
      "type": "ImportDeclaration",
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 25
        }
      },
      "specifiers": [
        {
          "type": "ImportDefaultSpecifier",
          "loc": {
            "start": {
              "line": 1,
              "column": 7
            },
            "end": {
              "line": 1,
              "column": 8
            }
          },
          "local": {
            "type": "Identifier",
            "loc": {
              "start": {
                "line": 1,
                "column": 7
              },
              "end": {
                "line": 1,
                "column": 8
              }
            },
            "name": "b"
          }
        },
        {
          "type": "ImportNamespaceSpecifier",
          "loc": {
            "start": {
              "line": 1,
              "column": 10
            },
            "end": {
              "line": 1,
              "column": 16
            }
          },
          "local": {
            "type": "Identifier",
            "loc": {
              "start": {
                "line": 1,
                "column": 15
              },
              "end": {
                "line": 1,
                "column": 16
              }
            },
            "name": "c"
          }
        }
      ],
      "source": {
        "type": "Literal",
        "loc": {
          "start": {
            "line": 1,
            "column": 22
          },
          "end": {
            "line": 1,
            "column": 25
          }
        },
        "value": "a"
      }
    }
  ]
}
	`, ast)
}

func Test341(t *testing.T) {
	ast, err := compile("export const a = 1;")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 19
    }
  },
  "body": [
    {
      "type": "ExportNamedDeclaration",
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 19
        }
      },
      "specifiers": [],
      "source": null,
      "declaration": {
        "type": "VariableDeclaration",
        "loc": {
          "start": {
            "line": 1,
            "column": 7
          },
          "end": {
            "line": 1,
            "column": 19
          }
        },
        "declarations": [
          {
            "type": "VariableDeclarator",
            "loc": {
              "start": {
                "line": 1,
                "column": 13
              },
              "end": {
                "line": 1,
                "column": 18
              }
            },
            "id": {
              "type": "Identifier",
              "loc": {
                "start": {
                  "line": 1,
                  "column": 13
                },
                "end": {
                  "line": 1,
                  "column": 14
                }
              },
              "name": "a"
            },
            "init": {
              "type": "Literal",
              "loc": {
                "start": {
                  "line": 1,
                  "column": 17
                },
                "end": {
                  "line": 1,
                  "column": 18
                }
              },
              "value": 1
            }
          }
        ],
        "kind": "const"
      }
    }
  ]
}
	`, ast)
}

func Test342(t *testing.T) {
	ast, err := compile("export { a as b  } from \"c\";")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 28
    }
  },
  "body": [
    {
      "type": "ExportNamedDeclaration",
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 28
        }
      },
      "specifiers": [
        {
          "type": "ExportSpecifier",
          "loc": {
            "start": {
              "line": 1,
              "column": 9
            },
            "end": {
              "line": 1,
              "column": 15
            }
          },
          "local": {
            "type": "Identifier",
            "loc": {
              "start": {
                "line": 1,
                "column": 9
              },
              "end": {
                "line": 1,
                "column": 10
              }
            },
            "name": "a"
          },
          "exported": {
            "type": "Identifier",
            "loc": {
              "start": {
                "line": 1,
                "column": 14
              },
              "end": {
                "line": 1,
                "column": 15
              }
            },
            "name": "b"
          }
        }
      ],
      "source": {
        "type": "Literal",
        "loc": {
          "start": {
            "line": 1,
            "column": 24
          },
          "end": {
            "line": 1,
            "column": 27
          }
        },
        "value": "c"
      },
      "declaration": null
    }
  ]
}
	`, ast)
}

func Test343(t *testing.T) {
	ast, err := compile("export async function a() {};")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 29
    }
  },
  "body": [
    {
      "type": "ExportNamedDeclaration",
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 28
        }
      },
      "specifiers": [],
      "source": null,
      "declaration": {
        "type": "FunctionDeclaration",
        "loc": {
          "start": {
            "line": 1,
            "column": 7
          },
          "end": {
            "line": 1,
            "column": 28
          }
        },
        "id": {
          "type": "Identifier",
          "loc": {
            "start": {
              "line": 1,
              "column": 22
            },
            "end": {
              "line": 1,
              "column": 23
            }
          },
          "name": "a"
        },
        "generator": false,
        "async": true,
        "params": [],
        "body": {
          "type": "BlockStatement",
          "loc": {
            "start": {
              "line": 1,
              "column": 26
            },
            "end": {
              "line": 1,
              "column": 28
            }
          },
          "body": []
        }
      }
    },
    {
      "type": "EmptyStatement",
      "loc": {
        "start": {
          "line": 1,
          "column": 28
        },
        "end": {
          "line": 1,
          "column": 29
        }
      }
    }
  ]
}
	`, ast)
}

func Test344(t *testing.T) {
	ast, err := compile("export * as b from 'a' ;")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 24
    }
  },
  "body": [
    {
      "type": "ExportAllDeclaration",
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 24
        }
      },
      "source": {
        "type": "Literal",
        "loc": {
          "start": {
            "line": 1,
            "column": 19
          },
          "end": {
            "line": 1,
            "column": 22
          }
        },
        "value": "a"
      },
      "exported": {
        "type": "Identifier",
        "loc": {
          "start": {
            "line": 1,
            "column": 12
          },
          "end": {
            "line": 1,
            "column": 13
          }
        },
        "name": "b"
      }
    }
  ]
}
	`, ast)
}

func Test345(t *testing.T) {
	ast, err := compile("const a = async () => {}")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 24
    }
  },
  "body": [
    {
      "type": "VariableDeclaration",
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 24
        }
      },
      "declarations": [
        {
          "type": "VariableDeclarator",
          "loc": {
            "start": {
              "line": 1,
              "column": 6
            },
            "end": {
              "line": 1,
              "column": 24
            }
          },
          "id": {
            "type": "Identifier",
            "loc": {
              "start": {
                "line": 1,
                "column": 6
              },
              "end": {
                "line": 1,
                "column": 7
              }
            },
            "name": "a"
          },
          "init": {
            "type": "ArrowFunctionExpression",
            "loc": {
              "start": {
                "line": 1,
                "column": 10
              },
              "end": {
                "line": 1,
                "column": 24
              }
            },
            "id": null,
            "generator": false,
            "async": true,
            "params": [],
            "body": {
              "type": "BlockStatement",
              "loc": {
                "start": {
                  "line": 1,
                  "column": 22
                },
                "end": {
                  "line": 1,
                  "column": 24
                }
              },
              "body": []
            }
          }
        }
      ],
      "kind": "const"
    }
  ]
}
	`, ast)
}

func Test346(t *testing.T) {
	ast, err := compile(`async function* a() {
  {
    yield 1
  }
}`)
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 5,
      "column": 1
    }
  },
  "body": [
    {
      "type": "FunctionDeclaration",
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 5,
          "column": 1
        }
      },
      "id": {
        "type": "Identifier",
        "loc": {
          "start": {
            "line": 1,
            "column": 16
          },
          "end": {
            "line": 1,
            "column": 17
          }
        },
        "name": "a"
      },
      "generator": true,
      "async": true,
      "params": [],
      "body": {
        "type": "BlockStatement",
        "loc": {
          "start": {
            "line": 1,
            "column": 20
          },
          "end": {
            "line": 5,
            "column": 1
          }
        },
        "body": [
          {
            "type": "BlockStatement",
            "loc": {
              "start": {
                "line": 2,
                "column": 2
              },
              "end": {
                "line": 4,
                "column": 3
              }
            },
            "body": [
              {
                "type": "ExpressionStatement",
                "loc": {
                  "start": {
                    "line": 3,
                    "column": 4
                  },
                  "end": {
                    "line": 3,
                    "column": 11
                  }
                },
                "expression": {
                  "type": "YieldExpression",
                  "loc": {
                    "start": {
                      "line": 3,
                      "column": 4
                    },
                    "end": {
                      "line": 3,
                      "column": 11
                    }
                  },
                  "delegate": false,
                  "argument": {
                    "type": "Literal",
                    "loc": {
                      "start": {
                        "line": 3,
                        "column": 10
                      },
                      "end": {
                        "line": 3,
                        "column": 11
                      }
                    },
                    "value": 1
                  }
                }
              }
            ]
          }
        ]
      }
    }
  ]
}
	`, ast)
}

func Test347(t *testing.T) {
	ast, err := compile("async function* a() { yield yield* 1 ? a : b }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 46
    }
  },
  "body": [
    {
      "type": "FunctionDeclaration",
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 46
        }
      },
      "id": {
        "type": "Identifier",
        "loc": {
          "start": {
            "line": 1,
            "column": 16
          },
          "end": {
            "line": 1,
            "column": 17
          }
        },
        "name": "a"
      },
      "generator": true,
      "async": true,
      "params": [],
      "body": {
        "type": "BlockStatement",
        "loc": {
          "start": {
            "line": 1,
            "column": 20
          },
          "end": {
            "line": 1,
            "column": 46
          }
        },
        "body": [
          {
            "type": "ExpressionStatement",
            "loc": {
              "start": {
                "line": 1,
                "column": 22
              },
              "end": {
                "line": 1,
                "column": 44
              }
            },
            "expression": {
              "type": "YieldExpression",
              "loc": {
                "start": {
                  "line": 1,
                  "column": 22
                },
                "end": {
                  "line": 1,
                  "column": 44
                }
              },
              "delegate": false,
              "argument": {
                "type": "YieldExpression",
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 28
                  },
                  "end": {
                    "line": 1,
                    "column": 44
                  }
                },
                "delegate": true,
                "argument": {
                  "type": "ConditionalExpression",
                  "loc": {
                    "start": {
                      "line": 1,
                      "column": 35
                    },
                    "end": {
                      "line": 1,
                      "column": 44
                    }
                  },
                  "test": {
                    "type": "Literal",
                    "loc": {
                      "start": {
                        "line": 1,
                        "column": 35
                      },
                      "end": {
                        "line": 1,
                        "column": 36
                      }
                    },
                    "value": 1
                  },
                  "consequent": {
                    "type": "Identifier",
                    "loc": {
                      "start": {
                        "line": 1,
                        "column": 39
                      },
                      "end": {
                        "line": 1,
                        "column": 40
                      }
                    },
                    "name": "a"
                  },
                  "alternate": {
                    "type": "Identifier",
                    "loc": {
                      "start": {
                        "line": 1,
                        "column": 43
                      },
                      "end": {
                        "line": 1,
                        "column": 44
                      }
                    },
                    "name": "b"
                  }
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

func Test348(t *testing.T) {
	ast, err := compile("+{} / 2")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 7,
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 7
    }
  },
  "body": [
    {
      "type": "ExpressionStatement",
      "start": 0,
      "end": 7,
      "loc": {
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
        "type": "BinaryExpression",
        "start": 0,
        "end": 7,
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 7
          }
        },
        "left": {
          "type": "UnaryExpression",
          "start": 0,
          "end": 3,
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 3
            }
          },
          "operator": "+",
          "prefix": true,
          "argument": {
            "type": "ObjectExpression",
            "start": 1,
            "end": 3,
            "loc": {
              "start": {
                "line": 1,
                "column": 1
              },
              "end": {
                "line": 1,
                "column": 3
              }
            },
            "properties": []
          }
        },
        "operator": "/",
        "right": {
          "type": "Literal",
          "start": 6,
          "end": 7,
          "loc": {
            "start": {
              "line": 1,
              "column": 6
            },
            "end": {
              "line": 1,
              "column": 7
            }
          },
          "value": 2
        }
      }
    }
  ]
}
	`, ast)
}

func Test349(t *testing.T) {
	ast, err := compile("{}\n/foo/")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 8,
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 2,
      "column": 5
    }
  },
  "body": [
    {
      "type": "BlockStatement",
      "start": 0,
      "end": 2,
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 2
        }
      },
      "body": [],
      "directives": []
    },
    {
      "type": "ExpressionStatement",
      "start": 3,
      "end": 8,
      "loc": {
        "start": {
          "line": 2,
          "column": 0
        },
        "end": {
          "line": 2,
          "column": 5
        }
      },
      "expression": {
        "type": "Literal",
        "start": 3,
        "end": 8,
        "loc": {
          "start": {
            "line": 2,
            "column": 0
          },
          "end": {
            "line": 2,
            "column": 5
          }
        },
        "regexp": {
          "pattern": "foo",
          "flags": ""
        }
      }
    }
  ]
}
	`, ast)
}

func Test350(t *testing.T) {
	ast, err := compile("({a: [1]}+[]) / 2")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 17,
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 17
    }
  },
  "body": [
    {
      "type": "ExpressionStatement",
      "start": 0,
      "end": 17,
      "loc": {
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
        "type": "BinaryExpression",
        "start": 0,
        "end": 17,
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 17
          }
        },
        "left": {
          "type": "BinaryExpression",
          "start": 1,
          "end": 12,
          "loc": {
            "start": {
              "line": 1,
              "column": 1
            },
            "end": {
              "line": 1,
              "column": 12
            }
          },
          "left": {
            "type": "ObjectExpression",
            "start": 1,
            "end": 9,
            "loc": {
              "start": {
                "line": 1,
                "column": 1
              },
              "end": {
                "line": 1,
                "column": 9
              }
            },
            "properties": [
              {
                "type": "Property",
                "start": 2,
                "end": 8,
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 2
                  },
                  "end": {
                    "line": 1,
                    "column": 8
                  }
                },
                "method": false,
                "key": {
                  "type": "Identifier",
                  "start": 2,
                  "end": 3,
                  "loc": {
                    "start": {
                      "line": 1,
                      "column": 2
                    },
                    "end": {
                      "line": 1,
                      "column": 3
                    }
                  },
                  "name": "a"
                },
                "computed": false,
                "shorthand": false,
                "value": {
                  "type": "ArrayExpression",
                  "start": 5,
                  "end": 8,
                  "loc": {
                    "start": {
                      "line": 1,
                      "column": 5
                    },
                    "end": {
                      "line": 1,
                      "column": 8
                    }
                  },
                  "elements": [
                    {
                      "type": "Literal",
                      "start": 6,
                      "end": 7,
                      "loc": {
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
                    }
                  ]
                }
              }
            ]
          },
          "operator": "+",
          "right": {
            "type": "ArrayExpression",
            "start": 10,
            "end": 12,
            "loc": {
              "start": {
                "line": 1,
                "column": 10
              },
              "end": {
                "line": 1,
                "column": 12
              }
            },
            "elements": []
          }
        },
        "operator": "/",
        "right": {
          "type": "Literal",
          "start": 16,
          "end": 17,
          "loc": {
            "start": {
              "line": 1,
              "column": 16
            },
            "end": {
              "line": 1,
              "column": 17
            }
          },
          "value": 2
        }
      }
    }
  ]
}
	`, ast)
}

func Test351(t *testing.T) {
	ast, err := compile("x = {foo: function x() {} / divide}")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 35,
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 35
    }
  },
  "body": [
    {
      "type": "ExpressionStatement",
      "start": 0,
      "end": 35,
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 35
        }
      },
      "expression": {
        "type": "AssignmentExpression",
        "start": 0,
        "end": 35,
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 35
          }
        },
        "operator": "=",
        "left": {
          "type": "Identifier",
          "start": 0,
          "end": 1,
          "loc": {
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
          "start": 4,
          "end": 35,
          "loc": {
            "start": {
              "line": 1,
              "column": 4
            },
            "end": {
              "line": 1,
              "column": 35
            }
          },
          "properties": [
            {
              "type": "Property",
              "start": 5,
              "end": 34,
              "loc": {
                "start": {
                  "line": 1,
                  "column": 5
                },
                "end": {
                  "line": 1,
                  "column": 34
                }
              },
              "method": false,
              "key": {
                "type": "Identifier",
                "start": 5,
                "end": 8,
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 5
                  },
                  "end": {
                    "line": 1,
                    "column": 8
                  }
                },
                "name": "foo"
              },
              "computed": false,
              "shorthand": false,
              "value": {
                "type": "BinaryExpression",
                "start": 10,
                "end": 34,
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 10
                  },
                  "end": {
                    "line": 1,
                    "column": 34
                  }
                },
                "left": {
                  "type": "FunctionExpression",
                  "start": 10,
                  "end": 25,
                  "loc": {
                    "start": {
                      "line": 1,
                      "column": 10
                    },
                    "end": {
                      "line": 1,
                      "column": 25
                    }
                  },
                  "id": {
                    "type": "Identifier",
                    "start": 19,
                    "end": 20,
                    "loc": {
                      "start": {
                        "line": 1,
                        "column": 19
                      },
                      "end": {
                        "line": 1,
                        "column": 20
                      }
                    },
                    "name": "x"
                  },
                  "generator": false,
                  "async": false,
                  "params": [],
                  "body": {
                    "type": "BlockStatement",
                    "start": 23,
                    "end": 25,
                    "loc": {
                      "start": {
                        "line": 1,
                        "column": 23
                      },
                      "end": {
                        "line": 1,
                        "column": 25
                      }
                    },
                    "body": [],
                    "directives": []
                  }
                },
                "operator": "/",
                "right": {
                  "type": "Identifier",
                  "start": 28,
                  "end": 34,
                  "loc": {
                    "start": {
                      "line": 1,
                      "column": 28
                    },
                    "end": {
                      "line": 1,
                      "column": 34
                    }
                  },
                  "name": "divide"
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

func Test352(t *testing.T) {
	ast, err := compile("switch(a) { case 1: {}\n/foo/ }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 30,
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 2,
      "column": 7
    }
  },
  "body": [
    {
      "type": "SwitchStatement",
      "start": 0,
      "end": 30,
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 2,
          "column": 7
        }
      },
      "discriminant": {
        "type": "Identifier",
        "start": 7,
        "end": 8,
        "loc": {
          "start": {
            "line": 1,
            "column": 7
          },
          "end": {
            "line": 1,
            "column": 8
          }
        },
        "name": "a"
      },
      "cases": [
        {
          "type": "SwitchCase",
          "start": 12,
          "end": 28,
          "loc": {
            "start": {
              "line": 1,
              "column": 12
            },
            "end": {
              "line": 2,
              "column": 5
            }
          },
          "consequent": [
            {
              "type": "BlockStatement",
              "start": 20,
              "end": 22,
              "loc": {
                "start": {
                  "line": 1,
                  "column": 20
                },
                "end": {
                  "line": 1,
                  "column": 22
                }
              },
              "body": [],
              "directives": []
            },
            {
              "type": "ExpressionStatement",
              "start": 23,
              "end": 28,
              "loc": {
                "start": {
                  "line": 2,
                  "column": 0
                },
                "end": {
                  "line": 2,
                  "column": 5
                }
              },
              "expression": {
                "type": "Literal",
                "start": 23,
                "end": 28,
                "loc": {
                  "start": {
                    "line": 2,
                    "column": 0
                  },
                  "end": {
                    "line": 2,
                    "column": 5
                  }
                },
                "regexp": {
                  "pattern": "foo",
                  "flags": ""
                }
              }
            }
          ],
          "test": {
            "type": "Literal",
            "start": 17,
            "end": 18,
            "loc": {
              "start": {
                "line": 1,
                "column": 17
              },
              "end": {
                "line": 1,
                "column": 18
              }
            },
            "value": 1
          }
        }
      ]
    }
  ]
}
	`, ast)
}

func Test353(t *testing.T) {
	ast, err := compile("+x++ / 2")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 8,
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 8
    }
  },
  "body": [
    {
      "type": "ExpressionStatement",
      "start": 0,
      "end": 8,
      "loc": {
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
        "type": "BinaryExpression",
        "start": 0,
        "end": 8,
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 8
          }
        },
        "left": {
          "type": "UnaryExpression",
          "start": 0,
          "end": 4,
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 4
            }
          },
          "operator": "+",
          "prefix": true,
          "argument": {
            "type": "UpdateExpression",
            "start": 1,
            "end": 4,
            "loc": {
              "start": {
                "line": 1,
                "column": 1
              },
              "end": {
                "line": 1,
                "column": 4
              }
            },
            "operator": "++",
            "prefix": false,
            "argument": {
              "type": "Identifier",
              "start": 1,
              "end": 2,
              "loc": {
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
            }
          }
        },
        "operator": "/",
        "right": {
          "type": "Literal",
          "start": 7,
          "end": 8,
          "loc": {
            "start": {
              "line": 1,
              "column": 7
            },
            "end": {
              "line": 1,
              "column": 8
            }
          },
          "value": 2
        }
      }
    }
  ]
}
	`, ast)
}

func Test354(t *testing.T) {
	ast, err := compile("+function f() {} / 3;")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 21,
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 21
    }
  },
  "body": [
    {
      "type": "ExpressionStatement",
      "start": 0,
      "end": 21,
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 21
        }
      },
      "expression": {
        "type": "BinaryExpression",
        "start": 0,
        "end": 20,
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 20
          }
        },
        "left": {
          "type": "UnaryExpression",
          "start": 0,
          "end": 16,
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 16
            }
          },
          "operator": "+",
          "prefix": true,
          "argument": {
            "type": "FunctionExpression",
            "start": 1,
            "end": 16,
            "loc": {
              "start": {
                "line": 1,
                "column": 1
              },
              "end": {
                "line": 1,
                "column": 16
              }
            },
            "id": {
              "type": "Identifier",
              "start": 10,
              "end": 11,
              "loc": {
                "start": {
                  "line": 1,
                  "column": 10
                },
                "end": {
                  "line": 1,
                  "column": 11
                }
              },
              "name": "f"
            },
            "generator": false,
            "async": false,
            "params": [],
            "body": {
              "type": "BlockStatement",
              "start": 14,
              "end": 16,
              "loc": {
                "start": {
                  "line": 1,
                  "column": 14
                },
                "end": {
                  "line": 1,
                  "column": 16
                }
              },
              "body": [],
              "directives": []
            }
          }
        },
        "operator": "/",
        "right": {
          "type": "Literal",
          "start": 19,
          "end": 20,
          "loc": {
            "start": {
              "line": 1,
              "column": 19
            },
            "end": {
              "line": 1,
              "column": 20
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

func Test355(t *testing.T) {
	ast, err := compile("foo; function f() {} /regexp/")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 29,
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 29
    }
  },
  "body": [
    {
      "type": "ExpressionStatement",
      "start": 0,
      "end": 4,
      "loc": {
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
        "type": "Identifier",
        "start": 0,
        "end": 3,
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 3
          }
        },
        "name": "foo"
      }
    },
    {
      "type": "FunctionDeclaration",
      "start": 5,
      "end": 20,
      "loc": {
        "start": {
          "line": 1,
          "column": 5
        },
        "end": {
          "line": 1,
          "column": 20
        }
      },
      "id": {
        "type": "Identifier",
        "start": 14,
        "end": 15,
        "loc": {
          "start": {
            "line": 1,
            "column": 14
          },
          "end": {
            "line": 1,
            "column": 15
          }
        },
        "name": "f"
      },
      "generator": false,
      "async": false,
      "params": [],
      "body": {
        "type": "BlockStatement",
        "start": 18,
        "end": 20,
        "loc": {
          "start": {
            "line": 1,
            "column": 18
          },
          "end": {
            "line": 1,
            "column": 20
          }
        },
        "body": [],
        "directives": []
      }
    },
    {
      "type": "ExpressionStatement",
      "start": 21,
      "end": 29,
      "loc": {
        "start": {
          "line": 1,
          "column": 21
        },
        "end": {
          "line": 1,
          "column": 29
        }
      },
      "expression": {
        "type": "Literal",
        "start": 21,
        "end": 29,
        "loc": {
          "start": {
            "line": 1,
            "column": 21
          },
          "end": {
            "line": 1,
            "column": 29
          }
        },
        "regexp": {
          "pattern": "regexp",
          "flags": ""
        }
      }
    }
  ]
}
	`, ast)
}

func Test356(t *testing.T) {
	ast, err := compile("{function f() {} /regexp/}")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 26,
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 26
    }
  },
  "body": [
    {
      "type": "BlockStatement",
      "start": 0,
      "end": 26,
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 26
        }
      },
      "body": [
        {
          "type": "FunctionDeclaration",
          "start": 1,
          "end": 16,
          "loc": {
            "start": {
              "line": 1,
              "column": 1
            },
            "end": {
              "line": 1,
              "column": 16
            }
          },
          "id": {
            "type": "Identifier",
            "start": 10,
            "end": 11,
            "loc": {
              "start": {
                "line": 1,
                "column": 10
              },
              "end": {
                "line": 1,
                "column": 11
              }
            },
            "name": "f"
          },
          "generator": false,
          "async": false,
          "params": [],
          "body": {
            "type": "BlockStatement",
            "start": 14,
            "end": 16,
            "loc": {
              "start": {
                "line": 1,
                "column": 14
              },
              "end": {
                "line": 1,
                "column": 16
              }
            },
            "body": [],
            "directives": []
          }
        },
        {
          "type": "ExpressionStatement",
          "start": 17,
          "end": 25,
          "loc": {
            "start": {
              "line": 1,
              "column": 17
            },
            "end": {
              "line": 1,
              "column": 25
            }
          },
          "expression": {
            "type": "Literal",
            "start": 17,
            "end": 25,
            "loc": {
              "start": {
                "line": 1,
                "column": 17
              },
              "end": {
                "line": 1,
                "column": 25
              }
            },
            "regexp": {
              "pattern": "regexp",
              "flags": ""
            }
          }
        }
      ]
    }
  ]
}
	`, ast)
}

func Test357(t *testing.T) {
	ast, err := compile(`function a() {
  return
  /a/
}`)
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 31,
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 4,
      "column": 1
    }
  },
  "body": [
    {
      "type": "FunctionDeclaration",
      "start": 0,
      "end": 31,
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 4,
          "column": 1
        }
      },
      "id": {
        "type": "Identifier",
        "start": 9,
        "end": 10,
        "loc": {
          "start": {
            "line": 1,
            "column": 9
          },
          "end": {
            "line": 1,
            "column": 10
          }
        },
        "name": "a"
      },
      "generator": false,
      "async": false,
      "params": [],
      "body": {
        "type": "BlockStatement",
        "start": 13,
        "end": 31,
        "loc": {
          "start": {
            "line": 1,
            "column": 13
          },
          "end": {
            "line": 4,
            "column": 1
          }
        },
        "body": [
          {
            "type": "ReturnStatement",
            "start": 17,
            "end": 23,
            "loc": {
              "start": {
                "line": 2,
                "column": 2
              },
              "end": {
                "line": 2,
                "column": 8
              }
            },
            "argument": null
          },
          {
            "type": "ExpressionStatement",
            "start": 26,
            "end": 29,
            "loc": {
              "start": {
                "line": 3,
                "column": 2
              },
              "end": {
                "line": 3,
                "column": 5
              }
            },
            "expression": {
              "type": "Literal",
              "start": 26,
              "end": 29,
              "loc": {
                "start": {
                  "line": 3,
                  "column": 2
                },
                "end": {
                  "line": 3,
                  "column": 5
                }
              },
              "regexp": {
                "pattern": "a",
                "flags": ""
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

func Test358(t *testing.T) {
	ast, err := compile("class a { async a() {} }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 24,
  "body": [
    {
      "type": "ClassDeclaration",
      "start": 0,
      "end": 24,
      "id": {
        "type": "Identifier",
        "start": 6,
        "end": 7,
        "name": "a"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "start": 8,
        "end": 24,
        "body": [
          {
            "type": "MethodDefinition",
            "start": 10,
            "end": 22,
            "static": false,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 16,
              "end": 17,
              "name": "a"
            },
            "kind": "method",
            "value": {
              "type": "FunctionExpression",
              "start": 17,
              "end": 22,
              "id": null,
              "expression": false,
              "generator": false,
              "async": true,
              "params": [],
              "body": {
                "type": "BlockStatement",
                "start": 20,
                "end": 22,
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

func Test359(t *testing.T) {
	ast, err := compile("class a { async* [1]() {} }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
        "name": "a"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "start": 8,
        "end": 27,
        "body": [
          {
            "type": "MethodDefinition",
            "start": 10,
            "end": 25,
            "static": false,
            "computed": true,
            "key": {
              "type": "Literal",
              "start": 18,
              "end": 19,
              "value": 1
            },
            "kind": "method",
            "value": {
              "type": "FunctionExpression",
              "start": 20,
              "end": 25,
              "id": null,
              "expression": false,
              "generator": true,
              "async": true,
              "params": [],
              "body": {
                "type": "BlockStatement",
                "start": 23,
                "end": 25,
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

func Test360(t *testing.T) {
	ast, err := compile("for (; function () {} / 1;);")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 28,
  "body": [
    {
      "type": "ForStatement",
      "start": 0,
      "end": 28,
      "init": null,
      "test": {
        "type": "BinaryExpression",
        "start": 7,
        "end": 25,
        "left": {
          "type": "FunctionExpression",
          "start": 7,
          "end": 21,
          "id": null,
          "expression": false,
          "generator": false,
          "async": false,
          "params": [],
          "body": {
            "type": "BlockStatement",
            "start": 19,
            "end": 21,
            "body": []
          }
        },
        "operator": "/",
        "right": {
          "type": "Literal",
          "start": 24,
          "end": 25,
          "value": 1
        }
      },
      "update": null,
      "body": {
        "type": "EmptyStatement",
        "start": 27,
        "end": 28
      }
    }
  ]
}
	`, ast)
}

func Test361(t *testing.T) {
	ast, err := compile("for (; class {} / 1;);")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 22,
  "body": [
    {
      "type": "ForStatement",
      "start": 0,
      "end": 22,
      "init": null,
      "test": {
        "type": "BinaryExpression",
        "start": 7,
        "end": 19,
        "left": {
          "type": "ClassExpression",
          "start": 7,
          "end": 15,
          "id": null,
          "superClass": null,
          "body": {
            "type": "ClassBody",
            "start": 13,
            "end": 15,
            "body": []
          }
        },
        "operator": "/",
        "right": {
          "type": "Literal",
          "start": 18,
          "end": 19,
          "value": 1
        }
      },
      "update": null,
      "body": {
        "type": "EmptyStatement",
        "start": 21,
        "end": 22
      }
    }
  ]
}
	`, ast)
}

func Test362(t *testing.T) {
	ast, err := compile("for (;; function () {} / 1);")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 28,
  "body": [
    {
      "type": "ForStatement",
      "start": 0,
      "end": 28,
      "init": null,
      "test": null,
      "update": {
        "type": "BinaryExpression",
        "start": 8,
        "end": 26,
        "left": {
          "type": "FunctionExpression",
          "start": 8,
          "end": 22,
          "id": null,
          "expression": false,
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
        "operator": "/",
        "right": {
          "type": "Literal",
          "start": 25,
          "end": 26,
          "value": 1
        }
      },
      "body": {
        "type": "EmptyStatement",
        "start": 27,
        "end": 28
      }
    }
  ]
}
	`, ast)
}

func Test363(t *testing.T) {
	ast, err := compile("for (;; class {} / 1);")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 22,
  "body": [
    {
      "type": "ForStatement",
      "start": 0,
      "end": 22,
      "init": null,
      "test": null,
      "update": {
        "type": "BinaryExpression",
        "start": 8,
        "end": 20,
        "left": {
          "type": "ClassExpression",
          "start": 8,
          "end": 16,
          "id": null,
          "superClass": null,
          "body": {
            "type": "ClassBody",
            "start": 14,
            "end": 16,
            "body": []
          }
        },
        "operator": "/",
        "right": {
          "type": "Literal",
          "start": 19,
          "end": 20,
          "value": 1
        }
      },
      "body": {
        "type": "EmptyStatement",
        "start": 21,
        "end": 22
      }
    }
  ]
}
	`, ast)
}

func Test364(t *testing.T) {
	ast, err := compile("class a { a = 1 }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 17,
  "body": [
    {
      "type": "ClassDeclaration",
      "start": 0,
      "end": 17,
      "id": {
        "type": "Identifier",
        "start": 6,
        "end": 7,
        "name": "a"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "start": 8,
        "end": 17,
        "body": [
          {
            "type": "PropertyDefinition",
            "start": 10,
            "end": 15,
            "static": false,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 10,
              "end": 11,
              "name": "a"
            },
            "value": {
              "type": "Literal",
              "start": 14,
              "end": 15,
              "value": 1
            }
          }
        ]
      }
    }
  ]
}
	`, ast)
}

func Test365(t *testing.T) {
	ast, err := compile("function a(a = 1) {}")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 20,
  "body": [
    {
      "type": "FunctionDeclaration",
      "start": 0,
      "end": 20,
      "id": {
        "type": "Identifier",
        "start": 9,
        "end": 10,
        "name": "a"
      },
      "generator": false,
      "async": false,
      "params": [
        {
          "type": "AssignmentPattern",
          "start": 11,
          "end": 16,
          "left": {
            "type": "Identifier",
            "start": 11,
            "end": 12,
            "name": "a"
          },
          "right": {
            "type": "Literal",
            "start": 15,
            "end": 16,
            "value": 1
          }
        }
      ],
      "body": {
        "type": "BlockStatement",
        "start": 18,
        "end": 20,
        "body": []
      }
    }
  ]
}
	`, ast)
}

func Test366(t *testing.T) {
	ast, err := compile("function a({ a: { b = 1 } }) {}")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 31,
  "body": [
    {
      "type": "FunctionDeclaration",
      "start": 0,
      "end": 31,
      "id": {
        "type": "Identifier",
        "start": 9,
        "end": 10,
        "name": "a"
      },
      "generator": false,
      "async": false,
      "params": [
        {
          "type": "ObjectPattern",
          "start": 11,
          "end": 27,
          "properties": [
            {
              "type": "Property",
              "start": 13,
              "end": 25,
              "method": false,
              "shorthand": false,
              "computed": false,
              "key": {
                "type": "Identifier",
                "start": 13,
                "end": 14,
                "name": "a"
              },
              "value": {
                "type": "ObjectPattern",
                "start": 16,
                "end": 25,
                "properties": [
                  {
                    "type": "Property",
                    "start": 18,
                    "end": 23,
                    "method": false,
                    "shorthand": true,
                    "computed": false,
                    "key": {
                      "type": "Identifier",
                      "start": 18,
                      "end": 19,
                      "name": "b"
                    },
                    "kind": "init",
                    "value": {
                      "type": "AssignmentPattern",
                      "start": 18,
                      "end": 23,
                      "left": {
                        "type": "Identifier",
                        "start": 18,
                        "end": 19,
                        "name": "b"
                      },
                      "right": {
                        "type": "Literal",
                        "start": 22,
                        "end": 23,
                        "value": 1
                      }
                    }
                  }
                ]
              },
              "kind": "init"
            }
          ]
        }
      ],
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

func Test367(t *testing.T) {
	ast, err := compile("function a({ a: { b = 1 } }) {}")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 31,
  "body": [
    {
      "type": "FunctionDeclaration",
      "start": 0,
      "end": 31,
      "id": {
        "type": "Identifier",
        "start": 9,
        "end": 10,
        "name": "a"
      },
      "generator": false,
      "async": false,
      "params": [
        {
          "type": "ObjectPattern",
          "start": 11,
          "end": 27,
          "properties": [
            {
              "type": "Property",
              "start": 13,
              "end": 25,
              "method": false,
              "shorthand": false,
              "computed": false,
              "key": {
                "type": "Identifier",
                "start": 13,
                "end": 14,
                "name": "a"
              },
              "value": {
                "type": "ObjectPattern",
                "start": 16,
                "end": 25,
                "properties": [
                  {
                    "type": "Property",
                    "start": 18,
                    "end": 23,
                    "method": false,
                    "shorthand": true,
                    "computed": false,
                    "key": {
                      "type": "Identifier",
                      "start": 18,
                      "end": 19,
                      "name": "b"
                    },
                    "kind": "init",
                    "value": {
                      "type": "AssignmentPattern",
                      "start": 18,
                      "end": 23,
                      "left": {
                        "type": "Identifier",
                        "start": 18,
                        "end": 19,
                        "name": "b"
                      },
                      "right": {
                        "type": "Literal",
                        "start": 22,
                        "end": 23,
                        "value": 1
                      }
                    }
                  }
                ]
              },
              "kind": "init"
            }
          ]
        }
      ],
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

func Test368(t *testing.T) {
	ast, err := compile("function a({ a: { b: { c = 1 } } }) {}")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 38,
  "body": [
    {
      "type": "FunctionDeclaration",
      "start": 0,
      "end": 38,
      "id": {
        "type": "Identifier",
        "start": 9,
        "end": 10,
        "name": "a"
      },
      "generator": false,
      "async": false,
      "params": [
        {
          "type": "ObjectPattern",
          "start": 11,
          "end": 34,
          "properties": [
            {
              "type": "Property",
              "start": 13,
              "end": 32,
              "method": false,
              "shorthand": false,
              "computed": false,
              "key": {
                "type": "Identifier",
                "start": 13,
                "end": 14,
                "name": "a"
              },
              "value": {
                "type": "ObjectPattern",
                "start": 16,
                "end": 32,
                "properties": [
                  {
                    "type": "Property",
                    "start": 18,
                    "end": 30,
                    "method": false,
                    "shorthand": false,
                    "computed": false,
                    "key": {
                      "type": "Identifier",
                      "start": 18,
                      "end": 19,
                      "name": "b"
                    },
                    "value": {
                      "type": "ObjectPattern",
                      "start": 21,
                      "end": 30,
                      "properties": [
                        {
                          "type": "Property",
                          "start": 23,
                          "end": 28,
                          "method": false,
                          "shorthand": true,
                          "computed": false,
                          "key": {
                            "type": "Identifier",
                            "start": 23,
                            "end": 24,
                            "name": "c"
                          },
                          "kind": "init",
                          "value": {
                            "type": "AssignmentPattern",
                            "start": 23,
                            "end": 28,
                            "left": {
                              "type": "Identifier",
                              "start": 23,
                              "end": 24,
                              "name": "c"
                            },
                            "right": {
                              "type": "Literal",
                              "start": 27,
                              "end": 28,
                              "value": 1
                            }
                          }
                        }
                      ]
                    },
                    "kind": "init"
                  }
                ]
              },
              "kind": "init"
            }
          ]
        }
      ],
      "body": {
        "type": "BlockStatement",
        "start": 36,
        "end": 38,
        "body": []
      }
    }
  ]
}
	`, ast)
}

func Test369(t *testing.T) {
	ast, err := compile("let a = (a = 1) => {}")
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
            "generator": false,
            "async": false,
            "params": [
              {
                "type": "AssignmentPattern",
                "start": 9,
                "end": 14,
                "left": {
                  "type": "Identifier",
                  "start": 9,
                  "end": 10,
                  "name": "a"
                },
                "right": {
                  "type": "Literal",
                  "start": 13,
                  "end": 14,
                  "value": 1
                }
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

func Test370(t *testing.T) {
	ast, err := compile("let a = ({ a = 1 }) => {}")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 25,
  "body": [
    {
      "type": "VariableDeclaration",
      "start": 0,
      "end": 25,
      "declarations": [
        {
          "type": "VariableDeclarator",
          "start": 4,
          "end": 25,
          "id": {
            "type": "Identifier",
            "start": 4,
            "end": 5,
            "name": "a"
          },
          "init": {
            "type": "ArrowFunctionExpression",
            "start": 8,
            "end": 25,
            "id": null,
            "generator": false,
            "async": false,
            "params": [
              {
                "type": "ObjectPattern",
                "start": 9,
                "end": 18,
                "properties": [
                  {
                    "type": "Property",
                    "start": 11,
                    "end": 16,
                    "method": false,
                    "shorthand": true,
                    "computed": false,
                    "key": {
                      "type": "Identifier",
                      "start": 11,
                      "end": 12,
                      "name": "a"
                    },
                    "kind": "init",
                    "value": {
                      "type": "AssignmentPattern",
                      "start": 11,
                      "end": 16,
                      "left": {
                        "type": "Identifier",
                        "start": 11,
                        "end": 12,
                        "name": "a"
                      },
                      "right": {
                        "type": "Literal",
                        "start": 15,
                        "end": 16,
                        "value": 1
                      }
                    }
                  }
                ]
              }
            ],
            "body": {
              "type": "BlockStatement",
              "start": 23,
              "end": 25,
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

func Test371(t *testing.T) {
	ast, err := compile("let a = ({ a: { b = 1 } }) => {}")
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
            "end": 5,
            "name": "a"
          },
          "init": {
            "type": "ArrowFunctionExpression",
            "start": 8,
            "end": 32,
            "id": null,
            "generator": false,
            "async": false,
            "params": [
              {
                "type": "ObjectPattern",
                "start": 9,
                "end": 25,
                "properties": [
                  {
                    "type": "Property",
                    "start": 11,
                    "end": 23,
                    "method": false,
                    "shorthand": false,
                    "computed": false,
                    "key": {
                      "type": "Identifier",
                      "start": 11,
                      "end": 12,
                      "name": "a"
                    },
                    "value": {
                      "type": "ObjectPattern",
                      "start": 14,
                      "end": 23,
                      "properties": [
                        {
                          "type": "Property",
                          "start": 16,
                          "end": 21,
                          "method": false,
                          "shorthand": true,
                          "computed": false,
                          "key": {
                            "type": "Identifier",
                            "start": 16,
                            "end": 17,
                            "name": "b"
                          },
                          "kind": "init",
                          "value": {
                            "type": "AssignmentPattern",
                            "start": 16,
                            "end": 21,
                            "left": {
                              "type": "Identifier",
                              "start": 16,
                              "end": 17,
                              "name": "b"
                            },
                            "right": {
                              "type": "Literal",
                              "start": 20,
                              "end": 21,
                              "value": 1
                            }
                          }
                        }
                      ]
                    },
                    "kind": "init"
                  }
                ]
              }
            ],
            "body": {
              "type": "BlockStatement",
              "start": 30,
              "end": 32,
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

func Test372(t *testing.T) {
	ast, err := compile("function a([ a = 1 ]) {}")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
        "start": 9,
        "end": 10,
        "name": "a"
      },
      "generator": false,
      "async": false,
      "params": [
        {
          "type": "ArrayPattern",
          "start": 11,
          "end": 20,
          "elements": [
            {
              "type": "AssignmentPattern",
              "start": 13,
              "end": 18,
              "left": {
                "type": "Identifier",
                "start": 13,
                "end": 14,
                "name": "a"
              },
              "right": {
                "type": "Literal",
                "start": 17,
                "end": 18,
                "value": 1
              }
            }
          ]
        }
      ],
      "body": {
        "type": "BlockStatement",
        "start": 22,
        "end": 24,
        "body": []
      }
    }
  ]
}
	`, ast)
}

func Test373(t *testing.T) {
	ast, err := compile("({...obj})")
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
        "type": "ObjectExpression",
        "start": 1,
        "end": 9,
        "properties": [
          {
            "type": "SpreadElement",
            "start": 2,
            "end": 8,
            "argument": {
              "type": "Identifier",
              "start": 5,
              "end": 8,
              "name": "obj"
            }
          }
        ]
      }
    }
  ]
}
	`, ast)
}

func Test374(t *testing.T) {
	ast, err := compile("({...obj1,})")
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
        "type": "ObjectExpression",
        "start": 1,
        "end": 11,
        "properties": [
          {
            "type": "SpreadElement",
            "start": 2,
            "end": 9,
            "argument": {
              "type": "Identifier",
              "start": 5,
              "end": 9,
              "name": "obj1"
            }
          }
        ]
      }
    }
  ]
}
	`, ast)
}

func Test375(t *testing.T) {
	ast, err := compile("({...obj1,...obj2})")
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
      "end": 19,
      "expression": {
        "type": "ObjectExpression",
        "start": 1,
        "end": 18,
        "properties": [
          {
            "type": "SpreadElement",
            "start": 2,
            "end": 9,
            "argument": {
              "type": "Identifier",
              "start": 5,
              "end": 9,
              "name": "obj1"
            }
          },
          {
            "type": "SpreadElement",
            "start": 10,
            "end": 17,
            "argument": {
              "type": "Identifier",
              "start": 13,
              "end": 17,
              "name": "obj2"
            }
          }
        ]
      }
    }
  ]
}
	`, ast)
}

func Test376(t *testing.T) {
	ast, err := compile("({a,...obj1,b:1,...obj2,c:2})")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 29,
  "body": [
    {
      "type": "ExpressionStatement",
      "start": 0,
      "end": 29,
      "expression": {
        "type": "ObjectExpression",
        "start": 1,
        "end": 28,
        "properties": [
          {
            "type": "Property",
            "start": 2,
            "end": 3,
            "method": false,
            "shorthand": true,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 2,
              "end": 3,
              "name": "a"
            },
            "kind": "init",
            "value": {
              "type": "Identifier",
              "start": 2,
              "end": 3,
              "name": "a"
            }
          },
          {
            "type": "SpreadElement",
            "start": 4,
            "end": 11,
            "argument": {
              "type": "Identifier",
              "start": 7,
              "end": 11,
              "name": "obj1"
            }
          },
          {
            "type": "Property",
            "start": 12,
            "end": 15,
            "method": false,
            "shorthand": false,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 12,
              "end": 13,
              "name": "b"
            },
            "value": {
              "type": "Literal",
              "start": 14,
              "end": 15,
              "value": 1
            },
            "kind": "init"
          },
          {
            "type": "SpreadElement",
            "start": 16,
            "end": 23,
            "argument": {
              "type": "Identifier",
              "start": 19,
              "end": 23,
              "name": "obj2"
            }
          },
          {
            "type": "Property",
            "start": 24,
            "end": 27,
            "method": false,
            "shorthand": false,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 24,
              "end": 25,
              "name": "c"
            },
            "value": {
              "type": "Literal",
              "start": 26,
              "end": 27,
              "value": 2
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

func Test377(t *testing.T) {
	ast, err := compile("({...(obj)})")
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
        "type": "ObjectExpression",
        "start": 1,
        "end": 11,
        "properties": [
          {
            "type": "SpreadElement",
            "start": 2,
            "end": 10,
            "argument": {
              "type": "Identifier",
              "start": 6,
              "end": 9,
              "name": "obj"
            }
          }
        ]
      }
    }
  ]
}
	`, ast)
}

func Test378(t *testing.T) {
	ast, err := compile("({...a,b,c})")
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
        "type": "ObjectExpression",
        "start": 1,
        "end": 11,
        "properties": [
          {
            "type": "SpreadElement",
            "start": 2,
            "end": 6,
            "argument": {
              "type": "Identifier",
              "start": 5,
              "end": 6,
              "name": "a"
            }
          },
          {
            "type": "Property",
            "start": 7,
            "end": 8,
            "method": false,
            "shorthand": true,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 7,
              "end": 8,
              "name": "b"
            },
            "kind": "init",
            "value": {
              "type": "Identifier",
              "start": 7,
              "end": 8,
              "name": "b"
            }
          },
          {
            "type": "Property",
            "start": 9,
            "end": 10,
            "method": false,
            "shorthand": true,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 9,
              "end": 10,
              "name": "c"
            },
            "kind": "init",
            "value": {
              "type": "Identifier",
              "start": 9,
              "end": 10,
              "name": "c"
            }
          }
        ]
      }
    }
  ]
}
	`, ast)
}

func Test379(t *testing.T) {
	ast, err := compile("({...(a,b),c})")
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
            "type": "SpreadElement",
            "start": 2,
            "end": 10,
            "argument": {
              "type": "SequenceExpression",
              "start": 5,
              "end": 10,
              "expressions": [
                {
                  "type": "Identifier",
                  "start": 6,
                  "end": 7,
                  "name": "a"
                },
                {
                  "type": "Identifier",
                  "start": 8,
                  "end": 9,
                  "name": "b"
                }
              ]
            }
          },
          {
            "type": "Property",
            "start": 11,
            "end": 12,
            "method": false,
            "shorthand": true,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 11,
              "end": 12,
              "name": "c"
            },
            "kind": "init",
            "value": {
              "type": "Identifier",
              "start": 11,
              "end": 12,
              "name": "c"
            }
          }
        ]
      }
    }
  ]
}
	`, ast)
}

func Test380(t *testing.T) {
	ast, err := compile("({...obj} = foo)")
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
        "type": "AssignmentExpression",
        "start": 1,
        "end": 15,
        "operator": "=",
        "left": {
          "type": "ObjectPattern",
          "start": 1,
          "end": 9,
          "properties": [
            {
              "type": "RestElement",
              "start": 2,
              "end": 8,
              "argument": {
                "type": "Identifier",
                "start": 5,
                "end": 8,
                "name": "obj"
              }
            }
          ]
        },
        "right": {
          "type": "Identifier",
          "start": 12,
          "end": 15,
          "name": "foo"
        }
      }
    }
  ]
}
	`, ast)
}

func Test381(t *testing.T) {
	ast, err := compile("({a,...obj} = foo)")
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
          "end": 11,
          "properties": [
            {
              "type": "Property",
              "start": 2,
              "end": 3,
              "method": false,
              "shorthand": true,
              "computed": false,
              "key": {
                "type": "Identifier",
                "start": 2,
                "end": 3,
                "name": "a"
              },
              "kind": "init",
              "value": {
                "type": "Identifier",
                "start": 2,
                "end": 3,
                "name": "a"
              }
            },
            {
              "type": "RestElement",
              "start": 4,
              "end": 10,
              "argument": {
                "type": "Identifier",
                "start": 7,
                "end": 10,
                "name": "obj"
              }
            }
          ]
        },
        "right": {
          "type": "Identifier",
          "start": 14,
          "end": 17,
          "name": "foo"
        }
      }
    }
  ]
}
	`, ast)
}

func Test382(t *testing.T) {
	ast, err := compile("({a:b,...obj} = foo)")
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
        "type": "AssignmentExpression",
        "start": 1,
        "end": 19,
        "operator": "=",
        "left": {
          "type": "ObjectPattern",
          "start": 1,
          "end": 13,
          "properties": [
            {
              "type": "Property",
              "start": 2,
              "end": 5,
              "method": false,
              "shorthand": false,
              "computed": false,
              "key": {
                "type": "Identifier",
                "start": 2,
                "end": 3,
                "name": "a"
              },
              "value": {
                "type": "Identifier",
                "start": 4,
                "end": 5,
                "name": "b"
              },
              "kind": "init"
            },
            {
              "type": "RestElement",
              "start": 6,
              "end": 12,
              "argument": {
                "type": "Identifier",
                "start": 9,
                "end": 12,
                "name": "obj"
              }
            }
          ]
        },
        "right": {
          "type": "Identifier",
          "start": 16,
          "end": 19,
          "name": "foo"
        }
      }
    }
  ]
}
	`, ast)
}

func Test383(t *testing.T) {
	ast, err := compile("({...obj}) => {}")
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
        "type": "ArrowFunctionExpression",
        "start": 0,
        "end": 16,
        "id": null,
        "expression": false,
        "generator": false,
        "async": false,
        "params": [
          {
            "type": "ObjectPattern",
            "start": 1,
            "end": 9,
            "properties": [
              {
                "type": "RestElement",
                "start": 2,
                "end": 8,
                "argument": {
                  "type": "Identifier",
                  "start": 5,
                  "end": 8,
                  "name": "obj"
                }
              }
            ]
          }
        ],
        "body": {
          "type": "BlockStatement",
          "start": 14,
          "end": 16,
          "body": []
        }
      }
    }
  ]
}
	`, ast)
}

func Test384(t *testing.T) {
	ast, err := compile("({...obj} = {}) => {}")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
        "type": "ArrowFunctionExpression",
        "start": 0,
        "end": 21,
        "id": null,
        "expression": false,
        "generator": false,
        "async": false,
        "params": [
          {
            "type": "AssignmentPattern",
            "start": 1,
            "end": 14,
            "left": {
              "type": "ObjectPattern",
              "start": 1,
              "end": 9,
              "properties": [
                {
                  "type": "RestElement",
                  "start": 2,
                  "end": 8,
                  "argument": {
                    "type": "Identifier",
                    "start": 5,
                    "end": 8,
                    "name": "obj"
                  }
                }
              ]
            },
            "right": {
              "type": "ObjectExpression",
              "start": 12,
              "end": 14,
              "properties": []
            }
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
  ]
}
	`, ast)
}

func Test385(t *testing.T) {
	ast, err := compile("({a,...obj}) => {}")
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
        "expression": false,
        "generator": false,
        "async": false,
        "params": [
          {
            "type": "ObjectPattern",
            "start": 1,
            "end": 11,
            "properties": [
              {
                "type": "Property",
                "start": 2,
                "end": 3,
                "method": false,
                "shorthand": true,
                "computed": false,
                "key": {
                  "type": "Identifier",
                  "start": 2,
                  "end": 3,
                  "name": "a"
                },
                "kind": "init",
                "value": {
                  "type": "Identifier",
                  "start": 2,
                  "end": 3,
                  "name": "a"
                }
              },
              {
                "type": "RestElement",
                "start": 4,
                "end": 10,
                "argument": {
                  "type": "Identifier",
                  "start": 7,
                  "end": 10,
                  "name": "obj"
                }
              }
            ]
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

func Test386(t *testing.T) {
	ast, err := compile("({a:b,...obj}) => {}")
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
        "type": "ArrowFunctionExpression",
        "start": 0,
        "end": 20,
        "id": null,
        "expression": false,
        "generator": false,
        "async": false,
        "params": [
          {
            "type": "ObjectPattern",
            "start": 1,
            "end": 13,
            "properties": [
              {
                "type": "Property",
                "start": 2,
                "end": 5,
                "method": false,
                "shorthand": false,
                "computed": false,
                "key": {
                  "type": "Identifier",
                  "start": 2,
                  "end": 3,
                  "name": "a"
                },
                "value": {
                  "type": "Identifier",
                  "start": 4,
                  "end": 5,
                  "name": "b"
                },
                "kind": "init"
              },
              {
                "type": "RestElement",
                "start": 6,
                "end": 12,
                "argument": {
                  "type": "Identifier",
                  "start": 9,
                  "end": 12,
                  "name": "obj"
                }
              }
            ]
          }
        ],
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
	`, ast)
}

func Test387(t *testing.T) {
	ast, err := compile("let {...obj1} = foo")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 19,
  "body": [
    {
      "type": "VariableDeclaration",
      "start": 0,
      "end": 19,
      "declarations": [
        {
          "type": "VariableDeclarator",
          "start": 4,
          "end": 19,
          "id": {
            "type": "ObjectPattern",
            "start": 4,
            "end": 13,
            "properties": [
              {
                "type": "RestElement",
                "start": 5,
                "end": 12,
                "argument": {
                  "type": "Identifier",
                  "start": 8,
                  "end": 12,
                  "name": "obj1"
                }
              }
            ]
          },
          "init": {
            "type": "Identifier",
            "start": 16,
            "end": 19,
            "name": "foo"
          }
        }
      ],
      "kind": "let"
    }
  ]
}
	`, ast)
}

func Test388(t *testing.T) {
	ast, err := compile("let {...obj1} =  {...a,}")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 24,
  "body": [
    {
      "type": "VariableDeclaration",
      "start": 0,
      "end": 24,
      "declarations": [
        {
          "type": "VariableDeclarator",
          "start": 4,
          "end": 24,
          "id": {
            "type": "ObjectPattern",
            "start": 4,
            "end": 13,
            "properties": [
              {
                "type": "RestElement",
                "start": 5,
                "end": 12,
                "argument": {
                  "type": "Identifier",
                  "start": 8,
                  "end": 12,
                  "name": "obj1"
                }
              }
            ]
          },
          "init": {
            "type": "ObjectExpression",
            "start": 17,
            "end": 24,
            "properties": [
              {
                "type": "SpreadElement",
                "start": 18,
                "end": 22,
                "argument": {
                  "type": "Identifier",
                  "start": 21,
                  "end": 22,
                  "name": "a"
                }
              }
            ]
          }
        }
      ],
      "kind": "let"
    }
  ]
}
	`, ast)
}

func Test389(t *testing.T) {
	// this sample raises error`Parenthesized pattern` in acorn
	// however it's legal in babel-parser, chrome and firefox
	ast, err := compile("({...(obj)} = foo)")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 18,
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 18
    }
  },
  "body": [
    {
      "type": "ExpressionStatement",
      "start": 0,
      "end": 18,
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 18
        }
      },
      "expression": {
        "type": "AssignmentExpression",
        "start": 1,
        "end": 17,
        "loc": {
          "start": {
            "line": 1,
            "column": 1
          },
          "end": {
            "line": 1,
            "column": 17
          }
        },
        "operator": "=",
        "left": {
          "type": "ObjectPattern",
          "start": 1,
          "end": 11,
          "loc": {
            "start": {
              "line": 1,
              "column": 1
            },
            "end": {
              "line": 1,
              "column": 11
            }
          },
          "properties": [
            {
              "type": "RestElement",
              "start": 2,
              "end": 10,
              "loc": {
                "start": {
                  "line": 1,
                  "column": 2
                },
                "end": {
                  "line": 1,
                  "column": 10
                }
              },
              "argument": {
                "type": "Identifier",
                "start": 6,
                "end": 9,
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 6
                  },
                  "end": {
                    "line": 1,
                    "column": 9
                  }
                },
                "name": "obj"
              }
            }
          ]
        },
        "right": {
          "type": "Identifier",
          "start": 14,
          "end": 17,
          "loc": {
            "start": {
              "line": 1,
              "column": 14
            },
            "end": {
              "line": 1,
              "column": 17
            }
          },
          "name": "foo"
        }
      }
    }
  ]
}
	`, ast)
}

func Test390(t *testing.T) {
	ast, err := compile("let z = {...x}")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 14,
  "body": [
    {
      "type": "VariableDeclaration",
      "start": 0,
      "end": 14,
      "declarations": [
        {
          "type": "VariableDeclarator",
          "start": 4,
          "end": 14,
          "id": {
            "type": "Identifier",
            "start": 4,
            "end": 5,
            "name": "z"
          },
          "init": {
            "type": "ObjectExpression",
            "start": 8,
            "end": 14,
            "properties": [
              {
                "type": "SpreadElement",
                "start": 9,
                "end": 13,
                "argument": {
                  "type": "Identifier",
                  "start": 12,
                  "end": 13,
                  "name": "x"
                }
              }
            ]
          }
        }
      ],
      "kind": "let"
    }
  ]
}
	`, ast)
}

func Test391(t *testing.T) {
	ast, err := compile("z = {x, ...y}")
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
        "type": "AssignmentExpression",
        "start": 0,
        "end": 13,
        "operator": "=",
        "left": {
          "type": "Identifier",
          "start": 0,
          "end": 1,
          "name": "z"
        },
        "right": {
          "type": "ObjectExpression",
          "start": 4,
          "end": 13,
          "properties": [
            {
              "type": "Property",
              "start": 5,
              "end": 6,
              "method": false,
              "shorthand": true,
              "computed": false,
              "key": {
                "type": "Identifier",
                "start": 5,
                "end": 6,
                "name": "x"
              },
              "kind": "init",
              "value": {
                "type": "Identifier",
                "start": 5,
                "end": 6,
                "name": "x"
              }
            },
            {
              "type": "SpreadElement",
              "start": 8,
              "end": 12,
              "argument": {
                "type": "Identifier",
                "start": 11,
                "end": 12,
                "name": "y"
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

func Test392(t *testing.T) {
	ast, err := compile("({x, ...y, a, ...b, c})")
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
            "start": 2,
            "end": 3,
            "method": false,
            "shorthand": true,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 2,
              "end": 3,
              "name": "x"
            },
            "kind": "init",
            "value": {
              "type": "Identifier",
              "start": 2,
              "end": 3,
              "name": "x"
            }
          },
          {
            "type": "SpreadElement",
            "start": 5,
            "end": 9,
            "argument": {
              "type": "Identifier",
              "start": 8,
              "end": 9,
              "name": "y"
            }
          },
          {
            "type": "Property",
            "start": 11,
            "end": 12,
            "method": false,
            "shorthand": true,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 11,
              "end": 12,
              "name": "a"
            },
            "kind": "init",
            "value": {
              "type": "Identifier",
              "start": 11,
              "end": 12,
              "name": "a"
            }
          },
          {
            "type": "SpreadElement",
            "start": 14,
            "end": 18,
            "argument": {
              "type": "Identifier",
              "start": 17,
              "end": 18,
              "name": "b"
            }
          },
          {
            "type": "Property",
            "start": 20,
            "end": 21,
            "method": false,
            "shorthand": true,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 20,
              "end": 21,
              "name": "c"
            },
            "kind": "init",
            "value": {
              "type": "Identifier",
              "start": 20,
              "end": 21,
              "name": "c"
            }
          }
        ]
      }
    }
  ]
}
	`, ast)
}

func Test393(t *testing.T) {
	ast, err := compile("var someObject = { someKey: { ...mapGetters([ 'some_val_1', 'some_val_2' ]) } }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 79,
  "body": [
    {
      "type": "VariableDeclaration",
      "start": 0,
      "end": 79,
      "declarations": [
        {
          "type": "VariableDeclarator",
          "start": 4,
          "end": 79,
          "id": {
            "type": "Identifier",
            "start": 4,
            "end": 14,
            "name": "someObject"
          },
          "init": {
            "type": "ObjectExpression",
            "start": 17,
            "end": 79,
            "properties": [
              {
                "type": "Property",
                "start": 19,
                "end": 77,
                "method": false,
                "shorthand": false,
                "computed": false,
                "key": {
                  "type": "Identifier",
                  "start": 19,
                  "end": 26,
                  "name": "someKey"
                },
                "value": {
                  "type": "ObjectExpression",
                  "start": 28,
                  "end": 77,
                  "properties": [
                    {
                      "type": "SpreadElement",
                      "start": 30,
                      "end": 75,
                      "argument": {
                        "type": "CallExpression",
                        "start": 33,
                        "end": 75,
                        "callee": {
                          "type": "Identifier",
                          "start": 33,
                          "end": 43,
                          "name": "mapGetters"
                        },
                        "arguments": [
                          {
                            "type": "ArrayExpression",
                            "start": 44,
                            "end": 74,
                            "elements": [
                              {
                                "type": "Literal",
                                "start": 46,
                                "end": 58,
                                "value": "some_val_1"
                              },
                              {
                                "type": "Literal",
                                "start": 60,
                                "end": 72,
                                "value": "some_val_2"
                              }
                            ]
                          }
                        ],
                        "optional": false
                      }
                    }
                  ]
                },
                "kind": "init"
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

func Test394(t *testing.T) {
	ast, err := compile("let {x, ...y} = v")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 17,
  "body": [
    {
      "type": "VariableDeclaration",
      "start": 0,
      "end": 17,
      "declarations": [
        {
          "type": "VariableDeclarator",
          "start": 4,
          "end": 17,
          "id": {
            "type": "ObjectPattern",
            "start": 4,
            "end": 13,
            "properties": [
              {
                "type": "Property",
                "start": 5,
                "end": 6,
                "method": false,
                "shorthand": true,
                "computed": false,
                "key": {
                  "type": "Identifier",
                  "start": 5,
                  "end": 6,
                  "name": "x"
                },
                "kind": "init",
                "value": {
                  "type": "Identifier",
                  "start": 5,
                  "end": 6,
                  "name": "x"
                }
              },
              {
                "type": "RestElement",
                "start": 8,
                "end": 12,
                "argument": {
                  "type": "Identifier",
                  "start": 11,
                  "end": 12,
                  "name": "y"
                }
              }
            ]
          },
          "init": {
            "type": "Identifier",
            "start": 16,
            "end": 17,
            "name": "v"
          }
        }
      ],
      "kind": "let"
    }
  ]
}
	`, ast)
}

func Test395(t *testing.T) {
	ast, err := compile("(function({x, ...y}) {})")
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
        "type": "FunctionExpression",
        "start": 1,
        "end": 23,
        "id": null,
        "expression": false,
        "generator": false,
        "async": false,
        "params": [
          {
            "type": "ObjectPattern",
            "start": 10,
            "end": 19,
            "properties": [
              {
                "type": "Property",
                "start": 11,
                "end": 12,
                "method": false,
                "shorthand": true,
                "computed": false,
                "key": {
                  "type": "Identifier",
                  "start": 11,
                  "end": 12,
                  "name": "x"
                },
                "kind": "init",
                "value": {
                  "type": "Identifier",
                  "start": 11,
                  "end": 12,
                  "name": "x"
                }
              },
              {
                "type": "RestElement",
                "start": 14,
                "end": 18,
                "argument": {
                  "type": "Identifier",
                  "start": 17,
                  "end": 18,
                  "name": "y"
                }
              }
            ]
          }
        ],
        "body": {
          "type": "BlockStatement",
          "start": 21,
          "end": 23,
          "body": []
        }
      }
    }
  ]}
	`, ast)
}

func Test396(t *testing.T) {
	ast, err := compile("const fn = ({text = \"default\", ...props}) => text + props.children")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 66,
  "body": [
    {
      "type": "VariableDeclaration",
      "start": 0,
      "end": 66,
      "declarations": [
        {
          "type": "VariableDeclarator",
          "start": 6,
          "end": 66,
          "id": {
            "type": "Identifier",
            "start": 6,
            "end": 8,
            "name": "fn"
          },
          "init": {
            "type": "ArrowFunctionExpression",
            "start": 11,
            "end": 66,
            "id": null,
            "expression": false,
            "generator": false,
            "async": false,
            "params": [
              {
                "type": "ObjectPattern",
                "start": 12,
                "end": 40,
                "properties": [
                  {
                    "type": "Property",
                    "start": 13,
                    "end": 29,
                    "method": false,
                    "shorthand": true,
                    "computed": false,
                    "key": {
                      "type": "Identifier",
                      "start": 13,
                      "end": 17,
                      "name": "text"
                    },
                    "kind": "init",
                    "value": {
                      "type": "AssignmentPattern",
                      "start": 13,
                      "end": 29,
                      "left": {
                        "type": "Identifier",
                        "start": 13,
                        "end": 17,
                        "name": "text"
                      },
                      "right": {
                        "type": "Literal",
                        "start": 20,
                        "end": 29,
                        "value": "default"
                      }
                    }
                  },
                  {
                    "type": "RestElement",
                    "start": 31,
                    "end": 39,
                    "argument": {
                      "type": "Identifier",
                      "start": 34,
                      "end": 39,
                      "name": "props"
                    }
                  }
                ]
              }
            ],
            "body": {
              "type": "BinaryExpression",
              "start": 45,
              "end": 66,
              "left": {
                "type": "Identifier",
                "start": 45,
                "end": 49,
                "name": "text"
              },
              "operator": "+",
              "right": {
                "type": "MemberExpression",
                "start": 52,
                "end": 66,
                "object": {
                  "type": "Identifier",
                  "start": 52,
                  "end": 57,
                  "name": "props"
                },
                "property": {
                  "type": "Identifier",
                  "start": 58,
                  "end": 66,
                  "name": "children"
                },
                "computed": false,
                "optional": false
              }
            }
          }
        }
      ],
      "kind": "const"
    }
  ]
}
	`, ast)
}

func Test397(t *testing.T) {
	ast, err := compile("`foo`")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
        "type": "TemplateLiteral",
        "start": 0,
        "end": 5,
        "expressions": [],
        "quasis": [
          {
            "type": "TemplateElement",
            "start": 1,
            "end": 4,
            "value": {
              "raw": "foo",
              "cooked": "foo"
            },
            "tail": true
          }
        ]
      }
    }
  ]
}
	`, ast)
}

func Test398(t *testing.T) {
	ast, err := compile("`foo\\u25a0`")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 11,
  "body": [
    {
      "type": "ExpressionStatement",
      "start": 0,
      "end": 11,
      "expression": {
        "type": "TemplateLiteral",
        "start": 0,
        "end": 11,
        "expressions": [],
        "quasis": [
          {
            "type": "TemplateElement",
            "start": 1,
            "end": 10,
            "value": {
              "raw": "foo\\u25a0",
              "cooked": "foo■"
            },
            "tail": true
          }
        ]
      }
    }
  ]
}
	`, ast)
}

func Test399(t *testing.T) {
	ast, err := compile("`foo${bar}\\u25a0`")
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
        "type": "TemplateLiteral",
        "start": 0,
        "end": 17,
        "expressions": [
          {
            "type": "Identifier",
            "start": 6,
            "end": 9,
            "name": "bar"
          }
        ],
        "quasis": [
          {
            "type": "TemplateElement",
            "start": 1,
            "end": 4,
            "value": {
              "raw": "foo",
              "cooked": "foo"
            },
            "tail": false
          },
          {
            "type": "TemplateElement",
            "start": 10,
            "end": 16,
            "value": {
              "raw": "\\u25a0",
              "cooked": "■"
            },
            "tail": true
          }
        ]
      }
    }
  ]
}
	`, ast)
}

func Test400(t *testing.T) {
	ast, err := compile("foo`\\u25a0`")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 11,
  "body": [
    {
      "type": "ExpressionStatement",
      "start": 0,
      "end": 11,
      "expression": {
        "type": "TaggedTemplateExpression",
        "start": 0,
        "end": 11,
        "tag": {
          "type": "Identifier",
          "start": 0,
          "end": 3,
          "name": "foo"
        },
        "quasi": {
          "type": "TemplateLiteral",
          "start": 3,
          "end": 11,
          "expressions": [],
          "quasis": [
            {
              "type": "TemplateElement",
              "start": 4,
              "end": 10,
              "value": {
                "raw": "\\u25a0",
                "cooked": "■"
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

func Test401(t *testing.T) {
	ast, err := compile("foo`foo${bar}\\u25a0`")
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
          "type": "Identifier",
          "start": 0,
          "end": 3,
          "name": "foo"
        },
        "quasi": {
          "type": "TemplateLiteral",
          "start": 3,
          "end": 20,
          "expressions": [
            {
              "type": "Identifier",
              "start": 9,
              "end": 12,
              "name": "bar"
            }
          ],
          "quasis": [
            {
              "type": "TemplateElement",
              "start": 4,
              "end": 7,
              "value": {
                "raw": "foo",
                "cooked": "foo"
              },
              "tail": false
            },
            {
              "type": "TemplateElement",
              "start": 13,
              "end": 19,
              "value": {
                "raw": "\\u25a0",
                "cooked": "■"
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

func Test402(t *testing.T) {
	ast, err := compile("foo`\\unicode`")
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
        "type": "TaggedTemplateExpression",
        "start": 0,
        "end": 13,
        "tag": {
          "type": "Identifier",
          "start": 0,
          "end": 3,
          "name": "foo"
        },
        "quasi": {
          "type": "TemplateLiteral",
          "start": 3,
          "end": 13,
          "expressions": [],
          "quasis": [
            {
              "type": "TemplateElement",
              "start": 4,
              "end": 12,
              "value": {
                "raw": "\\unicode",
                "cooked": ""
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

func Test403(t *testing.T) {
	ast, err := compile("foo`foo${bar}\\unicode`")
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
        "type": "TaggedTemplateExpression",
        "start": 0,
        "end": 22,
        "tag": {
          "type": "Identifier",
          "start": 0,
          "end": 3,
          "name": "foo"
        },
        "quasi": {
          "type": "TemplateLiteral",
          "start": 3,
          "end": 22,
          "expressions": [
            {
              "type": "Identifier",
              "start": 9,
              "end": 12,
              "name": "bar"
            }
          ],
          "quasis": [
            {
              "type": "TemplateElement",
              "start": 4,
              "end": 7,
              "value": {
                "raw": "foo",
                "cooked": "foo"
              },
              "tail": false
            },
            {
              "type": "TemplateElement",
              "start": 13,
              "end": 21,
              "value": {
                "raw": "\\unicode",
                "cooked": ""
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

func Test404(t *testing.T) {
	ast, err := compile("foo`\\u`")
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
        "type": "TaggedTemplateExpression",
        "start": 0,
        "end": 7,
        "tag": {
          "type": "Identifier",
          "start": 0,
          "end": 3,
          "name": "foo"
        },
        "quasi": {
          "type": "TemplateLiteral",
          "start": 3,
          "end": 7,
          "expressions": [],
          "quasis": [
            {
              "type": "TemplateElement",
              "start": 4,
              "end": 6,
              "value": {
                "raw": "\\u",
                "cooked": ""
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

func Test405(t *testing.T) {
	ast, err := compile("foo`\\u{`")
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
        "type": "TaggedTemplateExpression",
        "start": 0,
        "end": 8,
        "tag": {
          "type": "Identifier",
          "start": 0,
          "end": 3,
          "name": "foo"
        },
        "quasi": {
          "type": "TemplateLiteral",
          "start": 3,
          "end": 8,
          "expressions": [],
          "quasis": [
            {
              "type": "TemplateElement",
              "start": 4,
              "end": 7,
              "value": {
                "raw": "\\u{",
                "cooked": ""
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

func Test406(t *testing.T) {
	ast, err := compile("foo`\\u{abcdx`")
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
        "type": "TaggedTemplateExpression",
        "start": 0,
        "end": 13,
        "tag": {
          "type": "Identifier",
          "start": 0,
          "end": 3,
          "name": "foo"
        },
        "quasi": {
          "type": "TemplateLiteral",
          "start": 3,
          "end": 13,
          "expressions": [],
          "quasis": [
            {
              "type": "TemplateElement",
              "start": 4,
              "end": 12,
              "value": {
                "raw": "\\u{abcdx",
                "cooked": ""
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

func Test407(t *testing.T) {
	ast, err := compile("foo`\\u{abcdx}`")
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
        "type": "TaggedTemplateExpression",
        "start": 0,
        "end": 14,
        "tag": {
          "type": "Identifier",
          "start": 0,
          "end": 3,
          "name": "foo"
        },
        "quasi": {
          "type": "TemplateLiteral",
          "start": 3,
          "end": 14,
          "expressions": [],
          "quasis": [
            {
              "type": "TemplateElement",
              "start": 4,
              "end": 13,
              "value": {
                "raw": "\\u{abcdx}",
                "cooked": ""
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

func Test408(t *testing.T) {
	ast, err := compile("foo`\\unicode\\\\`")
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
        "type": "TaggedTemplateExpression",
        "start": 0,
        "end": 15,
        "tag": {
          "type": "Identifier",
          "start": 0,
          "end": 3,
          "name": "foo"
        },
        "quasi": {
          "type": "TemplateLiteral",
          "start": 3,
          "end": 15,
          "expressions": [],
          "quasis": [
            {
              "type": "TemplateElement",
              "start": 4,
              "end": 14,
              "value": {
                "raw": "\\unicode\\\\",
                "cooked": ""
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

func Test409(t *testing.T) {
	ast, err := compile("`${ {class: 1} }`")
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
        "type": "TemplateLiteral",
        "start": 0,
        "end": 17,
        "expressions": [
          {
            "type": "ObjectExpression",
            "start": 4,
            "end": 14,
            "properties": [
              {
                "type": "Property",
                "start": 5,
                "end": 13,
                "method": false,
                "shorthand": false,
                "computed": false,
                "key": {
                  "type": "Identifier",
                  "start": 5,
                  "end": 10,
                  "name": "class"
                },
                "value": {
                  "type": "Literal",
                  "start": 12,
                  "end": 13,
                  "value": 1
                },
                "kind": "init"
              }
            ]
          }
        ],
        "quasis": [
          {
            "type": "TemplateElement",
            "start": 1,
            "end": 1,
            "value": {
              "raw": "",
              "cooked": ""
            },
            "tail": false
          },
          {
            "type": "TemplateElement",
            "start": 16,
            "end": 16,
            "value": {
              "raw": "",
              "cooked": ""
            },
            "tail": true
          }
        ]
      }
    }
  ]
}
	`, ast)
}

func Test410(t *testing.T) {
	ast, err := compile("`${ {delete: 1} }`")
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
        "type": "TemplateLiteral",
        "start": 0,
        "end": 18,
        "expressions": [
          {
            "type": "ObjectExpression",
            "start": 4,
            "end": 15,
            "properties": [
              {
                "type": "Property",
                "start": 5,
                "end": 14,
                "method": false,
                "shorthand": false,
                "computed": false,
                "key": {
                  "type": "Identifier",
                  "start": 5,
                  "end": 11,
                  "name": "delete"
                },
                "value": {
                  "type": "Literal",
                  "start": 13,
                  "end": 14,
                  "value": 1
                },
                "kind": "init"
              }
            ]
          }
        ],
        "quasis": [
          {
            "type": "TemplateElement",
            "start": 1,
            "end": 1,
            "value": {
              "raw": "",
              "cooked": ""
            },
            "tail": false
          },
          {
            "type": "TemplateElement",
            "start": 17,
            "end": 17,
            "value": {
              "raw": "",
              "cooked": ""
            },
            "tail": true
          }
        ]
      }
    }
  ]
}
	`, ast)
}

func Test411(t *testing.T) {
	ast, err := compile("`${ {enum: 1} }`")
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
        "type": "TemplateLiteral",
        "start": 0,
        "end": 16,
        "expressions": [
          {
            "type": "ObjectExpression",
            "start": 4,
            "end": 13,
            "properties": [
              {
                "type": "Property",
                "start": 5,
                "end": 12,
                "method": false,
                "shorthand": false,
                "computed": false,
                "key": {
                  "type": "Identifier",
                  "start": 5,
                  "end": 9,
                  "name": "enum"
                },
                "value": {
                  "type": "Literal",
                  "start": 11,
                  "end": 12,
                  "value": 1
                },
                "kind": "init"
              }
            ]
          }
        ],
        "quasis": [
          {
            "type": "TemplateElement",
            "start": 1,
            "end": 1,
            "value": {
              "raw": "",
              "cooked": ""
            },
            "tail": false
          },
          {
            "type": "TemplateElement",
            "start": 15,
            "end": 15,
            "value": {
              "raw": "",
              "cooked": ""
            },
            "tail": true
          }
        ]
      }
    }
  ]
}
	`, ast)
}

func Test412(t *testing.T) {
	ast, err := compile("`${ {function: 1} }`")
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
        "type": "TemplateLiteral",
        "start": 0,
        "end": 20,
        "expressions": [
          {
            "type": "ObjectExpression",
            "start": 4,
            "end": 17,
            "properties": [
              {
                "type": "Property",
                "start": 5,
                "end": 16,
                "method": false,
                "shorthand": false,
                "computed": false,
                "key": {
                  "type": "Identifier",
                  "start": 5,
                  "end": 13,
                  "name": "function"
                },
                "value": {
                  "type": "Literal",
                  "start": 15,
                  "end": 16,
                  "value": 1
                },
                "kind": "init"
              }
            ]
          }
        ],
        "quasis": [
          {
            "type": "TemplateElement",
            "start": 1,
            "end": 1,
            "value": {
              "raw": "",
              "cooked": ""
            },
            "tail": false
          },
          {
            "type": "TemplateElement",
            "start": 19,
            "end": 19,
            "value": {
              "raw": "",
              "cooked": ""
            },
            "tail": true
          }
        ]
      }
    }
  ]
}
	`, ast)
}

func Test413(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test414(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test415(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test416(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test417(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test418(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test419(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test420(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test421(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test422(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test423(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test424(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test425(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test426(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test427(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test428(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test429(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test430(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test431(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test432(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test433(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test434(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test435(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test436(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test437(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test438(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test439(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test440(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test441(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test442(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test443(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test444(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test445(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test446(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test447(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test448(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test449(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test450(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test451(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test452(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test453(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test454(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test455(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test456(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test457(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test458(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test459(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test460(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test461(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test462(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test463(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test464(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test465(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test466(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test467(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test468(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test469(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test470(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test471(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test472(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test473(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test474(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test475(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test476(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test477(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test478(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test479(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test480(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test481(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test482(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test483(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test484(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test485(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test486(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test487(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test488(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test489(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test490(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test491(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test492(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test493(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test494(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test495(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test496(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test497(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test498(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test499(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test500(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test501(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test502(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test503(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test504(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test505(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test506(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test507(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test508(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test509(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test510(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test511(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test512(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test513(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test514(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test515(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test516(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test517(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test518(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test519(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test520(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test521(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test522(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test523(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test524(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test525(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test526(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func Test527(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}
