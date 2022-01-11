package estree_test

import (
	"testing"

	. "github.com/hsiaosiyuan0/mole/ecma/estree/test"
	. "github.com/hsiaosiyuan0/mole/internal"
)

func TestLogicAssign1(t *testing.T) {
	ast, err := Compile("a &&= b")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
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
        "type": "AssignmentExpression",
        "start": 0,
        "end": 7,
        "operator": "&&=",
        "left": {
          "type": "Identifier",
          "start": 0,
          "end": 1,
          "name": "a"
        },
        "right": {
          "type": "Identifier",
          "start": 6,
          "end": 7,
          "name": "b"
        }
      }
    }
  ]
}
`, ast)
}

func TestLogicAssign2(t *testing.T) {
	ast, err := Compile("a ||= b")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
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
        "type": "AssignmentExpression",
        "start": 0,
        "end": 7,
        "operator": "||=",
        "left": {
          "type": "Identifier",
          "start": 0,
          "end": 1,
          "name": "a"
        },
        "right": {
          "type": "Identifier",
          "start": 6,
          "end": 7,
          "name": "b"
        }
      }
    }
  ]
}
`, ast)
}

func TestLogicAssign3(t *testing.T) {
	ast, err := Compile("a ??= b")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
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
        "type": "AssignmentExpression",
        "start": 0,
        "end": 7,
        "operator": "??=",
        "left": {
          "type": "Identifier",
          "start": 0,
          "end": 1,
          "name": "a"
        },
        "right": {
          "type": "Identifier",
          "start": 6,
          "end": 7,
          "name": "b"
        }
      }
    }
  ]
}
`, ast)
}

func TestLogicAssign4(t *testing.T) {
	ast, err := Compile("a &&= b ||= c ??= d")
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
        "operator": "&&=",
        "left": {
          "type": "Identifier",
          "start": 0,
          "end": 1,
          "name": "a"
        },
        "right": {
          "type": "AssignmentExpression",
          "start": 6,
          "end": 19,
          "operator": "||=",
          "left": {
            "type": "Identifier",
            "start": 6,
            "end": 7,
            "name": "b"
          },
          "right": {
            "type": "AssignmentExpression",
            "start": 12,
            "end": 19,
            "operator": "??=",
            "left": {
              "type": "Identifier",
              "start": 12,
              "end": 13,
              "name": "c"
            },
            "right": {
              "type": "Identifier",
              "start": 18,
              "end": 19,
              "name": "d"
            }
          }
        }
      }
    }
  ]
}
`, ast)
}
