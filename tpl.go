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
	logger *log.Logger
}

func NewSqlTemplate(pattern string) *SqlTemplate {
	tpl := template.New("sqlt-template").Funcs(make(template.FuncMap))
	tpl = template.Must(tpl.ParseGlob(pattern))
	return &SqlTemplate{
		Template: tpl,
		Debug:    false,
		logger:   log.New(os.Stdout, "sqlt-std-template-", log.LstdFlags|log.Llongfile),
	}
}

func (t *SqlTemplate) MakeSql(id string, param interface{}) (string, error) {
	sb := new(strings.Builder)
	err := t.ExecuteTemplate(sb, id, param)

	if t.Debug && err == nil {
		t.logger.Println(sb)
	}

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
