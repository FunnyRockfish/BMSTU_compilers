package main

import (
	"fmt"
	"strconv"
	"strings"
)

func main() {
	inputString := "if if"
	tokens := strings.Fields(inputString)
	lex := NewLexer(tokens)
	parser := NewParser(lex)

	ast, err := parser.Parse()
	if err != nil {
		fmt.Println(err)
		return
	}

	PrintAST(ast)
}

type TokenType int

const (
	TOKEN_KEYWORD TokenType = iota
	TOKEN_NUMBER
	TOKEN_IDENTIFIER
)

type Token struct {
	Type  TokenType
	Value interface{}
}

func NewLexer(tokens []string) *Lexer {
	return &Lexer{
		Tokens:   tokens,
		KeyWords: []string{"define", "end", "if", "endif", "while", "do", "wend"},
	}
}

type Lexer struct {
	Tokens   []string
	KeyWords []string
}

func (l *Lexer) NextWord() string {
	currWord := l.Tokens[0]
	l.Tokens = l.Tokens[1:] // убираем первое слово

	return currWord
}

func (l *Lexer) IsItKeyWord(str string) bool {
	for _, word := range l.KeyWords {
		if word == str {
			return true
		}
	}
	return false
}

func (l *Lexer) Peek() string {
	if len(l.Tokens) > 0 {
		return l.Tokens[0]
	}
	return ""
}

func (l *Lexer) NextToken() *Token {
	if len(l.Tokens) == 0 {
		return nil
	}

	currToken := l.Tokens[0]
	l.Tokens = l.Tokens[1:]

	// Проверяем, является ли токен ключевым словом
	if l.IsItKeyWord(currToken) {
		return &Token{
			Type:  TOKEN_KEYWORD,
			Value: currToken,
		}
	}

	if num, err := strconv.Atoi(currToken); err == nil {
		return &Token{
			Type:  TOKEN_NUMBER,
			Value: num,
		}
	}

	return &Token{
		Type:  TOKEN_IDENTIFIER,
		Value: currToken,
	}
}
