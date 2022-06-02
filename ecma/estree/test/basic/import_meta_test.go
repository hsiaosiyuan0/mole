package estree_test

import (
	"testing"

	. "github.com/hsiaosiyuan0/mole/ecma/estree/test"
	. "github.com/hsiaosiyuan0/mole/util"
)

func TestImportMeta1(t *testing.T) {
	ast, err := Compile("import.meta")
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
        "type": "MetaProperty",
        "start": 0,
        "end": 11,
        "meta": {
          "type": "Identifier",
          "start": 0,
          "end": 6,
          "name": "import"
        },
        "property": {
          "type": "Identifier",
          "start": 7,
          "end": 11,
          "name": "meta"
        }
      }
    }
  ]
}
`, ast)
}

func TestImportMeta2(t *testing.T) {
	ast, err := Compile("import.meta.url")
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
        "type": "MemberExpression",
        "start": 0,
        "end": 15,
        "object": {
          "type": "MetaProperty",
          "start": 0,
          "end": 11,
          "meta": {
            "type": "Identifier",
            "start": 0,
            "end": 6,
            "name": "import"
          },
          "property": {
            "type": "Identifier",
            "start": 7,
            "end": 11,
            "name": "meta"
          }
        },
        "property": {
          "type": "Identifier",
          "start": 12,
          "end": 15,
          "name": "url"
        },
        "computed": false,
        "optional": false
      }
    }
  ]
}
`, ast)
}

func TestImportMeta3(t *testing.T) {
	ast, err := Compile("import.meta(s)")
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
        "type": "CallExpression",
        "start": 0,
        "end": 14,
        "callee": {
          "type": "MetaProperty",
          "start": 0,
          "end": 11,
          "meta": {
            "type": "Identifier",
            "start": 0,
            "end": 6,
            "name": "import"
          },
          "property": {
            "type": "Identifier",
            "start": 7,
            "end": 11,
            "name": "meta"
          }
        },
        "arguments": [
          {
            "type": "Identifier",
            "start": 12,
            "end": 13,
            "name": "s"
          }
        ],
        "optional": false
      }
    }
  ]
}
`, ast)
}

func TestImportMeta4(t *testing.T) {
	ast, err := Compile("import.meta?.(a)[b]")
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
        "type": "ChainExpression",
        "start": 0,
        "end": 19,
        "expression": {
          "type": "MemberExpression",
          "start": 0,
          "end": 19,
          "object": {
            "type": "CallExpression",
            "start": 0,
            "end": 16,
            "callee": {
              "type": "MetaProperty",
              "start": 0,
              "end": 11,
              "meta": {
                "type": "Identifier",
                "start": 0,
                "end": 6,
                "name": "import"
              },
              "property": {
                "type": "Identifier",
                "start": 7,
                "end": 11,
                "name": "meta"
              }
            },
            "arguments": [
              {
                "type": "Identifier",
                "start": 14,
                "end": 15,
                "name": "a"
              }
            ],
            "optional": true
          },
          "property": {
            "type": "Identifier",
            "start": 17,
            "end": 18,
            "name": "b"
          },
          "computed": true,
          "optional": false
        }
      }
    }
  ]
}
`, ast)
}

func TestImportMeta5(t *testing.T) {
	ast, err := Compile("import.meta?.a().c[d](e)")
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
        "type": "ChainExpression",
        "start": 0,
        "end": 24,
        "expression": {
          "type": "CallExpression",
          "start": 0,
          "end": 24,
          "callee": {
            "type": "MemberExpression",
            "start": 0,
            "end": 21,
            "object": {
              "type": "MemberExpression",
              "start": 0,
              "end": 18,
              "object": {
                "type": "CallExpression",
                "start": 0,
                "end": 16,
                "callee": {
                  "type": "MemberExpression",
                  "start": 0,
                  "end": 14,
                  "object": {
                    "type": "MetaProperty",
                    "start": 0,
                    "end": 11,
                    "meta": {
                      "type": "Identifier",
                      "start": 0,
                      "end": 6,
                      "name": "import"
                    },
                    "property": {
                      "type": "Identifier",
                      "start": 7,
                      "end": 11,
                      "name": "meta"
                    }
                  },
                  "property": {
                    "type": "Identifier",
                    "start": 13,
                    "end": 14,
                    "name": "a"
                  },
                  "computed": false,
                  "optional": true
                },
                "arguments": [],
                "optional": false
              },
              "property": {
                "type": "Identifier",
                "start": 17,
                "end": 18,
                "name": "c"
              },
              "computed": false,
              "optional": false
            },
            "property": {
              "type": "Identifier",
              "start": 19,
              "end": 20,
              "name": "d"
            },
            "computed": true,
            "optional": false
          },
          "arguments": [
            {
              "type": "Identifier",
              "start": 22,
              "end": 23,
              "name": "e"
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
