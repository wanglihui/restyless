package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/wanglihui/restyless/tpl"
)

type vistor struct {
	interName  string
	structName string
	fset       *token.FileSet
	dst        io.Writer
}

func (it *vistor) Visit(node ast.Node) ast.Visitor {
	switch node.(type) {
	case *ast.GenDecl:
		idx := -1
		genDecl := node.(*ast.GenDecl)
		// 处理import
		if genDecl.Tok == token.IMPORT {
			genDecl.Specs = append(genDecl.Specs, &ast.ImportSpec{
				Doc:     &ast.CommentGroup{},
				Name:    &ast.Ident{},
				Path:    &ast.BasicLit{ValuePos: 0, Kind: token.STRING, Value: `"github.com/go-resty/resty/v2"`},
				Comment: &ast.CommentGroup{},
				EndPos:  0,
			})
			genDecl.Specs = append(genDecl.Specs, &ast.ImportSpec{
				Doc:     &ast.CommentGroup{},
				Name:    &ast.Ident{},
				Path:    &ast.BasicLit{ValuePos: 0, Kind: token.STRING, Value: `"github.com/wanglihui/httperror"`},
				Comment: &ast.CommentGroup{},
				EndPos:  0,
			})
			if err := format.Node(it.dst, it.fset, genDecl); err != nil {
				panic(err)
			}
		}
		// 找到要处理的interface
		if genDecl.Tok == token.TYPE {
			for i, v := range genDecl.Specs {
				typeSepc, ok := v.(*ast.TypeSpec)
				if ok && typeSepc.Name.Name == it.interName {
					idx = i
					break
				}
			}
			if idx < 0 {
				return nil
			}
			var interDocParam docstore
			v := genDecl.Specs[idx].(*ast.TypeSpec)
			t, ok := v.Type.(*ast.InterfaceType)
			if !ok {
				return nil
			}
			// 写入structImp
			structData := tpl.StructData{
				TypeName: it.structName,
			}
			tpl.StructTemplate(it.dst, structData)
			if genDecl.Doc != nil && len(genDecl.Doc.List) > 0 {
				interDocParam = parseDoc(genDecl.Doc.List)
			}
			// 处理interface中定义的函数
			for _, v := range t.Methods.List {
				// 函数上注解信息，例如host=http://www, url=/token
				var funcDocParam docstore
				if v.Doc != nil && len(v.Doc.List) > 0 {
					funcDocParam = parseDoc(v.Doc.List)
				}
				funcDocParam = defaultParam(funcDocParam, interDocParam)
				funcParams := []tpl.Param{}
				params := v.Type.(*ast.FuncType).Params
				for _, v := range params.List {
					name := v.Names[0].Name
					var typ = ""
					switch v.Type.(type) {
					case *ast.SelectorExpr:
						typ = v.Type.(*ast.SelectorExpr).X.(*ast.Ident).Name + "." + v.Type.(*ast.SelectorExpr).Sel.Name
					case *ast.Ident:
						if v.Type.(*ast.Ident).Obj != nil {
							typ = v.Type.(*ast.Ident).Obj.Name
						} else {
							typ = v.Type.(*ast.Ident).Name
						}
					}
					if typ != "" {
						funcParams = append(funcParams, tpl.Param{Key: name, TypeVal: typ})
					}
				}
				returnParams := []string{}
				returns := v.Type.(*ast.FuncType).Results
				for _, v := range returns.List {
					typ := ""
					switch v.Type.(type) {
					case *ast.SelectorExpr:
						typ = v.Type.(*ast.SelectorExpr).X.(*ast.Ident).Name + "." + v.Type.(*ast.SelectorExpr).Sel.Name
					case *ast.Ident:
						if v.Type.(*ast.Ident) != nil {
							typ = v.Type.(*ast.Ident).Name
						}
					}
					if typ != "" {
						returnParams = append(returnParams, typ)
					}
				}
				// fmt.Println(funcParam)
				tpl.FuncTemplate(it.dst, tpl.FuncData{
					Name:          v.Names[0].Name,
					Params:        funcParams,
					Returns:       returnParams,
					InterfaceName: it.structName,
					URL:           funcDocParam["host"] + funcDocParam["url"],
				})
			}
		}
	}
	return it
}

type docstore = map[string]string

func parseDoc(docList []*ast.Comment) docstore {
	m := make(docstore)
	for _, v := range docList {
		doc := v
		reg, _ := regexp.Compile(`^//\s*`)
		text := string(reg.ReplaceAll([]byte(doc.Text), []byte("")))
		params := []string{text}
		if strings.Contains(text, ",") {
			params = strings.Split(text, ",")
		}
		for _, v := range params {
			keyval := strings.Split(v, "=")
			if len(keyval) == 2 {
				m[keyval[0]] = keyval[1]
			}
		}
	}
	return m
}
func defaultParam(funcParam, interParam docstore) docstore {
	for k, v := range interParam {
		if _, ok := funcParam[k]; !ok {
			funcParam[k] = v
		}
	}
	return funcParam
}

func main() {
	flag.Parse()
	var (
		pkgName    = os.Getenv("GOPACKAGE")
		fileName   = os.Getenv("GOFILE")
		line       = os.Getenv("GOLINE")
		typeName   = os.Args[1]
		structName = typeName + "Impl"
	)
	fmt.Println(os.Args)
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	// fmt.Println(pkgName, fileName, line, dir)
	//使用正则，找出接口名称
	fset := token.NewFileSet()
	fpath := strings.Join([]string{dir, fileName}, "/")
	f, err := parser.ParseFile(fset, fpath, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	// ast.Print(fset, f)
	fname := strings.ToLower(dir + "/" + typeName + ".gen.go")
	fs, err := os.OpenFile(fname, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer fs.Close()
	data := tpl.HeadData{
		PkgName:  pkgName,
		FileName: fileName,
		Line:     line,
	}
	tpl.FuncHead(fs, data)
	v := &vistor{
		interName:  typeName,
		structName: structName,
		fset:       fset,
		dst:        fs,
	}
	ast.Walk(v, f)
}
