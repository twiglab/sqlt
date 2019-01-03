package sqlt

import (
	"io"
	"strings"
	"text/template"
)

type (
	StdTemplateRender struct {
		pattern string
		funcMap template.FuncMap
		t       *template.Template
	}
)

func NewStdTemplateRenderDefault(pattern string) *StdTemplateRender {
	return NewStdTemplateRender(pattern, make(template.FuncMap))
}

func NewStdTemplateRender(pattern string, funcMap template.FuncMap) *StdTemplateRender {
	tpl := template.New("sqlt-template").Funcs(funcMap)
	tpl = template.Must(tpl.ParseGlob(pattern))
	return &StdTemplateRender{pattern: pattern, funcMap: funcMap, t: tpl}
}

func (st *StdTemplateRender) Render(w io.Writer, id string, param interface{}) error {
	return st.t.ExecuteTemplate(w, id, param)
}

type (
	StdSqlAssembler struct {
		Render *StdTemplateRender
	}
)

func (l *StdSqlAssembler) Sql(id string, data interface{}) (sql string, err error) {
	bs := new(strings.Builder)
	if err = l.Render.Render(bs, id, data); err != nil {
		sql = bs.String()
	}

	return
}

func (l *StdSqlAssembler) MustSql(id string, data interface{}) (sql string) {
	var err error
	if sql, err = l.Sql(id, data); err != nil {
		panic(err)
	}

	return
}

func NewStdSqlAssemblerDefault(pattern string) *StdSqlAssembler {
	return &StdSqlAssembler{
		Render: NewStdTemplateRenderDefault(pattern),
	}
}
