package main

import (
	"fmt"
	"strings"
)

type function struct {
	Typ  funcref     `json:"type,omitempty"` // 类型
	Body SymbolTable `json:"body,omitempty"`
}

func (m function) MarshalJSON() ([]byte, error) {
	return []byte("\"" + m.String() + "\""), nil
}

func (m function) String() string {
	msg := m.Typ.String() + " {\n"

	for _, v := range m.Body.Statements {
		msg += "\t" + v.String() + "\n"
	}

	vars := []string{}
	for name, _ := range m.Body.Variables {
		vars = append(vars, name)
	}

	types := []string{}
	for name, _ := range m.Body.Enums {
		types = append(types, name)
	}
	for name, _ := range m.Body.Types {
		types = append(types, name)
	}

	funcs := []string{}
	for name, _ := range m.Body.Functions {
		funcs = append(funcs, name)
	}

	msg += "}\n" + fmt.Sprintf("var:%v type:%v func:%v\n", vars, types, funcs)

	return msg
}

type funcref struct {
	Name   nameref      `json:"name,omitempty"`
	Self   funcself     `json:"self,omitempty"`
	Typ    functiontype `json:"type,omitempty"` // 类型
	Retval typedecl     `json:"ret,omitempty"`
	Params []vardecl    `json:"params,omitempty"`
}

func (m funcref) MarshalJSON() ([]byte, error) {
	return []byte("\"" + m.String() + "\""), nil
}

func (m funcref) String() string {
	params := ""
	switch m.Typ {
	case functionSelf:
		params = m.Self.String() + ", "
	}

	for _, v := range m.Params {
		params += v.String() + ", "
	}
	params = strings.TrimSuffix(params, ", ")

	return fmt.Sprintf("%s %s(%s)", m.Retval, m.Name.String(), params)
}

type funcself struct {
	Name string   `json:"name,omitempty"`
	Typ  typedecl `json:"type,omitempty"` // 类型
}

func (m funcself) MarshalJSON() ([]byte, error) {
	return []byte("\"" + m.String() + "\""), nil
}

func (m funcself) String() string {
	return m.Typ.String() + " " + m.Name
}

func (m funcref) GenUniqueName() string {
	prefix := ""
	switch m.Typ {
	case functionSelf:
		prefix = m.Params[0].RefType.Name
	}

	return prefix + m.Name.String()
}

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
