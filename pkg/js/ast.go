package js

// AST nodes is described as https://github.com/estree/estree/blob/master/es5.md but with
// one big difference - flatterned struct is used instead of inheritance

type Node interface {
	Type() NodeType
	Loc() *Loc
}

type NodeType int

const (
	N_ILLEGAL NodeType = iota
	N_PROG

	N_STMT_BEGIN
	N_STMT_EXPR
	N_STMT_EMPTY
	N_STMT_END

	N_EXPR_BEGIN

	N_LITERAL_BEGIN
	N_LITERAL_NULL
	N_LITERAL_BOOL
	N_LITERAL_NUMERIC
	N_LITERAL_STRING
	N_LITERAL_END

	N_EXPR_NEW
	N_EXPR_MEMBER
	N_EXPR_CALL
	N_EXPR_BIN
	N_EXPR_UNARY
	N_EXPR_UPDATE
	N_EXPR_COND
	N_EXPR_ASSIGN

	N_NAME

	N_EXPR_END
)

type Prog struct {
	typ   NodeType
	loc   *Loc
	stmts []Node
}

func NewProg() *Prog {
	return &Prog{N_PROG, &Loc{}, make([]Node, 0)}
}

func (n *Prog) Type() NodeType {
	return n.typ
}

func (n *Prog) Loc() *Loc {
	return n.loc
}

type ExprStmt struct {
	typ  NodeType
	loc  *Loc
	expr Node
}

func NewExprStmt() *ExprStmt {
	return &ExprStmt{N_STMT_EXPR, &Loc{}, nil}
}

func (n *ExprStmt) Type() NodeType {
	return n.typ
}

func (n *ExprStmt) Loc() *Loc {
	return n.loc
}

type EmptyStmt struct {
	typ NodeType
	loc *Loc
}

func NewEmptyStmt() *EmptyStmt {
	return &EmptyStmt{N_STMT_EMPTY, &Loc{}}
}

func (n *EmptyStmt) Type() NodeType {
	return n.typ
}

func (n *EmptyStmt) Loc() *Loc {
	return n.loc
}

type NullLit struct {
	typ NodeType
	loc *Loc
}

func NewNullLit() *NullLit {
	return &NullLit{N_LITERAL_NULL, &Loc{}}
}

type BoolLit struct {
	typ NodeType
	loc *Loc
}

func NewBoolLiteral() *BoolLit {
	return &BoolLit{N_LITERAL_BOOL, &Loc{}}
}

type NumLit struct {
	typ NodeType
	loc *Loc
	val *Token
}

func NewNumLit() *NumLit {
	return &NumLit{N_LITERAL_NUMERIC, &Loc{}, nil}
}

func (n *NumLit) Type() NodeType {
	return n.typ
}

func (n *NumLit) Loc() *Loc {
	return n.loc
}

type StrLit struct {
	typ NodeType
	loc *Loc
}

func NewStrLit() *StrLit {
	return &StrLit{N_LITERAL_BOOL, &Loc{}}
}

type Ident struct {
	typ     NodeType
	loc     *Loc
	val     *Token
	private bool
}

func NewIdent() *Ident {
	return &Ident{N_NAME, &Loc{}, nil, false}
}

func (n *Ident) Type() NodeType {
	return n.typ
}

func (n *Ident) Loc() *Loc {
	return n.loc
}

type NewExpr struct {
	typ  NodeType
	loc  *Loc
	expr Node
}

func NewNewExpr() *NewExpr {
	return &NewExpr{N_EXPR_NEW, &Loc{}, nil}
}

func (n *NewExpr) Type() NodeType {
	return n.typ
}

func (n *NewExpr) Loc() *Loc {
	return n.loc
}

type MemberExpr struct {
	typ      NodeType
	loc      *Loc
	obj      Node
	prop     Node
	compute  bool
	optional bool
}

func (n *MemberExpr) Type() NodeType {
	return n.typ
}

func (n *MemberExpr) Loc() *Loc {
	return n.loc
}

type CallExpr struct {
	typ  NodeType
	loc  *Loc
	expr Node
}

func (n *CallExpr) Type() NodeType {
	return n.typ
}

func (n *CallExpr) Loc() *Loc {
	return n.loc
}

type BinExpr struct {
	typ NodeType
	loc *Loc
	op  *Token
	lhs Node
	rhs Node
}

func NewBinExpr() *BinExpr {
	return &BinExpr{N_EXPR_BIN, nil, nil, nil, nil}
}

func (n *BinExpr) Type() NodeType {
	return n.typ
}

func (n *BinExpr) Loc() *Loc {
	return n.loc
}

type UnaryExpr struct {
	typ NodeType
	loc *Loc
	op  *Token
	arg Node
}

func NewUnaryExpr() *UnaryExpr {
	return &UnaryExpr{N_EXPR_UNARY, nil, nil, nil}
}

func (n *UnaryExpr) Type() NodeType {
	return n.typ
}

func (n *UnaryExpr) Loc() *Loc {
	return n.loc
}

type UpdateExpr struct {
	typ    NodeType
	loc    *Loc
	op     *Token
	prefix bool
	arg    Node
}

func NewUpdateExpr() *UpdateExpr {
	return &UpdateExpr{N_EXPR_UPDATE, nil, nil, false, nil}
}

func (n *UpdateExpr) Type() NodeType {
	return n.typ
}

func (n *UpdateExpr) Loc() *Loc {
	return n.loc
}

type CondExpr struct {
	typ  NodeType
	loc  *Loc
	test Node
	cons Node
	alt  Node
}

func NewCondExpr() *CondExpr {
	return &CondExpr{N_EXPR_COND, nil, nil, nil, nil}
}

func (n *CondExpr) Type() NodeType {
	return n.typ
}

func (n *CondExpr) Loc() *Loc {
	return n.loc
}

type AssignExpr struct {
	typ NodeType
	loc *Loc
	op  *Token
	lhs Node
	rhs Node
}

func NewAssignExpr() *AssignExpr {
	return &AssignExpr{N_EXPR_ASSIGN, nil, nil, nil, nil}
}

func (n *AssignExpr) Type() NodeType {
	return n.typ
}

func (n *AssignExpr) Loc() *Loc {
	return n.loc
}

type Loc struct {
	src   *Source
	begin Position
	end   Position
}

func (l *Loc) Clone() *Loc {
	return &Loc{
		src:   l.src,
		begin: l.begin.Clone(),
		end:   l.end.Clone(),
	}
}
