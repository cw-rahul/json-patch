package jsonpatch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func reformatJSON(j string) string {
	buf := new(bytes.Buffer)

	json.Indent(buf, []byte(j), "", "  ")

	return buf.String()
}

func compareJSON(a, b string) bool {
	// return Equal([]byte(a), []byte(b))

	var objA, objB interface{}
	json.Unmarshal([]byte(a), &objA)
	json.Unmarshal([]byte(b), &objB)

	// fmt.Printf("Comparing %#v\nagainst %#v\n", objA, objB)
	return reflect.DeepEqual(objA, objB)
}

func applyPatch(doc, patch string) (string, error) {
	obj, err := DecodePatch([]byte(patch))

	if err != nil {
		panic(err)
	}

	out, err := obj.Apply([]byte(doc))

	if err != nil {
		return "", err
	}

	return string(out), nil
}

func applyPatchImpl_allowRemove(doc, patch string) (string, error) {
	obj, err := DecodePatch([]byte(patch))

	if err != nil {
		panic(err)
	}

	out, err := obj.ApplyIndentAllowRemove([]byte(doc), "", true)

	if err != nil {
		return "", err
	}

	return string(out), nil
}

func applyPatchImpl_dontAllowRemove(doc, patch string) (string, error) {
	obj, err := DecodePatch([]byte(patch))

	if err != nil {
		panic(err)
	}

	out, err := obj.ApplyIndentAllowRemove([]byte(doc), "", false)

	if err != nil {
		return "", err
	}

	return string(out), nil
}

type Case struct {
	doc, patch, result string
}

func repeatedA(r int) string {
	var s string
	for i := 0; i < r; i++ {
		s += "A"
	}
	return s
}

