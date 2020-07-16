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
				Name:     "casbin",
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
	tables, err := db.Tables("casbin")
	if err != nil {
		panic(err)
	}
	fmt.Println(tables)
}
