package estree

type TSTypeAnnotation struct {
	Type           string  `json:"type"`
	Start          int     `json:"start"`
	End            int     `json:"end"`
	Loc            *SrcLoc `json:"loc"`
	TypeAnnotation Node    `json:"typeAnnotation"`
}

type TSIdentifier struct {
	Type           string  `json:"type"`
	Start          int     `json:"start"`
	End            int     `json:"end"`
	Loc            *SrcLoc `json:"loc"`
	Name           string  `json:"name"`
	Optional       bool    `json:"optional"`
	TypeAnnotation Node    `json:"typeAnnotation"`
}

type TSNumberKeyword struct {
	Type  string  `json:"type"`
	Start int     `json:"start"`
	End   int     `json:"end"`
	Loc   *SrcLoc `json:"loc"`
}

type TSStringKeyword struct {
	Type  string  `json:"type"`
	Start int     `json:"start"`
	End   int     `json:"end"`
	Loc   *SrcLoc `json:"loc"`
}

type TSVoidKeyword struct {
	Type  string  `json:"type"`
	Start int     `json:"start"`
	End   int     `json:"end"`
	Loc   *SrcLoc `json:"loc"`
}

type TSAnyKeyword struct {
	Type  string  `json:"type"`
	Start int     `json:"start"`
	End   int     `json:"end"`
	Loc   *SrcLoc `json:"loc"`
}

type TSBooleanKeyword struct {
	Type  string  `json:"type"`
	Start int     `json:"start"`
	End   int     `json:"end"`
	Loc   *SrcLoc `json:"loc"`
}

type TSThisType struct {
	Type  string  `json:"type"`
	Start int     `json:"start"`
	End   int     `json:"end"`
	Loc   *SrcLoc `json:"loc"`
}

type TSUnknownKeyword struct {
	Type  string  `json:"type"`
	Start int     `json:"start"`
	End   int     `json:"end"`
	Loc   *SrcLoc `json:"loc"`
}

type TSFunctionExpression struct {
	Type           string  `json:"type"`
	Start          int     `json:"start"`
	End            int     `json:"end"`
	Loc            *SrcLoc `json:"loc"`
	Id             Node    `json:"id"`
	Params         []Node  `json:"params"`
	Body           Node    `json:"body"`
	Generator      bool    `json:"generator"`
	Async          bool    `json:"async"`
	Expression     bool    `json:"expression"`
	TypeParameters Node    `json:"typeParameters"`
	ReturnType     Node    `json:"returnType"`
}

type TSFunctionDeclaration struct {
	Type           string  `json:"type"`
	Start          int     `json:"start"`
	End            int     `json:"end"`
	Loc            *SrcLoc `json:"loc"`
	Id             Node    `json:"id"`
	Params         []Node  `json:"params"`
	Body           Node    `json:"body"`
	Generator      bool    `json:"generator"`
	Async          bool    `json:"async"`
	TypeParameters Node    `json:"typeParameters"`
	ReturnType     Node    `json:"returnType"`
}

type TSArrowFunctionExpression struct {
	Type           string      `json:"type"`
	Start          int         `json:"start"`
	End            int         `json:"end"`
	Loc            *SrcLoc     `json:"loc"`
	Id             *Identifier `json:"id"`
	Params         []Node      `json:"params"`
	Body           Node        `json:"body"`
	Generator      bool        `json:"generator"`
	Async          bool        `json:"async"`
	Expression     bool        `json:"expression"`
	TypeParameters Node        `json:"typeParameters"`
	ReturnType     Node        `json:"returnType"`
}

type TSTypeReference struct {
	Type           string  `json:"type"`
	Start          int     `json:"start"`
	End            int     `json:"end"`
	Loc            *SrcLoc `json:"loc"`
	TypeName       Node    `json:"typeName"`
	TypeParameters Node    `json:"typeParameters"`
}

type TSTypeParameterDeclaration struct {
	Type   string  `json:"type"`
	Start  int     `json:"start"`
	End    int     `json:"end"`
	Loc    *SrcLoc `json:"loc"`
	Params []Node  `json:"params"`
}

type TSTypeParameterInstantiation struct {
	Type   string  `json:"type"`
	Start  int     `json:"start"`
	End    int     `json:"end"`
	Loc    *SrcLoc `json:"loc"`
	Params []Node  `json:"params"`
}

type TSTypeParameter struct {
	Type       string  `json:"type"`
	Start      int     `json:"start"`
	End        int     `json:"end"`
	Loc        *SrcLoc `json:"loc"`
	Name       Node    `json:"name"`
	Constraint Node    `json:"constraint"`
	Default    Node    `json:"default"`
}

type TSCallExpression struct {
	Type           string       `json:"type"`
	Start          int          `json:"start"`
	End            int          `json:"end"`
	Loc            *SrcLoc      `json:"loc"`
	Callee         Expression   `json:"callee"`
	Arguments      []Expression `json:"arguments"`
	Optional       bool         `json:"optional"`
	TypeParameters Node         `json:"typeParameters"`
}

