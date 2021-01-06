package mysql2

import "fmt"

// 代表一个可执行的sql的字符串部分
type SQLStr string

// 生成完整的sql
func (ss SQLStr) AddParams(params ...interface{}) SQL {
	sql := SQL{}
	sql.str = string(ss)
	if len(params) > 0 {
		sql.params = params
	} else {
		sql.params = make([]interface{}, 0)
	}

	return sql
}

// 不需要参数直接query
/*
func (s Sql) Query(source interface{}) *QueryRes {
	res, error := queryCommon(source, string(s.str), s.params)
	return NewQueryRes(res, error)
}
*/
func (ss SQLStr) Query(source interface{}) *QueryRes {
	return ss.AddParams().Query(source)
}

// 不需要参数直接exec
func (ss SQLStr) Exec(source interface{}) (int64, error) {
	return ss.AddParams().Exec(source)
}

// 代表一个可以执行的sql，一般由两部分组成，str，和变量
type SQL struct {
	str    string
	params []interface{}
}

func NewSQL(str string, params []interface{}) SQL {
	sql := SQL{}
	sql.str = str
	sql.params = params
	return sql
}

// 补上一个条件
func (s SQL) ConcatSQL(s2 SQL) SQL {
	res := NewSQL(s.str, s.params[:])
	res.str += s2.str
	res.params = append(res.params, s2.params...)
	return res
}

//补上一个 where in 语句
func (s SQL) In(key string, params []string) SQL {
	str, args := makeBatchSelectStr(params)

	sql2 := NewSQL(" where `"+key+"`"+" in "+str+" ", args)
	sql := s.ConcatSQL(sql2)
	return sql
}
func (s SQL) AndIn(key string, params []string) SQL {
	str, args := makeBatchSelectStr(params)

	sql2 := NewSQL(" and `"+key+"`"+" in "+str+" ", args)
	sql := s.ConcatSQL(sql2)
	return sql
}

func (s SQL) clone() SQL {
	return NewSQL(s.str, s.params[:])
}

func (s SQL) Limit(limit int) SQL {
	sql := s.clone()
	sql.str += " limit " + fmt.Sprint(limit)
	return sql
}

func (s SQL) Offset(offset int) SQL {
	sql := s.clone()
	sql.str += " offset " + fmt.Sprint(offset)
	return sql
}

func (s SQL) OrderBy(order string) SQL {
	sql := s.clone()
	sql.str += " order by " + fmt.Sprint(order)
	return sql
}
func (s SQL) Info() string {
	str := fmt.Sprint("str= ", s.str, "\n params=", s.params)
	return str
}
