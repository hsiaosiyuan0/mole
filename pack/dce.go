package pack

import (
	"encoding/json"
	"fmt"

	"container/list"

	"github.com/hsiaosiyuan0/mole/ecma/astutil"
	"github.com/hsiaosiyuan0/mole/ecma/parser"
	"github.com/hsiaosiyuan0/mole/ecma/walk"
	"github.com/hsiaosiyuan0/mole/span"
)

type TopmostDec struct {
	s          *span.Source
	Node       parser.Node
	Alive      bool
	SideEffect bool
	Owners     map[*TopmostDec]*TopmostDec
	Owned      map[*TopmostDec]*TopmostDec
}

func idOfDec(n parser.Node) uint64 {
	rng := n.Range()
	return uint64(rng.Lo)<<32 | uint64(rng.Hi)
}

func (d *TopmostDec) MarshalJSON() ([]byte, error) {
	owners := []uint64{}
	for _, d := range d.Owners {
		owners = append(owners, idOfDec(d.Node))
	}

	owned := []uint64{}
	for _, d := range d.Owned {
		owned = append(owned, idOfDec(d.Node))
	}

	src := ""
	switch d.Node.Type() {
	case parser.N_STMT_EXPORT:
		n := d.Node.(*parser.ExportDec)
		if n.Src() != nil {
			src = n.Src().(*parser.StrLit).Val()
		} else {
			names, _ := astutil.NamesInDecNode(n)
			if len(names) > 0 {
				src = names[0]
			}
		}
	case parser.N_STMT_IMPORT:
		n := d.Node.(*parser.ImportDec)
		if n.Src() != nil {
			src = n.Src().(*parser.StrLit).Val()
		}
	default:
		rng := d.Node.Range()
		rng.Hi = rng.Lo + 15
		if int(rng.Hi) > d.s.Len() {
			rng.Hi = uint32(d.s.Len())
		}
		src = d.s.RngText(rng)
	}

	return json.Marshal(struct {
		Id         uint64   `json:"id"`
		NodeType   string   `json:"nodeType"`
		Src        string   `json:"src"`
		Range      []uint32 `json:"range"`
		Alive      bool     `json:"alive"`
		SideEffect bool     `json:"sideEffect"`
		Owners     []uint64 `json:"owners"`
		Owned      []uint64 `json:"owned"`
	}{
		Id:         idOfDec(d.Node),
		NodeType:   d.Node.Type().String(),
		Src:        src,
		Range:      []uint32{d.Node.Range().Lo, d.Node.Range().Hi},
		Alive:      d.Alive,
		SideEffect: d.SideEffect,
		Owners:     owners,
		Owned:      owned,
	})
}

func isPure(node parser.Node, p *parser.Parser) bool {
	symtab := p.Symtab()
	nt := node.Type()
	if nt == parser.N_STMT_FN || nt == parser.N_EXPR_FN {
		return true
	}

	switch nt {
	case parser.N_STMT_FN, parser.N_EXPR_FN, parser.N_STMT_EXPORT:
		return true
	case parser.N_STMT_VAR_DEC:
		vdc := node.(*parser.VarDecStmt)
		if vdc.Kind() == "var" || len(vdc.DecList()) != 1 {
			return false
		}
		vd := vdc.DecList()[0].(*parser.VarDec)
		if vd.Init() == nil {
			return true
		}
		init := vd.Init()
		it := init.Type()
		if it == parser.N_STMT_FN || it == parser.N_EXPR_FN {
			return true
		}
		if it == parser.N_NAME {
			name := init.(*parser.Ident).Val()
			binding := symtab.Scopes[0].BindingOf(name)
			if binding != nil && binding.Typ == parser.RDT_FN {
				return true
			}
		}
		if astutil.IsPlainObj(init) {
			return true
		}
		return hasPureAnno(init, p)
	case parser.N_STMT_IMPORT:
		n := node.(*parser.ImportDec)
		for _, s := range n.Specs() {
			sp := s.(*parser.ImportSpec)
			if sp.NameSpace() {
				return false
			}
		}
		return true
	}

	return hasPureAnno(node, p)
}

