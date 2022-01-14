package estree_test

import (
	"testing"

	. "github.com/hsiaosiyuan0/mole/ecma/estree/test"
	"github.com/hsiaosiyuan0/mole/ecma/parser"
	. "github.com/hsiaosiyuan0/mole/fuzz"
)

// Public Class Field
func TestClassFeature1(t *testing.T) {
	ast, err := Compile("class C { aaa }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 15,
  "body": [
    {
      "type": "ClassDeclaration",
      "start": 0,
      "end": 15,
      "id": {
        "type": "Identifier",
        "start": 6,
        "end": 7,
        "name": "C"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "start": 8,
        "end": 15,
        "body": [
          {
            "type": "PropertyDefinition",
            "start": 10,
            "end": 13,
            "static": false,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 10,
              "end": 13,
              "name": "aaa"
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

func TestClassFeature2(t *testing.T) {
	ast, err := Compile("class C { aaa; }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 16,
  "body": [
    {
      "type": "ClassDeclaration",
      "start": 0,
      "end": 16,
      "id": {
        "type": "Identifier",
        "start": 6,
        "end": 7,
        "name": "C"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "start": 8,
        "end": 16,
        "body": [
          {
            "type": "PropertyDefinition",
            "start": 10,
            "end": 14,
            "static": false,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 10,
              "end": 13,
              "name": "aaa"
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

func TestClassFeature3(t *testing.T) {
	ast, err := Compile("class C { \\u0041 }")
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
        "name": "C"
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
              "name": "A"
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

func TestClassFeature4(t *testing.T) {
	ast, err := Compile("class C { '0' }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 15,
  "body": [
    {
      "type": "ClassDeclaration",
      "start": 0,
      "end": 15,
      "id": {
        "type": "Identifier",
        "start": 6,
        "end": 7,
        "name": "C"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "start": 8,
        "end": 15,
        "body": [
          {
            "type": "PropertyDefinition",
            "start": 10,
            "end": 13,
            "static": false,
            "computed": false,
            "key": {
              "type": "Literal",
              "start": 10,
              "end": 13,
              "value": "0",
              "raw": "'0'"
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

func TestClassFeature5(t *testing.T) {
	ast, err := Compile("class C { 1e2 }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 15,
  "body": [
    {
      "type": "ClassDeclaration",
      "start": 0,
      "end": 15,
      "id": {
        "type": "Identifier",
        "start": 6,
        "end": 7,
        "name": "C"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "start": 8,
        "end": 15,
        "body": [
          {
            "type": "PropertyDefinition",
            "start": 10,
            "end": 13,
            "static": false,
            "computed": false,
            "key": {
              "type": "Literal",
              "start": 10,
              "end": 13,
              "value": 100,
              "raw": "1e2"
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

func TestClassFeature6(t *testing.T) {
	ast, err := Compile("class C { [0] }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 15,
  "body": [
    {
      "type": "ClassDeclaration",
      "start": 0,
      "end": 15,
      "id": {
        "type": "Identifier",
        "start": 6,
        "end": 7,
        "name": "C"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "start": 8,
        "end": 15,
        "body": [
          {
            "type": "PropertyDefinition",
            "start": 10,
            "end": 13,
            "static": false,
            "computed": true,
            "key": {
              "type": "Literal",
              "start": 11,
              "end": 12,
              "value": 0,
              "raw": "0"
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

func TestClassFeature7(t *testing.T) {
	ast, err := Compile("class C { aaa = bbb }")
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
        "name": "C"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "start": 8,
        "end": 21,
        "body": [
          {
            "type": "PropertyDefinition",
            "start": 10,
            "end": 19,
            "static": false,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 10,
              "end": 13,
              "name": "aaa"
            },
            "value": {
              "type": "Identifier",
              "start": 16,
              "end": 19,
              "name": "bbb"
            }
          }
        ]
      }
    }
  ]
}
`, ast)
}

func TestClassFeature8(t *testing.T) {
	ast, err := Compile("class C { aaa = bbb; }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 22,
  "body": [
    {
      "type": "ClassDeclaration",
      "start": 0,
      "end": 22,
      "id": {
        "type": "Identifier",
        "start": 6,
        "end": 7,
        "name": "C"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "start": 8,
        "end": 22,
        "body": [
          {
            "type": "PropertyDefinition",
            "start": 10,
            "end": 20,
            "static": false,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 10,
              "end": 13,
              "name": "aaa"
            },
            "value": {
              "type": "Identifier",
              "start": 16,
              "end": 19,
              "name": "bbb"
            }
          }
        ]
      }
    }
  ]
}
`, ast)
}

func TestClassFeature9(t *testing.T) {
	ast, err := Compile("class C { aaa; bbb }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 20,
  "body": [
    {
      "type": "ClassDeclaration",
      "start": 0,
      "end": 20,
      "id": {
        "type": "Identifier",
        "start": 6,
        "end": 7,
        "name": "C"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "start": 8,
        "end": 20,
        "body": [
          {
            "type": "PropertyDefinition",
            "start": 10,
            "end": 14,
            "static": false,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 10,
              "end": 13,
              "name": "aaa"
            },
            "value": null
          },
          {
            "type": "PropertyDefinition",
            "start": 15,
            "end": 18,
            "static": false,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 15,
              "end": 18,
              "name": "bbb"
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

func TestClassFeature10(t *testing.T) {
	ast, err := Compile("class C { aaa\nbbb }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 19,
  "body": [
    {
      "type": "ClassDeclaration",
      "start": 0,
      "end": 19,
      "id": {
        "type": "Identifier",
        "start": 6,
        "end": 7,
        "name": "C"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "start": 8,
        "end": 19,
        "body": [
          {
            "type": "PropertyDefinition",
            "start": 10,
            "end": 13,
            "static": false,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 10,
              "end": 13,
              "name": "aaa"
            },
            "value": null
          },
          {
            "type": "PropertyDefinition",
            "start": 14,
            "end": 17,
            "static": false,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 14,
              "end": 17,
              "name": "bbb"
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

func TestClassFeature11(t *testing.T) {
	ast, err := Compile("class C { aaa\n bbb }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 20,
  "body": [
    {
      "type": "ClassDeclaration",
      "start": 0,
      "end": 20,
      "id": {
        "type": "Identifier",
        "start": 6,
        "end": 7,
        "name": "C"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "start": 8,
        "end": 20,
        "body": [
          {
            "type": "PropertyDefinition",
            "start": 10,
            "end": 13,
            "static": false,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 10,
              "end": 13,
              "name": "aaa"
            },
            "value": null
          },
          {
            "type": "PropertyDefinition",
            "start": 15,
            "end": 18,
            "static": false,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 15,
              "end": 18,
              "name": "bbb"
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

func TestClassFeature12(t *testing.T) {
	ast, err := Compile("class C { aaa=()=>0 }")
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
        "name": "C"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "start": 8,
        "end": 21,
        "body": [
          {
            "type": "PropertyDefinition",
            "start": 10,
            "end": 19,
            "static": false,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 10,
              "end": 13,
              "name": "aaa"
            },
            "value": {
              "type": "ArrowFunctionExpression",
              "start": 14,
              "end": 19,
              "id": null,
              "expression": true,
              "generator": false,
              "async": false,
              "params": [],
              "body": {
                "type": "Literal",
                "start": 18,
                "end": 19,
                "value": 0,
                "raw": "0"
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

func TestClassFeature13(t *testing.T) {
	ast, err := Compile("class C { static aaa }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 22,
  "body": [
    {
      "type": "ClassDeclaration",
      "start": 0,
      "end": 22,
      "id": {
        "type": "Identifier",
        "start": 6,
        "end": 7,
        "name": "C"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "start": 8,
        "end": 22,
        "body": [
          {
            "type": "PropertyDefinition",
            "start": 10,
            "end": 20,
            "static": true,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 17,
              "end": 20,
              "name": "aaa"
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

func TestClassFeature14(t *testing.T) {
	ast, err := Compile("class C { static \\u0041 }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 25,
  "body": [
    {
      "type": "ClassDeclaration",
      "start": 0,
      "end": 25,
      "id": {
        "type": "Identifier",
        "start": 6,
        "end": 7,
        "name": "C"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "start": 8,
        "end": 25,
        "body": [
          {
            "type": "PropertyDefinition",
            "start": 10,
            "end": 23,
            "static": true,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 17,
              "end": 23,
              "name": "A"
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

func TestClassFeature15(t *testing.T) {
	ast, err := Compile("class C { static '0' }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 22,
  "body": [
    {
      "type": "ClassDeclaration",
      "start": 0,
      "end": 22,
      "id": {
        "type": "Identifier",
        "start": 6,
        "end": 7,
        "name": "C"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "start": 8,
        "end": 22,
        "body": [
          {
            "type": "PropertyDefinition",
            "start": 10,
            "end": 20,
            "static": true,
            "computed": false,
            "key": {
              "type": "Literal",
              "start": 17,
              "end": 20,
              "value": "0",
              "raw": "'0'"
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

func TestClassFeature16(t *testing.T) {
	ast, err := Compile("class C { static 1e2 }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 22,
  "body": [
    {
      "type": "ClassDeclaration",
      "start": 0,
      "end": 22,
      "id": {
        "type": "Identifier",
        "start": 6,
        "end": 7,
        "name": "C"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "start": 8,
        "end": 22,
        "body": [
          {
            "type": "PropertyDefinition",
            "start": 10,
            "end": 20,
            "static": true,
            "computed": false,
            "key": {
              "type": "Literal",
              "start": 17,
              "end": 20,
              "value": 100,
              "raw": "1e2"
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

func TestClassFeature17(t *testing.T) {
	ast, err := Compile("class C { static [0] }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 22,
  "body": [
    {
      "type": "ClassDeclaration",
      "start": 0,
      "end": 22,
      "id": {
        "type": "Identifier",
        "start": 6,
        "end": 7,
        "name": "C"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "start": 8,
        "end": 22,
        "body": [
          {
            "type": "PropertyDefinition",
            "start": 10,
            "end": 20,
            "static": true,
            "computed": true,
            "key": {
              "type": "Literal",
              "start": 18,
              "end": 19,
              "value": 0,
              "raw": "0"
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

func TestClassFeature18(t *testing.T) {
	ast, err := Compile("class C { static aaa = bbb }")
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
        "name": "C"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "start": 8,
        "end": 28,
        "body": [
          {
            "type": "PropertyDefinition",
            "start": 10,
            "end": 26,
            "static": true,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 17,
              "end": 20,
              "name": "aaa"
            },
            "value": {
              "type": "Identifier",
              "start": 23,
              "end": 26,
              "name": "bbb"
            }
          }
        ]
      }
    }
  ]
}
`, ast)
}

func TestClassFeature19(t *testing.T) {
	ast, err := Compile("class C { aaa; aaa }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 20,
  "body": [
    {
      "type": "ClassDeclaration",
      "start": 0,
      "end": 20,
      "id": {
        "type": "Identifier",
        "start": 6,
        "end": 7,
        "name": "C"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "start": 8,
        "end": 20,
        "body": [
          {
            "type": "PropertyDefinition",
            "start": 10,
            "end": 14,
            "static": false,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 10,
              "end": 13,
              "name": "aaa"
            },
            "value": null
          },
          {
            "type": "PropertyDefinition",
            "start": 15,
            "end": 18,
            "static": false,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 15,
              "end": 18,
              "name": "aaa"
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

func TestClassFeature20(t *testing.T) {
	ast, err := Compile("class C { aaa=0\n['bbb'] }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 25,
  "body": [
    {
      "type": "ClassDeclaration",
      "start": 0,
      "end": 25,
      "id": {
        "type": "Identifier",
        "start": 6,
        "end": 7,
        "name": "C"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "start": 8,
        "end": 25,
        "body": [
          {
            "type": "PropertyDefinition",
            "start": 10,
            "end": 23,
            "static": false,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 10,
              "end": 13,
              "name": "aaa"
            },
            "value": {
              "type": "MemberExpression",
              "start": 14,
              "end": 23,
              "object": {
                "type": "Literal",
                "start": 14,
                "end": 15,
                "value": 0,
                "raw": "0"
              },
              "property": {
                "type": "Literal",
                "start": 17,
                "end": 22,
                "value": "bbb",
                "raw": "'bbb'"
              },
              "computed": true,
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

func TestClassFeature21(t *testing.T) {
	ast, err := Compile("class C { get; set; static; async }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 35,
  "body": [
    {
      "type": "ClassDeclaration",
      "start": 0,
      "end": 35,
      "id": {
        "type": "Identifier",
        "start": 6,
        "end": 7,
        "name": "C"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "start": 8,
        "end": 35,
        "body": [
          {
            "type": "PropertyDefinition",
            "start": 10,
            "end": 14,
            "static": false,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 10,
              "end": 13,
              "name": "get"
            },
            "value": null
          },
          {
            "type": "PropertyDefinition",
            "start": 15,
            "end": 19,
            "static": false,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 15,
              "end": 18,
              "name": "set"
            },
            "value": null
          },
          {
            "type": "PropertyDefinition",
            "start": 20,
            "end": 27,
            "static": false,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 20,
              "end": 26,
              "name": "static"
            },
            "value": null
          },
          {
            "type": "PropertyDefinition",
            "start": 28,
            "end": 33,
            "static": false,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 28,
              "end": 33,
              "name": "async"
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

func TestClassFeature22(t *testing.T) {
	ast, err := Compile("class C { static\nget\nfoo(){} }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 30,
  "body": [
    {
      "type": "ClassDeclaration",
      "start": 0,
      "end": 30,
      "id": {
        "type": "Identifier",
        "start": 6,
        "end": 7,
        "name": "C"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "start": 8,
        "end": 30,
        "body": [
          {
            "type": "MethodDefinition",
            "start": 10,
            "end": 28,
            "static": true,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 21,
              "end": 24,
              "name": "foo"
            },
            "kind": "get",
            "value": {
              "type": "FunctionExpression",
              "start": 24,
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
          }
        ]
      }
    }
  ]
}
`, ast)
}

func TestClassFeature23(t *testing.T) {
	ast, err := Compile("class C { async\n get(){} }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 26,
  "body": [
    {
      "type": "ClassDeclaration",
      "start": 0,
      "end": 26,
      "id": {
        "type": "Identifier",
        "start": 6,
        "end": 7,
        "name": "C"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "start": 8,
        "end": 26,
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
              "end": 15,
              "name": "async"
            },
            "value": null
          },
          {
            "type": "MethodDefinition",
            "start": 17,
            "end": 24,
            "static": false,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 17,
              "end": 20,
              "name": "get"
            },
            "kind": "method",
            "value": {
              "type": "FunctionExpression",
              "start": 20,
              "end": 24,
              "id": null,
              "expression": false,
              "generator": false,
              "async": false,
              "params": [],
              "body": {
                "type": "BlockStatement",
                "start": 22,
                "end": 24,
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

func TestClassFeature24(t *testing.T) {
	ast, err := Compile("class C { get\n *foo(){} }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 25,
  "body": [
    {
      "type": "ClassDeclaration",
      "start": 0,
      "end": 25,
      "id": {
        "type": "Identifier",
        "start": 6,
        "end": 7,
        "name": "C"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "start": 8,
        "end": 25,
        "body": [
          {
            "type": "PropertyDefinition",
            "start": 10,
            "end": 13,
            "static": false,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 10,
              "end": 13,
              "name": "get"
            },
            "value": null
          },
          {
            "type": "MethodDefinition",
            "start": 15,
            "end": 23,
            "static": false,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 16,
              "end": 19,
              "name": "foo"
            },
            "kind": "method",
            "value": {
              "type": "FunctionExpression",
              "start": 19,
              "end": 23,
              "id": null,
              "expression": false,
              "generator": true,
              "async": false,
              "params": [],
              "body": {
                "type": "BlockStatement",
                "start": 21,
                "end": 23,
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

func TestClassFeature25(t *testing.T) {
	ast, err := Compile("class C { set\n *foo(){} }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 25,
  "body": [
    {
      "type": "ClassDeclaration",
      "start": 0,
      "end": 25,
      "id": {
        "type": "Identifier",
        "start": 6,
        "end": 7,
        "name": "C"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "start": 8,
        "end": 25,
        "body": [
          {
            "type": "PropertyDefinition",
            "start": 10,
            "end": 13,
            "static": false,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 10,
              "end": 13,
              "name": "set"
            },
            "value": null
          },
          {
            "type": "MethodDefinition",
            "start": 15,
            "end": 23,
            "static": false,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 16,
              "end": 19,
              "name": "foo"
            },
            "kind": "method",
            "value": {
              "type": "FunctionExpression",
              "start": 19,
              "end": 23,
              "id": null,
              "expression": false,
              "generator": true,
              "async": false,
              "params": [],
              "body": {
                "type": "BlockStatement",
                "start": 21,
                "end": 23,
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

func TestClassFeature26(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_MODULE)
	ast, err := CompileWithOpts("async function f() { class C { aaa = await } }", opts)
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
        "start": 15,
        "end": 16,
        "name": "f"
      },
      "generator": false,
      "async": true,
      "params": [],
      "body": {
        "type": "BlockStatement",
        "start": 19,
        "end": 46,
        "body": [
          {
            "type": "ClassDeclaration",
            "start": 21,
            "end": 44,
            "id": {
              "type": "Identifier",
              "start": 27,
              "end": 28,
              "name": "C"
            },
            "superClass": null,
            "body": {
              "type": "ClassBody",
              "start": 29,
              "end": 44,
              "body": [
                {
                  "type": "PropertyDefinition",
                  "start": 31,
                  "end": 42,
                  "static": false,
                  "computed": false,
                  "key": {
                    "type": "Identifier",
                    "start": 31,
                    "end": 34,
                    "name": "aaa"
                  },
                  "value": {
                    "type": "Identifier",
                    "start": 37,
                    "end": 42,
                    "name": "await"
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

func TestClassFeature27(t *testing.T) {
	ast, err := Compile("class C { a = new.target }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 26,
  "body": [
    {
      "type": "ClassDeclaration",
      "start": 0,
      "end": 26,
      "id": {
        "type": "Identifier",
        "start": 6,
        "end": 7,
        "name": "C"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "start": 8,
        "end": 26,
        "body": [
          {
            "type": "PropertyDefinition",
            "start": 10,
            "end": 24,
            "static": false,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 10,
              "end": 11,
              "name": "a"
            },
            "value": {
              "type": "MetaProperty",
              "start": 14,
              "end": 24,
              "meta": {
                "type": "Identifier",
                "start": 14,
                "end": 17,
                "name": "new"
              },
              "property": {
                "type": "Identifier",
                "start": 18,
                "end": 24,
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

func TestClassFeature28(t *testing.T) {
	ast, err := Compile("class C { aaa = super.bbb }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
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
        "name": "C"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "start": 8,
        "end": 27,
        "body": [
          {
            "type": "PropertyDefinition",
            "start": 10,
            "end": 25,
            "static": false,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 10,
              "end": 13,
              "name": "aaa"
            },
            "value": {
              "type": "MemberExpression",
              "start": 16,
              "end": 25,
              "object": {
                "type": "Super",
                "start": 16,
                "end": 21
              },
              "property": {
                "type": "Identifier",
                "start": 22,
                "end": 25,
                "name": "bbb"
              },
              "computed": false,
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

func TestClassFeature29(t *testing.T) {
	ast, err := Compile("class C { aaa = () => super.bbb }")
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
        "name": "C"
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
            "end": 31,
            "static": false,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 10,
              "end": 13,
              "name": "aaa"
            },
            "value": {
              "type": "ArrowFunctionExpression",
              "start": 16,
              "end": 31,
              "id": null,
              "expression": true,
              "generator": false,
              "async": false,
              "params": [],
              "body": {
                "type": "MemberExpression",
                "start": 22,
                "end": 31,
                "object": {
                  "type": "Super",
                  "start": 22,
                  "end": 27
                },
                "property": {
                  "type": "Identifier",
                  "start": 28,
                  "end": 31,
                  "name": "bbb"
                },
                "computed": false,
                "optional": false
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

func TestClassFeature30(t *testing.T) {
	ast, err := Compile("class C { aaa = () => () => super.bbb }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 39,
  "body": [
    {
      "type": "ClassDeclaration",
      "start": 0,
      "end": 39,
      "id": {
        "type": "Identifier",
        "start": 6,
        "end": 7,
        "name": "C"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "start": 8,
        "end": 39,
        "body": [
          {
            "type": "PropertyDefinition",
            "start": 10,
            "end": 37,
            "static": false,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 10,
              "end": 13,
              "name": "aaa"
            },
            "value": {
              "type": "ArrowFunctionExpression",
              "start": 16,
              "end": 37,
              "id": null,
              "expression": true,
              "generator": false,
              "async": false,
              "params": [],
              "body": {
                "type": "ArrowFunctionExpression",
                "start": 22,
                "end": 37,
                "id": null,
                "expression": true,
                "generator": false,
                "async": false,
                "params": [],
                "body": {
                  "type": "MemberExpression",
                  "start": 28,
                  "end": 37,
                  "object": {
                    "type": "Super",
                    "start": 28,
                    "end": 33
                  },
                  "property": {
                    "type": "Identifier",
                    "start": 34,
                    "end": 37,
                    "name": "bbb"
                  },
                  "computed": false,
                  "optional": false
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

func TestClassFeature31(t *testing.T) {
	ast, err := Compile("class C { prototype }")
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
        "name": "C"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "start": 8,
        "end": 21,
        "body": [
          {
            "type": "PropertyDefinition",
            "start": 10,
            "end": 19,
            "static": false,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 10,
              "end": 19,
              "name": "prototype"
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

func TestClassFeature32(t *testing.T) {
	ast, err := Compile("class C { aaa = { arguments: 1 } }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 34,
  "body": [
    {
      "type": "ClassDeclaration",
      "start": 0,
      "end": 34,
      "id": {
        "type": "Identifier",
        "start": 6,
        "end": 7,
        "name": "C"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "start": 8,
        "end": 34,
        "body": [
          {
            "type": "PropertyDefinition",
            "start": 10,
            "end": 32,
            "static": false,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 10,
              "end": 13,
              "name": "aaa"
            },
            "value": {
              "type": "ObjectExpression",
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
                    "end": 27,
                    "name": "arguments"
                  },
                  "value": {
                    "type": "Literal",
                    "start": 29,
                    "end": 30,
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
    }
  ]
}
`, ast)
}

func TestClassFeature33(t *testing.T) {
	ast, err := Compile("class C { aaa = function(){ arguments } }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 41,
  "body": [
    {
      "type": "ClassDeclaration",
      "start": 0,
      "end": 41,
      "id": {
        "type": "Identifier",
        "start": 6,
        "end": 7,
        "name": "C"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "start": 8,
        "end": 41,
        "body": [
          {
            "type": "PropertyDefinition",
            "start": 10,
            "end": 39,
            "static": false,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 10,
              "end": 13,
              "name": "aaa"
            },
            "value": {
              "type": "FunctionExpression",
              "start": 16,
              "end": 39,
              "id": null,
              "expression": false,
              "generator": false,
              "async": false,
              "params": [],
              "body": {
                "type": "BlockStatement",
                "start": 26,
                "end": 39,
                "body": [
                  {
                    "type": "ExpressionStatement",
                    "start": 28,
                    "end": 37,
                    "expression": {
                      "type": "Identifier",
                      "start": 28,
                      "end": 37,
                      "name": "arguments"
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

func TestClassFeature34(t *testing.T) {
	ast, err := Compile("class C { [arguments] = 0 }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
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
        "name": "C"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "start": 8,
        "end": 27,
        "body": [
          {
            "type": "PropertyDefinition",
            "start": 10,
            "end": 25,
            "static": false,
            "computed": true,
            "key": {
              "type": "Identifier",
              "start": 11,
              "end": 20,
              "name": "arguments"
            },
            "value": {
              "type": "Literal",
              "start": 24,
              "end": 25,
              "value": 0,
              "raw": "0"
            }
          }
        ]
      }
    }
  ]
}
`, ast)
}

func TestClassFeature35(t *testing.T) {
	ast, err := Compile("class C { #aaa }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 16,
  "body": [
    {
      "type": "ClassDeclaration",
      "start": 0,
      "end": 16,
      "id": {
        "type": "Identifier",
        "start": 6,
        "end": 7,
        "name": "C"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "start": 8,
        "end": 16,
        "body": [
          {
            "type": "PropertyDefinition",
            "start": 10,
            "end": 14,
            "static": false,
            "computed": false,
            "key": {
              "type": "PrivateIdentifier",
              "start": 10,
              "end": 14,
              "name": "aaa"
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

func TestClassFeature36(t *testing.T) {
	ast, err := Compile("class C { #\\u0041 }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 19,
  "body": [
    {
      "type": "ClassDeclaration",
      "start": 0,
      "end": 19,
      "id": {
        "type": "Identifier",
        "start": 6,
        "end": 7,
        "name": "C"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "start": 8,
        "end": 19,
        "body": [
          {
            "type": "PropertyDefinition",
            "start": 10,
            "end": 17,
            "static": false,
            "computed": false,
            "key": {
              "type": "PrivateIdentifier",
              "start": 10,
              "end": 17,
              "name": "A"
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

func TestClassFeature37(t *testing.T) {
	ast, err := Compile("class C { static #aaa }")
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
        "name": "C"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "start": 8,
        "end": 23,
        "body": [
          {
            "type": "PropertyDefinition",
            "start": 10,
            "end": 21,
            "static": true,
            "computed": false,
            "key": {
              "type": "PrivateIdentifier",
              "start": 17,
              "end": 21,
              "name": "aaa"
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

func TestClassFeature38(t *testing.T) {
	ast, err := Compile("class C { static\n#\\u0041 }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 26,
  "body": [
    {
      "type": "ClassDeclaration",
      "start": 0,
      "end": 26,
      "id": {
        "type": "Identifier",
        "start": 6,
        "end": 7,
        "name": "C"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "start": 8,
        "end": 26,
        "body": [
          {
            "type": "PropertyDefinition",
            "start": 10,
            "end": 24,
            "static": true,
            "computed": false,
            "key": {
              "type": "PrivateIdentifier",
              "start": 17,
              "end": 24,
              "name": "A"
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

func TestClassFeature39(t *testing.T) {
	ast, err := Compile("class C { #aaa; }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
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
        "name": "C"
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
              "type": "PrivateIdentifier",
              "start": 10,
              "end": 14,
              "name": "aaa"
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

func TestClassFeature40(t *testing.T) {
	ast, err := Compile("class C { #aaa; #bbb }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 22,
  "body": [
    {
      "type": "ClassDeclaration",
      "start": 0,
      "end": 22,
      "id": {
        "type": "Identifier",
        "start": 6,
        "end": 7,
        "name": "C"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "start": 8,
        "end": 22,
        "body": [
          {
            "type": "PropertyDefinition",
            "start": 10,
            "end": 15,
            "static": false,
            "computed": false,
            "key": {
              "type": "PrivateIdentifier",
              "start": 10,
              "end": 14,
              "name": "aaa"
            },
            "value": null
          },
          {
            "type": "PropertyDefinition",
            "start": 16,
            "end": 20,
            "static": false,
            "computed": false,
            "key": {
              "type": "PrivateIdentifier",
              "start": 16,
              "end": 20,
              "name": "bbb"
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

func TestClassFeature41(t *testing.T) {
	ast, err := Compile("class C { #aaa\n#bbb }")
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
        "name": "C"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "start": 8,
        "end": 21,
        "body": [
          {
            "type": "PropertyDefinition",
            "start": 10,
            "end": 14,
            "static": false,
            "computed": false,
            "key": {
              "type": "PrivateIdentifier",
              "start": 10,
              "end": 14,
              "name": "aaa"
            },
            "value": null
          },
          {
            "type": "PropertyDefinition",
            "start": 15,
            "end": 19,
            "static": false,
            "computed": false,
            "key": {
              "type": "PrivateIdentifier",
              "start": 15,
              "end": 19,
              "name": "bbb"
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

func TestClassFeature42(t *testing.T) {
	ast, err := Compile("class C { #aaa\n #bbb }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 22,
  "body": [
    {
      "type": "ClassDeclaration",
      "start": 0,
      "end": 22,
      "id": {
        "type": "Identifier",
        "start": 6,
        "end": 7,
        "name": "C"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "start": 8,
        "end": 22,
        "body": [
          {
            "type": "PropertyDefinition",
            "start": 10,
            "end": 14,
            "static": false,
            "computed": false,
            "key": {
              "type": "PrivateIdentifier",
              "start": 10,
              "end": 14,
              "name": "aaa"
            },
            "value": null
          },
          {
            "type": "PropertyDefinition",
            "start": 16,
            "end": 20,
            "static": false,
            "computed": false,
            "key": {
              "type": "PrivateIdentifier",
              "start": 16,
              "end": 20,
              "name": "bbb"
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

func TestClassFeature43(t *testing.T) {
	ast, err := Compile("class C { #aaa = 0 }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 20,
  "body": [
    {
      "type": "ClassDeclaration",
      "start": 0,
      "end": 20,
      "id": {
        "type": "Identifier",
        "start": 6,
        "end": 7,
        "name": "C"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "start": 8,
        "end": 20,
        "body": [
          {
            "type": "PropertyDefinition",
            "start": 10,
            "end": 18,
            "static": false,
            "computed": false,
            "key": {
              "type": "PrivateIdentifier",
              "start": 10,
              "end": 14,
              "name": "aaa"
            },
            "value": {
              "type": "Literal",
              "start": 17,
              "end": 18,
              "value": 0,
              "raw": "0"
            }
          }
        ]
      }
    }
  ]
}
`, ast)
}

func TestClassFeature44(t *testing.T) {
	ast, err := Compile("class C { #𩸽 }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
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
        "name": "C"
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
              "type": "PrivateIdentifier",
              "start": 10,
              "end": 15,
              "name": "𩸽"
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

func TestClassFeature45(t *testing.T) {
	ast, err := Compile("class C { a; #a }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
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
        "name": "C"
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
            "end": 12,
            "static": false,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 10,
              "end": 11,
              "name": "a"
            },
            "value": null
          },
          {
            "type": "PropertyDefinition",
            "start": 13,
            "end": 15,
            "static": false,
            "computed": false,
            "key": {
              "type": "PrivateIdentifier",
              "start": 13,
              "end": 15,
              "name": "a"
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

func TestClassFeature46(t *testing.T) {
	ast, err := Compile("class C { #a; a }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
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
        "name": "C"
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
            "type": "PropertyDefinition",
            "start": 14,
            "end": 15,
            "static": false,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 14,
              "end": 15,
              "name": "a"
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

func TestClassFeature47(t *testing.T) {
	ast, err := Compile("class C { #aaa(){} }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 20,
  "body": [
    {
      "type": "ClassDeclaration",
      "start": 0,
      "end": 20,
      "id": {
        "type": "Identifier",
        "start": 6,
        "end": 7,
        "name": "C"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "start": 8,
        "end": 20,
        "body": [
          {
            "type": "MethodDefinition",
            "start": 10,
            "end": 18,
            "static": false,
            "computed": false,
            "key": {
              "type": "PrivateIdentifier",
              "start": 10,
              "end": 14,
              "name": "aaa"
            },
            "kind": "method",
            "value": {
              "type": "FunctionExpression",
              "start": 14,
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
          }
        ]
      }
    }
  ]
}
`, ast)
}

func TestClassFeature48(t *testing.T) {
	ast, err := Compile("class C { *#aaa(){} }")
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
        "name": "C"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "start": 8,
        "end": 21,
        "body": [
          {
            "type": "MethodDefinition",
            "start": 10,
            "end": 19,
            "static": false,
            "computed": false,
            "key": {
              "type": "PrivateIdentifier",
              "start": 11,
              "end": 15,
              "name": "aaa"
            },
            "kind": "method",
            "value": {
              "type": "FunctionExpression",
              "start": 15,
              "end": 19,
              "id": null,
              "expression": false,
              "generator": true,
              "async": false,
              "params": [],
              "body": {
                "type": "BlockStatement",
                "start": 17,
                "end": 19,
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

func TestClassFeature49(t *testing.T) {
	ast, err := Compile("class C { async#aaa(){} }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 25,
  "body": [
    {
      "type": "ClassDeclaration",
      "start": 0,
      "end": 25,
      "id": {
        "type": "Identifier",
        "start": 6,
        "end": 7,
        "name": "C"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "start": 8,
        "end": 25,
        "body": [
          {
            "type": "MethodDefinition",
            "start": 10,
            "end": 23,
            "static": false,
            "computed": false,
            "key": {
              "type": "PrivateIdentifier",
              "start": 15,
              "end": 19,
              "name": "aaa"
            },
            "kind": "method",
            "value": {
              "type": "FunctionExpression",
              "start": 19,
              "end": 23,
              "id": null,
              "expression": false,
              "generator": false,
              "async": true,
              "params": [],
              "body": {
                "type": "BlockStatement",
                "start": 21,
                "end": 23,
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

func TestClassFeature50(t *testing.T) {
	ast, err := Compile("class C { async*#aaa(){} }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 26,
  "body": [
    {
      "type": "ClassDeclaration",
      "start": 0,
      "end": 26,
      "id": {
        "type": "Identifier",
        "start": 6,
        "end": 7,
        "name": "C"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "start": 8,
        "end": 26,
        "body": [
          {
            "type": "MethodDefinition",
            "start": 10,
            "end": 24,
            "static": false,
            "computed": false,
            "key": {
              "type": "PrivateIdentifier",
              "start": 16,
              "end": 20,
              "name": "aaa"
            },
            "kind": "method",
            "value": {
              "type": "FunctionExpression",
              "start": 20,
              "end": 24,
              "id": null,
              "expression": false,
              "generator": true,
              "async": true,
              "params": [],
              "body": {
                "type": "BlockStatement",
                "start": 22,
                "end": 24,
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

func TestClassFeature51(t *testing.T) {
	ast, err := Compile("class C { static#aaa(){} }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 26,
  "body": [
    {
      "type": "ClassDeclaration",
      "start": 0,
      "end": 26,
      "id": {
        "type": "Identifier",
        "start": 6,
        "end": 7,
        "name": "C"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "start": 8,
        "end": 26,
        "body": [
          {
            "type": "MethodDefinition",
            "start": 10,
            "end": 24,
            "static": true,
            "computed": false,
            "key": {
              "type": "PrivateIdentifier",
              "start": 16,
              "end": 20,
              "name": "aaa"
            },
            "kind": "method",
            "value": {
              "type": "FunctionExpression",
              "start": 20,
              "end": 24,
              "id": null,
              "expression": false,
              "generator": false,
              "async": false,
              "params": [],
              "body": {
                "type": "BlockStatement",
                "start": 22,
                "end": 24,
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

func TestClassFeature52(t *testing.T) {
	ast, err := Compile("class C { static*#aaa(){} }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
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
        "name": "C"
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
            "static": true,
            "computed": false,
            "key": {
              "type": "PrivateIdentifier",
              "start": 17,
              "end": 21,
              "name": "aaa"
            },
            "kind": "method",
            "value": {
              "type": "FunctionExpression",
              "start": 21,
              "end": 25,
              "id": null,
              "expression": false,
              "generator": true,
              "async": false,
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

func TestClassFeature53(t *testing.T) {
	ast, err := Compile("class C { static async#aaa(){} }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 32,
  "body": [
    {
      "type": "ClassDeclaration",
      "start": 0,
      "end": 32,
      "id": {
        "type": "Identifier",
        "start": 6,
        "end": 7,
        "name": "C"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "start": 8,
        "end": 32,
        "body": [
          {
            "type": "MethodDefinition",
            "start": 10,
            "end": 30,
            "static": true,
            "computed": false,
            "key": {
              "type": "PrivateIdentifier",
              "start": 22,
              "end": 26,
              "name": "aaa"
            },
            "kind": "method",
            "value": {
              "type": "FunctionExpression",
              "start": 26,
              "end": 30,
              "id": null,
              "expression": false,
              "generator": false,
              "async": true,
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
    }
  ]
}
`, ast)
}

func TestClassFeature54(t *testing.T) {
	ast, err := Compile("class C { static async*#aaa(){} }")
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
        "name": "C"
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
            "computed": false,
            "key": {
              "type": "PrivateIdentifier",
              "start": 23,
              "end": 27,
              "name": "aaa"
            },
            "kind": "method",
            "value": {
              "type": "FunctionExpression",
              "start": 27,
              "end": 31,
              "id": null,
              "expression": false,
              "generator": true,
              "async": true,
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

func TestClassFeature55(t *testing.T) {
	ast, err := Compile("class C { static\n*\n#aaa(){} }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 29,
  "body": [
    {
      "type": "ClassDeclaration",
      "start": 0,
      "end": 29,
      "id": {
        "type": "Identifier",
        "start": 6,
        "end": 7,
        "name": "C"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "start": 8,
        "end": 29,
        "body": [
          {
            "type": "MethodDefinition",
            "start": 10,
            "end": 27,
            "static": true,
            "computed": false,
            "key": {
              "type": "PrivateIdentifier",
              "start": 19,
              "end": 23,
              "name": "aaa"
            },
            "kind": "method",
            "value": {
              "type": "FunctionExpression",
              "start": 23,
              "end": 27,
              "id": null,
              "expression": false,
              "generator": true,
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

func TestClassFeature56(t *testing.T) {
	ast, err := Compile("class C { static\nasync#aaa(){} }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 32,
  "body": [
    {
      "type": "ClassDeclaration",
      "start": 0,
      "end": 32,
      "id": {
        "type": "Identifier",
        "start": 6,
        "end": 7,
        "name": "C"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "start": 8,
        "end": 32,
        "body": [
          {
            "type": "MethodDefinition",
            "start": 10,
            "end": 30,
            "static": true,
            "computed": false,
            "key": {
              "type": "PrivateIdentifier",
              "start": 22,
              "end": 26,
              "name": "aaa"
            },
            "kind": "method",
            "value": {
              "type": "FunctionExpression",
              "start": 26,
              "end": 30,
              "id": null,
              "expression": false,
              "generator": false,
              "async": true,
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
    }
  ]
}
`, ast)
}

func TestClassFeature57(t *testing.T) {
	ast, err := Compile("class C { static\nasync*\n#aaa(){} }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 34,
  "body": [
    {
      "type": "ClassDeclaration",
      "start": 0,
      "end": 34,
      "id": {
        "type": "Identifier",
        "start": 6,
        "end": 7,
        "name": "C"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "start": 8,
        "end": 34,
        "body": [
          {
            "type": "MethodDefinition",
            "start": 10,
            "end": 32,
            "static": true,
            "computed": false,
            "key": {
              "type": "PrivateIdentifier",
              "start": 24,
              "end": 28,
              "name": "aaa"
            },
            "kind": "method",
            "value": {
              "type": "FunctionExpression",
              "start": 28,
              "end": 32,
              "id": null,
              "expression": false,
              "generator": true,
              "async": true,
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
`, ast)
}

func TestClassFeature58(t *testing.T) {
	ast, err := Compile("class C { static\nasync\n#aaa(){} }")
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
        "name": "C"
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
            "end": 22,
            "static": true,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 17,
              "end": 22,
              "name": "async"
            },
            "value": null
          },
          {
            "type": "MethodDefinition",
            "start": 23,
            "end": 31,
            "static": false,
            "computed": false,
            "key": {
              "type": "PrivateIdentifier",
              "start": 23,
              "end": 27,
              "name": "aaa"
            },
            "kind": "method",
            "value": {
              "type": "FunctionExpression",
              "start": 27,
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

func TestClassFeature59(t *testing.T) {
	ast, err := Compile("class C { static\nasync\n*\n#aaa(){} }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 35,
  "body": [
    {
      "type": "ClassDeclaration",
      "start": 0,
      "end": 35,
      "id": {
        "type": "Identifier",
        "start": 6,
        "end": 7,
        "name": "C"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "start": 8,
        "end": 35,
        "body": [
          {
            "type": "PropertyDefinition",
            "start": 10,
            "end": 22,
            "static": true,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 17,
              "end": 22,
              "name": "async"
            },
            "value": null
          },
          {
            "type": "MethodDefinition",
            "start": 23,
            "end": 33,
            "static": false,
            "computed": false,
            "key": {
              "type": "PrivateIdentifier",
              "start": 25,
              "end": 29,
              "name": "aaa"
            },
            "kind": "method",
            "value": {
              "type": "FunctionExpression",
              "start": 29,
              "end": 33,
              "id": null,
              "expression": false,
              "generator": true,
              "async": false,
              "params": [],
              "body": {
                "type": "BlockStatement",
                "start": 31,
                "end": 33,
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

func TestClassFeature60(t *testing.T) {
	ast, err := Compile("class C { static\n async\n #aaa(){} }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 35,
  "body": [
    {
      "type": "ClassDeclaration",
      "start": 0,
      "end": 35,
      "id": {
        "type": "Identifier",
        "start": 6,
        "end": 7,
        "name": "C"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "start": 8,
        "end": 35,
        "body": [
          {
            "type": "PropertyDefinition",
            "start": 10,
            "end": 23,
            "static": true,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 18,
              "end": 23,
              "name": "async"
            },
            "value": null
          },
          {
            "type": "MethodDefinition",
            "start": 25,
            "end": 33,
            "static": false,
            "computed": false,
            "key": {
              "type": "PrivateIdentifier",
              "start": 25,
              "end": 29,
              "name": "aaa"
            },
            "kind": "method",
            "value": {
              "type": "FunctionExpression",
              "start": 29,
              "end": 33,
              "id": null,
              "expression": false,
              "generator": false,
              "async": false,
              "params": [],
              "body": {
                "type": "BlockStatement",
                "start": 31,
                "end": 33,
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

func TestClassFeature61(t *testing.T) {
	ast, err := Compile("class C { static\n async\n *\n #aaa(){} }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 38,
  "body": [
    {
      "type": "ClassDeclaration",
      "start": 0,
      "end": 38,
      "id": {
        "type": "Identifier",
        "start": 6,
        "end": 7,
        "name": "C"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "start": 8,
        "end": 38,
        "body": [
          {
            "type": "PropertyDefinition",
            "start": 10,
            "end": 23,
            "static": true,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 18,
              "end": 23,
              "name": "async"
            },
            "value": null
          },
          {
            "type": "MethodDefinition",
            "start": 25,
            "end": 36,
            "static": false,
            "computed": false,
            "key": {
              "type": "PrivateIdentifier",
              "start": 28,
              "end": 32,
              "name": "aaa"
            },
            "kind": "method",
            "value": {
              "type": "FunctionExpression",
              "start": 32,
              "end": 36,
              "id": null,
              "expression": false,
              "generator": true,
              "async": false,
              "params": [],
              "body": {
                "type": "BlockStatement",
                "start": 34,
                "end": 36,
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

func TestClassFeature62(t *testing.T) {
	ast, err := Compile("class C { get #aaa(){} }")
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
        "name": "C"
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
              "type": "PrivateIdentifier",
              "start": 14,
              "end": 18,
              "name": "aaa"
            },
            "kind": "get",
            "value": {
              "type": "FunctionExpression",
              "start": 18,
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
            }
          }
        ]
      }
    }
  ]
}
`, ast)
}

func TestClassFeature63(t *testing.T) {
	ast, err := Compile("class C { set #aaa(x){} }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 25,
  "body": [
    {
      "type": "ClassDeclaration",
      "start": 0,
      "end": 25,
      "id": {
        "type": "Identifier",
        "start": 6,
        "end": 7,
        "name": "C"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "start": 8,
        "end": 25,
        "body": [
          {
            "type": "MethodDefinition",
            "start": 10,
            "end": 23,
            "static": false,
            "computed": false,
            "key": {
              "type": "PrivateIdentifier",
              "start": 14,
              "end": 18,
              "name": "aaa"
            },
            "kind": "set",
            "value": {
              "type": "FunctionExpression",
              "start": 18,
              "end": 23,
              "id": null,
              "expression": false,
              "generator": false,
              "async": false,
              "params": [
                {
                  "type": "Identifier",
                  "start": 19,
                  "end": 20,
                  "name": "x"
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
        ]
      }
    }
  ]
}
`, ast)
}

func TestClassFeature64(t *testing.T) {
	ast, err := Compile("class C { static get #aaa(){} }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 31,
  "body": [
    {
      "type": "ClassDeclaration",
      "start": 0,
      "end": 31,
      "id": {
        "type": "Identifier",
        "start": 6,
        "end": 7,
        "name": "C"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "start": 8,
        "end": 31,
        "body": [
          {
            "type": "MethodDefinition",
            "start": 10,
            "end": 29,
            "static": true,
            "computed": false,
            "key": {
              "type": "PrivateIdentifier",
              "start": 21,
              "end": 25,
              "name": "aaa"
            },
            "kind": "get",
            "value": {
              "type": "FunctionExpression",
              "start": 25,
              "end": 29,
              "id": null,
              "expression": false,
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
    }
  ]
}
`, ast)
}

func TestClassFeature65(t *testing.T) {
	ast, err := Compile("class C { static set #aaa(x){} }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 32,
  "body": [
    {
      "type": "ClassDeclaration",
      "start": 0,
      "end": 32,
      "id": {
        "type": "Identifier",
        "start": 6,
        "end": 7,
        "name": "C"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "start": 8,
        "end": 32,
        "body": [
          {
            "type": "MethodDefinition",
            "start": 10,
            "end": 30,
            "static": true,
            "computed": false,
            "key": {
              "type": "PrivateIdentifier",
              "start": 21,
              "end": 25,
              "name": "aaa"
            },
            "kind": "set",
            "value": {
              "type": "FunctionExpression",
              "start": 25,
              "end": 30,
              "id": null,
              "expression": false,
              "generator": false,
              "async": false,
              "params": [
                {
                  "type": "Identifier",
                  "start": 26,
                  "end": 27,
                  "name": "x"
                }
              ],
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
    }
  ]
}
`, ast)
}

func TestClassFeature66(t *testing.T) {
	ast, err := Compile("class C { static\nget\n#aaa(){} }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 31,
  "body": [
    {
      "type": "ClassDeclaration",
      "start": 0,
      "end": 31,
      "id": {
        "type": "Identifier",
        "start": 6,
        "end": 7,
        "name": "C"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "start": 8,
        "end": 31,
        "body": [
          {
            "type": "MethodDefinition",
            "start": 10,
            "end": 29,
            "static": true,
            "computed": false,
            "key": {
              "type": "PrivateIdentifier",
              "start": 21,
              "end": 25,
              "name": "aaa"
            },
            "kind": "get",
            "value": {
              "type": "FunctionExpression",
              "start": 25,
              "end": 29,
              "id": null,
              "expression": false,
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
    }
  ]
}
`, ast)
}

func TestClassFeature67(t *testing.T) {
	ast, err := Compile("class C { get; #aaa(){} }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 25,
  "body": [
    {
      "type": "ClassDeclaration",
      "start": 0,
      "end": 25,
      "id": {
        "type": "Identifier",
        "start": 6,
        "end": 7,
        "name": "C"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "start": 8,
        "end": 25,
        "body": [
          {
            "type": "PropertyDefinition",
            "start": 10,
            "end": 14,
            "static": false,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 10,
              "end": 13,
              "name": "get"
            },
            "value": null
          },
          {
            "type": "MethodDefinition",
            "start": 15,
            "end": 23,
            "static": false,
            "computed": false,
            "key": {
              "type": "PrivateIdentifier",
              "start": 15,
              "end": 19,
              "name": "aaa"
            },
            "kind": "method",
            "value": {
              "type": "FunctionExpression",
              "start": 19,
              "end": 23,
              "id": null,
              "expression": false,
              "generator": false,
              "async": false,
              "params": [],
              "body": {
                "type": "BlockStatement",
                "start": 21,
                "end": 23,
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

func TestClassFeature68(t *testing.T) {
	ast, err := Compile("class C { get #a(){} set #a(x){} }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 34,
  "body": [
    {
      "type": "ClassDeclaration",
      "start": 0,
      "end": 34,
      "id": {
        "type": "Identifier",
        "start": 6,
        "end": 7,
        "name": "C"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "start": 8,
        "end": 34,
        "body": [
          {
            "type": "MethodDefinition",
            "start": 10,
            "end": 20,
            "static": false,
            "computed": false,
            "key": {
              "type": "PrivateIdentifier",
              "start": 14,
              "end": 16,
              "name": "a"
            },
            "kind": "get",
            "value": {
              "type": "FunctionExpression",
              "start": 16,
              "end": 20,
              "id": null,
              "expression": false,
              "generator": false,
              "async": false,
              "params": [],
              "body": {
                "type": "BlockStatement",
                "start": 18,
                "end": 20,
                "body": []
              }
            }
          },
          {
            "type": "MethodDefinition",
            "start": 21,
            "end": 32,
            "static": false,
            "computed": false,
            "key": {
              "type": "PrivateIdentifier",
              "start": 25,
              "end": 27,
              "name": "a"
            },
            "kind": "set",
            "value": {
              "type": "FunctionExpression",
              "start": 27,
              "end": 32,
              "id": null,
              "expression": false,
              "generator": false,
              "async": false,
              "params": [
                {
                  "type": "Identifier",
                  "start": 28,
                  "end": 29,
                  "name": "x"
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
        ]
      }
    }
  ]
}
`, ast)
}

func TestClassFeature69(t *testing.T) {
	ast, err := Compile("class C { static get #a(){} static set #a(x){} }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 48,
  "body": [
    {
      "type": "ClassDeclaration",
      "start": 0,
      "end": 48,
      "id": {
        "type": "Identifier",
        "start": 6,
        "end": 7,
        "name": "C"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "start": 8,
        "end": 48,
        "body": [
          {
            "type": "MethodDefinition",
            "start": 10,
            "end": 27,
            "static": true,
            "computed": false,
            "key": {
              "type": "PrivateIdentifier",
              "start": 21,
              "end": 23,
              "name": "a"
            },
            "kind": "get",
            "value": {
              "type": "FunctionExpression",
              "start": 23,
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
          },
          {
            "type": "MethodDefinition",
            "start": 28,
            "end": 46,
            "static": true,
            "computed": false,
            "key": {
              "type": "PrivateIdentifier",
              "start": 39,
              "end": 41,
              "name": "a"
            },
            "kind": "set",
            "value": {
              "type": "FunctionExpression",
              "start": 41,
              "end": 46,
              "id": null,
              "expression": false,
              "generator": false,
              "async": false,
              "params": [
                {
                  "type": "Identifier",
                  "start": 42,
                  "end": 43,
                  "name": "x"
                }
              ],
              "body": {
                "type": "BlockStatement",
                "start": 44,
                "end": 46,
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

func TestClassFeature70(t *testing.T) {
	ast, err := Compile("class C { #aaa; f() { this.#aaa } }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 35,
  "body": [
    {
      "type": "ClassDeclaration",
      "start": 0,
      "end": 35,
      "id": {
        "type": "Identifier",
        "start": 6,
        "end": 7,
        "name": "C"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "start": 8,
        "end": 35,
        "body": [
          {
            "type": "PropertyDefinition",
            "start": 10,
            "end": 15,
            "static": false,
            "computed": false,
            "key": {
              "type": "PrivateIdentifier",
              "start": 10,
              "end": 14,
              "name": "aaa"
            },
            "value": null
          },
          {
            "type": "MethodDefinition",
            "start": 16,
            "end": 33,
            "static": false,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 16,
              "end": 17,
              "name": "f"
            },
            "kind": "method",
            "value": {
              "type": "FunctionExpression",
              "start": 17,
              "end": 33,
              "id": null,
              "expression": false,
              "generator": false,
              "async": false,
              "params": [],
              "body": {
                "type": "BlockStatement",
                "start": 20,
                "end": 33,
                "body": [
                  {
                    "type": "ExpressionStatement",
                    "start": 22,
                    "end": 31,
                    "expression": {
                      "type": "MemberExpression",
                      "start": 22,
                      "end": 31,
                      "object": {
                        "type": "ThisExpression",
                        "start": 22,
                        "end": 26
                      },
                      "property": {
                        "type": "PrivateIdentifier",
                        "start": 27,
                        "end": 31,
                        "name": "aaa"
                      },
                      "computed": false,
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

func TestClassFeature71(t *testing.T) {
	ast, err := Compile("class C { #aaa; f(obj) { obj.#aaa } }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 37,
  "body": [
    {
      "type": "ClassDeclaration",
      "start": 0,
      "end": 37,
      "id": {
        "type": "Identifier",
        "start": 6,
        "end": 7,
        "name": "C"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "start": 8,
        "end": 37,
        "body": [
          {
            "type": "PropertyDefinition",
            "start": 10,
            "end": 15,
            "static": false,
            "computed": false,
            "key": {
              "type": "PrivateIdentifier",
              "start": 10,
              "end": 14,
              "name": "aaa"
            },
            "value": null
          },
          {
            "type": "MethodDefinition",
            "start": 16,
            "end": 35,
            "static": false,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 16,
              "end": 17,
              "name": "f"
            },
            "kind": "method",
            "value": {
              "type": "FunctionExpression",
              "start": 17,
              "end": 35,
              "id": null,
              "expression": false,
              "generator": false,
              "async": false,
              "params": [
                {
                  "type": "Identifier",
                  "start": 18,
                  "end": 21,
                  "name": "obj"
                }
              ],
              "body": {
                "type": "BlockStatement",
                "start": 23,
                "end": 35,
                "body": [
                  {
                    "type": "ExpressionStatement",
                    "start": 25,
                    "end": 33,
                    "expression": {
                      "type": "MemberExpression",
                      "start": 25,
                      "end": 33,
                      "object": {
                        "type": "Identifier",
                        "start": 25,
                        "end": 28,
                        "name": "obj"
                      },
                      "property": {
                        "type": "PrivateIdentifier",
                        "start": 29,
                        "end": 33,
                        "name": "aaa"
                      },
                      "computed": false,
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

func TestClassFeature72(t *testing.T) {
	ast, err := Compile("class C { #aaa; f(obj) { obj?.#aaa } }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 38,
  "body": [
    {
      "type": "ClassDeclaration",
      "start": 0,
      "end": 38,
      "id": {
        "type": "Identifier",
        "start": 6,
        "end": 7,
        "name": "C"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "start": 8,
        "end": 38,
        "body": [
          {
            "type": "PropertyDefinition",
            "start": 10,
            "end": 15,
            "static": false,
            "computed": false,
            "key": {
              "type": "PrivateIdentifier",
              "start": 10,
              "end": 14,
              "name": "aaa"
            },
            "value": null
          },
          {
            "type": "MethodDefinition",
            "start": 16,
            "end": 36,
            "static": false,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 16,
              "end": 17,
              "name": "f"
            },
            "kind": "method",
            "value": {
              "type": "FunctionExpression",
              "start": 17,
              "end": 36,
              "id": null,
              "expression": false,
              "generator": false,
              "async": false,
              "params": [
                {
                  "type": "Identifier",
                  "start": 18,
                  "end": 21,
                  "name": "obj"
                }
              ],
              "body": {
                "type": "BlockStatement",
                "start": 23,
                "end": 36,
                "body": [
                  {
                    "type": "ExpressionStatement",
                    "start": 25,
                    "end": 34,
                    "expression": {
                      "type": "ChainExpression",
                      "start": 25,
                      "end": 34,
                      "expression": {
                        "type": "MemberExpression",
                        "start": 25,
                        "end": 34,
                        "object": {
                          "type": "Identifier",
                          "start": 25,
                          "end": 28,
                          "name": "obj"
                        },
                        "property": {
                          "type": "PrivateIdentifier",
                          "start": 30,
                          "end": 34,
                          "name": "aaa"
                        },
                        "computed": false,
                        "optional": true
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
  ]
}
`, ast)
}

func TestClassFeature73(t *testing.T) {
	ast, err := Compile("class C { #aaa; f(f) { f()?.#aaa } }")
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
        "name": "C"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "start": 8,
        "end": 36,
        "body": [
          {
            "type": "PropertyDefinition",
            "start": 10,
            "end": 15,
            "static": false,
            "computed": false,
            "key": {
              "type": "PrivateIdentifier",
              "start": 10,
              "end": 14,
              "name": "aaa"
            },
            "value": null
          },
          {
            "type": "MethodDefinition",
            "start": 16,
            "end": 34,
            "static": false,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 16,
              "end": 17,
              "name": "f"
            },
            "kind": "method",
            "value": {
              "type": "FunctionExpression",
              "start": 17,
              "end": 34,
              "id": null,
              "expression": false,
              "generator": false,
              "async": false,
              "params": [
                {
                  "type": "Identifier",
                  "start": 18,
                  "end": 19,
                  "name": "f"
                }
              ],
              "body": {
                "type": "BlockStatement",
                "start": 21,
                "end": 34,
                "body": [
                  {
                    "type": "ExpressionStatement",
                    "start": 23,
                    "end": 32,
                    "expression": {
                      "type": "ChainExpression",
                      "start": 23,
                      "end": 32,
                      "expression": {
                        "type": "MemberExpression",
                        "start": 23,
                        "end": 32,
                        "object": {
                          "type": "CallExpression",
                          "start": 23,
                          "end": 26,
                          "callee": {
                            "type": "Identifier",
                            "start": 23,
                            "end": 24,
                            "name": "f"
                          },
                          "arguments": [],
                          "optional": false
                        },
                        "property": {
                          "type": "PrivateIdentifier",
                          "start": 28,
                          "end": 32,
                          "name": "aaa"
                        },
                        "computed": false,
                        "optional": true
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
  ]
}
`, ast)
}

func TestClassFeature74(t *testing.T) {
	ast, err := Compile("class C { #aaa; f() { delete this.#aaa.foo } }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 46,
  "body": [
    {
      "type": "ClassDeclaration",
      "start": 0,
      "end": 46,
      "id": {
        "type": "Identifier",
        "start": 6,
        "end": 7,
        "name": "C"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "start": 8,
        "end": 46,
        "body": [
          {
            "type": "PropertyDefinition",
            "start": 10,
            "end": 15,
            "static": false,
            "computed": false,
            "key": {
              "type": "PrivateIdentifier",
              "start": 10,
              "end": 14,
              "name": "aaa"
            },
            "value": null
          },
          {
            "type": "MethodDefinition",
            "start": 16,
            "end": 44,
            "static": false,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 16,
              "end": 17,
              "name": "f"
            },
            "kind": "method",
            "value": {
              "type": "FunctionExpression",
              "start": 17,
              "end": 44,
              "id": null,
              "expression": false,
              "generator": false,
              "async": false,
              "params": [],
              "body": {
                "type": "BlockStatement",
                "start": 20,
                "end": 44,
                "body": [
                  {
                    "type": "ExpressionStatement",
                    "start": 22,
                    "end": 42,
                    "expression": {
                      "type": "UnaryExpression",
                      "start": 22,
                      "end": 42,
                      "operator": "delete",
                      "prefix": true,
                      "argument": {
                        "type": "MemberExpression",
                        "start": 29,
                        "end": 42,
                        "object": {
                          "type": "MemberExpression",
                          "start": 29,
                          "end": 38,
                          "object": {
                            "type": "ThisExpression",
                            "start": 29,
                            "end": 33
                          },
                          "property": {
                            "type": "PrivateIdentifier",
                            "start": 34,
                            "end": 38,
                            "name": "aaa"
                          },
                          "computed": false,
                          "optional": false
                        },
                        "property": {
                          "type": "Identifier",
                          "start": 39,
                          "end": 42,
                          "name": "foo"
                        },
                        "computed": false,
                        "optional": false
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
  ]
}
`, ast)
}

func TestClassFeature75(t *testing.T) {
	ast, err := Compile("class C { #aaa; f() { delete this.#aaa?.foo } }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 47,
  "body": [
    {
      "type": "ClassDeclaration",
      "start": 0,
      "end": 47,
      "id": {
        "type": "Identifier",
        "start": 6,
        "end": 7,
        "name": "C"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "start": 8,
        "end": 47,
        "body": [
          {
            "type": "PropertyDefinition",
            "start": 10,
            "end": 15,
            "static": false,
            "computed": false,
            "key": {
              "type": "PrivateIdentifier",
              "start": 10,
              "end": 14,
              "name": "aaa"
            },
            "value": null
          },
          {
            "type": "MethodDefinition",
            "start": 16,
            "end": 45,
            "static": false,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 16,
              "end": 17,
              "name": "f"
            },
            "kind": "method",
            "value": {
              "type": "FunctionExpression",
              "start": 17,
              "end": 45,
              "id": null,
              "expression": false,
              "generator": false,
              "async": false,
              "params": [],
              "body": {
                "type": "BlockStatement",
                "start": 20,
                "end": 45,
                "body": [
                  {
                    "type": "ExpressionStatement",
                    "start": 22,
                    "end": 43,
                    "expression": {
                      "type": "UnaryExpression",
                      "start": 22,
                      "end": 43,
                      "operator": "delete",
                      "prefix": true,
                      "argument": {
                        "type": "ChainExpression",
                        "start": 29,
                        "end": 43,
                        "expression": {
                          "type": "MemberExpression",
                          "start": 29,
                          "end": 43,
                          "object": {
                            "type": "MemberExpression",
                            "start": 29,
                            "end": 38,
                            "object": {
                              "type": "ThisExpression",
                              "start": 29,
                              "end": 33
                            },
                            "property": {
                              "type": "PrivateIdentifier",
                              "start": 34,
                              "end": 38,
                              "name": "aaa"
                            },
                            "computed": false,
                            "optional": false
                          },
                          "property": {
                            "type": "Identifier",
                            "start": 40,
                            "end": 43,
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
            }
          }
        ]
      }
    }
  ]
}
`, ast)
}

func TestClassFeature76(t *testing.T) {
	ast, err := Compile("class C { #a; a = this.#a; }")
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
        "name": "C"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "start": 8,
        "end": 28,
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
            "type": "PropertyDefinition",
            "start": 14,
            "end": 26,
            "static": false,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 14,
              "end": 15,
              "name": "a"
            },
            "value": {
              "type": "MemberExpression",
              "start": 18,
              "end": 25,
              "object": {
                "type": "ThisExpression",
                "start": 18,
                "end": 22
              },
              "property": {
                "type": "PrivateIdentifier",
                "start": 23,
                "end": 25,
                "name": "a"
              },
              "computed": false,
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

func TestClassFeature77(t *testing.T) {
	ast, err := Compile("class C { a = this.#a; #a; }")
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
        "name": "C"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "start": 8,
        "end": 28,
        "body": [
          {
            "type": "PropertyDefinition",
            "start": 10,
            "end": 22,
            "static": false,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 10,
              "end": 11,
              "name": "a"
            },
            "value": {
              "type": "MemberExpression",
              "start": 14,
              "end": 21,
              "object": {
                "type": "ThisExpression",
                "start": 14,
                "end": 18
              },
              "property": {
                "type": "PrivateIdentifier",
                "start": 19,
                "end": 21,
                "name": "a"
              },
              "computed": false,
              "optional": false
            }
          },
          {
            "type": "PropertyDefinition",
            "start": 23,
            "end": 26,
            "static": false,
            "computed": false,
            "key": {
              "type": "PrivateIdentifier",
              "start": 23,
              "end": 25,
              "name": "a"
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

func TestClassFeature78(t *testing.T) {
	ast, err := Compile("class C { #a; [this.#a]; }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 26,
  "body": [
    {
      "type": "ClassDeclaration",
      "start": 0,
      "end": 26,
      "id": {
        "type": "Identifier",
        "start": 6,
        "end": 7,
        "name": "C"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "start": 8,
        "end": 26,
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
            "type": "PropertyDefinition",
            "start": 14,
            "end": 24,
            "static": false,
            "computed": true,
            "key": {
              "type": "MemberExpression",
              "start": 15,
              "end": 22,
              "object": {
                "type": "ThisExpression",
                "start": 15,
                "end": 19
              },
              "property": {
                "type": "PrivateIdentifier",
                "start": 20,
                "end": 22,
                "name": "a"
              },
              "computed": false,
              "optional": false
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

func TestClassFeature79(t *testing.T) {
	ast, err := Compile("class C { [this.#a]; #a; }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 26,
  "body": [
    {
      "type": "ClassDeclaration",
      "start": 0,
      "end": 26,
      "id": {
        "type": "Identifier",
        "start": 6,
        "end": 7,
        "name": "C"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "start": 8,
        "end": 26,
        "body": [
          {
            "type": "PropertyDefinition",
            "start": 10,
            "end": 20,
            "static": false,
            "computed": true,
            "key": {
              "type": "MemberExpression",
              "start": 11,
              "end": 18,
              "object": {
                "type": "ThisExpression",
                "start": 11,
                "end": 15
              },
              "property": {
                "type": "PrivateIdentifier",
                "start": 16,
                "end": 18,
                "name": "a"
              },
              "computed": false,
              "optional": false
            },
            "value": null
          },
          {
            "type": "PropertyDefinition",
            "start": 21,
            "end": 24,
            "static": false,
            "computed": false,
            "key": {
              "type": "PrivateIdentifier",
              "start": 21,
              "end": 23,
              "name": "a"
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

func TestClassFeature80(t *testing.T) {
	ast, err := Compile("class C { #a; f(){ this.#a } }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 30,
  "body": [
    {
      "type": "ClassDeclaration",
      "start": 0,
      "end": 30,
      "id": {
        "type": "Identifier",
        "start": 6,
        "end": 7,
        "name": "C"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "start": 8,
        "end": 30,
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
              "start": 14,
              "end": 15,
              "name": "f"
            },
            "kind": "method",
            "value": {
              "type": "FunctionExpression",
              "start": 15,
              "end": 28,
              "id": null,
              "expression": false,
              "generator": false,
              "async": false,
              "params": [],
              "body": {
                "type": "BlockStatement",
                "start": 17,
                "end": 28,
                "body": [
                  {
                    "type": "ExpressionStatement",
                    "start": 19,
                    "end": 26,
                    "expression": {
                      "type": "MemberExpression",
                      "start": 19,
                      "end": 26,
                      "object": {
                        "type": "ThisExpression",
                        "start": 19,
                        "end": 23
                      },
                      "property": {
                        "type": "PrivateIdentifier",
                        "start": 24,
                        "end": 26,
                        "name": "a"
                      },
                      "computed": false,
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

func TestClassFeature81(t *testing.T) {
	ast, err := Compile("class C { f(){ this.#a } #a; }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 30,
  "body": [
    {
      "type": "ClassDeclaration",
      "start": 0,
      "end": 30,
      "id": {
        "type": "Identifier",
        "start": 6,
        "end": 7,
        "name": "C"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "start": 8,
        "end": 30,
        "body": [
          {
            "type": "MethodDefinition",
            "start": 10,
            "end": 24,
            "static": false,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 10,
              "end": 11,
              "name": "f"
            },
            "kind": "method",
            "value": {
              "type": "FunctionExpression",
              "start": 11,
              "end": 24,
              "id": null,
              "expression": false,
              "generator": false,
              "async": false,
              "params": [],
              "body": {
                "type": "BlockStatement",
                "start": 13,
                "end": 24,
                "body": [
                  {
                    "type": "ExpressionStatement",
                    "start": 15,
                    "end": 22,
                    "expression": {
                      "type": "MemberExpression",
                      "start": 15,
                      "end": 22,
                      "object": {
                        "type": "ThisExpression",
                        "start": 15,
                        "end": 19
                      },
                      "property": {
                        "type": "PrivateIdentifier",
                        "start": 20,
                        "end": 22,
                        "name": "a"
                      },
                      "computed": false,
                      "optional": false
                    }
                  }
                ]
              }
            }
          },
          {
            "type": "PropertyDefinition",
            "start": 25,
            "end": 28,
            "static": false,
            "computed": false,
            "key": {
              "type": "PrivateIdentifier",
              "start": 25,
              "end": 27,
              "name": "a"
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

func TestClassFeature82(t *testing.T) {
	ast, err := Compile("class Outer { #outer; Inner = class { #inner; f(obj) { obj.#outer + this.#inner } }; }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 86,
  "body": [
    {
      "type": "ClassDeclaration",
      "start": 0,
      "end": 86,
      "id": {
        "type": "Identifier",
        "start": 6,
        "end": 11,
        "name": "Outer"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "start": 12,
        "end": 86,
        "body": [
          {
            "type": "PropertyDefinition",
            "start": 14,
            "end": 21,
            "static": false,
            "computed": false,
            "key": {
              "type": "PrivateIdentifier",
              "start": 14,
              "end": 20,
              "name": "outer"
            },
            "value": null
          },
          {
            "type": "PropertyDefinition",
            "start": 22,
            "end": 84,
            "static": false,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 22,
              "end": 27,
              "name": "Inner"
            },
            "value": {
              "type": "ClassExpression",
              "start": 30,
              "end": 83,
              "id": null,
              "superClass": null,
              "body": {
                "type": "ClassBody",
                "start": 36,
                "end": 83,
                "body": [
                  {
                    "type": "PropertyDefinition",
                    "start": 38,
                    "end": 45,
                    "static": false,
                    "computed": false,
                    "key": {
                      "type": "PrivateIdentifier",
                      "start": 38,
                      "end": 44,
                      "name": "inner"
                    },
                    "value": null
                  },
                  {
                    "type": "MethodDefinition",
                    "start": 46,
                    "end": 81,
                    "static": false,
                    "computed": false,
                    "key": {
                      "type": "Identifier",
                      "start": 46,
                      "end": 47,
                      "name": "f"
                    },
                    "kind": "method",
                    "value": {
                      "type": "FunctionExpression",
                      "start": 47,
                      "end": 81,
                      "id": null,
                      "expression": false,
                      "generator": false,
                      "async": false,
                      "params": [
                        {
                          "type": "Identifier",
                          "start": 48,
                          "end": 51,
                          "name": "obj"
                        }
                      ],
                      "body": {
                        "type": "BlockStatement",
                        "start": 53,
                        "end": 81,
                        "body": [
                          {
                            "type": "ExpressionStatement",
                            "start": 55,
                            "end": 79,
                            "expression": {
                              "type": "BinaryExpression",
                              "start": 55,
                              "end": 79,
                              "left": {
                                "type": "MemberExpression",
                                "start": 55,
                                "end": 65,
                                "object": {
                                  "type": "Identifier",
                                  "start": 55,
                                  "end": 58,
                                  "name": "obj"
                                },
                                "property": {
                                  "type": "PrivateIdentifier",
                                  "start": 59,
                                  "end": 65,
                                  "name": "outer"
                                },
                                "computed": false,
                                "optional": false
                              },
                              "operator": "+",
                              "right": {
                                "type": "MemberExpression",
                                "start": 68,
                                "end": 79,
                                "object": {
                                  "type": "ThisExpression",
                                  "start": 68,
                                  "end": 72
                                },
                                "property": {
                                  "type": "PrivateIdentifier",
                                  "start": 73,
                                  "end": 79,
                                  "name": "inner"
                                },
                                "computed": false,
                                "optional": false
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
        ]
      }
    }
  ]
}
`, ast)
}

func TestClassFeature83(t *testing.T) {
	ast, err := Compile("class Outer { Inner = class { f(obj) { obj.#outer + this.#inner } #inner; }; #outer; }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 86,
  "body": [
    {
      "type": "ClassDeclaration",
      "start": 0,
      "end": 86,
      "id": {
        "type": "Identifier",
        "start": 6,
        "end": 11,
        "name": "Outer"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "start": 12,
        "end": 86,
        "body": [
          {
            "type": "PropertyDefinition",
            "start": 14,
            "end": 76,
            "static": false,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 14,
              "end": 19,
              "name": "Inner"
            },
            "value": {
              "type": "ClassExpression",
              "start": 22,
              "end": 75,
              "id": null,
              "superClass": null,
              "body": {
                "type": "ClassBody",
                "start": 28,
                "end": 75,
                "body": [
                  {
                    "type": "MethodDefinition",
                    "start": 30,
                    "end": 65,
                    "static": false,
                    "computed": false,
                    "key": {
                      "type": "Identifier",
                      "start": 30,
                      "end": 31,
                      "name": "f"
                    },
                    "kind": "method",
                    "value": {
                      "type": "FunctionExpression",
                      "start": 31,
                      "end": 65,
                      "id": null,
                      "expression": false,
                      "generator": false,
                      "async": false,
                      "params": [
                        {
                          "type": "Identifier",
                          "start": 32,
                          "end": 35,
                          "name": "obj"
                        }
                      ],
                      "body": {
                        "type": "BlockStatement",
                        "start": 37,
                        "end": 65,
                        "body": [
                          {
                            "type": "ExpressionStatement",
                            "start": 39,
                            "end": 63,
                            "expression": {
                              "type": "BinaryExpression",
                              "start": 39,
                              "end": 63,
                              "left": {
                                "type": "MemberExpression",
                                "start": 39,
                                "end": 49,
                                "object": {
                                  "type": "Identifier",
                                  "start": 39,
                                  "end": 42,
                                  "name": "obj"
                                },
                                "property": {
                                  "type": "PrivateIdentifier",
                                  "start": 43,
                                  "end": 49,
                                  "name": "outer"
                                },
                                "computed": false,
                                "optional": false
                              },
                              "operator": "+",
                              "right": {
                                "type": "MemberExpression",
                                "start": 52,
                                "end": 63,
                                "object": {
                                  "type": "ThisExpression",
                                  "start": 52,
                                  "end": 56
                                },
                                "property": {
                                  "type": "PrivateIdentifier",
                                  "start": 57,
                                  "end": 63,
                                  "name": "inner"
                                },
                                "computed": false,
                                "optional": false
                              }
                            }
                          }
                        ]
                      }
                    }
                  },
                  {
                    "type": "PropertyDefinition",
                    "start": 66,
                    "end": 73,
                    "static": false,
                    "computed": false,
                    "key": {
                      "type": "PrivateIdentifier",
                      "start": 66,
                      "end": 72,
                      "name": "inner"
                    },
                    "value": null
                  }
                ]
              }
            }
          },
          {
            "type": "PropertyDefinition",
            "start": 77,
            "end": 84,
            "static": false,
            "computed": false,
            "key": {
              "type": "PrivateIdentifier",
              "start": 77,
              "end": 83,
              "name": "outer"
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

func TestClassFeature84(t *testing.T) {
	ast, err := Compile("class C { static delete() {} }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 30,
  "body": [
    {
      "type": "ClassDeclaration",
      "start": 0,
      "end": 30,
      "id": {
        "type": "Identifier",
        "start": 6,
        "end": 7,
        "name": "C"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "start": 8,
        "end": 30,
        "body": [
          {
            "type": "MethodDefinition",
            "start": 10,
            "end": 28,
            "static": true,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 17,
              "end": 23,
              "name": "delete"
            },
            "kind": "method",
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
          }
        ]
      }
    }
  ]
}
`, ast)
}

func TestClassFeature85(t *testing.T) {
	ast, err := Compile(`class C {
  static x
  static y
  static {
    try {
      const obj = doSomethingWith(this.x)
      this.y = obj.y
      this.z = obj.z
    }
    catch {
      this.y = 0
      this.z = 0
    }
  }
}`)
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 200,
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 15,
      "column": 1
    }
  },
  "body": [
    {
      "type": "ClassDeclaration",
      "start": 0,
      "end": 200,
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 15,
          "column": 1
        }
      },
      "id": {
        "type": "Identifier",
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
        "name": "C"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "start": 8,
        "end": 200,
        "loc": {
          "start": {
            "line": 1,
            "column": 8
          },
          "end": {
            "line": 15,
            "column": 1
          }
        },
        "body": [
          {
            "type": "PropertyDefinition",
            "start": 12,
            "end": 20,
            "loc": {
              "start": {
                "line": 2,
                "column": 2
              },
              "end": {
                "line": 2,
                "column": 10
              }
            },
            "static": true,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 19,
              "end": 20,
              "loc": {
                "start": {
                  "line": 2,
                  "column": 9
                },
                "end": {
                  "line": 2,
                  "column": 10
                }
              },
              "name": "x"
            },
            "value": null
          },
          {
            "type": "PropertyDefinition",
            "start": 23,
            "end": 31,
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
            "static": true,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 30,
              "end": 31,
              "loc": {
                "start": {
                  "line": 3,
                  "column": 9
                },
                "end": {
                  "line": 3,
                  "column": 10
                }
              },
              "name": "y"
            },
            "value": null
          },
          {
            "type": "StaticBlock",
            "start": 34,
            "end": 198,
            "loc": {
              "start": {
                "line": 4,
                "column": 2
              },
              "end": {
                "line": 14,
                "column": 3
              }
            },
            "body": [
              {
                "type": "TryStatement",
                "start": 47,
                "end": 194,
                "loc": {
                  "start": {
                    "line": 5,
                    "column": 4
                  },
                  "end": {
                    "line": 13,
                    "column": 5
                  }
                },
                "block": {
                  "type": "BlockStatement",
                  "start": 51,
                  "end": 142,
                  "loc": {
                    "start": {
                      "line": 5,
                      "column": 8
                    },
                    "end": {
                      "line": 9,
                      "column": 5
                    }
                  },
                  "body": [
                    {
                      "type": "VariableDeclaration",
                      "start": 59,
                      "end": 94,
                      "loc": {
                        "start": {
                          "line": 6,
                          "column": 6
                        },
                        "end": {
                          "line": 6,
                          "column": 41
                        }
                      },
                      "declarations": [
                        {
                          "type": "VariableDeclarator",
                          "start": 65,
                          "end": 94,
                          "loc": {
                            "start": {
                              "line": 6,
                              "column": 12
                            },
                            "end": {
                              "line": 6,
                              "column": 41
                            }
                          },
                          "id": {
                            "type": "Identifier",
                            "start": 65,
                            "end": 68,
                            "loc": {
                              "start": {
                                "line": 6,
                                "column": 12
                              },
                              "end": {
                                "line": 6,
                                "column": 15
                              }
                            },
                            "name": "obj"
                          },
                          "init": {
                            "type": "CallExpression",
                            "start": 71,
                            "end": 94,
                            "loc": {
                              "start": {
                                "line": 6,
                                "column": 18
                              },
                              "end": {
                                "line": 6,
                                "column": 41
                              }
                            },
                            "callee": {
                              "type": "Identifier",
                              "start": 71,
                              "end": 86,
                              "loc": {
                                "start": {
                                  "line": 6,
                                  "column": 18
                                },
                                "end": {
                                  "line": 6,
                                  "column": 33
                                }
                              },
                              "name": "doSomethingWith"
                            },
                            "arguments": [
                              {
                                "type": "MemberExpression",
                                "start": 87,
                                "end": 93,
                                "loc": {
                                  "start": {
                                    "line": 6,
                                    "column": 34
                                  },
                                  "end": {
                                    "line": 6,
                                    "column": 40
                                  }
                                },
                                "object": {
                                  "type": "ThisExpression",
                                  "start": 87,
                                  "end": 91,
                                  "loc": {
                                    "start": {
                                      "line": 6,
                                      "column": 34
                                    },
                                    "end": {
                                      "line": 6,
                                      "column": 38
                                    }
                                  }
                                },
                                "property": {
                                  "type": "Identifier",
                                  "start": 92,
                                  "end": 93,
                                  "loc": {
                                    "start": {
                                      "line": 6,
                                      "column": 39
                                    },
                                    "end": {
                                      "line": 6,
                                      "column": 40
                                    }
                                  },
                                  "name": "x"
                                },
                                "computed": false,
                                "optional": false
                              }
                            ],
                            "optional": false
                          }
                        }
                      ],
                      "kind": "const"
                    },
                    {
                      "type": "ExpressionStatement",
                      "start": 101,
                      "end": 115,
                      "loc": {
                        "start": {
                          "line": 7,
                          "column": 6
                        },
                        "end": {
                          "line": 7,
                          "column": 20
                        }
                      },
                      "expression": {
                        "type": "AssignmentExpression",
                        "start": 101,
                        "end": 115,
                        "loc": {
                          "start": {
                            "line": 7,
                            "column": 6
                          },
                          "end": {
                            "line": 7,
                            "column": 20
                          }
                        },
                        "operator": "=",
                        "left": {
                          "type": "MemberExpression",
                          "start": 101,
                          "end": 107,
                          "loc": {
                            "start": {
                              "line": 7,
                              "column": 6
                            },
                            "end": {
                              "line": 7,
                              "column": 12
                            }
                          },
                          "object": {
                            "type": "ThisExpression",
                            "start": 101,
                            "end": 105,
                            "loc": {
                              "start": {
                                "line": 7,
                                "column": 6
                              },
                              "end": {
                                "line": 7,
                                "column": 10
                              }
                            }
                          },
                          "property": {
                            "type": "Identifier",
                            "start": 106,
                            "end": 107,
                            "loc": {
                              "start": {
                                "line": 7,
                                "column": 11
                              },
                              "end": {
                                "line": 7,
                                "column": 12
                              }
                            },
                            "name": "y"
                          },
                          "computed": false,
                          "optional": false
                        },
                        "right": {
                          "type": "MemberExpression",
                          "start": 110,
                          "end": 115,
                          "loc": {
                            "start": {
                              "line": 7,
                              "column": 15
                            },
                            "end": {
                              "line": 7,
                              "column": 20
                            }
                          },
                          "object": {
                            "type": "Identifier",
                            "start": 110,
                            "end": 113,
                            "loc": {
                              "start": {
                                "line": 7,
                                "column": 15
                              },
                              "end": {
                                "line": 7,
                                "column": 18
                              }
                            },
                            "name": "obj"
                          },
                          "property": {
                            "type": "Identifier",
                            "start": 114,
                            "end": 115,
                            "loc": {
                              "start": {
                                "line": 7,
                                "column": 19
                              },
                              "end": {
                                "line": 7,
                                "column": 20
                              }
                            },
                            "name": "y"
                          },
                          "computed": false,
                          "optional": false
                        }
                      }
                    },
                    {
                      "type": "ExpressionStatement",
                      "start": 122,
                      "end": 136,
                      "loc": {
                        "start": {
                          "line": 8,
                          "column": 6
                        },
                        "end": {
                          "line": 8,
                          "column": 20
                        }
                      },
                      "expression": {
                        "type": "AssignmentExpression",
                        "start": 122,
                        "end": 136,
                        "loc": {
                          "start": {
                            "line": 8,
                            "column": 6
                          },
                          "end": {
                            "line": 8,
                            "column": 20
                          }
                        },
                        "operator": "=",
                        "left": {
                          "type": "MemberExpression",
                          "start": 122,
                          "end": 128,
                          "loc": {
                            "start": {
                              "line": 8,
                              "column": 6
                            },
                            "end": {
                              "line": 8,
                              "column": 12
                            }
                          },
                          "object": {
                            "type": "ThisExpression",
                            "start": 122,
                            "end": 126,
                            "loc": {
                              "start": {
                                "line": 8,
                                "column": 6
                              },
                              "end": {
                                "line": 8,
                                "column": 10
                              }
                            }
                          },
                          "property": {
                            "type": "Identifier",
                            "start": 127,
                            "end": 128,
                            "loc": {
                              "start": {
                                "line": 8,
                                "column": 11
                              },
                              "end": {
                                "line": 8,
                                "column": 12
                              }
                            },
                            "name": "z"
                          },
                          "computed": false,
                          "optional": false
                        },
                        "right": {
                          "type": "MemberExpression",
                          "start": 131,
                          "end": 136,
                          "loc": {
                            "start": {
                              "line": 8,
                              "column": 15
                            },
                            "end": {
                              "line": 8,
                              "column": 20
                            }
                          },
                          "object": {
                            "type": "Identifier",
                            "start": 131,
                            "end": 134,
                            "loc": {
                              "start": {
                                "line": 8,
                                "column": 15
                              },
                              "end": {
                                "line": 8,
                                "column": 18
                              }
                            },
                            "name": "obj"
                          },
                          "property": {
                            "type": "Identifier",
                            "start": 135,
                            "end": 136,
                            "loc": {
                              "start": {
                                "line": 8,
                                "column": 19
                              },
                              "end": {
                                "line": 8,
                                "column": 20
                              }
                            },
                            "name": "z"
                          },
                          "computed": false,
                          "optional": false
                        }
                      }
                    }
                  ]
                },
                "handler": {
                  "type": "CatchClause",
                  "start": 147,
                  "end": 194,
                  "loc": {
                    "start": {
                      "line": 10,
                      "column": 4
                    },
                    "end": {
                      "line": 13,
                      "column": 5
                    }
                  },
                  "param": null,
                  "body": {
                    "type": "BlockStatement",
                    "start": 153,
                    "end": 194,
                    "loc": {
                      "start": {
                        "line": 10,
                        "column": 10
                      },
                      "end": {
                        "line": 13,
                        "column": 5
                      }
                    },
                    "body": [
                      {
                        "type": "ExpressionStatement",
                        "start": 161,
                        "end": 171,
                        "loc": {
                          "start": {
                            "line": 11,
                            "column": 6
                          },
                          "end": {
                            "line": 11,
                            "column": 16
                          }
                        },
                        "expression": {
                          "type": "AssignmentExpression",
                          "start": 161,
                          "end": 171,
                          "loc": {
                            "start": {
                              "line": 11,
                              "column": 6
                            },
                            "end": {
                              "line": 11,
                              "column": 16
                            }
                          },
                          "operator": "=",
                          "left": {
                            "type": "MemberExpression",
                            "start": 161,
                            "end": 167,
                            "loc": {
                              "start": {
                                "line": 11,
                                "column": 6
                              },
                              "end": {
                                "line": 11,
                                "column": 12
                              }
                            },
                            "object": {
                              "type": "ThisExpression",
                              "start": 161,
                              "end": 165,
                              "loc": {
                                "start": {
                                  "line": 11,
                                  "column": 6
                                },
                                "end": {
                                  "line": 11,
                                  "column": 10
                                }
                              }
                            },
                            "property": {
                              "type": "Identifier",
                              "start": 166,
                              "end": 167,
                              "loc": {
                                "start": {
                                  "line": 11,
                                  "column": 11
                                },
                                "end": {
                                  "line": 11,
                                  "column": 12
                                }
                              },
                              "name": "y"
                            },
                            "computed": false,
                            "optional": false
                          },
                          "right": {
                            "type": "Literal",
                            "start": 170,
                            "end": 171,
                            "loc": {
                              "start": {
                                "line": 11,
                                "column": 15
                              },
                              "end": {
                                "line": 11,
                                "column": 16
                              }
                            },
                            "value": 0,
                            "raw": "0"
                          }
                        }
                      },
                      {
                        "type": "ExpressionStatement",
                        "start": 178,
                        "end": 188,
                        "loc": {
                          "start": {
                            "line": 12,
                            "column": 6
                          },
                          "end": {
                            "line": 12,
                            "column": 16
                          }
                        },
                        "expression": {
                          "type": "AssignmentExpression",
                          "start": 178,
                          "end": 188,
                          "loc": {
                            "start": {
                              "line": 12,
                              "column": 6
                            },
                            "end": {
                              "line": 12,
                              "column": 16
                            }
                          },
                          "operator": "=",
                          "left": {
                            "type": "MemberExpression",
                            "start": 178,
                            "end": 184,
                            "loc": {
                              "start": {
                                "line": 12,
                                "column": 6
                              },
                              "end": {
                                "line": 12,
                                "column": 12
                              }
                            },
                            "object": {
                              "type": "ThisExpression",
                              "start": 178,
                              "end": 182,
                              "loc": {
                                "start": {
                                  "line": 12,
                                  "column": 6
                                },
                                "end": {
                                  "line": 12,
                                  "column": 10
                                }
                              }
                            },
                            "property": {
                              "type": "Identifier",
                              "start": 183,
                              "end": 184,
                              "loc": {
                                "start": {
                                  "line": 12,
                                  "column": 11
                                },
                                "end": {
                                  "line": 12,
                                  "column": 12
                                }
                              },
                              "name": "z"
                            },
                            "computed": false,
                            "optional": false
                          },
                          "right": {
                            "type": "Literal",
                            "start": 187,
                            "end": 188,
                            "loc": {
                              "start": {
                                "line": 12,
                                "column": 15
                              },
                              "end": {
                                "line": 12,
                                "column": 16
                              }
                            },
                            "value": 0,
                            "raw": "0"
                          }
                        }
                      }
                    ]
                  }
                },
                "finalizer": null
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

func TestClassFeature86(t *testing.T) {
	ast, err := Compile(`class C {
  static y
  static #z
  static {
    const obj = {}
    this.y = obj.y
    this.#z = obj.z
  }
}`)
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 107,
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 9,
      "column": 1
    }
  },
  "body": [
    {
      "type": "ClassDeclaration",
      "start": 0,
      "end": 107,
      "loc": {
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 9,
          "column": 1
        }
      },
      "id": {
        "type": "Identifier",
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
        "name": "C"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "start": 8,
        "end": 107,
        "loc": {
          "start": {
            "line": 1,
            "column": 8
          },
          "end": {
            "line": 9,
            "column": 1
          }
        },
        "body": [
          {
            "type": "PropertyDefinition",
            "start": 12,
            "end": 20,
            "loc": {
              "start": {
                "line": 2,
                "column": 2
              },
              "end": {
                "line": 2,
                "column": 10
              }
            },
            "static": true,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 19,
              "end": 20,
              "loc": {
                "start": {
                  "line": 2,
                  "column": 9
                },
                "end": {
                  "line": 2,
                  "column": 10
                }
              },
              "name": "y"
            },
            "value": null
          },
          {
            "type": "PropertyDefinition",
            "start": 23,
            "end": 32,
            "loc": {
              "start": {
                "line": 3,
                "column": 2
              },
              "end": {
                "line": 3,
                "column": 11
              }
            },
            "static": true,
            "computed": false,
            "key": {
              "type": "PrivateIdentifier",
              "start": 30,
              "end": 32,
              "loc": {
                "start": {
                  "line": 3,
                  "column": 9
                },
                "end": {
                  "line": 3,
                  "column": 11
                }
              },
              "name": "z"
            },
            "value": null
          },
          {
            "type": "StaticBlock",
            "start": 35,
            "end": 105,
            "loc": {
              "start": {
                "line": 4,
                "column": 2
              },
              "end": {
                "line": 8,
                "column": 3
              }
            },
            "body": [
              {
                "type": "VariableDeclaration",
                "start": 48,
                "end": 62,
                "loc": {
                  "start": {
                    "line": 5,
                    "column": 4
                  },
                  "end": {
                    "line": 5,
                    "column": 18
                  }
                },
                "declarations": [
                  {
                    "type": "VariableDeclarator",
                    "start": 54,
                    "end": 62,
                    "loc": {
                      "start": {
                        "line": 5,
                        "column": 10
                      },
                      "end": {
                        "line": 5,
                        "column": 18
                      }
                    },
                    "id": {
                      "type": "Identifier",
                      "start": 54,
                      "end": 57,
                      "loc": {
                        "start": {
                          "line": 5,
                          "column": 10
                        },
                        "end": {
                          "line": 5,
                          "column": 13
                        }
                      },
                      "name": "obj"
                    },
                    "init": {
                      "type": "ObjectExpression",
                      "start": 60,
                      "end": 62,
                      "loc": {
                        "start": {
                          "line": 5,
                          "column": 16
                        },
                        "end": {
                          "line": 5,
                          "column": 18
                        }
                      },
                      "properties": []
                    }
                  }
                ],
                "kind": "const"
              },
              {
                "type": "ExpressionStatement",
                "start": 67,
                "end": 81,
                "loc": {
                  "start": {
                    "line": 6,
                    "column": 4
                  },
                  "end": {
                    "line": 6,
                    "column": 18
                  }
                },
                "expression": {
                  "type": "AssignmentExpression",
                  "start": 67,
                  "end": 81,
                  "loc": {
                    "start": {
                      "line": 6,
                      "column": 4
                    },
                    "end": {
                      "line": 6,
                      "column": 18
                    }
                  },
                  "operator": "=",
                  "left": {
                    "type": "MemberExpression",
                    "start": 67,
                    "end": 73,
                    "loc": {
                      "start": {
                        "line": 6,
                        "column": 4
                      },
                      "end": {
                        "line": 6,
                        "column": 10
                      }
                    },
                    "object": {
                      "type": "ThisExpression",
                      "start": 67,
                      "end": 71,
                      "loc": {
                        "start": {
                          "line": 6,
                          "column": 4
                        },
                        "end": {
                          "line": 6,
                          "column": 8
                        }
                      }
                    },
                    "property": {
                      "type": "Identifier",
                      "start": 72,
                      "end": 73,
                      "loc": {
                        "start": {
                          "line": 6,
                          "column": 9
                        },
                        "end": {
                          "line": 6,
                          "column": 10
                        }
                      },
                      "name": "y"
                    },
                    "computed": false,
                    "optional": false
                  },
                  "right": {
                    "type": "MemberExpression",
                    "start": 76,
                    "end": 81,
                    "loc": {
                      "start": {
                        "line": 6,
                        "column": 13
                      },
                      "end": {
                        "line": 6,
                        "column": 18
                      }
                    },
                    "object": {
                      "type": "Identifier",
                      "start": 76,
                      "end": 79,
                      "loc": {
                        "start": {
                          "line": 6,
                          "column": 13
                        },
                        "end": {
                          "line": 6,
                          "column": 16
                        }
                      },
                      "name": "obj"
                    },
                    "property": {
                      "type": "Identifier",
                      "start": 80,
                      "end": 81,
                      "loc": {
                        "start": {
                          "line": 6,
                          "column": 17
                        },
                        "end": {
                          "line": 6,
                          "column": 18
                        }
                      },
                      "name": "y"
                    },
                    "computed": false,
                    "optional": false
                  }
                }
              },
              {
                "type": "ExpressionStatement",
                "start": 86,
                "end": 101,
                "loc": {
                  "start": {
                    "line": 7,
                    "column": 4
                  },
                  "end": {
                    "line": 7,
                    "column": 19
                  }
                },
                "expression": {
                  "type": "AssignmentExpression",
                  "start": 86,
                  "end": 101,
                  "loc": {
                    "start": {
                      "line": 7,
                      "column": 4
                    },
                    "end": {
                      "line": 7,
                      "column": 19
                    }
                  },
                  "operator": "=",
                  "left": {
                    "type": "MemberExpression",
                    "start": 86,
                    "end": 93,
                    "loc": {
                      "start": {
                        "line": 7,
                        "column": 4
                      },
                      "end": {
                        "line": 7,
                        "column": 11
                      }
                    },
                    "object": {
                      "type": "ThisExpression",
                      "start": 86,
                      "end": 90,
                      "loc": {
                        "start": {
                          "line": 7,
                          "column": 4
                        },
                        "end": {
                          "line": 7,
                          "column": 8
                        }
                      }
                    },
                    "property": {
                      "type": "PrivateIdentifier",
                      "start": 91,
                      "end": 93,
                      "loc": {
                        "start": {
                          "line": 7,
                          "column": 9
                        },
                        "end": {
                          "line": 7,
                          "column": 11
                        }
                      },
                      "name": "z"
                    },
                    "computed": false,
                    "optional": false
                  },
                  "right": {
                    "type": "MemberExpression",
                    "start": 96,
                    "end": 101,
                    "loc": {
                      "start": {
                        "line": 7,
                        "column": 14
                      },
                      "end": {
                        "line": 7,
                        "column": 19
                      }
                    },
                    "object": {
                      "type": "Identifier",
                      "start": 96,
                      "end": 99,
                      "loc": {
                        "start": {
                          "line": 7,
                          "column": 14
                        },
                        "end": {
                          "line": 7,
                          "column": 17
                        }
                      },
                      "name": "obj"
                    },
                    "property": {
                      "type": "Identifier",
                      "start": 100,
                      "end": 101,
                      "loc": {
                        "start": {
                          "line": 7,
                          "column": 18
                        },
                        "end": {
                          "line": 7,
                          "column": 19
                        }
                      },
                      "name": "z"
                    },
                    "computed": false,
                    "optional": false
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

func TestClassFeature87(t *testing.T) {
	ast, err := Compile(`let zRead
class C {
  static #z
  static {
    zRead = () => this.#z
  }
}`)
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 74,
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 7,
      "column": 1
    }
  },
  "body": [
    {
      "type": "VariableDeclaration",
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
      "declarations": [
        {
          "type": "VariableDeclarator",
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
            "name": "zRead"
          },
          "init": null
        }
      ],
      "kind": "let"
    },
    {
      "type": "ClassDeclaration",
      "start": 10,
      "end": 74,
      "loc": {
        "start": {
          "line": 2,
          "column": 0
        },
        "end": {
          "line": 7,
          "column": 1
        }
      },
      "id": {
        "type": "Identifier",
        "start": 16,
        "end": 17,
        "loc": {
          "start": {
            "line": 2,
            "column": 6
          },
          "end": {
            "line": 2,
            "column": 7
          }
        },
        "name": "C"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "start": 18,
        "end": 74,
        "loc": {
          "start": {
            "line": 2,
            "column": 8
          },
          "end": {
            "line": 7,
            "column": 1
          }
        },
        "body": [
          {
            "type": "PropertyDefinition",
            "start": 22,
            "end": 31,
            "loc": {
              "start": {
                "line": 3,
                "column": 2
              },
              "end": {
                "line": 3,
                "column": 11
              }
            },
            "static": true,
            "computed": false,
            "key": {
              "type": "PrivateIdentifier",
              "start": 29,
              "end": 31,
              "loc": {
                "start": {
                  "line": 3,
                  "column": 9
                },
                "end": {
                  "line": 3,
                  "column": 11
                }
              },
              "name": "z"
            },
            "value": null
          },
          {
            "type": "StaticBlock",
            "start": 34,
            "end": 72,
            "loc": {
              "start": {
                "line": 4,
                "column": 2
              },
              "end": {
                "line": 6,
                "column": 3
              }
            },
            "body": [
              {
                "type": "ExpressionStatement",
                "start": 47,
                "end": 68,
                "loc": {
                  "start": {
                    "line": 5,
                    "column": 4
                  },
                  "end": {
                    "line": 5,
                    "column": 25
                  }
                },
                "expression": {
                  "type": "AssignmentExpression",
                  "start": 47,
                  "end": 68,
                  "loc": {
                    "start": {
                      "line": 5,
                      "column": 4
                    },
                    "end": {
                      "line": 5,
                      "column": 25
                    }
                  },
                  "operator": "=",
                  "left": {
                    "type": "Identifier",
                    "start": 47,
                    "end": 52,
                    "loc": {
                      "start": {
                        "line": 5,
                        "column": 4
                      },
                      "end": {
                        "line": 5,
                        "column": 9
                      }
                    },
                    "name": "zRead"
                  },
                  "right": {
                    "type": "ArrowFunctionExpression",
                    "start": 55,
                    "end": 68,
                    "loc": {
                      "start": {
                        "line": 5,
                        "column": 12
                      },
                      "end": {
                        "line": 5,
                        "column": 25
                      }
                    },
                    "id": null,
                    "expression": true,
                    "generator": false,
                    "async": false,
                    "params": [],
                    "body": {
                      "type": "MemberExpression",
                      "start": 61,
                      "end": 68,
                      "loc": {
                        "start": {
                          "line": 5,
                          "column": 18
                        },
                        "end": {
                          "line": 5,
                          "column": 25
                        }
                      },
                      "object": {
                        "type": "ThisExpression",
                        "start": 61,
                        "end": 65,
                        "loc": {
                          "start": {
                            "line": 5,
                            "column": 18
                          },
                          "end": {
                            "line": 5,
                            "column": 22
                          }
                        }
                      },
                      "property": {
                        "type": "PrivateIdentifier",
                        "start": 66,
                        "end": 68,
                        "loc": {
                          "start": {
                            "line": 5,
                            "column": 23
                          },
                          "end": {
                            "line": 5,
                            "column": 25
                          }
                        },
                        "name": "z"
                      },
                      "computed": false,
                      "optional": false
                    }
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

func TestClassFeature88(t *testing.T) {
	ast, err := Compile(`let zRead
class C {
  static #z
  static {
    zRead = (obj) => obj.#z
  }
}
zRead(new C())`)
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 91,
  "loc": {
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 8,
      "column": 14
    }
  },
  "body": [
    {
      "type": "VariableDeclaration",
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
      "declarations": [
        {
          "type": "VariableDeclarator",
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
            "name": "zRead"
          },
          "init": null
        }
      ],
      "kind": "let"
    },
    {
      "type": "ClassDeclaration",
      "start": 10,
      "end": 76,
      "loc": {
        "start": {
          "line": 2,
          "column": 0
        },
        "end": {
          "line": 7,
          "column": 1
        }
      },
      "id": {
        "type": "Identifier",
        "start": 16,
        "end": 17,
        "loc": {
          "start": {
            "line": 2,
            "column": 6
          },
          "end": {
            "line": 2,
            "column": 7
          }
        },
        "name": "C"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "start": 18,
        "end": 76,
        "loc": {
          "start": {
            "line": 2,
            "column": 8
          },
          "end": {
            "line": 7,
            "column": 1
          }
        },
        "body": [
          {
            "type": "PropertyDefinition",
            "start": 22,
            "end": 31,
            "loc": {
              "start": {
                "line": 3,
                "column": 2
              },
              "end": {
                "line": 3,
                "column": 11
              }
            },
            "static": true,
            "computed": false,
            "key": {
              "type": "PrivateIdentifier",
              "start": 29,
              "end": 31,
              "loc": {
                "start": {
                  "line": 3,
                  "column": 9
                },
                "end": {
                  "line": 3,
                  "column": 11
                }
              },
              "name": "z"
            },
            "value": null
          },
          {
            "type": "StaticBlock",
            "start": 34,
            "end": 74,
            "loc": {
              "start": {
                "line": 4,
                "column": 2
              },
              "end": {
                "line": 6,
                "column": 3
              }
            },
            "body": [
              {
                "type": "ExpressionStatement",
                "start": 47,
                "end": 70,
                "loc": {
                  "start": {
                    "line": 5,
                    "column": 4
                  },
                  "end": {
                    "line": 5,
                    "column": 27
                  }
                },
                "expression": {
                  "type": "AssignmentExpression",
                  "start": 47,
                  "end": 70,
                  "loc": {
                    "start": {
                      "line": 5,
                      "column": 4
                    },
                    "end": {
                      "line": 5,
                      "column": 27
                    }
                  },
                  "operator": "=",
                  "left": {
                    "type": "Identifier",
                    "start": 47,
                    "end": 52,
                    "loc": {
                      "start": {
                        "line": 5,
                        "column": 4
                      },
                      "end": {
                        "line": 5,
                        "column": 9
                      }
                    },
                    "name": "zRead"
                  },
                  "right": {
                    "type": "ArrowFunctionExpression",
                    "start": 55,
                    "end": 70,
                    "loc": {
                      "start": {
                        "line": 5,
                        "column": 12
                      },
                      "end": {
                        "line": 5,
                        "column": 27
                      }
                    },
                    "id": null,
                    "expression": true,
                    "generator": false,
                    "async": false,
                    "params": [
                      {
                        "type": "Identifier",
                        "start": 56,
                        "end": 59,
                        "loc": {
                          "start": {
                            "line": 5,
                            "column": 13
                          },
                          "end": {
                            "line": 5,
                            "column": 16
                          }
                        },
                        "name": "obj"
                      }
                    ],
                    "body": {
                      "type": "MemberExpression",
                      "start": 64,
                      "end": 70,
                      "loc": {
                        "start": {
                          "line": 5,
                          "column": 21
                        },
                        "end": {
                          "line": 5,
                          "column": 27
                        }
                      },
                      "object": {
                        "type": "Identifier",
                        "start": 64,
                        "end": 67,
                        "loc": {
                          "start": {
                            "line": 5,
                            "column": 21
                          },
                          "end": {
                            "line": 5,
                            "column": 24
                          }
                        },
                        "name": "obj"
                      },
                      "property": {
                        "type": "PrivateIdentifier",
                        "start": 68,
                        "end": 70,
                        "loc": {
                          "start": {
                            "line": 5,
                            "column": 25
                          },
                          "end": {
                            "line": 5,
                            "column": 27
                          }
                        },
                        "name": "z"
                      },
                      "computed": false,
                      "optional": false
                    }
                  }
                }
              }
            ]
          }
        ]
      }
    },
    {
      "type": "ExpressionStatement",
      "start": 77,
      "end": 91,
      "loc": {
        "start": {
          "line": 8,
          "column": 0
        },
        "end": {
          "line": 8,
          "column": 14
        }
      },
      "expression": {
        "type": "CallExpression",
        "start": 77,
        "end": 91,
        "loc": {
          "start": {
            "line": 8,
            "column": 0
          },
          "end": {
            "line": 8,
            "column": 14
          }
        },
        "callee": {
          "type": "Identifier",
          "start": 77,
          "end": 82,
          "loc": {
            "start": {
              "line": 8,
              "column": 0
            },
            "end": {
              "line": 8,
              "column": 5
            }
          },
          "name": "zRead"
        },
        "arguments": [
          {
            "type": "NewExpression",
            "start": 83,
            "end": 90,
            "loc": {
              "start": {
                "line": 8,
                "column": 6
              },
              "end": {
                "line": 8,
                "column": 13
              }
            },
            "callee": {
              "type": "Identifier",
              "start": 87,
              "end": 88,
              "loc": {
                "start": {
                  "line": 8,
                  "column": 10
                },
                "end": {
                  "line": 8,
                  "column": 11
                }
              },
              "name": "C"
            },
            "arguments": []
          }
        ],
        "optional": false
      }
    }
  ]
}
`, ast)
}

func TestClassFeature89(t *testing.T) {
	ast, err := Compile("class C { async\nget(){} }")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 25,
  "body": [
    {
      "type": "ClassDeclaration",
      "start": 0,
      "end": 25,
      "id": {
        "type": "Identifier",
        "start": 6,
        "end": 7,
        "name": "C"
      },
      "superClass": null,
      "body": {
        "type": "ClassBody",
        "start": 8,
        "end": 25,
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
              "end": 15,
              "name": "async"
            },
            "value": null
          },
          {
            "type": "MethodDefinition",
            "start": 16,
            "end": 23,
            "static": false,
            "computed": false,
            "key": {
              "type": "Identifier",
              "start": 16,
              "end": 19,
              "name": "get"
            },
            "kind": "method",
            "value": {
              "type": "FunctionExpression",
              "start": 19,
              "end": 23,
              "id": null,
              "expression": false,
              "generator": false,
              "async": false,
              "params": [],
              "body": {
                "type": "BlockStatement",
                "start": 21,
                "end": 23,
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
