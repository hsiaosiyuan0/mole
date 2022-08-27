package resolver

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strings"
)

type SpecifierKind uint8

const (
	SPK_UNKNOWN SpecifierKind = iota
	SPK_FILE
	SPK_IMPT
	SPK_BARE
	SPK_DATA
	SPK_NODE
)

type SpecifierDataKind uint8

const (
	SPDK_UNKNOWN SpecifierDataKind = iota
	SPDK_JS
	SPDK_JSON
	SPDK_WASM
)

// https://nodejs.org/api/esm.html#terminology
type Specifier struct {
	kind SpecifierKind
	raw  string
	s    string
	ss   string
	u    *url.URL
	d    string
	dk   SpecifierDataKind
}

var fileSchema = "file:"
var dataSchema = "data:"
var nodeSchema = "node:"

var dataJsMime = "text/javascript"
var dataJsonMime = "application/json"
var dataWasmMime = "application/wasm"

func NewSpecifier(sp, cw string) (*Specifier, error) {
	if len(sp) == 0 {
		return nil, newInvalidSpecifierErr(sp, "deformed specifier")
	}

	s := &Specifier{SPK_UNKNOWN, sp, "", "", nil, "", SPDK_UNKNOWN}
	c := sp[0]
	if c == '.' || c == '/' { // implicit `file:`
		if c == '.' {
			sp = filepath.Join(cw, sp)
		}
		s.kind = SPK_FILE
	} else if strings.HasPrefix(sp, fileSchema) { // explicit `file:`
		sp = sp[len(fileSchema):]
		if strings.HasPrefix(sp, "//") {
			sp = sp[2:]
		}
		s.kind = SPK_FILE
	} else if strings.HasPrefix(sp, nodeSchema) { // `node:`
		s.kind = SPK_NODE
		sp = sp[len(nodeSchema):]
		s.s = sp
	} else if strings.HasPrefix(sp, dataSchema) { // `data:`
		s.kind = SPK_DATA
		sp = sp[len(dataSchema):]
		s.d, s.dk = parseMine(sp)
		if s.dk == SPDK_UNKNOWN {
			return nil, newInvalidSpecifierErr(sp, "deformed specifier")
		}
	} else if c == '#' {
		s.kind = SPK_IMPT
		s.s = sp
	} else if strings.IndexByte(sp, ':') == -1 {
		s.kind = SPK_BARE
		s.s, s.ss = subpathOf(sp)
	} else {
		return nil, newInvalidSpecifierErr(sp, "unsupported specifier")
	}

	if s.kind == SPK_FILE {
		u, err := url.Parse(sp)
		if err != nil {
			return nil, newInvalidSpecifierErr(sp, err.Error())
		}
		s.u = u
		s.s = u.Path
	}

	return s, nil
}

func subpathOf(s string) (string, string) {
	i := strings.IndexRune(s, '/')
	if i == -1 {
		return s, "."
	}

	name := s[0:i]
	subpath := s[i+1:]
	if name[0] == '@' {
		j := strings.IndexRune(subpath, '/')
		if j == -1 {
			return s, "."
		}
		k := i + j + 1
		name = s[0:k]
		subpath = s[k:]
	} else {
		subpath = "/" + subpath
	}

	if name[0] == '.' || subpath[len(subpath)-1] == '/' {
		return "", ""
	}

	return name, "." + subpath
}

func parseMine(data string) (string, SpecifierDataKind) {
	kind := SPDK_UNKNOWN
	raw := ""
	if strings.HasPrefix(data, dataJsMime) {
		raw = data[len(dataJsMime):]
		kind = SPDK_JS
	} else if strings.HasPrefix(data, dataJsonMime) {
		raw = data[len(dataJsonMime):]
		kind = SPDK_JSON
	} else if strings.HasPrefix(data, dataWasmMime) {
		raw = data[len(dataWasmMime):]
		kind = SPDK_WASM
	}
	if len(raw) == 0 {
		return "", SPDK_UNKNOWN
	}
	if raw[0] != ',' {
		return "", SPDK_UNKNOWN
	}
	return raw[1:], kind
}

type InvalidSpecifierErr struct {
	Specifier string
	Msg       string
}

func newInvalidSpecifierErr(specifier, msg string) *InvalidSpecifierErr {
	return &InvalidSpecifierErr{specifier, msg}
}

func (e *InvalidSpecifierErr) Error() string {
	return fmt.Sprintf(`Module specifier is an invalid URL, reason: %s specifier: "%s"  `, e.Msg, e.Specifier)
}
