package parser

import (
	"fmt"

	"github.com/hsiaosiyuan0/mole/span"
)

func (p *Parser) jsxMember(obj Node, path string) (Node, string, error) {
	for {
		ahead := p.lexer.Peek()
		av := ahead.value
		if av == T_DOT {
			p.lexer.Next()
			prop, err := p.ident(nil, false)
			if err != nil {
				return nil, "", err
			}
			path = path + "." + prop.val
			obj = &JsxMember{N_JSX_MEMBER, p.finRng(obj.Range()), obj, prop, p.newTypInfo(N_JSX_MEMBER)}
		} else {
			break
		}
	}
	return obj, path, nil
}

func (p *Parser) jsxNsExpr(ns Node, path string) (Node, string, error) {
	if _, err := p.nextMustTok(T_COLON); err != nil {
		return nil, "", err
	}
	name, pth, err := p.jsxId()
	if err != nil {
		return nil, "", err
	}
	return &JsxNsName{N_JSX_NS, p.finRng(ns.Range()), ns, name}, path + ":" + pth, nil
}

func (p *Parser) jsxId() (Node, string, error) {
	tok := p.lexer.Next()
	tv := tok.value
	if tv != T_NAME {
		return nil, "", p.errorTok(tok)
	}

	t := p.TokText(tok)
	return &JsxIdent{N_JSX_ID, p.finRng(tok.rng), t, span.Range{}, p.newTypInfo(N_JSX_ID)}, t, nil
}

func (p *Parser) jsxName() (Node, string, error) {
	id, jsxName, err := p.jsxId()
	if err != nil {
		return nil, "", err
	}

	ahead := p.lexer.Peek()
	av := ahead.value

	var name Node
	var pth string
	if av == T_DOT {
		name, pth, err = p.jsxMember(id, jsxName)
		if err != nil {
			return nil, "", err
		}
	} else if av == T_COLON && p.feat&FEAT_JSX_NS != 0 {
		name, pth, err = p.jsxNsExpr(id, jsxName)
		if err != nil {
			return nil, "", err
		}
	} else {
		name = id
		pth = jsxName
	}
	typArgs, err := p.tsTryTypArgs(span.Range{}, true)
	if err != nil {
		return nil, "", err
	}
	if wt, ok := name.(NodeWithTypInfo); ok {
		ti := wt.TypInfo()
		if ti != nil {
			ti.SetTypArgs(typArgs)
		}
	} else if typArgs != nil {
		return nil, "", p.errorAtLoc(typArgs.Range(), ERR_UNEXPECTED_TOKEN)
	}

	return name, pth, nil
}

func (p *Parser) jsxAttr() (Node, error) {
	ahead := p.lexer.Peek()
	if ahead.value == T_BRACE_L {
		rng := p.lexer.Next().rng
		val, err := p.jsxExpr(rng)
		if err != nil {
			return nil, err
		}
		if val.Type() != N_SPREAD {
			return nil, p.errorAtLoc(rng, ERR_UNEXPECTED_TOKEN)
		}
		return &JsxSpreadAttr{N_JSX_ATTR_SPREAD, p.finRng(rng), val}, nil
	}

	id, name, err := p.jsxName()
	if err != nil {
		return nil, err
	}
	p.lexer.PushMode(LM_JSX_ATTR, true)
	attr := &JsxAttr{N_JSX_ATTR, p.finRng(id.Range()), id, name, nil}
	if p.lexer.Peek().value != T_ASSIGN {
		p.lexer.PopMode()
		return attr, nil
	}
	p.lexer.Next()

	ahead = p.lexer.Peek()
	av := ahead.value
	var val Node
	if av == T_BRACE_L {
		rng := p.lexer.Next().rng
		val, err = p.jsxExpr(rng)
		if err != nil {
			return nil, err
		}
	} else if av == T_STRING {
		tok := p.lexer.Next()
		rng := tok.rng
		val = &StrLit{N_LIT_STR, p.finRng(rng), p.TokText(tok), tok.HasLegacyOctalEscapeSeq(), span.Range{}, nil}
	} else if av == T_LT {
		val, err = p.jsx(true, false)
		if err != nil {
			return nil, err
		}
	}
	attr.rng = p.finRng(attr.rng)
	attr.val = val
	p.lexer.PopMode()
	return attr, nil
}

