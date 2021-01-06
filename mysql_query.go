package mysql2

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
)

type QueryRes struct {
	res []map[string]interface{}
	err error
}

func NewQueryRes(res []map[string]interface{}, err error) *QueryRes {
	q := new(QueryRes)
	q.res = res
	q.err = err
	return q
}
func (r *QueryRes) Error() error {
	return r.err
}
func (r *QueryRes) Err() error {
	return r.err
}

// 还原为数组
func (r *QueryRes) Data() []map[string]interface{} {
	return r.res
}

// 用函数遍历内部的数据, 用来修改自己本身
func (r *QueryRes) Map(f func(map[string]interface{})) {
	for _, v := range r.res {
		f(v)
	}
}

// 过滤掉一部分数据
func (r *QueryRes) Erase(f func(map[string]interface{}) bool) *QueryRes {
	newRes := make([]map[string]interface{}, 0, 10)
	for _, v := range r.res {
		if !f(v) {
			newRes = append(newRes, v)
		}
	}
	r.res = newRes
	return r
}

// 保留一部分数据
func (r *QueryRes) Keep(f func(map[string]interface{}) bool) *QueryRes {
	newRes := make([]map[string]interface{}, 0, 10)
	for _, v := range r.res {
		if f(v) {
			newRes = append(newRes, v)
		}
	}
	r.res = newRes
	return r
}

/********************************************************************/
// 执行exec   参数是*DB  or *DbTx
func (s SQL) Query(source interface{}) *QueryRes {
	if s.str == "" {
		return NewQueryRes(nil, errors.New("blank sql"))
	}
	res, err := queryCommon(source, s.str, s.params)
	return NewQueryRes(res, err)
}

// 统一处理事务内，和非事务内query
func queryCommon(source interface{}, sqlStr string, args []interface{}) ([]map[string]interface{}, error) {

	if Conf.Log {
		fmt.Println("running....query sql = ", sqlStr, "\n args=", args)
	}
	p, ok := source.(*DB)
	if ok {
		if p.err != nil {
			return nil, p.err
		}
		return queryResFromRows(p.realPool.Query(sqlStr, args...))
	}
	// 多个sql事务
	t, ok := source.(*DBTx)
	if ok {
		if t.err != nil {
			return nil, t.err
		}
		return queryResFromRows(t.realtx.Query(sqlStr, args...))
	}
	return nil, errors.New("only support DbPool , DbTx")
}

//  scan的行为null 对应nil  数字对数字  其他对字符串 ,所以所有的字段数据类型归结为简单的几种。这可能不能处理非常规情况。
// 联表查询，如果两个表中有同名字段的时候，不会报错，会忠实的输出数据
// 另外如果数据库里是null，怎会被转换成0，空字符串，可能会影响业务逻辑，需要开发者自己注意
func queryResFromRows(rows *sql.Rows, e error) ([]map[string]interface{}, error) {
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	res := make([]map[string]interface{}, 0, 100)
	columns, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}
	fields, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	lengthRow := len(fields)
	for {
		if result := rows.Next(); result {
			scanRes := make([]sql.Scanner, lengthRow)
			for i := 0; i < lengthRow; i++ {
				vType := columns[i].DatabaseTypeName()

				switch vType {
				case "INT", "BIGINT", "TINYINT", "MEDIUMINT":
					scanRes[i] = &sql.NullInt64{}
				case "FLOAT", "DOUBLE":
					scanRes[i] = &sql.NullFloat64{}
				case "CHAR", "VARCHAR", "TIME", "TEXT", "DECIMAL", "BLOB", "GEOMETRY", "BIT", "DATETIME", "DATE", "TIMESTAMP":
					scanRes[i] = &sql.NullString{}
				default:
					scanRes[i] = &sql.NullString{}
				}
			}
			resultData := make(map[string]interface{}, lengthRow)
			vScanRes := reflect.ValueOf(&scanRes)
			fn := reflect.ValueOf(rows.Scan)
			fnParams := make([]reflect.Value, lengthRow)
			for i := 0; i < lengthRow; i++ {
				fnParams[i] = vScanRes.Elem().Index(i)
			}
			callResult := fn.Call(fnParams)
			if callResult[0].Interface() != nil {
				return nil, callResult[0].Interface().(error)
			}
			for i := 0; i < lengthRow; i++ {
				resultData[fields[i]] = fetchFromScanner(scanRes[i])
			}
			res = append(res, resultData)
		} else {
			break
		}
	}
	return res, nil
}

// 把诸如*sql.NullXX  转化为正常的XX值，null一般转化为XX的零值
// 所以设计数据库的时候，要注意这套代码实际上是无法区分null和该字段类型的零值的
func fetchFromScanner(data sql.Scanner) interface{} {
	switch v := data.(type) { // v表示b1 接口转换成Bag对象的值
	case *sql.NullInt64:
		if v.Valid {
			return v.Int64
		}
		return int64(0)

	case *sql.NullFloat64:
		if v.Valid {
			return v.Float64
		}
		return float64(0)

	case *sql.NullString:
		if v.Valid {
			return v.String
		}
		return ""

	default:
		// 不可能会运行到这里
		return nil
	}
}
