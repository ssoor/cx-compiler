package main

func NewFuncSymbolTable(fun function, stmts ...statement) SymbolTable {
	st := SymbolTable{
		Parent:    nil,
		Enums:     make(map[string]enumdecl),
		Types:     make(map[string]typedecl),
		Variables: make(map[string]vardecl),
		Functions: make(map[string]function),
	}

	st.AddVar(vardecl{
		Typ:     varType,
		RefType: fun.Typ.Self.Typ,
		Values: []vardeclval{
			{
				Name: fun.Typ.Self.Name,
			},
		},
	})

	for _, v := range fun.Typ.Params {
		st.AddVar(v)
	}

	st.AddStmts(stmts...)

	return st
}

type SymbolTable struct {
	Parent     *SymbolTable   `json:"-"`
	Childs     []*SymbolTable `json:"-"`
	Statements []statement    `json:"block,omitempty"`

	Enums     map[string]enumdecl `json:"enums,omitempty"`
	Types     map[string]typedecl `json:"types,omitempty"`
	Variables map[string]vardecl  `json:"vars,omitempty"`
	Functions map[string]function `json:"funs,omitempty"`
}

func (m *SymbolTable) AddStmts(stmt ...statement) {
	for _, s := range stmt {
		m.AddStmt(s)
	}
}

func (m *SymbolTable) AddStmt(stmt statement) {

	switch stmt.Typ {
	case stmtIf:
	case stmtFor:
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
					Values: []vardeclval{
						{
							Name:  opStmt.R.Expr.(varref).String(),
							Value: varStmt.R,
						},
					},
				}
				m.AddVar(decl)
				stmt.Stmt = decl
			}
			// m.AddFun(stmt)
		}
	case stmtEnumDecl:
		m.AddEnum(stmt.Stmt.(enumdecl))
	case stmtTypeDecl:
		m.AddType(stmt.Stmt.(typedecl))
	default:
		panic("unknown statement")
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
		m.Variables[v.Name] = decl
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
