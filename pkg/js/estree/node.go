package estree

// https://github.com/estree/estree/blob/master/es5.md

type Node interface{}

type Position struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

type SrcLoc struct {
	Source string    `json:"source"`
	Start  *Position `json:"start"`
	End    *Position `json:"end"`
	Range  *SrcRange `json:"range"`
}

type SrcRange struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

// https://github.com/estree/estree/blob/master/es5.md#programs
type Program struct {
	Type       string  `json:"type"`
	Loc        *SrcLoc `json:"loc"`
	SourceType string  `json:"sourceType"` // "script" | "module"
	Body       []Node  `json:"body"`       // [ Directive | Statement ]
}

// https://github.com/estree/estree/blob/master/es5.md#identifier
type Identifier struct {
	Type string  `json:"type"`
	Loc  *SrcLoc `json:"loc"`
	Name string  `json:"name"`
}

// https://github.com/estree/estree/blob/master/es5.md#literal
type Literal struct {
	Type  string      `json:"type"`
	Loc   *SrcLoc     `json:"loc"`
	Value interface{} `json:"value"` // string | boolean | null | number | RegExp | bigint(es2020)
}

type Regexp struct {
	Pattern string `json:"pattern"`
	Flags   string `json:"flags"`
}

// https://github.com/estree/estree/blob/master/es5.md#regexpliteral
type RegExpLiteral struct {
	Type   string      `json:"type"`
	Loc    *SrcLoc     `json:"loc"`
	Value  interface{} `json:"value"`
	Regexp *Regexp     `json:"regexp"`
}

// https://github.com/estree/estree/blob/master/es2020.md#bigintliteral
type BigIntLiteral struct {
	Type   string      `json:"type"`
	Loc    *SrcLoc     `json:"loc"`
	Value  interface{} `json:"value"`
	Bigint string      `json:"bigint"`
}

type Expression interface{}

type Statement interface{}

// https://github.com/estree/estree/blob/master/es5.md#expressionstatement
type ExpressionStatement struct {
	Type       string     `json:"type"`
	Loc        *SrcLoc    `json:"loc"`
	Expression Expression `json:"expression"`
}

// https://github.com/estree/estree/blob/master/es5.md#directive
type Directive struct {
	Type      string  `json:"type"`
	Loc       *SrcLoc `json:"loc"`
	Directive string  `json:"directive"`
}

type EmptyStatement struct {
	Type string  `json:"type"`
	Loc  *SrcLoc `json:"loc"`
}

type DebuggerStatement struct {
	Type string  `json:"type"`
	Loc  *SrcLoc `json:"loc"`
}

// https://github.com/estree/estree/blob/master/es5.md#blockstatement
type BlockStatement struct {
	Type string      `json:"type"`
	Loc  *SrcLoc     `json:"loc"`
	Body []Statement `json:"body"`
}

// https://github.com/estree/estree/blob/master/es5.md#withstatement
type WithStatement struct {
	Type   string     `json:"type"`
	Loc    *SrcLoc    `json:"loc"`
	Object Expression `json:"object"`
	Body   Statement  `json:"body"`
}

// https://github.com/estree/estree/blob/master/es5.md#returnstatement
type ReturnStatement struct {
	Type     string     `json:"type"`
	Loc      *SrcLoc    `json:"loc"`
	Argument Expression `json:"argument"`
}

// https://github.com/estree/estree/blob/master/es5.md#labeledstatement
type LabeledStatement struct {
	Type  string     `json:"type"`
	Loc   *SrcLoc    `json:"loc"`
	Label Expression `json:"label"`
	Body  Statement  `json:"body"`
}

// https://github.com/estree/estree/blob/master/es5.md#breakstatement
type BreakStatement struct {
	Type  string     `json:"type"`
	Loc   *SrcLoc    `json:"loc"`
	Label Expression `json:"label"`
}

// https://github.com/estree/estree/blob/master/es5.md#continuestatement
type ContinueStatement struct {
	Type  string     `json:"type"`
	Loc   *SrcLoc    `json:"loc"`
	Label Expression `json:"label"`
}

// https://github.com/estree/estree/blob/master/es5.md#ifstatement
type IfStatement struct {
	Type       string     `json:"type"`
	Loc        *SrcLoc    `json:"loc"`
	Test       Expression `json:"test"`
	Consequent Statement  `json:"consequent"`
	Alternate  Statement  `json:"alternate"`
}