func hasPureAnno(n parser.Node, p *parser.Parser) bool {
	cmts := p.PrevCmts(n)
	if len(cmts) != 1 {
		return false
	}
	return p.RngText(cmts[0]) == "/*#__PURE__*/"
}

func isRefDefInNode(ref *parser.Ref, n parser.Node) bool {
	nt := n.Type()
	if nt == parser.N_STMT_FN || nt == parser.N_EXPR_FN {
		f := n.(*parser.FnDec)
		return f.Id() == ref.Id
	}
	if nt == parser.N_STMT_VAR_DEC {
		vds := n.(*parser.VarDecStmt).DecList()
		for _, vd := range vds {
			if vd.(*parser.VarDec).Id() == ref.Id {
				return true
			}
		}
	}
	return false
}

func resolveTopmostDecs(p *parser.Parser) (tds map[parser.Node]*TopmostDec, exports map[string]*TopmostDec, exportAll []*TopmostDec) {
	ast := p.Ast()
	stmts := ast.(*parser.Prog).Body()
	tds = map[parser.Node]*TopmostDec{}

	symtab := p.Symtab()
	ref2td := map[*parser.Ref]*TopmostDec{}
	for _, n := range stmts {
		td := &TopmostDec{p.Source(), n, false, !isPure(n, p), map[*TopmostDec]*TopmostDec{}, map[*TopmostDec]*TopmostDec{}}
		name := astutil.GetNodeName(n)
		ref := symtab.Scopes[0].BindingOf(name)
		if ref != nil {
			ref2td[ref] = td
		}
		tds[n] = td
	}

	ctx := walk.NewWalkCtx(ast, p.Symtab())

	// build the relations of the topmostDecs
	handleName := func(node parser.Node, key string, ctx *walk.VisitorCtx) {
		if key == "Id" {
			return
		}

		scope := symtab.Scopes[ctx.ScopeId()]
		ref := scope.BindingOf(astutil.GetNodeName(node))
		if ref == nil || ref.Scope.Id != 0 {
			return
		}

		// ref here is a topmost ref and `td` will be its target topmostDec
		if td := tds[ref.Dec]; td != nil {

			// find the topmostDec which encapsulates this reference point then
			// set it ast the the owner of `td`
			ctx = ctx.Parent
			for {
				if ctx == nil {
					break
				}
				n := ctx.Node
				if n == nil {
					break
				}

				if otd := tds[n]; otd != nil {
					if otd != td {
						td.Owners[otd] = otd
						otd.Owned[td] = td
					}
					break
				}
				ctx = ctx.Parent
			}
		}
	}

	walk.AddNodeAfterListener(&ctx.Listeners, parser.N_NAME, &walk.Listener{
		Id:     "N_NAME",
		Handle: handleName,
	})
	walk.AddNodeAfterListener(&ctx.Listeners, parser.N_JSX_ID, &walk.Listener{
		Id:     "N_JSX_ID",
		Handle: handleName,
	})

	// build the map for routing the export name to its topmostDec
	exports = map[string]*TopmostDec{}
	exportAll = []*TopmostDec{}
	walk.AddNodeAfterListener(&ctx.Listeners, parser.N_STMT_EXPORT, &walk.Listener{
		Id: "parseDep",
		Handle: func(node parser.Node, key string, ctx *walk.VisitorCtx) {
			n := node.(*parser.ExportDec)

			td := tds[n]
			if n.Default() {
				exports["default"] = td
			} else if n.All() {
				exportAll = append(exportAll, td)
			} else if dec := n.Dec(); dec != nil {
				name := astutil.GetNodeName(dec)
				if name != "" {
					exports[name] = td
				}
			} else {
				scope := symtab.Scopes[ctx.ScopeId()]
				ext := n.Src() != nil
				for _, s := range n.Specs() {
					sp := s.(*parser.ExportSpec)
					name := sp.Id().(*parser.Ident).Val()
					if ext {
						exports[name] = td
					} else {
						local := sp.Local().(*parser.Ident).Val()
						ref := scope.BindingOf(local)
						exports[name] = tds[ref.Dec]
					}
				}
			}
		},
	})

	walk.VisitNode(ast, "", ctx.VisitorCtx())

	return
}

