package main

import (
	"fmt"
	"reflect"
	"strings"
)

type typedefine struct {
	Name string      `json:"name,omitempty"`
	Typ  interface{} `json:"type,omitempty"`
}

func (m typedefine) MarshalJSON() ([]byte, error) {
	return []byte("\"" + m.String() + "\""), nil
}

func (m typedefine) String() string {
	return fmt.Sprintf("typedef %s %s", infa2str(m.Typ), m.Name)
}

type enumdecl struct {
	Name   string       `json:"name,omitempty"`
	Fields []vardeclval `json:"field,omitempty"`
}

func (m enumdecl) MarshalJSON() ([]byte, error) {
	return []byte("\"" + m.String() + "\""), nil
}

func (m enumdecl) String() string {
	msg := "enum " + m.Name + "{\n"
	for _, v := range m.Fields {
		msg += "\t" + v.Name
		if v.Value != nil {
			msg += " = " + infa2str(v.Value)
		}
		msg += ",\n"
	}

	return msg + "}"
}

type typedecl struct {
	Name   string       `json:"name,omitempty"`
	Typ    typedecltype `json:"type,omitempty"` // 类型
	Ref    []string     `json:"ref,omitempty"`
	Fields []vardecl    `json:"field,omitempty"`
}

func (m typedecl) MarshalJSON() ([]byte, error) {
	return []byte("\"" + m.String() + "\""), nil
}

func (m typedecl) String() string {
	switch m.Typ {
	case typedeclRef:
		msg := ""
		for _, v := range m.Ref {
			switch v {
			case "*":
				msg += "*"
			default:
				msg += " " + v
			}
		}
		msg += " " + m.Name

		return strings.TrimSuffix(strings.TrimPrefix(msg, " "), " ")
	case typedeclStruct:
		msg := "struct " + m.Name + "{\n"
		for _, v := range m.Fields {
			msg += "\t" + v.String() + ";\n"
		}

		return msg + "}"
	case typedeclUnion:
		msg := "union " + m.Name + "{\n"
		for _, v := range m.Fields {
			msg += "\t" + v.String() + ";\n"
		}

		return msg + "}"
	}

	return fmt.Sprintf("!panic(%s)", reflect.TypeOf(m).String())
}

// Storage? Typ+ Name([Arr])? (= Value)?
type vardecl struct {
	Name    string       `json:"name,omitempty"`
	Typ     vardecltype  `json:"type,omitempty"`
	RefType typedecl     `json:"tref,omitempty"`
	RefFunc function     `json:"fref,omitempty"`
	Storage string       `json:"storage,omitempty"`
	Values  []vardeclval `json:"vals,omitempty"`
}

func (m vardecl) MarshalJSON() ([]byte, error) {
	return []byte("\"" + m.String() + "\""), nil
}

func (m vardecl) String() string {
	msg := ""
	switch m.Typ {
	case varSelf:
		msg = m.RefType.String() + " " + m.Name
	case varFunc:
		// return m.RefFunc.String()
	case varType:
		valuesMsg := ""
		for _, v := range m.Values {
			valuesMsg += v.String() + ", "
		}
		valuesMsg = strings.TrimSuffix(valuesMsg, ", ")

		msg = m.RefType.String() + " " + valuesMsg
	default:
		msg = fmt.Sprintf("!panic(%s)", reflect.TypeOf(m).String())
	}

	if m.Storage != "" {
		msg = m.Storage + " " + msg
	}

	return msg
}

type vardeclval struct {
	Name  string       `json:"name,omitempty"`
	Arr   []vararrdecl `json:"arr,omitempty"`
	Value interface{}  `json:"val,omitempty"` // expression
}

func (m vardeclval) MarshalJSON() ([]byte, error) {
	return []byte("\"" + m.String() + "\""), nil
}

func (m vardeclval) String() string {
	msg := m.Name
	for _, v := range m.Arr {
		msg += v.String()
	}

	if m.Value != nil {
		msg += " = " + infa2str(m.Value)
	}

	return msg
}

type varblockdecl struct {
	Arr    bool               `json:"arr,omitempty"`
	Fields []varblocksubfield `json:"fields,omitempty"`
}

