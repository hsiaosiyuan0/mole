package estree_test

import (
	"testing"

	. "github.com/hsiaosiyuan0/mole/ecma/estree/test"
	. "github.com/hsiaosiyuan0/mole/util"
)

func TestJSX1(t *testing.T) {
	ast, err := Compile("<a>{/* foo */}</a>")
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
        "type": "JSXElement",
        "start": 0,
        "end": 18,
        "openingElement": {
          "type": "JSXOpeningElement",
          "start": 0,
          "end": 3,
          "attributes": [],
          "name": {
            "type": "JSXIdentifier",
            "start": 1,
            "end": 2,
            "name": "a"
          },
          "selfClosing": false
        },
        "closingElement": {
          "type": "JSXClosingElement",
          "start": 14,
          "end": 18,
          "name": {
            "type": "JSXIdentifier",
            "start": 16,
            "end": 17,
            "name": "a"
          }
        },
        "children": [
          {
            "type": "JSXExpressionContainer",
            "start": 3,
            "end": 14,
            "expression": {
              "type": "JSXEmptyExpression",
              "start": 4,
              "end": 13
            }
          }
        ]
      }
    }
  ]
}
	`, ast)
}

func TestJSX2(t *testing.T) {
	ast, err := Compile("<A>foo{\"{\"}</A>")
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
        "type": "JSXElement",
        "start": 0,
        "end": 15,
        "openingElement": {
          "type": "JSXOpeningElement",
          "start": 0,
          "end": 3,
          "attributes": [],
          "name": {
            "type": "JSXIdentifier",
            "start": 1,
            "end": 2,
            "name": "A"
          },
          "selfClosing": false
        },
        "closingElement": {
          "type": "JSXClosingElement",
          "start": 11,
          "end": 15,
          "name": {
            "type": "JSXIdentifier",
            "start": 13,
            "end": 14,
            "name": "A"
          }
        },
        "children": [
          {
            "type": "JSXText",
            "start": 3,
            "end": 6,
            "value": "foo",
            "raw": "foo"
          },
          {
            "type": "JSXExpressionContainer",
            "start": 6,
            "end": 11,
            "expression": {
              "type": "Literal",
              "start": 7,
              "end": 10,
              "value": "{",
              "raw": "\"{\""
            }
          }
        ]
      }
    }
  ]
}
	`, ast)
}

func TestJSX3(t *testing.T) {
	ast, err := Compile("<A>foo{\"<\"}</A>")
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
        "type": "JSXElement",
        "start": 0,
        "end": 15,
        "openingElement": {
          "type": "JSXOpeningElement",
          "start": 0,
          "end": 3,
          "attributes": [],
          "name": {
            "type": "JSXIdentifier",
            "start": 1,
            "end": 2,
            "name": "A"
          },
          "selfClosing": false
        },
        "closingElement": {
          "type": "JSXClosingElement",
          "start": 11,
          "end": 15,
          "name": {
            "type": "JSXIdentifier",
            "start": 13,
            "end": 14,
            "name": "A"
          }
        },
        "children": [
          {
            "type": "JSXText",
            "start": 3,
            "end": 6,
            "value": "foo",
            "raw": "foo"
          },
          {
            "type": "JSXExpressionContainer",
            "start": 6,
            "end": 11,
            "expression": {
              "type": "Literal",
              "start": 7,
              "end": 10,
              "value": "<",
              "raw": "\"<\""
            }
          }
        ]
      }
    }
  ]
}
	`, ast)
}

func TestJSX4(t *testing.T) {
	ast, err := Compile("<A>foo{\"}\"}</A>")
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
        "type": "JSXElement",
        "start": 0,
        "end": 15,
        "openingElement": {
          "type": "JSXOpeningElement",
          "start": 0,
          "end": 3,
          "attributes": [],
          "name": {
            "type": "JSXIdentifier",
            "start": 1,
            "end": 2,
            "name": "A"
          },
          "selfClosing": false
        },
        "closingElement": {
          "type": "JSXClosingElement",
          "start": 11,
          "end": 15,
          "name": {
            "type": "JSXIdentifier",
            "start": 13,
            "end": 14,
            "name": "A"
          }
        },
        "children": [
          {
            "type": "JSXText",
            "start": 3,
            "end": 6,
            "value": "foo",
            "raw": "foo"
          },
          {
            "type": "JSXExpressionContainer",
            "start": 6,
            "end": 11,
            "expression": {
              "type": "Literal",
              "start": 7,
              "end": 10,
              "value": "}",
              "raw": "\"}\""
            }
          }
        ]
      }
    }
  ]
}
	`, ast)
}

