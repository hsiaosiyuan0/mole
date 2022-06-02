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

	"github.com/hsiaosiyuan0/mole/script/macro"
	"github.com/hsiaosiyuan0/mole/util"
)

type MethodInfo struct {
	Name  string
	Nodes bool // whether the getter returns `[]parser.Node` or not
	Dec   *ast.FuncType
}

type StructInfo struct {
	Name      string
	TypName   string
	PushScope bool
	Methods   []*MethodInfo
}

func genVisitor(output io.Writer, s *StructInfo) error {
	name := fmt.Sprintf("visitor_%s", s.Name)
	// use tpl instead of manually construct the ast
	fnMap := template.FuncMap{
		"UnPrefix": unPrefix,
	}
	tpl, err := template.New(name).Funcs(fnMap).Parse(`
func Visit{{ .Name }}(node parser.Node, key string, ctx *VisitorCtx) {
  {{- if ge (len .Methods) 1 }}
    n := node.(*parser.{{ .Name }})

    {{- range $key, $value := .Methods }}
      {{ if eq $value.Name "PUSH_SCOPE" }}
        ctx.WalkCtx.PushScope()
        defer ctx.WalkCtx.PopScope()
      {{- end }}

      {{ if eq $key 0 }}
        CallVisitor(N_{{ $.TypName | UnPrefix }}_BEFORE, n, key, ctx)
        defer CallVisitor(N_{{ $.TypName | UnPrefix }}_AFTER, n, key, ctx)
      {{- end }}

      {{ if ne $value.Name "PUSH_SCOPE" }}
        {{ if .Nodes }}
          VisitNodes(n, n.{{ $value.Name }}(), "{{ $value.Name }}", ctx)
        {{- else }}
          VisitNode(n.{{ $value.Name }}(), "{{ $value.Name }}", ctx)
        {{- end }}
        if ctx.WalkCtx.Stopped() {
          return
        }
      {{- end }}
    {{- end }}
  {{- else}}
    CallListener(N_{{ $.TypName | UnPrefix }}_BEFORE, node, key, ctx)
    CallListener(N_{{ $.TypName | UnPrefix }}_AFTER, node, key, ctx)
  {{- end }}
}

{{ if ge (len .Methods) 1 }}
func Visit{{ .Name }}Before(node parser.Node, key string, ctx *VisitorCtx) {
  CallListener(N_{{ $.TypName | UnPrefix }}_BEFORE, node, key, ctx)
}

func Visit{{ .Name }}After(node parser.Node, key string, ctx *VisitorCtx) {
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
type Visitor = func(node parser.Node, key string, ctx *VisitorCtx)
type Visitors = [N_BEFORE_AFTER_DEF_END]Visitor

// replace the default visitor with the specified one
func SetVisitor(vs *Visitors, t parser.NodeType, impl Visitor) {
  vs[t] = impl
}

type Listener = func(node parser.Node, key string, ctx *VisitorCtx)
type Listeners = [N_BEFORE_AFTER_DEF_END][]Listener

func AddListener(ls *Listeners, t parser.NodeType, impl Listener) {
	ls[t] = append(ls[t], impl)
}

func NodeBeforeEvent(t parser.NodeType) parser.NodeType  {
  return N_BEFORE_AFTER_DEF_BEGIN + (parser.N_NODE_DEF_END - t) * 2 - 1
}

func NodeAfterEvent(t parser.NodeType) parser.NodeType  {
  return N_BEFORE_AFTER_DEF_BEGIN + (parser.N_NODE_DEF_END - t) * 2
}

func AddNodeBeforeListener(ls *Listeners, t parser.NodeType, impl Listener) {
	AddListener(ls, NodeBeforeEvent(t),impl)
}

func AddNodeAfterListener(ls *Listeners, t parser.NodeType, impl Listener) {
	AddListener(ls, NodeAfterEvent(t),impl)
}


func AddBeforeListener(ls *Listeners, impl Listener) {
	for t := range NodeTypes {
		AddNodeBeforeListener(ls, t, impl)
	}
}

func AddAfterListener(ls *Listeners, impl Listener) {
	for t := range NodeTypes {
		AddNodeAfterListener(ls, t, impl)
	}
}

func AddAtomListener(ls *Listeners, impl Listener) {
	for t := range AtomNodeTypes {
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

  {{- range $key, $value := .NodeTypStruct }}
    N_{{ $key | UnPrefix | ToUpper }}_BEFORE = N_BEFORE_AFTER_DEF_BEGIN + (parser.N_NODE_DEF_END - N_{{ $key | UnPrefix | ToUpper }}) * 2 - 1
    N_{{ $key | UnPrefix | ToUpper }}_AFTER =  N_BEFORE_AFTER_DEF_BEGIN + (parser.N_NODE_DEF_END - N_{{ $key | UnPrefix | ToUpper }}) * 2
  {{- end }}

  N_BEFORE_AFTER_DEF_END = N_BEFORE_AFTER_DEF_BEGIN + parser.N_NODE_DEF_END * 2
)

var AtomNodeTypes = map[parser.NodeType]bool{
  {{- range $key, $value := .NodeTypStruct }}
    {{- if eq (len (index $.StructColl $value).Methods) 0  }}
      N_{{ $key | UnPrefix | ToUpper }}: true,
    {{- end }}
  {{- end }}
}

var NodeTypes = map[parser.NodeType]bool{
  {{- range $key, $value := .NodeTypStruct }}
    N_{{ $key | UnPrefix | ToUpper }}: true,
  {{- end }}
}

var NodeBeforeEvents = map[parser.NodeType]bool {
  {{- range $key, $value := .NodeTypStruct }}
    N_BEFORE_AFTER_DEF_BEGIN + (parser.N_NODE_DEF_END - N_{{ $key | UnPrefix | ToUpper }}) * 2 - 1: true,
  {{- end }}
}

var NodeAfterEvents = map[parser.NodeType]bool {
  {{- range $key, $value := .NodeTypStruct }}
    N_BEFORE_AFTER_DEF_BEGIN + (parser.N_NODE_DEF_END - N_{{ $key | UnPrefix | ToUpper }}) * 2: true,
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
		structColl[name] = &StructInfo{name, "", false, []*MethodInfo{}}
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
				if n == "PUSH_SCOPE" {
					s.PushScope = true
				}
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
	util.Shell("gofmt", "-w", distFile)
}
