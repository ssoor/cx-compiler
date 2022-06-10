package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// The parser expects the lexer to return 0 on EOF.  Give it a name
// for clarity.
const eof = 0

type token int

func (m token) String() string {
	return exprToknameMap[int(m)]
}

type lexVal struct {
	token  token
	length int
	column int
	text   string
}

// The parser uses the type <prefix>Lex as a lexer. It must provide
// the methods Lex(*<prefix>SymType) int and Error(string).
type exprLex struct {
	readStr string
	reader  *bufio.Reader
	val     lexVal
}

func newLex(fileName string) (*exprLex, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	return &exprLex{
		readStr: "",
		reader:  bufio.NewReader(bytes.NewBuffer(data)),
	}, nil
}

// The parser calls this method to get each new token. This
// implementation returns operators and NUM.
func (x *exprLex) Lex(yylval *exprSymType) int {
	for {
		x.val.token = 0
		x.val.length = 0
		x.val.column = 0
		x.val.text = ""

		x.val.token = token(x.token())
		if x.val.token == -1 {
			return eof
		}

		x.val.column = x.number()
		x.val.length = x.number()
		x.val.text += x.string(x.val.length)

		x.readStr += x.val.text
		switch x.val.token {
		case LINEEND:
			x.val.token = 0
		case COMMENT, BLOCK_COMMENT:
			x.val.token = 0
		case IGNORE:
			x.val.token = 0
			// ST().AddStmt(statement{Typ: stmtLineEnd, Stmt: x.val.text})
		}

		if x.val.token != 0 {
			break
		}
	}

	yylval.lval = x.val

	return int(x.val.token)
}

// Lex a number.
func (x *exprLex) number() int {
	line, err := x.reader.ReadString('\n')
	if err != nil {
		if err == io.EOF {
			return -1
		}
		x.Error(err.Error())
		return -1
	}

	line = strings.TrimRight(line, "\n")
	token, err := strconv.Atoi(string(line))
	if err != nil {
		x.Error(err.Error())
		return -1
	}
	return token
}

// Lex a number.
func (x *exprLex) token() int {
	return x.number()
}

// Return the next rune for the lexer.
func (x *exprLex) string(size int) string {
	if size <= 0 {
		return ""
	}

	data := make([]byte, size)
	rlen, err := io.ReadFull(x.reader, data)
	if err != nil {
		x.Error(err.Error())
		return ""
	}

	if rlen != len(data) {
		x.Error(fmt.Sprintf("%s: %d-%d %d", "读取长度不对", size, len(data), rlen))
		return ""
	}

	return string(data)
}

// The parser calls this method on a parse error.
func (x *exprLex) Error(s string) {
	blankn := x.val.column
	placen := x.val.length
	if blankn > placen {
		blankn -= placen
	} else {
		blankn = 0
		placen = blankn
	}

	msg := string(x.readStr) + "\n"
	for i := 0; i < blankn; i++ {
		msg += " "
	}
	for i := 0; i < placen; i++ {
		msg += "^"
	}
	msg += fmt.Sprintf(" [%d] %s\n", x.val.token, s)

	fmt.Printf("%s", msg)
}