var ExprCases = []Case{
	{
		`[ {"foo": [{"key":"bar","key2":"bar2"},{"key":"qux"},{"key":"baz"}]}]`,
		`[ { "op": "replace", "path": "/0/foo/{'key':'bar'}", "value": {"key2":"bar2"}}]`,
		`[ {"foo": [{"key2":"bar2"},{"key":"qux"},{"key":"baz"}]}]`,
	},
	{
		`[{"key":"bar","key2":"bar2"},{"key":"qux"},{"key":"baz"}]`,
		`[ { "op": "replace", "path": "/{'key':'bar'}", "value": {"key2":"bar2"}}]`,
		`[{"key2":"bar2"},{"key":"qux"},{"key":"baz"}]`,
	},
	{
		`[{"key":"bar","key2":"bar2"},{"key":"qux"},{"key":"baz"}]`,
		`[ { "op": "add", "path": "/{'key':'qux'}", "value": {"key2":"bar2"}}]`,
		`[{"key":"bar","key2":"bar2"},{"key2":"bar2"},{"key":"qux"},{"key":"baz"}]`,
	},
	{
		`[{"key":"bar","key2":"bar2"},{"key":"qux"},{"key":"baz"}]`,
		`[ { "op": "remove", "path": "/{'key':'bar'}"}]`,
		`[{"key":"qux"},{"key":"baz"}]`,
	},
	{
		`[{"key":"bar","key2":"bar2"},{"key":"qux"},{"key":"baz"}]`,
		`[ { "op": "remove", "path": "/{'key': 'bar'}"}]`,
		`[{"key":"qux"},{"key":"baz"}]`,
	},
	{
		`[{"key":"bar","key2":"bar2"},{"key":"qux"},{"key":"baz"}]`,
		`[ { "op": "remove", "path": "/{'key' :'bar'}"}]`,
		`[{"key":"qux"},{"key":"baz"}]`,
	},
	{
		`[{"key":"bar","key2":"bar2"},{"key":"qux"},{"key":"baz"}]`,
		`[ { "op": "test", "path": "/{'key':'bar'}/key2","value":"bar2"},{ "op": "remove", "path": "/{'key':'bar'}"}]`,
		`[{"key":"qux"},{"key":"baz"}]`,
	},
}
var Cases = []Case{
	{
		`{ "foo": "bar"}`,
		`[
         { "op": "add", "path": "/baz", "value": "qux" }
     ]`,
		`{
       "baz": "qux",
       "foo": "bar"
     }`,
	},
	{
		`{ "foo": [ "bar", "baz" ] }`,
		`[
     { "op": "add", "path": "/foo/1", "value": "qux" }
    ]`,
		`{ "foo": [ "bar", "qux", "baz" ] }`,
	},
	{
		`{ "foo": [ "bar", "baz" ] }`,
		`[
     { "op": "add", "path": "/foo/-1", "value": "qux" }
    ]`,
		`{ "foo": [ "bar", "baz", "qux" ] }`,
	},
	{
		`{ "baz": "qux", "foo": "bar" }`,
		`[ { "op": "remove", "path": "/baz" } ]`,
		`{ "foo": "bar" }`,
	},
	{
		`{ "foo": [ "bar", "qux", "baz" ] }`,
		`[ { "op": "remove", "path": "/foo/1" } ]`,
		`{ "foo": [ "bar", "baz" ] }`,
	},
	{
		`{ "baz": "qux", "foo": "bar" }`,
		`[ { "op": "replace", "path": "/baz", "value": "boo" } ]`,
		`{ "baz": "boo", "foo": "bar" }`,
	},
	{
		`{
     "foo": {
       "bar": "baz",
       "waldo": "fred"
     },
     "qux": {
       "corge": "grault"
     }
   }`,
		`[ { "op": "move", "from": "/foo/waldo", "path": "/qux/thud" } ]`,
		`{
     "foo": {
       "bar": "baz"
     },
     "qux": {
       "corge": "grault",
       "thud": "fred"
     }
   }`,
	},
	{
		`{ "foo": [ "all", "grass", "cows", "eat" ] }`,
		`[ { "op": "move", "from": "/foo/1", "path": "/foo/3" } ]`,
		`{ "foo": [ "all", "cows", "eat", "grass" ] }`,
	},
	{
		`{ "foo": [ "all", "grass", "cows", "eat" ] }`,
		`[ { "op": "move", "from": "/foo/1", "path": "/foo/2" } ]`,
		`{ "foo": [ "all", "cows", "grass", "eat" ] }`,
	},
	{
		`{ "foo": "bar" }`,
		`[ { "op": "add", "path": "/child", "value": { "grandchild": { } } } ]`,
		`{ "foo": "bar", "child": { "grandchild": { } } }`,
	},
	{
		`{ "foo": ["bar"] }`,
		`[ { "op": "add", "path": "/foo/-", "value": ["abc", "def"] } ]`,
		`{ "foo": ["bar", ["abc", "def"]] }`,
	},
	{
		`{ "foo": "bar", "qux": { "baz": 1, "bar": null } }`,
		`[ { "op": "remove", "path": "/qux/bar" } ]`,
		`{ "foo": "bar", "qux": { "baz": 1 } }`,
	},
	{
		`{ "foo": "bar" }`,
		`[ { "op": "add", "path": "/baz", "value": null } ]`,
		`{ "baz": null, "foo": "bar" }`,
	},
	{
		`{ "foo": ["bar"]}`,
		`[ { "op": "replace", "path": "/foo/0", "value": "baz"}]`,
		`{ "foo": ["baz"]}`,
	},
	{
		`{ "foo": ["bar","baz"]}`,
		`[ { "op": "replace", "path": "/foo/0", "value": "bum"}]`,
		`{ "foo": ["bum","baz"]}`,
	},
	{
		`{ "foo": ["bar","qux","baz"]}`,
		`[ { "op": "replace", "path": "/foo/1", "value": "bum"}]`,
		`{ "foo": ["bar", "bum","baz"]}`,
	},
	{
		`[ {"foo": ["bar","qux","baz"]}]`,
		`[ { "op": "replace", "path": "/0/foo/0", "value": "bum"}]`,
		`[ {"foo": ["bum","qux","baz"]}]`,
	},
	{
		`[ {"foo": ["bar","qux","baz"], "bar": ["qux","baz"]}]`,
		`[ { "op": "copy", "from": "/0/foo/0", "path": "/0/bar/0"}]`,
		`[ {"foo": ["bar","qux","baz"], "bar": ["bar", "qux", "baz"]}]`,
	},
	{
		`[ {"foo": ["bar","qux","baz"], "bar": ["qux","baz"]}]`,
		`[ { "op": "copy", "from": "/0/foo/0", "path": "/0/bar"}]`,
		`[ {"foo": ["bar","qux","baz"], "bar": "bar"}]`,
	},
	{
		`[ { "foo": {"bar": ["qux","baz"]}, "baz": {"qux": "bum"}}]`,
		`[ { "op": "copy", "from": "/0/foo/bar", "path": "/0/baz/bar"}]`,
		`[ { "baz": {"bar": ["qux","baz"], "qux":"bum"}, "foo": {"bar": ["qux","baz"]}}]`,
	},
	{
		`{ "foo": ["bar"]}`,
		`[{"op": "copy", "path": "/foo/0", "from": "/foo"}]`,
		`{ "foo": [["bar"], "bar"]}`,
	},
	{
		`{ "foo": ["bar","qux","baz"]}`,
		`[ { "op": "remove", "path": "/foo/-2"}]`,
		`{ "foo": ["bar", "baz"]}`,
	},
	{
		`{ "foo": []}`,
		`[ { "op": "add", "path": "/foo/-1", "value": "qux"}]`,
		`{ "foo": ["qux"]}`,
	},
	{
		`{ "bar": [{"baz": null}]}`,
		`[ { "op": "replace", "path": "/bar/0/baz", "value": 1 } ]`,
		`{ "bar": [{"baz": 1}]}`,
	},
	{
		`{ "bar": [{"baz": 1}]}`,
		`[ { "op": "replace", "path": "/bar/0/baz", "value": null } ]`,
		`{ "bar": [{"baz": null}]}`,
	},
	{
		`{ "bar": [null]}`,
		`[ { "op": "replace", "path": "/bar/0", "value": 1 } ]`,
		`{ "bar": [1]}`,
	},
	{
		`{ "bar": [1]}`,
		`[ { "op": "replace", "path": "/bar/0", "value": null } ]`,
		`{ "bar": [null]}`,
	},
	{
		fmt.Sprintf(`{ "foo": ["A", %q] }`, repeatedA(48)),
		// The wrapping quotes around 'A's are included in the copy
		// size, so each copy operation increases the size by 50 bytes.
		`[ { "op": "copy", "path": "/foo/-", "from": "/foo/1" },
		   { "op": "copy", "path": "/foo/-", "from": "/foo/1" }]`,
		fmt.Sprintf(`{ "foo": ["A", %q] }`, repeatedA(48)),
	},
	{
		`{"foo": ["bar", "baz", "qux"]}`,
		`[{"op": "add", "path": "/foo/-", "value": "baz"}]`,
		`{"foo": ["bar", "baz", "qux"]}`,
	},
	{
		`{"foo": []}`,
		`[{"op": "add", "path": "/foo/-", "value": "bar"}, {"op": "add", "path": "/foo/-", "value": "bar"}, {"op": "add", "path": "/foo/-", "value": "baz"}]`,
		`{"foo": ["bar", "baz"]}`,
	},
	{
		`{"foo": ["bar", "baz", "qux"]}`,
		`[{"op": "replace", "path": "/foo/1", "value": "bar"}]`,
		`{"foo": ["bar", "qux"]}`,
	},
	{
		`{"foo": ["bar", "baz", "qux"], "one": "bar"}`,
		`[{"op": "move", "from": "/one", "path": "/foo/-"}]`,
		`{"foo": ["bar", "baz", "qux"]}`,
	},
	{
		`{"foo": ["bar", "baz", "qux"], "one": "bar"}`,
		`[{"op": "copy", "from": "/one", "path": "/foo/-"}]`,
		`{"foo": ["bar", "baz", "qux"], "one": "bar"}`,
	},
	{
		`{"foo": [{"bar": "baz", "one": "two"}]}`,
		`[{"op": "add", "path": "/foo/-", "value": {"one": "two", "bar": "baz"}}]`,
		`{"foo": [{"bar": "baz", "one": "two"}]}`,
	},
	{
		`{"foo": [{"bar": "baz", "one": "two", "timestampUTC": "2019-08-26T00:00:00Z"}]}`,
		`[{"op": "add", "path": "/foo/-", "value": {"one": "two", "bar": "baz", "timestampUTC": "2019-11-08T00:00:00Z"}}]`,
		`{"foo": [{"bar": "baz", "one": "two", "timestampUTC": "2019-08-26T00:00:00Z"}]}`,
	},
	{
		`{"foo": [{"bar": "baz", "one": "two", "timestampUTC": "2019-08-26T00:00:00Z"}], "more": {"one": "two", "bar": "baz", "timestampUTC": "2019-11-08T00:00:00Z"}}`,
		`[{"op": "copy", "from": "/more", "path": "/foo/0"}]`,
		`{"foo": [{"bar": "baz", "one": "two", "timestampUTC": "2019-08-26T00:00:00Z"}], "more": {"one": "two", "bar": "baz", "timestampUTC": "2019-11-08T00:00:00Z"}}`,
	},
	{
		`{"foo": [{"bar": "baz", "qux": {"one": "two", "timestampUTC": "2019-08-26T00:00:00Z"}}]}`,
		`[{"op": "add", "path": "/foo/-", "value": {"bar": "baz", "qux": {"one": "two", "timestampUTC": "2019-11-08T00:00:00Z"}}}]`,
		`{"foo": [{"bar": "baz", "qux": {"one": "two", "timestampUTC": "2019-08-26T00:00:00Z"}}]}`,
	},
	{
		`[{
			"name": "ProhibitRemovableDevice",
			"type": "SCHEDULE",
			"version": "1.0",
			"timestampUTC": "2019-08-26T00:00:00Z",
			"path": "/scripting/scripting/command",
			"task": "/schedule/",
			"taskInput": "SKIPPED",
			"executeNow": "true",
			"schedule": "@every 7m45s"
		}]`,
		`[
			{
				"op": "add",
				"path": "/-",
				"value": {
					"version": "1.0",
					"type": "SCHEDULE",
					"name": "ProhibitRemovableDevice",
					"timestampUTC": "2019-08-27T00:00:00Z",
					"path": "/scripting/scripting/command",
					"task": "/schedule/",
					"taskInput": "SKIPPED",
					"executeNow": "true",
					"schedule": "@every 7m45s"
				}
			},
			{
				"op": "add",
				"path": "/-",
				"value": {
					"version": "1.0",
					"type": "SCHEDULE",
					"name": "ProhibitRemovableDevice",
					"timestampUTC": "2019-08-28T00:00:00Z",
					"path": "/scripting/scripting/command",
					"task": "/schedule/",
					"taskInput": "SKIPPED",
					"executeNow": "true",
					"schedule": "@every 10m"
				}
			}
		]`,
		`[
			{
				"name": "ProhibitRemovableDevice",
				"type": "SCHEDULE",
				"version": "1.0",
				"timestampUTC": "2019-08-26T00:00:00Z",
				"path": "/scripting/scripting/command",
				"task": "/schedule/",
				"taskInput": "SKIPPED",
				"executeNow": "true",
				"schedule": "@every 7m45s"
			},
			{
				"version": "1.0",
				"type": "SCHEDULE",
				"name": "ProhibitRemovableDevice",
				"timestampUTC": "2019-08-28T00:00:00Z",
				"path": "/scripting/scripting/command",
				"task": "/schedule/",
				"taskInput": "SKIPPED",
				"executeNow": "true",
				"schedule": "@every 10m"
			}
		]`,
	},
}