func TestJSX5(t *testing.T) {
	ast, err := Compile("<A>foo&rbrace;</A>")
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
        "type": "JSXElement",
        "start": 0,
        "end": 18,
        "openingElement": {
          "type": "JSXOpeningElement",
          "start": 0,
          "end": 3,
          "attributes": [],
          "name": {
            "type": "JSXIdentifier",
            "start": 1,
            "end": 2,
            "name": "A"
          },
          "selfClosing": false
        },
        "closingElement": {
          "type": "JSXClosingElement",
          "start": 14,
          "end": 18,
          "name": {
            "type": "JSXIdentifier",
            "start": 16,
            "end": 17,
            "name": "A"
          }
        },
        "children": [
          {
            "type": "JSXText",
            "start": 3,
            "end": 14,
            "value": "foo}",
            "raw": "foo&rbrace;"
          }
        ]
      }
    }
  ]
}
	`, ast)
}

func TestJSX6(t *testing.T) {
	ast, err := Compile("<A>foo{\">\"}</A>")
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
        "type": "JSXElement",
        "start": 0,
        "end": 15,
        "openingElement": {
          "type": "JSXOpeningElement",
          "start": 0,
          "end": 3,
          "attributes": [],
          "name": {
            "type": "JSXIdentifier",
            "start": 1,
            "end": 2,
            "name": "A"
          },
          "selfClosing": false
        },
        "closingElement": {
          "type": "JSXClosingElement",
          "start": 11,
          "end": 15,
          "name": {
            "type": "JSXIdentifier",
            "start": 13,
            "end": 14,
            "name": "A"
          }
        },
        "children": [
          {
            "type": "JSXText",
            "start": 3,
            "end": 6,
            "value": "foo",
            "raw": "foo"
          },
          {
            "type": "JSXExpressionContainer",
            "start": 6,
            "end": 11,
            "expression": {
              "type": "Literal",
              "start": 7,
              "end": 10,
              "value": ">",
              "raw": "\">\""
            }
          }
        ]
      }
    }
  ]
}
	`, ast)
}

func TestJSX7(t *testing.T) {
	ast, err := Compile("<A>foo&gt;</A>")
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
        "type": "JSXElement",
        "start": 0,
        "end": 14,
        "openingElement": {
          "type": "JSXOpeningElement",
          "start": 0,
          "end": 3,
          "attributes": [],
          "name": {
            "type": "JSXIdentifier",
            "start": 1,
            "end": 2,
            "name": "A"
          },
          "selfClosing": false
        },
        "closingElement": {
          "type": "JSXClosingElement",
          "start": 10,
          "end": 14,
          "name": {
            "type": "JSXIdentifier",
            "start": 12,
            "end": 13,
            "name": "A"
          }
        },
        "children": [
          {
            "type": "JSXText",
            "start": 3,
            "end": 10,
            "value": "foo>",
            "raw": "foo&gt;"
          }
        ]
      }
    }
  ]
}
	`, ast)
}

func TestJSX8(t *testing.T) {
	ast, err := Compile("function*it(){yield <a></a>}")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 28,
  "body": [
    {
      "type": "FunctionDeclaration",
      "start": 0,
      "end": 28,
      "id": {
        "type": "Identifier",
        "start": 9,
        "end": 11,
        "name": "it"
      },
      "generator": true,
      "async": false,
      "params": [],
      "body": {
        "type": "BlockStatement",
        "start": 13,
        "end": 28,
        "body": [
          {
            "type": "ExpressionStatement",
            "start": 14,
            "end": 27,
            "expression": {
              "type": "YieldExpression",
              "start": 14,
              "end": 27,
              "delegate": false,
              "argument": {
                "type": "JSXElement",
                "start": 20,
                "end": 27,
                "openingElement": {
                  "type": "JSXOpeningElement",
                  "start": 20,
                  "end": 23,
                  "attributes": [],
                  "name": {
                    "type": "JSXIdentifier",
                    "start": 21,
                    "end": 22,
                    "name": "a"
                  },
                  "selfClosing": false
                },
                "closingElement": {
                  "type": "JSXClosingElement",
                  "start": 23,
                  "end": 27,
                  "name": {
                    "type": "JSXIdentifier",
                    "start": 25,
                    "end": 26,
                    "name": "a"
                  }
                },
                "children": []
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

func TestJSX9(t *testing.T) {
	ast, err := Compile("<A>foo}</A>")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 11,
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
      "start": 0,
      "end": 11,
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
        "type": "JSXElement",
        "start": 0,
        "end": 11,
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
        "openingElement": {
          "type": "JSXOpeningElement",
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
          "name": {
            "type": "JSXIdentifier",
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
            "name": "A"
          },
          "attributes": [],
          "selfClosing": false
        },
        "closingElement": {
          "type": "JSXClosingElement",
          "start": 7,
          "end": 11,
          "loc": {
            "start": {
              "line": 1,
              "column": 7
            },
            "end": {
              "line": 1,
              "column": 11
            }
          },
          "name": {
            "type": "JSXIdentifier",
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
            "name": "A"
          }
        },
        "children": [
          {
            "type": "JSXText",
            "start": 3,
            "end": 7,
            "loc": {
              "start": {
                "line": 1,
                "column": 3
              },
              "end": {
                "line": 1,
                "column": 7
              }
            },
            "raw": "foo}",
            "value": "foo}"
          }
        ]
      }
    }
  ]
}
	`, ast)
}

