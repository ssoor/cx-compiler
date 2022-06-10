package main

import (
	"fmt"
	"reflect"
	"strings"
)

// Storage? Typ+ Name([Arr])? (= Value)?
type vardecl struct {
	Storage string           `json:"storage,omitempty"`
	Typ     vardecltype      `json:"type,omitempty"`
	RefFunc funcref          `json:"fref,omitempty"`
	RefType typedecl         `json:"tref,omitempty"`
	Values  []vardeclvaldecl `json:"vals,omitempty"`
}

func (m vardecl) MarshalJSON() ([]byte, error) {
	return []byte("\"" + m.String() + "\""), nil
}

func (m vardecl) String() string {
	msg := ""
	switch m.Typ {
	case varFunc:
		return m.RefFunc.String()
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

type vardeclvaldecl struct {
	Name  vardeclnametype `json:"name,omitempty"`
	Value interface{}     `json:"val,omitempty"` // expression
}

func (m vardeclvaldecl) MarshalJSON() ([]byte, error) {
	return []byte("\"" + m.String() + "\""), nil
}

func (m vardeclvaldecl) String() string {
	msg := m.Name.String()

	if m.Value != nil {
		msg += " = " + infa2str(m.Value)
	}

	return msg
}

type varpkgname struct {
	Pkgs []string        `json:"pkgs,omitempty"`
	Name vardeclnametype `json:"name,omitempty"`
}

func (m varpkgname) MarshalJSON() ([]byte, error) {
	return []byte("\"" + m.String() + "\""), nil
}

func (m varpkgname) String() string {
	msg := ""
	if len(m.Pkgs) > 0 {
		msg += strings.Join(m.Pkgs, "_") + "_"
	}

	return msg + m.Name.String()
}

type vardeclnametype struct {
	Name string           `json:"name,omitempty"`
	Ops  []string         `json:"ops,omitempty"`
	Arr  []vardeclarrtype `json:"arr,omitempty"`
}

func (m vardeclnametype) MarshalJSON() ([]byte, error) {
	return []byte("\"" + m.String() + "\""), nil
}

func (m vardeclnametype) String() string {
	msg := ""
	for _, v := range m.Ops {
		msg += v
	}

	msg += m.Name
	for _, v := range m.Arr {
		msg += v.String()
	}

	return msg
}

type vardeclarrtype struct {
	Arr   bool   `json:"arr,omitempty"`   // 是否是数组
	Count string `json:"count,omitempty"` // 数组长度, -1 代表不定长
}

func (m vardeclarrtype) MarshalJSON() ([]byte, error) {
	return []byte("\"" + m.String() + "\""), nil
}

func (m vardeclarrtype) String() string {
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
	Fields []varrefField `json:"field,omitempty"`
}

func (m varref) MarshalJSON() ([]byte, error) {
	return []byte("\"" + m.String() + "\""), nil
}

func (m varref) String() string {
	msg := infa2str(m.Parent)

	for _, v := range m.Fields {
		msg += v.String()
	}

	return msg + ""
}

type varrefField struct {
	Ptr  bool       `json:"ptr,omitempty"`
	Name varpkgname `json:"name,omitempty"`
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
