package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	consoleMode := flag.Bool("c", false, "Однопроходный")
	fileMode := flag.Bool("f", false, "Разбор файла")

	flag.Parse()
	var lexer *Lexer

	if *consoleMode {
		reader := bufio.NewReader(os.Stdin)
		lexer = NewLexer(reader)
	} else if *fileMode {
		inputStr := ReadFile()
		lexer = NewLexer(inputStr)
	}
	for {
		tok := lexer.NextToken()
		if tok.Type == TOKEN_EOF {
			PrintTokens(tok)
			break
		}
		PrintTokens(tok)

		if tok.Type == TOKEN_IDENT && lexer.reader != nil {
			fmt.Println("Таблица идентификаторов:")
			for key, val := range lexer.identMap {
				fmt.Println(key, "->", val)
			}
		}

		if len(lexer.errors) != 0 && lexer.reader != nil {
			fmt.Println("Ошибки разбора:")
			for _, err := range lexer.errors {
				fmt.Println(err)
			}
			fmt.Println()

			// Важно: очищаем ошибки после вывода!
			lexer.errors = nil
		}
	}

	fmt.Println()
	if len(lexer.errors) != 0 {
		fmt.Println("Ошибки разбора: ")
		for _, err := range lexer.errors {
			fmt.Println(err)
		}
		fmt.Println()
	}

	fmt.Println("Таблица идентификаторов:")
	for key, val := range lexer.identMap {
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