// https://github.com/estree/estree/blob/master/es5.md#switchstatement
type SwitchStatement struct {
	Type         string        `json:"type"`
	Loc          *SrcLoc       `json:"loc"`
	Discriminant Expression    `json:"discriminant"`
	Cases        []*SwitchCase `json:"cases"`
}

// https://github.com/estree/estree/blob/master/es5.md#switchcase
type SwitchCase struct {
	Type       string      `json:"type"`
	Loc        *SrcLoc     `json:"loc"`
	Test       Expression  `json:"test"`
	Consequent []Statement `json:"consequent"`
}

// https://github.com/estree/estree/blob/master/es5.md#throwstatement
type ThrowStatement struct {
	Type     string     `json:"type"`
	Loc      *SrcLoc    `json:"loc"`
	Argument Expression `json:"argument"`
}

// https://github.com/estree/estree/blob/master/es5.md#trystatement
type TryStatement struct {
	Type      string     `json:"type"`
	Loc       *SrcLoc    `json:"loc"`
	Block     Statement  `json:"block"`
	Handler   Expression `json:"handler"`
	Finalizer Statement  `json:"finalizer"`
}

// https://github.com/estree/estree/blob/master/es5.md#catchclause
type CatchClause struct {
	Type  string    `json:"type"`
	Loc   *SrcLoc   `json:"loc"`
	Param Pattern   `json:"param"` // `Pattern | null` from es2019
	Body  Statement `json:"body"`
}

// https://github.com/estree/estree/blob/master/es5.md#whilestatement
type WhileStatement struct {
	Type string     `json:"type"`
	Loc  *SrcLoc    `json:"loc"`
	Test Expression `json:"test"`
	Body Statement  `json:"body"`
}

// https://github.com/estree/estree/blob/master/es5.md#dowhilestatement
type DoWhileStatement struct {
	Type string     `json:"type"`
	Loc  *SrcLoc    `json:"loc"`
	Test Expression `json:"test"`
	Body Statement  `json:"body"`
}

// https://github.com/estree/estree/blob/master/es5.md#forstatement
type ForStatement struct {
	Type   string     `json:"type"`
	Loc    *SrcLoc    `json:"loc"`
	Init   Node       `json:"init"` // VariableDeclaration | Expression | null
	Test   Expression `json:"test"`
	Update Expression `json:"update"`
	Body   Statement  `json:"body"`
}

// https://github.com/estree/estree/blob/master/es5.md#forinstatement
type ForInStatement struct {
	Type  string     `json:"type"`
	Loc   *SrcLoc    `json:"loc"`
	Left  Node       `json:"left"`
	Right Expression `json:"right"`
	Body  Statement  `json:"body"`
}

// https://github.com/estree/estree/blob/master/es2015.md#forofstatement
type ForOfStatement struct {
	Type  string     `json:"type"`
	Loc   *SrcLoc    `json:"loc"`
	Left  Node       `json:"left"`
	Right Expression `json:"right"`
	Body  Statement  `json:"body"`
	Await bool       `json:"await"`
}

type Declaration interface{}

// https://github.com/estree/estree/blob/master/es5.md#functiondeclaration
type FunctionDeclaration struct {
	Type      string  `json:"type"`
	Loc       *SrcLoc `json:"loc"`
	Id        Node    `json:"id"`
	Params    []Node  `json:"params"`
	Body      Node    `json:"body"`
	Generator bool    `json:"generator"`
	Async     bool    `json:"async"`
}

// https://github.com/estree/estree/blob/master/es5.md#variabledeclaration
type VariableDeclaration struct {
	Type         string                `json:"type"`
	Loc          *SrcLoc               `json:"loc"`
	Kind         string                `json:"kind"`
	Declarations []*VariableDeclarator `json:"declarations"`
}

// https://github.com/estree/estree/blob/master/es5.md#variabledeclarator
type VariableDeclarator struct {
	Type string     `json:"type"`
	Loc  *SrcLoc    `json:"loc"`
	Id   Pattern    `json:"id"`
	Init Expression `json:"init"`
}

type ThisExpression struct {
	Type string  `json:"type"`
	Loc  *SrcLoc `json:"loc"`
}

// https://github.com/estree/estree/blob/master/es5.md#arrayexpression
type ArrayExpression struct {
	Type     string       `json:"type"`
	Loc      *SrcLoc      `json:"loc"`
	Elements []Expression `json:"elements"`
}

// https://github.com/estree/estree/blob/master/es5.md#objectexpression
type ObjectExpression struct {
	Type       string      `json:"type"`
	Loc        *SrcLoc     `json:"loc"`
	Properties []*Property `json:"properties"`
}