type BadCase struct {
	doc, patch string
}
type BadCases_remove struct {
	doc, patch string
}

var MutationTestCases = []BadCase{
	{
		`{ "foo": "bar", "qux": { "baz": 1, "bar": null } }`,
		`[ { "op": "remove", "path": "/qux/bar" } ]`,
	},
	{
		`{ "foo": "bar", "qux": { "baz": 1, "bar": null } }`,
		`[ { "op": "replace", "path": "/qux/baz", "value": null } ]`,
	},
}

var BadCases = []BadCase{
	{
		`{ "foo": "bar" }`,
		`[ { "op": "add", "path": "/baz/bat", "value": "qux" } ]`,
	},
	{
		`{ "a": { "b": { "d": 1 } } }`,
		`[ { "op": "remove", "path": "/a/b/c" } ]`,
	},
	{
		`{ "a": { "b": { "d": 1 } } }`,
		`[ { "op": "move", "from": "/a/b/c", "path": "/a/b/e" } ]`,
	},
	{
		`{ "a": { "b": [1] } }`,
		`[ { "op": "remove", "path": "/a/b/1" } ]`,
	},
	{
		`{ "a": { "b": [1] } }`,
		`[ { "op": "move", "from": "/a/b/1", "path": "/a/b/2" } ]`,
	},
	{
		`{ "foo": "bar" }`,
		`[ { "op": "add", "pathz": "/baz", "value": "qux" } ]`,
	},
	{
		`{ "foo": "bar" }`,
		`[ { "op": "add", "path": "", "value": "qux" } ]`,
	},
	{
		`{ "foo": ["bar","baz"]}`,
		`[ { "op": "replace", "path": "/foo/2", "value": "bum"}]`,
	},
	{
		`{ "foo": ["bar","baz"]}`,
		`[ { "op": "add", "path": "/foo/-4", "value": "bum"}]`,
	},
	{
		`{ "name":{ "foo": "bat", "qux": "bum"}}`,
		`[ { "op": "replace", "path": "/foo/bar", "value":"baz"}]`,
	},
	{
		`{ "foo": ["bar"]}`,
		`[ {"op": "add", "path": "/foo/2", "value": "bum"}]`,
	},
	{
		`{ "foo": []}`,
		`[ {"op": "remove", "path": "/foo/-"}]`,
	},
	{
		`{ "foo": []}`,
		`[ {"op": "remove", "path": "/foo/-1"}]`,
	},
	{
		`{ "foo": ["bar"]}`,
		`[ {"op": "remove", "path": "/foo/-2"}]`,
	},
	{
		`{}`,
		`[ {"op":null,"path":""} ]`,
	},
	{
		`{}`,
		`[ {"op":"add","path":null} ]`,
	},
	{
		`{}`,
		`[ { "op": "copy", "from": null }]`,
	},
	{
		`{ "foo": ["bar"]}`,
		`[{"op": "copy", "path": "/foo/6666666666", "from": "/"}]`,
	},
	// Accumulated copy size cannot exceed AccumulatedCopySizeLimit.
	{
		fmt.Sprintf(`{ "foo": ["A", %q] }`, repeatedA(49)),
		// The wrapping quotes around 'A's are included in the copy
		// size, so each copy operation increases the size by 51 bytes.
		`[ { "op": "copy", "path": "/foo/-", "from": "/foo/1" },
		   { "op": "copy", "path": "/foo/-", "from": "/foo/1" }]`,
	},
	// Can't move into an index greater than or equal to the size of the array
	{
		`{ "foo": [ "all", "grass", "cows", "eat" ] }`,
		`[ { "op": "move", "from": "/foo/1", "path": "/foo/4" } ]`,
	},
}

