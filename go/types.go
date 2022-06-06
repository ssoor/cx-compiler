package main

type typedecl struct {
	Name   string       `json:"name,omitempty"`
	Typ    typedecltype `json:"type,omitempty"` // 类型
	Ref    typeref      `json:"ref,omitempty"`
	Fields []vardecl    `json:"field,omitempty"`
}

type typeref []string

// Storage? Typ+ Name([Arr])? (= Value)?
type vardecl struct {
	Name    string      `json:"name,omitempty"`
	Typ     vardecltype `json:"type,omitempty"`
	RefType typeref     `json:"tref,omitempty"`
	RefFunc function    `json:"fref,omitempty"`
	Arr     vararrdecl  `json:"arr,omitempty"`
	Value   interface{} `json:"val,omitempty"` // expression
	Storage string      `json:"storage,omitempty"`
}
type vararrdecl struct {
	Arr   bool   `json:"arr,omitempty"`   // 是否是数组
	Count string `json:"count,omitempty"` // 数组长度, -1 代表不定长
}

type nameref struct {
	Name string   `json:"name,omitempty"`
	Pkgs []string `json:"pkgs,omitempty"`
}

type varref struct {
	Name   nameref       `json:"name,omitempty"`
	Ops    []string      `json:"ops,omitempty"`
	Fields []varrefField `json:"field,omitempty"`
}

type varrefField struct {
	Ptr  bool    `json:"ptr,omitempty"`
	Name nameref `json:"name,omitempty"`
}

type funcref struct {
	Name   nameref      `json:"name,omitempty"`
	Typ    functiontype `json:"type,omitempty"` // 类型
	Retval typeref      `json:"ret,omitempty"`
	Params []vardecl    `json:"params,omitempty"`
}

type function struct {
	Typ   funcref     `json:"type,omitempty"` // 类型
	Block []statement `json:"block,omitempty"`
}

type statement struct {
	Typ  stmttype    `json:"type,omitempty"`
	Stmt interface{} `json:"stmt,omitempty"`
}

type returnstmt struct {
	Expr interface{} `json:"expr,omitempty"`
}

type expression struct {
	Typ  exprtype    `json:"type,omitempty"`
	Expr interface{} `json:"expr,omitempty"`
}

type opexpr struct {
	Op    string     `json:"op,omitempty"`
	L     expression `json:"L,omitempty"`
	R     expression `json:"R,omitempty"`
	There expression `json:"there,omitempty"`
}

type callexpr struct {
	Name   nameref      `json:"name,omitempty"`
	Params []expression `json:"params,omitempty"`
}

type vardecltype int

func (m vardecltype) MarshalJSON() ([]byte, error) {
	return []byte("\"" + m.String() + "\""), nil
}

func (m vardecltype) String() string {
	return map[vardecltype]string{
		varType: "type",
		varFunc: "func",
	}[m]
}

const (
	varUnknown vardecltype = iota
	varType                // 结构变量
	varFunc                // 方法变量
)

type functiontype int

func (m functiontype) MarshalJSON() ([]byte, error) {
	return []byte("\"" + m.String() + "\""), nil
}

func (m functiontype) String() string {
	return map[functiontype]string{
		functionNormal: "function",
		functionSelf:   "func_self",
	}[m]
}

const (
	functionNormal functiontype = iota // 常规方法
	functionSelf                       // 成员方法
)

type typedecltype int

func (m typedecltype) MarshalJSON() ([]byte, error) {
	return []byte("\"" + m.String() + "\""), nil
}

func (m typedecltype) String() string {
	return map[typedecltype]string{
		typedeclStruct: "struct",
		typedeclEnum:   "enum",
		typedeclUnion:  "union",
	}[m]
}

const (
	typedeclRef    typedecltype = iota // 引用
	typedeclStruct                     // 结构
	typedeclEnum                       // 枚举
	typedeclUnion                      // 联合
)

type stmttype int

func (m stmttype) MarshalJSON() ([]byte, error) {
	return []byte("\"" + m.String() + "\""), nil
}

func (m stmttype) String() string {
	return map[stmttype]string{
		stmtUnknown:  "unknown",
		stmtIf:       "if",
		stmtFor:      "for",
		stmtSwitch:   "switch",
		stmtReturn:   "return",
		stmtExpr:     "stmt_expr",
		stmtVarDecl:  "vardecl",
		stmtTypeDecl: "typedecl",
	}[m]
}

const (
	stmtUnknown stmttype = iota
	stmtIf
	stmtFor
	stmtWhile
	stmtSwitch
	stmtReturn
	stmtExpr
	stmtVarDecl
	stmtTypeDecl
)

type exprtype int

func (m exprtype) MarshalJSON() ([]byte, error) {
	return []byte("\"" + m.String() + "\""), nil
}

func (m exprtype) String() string {
	return map[exprtype]string{
		exprUnknown:    "unknown",
		exprOp:         "op",
		exprVar:        "var",
		exprCall:       "call",
		exprConstant:   "const",
		exprParenthese: "parent",
	}[m]
}

const (
	exprUnknown    exprtype = iota // 常规方法
	exprOp                         // 计算
	exprVar                        // 变量引用
	exprCall                       // 调用
	exprConstant                   // 常量
	exprParenthese                 // 圆括号
)
