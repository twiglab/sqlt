package sqlt

import (
	"strings"
	"text/template"
)

type SqlTemplate struct {
	*template.Template
}

func NewSqlTemplate(pattern string, funcMap template.FuncMap) *SqlTemplate {
	tpl := template.New("sqlt-template").Funcs(funcMap)
	tpl = template.Must(tpl.ParseGlob(pattern))
	return &SqlTemplate{Template: tpl}
}

func (t *SqlTemplate) MakeSql(id string, param interface{}) (string, error) {
	sb := new(strings.Builder)
	err := t.ExecuteTemplate(sb, id, param)
	return sb.String(), err
}

type Maker interface {
	MakeSql(string, interface{}) (string, error)
}

func Sql(maker Maker, id string, param interface{}) (string, error) {
	return maker.MakeSql(id, param)
}

func MustSql(maker Maker, id string, param interface{}) (sql string) {
	var e error
	if sql, e = Sql(maker, id, param); e != nil {
		panic(e)
	}

	return
}
