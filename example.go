package mysql2

/**
package main

import(
	"fmt"
	"github.com/rz1226/mysql2"
	"os"
	"strings"
)

var DB *mysql2.DB

func init(){
	var MYSQL_HOST = strings.TrimSpace(os.Getenv("XX_MYSQL_HOST"))
	var MYSQL_PORT = strings.TrimSpace(os.Getenv("XX_MYSQL_PORT"))
	var MYSQL_USERNAME = strings.TrimSpace(os.Getenv("XX_MYSQL_USERNAME"))
	var MYSQL_PASSWORD = strings.TrimSpace(os.Getenv("XX_MYSQL_PASSWORD"))

	conf := mysql2.NewDBConf(MYSQL_USERNAME,MYSQL_PASSWORD,MYSQL_HOST,MYSQL_PORT,"fenxi",5)
	fmt.Println(conf.Str())
	DB = mysql2.Connect( conf )
}

type Data struct{
	Id int64 `orm:"id" auto:"1"`
	Name string `orm:"order_name"`
}

func main(){
	sql := mysql2.NewSQL("select * from fenxi where id in (?,?) limit 100 ", 1,2 )
	var data []*Data
	err := sql.Query(DB).ToStruct(&data)

	fmt.Println(err )
	for _, v := range data{
		fmt.Println(v.Name)
	}
	fmt.Println(sql.Query(DB).Data())
	fmt.Println(sql.Query(DB).ToInsertSQL("fenxi",map[string]int{"id":1}).Info())
}


*/