var BadCasesRemove = []BadCases_remove{
	{
		`{ "a": { "b": { "d": 1 } } }`,
		`[ { "op": "remove", "path": "/a/b/c" } ]`,
	},
	{
		`{ "a": { "b": [1] } }`,
		`[ { "op": "remove", "path": "/a/b/1" } ]`,
	},
	{
		`{ "foo": []}`,
		`[ {"op": "remove", "path": "/foo/-"}]`,
	},
	{
		`{ "foo": []}`,
		`[ {"op": "remove", "path": "/foo/-1"}]`,
	},
	{
		`{ "foo": ["bar"]}`,
		`[ {"op": "remove", "path": "/foo/-2"}]`,
	},
}

// This is not thread safe, so we cannot run patch tests in parallel.
func configureGlobals(accumulatedCopySizeLimit int64) func() {
	oldAccumulatedCopySizeLimit := AccumulatedCopySizeLimit
	AccumulatedCopySizeLimit = accumulatedCopySizeLimit
	return func() {
		AccumulatedCopySizeLimit = oldAccumulatedCopySizeLimit
	}
}

func TestAllCases(t *testing.T) {
	defer configureGlobals(int64(100))()
	for _, c := range ExprCases {
		out, err := applyPatch(c.doc, c.patch)

		if err != nil {
			t.Errorf("Unable to apply patch: %s", err)
		}

		if !compareJSON(out, c.result) {
			t.Errorf("Patch did not apply. Expected:\n%s\n\nActual:\n%s\nPatch:%s\n",
				reformatJSON(c.result), reformatJSON(out), reformatJSON(c.patch))
		}
	}

	for _, c := range Cases {
		out, err := applyPatch(c.doc, c.patch)

		if err != nil {
			t.Errorf("Unable to apply patch: %s", err)
		}

		if !compareJSON(out, c.result) {
			t.Errorf("Patch did not apply. Expected:\n%s\n\nActual:\n%s\nPatch:%s\n",
				reformatJSON(c.result), reformatJSON(out), reformatJSON(c.patch))
		}
	}

	for _, c := range MutationTestCases {
		out, err := applyPatch(c.doc, c.patch)

		if err != nil {
			t.Errorf("Unable to apply patch: %s", err)
		}

		if compareJSON(out, c.doc) {
			t.Errorf("Patch did not apply. Original:\n%s\n\nPatched:\n%s",
				reformatJSON(c.doc), reformatJSON(out))
		}
	}

	for _, c := range BadCases {
		_, err := applyPatch(c.doc, c.patch)

		if err == nil {
			t.Errorf("Patch %q should have failed to apply but it did not", c.patch)
		}
	}
}