func namesOfExport(n parser.ExportDec) (ret []string, all bool) {
	ret = []string{}
	if n.Default() {
		ret = append(ret, "default")
	} else if n.All() {
		all = true
	} else if n.Dec() != nil {

	} else {
		for _, spec := range n.Specs() {
			sp := spec.(*parser.ExportSpec)
			if sp.NameSpace() {
				all = true
			} else {
				ret = append(ret, sp.Id().(*parser.Ident).Val())
			}
		}
	}
	return
}

type Import struct {
	name     string
	from     int64
	delegate bool
}

func (t *Import) key() string {
	return fmt.Sprintf("%s:%d", t.name, t.from)
}

func importsOfNode(node parser.Node, m *Module) *list.List {
	imports := list.New()
	switch node.Type() {
	case parser.N_STMT_IMPORT:
		n := node.(*parser.ImportDec)
		from := m.extsMap[n.Src().(*parser.StrLit).Val()]
		if from == 0 {
			break
		}

		if len(n.Specs()) > 0 {
			for _, s := range n.Specs() {
				spec := s.(*parser.ImportSpec)
				name := ""
				if spec.Default() {
					name = "default"
				} else if spec.NameSpace() {
					name = "#all"
				} else {
					name = spec.Id().(*parser.Ident).Val()
				}
				imports.PushBack(&Import{name, from, false})
			}
		} else {
			imports.PushBack(&Import{"#all", from, false})
		}
	case parser.N_STMT_EXPORT:
		n := node.(*parser.ExportDec)
		if n.Src() != nil {
			from := m.extsMap[n.Src().(*parser.StrLit).Val()]
			if from == 0 {
				break
			}
			if n.Default() {
				imports.PushBack(&Import{"default", from, false})
			} else if n.All() {
				imports.PushBack(&Import{"#all", from, false})
			} else {
				for _, spec := range n.Specs() {
					sp := spec.(*parser.ExportSpec)
					if sp.NameSpace() {
						imports.PushBack(&Import{"#all", from, false})
					} else {
						imports.PushBack(&Import{sp.Id().(*parser.Ident).Val(), from, false})
					}
				}
			}
		}
	default:
		walkRequireCall(node, func(str string) {
			from := m.extsMap[str]
			imports.PushBack(&Import{"#all", from, false})
		})
	}
	return imports
}

type RequireCb = func(str string)

func walkRequireCall(node parser.Node, cb RequireCb) {
	wc := walk.NewWalkCtx(node, nil)
	walk.AddNodeAfterListener(&wc.Listeners, parser.N_EXPR_CALL, &walk.Listener{
		Id: "N_EXPR_CALL",
		Handle: func(node parser.Node, key string, ctx *walk.VisitorCtx) {
			n := node.(*parser.CallExpr)
			s := ctx.WalkCtx.Scope()
			callee := n.Callee()
			args := n.Args()

			isRequire := astutil.GetName(callee) == "require" && s.BindingOf("require") == nil &&
				len(args) == 1 && args[0].Type() == parser.N_LIT_STR

			if isRequire {
				cb(args[0].(*parser.StrLit).Val())
			}
		},
	})
}

func importsOfModule(m *Module) *list.List {
	imports := list.New()
	for _, d := range m.tds {
		imports.PushBackList(importsOfNode(d.Node, m))
	}
	return imports
}

type markTopmostDecCb = func(*TopmostDec)

func markTopmostDec(td *TopmostDec, cb markTopmostDecCb) {
	tds := []*TopmostDec{td}
	unique := map[*TopmostDec]bool{}
	for {
		if len(tds) == 0 {
			break
		}

		first, remain := tds[0], tds[1:]
		cb(first)

		for _, owd := range first.Owned {
			if !unique[owd] {
				unique[owd] = true
				remain = append(remain, owd)
			}
		}
		tds = remain
	}
}

