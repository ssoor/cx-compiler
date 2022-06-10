package main

import "strings"

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
