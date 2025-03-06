package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"unicode/utf8"
)

type Token struct {
	Tag   string
	Value string
	Line  int
	Col   int
}

type Lexer struct {
	text string
	pos  int
	line int
	col  int
}

var reIdent = regexp.MustCompile(`^[A-Z]+(?:[+\-\*]+)*`)
var reNumber = regexp.MustCompile(`^(?:\*|[+]+|[-]+)`)

func (lx *Lexer) nextToken() *Token {
	for lx.pos < len(lx.text) {
		lx.skipWhitespace()
		if lx.pos >= len(lx.text) {
			break
		}

		startLine, startCol := lx.line, lx.col
		substr := lx.text[lx.pos:]

		type candidate struct {
			tag      string
			value    string
			length   int
			priority int
		}
		var candidates []candidate

		addCandidate := func(tag, val string, priority int) {
			candidates = append(candidates, candidate{tag, val, len(val), priority})
		}

		keywords := []struct {
			tag           string
			word          string
			checkBoundary bool // нужна ли проверка на границу, чтобы не распознавать как преф идент
		}{
			{"ON", "ON", true},
			{"OFF", "OFF", true},
			{"**", "**", false},
		}

		for _, kw := range keywords {
			if kw.checkBoundary {
				if match, ok := matchKeyword(lx.text, lx.pos, kw.word); ok {
					addCandidate(kw.tag, match, 3)
				}
			} else {
				if strings.HasPrefix(substr, kw.word) {
					addCandidate(kw.tag, kw.word, 3)
				}
			}
		}

		if m := reIdent.FindString(substr); m != "" {
			addCandidate("IDENT", m, 2)
		}

		if m := reNumber.FindString(substr); m != "" {
			addCandidate("NUMBER", m, 1)
		}

		if len(candidates) > 0 {
			best := candidates[0]
			for _, cand := range candidates[1:] {
				if cand.length > best.length || (cand.length == best.length && cand.priority > best.priority) {
					best = cand
				}
			}
			token := &Token{best.tag, best.value, startLine, startCol}
			lx.advance(best.value)
			return token
		}

		fmt.Printf("syntax error (%d,%d)\n", lx.line, lx.col)
		lx.errorRecovery()
	}
	return nil
}

func (lx *Lexer) skipWhitespace() {
	for lx.pos < len(lx.text) {
		r, _ := utf8.DecodeRuneInString(lx.text[lx.pos:])
		if !isWhitespace(r) {
			break
		}
		lx.advance(string(r))
	}
}

func (lx *Lexer) advance(s string) {
	for _, r := range s {
		if r == '\n' {
			lx.line++
			lx.col = 1
		} else {
			lx.col++
		}
	}
	lx.pos += len(s)
}

func (lx *Lexer) errorRecovery() {
	for lx.pos < len(lx.text) {
		r, _ := utf8.DecodeRuneInString(lx.text[lx.pos:])
		if isWhitespace(r) || isValidTokenStart(r) {
			break
		}
		lx.advance(string(r))
	}
}

func matchKeyword(text string, pos int, keyword string) (string, bool) {
	if !strings.HasPrefix(text[pos:], keyword) {
		return "", false
	}
	if keyword == "**" {
		return keyword, true
	}
	if pos+len(keyword) < len(text) {
		r, _ := utf8.DecodeRuneInString(text[pos+len(keyword):])
		if isIdentContinuation(r) {
			return "", false
		}
	}
	return keyword, true
}

func isWhitespace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n' || r == '\r'
}

func isIdentContinuation(r rune) bool {
	return (r >= 'A' && r <= 'Z') || r == '+' || r == '-' || r == '*'
}

func isValidTokenStart(r rune) bool {
	return (r >= 'A' && r <= 'Z') || r == '*' || r == '+' || r == '-'
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Использование: go run main.go <input_file>")
		return
	}

	filename := os.Args[1]
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Ошибка чтения файла:", err)
		return
	}

	lexer := Lexer{text: string(data), pos: 0, line: 1, col: 1}

	for {
		token := lexer.nextToken()
		if token == nil {
			break
		}
		fmt.Printf("%s (%d, %d): %s\n", token.Tag, token.Line, token.Col, token.Value)
	}
}
