package estree_test

import (
	"testing"

	. "github.com/hsiaosiyuan0/mole/ecma/estree/test"
	"github.com/hsiaosiyuan0/mole/ecma/parser"
	. "github.com/hsiaosiyuan0/mole/fuzz"
)

func TestJsonSuperSet1(t *testing.T) {
	ast, err := Compile("'\u2029'")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
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
        "type": "Literal",
        "start": 0,
        "end": 5,
        "value": " ",
        "raw": "' '"
      },
      "directive": " "
    }
  ]
}
`, ast)
}

func TestJsonSuperSet2(t *testing.T) {
	ast, err := Compile("'\\u2028'")
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
        "value": " ",
        "raw": "'\\u2028'"
      },
      "directive": "\\u2028"
    }
  ]
}
`, ast)
}

func TestJsonSuperSet3(t *testing.T) {
	ast, err := Compile("\"\u2028\"")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 5,
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
  "body": [
    {
      "type": "ExpressionStatement",
      "start": 0,
      "end": 5,
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
      "expression": {
        "type": "Literal",
        "start": 0,
        "end": 5,
        "value": " ",
        "raw": "\" \""
      },
      "directive": " "
    }
  ]
}
`, ast)
}

func TestJsonSuperSet4(t *testing.T) {
	ast, err := Compile("`\u2029`")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 5,
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
  "body": [
    {
      "type": "ExpressionStatement",
      "start": 0,
      "end": 5,
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
            "value": {
              "raw": " ",
              "cooked": " "
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

func TestJsonSuperSet5(t *testing.T) {
	ast, err := Compile("`\u2028`")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
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
              "raw": " ",
              "cooked": " "
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

func TestJsonSuperSet6(t *testing.T) {
	ast, err := Compile("\"\u2029\"")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
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
        "type": "Literal",
        "start": 0,
        "end": 5,
        "value": " ",
        "raw": "\" \""
      },
      "directive": " "
    }
  ]
}
`, ast)
}

func TestJsonSupersetFail1(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_JSON_SUPER_SET)
	TestFail(t, "\"\u2029\"", "Unterminated string constant at (1:0)", opts)
}

func TestJsonSupersetFail2(t *testing.T) {
	TestFail(t, "/\u2029/", "Unterminated regular expression at (1:1)", nil)
}

func TestJsonSupersetFail3(t *testing.T) {
	TestFail(t, "/\u2028/", "Unterminated regular expression at (1:1)", nil)
}
