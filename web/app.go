package main

import (
	"fmt"
	"git.lcc.lib/orm"

	// "math/rand"
	"os"
	"time"
)

func main() {
	dbQuery := orm.NewOrmQuery("mysql", "root:@tcp(127.0.0.1:3306)/test?charset=utf8")

	//dbQuery.Table("user")
	user := dbQuery.Table("user").Where("id", "2147483647").GetRow()

	//list := dbQuery.Clear("where", "table").Table("user").OrderBy("id", orm.ASC).Field("id").GetList(0, -1)

	cls := dbQuery.Clear().Table("user").OrderBy("id", orm.ASC).GroupBy("id").Having("id=1").Field("id").GetColumn(0, -1)

	nrow := dbQuery.Clear().Table("user").Where("id", 5, orm.OP_GT).Delete()

	nid := dbQuery.Clear().Table("User").Insert(map[string]interface{}{"id": 0, "user": "leichenchun", "age": 30, "ip": "192.168.19.53", "stime": time.Now()})

	arow := dbQuery.Clear().Table("User").Where("id", 2).Update(map[string]interface{}{"user": "lchenchun", "age": 29, "ip": "192.168.19.79", "stime": time.Now()})

	aval := dbQuery.Clear().Table("User").Where("id", 2).Field("user").GetValue()
	fmt.Println(user, cls, nrow, aval, nid, arow, dbQuery)
	os.Exit(0)

}