// https://github.com/estree/estree/blob/master/es5.md#functionexpression
type FunctionExpression struct {
	Type       string  `json:"type"`
	Loc        *SrcLoc `json:"loc"`
	Id         Node    `json:"id"`
	Params     []Node  `json:"params"`
	Body       Node    `json:"body"`
	Generator  bool    `json:"generator"`
	Async      bool    `json:"async"`
	Expression bool    `json:"expression"`
}

// https://github.com/estree/estree/blob/master/es5.md#unaryexpression
type UnaryExpression struct {
	Type     string     `json:"type"`
	Loc      *SrcLoc    `json:"loc"`
	Operator string     `json:"operator"` //  "-" | "+" | "!" | "~" | "typeof" | "void" | "delete"
	Prefix   bool       `json:"prefix"`
	Argument Expression `json:"argument"`
}

// https://github.com/estree/estree/blob/master/es5.md#updateexpression
type UpdateExpression struct {
	Type     string     `json:"type"`
	Loc      *SrcLoc    `json:"loc"`
	Operator string     `json:"operator"` // "++" | "--"
	Argument Expression `json:"argument"`
	Prefix   bool       `json:"prefix"`
}

// https://github.com/estree/estree/blob/master/es5.md#binaryexpression
type BinaryExpression struct {
	Type     string     `json:"type"`
	Loc      *SrcLoc    `json:"loc"`
	Operator string     `json:"operator"`
	Left     Expression `json:"left"`
	Right    Expression `json:"right"`
}

// https://github.com/estree/estree/blob/master/es5.md#assignmentexpression
type AssignmentExpression struct {
	Type     string     `json:"type"`
	Loc      *SrcLoc    `json:"loc"`
	Operator string     `json:"operator"`
	Left     Node       `json:"left"`
	Right    Expression `json:"right"`
}

// https://github.com/estree/estree/blob/master/es5.md#logicalexpression
type LogicalExpression struct {
	Type     string     `json:"type"`
	Loc      *SrcLoc    `json:"loc"`
	Operator string     `json:"operator"`
	Left     Expression `json:"left"`
	Right    Expression `json:"right"`
}

// https://github.com/estree/estree/blob/master/es5.md#memberexpression
type MemberExpression struct {
	Type     string     `json:"type"`
	Loc      *SrcLoc    `json:"loc"`
	Object   Expression `json:"object"`
	Property Expression `json:"property"`
	Computed bool       `json:"computed"`
	Optional bool       `json:"optional"`
}

// https://github.com/estree/estree/blob/master/es5.md#conditionalexpression
type ConditionalExpression struct {
	Type       string     `json:"type"`
	Loc        *SrcLoc    `json:"loc"`
	Test       Expression `json:"test"`
	Consequent Expression `json:"consequent"`
	Alternate  Expression `json:"alternate"`
}

// https://github.com/estree/estree/blob/master/es5.md#callexpression
type CallExpression struct {
	Type      string       `json:"type"`
	Loc       *SrcLoc      `json:"loc"`
	Callee    Expression   `json:"callee"`
	Arguments []Expression `json:"arguments"`
	Optional  bool         `json:"optional"`
}

// https://github.com/estree/estree/blob/master/es5.md#newexpression
type NewExpression struct {
	Type      string       `json:"type"`
	Loc       *SrcLoc      `json:"loc"`
	Callee    Expression   `json:"callee"`
	Arguments []Expression `json:"arguments"`
}

// https://github.com/estree/estree/blob/master/es5.md#sequenceexpression
type SequenceExpression struct {
	Type        string       `json:"type"`
	Loc         *SrcLoc      `json:"loc"`
	Expressions []Expression `json:"expressions"`
}

// https://github.com/estree/estree/blob/master/es2015.md#expressions
type Super struct {
	Type string  `json:"type"`
	Loc  *SrcLoc `json:"loc"`
}

// https://github.com/estree/estree/blob/master/es2015.md#expressions
type SpreadElement struct {
	Type     string     `json:"type"`
	Loc      *SrcLoc    `json:"loc"`
	Argument Expression `json:"argument"`
}

// https://github.com/estree/estree/blob/master/es2015.md#arrowfunctionexpression
type ArrowFunctionExpression struct {
	Type       string      `json:"type"`
	Loc        *SrcLoc     `json:"loc"`
	Id         *Identifier `json:"id"`
	Params     []Pattern   `json:"params"`
	Body       Node        `json:"body"`
	Generator  bool        `json:"generator"`
	Async      bool        `json:"async"`
	Expression bool        `json:"expression"`
}

