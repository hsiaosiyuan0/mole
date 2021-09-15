#!/usr/bin/env bash

except=$(cat <<END
{
  "type": "Program",
  "loc": {
    "source": "",
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 4
    },
    "range": {
      "start": 0,
      "end": 10
    }
  },
  "sourceType": "",
  "body": [
    {
      "type": "ExpressionStatement",
      "loc": {
        "source": "",
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 4
        },
        "range": {
          "start": 0,
          "end": 10
        }
      },
      "expression": {
        "type": "NewExpression",
        "loc": {
          "source": "",
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 4
          },
          "range": {
            "start": 0,
            "end": 10
          }
        },
        "callee": {
          "type": "Identifier",
          "loc": null,
          "name": "Object"
        },
        "arguments": null
      }
    }
  ]
}
END
)

actual=$(cat <<END
{
  "type": "Program",
  "loc": {
    "source": "",
    "start": {
      "line": 1,
      "column": 0
    },
    "end": {
      "line": 1,
      "column": 4
    },
    "range": {
      "start": 0,
      "end": 10
    }
  },
  "sourceType": "",
  "body": [
    {
      "type": "ExpressionStatement",
      "loc": {
        "source": "",
        "start": {
          "line": 1,
          "column": 0
        },
        "end": {
          "line": 1,
          "column": 4
        },
        "range": {
          "start": 0,
          "end": 10
        }
      },
      "expression": {
        "type": "NewExpression",
        "loc": {
          "source": "",
          "start": {
            "line": 1,
            "column": 0
          },
          "end": {
            "line": 1,
            "column": 4
          },
          "range": {
            "start": 0,
            "end": 10
          }
        },
        "callee": {
          "type": "Identifier",
          "loc": null,
          "name": "Object"
        },
        "arguments": null1
      }
    }
  ]
}
END
)

/usr/local/Cellar/diffutils/3.8/bin/diff --color -u <(echo "$except") <(echo "$actual")

# echo "$actual" | git diff --no-index aa -