func (p *Parser) jsxAttrs() ([]Node, error) {
	attrs := make([]Node, 0)
	for {
		ahead := p.lexer.Peek()
		av := ahead.value
		if av == T_GT || av == T_DIV || av == T_EOF {
			break
		}
		attr, err := p.jsxAttr()
		if err != nil {
			return nil, err
		}
		attrs = append(attrs, attr)
	}
	return attrs, nil
}

func (p *Parser) jsxOpen(tok *Token) (Node, error) {
	rng := tok.rng
	// fragment
	if p.lexer.Peek().value == T_GT {
		p.lexer.Next()
		return &JsxOpen{N_JSX_OPEN, p.finRng(rng), nil, "", nil, false}, nil
	}

	id, name, err := p.jsxName()
	if err != nil {
		return nil, err
	}
	attrs, err := p.jsxAttrs()
	if err != nil {
		return nil, err
	}
	closed := p.lexer.Peek().value == T_DIV
	if closed {
		p.lexer.Next()
	}
	if _, err := p.nextMustTok(T_GT); err != nil {
		return nil, err
	}
	return &JsxOpen{N_JSX_OPEN, p.finRng(rng), id, name, attrs, closed}, nil
}

func (p *Parser) jsxExpr(rng span.Range) (Node, error) {
	p.lexer.PushMode(LM_NONE, true)

	locAfterBrace := p.rng()

	var expr Node
	var empty Node
	var err error

	ahead := p.lexer.Peek()
	av := ahead.value

	if av == T_DOT_TRI {
		if expr, err = p.spread(); err != nil {
			return nil, err
		}
		if _, err := p.nextMustTok(T_BRACE_R); err != nil {
			return nil, err
		}
	} else if av == T_BRACE_R {
		tok := p.lexer.Next()
		// adjust loc of the empty node
		locAfterBrace.Hi = tok.rng.Lo
		empty = &JsxEmpty{N_JSX_EMPTY, locAfterBrace}
		expr = &JsxExprSpan{N_JSX_EXPR_SPAN, p.finRng(rng), empty}
	} else {
		if expr, err = p.expr(); err != nil {
			return nil, err
		}
		if _, err := p.nextMustTok(T_BRACE_R); err != nil {
			return nil, err
		}
		expr = &JsxExprSpan{N_JSX_EXPR_SPAN, p.finRng(rng), expr}
	}

	p.lexer.PopMode()
	return expr, nil
}

func (p *Parser) jsxClose(rng span.Range) (Node, error) {
	p.lexer.Next() // `/`

	// fragment
	if p.lexer.Peek().value == T_GT {
		p.lexer.Next()
		p.lexer.PopMode()
		return &JsxClose{N_JSX_CLOSE, p.finRng(rng), nil, ""}, nil
	}

	id, name, err := p.jsxName()
	if err != nil {
		return nil, err
	}
	if _, err := p.nextMustTok(T_GT); err != nil {
		return nil, err
	}
	// balance the `pushMode` at the beginning of the `p.jsx()`
	p.lexer.PopMode()
	return &JsxClose{N_JSX_CLOSE, p.finRng(rng), id, name}, nil
}

func (p *Parser) isCloseTag(open Node, close Node) bool {
	if close.Type() != N_JSX_CLOSE {
		return false
	}
	on := open.(*JsxOpen)
	cn := close.(*JsxClose)
	return on.nameStr == cn.nameStr
}

