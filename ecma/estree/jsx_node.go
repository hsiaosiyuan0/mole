package estree

// grammar: https://github.com/facebook/jsx
// ast: https://github.com/facebook/jsx/blob/main/AST.md

type JSXElement struct {
	Type           string  `json:"type"`
	Start          int     `json:"start"`
	End            int     `json:"end"`
	Loc            *SrcLoc `json:"loc"`
	OpeningElement Node    `json:"openingElement"`
	Children       []Node  `json:"children"`
	ClosingElement Node    `json:"closingElement"`
}

type JSXOpeningElement struct {
	Type        string  `json:"type"`
	Start       int     `json:"start"`
	End         int     `json:"end"`
	Loc         *SrcLoc `json:"loc"`
	Name        Node    `json:"name"`
	Attributes  []Node  `json:"attributes"`
	SelfClosing bool    `json:"selfClosing"`
}

type JSXIdentifier struct {
	Type  string  `json:"type"`
	Start int     `json:"start"`
	End   int     `json:"end"`
	Loc   *SrcLoc `json:"loc"`
	Name  string  `json:"name"`
}

type JSXNamespacedName struct {
	Type      string  `json:"type"`
	Start     int     `json:"start"`
	End       int     `json:"end"`
	Loc       *SrcLoc `json:"loc"`
	Namespace string  `json:"namespace"`
	Name      string  `json:"name"`
}

type JSXMemberExpression struct {
	Type     string  `json:"type"`
	Start    int     `json:"start"`
	End      int     `json:"end"`
	Loc      *SrcLoc `json:"loc"`
	Object   Node    `json:"object"`
	Property Node    `json:"property"`
}

type JSXClosingElement struct {
	Type  string  `json:"type"`
	Start int     `json:"start"`
	End   int     `json:"end"`
	Loc   *SrcLoc `json:"loc"`
	Name  Node    `json:"name"`
}

type JSXFragment struct {
	Type            string  `json:"type"`
	Start           int     `json:"start"`
	End             int     `json:"end"`
	Loc             *SrcLoc `json:"loc"`
	OpeningFragment Node    `json:"openingFragment"`
	Children        []Node  `json:"children"`
	ClosingFragment Node    `json:"closingFragment"`
}

type JSXOpeningFragment struct {
	Type  string  `json:"type"`
	Start int     `json:"start"`
	End   int     `json:"end"`
	Loc   *SrcLoc `json:"loc"`
}

type JSXClosingFragment struct {
	Type  string  `json:"type"`
	Start int     `json:"start"`
	End   int     `json:"end"`
	Loc   *SrcLoc `json:"loc"`
}

type JSXText struct {
	Type  string  `json:"type"`
	Start int     `json:"start"`
	End   int     `json:"end"`
	Loc   *SrcLoc `json:"loc"`
	Value string  `json:"value"`
	Raw   string  `json:"raw"`
}

type JSXExpressionContainer struct {
	Type       string  `json:"type"`
	Start      int     `json:"start"`
	End        int     `json:"end"`
	Loc        *SrcLoc `json:"loc"`
	Expression Node    `json:"expression"`
}

type JSXSpreadAttribute struct {
	Type     string     `json:"type"`
	Start    int        `json:"start"`
	End      int        `json:"end"`
	Loc      *SrcLoc    `json:"loc"`
	Argument Expression `json:"argument"`
}

type JSXSpreadChild struct {
	Type       string  `json:"type"`
	Start      int     `json:"start"`
	End        int     `json:"end"`
	Loc        *SrcLoc `json:"loc"`
	Expression Node    `json:"expression"`
}

type JSXAttribute struct {
	Type  string  `json:"type"`
	Start int     `json:"start"`
	End   int     `json:"end"`
	Loc   *SrcLoc `json:"loc"`
	Name  Node    `json:"name"`
	Value Node    `json:"value"`
}

type JSXEmptyExpression struct {
	Type  string  `json:"type"`
	Start int     `json:"start"`
	End   int     `json:"end"`
	Loc   *SrcLoc `json:"loc"`
}
