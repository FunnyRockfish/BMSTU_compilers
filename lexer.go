package main

import (
	"fmt"
	"strconv"
	"unicode"
)

type TokenType int

const (
	TOKEN_EOF TokenType = iota
	TOKEN_STRING
	TOKEN_REAL
	TOKEN_IDENT
)

type Token struct {
	Type  TokenType
	Value interface{}
	Pos   TokenPosition
}

type Position struct {
	Line  int
	Col   int
	Index int
}

type TokenPosition struct {
	PosStart Position
	PosEnd   Position
}

type Lexer struct {
	position     Position
	inputRunes   []rune
	identMap     map[int]string
	identCounter int
	errors       []string
}

func NewLexer(input string) *Lexer {
	return &Lexer{
		inputRunes:   []rune(input),
		position:     Position{Line: 1, Col: 1, Index: 0},
		identMap:     map[int]string{},
		identCounter: 0,
		errors:       make([]string, 0),
	}
}

func (l *Lexer) NextToken() Token {
	l.SkipWhitespace()
	tok := Token{}

	currChar := l.GetCurrentChar()

	switch {
	case currChar == -1:
		tok.Pos.PosStart = l.position
		tok.Type = TOKEN_EOF
		tok.Value = ""
		tok.Pos.PosEnd = l.position
	case currChar == '\'':
		tok.Pos.PosStart = l.position
		tok.Value = l.RecognizeStringLiteral()
		tok.Type = TOKEN_STRING
		tok.Pos.PosEnd = l.position
		tok.Pos.PosEnd.Col--
	case unicode.IsDigit(currChar) || currChar == '-':
		var err error
		tok.Pos.PosStart = l.position
		tok.Value, err = strconv.ParseFloat(l.RecognizeRealNumber(), 64)
		tok.Type = TOKEN_REAL
		tok.Pos.PosEnd = l.position
		tok.Pos.PosEnd.Col--
		if err != nil {
			fmt.Println(err)
		}
	case unicode.IsLetter(currChar):
		tok.Pos.PosStart = l.position
		tok.Value = l.RecognizeIdent()
		tok.Type = TOKEN_IDENT
		tok.Pos.PosEnd = l.position
		tok.Pos.PosEnd.Col--
	default:
		err := fmt.Sprintf("неизвестный символ %c на позиции %d, %d", currChar, l.position.Line, l.position.Col)
		l.errors = append(l.errors, err)
		l.ConsumeSymbol()
		return l.NextToken()
		//fmt.Println("неизвестный символ: ", string(currChar))
	}

	return tok
}

func (l *Lexer) RecognizeStringLiteral() string {
	l.ConsumeSymbol()
	str := ""
	for l.position.Index < len(l.inputRunes) {
		currChar := l.GetCurrentChar()
		if currChar == '\'' {
			if l.PeekNextChar() == '\'' {
				l.ConsumeSymbol()
				str += "'"
			} else {
				l.ConsumeSymbol()
				return str
			}
		} else if currChar == '\n' {
			err := fmt.Sprintf("Ошибка (%d,%d): Строка не может начинаться на одном лайне и заканчиваться на другом!", l.position.Line, l.position.Col)
			l.errors = append(l.errors, err)
		} else {
			str += string(currChar)
		}
		l.ConsumeSymbol()
	}
	return str
}

func (l *Lexer) RecognizeRealNumber() string {
	var realNumStr string

	if l.position.Index < len(l.inputRunes) && l.GetCurrentChar() == '-' {
		realNumStr += string(l.GetCurrentChar())
		l.ConsumeSymbol()
	}

	for l.position.Index < len(l.inputRunes) && unicode.IsDigit(l.GetCurrentChar()) {
		realNumStr += string(l.GetCurrentChar())
		l.ConsumeSymbol()
	}

	if l.position.Index < len(l.inputRunes) && l.GetCurrentChar() == '.' {
		realNumStr += string(l.GetCurrentChar())
		l.ConsumeSymbol()
		for l.position.Index < len(l.inputRunes) && unicode.IsDigit(l.GetCurrentChar()) {
			realNumStr += string(l.GetCurrentChar())
			l.ConsumeSymbol()
		}
	}

	return realNumStr
}

func (l *Lexer) RecognizeIdent() int {
	ident := ""
	for l.position.Index < len(l.inputRunes) {
		currChar := l.GetCurrentChar()
		if unicode.IsLetter(currChar) || currChar == '.' || unicode.IsDigit(currChar) {
			ident += string(currChar)
		} else {
			break
		}
		l.ConsumeSymbol()
	}

	if ident != "" {

		identIdx := l.FindIdentIdx(ident)
		if identIdx == -1 {
			l.identMap[l.identCounter] = ident
			l.identCounter++
			return l.identCounter - 1
		} else {
			return identIdx
		}

	}
	return -1
}

func (l *Lexer) FindIdentIdx(ident string) int {
	for key, val := range l.identMap {
		if val == ident {
			return key
		}
	}
	return -1
}

func (l *Lexer) SkipWhitespace() {
	for unicode.IsSpace(l.GetCurrentChar()) {
		l.ConsumeSymbol()
	}
}

func (l *Lexer) PeekNextChar() rune {
	nextIdx := l.position.Index + 1
	if nextIdx < len(l.inputRunes) {
		return l.inputRunes[nextIdx]
	}
	return 0
}

func (l *Lexer) GetCurrentChar() rune {
	if l.position.Index >= len(l.inputRunes) {
		return -1
	}
	currChar := l.inputRunes[l.position.Index]
	return currChar
}

func (l *Lexer) ConsumeSymbol() {
	currChar := l.inputRunes[l.position.Index]
	if currChar == '\n' {
		l.position.Col = 1
		l.position.Line++
	} else {
		l.position.Col++
	}
	l.position.Index++
}

func (l *Lexer) GetAndConsumeChar() rune {
	currChar := l.GetCurrentChar()
	l.ConsumeSymbol()
	return currChar
}
