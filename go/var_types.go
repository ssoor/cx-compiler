package main

import (
	"fmt"
	"reflect"
	"strings"
)

// Storage? Typ+ Name([Arr])? (= Value)?
type vardecl struct {
	Storage string       `json:"storage,omitempty"`
	Typ     vardecltype  `json:"type,omitempty"`
	RefType typedecl     `json:"tref,omitempty"`
	RefFunc function     `json:"fref,omitempty"`
	Values  []vardeclval `json:"vals,omitempty"`
}

func (m vardecl) MarshalJSON() ([]byte, error) {
	return []byte("\"" + m.String() + "\""), nil
}

func (m vardecl) String() string {
	msg := ""
	switch m.Typ {
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