type TSRestElement struct {
	Type           string  `json:"type"`
	Start          int     `json:"start"`
	End            int     `json:"end"`
	Loc            *SrcLoc `json:"loc"`
	Argument       Pattern `json:"argument"`
	Optional       bool    `json:"optional"`
	TypeAnnotation Node    `json:"typeAnnotation"`
}

type TSArrayType struct {
	Type        string  `json:"type"`
	Start       int     `json:"start"`
	End         int     `json:"end"`
	Loc         *SrcLoc `json:"loc"`
	ElementType Node    `json:"elementType"`
}

type TSObjectPattern struct {
	Type           string  `json:"type"`
	Start          int     `json:"start"`
	End            int     `json:"end"`
	Loc            *SrcLoc `json:"loc"`
	Properties     []Node  `json:"properties"`
	Optional       bool    `json:"optional"`
	TypeAnnotation Node    `json:"typeAnnotation"`
}

type TSTypePredicate struct {
	Type           string  `json:"type"`
	Start          int     `json:"start"`
	End            int     `json:"end"`
	Loc            *SrcLoc `json:"loc"`
	ParameterName  Node    `json:"parameterName"`
	TypeAnnotation Node    `json:"typeAnnotation"`
	Asserts        bool    `json:"asserts"`
}

type TSDeclareFunction struct {
	Type           string  `json:"type"`
	Start          int     `json:"start"`
	End            int     `json:"end"`
	Loc            *SrcLoc `json:"loc"`
	Id             Node    `json:"id"`
	Params         []Node  `json:"params"`
	Body           Node    `json:"body"`
	Generator      bool    `json:"generator"`
	Async          bool    `json:"async"`
	TypeParameters Node    `json:"typeParameters"`
	ReturnType     Node    `json:"returnType"`
}

type TSMethodDefinition struct {
	Type           string     `json:"type"`
	Start          int        `json:"start"`
	End            int        `json:"end"`
	Loc            *SrcLoc    `json:"loc"`
	Key            Expression `json:"key"`
	Value          Expression `json:"value"`
	Kind           string     `json:"kind"` // "constructor" | "method" | "get" | "set"
	Computed       bool       `json:"computed"`
	Static         bool       `json:"static"`
	TypeParameters Node       `json:"typeParameters"`
	ReturnType     Node       `json:"returnType"`
	Optional       bool       `json:"optional"`
	Abstract       bool       `json:"abstract"`
}

type TSParameterProperty struct {
	Type      string  `json:"type"`
	Start     int     `json:"start"`
	End       int     `json:"end"`
	Loc       *SrcLoc `json:"loc"`
	Parameter Node    `json:"parameter"`
	Readonly  bool    `json:"readonly"`
	Override  bool    `json:"override"`
}

type TSPropertyDefinition struct {
	Type           string     `json:"type"`
	Start          int        `json:"start"`
	End            int        `json:"end"`
	Loc            *SrcLoc    `json:"loc"`
	Key            Expression `json:"key"`
	Value          Expression `json:"value"`
	Computed       bool       `json:"computed"`
	Static         bool       `json:"static"`
	Abstract       bool       `json:"abstract"`
	Optional       bool       `json:"optional"`
	Declare        bool       `json:"declare"`
	Accessibility  string     `json:"accessibility"`
	TypeAnnotation Node       `json:"typeAnnotation"`
}

type TSIndexSignature struct {
	Type           string  `json:"type"`
	Start          int     `json:"start"`
	End            int     `json:"end"`
	Loc            *SrcLoc `json:"loc"`
	Parameters     []Node  `json:"parameters"`
	TypeAnnotation Node    `json:"typeAnnotation"`
}

type TSAsExpression struct {
	Type           string  `json:"type"`
	Start          int     `json:"start"`
	End            int     `json:"end"`
	Loc            *SrcLoc `json:"loc"`
	Expression     Node    `json:"expression"`
	TypeAnnotation Node    `json:"typeAnnotation"`
}

type TSTypeAssertion struct {
	Type           string  `json:"type"`
	Start          int     `json:"start"`
	End            int     `json:"end"`
	Loc            *SrcLoc `json:"loc"`
	Expression     Node    `json:"expression"`
	TypeAnnotation Node    `json:"typeAnnotation"`
}

type TSNonNullExpression struct {
	Type       string  `json:"type"`
	Start      int     `json:"start"`
	End        int     `json:"end"`
	Loc        *SrcLoc `json:"loc"`
	Expression Node    `json:"expression"`
}

type TSUnionType struct {
	Type  string  `json:"type"`
	Start int     `json:"start"`
	End   int     `json:"end"`
	Loc   *SrcLoc `json:"loc"`
	Types []Node  `json:"types"`
}

type TSIntersectionType struct {
	Type  string  `json:"type"`
	Start int     `json:"start"`
	End   int     `json:"end"`
	Loc   *SrcLoc `json:"loc"`
	Types []Node  `json:"types"`
}

type TSClassDeclaration struct {
	Type       string     `json:"type"`
	Start      int        `json:"start"`
	End        int        `json:"end"`
	Loc        *SrcLoc    `json:"loc"`
	Id         Expression `json:"id"`
	SuperClass Expression `json:"superClass"`
	Body       Expression `json:"body"`
	Declare    bool       `json:"declare"`
	Abstract   bool       `json:"abstract"`
}