func (m varblockdecl) MarshalJSON() ([]byte, error) {
	return []byte("\"" + m.String() + "\""), nil
}

func (m varblockdecl) String() string {
	msg := "{"
	if m.Arr {
		msg = "["
	}

	for _, v := range m.Fields {
		msg += v.String() + ", "
	}
	msg = strings.TrimSuffix(msg, ", ")

	if m.Arr {
		msg += "]"
	} else {
		msg += "}"
	}

	return msg
}

type varblocksubfield struct {
	Name  string      `json:"name,omitempty"`
	Value interface{} `json:"val,omitempty"`
}

func (m varblocksubfield) MarshalJSON() ([]byte, error) {
	return []byte("\"" + m.String() + "\""), nil
}

func (m varblocksubfield) String() string {
	if m.Name != "" {
		return "." + m.Name + " = " + infa2str(m.Value)
	}

	return infa2str(m.Value)
}

type vararrdecl struct {
	Arr   bool   `json:"arr,omitempty"`   // 是否是数组
	Count string `json:"count,omitempty"` // 数组长度, -1 代表不定长
}

func (m vararrdecl) MarshalJSON() ([]byte, error) {
	return []byte("\"" + m.String() + "\""), nil
}

func (m vararrdecl) String() string {
	if !m.Arr {
		return ""
	}
	if m.Count == "-1" {
		return "[]"
	}

	return "[" + m.Count + "]"
}

type nameref struct {
	Name string   `json:"name,omitempty"`
	Pkgs []string `json:"pkgs,omitempty"`
}

func (m nameref) MarshalJSON() ([]byte, error) {
	return []byte("\"" + m.String() + "\""), nil
}

func (m nameref) String() string {
	pkgs := append([]string{}, m.Pkgs...)
	return strings.Join(append(pkgs, m.Name), ":")
}

type varref struct {
	Parent interface{}   `json:"parent,omitempty"`
	Ops    []string      `json:"ops,omitempty"`
	Fields []varrefField `json:"field,omitempty"`
}

func (m varref) MarshalJSON() ([]byte, error) {
	return []byte("\"" + m.String() + "\""), nil
}

func (m varref) String() string {
	msg := ""
	for _, v := range m.Ops {
		msg += v
	}

	msg += infa2str(m.Parent)

	for _, v := range m.Fields {
		msg += v.String()
	}

	return msg + ""
}

type varrefField struct {
	Ptr  bool    `json:"ptr,omitempty"`
	Name nameref `json:"name,omitempty"`
}

func (m varrefField) MarshalJSON() ([]byte, error) {
	return []byte("\"" + m.String() + "\""), nil
}

func (m varrefField) String() string {
	op := "."
	if m.Ptr {
		op = "->"
	}

	return fmt.Sprintf("%s%s", op, m.Name.String())
}

type funcref struct {
	Name   nameref      `json:"name,omitempty"`
	Typ    functiontype `json:"type,omitempty"` // 类型
	Retval typedecl     `json:"ret,omitempty"`
	Params []vardecl    `json:"params,omitempty"`
}

func (m funcref) MarshalJSON() ([]byte, error) {
	return []byte("\"" + m.String() + "\""), nil
}

func (m funcref) String() string {
	params := ""
	for _, v := range m.Params {
		params += v.String() + ", "
	}
	params = strings.TrimSuffix(params, ", ")

	return fmt.Sprintf("%s %s(%s)", m.Retval, m.Name.String(), params)
}

type function struct {
	Typ   funcref     `json:"type,omitempty"` // 类型
	Block []statement `json:"block,omitempty"`
}

func (m function) MarshalJSON() ([]byte, error) {
	return []byte("\"" + m.String() + "\""), nil
}

func (m function) String() string {
	msg := m.Typ.String() + " {\n"

	for _, v := range m.Block {
		msg += "\t" + v.String() + "\n"
	}

	msg += "}\n"

	return msg
}

type statement struct {
	Typ  stmttype    `json:"type,omitempty"`
	Stmt interface{} `json:"stmt,omitempty"`
}

func (m statement) MarshalJSON() ([]byte, error) {
	return []byte("\"" + m.String() + "\""), nil
}