func TestJSX10(t *testing.T) {
	ast, err := Compile("<A>foo></A>")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 11,
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
      "start": 0,
      "end": 11,
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
        "type": "JSXElement",
        "start": 0,
        "end": 11,
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
        "openingElement": {
          "type": "JSXOpeningElement",
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
          "name": {
            "type": "JSXIdentifier",
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
            "name": "A"
          },
          "attributes": [],
          "selfClosing": false
        },
        "closingElement": {
          "type": "JSXClosingElement",
          "start": 7,
          "end": 11,
          "loc": {
            "start": {
              "line": 1,
              "column": 7
            },
            "end": {
              "line": 1,
              "column": 11
            }
          },
          "name": {
            "type": "JSXIdentifier",
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
            "name": "A"
          }
        },
        "children": [
          {
            "type": "JSXText",
            "start": 3,
            "end": 7,
            "loc": {
              "start": {
                "line": 1,
                "column": 3
              },
              "end": {
                "line": 1,
                "column": 7
              }
            },
            "raw": "foo>",
            "value": "foo>"
          }
        ]
      }
    }
  ]
}
	`, ast)
}

func TestJSX11(t *testing.T) {
	ast, err := Compile("<LeftRight left=<a /> right=<b>monkeys /> gorillas</b> />")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 57,
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 57
    }
  },
  "body": [
    {
      "type": "ExpressionStatement",
      "start": 0,
      "end": 57,
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 57
        }
      },
      "expression": {
        "type": "JSXElement",
        "start": 0,
        "end": 57,
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 57
          }
        },
        "openingElement": {
          "type": "JSXOpeningElement",
          "start": 0,
          "end": 57,
          "loc": {
            "start": {
              "line": 1,
              "column": 0
            },
            "end": {
              "line": 1,
              "column": 57
            }
          },
          "name": {
            "type": "JSXIdentifier",
            "start": 1,
            "end": 10,
            "loc": {
              "start": {
                "line": 1,
                "column": 1
              },
              "end": {
                "line": 1,
                "column": 10
              }
            },
            "name": "LeftRight"
          },
          "attributes": [
            {
              "type": "JSXAttribute",
              "start": 11,
              "end": 21,
              "loc": {
                "start": {
                  "line": 1,
                  "column": 11
                },
                "end": {
                  "line": 1,
                  "column": 21
                }
              },
              "name": {
                "type": "JSXIdentifier",
                "start": 11,
                "end": 15,
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 11
                  },
                  "end": {
                    "line": 1,
                    "column": 15
                  }
                },
                "name": "left"
              },
              "value": {
                "type": "JSXElement",
                "start": 16,
                "end": 21,
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
                "openingElement": {
                  "type": "JSXOpeningElement",
                  "start": 16,
                  "end": 21,
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
                  "name": {
                    "type": "JSXIdentifier",
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
                    "name": "a"
                  },
                  "attributes": [],
                  "selfClosing": true
                },
                "closingElement": null,
                "children": []
              }
            },
            {
              "type": "JSXAttribute",
              "start": 22,
              "end": 54,
              "loc": {
                "start": {
                  "line": 1,
                  "column": 22
                },
                "end": {
                  "line": 1,
                  "column": 54
                }
              },
              "name": {
                "type": "JSXIdentifier",
                "start": 22,
                "end": 27,
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
                "name": "right"
              },
              "value": {
                "type": "JSXElement",
                "start": 28,
                "end": 54,
                "loc": {
                  "start": {
                    "line": 1,
                    "column": 28
                  },
                  "end": {
                    "line": 1,
                    "column": 54
                  }
                },
                "openingElement": {
                  "type": "JSXOpeningElement",
                  "start": 28,
                  "end": 31,
                  "loc": {
                    "start": {
                      "line": 1,
                      "column": 28
                    },
                    "end": {
                      "line": 1,
                      "column": 31
                    }
                  },
                  "name": {
                    "type": "JSXIdentifier",
                    "start": 29,
                    "end": 30,
                    "loc": {
                      "start": {
                        "line": 1,
                        "column": 29
                      },
                      "end": {
                        "line": 1,
                        "column": 30
                      }
                    },
                    "name": "b"
                  },
                  "attributes": [],
                  "selfClosing": false
                },
                "closingElement": {
                  "type": "JSXClosingElement",
                  "start": 50,
                  "end": 54,
                  "loc": {
                    "start": {
                      "line": 1,
                      "column": 50
                    },
                    "end": {
                      "line": 1,
                      "column": 54
                    }
                  },
                  "name": {
                    "type": "JSXIdentifier",
                    "start": 52,
                    "end": 53,
                    "loc": {
                      "start": {
                        "line": 1,
                        "column": 52
                      },
                      "end": {
                        "line": 1,
                        "column": 53
                      }
                    },
                    "name": "b"
                  }
                },
                "children": [
                  {
                    "type": "JSXText",
                    "start": 31,
                    "end": 50,
                    "loc": {
                      "start": {
                        "line": 1,
                        "column": 31
                      },
                      "end": {
                        "line": 1,
                        "column": 50
                      }
                    },
                    "raw": "monkeys /> gorillas",
                    "value": "monkeys /> gorillas"
                  }
                ]
              }
            }
          ],
          "selfClosing": true
        },
        "closingElement": null,
        "children": []
      }
    }
  ]
}
	`, ast)
}

