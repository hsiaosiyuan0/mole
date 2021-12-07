package estree_test

import (
	"testing"

	"github.com/hsiaosiyuan0/mole/pkg/assert"
)

func TestNumSep1(t *testing.T) {
	ast, err := compile("123_456")
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
	ast, err := compile("123_456.123_456e+123_456")
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
	ast, err := compile("0b1010_0001")
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
	ast, err := compile("0xDEAD_BEAF")
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
	ast, err := compile("0o755_666")
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
	ast, err := compile("123_456n")
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
	ast, err := compile(".012_345")
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
