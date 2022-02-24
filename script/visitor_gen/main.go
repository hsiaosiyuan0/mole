package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/hsiaosiyuan0/mole/fuzz"
	"github.com/hsiaosiyuan0/mole/script/macro"
)

type MethodInfo struct {
	Name  string
	Nodes bool // whether the getter returns `[]parser.Node` or not
	Dec   *ast.FuncType
}

type StructInfo struct {
	Name    string
	TypName string
	Methods []*MethodInfo
}

func genVisitor(output io.Writer, s *StructInfo) error {
	name := fmt.Sprintf("visitor_%s", s.Name)
	// use tpl instead of manually construct the ast
	fnMap := template.FuncMap{
		"UnPrefix": unPrefix,
	}
	tpl, err := template.New(name).Funcs(fnMap).Parse(`
func Visit{{ .Name }}(node parser.Node, key string, ctx *WalkCtx) {
  {{- if ge (len .Methods) 1 }}
    n := node.(*parser.{{ .Name }})
    ctx.PushVisitorCtx(n, key)
    defer ctx.PopVisitorCtx()

    {{- range $key, $value := .Methods }}
      {{ if eq $value.Name "PUSH_SCOPE" }}
        ctx.PushScope()
        defer ctx.PopScope()

        CallVisitor(N_{{ $.TypName | UnPrefix }}_BEFORE, n, key, ctx)
        defer CallVisitor(N_{{ $.TypName | UnPrefix }}_AFTER, n, key, ctx)
      {{- else }}
        {{ if eq $key 0 }}
          CallVisitor(N_{{ $.TypName | UnPrefix }}_BEFORE, n, key, ctx)
          defer CallVisitor(N_{{ $.TypName | UnPrefix }}_AFTER, n, key, ctx)
        {{- end }}
        {{ if .Nodes }}
          VisitNodes(n, n.{{ $value.Name }}(), "{{ $value.Name }}", ctx)
        {{- else }}
          VisitNode(n.{{ $value.Name }}(), "{{ $value.Name }}", ctx)
        {{- end }}
        if ctx.stop {
          return
        }
      {{- end }}
    {{- end }}
  {{- else}}
    CallListener(parser.N_{{ $.TypName | UnPrefix }}, node, key, ctx)
  {{- end }}
}

{{ if ge (len .Methods) 1 }}
func Visit{{ .Name }}Before(node parser.Node, key string, ctx *WalkCtx) {
  CallListener(N_{{ $.TypName | UnPrefix }}_BEFORE, node, key, ctx)
}

func Visit{{ .Name }}After(node parser.Node, key string, ctx *WalkCtx) {
  CallListener(N_{{ $.TypName | UnPrefix }}_AFTER, node, key, ctx)
}
{{- end}}
  `)
	if err != nil {
		return err
	}
	return tpl.Execute(output, s)
}

func genVisitors(output io.Writer, nodeTypStruct map[string]string, structColl map[string]*StructInfo) {
	output.Write([]byte(`
type Visitor = func(node parser.Node, key string, ctx *WalkCtx)
type Visitors = [N_BEFORE_AFTER_DEF_END]Visitor

// replace the default visitor with the specified one
func SetVisitor(vs *Visitors, t parser.NodeType, impl Visitor) {
  vs[t] = impl
}

type Listener = func(node parser.Node, key string, ctx *WalkCtx)
type Listeners = [N_BEFORE_AFTER_DEF_END][]Listener

func AddListener(ls *Listeners, t parser.NodeType, impl Listener) {
	ls[t] = append(ls[t], impl)
}

func AddBeforeListener(ls *Listeners, impl Listener) {
	for i := N_BEFORE_DEF_BEGIN + 1; i < N_BEFORE_DEF_END; i++ {
		ls[i] = append(ls[i], impl)
	}
}

func AddAfterListener(ls *Listeners, impl Listener) {
	for i := N_AFTER_DEF_BEGIN + 1; i < N_AFTER_DEF_END; i++ {
		ls[i] = append(ls[i], impl)
	}
}

func AddAtomListener(ls *Listeners, impl Listener) {
  for _, t:= range atomNodeTypes {
    ls[t] = append(ls[t], impl)
  }
}
  `))

	processedVisitors := map[string]bool{}
	for _, structName := range nodeTypStruct {
		if _, ok := processedVisitors[structName]; ok {
			continue
		}
		err := genVisitor(output, structColl[structName])
		if err != nil {
			log.Fatal(err)
		}
		processedVisitors[structName] = true
	}
}

func unPrefix(s string) string {
	idx := strings.IndexByte(s, '_')
	return s[idx+1:]
}

type TplParamsGenVisitorKinds struct {
	NodeTypStruct map[string]string
	StructColl    map[string]*StructInfo
}

func genVisitorKinds(output io.Writer, nodeTypStruct map[string]string, structColl map[string]*StructInfo) error {
	fnMap := template.FuncMap{
		"ToUpper":  strings.ToUpper,
		"UnPrefix": unPrefix,
	}
	tpl, err := template.New("visitor types").Funcs(fnMap).Parse(`
const (
  {{- range $key, $value := .NodeTypStruct }}
    N_{{- $key | UnPrefix | ToUpper }} = parser.{{ $key }}
  {{- end }}
)

const (
  N_BEFORE_AFTER_DEF_BEGIN = parser.NodeType(parser.N_NODE_DEF_END + iota)
  N_BEFORE_DEF_BEGIN
  {{- range $key, $value := .NodeTypStruct }}
    {{- if ge (len (index $.StructColl $value).Methods) 1  }}
      N_{{ $key | UnPrefix | ToUpper }}_BEFORE
    {{- end }}
  {{- end }}
  N_BEFORE_DEF_END

  N_AFTER_DEF_BEGIN
  {{- range $key, $value := .NodeTypStruct }}
    {{- if ge (len (index $.StructColl $value).Methods) 1  }}
      N_{{ $key | UnPrefix | ToUpper }}_AFTER
    {{- end }}
  {{- end }}
  N_AFTER_DEF_END

  N_BEFORE_AFTER_DEF_END
)

var atomNodeTypes = []parser.NodeType{
  {{- range $key, $value := .NodeTypStruct }}
    {{- if eq (len (index $.StructColl $value).Methods) 0  }}
      N_{{ $key | UnPrefix | ToUpper }},
    {{- end }}
  {{- end }}
}
  `)
	if err != nil {
		return err
	}
	return tpl.Execute(output, &TplParamsGenVisitorKinds{nodeTypStruct, structColl})
}