// https://github.com/estree/estree/blob/master/es2015.md#yieldexpression
type YieldExpression struct {
	Type     string     `json:"type"`
	Loc      *SrcLoc    `json:"loc"`
	Argument Expression `json:"argument"`
	Delegate bool       `json:"delegate"`
}

// https://github.com/estree/estree/blob/master/es2017.md#awaitexpression
type AwaitExpression struct {
	Type     string     `json:"type"`
	Loc      *SrcLoc    `json:"loc"`
	Argument Expression `json:"argument"`
}

// https://github.com/estree/estree/blob/master/es2015.md#templateliteral
type TemplateLiteral struct {
	Type        string       `json:"type"`
	Loc         *SrcLoc      `json:"loc"`
	Quasis      []Expression `json:"quasis"`
	Expressions []Expression `json:"expressions"`
}

// https://github.com/estree/estree/blob/master/es2015.md#taggedtemplateexpression
type TaggedTemplateExpression struct {
	Type  string     `json:"type"`
	Loc   *SrcLoc    `json:"loc"`
	Tag   Expression `json:"tag"`
	Quasi Expression `json:"quasi"`
}

// https://github.com/estree/estree/blob/master/es2015.md#templateelement
type TemplateElement struct {
	Type  string                `json:"type"`
	Loc   *SrcLoc               `json:"loc"`
	Tail  bool                  `json:"tail"`
	Value *TemplateElementValue `json:"value"`
}

type TemplateElementValue struct {
	Cooked string `json:"cooked"`
	Raw    string `json:"raw"`
}

// https://github.com/estree/estree/blob/master/es2020.md#chainexpression
type ChainExpression struct {
	Type       string     `json:"type"`
	Loc        *SrcLoc    `json:"loc"`
	Expression Expression `json:"expression"` // CallExpression | MemberExpression
}

// https://github.com/estree/estree/blob/master/es2020.md#importexpression
type ImportExpression struct {
	Type   string     `json:"type"`
	Loc    *SrcLoc    `json:"loc"`
	Source Expression `json:"source"`
}

// https://github.com/estree/estree/blob/master/es2015.md#patterns
type Pattern interface{}

// https://github.com/estree/estree/blob/master/es2015.md#expressions
type Property struct {
	Type      string     `json:"type"`
	Loc       *SrcLoc    `json:"loc"`
	Key       Expression `json:"key"`
	Value     Expression `json:"value"`
	Kind      string     `json:"kind"`
	Method    bool       `json:"method"`
	Shorthand bool       `json:"shorthand"`
	Computed  bool       `json:"computed"`
}

// https://github.com/estree/estree/blob/master/es2015.md#objectpattern
type AssignmentProperty struct {
	Type      string     `json:"type"`
	Loc       *SrcLoc    `json:"loc"`
	Key       Expression `json:"key"`
	Method    bool       `json:"method"`
	Shorthand bool       `json:"shorthand"`
	Computed  bool       `json:"computed"`
	Value     Pattern    `json:"value"`
	Kind      string     `json:"kind"`
}

// https://github.com/estree/estree/blob/master/es2015.md#objectpattern
type ObjectPattern struct {
	Type       string               `json:"type"`
	Loc        *SrcLoc              `json:"loc"`
	Properties []AssignmentProperty `json:"properties"`
}

// https://github.com/estree/estree/blob/master/es2015.md#arraypattern
type ArrayPattern struct {
	Type     string    `json:"type"`
	Loc      *SrcLoc   `json:"loc"`
	Elements []Pattern `json:"elements"`
}

// https://github.com/estree/estree/blob/master/es2015.md#restelement
type RestElement struct {
	Type     string  `json:"type"`
	Loc      *SrcLoc `json:"loc"`
	Argument Pattern `json:"argument"`
}

// https://github.com/estree/estree/blob/master/es2015.md#assignmentpattern
type AssignmentPattern struct {
	Type  string     `json:"type"`
	Loc   *SrcLoc    `json:"loc"`
	Left  Pattern    `json:"left"`
	Right Expression `json:"right"`
}

// https://github.com/estree/estree/blob/master/es2015.md#classbody
type ClassBody struct {
	Type string       `json:"type"`
	Loc  *SrcLoc      `json:"loc"`
	Body []Expression `json:"body"` // MethodDefinition | PropertyDefinition | StaticBlock
}

