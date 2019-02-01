package sqlt

import (
	"log"
	"os"
	"strings"
	"text/template"
)

type SqlTemplate struct {
	*template.Template
	Debug  bool
	Logger *log.Logger
}

func NewSqlTemplate(pattern string) *SqlTemplate {
	return NewSqlTemplateWithFuncs(pattern, make(template.FuncMap))
}

func NewSqlTemplateWithFuncs(pattern string, funcs template.FuncMap) *SqlTemplate {
	tpl := template.Must(template.New("sqlt-template").Funcs(funcs).ParseGlob(pattern))
	return &SqlTemplate{
		Template: tpl,
		Debug:    false,
		Logger:   log.New(os.Stdout, "sqlt-std-maker-", log.LstdFlags|log.Lshortfile),
	}
}

func (t *SqlTemplate) MakeSql(id string, param interface{}) (string, error) {
	sb := new(strings.Builder)
	err := t.ExecuteTemplate(sb, id, param)

	if t.Debug && err == nil {
		t.Logger.Println(sb)
	}

	return sb.String(), err
}

func (t *SqlTemplate) SetDebug(b bool) {
	t.Debug = b
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
