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
	Decorators     []Node  `json:"decorators"`
}

type TSNumberKeyword struct {
	Type  string  `json:"type"`
	Start int     `json:"start"`
	End   int     `json:"end"`
	Loc   *SrcLoc `json:"loc"`
}

type TSObjectKeyword struct {
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

type TSIntrinsicKeyword struct {
	Type  string  `json:"type"`
	Start int     `json:"start"`
	End   int     `json:"end"`
	Loc   *SrcLoc `json:"loc"`
}

type TSNeverKeyword struct {
	Type  string  `json:"type"`
	Start int     `json:"start"`
	End   int     `json:"end"`
	Loc   *SrcLoc `json:"loc"`
}

type TSSymbolKeyword struct {
	Type  string  `json:"type"`
	Start int     `json:"start"`
	End   int     `json:"end"`
	Loc   *SrcLoc `json:"loc"`
}

type TSUndefinedKeyword struct {
	Type  string  `json:"type"`
	Start int     `json:"start"`
	End   int     `json:"end"`
	Loc   *SrcLoc `json:"loc"`
}

type TSBigIntKeyword struct {
	Type  string  `json:"type"`
	Start int     `json:"start"`
	End   int     `json:"end"`
	Loc   *SrcLoc `json:"loc"`
}

type TSNullKeyword struct {
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

type TSCallSignatureDeclaration struct {
	Type           string  `json:"type"`
	Start          int     `json:"start"`
	End            int     `json:"end"`
	Loc            *SrcLoc `json:"loc"`
	Id             Node    `json:"id"`
	Params         []Node  `json:"params"`
	Generator      bool    `json:"generator"`
	Async          bool    `json:"async"`
	TypeParameters Node    `json:"typeParameters"`
	ReturnType     Node    `json:"returnType"`
}

type TSConstructSignatureDeclaration struct {
	Type           string  `json:"type"`
	Start          int     `json:"start"`
	End            int     `json:"end"`
	Loc            *SrcLoc `json:"loc"`
	Id             Node    `json:"id"`
	Params         []Node  `json:"params"`
	Generator      bool    `json:"generator"`
	Async          bool    `json:"async"`
	TypeParameters Node    `json:"typeParameters"`
	ReturnType     Node    `json:"returnType"`
	Abstract       bool    `json:"abstract"`
}

type TSConstructorType struct {
	Type           string  `json:"type"`
	Start          int     `json:"start"`
	End            int     `json:"end"`
	Loc            *SrcLoc `json:"loc"`
	Id             Node    `json:"id"`
	Params         []Node  `json:"params"`
	Generator      bool    `json:"generator"`
	Async          bool    `json:"async"`
	TypeParameters Node    `json:"typeParameters"`
	ReturnType     Node    `json:"returnType"`
	Abstract       bool    `json:"abstract"`
}

type TSFunctionType struct {
	Type           string  `json:"type"`
	Start          int     `json:"start"`
	End            int     `json:"end"`
	Loc            *SrcLoc `json:"loc"`
	Id             Node    `json:"id"`
	Params         []Node  `json:"params"`
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

type TSNewExpression struct {
	Type           string       `json:"type"`
	Start          int          `json:"start"`
	End            int          `json:"end"`
	Loc            *SrcLoc      `json:"loc"`
	Callee         Expression   `json:"callee"`
	Arguments      []Expression `json:"arguments"`
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

type TSTypeLiteral struct {
	Type    string  `json:"type"`
	Start   int     `json:"start"`
	End     int     `json:"end"`
	Loc     *SrcLoc `json:"loc"`
	Members Node    `json:"members"`
}

// used as the member of `TSTypeLiteral`
type TSPropertySignature struct {
	Type           string  `json:"type"`
	Start          int     `json:"start"`
	End            int     `json:"end"`
	Loc            *SrcLoc `json:"loc"`
	Key            Node    `json:"key"`
	Computed       bool    `json:"computed"`
	Optional       bool    `json:"optional"`
	TypeAnnotation Node    `json:"typeAnnotation"`
	Kind           string  `json:"kind"`
	Readonly       bool    `json:"readonly"`
}

type TSMethodSignature struct {
	Type     string     `json:"type"`
	Start    int        `json:"start"`
	End      int        `json:"end"`
	Loc      *SrcLoc    `json:"loc"`
	Key      Node       `json:"key"`
	Value    Expression `json:"value"`
	Computed bool       `json:"computed"`
	Optional bool       `json:"optional"`
	Kind     string     `json:"kind"`
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
	Type          string     `json:"type"`
	Start         int        `json:"start"`
	End           int        `json:"end"`
	Loc           *SrcLoc    `json:"loc"`
	Key           Expression `json:"key"`
	Value         Expression `json:"value"`
	Kind          string     `json:"kind"` // "constructor" | "method" | "get" | "set"
	Computed      bool       `json:"computed"`
	Static        bool       `json:"static"`
	Optional      bool       `json:"optional"`
	Definite      bool       `json:"definite"`
	Override      bool       `json:"override"`
	Abstract      bool       `json:"abstract"`
	Readonly      bool       `json:"readonly"`
	Accessibility string     `json:"accessibility"`
	Decorators    []Node     `json:"decorators"`
}

// represets the properties defined via constructor params
type TSParameterProperty struct {
	Type          string  `json:"type"`
	Start         int     `json:"start"`
	End           int     `json:"end"`
	Loc           *SrcLoc `json:"loc"`
	Parameter     Node    `json:"parameter"`
	Readonly      bool    `json:"readonly"`
	Accessibility string  `json:"accessibility"`
	Override      bool    `json:"override"`
	Decorators    []Node  `json:"decorators"`
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
	Definite       bool       `json:"definite"`
	Readonly       bool       `json:"readonly"`
	Override       bool       `json:"override"`
	Declare        bool       `json:"declare"`
	Accessibility  string     `json:"accessibility"`
	TypeAnnotation Node       `json:"typeAnnotation"`
	Decorators     []Node     `json:"decorators"`
}

type TSIndexSignature struct {
	Type           string  `json:"type"`
	Start          int     `json:"start"`
	End            int     `json:"end"`
	Loc            *SrcLoc `json:"loc"`
	Static         bool    `json:"static"`
	Abstract       bool    `json:"abstract"`
	Optional       bool    `json:"optional"`
	Readonly       bool    `json:"readonly"`
	Declare        bool    `json:"declare"`
	Accessibility  string  `json:"accessibility"`
	Parameters     []Node  `json:"parameters"`
	TypeAnnotation Node    `json:"typeAnnotation"`
	Decorators     []Node  `json:"decorators"`
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
	Type                string     `json:"type"`
	Start               int        `json:"start"`
	End                 int        `json:"end"`
	Loc                 *SrcLoc    `json:"loc"`
	Id                  Expression `json:"id"`
	TypeParameters      Node       `json:"typeParameters"`
	SuperClass          Expression `json:"superClass"`
	SuperTypeParameters Node       `json:"superTypeParameters"`
	Implements          []Node     `json:"implements"`
	Body                Expression `json:"body"`
	Declare             bool       `json:"declare"`
	Decorators          []Node     `json:"decorators"`
	Abstract            bool       `json:"abstract"`
}

type TSClassExpression struct {
	Type                string     `json:"type"`
	Start               int        `json:"start"`
	End                 int        `json:"end"`
	Loc                 *SrcLoc    `json:"loc"`
	Id                  Expression `json:"id"`
	TypeParameters      Node       `json:"typeParameters"`
	SuperClass          Expression `json:"superClass"`
	SuperTypeParameters Node       `json:"superTypeParameters"`
	Implements          []Node     `json:"implements"`
	Body                Expression `json:"body"`
	Decorators          []Node     `json:"decorators"`
	Abstract            bool       `json:"abstract"`
}

type TSQualifiedName struct {
	Type  string     `json:"type"`
	Start int        `json:"start"`
	End   int        `json:"end"`
	Loc   *SrcLoc    `json:"loc"`
	Left  Expression `json:"left"`
	Right Expression `json:"right"`
}

type TSVariableDeclaration struct {
	Type         string                `json:"type"`
	Start        int                   `json:"start"`
	End          int                   `json:"end"`
	Loc          *SrcLoc               `json:"loc"`
	Kind         string                `json:"kind"`
	Declarations []*VariableDeclarator `json:"declarations"`
	Declare      bool                  `json:"declare"`
}

type TSInterfaceDeclaration struct {
	Type           string     `json:"type"`
	Start          int        `json:"start"`
	End            int        `json:"end"`
	Loc            *SrcLoc    `json:"loc"`
	Id             Expression `json:"id"`
	TypeParameters Node       `json:"typeParameters"`
	Extends        []Node     `json:"extends"`
	Body           Expression `json:"body"`
	Declare        bool       `json:"declare"`
}

type TSInterfaceBody struct {
	Type  string  `json:"type"`
	Start int     `json:"start"`
	End   int     `json:"end"`
	Loc   *SrcLoc `json:"loc"`
	Body  []Node  `json:"body"`
}

type TSExpressionWithTypeArguments struct {
	Type           string  `json:"type"`
	Start          int     `json:"start"`
	End            int     `json:"end"`
	Loc            *SrcLoc `json:"loc"`
	Expression     Node    `json:"expression"`
	TypeParameters Node    `json:"typeParameters"`
}

type TSArrayPattern struct {
	Type           string  `json:"type"`
	Start          int     `json:"start"`
	End            int     `json:"end"`
	Loc            *SrcLoc `json:"loc"`
	Elements       []Node  `json:"elements"`
	Optional       bool    `json:"optional"`
	TypeAnnotation Node    `json:"typeAnnotation"`
}

type TSEnumDeclaration struct {
	Type    string     `json:"type"`
	Start   int        `json:"start"`
	End     int        `json:"end"`
	Loc     *SrcLoc    `json:"loc"`
	Id      Expression `json:"id"`
	Members []Node     `json:"members"`
	Const   bool       `json:"const"`
	Declare bool       `json:"declare"`
}

type TSEnumMember struct {
	Type        string     `json:"type"`
	Start       int        `json:"start"`
	End         int        `json:"end"`
	Loc         *SrcLoc    `json:"loc"`
	Id          Expression `json:"id"`
	Initializer Node       `json:"initializer"`
}

type TSTypeAliasDeclaration struct {
	Type           string     `json:"type"`
	Start          int        `json:"start"`
	End            int        `json:"end"`
	Loc            *SrcLoc    `json:"loc"`
	Id             Expression `json:"id"`
	TypeParameters Node       `json:"typeParameters"`
	TypeAnnotation Node       `json:"typeAnnotation"`
	Declare        bool       `json:"declare"`
}

type TSModuleDeclaration struct {
	Type    string     `json:"type"`
	Start   int        `json:"start"`
	End     int        `json:"end"`
	Loc     *SrcLoc    `json:"loc"`
	Id      Expression `json:"id"`
	Body    Node       `json:"body"`
	Declare bool       `json:"declare"`
	Global  bool       `json:"global"`
}

type TSNamespaceExportDeclaration struct {
	Type  string     `json:"type"`
	Start int        `json:"start"`
	End   int        `json:"end"`
	Loc   *SrcLoc    `json:"loc"`
	Id    Expression `json:"id"`
}

type TSExportAssignment struct {
	Type       string     `json:"type"`
	Start      int        `json:"start"`
	End        int        `json:"end"`
	Loc        *SrcLoc    `json:"loc"`
	Expression Expression `json:"expression"`
}

type TSLiteralType struct {
	Type    string  `json:"type"`
	Start   int     `json:"start"`
	End     int     `json:"end"`
	Loc     *SrcLoc `json:"loc"`
	Literal Node    `json:"literal"`
}

type TSImportEqualsDeclaration struct {
	Type            string  `json:"type"`
	Start           int     `json:"start"`
	End             int     `json:"end"`
	Loc             *SrcLoc `json:"loc"`
	Id              Node    `json:"id"`
	ModuleReference Node    `json:"moduleReference"`
	IsExport        bool    `json:"isExport"`
}

type TSExternalModuleReference struct {
	Type       string  `json:"type"`
	Start      int     `json:"start"`
	End        int     `json:"end"`
	Loc        *SrcLoc `json:"loc"`
	Expression Node    `json:"expression"`
}

type TSTaggedTemplateExpression struct {
	Type           string     `json:"type"`
	Start          int        `json:"start"`
	End            int        `json:"end"`
	Loc            *SrcLoc    `json:"loc"`
	Tag            Expression `json:"tag"`
	Quasi          Expression `json:"quasi"`
	TypeParameters Node       `json:"typeParameters"`
}

type TSXOpeningElement struct {
	Type           string  `json:"type"`
	Start          int     `json:"start"`
	End            int     `json:"end"`
	Loc            *SrcLoc `json:"loc"`
	Name           Node    `json:"name"`
	Attributes     []Node  `json:"attributes"`
	SelfClosing    bool    `json:"selfClosing"`
	TypeParameters Node    `json:"typeParameters"`
}

type TSImportType struct {
	Type           string  `json:"type"`
	Start          int     `json:"start"`
	End            int     `json:"end"`
	Loc            *SrcLoc `json:"loc"`
	Argument       Node    `json:"argument"`
	Qualifier      Node    `json:"qualifier"`
	TypeParameters Node    `json:"typeParameters"`
}

type TSTypeQuery struct {
	Type     string  `json:"type"`
	Start    int     `json:"start"`
	End      int     `json:"end"`
	Loc      *SrcLoc `json:"loc"`
	ExprName Node    `json:"exprName"`
}

type TSConditionalType struct {
	Type        string  `json:"type"`
	Start       int     `json:"start"`
	End         int     `json:"end"`
	Loc         *SrcLoc `json:"loc"`
	CheckType   Node    `json:"checkType"`
	ExtendsType Node    `json:"extendsType"`
	TrueType    Node    `json:"trueType"`
	FalseType   Node    `json:"falseType"`
}

type TSInferType struct {
	Type          string  `json:"type"`
	Start         int     `json:"start"`
	End           int     `json:"end"`
	Loc           *SrcLoc `json:"loc"`
	TypeParameter Node    `json:"typeParameter"`
}

type TSParenthesizedType struct {
	Type           string  `json:"type"`
	Start          int     `json:"start"`
	End            int     `json:"end"`
	Loc            *SrcLoc `json:"loc"`
	TypeAnnotation Node    `json:"typeAnnotation"`
}

type TSIndexedAccessType struct {
	Type       string  `json:"type"`
	Start      int     `json:"start"`
	End        int     `json:"end"`
	Loc        *SrcLoc `json:"loc"`
	ObjectType Node    `json:"objectType"`
	IndexType  Node    `json:"indexType"`
}

type TSMappedType struct {
	Type           string      `json:"type"`
	Start          int         `json:"start"`
	End            int         `json:"end"`
	Loc            *SrcLoc     `json:"loc"`
	Readonly       interface{} `json:"readonly"`
	Optional       interface{} `json:"optional"`
	TypeParameter  Node        `json:"typeParameter"`
	NameType       Node        `json:"nameType"`
	TypeAnnotation Node        `json:"typeAnnotation"`
}

type TSTypeOperator struct {
	Type           string      `json:"type"`
	Start          int         `json:"start"`
	End            int         `json:"end"`
	Loc            *SrcLoc     `json:"loc"`
	Operator       interface{} `json:"operator"`
	TypeAnnotation Node        `json:"typeAnnotation"`
}

type TSTupleType struct {
	Type         string  `json:"type"`
	Start        int     `json:"start"`
	End          int     `json:"end"`
	Loc          *SrcLoc `json:"loc"`
	ElementTypes []Node  `json:"elementTypes"`
}

type TSRestType struct {
	Type           string  `json:"type"`
	Start          int     `json:"start"`
	End            int     `json:"end"`
	Loc            *SrcLoc `json:"loc"`
	TypeAnnotation Node    `json:"typeAnnotation"`
}

type TSNamedTupleMember struct {
	Type        string  `json:"type"`
	Start       int     `json:"start"`
	End         int     `json:"end"`
	Loc         *SrcLoc `json:"loc"`
	Optional    bool    `json:"optional"`
	Label       Node    `json:"label"`
	ElementType Node    `json:"elementType"`
}

type TSOptionalType struct {
	Type           string  `json:"type"`
	Start          int     `json:"start"`
	End            int     `json:"end"`
	Loc            *SrcLoc `json:"loc"`
	TypeAnnotation Node    `json:"typeAnnotation"`
}