func TestJSX12(t *testing.T) {
	ast, err := Compile("<>{tips.map((tip, i) => <div key={i}>{`Tip ${i}:` + tip}</div>)}</>")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 67,
  "body": [
    {
      "type": "ExpressionStatement",
      "start": 0,
      "end": 67,
      "expression": {
        "type": "JSXFragment",
        "start": 0,
        "end": 67,
        "openingFragment": {
          "type": "JSXOpeningFragment",
          "start": 0,
          "end": 2,
          "attributes": []
        },
        "closingFragment": {
          "type": "JSXClosingFragment",
          "start": 64,
          "end": 67
        },
        "children": [
          {
            "type": "JSXExpressionContainer",
            "start": 2,
            "end": 64,
            "expression": {
              "type": "CallExpression",
              "start": 3,
              "end": 63,
              "callee": {
                "type": "MemberExpression",
                "start": 3,
                "end": 11,
                "object": {
                  "type": "Identifier",
                  "start": 3,
                  "end": 7,
                  "name": "tips"
                },
                "property": {
                  "type": "Identifier",
                  "start": 8,
                  "end": 11,
                  "name": "map"
                },
                "computed": false,
                "optional": false
              },
              "arguments": [
                {
                  "type": "ArrowFunctionExpression",
                  "start": 12,
                  "end": 62,
                  "id": null,
                  "expression": true,
                  "generator": false,
                  "async": false,
                  "params": [
                    {
                      "type": "Identifier",
                      "start": 13,
                      "end": 16,
                      "name": "tip"
                    },
                    {
                      "type": "Identifier",
                      "start": 18,
                      "end": 19,
                      "name": "i"
                    }
                  ],
                  "body": {
                    "type": "JSXElement",
                    "start": 24,
                    "end": 62,
                    "openingElement": {
                      "type": "JSXOpeningElement",
                      "start": 24,
                      "end": 37,
                      "attributes": [
                        {
                          "type": "JSXAttribute",
                          "start": 29,
                          "end": 36,
                          "name": {
                            "type": "JSXIdentifier",
                            "start": 29,
                            "end": 32,
                            "name": "key"
                          },
                          "value": {
                            "type": "JSXExpressionContainer",
                            "start": 33,
                            "end": 36,
                            "expression": {
                              "type": "Identifier",
                              "start": 34,
                              "end": 35,
                              "name": "i"
                            }
                          }
                        }
                      ],
                      "name": {
                        "type": "JSXIdentifier",
                        "start": 25,
                        "end": 28,
                        "name": "div"
                      },
                      "selfClosing": false
                    },
                    "closingElement": {
                      "type": "JSXClosingElement",
                      "start": 56,
                      "end": 62,
                      "name": {
                        "type": "JSXIdentifier",
                        "start": 58,
                        "end": 61,
                        "name": "div"
                      }
                    },
                    "children": [
                      {
                        "type": "JSXExpressionContainer",
                        "start": 37,
                        "end": 56,
                        "expression": {
                          "type": "BinaryExpression",
                          "start": 38,
                          "end": 55,
                          "left": {
                            "type": "TemplateLiteral",
                            "start": 38,
                            "end": 49,
                            "expressions": [
                              {
                                "type": "Identifier",
                                "start": 45,
                                "end": 46,
                                "name": "i"
                              }
                            ],
                            "quasis": [
                              {
                                "type": "TemplateElement",
                                "start": 39,
                                "end": 43,
                                "value": {
                                  "raw": "Tip ",
                                  "cooked": "Tip "
                                },
                                "tail": false
                              },
                              {
                                "type": "TemplateElement",
                                "start": 47,
                                "end": 48,
                                "value": {
                                  "raw": ":",
                                  "cooked": ":"
                                },
                                "tail": true
                              }
                            ]
                          },
                          "operator": "+",
                          "right": {
                            "type": "Identifier",
                            "start": 52,
                            "end": 55,
                            "name": "tip"
                          }
                        }
                      }
                    ]
                  }
                }
              ],
              "optional": false
            }
          }
        ]
      }
    }
  ]
}
	`, ast)
}

func TestJSX13(t *testing.T) {
	ast, err := Compile(`<div>
  <h1>Hello</h1>
