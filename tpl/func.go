package tpl

import (
	"fmt"
	"html/template"
	"io"
	"strings"
)

type Param struct {
	Key     string
	TypeVal string
	Star    bool
}

type FuncData struct {
	Name          string
	URL           string
	InterfaceName string
	Params        []Param
	Returns       []Param
}

func FuncTemplate(dst io.Writer, data FuncData) {
	var tpl = `
func (it *{{.InterfaceName}}Impl) {{.Name}} ({{.ParamStr}}) ({{.ReturnStr}}) {
	var e httperror.HTTPError
	r := it.r.R().SetError(e)
	{{with .Headers -}}{{range .}}
	r = r.SetHeader("{{.Key}}", string({{.Key}}))
	{{end}}{{end}}
	{{with .Query -}}{{range .}}
	r=r.SetQueryParam("{{.Key}}", string({{.Key}}))
	{{end}}{{end}}
	{{with .PathParams -}}{{range .}}
	r=r.SetPathParam("{{.Key}}", string({{.Key}}))
	{{end}}{{end}}
	{{if .Body}}{{if .Body.Key}}r=r.SetBody({{.Body.Key}}){{end}}{{end}}
	{{if .ReturnStruct.TypeVal}}
	{{if .ReturnStruct.Star}}
	var ret = new({{.ReturnStruct.TypeVal}})
	{{else}}
	var ret {{.ReturnStruct.TypeVal}}
	{{end}}
	r = r.SetResult(&ret)
	_, err := r.{{.Method}}("{{.URL}}")
	return ret, err
	{{else}}
	_, err := r.{{.Method}}("{{.URL}}")
	return err
	{{end}}
}
`
	type comData struct {
		FuncData
		ReturnStr    string
		ParamStr     string
		ReturnStruct Param
		Headers      []Param
		Query        []Param
		PathParams   []Param
		Body         Param
		Method       string
	}
	paramStr := []string{}
	returnStr := []string{}
	for _, v := range data.Params {
		var fullType = v.Key + " " + v.TypeVal
		if v.Star {
			fullType = v.Key + " *" + v.TypeVal
		}
		paramStr = append(paramStr, fullType)
	}
	for _, v := range data.Returns {
		var fullType = v.Key + " " + v.TypeVal
		if v.Star {
			fullType = v.Key + " *" + v.TypeVal
		}
		returnStr = append(returnStr, fullType)
	}
	d := comData{
		FuncData:  data,
		ReturnStr: strings.Join(returnStr, ","),
		ParamStr:  strings.Join(paramStr, ","),
	}
	if len(data.Returns) == 2 {
		d.ReturnStruct = data.Returns[0]
	}
	method := "Get"
	if strings.Contains(data.Name, "Post") {
		method = "Post"
	} else if strings.Contains(data.Name, "Del") {
		method = "Delete"
	} else if strings.Contains(data.Name, "Put") {
		method = "Put"
	}
	d.Method = method
	for _, v := range data.Params {
		if strings.Contains(v.TypeVal, "context.Context") {
			continue
		}
		fmt.Println(v)
		if strings.Contains(v.TypeVal, "HeaderParam") {
			d.Headers = append(d.Headers, v)
		} else if strings.Contains(v.TypeVal, "QueryParam") {
			d.Query = append(d.Query, v)
		} else if strings.Contains(v.TypeVal, "PathParam") {
			d.PathParams = append(d.PathParams, v)
		} else if method != "Get" && (strings.Contains(v.TypeVal, "BodyParam") || !isSimpleType(v.TypeVal)) {
			d.Body = v
		}
	}
	t, err := template.New("tpl").Parse(tpl)
	if err != nil {
		panic(err)
	}
	if err := t.Execute(dst, d); err != nil {
		panic(err)
	}
}

func isSimpleType(typeVal string) bool {
	return (!strings.Contains(typeVal, "map") &&
		!strings.Contains(typeVal, "[]") &&
		!strings.Contains(typeVal, "*")) &&
		(strings.Contains(typeVal, "string") ||
			strings.Contains(typeVal, "int") ||
			strings.Contains(typeVal, "float64") ||
			strings.Contains(typeVal, "bool"))
}

type HeadData struct {
	Line     string
	PkgName  string
	FileName string
}

func FuncHead(dst io.Writer, data HeadData) {
	const tpl = `
// Code generated by file {{.FileName}} and line {{.Line}} DO NOT EDIT.
// For more detail see https://github.com/wanglihui/restyless
package {{.PkgName}}
`
	t, err := template.New("tpl").Parse(tpl)
	if err != nil {
		panic(err)
	}
	if err := t.Execute(dst, data); err != nil {
		panic(err)
	}
}

type StructData struct {
	TypeName string
}

func StructTemplate(dst io.Writer, data StructData) {
	const tpl = `
func New{{.TypeName}}Impl(r *resty.Client) {{.TypeName}} {
	return &{{.TypeName}}Impl{
		r : r,
	}
}

type {{.TypeName}}Impl	struct {
	r *resty.Client
}
`
	t, err := template.New("tpl").Parse(tpl)
	if err != nil {
		panic(err)
	}
	if err := t.Execute(dst, data); err != nil {
		panic(err)
	}
}
