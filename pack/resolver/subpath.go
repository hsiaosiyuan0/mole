package resolver

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/hsiaosiyuan0/mole/util"
)

// `Subpath` is used to handle these concepts:
// - [subpath patterns](https://nodejs.org/api/packages.html#subpath-patterns)
// - [browser](https://github.com/defunctzombie/package-browser-field-spec)
//
// the quick view of the cases handled by this module can be found at:
// https://github.com/hsiaosiyuan0/mole/issues/23
type Subpath struct {
	pat  interface{}            // string|*regexp.Regexp|false
	cond map[string]interface{} // condition => replacement
}

func NewSubpath(src string, cond interface{}) (*Subpath, error) {
	switch cv := cond.(type) {
	case string:
		// normalize input like `{ ".": "./index.js" }`
		return NewSubpath(src, map[string]interface{}{
			"default": cv,
		})

	case map[string]interface{}:
		if err := normalizeCond(cv); err != nil {
			return nil, err
		}

		pat, err := compileSubpath(src)
		if err != nil {
			return nil, err
		}
		return &Subpath{pat, cv}, nil

	case []interface{}:
		if len(cv) == 0 {
			return nil, errors.New(fmt.Sprintf("deformed condition: %v", cond))
		}
		return NewSubpath(src, cv[0])

	case nil, bool:
		// the browser spec use `false` to indicate the module should be ignored:
		//
		// "browser": {
		//   "module-a": false, // same as the `null` in subpath patterns
		//   "./server/only.js": "./shims/server-only.js"
		// }
		pat, err := compileSubpath(src)
		if err != nil {
			return nil, err
		}
		return &Subpath{pat, nil}, nil

	default:
		return nil, errors.New(fmt.Sprintf("deformed condition: %v", cond))
	}
}

func (m *Subpath) Match(t string, conditions [][]string) (bool, string) {
	mc := false
	var mcs []string
	switch v := m.pat.(type) {
	case string:
		mc = t == v
	case *regexp.Regexp:
		mcs = v.FindStringSubmatch(t)
		mc = len(mcs) > 0
	}
	if !mc {
		return false, ""
	}

	if m.cond == nil {
		return true, ""
	}

	var rhs interface{}
	for _, cond := range conditions {
		rhs = util.GetByPath(m.cond, cond)
		if rhs != nil {
			break
		}
	}
	if rhs == nil {
		return false, ""
	}

	rv := rhs.(*CondRhs)
	if rv.glob {
		if len(mcs) != 2 {
			return false, ""
		}
		return true, strings.Replace(rv.sp, "*", mcs[1], -1)
	}
	return true, rv.sp
}

func compileSubpath(p string) (interface{}, error) {
	if strings.Index(p, "*") == -1 {
		return p, nil
	}
	pat := strings.ReplaceAll(p, "*", "(.*?)")
	pat = strings.ReplaceAll(pat, "$", "\\$")
	pat = "^" + pat + "$"
	return regexp.Compile(pat)
}

type CondRhs struct {
	glob bool
	sp   string
}

func normalizeCond(cond map[string]interface{}) error {
	for k, v := range cond {
		switch vv := v.(type) {
		case string:
			cond[k] = &CondRhs{strings.Index(vv, "*") != -1, vv}
		case nil:
			cond[k] = &CondRhs{false, ""}
		case map[string]interface{}:
			if err := normalizeCond(vv); err != nil {
				return err
			}
		default:
			return errors.New(fmt.Sprintf("deformed condition: %v", vv))
		}
	}
	return nil
}

type SubpathGrp struct {
	neg []*Subpath
	pos []*Subpath
}

// tacit-subpath means the lhs of the subpath is not a relative path, eg:
//
// ```
// {
//   "exports": {
//     "import": "./index-module.js",
//     "require": "./index-require.cjs"
//   }
// }
// ```
//
// should be normalized to
//
// {
//   "exports": {
//     ".": {
//       "import": "./index-module.js",
//       "require": "./index-require.cjs"
//     }
//   }
// }

var conditions = []string{"node-addons", "node", "import", "require", "default", "browser", "electron"}

func isTacitSubpath(c map[string]interface{}) bool {
	for key := range c {
		if util.Includes(conditions, key) {
			return true
		}
	}
	return false
}

func NormalizeSubpath(c interface{}) (map[string]interface{}, error) {
	var cm map[string]interface{}

	switch vc := c.(type) {
	case string:
		cm = map[string]interface{}{
			".": map[string]interface{}{
				"default": vc,
			},
		}
	case map[string]interface{}:
		if isTacitSubpath(vc) {
			cm = map[string]interface{}{
				".": vc,
			}
		} else {
			cm = vc
		}
	default:
		return nil, errors.New(fmt.Sprintf("deformed subpath group: %v", c))
	}

	return cm, nil
}

func NewSubpathGrp(c interface{}) (*SubpathGrp, error) {
	cm, err := NormalizeSubpath(c)
	if err != nil {
		return nil, err
	}

	sg := &SubpathGrp{
		neg: []*Subpath{},
		pos: []*Subpath{},
	}

	for src, cond := range cm {
		s, err := NewSubpath(src, cond)
		if err != nil {
			return nil, err
		}

		if cond == nil || cond == false {
			sg.neg = append(sg.neg, s)
		} else {
			sg.pos = append(sg.pos, s)
		}
	}

	return sg, nil
}

func (sg *SubpathGrp) IsEmpty() bool {
	return len(sg.pos) == 0 && len(sg.neg) == 0
}

func (sg *SubpathGrp) Match(t string, conditions [][]string) (pos, neg bool, m string) {
	for _, s := range sg.neg {
		ok, _ := s.Match(t, conditions)
		if ok {
			neg = true
			return
		}
	}
	for _, s := range sg.pos {
		ok, mm := s.Match(t, conditions)
		if ok {
			pos = true
			m = mm
			break
		}
	}
	return
}
