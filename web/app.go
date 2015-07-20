package main

import (
	"fmt"
	"git.lcc.lib/core"
	"git.lcc.lib/orm"
	"os"
	"sync"
	"time"
)

var wg sync.WaitGroup

func main() {
	idcard := "3522251986111620141"
	fmt.Println(core.IsIdCard(idcard))
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	wg.Add(1)
	go dbTest("1111")
	wg.Add(1)
	go dbTest("2222")
	wg.Wait()
	os.Exit(0)
}

func dbTest(seg string) {
	fmt.Println("-----------------------" + seg + "------------------------------")
	db := orm.DB("DB-TEST")
	//dbQuery.Table("user")
	user := db.Table("user").Where("id", "2147483647").GetRow()
	//list := dbQuery.Clear("where", "table").Table("user").OrderBy("id", orm.ASC).Field("id").GetList(0, -1)

	cls := db.Clear().Table("user").OrderBy("id", orm.ASC).GroupBy("id").Having("id=1").Field("id").GetColumn(0, -1)

	nid := db.Clear().Table("User").Insert(map[string]interface{}{"id": 0, "user": "leichenchun", "age": 30, "ip": "192.168.19.53", "stime": time.Now()})

	aval := db.Clear().Table("User").Where("id", 2).Field("user").GetValue()

	db.GetParse().Begin()

	nrow := db.Clear().Table("user").Where("id", 5, orm.OP_GT).Delete()

	//nid := dbQuery.Clear().Table("User").Insert(map[string]interface{}{"id": 0, "user": "leichenchun", "age": 30, "ip": "192.168.19.53", "stime": time.Now()})

	arow := db.Clear().Table("User").Where("id", 2).Update(map[string]interface{}{"user": "lchenchun", "age": 29, "ip": "192.168.19.79", "stime": time.Now()})
	db.GetParse().Commit()

	fmt.Println(user, cls, nrow, aval, nid, arow, db)
	fmt.Println("-----------------------end " + seg + "------------------------------")
	wg.Done()
}
