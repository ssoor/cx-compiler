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

type normaltypedecl struct {
	Name string   `json:"name,omitempty"`
	Ref  []string `json:"ref,omitempty"`
}

func (m normaltypedecl) MarshalJSON() ([]byte, error) {
	return []byte("\"" + m.String() + "\""), nil
}

func (m normaltypedecl) String() string {
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
}

type blocktypedecl struct {
	Name   string       `json:"name,omitempty"`
	Typ    typedecltype `json:"type,omitempty"` // 类型
	Fields []vardecl    `json:"field,omitempty"`
}

func (m blocktypedecl) MarshalJSON() ([]byte, error) {
	return []byte("\"" + m.String() + "\""), nil
}

func (m blocktypedecl) String() string {
	switch m.Typ {
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

type typedecl struct {
	Name   string       `json:"name,omitempty"`
	Typ    typedecltype `json:"type,omitempty"` // 类型
	Ref    []string     `json:"ref,omitempty"`
	Fields []vardecl    `json:"field,omitempty"`
	Values []vardeclval `json:"vals,omitempty"`
}

func (m typedecl) MarshalJSON() ([]byte, error) {
	return []byte("\"" + m.String() + "\""), nil
}

func (m typedecl) String() string {
	msg := ""
	switch m.Typ {
	case typedeclRef:
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
		msg = "struct " + m.Name + "{\n"
		for _, v := range m.Fields {
			msg += "\t" + v.String() + ";\n"
		}

		msg += "}"
	case typedeclUnion:
		msg = "union " + m.Name + "{\n"
		for _, v := range m.Fields {
			msg += "\t" + v.String() + ";\n"
		}

		msg += "}"
	default:
		msg = fmt.Sprintf("!panic(%s)", reflect.TypeOf(m).String())
	}

	valuesMsg := ""
	for _, v := range m.Values {
		valuesMsg += v.String() + ", "
	}
	valuesMsg = strings.TrimSuffix(valuesMsg, ", ")

	if len(valuesMsg) != 0 {
		msg += " "
	}

	return msg + valuesMsg
}

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
	typedeclEnum                 // 枚举
	typedeclRef                  // 引用
	typedeclStruct               // 结构
	typedeclUnion                // 联合
)
