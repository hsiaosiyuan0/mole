package linter

import (
	"testing"

	"github.com/hsiaosiyuan0/mole/util"
)

func TestNoEmpty1(t *testing.T) {
	r := lint(t, "try {} catch (ex) {throw ex}", &NoEmpty{})
	util.AssertEqual(t, 0, len(r.Abnormals), "should be ok")
	util.AssertEqual(t, 1, len(r.Diagnoses), "should be ok")
	util.AssertEqual(t, "disallow empty block statements", r.Diagnoses[0].Msg, "should be ok")
}

func TestNoEmpty2(t *testing.T) {
	r := lint(t, "try { foo() } catch (ex) {}", &NoEmpty{})
	util.AssertEqual(t, 0, len(r.Abnormals), "should be ok")
	util.AssertEqual(t, 1, len(r.Diagnoses), "should be ok")
	util.AssertEqual(t, "disallow empty block statements", r.Diagnoses[0].Msg, "should be ok")
}

func TestNoEmpty3(t *testing.T) {
	r := lint(t, "try { foo() } catch (ex) {throw ex} finally {}", &NoEmpty{})
	util.AssertEqual(t, 0, len(r.Abnormals), "should be ok")
	util.AssertEqual(t, 1, len(r.Diagnoses), "should be ok")
	util.AssertEqual(t, "disallow empty block statements", r.Diagnoses[0].Msg, "should be ok")
}

func TestNoEmpty4(t *testing.T) {
	r := lint(t, "if (foo) {}", &NoEmpty{})
	util.AssertEqual(t, 0, len(r.Abnormals), "should be ok")
	util.AssertEqual(t, 1, len(r.Diagnoses), "should be ok")
	util.AssertEqual(t, "disallow empty block statements", r.Diagnoses[0].Msg, "should be ok")
}

func TestNoEmpty5(t *testing.T) {
	r := lint(t, "while (foo) {}", &NoEmpty{})
	util.AssertEqual(t, 0, len(r.Abnormals), "should be ok")
	util.AssertEqual(t, 1, len(r.Diagnoses), "should be ok")
	util.AssertEqual(t, "disallow empty block statements", r.Diagnoses[0].Msg, "should be ok")
}

func TestNoEmpty6(t *testing.T) {
	r := lint(t, "for (;foo;) {}", &NoEmpty{})
	util.AssertEqual(t, 0, len(r.Abnormals), "should be ok")
	util.AssertEqual(t, 1, len(r.Diagnoses), "should be ok")
	util.AssertEqual(t, "disallow empty block statements", r.Diagnoses[0].Msg, "should be ok")
}

func TestNoEmpty7(t *testing.T) {
	r := lint(t, "switch(foo) {}", &NoEmpty{})
	util.AssertEqual(t, 0, len(r.Abnormals), "should be ok")
	util.AssertEqual(t, 1, len(r.Diagnoses), "should be ok")
	util.AssertEqual(t, "disallow empty block statements", r.Diagnoses[0].Msg, "should be ok")
}

func TestNoEmpty8(t *testing.T) {
	r := lintCb(t, `
  try { /* 1 */ } catch (ex) {}
  `, func(c *Config) {
		c.Rules["no-empty"] = []interface{}{
			2,
			map[string]interface{}{
				"allowEmptyCatch": true,
			},
		}
	}, &NoEmpty{})
	util.AssertEqual(t, 0, len(r.Abnormals), "should be ok")
	util.AssertEqual(t, 0, len(r.Diagnoses), "should be ok")
}

func TestNoEmpty9(t *testing.T) {
	r := lintCb(t, `
  try { foo(); } catch (ex) {} finally { /* 1 */ }
  `, func(c *Config) {
		c.Rules["no-empty"] = []interface{}{
			2,
			map[string]interface{}{
				"allowEmptyCatch": true,
			},
		}
	}, &NoEmpty{})
	util.AssertEqual(t, 0, len(r.Abnormals), "should be ok")
	util.AssertEqual(t, 0, len(r.Diagnoses), "should be ok")
}

func TestNoEmpty10(t *testing.T) {
	r := lintCb(t, `
  try {} catch (ex) {} finally {}
  `, func(c *Config) {
		c.Rules["no-empty"] = []interface{}{
			2,
			map[string]interface{}{
				"allowEmptyCatch": true,
			},
		}
	}, &NoEmpty{})
	util.AssertEqual(t, 0, len(r.Abnormals), "should be ok")
	util.AssertEqual(t, 2, len(r.Diagnoses), "should be ok")
}

func TestNoEmpty11(t *testing.T) {
	r := lintCb(t, `
  try {} catch (ex) {} finally {}
  `, func(c *Config) {
		c.Rules["no-empty"] = []interface{}{
			2,
			map[string]interface{}{
				"allowEmptyCatch": false,
			},
		}
	}, &NoEmpty{})
	util.AssertEqual(t, 0, len(r.Abnormals), "should be ok")
	util.AssertEqual(t, 3, len(r.Diagnoses), "should be ok")
}

func TestNoEmpty12(t *testing.T) {
	r := lint(t, "try { foo(); } catch (ex) {} finally {}", &NoEmpty{})
	util.AssertEqual(t, 0, len(r.Abnormals), "should be ok")
	util.AssertEqual(t, 2, len(r.Diagnoses), "should be ok")
	util.AssertEqual(t, "disallow empty block statements", r.Diagnoses[0].Msg, "should be ok")
}
