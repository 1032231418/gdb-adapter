package main

import (
	"fmt"
	"github.com/gogf/gf/database/gdb"
)

func main() {
	gdb.SetConfig(gdb.Config{
		"casbin": gdb.ConfigGroup{
			gdb.ConfigNode{
				Host:     "127.0.0.1",
				Port:     "3306",
				User:     "root",
				Pass:     "gdkid,,..",
				Name:     "gf-vue-admin",
				Type:     "mysql",
				Role:     "master",
				Weight:   100,
			},
		},
	})
	db, err := gdb.New("casbin")
	if err != nil {
		panic(err)
	}
	db.SetDebug(true)
	r , err := db.Table("jwts").FindOne()
	if err != nil {
		panic(err)
	}
	fmt.Println(r)
	if r, err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS `%v`;", "jwts")); err != nil { // 判断表是否存在
		fmt.Println(err)
		panic(err)
	}else {
		i1, _ := r.RowsAffected()
		i2, _ := r.LastInsertId()
		fmt.Println(i1, i2)
	}
}
