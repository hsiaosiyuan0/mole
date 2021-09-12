package js

// AST nodes is described as: https://github.com/estree/estree/blob/master/es5.md
// flatterned struct is used instead of inheritance

type Node interface {
	Type() NodeType
	Loc() *Loc
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

type NodeType int

const (
	N_ILLEGAL NodeType = iota
	N_PROG

	N_STMT_BEGIN
	N_STMT_EXPR
	N_STMT_EMPTY
	N_STMT_VAR_DEC
	N_STMT_FN
	N_STMT_BLOCK
	N_STMT_DO_WHILE
	N_STMT_WHILE
	N_STMT_FOR
	N_STMT_FOR_IN_OF
	N_STMT_IF
	N_STMT_SWITCH
	N_STMT_BRK
	N_STMT_CONT
	N_STMT_LABEL
	N_STMT_RET
	N_STMT_THROW
	N_STMT_TRY
	N_STMT_DEBUG
	N_STMT_CLASS
	N_STMT_END

	N_EXPR_BEGIN
	N_LIT_BEGIN
	N_LIT_NULL
	N_LIT_BOOL
	N_LIT_NUM
	N_LIT_STR
	N_LIT_ARR
	N_LIT_OBJ
	N_LIT_REGEXP
	N_LIT_END

	N_EXPR_NEW
	N_EXPR_MEMBER
	N_EXPR_CALL
	N_EXPR_BIN
	N_EXPR_UNARY
	N_EXPR_UPDATE
	N_EXPR_COND
	N_EXPR_ASSIGN
	N_EXPR_FN
	N_EXPR_THIS
	N_EXPR_PAREN
	N_EXPR_ARROW
	N_EXPR_SEQ
	N_EXPR_CLASS
	N_EXPR_TPL
	N_NAME

	N_VAR_DEC
	N_PATTERN_REST
	N_PATTERN_ARRAY
	N_PATTERN_ASSIGN
	N_PATTERN_OBJ
	N_PROP
	N_SPREAD
	N_SWITCH_CASE
	N_CATCH
	N_ClASS_BODY
	N_STATIC_BLOCK
	N_METHOD
	N_FIELD

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
	return &NullLit{N_LIT_NULL, &Loc{}}
}

type BoolLit struct {
	typ NodeType
	loc *Loc
}

func NewBoolLiteral() *BoolLit {
	return &BoolLit{N_LIT_BOOL, &Loc{}}
}

type NumLit struct {
	typ NodeType
	loc *Loc
	val *Token
}

func NewNumLit() *NumLit {
	return &NumLit{N_LIT_NUM, &Loc{}, nil}
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
	val *Token
}

func NewStrLit() *StrLit {
	return &StrLit{N_LIT_STR, &Loc{}, nil}
}

func (n *StrLit) Type() NodeType {
	return n.typ
}

func (n *StrLit) Loc() *Loc {
	return n.loc
}

type RegexpLit struct {
	typ NodeType
	loc *Loc
	val *Token
}

func (n *RegexpLit) Type() NodeType {
	return n.typ
}

func (n *RegexpLit) Loc() *Loc {
	return n.loc
}

type ArrLit struct {
	typ   NodeType
	loc   *Loc
	elems []Node
}

func (n *ArrLit) Type() NodeType {
	return n.typ
}

func (n *ArrLit) Loc() *Loc {
	return n.loc
}

type Spread struct {
	typ NodeType
	loc *Loc
	arg Node
}

func (n *Spread) Type() NodeType {
	return n.typ
}

func (n *Spread) Loc() *Loc {
	return n.loc
}

type ObjLit struct {
	typ   NodeType
	loc   *Loc
	props []Node
}

func (n *ObjLit) Type() NodeType {
	return n.typ
}

func (n *ObjLit) Loc() *Loc {
	return n.loc
}

type Ident struct {
	typ NodeType
	loc *Loc
	val *Token
	pvt bool
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
	typ    NodeType
	loc    *Loc
	callee Node
	args   []Node
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

type VarDecStmt struct {
	typ     NodeType
	loc     *Loc
	kind    TokenValue
	decList []*VarDec
}

func NewVarDecStmt() *VarDecStmt {
	return &VarDecStmt{N_STMT_VAR_DEC, nil, T_ILLEGAL, make([]*VarDec, 0, 1)}
}

func (n *VarDecStmt) Type() NodeType {
	return n.typ
}

func (n *VarDecStmt) Loc() *Loc {
	return n.loc
}

type VarDec struct {
	typ  NodeType
	loc  *Loc
	id   Node
	init Node
}

func (n *VarDec) Type() NodeType {
	return n.typ
}

func (n *VarDec) Loc() *Loc {
	return n.loc
}

type ArrayPattern struct {
	typ   NodeType
	loc   *Loc
	elems []Node
}

func (n *ArrayPattern) Type() NodeType {
	return n.typ
}

func (n *ArrayPattern) Loc() *Loc {
	return n.loc
}

type AssignPattern struct {
	typ   NodeType
	loc   *Loc
	left  Node
	right Node
}

func (n *AssignPattern) Type() NodeType {
	return n.typ
}

func (n *AssignPattern) Loc() *Loc {
	return n.loc
}

type RestPattern struct {
	typ NodeType
	loc *Loc
	arg Node
}

func (n *RestPattern) Type() NodeType {
	return n.typ
}

func (n *RestPattern) Loc() *Loc {
	return n.loc
}

type ObjPattern struct {
	typ   NodeType
	loc   *Loc
	props []Node
}

func (n *ObjPattern) Type() NodeType {
	return n.typ
}

func (n *ObjPattern) Loc() *Loc {
	return n.loc
}

type Prop struct {
	typ      NodeType
	loc      *Loc
	key      Node
	value    Node
	computed bool
}

func (n *Prop) Type() NodeType {
	return n.typ
}

func (n *Prop) Loc() *Loc {
	return n.loc
}

type FnDec struct {
	typ       NodeType
	loc       *Loc
	id        Node
	generator bool
	async     bool
	params    []Node
	body      Node
}

func (n *FnDec) Type() NodeType {
	return n.typ
}

func (n *FnDec) Loc() *Loc {
	return n.loc
}

type ArrowFn struct {
	typ       NodeType
	loc       *Loc
	generator bool
	async     bool
	params    []Node
	body      Node
}

func (n *ArrowFn) Type() NodeType {
	return n.typ
}

func (n *ArrowFn) Loc() *Loc {
	return n.loc
}

type BlockStmt struct {
	typ  NodeType
	loc  *Loc
	body []Node
}

func (n *BlockStmt) Type() NodeType {
	return n.typ
}

func (n *BlockStmt) Loc() *Loc {
	return n.loc
}

type DoWhileStmt struct {
	typ  NodeType
	loc  *Loc
	test Node
	body Node
}

func (n *DoWhileStmt) Type() NodeType {
	return n.typ
}

func (n *DoWhileStmt) Loc() *Loc {
	return n.loc
}

type WhileStmt struct {
	typ  NodeType
	loc  *Loc
	test Node
	body Node
}

func (n *WhileStmt) Type() NodeType {
	return n.typ
}

func (n *WhileStmt) Loc() *Loc {
	return n.loc
}

type ForStmt struct {
	typ    NodeType
	loc    *Loc
	init   Node
	test   Node
	update Node
	body   Node
}

func (n *ForStmt) Type() NodeType {
	return n.typ
}

func (n *ForStmt) Loc() *Loc {
	return n.loc
}

type ForInOfStmt struct {
	typ   NodeType
	loc   *Loc
	in    bool
	await bool
	left  Node
	right Node
	body  Node
}

func (n *ForInOfStmt) Type() NodeType {
	return n.typ
}

func (n *ForInOfStmt) Loc() *Loc {
	return n.loc
}

type IfStmt struct {
	typ  NodeType
	loc  *Loc
	test Node
	cons Node
	alt  Node
}

func (n *IfStmt) Type() NodeType {
	return n.typ
}

func (n *IfStmt) Loc() *Loc {
	return n.loc
}

type SwitchStmt struct {
	typ   NodeType
	loc   *Loc
	test  Node
	cases []*SwitchCase
}

func (n *SwitchStmt) Type() NodeType {
	return n.typ
}

func (n *SwitchStmt) Loc() *Loc {
	return n.loc
}

type SwitchCase struct {
	typ  NodeType
	loc  *Loc
	test Node // nil in default clause
	cons []Node
}

func (n *SwitchCase) Type() NodeType {
	return n.typ
}

func (n *SwitchCase) Loc() *Loc {
	return n.loc
}

type BrkStmt struct {
	typ   NodeType
	loc   *Loc
	label *Ident
}

func (n *BrkStmt) Type() NodeType {
	return n.typ
}

func (n *BrkStmt) Loc() *Loc {
	return n.loc
}

type ContStmt struct {
	typ   NodeType
	loc   *Loc
	label *Ident
}

func (n *ContStmt) Type() NodeType {
	return n.typ
}

func (n *ContStmt) Loc() *Loc {
	return n.loc
}

type LabelStmt struct {
	typ   NodeType
	loc   *Loc
	label *Ident
	body  Node
}

func (n *LabelStmt) Type() NodeType {
	return n.typ
}

func (n *LabelStmt) Loc() *Loc {
	return n.loc
}

type RetStmt struct {
	typ NodeType
	loc *Loc
	arg Node
}

func (n *RetStmt) Type() NodeType {
	return n.typ
}

func (n *RetStmt) Loc() *Loc {
	return n.loc
}

type ThrowStmt struct {
	typ NodeType
	loc *Loc
	arg Node
}

func (n *ThrowStmt) Type() NodeType {
	return n.typ
}

func (n *ThrowStmt) Loc() *Loc {
	return n.loc
}

type Catch struct {
	typ   NodeType
	loc   *Loc
	param Node
	body  Node
}

func (n *Catch) Type() NodeType {
	return n.typ
}

func (n *Catch) Loc() *Loc {
	return n.loc
}

type TryStmt struct {
	typ   NodeType
	loc   *Loc
	try   *BlockStmt
	catch *Catch
	fin   *BlockStmt
}

func (n *TryStmt) Type() NodeType {
	return n.typ
}

func (n *TryStmt) Loc() *Loc {
	return n.loc
}

type DebugStmt struct {
	typ NodeType
	loc *Loc
}

func (n *DebugStmt) Type() NodeType {
	return n.typ
}

func (n *DebugStmt) Loc() *Loc {
	return n.loc
}

type ClassDec struct {
	typ   NodeType
	loc   *Loc
	id    Node
	super Node
	body  *ClassBody
}

func (n *ClassDec) Type() NodeType {
	return n.typ
}

func (n *ClassDec) Loc() *Loc {
	return n.loc
}

type ClassBody struct {
	typ   NodeType
	loc   *Loc
	elems []Node
}

func (n *ClassBody) Type() NodeType {
	return n.typ
}

func (n *ClassBody) Loc() *Loc {
	return n.loc
}

type Method struct {
	typ      NodeType
	loc      *Loc
	key      Node
	static   bool
	computed bool
	kind     *Token
	value    Node
}

func (n *Method) Type() NodeType {
	return n.typ
}

func (n *Method) Loc() *Loc {
	return n.loc
}

type Field struct {
	typ      NodeType
	loc      *Loc
	key      Node
	static   bool
	computed bool
	value    Node
}

func (n *Field) Type() NodeType {
	return n.typ
}

func (n *Field) Loc() *Loc {
	return n.loc
}

type StaticBlock struct {
	typ  NodeType
	loc  *Loc
	body []Node
}

func (n *StaticBlock) Type() NodeType {
	return n.typ
}

func (n *StaticBlock) Loc() *Loc {
	return n.loc
}

type ThisExpr struct {
	typ NodeType
	loc *Loc
}

func (n *ThisExpr) Type() NodeType {
	return n.typ
}

func (n *ThisExpr) Loc() *Loc {
	return n.loc
}

type ParenExpr struct {
	typ  NodeType
	loc  *Loc
	expr Node
}

func (n *ParenExpr) Type() NodeType {
	return n.typ
}

func (n *ParenExpr) Loc() *Loc {
	return n.loc
}

type SeqExpr struct {
	typ   NodeType
	loc   *Loc
	elems []Node
}

func (n *SeqExpr) Type() NodeType {
	return n.typ
}

func (n *SeqExpr) Loc() *Loc {
	return n.loc
}

type TplExpr struct {
	typ   NodeType
	loc   *Loc
	tag   Node
	elems []Node
}

func (n *TplExpr) Type() NodeType {
	return n.typ
}

func (n *TplExpr) Loc() *Loc {
	return n.loc
}
