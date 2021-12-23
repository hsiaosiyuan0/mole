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

type TSFunctionExpression struct {
	Type       string  `json:"type"`
	Start      int     `json:"start"`
	End        int     `json:"end"`
	Loc        *SrcLoc `json:"loc"`
	Id         Node    `json:"id"`
	Params     []Node  `json:"params"`
	Body       Node    `json:"body"`
	Generator  bool    `json:"generator"`
	Async      bool    `json:"async"`
	Expression bool    `json:"expression"`
	ReturnType Node    `json:"returnType"`
}

type TSFunctionDeclaration struct {
	Type       string  `json:"type"`
	Start      int     `json:"start"`
	End        int     `json:"end"`
	Loc        *SrcLoc `json:"loc"`
	Id         Node    `json:"id"`
	Params     []Node  `json:"params"`
	Body       Node    `json:"body"`
	Generator  bool    `json:"generator"`
	Async      bool    `json:"async"`
	ReturnType Node    `json:"returnType"`
}

type TSArrowFunctionExpression struct {
	Type       string      `json:"type"`
	Start      int         `json:"start"`
	End        int         `json:"end"`
	Loc        *SrcLoc     `json:"loc"`
	Id         *Identifier `json:"id"`
	Params     []Node      `json:"params"`
	Body       Node        `json:"body"`
	Generator  bool        `json:"generator"`
	Async      bool        `json:"async"`
	Expression bool        `json:"expression"`
	ReturnType Node        `json:"returnType"`
}

type TSTypeReference struct {
	Type           string  `json:"type"`
	Start          int     `json:"start"`
	End            int     `json:"end"`
	Loc            *SrcLoc `json:"loc"`
	TypeName       Node    `json:"typeName"`
	TypeParameters Node    `json:"typeParameters"`
}

type TSTypeParameterInstantiation struct {
	Type   string  `json:"type"`
	Start  int     `json:"start"`
	End    int     `json:"end"`
	Loc    *SrcLoc `json:"loc"`
	Params []Node  `json:"params"`
}
