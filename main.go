package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	inputStr := ReadFile()
	lex := NewLexer(inputStr)
	for {
		tok := lex.NextToken()
		if tok.Type == TOKEN_EOF {
			PrintTokens(tok)
			break
		}
		PrintTokens(tok)
	}

	fmt.Println()
	if len(lex.errors) != 0 {
		fmt.Println("Ошибки разбора: ")
		for _, err := range lex.errors {
			fmt.Println(err)
		}
		fmt.Println()
	}

	fmt.Println("Таблица идентификаторов:")
	for key, val := range lex.identMap {
		fmt.Println(key, "->", val)
	}
}

func PrintTokens(tok Token) {
	switch tok.Type {
	case TOKEN_STRING:
		fmt.Printf("STRING (%d,%d)-(%d,%d): %s\n",
			tok.Pos.PosStart.Line, tok.Pos.PosStart.Col,
			tok.Pos.PosEnd.Line, tok.Pos.PosEnd.Col,
			tok.Value.(string))
	case TOKEN_REAL:
		fmt.Printf("REAL (%d,%d)-(%d,%d): %f\n",
			tok.Pos.PosStart.Line, tok.Pos.PosStart.Col,
			tok.Pos.PosEnd.Line, tok.Pos.PosEnd.Col,
			tok.Value.(float64)) // кастим к float64
	case TOKEN_IDENT:
		fmt.Printf("IDENT (%d,%d)-(%d,%d): %d\n",
			tok.Pos.PosStart.Line, tok.Pos.PosStart.Col,
			tok.Pos.PosEnd.Line, tok.Pos.PosEnd.Col,
			tok.Value.(int)) // кастим к int (индекс в таблице)
	case TOKEN_EOF:
		fmt.Printf("EOF (%d,%d)-(%d,%d)",
			tok.Pos.PosStart.Line, tok.Pos.PosStart.Col,
			tok.Pos.PosEnd.Line, tok.Pos.PosEnd.Col)
		fmt.Println()
	default:
		fmt.Println("неизвестный тип токена")
	}
}

func ReadFile() string {
	data, err := os.ReadFile("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	return string(data)
}
