## sqlt的历史和由来

java的数据库访问工具MyBatis给大家留下了深刻的印象，早在几年前，我刚刚接触golang的时候，也希望golang也有类似的工具，对golang稍微熟悉后，发现golang自带模板功能（text/template），于是在另外一个开源库sqlx的基础上，增加模板拼接sql的功能，所以 sqlt 就诞生了。

最早sqlt是托管在oschina，后续迁移到github.com/it512/sqlt，之后加入twiglab

## 安装

```
go get github.com/twiglab/sqlt
```

sqlt 也支持 go mod

## sqlt 架构简要说明

sqlt 深度依赖sqlx， 是在sqlx的基础上增加了模板功能，底层的数据库方法全部通过sqlx的NamedStmt和PrepareName完成对数据库的访问。

sqlt 对外提供的所有操作全部通过 `Dbop` struct 提供，Dbop struct 组合了sqlx.DB和模板（由Maker接口定义）

```go
type Dbop struct {
	Maker
	*sqlx.DB
}
```

- sqlt 没有隐藏任何使用sqlx的细节，Dbop对外直接暴露sqlx.DB，任何sqlt.DB的方法均可以直接使用，请参考sqlx的文档 
- sqlt 全部采用Parper和NamedStmt 完成对数据库的访问，所以也受到数据库驱动的限制，请详细参考数据库驱动的文档

目前sqlt自带的模板为`text/template`，任何Maker接口的实现都可以作为sqlt的模板使用。
```go
type Maker interface {
	MakeSql(string, interface{}) (string, error)
}
```
（*golang 自带的模板未必最好，欢迎pr更好模板实现*）

## 使用说明

### Dbop的创建
最简单的创建方式为：

```go
dbop := sqlt.Default("postgres", "dbname=testdb sslmode=disable", "tpl/*.tpl")
```
**注意：不要忘记引入数据库驱动**

如果你有现成的数据库链接，或者对模板有特殊的要求，有可以用使用`sqlt.New`方法创建

```go
dbx := sqlt.MustConnect("postgres", "dbname=testdb sslmode=disable")
tpl := sqlt.NewSqlTemplate("tpl/*.tpl")
tpl.SetDebug(true)

dbop := sqlt.New(dbx, tpl)
```

如果你的模板方法中用到了自定义的函数，sqlt也提供了一个 `NewSqlTemplateWithFuncs` 的方法用于创建带有自定义函数的模板 （位于`tpl.go`中）

### 模板

再次说明sqlt默认自带的模板是text/template的封装实现，详细的用法请参考text/template

(*example目录中有完整的例子*)

### sqlt的方法

最新版本中，所有的sqlt的方法都可以直接调用，分别为：

- func Query(execer TExecer, ctx context.Context, id string, param interface{}, h RowsExtractor) (err error) 
- func Exec(execer TExecer, ctx context.Context, id string, param interface{}) (r sql.Result, err error) 

以及对应的Must版本

Texcer 为 `*Dbop`，所有的方法都支持`context.Context`，id 为模板的id，param为传递给模板和用于Prepare的参数（用于构建条件，和sql拼接）

Exec方法对应执行无返回的sql语句，如：insert， update， delete和存储过程。

Query方法用于执行带有返回结果的sql语句，如：select，和带有returnning子句的insert， update （*returnning子句需要数据库和驱动的支持*）

#### Query 的结果集处理

sqlt 提供了 `RowsExtractor` 接口处理结果集，Query的最后一个参数中传入RowsExtractor的实现。

例子：
```go
type Staff struct {
	StaffId   int       `db:"staff_id"`   //
	StaffName string    `db:"staff_name"` //
	CreatedAt time.Time `db:"created_at"` //
	UpdatedAt time.Time `db:"updated_at"` //
	Age       int       `db:"age"`        //
}

type StaffHandler struct {
	Staffs []*Staff
}

func (sh *StaffHandler) Extract(rs sqlt.Rows) (err error) {
	for rs.Next() {
		staff := new(Staff)
		if err = rs.StructScan(staff); err != nil {
			return
		}

		sh.Staffs = append(sh.Staffs, staff)
	}

	return rs.Err()
}

staff := new(Staff)
staff.StaffId = 67890
h := new(StaffHandler)
sqlt.MustQuery(dbop, context.Background(), "Staff.select", staff, h)
```

**sqlt *不要求*返回字段和struct字段一一对应，struct映射是按照row和struct公共部分映射的**


staff 最为查询条件传入模板，模板会根据staff字段构建查询条件，然后通过sqlx执行查询，返回结果由RowsExtractor的实现StaffHandler处理
模板：
```template
{{ define "Staff.select"}}
select
	staff_id,
	staff_name,
	created_at
from
	Staff
where
	{{if .StaffId}} staff_id = :staff_id {{end}}
	{{if .StaffName}}and  staff_name = :staff_name {{end}}
{{end}}
```

## Gen （代码生成器）

cmd目录下的pggen是postgresql数据库的代码生成器，目前sqlt只提供了pg代码生成

首先需要说明的是：
- 生成器以表为单位生成对应的struct，常量，和增删改查通用的模板
- 生成器生成的内容直接输出到屏幕上（stdout），需要重定向到文件
- 生成器生成的代码片段，并不是可以直接使用的结果

**生成器生成的是代码片段，并不是可以直接使用的结果**

开发人员需要从生成结果中复制有用的片段到程序中

## 一些规则

这些规则不是强制的，但是如果不遵循这些规则，使用sqlt的工作量会增大，从而失去价值

- 数据库对象（表，字段）的命名，采用下划线格式（_) 如： staff_name, staff_id
- struct 中的命名采用golang推荐的snake风格，如：StaffName, StaffId
- 每个表里面最好都加上 created_at和updated_at这2个字段，用于记录创建时间和修改时间（后面生成模板的时候会简单一些）
- 模板里面，逗号，and 都放在字段前面
- 模板的名称要能顾名思义
- sqlt底层用为sqlx，所以struct的处理也必须符合sqlx的要求（增加db tag）

第一条和第二天规则方便生成和struct和row直接的映射
第三条和第四条方便模板的编写，确保不会生成没有字段的错误sql

## sqlt已知的一些问题和避免方法

较之 Mybatis，由于golang和模板的限制，sqlt存在下列问题
- text/template是没有上下文的， 所以无法帮助你处理结尾逗号(,)问题，避免的方法参见规则3,4
- 0值问题，golang默认0值，对于0值模板会排除，需要用自定义函数去处理

## 特别提醒

**模板只是帮您拼接了sql， sql是否正确，以及效率需要开发人员保证。**

## 版权
MIT
