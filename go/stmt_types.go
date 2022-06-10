package main

import (
	"fmt"
	"reflect"
	"strings"
)

type codeblock struct {
	Body statement `json:"stmt,omitempty"`
}

func (m codeblock) String() string {
	msg := infa2str(m.Body)
	switch m.Body.Typ {
	case stmtBlock:
	default:
		msg = indentation + strings.ReplaceAll(msg, "\n", "\n"+indentation)
	}

	if count := strings.Count(msg, "\n"); count > 1 {
		msg = strings.Replace(msg, "\n", "\n"+indentation, count-1)
	}

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
	switch m.Typ {
	case stmtNone:
		return ";"
	case stmtBreak:
		msg := ""
		if m.Stmt != nil {
			msg += infa2str(m.Stmt) + "\n"
		}
		return msg + "break;"
	case stmtBlock:
		st := m.Stmt.(SymbolTable)
		return st.String() + st.Debug()
	case stmtVarDecl, stmtEnumDecl, stmtTypeDef, stmtTypeDecl, stmtExpr:
		return infa2str(m.Stmt) + ";"
	default:
		return infa2str(m.Stmt)
	}
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

	return fmt.Sprintf("return %s;", strm.String())
}

type ifstmt struct {
	Expr interface{} `json:"expr,omitempty"`
	Body codeblock   `json:"body,omitempty"`
	Else interface{} `json:"else,omitempty"`
}

func (m ifstmt) MarshalJSON() ([]byte, error) {
	return []byte("\"" + m.String() + "\""), nil
}

func (m ifstmt) String() string {
	msg := "if (" + infa2str(m.Expr) + ") " + m.Body.String()

	if m.Else != nil {
		elseBody := m.Else.(codeblock)

		if elseBody.Body.Typ == stmtIf {
			msg += " else " + elseBody.Body.String()
		} else {
			msg += " else " + infa2str(m.Else)
		}
	}

	return msg
}

type forstmt struct {
	Init interface{} `json:"expr,omitempty"`
	Expr expression  `json:"expr,omitempty"`
	Incr expression  `json:"expr,omitempty"`
	Body codeblock   `json:"body,omitempty"`
}

func (m forstmt) MarshalJSON() ([]byte, error) {
	return []byte("\"" + m.String() + "\""), nil
}

func (m forstmt) String() string {
	msg := "for (" + infa2str(m.Init) + ";" + m.Expr.String() + ";" + m.Incr.String() + ") " + m.Body.String()

	return msg
}

type whilestmt struct {
	Do   bool        `json:"def,omitempty"`
	Expr interface{} `json:"expr,omitempty"`
	Body codeblock   `json:"body,omitempty"`
}

func (m whilestmt) MarshalJSON() ([]byte, error) {
	return []byte("\"" + m.String() + "\""), nil
}

func (m whilestmt) String() string {
	msg := "default"
	if m.Do {
		msg = "do " + m.Body.String() + "while (" + infa2str(m.Expr) + ");"
	} else {
		msg = "while (" + infa2str(m.Expr) + ") " + m.Body.String()
	}

	return msg
}

type casestmt struct {
	Default bool        `json:"def,omitempty"`
	Expr    interface{} `json:"expr,omitempty"`
	Body    codeblock   `json:"body,omitempty"`
}

func (m casestmt) MarshalJSON() ([]byte, error) {
	return []byte("\"" + m.String() + "\""), nil
}

func (m casestmt) String() string {
	msg := "default"
	if !m.Default {
		msg = "case " + infa2str(m.Expr)
	}

	if m.Body.Body.Typ == stmtCase {
		return msg + ":\n" + m.Body.Body.String() + ""
	}

	return msg + ":\n" + m.Body.String() + ""
}

type switchstmt struct {
	Expr expression `json:"expr,omitempty"`
	Body codeblock  `json:"body,omitempty"`
}

func (m switchstmt) MarshalJSON() ([]byte, error) {
	return []byte("\"" + m.String() + "\""), nil
}

func (m switchstmt) String() string {
	msg := "switch(" + m.Expr.String() + ") " + m.Body.String()

	return msg
}

type stmttype int

func (m stmttype) MarshalJSON() ([]byte, error) {
	return []byte("\"" + m.String() + "\""), nil
}

func (m stmttype) String() string {
	return map[stmttype]string{
		stmtUnknown:  "unknown",
		stmtNone:     "none",
		stmtBlock:    "block",
		stmtBreak:    "break",
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
		stmtTypeDef:  "typedef",
		stmtFuncDecl: "funcdecl",
		stmtComment:  "comment",
		stmtLineEnd:  "\\n",
	}[m]
}

const (
	stmtUnknown stmttype = iota
	stmtNone
	stmtBlock
	stmtBreak
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
	stmtTypeDef
	stmtTypeDecl
	stmtFuncDecl
	stmtComment
	stmtLineEnd
)
