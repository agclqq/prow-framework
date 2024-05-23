package make

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
)

func parserDbConfig(filename string) {
	// 指定要解析的文件

	// 创建文件集
	fset := token.NewFileSet()

	// 解析文件
	file, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		fmt.Println("Error parsing file:", err)
		os.Exit(1)
	}

	// 寻找并提取 db 变量
	var foundDb map[string]map[string]string
	ast.Inspect(file, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.AssignStmt:
			// 检查是否是对 db 变量的赋值语句
			for _, expr := range x.Lhs {
				if ident, ok := expr.(*ast.Ident); ok && ident.Name == "db" {
					// 获取赋值的表达式
					if mapExpr, ok := x.Rhs[0].(*ast.CompositeLit); ok {
						foundDb = make(map[string]map[string]string)
						// 遍历 map 的键值对
						for _, elem := range mapExpr.Elts {
							if kv, ok := elem.(*ast.KeyValueExpr); ok {
								if keyIdent, ok := kv.Key.(*ast.BasicLit); ok && keyIdent.Kind == token.STRING {
									key := keyIdent.Value[1 : len(keyIdent.Value)-1] // 去除引号
									// 获取值的表达式
									if valueIdent, ok := kv.Value.(*ast.CompositeLit); ok {
										valueMap := make(map[string]string)
										// 遍历嵌套 map 的键值对
										for _, valueElem := range valueIdent.Elts {
											if valueKv, ok := valueElem.(*ast.KeyValueExpr); ok {
												if valueKeyIdent, ok := valueKv.Key.(*ast.Ident); ok {
													if valueLit, ok := valueKv.Value.(*ast.BasicLit); ok && valueLit.Kind == token.STRING {
														valueMap[valueKeyIdent.Name] = valueLit.Value[1 : len(valueLit.Value)-1] // 去除引号
													}
												}
											}
										}
										foundDb[key] = valueMap
									}
								}
							}
						}
					}
					return false
				}
			}
		}
		return true
	})

	// 打印找到的 db 变量
	fmt.Println("db from config.go:", foundDb)
}
