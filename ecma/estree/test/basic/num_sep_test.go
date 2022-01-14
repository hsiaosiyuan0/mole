package estree_test

import (
	"testing"

	. "github.com/hsiaosiyuan0/mole/ecma/estree/test"
	. "github.com/hsiaosiyuan0/mole/fuzz"
)

func TestNumSep1(t *testing.T) {
	ast, err := Compile("123_456")
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
        "type": "Literal",
        "start": 0,
        "end": 7,
        "value": 123456,
        "raw": "123_456"
      }
    }
  ]
}
`, ast)
}

func TestNumSep2(t *testing.T) {
	ast, err := Compile("123_456.123_456e+123_456")
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
        "type": "Literal",
        "start": 0,
        "end": 24,
        "value": 0,
        "raw": "123_456.123_456e+123_456"
      }
    }
  ]
}
`, ast)
}

func TestNumSep3(t *testing.T) {
	ast, err := Compile("0b1010_0001")
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
        "type": "Literal",
        "start": 0,
        "end": 11,
        "value": 161,
        "raw": "0b1010_0001"
      }
    }
  ]
}
`, ast)
}

func TestNumSep4(t *testing.T) {
	ast, err := Compile("0xDEAD_BEAF")
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
        "type": "Literal",
        "start": 0,
        "end": 11,
        "value": 3735928495,
        "raw": "0xDEAD_BEAF"
      }
    }
  ]
}
`, ast)
}

func TestNumSep5(t *testing.T) {
	ast, err := Compile("0o755_666")
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
        "type": "Literal",
        "start": 0,
        "end": 9,
        "value": 252854,
        "raw": "0o755_666"
      }
    }
  ]
}
`, ast)
}

func TestNumSep6(t *testing.T) {
	ast, err := Compile("123_456n")
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
        "value": 123456,
        "raw": "123_456n"
      }
    }
  ]
}
`, ast)
}

func TestNumSep7(t *testing.T) {
	ast, err := Compile(".012_345")
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
        "value": 0.012345,
        "raw": ".012_345"
      }
    }
  ]
}
`, ast)
}