func genDefaultVisitors(output io.Writer, nodeTypStruct map[string]string, structColl map[string]*StructInfo) error {
	fnMap := template.FuncMap{
		"ToUpper":  strings.ToUpper,
		"UnPrefix": unPrefix,
	}
	tpl, err := template.New("visitor types").Funcs(fnMap).Parse(`
var DefaultVisitors Visitors = [N_BEFORE_AFTER_DEF_END]Visitor{}
var DefaultListeners Listeners = [N_BEFORE_AFTER_DEF_END][]Listener{}

func init() {
  {{- range $key, $value := .NodeTypStruct }}
    DefaultVisitors[N_{{ $key | UnPrefix | ToUpper }}] = Visit{{ $value }}
    {{- if ge (len (index $.StructColl $value).Methods) 1  }}
      DefaultVisitors[N_{{ $key | UnPrefix | ToUpper }}_BEFORE] = Visit{{ $value }}Before
      DefaultVisitors[N_{{ $key | UnPrefix | ToUpper }}_AFTER] = Visit{{ $value }}After
    {{- end }}
  {{- end }}

  {{ range $key, $value := .NodeTypStruct }}
    DefaultListeners[N_{{ $key | UnPrefix | ToUpper }}] = []Listener{}
  {{- end }}
}
  `)
	if err != nil {
		return err
	}
	return tpl.Execute(output, &TplParamsGenVisitorKinds{nodeTypStruct, structColl})
}

func IfNodesReurned(f *ast.FuncType, name string) bool {
	if f.Results == nil || len(f.Results.List) != 1 {
		log.Fatalf("%s should return only one value", name)
	}
	r := f.Results.List[0]
	a, ok := r.Type.(*ast.ArrayType)
	if !ok {
		return false
	}
	i, ok := a.Elt.(*ast.Ident)
	if !ok {
		return false
	}
	return i.Name == "Node"
}

func main() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	distFile := path.Join(wd, "visitor.go")
	_, err = os.Stat(distFile)
	if err != nil {
		return
	}

	var defDir string
	flag.StringVar(&defDir, "d", "", "the AST definition directory, relative with current file")
	flag.Parse()

	ctxs, procCtx, err := macro.MacroCtxsOfWorkingDir(wd, defDir)
	if err != nil {
		log.Fatal(err)
	}

	nodeTypStruct := map[string]string{}
	structColl := map[string]*StructInfo{}

	// upsert struct by its name
	getStruct := func(name string) *StructInfo {
		s, ok := structColl[name]
		if ok {
			return s
		}
		structColl[name] = &StructInfo{name, "", []*MethodInfo{}}
		return structColl[name]
	}

	// fullfill `nodeTypStruct` and `structColl`
	for _, ctx := range ctxs {
		if ctx.Name != "visitor" {
			continue
		}
		if v, ok := ctx.Node.(*ast.ValueSpec); ok {
			nodeTyp := v.Names[0].Name
			structName := ctx.Args[0].(string)
			nodeTypStruct[nodeTyp] = structName

			s := getStruct(structName)
			s.TypName = nodeTyp
		} else if name, _, ok := macro.IsStructDec(ctx.Node); ok {
			s := getStruct(name)
			for _, n := range ctx.Args {
				s.Methods = append(s.Methods, &MethodInfo{Name: n.(string), Dec: nil})
			}
		}
	}

	// walk the pkgs to find out the type of `MethodInfo.dec`
	macro.WalkPkgs(procCtx.Pkgs, func(f *ast.File, s string, pc *macro.ProcCtx) error {
		for _, dec := range f.Decls {
			if v, ok := dec.(*ast.FuncDecl); ok {
				if v.Recv == nil {
					continue
				}
				name := v.Name.Name
				fnTyp := v.Type
				recvName := macro.RecvName(v)
				if recvName != "" {
					if structColl[recvName] == nil {
						continue
					}
					sm := structColl[recvName].Methods
					for _, m := range sm {
						if m.Name == name {
							m.Dec = fnTyp
							m.Nodes = IfNodesReurned(fnTyp, name)
							break
						}
					}
				}
			}
		}
		return nil
	}, procCtx)

	var buf bytes.Buffer
	buf.WriteString(`// Code generated by script/visitor_gen. DO NOT EDIT.

//go:generate go run github.com/hsiaosiyuan0/mole/script/visitor_gen -d=../parser

package walk

import "github.com/hsiaosiyuan0/mole/ecma/parser"

  `)

	// genereate visitor kinds
	err = genVisitorKinds(&buf, nodeTypStruct, structColl)
	if err != nil {
		log.Fatal(err)
	}

	// generate visitors
	genVisitors(&buf, nodeTypStruct, structColl)

	// generate defalt visitors
	err = genDefaultVisitors(&buf, nodeTypStruct, structColl)
	if err != nil {
		log.Fatal(err)
	}

	ioutil.WriteFile(distFile, buf.Bytes(), 0644)
	fuzz.Shell("gofmt", "-w", distFile)
}
