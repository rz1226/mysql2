package mysql2

import (
	"bytes"
	"sort"
	"strings"
)

func (q *QueryRes) ToInsertSQL(tableName string, fieldsExclude map[string]int) SQL {
	if q.err != nil {
		return NewSQL("")
	}
	return sqlFromDatasForInsert(q.Data(), tableName, fieldsExclude)
}

func (q *QueryRes) ToUpdateSQL(tableName, condition string, updateFields map[string]int) SQL {
	if q.err != nil {
		return NewSQL("")
	}
	return sqlFromDataForUpdate(q.Data()[0], tableName, condition, updateFields)
}

// 第二个参数,第三个参数是指保留或者排除的字段，include=true为保留，false为排除，如果是nil，都保留, 用map格式，是因为查找键更方便快速,值无用
func _sqlFromQueryRes(data map[string]interface{}, fields map[string]int, include bool) (strFieldList, strInsert string, dataParams []interface{}, strUpdate string) {
	insertFieldList := "("
	insertMarksStr := "("
	updateStr := ""
	insertValuesSli := make([]interface{}, 0, 30)
	// type lineData map[string]*fieldData
	//map用前要排序，否则会出错，导致多个insert生成的时候，字段名和参数对不上
	keys := make([]string, 0, len(data))
	for k, _ := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := data[k]
		// 略过过滤
		lengthOfFields := len(fields)
		if lengthOfFields > 0 {
			_, ok := fields[k]
			if ok {
				if include == false {
					continue
				}
			} else {
				if include == true {
					continue
				}
			}
		}
		insertFieldList += "`" + k + "`" + ","
		insertMarksStr += "?,"
		insertValuesSli = append(insertValuesSli, v)
		updateStr += "`" + k + "` = ?" + ","

	}
	return strings.TrimRight(insertFieldList, ",") + ")",
		strings.TrimRight(insertMarksStr, ",") + ")",
		insertValuesSli, strings.TrimRight(updateStr, ",")
}

// 生成一个insert语句
func sqlFromDataForInsert(data map[string]interface{}, tableName string, fieldsExclude map[string]int) SQL {
	insertFields, insertMarks, insertParams, _ := _sqlFromQueryRes(data, fieldsExclude, false)
	insertSQL := "insert into " + tableName + "  " + insertFields + " values " + insertMarks
	return NewSQL(insertSQL, insertParams...)
}

// 生成一个update语句
func sqlFromDataForUpdate(data map[string]interface{}, tableName, condition string, updateFields map[string]int) SQL {
	_, _, insertParams, updateStr := _sqlFromQueryRes(data, updateFields, true)
	var insertSQL string
	if strings.Trim(condition, " ") == "" {
		insertSQL = "update " + tableName + " set " + updateStr
	} else {
		insertSQL = "update " + tableName + " set " + updateStr + " where " + condition
	}

	return NewSQL(insertSQL, insertParams...)
}

func sqlFromDatasForInsert(data []map[string]interface{}, tableName string, fieldsExclude map[string]int) SQL {
	length := len(data)
	if length == 0 {
		return NewSQL("")
	}
	var marksBuf bytes.Buffer
	insertFields, insertMarks, insertParams, _ := _sqlFromQueryRes(data[0], fieldsExclude, false)
	marksBuf.WriteString(insertMarks)
	marksBuf.WriteString(",")
	for i := 1; i < length; i++ {
		_, marks, params, _ := _sqlFromQueryRes(data[i], fieldsExclude, false)
		marksBuf.WriteString(marks)
		marksBuf.WriteString(",")
		insertParams = append(insertParams, params...)
	}
	insertSQL := "insert into " + tableName + "  " + insertFields + " values " + strings.Trim(marksBuf.String(), ",")
	return NewSQL(insertSQL, insertParams...)
}

func concatStr(strs ...string) string {
	var buf bytes.Buffer
	for _, v := range strs {
		buf.WriteString(v)
	}
	return buf.String()
}