func (m statement) String() string {
	if str, ok := m.Stmt.(string); ok {
		return str
	}

	strm, ok := m.Stmt.(fmt.Stringer)
	if !ok {
		if m.Stmt == nil {
			return "<nil>"
		}
		return fmt.Sprintf("!panic(%s)", reflect.TypeOf(m.Stmt).String())
	}

	return strm.String() + ";"
}

type returnstmt struct {
	Expr interface{} `json:"expr,omitempty"`
}

func (m returnstmt) MarshalJSON() ([]byte, error) {
	return []byte("\"" + m.String() + "\""), nil
}

func (m returnstmt) String() string {
	strm, ok := m.Expr.(fmt.Stringer)
	if !ok {
		if m.Expr == nil {
			return "<nil>"
		}
		return fmt.Sprintf("!panic(%s)", reflect.TypeOf(m.Expr).String())
	}

	return fmt.Sprintf("return %s", strm.String())
}

type expression struct {
	Typ  exprtype    `json:"type,omitempty"`
	Expr interface{} `json:"expr,omitempty"`
}

func (m expression) MarshalJSON() ([]byte, error) {
	return []byte("\"" + m.String() + "\""), nil
}

func (m expression) String() string {
	switch m.Typ {
	case exprOp:
		return infa2str(m.Expr)
	case exprVar:
		return infa2str(m.Expr)
	case exprVarBlock:
		return infa2str(m.Expr)
	case exprCall:
		return infa2str(m.Expr)
	case exprConstant:
		return infa2str(m.Expr)
	case exprParenthese:
		return fmt.Sprintf("(%s)", m.Expr.(expression).String())
	}

	if m.Expr == nil {
		return "<nil>"
	}
	return fmt.Sprintf("!panic(%s)", reflect.TypeOf(m.Expr).String())
}

type opexpr struct {
	Op    string     `json:"op,omitempty"`
	L     expression `json:"L,omitempty"`
	R     expression `json:"R,omitempty"`
	There expression `json:"there,omitempty"`
}

func (m opexpr) MarshalJSON() ([]byte, error) {
	return []byte("\"" + m.String() + "\""), nil
}

func (m opexpr) String() string {
	if m.Op == "?" {
		return fmt.Sprintf("%s ? %s : %s", m.L.String(), m.R.String(), m.There.String())
	}

	if m.L.Expr == nil {
		return fmt.Sprintf("%s%s", m.Op, m.R.String())
	}

	if m.R.Expr == nil {
		return fmt.Sprintf("%s%s", m.L.String(), m.Op)
	}

	return fmt.Sprintf("%s %s %s", m.L.String(), m.Op, m.R.String())
}

type callexpr struct {
	Var    varref       `json:"var,omitempty"`
	Params []expression `json:"params,omitempty"`
}

func (m callexpr) MarshalJSON() ([]byte, error) {
	return []byte("\"" + m.String() + "\""), nil
}

func (m callexpr) String() string {
	params := ""
	for _, v := range m.Params {
		params += v.String() + ", "
	}
	params = strings.TrimSuffix(params, ", ")
	return fmt.Sprintf("%s(%s)", m.Var.String(), params)
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
	varSelf                // 成员变量
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
		typedeclUnion:  "union",
	}[m]
}

const (
	typedeclUnknown typedecltype = iota
	typedeclRef                  // 引用
	typedeclStruct               // 结构
	typedeclUnion                // 联合
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
		stmtAssign:   "=",
		stmtVarDecl:  "vardecl",
		stmtEnumDecl: "enumdecl",
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
	stmtAssign
	stmtVarDecl
	stmtEnumDecl
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
	exprVarBlock                   // 块变量定义
	exprCall                       // 调用
	exprConstant                   // 常量
	exprParenthese                 // 圆括号
)

func infa2str(in interface{}) string {
	msg := ""

	if str, ok := in.(string); ok {
		msg = str
	} else if strm, ok := in.(fmt.Stringer); ok {
		msg = strm.String()
	} else if in == nil {
		msg = "<nil>"
	} else {
		msg = fmt.Sprintf("!panic(%s)", reflect.TypeOf(in).String())
	}

	return msg
}
