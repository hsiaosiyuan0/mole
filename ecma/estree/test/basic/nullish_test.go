package estree_test

import (
	"testing"

	. "github.com/hsiaosiyuan0/mole/ecma/estree/test"
	. "github.com/hsiaosiyuan0/mole/internal"
)

func TestNullish1(t *testing.T) {
	ast, err := Compile("a ?? b")
	AssertEqual(t, nil, err, "should be prog ok")

	AssertEqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 6,
  "body": [
    {
      "type": "ExpressionStatement",
      "start": 0,
      "end": 6,
      "expression": {
        "type": "LogicalExpression",
        "start": 0,
        "end": 6,
        "left": {
          "type": "Identifier",
          "start": 0,
          "end": 1,
          "name": "a"
        },
        "operator": "??",
        "right": {
          "type": "Identifier",
          "start": 5,
          "end": 6,
          "name": "b"
        }
      }
    }
  ]
}
`, ast)
}

func TestNullish2(t *testing.T) {
	ast, err := Compile("a ?? b ?? c")
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
        "type": "LogicalExpression",
        "start": 0,
        "end": 11,
        "left": {
          "type": "LogicalExpression",
          "start": 0,
          "end": 6,
          "left": {
            "type": "Identifier",
            "start": 0,
            "end": 1,
            "name": "a"
          },
          "operator": "??",
          "right": {
            "type": "Identifier",
            "start": 5,
            "end": 6,
            "name": "b"
          }
        },
        "operator": "??",
        "right": {
          "type": "Identifier",
          "start": 10,
          "end": 11,
          "name": "c"
        }
      }
    }
  ]
}
`, ast)
}

func TestNullish3(t *testing.T) {
	ast, err := Compile("a | b ?? c | d")
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
        "type": "LogicalExpression",
        "start": 0,
        "end": 14,
        "left": {
          "type": "BinaryExpression",
          "start": 0,
          "end": 5,
          "left": {
            "type": "Identifier",
            "start": 0,
            "end": 1,
            "name": "a"
          },
          "operator": "|",
          "right": {
            "type": "Identifier",
            "start": 4,
            "end": 5,
            "name": "b"
          }
        },
        "operator": "??",
        "right": {
          "type": "BinaryExpression",
          "start": 9,
          "end": 14,
          "left": {
            "type": "Identifier",
            "start": 9,
            "end": 10,
            "name": "c"
          },
          "operator": "|",
          "right": {
            "type": "Identifier",
            "start": 13,
            "end": 14,
            "name": "d"
          }
        }
      }
    }
  ]
}
`, ast)
}

func TestNullish4(t *testing.T) {
	ast, err := Compile("a ?? b ? c : d")
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
        "type": "ConditionalExpression",
        "start": 0,
        "end": 14,
        "test": {
          "type": "LogicalExpression",
          "start": 0,
          "end": 6,
          "left": {
            "type": "Identifier",
            "start": 0,
            "end": 1,
            "name": "a"
          },
          "operator": "??",
          "right": {
            "type": "Identifier",
            "start": 5,
            "end": 6,
            "name": "b"
          }
        },
        "consequent": {
          "type": "Identifier",
          "start": 9,
          "end": 10,
          "name": "c"
        },
        "alternate": {
          "type": "Identifier",
          "start": 13,
          "end": 14,
          "name": "d"
        }
      }
    }
  ]
}
`, ast)
}

func TestNullish5(t *testing.T) {
	ast, err := Compile("(a || b) ?? c")
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
      "end": 13,
      "expression": {
        "type": "LogicalExpression",
        "start": 0,
        "end": 13,
        "left": {
          "type": "LogicalExpression",
          "start": 1,
          "end": 7,
          "left": {
            "type": "Identifier",
            "start": 1,
            "end": 2,
            "name": "a"
          },
          "operator": "||",
          "right": {
            "type": "Identifier",
            "start": 6,
            "end": 7,
            "name": "b"
          }
        },
        "operator": "??",
        "right": {
          "type": "Identifier",
          "start": 12,
          "end": 13,
          "name": "c"
        }
      }
    }
  ]
}
`, ast)
}

func TestNullish6(t *testing.T) {
	ast, err := Compile("a || (b ?? c)")
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
      "end": 13,
      "expression": {
        "type": "LogicalExpression",
        "start": 0,
        "end": 13,
        "left": {
          "type": "Identifier",
          "start": 0,
          "end": 1,
          "name": "a"
        },
        "operator": "||",
        "right": {
          "type": "LogicalExpression",
          "start": 6,
          "end": 12,
          "left": {
            "type": "Identifier",
            "start": 6,
            "end": 7,
            "name": "b"
          },
          "operator": "??",
          "right": {
            "type": "Identifier",
            "start": 11,
            "end": 12,
            "name": "c"
          }
        }
      }
    }
  ]
}
`, ast)
}