func TestAllCasesImpl_allowRemove(t *testing.T) {
	defer configureGlobals(int64(100))()
	for _, c := range ExprCases {
		out, err := applyPatchImpl_allowRemove(c.doc, c.patch)

		if err != nil {
			t.Errorf("Unable to apply patch: %s", err)
		}

		if !compareJSON(out, c.result) {
			t.Errorf("Patch did not apply. Expected:\n%s\n\nActual:\n%s\nPatch:%s\n",
				reformatJSON(c.result), reformatJSON(out), reformatJSON(c.patch))
		}
	}

	for _, c := range Cases {
		out, err := applyPatchImpl_allowRemove(c.doc, c.patch)

		if err != nil {
			t.Errorf("Unable to apply patch: %s", err)
		}

		if !compareJSON(out, c.result) {
			t.Errorf("Patch did not apply. Expected:\n%s\n\nActual:\n%s\nPatch:%s\n",
				reformatJSON(c.result), reformatJSON(out), reformatJSON(c.patch))
		}
	}

	for _, c := range MutationTestCases {
		out, err := applyPatchImpl_allowRemove(c.doc, c.patch)

		if err != nil {
			t.Errorf("Unable to apply patch: %s", err)
		}

		if compareJSON(out, c.doc) {
			t.Errorf("Patch did not apply. Original:\n%s\n\nPatched:\n%s",
				reformatJSON(c.doc), reformatJSON(out))
		}
	}

	for _, c := range BadCasesRemove {
		_, err := applyPatchImpl_allowRemove(c.doc, c.patch)

		if err != nil {
			t.Errorf("Patch %q should have failed to apply but it did not", c.patch)
		}
	}
}

