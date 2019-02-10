## sqlt的历史和由来

java的数据库访问工具MyBatis给大家留下了深刻的印象，早在几年前，我刚刚接触golang的时候，也希望golang也有类似的工具，对golang稍微熟悉后，发现golang自带模板功能（text/template），于是在另外一个开源库sqlx的基础上，增加模板拼接sql的功能，所以 sqlt 就诞生了。

## 安装

```
go get github.com/twiglab/sqlt
```

sqlt 也支持 go mod

## sqlt 架构简要说明

sqlt 深度依赖sqlx， 是在sqlx的基础上增加了模板功能，底层的数据库方法全部通过sqlx的NamedStmt和PrepareName完成对数据库的访问。

sqlt 对外提供的所有操作全部通过 `Dbop` struct 提供，Dbop struct 组合了sqlx.DB和模板（由Maker）

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

