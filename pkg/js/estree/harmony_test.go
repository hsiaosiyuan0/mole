package estree

import (
	"testing"

	"github.com/hsiaosiyuan0/mole/pkg/assert"
)

// below tests follow the copyright declaration in the head of file:
// https://github.com/acornjs/acorn/blob/f85a712661fe2b92dbd73813d0cae37dc920fe6d/test/tests-harmony.js

// ES6 Unicode Code Point Escape Sequence
func TestHarmony1(t *testing.T) {
	ast, err := compile("\"\\u{714E}\\u{8336}\"")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "Literal",
        "value": "煎茶",
        "raw": "\"\\u{714E}\\u{8336}\"",
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

func TestHarmony2(t *testing.T) {
	ast, err := compile("\"\\u{20BB7}\\u{91CE}\\u{5BB6}\"")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "Literal",
        "value": "𠮷野家",
        "raw": "\"\\u{20BB7}\\u{91CE}\\u{5BB6}\"",
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

// ES6: Numeric Literal

func TestHarmony3(t *testing.T) {
	ast, err := compile("00")
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
        "raw": "00",
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

func TestHarmony4(t *testing.T) {
	ast, err := compile("0o0")
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
        "raw": "0o0",
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

func TestHarmony5(t *testing.T) {
	ast, err := compile("function test() {'use strict'; 0o0; }")
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
      "params": [],
      "body": {
        "type": "BlockStatement",
        "body": [
          {
            "type": "ExpressionStatement",
            "expression": {
              "type": "Literal",
              "value": "use strict",
              "raw": "'use strict'",
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
                "column": 17
              },
              "end": {
                "line": 1,
                "column": 30
              }
            }
          },
          {
            "type": "ExpressionStatement",
            "expression": {
              "type": "Literal",
              "value": 0,
              "raw": "0o0",
              "loc": {
                "start": {
                  "line": 1,
                  "column": 31
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
                "column": 31
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
            "column": 16
          },
          "end": {
            "line": 1,
            "column": 37
          }
        }
      },
      "generator": false,
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

func TestHarmony6(t *testing.T) {
	ast, err := compile("0o2")
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
        "raw": "0o2",
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

func TestHarmony7(t *testing.T) {
	ast, err := compile("0o12")
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
        "raw": "0o12",
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

func TestHarmony8(t *testing.T) {
	ast, err := compile("0O0")
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
        "raw": "0O0",
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

func TestHarmony9(t *testing.T) {
	ast, err := compile("function test() {'use strict'; 0O0; }")
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
      "params": [],
      "body": {
        "type": "BlockStatement",
        "body": [
          {
            "type": "ExpressionStatement",
            "expression": {
              "type": "Literal",
              "value": "use strict",
              "raw": "'use strict'",
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
                "column": 17
              },
              "end": {
                "line": 1,
                "column": 30
              }
            }
          },
          {
            "type": "ExpressionStatement",
            "expression": {
              "type": "Literal",
              "value": 0,
              "raw": "0O0",
              "loc": {
                "start": {
                  "line": 1,
                  "column": 31
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
                "column": 31
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
            "column": 16
          },
          "end": {
            "line": 1,
            "column": 37
          }
        }
      },
      "generator": false,
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

func TestHarmony10(t *testing.T) {
	ast, err := compile("0O2")
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
        "raw": "0O2",
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

func TestHarmony11(t *testing.T) {
	ast, err := compile("0O12")
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
        "raw": "0O12",
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

func TestHarmony12(t *testing.T) {
	ast, err := compile("0b0")
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
        "raw": "0b0",
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

func TestHarmony13(t *testing.T) {
	ast, err := compile("0b1")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "Literal",
        "value": 1,
        "raw": "0b1",
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

func TestHarmony14(t *testing.T) {
	ast, err := compile("0b10")
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
        "raw": "0b10",
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

func TestHarmony15(t *testing.T) {
	ast, err := compile("0B0")
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
        "raw": "0B0",
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

func TestHarmony16(t *testing.T) {
	ast, err := compile("0B1")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "Literal",
        "value": 1,
        "raw": "0B1",
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

func TestHarmony17(t *testing.T) {
	ast, err := compile("0B10")
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
        "raw": "0B10",
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

// ES6 Template Strings

func TestHarmony18(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony19(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony20(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony21(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony22(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony23(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony24(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony25(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony26(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony27(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony28(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony29(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony30(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony31(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony32(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony33(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony34(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony35(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony36(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony37(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony38(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony39(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony40(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony41(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony42(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony43(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony44(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony45(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony46(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony47(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony48(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony49(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony50(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony51(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony52(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony53(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony54(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony55(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony56(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony57(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony58(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony59(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony60(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony61(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony62(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony63(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony64(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony65(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony66(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony67(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony68(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony69(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony70(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony71(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony72(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony73(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony74(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony75(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony76(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony77(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony78(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony79(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony80(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony81(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony82(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony83(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony84(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony85(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony86(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony87(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony88(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony89(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony90(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony91(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony92(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony93(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony94(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony95(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony96(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony97(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony98(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony99(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony100(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony101(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony102(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony103(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony104(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony105(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony106(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony107(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony108(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony109(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony110(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony111(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony112(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony113(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony114(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony115(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony116(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony117(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony118(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony119(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony120(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony121(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony122(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony123(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony124(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony125(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony126(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony127(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony128(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony129(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony130(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony131(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony132(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony133(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony134(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony135(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony136(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony137(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony138(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony139(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony140(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony141(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony142(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony143(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony144(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony145(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony146(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony147(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony148(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony149(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony150(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony151(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony152(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony153(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony154(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony155(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony156(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony157(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony158(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony159(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony160(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony161(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony162(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony163(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony164(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony165(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony166(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony167(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony168(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony169(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony170(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony171(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony172(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony173(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony174(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony175(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony176(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony177(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony178(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony179(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony180(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony181(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony182(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony183(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony184(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony185(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony186(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony187(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony188(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony189(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony190(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony191(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony192(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony193(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony194(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony195(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony196(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony197(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony198(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony199(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony200(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony201(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony202(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony203(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony204(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony205(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony206(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony207(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony208(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony209(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony210(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony211(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony212(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony213(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony214(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony215(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony216(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony217(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony218(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony219(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony220(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony221(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony222(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony223(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony224(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony225(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony226(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony227(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony228(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony229(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony230(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony231(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony232(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony233(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony234(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony235(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony236(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony237(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony238(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony239(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony240(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony241(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony242(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony243(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony244(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony245(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony246(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony247(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony248(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony249(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony250(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony251(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony252(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony253(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony254(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony255(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony256(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony257(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony258(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony259(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony260(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony261(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony262(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony263(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony264(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony265(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony266(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony267(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony268(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony269(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony270(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony271(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony272(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony273(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony274(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony275(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony276(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony277(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony278(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony279(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony280(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony281(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony282(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony283(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony284(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony285(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony286(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony287(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony288(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony289(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony290(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony291(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony292(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony293(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony294(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony295(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony296(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony297(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony298(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony299(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony300(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}