func TestNullish7(t *testing.T) {
	ast, err := Compile("(a && b) ?? c")
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
      "end": 13,
      "expression": {
        "type": "LogicalExpression",
        "start": 0,
        "end": 13,
        "left": {
          "type": "LogicalExpression",
          "start": 1,
          "end": 7,
          "left": {
            "type": "Identifier",
            "start": 1,
            "end": 2,
            "name": "a"
          },
          "operator": "&&",
          "right": {
            "type": "Identifier",
            "start": 6,
            "end": 7,
            "name": "b"
          }
        },
        "operator": "??",
        "right": {
          "type": "Identifier",
          "start": 12,
          "end": 13,
          "name": "c"
        }
      }
    }
  ]
}
`, ast)
}

func TestNullish8(t *testing.T) {
	ast, err := Compile("a && (b ?? c)")
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
      "end": 13,
      "expression": {
        "type": "LogicalExpression",
        "start": 0,
        "end": 13,
        "left": {
          "type": "Identifier",
          "start": 0,
          "end": 1,
          "name": "a"
        },
        "operator": "&&",
        "right": {
          "type": "LogicalExpression",
          "start": 6,
          "end": 12,
          "left": {
            "type": "Identifier",
            "start": 6,
            "end": 7,
            "name": "b"
          },
          "operator": "??",
          "right": {
            "type": "Identifier",
            "start": 11,
            "end": 12,
            "name": "c"
          }
        }
      }
    }
  ]
}
`, ast)
}

func TestNullish9(t *testing.T) {
	ast, err := Compile("(a ?? b) || c")
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
      "end": 13,
      "expression": {
        "type": "LogicalExpression",
        "start": 0,
        "end": 13,
        "left": {
          "type": "LogicalExpression",
          "start": 1,
          "end": 7,
          "left": {
            "type": "Identifier",
            "start": 1,
            "end": 2,
            "name": "a"
          },
          "operator": "??",
          "right": {
            "type": "Identifier",
            "start": 6,
            "end": 7,
            "name": "b"
          }
        },
        "operator": "||",
        "right": {
          "type": "Identifier",
          "start": 12,
          "end": 13,
          "name": "c"
        }
      }
    }
  ]
}
`, ast)
}

func TestNullish10(t *testing.T) {
	ast, err := Compile("a ?? (b || c)")
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
      "end": 13,
      "expression": {
        "type": "LogicalExpression",
        "start": 0,
        "end": 13,
        "left": {
          "type": "Identifier",
          "start": 0,
          "end": 1,
          "name": "a"
        },
        "operator": "??",
        "right": {
          "type": "LogicalExpression",
          "start": 6,
          "end": 12,
          "left": {
            "type": "Identifier",
            "start": 6,
            "end": 7,
            "name": "b"
          },
          "operator": "||",
          "right": {
            "type": "Identifier",
            "start": 11,
            "end": 12,
            "name": "c"
          }
        }
      }
    }
  ]
}
`, ast)
}

func TestNullish11(t *testing.T) {
	ast, err := Compile("(a ?? b) && c")
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
      "end": 13,
      "expression": {
        "type": "LogicalExpression",
        "start": 0,
        "end": 13,
        "left": {
          "type": "LogicalExpression",
          "start": 1,
          "end": 7,
          "left": {
            "type": "Identifier",
            "start": 1,
            "end": 2,
            "name": "a"
          },
          "operator": "??",
          "right": {
            "type": "Identifier",
            "start": 6,
            "end": 7,
            "name": "b"
          }
        },
        "operator": "&&",
        "right": {
          "type": "Identifier",
          "start": 12,
          "end": 13,
          "name": "c"
        }
      }
    }
  ]
}
`, ast)
}

func TestNullish12(t *testing.T) {
	ast, err := Compile("a ?? (b && c)")
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
      "end": 13,
      "expression": {
        "type": "LogicalExpression",
        "start": 0,
        "end": 13,
        "left": {
          "type": "Identifier",
          "start": 0,
          "end": 1,
          "name": "a"
        },
        "operator": "??",
        "right": {
          "type": "LogicalExpression",
          "start": 6,
          "end": 12,
          "left": {
            "type": "Identifier",
            "start": 6,
            "end": 7,
            "name": "b"
          },
          "operator": "&&",
          "right": {
            "type": "Identifier",
            "start": 11,
            "end": 12,
            "name": "c"
          }
        }
      }
    }
  ]
}
`, ast)
}
