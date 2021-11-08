package estree

import (
	"testing"

	"github.com/hsiaosiyuan0/mole/pkg/assert"
	"github.com/hsiaosiyuan0/mole/pkg/js/parser"
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
	ast, err := compile("`42`")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
	ast, err := compile("raw`42`")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
	ast, err := compile("raw`hello ${name}`")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
	ast, err := compile("`$`")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
	ast, err := compile("`\\n\\r\\b\\v\\t\\f\\\n\\\r\n\\\u2028\\\u2029`")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
	ast, err := compile("`\n\r\n\r`")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
	ast, err := compile("`\\u{000042}\\u0042\\x42u0\\A`")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
	ast, err := compile("new raw`42`")
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
	ast, err := compile("`outer${{x: {y: 10}}}bar${`nested${function(){return 1;}}endnest`}end`")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
	ast, err := compile("switch (answer) { case 42: let t = 42; break; }")
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
	ast, err := compile("() => \"test\"")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
	ast, err := compile("e => \"test\"")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
	ast, err := compile("(e) => \"test\"")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
	ast, err := compile("(a, b) => \"test\"")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
	ast, err := compile("e => { 42; }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
	ast, err := compile("e => ({ property: 42 })")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
	ast, err := compile("e => { label: 42 }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
	ast, err := compile("(a, b) => { 42; }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
	ast, err := compile("([a, , b]) => 42")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
	ast, err := compile("(() => {})()")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
	ast, err := compile("((() => {}))()")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
	ast, err := compile("(x=1) => x * x")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
	ast, err := compileWithOpts("eval => 42", opts)
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
	ast, err := compile("(a) => 00")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
	ast, err := compileWithOpts("(eval, a) => 42", opts)
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
	ast, err := compileWithOpts("(eval = 10) => 42", opts)
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
	ast, err := compileWithOpts("(eval, a = 10) => 42", opts)
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
	ast, err := compile("(x => x)")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
	ast, err := compile("x => y => 42")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
	ast, err := compile("(x) => ((y, z) => (x, y, z))")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
	ast, err := compile("foo(() => {})")
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
	ast, err := compile("foo((x, y) => {})")
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
	ast, err := compile("x = { method() { } }")
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
	ast, err := compile("x = { method(test) { } }")
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
	ast, err := compile("x = { 'method'() { } }")
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
	ast, err := compile("x = { get() { } }")
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
	ast, err := compile("x = { set() { } }")
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
	ast, err := compile("x = { method() { super.a(); } }")
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
	ast, err := compile("x = { y, z }")
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
	ast, err := compile("[a, b] = [b, a]")
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
	ast, err := compile("[a.r] = b")
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
	ast, err := compile("let [a,,b] = c")
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
	ast, err := compile("({ responseText: text } = res)")
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
	ast, err := compile("const {a} = {}")
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
	ast, err := compile("const [a] = []")
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
	ast, err := compile("let {a} = {}")
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
	ast, err := compile("let [a] = []")
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
	ast, err := compile("var {a} = {}")
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
	ast, err := compile("var [a] = []")
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
	ast, err := compile("const {a:b} = {}")
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
	ast, err := compile("let {a:b} = {}")
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
	ast, err := compile("var {a:b} = {}")
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
	ast, err := compile("export var document")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
	ast, err := compile("export var document = { }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
	ast, err := compile("export let document")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
	ast, err := compile("export let document = { }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
	ast, err := compile("export const document = { }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
	ast, err := compile("export function parse() { }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
	ast, err := compile("export class Class {}")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
	ast, err := compile("export default 42")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
	ast, err := compile("export default function () {}")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
	ast, err := compile("export default function f() {}")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
	ast, err := compile("export default class {}")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
	ast, err := compile("export default class A {}")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
	ast, err := compile("export default (class{});")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
	ast, err := compile("export * from \"crypto\"")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
	ast, err := compile("export { encrypt }\nvar encrypt")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
	ast, err := compile("function encrypt() {} let decrypt; export { encrypt, decrypt }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
	ast, err := compile("export default class Test {}; export { Test }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{}
	`, ast)
}

func TestHarmony87(t *testing.T) {
	ast, err := compile("{ var encrypt } export { encrypt }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{}
	`, ast)
}

func TestHarmony88(t *testing.T) {
	ast, err := compile("export { encrypt as default }; function* encrypt() {}")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
	ast, err := compile("export { encrypt, decrypt as dec }; let encrypt, decrypt")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
	ast, err := compile("export { default } from \"other\"")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
	ast, err := compile("import \"jquery\"")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
	ast, err := compile("import $ from \"jquery\"")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
	ast, err := compile("import { encrypt, decrypt } from \"crypto\"")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
	ast, err := compile("import { encrypt as enc } from \"crypto\"")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
	ast, err := compile("import crypto, { decrypt, encrypt as enc } from \"crypto\"")
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
	ast, err := compile("import { null as nil } from \"bar\"")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
	ast, err := compile("import * as crypto from \"crypto\"")
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
	ast, err := compile("(function* () { yield v })")
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
	ast, err := compile("(function* () { yield\nv })")
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
	ast, err := compile("(function* () { yield *v })")
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
	ast, err := compile("function* test () { yield *v }")
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
	ast, err := compile("var x = { *test () { yield *v } };")
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
	ast, err := compile("function* foo() { console.log(yield); }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
	ast, err := compile("function* t() {}")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
	ast, err := compile("(function* () { yield yield 10 })")
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
	ast, err := compile("for(x of list) process(x);")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
	ast, err := compile("for (var x of list) process(x);")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
	ast, err := compile("for (let x of list) process(x);")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
	ast, err := compile("for (let\n{x} of list) process(x);")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
	ast, err := compile("var A = class extends B {}")
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
	ast, err := compile("class A extends class B extends C {} {}")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
	ast, err := compile("class A {get() {}}")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
	ast, err := compile("class A { static get() {}}")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
	ast, err := compile("class A extends B {get foo() {}}")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
	ast, err := compile("class A extends B { static get foo() {}}")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
	ast, err := compile("class A {set a(v) {}}")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
	ast, err := compile("class A { static set a(v) {}}")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
	ast, err := compile("class A {set(v) {};}")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
	ast, err := compile("class A { static set(v) {};}")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
	ast, err := compile("class A {*gen(v) { yield v; }}")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
	ast, err := compile("class A { static *gen(v) { yield v; }}")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
	ast, err := compile("(class { *static() {} })")
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
	ast, err := compile("\"use strict\"; (class A extends B {constructor() { super() }})")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
	ast, err := compile("(class A extends B { constructor() { (() => { super() }); } })")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{}
	`, ast)
}

func TestHarmony125(t *testing.T) {
	ast, err := compile("class A {'constructor'() {}}")
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
	ast, err := compile("class A { get ['constructor']() {} }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
	ast, err := compile("class A {static foo() {}}")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
	ast, err := compile("class A {foo() {} static bar() {}}")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
	ast, err := compile("class A { foo() {} bar() {}}")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
	ast, err := compile("class A { static get foo() {} get foo() {}}")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
	ast, err := compile("class A { static get foo() {} static get bar() {} }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
	ast, err := compile("class A { static get foo() {} static set foo(v) {} get foo() {} set foo(v) {}}")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
	ast, err := compile("class A { static [foo]() {} }")
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
	ast, err := compile("class A { static get [foo]() {} }")
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
	ast, err := compile("class A { set foo(v) {} get foo() {} }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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
	ast, err := compile("class A { foo() {} get foo() {} }")
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
	ast, err := compile("class Semicolon { ; }")
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
	ast, err := compile("class a { static }")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
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

func TestHarmony301(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony302(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony303(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony304(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony305(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony306(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony307(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony308(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony309(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony310(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony311(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony312(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony313(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony314(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony315(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony316(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony317(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony318(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony319(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}

func TestHarmony320(t *testing.T) {
	// ast, err := compile("x = { false: 42 }")
	// assert.Equal(t, nil, err, "should be prog ok")

	// assert.EqualJson(t, `

	// `, ast)
}
