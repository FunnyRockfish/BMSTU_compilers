package main

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"log"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("usage: astprint <filename.go>\n")
		return
	}

	// Создаем хранилище данных об исходных файлах
	fset := token.NewFileSet()

	// Парсим исходный файл
	file, err := parser.ParseFile(
		fset,                 // данные об исходниках
		os.Args[1],           // имя файла с исходником программы
		nil,                  // пусть парсер сам загрузит исходник
		parser.ParseComments, // сохраняем комментарии
	)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Для отладки: выводим AST в стандартный вывод
	ast.Fprint(os.Stdout, fset, file, nil)

	// Здесь вызываем функцию, которая изменяет AST (заменяет var на :=)
	changeDeclaration(file)

	// Теперь генерируем исходный код из модифицированного AST
	outFile, err := os.Create("new_code.go")
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()

	if err := format.Node(outFile, fset, file); err != nil {
		log.Fatalf("Ошибка при генерации кода: %v", err)
	}
}

func changeDeclaration(file *ast.File) {
	ast.Inspect(file, func(node ast.Node) bool {
		if funcDecl, ok := node.(*ast.FuncDecl); ok {
			fmt.Printf("Найдена функция: %s\n", funcDecl.Name.Name)
			for i, stmt := range funcDecl.Body.List {
				declStmt, ok := stmt.(*ast.DeclStmt)
				if !ok {
					continue
				}
				genDecl, ok := declStmt.Decl.(*ast.GenDecl)
				if !ok || genDecl.Tok != token.VAR {
					continue
				}

				if len(genDecl.Specs) != 1 {
					continue
				}
				valueSpec, ok := genDecl.Specs[0].(*ast.ValueSpec)
				if !ok {
					continue
				}
				if len(valueSpec.Values) != 1 || len(valueSpec.Names) != 1 || valueSpec.Type != nil {
					continue
				}
				asignStmt := &ast.AssignStmt{
					Lhs: []ast.Expr{valueSpec.Names[0]},
					Rhs: []ast.Expr{valueSpec.Values[0]},
					Tok: token.DEFINE,
				}

				funcDecl.Body.List[i] = asignStmt

			}
		}
		return true
	})
}