func TestAllCasesImpl_dontAllow(t *testing.T) {
	defer configureGlobals(int64(100))()
	for _, c := range ExprCases {
		out, err := applyPatchImpl_dontAllowRemove(c.doc, c.patch)

		if err != nil {
			t.Errorf("Unable to apply patch: %s", err)
		}

		if !compareJSON(out, c.result) {
			t.Errorf("Patch did not apply. Expected:\n%s\n\nActual:\n%s\nPatch:%s\n",
				reformatJSON(c.result), reformatJSON(out), reformatJSON(c.patch))
		}
	}

	for _, c := range Cases {
		out, err := applyPatchImpl_dontAllowRemove(c.doc, c.patch)

		if err != nil {
			t.Errorf("Unable to apply patch: %s", err)
		}

		if !compareJSON(out, c.result) {
			t.Errorf("Patch did not apply. Expected:\n%s\n\nActual:\n%s\nPatch:%s\n",
				reformatJSON(c.result), reformatJSON(out), reformatJSON(c.patch))
		}
	}

	for _, c := range MutationTestCases {
		out, err := applyPatchImpl_dontAllowRemove(c.doc, c.patch)

		if err != nil {
			t.Errorf("Unable to apply patch: %s", err)
		}

		if compareJSON(out, c.doc) {
			t.Errorf("Patch did not apply. Original:\n%s\n\nPatched:\n%s",
				reformatJSON(c.doc), reformatJSON(out))
		}
	}

	for _, c := range BadCases {
		_, err := applyPatchImpl_dontAllowRemove(c.doc, c.patch)

		if err == nil {
			t.Errorf("Patch %q should have failed to apply but it did not", c.patch)
		}
	}
}