// https://github.com/estree/estree/blob/master/es2015.md#methoddefinition
type MethodDefinition struct {
	Type     string     `json:"type"`
	Loc      *SrcLoc    `json:"loc"`
	Key      Expression `json:"key"`
	Value    Expression `json:"value"`
	Kind     string     `json:"kind"` // "constructor" | "method" | "get" | "set"
	Computed bool       `json:"computed"`
	Static   bool       `json:"static"`
}

// https://github.com/estree/estree/blob/master/es2022.md#propertydefinition
type PropertyDefinition struct {
	Type     string     `json:"type"`
	Loc      *SrcLoc    `json:"loc"`
	Key      Expression `json:"key"`
	Value    Expression `json:"value"`
	Computed bool       `json:"computed"`
	Static   bool       `json:"static"`
}

// https://github.com/estree/estree/blob/master/es2015.md#classdeclaration
type ClassDeclaration struct {
	Type       string     `json:"type"`
	Loc        *SrcLoc    `json:"loc"`
	Id         Expression `json:"id"`
	SuperClass Expression `json:"superClass"`
	Body       Expression `json:"body"`
}

// https://github.com/estree/estree/blob/master/es2015.md#classexpression
type ClassExpression struct {
	Type       string     `json:"type"`
	Loc        *SrcLoc    `json:"loc"`
	Id         Expression `json:"id"`
	SuperClass Expression `json:"superClass"`
	Body       Expression `json:"body"`
}

// https://github.com/estree/estree/blob/master/es2022.md#privateidentifier
type PrivateIdentifier struct {
	Type string  `json:"type"`
	Loc  *SrcLoc `json:"loc"`
	Name string  `json:"name"`
}

// https://github.com/estree/estree/blob/master/es2022.md#staticblock
type StaticBlock struct {
	Type string      `json:"type"`
	Loc  *SrcLoc     `json:"loc"`
	Body []Statement `json:"body"`
}

// https://github.com/estree/estree/blob/master/es2015.md#metaproperty
type MetaProperty struct {
	Type     string     `json:"type"`
	Loc      *SrcLoc    `json:"loc"`
	Meta     Expression `json:"meta"`
	Property Expression `json:"property"`
}

type ModuleDeclaration interface{}

// https://github.com/estree/estree/blob/master/es2015.md#importdeclaration
type ImportDeclaration struct {
	Type       string  `json:"type"`
	Loc        *SrcLoc `json:"loc"`
	Specifiers []Node  `json:"specifiers"` // [ ImportSpecifier | ImportDefaultSpecifier | ImportNamespaceSpecifier ]
	Source     Literal `json:"source"`
}

// https://github.com/estree/estree/blob/master/es2015.md#importspecifier
type ImportSpecifier struct {
	Type     string      `json:"type"`
	Loc      *SrcLoc     `json:"loc"`
	Local    *Identifier `json:"local"`
	Imported *Identifier `json:"imported"`
}

// https://github.com/estree/estree/blob/master/es2015.md#importdefaultspecifier
type ImportDefaultSpecifier struct {
	Type  string      `json:"type"`
	Loc   *SrcLoc     `json:"loc"`
	Local *Identifier `json:"local"`
}

// https://github.com/estree/estree/blob/master/es2015.md#importnamespacespecifier
type ImportNamespaceSpecifier struct {
	Type  string      `json:"type"`
	Loc   *SrcLoc     `json:"loc"`
	Local *Identifier `json:"local"`
}

// https://github.com/estree/estree/blob/master/es2015.md#exportnameddeclaration
type ExportNamedDeclaration struct {
	Type        string             `json:"type"`
	Loc         *SrcLoc            `json:"loc"`
	Declaration Declaration        `json:"declaration"` // Declaration | null
	Specifiers  []*ExportSpecifier `json:"specifiers"`
	Source      Literal            `json:"source"` // Literal | null
}

// https://github.com/estree/estree/blob/master/es2015.md#exportspecifier
type ExportSpecifier struct {
	Type     string      `json:"type"`
	Loc      *SrcLoc     `json:"loc"`
	Local    *Identifier `json:"local"`
	Exported *Identifier `json:"exported"`
}

type ExportDefaultDeclaration struct {
	Type string  `json:"type"`
	Loc  *SrcLoc `json:"loc"`

	// AnonymousDefaultExportedFunctionDeclaration | FunctionDeclaration | AnonymousDefaultExportedClassDeclaration | ClassDeclaration | Expression
	Declaration Node `json:"declaration"`
}

// https://github.com/estree/estree/blob/master/es2015.md#exportalldeclaration
type ExportAllDeclaration struct {
	Type   string  `json:"type"`
	Loc    *SrcLoc `json:"loc"`
	Source Literal `json:"source"`
}
