package jsonpatch

import (
	"testing"
)

var N = 100 // number of patch apply iterations

func TestAdd(t *testing.T) {
	reference, err := apply(original, addPatches)
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < N; i++ {
		p, err := apply(original, addPatches)
		if err != nil {
			t.Fatal(err)
		}

		if p != reference {
			t.Errorf("%s should be equal to %s", p, reference)
		}
	}
}

func apply(original, p string) (string, error) {
	patch, err := DecodePatch([]byte(p))
	if err != nil {
		return "", err
	}

	patched, err := patch.Apply([]byte(original))
	if err != nil {
		return "", err
	}

	return string(patched), nil
}

var original = `{"rules": []}`

var addPatches = `
[
  {
    "op": "add",
    "path": "/rules/-",
    "value": {
      "collector": "winlog",
      "description": "Uninterruptible Power Supply - AVR Trim Active.",
      "eventDetails": {
        "facility": 18,
        "source": "APCPBEAgent",
        "eventIDs": [
          "2007"
        ],
        "severity": [
          3,
          4,
          5
        ]
      }
    }
  },
  {
    "op": "add",
    "path": "/rules/-",
    "value": {
      "collector": "winlog",
      "description": "Uninterruptible Power Supply - AVR Trim Active_closing",
      "eventDetails": {
        "facility": 18,
        "source": "APCPBEAgent",
        "eventIDs": [
          "1056",
          "1060"
        ],
        "severity": [
          5
        ]
      }
    }
  },
  {
    "op": "add",
    "path": "/rules/-",
    "value": {
      "collector": "winlog",
      "description": "Uninterruptible Power Supply - AVR Trim Active.",
      "eventDetails": {
        "facility": 18,
        "source": "APCPBEAgent",
        "eventIDs": [
          "2007"
        ],
        "severity": [
          3,
          4,
          5
        ]
      }
    }
  },
  {
    "op": "add",
    "path": "/rules/-",
    "value": {
      "collector": "winlog",
      "description": "MSSQL - Failed to complete Backup of Database_Rule",
      "eventDetails": {
        "facility": 18,
        "source": "MSSQLSERVER",
        "eventIDs": [
          "3041",
          "18210",
          "18204",
          "12291"
        ],
        "severity": [
          3
        ]
      }
    }
  },
  {
    "op": "add",
    "path": "/rules/-",
    "value": {
      "collector": "winlog",
      "description": "MSSQL - Failed to complete Backup of Database_Rule",
      "eventDetails": {
        "facility": 18,
        "source": "MSSQLSERVER",
        "eventIDs": [
          "12289"
        ],
        "severity": [
          3,
          5
        ]
      }
    }
  },
  {
    "op": "add",
    "path": "/rules/-",
    "value": {
      "collector": "winlog",
      "description": "MSSQL - Another one",
      "eventDetails": {
        "facility": 18,
        "source": "MSSQLSERVER",
        "eventIDs": [
          "12289",
          "2"
        ],
        "severity": [
          3,
          5
        ]
      }
    }
  },
  {
    "op": "add",
    "path": "/rules/-",
    "value": {
      "collector": "winlog",
      "description": "MSSQL - Another one",
      "eventDetails": {
        "facility": 18,
        "foo": "bar",
        "source": "MSSQLSERVER",
        "eventIDs": [
          "12289",
          "2"
        ],
        "severity": [
          3,
          5
        ]
      }
    }
  },
  {
  	"op": "copy",
  	"from": "/rules/0",
  	"path": "/rules/-"
  },
  {
  	"op": "replace",
  	"path": "/rules/1",
  	"value": {
      "collector": "winlog",
      "description": "Uninterruptible Power Supply - AVR Trim Active_closing",
      "eventDetails": {
        "facility": 18,
        "source": "APCPBEAgent",
        "eventIDs": [
          "1056",
          "1060",
          "42"
        ],
        "severity": [
          5
        ]
      }
    }
  },
  {
  	"op": "move",
  	"from": "/rules/0",
  	"path": "/rules/-"
  },
  {
  	"op": "remove",
  	"path": "/rules/0"
  }
]
`
