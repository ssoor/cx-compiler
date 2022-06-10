package main

import (
	"fmt"
	"strings"
)

type function struct {
	Typ  funcref   `json:"type,omitempty"` // 类型
	Body codeblock `json:"body,omitempty"`
}

func (m function) MarshalJSON() ([]byte, error) {
	return []byte("\"" + m.String() + "\""), nil
}

func (m function) String() string {
	msg := m.Typ.String() + " "
	return msg + m.Body.String()
}

type funcref struct {
	Name   varpkgname     `json:"name,omitempty"`
	Ref    vardeclvaldecl `json:"ref,omitempty"`
	Self   funcself       `json:"self,omitempty"`
	Typ    functiontype   `json:"type,omitempty"` // 类型
	Retval typedecl       `json:"ret,omitempty"`
	Params []vardecl      `json:"params,omitempty"`
}

func (m funcref) MarshalJSON() ([]byte, error) {
	return []byte("\"" + m.String() + "\""), nil
}

func (m funcref) String() string {
	name := m.Name.String()

	params := ""
	for _, v := range m.Params {
		params += v.String() + ", "
	}

	extendMsg := ""
	switch m.Typ {
	case functionRef:
		name = "(" + m.Ref.Name.String() + ")" + name
		if m.Ref.Value != nil {
			extendMsg = " = " + infa2str(m.Ref.Value)
		}
	case functionSelf:
		typName := ""
		for _, v := range m.Self.Typ.Ref {
			switch v {
			case "*":
			case "&":
			default:
				typName = v
			}

			if len(typName) != 0 {
				break
			}
		}

		name = typName + "_" + name
		params = m.Self.String() + ", " + params
	}

	params = strings.TrimSuffix(params, ", ")
	return fmt.Sprintf("%s %s(%s)%s", m.Retval, name, params, extendMsg)
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
		functionRef:    "func_ref",
		functionSelf:   "func_self",
	}[m]
}

const (
	functionNormal functiontype = iota // 常规方法
	functionRef                        // 引用方法
	functionSelf                       // 成员方法
)
