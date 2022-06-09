package main

import (
	"fmt"
	"reflect"
	"strings"
)

type nameref struct {
	Name string   `json:"name,omitempty"`
	Pkgs []string `json:"pkgs,omitempty"`
}

func (m nameref) MarshalJSON() ([]byte, error) {
	return []byte("\"" + m.String() + "\""), nil
}

func (m nameref) String() string {
	pkgs := append([]string{}, m.Pkgs...)
	return strings.Join(append(pkgs, m.Name), "_")
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

type stmttype int

func (m stmttype) MarshalJSON() ([]byte, error) {
	return []byte("\"" + m.String() + "\""), nil
}

func (m stmttype) String() string {
	return map[stmttype]string{
		stmtUnknown:  "unknown",
		stmtIf:       "if",
		stmtFor:      "for",
		stmtCase:     "case",
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
	stmtCase
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
