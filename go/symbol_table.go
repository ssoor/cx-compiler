package main

import (
	"fmt"
)

func NewSymbolTable() SymbolTable {
	st := SymbolTable{
		Parent:    nil,
		Enums:     make(map[string]enumdecl),
		Types:     make(map[string]typedecl),
		Variables: make(map[string]vardecl),
		Functions: make(map[string]function),
	}

	return st
}

func NewBodySymbolTable(stmts ...statement) SymbolTable {
	st := NewSymbolTable()

	for _, s := range stmts {
		st.AddStmt(s)
	}

	return st
}

func NewFuncSymbolTable(fun function, stmts ...statement) SymbolTable {
	st := NewSymbolTable()

	st.AddVar(vardecl{
		Typ:     varType,
		RefType: fun.Typ.Self.Typ,
		Values: []vardeclvaldecl{
			{
				Name: vardeclnametype{
					Name: fun.Typ.Self.Name,
				},
			},
		},
	})

	for _, v := range fun.Typ.Params {
		st.AddVar(v)
	}

	for _, s := range stmts {
		st.AddStmt(s)
	}

	return st
}

type SymbolTable struct {
	Enums     map[string]enumdecl `json:"enums,omitempty"`
	Types     map[string]typedecl `json:"types,omitempty"`
	Variables map[string]vardecl  `json:"vars,omitempty"`
	Functions map[string]function `json:"funs,omitempty"`

	Parent     *SymbolTable   `json:"-"`
	Childs     []*SymbolTable `json:"-"`
	Statements []statement    `json:"block,omitempty"`
}

func (m SymbolTable) Source() string {
	msg := ""
	for _, stmt := range m.Statements {
		msg += stmt.String() + "\n"
	}

	return msg
}

func (m SymbolTable) String() string {
	msg := "{\n"

	for _, stmt := range m.Statements {
		msg += stmt.String() + "\n"
	}

	return msg + "}"
}

func (m SymbolTable) Debug() string {
	msg := ""
	vars := []string{}
	for name, _ := range m.Variables {
		vars = append(vars, name)
	}

	types := []string{}
	for name, _ := range m.Enums {
		types = append(types, name)
	}
	for name, _ := range m.Types {
		types = append(types, name)
	}

	funcs := []string{}
	for name, _ := range m.Functions {
		funcs = append(funcs, name)
	}

	if len(vars) != 0 {
		msg += fmt.Sprintf(" var:%v", vars)
	}

	if len(types) != 0 {
		msg += fmt.Sprintf(" type:%v", types)
	}

	if len(funcs) != 0 {
		msg += fmt.Sprintf(" func:%v", funcs)
	}

	if len(msg) != 0 {
		msg = fmt.Sprintf(" /*%s */ ", msg)
	}

	return msg
}

func (m *SymbolTable) AddStmt(stmt statement) {
	switch stmt.Typ {
	case stmtIf:
	case stmtFor:
		forStmt := stmt.Stmt.(forstmt)

		if decl, ok := forStmt.Init.(vardecl); ok {
			if forStmt.Body.Body.Typ == stmtBlock {
				st := forStmt.Body.Body.Stmt.(SymbolTable)

				st.AddVar(decl)
			}
		}
	case stmtCase:
	case stmtWhile:
	case stmtSwitch:
	case stmtReturn:
	case stmtExpr:
	case stmtAssign:
	case stmtVarDecl:
		switch varStmt := stmt.Stmt.(type) {
		case vardecl:
			m.AddVar(varStmt)
		case opexpr:
			if opStmt, ok := varStmt.L.Expr.(opexpr); ok && opStmt.Op == "*" {
				decl := vardecl{
					Typ: varType,
					RefType: typedecl{
						Typ: typedeclRef,
						Ref: []string{opStmt.L.Expr.(varref).String(), opStmt.Op},
					},
					Values: []vardeclvaldecl{
						{
							Name: vardeclnametype{
								Name: opStmt.R.Expr.(varref).String(),
							},
							Value: varStmt.R,
						},
					},
				}
				m.AddVar(decl)
				stmt.Stmt = decl
			}
		}
	case stmtEnumDecl:
		m.AddEnum(stmt.Stmt.(enumdecl))
	case stmtTypeDecl:
		m.AddType(stmt.Stmt.(typedecl))
	case stmtBlock:
	case stmtNone:
	case stmtBreak:
	case stmtTypeDef:
	case stmtFuncDecl:
	case stmtComment:
	case stmtLineEnd:
	default:
		panic(fmt.Sprintf("unknown statement [%+v]", stmt.Typ))
	}

	m.Statements = append(m.Statements, stmt)
}

func (m *SymbolTable) AddEnum(decl enumdecl) {
	if m.Enums == nil {
		m.Enums = make(map[string]enumdecl)
	}

	m.Enums[decl.Name] = decl
}

func (m *SymbolTable) AddType(decl typedecl) {
	if m.Types == nil {
		m.Types = make(map[string]typedecl)
	}

	m.Types[decl.Name] = decl

	m.AddVar(vardecl{
		Typ:     varType,
		RefType: decl,
		Values:  decl.Values,
	})
}

func (m *SymbolTable) AddTypedef(decl typedefine) {
	switch decl := decl.Typ.(type) {
	case enumdecl:
		m.AddEnum(decl)
	case typedecl:
		m.AddType(decl)
	default:
		panic("unknown type")
	}
}

func (m *SymbolTable) AddVar(decl vardecl) {
	if m.Variables == nil {
		m.Variables = make(map[string]vardecl)
	}

	for _, v := range decl.Values {
		m.Variables[v.Name.Name] = decl
	}
}

func (m *SymbolTable) AddFun(decl function) {
	if m.Functions == nil {
		m.Functions = make(map[string]function)
	}

	m.Functions[decl.Typ.GenUniqueName()] = decl
}

func (m *SymbolTable) CheckExist(name string) bool {
	if m.CheckEnumExist(name) {
		return true
	}

	if m.CheckTypeExist(name) {
		return true
	}

	if m.CheckVariableExist(name) {
		return true
	}

	if m.CheckFunctionExist(name) {
		return true
	}

	if m.Parent == nil {
		return false
	}

	return m.Parent.CheckExist(name)
}

func (m *SymbolTable) CheckEnumExist(name string) bool {
	if _, exist := m.Enums[name]; exist {
		return true
	}

	if m.Parent == nil {
		return false
	}

	return m.Parent.CheckEnumExist(name)
}

func (m *SymbolTable) CheckTypeExist(name string) bool {
	if _, exist := m.Types[name]; exist {
		return true
	}

	if m.Parent == nil {
		return false
	}

	return m.Parent.CheckEnumExist(name)
}

func (m *SymbolTable) CheckVariableExist(name string) bool {
	if _, exist := m.Variables[name]; exist {
		return true
	}

	if m.Parent == nil {
		return false
	}

	return m.Parent.CheckEnumExist(name)
}

func (m *SymbolTable) CheckFunctionExist(name string) bool {
	if _, exist := m.Functions[name]; exist {
		return true
	}

	if m.Parent == nil {
		return false
	}

	return m.Parent.CheckEnumExist(name)
}