type TestCase struct {
	doc, patch string
	result     bool
	failedPath string
}

var TestCases = []TestCase{
	{
		`{
      "baz": "qux",
      "foo": [ "a", 2, "c" ]
    }`,
		`[
      { "op": "test", "path": "/baz", "value": "qux" },
      { "op": "test", "path": "/foo/1", "value": 2 }
    ]`,
		true,
		"",
	},
	{
		`{ "baz": "qux" }`,
		`[ { "op": "test", "path": "/baz", "value": "bar" } ]`,
		false,
		"/baz",
	},
	{
		`{
      "baz": "qux",
      "foo": ["a", 2, "c"]
    }`,
		`[
      { "op": "test", "path": "/baz", "value": "qux" },
      { "op": "test", "path": "/foo/1", "value": "c" }
    ]`,
		false,
		"/foo/1",
	},
	{
		`{ "baz": "qux" }`,
		`[ { "op": "test", "path": "/foo", "value": 42 } ]`,
		false,
		"/foo",
	},
	{
		`{ "baz": "qux" }`,
		`[ { "op": "test", "path": "/foo", "value": null } ]`,
		true,
		"",
	},
	{
		`{ "foo": null }`,
		`[ { "op": "test", "path": "/foo", "value": null } ]`,
		true,
		"",
	},
	{
		`{ "foo": {} }`,
		`[ { "op": "test", "path": "/foo", "value": null } ]`,
		false,
		"/foo",
	},
	{
		`{ "foo": [] }`,
		`[ { "op": "test", "path": "/foo", "value": null } ]`,
		false,
		"/foo",
	},
	{
		`{ "baz/foo": "qux" }`,
		`[ { "op": "test", "path": "/baz~1foo", "value": "qux"} ]`,
		true,
		"",
	},
	{
		`{ "foo": [] }`,
		`[ { "op": "test", "path": "/foo"} ]`,
		false,
		"/foo",
	},
}

func TestAllTest(t *testing.T) {
	for _, c := range TestCases {
		_, err := applyPatch(c.doc, c.patch)

		if c.result && err != nil {
			t.Errorf("Testing failed when it should have passed: %s", err)
		} else if !c.result && err == nil {
			t.Errorf("Testing passed when it should have faild: %s", err)
		} else if !c.result {
			expected := fmt.Sprintf("testing value %s failed: test failed", c.failedPath)
			if err.Error() != expected {
				t.Errorf("Testing failed as expected but invalid message: expected [%s], got [%s]", expected, err)
			}
		}
	}
}

func TestAllTest1(t *testing.T) {
	for _, c := range TestCases {
		_, err := applyPatchImpl_allowRemove(c.doc, c.patch)
		if c.result && err != nil {
			t.Errorf("Testing failed when it should have passed: %s", err)
		} else if !c.result && err == nil {
			t.Errorf("Testing passed when it should have faild: %s", err)
		} else if !c.result {
			expected := fmt.Sprintf("testing value %s failed: test failed", c.failedPath)
			if err.Error() != expected {
				t.Errorf("Testing failed as expected but invalid message: expected [%s], got [%s]", expected, err)
			}
		}
	}
}
