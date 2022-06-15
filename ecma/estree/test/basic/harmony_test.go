package estree_test

import (
	"testing"

	. "github.com/hsiaosiyuan0/mole/ecma/estree/test"
	"github.com/hsiaosiyuan0/mole/ecma/parser"
	. "github.com/hsiaosiyuan0/mole/util"
)

// below tests follow the copyright declaration in the head of file:
// https://github.com/acornjs/acorn/blob/f85a712661fe2b92dbd73813d0cae37dc920fe6d/test/tests-harmony.js

// ES6 Unicode Code Point Escape Sequence
func TestHarmony1(t *testing.T) {
	ast, err := Compile("\"\\u{714E}\\u{8336}\"")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
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
	ast, err := Compile("\"\\u{20BB7}\\u{91CE}\\u{5BB6}\"")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
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
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_STRICT)
	ast, err := CompileWithOpts("00", opts)
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
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
	ast, err := Compile("0o0")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
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
	ast, err := Compile("function test() {'use strict'; 0o0; }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
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
	ast, err := Compile("0o2")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
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
	ast, err := Compile("0o12")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
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
	ast, err := Compile("0O0")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
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
	ast, err := Compile("function test() {'use strict'; 0O0; }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
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
	ast, err := Compile("0O2")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
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
	ast, err := Compile("0O12")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
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
	ast, err := Compile("0b0")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
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
	ast, err := Compile("0b1")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
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
	ast, err := Compile("0b10")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
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
	ast, err := Compile("0B0")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
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
	ast, err := Compile("0B1")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
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
	ast, err := Compile("0B10")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
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
	ast, err := Compile("`42`")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "TemplateLiteral",
        "quasis": [
          {
            "type": "TemplateElement",
            "value": {
              "raw": "42",
              "cooked": "42"
            },
            "tail": true,
            "loc": {
              "start": {
                "line": 1,
                "column": 1
              },
              "end": {
                "line": 1,
                "column": 3
              }
            }
          }
        ],
        "expressions": [],
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

func TestHarmony19(t *testing.T) {
	ast, err := Compile("raw`42`")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "TaggedTemplateExpression",
        "tag": {
          "type": "Identifier",
          "name": "raw",
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
        "quasi": {
          "type": "TemplateLiteral",
          "quasis": [
            {
              "type": "TemplateElement",
              "value": {
                "raw": "42",
                "cooked": "42"
              },
              "tail": true,
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
          "expressions": [],
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

func TestHarmony20(t *testing.T) {
	ast, err := Compile("raw`hello ${name}`")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "TaggedTemplateExpression",
        "tag": {
          "type": "Identifier",
          "name": "raw",
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
        "quasi": {
          "type": "TemplateLiteral",
          "quasis": [
            {
              "type": "TemplateElement",
              "value": {
                "raw": "hello ",
                "cooked": "hello "
              },
              "tail": false,
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
              "type": "TemplateElement",
              "value": {
                "raw": "",
                "cooked": ""
              },
              "tail": true,
              "loc": {
                "start": {
                  "line": 1,
                  "column": 17
                },
                "end": {
                  "line": 1,
                  "column": 17
                }
              }
            }
          ],
          "expressions": [
            {
              "type": "Identifier",
              "name": "name",
              "loc": {
                "start": {
                  "line": 1,
                  "column": 12
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
              "column": 3
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

func TestHarmony21(t *testing.T) {
	ast, err := Compile("`$`")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "TemplateLiteral",
        "quasis": [
          {
            "type": "TemplateElement",
            "value": {
              "raw": "$",
              "cooked": "$"
            },
            "tail": true,
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
          }
        ],
        "expressions": [],
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

func TestHarmony22(t *testing.T) {
	ast, err := Compile("`\\n\\r\\b\\v\\t\\f\\\n\\\r\n\\\u2028\\\u2029`")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "TemplateLiteral",
        "quasis": [
          {
            "type": "TemplateElement",
            "value": {
              "raw": "\\n\\r\\b\\v\\t\\f\\\n\\\r\n\\\u2028\\\u2029",
              "cooked": "\n\r\b\u000b\t\f"
            },
            "tail": true,
            "loc": {
              "start": {
                "line": 1,
                "column": 1
              },
              "end": {
                "line": 3,
                "column": 4
              }
            }
          }
        ],
        "expressions": [],
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 3,
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
          "line": 3,
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
      "line": 3,
      "column": 5
    }
  }
}
	`, ast)
}

func TestHarmony23(t *testing.T) {
	ast, err := Compile("`\n\r\n\r`")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "TemplateLiteral",
        "quasis": [
          {
            "type": "TemplateElement",
            "value": {
              "raw": "\n\r\n\r",
              "cooked": "\n\n\n"
            },
            "tail": true,
            "loc": {
              "start": {
                "line": 1,
                "column": 1
              },
              "end": {
                "line": 4,
                "column": 0
              }
            }
          }
        ],
        "expressions": [],
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 4,
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
          "line": 4,
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
      "line": 4,
      "column": 1
    }
  }
}
	`, ast)
}

func TestHarmony24(t *testing.T) {
	ast, err := Compile("`\\u{000042}\\u0042\\x42u0\\A`")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "TemplateLiteral",
        "quasis": [
          {
            "type": "TemplateElement",
            "value": {
              "raw": "\\u{000042}\\u0042\\x42u0\\A",
              "cooked": "BBBu0A"
            },
            "tail": true,
            "loc": {
              "start": {
                "line": 1,
                "column": 1
              },
              "end": {
                "line": 1,
                "column": 25
              }
            }
          }
        ],
        "expressions": [],
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

func TestHarmony25(t *testing.T) {
	ast, err := Compile("new raw`42`")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "NewExpression",
        "callee": {
          "type": "TaggedTemplateExpression",
          "tag": {
            "type": "Identifier",
            "name": "raw",
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
          "quasi": {
            "type": "TemplateLiteral",
            "quasis": [
              {
                "type": "TemplateElement",
                "value": {
                  "raw": "42",
                  "cooked": "42"
                },
                "tail": true,
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
              }
            ],
            "expressions": [],
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

func TestHarmony26(t *testing.T) {
	ast, err := Compile("`outer${{x: {y: 10}}}bar${`nested${function(){return 1;}}endnest`}end`")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "TemplateLiteral",
        "expressions": [
          {
            "type": "ObjectExpression",
            "properties": [
              {
                "type": "Property",
                "method": false,
                "shorthand": false,
                "computed": false,
                "key": {
                  "type": "Identifier",
                  "name": "x"
                },
                "value": {
                  "type": "ObjectExpression",
                  "properties": [
                    {
                      "type": "Property",
                      "method": false,
                      "shorthand": false,
                      "computed": false,
                      "key": {
                        "type": "Identifier",
                        "name": "y"
                      },
                      "value": {
                        "type": "Literal",
                        "value": 10,
                        "raw": "10"
                      },
                      "kind": "init"
                    }
                  ]
                },
                "kind": "init"
              }
            ]
          },
          {
            "type": "TemplateLiteral",
            "expressions": [
              {
                "type": "FunctionExpression",
                "id": null,
                "params": [],
                "generator": false,
                "body": {
                  "type": "BlockStatement",
                  "body": [
                    {
                      "type": "ReturnStatement",
                      "argument": {
                        "type": "Literal",
                        "value": 1,
                        "raw": "1"
                      }
                    }
                  ]
                },
                "expression": false
              }
            ],
            "quasis": [
              {
                "type": "TemplateElement",
                "value": {
                  "cooked": "nested",
                  "raw": "nested"
                },
                "tail": false
              },
              {
                "type": "TemplateElement",
                "value": {
                  "cooked": "endnest",
                  "raw": "endnest"
                },
                "tail": true
              }
            ]
          }
        ],
        "quasis": [
          {
            "type": "TemplateElement",
            "value": {
              "cooked": "outer",
              "raw": "outer"
            },
            "tail": false
          },
          {
            "type": "TemplateElement",
            "value": {
              "cooked": "bar",
              "raw": "bar"
            },
            "tail": false
          },
          {
            "type": "TemplateElement",
            "value": {
              "cooked": "end",
              "raw": "end"
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

// ES6: Switch Case Declaration

func TestHarmony27(t *testing.T) {
	ast, err := Compile("switch (answer) { case 42: let t = 42; break; }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
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
          "test": {
            "type": "Literal",
            "value": 42,
            "raw": "42",
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
          "consequent": [
            {
              "type": "VariableDeclaration",
              "declarations": [
                {
                  "type": "VariableDeclarator",
                  "id": {
                    "type": "Identifier",
                    "name": "t",
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
                  "init": {
                    "type": "Literal",
                    "value": 42,
                    "raw": "42",
                    "loc": {
                      "start": {
                        "line": 1,
                        "column": 35
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
                      "column": 31
                    },
                    "end": {
                      "line": 1,
                      "column": 37
                    }
                  }
                }
              ],
              "kind": "let",
              "loc": {
                "start": {
                  "line": 1,
                  "column": 27
                },
                "end": {
                  "line": 1,
                  "column": 38
                }
              }
            },
            {
              "type": "BreakStatement",
              "label": null,
              "loc": {
                "start": {
                  "line": 1,
                  "column": 39
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
              "column": 18
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
          "column": 47
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 47
    }
  }
}
	`, ast)
}

// ES6: Arrow Function

func TestHarmony28(t *testing.T) {
	ast, err := Compile("() => \"test\"")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "ArrowFunctionExpression",
        "id": null,
        "params": [],
        "body": {
          "type": "Literal",
          "value": "test",
          "raw": "\"test\"",
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
        "generator": false,
        "expression": true,
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

func TestHarmony29(t *testing.T) {
	ast, err := Compile("e => \"test\"")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "ArrowFunctionExpression",
        "id": null,
        "params": [
          {
            "type": "Identifier",
            "name": "e",
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
        "body": {
          "type": "Literal",
          "value": "test",
          "raw": "\"test\"",
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
        "generator": false,
        "expression": true,
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

func TestHarmony30(t *testing.T) {
	ast, err := Compile("(e) => \"test\"")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "ArrowFunctionExpression",
        "id": null,
        "params": [
          {
            "type": "Identifier",
            "name": "e",
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
          }
        ],
        "body": {
          "type": "Literal",
          "value": "test",
          "raw": "\"test\"",
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
        "generator": false,
        "expression": true,
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

func TestHarmony31(t *testing.T) {
	ast, err := Compile("(a, b) => \"test\"")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "ArrowFunctionExpression",
        "id": null,
        "params": [
          {
            "type": "Identifier",
            "name": "a",
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
          {
            "type": "Identifier",
            "name": "b",
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
        "body": {
          "type": "Literal",
          "value": "test",
          "raw": "\"test\"",
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
        "generator": false,
        "expression": true,
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

func TestHarmony32(t *testing.T) {
	ast, err := Compile("e => { 42; }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "ArrowFunctionExpression",
        "id": null,
        "params": [
          {
            "type": "Identifier",
            "name": "e",
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
        "body": {
          "type": "BlockStatement",
          "body": [
            {
              "type": "ExpressionStatement",
              "expression": {
                "type": "Literal",
                "value": 42,
                "raw": "42",
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
                  "column": 7
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
              "column": 5
            },
            "end": {
              "line": 1,
              "column": 12
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

func TestHarmony33(t *testing.T) {
	ast, err := Compile("e => ({ property: 42 })")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "ArrowFunctionExpression",
        "id": null,
        "params": [
          {
            "type": "Identifier",
            "name": "e",
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
        "body": {
          "type": "ObjectExpression",
          "properties": [
            {
              "type": "Property",
              "key": {
                "type": "Identifier",
                "name": "property",
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
              "value": {
                "type": "Literal",
                "value": 42,
                "raw": "42",
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
              "kind": "init",
              "method": false,
              "shorthand": false,
              "computed": false,
              "loc": {
                "start": {
                  "line": 1,
                  "column": 8
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
              "column": 6
            },
            "end": {
              "line": 1,
              "column": 22
            }
          }
        },
        "generator": false,
        "expression": true,
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

func TestHarmony34(t *testing.T) {
	ast, err := Compile("e => { label: 42 }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "ArrowFunctionExpression",
        "id": null,
        "params": [
          {
            "type": "Identifier",
            "name": "e",
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
        "body": {
          "type": "BlockStatement",
          "body": [
            {
              "type": "LabeledStatement",
              "label": {
                "type": "Identifier",
                "name": "label",
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 7
                  },
                  "end": {
                    "line": 1,
                    "column": 12
                  }
                }
              },
              "body": {
                "type": "ExpressionStatement",
                "expression": {
                  "type": "Literal",
                  "value": 42,
                  "raw": "42",
                  "loc": {
                    "start": {
                      "line": 1,
                      "column": 14
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
                    "column": 14
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
                  "column": 7
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
              "column": 5
            },
            "end": {
              "line": 1,
              "column": 18
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

func TestHarmony35(t *testing.T) {
	ast, err := Compile("(a, b) => { 42; }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "ArrowFunctionExpression",
        "id": null,
        "params": [
          {
            "type": "Identifier",
            "name": "a",
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
          {
            "type": "Identifier",
            "name": "b",
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
        "body": {
          "type": "BlockStatement",
          "body": [
            {
              "type": "ExpressionStatement",
              "expression": {
                "type": "Literal",
                "value": 42,
                "raw": "42",
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
                  "column": 12
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
              "column": 10
            },
            "end": {
              "line": 1,
              "column": 17
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

func TestHarmony36(t *testing.T) {
	ast, err := Compile("([a, , b]) => 42")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "ArrowFunctionExpression",
        "id": null,
        "params": [
          {
            "type": "ArrayPattern",
            "elements": [
              {
                "type": "Identifier",
                "name": "a",
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
              null,
              {
                "type": "Identifier",
                "name": "b",
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
            "loc": {
              "start": {
                "line": 1,
                "column": 1
              },
              "end": {
                "line": 1,
                "column": 9
              }
            }
          }
        ],
        "body": {
          "type": "Literal",
          "value": 42,
          "raw": "42",
          "loc": {
            "start": {
              "line": 1,
              "column": 14
            },
            "end": {
              "line": 1,
              "column": 16
            }
          }
        },
        "generator": false,
        "expression": true,
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

func TestHarmony37(t *testing.T) {
	ast, err := Compile("(() => {})()")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "CallExpression",
        "start": 0,
        "end": 12,
        "callee": {
          "type": "ArrowFunctionExpression",
          "id": null,
          "params": [],
          "body": {
            "type": "BlockStatement",
            "body": [],
            "start": 7,
            "end": 9,
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
          "generator": false,
          "loc": {
            "start": {
              "line": 1,
              "column": 1
            },
            "end": {
              "line": 1,
              "column": 9
            }
          }
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

func TestHarmony38(t *testing.T) {
	ast, err := Compile("((() => {}))()")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "CallExpression",
        "start": 0,
        "end": 14,
        "callee": {
          "type": "ArrowFunctionExpression",
          "id": null,
          "params": [],
          "body": {
            "type": "BlockStatement",
            "body": [],
            "start": 8,
            "end": 10,
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
          "generator": false,
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

func TestHarmony39(t *testing.T) {
	ast, err := Compile("(x=1) => x * x")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "ArrowFunctionExpression",
        "id": null,
        "params": [
          {
            "type": "AssignmentPattern",
            "left": {
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
            "right": {
              "type": "Literal",
              "value": 1,
              "raw": "1",
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
            },
            "loc": {
              "start": {
                "line": 1,
                "column": 1
              },
              "end": {
                "line": 1,
                "column": 4
              }
            }
          }
        ],
        "body": {
          "type": "BinaryExpression",
          "operator": "*",
          "left": {
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
          "right": {
            "type": "Identifier",
            "name": "x",
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
              "column": 9
            },
            "end": {
              "line": 1,
              "column": 14
            }
          }
        },
        "generator": false,
        "expression": true,
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

func TestHarmony40(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_STRICT)
	ast, err := CompileWithOpts("eval => 42", opts)
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "ArrowFunctionExpression",
        "id": null,
        "params": [
          {
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
          }
        ],
        "body": {
          "type": "Literal",
          "value": 42,
          "raw": "42",
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
        "generator": false,
        "expression": true,
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

func TestHarmony41(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_STRICT)
	ast, err := CompileWithOpts("(a) => 00", opts)
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "ArrowFunctionExpression",
        "id": null,
        "params": [
          {
            "type": "Identifier",
            "name": "a",
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
          }
        ],
        "body": {
          "type": "Literal",
          "value": 0,
          "raw": "00",
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
        "generator": false,
        "expression": true,
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

func TestHarmony42(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_STRICT)
	ast, err := CompileWithOpts("(eval, a) => 42", opts)
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "ArrowFunctionExpression",
        "id": null,
        "params": [
          {
            "type": "Identifier",
            "name": "eval",
            "loc": {
              "start": {
                "line": 1,
                "column": 1
              },
              "end": {
                "line": 1,
                "column": 5
              }
            }
          },
          {
            "type": "Identifier",
            "name": "a",
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
        "body": {
          "type": "Literal",
          "value": 42,
          "raw": "42",
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
        "generator": false,
        "expression": true,
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

func TestHarmony43(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_STRICT)
	ast, err := CompileWithOpts("(eval = 10) => 42", opts)
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "ArrowFunctionExpression",
        "id": null,
        "params": [
          {
            "type": "AssignmentPattern",
            "left": {
              "type": "Identifier",
              "name": "eval",
              "loc": {
                "start": {
                  "line": 1,
                  "column": 1
                },
                "end": {
                  "line": 1,
                  "column": 5
                }
              }
            },
            "right": {
              "type": "Literal",
              "value": 10,
              "raw": "10",
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
                "column": 1
              },
              "end": {
                "line": 1,
                "column": 10
              }
            }
          }
        ],
        "body": {
          "type": "Literal",
          "value": 42,
          "raw": "42",
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
        "generator": false,
        "expression": true,
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

func TestHarmony44(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_STRICT)
	ast, err := CompileWithOpts("(eval, a = 10) => 42", opts)
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "ArrowFunctionExpression",
        "id": null,
        "params": [
          {
            "type": "Identifier",
            "name": "eval",
            "loc": {
              "start": {
                "line": 1,
                "column": 1
              },
              "end": {
                "line": 1,
                "column": 5
              }
            }
          },
          {
            "type": "AssignmentPattern",
            "left": {
              "type": "Identifier",
              "name": "a",
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
            "right": {
              "type": "Literal",
              "value": 10,
              "raw": "10",
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
          }
        ],
        "body": {
          "type": "Literal",
          "value": 42,
          "raw": "42",
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
        "generator": false,
        "expression": true,
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

func TestHarmony45(t *testing.T) {
	ast, err := Compile("(x => x)")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "ArrowFunctionExpression",
        "id": null,
        "params": [
          {
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
          }
        ],
        "body": {
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
        "generator": false,
        "expression": true,
        "loc": {
          "start": {
            "line": 1,
            "column": 1
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

func TestHarmony46(t *testing.T) {
	ast, err := Compile("x => y => 42")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "ArrowFunctionExpression",
        "id": null,
        "params": [
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
          }
        ],
        "body": {
          "type": "ArrowFunctionExpression",
          "id": null,
          "params": [
            {
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
            }
          ],
          "body": {
            "type": "Literal",
            "value": 42,
            "raw": "42",
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
          "generator": false,
          "expression": true,
          "loc": {
            "start": {
              "line": 1,
              "column": 5
            },
            "end": {
              "line": 1,
              "column": 12
            }
          }
        },
        "generator": false,
        "expression": true,
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

func TestHarmony47(t *testing.T) {
	ast, err := Compile("(x) => ((y, z) => (x, y, z))")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "ArrowFunctionExpression",
        "id": null,
        "params": [
          {
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
          }
        ],
        "body": {
          "type": "ArrowFunctionExpression",
          "id": null,
          "params": [
            {
              "type": "Identifier",
              "name": "y",
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
            {
              "type": "Identifier",
              "name": "z",
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
            }
          ],
          "body": {
            "type": "SequenceExpression",
            "expressions": [
              {
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
              {
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
              {
                "type": "Identifier",
                "name": "z",
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
              }
            ],
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
          "generator": false,
          "expression": true,
          "loc": {
            "start": {
              "line": 1,
              "column": 8
            },
            "end": {
              "line": 1,
              "column": 27
            }
          }
        },
        "generator": false,
        "expression": true,
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

func TestHarmony48(t *testing.T) {
	ast, err := Compile("foo(() => {})")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
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
            "type": "ArrowFunctionExpression",
            "id": null,
            "params": [],
            "body": {
              "type": "BlockStatement",
              "body": [],
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
            "generator": false,
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

func TestHarmony49(t *testing.T) {
	ast, err := Compile("foo((x, y) => {})")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
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
            "type": "ArrowFunctionExpression",
            "id": null,
            "params": [
              {
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
              {
                "type": "Identifier",
                "name": "y",
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
              }
            ],
            "body": {
              "type": "BlockStatement",
              "body": [],
              "loc": {
                "start": {
                  "line": 1,
                  "column": 14
                },
                "end": {
                  "line": 1,
                  "column": 16
                }
              }
            },
            "generator": false,
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

// ES6: Method Definition

func TestHarmony50(t *testing.T) {
	ast, err := Compile("x = { method() { } }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
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
                "name": "method",
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
                      "column": 18
                    }
                  }
                },
                "generator": false,
                "expression": false,
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
              "kind": "init",
              "method": true,
              "shorthand": false,
              "computed": false,
              "loc": {
                "start": {
                  "line": 1,
                  "column": 6
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

func TestHarmony51(t *testing.T) {
	ast, err := Compile("x = { method(test) { } }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
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
                "name": "method",
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
              "value": {
                "type": "FunctionExpression",
                "id": null,
                "params": [
                  {
                    "type": "Identifier",
                    "name": "test",
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
                  }
                ],
                "body": {
                  "type": "BlockStatement",
                  "body": [],
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
                "generator": false,
                "expression": false,
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 12
                  },
                  "end": {
                    "line": 1,
                    "column": 22
                  }
                }
              },
              "kind": "init",
              "method": true,
              "shorthand": false,
              "computed": false,
              "loc": {
                "start": {
                  "line": 1,
                  "column": 6
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

func TestHarmony52(t *testing.T) {
	ast, err := Compile("x = { 'method'() { } }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
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
                "value": "method",
                "raw": "'method'",
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
                      "column": 20
                    }
                  }
                },
                "generator": false,
                "expression": false,
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 14
                  },
                  "end": {
                    "line": 1,
                    "column": 20
                  }
                }
              },
              "kind": "init",
              "method": true,
              "shorthand": false,
              "computed": false,
              "loc": {
                "start": {
                  "line": 1,
                  "column": 6
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

func TestHarmony53(t *testing.T) {
	ast, err := Compile("x = { get() { } }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
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
                "type": "FunctionExpression",
                "id": null,
                "params": [],
                "body": {
                  "type": "BlockStatement",
                  "body": [],
                  "loc": {
                    "start": {
                      "line": 1,
                      "column": 12
                    },
                    "end": {
                      "line": 1,
                      "column": 15
                    }
                  }
                },
                "generator": false,
                "expression": false,
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
              },
              "kind": "init",
              "method": true,
              "shorthand": false,
              "computed": false,
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

func TestHarmony54(t *testing.T) {
	ast, err := Compile("x = { set() { } }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
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
                "type": "FunctionExpression",
                "id": null,
                "params": [],
                "body": {
                  "type": "BlockStatement",
                  "body": [],
                  "loc": {
                    "start": {
                      "line": 1,
                      "column": 12
                    },
                    "end": {
                      "line": 1,
                      "column": 15
                    }
                  }
                },
                "generator": false,
                "expression": false,
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
              },
              "kind": "init",
              "method": true,
              "shorthand": false,
              "computed": false,
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

func TestHarmony55(t *testing.T) {
	ast, err := Compile("x = { method() { super.a(); } }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 31
    }
  },
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "AssignmentExpression",
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 31
          }
        },
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
          "loc": {
            "start": {
              "line": 1,
              "column": 4
            },
            "end": {
              "line": 1,
              "column": 31
            }
          },
          "properties": [
            {
              "type": "Property",
              "loc": {
                "start": {
                  "line": 1,
                  "column": 6
                },
                "end": {
                  "line": 1,
                  "column": 29
                }
              },
              "method": true,
              "shorthand": false,
              "computed": false,
              "key": {
                "type": "Identifier",
                "name": "method"
              },
              "kind": "init",
              "value": {
                "type": "FunctionExpression",
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 12
                  },
                  "end": {
                    "line": 1,
                    "column": 29
                  }
                },
                "id": null,
                "expression": false,
                "generator": false,
                "params": [],
                "body": {
                  "type": "BlockStatement",
                  "loc": {
                    "start": {
                      "line": 1,
                      "column": 15
                    },
                    "end": {
                      "line": 1,
                      "column": 29
                    }
                  },
                  "body": [
                    {
                      "type": "ExpressionStatement",
                      "loc": {
                        "start": {
                          "line": 1,
                          "column": 17
                        },
                        "end": {
                          "line": 1,
                          "column": 27
                        }
                      },
                      "expression": {
                        "type": "CallExpression",
                        "loc": {
                          "start": {
                            "line": 1,
                            "column": 17
                          },
                          "end": {
                            "line": 1,
                            "column": 26
                          }
                        },
                        "callee": {
                          "type": "MemberExpression",
                          "loc": {
                            "start": {
                              "line": 1,
                              "column": 17
                            },
                            "end": {
                              "line": 1,
                              "column": 24
                            }
                          },
                          "object": {
                            "type": "Super",
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
                          },
                          "property": {
                            "type": "Identifier",
                            "loc": {
                              "start": {
                                "line": 1,
                                "column": 23
                              },
                              "end": {
                                "line": 1,
                                "column": 24
                              }
                            },
                            "name": "a"
                          },
                          "computed": false
                        },
                        "arguments": []
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

// Harmony: Object Literal Property Value Shorthand

func TestHarmony56(t *testing.T) {
	ast, err := Compile("x = { y, z }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
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
              "value": {
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
              "kind": "init",
              "method": false,
              "shorthand": true,
              "computed": false,
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
            {
              "type": "Property",
              "key": {
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
              "value": {
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
              "kind": "init",
              "method": false,
              "shorthand": true,
              "computed": false,
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

// Harmony: Destructuring

func TestHarmony57(t *testing.T) {
	ast, err := Compile("[a, b] = [b, a]")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "AssignmentExpression",
        "operator": "=",
        "left": {
          "type": "ArrayPattern",
          "elements": [
            {
              "type": "Identifier",
              "name": "a",
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
            {
              "type": "Identifier",
              "name": "b",
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
        "right": {
          "type": "ArrayExpression",
          "elements": [
            {
              "type": "Identifier",
              "name": "b",
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
            {
              "type": "Identifier",
              "name": "a",
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

func TestHarmony58(t *testing.T) {
	ast, err := Compile("[a.r] = b")
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
          "type": "ArrayPattern",
          "start": 0,
          "end": 5,
          "elements": [
            {
              "type": "MemberExpression",
              "start": 1,
              "end": 4,
              "object": {
                "type": "Identifier",
                "start": 1,
                "end": 2,
                "name": "a"
              },
              "property": {
                "type": "Identifier",
                "start": 3,
                "end": 4,
                "name": "r"
              },
              "computed": false
            }
          ]
        },
        "right": {
          "type": "Identifier",
          "start": 8,
          "end": 9,
          "name": "b"
        }
      }
    }
  ]
}
	`, ast)
}

func TestHarmony59(t *testing.T) {
	ast, err := Compile("let [a,,b] = c")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
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
            "type": "ArrayPattern",
            "start": 4,
            "end": 10,
            "elements": [
              {
                "type": "Identifier",
                "start": 5,
                "end": 6,
                "name": "a"
              },
              null,
              {
                "type": "Identifier",
                "start": 8,
                "end": 9,
                "name": "b"
              }
            ]
          },
          "init": {
            "type": "Identifier",
            "start": 13,
            "end": 14,
            "name": "c"
          }
        }
      ],
      "kind": "let"
    }
  ]
}
	`, ast)
}

func TestHarmony60(t *testing.T) {
	ast, err := Compile("({ responseText: text } = res)")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "AssignmentExpression",
        "operator": "=",
        "left": {
          "type": "ObjectPattern",
          "properties": [
            {
              "type": "Property",
              "key": {
                "type": "Identifier",
                "name": "responseText",
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 3
                  },
                  "end": {
                    "line": 1,
                    "column": 15
                  }
                }
              },
              "value": {
                "type": "Identifier",
                "name": "text",
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
              "kind": "init",
              "method": false,
              "shorthand": false,
              "computed": false,
              "loc": {
                "start": {
                  "line": 1,
                  "column": 3
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
              "column": 1
            },
            "end": {
              "line": 1,
              "column": 23
            }
          }
        },
        "right": {
          "type": "Identifier",
          "name": "res",
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
            "column": 1
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

func TestHarmony61(t *testing.T) {
	ast, err := Compile("const {a} = {}")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "VariableDeclaration",
      "declarations": [
        {
          "type": "VariableDeclarator",
          "id": {
            "type": "ObjectPattern",
            "properties": [
              {
                "type": "Property",
                "key": {
                  "type": "Identifier",
                  "name": "a",
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
                "value": {
                  "type": "Identifier",
                  "name": "a",
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
                "kind": "init",
                "method": false,
                "shorthand": true,
                "computed": false,
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
          "init": {
            "type": "ObjectExpression",
            "properties": [],
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
              "column": 6
            },
            "end": {
              "line": 1,
              "column": 14
            }
          }
        }
      ],
      "kind": "const",
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

func TestHarmony62(t *testing.T) {
	ast, err := Compile("const [a] = []")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "VariableDeclaration",
      "declarations": [
        {
          "type": "VariableDeclarator",
          "id": {
            "type": "ArrayPattern",
            "elements": [
              {
                "type": "Identifier",
                "name": "a",
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
          "init": {
            "type": "ArrayExpression",
            "elements": [],
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
              "column": 6
            },
            "end": {
              "line": 1,
              "column": 14
            }
          }
        }
      ],
      "kind": "const",
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

func TestHarmony63(t *testing.T) {
	ast, err := Compile("let {a} = {}")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "VariableDeclaration",
      "declarations": [
        {
          "type": "VariableDeclarator",
          "id": {
            "type": "ObjectPattern",
            "properties": [
              {
                "type": "Property",
                "key": {
                  "type": "Identifier",
                  "name": "a",
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
                "value": {
                  "type": "Identifier",
                  "name": "a",
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
                "kind": "init",
                "method": false,
                "shorthand": true,
                "computed": false,
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
              }
            ],
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
          "init": {
            "type": "ObjectExpression",
            "properties": [],
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
              "column": 4
            },
            "end": {
              "line": 1,
              "column": 12
            }
          }
        }
      ],
      "kind": "let",
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

func TestHarmony64(t *testing.T) {
	ast, err := Compile("let [a] = []")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "VariableDeclaration",
      "declarations": [
        {
          "type": "VariableDeclarator",
          "id": {
            "type": "ArrayPattern",
            "elements": [
              {
                "type": "Identifier",
                "name": "a",
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
              }
            ],
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
          "init": {
            "type": "ArrayExpression",
            "elements": [],
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
              "column": 4
            },
            "end": {
              "line": 1,
              "column": 12
            }
          }
        }
      ],
      "kind": "let",
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

func TestHarmony65(t *testing.T) {
	ast, err := Compile("var {a} = {}")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "VariableDeclaration",
      "declarations": [
        {
          "type": "VariableDeclarator",
          "id": {
            "type": "ObjectPattern",
            "properties": [
              {
                "type": "Property",
                "key": {
                  "type": "Identifier",
                  "name": "a",
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
                "value": {
                  "type": "Identifier",
                  "name": "a",
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
                "kind": "init",
                "method": false,
                "shorthand": true,
                "computed": false,
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
              }
            ],
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
          "init": {
            "type": "ObjectExpression",
            "properties": [],
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
              "column": 4
            },
            "end": {
              "line": 1,
              "column": 12
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

func TestHarmony66(t *testing.T) {
	ast, err := Compile("var [a] = []")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "VariableDeclaration",
      "declarations": [
        {
          "type": "VariableDeclarator",
          "id": {
            "type": "ArrayPattern",
            "elements": [
              {
                "type": "Identifier",
                "name": "a",
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
              }
            ],
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
          "init": {
            "type": "ArrayExpression",
            "elements": [],
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
              "column": 4
            },
            "end": {
              "line": 1,
              "column": 12
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

func TestHarmony67(t *testing.T) {
	ast, err := Compile("const {a:b} = {}")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "VariableDeclaration",
      "declarations": [
        {
          "type": "VariableDeclarator",
          "id": {
            "type": "ObjectPattern",
            "properties": [
              {
                "type": "Property",
                "key": {
                  "type": "Identifier",
                  "name": "a",
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
                "value": {
                  "type": "Identifier",
                  "name": "b",
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
                "kind": "init",
                "method": false,
                "shorthand": false,
                "computed": false,
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 7
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
                "column": 6
              },
              "end": {
                "line": 1,
                "column": 11
              }
            }
          },
          "init": {
            "type": "ObjectExpression",
            "properties": [],
            "loc": {
              "start": {
                "line": 1,
                "column": 14
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
              "column": 6
            },
            "end": {
              "line": 1,
              "column": 16
            }
          }
        }
      ],
      "kind": "const",
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

func TestHarmony68(t *testing.T) {
	ast, err := Compile("let {a:b} = {}")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "VariableDeclaration",
      "declarations": [
        {
          "type": "VariableDeclarator",
          "id": {
            "type": "ObjectPattern",
            "properties": [
              {
                "type": "Property",
                "key": {
                  "type": "Identifier",
                  "name": "a",
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
                "value": {
                  "type": "Identifier",
                  "name": "b",
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
                "kind": "init",
                "method": false,
                "shorthand": false,
                "computed": false,
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
              }
            ],
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
            "type": "ObjectExpression",
            "properties": [],
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
              "column": 4
            },
            "end": {
              "line": 1,
              "column": 14
            }
          }
        }
      ],
      "kind": "let",
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

func TestHarmony69(t *testing.T) {
	ast, err := Compile("var {a:b} = {}")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "VariableDeclaration",
      "declarations": [
        {
          "type": "VariableDeclarator",
          "id": {
            "type": "ObjectPattern",
            "properties": [
              {
                "type": "Property",
                "key": {
                  "type": "Identifier",
                  "name": "a",
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
                "value": {
                  "type": "Identifier",
                  "name": "b",
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
                "kind": "init",
                "method": false,
                "shorthand": false,
                "computed": false,
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
              }
            ],
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
            "type": "ObjectExpression",
            "properties": [],
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
              "column": 4
            },
            "end": {
              "line": 1,
              "column": 14
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

// Harmony: Modules

func TestHarmony70(t *testing.T) {
	ast, err := Compile("export var document")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExportNamedDeclaration",
      "declaration": {
        "type": "VariableDeclaration",
        "declarations": [
          {
            "type": "VariableDeclarator",
            "id": {
              "type": "Identifier",
              "name": "document",
              "loc": {
                "start": {
                  "line": 1,
                  "column": 11
                },
                "end": {
                  "line": 1,
                  "column": 19
                }
              }
            },
            "init": null,
            "loc": {
              "start": {
                "line": 1,
                "column": 11
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
            "column": 7
          },
          "end": {
            "line": 1,
            "column": 19
          }
        }
      },
      "specifiers": [],
      "source": null,
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

func TestHarmony71(t *testing.T) {
	ast, err := Compile("export var document = { }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExportNamedDeclaration",
      "declaration": {
        "type": "VariableDeclaration",
        "declarations": [
          {
            "type": "VariableDeclarator",
            "id": {
              "type": "Identifier",
              "name": "document",
              "loc": {
                "start": {
                  "line": 1,
                  "column": 11
                },
                "end": {
                  "line": 1,
                  "column": 19
                }
              }
            },
            "init": {
              "type": "ObjectExpression",
              "properties": [],
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
                "column": 11
              },
              "end": {
                "line": 1,
                "column": 25
              }
            }
          }
        ],
        "kind": "var",
        "loc": {
          "start": {
            "line": 1,
            "column": 7
          },
          "end": {
            "line": 1,
            "column": 25
          }
        }
      },
      "specifiers": [],
      "source": null,
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

func TestHarmony72(t *testing.T) {
	ast, err := Compile("export let document")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExportNamedDeclaration",
      "declaration": {
        "type": "VariableDeclaration",
        "declarations": [
          {
            "type": "VariableDeclarator",
            "id": {
              "type": "Identifier",
              "name": "document",
              "loc": {
                "start": {
                  "line": 1,
                  "column": 11
                },
                "end": {
                  "line": 1,
                  "column": 19
                }
              }
            },
            "init": null,
            "loc": {
              "start": {
                "line": 1,
                "column": 11
              },
              "end": {
                "line": 1,
                "column": 19
              }
            }
          }
        ],
        "kind": "let",
        "loc": {
          "start": {
            "line": 1,
            "column": 7
          },
          "end": {
            "line": 1,
            "column": 19
          }
        }
      },
      "specifiers": [],
      "source": null,
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

func TestHarmony73(t *testing.T) {
	ast, err := Compile("export let document = { }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExportNamedDeclaration",
      "declaration": {
        "type": "VariableDeclaration",
        "declarations": [
          {
            "type": "VariableDeclarator",
            "id": {
              "type": "Identifier",
              "name": "document",
              "loc": {
                "start": {
                  "line": 1,
                  "column": 11
                },
                "end": {
                  "line": 1,
                  "column": 19
                }
              }
            },
            "init": {
              "type": "ObjectExpression",
              "properties": [],
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
                "column": 11
              },
              "end": {
                "line": 1,
                "column": 25
              }
            }
          }
        ],
        "kind": "let",
        "loc": {
          "start": {
            "line": 1,
            "column": 7
          },
          "end": {
            "line": 1,
            "column": 25
          }
        }
      },
      "specifiers": [],
      "source": null,
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

func TestHarmony74(t *testing.T) {
	ast, err := Compile("export const document = { }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExportNamedDeclaration",
      "declaration": {
        "type": "VariableDeclaration",
        "declarations": [
          {
            "type": "VariableDeclarator",
            "id": {
              "type": "Identifier",
              "name": "document",
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
            "init": {
              "type": "ObjectExpression",
              "properties": [],
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
                "column": 13
              },
              "end": {
                "line": 1,
                "column": 27
              }
            }
          }
        ],
        "kind": "const",
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
      "specifiers": [],
      "source": null,
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

func TestHarmony75(t *testing.T) {
	ast, err := Compile("export function parse() { }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExportNamedDeclaration",
      "declaration": {
        "type": "FunctionDeclaration",
        "id": {
          "type": "Identifier",
          "name": "parse",
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
              "column": 27
            }
          }
        },
        "generator": false,
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
      "specifiers": [],
      "source": null,
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

func TestHarmony76(t *testing.T) {
	ast, err := Compile("export class Class {}")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExportNamedDeclaration",
      "declaration": {
        "type": "ClassDeclaration",
        "id": {
          "type": "Identifier",
          "name": "Class",
          "loc": {
            "start": {
              "line": 1,
              "column": 13
            },
            "end": {
              "line": 1,
              "column": 18
            }
          }
        },
        "superClass": null,
        "body": {
          "type": "ClassBody",
          "body": [],
          "loc": {
            "start": {
              "line": 1,
              "column": 19
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
            "column": 7
          },
          "end": {
            "line": 1,
            "column": 21
          }
        }
      },
      "specifiers": [],
      "source": null,
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

func TestHarmony77(t *testing.T) {
	ast, err := Compile("export default 42")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExportDefaultDeclaration",
      "declaration": {
        "type": "Literal",
        "value": 42,
        "raw": "42",
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

func TestHarmony78(t *testing.T) {
	ast, err := Compile("export default function () {}")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 29,
  "body": [
    {
      "type": "ExportDefaultDeclaration",
      "start": 0,
      "end": 29,
      "declaration": {
        "type": "FunctionDeclaration",
        "start": 15,
        "end": 29,
        "id": null,
        "generator": false,
        "async": false,
        "params": [],
        "body": {
          "type": "BlockStatement",
          "start": 27,
          "end": 29,
          "body": []
        }
      }
    }
  ]
}
	`, ast)
}

func TestHarmony79(t *testing.T) {
	ast, err := Compile("export default function f() {}")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 30,
  "body": [
    {
      "type": "ExportDefaultDeclaration",
      "start": 0,
      "end": 30,
      "declaration": {
        "type": "FunctionDeclaration",
        "start": 15,
        "end": 30,
        "id": {
          "type": "Identifier",
          "start": 24,
          "end": 25,
          "name": "f"
        },
        "generator": false,
        "async": false,
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
	`, ast)
}

func TestHarmony80(t *testing.T) {
	ast, err := Compile("export default class {}")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 23,
  "body": [
    {
      "type": "ExportDefaultDeclaration",
      "start": 0,
      "end": 23,
      "declaration": {
        "type": "ClassDeclaration",
        "start": 15,
        "end": 23,
        "id": null,
        "superClass": null,
        "body": {
          "type": "ClassBody",
          "start": 21,
          "end": 23,
          "body": []
        }
      }
    }
  ]
}
	`, ast)
}

func TestHarmony81(t *testing.T) {
	ast, err := Compile("export default class A {}")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 25,
  "body": [
    {
      "type": "ExportDefaultDeclaration",
      "start": 0,
      "end": 25,
      "declaration": {
        "type": "ClassDeclaration",
        "start": 15,
        "end": 25,
        "id": {
          "type": "Identifier",
          "start": 21,
          "end": 22,
          "name": "A"
        },
        "superClass": null,
        "body": {
          "type": "ClassBody",
          "start": 23,
          "end": 25,
          "body": []
        }
      }
    }
  ]
}
	`, ast)
}

func TestHarmony82(t *testing.T) {
	ast, err := Compile("export default (class{});")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExportDefaultDeclaration",
      "declaration": {
        "type": "ClassExpression",
        "id": null,
        "superClass": null,
        "body": {
          "type": "ClassBody",
          "body": []
        }
      }
    }
  ]
}
	`, ast)
}

func TestHarmony83(t *testing.T) {
	ast, err := Compile("export * from \"crypto\"")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExportAllDeclaration",
      "source": {
        "type": "Literal",
        "value": "crypto",
        "raw": "\"crypto\"",
        "loc": {
          "start": {
            "line": 1,
            "column": 14
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

func TestHarmony84(t *testing.T) {
	ast, err := Compile("export { encrypt }\nvar encrypt")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExportNamedDeclaration",
      "declaration": null,
      "specifiers": [
        {
          "type": "ExportSpecifier",
          "exported": {
            "type": "Identifier",
            "name": "encrypt",
            "loc": {
              "start": {
                "line": 1,
                "column": 9
              },
              "end": {
                "line": 1,
                "column": 16
              }
            }
          },
          "local": {
            "type": "Identifier",
            "name": "encrypt",
            "loc": {
              "start": {
                "line": 1,
                "column": 9
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
              "column": 9
            },
            "end": {
              "line": 1,
              "column": 16
            }
          }
        }
      ],
      "source": null,
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
    {
      "type": "VariableDeclaration",
      "declarations": [
        {
          "type": "VariableDeclarator",
          "id": {
            "type": "Identifier",
            "name": "encrypt",
            "loc": {
              "start": {
                "line": 2,
                "column": 4
              },
              "end": {
                "line": 2,
                "column": 11
              }
            }
          },
          "init": null,
          "loc": {
            "start": {
              "line": 2,
              "column": 4
            },
            "end": {
              "line": 2,
              "column": 11
            }
          }
        }
      ],
      "kind": "var",
      "loc": {
        "start": {
          "line": 2,
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

func TestHarmony85(t *testing.T) {
	ast, err := Compile("function encrypt() {} let decrypt; export { encrypt, decrypt }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "FunctionDeclaration",
      "id": {
        "type": "Identifier",
        "name": "encrypt",
        "loc": {
          "start": {
            "line": 1,
            "column": 9
          },
          "end": {
            "line": 1,
            "column": 16
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
            "column": 19
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
    {
      "type": "VariableDeclaration",
      "declarations": [
        {
          "type": "VariableDeclarator",
          "id": {
            "type": "Identifier",
            "name": "decrypt",
            "loc": {
              "start": {
                "line": 1,
                "column": 26
              },
              "end": {
                "line": 1,
                "column": 33
              }
            }
          },
          "init": null,
          "loc": {
            "start": {
              "line": 1,
              "column": 26
            },
            "end": {
              "line": 1,
              "column": 33
            }
          }
        }
      ],
      "kind": "let",
      "loc": {
        "start": {
          "line": 1,
          "column": 22
        },
        "end": {
          "line": 1,
          "column": 34
        }
      }
    },
    {
      "type": "ExportNamedDeclaration",
      "declaration": null,
      "specifiers": [
        {
          "type": "ExportSpecifier",
          "exported": {
            "type": "Identifier",
            "name": "encrypt",
            "loc": {
              "start": {
                "line": 1,
                "column": 44
              },
              "end": {
                "line": 1,
                "column": 51
              }
            }
          },
          "local": {
            "type": "Identifier",
            "name": "encrypt",
            "loc": {
              "start": {
                "line": 1,
                "column": 44
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
              "column": 44
            },
            "end": {
              "line": 1,
              "column": 51
            }
          }
        },
        {
          "type": "ExportSpecifier",
          "exported": {
            "type": "Identifier",
            "name": "decrypt",
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
          "local": {
            "type": "Identifier",
            "name": "decrypt",
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
        }
      ],
      "source": null,
      "loc": {
        "start": {
          "line": 1,
          "column": 35
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
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 62
    }
  }
}
	`, ast)
}

func TestHarmony86(t *testing.T) {
	ast, err := Compile("export default class Test {}; export { Test }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{}
	`, ast)
}

func TestHarmony87(t *testing.T) {
	ast, err := Compile("{ var encrypt } export { encrypt }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{}
	`, ast)
}

func TestHarmony88(t *testing.T) {
	ast, err := Compile("export { encrypt as default }; function* encrypt() {}")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExportNamedDeclaration",
      "declaration": null,
      "specifiers": [
        {
          "type": "ExportSpecifier",
          "exported": {
            "type": "Identifier",
            "name": "default",
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
          "local": {
            "type": "Identifier",
            "name": "encrypt",
            "loc": {
              "start": {
                "line": 1,
                "column": 9
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
              "column": 9
            },
            "end": {
              "line": 1,
              "column": 27
            }
          }
        }
      ],
      "source": null,
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
    {
      "type": "FunctionDeclaration",
      "generator": true,
      "params": [],
      "body": {
        "type": "BlockStatement",
        "body": [],
        "loc": {
          "start": {
            "line": 1,
            "column": 51
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
          "column": 31
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

func TestHarmony89(t *testing.T) {
	ast, err := Compile("export { encrypt, decrypt as dec }; let encrypt, decrypt")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExportNamedDeclaration",
      "declaration": null,
      "specifiers": [
        {
          "type": "ExportSpecifier",
          "exported": {
            "type": "Identifier",
            "name": "encrypt",
            "loc": {
              "start": {
                "line": 1,
                "column": 9
              },
              "end": {
                "line": 1,
                "column": 16
              }
            }
          },
          "local": {
            "type": "Identifier",
            "name": "encrypt",
            "loc": {
              "start": {
                "line": 1,
                "column": 9
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
              "column": 9
            },
            "end": {
              "line": 1,
              "column": 16
            }
          }
        },
        {
          "type": "ExportSpecifier",
          "exported": {
            "type": "Identifier",
            "name": "dec",
            "loc": {
              "start": {
                "line": 1,
                "column": 29
              },
              "end": {
                "line": 1,
                "column": 32
              }
            }
          },
          "local": {
            "type": "Identifier",
            "name": "decrypt",
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
      "source": null,
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
    {
      "type": "VariableDeclaration",
      "declarations": [
        {
          "type": "VariableDeclarator",
          "id": {
            "type": "Identifier",
            "name": "encrypt",
            "loc": {
              "start": {
                "line": 1,
                "column": 40
              },
              "end": {
                "line": 1,
                "column": 47
              }
            }
          },
          "init": null,
          "loc": {
            "start": {
              "line": 1,
              "column": 40
            },
            "end": {
              "line": 1,
              "column": 47
            }
          }
        },
        {
          "type": "VariableDeclarator",
          "id": {
            "type": "Identifier",
            "name": "decrypt",
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
          "init": null,
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
        }
      ],
      "kind": "let",
      "loc": {
        "start": {
          "line": 1,
          "column": 36
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

func TestHarmony90(t *testing.T) {
	ast, err := Compile("export { default } from \"other\"")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExportNamedDeclaration",
      "declaration": null,
      "specifiers": [
        {
          "type": "ExportSpecifier",
          "exported": {
            "type": "Identifier",
            "name": "default",
            "loc": {
              "start": {
                "line": 1,
                "column": 9
              },
              "end": {
                "line": 1,
                "column": 16
              }
            }
          },
          "local": {
            "type": "Identifier",
            "name": "default",
            "loc": {
              "start": {
                "line": 1,
                "column": 9
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
              "column": 9
            },
            "end": {
              "line": 1,
              "column": 16
            }
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
            "column": 31
          }
        },
        "value": "other",
        "raw": "\"other\""
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

func TestHarmony91(t *testing.T) {
	ast, err := Compile("import \"jquery\"")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ImportDeclaration",
      "specifiers": [],
      "source": {
        "type": "Literal",
        "value": "jquery",
        "raw": "\"jquery\"",
        "loc": {
          "start": {
            "line": 1,
            "column": 7
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

func TestHarmony92(t *testing.T) {
	ast, err := Compile("import $ from \"jquery\"")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ImportDeclaration",
      "specifiers": [
        {
          "type": "ImportDefaultSpecifier",
          "local": {
            "type": "Identifier",
            "name": "$",
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
              "column": 7
            },
            "end": {
              "line": 1,
              "column": 8
            }
          }
        }
      ],
      "source": {
        "type": "Literal",
        "value": "jquery",
        "raw": "\"jquery\"",
        "loc": {
          "start": {
            "line": 1,
            "column": 14
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

func TestHarmony93(t *testing.T) {
	ast, err := Compile("import { encrypt, decrypt } from \"crypto\"")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ImportDeclaration",
      "specifiers": [
        {
          "type": "ImportSpecifier",
          "imported": {
            "type": "Identifier",
            "name": "encrypt",
            "loc": {
              "start": {
                "line": 1,
                "column": 9
              },
              "end": {
                "line": 1,
                "column": 16
              }
            }
          },
          "local": {
            "type": "Identifier",
            "name": "encrypt",
            "loc": {
              "start": {
                "line": 1,
                "column": 9
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
              "column": 9
            },
            "end": {
              "line": 1,
              "column": 16
            }
          }
        },
        {
          "type": "ImportSpecifier",
          "imported": {
            "type": "Identifier",
            "name": "decrypt",
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
          "local": {
            "type": "Identifier",
            "name": "decrypt",
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
        }
      ],
      "source": {
        "type": "Literal",
        "value": "crypto",
        "raw": "\"crypto\"",
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

func TestHarmony94(t *testing.T) {
	ast, err := Compile("import { encrypt as enc } from \"crypto\"")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ImportDeclaration",
      "specifiers": [
        {
          "type": "ImportSpecifier",
          "imported": {
            "type": "Identifier",
            "name": "encrypt",
            "loc": {
              "start": {
                "line": 1,
                "column": 9
              },
              "end": {
                "line": 1,
                "column": 16
              }
            }
          },
          "local": {
            "type": "Identifier",
            "name": "enc",
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
              "column": 9
            },
            "end": {
              "line": 1,
              "column": 23
            }
          }
        }
      ],
      "source": {
        "type": "Literal",
        "value": "crypto",
        "raw": "\"crypto\"",
        "loc": {
          "start": {
            "line": 1,
            "column": 31
          },
          "end": {
            "line": 1,
            "column": 39
          }
        }
      },
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
    }
  ],
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
}
	`, ast)
}

func TestHarmony95(t *testing.T) {
	ast, err := Compile("import crypto, { decrypt, encrypt as enc } from \"crypto\"")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 56
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
          "column": 56
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
              "column": 13
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
                "column": 13
              }
            },
            "name": "crypto"
          }
        },
        {
          "type": "ImportSpecifier",
          "loc": {
            "start": {
              "line": 1,
              "column": 17
            },
            "end": {
              "line": 1,
              "column": 24
            }
          },
          "imported": {
            "type": "Identifier",
            "loc": {
              "start": {
                "line": 1,
                "column": 17
              },
              "end": {
                "line": 1,
                "column": 24
              }
            },
            "name": "decrypt"
          },
          "local": {
            "type": "Identifier",
            "loc": {
              "start": {
                "line": 1,
                "column": 17
              },
              "end": {
                "line": 1,
                "column": 24
              }
            },
            "name": "decrypt"
          }
        },
        {
          "type": "ImportSpecifier",
          "loc": {
            "start": {
              "line": 1,
              "column": 26
            },
            "end": {
              "line": 1,
              "column": 40
            }
          },
          "imported": {
            "type": "Identifier",
            "loc": {
              "start": {
                "line": 1,
                "column": 26
              },
              "end": {
                "line": 1,
                "column": 33
              }
            },
            "name": "encrypt"
          },
          "local": {
            "type": "Identifier",
            "loc": {
              "start": {
                "line": 1,
                "column": 37
              },
              "end": {
                "line": 1,
                "column": 40
              }
            },
            "name": "enc"
          }
        }
      ],
      "source": {
        "type": "Literal",
        "loc": {
          "start": {
            "line": 1,
            "column": 48
          },
          "end": {
            "line": 1,
            "column": 56
          }
        },
        "value": "crypto",
        "raw": "\"crypto\""
      }
    }
  ]
}
	`, ast)
}

func TestHarmony96(t *testing.T) {
	ast, err := Compile("import { null as nil } from \"bar\"")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ImportDeclaration",
      "specifiers": [
        {
          "type": "ImportSpecifier",
          "imported": {
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
          "local": {
            "type": "Identifier",
            "name": "nil",
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
              "column": 9
            },
            "end": {
              "line": 1,
              "column": 20
            }
          }
        }
      ],
      "source": {
        "type": "Literal",
        "value": "bar",
        "raw": "\"bar\"",
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

func TestHarmony97(t *testing.T) {
	ast, err := Compile("import * as crypto from \"crypto\"")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 32
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
          "column": 32
        }
      },
      "specifiers": [
        {
          "type": "ImportNamespaceSpecifier",
          "loc": {
            "start": {
              "line": 1,
              "column": 7
            },
            "end": {
              "line": 1,
              "column": 18
            }
          },
          "local": {
            "type": "Identifier",
            "loc": {
              "start": {
                "line": 1,
                "column": 12
              },
              "end": {
                "line": 1,
                "column": 18
              }
            },
            "name": "crypto"
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
            "column": 32
          }
        },
        "value": "crypto",
        "raw": "\"crypto\""
      }
    }
  ]
}
	`, ast)
}

// Harmony: Yield Expression

func TestHarmony98(t *testing.T) {
	ast, err := Compile("(function* () { yield v })")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
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
              "type": "ExpressionStatement",
              "expression": {
                "type": "YieldExpression",
                "argument": {
                  "type": "Identifier",
                  "name": "v",
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
                "delegate": false,
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 16
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
                  "column": 16
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
              "column": 14
            },
            "end": {
              "line": 1,
              "column": 25
            }
          }
        },
        "generator": true,
        "expression": false,
        "loc": {
          "start": {
            "line": 1,
            "column": 1
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

func TestHarmony99(t *testing.T) {
	ast, err := Compile("(function* () { yield\nv })")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
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
              "type": "ExpressionStatement",
              "expression": {
                "type": "YieldExpression",
                "argument": null,
                "delegate": false,
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
                  "column": 16
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
                "type": "Identifier",
                "name": "v",
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
                  "column": 1
                }
              }
            }
          ],
          "loc": {
            "start": {
              "line": 1,
              "column": 14
            },
            "end": {
              "line": 2,
              "column": 3
            }
          }
        },
        "generator": true,
        "expression": false,
        "loc": {
          "start": {
            "line": 1,
            "column": 1
          },
          "end": {
            "line": 2,
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

func TestHarmony100(t *testing.T) {
	ast, err := Compile("(function* () { yield *v })")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
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
              "type": "ExpressionStatement",
              "expression": {
                "type": "YieldExpression",
                "argument": {
                  "type": "Identifier",
                  "name": "v",
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
                "delegate": true,
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 16
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
                  "column": 16
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
              "column": 14
            },
            "end": {
              "line": 1,
              "column": 26
            }
          }
        },
        "generator": true,
        "expression": false,
        "loc": {
          "start": {
            "line": 1,
            "column": 1
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

func TestHarmony101(t *testing.T) {
	ast, err := Compile("function* test () { yield *v }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
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
            "column": 10
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
              "type": "YieldExpression",
              "argument": {
                "type": "Identifier",
                "name": "v",
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
              "delegate": true,
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
      "generator": true,
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

func TestHarmony102(t *testing.T) {
	ast, err := Compile("var x = { *test () { yield *v } };")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
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
            "type": "ObjectExpression",
            "properties": [
              {
                "type": "Property",
                "key": {
                  "type": "Identifier",
                  "name": "test",
                  "loc": {
                    "start": {
                      "line": 1,
                      "column": 11
                    },
                    "end": {
                      "line": 1,
                      "column": 15
                    }
                  }
                },
                "value": {
                  "type": "FunctionExpression",
                  "id": null,
                  "params": [],
                  "body": {
                    "type": "BlockStatement",
                    "body": [
                      {
                        "type": "ExpressionStatement",
                        "expression": {
                          "type": "YieldExpression",
                          "argument": {
                            "type": "Identifier",
                            "name": "v",
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
                          },
                          "delegate": true,
                          "loc": {
                            "start": {
                              "line": 1,
                              "column": 21
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
                            "column": 21
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
                        "column": 19
                      },
                      "end": {
                        "line": 1,
                        "column": 31
                      }
                    }
                  },
                  "generator": true,
                  "expression": false,
                  "loc": {
                    "start": {
                      "line": 1,
                      "column": 16
                    },
                    "end": {
                      "line": 1,
                      "column": 31
                    }
                  }
                },
                "kind": "init",
                "method": true,
                "shorthand": false,
                "computed": false,
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 10
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
                "column": 8
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

func TestHarmony103(t *testing.T) {
	ast, err := Compile("function* foo() { console.log(yield); }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "body": [
    {
      "id": {
        "name": "foo",
        "type": "Identifier"
      },
      "generator": true,
      "params": [],
      "body": {
        "body": [
          {
            "expression": {
              "callee": {
                "object": {
                  "name": "console",
                  "type": "Identifier"
                },
                "property": {
                  "name": "log",
                  "type": "Identifier"
                },
                "computed": false,
                "type": "MemberExpression"
              },
              "arguments": [
                {
                  "delegate": false,
                  "argument": null,
                  "type": "YieldExpression"
                }
              ],
              "type": "CallExpression"
            },
            "type": "ExpressionStatement"
          }
        ],
        "type": "BlockStatement"
      },
      "type": "FunctionDeclaration"
    }
  ],
  "type": "Program"
}
	`, ast)
}

func TestHarmony104(t *testing.T) {
	ast, err := Compile("function* t() {}")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "FunctionDeclaration",
      "id": {
        "type": "Identifier",
        "name": "t",
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
      "params": [],
      "body": {
        "type": "BlockStatement",
        "body": [],
        "loc": {
          "start": {
            "line": 1,
            "column": 14
          },
          "end": {
            "line": 1,
            "column": 16
          }
        }
      },
      "generator": true,
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

func TestHarmony105(t *testing.T) {
	ast, err := Compile("(function* () { yield yield 10 })")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
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
              "type": "ExpressionStatement",
              "expression": {
                "type": "YieldExpression",
                "argument": {
                  "type": "YieldExpression",
                  "argument": {
                    "type": "Literal",
                    "value": 10,
                    "raw": "10",
                    "loc": {
                      "start": {
                        "line": 1,
                        "column": 28
                      },
                      "end": {
                        "line": 1,
                        "column": 30
                      }
                    }
                  },
                  "delegate": false,
                  "loc": {
                    "start": {
                      "line": 1,
                      "column": 22
                    },
                    "end": {
                      "line": 1,
                      "column": 30
                    }
                  }
                },
                "delegate": false,
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
                  "column": 16
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
              "column": 14
            },
            "end": {
              "line": 1,
              "column": 32
            }
          }
        },
        "generator": true,
        "expression": false,
        "loc": {
          "start": {
            "line": 1,
            "column": 1
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

// Harmony: Iterators

func TestHarmony106(t *testing.T) {
	ast, err := Compile("for(x of list) process(x);")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ForOfStatement",
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

func TestHarmony107(t *testing.T) {
	ast, err := Compile("for (var x of list) process(x);")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ForOfStatement",
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

func TestHarmony108(t *testing.T) {
	ast, err := Compile("for (let x of list) process(x);")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ForOfStatement",
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
        "kind": "let",
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

func TestHarmony109(t *testing.T) {
	ast, err := Compile("for (let\n{x} of list) process(x);")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ForOfStatement",
      "left": {
        "type": "VariableDeclaration",
        "declarations": [
          {
            "type": "VariableDeclarator",
            "id": {
              "type": "ObjectPattern",
              "loc": {
                "start": {
                  "line": 2,
                  "column": 0
                },
                "end": {
                  "line": 2,
                  "column": 3
                }
              },
              "properties": [
                {
                  "type": "Property",
                  "kind": "init",
                  "key": {
                    "type": "Identifier",
                    "name": "x",
                    "loc": {
                      "start": {
                        "line": 2,
                        "column": 1
                      },
                      "end": {
                        "line": 2,
                        "column": 2
                      }
                    }
                  },
                  "value": {
                    "type": "Identifier",
                    "name": "x",
                    "loc": {
                      "start": {
                        "line": 2,
                        "column": 1
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
                      "column": 1
                    },
                    "end": {
                      "line": 2,
                      "column": 2
                    }
                  }
                }
              ]
            },
            "init": null,
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
        "kind": "let",
        "loc": {
          "start": {
            "line": 1,
            "column": 5
          },
          "end": {
            "line": 2,
            "column": 3
          }
        }
      },
      "right": {
        "type": "Identifier",
        "name": "list",
        "loc": {
          "start": {
            "line": 2,
            "column": 7
          },
          "end": {
            "line": 2,
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
            "name": "process",
            "loc": {
              "start": {
                "line": 2,
                "column": 13
              },
              "end": {
                "line": 2,
                "column": 20
              }
            }
          },
          "arguments": [
            {
              "type": "Identifier",
              "name": "x",
              "loc": {
                "start": {
                  "line": 2,
                  "column": 21
                },
                "end": {
                  "line": 2,
                  "column": 22
                }
              }
            }
          ],
          "loc": {
            "start": {
              "line": 2,
              "column": 13
            },
            "end": {
              "line": 2,
              "column": 23
            }
          }
        },
        "loc": {
          "start": {
            "line": 2,
            "column": 13
          },
          "end": {
            "line": 2,
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
          "line": 2,
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
      "line": 2,
      "column": 24
    }
  }
}
	`, ast)
}

// Harmony: Class

func TestHarmony110(t *testing.T) {
	ast, err := Compile("var A = class extends B {}")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
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
            "name": "A",
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
            "type": "ClassExpression",
            "superClass": {
              "type": "Identifier",
              "name": "B",
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
            "body": {
              "type": "ClassBody",
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
                "column": 8
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
              "column": 4
            },
            "end": {
              "line": 1,
              "column": 26
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

func TestHarmony111(t *testing.T) {
	ast, err := Compile("class A extends class B extends C {} {}")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ClassDeclaration",
      "id": {
        "type": "Identifier",
        "name": "A",
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
      "superClass": {
        "type": "ClassExpression",
        "id": {
          "type": "Identifier",
          "name": "B",
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
        "superClass": {
          "type": "Identifier",
          "name": "C",
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
        },
        "body": {
          "type": "ClassBody",
          "body": [],
          "loc": {
            "start": {
              "line": 1,
              "column": 34
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
            "column": 16
          },
          "end": {
            "line": 1,
            "column": 36
          }
        }
      },
      "body": {
        "type": "ClassBody",
        "body": [],
        "loc": {
          "start": {
            "line": 1,
            "column": 37
          },
          "end": {
            "line": 1,
            "column": 39
          }
        }
      },
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
    }
  ],
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
}
	`, ast)
}

func TestHarmony112(t *testing.T) {
	ast, err := Compile("class A {get() {}}")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ClassDeclaration",
      "id": {
        "type": "Identifier",
        "name": "A",
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
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "body": [
          {
            "type": "MethodDefinition",
            "computed": false,
            "key": {
              "type": "Identifier",
              "name": "get",
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
              "generator": false,
              "expression": false,
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
            "kind": "method",
            "static": false,
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
          }
        ],
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

func TestHarmony113(t *testing.T) {
	ast, err := Compile("class A { static get() {}}")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ClassDeclaration",
      "id": {
        "type": "Identifier",
        "name": "A",
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
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "body": [
          {
            "type": "MethodDefinition",
            "computed": false,
            "key": {
              "type": "Identifier",
              "name": "get",
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
                    "column": 23
                  },
                  "end": {
                    "line": 1,
                    "column": 25
                  }
                }
              },
              "generator": false,
              "expression": false,
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
            "kind": "method",
            "static": true,
            "loc": {
              "start": {
                "line": 1,
                "column": 10
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
            "column": 8
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

func TestHarmony114(t *testing.T) {
	ast, err := Compile("class A extends B {get foo() {}}")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ClassDeclaration",
      "id": {
        "type": "Identifier",
        "name": "A",
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
      "superClass": {
        "type": "Identifier",
        "name": "B",
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
      "body": {
        "type": "ClassBody",
        "body": [
          {
            "type": "MethodDefinition",
            "computed": false,
            "key": {
              "type": "Identifier",
              "name": "foo",
              "loc": {
                "start": {
                  "line": 1,
                  "column": 23
                },
                "end": {
                  "line": 1,
                  "column": 26
                }
              }
            },
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
                    "column": 29
                  },
                  "end": {
                    "line": 1,
                    "column": 31
                  }
                }
              },
              "generator": false,
              "expression": false,
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
            "kind": "get",
            "static": false,
            "loc": {
              "start": {
                "line": 1,
                "column": 19
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

func TestHarmony115(t *testing.T) {
	ast, err := Compile("class A extends B { static get foo() {}}")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ClassDeclaration",
      "id": {
        "type": "Identifier",
        "name": "A",
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
      "superClass": {
        "type": "Identifier",
        "name": "B",
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
      "body": {
        "type": "ClassBody",
        "body": [
          {
            "type": "MethodDefinition",
            "computed": false,
            "key": {
              "type": "Identifier",
              "name": "foo",
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
                    "column": 37
                  },
                  "end": {
                    "line": 1,
                    "column": 39
                  }
                }
              },
              "generator": false,
              "expression": false,
              "loc": {
                "start": {
                  "line": 1,
                  "column": 34
                },
                "end": {
                  "line": 1,
                  "column": 39
                }
              }
            },
            "kind": "get",
            "static": true,
            "loc": {
              "start": {
                "line": 1,
                "column": 20
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
            "column": 18
          },
          "end": {
            "line": 1,
            "column": 40
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 40
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 40
    }
  }
}
	`, ast)
}

func TestHarmony116(t *testing.T) {
	ast, err := Compile("class A {set a(v) {}}")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ClassDeclaration",
      "id": {
        "type": "Identifier",
        "name": "A",
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
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "body": [
          {
            "type": "MethodDefinition",
            "computed": false,
            "key": {
              "type": "Identifier",
              "name": "a",
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
            "value": {
              "type": "FunctionExpression",
              "id": null,
              "params": [
                {
                  "type": "Identifier",
                  "name": "v",
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
              "generator": false,
              "expression": false,
              "loc": {
                "start": {
                  "line": 1,
                  "column": 14
                },
                "end": {
                  "line": 1,
                  "column": 20
                }
              }
            },
            "kind": "set",
            "static": false,
            "loc": {
              "start": {
                "line": 1,
                "column": 9
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

func TestHarmony117(t *testing.T) {
	ast, err := Compile("class A { static set a(v) {}}")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ClassDeclaration",
      "id": {
        "type": "Identifier",
        "name": "A",
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
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "body": [
          {
            "type": "MethodDefinition",
            "computed": false,
            "key": {
              "type": "Identifier",
              "name": "a",
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
            "value": {
              "type": "FunctionExpression",
              "id": null,
              "params": [
                {
                  "type": "Identifier",
                  "name": "v",
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
                    "column": 28
                  }
                }
              },
              "generator": false,
              "expression": false,
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
            "kind": "set",
            "static": true,
            "loc": {
              "start": {
                "line": 1,
                "column": 10
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
            "column": 8
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

func TestHarmony118(t *testing.T) {
	ast, err := Compile("class A {set(v) {};}")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ClassDeclaration",
      "id": {
        "type": "Identifier",
        "name": "A",
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
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "body": [
          {
            "type": "MethodDefinition",
            "computed": false,
            "key": {
              "type": "Identifier",
              "name": "set",
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
            "value": {
              "type": "FunctionExpression",
              "id": null,
              "params": [
                {
                  "type": "Identifier",
                  "name": "v",
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
                "body": [],
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
              "generator": false,
              "expression": false,
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
            "kind": "method",
            "static": false,
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
          }
        ],
        "loc": {
          "start": {
            "line": 1,
            "column": 8
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

func TestHarmony119(t *testing.T) {
	ast, err := Compile("class A { static set(v) {};}")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ClassDeclaration",
      "id": {
        "type": "Identifier",
        "name": "A",
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
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "body": [
          {
            "type": "MethodDefinition",
            "computed": false,
            "key": {
              "type": "Identifier",
              "name": "set",
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
            "value": {
              "type": "FunctionExpression",
              "id": null,
              "params": [
                {
                  "type": "Identifier",
                  "name": "v",
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
                    "column": 26
                  }
                }
              },
              "generator": false,
              "expression": false,
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
            "kind": "method",
            "static": true,
            "loc": {
              "start": {
                "line": 1,
                "column": 10
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
            "column": 8
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

func TestHarmony120(t *testing.T) {
	ast, err := Compile("class A {*gen(v) { yield v; }}")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ClassDeclaration",
      "id": {
        "type": "Identifier",
        "name": "A",
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
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "body": [
          {
            "type": "MethodDefinition",
            "computed": false,
            "key": {
              "type": "Identifier",
              "name": "gen",
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
            "value": {
              "type": "FunctionExpression",
              "id": null,
              "params": [
                {
                  "type": "Identifier",
                  "name": "v",
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
                }
              ],
              "body": {
                "type": "BlockStatement",
                "body": [
                  {
                    "type": "ExpressionStatement",
                    "expression": {
                      "type": "YieldExpression",
                      "argument": {
                        "type": "Identifier",
                        "name": "v",
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
                      "delegate": false,
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
              "generator": true,
              "expression": false,
              "loc": {
                "start": {
                  "line": 1,
                  "column": 13
                },
                "end": {
                  "line": 1,
                  "column": 29
                }
              }
            },
            "kind": "method",
            "static": false,
            "loc": {
              "start": {
                "line": 1,
                "column": 9
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
            "column": 8
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

func TestHarmony121(t *testing.T) {
	ast, err := Compile("class A { static *gen(v) { yield v; }}")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ClassDeclaration",
      "id": {
        "type": "Identifier",
        "name": "A",
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
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "body": [
          {
            "type": "MethodDefinition",
            "computed": false,
            "key": {
              "type": "Identifier",
              "name": "gen",
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
            "value": {
              "type": "FunctionExpression",
              "id": null,
              "params": [
                {
                  "type": "Identifier",
                  "name": "v",
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
              ],
              "body": {
                "type": "BlockStatement",
                "body": [
                  {
                    "type": "ExpressionStatement",
                    "expression": {
                      "type": "YieldExpression",
                      "argument": {
                        "type": "Identifier",
                        "name": "v",
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
                      },
                      "delegate": false,
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
                        "column": 35
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
                    "column": 37
                  }
                }
              },
              "generator": true,
              "expression": false,
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
            },
            "kind": "method",
            "static": true,
            "loc": {
              "start": {
                "line": 1,
                "column": 10
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
            "column": 8
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

func TestHarmony122(t *testing.T) {
	ast, err := Compile("(class { *static() {} })")
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
      "end": 24,
      "expression": {
        "type": "ClassExpression",
        "start": 1,
        "end": 23,
        "id": null,
        "superClass": null,
        "body": {
          "type": "ClassBody",
          "start": 7,
          "end": 23,
          "body": [
            {
              "type": "MethodDefinition",
              "start": 9,
              "end": 21,
              "computed": false,
              "key": {
                "type": "Identifier",
                "start": 10,
                "end": 16,
                "name": "static"
              },
              "static": false,
              "kind": "method",
              "value": {
                "type": "FunctionExpression",
                "start": 16,
                "end": 21,
                "id": null,
                "generator": true,
                "expression": false,
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
      }
    }
  ]
}
	`, ast)
}

func TestHarmony123(t *testing.T) {
	ast, err := Compile("\"use strict\"; (class A extends B {constructor() { super() }})")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "Literal",
        "value": "use strict",
        "raw": "\"use strict\"",
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
          "column": 13
        }
      }
    },
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "ClassExpression",
        "id": {
          "type": "Identifier",
          "name": "A",
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
        "superClass": {
          "type": "Identifier",
          "name": "B",
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
        "body": {
          "type": "ClassBody",
          "body": [
            {
              "type": "MethodDefinition",
              "computed": false,
              "key": {
                "type": "Identifier",
                "name": "constructor",
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 34
                  },
                  "end": {
                    "line": 1,
                    "column": 45
                  }
                }
              },
              "value": {
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
                          "type": "Super",
                          "loc": {
                            "start": {
                              "line": 1,
                              "column": 50
                            },
                            "end": {
                              "line": 1,
                              "column": 55
                            }
                          }
                        },
                        "arguments": [],
                        "loc": {
                          "start": {
                            "line": 1,
                            "column": 50
                          },
                          "end": {
                            "line": 1,
                            "column": 57
                          }
                        }
                      },
                      "loc": {
                        "start": {
                          "line": 1,
                          "column": 50
                        },
                        "end": {
                          "line": 1,
                          "column": 57
                        }
                      }
                    }
                  ],
                  "loc": {
                    "start": {
                      "line": 1,
                      "column": 48
                    },
                    "end": {
                      "line": 1,
                      "column": 59
                    }
                  }
                },
                "generator": false,
                "expression": false,
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 45
                  },
                  "end": {
                    "line": 1,
                    "column": 59
                  }
                }
              },
              "kind": "constructor",
              "static": false,
              "loc": {
                "start": {
                  "line": 1,
                  "column": 34
                },
                "end": {
                  "line": 1,
                  "column": 59
                }
              }
            }
          ],
          "loc": {
            "start": {
              "line": 1,
              "column": 33
            },
            "end": {
              "line": 1,
              "column": 60
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
            "column": 60
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
          "column": 61
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 61
    }
  }
}
	`, ast)
}

func TestHarmony124(t *testing.T) {
	ast, err := Compile("(class A extends B { constructor() { (() => { super() }); } })")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{}
	`, ast)
}

func TestHarmony125(t *testing.T) {
	ast, err := Compile("class A {'constructor'() {}}")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ClassDeclaration",
      "id": {
        "type": "Identifier",
        "name": "A"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "body": [
          {
            "type": "MethodDefinition",
            "computed": false,
            "key": {
              "type": "Literal",
              "value": "constructor"
            },
            "static": false,
            "kind": "constructor",
            "value": {
              "type": "FunctionExpression",
              "id": null,
              "generator": false,
              "expression": false,
              "params": [],
              "body": {
                "type": "BlockStatement",
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

func TestHarmony126(t *testing.T) {
	ast, err := Compile("class A { get ['constructor']() {} }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 36,
  "body": [
    {
      "type": "ClassDeclaration",
      "start": 0,
      "end": 36,
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
        "end": 36,
        "body": [
          {
            "type": "MethodDefinition",
            "start": 10,
            "end": 34,
            "static": false,
            "computed": true,
            "key": {
              "type": "Literal",
              "start": 15,
              "end": 28,
              "value": "constructor",
              "raw": "'constructor'"
            },
            "kind": "get",
            "value": {
              "type": "FunctionExpression",
              "start": 29,
              "end": 34,
              "id": null,
              "params": [],
              "generator": false,
              "expression": false,
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
  ]
}
	`, ast)
}

func TestHarmony127(t *testing.T) {
	ast, err := Compile("class A {static foo() {}}")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ClassDeclaration",
      "id": {
        "type": "Identifier",
        "name": "A",
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
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "body": [
          {
            "type": "MethodDefinition",
            "computed": false,
            "key": {
              "type": "Identifier",
              "name": "foo",
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
                    "column": 22
                  },
                  "end": {
                    "line": 1,
                    "column": 24
                  }
                }
              },
              "generator": false,
              "expression": false,
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
            "kind": "method",
            "static": true,
            "loc": {
              "start": {
                "line": 1,
                "column": 9
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
            "column": 8
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

func TestHarmony128(t *testing.T) {
	ast, err := Compile("class A {foo() {} static bar() {}}")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ClassDeclaration",
      "id": {
        "type": "Identifier",
        "name": "A",
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
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "body": [
          {
            "type": "MethodDefinition",
            "computed": false,
            "key": {
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
              "generator": false,
              "expression": false,
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
            "kind": "method",
            "static": false,
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
          {
            "type": "MethodDefinition",
            "computed": false,
            "key": {
              "type": "Identifier",
              "name": "bar",
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
                    "column": 31
                  },
                  "end": {
                    "line": 1,
                    "column": 33
                  }
                }
              },
              "generator": false,
              "expression": false,
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
            "kind": "method",
            "static": true,
            "loc": {
              "start": {
                "line": 1,
                "column": 18
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
            "column": 8
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

func TestHarmony129(t *testing.T) {
	ast, err := Compile("class A { foo() {} bar() {}}")
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
            "start": 10,
            "end": 18,
            "static": false,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 10,
              "end": 13,
              "name": "foo"
            },
            "kind": "method",
            "value": {
              "type": "FunctionExpression",
              "start": 13,
              "end": 18,
              "id": null,
              "expression": false,
              "generator": false,
              "async": false,
              "params": [],
              "body": {
                "type": "BlockStatement",
                "start": 16,
                "end": 18,
                "body": []
              }
            }
          },
          {
            "type": "MethodDefinition",
            "start": 19,
            "end": 27,
            "static": false,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 19,
              "end": 22,
              "name": "bar"
            },
            "kind": "method",
            "value": {
              "type": "FunctionExpression",
              "start": 22,
              "end": 27,
              "id": null,
              "expression": false,
              "generator": false,
              "async": false,
              "params": [],
              "body": {
                "type": "BlockStatement",
                "start": 25,
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

func TestHarmony130(t *testing.T) {
	ast, err := Compile("class A { static get foo() {} get foo() {}}")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ClassDeclaration",
      "id": {
        "type": "Identifier",
        "name": "A",
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
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "body": [
          {
            "type": "MethodDefinition",
            "computed": false,
            "key": {
              "type": "Identifier",
              "name": "foo",
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
                    "column": 27
                  },
                  "end": {
                    "line": 1,
                    "column": 29
                  }
                }
              },
              "generator": false,
              "expression": false,
              "loc": {
                "start": {
                  "line": 1,
                  "column": 24
                },
                "end": {
                  "line": 1,
                  "column": 29
                }
              }
            },
            "kind": "get",
            "static": true,
            "loc": {
              "start": {
                "line": 1,
                "column": 10
              },
              "end": {
                "line": 1,
                "column": 29
              }
            }
          },
          {
            "type": "MethodDefinition",
            "computed": false,
            "key": {
              "type": "Identifier",
              "name": "foo",
              "loc": {
                "start": {
                  "line": 1,
                  "column": 34
                },
                "end": {
                  "line": 1,
                  "column": 37
                }
              }
            },
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
                    "column": 40
                  },
                  "end": {
                    "line": 1,
                    "column": 42
                  }
                }
              },
              "generator": false,
              "expression": false,
              "loc": {
                "start": {
                  "line": 1,
                  "column": 37
                },
                "end": {
                  "line": 1,
                  "column": 42
                }
              }
            },
            "kind": "get",
            "static": false,
            "loc": {
              "start": {
                "line": 1,
                "column": 30
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
            "column": 8
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
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 43
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 43
    }
  }
}
	`, ast)
}

func TestHarmony131(t *testing.T) {
	ast, err := Compile("class A { static get foo() {} static get bar() {} }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ClassDeclaration",
      "id": {
        "type": "Identifier",
        "name": "A",
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
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "body": [
          {
            "type": "MethodDefinition",
            "computed": false,
            "key": {
              "type": "Identifier",
              "name": "foo",
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
                    "column": 27
                  },
                  "end": {
                    "line": 1,
                    "column": 29
                  }
                }
              },
              "generator": false,
              "expression": false,
              "loc": {
                "start": {
                  "line": 1,
                  "column": 24
                },
                "end": {
                  "line": 1,
                  "column": 29
                }
              }
            },
            "kind": "get",
            "static": true,
            "loc": {
              "start": {
                "line": 1,
                "column": 10
              },
              "end": {
                "line": 1,
                "column": 29
              }
            }
          },
          {
            "type": "MethodDefinition",
            "computed": false,
            "key": {
              "type": "Identifier",
              "name": "bar",
              "loc": {
                "start": {
                  "line": 1,
                  "column": 41
                },
                "end": {
                  "line": 1,
                  "column": 44
                }
              }
            },
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
                    "column": 47
                  },
                  "end": {
                    "line": 1,
                    "column": 49
                  }
                }
              },
              "generator": false,
              "expression": false,
              "loc": {
                "start": {
                  "line": 1,
                  "column": 44
                },
                "end": {
                  "line": 1,
                  "column": 49
                }
              }
            },
            "kind": "get",
            "static": true,
            "loc": {
              "start": {
                "line": 1,
                "column": 30
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
            "column": 8
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
          "column": 0
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
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 51
    }
  }
}
	`, ast)
}

func TestHarmony132(t *testing.T) {
	ast, err := Compile("class A { static get foo() {} static set foo(v) {} get foo() {} set foo(v) {}}")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ClassDeclaration",
      "id": {
        "type": "Identifier",
        "name": "A",
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
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "body": [
          {
            "type": "MethodDefinition",
            "computed": false,
            "key": {
              "type": "Identifier",
              "name": "foo",
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
                    "column": 27
                  },
                  "end": {
                    "line": 1,
                    "column": 29
                  }
                }
              },
              "generator": false,
              "expression": false,
              "loc": {
                "start": {
                  "line": 1,
                  "column": 24
                },
                "end": {
                  "line": 1,
                  "column": 29
                }
              }
            },
            "kind": "get",
            "static": true,
            "loc": {
              "start": {
                "line": 1,
                "column": 10
              },
              "end": {
                "line": 1,
                "column": 29
              }
            }
          },
          {
            "type": "MethodDefinition",
            "computed": false,
            "key": {
              "type": "Identifier",
              "name": "foo",
              "loc": {
                "start": {
                  "line": 1,
                  "column": 41
                },
                "end": {
                  "line": 1,
                  "column": 44
                }
              }
            },
            "value": {
              "type": "FunctionExpression",
              "id": null,
              "params": [
                {
                  "type": "Identifier",
                  "name": "v",
                  "loc": {
                    "start": {
                      "line": 1,
                      "column": 45
                    },
                    "end": {
                      "line": 1,
                      "column": 46
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
                    "column": 48
                  },
                  "end": {
                    "line": 1,
                    "column": 50
                  }
                }
              },
              "generator": false,
              "expression": false,
              "loc": {
                "start": {
                  "line": 1,
                  "column": 44
                },
                "end": {
                  "line": 1,
                  "column": 50
                }
              }
            },
            "kind": "set",
            "static": true,
            "loc": {
              "start": {
                "line": 1,
                "column": 30
              },
              "end": {
                "line": 1,
                "column": 50
              }
            }
          },
          {
            "type": "MethodDefinition",
            "computed": false,
            "key": {
              "type": "Identifier",
              "name": "foo",
              "loc": {
                "start": {
                  "line": 1,
                  "column": 55
                },
                "end": {
                  "line": 1,
                  "column": 58
                }
              }
            },
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
                    "column": 61
                  },
                  "end": {
                    "line": 1,
                    "column": 63
                  }
                }
              },
              "generator": false,
              "expression": false,
              "loc": {
                "start": {
                  "line": 1,
                  "column": 58
                },
                "end": {
                  "line": 1,
                  "column": 63
                }
              }
            },
            "kind": "get",
            "static": false,
            "loc": {
              "start": {
                "line": 1,
                "column": 51
              },
              "end": {
                "line": 1,
                "column": 63
              }
            }
          },
          {
            "type": "MethodDefinition",
            "computed": false,
            "key": {
              "type": "Identifier",
              "name": "foo",
              "loc": {
                "start": {
                  "line": 1,
                  "column": 68
                },
                "end": {
                  "line": 1,
                  "column": 71
                }
              }
            },
            "value": {
              "type": "FunctionExpression",
              "id": null,
              "params": [
                {
                  "type": "Identifier",
                  "name": "v",
                  "loc": {
                    "start": {
                      "line": 1,
                      "column": 72
                    },
                    "end": {
                      "line": 1,
                      "column": 73
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
                    "column": 75
                  },
                  "end": {
                    "line": 1,
                    "column": 77
                  }
                }
              },
              "generator": false,
              "expression": false,
              "loc": {
                "start": {
                  "line": 1,
                  "column": 71
                },
                "end": {
                  "line": 1,
                  "column": 77
                }
              }
            },
            "kind": "set",
            "static": false,
            "loc": {
              "start": {
                "line": 1,
                "column": 64
              },
              "end": {
                "line": 1,
                "column": 77
              }
            }
          }
        ],
        "loc": {
          "start": {
            "line": 1,
            "column": 8
          },
          "end": {
            "line": 1,
            "column": 78
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 78
        }
      }
    }
  ],
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 78
    }
  }
}
	`, ast)
}

func TestHarmony133(t *testing.T) {
	ast, err := Compile("class A { static [foo]() {} }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
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
      "type": "ClassDeclaration",
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
        "name": "A"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "loc": {
          "start": {
            "line": 1,
            "column": 8
          },
          "end": {
            "line": 1,
            "column": 29
          }
        },
        "body": [
          {
            "type": "MethodDefinition",
            "loc": {
              "start": {
                "line": 1,
                "column": 10
              },
              "end": {
                "line": 1,
                "column": 27
              }
            },
            "static": true,
            "computed": true,
            "key": {
              "type": "Identifier",
              "loc": {
                "start": {
                  "line": 1,
                  "column": 18
                },
                "end": {
                  "line": 1,
                  "column": 21
                }
              },
              "name": "foo"
            },
            "kind": "method",
            "value": {
              "type": "FunctionExpression",
              "loc": {
                "start": {
                  "line": 1,
                  "column": 22
                },
                "end": {
                  "line": 1,
                  "column": 27
                }
              },
              "id": null,
              "params": [],
              "generator": false,
              "body": {
                "type": "BlockStatement",
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 25
                  },
                  "end": {
                    "line": 1,
                    "column": 27
                  }
                },
                "body": []
              },
              "expression": false
            }
          }
        ]
      }
    }
  ]
}
	`, ast)
}

func TestHarmony134(t *testing.T) {
	ast, err := Compile("class A { static get [foo]() {} }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
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
            "type": "MethodDefinition",
            "start": 10,
            "end": 31,
            "static": true,
            "computed": true,
            "key": {
              "type": "Identifier",
              "start": 22,
              "end": 25,
              "name": "foo"
            },
            "kind": "get",
            "value": {
              "type": "FunctionExpression",
              "start": 26,
              "end": 31,
              "id": null,
              "expression": false,
              "generator": false,
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
        ]
      }
    }
  ]
}
	`, ast)
}

func TestHarmony135(t *testing.T) {
	ast, err := Compile("class A { set foo(v) {} get foo() {} }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ClassDeclaration",
      "id": {
        "type": "Identifier",
        "name": "A",
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
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "body": [
          {
            "type": "MethodDefinition",
            "computed": false,
            "key": {
              "type": "Identifier",
              "name": "foo",
              "loc": {
                "start": {
                  "line": 1,
                  "column": 14
                },
                "end": {
                  "line": 1,
                  "column": 17
                }
              }
            },
            "value": {
              "type": "FunctionExpression",
              "id": null,
              "params": [
                {
                  "type": "Identifier",
                  "name": "v",
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
                    "column": 23
                  }
                }
              },
              "generator": false,
              "expression": false,
              "loc": {
                "start": {
                  "line": 1,
                  "column": 17
                },
                "end": {
                  "line": 1,
                  "column": 23
                }
              }
            },
            "kind": "set",
            "static": false,
            "loc": {
              "start": {
                "line": 1,
                "column": 10
              },
              "end": {
                "line": 1,
                "column": 23
              }
            }
          },
          {
            "type": "MethodDefinition",
            "computed": false,
            "key": {
              "type": "Identifier",
              "name": "foo",
              "loc": {
                "start": {
                  "line": 1,
                  "column": 28
                },
                "end": {
                  "line": 1,
                  "column": 31
                }
              }
            },
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
                    "column": 34
                  },
                  "end": {
                    "line": 1,
                    "column": 36
                  }
                }
              },
              "generator": false,
              "expression": false,
              "loc": {
                "start": {
                  "line": 1,
                  "column": 31
                },
                "end": {
                  "line": 1,
                  "column": 36
                }
              }
            },
            "kind": "get",
            "static": false,
            "loc": {
              "start": {
                "line": 1,
                "column": 24
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
            "column": 8
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

func TestHarmony136(t *testing.T) {
	ast, err := Compile("class A { foo() {} get foo() {} }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 33
    }
  },
  "body": [
    {
      "type": "ClassDeclaration",
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 33
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
        "name": "A"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "loc": {
          "start": {
            "line": 1,
            "column": 8
          },
          "end": {
            "line": 1,
            "column": 33
          }
        },
        "body": [
          {
            "type": "MethodDefinition",
            "loc": {
              "start": {
                "line": 1,
                "column": 10
              },
              "end": {
                "line": 1,
                "column": 18
              }
            },
            "static": false,
            "computed": false,
            "key": {
              "type": "Identifier",
              "loc": {
                "start": {
                  "line": 1,
                  "column": 10
                },
                "end": {
                  "line": 1,
                  "column": 13
                }
              },
              "name": "foo"
            },
            "kind": "method",
            "value": {
              "type": "FunctionExpression",
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
              "id": null,
              "params": [],
              "generator": false,
              "body": {
                "type": "BlockStatement",
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
                "body": []
              },
              "expression": false
            }
          },
          {
            "type": "MethodDefinition",
            "loc": {
              "start": {
                "line": 1,
                "column": 19
              },
              "end": {
                "line": 1,
                "column": 31
              }
            },
            "static": false,
            "computed": false,
            "key": {
              "type": "Identifier",
              "loc": {
                "start": {
                  "line": 1,
                  "column": 23
                },
                "end": {
                  "line": 1,
                  "column": 26
                }
              },
              "name": "foo"
            },
            "kind": "get",
            "value": {
              "type": "FunctionExpression",
              "loc": {
                "start": {
                  "line": 1,
                  "column": 26
                },
                "end": {
                  "line": 1,
                  "column": 31
                }
              },
              "id": null,
              "params": [],
              "generator": false,
              "body": {
                "type": "BlockStatement",
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 29
                  },
                  "end": {
                    "line": 1,
                    "column": 31
                  }
                },
                "body": []
              },
              "expression": false
            }
          }
        ]
      }
    }
  ]
}
	`, ast)
}

func TestHarmony137(t *testing.T) {
	ast, err := Compile("class Semicolon { ; }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
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
      "type": "ClassDeclaration",
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
      "id": {
        "type": "Identifier",
        "loc": {
          "start": {
            "line": 1,
            "column": 6
          },
          "end": {
            "line": 1,
            "column": 15
          }
        },
        "name": "Semicolon"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "loc": {
          "start": {
            "line": 1,
            "column": 16
          },
          "end": {
            "line": 1,
            "column": 21
          }
        },
        "body": []
      }
    }
  ]
}
	`, ast)
}

func TestHarmony138(t *testing.T) {
	ast, err := Compile("class a { static }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 18,
  "body": [
    {
      "type": "ClassDeclaration",
      "start": 0,
      "end": 18,
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
        "end": 18,
        "body": [
          {
            "type": "PropertyDefinition",
            "start": 10,
            "end": 16,
            "static": false,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 10,
              "end": 16,
              "name": "static"
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

// ES6: Computed Properties

func TestHarmony139(t *testing.T) {
	ast, err := Compile("({[x]: 10})")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
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
        "type": "ObjectExpression",
        "start": 1,
        "end": 10,
        "properties": [
          {
            "type": "Property",
            "start": 2,
            "end": 9,
            "method": false,
            "shorthand": false,
            "computed": true,
            "key": {
              "type": "Identifier",
              "start": 3,
              "end": 4,
              "name": "x"
            },
            "value": {
              "type": "Literal",
              "start": 7,
              "end": 9,
              "value": 10,
              "raw": "10"
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

func TestHarmony140(t *testing.T) {
	ast, err := Compile("({[\"x\" + \"y\"]: 10})")
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
            "method": false,
            "shorthand": false,
            "computed": true,
            "key": {
              "type": "BinaryExpression",
              "start": 3,
              "end": 12,
              "left": {
                "type": "Literal",
                "start": 3,
                "end": 6,
                "value": "x",
                "raw": "\"x\""
              },
              "operator": "+",
              "right": {
                "type": "Literal",
                "start": 9,
                "end": 12,
                "value": "y",
                "raw": "\"y\""
              }
            },
            "value": {
              "type": "Literal",
              "start": 15,
              "end": 17,
              "value": 10,
              "raw": "10"
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

func TestHarmony141(t *testing.T) {
	ast, err := Compile("({[x]: function() {}})")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
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
                  "column": 3
                },
                "end": {
                  "line": 1,
                  "column": 4
                }
              }
            },
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
              "generator": false,
              "expression": false,
              "loc": {
                "start": {
                  "line": 1,
                  "column": 7
                },
                "end": {
                  "line": 1,
                  "column": 20
                }
              }
            },
            "kind": "init",
            "method": false,
            "shorthand": false,
            "computed": true,
            "loc": {
              "start": {
                "line": 1,
                "column": 2
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

func TestHarmony142(t *testing.T) {
	ast, err := Compile("({[x]: 10, y: 20})")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
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
                  "column": 3
                },
                "end": {
                  "line": 1,
                  "column": 4
                }
              }
            },
            "value": {
              "type": "Literal",
              "value": 10,
              "raw": "10",
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
            "kind": "init",
            "method": false,
            "shorthand": false,
            "computed": true,
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
          {
            "type": "Property",
            "key": {
              "type": "Identifier",
              "name": "y",
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
            "value": {
              "type": "Literal",
              "value": 20,
              "raw": "20",
              "loc": {
                "start": {
                  "line": 1,
                  "column": 14
                },
                "end": {
                  "line": 1,
                  "column": 16
                }
              }
            },
            "kind": "init",
            "method": false,
            "shorthand": false,
            "computed": false,
            "loc": {
              "start": {
                "line": 1,
                "column": 11
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
            "column": 1
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

func TestHarmony143(t *testing.T) {
	ast, err := Compile("({get [x]() {}, set [x](v) {}})")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
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
                  "column": 7
                },
                "end": {
                  "line": 1,
                  "column": 8
                }
              }
            },
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
                    "column": 12
                  },
                  "end": {
                    "line": 1,
                    "column": 14
                  }
                }
              },
              "generator": false,
              "expression": false,
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
            "kind": "get",
            "method": false,
            "shorthand": false,
            "computed": true,
            "loc": {
              "start": {
                "line": 1,
                "column": 2
              },
              "end": {
                "line": 1,
                "column": 14
              }
            }
          },
          {
            "type": "Property",
            "key": {
              "type": "Identifier",
              "name": "x",
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
            "value": {
              "type": "FunctionExpression",
              "id": null,
              "params": [
                {
                  "type": "Identifier",
                  "name": "v",
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
                    "column": 29
                  }
                }
              },
              "generator": false,
              "expression": false,
              "loc": {
                "start": {
                  "line": 1,
                  "column": 23
                },
                "end": {
                  "line": 1,
                  "column": 29
                }
              }
            },
            "kind": "set",
            "method": false,
            "shorthand": false,
            "computed": true,
            "loc": {
              "start": {
                "line": 1,
                "column": 16
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
            "column": 1
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

func TestHarmony144(t *testing.T) {
	ast, err := Compile("({[x]() {}})")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
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
                  "column": 3
                },
                "end": {
                  "line": 1,
                  "column": 4
                }
              }
            },
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
                    "column": 8
                  },
                  "end": {
                    "line": 1,
                    "column": 10
                  }
                }
              },
              "generator": false,
              "expression": false,
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
            "kind": "init",
            "method": true,
            "shorthand": false,
            "computed": true,
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
          }
        ],
        "loc": {
          "start": {
            "line": 1,
            "column": 1
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

func TestHarmony145(t *testing.T) {
	ast, err := Compile("var {[x]: y} = {y}")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "VariableDeclaration",
      "declarations": [
        {
          "type": "VariableDeclarator",
          "id": {
            "type": "ObjectPattern",
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
                "kind": "init",
                "method": false,
                "shorthand": false,
                "computed": true,
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
              }
            ],
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
          "init": {
            "type": "ObjectExpression",
            "properties": [
              {
                "type": "Property",
                "key": {
                  "type": "Identifier",
                  "name": "y",
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
                "value": {
                  "type": "Identifier",
                  "name": "y",
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
                "kind": "init",
                "method": false,
                "shorthand": true,
                "computed": false,
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

func TestHarmony146(t *testing.T) {
	ast, err := Compile("function f({[x]: y}) {}")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "FunctionDeclaration",
      "id": {
        "type": "Identifier",
        "name": "f",
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
      "params": [
        {
          "type": "ObjectPattern",
          "properties": [
            {
              "type": "Property",
              "key": {
                "type": "Identifier",
                "name": "x",
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
              "value": {
                "type": "Identifier",
                "name": "y",
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
              "kind": "init",
              "method": false,
              "shorthand": false,
              "computed": true,
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
            }
          ],
          "loc": {
            "start": {
              "line": 1,
              "column": 11
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
            "column": 23
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

func TestHarmony147(t *testing.T) {
	ast, err := Compile("var x = {*[test]() { yield *v; }}")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
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
            "type": "ObjectExpression",
            "properties": [
              {
                "type": "Property",
                "key": {
                  "type": "Identifier",
                  "name": "test",
                  "loc": {
                    "start": {
                      "line": 1,
                      "column": 11
                    },
                    "end": {
                      "line": 1,
                      "column": 15
                    }
                  }
                },
                "value": {
                  "type": "FunctionExpression",
                  "id": null,
                  "params": [],
                  "body": {
                    "type": "BlockStatement",
                    "body": [
                      {
                        "type": "ExpressionStatement",
                        "expression": {
                          "type": "YieldExpression",
                          "argument": {
                            "type": "Identifier",
                            "name": "v",
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
                          },
                          "delegate": true,
                          "loc": {
                            "start": {
                              "line": 1,
                              "column": 21
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
                            "column": 21
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
                        "column": 19
                      },
                      "end": {
                        "line": 1,
                        "column": 32
                      }
                    }
                  },
                  "generator": true,
                  "expression": false,
                  "loc": {
                    "start": {
                      "line": 1,
                      "column": 16
                    },
                    "end": {
                      "line": 1,
                      "column": 32
                    }
                  }
                },
                "kind": "init",
                "method": true,
                "shorthand": false,
                "computed": true,
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 9
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
                "column": 8
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

func TestHarmony148(t *testing.T) {
	ast, err := Compile("class A {[x]() {}}")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
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
      "type": "ClassDeclaration",
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
        "name": "A"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
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
        "body": [
          {
            "type": "MethodDefinition",
            "loc": {
              "start": {
                "line": 1,
                "column": 9
              },
              "end": {
                "line": 1,
                "column": 17
              }
            },
            "static": false,
            "computed": true,
            "key": {
              "type": "Identifier",
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
              "name": "x"
            },
            "kind": "method",
            "value": {
              "type": "FunctionExpression",
              "loc": {
                "start": {
                  "line": 1,
                  "column": 12
                },
                "end": {
                  "line": 1,
                  "column": 17
                }
              },
              "id": null,
              "params": [],
              "generator": false,
              "body": {
                "type": "BlockStatement",
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 15
                  },
                  "end": {
                    "line": 1,
                    "column": 17
                  }
                },
                "body": []
              },
              "expression": false
            }
          }
        ]
      }
    }
  ]
}
	`, ast)
}

// ES6: Default parameters

func TestHarmony149(t *testing.T) {
	ast, err := Compile("function f([x] = [1]) {}")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "FunctionDeclaration",
      "id": {
        "type": "Identifier",
        "name": "f",
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
      "params": [
        {
          "type": "AssignmentPattern",
          "left": {
            "type": "ArrayPattern",
            "elements": [
              {
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
              }
            ],
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
            "type": "ArrayExpression",
            "elements": [
              {
                "type": "Literal",
                "value": 1,
                "raw": "1",
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
        }
      ],
      "body": {
        "type": "BlockStatement",
        "body": [],
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
      "generator": false,
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

func TestHarmony150(t *testing.T) {
	ast, err := Compile("function f([x] = [1]) {  }")
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
        "start": 9,
        "end": 10,
        "name": "f"
      },
      "generator": false,
      "async": false,
      "params": [
        {
          "type": "AssignmentPattern",
          "start": 11,
          "end": 20,
          "left": {
            "type": "ArrayPattern",
            "start": 11,
            "end": 14,
            "elements": [
              {
                "type": "Identifier",
                "start": 12,
                "end": 13,
                "name": "x"
              }
            ]
          },
          "right": {
            "type": "ArrayExpression",
            "start": 17,
            "end": 20,
            "elements": [
              {
                "type": "Literal",
                "start": 18,
                "end": 19,
                "value": 1,
                "raw": "1"
              }
            ]
          }
        }
      ],
      "body": {
        "type": "BlockStatement",
        "start": 22,
        "end": 26,
        "body": []
      }
    }
  ]
}
	`, ast)
}

func TestHarmony151(t *testing.T) {
	ast, err := Compile("function f({x} = {x: 10}) {}")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "FunctionDeclaration",
      "id": {
        "type": "Identifier",
        "name": "f",
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
      "params": [
        {
          "type": "AssignmentPattern",
          "left": {
            "type": "ObjectPattern",
            "properties": [
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
                "kind": "init",
                "method": false,
                "shorthand": true,
                "computed": false,
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
              }
            ],
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
                      "column": 18
                    },
                    "end": {
                      "line": 1,
                      "column": 19
                    }
                  }
                },
                "value": {
                  "type": "Literal",
                  "value": 10,
                  "raw": "10",
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
                "kind": "init",
                "method": false,
                "shorthand": false,
                "computed": false,
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 18
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
                "column": 17
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
              "column": 11
            },
            "end": {
              "line": 1,
              "column": 24
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
            "column": 26
          },
          "end": {
            "line": 1,
            "column": 28
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

func TestHarmony152(t *testing.T) {
	ast, err := Compile("f = function({x} = {x: 10}) {}")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
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
          "name": "f",
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
          "type": "FunctionExpression",
          "id": null,
          "params": [
            {
              "type": "AssignmentPattern",
              "left": {
                "type": "ObjectPattern",
                "properties": [
                  {
                    "type": "Property",
                    "key": {
                      "type": "Identifier",
                      "name": "x",
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
                    "value": {
                      "type": "Identifier",
                      "name": "x",
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
                    "kind": "init",
                    "method": false,
                    "shorthand": true,
                    "computed": false,
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
                  }
                ],
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 13
                  },
                  "end": {
                    "line": 1,
                    "column": 16
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
                          "column": 20
                        },
                        "end": {
                          "line": 1,
                          "column": 21
                        }
                      }
                    },
                    "value": {
                      "type": "Literal",
                      "value": 10,
                      "raw": "10",
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
                    "kind": "init",
                    "method": false,
                    "shorthand": false,
                    "computed": false,
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
                  }
                ],
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
                  "column": 13
                },
                "end": {
                  "line": 1,
                  "column": 26
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
                "column": 28
              },
              "end": {
                "line": 1,
                "column": 30
              }
            }
          },
          "generator": false,
          "expression": false,
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

func TestHarmony153(t *testing.T) {
	ast, err := Compile("({f: function({x} = {x: 10}) {}})")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "ObjectExpression",
        "properties": [
          {
            "type": "Property",
            "key": {
              "type": "Identifier",
              "name": "f",
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
            "value": {
              "type": "FunctionExpression",
              "id": null,
              "params": [
                {
                  "type": "AssignmentPattern",
                  "left": {
                    "type": "ObjectPattern",
                    "properties": [
                      {
                        "type": "Property",
                        "key": {
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
                        "value": {
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
                        "kind": "init",
                        "method": false,
                        "shorthand": true,
                        "computed": false,
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
                    "loc": {
                      "start": {
                        "line": 1,
                        "column": 14
                      },
                      "end": {
                        "line": 1,
                        "column": 17
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
                              "column": 21
                            },
                            "end": {
                              "line": 1,
                              "column": 22
                            }
                          }
                        },
                        "value": {
                          "type": "Literal",
                          "value": 10,
                          "raw": "10",
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
                        "kind": "init",
                        "method": false,
                        "shorthand": false,
                        "computed": false,
                        "loc": {
                          "start": {
                            "line": 1,
                            "column": 21
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
                      "column": 14
                    },
                    "end": {
                      "line": 1,
                      "column": 27
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
                    "column": 29
                  },
                  "end": {
                    "line": 1,
                    "column": 31
                  }
                }
              },
              "generator": false,
              "expression": false,
              "loc": {
                "start": {
                  "line": 1,
                  "column": 5
                },
                "end": {
                  "line": 1,
                  "column": 31
                }
              }
            },
            "kind": "init",
            "method": false,
            "shorthand": false,
            "computed": false,
            "loc": {
              "start": {
                "line": 1,
                "column": 2
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
            "column": 1
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

func TestHarmony154(t *testing.T) {
	ast, err := Compile("({f({x} = {x: 10}) {}})")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "ObjectExpression",
        "properties": [
          {
            "type": "Property",
            "key": {
              "type": "Identifier",
              "name": "f",
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
            "value": {
              "type": "FunctionExpression",
              "id": null,
              "params": [
                {
                  "type": "AssignmentPattern",
                  "left": {
                    "type": "ObjectPattern",
                    "properties": [
                      {
                        "type": "Property",
                        "key": {
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
                        "value": {
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
                        "kind": "init",
                        "method": false,
                        "shorthand": true,
                        "computed": false,
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
                      }
                    ],
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
                              "column": 11
                            },
                            "end": {
                              "line": 1,
                              "column": 12
                            }
                          }
                        },
                        "value": {
                          "type": "Literal",
                          "value": 10,
                          "raw": "10",
                          "loc": {
                            "start": {
                              "line": 1,
                              "column": 14
                            },
                            "end": {
                              "line": 1,
                              "column": 16
                            }
                          }
                        },
                        "kind": "init",
                        "method": false,
                        "shorthand": false,
                        "computed": false,
                        "loc": {
                          "start": {
                            "line": 1,
                            "column": 11
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
                        "column": 10
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
              "body": {
                "type": "BlockStatement",
                "body": [],
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 19
                  },
                  "end": {
                    "line": 1,
                    "column": 21
                  }
                }
              },
              "generator": false,
              "expression": false,
              "loc": {
                "start": {
                  "line": 1,
                  "column": 3
                },
                "end": {
                  "line": 1,
                  "column": 21
                }
              }
            },
            "kind": "init",
            "method": true,
            "shorthand": false,
            "computed": false,
            "loc": {
              "start": {
                "line": 1,
                "column": 2
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

func TestHarmony155(t *testing.T) {
	ast, err := Compile("(class {f({x} = {x: 10}) {}})")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "ClassExpression",
        "superClass": null,
        "body": {
          "type": "ClassBody",
          "body": [
            {
              "type": "MethodDefinition",
              "computed": false,
              "key": {
                "type": "Identifier",
                "name": "f",
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
              "value": {
                "type": "FunctionExpression",
                "id": null,
                "params": [
                  {
                    "type": "AssignmentPattern",
                    "left": {
                      "type": "ObjectPattern",
                      "properties": [
                        {
                          "type": "Property",
                          "key": {
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
                          "value": {
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
                          "kind": "init",
                          "method": false,
                          "shorthand": true,
                          "computed": false,
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
                        }
                      ],
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
                                "column": 17
                              },
                              "end": {
                                "line": 1,
                                "column": 18
                              }
                            }
                          },
                          "value": {
                            "type": "Literal",
                            "value": 10,
                            "raw": "10",
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
                          "kind": "init",
                          "method": false,
                          "shorthand": false,
                          "computed": false,
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
                      "loc": {
                        "start": {
                          "line": 1,
                          "column": 16
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
                        "column": 10
                      },
                      "end": {
                        "line": 1,
                        "column": 23
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
                      "column": 25
                    },
                    "end": {
                      "line": 1,
                      "column": 27
                    }
                  }
                },
                "generator": false,
                "expression": false,
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 9
                  },
                  "end": {
                    "line": 1,
                    "column": 27
                  }
                }
              },
              "kind": "method",
              "static": false,
              "loc": {
                "start": {
                  "line": 1,
                  "column": 8
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
              "column": 7
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
            "column": 1
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

func TestHarmony156(t *testing.T) {
	ast, err := Compile("(({x} = {x: 10}) => {})")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "ArrowFunctionExpression",
        "id": null,
        "params": [
          {
            "type": "AssignmentPattern",
            "left": {
              "type": "ObjectPattern",
              "properties": [
                {
                  "type": "Property",
                  "key": {
                    "type": "Identifier",
                    "name": "x",
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
                  },
                  "value": {
                    "type": "Identifier",
                    "name": "x",
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
                  },
                  "kind": "init",
                  "method": false,
                  "shorthand": true,
                  "computed": false,
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
                  "column": 2
                },
                "end": {
                  "line": 1,
                  "column": 5
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
                        "column": 9
                      },
                      "end": {
                        "line": 1,
                        "column": 10
                      }
                    }
                  },
                  "value": {
                    "type": "Literal",
                    "value": 10,
                    "raw": "10",
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
                  "kind": "init",
                  "method": false,
                  "shorthand": false,
                  "computed": false,
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
                }
              ],
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
            "loc": {
              "start": {
                "line": 1,
                "column": 2
              },
              "end": {
                "line": 1,
                "column": 15
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
              "column": 22
            }
          }
        },
        "generator": false,
        "expression": false,
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

func TestHarmony157(t *testing.T) {
	ast, err := Compile("x = function(y = 1) {}")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
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
          "type": "FunctionExpression",
          "id": null,
          "params": [
            {
              "type": "AssignmentPattern",
              "left": {
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
              "right": {
                "type": "Literal",
                "value": 1,
                "raw": "1",
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
                  "column": 13
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
                "column": 22
              }
            }
          },
          "generator": false,
          "expression": false,
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

func TestHarmony158(t *testing.T) {
	ast, err := Compile("function f(a = 1) {}")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "FunctionDeclaration",
      "id": {
        "type": "Identifier",
        "name": "f",
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
      "params": [
        {
          "type": "AssignmentPattern",
          "left": {
            "type": "Identifier",
            "name": "a",
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
          "right": {
            "type": "Literal",
            "value": 1,
            "raw": "1",
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
          "loc": {
            "start": {
              "line": 1,
              "column": 11
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
      "generator": false,
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

func TestHarmony159(t *testing.T) {
	ast, err := Compile("x = { f: function(a=1) {} }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
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
                "name": "f",
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
                "type": "FunctionExpression",
                "id": null,
                "params": [
                  {
                    "type": "AssignmentPattern",
                    "left": {
                      "type": "Identifier",
                      "name": "a",
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
                    "right": {
                      "type": "Literal",
                      "value": 1,
                      "raw": "1",
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
                        "column": 18
                      },
                      "end": {
                        "line": 1,
                        "column": 21
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
                      "column": 23
                    },
                    "end": {
                      "line": 1,
                      "column": 25
                    }
                  }
                },
                "generator": false,
                "expression": false,
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 9
                  },
                  "end": {
                    "line": 1,
                    "column": 25
                  }
                }
              },
              "kind": "init",
              "method": false,
              "shorthand": false,
              "computed": false,
              "loc": {
                "start": {
                  "line": 1,
                  "column": 6
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
              "column": 4
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

func TestHarmony160(t *testing.T) {
	ast, err := Compile("x = { f(a=1) {} }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
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
                "name": "f",
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
                "type": "FunctionExpression",
                "id": null,
                "params": [
                  {
                    "type": "AssignmentPattern",
                    "left": {
                      "type": "Identifier",
                      "name": "a",
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
                    "right": {
                      "type": "Literal",
                      "value": 1,
                      "raw": "1",
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
                        "column": 8
                      },
                      "end": {
                        "line": 1,
                        "column": 11
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
                      "column": 13
                    },
                    "end": {
                      "line": 1,
                      "column": 15
                    }
                  }
                },
                "generator": false,
                "expression": false,
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 7
                  },
                  "end": {
                    "line": 1,
                    "column": 15
                  }
                }
              },
              "kind": "init",
              "method": true,
              "shorthand": false,
              "computed": false,
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

// ES6: Rest parameters

func TestHarmony161(t *testing.T) {
	ast, err := Compile("function f(a, ...b) {}")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "FunctionDeclaration",
      "id": {
        "type": "Identifier",
        "name": "f",
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
      "params": [
        {
          "type": "Identifier",
          "name": "a",
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
        {
          "type": "RestElement",
          "argument": {
            "type": "Identifier",
            "name": "b",
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
            "column": 22
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

// ES6: Destructured Parameters

func TestHarmony162(t *testing.T) {
	ast, err := Compile("function x([ a, b ]){}")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "FunctionDeclaration",
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
      "params": [
        {
          "type": "ArrayPattern",
          "elements": [
            {
              "type": "Identifier",
              "name": "a",
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
            {
              "type": "Identifier",
              "name": "b",
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
          "loc": {
            "start": {
              "line": 1,
              "column": 11
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
            "column": 20
          },
          "end": {
            "line": 1,
            "column": 22
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

func TestHarmony163(t *testing.T) {
	ast, err := Compile("function x({ a, b }){}")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "FunctionDeclaration",
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
      "params": [
        {
          "type": "ObjectPattern",
          "properties": [
            {
              "type": "Property",
              "key": {
                "type": "Identifier",
                "name": "a",
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
              "value": {
                "type": "Identifier",
                "name": "a",
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
              "kind": "init",
              "method": false,
              "shorthand": true,
              "computed": false,
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
            {
              "type": "Property",
              "key": {
                "type": "Identifier",
                "name": "b",
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
              "value": {
                "type": "Identifier",
                "name": "b",
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
              "kind": "init",
              "method": false,
              "shorthand": true,
              "computed": false,
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
          "loc": {
            "start": {
              "line": 1,
              "column": 11
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
            "column": 20
          },
          "end": {
            "line": 1,
            "column": 22
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

func TestHarmony164(t *testing.T) {
	ast, err := Compile("(function x([ a, b ]){})")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "FunctionExpression",
        "id": {
          "type": "Identifier",
          "name": "x",
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
        "params": [
          {
            "type": "ArrayPattern",
            "elements": [
              {
                "type": "Identifier",
                "name": "a",
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
                "name": "b",
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
              "column": 23
            }
          }
        },
        "generator": false,
        "expression": false,
        "loc": {
          "start": {
            "line": 1,
            "column": 1
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

func TestHarmony165(t *testing.T) {
	ast, err := Compile("(function x({ a, b }){})")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "FunctionExpression",
        "id": {
          "type": "Identifier",
          "name": "x",
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
        "params": [
          {
            "type": "ObjectPattern",
            "properties": [
              {
                "type": "Property",
                "key": {
                  "type": "Identifier",
                  "name": "a",
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
                "value": {
                  "type": "Identifier",
                  "name": "a",
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
                "kind": "init",
                "method": false,
                "shorthand": true,
                "computed": false,
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
                "type": "Property",
                "key": {
                  "type": "Identifier",
                  "name": "b",
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
                "value": {
                  "type": "Identifier",
                  "name": "b",
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
                "kind": "init",
                "method": false,
                "shorthand": true,
                "computed": false,
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
              "column": 23
            }
          }
        },
        "generator": false,
        "expression": false,
        "loc": {
          "start": {
            "line": 1,
            "column": 1
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

func TestHarmony166(t *testing.T) {
	ast, err := Compile("({ x([ a, b ]){} })")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
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
                  "column": 3
                },
                "end": {
                  "line": 1,
                  "column": 4
                }
              }
            },
            "value": {
              "type": "FunctionExpression",
              "id": null,
              "params": [
                {
                  "type": "ArrayPattern",
                  "elements": [
                    {
                      "type": "Identifier",
                      "name": "a",
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
                    {
                      "type": "Identifier",
                      "name": "b",
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
                    }
                  ],
                  "loc": {
                    "start": {
                      "line": 1,
                      "column": 5
                    },
                    "end": {
                      "line": 1,
                      "column": 13
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
                    "column": 14
                  },
                  "end": {
                    "line": 1,
                    "column": 16
                  }
                }
              },
              "generator": false,
              "expression": false,
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
            "kind": "init",
            "method": true,
            "shorthand": false,
            "computed": false,
            "loc": {
              "start": {
                "line": 1,
                "column": 3
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
            "column": 1
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

func TestHarmony167(t *testing.T) {
	ast, err := Compile("({ x(...[ a, b ]){} })")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
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
                  "column": 3
                },
                "end": {
                  "line": 1,
                  "column": 4
                }
              }
            },
            "value": {
              "type": "FunctionExpression",
              "id": null,
              "params": [
                {
                  "type": "RestElement",
                  "argument": {
                    "type": "ArrayPattern",
                    "elements": [
                      {
                        "type": "Identifier",
                        "name": "a",
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
                      {
                        "type": "Identifier",
                        "name": "b",
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
                  }
                }
              ],
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
              "generator": false,
              "expression": false,
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
            "kind": "init",
            "method": true,
            "shorthand": false,
            "computed": false,
            "loc": {
              "start": {
                "line": 1,
                "column": 3
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

func TestHarmony168(t *testing.T) {
	ast, err := Compile("({ x({ a: { w, x }, b: [y, z] }, ...[a, b, c]){} })")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
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
                  "column": 3
                },
                "end": {
                  "line": 1,
                  "column": 4
                }
              }
            },
            "value": {
              "type": "FunctionExpression",
              "id": null,
              "params": [
                {
                  "type": "ObjectPattern",
                  "properties": [
                    {
                      "type": "Property",
                      "key": {
                        "type": "Identifier",
                        "name": "a",
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
                      "value": {
                        "type": "ObjectPattern",
                        "properties": [
                          {
                            "type": "Property",
                            "key": {
                              "type": "Identifier",
                              "name": "w",
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
                              "type": "Identifier",
                              "name": "w",
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
                            "kind": "init",
                            "method": false,
                            "shorthand": true,
                            "computed": false,
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
                          {
                            "type": "Property",
                            "key": {
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
                            "value": {
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
                            "kind": "init",
                            "method": false,
                            "shorthand": true,
                            "computed": false,
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
                        "loc": {
                          "start": {
                            "line": 1,
                            "column": 10
                          },
                          "end": {
                            "line": 1,
                            "column": 18
                          }
                        }
                      },
                      "kind": "init",
                      "method": false,
                      "shorthand": false,
                      "computed": false,
                      "loc": {
                        "start": {
                          "line": 1,
                          "column": 7
                        },
                        "end": {
                          "line": 1,
                          "column": 18
                        }
                      }
                    },
                    {
                      "type": "Property",
                      "key": {
                        "type": "Identifier",
                        "name": "b",
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
                      "value": {
                        "type": "ArrayPattern",
                        "elements": [
                          {
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
                          {
                            "type": "Identifier",
                            "name": "z",
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
                          }
                        ],
                        "loc": {
                          "start": {
                            "line": 1,
                            "column": 23
                          },
                          "end": {
                            "line": 1,
                            "column": 29
                          }
                        }
                      },
                      "kind": "init",
                      "method": false,
                      "shorthand": false,
                      "computed": false,
                      "loc": {
                        "start": {
                          "line": 1,
                          "column": 20
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
                      "column": 5
                    },
                    "end": {
                      "line": 1,
                      "column": 31
                    }
                  }
                },
                {
                  "type": "RestElement",
                  "argument": {
                    "type": "ArrayPattern",
                    "elements": [
                      {
                        "type": "Identifier",
                        "name": "a",
                        "loc": {
                          "start": {
                            "line": 1,
                            "column": 37
                          },
                          "end": {
                            "line": 1,
                            "column": 38
                          }
                        }
                      },
                      {
                        "type": "Identifier",
                        "name": "b",
                        "loc": {
                          "start": {
                            "line": 1,
                            "column": 40
                          },
                          "end": {
                            "line": 1,
                            "column": 41
                          }
                        }
                      },
                      {
                        "type": "Identifier",
                        "name": "c",
                        "loc": {
                          "start": {
                            "line": 1,
                            "column": 43
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
                        "column": 36
                      },
                      "end": {
                        "line": 1,
                        "column": 45
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
                    "column": 46
                  },
                  "end": {
                    "line": 1,
                    "column": 48
                  }
                }
              },
              "generator": false,
              "expression": false,
              "loc": {
                "start": {
                  "line": 1,
                  "column": 4
                },
                "end": {
                  "line": 1,
                  "column": 48
                }
              }
            },
            "kind": "init",
            "method": true,
            "shorthand": false,
            "computed": false,
            "loc": {
              "start": {
                "line": 1,
                "column": 3
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
            "column": 1
          },
          "end": {
            "line": 1,
            "column": 50
          }
        }
      },
      "loc": {
        "start": {
          "line": 1,
          "column": 0
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
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 51
    }
  }
}
	`, ast)
}

func TestHarmony169(t *testing.T) {
	ast, err := Compile("(...a) => {}")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "ArrowFunctionExpression",
        "id": null,
        "params": [
          {
            "type": "RestElement",
            "argument": {
              "type": "Identifier",
              "name": "a",
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
          }
        ],
        "body": {
          "type": "BlockStatement",
          "body": [],
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
        "generator": false,
        "expression": false,
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

func TestHarmony170(t *testing.T) {
	ast, err := Compile("(a, ...b) => {}")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "ArrowFunctionExpression",
        "id": null,
        "params": [
          {
            "type": "Identifier",
            "name": "a",
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
          {
            "type": "RestElement",
            "argument": {
              "type": "Identifier",
              "name": "b",
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
          }
        ],
        "body": {
          "type": "BlockStatement",
          "body": [],
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
        "generator": false,
        "expression": false,
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

func TestHarmony171(t *testing.T) {
	ast, err := Compile("({ a }) => {}")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "ArrowFunctionExpression",
        "id": null,
        "params": [
          {
            "type": "ObjectPattern",
            "properties": [
              {
                "type": "Property",
                "key": {
                  "type": "Identifier",
                  "name": "a",
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
                },
                "value": {
                  "type": "Identifier",
                  "name": "a",
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
                },
                "kind": "init",
                "method": false,
                "shorthand": true,
                "computed": false,
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
                "column": 1
              },
              "end": {
                "line": 1,
                "column": 6
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
              "column": 11
            },
            "end": {
              "line": 1,
              "column": 13
            }
          }
        },
        "generator": false,
        "expression": false,
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

func TestHarmony172(t *testing.T) {
	ast, err := Compile("({ a }, ...b) => {}")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "ArrowFunctionExpression",
        "id": null,
        "params": [
          {
            "type": "ObjectPattern",
            "properties": [
              {
                "type": "Property",
                "key": {
                  "type": "Identifier",
                  "name": "a",
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
                },
                "value": {
                  "type": "Identifier",
                  "name": "a",
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
                },
                "kind": "init",
                "method": false,
                "shorthand": true,
                "computed": false,
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
                "column": 1
              },
              "end": {
                "line": 1,
                "column": 6
              }
            }
          },
          {
            "type": "RestElement",
            "argument": {
              "type": "Identifier",
              "name": "b",
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
            }
          }
        ],
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
        "generator": false,
        "expression": false,
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

func TestHarmony173(t *testing.T) {
	ast, err := Compile("({ a: [a, b] }, ...c) => {}")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "ArrowFunctionExpression",
        "id": null,
        "params": [
          {
            "type": "ObjectPattern",
            "properties": [
              {
                "type": "Property",
                "key": {
                  "type": "Identifier",
                  "name": "a",
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
                },
                "value": {
                  "type": "ArrayPattern",
                  "elements": [
                    {
                      "type": "Identifier",
                      "name": "a",
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
                    {
                      "type": "Identifier",
                      "name": "b",
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
                    }
                  ],
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
                "kind": "init",
                "method": false,
                "shorthand": false,
                "computed": false,
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 3
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
                "column": 1
              },
              "end": {
                "line": 1,
                "column": 14
              }
            }
          },
          {
            "type": "RestElement",
            "argument": {
              "type": "Identifier",
              "name": "c",
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
            }
          }
        ],
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
        "generator": false,
        "expression": false,
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

func TestHarmony174(t *testing.T) {
	ast, err := Compile("({ a: b, c }, [d, e], ...f) => {}")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "ArrowFunctionExpression",
        "id": null,
        "params": [
          {
            "type": "ObjectPattern",
            "properties": [
              {
                "type": "Property",
                "key": {
                  "type": "Identifier",
                  "name": "a",
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
                },
                "value": {
                  "type": "Identifier",
                  "name": "b",
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
                "kind": "init",
                "method": false,
                "shorthand": false,
                "computed": false,
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
              {
                "type": "Property",
                "key": {
                  "type": "Identifier",
                  "name": "c",
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
                "value": {
                  "type": "Identifier",
                  "name": "c",
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
                "kind": "init",
                "method": false,
                "shorthand": true,
                "computed": false,
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
            "loc": {
              "start": {
                "line": 1,
                "column": 1
              },
              "end": {
                "line": 1,
                "column": 12
              }
            }
          },
          {
            "type": "ArrayPattern",
            "elements": [
              {
                "type": "Identifier",
                "name": "d",
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
                "name": "e",
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
            "loc": {
              "start": {
                "line": 1,
                "column": 14
              },
              "end": {
                "line": 1,
                "column": 20
              }
            }
          },
          {
            "type": "RestElement",
            "argument": {
              "type": "Identifier",
              "name": "f",
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
            }
          }
        ],
        "body": {
          "type": "BlockStatement",
          "body": [],
          "loc": {
            "start": {
              "line": 1,
              "column": 31
            },
            "end": {
              "line": 1,
              "column": 33
            }
          }
        },
        "generator": false,
        "expression": false,
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

// ES6: SpreadElement

func TestHarmony175(t *testing.T) {
	ast, err := Compile("[...a] = b")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "AssignmentExpression",
        "operator": "=",
        "left": {
          "type": "ArrayPattern",
          "elements": [
            {
              "type": "RestElement",
              "argument": {
                "type": "Identifier",
                "name": "a",
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
                  "column": 1
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
              "column": 6
            }
          }
        },
        "right": {
          "type": "Identifier",
          "name": "b",
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

func TestHarmony176(t *testing.T) {
	ast, err := Compile("[a, ...b] = c")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "AssignmentExpression",
        "operator": "=",
        "left": {
          "type": "ArrayPattern",
          "elements": [
            {
              "type": "Identifier",
              "name": "a",
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
            {
              "type": "RestElement",
              "argument": {
                "type": "Identifier",
                "name": "b",
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
                  "column": 4
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
              "column": 9
            }
          }
        },
        "right": {
          "type": "Identifier",
          "name": "c",
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

func TestHarmony177(t *testing.T) {
	ast, err := Compile("[{ a, b }, ...c] = d")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "AssignmentExpression",
        "operator": "=",
        "left": {
          "type": "ArrayPattern",
          "elements": [
            {
              "type": "ObjectPattern",
              "properties": [
                {
                  "type": "Property",
                  "key": {
                    "type": "Identifier",
                    "name": "a",
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
                  },
                  "value": {
                    "type": "Identifier",
                    "name": "a",
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
                  },
                  "kind": "init",
                  "method": false,
                  "shorthand": true,
                  "computed": false,
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
                },
                {
                  "type": "Property",
                  "key": {
                    "type": "Identifier",
                    "name": "b",
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
                    "type": "Identifier",
                    "name": "b",
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
                  "kind": "init",
                  "method": false,
                  "shorthand": true,
                  "computed": false,
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
                }
              ],
              "loc": {
                "start": {
                  "line": 1,
                  "column": 1
                },
                "end": {
                  "line": 1,
                  "column": 9
                }
              }
            },
            {
              "type": "RestElement",
              "argument": {
                "type": "Identifier",
                "name": "c",
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
              "loc": {
                "start": {
                  "line": 1,
                  "column": 11
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
              "column": 16
            }
          }
        },
        "right": {
          "type": "Identifier",
          "name": "d",
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

func TestHarmony178(t *testing.T) {
	ast, err := Compile("[a, ...[b, c]] = d")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "AssignmentExpression",
        "operator": "=",
        "left": {
          "type": "ArrayPattern",
          "elements": [
            {
              "type": "Identifier",
              "name": "a",
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
            {
              "type": "RestElement",
              "argument": {
                "type": "ArrayPattern",
                "elements": [
                  {
                    "type": "Identifier",
                    "name": "b",
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
                  {
                    "type": "Identifier",
                    "name": "c",
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
                  }
                ],
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
        },
        "right": {
          "type": "Identifier",
          "name": "d",
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

func TestHarmony179(t *testing.T) {
	ast, err := Compile("var [...a] = b")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "VariableDeclaration",
      "declarations": [
        {
          "type": "VariableDeclarator",
          "id": {
            "type": "ArrayPattern",
            "elements": [
              {
                "type": "RestElement",
                "argument": {
                  "type": "Identifier",
                  "name": "a",
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
                    "column": 5
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
                "column": 4
              },
              "end": {
                "line": 1,
                "column": 10
              }
            }
          },
          "init": {
            "type": "Identifier",
            "name": "b",
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
              "column": 4
            },
            "end": {
              "line": 1,
              "column": 14
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

func TestHarmony180(t *testing.T) {
	ast, err := Compile("var [a, ...b] = c")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "VariableDeclaration",
      "declarations": [
        {
          "type": "VariableDeclarator",
          "id": {
            "type": "ArrayPattern",
            "elements": [
              {
                "type": "Identifier",
                "name": "a",
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
              {
                "type": "RestElement",
                "argument": {
                  "type": "Identifier",
                  "name": "b",
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
                    "column": 8
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
                "column": 4
              },
              "end": {
                "line": 1,
                "column": 13
              }
            }
          },
          "init": {
            "type": "Identifier",
            "name": "c",
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

func TestHarmony181(t *testing.T) {
	ast, err := Compile("var [{ a, b }, ...c] = d")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "VariableDeclaration",
      "declarations": [
        {
          "type": "VariableDeclarator",
          "id": {
            "type": "ArrayPattern",
            "elements": [
              {
                "type": "ObjectPattern",
                "properties": [
                  {
                    "type": "Property",
                    "key": {
                      "type": "Identifier",
                      "name": "a",
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
                    "value": {
                      "type": "Identifier",
                      "name": "a",
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
                    "kind": "init",
                    "method": false,
                    "shorthand": true,
                    "computed": false,
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
                  {
                    "type": "Property",
                    "key": {
                      "type": "Identifier",
                      "name": "b",
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
                    "value": {
                      "type": "Identifier",
                      "name": "b",
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
                    "kind": "init",
                    "method": false,
                    "shorthand": true,
                    "computed": false,
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
                  }
                ],
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 5
                  },
                  "end": {
                    "line": 1,
                    "column": 13
                  }
                }
              },
              {
                "type": "RestElement",
                "argument": {
                  "type": "Identifier",
                  "name": "c",
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
                    "column": 15
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
                "column": 4
              },
              "end": {
                "line": 1,
                "column": 20
              }
            }
          },
          "init": {
            "type": "Identifier",
            "name": "d",
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
              "column": 4
            },
            "end": {
              "line": 1,
              "column": 24
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

func TestHarmony182(t *testing.T) {
	ast, err := Compile("var [a, ...[b, c]] = d")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "VariableDeclaration",
      "declarations": [
        {
          "type": "VariableDeclarator",
          "id": {
            "type": "ArrayPattern",
            "elements": [
              {
                "type": "Identifier",
                "name": "a",
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
              {
                "type": "RestElement",
                "argument": {
                  "type": "ArrayPattern",
                  "elements": [
                    {
                      "type": "Identifier",
                      "name": "b",
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
                    {
                      "type": "Identifier",
                      "name": "c",
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
          "init": {
            "type": "Identifier",
            "name": "d",
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
              "column": 4
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

func TestHarmony183(t *testing.T) {
	ast, err := Compile("func(...a)")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "CallExpression",
        "callee": {
          "type": "Identifier",
          "name": "func",
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
        "arguments": [
          {
            "type": "SpreadElement",
            "argument": {
              "type": "Identifier",
              "name": "a",
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
                "column": 5
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

func TestHarmony184(t *testing.T) {
	ast, err := Compile("func(a, ...b)")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "CallExpression",
        "callee": {
          "type": "Identifier",
          "name": "func",
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
        "arguments": [
          {
            "type": "Identifier",
            "name": "a",
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
          {
            "type": "SpreadElement",
            "argument": {
              "type": "Identifier",
              "name": "b",
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
                "column": 8
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

func TestHarmony185(t *testing.T) {
	ast, err := Compile("func(...a, b)")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 13
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
          "column": 13
        }
      },
      "expression": {
        "type": "CallExpression",
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 13
          }
        },
        "callee": {
          "type": "Identifier",
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
          "name": "func"
        },
        "arguments": [
          {
            "type": "SpreadElement",
            "loc": {
              "start": {
                "line": 1,
                "column": 5
              },
              "end": {
                "line": 1,
                "column": 9
              }
            },
            "argument": {
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
              "name": "a"
            }
          },
          {
            "type": "Identifier",
            "loc": {
              "start": {
                "line": 1,
                "column": 11
              },
              "end": {
                "line": 1,
                "column": 12
              }
            },
            "name": "b"
          }
        ]
      }
    }
  ]
}
	`, ast)
}

func TestHarmony186(t *testing.T) {
	ast, err := Compile("/[a-z]/u")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "Literal",
        "regexp": {
          "pattern": "[a-z]",
          "flags": "u"
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
    }
  ]
}
	`, ast)
}

func TestHarmony187(t *testing.T) {
	ast, err := Compile("/[\\uD834\\uDF06-\\uD834\\uDF08a-z]/u")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "Literal",
        "regexp": {
          "pattern": "[\\uD834\\uDF06-\\uD834\\uDF08a-z]",
          "flags": "u"
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
    }
  ]
}
	`, ast)
}

func TestHarmony188(t *testing.T) {
	ast, err := Compile("do {} while (false) foo();")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 26,
  "body": [
    {
      "type": "DoWhileStatement",
      "start": 0,
      "end": 19,
      "body": {
        "type": "BlockStatement",
        "start": 3,
        "end": 5,
        "body": []
      },
      "test": {
        "type": "Literal",
        "start": 13,
        "end": 18,
        "value": false
      }
    },
    {
      "type": "ExpressionStatement",
      "start": 20,
      "end": 26,
      "expression": {
        "type": "CallExpression",
        "start": 20,
        "end": 25,
        "callee": {
          "type": "Identifier",
          "start": 20,
          "end": 23,
          "name": "foo"
        },
        "arguments": []
      }
    }
  ]
}
	`, ast)
}

func TestHarmony189(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_STRICT)
	ast, err := CompileWithOpts("let + 1", opts)
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "BinaryExpression",
        "left": {
          "type": "Identifier",
          "name": "let"
        },
        "operator": "+",
        "right": {
          "type": "Literal",
          "value": 1,
          "raw": "1"
        }
      }
    }
  ]
}
	`, ast)
}

func TestHarmony190(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_STRICT)
	ast, err := CompileWithOpts("var let = 1", opts)
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
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
            "name": "let"
          },
          "init": {
            "type": "Literal",
            "value": 1,
            "raw": "1"
          }
        }
      ],
      "kind": "var"
    }
  ]
}
	`, ast)
}

func TestHarmony191(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_STRICT)
	ast, err := CompileWithOpts("e => yield* 10", opts)
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "ArrowFunctionExpression",
        "id": null,
        "params": [
          {
            "type": "Identifier",
            "name": "e",
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
        "body": {
          "type": "BinaryExpression",
          "operator": "*",
          "left": {
            "type": "Identifier",
            "name": "yield",
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
            "type": "Literal",
            "value": 10,
            "raw": "10",
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
              "column": 5
            },
            "end": {
              "line": 1,
              "column": 14
            }
          }
        },
        "generator": false,
        "expression": true,
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

func TestHarmony192(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_STRICT)
	ast, err := CompileWithOpts("(function () { yield* 10 })", opts)
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
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
              "type": "ExpressionStatement",
              "expression": {
                "type": "BinaryExpression",
                "operator": "*",
                "left": {
                  "type": "Identifier",
                  "name": "yield",
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
                "right": {
                  "type": "Literal",
                  "value": 10,
                  "raw": "10",
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
        "generator": false,
        "expression": false,
        "loc": {
          "start": {
            "line": 1,
            "column": 1
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

func TestHarmony193(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_STRICT)
	ast, err := CompileWithOpts("if (1) let\n{}", opts)
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "IfStatement",
      "test": {
        "type": "Literal",
        "value": 1,
        "raw": "1"
      },
      "consequent": {
        "type": "ExpressionStatement",
        "expression": {
          "type": "Identifier",
          "name": "let"
        }
      },
      "alternate": null
    },
    {
      "type": "BlockStatement",
      "body": []
    }
  ]
}
	`, ast)
}

func TestHarmony194(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_STRICT)
	ast, err := CompileWithOpts("var yield = 2", opts)
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
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
            "name": "yield"
          },
          "init": {
            "type": "Literal",
            "value": 2,
            "raw": "2"
          }
        }
      ],
      "kind": "var"
    }
  ]
}
	`, ast)
}

func TestHarmony195(t *testing.T) {
	ast, err := Compile("[...{ a }] = b")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "AssignmentExpression",
        "operator": "=",
        "left": {
          "type": "ArrayPattern",
          "elements": [
            {
              "type": "RestElement",
              "argument": {
                "type": "ObjectPattern",
                "properties": [
                  {
                    "type": "Property",
                    "key": {
                      "type": "Identifier",
                      "name": "a"
                    },
                    "computed": false,
                    "value": {
                      "type": "Identifier",
                      "name": "a"
                    },
                    "kind": "init",
                    "method": false,
                    "shorthand": true
                  }
                ]
              }
            }
          ]
        },
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

/* Regression tests */

func TestHarmony196(t *testing.T) {
	ast, err := Compile("doSth(`${x} + ${y} = ${x + y}`)")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "CallExpression",
        "callee": {
          "type": "Identifier",
          "name": "doSth"
        },
        "arguments": [
          {
            "type": "TemplateLiteral",
            "quasis": [
              {
                "type": "TemplateElement",
                "value": {
                  "raw": "",
                  "cooked": ""
                },
                "tail": false
              },
              {
                "type": "TemplateElement",
                "value": {
                  "raw": " + ",
                  "cooked": " + "
                },
                "tail": false
              },
              {
                "type": "TemplateElement",
                "value": {
                  "raw": " = ",
                  "cooked": " = "
                },
                "tail": false
              },
              {
                "type": "TemplateElement",
                "value": {
                  "raw": "",
                  "cooked": ""
                },
                "tail": true
              }
            ],
            "expressions": [
              {
                "type": "Identifier",
                "name": "x"
              },
              {
                "type": "Identifier",
                "name": "y"
              },
              {
                "type": "BinaryExpression",
                "operator": "+",
                "left": {
                  "type": "Identifier",
                  "name": "x"
                },
                "right": {
                  "type": "Identifier",
                  "name": "y"
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

func TestHarmony197(t *testing.T) {
	ast, err := Compile("function normal(x, y = 10) {}")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "FunctionDeclaration",
      "id": {
        "type": "Identifier",
        "name": "normal"
      },
      "params": [
        {
          "type": "Identifier",
          "name": "x"
        },
        {
          "type": "AssignmentPattern",
          "left": {
            "type": "Identifier",
            "name": "y"
          },
          "right": {
            "type": "Literal",
            "value": 10,
            "raw": "10"
          }
        }
      ],
      "generator": false,
      "body": {
        "type": "BlockStatement",
        "body": []
      }
    }
  ]
}
	`, ast)
}

func TestHarmony198(t *testing.T) {
	ast, err := Compile("() => 42")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "ArrowFunctionExpression",
        "expression": true
      }
    }
  ]
}
	`, ast)
}

func TestHarmony199(t *testing.T) {
	ast, err := Compile("import foo, * as bar from 'baz';")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ImportDeclaration",
      "specifiers": [
        {
          "type": "ImportDefaultSpecifier",
          "local": {
            "type": "Identifier",
            "name": "foo"
          }
        },
        {
          "type": "ImportNamespaceSpecifier",
          "local": {
            "type": "Identifier",
            "name": "bar"
          }
        }
      ],
      "source": {
        "type": "Literal",
        "value": "baz",
        "raw": "'baz'"
      }
    }
  ]
}
	`, ast)
}

func TestHarmony200(t *testing.T) {
	ast, err := Compile("`{${x}}`, `}`")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "SequenceExpression",
        "expressions": [
          {
            "type": "TemplateLiteral",
            "expressions": [
              {
                "type": "Identifier",
                "name": "x"
              }
            ],
            "quasis": [
              {
                "type": "TemplateElement",
                "value": {
                  "cooked": "{",
                  "raw": "{"
                },
                "tail": false
              },
              {
                "type": "TemplateElement",
                "value": {
                  "cooked": "}",
                  "raw": "}"
                },
                "tail": true
              }
            ]
          },
          {
            "type": "TemplateLiteral",
            "expressions": [],
            "quasis": [
              {
                "type": "TemplateElement",
                "value": {
                  "cooked": "}",
                  "raw": "}"
                },
                "tail": true
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

func TestHarmony201(t *testing.T) {
	ast, err := Compile("var {get} = obj;")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "VariableDeclaration",
      "declarations": [
        {
          "type": "VariableDeclarator",
          "id": {
            "type": "ObjectPattern",
            "properties": [
              {
                "type": "Property",
                "method": false,
                "shorthand": true,
                "computed": false,
                "key": {
                  "type": "Identifier",
                  "name": "get"
                },
                "kind": "init",
                "value": {
                  "type": "Identifier",
                  "name": "get"
                }
              }
            ]
          },
          "init": {
            "type": "Identifier",
            "name": "obj"
          }
        }
      ],
      "kind": "var"
    }
  ]
}
	`, ast)
}

func TestHarmony202(t *testing.T) {
	ast, err := Compile("var {propName: localVar = defaultValue} = obj")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 45,
  "body": [
    {
      "type": "VariableDeclaration",
      "start": 0,
      "end": 45,
      "declarations": [
        {
          "type": "VariableDeclarator",
          "start": 4,
          "end": 45,
          "id": {
            "type": "ObjectPattern",
            "start": 4,
            "end": 39,
            "properties": [
              {
                "type": "Property",
                "start": 5,
                "end": 38,
                "method": false,
                "shorthand": false,
                "computed": false,
                "key": {
                  "type": "Identifier",
                  "start": 5,
                  "end": 13,
                  "name": "propName"
                },
                "value": {
                  "type": "AssignmentPattern",
                  "start": 15,
                  "end": 38,
                  "left": {
                    "type": "Identifier",
                    "start": 15,
                    "end": 23,
                    "name": "localVar"
                  },
                  "right": {
                    "type": "Identifier",
                    "start": 26,
                    "end": 38,
                    "name": "defaultValue"
                  }
                },
                "kind": "init"
              }
            ]
          },
          "init": {
            "type": "Identifier",
            "start": 42,
            "end": 45,
            "name": "obj"
          }
        }
      ],
      "kind": "var"
    }
  ]
}
	`, ast)
}

func TestHarmony203(t *testing.T) {
	ast, err := Compile("var {propName = defaultValue} = obj")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 35,
  "body": [
    {
      "type": "VariableDeclaration",
      "start": 0,
      "end": 35,
      "declarations": [
        {
          "type": "VariableDeclarator",
          "start": 4,
          "end": 35,
          "id": {
            "type": "ObjectPattern",
            "start": 4,
            "end": 29,
            "properties": [
              {
                "type": "Property",
                "start": 5,
                "end": 28,
                "method": false,
                "shorthand": true,
                "computed": false,
                "key": {
                  "type": "Identifier",
                  "start": 5,
                  "end": 13,
                  "name": "propName"
                },
                "kind": "init",
                "value": {
                  "type": "AssignmentPattern",
                  "start": 5,
                  "end": 28,
                  "left": {
                    "type": "Identifier",
                    "start": 5,
                    "end": 13,
                    "name": "propName"
                  },
                  "right": {
                    "type": "Identifier",
                    "start": 16,
                    "end": 28,
                    "name": "defaultValue"
                  }
                }
              }
            ]
          },
          "init": {
            "type": "Identifier",
            "start": 32,
            "end": 35,
            "name": "obj"
          }
        }
      ],
      "kind": "var"
    }
  ]
}
	`, ast)
}

func TestHarmony204(t *testing.T) {
	ast, err := Compile("var {get = defaultValue} = obj")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 30,
  "body": [
    {
      "type": "VariableDeclaration",
      "start": 0,
      "end": 30,
      "declarations": [
        {
          "type": "VariableDeclarator",
          "start": 4,
          "end": 30,
          "id": {
            "type": "ObjectPattern",
            "start": 4,
            "end": 24,
            "properties": [
              {
                "type": "Property",
                "start": 5,
                "end": 23,
                "method": false,
                "shorthand": true,
                "computed": false,
                "key": {
                  "type": "Identifier",
                  "start": 5,
                  "end": 8,
                  "name": "get"
                },
                "kind": "init",
                "value": {
                  "type": "AssignmentPattern",
                  "start": 5,
                  "end": 23,
                  "left": {
                    "type": "Identifier",
                    "start": 5,
                    "end": 8,
                    "name": "get"
                  },
                  "right": {
                    "type": "Identifier",
                    "start": 11,
                    "end": 23,
                    "name": "defaultValue"
                  }
                }
              }
            ]
          },
          "init": {
            "type": "Identifier",
            "start": 27,
            "end": 30,
            "name": "obj"
          }
        }
      ],
      "kind": "var"
    }
  ]
}
	`, ast)
}

func TestHarmony205(t *testing.T) {
	ast, err := Compile("var [localVar = defaultValue] = obj")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 35,
  "body": [
    {
      "type": "VariableDeclaration",
      "start": 0,
      "end": 35,
      "declarations": [
        {
          "type": "VariableDeclarator",
          "start": 4,
          "end": 35,
          "id": {
            "type": "ArrayPattern",
            "start": 4,
            "end": 29,
            "elements": [
              {
                "type": "AssignmentPattern",
                "start": 5,
                "end": 28,
                "left": {
                  "type": "Identifier",
                  "start": 5,
                  "end": 13,
                  "name": "localVar"
                },
                "right": {
                  "type": "Identifier",
                  "start": 16,
                  "end": 28,
                  "name": "defaultValue"
                }
              }
            ]
          },
          "init": {
            "type": "Identifier",
            "start": 32,
            "end": 35,
            "name": "obj"
          }
        }
      ],
      "kind": "var"
    }
  ]
}
	`, ast)
}

func TestHarmony206(t *testing.T) {
	ast, err := Compile("({x = 0} = obj)")
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
        "type": "AssignmentExpression",
        "start": 1,
        "end": 14,
        "operator": "=",
        "left": {
          "type": "ObjectPattern",
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
                "end": 3,
                "name": "x"
              },
              "kind": "init",
              "value": {
                "type": "AssignmentPattern",
                "start": 2,
                "end": 7,
                "left": {
                  "type": "Identifier",
                  "start": 2,
                  "end": 3,
                  "name": "x"
                },
                "right": {
                  "type": "Literal",
                  "start": 6,
                  "end": 7,
                  "value": 0,
                  "raw": "0"
                }
              }
            }
          ]
        },
        "right": {
          "type": "Identifier",
          "start": 11,
          "end": 14,
          "name": "obj"
        }
      }
    }
  ]
}
	`, ast)
}

func TestHarmony207(t *testing.T) {
	ast, err := Compile("({x = 0}) => x")
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
        "type": "ArrowFunctionExpression",
        "start": 0,
        "end": 14,
        "id": null,
        "expression": true,
        "generator": false,
        "async": false,
        "params": [
          {
            "type": "ObjectPattern",
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
                  "end": 3,
                  "name": "x"
                },
                "kind": "init",
                "value": {
                  "type": "AssignmentPattern",
                  "start": 2,
                  "end": 7,
                  "left": {
                    "type": "Identifier",
                    "start": 2,
                    "end": 3,
                    "name": "x"
                  },
                  "right": {
                    "type": "Literal",
                    "start": 6,
                    "end": 7,
                    "value": 0,
                    "raw": "0"
                  }
                }
              }
            ]
          }
        ],
        "body": {
          "type": "Identifier",
          "start": 13,
          "end": 14,
          "name": "x"
        }
      }
    }
  ]
}
	`, ast)
}

func TestHarmony208(t *testing.T) {
	ast, err := Compile("[a, {b: {c = 1}}] = arr")
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
        "type": "AssignmentExpression",
        "start": 0,
        "end": 23,
        "operator": "=",
        "left": {
          "type": "ArrayPattern",
          "start": 0,
          "end": 17,
          "elements": [
            {
              "type": "Identifier",
              "start": 1,
              "end": 2,
              "name": "a"
            },
            {
              "type": "ObjectPattern",
              "start": 4,
              "end": 16,
              "properties": [
                {
                  "type": "Property",
                  "start": 5,
                  "end": 15,
                  "method": false,
                  "shorthand": false,
                  "computed": false,
                  "key": {
                    "type": "Identifier",
                    "start": 5,
                    "end": 6,
                    "name": "b"
                  },
                  "value": {
                    "type": "ObjectPattern",
                    "start": 8,
                    "end": 15,
                    "properties": [
                      {
                        "type": "Property",
                        "start": 9,
                        "end": 14,
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
                          "type": "AssignmentPattern",
                          "start": 9,
                          "end": 14,
                          "left": {
                            "type": "Identifier",
                            "start": 9,
                            "end": 10,
                            "name": "c"
                          },
                          "right": {
                            "type": "Literal",
                            "start": 13,
                            "end": 14,
                            "value": 1,
                            "raw": "1"
                          }
                        }
                      }
                    ]
                  },
                  "kind": "init"
                }
              ]
            }
          ]
        },
        "right": {
          "type": "Identifier",
          "start": 20,
          "end": 23,
          "name": "arr"
        }
      }
    }
  ]
}
	`, ast)
}

func TestHarmony209(t *testing.T) {
	ast, err := Compile("for ({x = 0} in arr);")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 21,
  "body": [
    {
      "type": "ForInStatement",
      "start": 0,
      "end": 21,
      "left": {
        "type": "ObjectPattern",
        "start": 5,
        "end": 12,
        "properties": [
          {
            "type": "Property",
            "start": 6,
            "end": 11,
            "method": false,
            "shorthand": true,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 6,
              "end": 7,
              "name": "x"
            },
            "kind": "init",
            "value": {
              "type": "AssignmentPattern",
              "start": 6,
              "end": 11,
              "left": {
                "type": "Identifier",
                "start": 6,
                "end": 7,
                "name": "x"
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
        "type": "Identifier",
        "start": 16,
        "end": 19,
        "name": "arr"
      },
      "body": {
        "type": "EmptyStatement",
        "start": 20,
        "end": 21
      }
    }
  ]
}
	`, ast)
}

func TestHarmony210(t *testing.T) {
	ast, err := Compile("try {} catch ({message}) {}")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 27,
  "body": [
    {
      "type": "TryStatement",
      "start": 0,
      "end": 27,
      "block": {
        "type": "BlockStatement",
        "start": 4,
        "end": 6,
        "body": []
      },
      "handler": {
        "type": "CatchClause",
        "start": 7,
        "end": 27,
        "param": {
          "type": "ObjectPattern",
          "start": 14,
          "end": 23,
          "properties": [
            {
              "type": "Property",
              "start": 15,
              "end": 22,
              "method": false,
              "shorthand": true,
              "computed": false,
              "key": {
                "type": "Identifier",
                "start": 15,
                "end": 22,
                "name": "message"
              },
              "kind": "init",
              "value": {
                "type": "Identifier",
                "start": 15,
                "end": 22,
                "name": "message"
              }
            }
          ]
        },
        "body": {
          "type": "BlockStatement",
          "start": 25,
          "end": 27,
          "body": []
        }
      },
      "finalizer": null
    }
  ]
}
	`, ast)
}

func TestHarmony211(t *testing.T) {
	ast, err := Compile("class A { static() {} }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 23,
  "body": [
    {
      "type": "ClassDeclaration",
      "start": 0,
      "end": 23,
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
        "end": 23,
        "body": [
          {
            "type": "MethodDefinition",
            "start": 10,
            "end": 21,
            "static": false,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 10,
              "end": 16,
              "name": "static"
            },
            "kind": "method",
            "value": {
              "type": "FunctionExpression",
              "start": 16,
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
            }
          }
        ]
      }
    }
  ]
}
	`, ast)
}

func TestHarmony212(t *testing.T) {
	ast, err := Compile("for (const x of list) process(x);")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 33,
  "body": [
    {
      "type": "ForOfStatement",
      "start": 0,
      "end": 33,
      "await": false,
      "left": {
        "type": "VariableDeclaration",
        "start": 5,
        "end": 12,
        "declarations": [
          {
            "type": "VariableDeclarator",
            "start": 11,
            "end": 12,
            "id": {
              "type": "Identifier",
              "start": 11,
              "end": 12,
              "name": "x"
            },
            "init": null
          }
        ],
        "kind": "const"
      },
      "right": {
        "type": "Identifier",
        "start": 16,
        "end": 20,
        "name": "list"
      },
      "body": {
        "type": "ExpressionStatement",
        "start": 22,
        "end": 33,
        "expression": {
          "type": "CallExpression",
          "start": 22,
          "end": 32,
          "callee": {
            "type": "Identifier",
            "start": 22,
            "end": 29,
            "name": "process"
          },
          "arguments": [
            {
              "type": "Identifier",
              "start": 30,
              "end": 31,
              "name": "x"
            }
          ],
          "optional": false
        }
      }
    }
  ]
}
	`, ast)
}

func TestHarmony213(t *testing.T) {
	ast, err := Compile("class A { *static() {} }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
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
        "name": "A"
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
              "start": 11,
              "end": 17,
              "name": "static"
            },
            "kind": "method",
            "value": {
              "type": "FunctionExpression",
              "start": 17,
              "end": 22,
              "id": null,
              "expression": false,
              "generator": true,
              "async": false,
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

func TestHarmony214(t *testing.T) {
	ast, err := Compile("`${/\\d/.exec('1')[0]}`")
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
        "type": "TemplateLiteral",
        "start": 0,
        "end": 22,
        "expressions": [
          {
            "type": "MemberExpression",
            "start": 3,
            "end": 20,
            "object": {
              "type": "CallExpression",
              "start": 3,
              "end": 17,
              "callee": {
                "type": "MemberExpression",
                "start": 3,
                "end": 12,
                "object": {
                  "type": "Literal",
                  "start": 3,
                  "end": 7,
                  "value": null,
                  "regexp": {
                    "pattern": "\\d",
                    "flags": ""
                  }
                },
                "property": {
                  "type": "Identifier",
                  "start": 8,
                  "end": 12,
                  "name": "exec"
                },
                "computed": false,
                "optional": false
              },
              "arguments": [
                {
                  "type": "Literal",
                  "start": 13,
                  "end": 16,
                  "value": "1",
                  "raw": "'1'"
                }
              ],
              "optional": false
            },
            "property": {
              "type": "Literal",
              "start": 18,
              "end": 19,
              "value": 0,
              "raw": "0"
            },
            "computed": true,
            "optional": false
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
            "start": 21,
            "end": 21,
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

func TestHarmony215(t *testing.T) {
	ast, err := Compile("var _𐒦 = 10;")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 15,
  "body": [
    {
      "type": "VariableDeclaration",
      "start": 0,
      "end": 15,
      "declarations": [
        {
          "type": "VariableDeclarator",
          "start": 4,
          "end": 14,
          "id": {
            "type": "Identifier",
            "start": 4,
            "end": 9,
            "name": "_𐒦"
          },
          "init": {
            "type": "Literal",
            "start": 12,
            "end": 14,
            "value": 10,
            "raw": "10"
          }
        }
      ],
      "kind": "var"
    }
  ]
}
	`, ast)
}

func TestHarmony216(t *testing.T) {
	ast, err := Compile("var 𫠝_ = 10;")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 15,
  "body": [
    {
      "type": "VariableDeclaration",
      "start": 0,
      "end": 15,
      "declarations": [
        {
          "type": "VariableDeclarator",
          "start": 4,
          "end": 14,
          "id": {
            "type": "Identifier",
            "start": 4,
            "end": 9,
            "name": "𫠝_"
          },
          "init": {
            "type": "Literal",
            "start": 12,
            "end": 14,
            "value": 10,
            "raw": "10"
          }
        }
      ],
      "kind": "var"
    }
  ]
}
	`, ast)
}

func TestHarmony217(t *testing.T) {
	ast, err := Compile("var _\\u{104A6} = 10;")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 20,
  "body": [
    {
      "type": "VariableDeclaration",
      "start": 0,
      "end": 20,
      "declarations": [
        {
          "type": "VariableDeclarator",
          "start": 4,
          "end": 19,
          "id": {
            "type": "Identifier",
            "start": 4,
            "end": 14,
            "name": "_𐒦"
          },
          "init": {
            "type": "Literal",
            "start": 17,
            "end": 19,
            "value": 10,
            "raw": "10"
          }
        }
      ],
      "kind": "var"
    }
  ]
}
	`, ast)
}

func TestHarmony218(t *testing.T) {
	ast, err := Compile("let [x,] = [1]")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "start": 0,
  "body": [
    {
      "start": 0,
      "declarations": [
        {
          "start": 4,
          "id": {
            "start": 4,
            "elements": [
              {
                "start": 5,
                "name": "x",
                "type": "Identifier",
                "end": 6
              }
            ],
            "type": "ArrayPattern",
            "end": 8
          },
          "init": {
            "start": 11,
            "elements": [
              {
                "start": 12,
                "value": 1,
                "raw": "1",
                "type": "Literal",
                "end": 13
              }
            ],
            "type": "ArrayExpression",
            "end": 14
          },
          "type": "VariableDeclarator",
          "end": 14
        }
      ],
      "kind": "let",
      "type": "VariableDeclaration",
      "end": 14
    }
  ],
  "type": "Program",
  "end": 14
}
	`, ast)
}

func TestHarmony219(t *testing.T) {
	ast, err := Compile("let {x} = y")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "start": 0,
  "body": [
    {
      "start": 0,
      "declarations": [
        {
          "start": 4,
          "id": {
            "start": 4,
            "properties": [
              {
                "start": 5,
                "method": false,
                "shorthand": true,
                "computed": false,
                "key": {
                  "start": 5,
                  "name": "x",
                  "type": "Identifier",
                  "end": 6
                },
                "kind": "init",
                "value": {
                  "start": 5,
                  "name": "x",
                  "type": "Identifier",
                  "end": 6
                },
                "type": "Property",
                "end": 6
              }
            ],
            "type": "ObjectPattern",
            "end": 7
          },
          "init": {
            "start": 10,
            "name": "y",
            "type": "Identifier",
            "end": 11
          },
          "type": "VariableDeclarator",
          "end": 11
        }
      ],
      "kind": "let",
      "type": "VariableDeclaration",
      "end": 11
    }
  ],
  "type": "Program",
  "end": 11
}
	`, ast)
}

func TestHarmony220(t *testing.T) {
	ast, err := Compile("[x,,] = 1")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "AssignmentExpression",
        "operator": "=",
        "left": {
          "type": "ArrayPattern",
          "elements": [
            {
              "type": "Identifier",
              "name": "x"
            },
            null
          ]
        },
        "right": {
          "type": "Literal",
          "value": 1,
          "raw": "1"
        }
      }
    }
  ]
}
	`, ast)
}

func TestHarmony221(t *testing.T) {
	ast, err := Compile("for (var [name, value] in obj) {}")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 33,
  "body": [
    {
      "type": "ForInStatement",
      "start": 0,
      "end": 33,
      "left": {
        "type": "VariableDeclaration",
        "start": 5,
        "end": 22,
        "declarations": [
          {
            "type": "VariableDeclarator",
            "start": 9,
            "end": 22,
            "id": {
              "type": "ArrayPattern",
              "start": 9,
              "end": 22,
              "elements": [
                {
                  "type": "Identifier",
                  "start": 10,
                  "end": 14,
                  "name": "name"
                },
                {
                  "type": "Identifier",
                  "start": 16,
                  "end": 21,
                  "name": "value"
                }
              ]
            },
            "init": null
          }
        ],
        "kind": "var"
      },
      "right": {
        "type": "Identifier",
        "start": 26,
        "end": 29,
        "name": "obj"
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

func TestHarmony222(t *testing.T) {
	ast, err := Compile("function foo() { new.target; }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 30,
  "body": [
    {
      "type": "FunctionDeclaration",
      "start": 0,
      "end": 30,
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
        "end": 30,
        "body": [
          {
            "type": "ExpressionStatement",
            "start": 17,
            "end": 28,
            "expression": {
              "type": "MetaProperty",
              "start": 17,
              "end": 27,
              "meta": {
                "type": "Identifier",
                "start": 17,
                "end": 20,
                "name": "new"
              },
              "property": {
                "type": "Identifier",
                "start": 21,
                "end": 27,
                "name": "target"
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

func TestHarmony223(t *testing.T) {
	ast, err := Compile("function x() { return () => new.target }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
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
        "end": 10,
        "name": "x"
      },
      "generator": false,
      "async": false,
      "params": [],
      "body": {
        "type": "BlockStatement",
        "start": 13,
        "end": 40,
        "body": [
          {
            "type": "ReturnStatement",
            "start": 15,
            "end": 38,
            "argument": {
              "type": "ArrowFunctionExpression",
              "start": 22,
              "end": 38,
              "id": null,
              "expression": true,
              "generator": false,
              "async": false,
              "params": [],
              "body": {
                "type": "MetaProperty",
                "start": 28,
                "end": 38,
                "meta": {
                  "type": "Identifier",
                  "start": 28,
                  "end": 31,
                  "name": "new"
                },
                "property": {
                  "type": "Identifier",
                  "start": 32,
                  "end": 38,
                  "name": "target"
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

func TestHarmony224(t *testing.T) {
	ast, err := Compile("export default function foo() {} false")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 38,
  "body": [
    {
      "type": "ExportDefaultDeclaration",
      "start": 0,
      "end": 32,
      "declaration": {
        "type": "FunctionDeclaration",
        "start": 15,
        "end": 32,
        "id": {
          "type": "Identifier",
          "start": 24,
          "end": 27,
          "name": "foo"
        },
        "generator": false,
        "async": false,
        "params": [],
        "body": {
          "type": "BlockStatement",
          "start": 30,
          "end": 32,
          "body": []
        }
      }
    },
    {
      "type": "ExpressionStatement",
      "start": 33,
      "end": 38,
      "expression": {
        "type": "Literal",
        "start": 33,
        "end": 38,
        "value": false
      }
    }
  ]
}
	`, ast)
}

func TestHarmony225(t *testing.T) {
	ast, err := Compile("({ ['__proto__']: 1, __proto__: 2 })")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 36,
  "body": [
    {
      "type": "ExpressionStatement",
      "start": 0,
      "end": 36,
      "expression": {
        "type": "ObjectExpression",
        "start": 1,
        "end": 35,
        "properties": [
          {
            "type": "Property",
            "start": 3,
            "end": 19,
            "method": false,
            "shorthand": false,
            "computed": true,
            "key": {
              "type": "Literal",
              "start": 4,
              "end": 15,
              "value": "__proto__",
              "raw": "'__proto__'"
            },
            "value": {
              "type": "Literal",
              "start": 18,
              "end": 19,
              "value": 1,
              "raw": "1"
            },
            "kind": "init"
          },
          {
            "type": "Property",
            "start": 21,
            "end": 33,
            "method": false,
            "shorthand": false,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 21,
              "end": 30,
              "name": "__proto__"
            },
            "value": {
              "type": "Literal",
              "start": 32,
              "end": 33,
              "value": 2,
              "raw": "2"
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

func TestHarmony226(t *testing.T) {
	ast, err := Compile("({ __proto__() { return 1 }, __proto__: 2 })")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
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
        "type": "ObjectExpression",
        "start": 1,
        "end": 43,
        "properties": [
          {
            "type": "Property",
            "start": 3,
            "end": 27,
            "method": true,
            "shorthand": false,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 3,
              "end": 12,
              "name": "__proto__"
            },
            "kind": "init",
            "value": {
              "type": "FunctionExpression",
              "start": 12,
              "end": 27,
              "id": null,
              "expression": false,
              "generator": false,
              "async": false,
              "params": [],
              "body": {
                "type": "BlockStatement",
                "start": 15,
                "end": 27,
                "body": [
                  {
                    "type": "ReturnStatement",
                    "start": 17,
                    "end": 25,
                    "argument": {
                      "type": "Literal",
                      "start": 24,
                      "end": 25,
                      "value": 1,
                      "raw": "1"
                    }
                  }
                ]
              }
            }
          },
          {
            "type": "Property",
            "start": 29,
            "end": 41,
            "method": false,
            "shorthand": false,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 29,
              "end": 38,
              "name": "__proto__"
            },
            "value": {
              "type": "Literal",
              "start": 40,
              "end": 41,
              "value": 2,
              "raw": "2"
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

func TestHarmony227(t *testing.T) {
	ast, err := Compile("({ get __proto__() { return 1 }, __proto__: 2 })")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 48,
  "body": [
    {
      "type": "ExpressionStatement",
      "start": 0,
      "end": 48,
      "expression": {
        "type": "ObjectExpression",
        "start": 1,
        "end": 47,
        "properties": [
          {
            "type": "Property",
            "start": 3,
            "end": 31,
            "method": false,
            "shorthand": false,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 7,
              "end": 16,
              "name": "__proto__"
            },
            "kind": "get",
            "value": {
              "type": "FunctionExpression",
              "start": 16,
              "end": 31,
              "id": null,
              "expression": false,
              "generator": false,
              "async": false,
              "params": [],
              "body": {
                "type": "BlockStatement",
                "start": 19,
                "end": 31,
                "body": [
                  {
                    "type": "ReturnStatement",
                    "start": 21,
                    "end": 29,
                    "argument": {
                      "type": "Literal",
                      "start": 28,
                      "end": 29,
                      "value": 1,
                      "raw": "1"
                    }
                  }
                ]
              }
            }
          },
          {
            "type": "Property",
            "start": 33,
            "end": 45,
            "method": false,
            "shorthand": false,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 33,
              "end": 42,
              "name": "__proto__"
            },
            "value": {
              "type": "Literal",
              "start": 44,
              "end": 45,
              "value": 2,
              "raw": "2"
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

func TestHarmony228(t *testing.T) {
	ast, err := Compile("({ __proto__, __proto__: 2 })")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
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
            "start": 3,
            "end": 12,
            "method": false,
            "shorthand": true,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 3,
              "end": 12,
              "name": "__proto__"
            },
            "kind": "init",
            "value": {
              "type": "Identifier",
              "start": 3,
              "end": 12,
              "name": "__proto__"
            }
          },
          {
            "type": "Property",
            "start": 14,
            "end": 26,
            "method": false,
            "shorthand": false,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 14,
              "end": 23,
              "name": "__proto__"
            },
            "value": {
              "type": "Literal",
              "start": 25,
              "end": 26,
              "value": 2,
              "raw": "2"
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

func TestHarmony229(t *testing.T) {
	ast, err := Compile("({__proto__: a, __proto__: b} = {})")
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
        "type": "AssignmentExpression",
        "start": 1,
        "end": 34,
        "operator": "=",
        "left": {
          "type": "ObjectPattern",
          "start": 1,
          "end": 29,
          "properties": [
            {
              "type": "Property",
              "start": 2,
              "end": 14,
              "method": false,
              "shorthand": false,
              "computed": false,
              "key": {
                "type": "Identifier",
                "start": 2,
                "end": 11,
                "name": "__proto__"
              },
              "value": {
                "type": "Identifier",
                "start": 13,
                "end": 14,
                "name": "a"
              },
              "kind": "init"
            },
            {
              "type": "Property",
              "start": 16,
              "end": 28,
              "method": false,
              "shorthand": false,
              "computed": false,
              "key": {
                "type": "Identifier",
                "start": 16,
                "end": 25,
                "name": "__proto__"
              },
              "value": {
                "type": "Identifier",
                "start": 27,
                "end": 28,
                "name": "b"
              },
              "kind": "init"
            }
          ]
        },
        "right": {
          "type": "ObjectExpression",
          "start": 32,
          "end": 34,
          "properties": []
        }
      }
    }
  ]
}
	`, ast)
}

func TestHarmony230(t *testing.T) {
	ast, err := Compile("export default /foo/")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 20,
  "body": [
    {
      "type": "ExportDefaultDeclaration",
      "start": 0,
      "end": 20,
      "declaration": {
        "type": "Literal",
        "start": 15,
        "end": 20,
        "value": {},
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

func TestHarmony231(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_STRICT)
	ast, err := CompileWithOpts("l\\u0065t\na", opts)
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "Identifier",
        "name": "let"
      }
    },
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "Identifier",
        "name": "a"
      }
    }
  ]
}
	`, ast)
}

func TestHarmony232(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_MODULE).Off(parser.FEAT_GLOBAL_ASYNC)
	ast, err := CompileWithOpts("var await = 0", opts)
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 13,
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 13
    }
  },
  "body": [
    {
      "type": "VariableDeclaration",
      "start": 0,
      "end": 13,
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 13
        }
      },
      "declarations": [
        {
          "type": "VariableDeclarator",
          "start": 4,
          "end": 13,
          "loc": {
            "start": {
              "line": 1,
              "column": 4
            },
            "end": {
              "line": 1,
              "column": 13
            }
          },
          "id": {
            "type": "Identifier",
            "start": 4,
            "end": 9,
            "loc": {
              "start": {
                "line": 1,
                "column": 4
              },
              "end": {
                "line": 1,
                "column": 9
              }
            },
            "name": "await"
          },
          "init": {
            "type": "Literal",
            "start": 12,
            "end": 13,
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
            "value": 0,
            "raw": "0"
          }
        }
      ],
      "kind": "var"
    }
  ]
}
	`, ast)
}

func TestHarmony233(t *testing.T) {
	ast, err := Compile("/[a-z]/gimuy")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "Literal",
        "regexp": {
          "pattern": "[a-z]",
          "flags": "gimuy"
        }
      }
    }
  ]
}
	`, ast)
}

func TestHarmony234(t *testing.T) {
	ast, err := Compile("/[a-z]/s")
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
        "type": "Literal",
        "start": 0,
        "end": 8,
        "regexp": {
          "pattern": "[a-z]",
          "flags": "s"
        }
      }
    }
  ]
}
	`, ast)
}

func TestHarmony235(t *testing.T) {
	ast, err := Compile("(([,]) => 0)")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "ArrowFunctionExpression",
        "params": [
          {
            "type": "ArrayPattern",
            "elements": [
              null
            ]
          }
        ],
        "body": {
          "type": "Literal",
          "value": 0,
          "raw": "0"
        },
        "expression": true
      }
    }
  ]
}
	`, ast)
}

func TestHarmony236(t *testing.T) {
	ast, err := Compile("function foo() { return {arguments} }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
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
        "end": 12,
        "name": "foo"
      },
      "generator": false,
      "async": false,
      "params": [],
      "body": {
        "type": "BlockStatement",
        "start": 15,
        "end": 37,
        "body": [
          {
            "type": "ReturnStatement",
            "start": 17,
            "end": 35,
            "argument": {
              "type": "ObjectExpression",
              "start": 24,
              "end": 35,
              "properties": [
                {
                  "type": "Property",
                  "start": 25,
                  "end": 34,
                  "method": false,
                  "shorthand": true,
                  "computed": false,
                  "key": {
                    "type": "Identifier",
                    "start": 25,
                    "end": 34,
                    "name": "arguments"
                  },
                  "kind": "init",
                  "value": {
                    "type": "Identifier",
                    "start": 25,
                    "end": 34,
                    "name": "arguments"
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

func TestHarmony237(t *testing.T) {
	ast, err := Compile("function foo() { return {eval} }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 32,
  "body": [
    {
      "type": "FunctionDeclaration",
      "start": 0,
      "end": 32,
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
        "end": 32,
        "body": [
          {
            "type": "ReturnStatement",
            "start": 17,
            "end": 30,
            "argument": {
              "type": "ObjectExpression",
              "start": 24,
              "end": 30,
              "properties": [
                {
                  "type": "Property",
                  "start": 25,
                  "end": 29,
                  "method": false,
                  "shorthand": true,
                  "computed": false,
                  "key": {
                    "type": "Identifier",
                    "start": 25,
                    "end": 29,
                    "name": "eval"
                  },
                  "kind": "init",
                  "value": {
                    "type": "Identifier",
                    "start": 25,
                    "end": 29,
                    "name": "eval"
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

func TestHarmony238(t *testing.T) {
	ast, err := Compile("function foo() { 'use strict'; return {arguments} }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 51,
  "body": [
    {
      "type": "FunctionDeclaration",
      "start": 0,
      "end": 51,
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
        "end": 51,
        "body": [
          {
            "type": "ExpressionStatement",
            "start": 17,
            "end": 30,
            "expression": {
              "type": "Literal",
              "start": 17,
              "end": 29,
              "value": "use strict",
              "raw": "'use strict'"
            }
          },
          {
            "type": "ReturnStatement",
            "start": 31,
            "end": 49,
            "argument": {
              "type": "ObjectExpression",
              "start": 38,
              "end": 49,
              "properties": [
                {
                  "type": "Property",
                  "start": 39,
                  "end": 48,
                  "method": false,
                  "shorthand": true,
                  "computed": false,
                  "key": {
                    "type": "Identifier",
                    "start": 39,
                    "end": 48,
                    "name": "arguments"
                  },
                  "kind": "init",
                  "value": {
                    "type": "Identifier",
                    "start": 39,
                    "end": 48,
                    "name": "arguments"
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

func TestHarmony239(t *testing.T) {
	ast, err := Compile("function foo() { 'use strict'; return {eval} }")
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
        "end": 46,
        "body": [
          {
            "type": "ExpressionStatement",
            "start": 17,
            "end": 30,
            "expression": {
              "type": "Literal",
              "start": 17,
              "end": 29,
              "value": "use strict",
              "raw": "'use strict'"
            }
          },
          {
            "type": "ReturnStatement",
            "start": 31,
            "end": 44,
            "argument": {
              "type": "ObjectExpression",
              "start": 38,
              "end": 44,
              "properties": [
                {
                  "type": "Property",
                  "start": 39,
                  "end": 43,
                  "method": false,
                  "shorthand": true,
                  "computed": false,
                  "key": {
                    "type": "Identifier",
                    "start": 39,
                    "end": 43,
                    "name": "eval"
                  },
                  "kind": "init",
                  "value": {
                    "type": "Identifier",
                    "start": 39,
                    "end": 43,
                    "name": "eval"
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

func TestHarmony240(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_STRICT)
	ast, err := CompileWithOpts("function foo() { return {yield} }", opts)
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 33,
  "body": [
    {
      "type": "FunctionDeclaration",
      "start": 0,
      "end": 33,
      "id": {
        "type": "Identifier",
        "start": 9,
        "end": 12,
        "name": "foo"
      },
      "generator": false,
      "params": [],
      "body": {
        "type": "BlockStatement",
        "start": 15,
        "end": 33,
        "body": [
          {
            "type": "ReturnStatement",
            "start": 17,
            "end": 31,
            "argument": {
              "type": "ObjectExpression",
              "start": 24,
              "end": 31,
              "properties": [
                {
                  "type": "Property",
                  "start": 25,
                  "end": 30,
                  "method": false,
                  "shorthand": true,
                  "computed": false,
                  "key": {
                    "type": "Identifier",
                    "start": 25,
                    "end": 30,
                    "name": "yield"
                  },
                  "kind": "init",
                  "value": {
                    "type": "Identifier",
                    "start": 25,
                    "end": 30,
                    "name": "yield"
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

func TestHarmony241(t *testing.T) {
	ast, err := Compile("function* foo(a = function*(b) { yield b }) { }")
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
        "start": 10,
        "end": 13,
        "name": "foo"
      },
      "generator": true,
      "async": false,
      "params": [
        {
          "type": "AssignmentPattern",
          "start": 14,
          "end": 42,
          "left": {
            "type": "Identifier",
            "start": 14,
            "end": 15,
            "name": "a"
          },
          "right": {
            "type": "FunctionExpression",
            "start": 18,
            "end": 42,
            "id": null,
            "expression": false,
            "generator": true,
            "async": false,
            "params": [
              {
                "type": "Identifier",
                "start": 28,
                "end": 29,
                "name": "b"
              }
            ],
            "body": {
              "type": "BlockStatement",
              "start": 31,
              "end": 42,
              "body": [
                {
                  "type": "ExpressionStatement",
                  "start": 33,
                  "end": 40,
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
      ],
      "body": {
        "type": "BlockStatement",
        "start": 44,
        "end": 47,
        "body": []
      }
    }
  ]
}
	`, ast)
}

// 'yield' as function names.

func TestHarmony242(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_STRICT)
	ast, err := CompileWithOpts("function* yield() {}", opts)
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
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
        "start": 10,
        "end": 15,
        "name": "yield"
      },
      "generator": true,
      "params": [],
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

func TestHarmony243(t *testing.T) {
	ast, err := Compile("({*yield() {}})")
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
              "start": 3,
              "end": 8,
              "name": "yield"
            },
            "kind": "init",
            "value": {
              "type": "FunctionExpression",
              "start": 8,
              "end": 13,
              "id": null,
              "generator": true,
              "expression": false,
              "params": [],
              "body": {
                "type": "BlockStatement",
                "start": 11,
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

func TestHarmony244(t *testing.T) {
	ast, err := Compile("class A {*yield() {}}")
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
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 10,
              "end": 15,
              "name": "yield"
            },
            "static": false,
            "kind": "method",
            "value": {
              "type": "FunctionExpression",
              "start": 15,
              "end": 20,
              "id": null,
              "generator": true,
              "expression": false,
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

func TestHarmony245(t *testing.T) {
	ast, err := Compile("function* wrap() {\n({*yield() {}})\n}")
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
        "start": 10,
        "end": 14,
        "name": "wrap"
      },
      "generator": true,
      "async": false,
      "params": [],
      "body": {
        "type": "BlockStatement",
        "start": 17,
        "end": 36,
        "body": [
          {
            "type": "ExpressionStatement",
            "start": 19,
            "end": 34,
            "expression": {
              "type": "ObjectExpression",
              "start": 20,
              "end": 33,
              "properties": [
                {
                  "type": "Property",
                  "start": 21,
                  "end": 32,
                  "method": true,
                  "shorthand": false,
                  "computed": false,
                  "key": {
                    "type": "Identifier",
                    "start": 22,
                    "end": 27,
                    "name": "yield"
                  },
                  "kind": "init",
                  "value": {
                    "type": "FunctionExpression",
                    "start": 27,
                    "end": 32,
                    "id": null,
                    "expression": false,
                    "generator": true,
                    "async": false,
                    "params": [],
                    "body": {
                      "type": "BlockStatement",
                      "start": 30,
                      "end": 32,
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

func TestHarmony246(t *testing.T) {
	ast, err := Compile("function* wrap() {\nclass A {*yield() {}}\n}")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 42,
  "body": [
    {
      "type": "FunctionDeclaration",
      "start": 0,
      "end": 42,
      "id": {
        "type": "Identifier",
        "start": 10,
        "end": 14,
        "name": "wrap"
      },
      "generator": true,
      "async": false,
      "params": [],
      "body": {
        "type": "BlockStatement",
        "start": 17,
        "end": 42,
        "body": [
          {
            "type": "ClassDeclaration",
            "start": 19,
            "end": 40,
            "id": {
              "type": "Identifier",
              "start": 25,
              "end": 26,
              "name": "A"
            },
            "superClass": null,
            "body": {
              "type": "ClassBody",
              "start": 27,
              "end": 40,
              "body": [
                {
                  "type": "MethodDefinition",
                  "start": 28,
                  "end": 39,
                  "static": false,
                  "computed": false,
                  "key": {
                    "type": "Identifier",
                    "start": 29,
                    "end": 34,
                    "name": "yield"
                  },
                  "kind": "method",
                  "value": {
                    "type": "FunctionExpression",
                    "start": 34,
                    "end": 39,
                    "id": null,
                    "expression": false,
                    "generator": true,
                    "async": false,
                    "params": [],
                    "body": {
                      "type": "BlockStatement",
                      "start": 37,
                      "end": 39,
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

// Allow yield expressions inside functions in default parameters:

func TestHarmony247(t *testing.T) {
	ast, err := Compile("function* foo(a = function* foo() { yield b }) {}")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 49,
  "body": [
    {
      "type": "FunctionDeclaration",
      "start": 0,
      "end": 49,
      "id": {
        "type": "Identifier",
        "start": 10,
        "end": 13,
        "name": "foo"
      },
      "generator": true,
      "params": [
        {
          "type": "AssignmentPattern",
          "start": 14,
          "end": 45,
          "left": {
            "type": "Identifier",
            "start": 14,
            "end": 15,
            "name": "a"
          },
          "right": {
            "type": "FunctionExpression",
            "start": 18,
            "end": 45,
            "id": {
              "type": "Identifier",
              "start": 28,
              "end": 31,
              "name": "foo"
            },
            "generator": true,
            "expression": false,
            "params": [],
            "body": {
              "type": "BlockStatement",
              "start": 34,
              "end": 45,
              "body": [
                {
                  "type": "ExpressionStatement",
                  "start": 36,
                  "end": 43,
                  "expression": {
                    "type": "YieldExpression",
                    "start": 36,
                    "end": 43,
                    "delegate": false,
                    "argument": {
                      "type": "Identifier",
                      "start": 42,
                      "end": 43,
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
        "start": 47,
        "end": 49,
        "body": []
      }
    }
  ]
}
	`, ast)
}

func TestHarmony248(t *testing.T) {
	ast, err := Compile("function* foo(a = {*bar() { yield b }}) {}")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 42,
  "body": [
    {
      "type": "FunctionDeclaration",
      "start": 0,
      "end": 42,
      "id": {
        "type": "Identifier",
        "start": 10,
        "end": 13,
        "name": "foo"
      },
      "generator": true,
      "params": [
        {
          "type": "AssignmentPattern",
          "start": 14,
          "end": 38,
          "left": {
            "type": "Identifier",
            "start": 14,
            "end": 15,
            "name": "a"
          },
          "right": {
            "type": "ObjectExpression",
            "start": 18,
            "end": 38,
            "properties": [
              {
                "type": "Property",
                "start": 19,
                "end": 37,
                "method": true,
                "shorthand": false,
                "computed": false,
                "key": {
                  "type": "Identifier",
                  "start": 20,
                  "end": 23,
                  "name": "bar"
                },
                "kind": "init",
                "value": {
                  "type": "FunctionExpression",
                  "start": 23,
                  "end": 37,
                  "id": null,
                  "generator": true,
                  "expression": false,
                  "params": [],
                  "body": {
                    "type": "BlockStatement",
                    "start": 26,
                    "end": 37,
                    "body": [
                      {
                        "type": "ExpressionStatement",
                        "start": 28,
                        "end": 35,
                        "expression": {
                          "type": "YieldExpression",
                          "start": 28,
                          "end": 35,
                          "delegate": false,
                          "argument": {
                            "type": "Identifier",
                            "start": 34,
                            "end": 35,
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
        "start": 40,
        "end": 42,
        "body": []
      }
    }
  ]
}
	`, ast)
}

func TestHarmony249(t *testing.T) {
	ast, err := Compile("function* foo(a = class {*bar() { yield b }}) {}")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
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
        "start": 10,
        "end": 13,
        "name": "foo"
      },
      "generator": true,
      "params": [
        {
          "type": "AssignmentPattern",
          "start": 14,
          "end": 44,
          "left": {
            "type": "Identifier",
            "start": 14,
            "end": 15,
            "name": "a"
          },
          "right": {
            "type": "ClassExpression",
            "start": 18,
            "end": 44,
            "id": null,
            "superClass": null,
            "body": {
              "type": "ClassBody",
              "start": 24,
              "end": 44,
              "body": [
                {
                  "type": "MethodDefinition",
                  "start": 25,
                  "end": 43,
                  "computed": false,
                  "key": {
                    "type": "Identifier",
                    "start": 26,
                    "end": 29,
                    "name": "bar"
                  },
                  "static": false,
                  "kind": "method",
                  "value": {
                    "type": "FunctionExpression",
                    "start": 29,
                    "end": 43,
                    "id": null,
                    "generator": true,
                    "expression": false,
                    "params": [],
                    "body": {
                      "type": "BlockStatement",
                      "start": 32,
                      "end": 43,
                      "body": [
                        {
                          "type": "ExpressionStatement",
                          "start": 34,
                          "end": 41,
                          "expression": {
                            "type": "YieldExpression",
                            "start": 34,
                            "end": 41,
                            "delegate": false,
                            "argument": {
                              "type": "Identifier",
                              "start": 40,
                              "end": 41,
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
        "start": 46,
        "end": 48,
        "body": []
      }
    }
  ]
}
	`, ast)
}

// Distinguish ParenthesizedExpression or ArrowFunctionExpression

func TestHarmony250(t *testing.T) {
	ast, err := Compile("function* wrap() {\n(a = yield b)\n}")
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
        "start": 10,
        "end": 14,
        "name": "wrap"
      },
      "generator": true,
      "params": [],
      "body": {
        "type": "BlockStatement",
        "start": 17,
        "end": 34,
        "body": [
          {
            "type": "ExpressionStatement",
            "start": 19,
            "end": 32,
            "expression": {
              "type": "AssignmentExpression",
              "start": 20,
              "end": 31,
              "operator": "=",
              "left": {
                "type": "Identifier",
                "start": 20,
                "end": 21,
                "name": "a"
              },
              "right": {
                "type": "YieldExpression",
                "start": 24,
                "end": 31,
                "delegate": false,
                "argument": {
                  "type": "Identifier",
                  "start": 30,
                  "end": 31,
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

func TestHarmony251(t *testing.T) {
	ast, err := Compile("function* wrap() {\n({a = yield b} = obj)\n}")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
 {
  "type": "Program",
  "start": 0,
  "end": 42,
  "body": [
    {
      "type": "FunctionDeclaration",
      "start": 0,
      "end": 42,
      "id": {
        "type": "Identifier",
        "start": 10,
        "end": 14,
        "name": "wrap"
      },
      "params": [],
      "generator": true,
      "body": {
        "type": "BlockStatement",
        "start": 17,
        "end": 42,
        "body": [
          {
            "type": "ExpressionStatement",
            "start": 19,
            "end": 40,
            "expression": {
              "type": "AssignmentExpression",
              "start": 20,
              "end": 39,
              "operator": "=",
              "left": {
                "type": "ObjectPattern",
                "start": 20,
                "end": 33,
                "properties": [
                  {
                    "type": "Property",
                    "start": 21,
                    "end": 32,
                    "method": false,
                    "shorthand": true,
                    "computed": false,
                    "key": {
                      "type": "Identifier",
                      "start": 21,
                      "end": 22,
                      "name": "a"
                    },
                    "kind": "init",
                    "value": {
                      "type": "AssignmentPattern",
                      "start": 21,
                      "end": 32,
                      "left": {
                        "type": "Identifier",
                        "start": 21,
                        "end": 22,
                        "name": "a"
                      },
                      "right": {
                        "type": "YieldExpression",
                        "start": 25,
                        "end": 32,
                        "delegate": false,
                        "argument": {
                          "type": "Identifier",
                          "start": 31,
                          "end": 32,
                          "name": "b"
                        }
                      }
                    }
                  }
                ]
              },
              "right": {
                "type": "Identifier",
                "start": 36,
                "end": 39,
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

func TestHarmony252(t *testing.T) {
	ast, err := Compile("export default class Foo {}++x")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExportDefaultDeclaration",
      "declaration": {
        "type": "ClassDeclaration",
        "id": {
          "type": "Identifier",
          "name": "Foo"
        },
        "superClass": null,
        "body": {
          "type": "ClassBody",
          "body": []
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
          "name": "x"
        }
      }
    }
  ]
}
	`, ast)
}

func TestHarmony253(t *testing.T) {
	ast, err := Compile("function *f() { yield\n{}/1/g\n}")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "FunctionDeclaration",
      "id": {
        "type": "Identifier",
        "name": "f"
      },
      "body": {
        "type": "BlockStatement",
        "body": [
          {
            "type": "ExpressionStatement",
            "expression": {
              "type": "YieldExpression",
              "argument": null,
              "delegate": false
            }
          },
          {
            "type": "BlockStatement",
            "body": []
          },
          {
            "type": "ExpressionStatement",
            "expression": {
              "type": "Literal",
              "regexp": {
                "pattern": "1",
                "flags": "g"
              }
            }
          }
        ]
      },
      "generator": true
    }
  ]
}
	`, ast)
}

func TestHarmony254(t *testing.T) {
	ast, err := Compile("class B extends A { constructor(a = super()) { return a }}")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 58,
  "body": [
    {
      "type": "ClassDeclaration",
      "start": 0,
      "end": 58,
      "id": {
        "type": "Identifier",
        "start": 6,
        "end": 7,
        "name": "B"
      },
      "superClass": {
        "type": "Identifier",
        "start": 16,
        "end": 17,
        "name": "A"
      },
      "body": {
        "type": "ClassBody",
        "start": 18,
        "end": 58,
        "body": [
          {
            "type": "MethodDefinition",
            "start": 20,
            "end": 57,
            "static": false,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 20,
              "end": 31,
              "name": "constructor"
            },
            "kind": "constructor",
            "value": {
              "type": "FunctionExpression",
              "start": 31,
              "end": 57,
              "id": null,
              "expression": false,
              "generator": false,
              "async": false,
              "params": [
                {
                  "type": "AssignmentPattern",
                  "start": 32,
                  "end": 43,
                  "left": {
                    "type": "Identifier",
                    "start": 32,
                    "end": 33,
                    "name": "a"
                  },
                  "right": {
                    "type": "CallExpression",
                    "start": 36,
                    "end": 43,
                    "callee": {
                      "type": "Super",
                      "start": 36,
                      "end": 41
                    },
                    "arguments": [],
                    "optional": false
                  }
                }
              ],
              "body": {
                "type": "BlockStatement",
                "start": 45,
                "end": 57,
                "body": [
                  {
                    "type": "ReturnStatement",
                    "start": 47,
                    "end": 55,
                    "argument": {
                      "type": "Identifier",
                      "start": 54,
                      "end": 55,
                      "name": "a"
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

func TestHarmony255(t *testing.T) {
	ast, err := Compile("class B { foo(a = super.foo()) { return a }}")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
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
        "name": "B"
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
            "end": 43,
            "static": false,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 10,
              "end": 13,
              "name": "foo"
            },
            "kind": "method",
            "value": {
              "type": "FunctionExpression",
              "start": 13,
              "end": 43,
              "id": null,
              "expression": false,
              "generator": false,
              "async": false,
              "params": [
                {
                  "type": "AssignmentPattern",
                  "start": 14,
                  "end": 29,
                  "left": {
                    "type": "Identifier",
                    "start": 14,
                    "end": 15,
                    "name": "a"
                  },
                  "right": {
                    "type": "CallExpression",
                    "start": 18,
                    "end": 29,
                    "callee": {
                      "type": "MemberExpression",
                      "start": 18,
                      "end": 27,
                      "object": {
                        "type": "Super",
                        "start": 18,
                        "end": 23
                      },
                      "property": {
                        "type": "Identifier",
                        "start": 24,
                        "end": 27,
                        "name": "foo"
                      },
                      "computed": false,
                      "optional": false
                    },
                    "arguments": [],
                    "optional": false
                  }
                }
              ],
              "body": {
                "type": "BlockStatement",
                "start": 31,
                "end": 43,
                "body": [
                  {
                    "type": "ReturnStatement",
                    "start": 33,
                    "end": 41,
                    "argument": {
                      "type": "Identifier",
                      "start": 40,
                      "end": 41,
                      "name": "a"
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

func TestHarmony256(t *testing.T) {
	ast, err := Compile("export { x as y } from './y.js';\nexport { x as z } from './z.js';")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 65,
  "body": [
    {
      "type": "ExportNamedDeclaration",
      "start": 0,
      "end": 32,
      "declaration": null,
      "specifiers": [
        {
          "type": "ExportSpecifier",
          "start": 9,
          "end": 15,
          "local": {
            "type": "Identifier",
            "start": 9,
            "end": 10,
            "name": "x"
          },
          "exported": {
            "type": "Identifier",
            "start": 14,
            "end": 15,
            "name": "y"
          }
        }
      ],
      "source": {
        "type": "Literal",
        "start": 23,
        "end": 31,
        "value": "./y.js",
        "raw": "'./y.js'"
      }
    },
    {
      "type": "ExportNamedDeclaration",
      "start": 33,
      "end": 65,
      "declaration": null,
      "specifiers": [
        {
          "type": "ExportSpecifier",
          "start": 42,
          "end": 48,
          "local": {
            "type": "Identifier",
            "start": 42,
            "end": 43,
            "name": "x"
          },
          "exported": {
            "type": "Identifier",
            "start": 47,
            "end": 48,
            "name": "z"
          }
        }
      ],
      "source": {
        "type": "Literal",
        "start": 56,
        "end": 64,
        "value": "./z.js",
        "raw": "'./z.js'"
      }
    }
  ]
}
	`, ast)
}

func TestHarmony257(t *testing.T) {
	ast, err := Compile("export { default as y } from './y.js';\nexport default 42;")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 57,
  "body": [
    {
      "type": "ExportNamedDeclaration",
      "start": 0,
      "end": 38,
      "declaration": null,
      "specifiers": [
        {
          "type": "ExportSpecifier",
          "start": 9,
          "end": 21,
          "local": {
            "type": "Identifier",
            "start": 9,
            "end": 16,
            "name": "default"
          },
          "exported": {
            "type": "Identifier",
            "start": 20,
            "end": 21,
            "name": "y"
          }
        }
      ],
      "source": {
        "type": "Literal",
        "start": 29,
        "end": 37,
        "value": "./y.js",
        "raw": "'./y.js'"
      }
    },
    {
      "type": "ExportDefaultDeclaration",
      "start": 39,
      "end": 57,
      "declaration": {
        "type": "Literal",
        "start": 54,
        "end": 56,
        "value": 42,
        "raw": "42"
      }
    }
  ]
}
	`, ast)
}

func TestHarmony258(t *testing.T) {
	ast, err := Compile("[x, (y), {z, u: (v)}] = foo")
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
        "type": "AssignmentExpression",
        "start": 0,
        "end": 27,
        "operator": "=",
        "left": {
          "type": "ArrayPattern",
          "start": 0,
          "end": 21,
          "elements": [
            {
              "type": "Identifier",
              "start": 1,
              "end": 2,
              "name": "x"
            },
            {
              "type": "Identifier",
              "start": 5,
              "end": 6,
              "name": "y"
            },
            {
              "type": "ObjectPattern",
              "start": 9,
              "end": 20,
              "properties": [
                {
                  "type": "Property",
                  "start": 10,
                  "end": 11,
                  "method": false,
                  "shorthand": true,
                  "computed": false,
                  "key": {
                    "type": "Identifier",
                    "start": 10,
                    "end": 11,
                    "name": "z"
                  },
                  "kind": "init",
                  "value": {
                    "type": "Identifier",
                    "start": 10,
                    "end": 11,
                    "name": "z"
                  }
                },
                {
                  "type": "Property",
                  "start": 13,
                  "end": 19,
                  "method": false,
                  "shorthand": false,
                  "computed": false,
                  "key": {
                    "type": "Identifier",
                    "start": 13,
                    "end": 14,
                    "name": "u"
                  },
                  "value": {
                    "type": "Identifier",
                    "start": 17,
                    "end": 18,
                    "name": "v"
                  },
                  "kind": "init"
                }
              ]
            }
          ]
        },
        "right": {
          "type": "Identifier",
          "start": 24,
          "end": 27,
          "name": "foo"
        }
      }
    }
  ]
}
	`, ast)
}

func TestHarmony259(t *testing.T) {
	ast, err := Compile("export default function(x) {};")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 30,
  "body": [
    {
      "type": "ExportDefaultDeclaration",
      "start": 0,
      "end": 29,
      "declaration": {
        "type": "FunctionDeclaration",
        "start": 15,
        "end": 29,
        "id": null,
        "generator": false,
        "async": false,
        "params": [
          {
            "type": "Identifier",
            "start": 24,
            "end": 25,
            "name": "x"
          }
        ],
        "body": {
          "type": "BlockStatement",
          "start": 27,
          "end": 29,
          "body": []
        }
      }
    },
    {
      "type": "EmptyStatement",
      "start": 29,
      "end": 30
    }
  ]
}
	`, ast)
}

func TestHarmony260(t *testing.T) {
	ast, err := Compile("var foo = 1; var foo = 1;")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 25,
  "body": [
    {
      "type": "VariableDeclaration",
      "start": 0,
      "end": 12,
      "declarations": [
        {
          "type": "VariableDeclarator",
          "start": 4,
          "end": 11,
          "id": {
            "type": "Identifier",
            "start": 4,
            "end": 7,
            "name": "foo"
          },
          "init": {
            "type": "Literal",
            "start": 10,
            "end": 11,
            "value": 1,
            "raw": "1"
          }
        }
      ],
      "kind": "var"
    },
    {
      "type": "VariableDeclaration",
      "start": 13,
      "end": 25,
      "declarations": [
        {
          "type": "VariableDeclarator",
          "start": 17,
          "end": 24,
          "id": {
            "type": "Identifier",
            "start": 17,
            "end": 20,
            "name": "foo"
          },
          "init": {
            "type": "Literal",
            "start": 23,
            "end": 24,
            "value": 1,
            "raw": "1"
          }
        }
      ],
      "kind": "var"
    }
  ]
}
	`, ast)
}

func TestHarmony261(t *testing.T) {
	ast, err := Compile("if (x) var foo = 1; var foo = 1;")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 32,
  "body": [
    {
      "type": "IfStatement",
      "start": 0,
      "end": 19,
      "test": {
        "type": "Identifier",
        "start": 4,
        "end": 5,
        "name": "x"
      },
      "consequent": {
        "type": "VariableDeclaration",
        "start": 7,
        "end": 19,
        "declarations": [
          {
            "type": "VariableDeclarator",
            "start": 11,
            "end": 18,
            "id": {
              "type": "Identifier",
              "start": 11,
              "end": 14,
              "name": "foo"
            },
            "init": {
              "type": "Literal",
              "start": 17,
              "end": 18,
              "value": 1,
              "raw": "1"
            }
          }
        ],
        "kind": "var"
      },
      "alternate": null
    },
    {
      "type": "VariableDeclaration",
      "start": 20,
      "end": 32,
      "declarations": [
        {
          "type": "VariableDeclarator",
          "start": 24,
          "end": 31,
          "id": {
            "type": "Identifier",
            "start": 24,
            "end": 27,
            "name": "foo"
          },
          "init": {
            "type": "Literal",
            "start": 30,
            "end": 31,
            "value": 1,
            "raw": "1"
          }
        }
      ],
      "kind": "var"
    }
  ]
}
	`, ast)
}

func TestHarmony262(t *testing.T) {
	ast, err := Compile("function x() { var foo = 1; } let foo = 1;")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 42,
  "body": [
    {
      "type": "FunctionDeclaration",
      "start": 0,
      "end": 29,
      "id": {
        "type": "Identifier",
        "start": 9,
        "end": 10,
        "name": "x"
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
            "type": "VariableDeclaration",
            "start": 15,
            "end": 27,
            "declarations": [
              {
                "type": "VariableDeclarator",
                "start": 19,
                "end": 26,
                "id": {
                  "type": "Identifier",
                  "start": 19,
                  "end": 22,
                  "name": "foo"
                },
                "init": {
                  "type": "Literal",
                  "start": 25,
                  "end": 26,
                  "value": 1,
                  "raw": "1"
                }
              }
            ],
            "kind": "var"
          }
        ]
      }
    },
    {
      "type": "VariableDeclaration",
      "start": 30,
      "end": 42,
      "declarations": [
        {
          "type": "VariableDeclarator",
          "start": 34,
          "end": 41,
          "id": {
            "type": "Identifier",
            "start": 34,
            "end": 37,
            "name": "foo"
          },
          "init": {
            "type": "Literal",
            "start": 40,
            "end": 41,
            "value": 1,
            "raw": "1"
          }
        }
      ],
      "kind": "let"
    }
  ]
}
	`, ast)
}

func TestHarmony263(t *testing.T) {
	ast, err := Compile("function foo() { let foo = 1; }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
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
        "end": 12,
        "name": "foo"
      },
      "generator": false,
      "async": false,
      "params": [],
      "body": {
        "type": "BlockStatement",
        "start": 15,
        "end": 31,
        "body": [
          {
            "type": "VariableDeclaration",
            "start": 17,
            "end": 29,
            "declarations": [
              {
                "type": "VariableDeclarator",
                "start": 21,
                "end": 28,
                "id": {
                  "type": "Identifier",
                  "start": 21,
                  "end": 24,
                  "name": "foo"
                },
                "init": {
                  "type": "Literal",
                  "start": 27,
                  "end": 28,
                  "value": 1,
                  "raw": "1"
                }
              }
            ],
            "kind": "let"
          }
        ]
      }
    }
  ]
}
	`, ast)
}

func TestHarmony264(t *testing.T) {
	ast, err := Compile("var foo = 1; { let foo = 1; }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 29,
  "body": [
    {
      "type": "VariableDeclaration",
      "start": 0,
      "end": 12,
      "declarations": [
        {
          "type": "VariableDeclarator",
          "start": 4,
          "end": 11,
          "id": {
            "type": "Identifier",
            "start": 4,
            "end": 7,
            "name": "foo"
          },
          "init": {
            "type": "Literal",
            "start": 10,
            "end": 11,
            "value": 1,
            "raw": "1"
          }
        }
      ],
      "kind": "var"
    },
    {
      "type": "BlockStatement",
      "start": 13,
      "end": 29,
      "body": [
        {
          "type": "VariableDeclaration",
          "start": 15,
          "end": 27,
          "declarations": [
            {
              "type": "VariableDeclarator",
              "start": 19,
              "end": 26,
              "id": {
                "type": "Identifier",
                "start": 19,
                "end": 22,
                "name": "foo"
              },
              "init": {
                "type": "Literal",
                "start": 25,
                "end": 26,
                "value": 1,
                "raw": "1"
              }
            }
          ],
          "kind": "let"
        }
      ]
    }
  ]
}
	`, ast)
}

func TestHarmony265(t *testing.T) {
	ast, err := Compile("{ let foo = 1; { let foo = 2; } }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 33,
  "body": [
    {
      "type": "BlockStatement",
      "start": 0,
      "end": 33,
      "body": [
        {
          "type": "VariableDeclaration",
          "start": 2,
          "end": 14,
          "declarations": [
            {
              "type": "VariableDeclarator",
              "start": 6,
              "end": 13,
              "id": {
                "type": "Identifier",
                "start": 6,
                "end": 9,
                "name": "foo"
              },
              "init": {
                "type": "Literal",
                "start": 12,
                "end": 13,
                "value": 1,
                "raw": "1"
              }
            }
          ],
          "kind": "let"
        },
        {
          "type": "BlockStatement",
          "start": 15,
          "end": 31,
          "body": [
            {
              "type": "VariableDeclaration",
              "start": 17,
              "end": 29,
              "declarations": [
                {
                  "type": "VariableDeclarator",
                  "start": 21,
                  "end": 28,
                  "id": {
                    "type": "Identifier",
                    "start": 21,
                    "end": 24,
                    "name": "foo"
                  },
                  "init": {
                    "type": "Literal",
                    "start": 27,
                    "end": 28,
                    "value": 2,
                    "raw": "2"
                  }
                }
              ],
              "kind": "let"
            }
          ]
        }
      ]
    }
  ]
}
	`, ast)
}

func TestHarmony266(t *testing.T) {
	ast, err := Compile("var foo; try {} catch (_) { let foo; }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 38,
  "body": [
    {
      "type": "VariableDeclaration",
      "start": 0,
      "end": 8,
      "declarations": [
        {
          "type": "VariableDeclarator",
          "start": 4,
          "end": 7,
          "id": {
            "type": "Identifier",
            "start": 4,
            "end": 7,
            "name": "foo"
          },
          "init": null
        }
      ],
      "kind": "var"
    },
    {
      "type": "TryStatement",
      "start": 9,
      "end": 38,
      "block": {
        "type": "BlockStatement",
        "start": 13,
        "end": 15,
        "body": []
      },
      "handler": {
        "type": "CatchClause",
        "start": 16,
        "end": 38,
        "param": {
          "type": "Identifier",
          "start": 23,
          "end": 24,
          "name": "_"
        },
        "body": {
          "type": "BlockStatement",
          "start": 26,
          "end": 38,
          "body": [
            {
              "type": "VariableDeclaration",
              "start": 28,
              "end": 36,
              "declarations": [
                {
                  "type": "VariableDeclarator",
                  "start": 32,
                  "end": 35,
                  "id": {
                    "type": "Identifier",
                    "start": 32,
                    "end": 35,
                    "name": "foo"
                  },
                  "init": null
                }
              ],
              "kind": "let"
            }
          ]
        }
      },
      "finalizer": null
    }
  ]
}
	`, ast)
}

func TestHarmony267(t *testing.T) {
	ast, err := Compile("let x = 1; function foo(x) {}")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 29,
  "body": [
    {
      "type": "VariableDeclaration",
      "start": 0,
      "end": 10,
      "declarations": [
        {
          "type": "VariableDeclarator",
          "start": 4,
          "end": 9,
          "id": {
            "type": "Identifier",
            "start": 4,
            "end": 5,
            "name": "x"
          },
          "init": {
            "type": "Literal",
            "start": 8,
            "end": 9,
            "value": 1,
            "raw": "1"
          }
        }
      ],
      "kind": "let"
    },
    {
      "type": "FunctionDeclaration",
      "start": 11,
      "end": 29,
      "id": {
        "type": "Identifier",
        "start": 20,
        "end": 23,
        "name": "foo"
      },
      "generator": false,
      "async": false,
      "params": [
        {
          "type": "Identifier",
          "start": 24,
          "end": 25,
          "name": "x"
        }
      ],
      "body": {
        "type": "BlockStatement",
        "start": 27,
        "end": 29,
        "body": []
      }
    }
  ]
}
	`, ast)
}

func TestHarmony268(t *testing.T) {
	ast, err := Compile("for (let i = 0;;); for (let i = 0;;);")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 37,
  "body": [
    {
      "type": "ForStatement",
      "start": 0,
      "end": 18,
      "init": {
        "type": "VariableDeclaration",
        "start": 5,
        "end": 14,
        "declarations": [
          {
            "type": "VariableDeclarator",
            "start": 9,
            "end": 14,
            "id": {
              "type": "Identifier",
              "start": 9,
              "end": 10,
              "name": "i"
            },
            "init": {
              "type": "Literal",
              "start": 13,
              "end": 14,
              "value": 0,
              "raw": "0"
            }
          }
        ],
        "kind": "let"
      },
      "test": null,
      "update": null,
      "body": {
        "type": "EmptyStatement",
        "start": 17,
        "end": 18
      }
    },
    {
      "type": "ForStatement",
      "start": 19,
      "end": 37,
      "init": {
        "type": "VariableDeclaration",
        "start": 24,
        "end": 33,
        "declarations": [
          {
            "type": "VariableDeclarator",
            "start": 28,
            "end": 33,
            "id": {
              "type": "Identifier",
              "start": 28,
              "end": 29,
              "name": "i"
            },
            "init": {
              "type": "Literal",
              "start": 32,
              "end": 33,
              "value": 0,
              "raw": "0"
            }
          }
        ],
        "kind": "let"
      },
      "test": null,
      "update": null,
      "body": {
        "type": "EmptyStatement",
        "start": 36,
        "end": 37
      }
    }
  ]
}
	`, ast)
}

func TestHarmony269(t *testing.T) {
	ast, err := Compile("for (const foo of bar); for (const foo of bar);")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 47,
  "body": [
    {
      "type": "ForOfStatement",
      "start": 0,
      "end": 23,
      "await": false,
      "left": {
        "type": "VariableDeclaration",
        "start": 5,
        "end": 14,
        "declarations": [
          {
            "type": "VariableDeclarator",
            "start": 11,
            "end": 14,
            "id": {
              "type": "Identifier",
              "start": 11,
              "end": 14,
              "name": "foo"
            },
            "init": null
          }
        ],
        "kind": "const"
      },
      "right": {
        "type": "Identifier",
        "start": 18,
        "end": 21,
        "name": "bar"
      },
      "body": {
        "type": "EmptyStatement",
        "start": 22,
        "end": 23
      }
    },
    {
      "type": "ForOfStatement",
      "start": 24,
      "end": 47,
      "await": false,
      "left": {
        "type": "VariableDeclaration",
        "start": 29,
        "end": 38,
        "declarations": [
          {
            "type": "VariableDeclarator",
            "start": 35,
            "end": 38,
            "id": {
              "type": "Identifier",
              "start": 35,
              "end": 38,
              "name": "foo"
            },
            "init": null
          }
        ],
        "kind": "const"
      },
      "right": {
        "type": "Identifier",
        "start": 42,
        "end": 45,
        "name": "bar"
      },
      "body": {
        "type": "EmptyStatement",
        "start": 46,
        "end": 47
      }
    }
  ]
}
	`, ast)
}

func TestHarmony270(t *testing.T) {
	ast, err := Compile("for (const foo in bar); for (const foo in bar);")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 47,
  "body": [
    {
      "type": "ForInStatement",
      "start": 0,
      "end": 23,
      "left": {
        "type": "VariableDeclaration",
        "start": 5,
        "end": 14,
        "declarations": [
          {
            "type": "VariableDeclarator",
            "start": 11,
            "end": 14,
            "id": {
              "type": "Identifier",
              "start": 11,
              "end": 14,
              "name": "foo"
            },
            "init": null
          }
        ],
        "kind": "const"
      },
      "right": {
        "type": "Identifier",
        "start": 18,
        "end": 21,
        "name": "bar"
      },
      "body": {
        "type": "EmptyStatement",
        "start": 22,
        "end": 23
      }
    },
    {
      "type": "ForInStatement",
      "start": 24,
      "end": 47,
      "left": {
        "type": "VariableDeclaration",
        "start": 29,
        "end": 38,
        "declarations": [
          {
            "type": "VariableDeclarator",
            "start": 35,
            "end": 38,
            "id": {
              "type": "Identifier",
              "start": 35,
              "end": 38,
              "name": "foo"
            },
            "init": null
          }
        ],
        "kind": "const"
      },
      "right": {
        "type": "Identifier",
        "start": 42,
        "end": 45,
        "name": "bar"
      },
      "body": {
        "type": "EmptyStatement",
        "start": 46,
        "end": 47
      }
    }
  ]
}
	`, ast)
}

func TestHarmony271(t *testing.T) {
	ast, err := Compile("for (let foo in bar) { let foo = 1; }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 37,
  "body": [
    {
      "type": "ForInStatement",
      "start": 0,
      "end": 37,
      "left": {
        "type": "VariableDeclaration",
        "start": 5,
        "end": 12,
        "declarations": [
          {
            "type": "VariableDeclarator",
            "start": 9,
            "end": 12,
            "id": {
              "type": "Identifier",
              "start": 9,
              "end": 12,
              "name": "foo"
            },
            "init": null
          }
        ],
        "kind": "let"
      },
      "right": {
        "type": "Identifier",
        "start": 16,
        "end": 19,
        "name": "bar"
      },
      "body": {
        "type": "BlockStatement",
        "start": 21,
        "end": 37,
        "body": [
          {
            "type": "VariableDeclaration",
            "start": 23,
            "end": 35,
            "declarations": [
              {
                "type": "VariableDeclarator",
                "start": 27,
                "end": 34,
                "id": {
                  "type": "Identifier",
                  "start": 27,
                  "end": 30,
                  "name": "foo"
                },
                "init": {
                  "type": "Literal",
                  "start": 33,
                  "end": 34,
                  "value": 1,
                  "raw": "1"
                }
              }
            ],
            "kind": "let"
          }
        ]
      }
    }
  ]
}
	`, ast)
}

func TestHarmony272(t *testing.T) {
	ast, err := Compile("for (let foo of bar) { let foo = 1; }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 37,
  "body": [
    {
      "type": "ForOfStatement",
      "start": 0,
      "end": 37,
      "await": false,
      "left": {
        "type": "VariableDeclaration",
        "start": 5,
        "end": 12,
        "declarations": [
          {
            "type": "VariableDeclarator",
            "start": 9,
            "end": 12,
            "id": {
              "type": "Identifier",
              "start": 9,
              "end": 12,
              "name": "foo"
            },
            "init": null
          }
        ],
        "kind": "let"
      },
      "right": {
        "type": "Identifier",
        "start": 16,
        "end": 19,
        "name": "bar"
      },
      "body": {
        "type": "BlockStatement",
        "start": 21,
        "end": 37,
        "body": [
          {
            "type": "VariableDeclaration",
            "start": 23,
            "end": 35,
            "declarations": [
              {
                "type": "VariableDeclarator",
                "start": 27,
                "end": 34,
                "id": {
                  "type": "Identifier",
                  "start": 27,
                  "end": 30,
                  "name": "foo"
                },
                "init": {
                  "type": "Literal",
                  "start": 33,
                  "end": 34,
                  "value": 1,
                  "raw": "1"
                }
              }
            ],
            "kind": "let"
          }
        ]
      }
    }
  ]
}
	`, ast)
}

func TestHarmony273(t *testing.T) {
	ast, err := Compile("class Foo { method(foo) {} method2() { let foo; } }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
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
        "end": 9,
        "name": "Foo"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "start": 10,
        "end": 51,
        "body": [
          {
            "type": "MethodDefinition",
            "start": 12,
            "end": 26,
            "static": false,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 12,
              "end": 18,
              "name": "method"
            },
            "kind": "method",
            "value": {
              "type": "FunctionExpression",
              "start": 18,
              "end": 26,
              "id": null,
              "expression": false,
              "generator": false,
              "async": false,
              "params": [
                {
                  "type": "Identifier",
                  "start": 19,
                  "end": 22,
                  "name": "foo"
                }
              ],
              "body": {
                "type": "BlockStatement",
                "start": 24,
                "end": 26,
                "body": []
              }
            }
          },
          {
            "type": "MethodDefinition",
            "start": 27,
            "end": 49,
            "static": false,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 27,
              "end": 34,
              "name": "method2"
            },
            "kind": "method",
            "value": {
              "type": "FunctionExpression",
              "start": 34,
              "end": 49,
              "id": null,
              "expression": false,
              "generator": false,
              "async": false,
              "params": [],
              "body": {
                "type": "BlockStatement",
                "start": 37,
                "end": 49,
                "body": [
                  {
                    "type": "VariableDeclaration",
                    "start": 39,
                    "end": 47,
                    "declarations": [
                      {
                        "type": "VariableDeclarator",
                        "start": 43,
                        "end": 46,
                        "id": {
                          "type": "Identifier",
                          "start": 43,
                          "end": 46,
                          "name": "foo"
                        },
                        "init": null
                      }
                    ],
                    "kind": "let"
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

func TestHarmony274(t *testing.T) {
	ast, err := Compile("() => { let foo; }; foo => {}")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 29,
  "body": [
    {
      "type": "ExpressionStatement",
      "start": 0,
      "end": 19,
      "expression": {
        "type": "ArrowFunctionExpression",
        "start": 0,
        "end": 18,
        "id": null,
        "expression": false,
        "generator": false,
        "async": false,
        "params": [],
        "body": {
          "type": "BlockStatement",
          "start": 6,
          "end": 18,
          "body": [
            {
              "type": "VariableDeclaration",
              "start": 8,
              "end": 16,
              "declarations": [
                {
                  "type": "VariableDeclarator",
                  "start": 12,
                  "end": 15,
                  "id": {
                    "type": "Identifier",
                    "start": 12,
                    "end": 15,
                    "name": "foo"
                  },
                  "init": null
                }
              ],
              "kind": "let"
            }
          ]
        }
      }
    },
    {
      "type": "ExpressionStatement",
      "start": 20,
      "end": 29,
      "expression": {
        "type": "ArrowFunctionExpression",
        "start": 20,
        "end": 29,
        "id": null,
        "expression": false,
        "generator": false,
        "async": false,
        "params": [
          {
            "type": "Identifier",
            "start": 20,
            "end": 23,
            "name": "foo"
          }
        ],
        "body": {
          "type": "BlockStatement",
          "start": 27,
          "end": 29,
          "body": []
        }
      }
    }
  ]
}
	`, ast)
}

func TestHarmony275(t *testing.T) {
	ast, err := Compile("() => { let foo; }; () => { let foo; }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 38,
  "body": [
    {
      "type": "ExpressionStatement",
      "start": 0,
      "end": 19,
      "expression": {
        "type": "ArrowFunctionExpression",
        "start": 0,
        "end": 18,
        "id": null,
        "expression": false,
        "generator": false,
        "async": false,
        "params": [],
        "body": {
          "type": "BlockStatement",
          "start": 6,
          "end": 18,
          "body": [
            {
              "type": "VariableDeclaration",
              "start": 8,
              "end": 16,
              "declarations": [
                {
                  "type": "VariableDeclarator",
                  "start": 12,
                  "end": 15,
                  "id": {
                    "type": "Identifier",
                    "start": 12,
                    "end": 15,
                    "name": "foo"
                  },
                  "init": null
                }
              ],
              "kind": "let"
            }
          ]
        }
      }
    },
    {
      "type": "ExpressionStatement",
      "start": 20,
      "end": 38,
      "expression": {
        "type": "ArrowFunctionExpression",
        "start": 20,
        "end": 38,
        "id": null,
        "expression": false,
        "generator": false,
        "async": false,
        "params": [],
        "body": {
          "type": "BlockStatement",
          "start": 26,
          "end": 38,
          "body": [
            {
              "type": "VariableDeclaration",
              "start": 28,
              "end": 36,
              "declarations": [
                {
                  "type": "VariableDeclarator",
                  "start": 32,
                  "end": 35,
                  "id": {
                    "type": "Identifier",
                    "start": 32,
                    "end": 35,
                    "name": "foo"
                  },
                  "init": null
                }
              ],
              "kind": "let"
            }
          ]
        }
      }
    }
  ]
}
	`, ast)
}

func TestHarmony276(t *testing.T) {
	ast, err := Compile("switch(x) { case 1: let foo = 1; } let foo = 1;")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 47,
  "body": [
    {
      "type": "SwitchStatement",
      "start": 0,
      "end": 34,
      "discriminant": {
        "type": "Identifier",
        "start": 7,
        "end": 8,
        "name": "x"
      },
      "cases": [
        {
          "type": "SwitchCase",
          "start": 12,
          "end": 32,
          "consequent": [
            {
              "type": "VariableDeclaration",
              "start": 20,
              "end": 32,
              "declarations": [
                {
                  "type": "VariableDeclarator",
                  "start": 24,
                  "end": 31,
                  "id": {
                    "type": "Identifier",
                    "start": 24,
                    "end": 27,
                    "name": "foo"
                  },
                  "init": {
                    "type": "Literal",
                    "start": 30,
                    "end": 31,
                    "value": 1,
                    "raw": "1"
                  }
                }
              ],
              "kind": "let"
            }
          ],
          "test": {
            "type": "Literal",
            "start": 17,
            "end": 18,
            "value": 1,
            "raw": "1"
          }
        }
      ]
    },
    {
      "type": "VariableDeclaration",
      "start": 35,
      "end": 47,
      "declarations": [
        {
          "type": "VariableDeclarator",
          "start": 39,
          "end": 46,
          "id": {
            "type": "Identifier",
            "start": 39,
            "end": 42,
            "name": "foo"
          },
          "init": {
            "type": "Literal",
            "start": 45,
            "end": 46,
            "value": 1,
            "raw": "1"
          }
        }
      ],
      "kind": "let"
    }
  ]
}
	`, ast)
}

func TestHarmony277(t *testing.T) {
	ast, err := Compile("'use strict'; function foo() { let foo = 1; }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 45,
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
      }
    },
    {
      "type": "FunctionDeclaration",
      "start": 14,
      "end": 45,
      "id": {
        "type": "Identifier",
        "start": 23,
        "end": 26,
        "name": "foo"
      },
      "generator": false,
      "async": false,
      "params": [],
      "body": {
        "type": "BlockStatement",
        "start": 29,
        "end": 45,
        "body": [
          {
            "type": "VariableDeclaration",
            "start": 31,
            "end": 43,
            "declarations": [
              {
                "type": "VariableDeclarator",
                "start": 35,
                "end": 42,
                "id": {
                  "type": "Identifier",
                  "start": 35,
                  "end": 38,
                  "name": "foo"
                },
                "init": {
                  "type": "Literal",
                  "start": 41,
                  "end": 42,
                  "value": 1,
                  "raw": "1"
                }
              }
            ],
            "kind": "let"
          }
        ]
      }
    }
  ]
}
	`, ast)
}

func TestHarmony278(t *testing.T) {
	ast, err := Compile("let foo = 1; function x() { var foo = 1; }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 42,
  "body": [
    {
      "type": "VariableDeclaration",
      "start": 0,
      "end": 12,
      "declarations": [
        {
          "type": "VariableDeclarator",
          "start": 4,
          "end": 11,
          "id": {
            "type": "Identifier",
            "start": 4,
            "end": 7,
            "name": "foo"
          },
          "init": {
            "type": "Literal",
            "start": 10,
            "end": 11,
            "value": 1,
            "raw": "1"
          }
        }
      ],
      "kind": "let"
    },
    {
      "type": "FunctionDeclaration",
      "start": 13,
      "end": 42,
      "id": {
        "type": "Identifier",
        "start": 22,
        "end": 23,
        "name": "x"
      },
      "generator": false,
      "async": false,
      "params": [],
      "body": {
        "type": "BlockStatement",
        "start": 26,
        "end": 42,
        "body": [
          {
            "type": "VariableDeclaration",
            "start": 28,
            "end": 40,
            "declarations": [
              {
                "type": "VariableDeclarator",
                "start": 32,
                "end": 39,
                "id": {
                  "type": "Identifier",
                  "start": 32,
                  "end": 35,
                  "name": "foo"
                },
                "init": {
                  "type": "Literal",
                  "start": 38,
                  "end": 39,
                  "value": 1,
                  "raw": "1"
                }
              }
            ],
            "kind": "var"
          }
        ]
      }
    }
  ]
}
	`, ast)
}

func TestHarmony279(t *testing.T) {
	ast, err := Compile("[...foo, bar = 1]")
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
        "type": "ArrayExpression",
        "start": 0,
        "end": 17,
        "elements": [
          {
            "type": "SpreadElement",
            "start": 1,
            "end": 7,
            "argument": {
              "type": "Identifier",
              "start": 4,
              "end": 7,
              "name": "foo"
            }
          },
          {
            "type": "AssignmentExpression",
            "start": 9,
            "end": 16,
            "operator": "=",
            "left": {
              "type": "Identifier",
              "start": 9,
              "end": 12,
              "name": "bar"
            },
            "right": {
              "type": "Literal",
              "start": 15,
              "end": 16,
              "value": 1,
              "raw": "1"
            }
          }
        ]
      }
    }
  ]
}
	`, ast)
}

func TestHarmony280(t *testing.T) {
	ast, err := Compile("for (var a of /b/) {}")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 21,
  "body": [
    {
      "type": "ForOfStatement",
      "start": 0,
      "end": 21,
      "await": false,
      "left": {
        "type": "VariableDeclaration",
        "start": 5,
        "end": 10,
        "declarations": [
          {
            "type": "VariableDeclarator",
            "start": 9,
            "end": 10,
            "id": {
              "type": "Identifier",
              "start": 9,
              "end": 10,
              "name": "a"
            },
            "init": null
          }
        ],
        "kind": "var"
      },
      "right": {
        "type": "Literal",
        "start": 14,
        "end": 17,
        "value": {},
        "regexp": {
          "pattern": "b",
          "flags": ""
        }
      },
      "body": {
        "type": "BlockStatement",
        "start": 19,
        "end": 21,
        "body": []
      }
    }
  ]
}
	`, ast)
}

func TestHarmony281(t *testing.T) {
	ast, err := Compile("for (var {a} of /b/) {}")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 23,
  "body": [
    {
      "type": "ForOfStatement",
      "start": 0,
      "end": 23,
      "await": false,
      "left": {
        "type": "VariableDeclaration",
        "start": 5,
        "end": 12,
        "declarations": [
          {
            "type": "VariableDeclarator",
            "start": 9,
            "end": 12,
            "id": {
              "type": "ObjectPattern",
              "start": 9,
              "end": 12,
              "properties": [
                {
                  "type": "Property",
                  "start": 10,
                  "end": 11,
                  "method": false,
                  "shorthand": true,
                  "computed": false,
                  "key": {
                    "type": "Identifier",
                    "start": 10,
                    "end": 11,
                    "name": "a"
                  },
                  "kind": "init",
                  "value": {
                    "type": "Identifier",
                    "start": 10,
                    "end": 11,
                    "name": "a"
                  }
                }
              ]
            },
            "init": null
          }
        ],
        "kind": "var"
      },
      "right": {
        "type": "Literal",
        "start": 16,
        "end": 19,
        "value": {},
        "regexp": {
          "pattern": "b",
          "flags": ""
        }
      },
      "body": {
        "type": "BlockStatement",
        "start": 21,
        "end": 23,
        "body": []
      }
    }
  ]
}
	`, ast)
}

func TestHarmony282(t *testing.T) {
	ast, err := Compile("for (let {a} of /b/) {}")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 23,
  "body": [
    {
      "type": "ForOfStatement",
      "start": 0,
      "end": 23,
      "await": false,
      "left": {
        "type": "VariableDeclaration",
        "start": 5,
        "end": 12,
        "declarations": [
          {
            "type": "VariableDeclarator",
            "start": 9,
            "end": 12,
            "id": {
              "type": "ObjectPattern",
              "start": 9,
              "end": 12,
              "properties": [
                {
                  "type": "Property",
                  "start": 10,
                  "end": 11,
                  "method": false,
                  "shorthand": true,
                  "computed": false,
                  "key": {
                    "type": "Identifier",
                    "start": 10,
                    "end": 11,
                    "name": "a"
                  },
                  "kind": "init",
                  "value": {
                    "type": "Identifier",
                    "start": 10,
                    "end": 11,
                    "name": "a"
                  }
                }
              ]
            },
            "init": null
          }
        ],
        "kind": "let"
      },
      "right": {
        "type": "Literal",
        "start": 16,
        "end": 19,
        "value": {},
        "regexp": {
          "pattern": "b",
          "flags": ""
        }
      },
      "body": {
        "type": "BlockStatement",
        "start": 21,
        "end": 23,
        "body": []
      }
    }
  ]
}
	`, ast)
}

func TestHarmony283(t *testing.T) {
	ast, err := Compile("for (const {a} of /b/) {}")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 25,
  "body": [
    {
      "type": "ForOfStatement",
      "start": 0,
      "end": 25,
      "await": false,
      "left": {
        "type": "VariableDeclaration",
        "start": 5,
        "end": 14,
        "declarations": [
          {
            "type": "VariableDeclarator",
            "start": 11,
            "end": 14,
            "id": {
              "type": "ObjectPattern",
              "start": 11,
              "end": 14,
              "properties": [
                {
                  "type": "Property",
                  "start": 12,
                  "end": 13,
                  "method": false,
                  "shorthand": true,
                  "computed": false,
                  "key": {
                    "type": "Identifier",
                    "start": 12,
                    "end": 13,
                    "name": "a"
                  },
                  "kind": "init",
                  "value": {
                    "type": "Identifier",
                    "start": 12,
                    "end": 13,
                    "name": "a"
                  }
                }
              ]
            },
            "init": null
          }
        ],
        "kind": "const"
      },
      "right": {
        "type": "Literal",
        "start": 18,
        "end": 21,
        "value": {},
        "regexp": {
          "pattern": "b",
          "flags": ""
        }
      },
      "body": {
        "type": "BlockStatement",
        "start": 23,
        "end": 25,
        "body": []
      }
    }
  ]
}
	`, ast)
}

func TestHarmony284(t *testing.T) {
	ast, err := Compile("function* bar() { yield /re/ }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 30,
  "body": [
    {
      "type": "FunctionDeclaration",
      "start": 0,
      "end": 30,
      "id": {
        "type": "Identifier",
        "start": 10,
        "end": 13,
        "name": "bar"
      },
      "generator": true,
      "async": false,
      "params": [],
      "body": {
        "type": "BlockStatement",
        "start": 16,
        "end": 30,
        "body": [
          {
            "type": "ExpressionStatement",
            "start": 18,
            "end": 28,
            "expression": {
              "type": "YieldExpression",
              "start": 18,
              "end": 28,
              "delegate": false,
              "argument": {
                "type": "Literal",
                "start": 24,
                "end": 28,
                "value": {},
                "regexp": {
                  "pattern": "re",
                  "flags": ""
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

func TestHarmony285(t *testing.T) {
	ast, err := Compile("function* bar() { yield class {} }")
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
        "start": 10,
        "end": 13,
        "name": "bar"
      },
      "generator": true,
      "async": false,
      "params": [],
      "body": {
        "type": "BlockStatement",
        "start": 16,
        "end": 34,
        "body": [
          {
            "type": "ExpressionStatement",
            "start": 18,
            "end": 32,
            "expression": {
              "type": "YieldExpression",
              "start": 18,
              "end": 32,
              "delegate": false,
              "argument": {
                "type": "ClassExpression",
                "start": 24,
                "end": 32,
                "id": null,
                "superClass": null,
                "body": {
                  "type": "ClassBody",
                  "start": 30,
                  "end": 32,
                  "body": []
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

func TestHarmony286(t *testing.T) {
	ast, err := Compile("() => {}\n/re/")
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
      "end": 8,
      "expression": {
        "type": "ArrowFunctionExpression",
        "start": 0,
        "end": 8,
        "id": null,
        "expression": false,
        "generator": false,
        "async": false,
        "params": [],
        "body": {
          "type": "BlockStatement",
          "start": 6,
          "end": 8,
          "body": []
        }
      }
    },
    {
      "type": "ExpressionStatement",
      "start": 9,
      "end": 13,
      "expression": {
        "type": "Literal",
        "start": 9,
        "end": 13,
        "value": {},
        "regexp": {
          "pattern": "re",
          "flags": ""
        }
      }
    }
  ]
}
	`, ast)
}

func TestHarmony287(t *testing.T) {
	ast, err := Compile("(() => {}) + 2")
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
        "type": "BinaryExpression",
        "start": 0,
        "end": 14,
        "left": {
          "type": "ArrowFunctionExpression",
          "start": 1,
          "end": 9,
          "id": null,
          "expression": false,
          "generator": false,
          "async": false,
          "params": [],
          "body": {
            "type": "BlockStatement",
            "start": 7,
            "end": 9,
            "body": []
          }
        },
        "operator": "+",
        "right": {
          "type": "Literal",
          "start": 13,
          "end": 14,
          "value": 2,
          "raw": "2"
        }
      }
    }
  ]
}
	`, ast)
}

func TestHarmony288(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_STRICT)
	ast, err := CompileWithOpts("function *f1() { function g() { return yield / 1 } }", opts)
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "FunctionDeclaration",
      "id": {
        "type": "Identifier",
        "name": "f1"
      },
      "params": [],
      "body": {
        "type": "BlockStatement",
        "body": [
          {
            "type": "FunctionDeclaration",
            "id": {
              "type": "Identifier",
              "name": "g"
            },
            "params": [],
            "body": {
              "type": "BlockStatement",
              "body": [
                {
                  "type": "ReturnStatement",
                  "argument": {
                    "type": "BinaryExpression",
                    "operator": "/",
                    "left": {
                      "type": "Identifier",
                      "name": "yield"
                    },
                    "right": {
                      "type": "Literal",
                      "value": 1,
                      "raw": "1"
                    }
                  }
                }
              ]
            },
            "generator": false,
            "async": false
          }
        ]
      },
      "generator": true,
      "async": false
    }
  ]
}
	`, ast)
}

func TestHarmony289(t *testing.T) {
	// Annex B allows function redeclaration for plain functions in sloppy mode
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_STRICT)
	ast, err := CompileWithOpts("function f() {} function f() {}", opts)
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "FunctionDeclaration",
      "id": {
        "type": "Identifier",
        "name": "f"
      },
      "params": [],
      "body": {
        "type": "BlockStatement",
        "body": []
      },
      "generator": false,
      "async": false
    },
    {
      "type": "FunctionDeclaration",
      "id": {
        "type": "Identifier",
        "name": "f"
      },
      "params": [],
      "body": {
        "type": "BlockStatement",
        "body": []
      },
      "generator": false,
      "async": false
    }
  ]
}
	`, ast)
}

func TestHarmony290(t *testing.T) {
	ast, err := Compile("class Foo {} /regexp/")
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
      "end": 12,
      "id": {
        "type": "Identifier",
        "start": 6,
        "end": 9,
        "name": "Foo"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "start": 10,
        "end": 12,
        "body": []
      }
    },
    {
      "type": "ExpressionStatement",
      "start": 13,
      "end": 21,
      "expression": {
        "type": "Literal",
        "start": 13,
        "end": 21,
        "value": {},
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

func TestHarmony291(t *testing.T) {
	ast, err := Compile("(class Foo {} / 2)")
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
        "type": "BinaryExpression",
        "start": 1,
        "end": 17,
        "left": {
          "type": "ClassExpression",
          "start": 1,
          "end": 13,
          "id": {
            "type": "Identifier",
            "start": 7,
            "end": 10,
            "name": "Foo"
          },
          "superClass": null,
          "body": {
            "type": "ClassBody",
            "start": 11,
            "end": 13,
            "body": []
          }
        },
        "operator": "/",
        "right": {
          "type": "Literal",
          "start": 16,
          "end": 17,
          "value": 2,
          "raw": "2"
        }
      }
    }
  ]
}
	`, ast)
}

func TestHarmony292(t *testing.T) {
	ast, err := Compile("1 <!--b")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "BinaryExpression",
        "operator": "<"
      }
    }
  ]
}
	`, ast)
}

func TestHarmony293(t *testing.T) {
	ast, err := Compile("({super: 1})")
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
        "type": "ObjectExpression",
        "start": 1,
        "end": 11,
        "properties": [
          {
            "type": "Property",
            "start": 2,
            "end": 10,
            "method": false,
            "shorthand": false,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 2,
              "end": 7,
              "name": "super"
            },
            "value": {
              "type": "Literal",
              "start": 9,
              "end": 10,
              "value": 1,
              "raw": "1"
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

func TestHarmony294(t *testing.T) {
	ast, err := Compile("import {super as a} from 'a'")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 28,
  "body": [
    {
      "type": "ImportDeclaration",
      "start": 0,
      "end": 28,
      "specifiers": [
        {
          "type": "ImportSpecifier",
          "start": 8,
          "end": 18,
          "imported": {
            "type": "Identifier",
            "start": 8,
            "end": 13,
            "name": "super"
          },
          "local": {
            "type": "Identifier",
            "start": 17,
            "end": 18,
            "name": "a"
          }
        }
      ],
      "source": {
        "type": "Literal",
        "start": 25,
        "end": 28,
        "value": "a",
        "raw": "'a'"
      }
    }
  ]
}
	`, ast)
}

func TestHarmony295(t *testing.T) {
	ast, err := Compile("function a() {} export {a as super}")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 35,
  "body": [
    {
      "type": "FunctionDeclaration",
      "start": 0,
      "end": 15,
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
        "end": 15,
        "body": []
      }
    },
    {
      "type": "ExportNamedDeclaration",
      "start": 16,
      "end": 35,
      "declaration": null,
      "specifiers": [
        {
          "type": "ExportSpecifier",
          "start": 24,
          "end": 34,
          "local": {
            "type": "Identifier",
            "start": 24,
            "end": 25,
            "name": "a"
          },
          "exported": {
            "type": "Identifier",
            "start": 29,
            "end": 34,
            "name": "super"
          }
        }
      ],
      "source": null
    }
  ]
}
	`, ast)
}

func TestHarmony296(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_STRICT)
	ast, err := CompileWithOpts("let instanceof Foo", opts)
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
        "type": "BinaryExpression",
        "start": 0,
        "end": 18,
        "left": {
          "type": "Identifier",
          "start": 0,
          "end": 3,
          "name": "let"
        },
        "operator": "instanceof",
        "right": {
          "type": "Identifier",
          "start": 15,
          "end": 18,
          "name": "Foo"
        }
      }
    }
  ]
}
	`, ast)
}

func TestHarmony297(t *testing.T) {
	ast, err := Compile("function fn({__proto__: a, __proto__: b}) {}")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 44,
  "body": [
    {
      "type": "FunctionDeclaration",
      "start": 0,
      "end": 44,
      "id": {
        "type": "Identifier",
        "start": 9,
        "end": 11,
        "name": "fn"
      },
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
              "end": 25,
              "method": false,
              "shorthand": false,
              "computed": false,
              "key": {
                "type": "Identifier",
                "start": 13,
                "end": 22,
                "name": "__proto__"
              },
              "value": {
                "type": "Identifier",
                "start": 24,
                "end": 25,
                "name": "a"
              },
              "kind": "init"
            },
            {
              "type": "Property",
              "start": 27,
              "end": 39,
              "method": false,
              "shorthand": false,
              "computed": false,
              "key": {
                "type": "Identifier",
                "start": 27,
                "end": 36,
                "name": "__proto__"
              },
              "value": {
                "type": "Identifier",
                "start": 38,
                "end": 39,
                "name": "b"
              },
              "kind": "init"
            }
          ]
        }
      ],
      "body": {
        "type": "BlockStatement",
        "start": 42,
        "end": 44,
        "body": []
      }
    }
  ]
}
	`, ast)
}

func TestHarmony298(t *testing.T) {
	ast, err := Compile("[...a, x][1] = b")
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
        "type": "AssignmentExpression",
        "start": 0,
        "end": 16,
        "operator": "=",
        "left": {
          "type": "MemberExpression",
          "start": 0,
          "end": 12,
          "object": {
            "type": "ArrayExpression",
            "start": 0,
            "end": 9,
            "elements": [
              {
                "type": "SpreadElement",
                "start": 1,
                "end": 5,
                "argument": {
                  "type": "Identifier",
                  "start": 4,
                  "end": 5,
                  "name": "a"
                }
              },
              {
                "type": "Identifier",
                "start": 7,
                "end": 8,
                "name": "x"
              }
            ]
          },
          "property": {
            "type": "Literal",
            "start": 10,
            "end": 11,
            "value": 1,
            "raw": "1"
          },
          "computed": true,
          "optional": false
        },
        "right": {
          "type": "Identifier",
          "start": 15,
          "end": 16,
          "name": "b"
        }
      }
    }
  ]
}
	`, ast)
}

func TestHarmony299(t *testing.T) {
	ast, err := Compile("for ([...foo, bar].baz in qux);")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 31,
  "body": [
    {
      "type": "ForInStatement",
      "start": 0,
      "end": 31,
      "left": {
        "type": "MemberExpression",
        "start": 5,
        "end": 22,
        "object": {
          "type": "ArrayExpression",
          "start": 5,
          "end": 18,
          "elements": [
            {
              "type": "SpreadElement",
              "start": 6,
              "end": 12,
              "argument": {
                "type": "Identifier",
                "start": 9,
                "end": 12,
                "name": "foo"
              }
            },
            {
              "type": "Identifier",
              "start": 14,
              "end": 17,
              "name": "bar"
            }
          ]
        },
        "property": {
          "type": "Identifier",
          "start": 19,
          "end": 22,
          "name": "baz"
        },
        "computed": false,
        "optional": false
      },
      "right": {
        "type": "Identifier",
        "start": 26,
        "end": 29,
        "name": "qux"
      },
      "body": {
        "type": "EmptyStatement",
        "start": 30,
        "end": 31
      }
    }
  ]
}
	`, ast)
}

func TestHarmony300(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_STRICT)
	ast, err := CompileWithOpts("function f() { var x; function x() {} }", opts)
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
        "start": 9,
        "end": 10,
        "name": "f"
      },
      "generator": false,
      "async": false,
      "params": [],
      "body": {
        "type": "BlockStatement",
        "start": 13,
        "end": 39,
        "body": [
          {
            "type": "VariableDeclaration",
            "start": 15,
            "end": 21,
            "declarations": [
              {
                "type": "VariableDeclarator",
                "start": 19,
                "end": 20,
                "id": {
                  "type": "Identifier",
                  "start": 19,
                  "end": 20,
                  "name": "x"
                },
                "init": null
              }
            ],
            "kind": "var"
          },
          {
            "type": "FunctionDeclaration",
            "start": 22,
            "end": 37,
            "id": {
              "type": "Identifier",
              "start": 31,
              "end": 32,
              "name": "x"
            },
            "generator": false,
            "async": false,
            "params": [],
            "body": {
              "type": "BlockStatement",
              "start": 35,
              "end": 37,
              "body": []
            }
          }
        ]
      }
    }
  ]
}
	`, ast)
}

func TestHarmony301(t *testing.T) {
	ast, err := Compile("a.of / 2")
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
          "type": "MemberExpression",
          "start": 0,
          "end": 4,
          "object": {
            "type": "Identifier",
            "start": 0,
            "end": 1,
            "name": "a"
          },
          "property": {
            "type": "Identifier",
            "start": 2,
            "end": 4,
            "name": "of"
          },
          "computed": false,
          "optional": false
        },
        "operator": "/",
        "right": {
          "type": "Literal",
          "start": 7,
          "end": 8,
          "value": 2,
          "raw": "2"
        }
      }
    }
  ]
}
	`, ast)
}

func TestHarmony302(t *testing.T) {
	ast, err := Compile("let x = 1; x = 2")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 16,
  "body": [
    {
      "type": "VariableDeclaration",
      "start": 0,
      "end": 10,
      "declarations": [
        {
          "type": "VariableDeclarator",
          "start": 4,
          "end": 9,
          "id": {
            "type": "Identifier",
            "start": 4,
            "end": 5,
            "name": "x"
          },
          "init": {
            "type": "Literal",
            "start": 8,
            "end": 9,
            "value": 1,
            "raw": "1"
          }
        }
      ],
      "kind": "let"
    },
    {
      "type": "ExpressionStatement",
      "start": 11,
      "end": 16,
      "expression": {
        "type": "AssignmentExpression",
        "start": 11,
        "end": 16,
        "operator": "=",
        "left": {
          "type": "Identifier",
          "start": 11,
          "end": 12,
          "name": "x"
        },
        "right": {
          "type": "Literal",
          "start": 15,
          "end": 16,
          "value": 2,
          "raw": "2"
        }
      }
    }
  ]
}
	`, ast)
}

func TestHarmony303(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_STRICT)
	ast, err := CompileWithOpts("function *f2() { () => yield / 1 }", opts)
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
        "start": 10,
        "end": 12,
        "name": "f2"
      },
      "generator": true,
      "async": false,
      "params": [],
      "body": {
        "type": "BlockStatement",
        "start": 15,
        "end": 34,
        "body": [
          {
            "type": "ExpressionStatement",
            "start": 17,
            "end": 32,
            "expression": {
              "type": "ArrowFunctionExpression",
              "start": 17,
              "end": 32,
              "id": null,
              "expression": true,
              "generator": false,
              "async": false,
              "params": [],
              "body": {
                "type": "BinaryExpression",
                "start": 23,
                "end": 32,
                "left": {
                  "type": "Identifier",
                  "start": 23,
                  "end": 28,
                  "name": "yield"
                },
                "operator": "/",
                "right": {
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

func TestHarmony304(t *testing.T) {
	ast, err := Compile("({ a = 42, b: c.d } = e)")
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
      "end": 24,
      "expression": {
        "type": "AssignmentExpression",
        "start": 1,
        "end": 23,
        "operator": "=",
        "left": {
          "type": "ObjectPattern",
          "start": 1,
          "end": 19,
          "properties": [
            {
              "type": "Property",
              "start": 3,
              "end": 9,
              "method": false,
              "shorthand": true,
              "computed": false,
              "key": {
                "type": "Identifier",
                "start": 3,
                "end": 4,
                "name": "a"
              },
              "kind": "init",
              "value": {
                "type": "AssignmentPattern",
                "start": 3,
                "end": 9,
                "left": {
                  "type": "Identifier",
                  "start": 3,
                  "end": 4,
                  "name": "a"
                },
                "right": {
                  "type": "Literal",
                  "start": 7,
                  "end": 9,
                  "value": 42,
                  "raw": "42"
                }
              }
            },
            {
              "type": "Property",
              "start": 11,
              "end": 17,
              "method": false,
              "shorthand": false,
              "computed": false,
              "key": {
                "type": "Identifier",
                "start": 11,
                "end": 12,
                "name": "b"
              },
              "value": {
                "type": "MemberExpression",
                "start": 14,
                "end": 17,
                "object": {
                  "type": "Identifier",
                  "start": 14,
                  "end": 15,
                  "name": "c"
                },
                "property": {
                  "type": "Identifier",
                  "start": 16,
                  "end": 17,
                  "name": "d"
                },
                "computed": false,
                "optional": false
              },
              "kind": "init"
            }
          ]
        },
        "right": {
          "type": "Identifier",
          "start": 22,
          "end": 23,
          "name": "e"
        }
      }
    }
  ]
}
	`, ast)
}

func TestHarmony305(t *testing.T) {
	ast, err := Compile("({ __proto__: x, __proto__: y, __proto__: z }) => {}")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 52,
  "body": [
    {
      "type": "ExpressionStatement",
      "start": 0,
      "end": 52,
      "expression": {
        "type": "ArrowFunctionExpression",
        "start": 0,
        "end": 52,
        "id": null,
        "expression": false,
        "generator": false,
        "async": false,
        "params": [
          {
            "type": "ObjectPattern",
            "start": 1,
            "end": 45,
            "properties": [
              {
                "type": "Property",
                "start": 3,
                "end": 15,
                "method": false,
                "shorthand": false,
                "computed": false,
                "key": {
                  "type": "Identifier",
                  "start": 3,
                  "end": 12,
                  "name": "__proto__"
                },
                "value": {
                  "type": "Identifier",
                  "start": 14,
                  "end": 15,
                  "name": "x"
                },
                "kind": "init"
              },
              {
                "type": "Property",
                "start": 17,
                "end": 29,
                "method": false,
                "shorthand": false,
                "computed": false,
                "key": {
                  "type": "Identifier",
                  "start": 17,
                  "end": 26,
                  "name": "__proto__"
                },
                "value": {
                  "type": "Identifier",
                  "start": 28,
                  "end": 29,
                  "name": "y"
                },
                "kind": "init"
              },
              {
                "type": "Property",
                "start": 31,
                "end": 43,
                "method": false,
                "shorthand": false,
                "computed": false,
                "key": {
                  "type": "Identifier",
                  "start": 31,
                  "end": 40,
                  "name": "__proto__"
                },
                "value": {
                  "type": "Identifier",
                  "start": 42,
                  "end": 43,
                  "name": "z"
                },
                "kind": "init"
              }
            ]
          }
        ],
        "body": {
          "type": "BlockStatement",
          "start": 50,
          "end": 52,
          "body": []
        }
      }
    }
  ]
}
	`, ast)
}

func TestHarmony306(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_STRICT)
	ast, err := CompileWithOpts("class x {}\n05", opts)
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "ClassDeclaration",
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
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "body": [],
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
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 10
        }
      }
    },
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "Literal",
        "value": 5,
        "raw": "05",
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

func TestHarmony307(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_STRICT)
	ast, err := CompileWithOpts("function x() { 'use strict' }\n05", opts)
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "body": [
    {
      "type": "FunctionDeclaration",
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
                  "column": 15
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
                "column": 15
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
            "column": 13
          },
          "end": {
            "line": 1,
            "column": 29
          }
        }
      },
      "generator": false,
      "async": false,
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
    {
      "type": "ExpressionStatement",
      "expression": {
        "type": "Literal",
        "value": 5,
        "raw": "05",
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

func TestHarmony308(t *testing.T) {
	ast, err := Compile("const myFn = ({ set = '' }) => {};")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 34,
  "body": [
    {
      "type": "VariableDeclaration",
      "start": 0,
      "end": 34,
      "declarations": [
        {
          "type": "VariableDeclarator",
          "start": 6,
          "end": 33,
          "id": {
            "type": "Identifier",
            "start": 6,
            "end": 10,
            "name": "myFn"
          },
          "init": {
            "type": "ArrowFunctionExpression",
            "start": 13,
            "end": 33,
            "id": null,
            "expression": false,
            "generator": false,
            "async": false,
            "params": [
              {
                "type": "ObjectPattern",
                "start": 14,
                "end": 26,
                "properties": [
                  {
                    "type": "Property",
                    "start": 16,
                    "end": 24,
                    "method": false,
                    "shorthand": true,
                    "computed": false,
                    "key": {
                      "type": "Identifier",
                      "start": 16,
                      "end": 19,
                      "name": "set"
                    },
                    "kind": "init",
                    "value": {
                      "type": "AssignmentPattern",
                      "start": 16,
                      "end": 24,
                      "left": {
                        "type": "Identifier",
                        "start": 16,
                        "end": 19,
                        "name": "set"
                      },
                      "right": {
                        "type": "Literal",
                        "start": 22,
                        "end": 24,
                        "value": "",
                        "raw": "''"
                      }
                    }
                  }
                ]
              }
            ],
            "body": {
              "type": "BlockStatement",
              "start": 31,
              "end": 33,
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

func TestHarmony309(t *testing.T) {
	ast, err := Compile("[[...[], 0].x] = []")
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
        "type": "AssignmentExpression",
        "start": 0,
        "end": 19,
        "operator": "=",
        "left": {
          "type": "ArrayPattern",
          "start": 0,
          "end": 14,
          "elements": [
            {
              "type": "MemberExpression",
              "start": 1,
              "end": 13,
              "object": {
                "type": "ArrayExpression",
                "start": 1,
                "end": 11,
                "elements": [
                  {
                    "type": "SpreadElement",
                    "start": 2,
                    "end": 7,
                    "argument": {
                      "type": "ArrayExpression",
                      "start": 5,
                      "end": 7,
                      "elements": []
                    }
                  },
                  {
                    "type": "Literal",
                    "start": 9,
                    "end": 10,
                    "value": 0,
                    "raw": "0"
                  }
                ]
              },
              "property": {
                "type": "Identifier",
                "start": 12,
                "end": 13,
                "name": "x"
              },
              "computed": false,
              "optional": false
            }
          ]
        },
        "right": {
          "type": "ArrayExpression",
          "start": 17,
          "end": 19,
          "elements": []
        }
      }
    }
  ]
}
	`, ast)
}

func TestHarmony310(t *testing.T) {
	ast, err := Compile("let \\u0061;")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 11,
  "body": [
    {
      "type": "VariableDeclaration",
      "start": 0,
      "end": 11,
      "declarations": [
        {
          "type": "VariableDeclarator",
          "start": 4,
          "end": 10,
          "id": {
            "type": "Identifier",
            "start": 4,
            "end": 10,
            "name": "a"
          },
          "init": null
        }
      ],
      "kind": "let"
    }
  ]
}
	`, ast)
}

func TestHarmony311(t *testing.T) {
	ast, err := Compile("let in\\u0061;")
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
          "end": 12,
          "id": {
            "type": "Identifier",
            "start": 4,
            "end": 12,
            "name": "ina"
          },
          "init": null
        }
      ],
      "kind": "let"
    }
  ]
}
	`, ast)
}

func TestHarmony312(t *testing.T) {
	ast, err := Compile("let in𝐬𝐭𝐚𝐧𝐜𝐞𝐨𝐟;")
	AssertEqual(t, nil, err, "should be prog ok")

	// `start` and `end` are utf16 based position of the source string
	// used in the javascript implemented parsers, however they are utf8
	// based position in mole
	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 39,
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
      "type": "VariableDeclaration",
      "start": 0,
      "end": 39,
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
      "declarations": [
        {
          "type": "VariableDeclarator",
          "start": 4,
          "end": 38,
          "loc": {
            "start": {
              "line": 1,
              "column": 4
            },
            "end": {
              "line": 1,
              "column": 14
            }
          },
          "id": {
            "type": "Identifier",
            "start": 4,
            "end": 38,
            "loc": {
              "start": {
                "line": 1,
                "column": 4
              },
              "end": {
                "line": 1,
                "column": 14
              }
            },
            "name": "in𝐬𝐭𝐚𝐧𝐜𝐞𝐨𝐟"
          },
          "init": null
        }
      ],
      "kind": "let"
    }
  ]
}
	`, ast)
}

func TestHarmony313(t *testing.T) {
	ast, err := Compile("let 𝐢𝐧;")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 13,
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
      "type": "VariableDeclaration",
      "start": 0,
      "end": 13,
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
      "kind": "let",
      "declarations": [
        {
          "type": "VariableDeclarator",
          "start": 4,
          "end": 12,
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
          "id": {
            "type": "Identifier",
            "start": 4,
            "end": 12,
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
            "name": "𝐢𝐧"
          },
          "init": null
        }
      ]
    }
  ]
}
	`, ast)
}

func TestHarmony314(t *testing.T) {
	ast, err := Compile("for ((a in b);;);")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 17,
  "body": [
    {
      "type": "ForStatement",
      "start": 0,
      "end": 17,
      "init": {
        "type": "BinaryExpression",
        "start": 6,
        "end": 12,
        "left": {
          "type": "Identifier",
          "start": 6,
          "end": 7,
          "name": "a"
        },
        "operator": "in",
        "right": {
          "type": "Identifier",
          "start": 11,
          "end": 12,
          "name": "b"
        }
      },
      "test": null,
      "update": null,
      "body": {
        "type": "EmptyStatement",
        "start": 16,
        "end": 17
      }
    }
  ]
}
	`, ast)
}

func TestHarmony315(t *testing.T) {
	ast, err := Compile("for (function (){ a in b };;);")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 30,
  "body": [
    {
      "type": "ForStatement",
      "start": 0,
      "end": 30,
      "init": {
        "type": "FunctionExpression",
        "start": 5,
        "end": 26,
        "id": null,
        "expression": false,
        "generator": false,
        "async": false,
        "params": [],
        "body": {
          "type": "BlockStatement",
          "start": 16,
          "end": 26,
          "body": [
            {
              "type": "ExpressionStatement",
              "start": 18,
              "end": 24,
              "expression": {
                "type": "BinaryExpression",
                "start": 18,
                "end": 24,
                "left": {
                  "type": "Identifier",
                  "start": 18,
                  "end": 19,
                  "name": "a"
                },
                "operator": "in",
                "right": {
                  "type": "Identifier",
                  "start": 23,
                  "end": 24,
                  "name": "b"
                }
              }
            }
          ]
        }
      },
      "test": null,
      "update": null,
      "body": {
        "type": "EmptyStatement",
        "start": 29,
        "end": 30
      }
    }
  ]
}
	`, ast)
}

func TestHarmony316(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_STRICT)
	ast, err := CompileWithOpts("let a = yield + 1", opts)
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
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
            "type": "Identifier",
            "start": 4,
            "end": 5,
            "name": "a"
          },
          "init": {
            "type": "BinaryExpression",
            "start": 8,
            "end": 17,
            "left": {
              "type": "Identifier",
              "start": 8,
              "end": 13,
              "name": "yield"
            },
            "operator": "+",
            "right": {
              "type": "Literal",
              "start": 16,
              "end": 17,
              "value": 1,
              "raw": "1"
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

func TestHarmony317(t *testing.T) {
	ast, err := Compile("try {} catch {}")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 15,
  "body": [
    {
      "type": "TryStatement",
      "start": 0,
      "end": 15,
      "block": {
        "type": "BlockStatement",
        "start": 4,
        "end": 6,
        "body": []
      },
      "handler": {
        "type": "CatchClause",
        "start": 7,
        "end": 15,
        "param": null,
        "body": {
          "type": "BlockStatement",
          "start": 13,
          "end": 15,
          "body": []
        }
      },
      "finalizer": null
    }
  ]
}
	`, ast)
}

func TestHarmony318(t *testing.T) {
	ast, err := Compile("function foo() { return; x = 1; }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 33,
  "body": [
    {
      "type": "FunctionDeclaration",
      "start": 0,
      "end": 33,
      "id": {
        "type": "Identifier",
        "start": 9,
        "end": 12,
        "name": "foo"
      },
      "params": [],
      "body": {
        "type": "BlockStatement",
        "start": 15,
        "end": 33,
        "body": [
          {
            "type": "ReturnStatement",
            "start": 17,
            "end": 24,
            "argument": null
          },
          {
            "type": "ExpressionStatement",
            "start": 25,
            "end": 31,
            "expression": {
              "type": "AssignmentExpression",
              "start": 25,
              "end": 30,
              "operator": "=",
              "left": {
                "type": "Identifier",
                "start": 25,
                "end": 26,
                "name": "x"
              },
              "right": {
                "type": "Literal",
                "start": 29,
                "end": 30,
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

func TestHarmony319(t *testing.T) {
	ast, err := Compile("class foo { get bar(){ ~function () { return true; }()}}")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 56,
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 56
    }
  },
  "body": [
    {
      "type": "ClassDeclaration",
      "start": 0,
      "end": 56,
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 56
        }
      },
      "id": {
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
        "name": "foo"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "start": 10,
        "end": 56,
        "loc": {
          "start": {
            "line": 1,
            "column": 10
          },
          "end": {
            "line": 1,
            "column": 56
          }
        },
        "body": [
          {
            "type": "MethodDefinition",
            "start": 12,
            "end": 55,
            "loc": {
              "start": {
                "line": 1,
                "column": 12
              },
              "end": {
                "line": 1,
                "column": 55
              }
            },
            "static": false,
            "key": {
              "type": "Identifier",
              "start": 16,
              "end": 19,
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
              "name": "bar"
            },
            "computed": false,
            "kind": "get",
            "id": null,
            "value": {
              "type": "FunctionExpression",
              "params": [],
              "body": {
                "type": "BlockStatement",
                "start": 21,
                "end": 55,
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 21
                  },
                  "end": {
                    "line": 1,
                    "column": 55
                  }
                },
                "body": [
                  {
                    "type": "ExpressionStatement",
                    "start": 23,
                    "end": 54,
                    "loc": {
                      "start": {
                        "line": 1,
                        "column": 23
                      },
                      "end": {
                        "line": 1,
                        "column": 54
                      }
                    },
                    "expression": {
                      "type": "UnaryExpression",
                      "start": 23,
                      "end": 54,
                      "loc": {
                        "start": {
                          "line": 1,
                          "column": 23
                        },
                        "end": {
                          "line": 1,
                          "column": 54
                        }
                      },
                      "operator": "~",
                      "prefix": true,
                      "argument": {
                        "type": "CallExpression",
                        "start": 24,
                        "end": 54,
                        "loc": {
                          "start": {
                            "line": 1,
                            "column": 24
                          },
                          "end": {
                            "line": 1,
                            "column": 54
                          }
                        },
                        "callee": {
                          "type": "FunctionExpression",
                          "start": 24,
                          "end": 52,
                          "loc": {
                            "start": {
                              "line": 1,
                              "column": 24
                            },
                            "end": {
                              "line": 1,
                              "column": 52
                            }
                          },
                          "id": null,
                          "generator": false,
                          "async": false,
                          "params": [],
                          "body": {
                            "type": "BlockStatement",
                            "start": 36,
                            "end": 52,
                            "loc": {
                              "start": {
                                "line": 1,
                                "column": 36
                              },
                              "end": {
                                "line": 1,
                                "column": 52
                              }
                            },
                            "body": [
                              {
                                "type": "ReturnStatement",
                                "start": 38,
                                "end": 50,
                                "loc": {
                                  "start": {
                                    "line": 1,
                                    "column": 38
                                  },
                                  "end": {
                                    "line": 1,
                                    "column": 50
                                  }
                                },
                                "argument": {
                                  "type": "Literal",
                                  "start": 45,
                                  "end": 49,
                                  "loc": {
                                    "start": {
                                      "line": 1,
                                      "column": 45
                                    },
                                    "end": {
                                      "line": 1,
                                      "column": 49
                                    }
                                  },
                                  "value": true
                                }
                              }
                            ],
                            "directives": []
                          }
                        },
                        "arguments": []
                      }
                    }
                  }
                ],
                "directives": []
              },
              "async": false,
              "generator": false
            }
          }
        ]
      }
    }
  ],
  "directives": []
}
	`, ast)
}
