package main

import (
	"errors"
	"fmt"
	"strings"
)

type Parser struct {
	tokens []*Token
	pos    int
}

func NewParser(l *Lexer) *Parser {
	var allTokens []*Token
	for {
		t := l.NextToken()
		if t == nil {
			break
		}
		allTokens = append(allTokens, t)
	}
	return &Parser{tokens: allTokens, pos: 0}
}

type AST struct {
	Articles map[string][]interface{}
	Body     []interface{}
}

func (p *Parser) Advance() {
	if p.pos < len(p.tokens) {
		p.pos++
	}
}

func (p *Parser) CurrToken() *Token {
	if p.pos < len(p.tokens) {
		return p.tokens[p.pos]
	}
	return nil
}

func (p *Parser) Parse() (*AST, error) {
	articles, err := p.ParseArticles()
	if err != nil {
		return nil, err
	}
	body, err := p.parseBody(nil)
	if err != nil {
		return nil, err
	}
	return &AST{
		Articles: articles,
		Body:     body,
	}, nil
}

// ParseArticles: <Articles> ::= <Article> <Articles> | ε
func (p *Parser) ParseArticles() (map[string][]interface{}, error) {
	articles := make(map[string][]interface{})
	for {
		currToken := p.CurrToken()
		if currToken == nil {
			break
		}
		// Если следующий токен не "define", статьи закончились
		if currToken.Type == TOKEN_KEYWORD && currToken.Value == "define" {
			name, body, err := p.ParseArticle()
			if err != nil {
				return nil, err
			}
			articles[name] = body
		} else {
			break
		}
	}
	return articles, nil
}

// ParseArticle: <Article> ::= define word <Body> end
func (p *Parser) ParseArticle() (string, []interface{}, error) {
	tok := p.CurrToken()
	if tok == nil || tok.Value != "define" {
		return "", nil, fmt.Errorf("Ожидалось 'define', получено %v", tok)
	}
	p.Advance() // пропускаем "define"

	tok = p.CurrToken()
	if tok == nil {
		return "", nil, fmt.Errorf("Ожидалось слово после 'define'")
	}
	if tok.Type != TOKEN_IDENTIFIER {
		return "", nil, fmt.Errorf("Ожидалось слово, получили %v", tok.Type)
	}
	articleName, ok := tok.Value.(string)
	if !ok {
		return "", nil, fmt.Errorf("Некорректное значение имени статьи")
	}
	p.Advance()
	stopTokens := map[interface{}]bool{"end": true}
	body, err := p.parseBody(stopTokens)
	if err != nil {
		return "", nil, err
	}

	tok = p.CurrToken()
	if tok == nil || tok.Value != "end" {
		return "", nil, errors.New("Ожидалось 'end' после тела статьи")
	}
	p.Advance()
	return articleName, body, nil
}

// parseBody: <Body> ::= if <Body> endif <Body>
//
//	| while <Body> do <Body> wend <Body>
//	| integer <Body> | word <Body> | ε
//
// stopTokens – токены, при встрече которых разбор тела прекращается.
func (p *Parser) parseBody(stopTokens map[interface{}]bool) ([]interface{}, error) {
	var result []interface{}

	for {
		tok := p.CurrToken()
		if tok == nil {
			break
		}
		if stopTokens != nil && stopTokens[tok.Value] {
			break
		}

		switch {
		case tok.Type == TOKEN_KEYWORD && tok.Value == "if":
			p.Advance() // пропускаем if

			ifStop := map[interface{}]bool{"endif": true}
			ifBody, err := p.parseBody(ifStop)
			if err != nil {
				return nil, err
			}

			tok = p.CurrToken()
			if tok == nil || tok.Value != "endif" {
				return nil, errors.New("Ожидалось 'endif' после if-тела")
			}
			p.Advance() // пропускаем endif

			result = append(result, []interface{}{"if", ifBody})
			continue

		case tok.Type == TOKEN_KEYWORD && tok.Value == "while":
			p.Advance() // пропускаем while

			// Парсим условие до do
			conditionStop := map[interface{}]bool{"do": true}
			condition, err := p.parseBody(conditionStop)
			if err != nil {
				return nil, err
			}

			tok = p.CurrToken()
			if tok == nil || tok.Value != "do" {
				return nil, errors.New("Ожидалось 'do' после условия while")
			}
			p.Advance() // do

			// Парсим тело до wend
			bodyStop := map[interface{}]bool{"wend": true}
			loopBody, err := p.parseBody(bodyStop)
			if err != nil {
				return nil, err
			}

			tok = p.CurrToken()
			if tok == nil || tok.Value != "wend" {
				return nil, errors.New("Ожидалось 'wend' после тела цикла while")
			}
			p.Advance() // wend

			result = append(result, []interface{}{"while", condition, loopBody})
			continue

		case tok.Type == TOKEN_NUMBER:
			strVal, ok := tok.Value.(int)
			if !ok {
				fmt.Println(tok.Value.(int))
				return nil, errors.New("Некорректное числовое значение")
			}
			num := strVal
			result = append(result, num)
			p.Advance()
			continue

		case tok.Type == TOKEN_IDENTIFIER:
			result = append(result, tok.Value)
			p.Advance()
			continue

		default:
			return nil, errors.New("Неожиданный токен: " + fmt.Sprintf("%v", tok.Value))
		}
	}

	return result, nil
}

func prettyPrintValue(val interface{}, indent int) string {
	indentStr := strings.Repeat(" ", indent)
	switch v := val.(type) {
	case []interface{}:
		var parts []string
		parts = append(parts, "(")
		for _, elem := range v {
			parts = append(parts, "\n"+indentStr+prettyPrintValue(elem, indent+2))
		}
		parts = append(parts, "\n"+strings.Repeat(" ", indent-2)+")")
		return strings.Join(parts, "")
	case string:
		return v
	case int:
		return fmt.Sprintf("%d", v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

func PrintAST(ast *AST) {
	fmt.Println("Articles:")
	for name, body := range ast.Articles {
		fmt.Printf("  %s: %s\n", name, prettyPrintValue(body, 4))
	}
	fmt.Println("Body:")
	fmt.Println(prettyPrintValue(ast.Body, 2))
}