</div>;`)
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
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
      "line": 3,
      "column": 7
    }
  },
  "body": [
    {
      "type": "ExpressionStatement",
      "start": 0,
      "end": 30,
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 3,
          "column": 7
        }
      },
      "expression": {
        "type": "JSXElement",
        "start": 0,
        "end": 29,
        "loc": {
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 3,
            "column": 6
          }
        },
        "openingElement": {
          "type": "JSXOpeningElement",
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
          "name": {
            "type": "JSXIdentifier",
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
            "name": "div"
          },
          "attributes": [],
          "selfClosing": false
        },
        "closingElement": {
          "type": "JSXClosingElement",
          "start": 23,
          "end": 29,
          "loc": {
            "start": {
              "line": 3,
              "column": 0
            },
            "end": {
              "line": 3,
              "column": 6
            }
          },
          "name": {
            "type": "JSXIdentifier",
            "start": 25,
            "end": 28,
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
            "name": "div"
          }
        },
        "children": [
          {
            "type": "JSXText",
            "start": 5,
            "end": 8,
            "loc": {
              "start": {
                "line": 1,
                "column": 5
              },
              "end": {
                "line": 2,
                "column": 2
              }
            },
            "value": "\n  "
          },
          {
            "type": "JSXElement",
            "start": 8,
            "end": 22,
            "loc": {
              "start": {
                "line": 2,
                "column": 2
              },
              "end": {
                "line": 2,
                "column": 16
              }
            },
            "openingElement": {
              "type": "JSXOpeningElement",
              "start": 8,
              "end": 12,
              "loc": {
                "start": {
                  "line": 2,
                  "column": 2
                },
                "end": {
                  "line": 2,
                  "column": 6
                }
              },
              "name": {
                "type": "JSXIdentifier",
                "start": 9,
                "end": 11,
                "loc": {
                  "start": {
                    "line": 2,
                    "column": 3
                  },
                  "end": {
                    "line": 2,
                    "column": 5
                  }
                },
                "name": "h1"
              },
              "attributes": [],
              "selfClosing": false
            },
            "closingElement": {
              "type": "JSXClosingElement",
              "start": 17,
              "end": 22,
              "loc": {
                "start": {
                  "line": 2,
                  "column": 11
                },
                "end": {
                  "line": 2,
                  "column": 16
                }
              },
              "name": {
                "type": "JSXIdentifier",
                "start": 19,
                "end": 21,
                "loc": {
                  "start": {
                    "line": 2,
                    "column": 13
                  },
                  "end": {
                    "line": 2,
                    "column": 15
                  }
                },
                "name": "h1"
              }
            },
            "children": [
              {
                "type": "JSXText",
                "start": 12,
                "end": 17,
                "loc": {
                  "start": {
                    "line": 2,
                    "column": 6
                  },
                  "end": {
                    "line": 2,
                    "column": 11
                  }
                },
                "value": "Hello"
              }
            ]
          },
          {
            "type": "JSXText",
            "start": 22,
            "end": 23,
            "loc": {
              "start": {
                "line": 2,
                "column": 16
              },
              "end": {
                "line": 3,
                "column": 0
              }
            },
            "value": "\n"
          }
        ]
      }
    }
  ]
}
	`, ast)
}