// whitespace before `LT`
func (p *Parser) jsxWsTxt() Node {
	if p.lexer.state.prevWs.len == 0 {
		return nil
	}

	rng := p.rng()
	prevWs := &p.lexer.state.prevWs
	rng.Lo = prevWs.rng.Lo
	rng.Hi = prevWs.rng.Hi
	prevWs.len = 0
	return &JsxText{N_JSX_TXT, rng, p.RngText(rng)}
}

// `opening` indicates the opening of the tag has presented, so the
// closing tag is deserved to be appearing
func (p *Parser) jsx(root bool, opening bool) (Node, error) {
	tok := p.lexer.Next() // `<`

	p.lexer.PushMode(LM_JSX, true)

	ahead := p.lexer.Peek()
	if ahead.value == T_DIV {
		if !opening {
			return nil, p.errorAt(ahead.value, ahead.rng, ERR_UNEXPECTED_TOKEN)
		}
		return p.jsxClose(tok.rng)
	}

	open, err := p.jsxOpen(tok)
	if err != nil {
		return nil, err
	}

	var close Node
	var children []Node
	var openTag = open.(*JsxOpen)
	var child Node
	if !openTag.closed {
		p.lexer.PushMode(LM_JSX_CHILD, true)
		children = make([]Node, 0)
		for {
			ahead := p.lexer.Peek()
			av := ahead.value
			if av == T_BRACE_L {
				if ws := p.jsxWsTxt(); ws != nil {
					children = append(children, ws)
				}

				tok := p.lexer.Next() // `{`
				child, err = p.jsxExpr(tok.rng)
				if err != nil {
					return nil, err
				}
				if child.Type() == N_SPREAD {
					s := child.(*Spread)
					child = &JsxSpreadChild{N_JSX_CHILD_SPREAD, p.finRng(tok.rng), s.arg}
				}
				children = append(children, child)

				// whitespace after tag
				ahead := p.lexer.Peek()
				av := ahead.value
				if av == T_GT || av == T_BRACE_L {
					if ws := p.jsxWsTxt(); ws != nil {
						children = append(children, ws)
					}
				}
			} else if av == T_LT {
				if ws := p.jsxWsTxt(); ws != nil {
					children = append(children, ws)
				}
				tag, err := p.jsx(false, true)
				if err != nil {
					return nil, err
				}
				if tag.Type() == N_JSX_CLOSE {
					if p.isCloseTag(open, tag) {
						close = tag
						break
					}
					return nil, p.errorAtLoc(tag.Range(), fmt.Sprintf(ERR_TPL_UNBALANCED_JSX_TAG, openTag.nameStr))
				}
				children = append(children, tag)

				// whitespace after tag
				ahead := p.lexer.Peek()
				av := ahead.value
				if av == T_GT || av == T_BRACE_L {
					if ws := p.jsxWsTxt(); ws != nil {
						children = append(children, ws)
					}
				}
			} else if av == T_JSX_TXT {
				tok := p.lexer.Next()
				child := &JsxText{N_JSX_TXT, p.finRng(tok.rng), tok.ext.(string)}
				children = append(children, child)
			} else if av == T_EOF {
				return nil, p.errorAtLoc(ahead.rng, ERR_UNTERMINATED_JSX_CONTENTS)
			} else if av == T_ILLEGAL {
				return nil, p.errorTok(ahead)
			}
		}
		p.lexer.PopMode()
	}

	// element is closed
	p.lexer.PopMode()
	ahead = p.lexer.Peek()
	// here `T_LT` is not say that the followed node is a jsx-open tag since the close-tag also starts with `<`
	// however if we combined the `is root` condition with is `is LI` then we can report the error `ERR_JSX_ADJACENT_ELEM_SHOULD_BE_WRAPPED`
	// correctly since the root jsx element must stand alone
	if ahead.value == T_LT && root {
		return nil, p.errorAt(ahead.value, ahead.rng, ERR_JSX_ADJACENT_ELEM_SHOULD_BE_WRAPPED)
	}
	return &JsxElem{N_JSX_ELEM, p.finRng(open.Range()), open, close, children}, nil
}
