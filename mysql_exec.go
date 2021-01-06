package mysql2

import (
	dsql "database/sql"
	"errors"
	"fmt"
	"strings"
)

// 执行exec   参数是*DB  or *DbTx
func (s SQL) Exec(source interface{}) (int64, error) {
	if s.str == "" {
		return 0, errors.New("blank sql")
	}
	n, err := execCommon(source, s.str, s.params)
	return n, err
}

func execCommon(source interface{}, sqlStr string, args []interface{}) (int64, error) {

	if Conf.Log {
		fmt.Println("running.... exec sql = ", sqlStr, "\n args=", args)
	}
	p, ok := source.(*DB)
	if ok {
		if p.err != nil {
			return 0, p.err
		}
		result, err := p.realPool.Exec(sqlStr, args...)
		if err != nil {
			return int64(0), err
		}
		return affectedResult(sqlStr, result)
	}
	t, ok := source.(*DBTx)
	if ok {
		if t.err != nil {
			return 0, t.err
		}
		result, err := t.realtx.Exec(sqlStr, args...)
		if err != nil {
			return int64(0), err
		}
		return affectedResult(sqlStr, result)
	}
	return int64(0), errors.New("only support DbPool , DbTx")
}

// 从exec的result获取   当insert获取最后一个id， update，delete获取影响行数，replace获取最后一个id
func affectedResult(sqlStr string, result dsql.Result) (int64, error) {
	if isSQLUpdate(sqlStr) || isSQLDelete(sqlStr) {
		return result.RowsAffected() // 本身就是多个返回值
	}
	if isSQLInsert(sqlStr) {
		return result.LastInsertId() // 本身就是多个返回值
	}
	if isSQLReplace(sqlStr) {
		return result.LastInsertId() // 本身就是多个返回值
	}
	return int64(0), errors.New("only support update insert delete replace")
}

func isSQLReplace(sqlStr string) bool {
	str := strings.TrimSpace(strings.ToLower(sqlStr))
	return strings.HasPrefix(str, "replace")
}
func isSQLInsert(sqlStr string) bool {
	str := strings.TrimSpace(strings.ToLower(sqlStr))
	return strings.HasPrefix(str, "insert")
}

func isSQLUpdate(sqlStr string) bool {
	str := strings.TrimSpace(strings.ToLower(sqlStr))
	return strings.HasPrefix(str, "update")
}

func isSQLDelete(sqlStr string) bool {
	str := strings.TrimSpace(strings.ToLower(sqlStr))
	return strings.HasPrefix(str, "delete")
}