func importSrcOf(node parser.Node) ([]string, bool) {
	switch node.Type() {
	case parser.N_STMT_IMPORT:
		n := node.(*parser.ImportDec)
		if n.Src() != nil {
			return []string{n.Src().(*parser.StrLit).Val()}, false
		}
	case parser.N_STMT_EXPORT:
		n := node.(*parser.ExportDec)
		if n.Src() != nil {
			return []string{n.Src().(*parser.StrLit).Val()}, false
		}
	default:
		ret := []string{}
		walkRequireCall(node, func(str string) {
			ret = append(ret, str)
		})
		return ret, true
	}
	return []string{}, false
}

func importsOfTopmostDec(td *TopmostDec, m *Module) *list.List {
	imports := list.New()
	owned := []*TopmostDec{td}

	unique := map[parser.Node]bool{}
	for {
		if len(owned) == 0 {
			break
		}

		first, remain := owned[0], owned[1:]
		if unique[first.Node] {
			owned = remain
			continue
		}

		unique[first.Node] = true
		imports.PushBackList(importsOfNode(first.Node, m))
		for _, td := range first.Owners {
			remain = append(remain, td)
		}
		owned = remain
	}
	return imports
}

func (s *DepScanner) DCE() {
	imports := list.New()

	imported := map[int64]bool{}
	for _, mid := range s.entries {
		m := s.allModules[mid]

		// every export in entry has side-effect
		for _, exp := range m.exports {
			exp.SideEffect = true
			markTopmostDec(exp, func(td *TopmostDec) {
				td.Alive = true
			})
		}
		for _, exp := range m.exportAll {
			exp.SideEffect = true
			markTopmostDec(exp, func(td *TopmostDec) {
				td.Alive = true
			})
		}

		imports.PushBackList(importsOfModule(m))
	}

	unique := map[string]bool{}
	pushImport := func(ipt *Import) {
		if !unique[ipt.key()] {
			unique[ipt.key()] = true
			imports.PushBack(ipt)
		}
	}

	for {
		if imports.Len() == 0 {
			break
		}

		ipt := imports.Remove(imports.Front()).(*Import)
		m := s.allModules[ipt.from]

		imported[ipt.from] = true

		// apply affects
		all := ipt.name == "#all"

		if !all {
			td := m.exports[ipt.name]
			if td != nil {
				td.Alive = true
				markTopmostDec(td, func(td *TopmostDec) {
					td.Alive = true
					imports := importsOfTopmostDec(td, m)
					next := imports.Front()
					for {
						if next == nil {
							break
						}
						ipt := next.Value.(*Import)
						pushImport(ipt)
						next = next.Next()
					}
				})
			} else {
				// delegate the import to the modules which are imported by `*`
				for _, exp := range m.exportAll {
					var from int64
					switch exp.Node.Type() {
					case parser.N_STMT_IMPORT:
						n := exp.Node.(*parser.ImportDec)
						from = m.extsMap[n.Src().(*parser.StrLit).Val()]
					case parser.N_STMT_EXPORT:
						n := exp.Node.(*parser.ExportDec)
						if n.Src() != nil {
							from = m.extsMap[n.Src().(*parser.StrLit).Val()]
						}
					}
					if from != 0 {
						imports.PushBack(&Import{ipt.name, from, true})
					}
				}
			}
		}

		for _, td := range m.tds {
			td.Alive = td.Alive || td.SideEffect || all
			if td.Alive {
				markTopmostDec(td, func(td *TopmostDec) {
					td.Alive = true
					srcList, sideEffect := importSrcOf(td.Node)
					for _, src := range srcList {
						if src != "" {
							// does not eliminate the unused named-import yet
							names, all := astutil.NamesInDecNode(td.Node)
							from := m.extsMap[src]
							if all || sideEffect {
								pushImport(&Import{"#all", from, false})
							} else if from != 0 {
								for _, name := range names {
									pushImport(&Import{name, from, false})
								}
							}
						}
					}
				})
			}
		}
	}

	// recalculate the size of umbrella modules
	for mid := range imported {
		m := s.allModules[mid]

		if !m.IsUmbrella() {
			um := s.allModules[m.umbrella]
			if um != nil {
				um.dceSize += m.calcDceSize()
			}
		}
	}
}
