package pack

import (
	"testing"

	"github.com/hsiaosiyuan0/mole/util"
)

func TestSubpath(t *testing.T) {
	// 1
	sp, err := NewSubpath("./features/*.js", "./src/features/*.js")
	util.AssertEqual(t, nil, err, "should be ok")

	ok, m := sp.Match("./features/a.js", [][]string{{"default"}})
	util.AssertEqual(t, true, ok, "should be ok")
	util.AssertEqual(t, "./src/features/a.js", m, "should be ok")

	// 2
	sp, err = NewSubpath("./features/*.js", map[string]interface{}{
		"node":    "./feature-node.js",
		"default": "./feature.js",
	})
	util.AssertEqual(t, nil, err, "should be ok")

	ok, m = sp.Match("./features/a.js", [][]string{{"node"}})
	util.AssertEqual(t, true, ok, "should be ok")
	util.AssertEqual(t, "./feature-node.js", m, "should be ok")

	// 3
	sp, err = NewSubpath("./features/private-internal/*", nil)
	util.AssertEqual(t, nil, err, "should be ok")

	ok, m = sp.Match("./features/private-internal/a.js", [][]string{{"node"}})
	util.AssertEqual(t, true, ok, "should be ok")
	util.AssertEqual(t, "", m, "should be ok")
}

func TestSubpathGrp(t *testing.T) {
	// 1
	sg, err := NewSubpathGrp(map[string]interface{}{
		"node": map[string]interface{}{
			"import":  "./feature-node.mjs",
			"require": "./feature-node.cjs",
		},
		"default": "./feature.mjs",
	})
	util.AssertEqual(t, nil, err, "should be ok")

	ok, m := sg.Match(".", [][]string{{"node", "require"}})
	util.AssertEqual(t, true, ok, "should be ok")
	util.AssertEqual(t, "./feature-node.cjs", m, "should be ok")

	ok, m = sg.Match(".", [][]string{{"default"}})
	util.AssertEqual(t, true, ok, "should be ok")
	util.AssertEqual(t, "./feature.mjs", m, "should be ok")

	// 2
	sg, err = NewSubpathGrp(map[string]interface{}{
		"./features/*.js":               "./src/features/*.js",
		"./features/private-internal/*": nil,
	})
	util.AssertEqual(t, nil, err, "should be ok")

	ok, m = sg.Match("./features/m.js", [][]string{{"default"}})
	util.AssertEqual(t, true, ok, "should be ok")
	util.AssertEqual(t, "./src/features/m.js", m, "should be ok")

	ok, m = sg.Match("./features/private-internal/m.js", [][]string{{"default"}})
	util.AssertEqual(t, false, ok, "should be ok")

	// 3
	sg, err = NewSubpathGrp(map[string]interface{}{
		"#dep": map[string]interface{}{
			"node":    "dep-node-native",
			"default": "./dep-polyfill.js",
		},
	})
	util.AssertEqual(t, nil, err, "should be ok")

	ok, m = sg.Match("#dep", [][]string{{"default"}})
	util.AssertEqual(t, true, ok, "should be ok")
	util.AssertEqual(t, "./dep-polyfill.js", m, "should be ok")
}
