package mysql2

import "bytes"

//辅助生成类似  in(?,?,?,?) 批量查询的sql
func makeBatchSelectStr(data []string) (string, []interface{}) {
	length := len(data)
	if length == 0 {
		return "", nil
	}

	params := make([]interface{}, 0, length)

	sqlStringBuffer := bytes.Buffer{}
	sqlStringBuffer.WriteString("(")

	for k, v := range data {
		params = append(params, v)
		if length == k+1 {
			sqlStringBuffer.WriteString("?")
		} else {
			sqlStringBuffer.WriteString("?,")
		}
	}
	sqlStringBuffer.WriteString(")")

	return sqlStringBuffer.String(), params

}
