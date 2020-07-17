// Copyright 2020 The casbin Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"github.com/gogf/gf/database/gdb"
	"runtime"

	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
)

type CasbinRule struct {
	PType string `json:"ptype"`
	V0    string `json:"v0"`
	V1    string `json:"v1"`
	V2    string `json:"v2"`
	V3    string `json:"v3"`
	V4    string `json:"v4"`
	V5    string `json:"v5"`
}

// Adapter represents the gdb adapter for policy storage.
// Adapter 代表用于策略存储的gdb适配器。
type Adapter struct {
	driverName     string
	dataSourceName string
	dbSpecified    bool
	tableName      string
	db             gdb.DB
	isFiltered     bool
}

// finalizer is the destructor for Adapter.
// finalizer是Adapter的析构函数。
func finalizer(a *Adapter) {
	// 注意不用的时候不需要使用Close方法关闭数据库连接(并且gdb也没有提供Close方法)，
	// 数据库引擎底层采用了链接池设计，当链接不再使用时会自动关闭
	a.db = nil
}

// NewAdapter is the constructor for Adapter.
// dbSpecified is an optional bool parameter. The default value is false.
// It's up to whether you have specified an existing DB in dataSourceName.
// If dbSpecified == true, you need to make sure the DB in dataSourceName exists.
// If dbSpecified == false, the adapter will automatically create a DB named "casbin".

// NewAdapter是Adapter的构造函数。
// dbSpecified是可选的bool参数。 默认值为false。
// 如果dbSpecified == true，则需要确保dataSourceName中的数据库存在。
// 如果dbSpecified == false，则适配器将自动创建一个名为“ casbin”的数据库。
func NewAdapter(driverName string, dataSourceName string) (*Adapter, error) {
	a := &Adapter{}
	a.driverName = driverName
	a.dataSourceName = dataSourceName
	a.tableName = "casbin_rule"

	// Open the DB, create it if not existed.
	if err := a.open(); err != nil {
		return nil, err
	}

	// Call the destructor when the object is released.
	runtime.SetFinalizer(a, finalizer)

	return a, nil
}

// NewAdapterByDB is the constructor for Adapter.Need to pass in gdb.DB
// NewAdapterByDB 是Adapter的构造函数,需要传入gdb.DB
func NewAdapterByDB(db gdb.DB) (*Adapter, error) {
	gdb.New()
	a := &Adapter{
		db: db,
	}
	if err := a.createTable(); err != nil {
		return nil, err
	}
	return a, nil
}

// NewAdapterFromOptions is the constructor for Adapter with existed connection
// NewAdapterFromOptions 适配器的构造函数是否具有已存在的连接
func NewAdapterFromOptions(adapter *Adapter) (*Adapter, error) {

	if adapter.tableName == "" {
		adapter.tableName = "casbin_rule"
	}
	if adapter.db == nil {
		err := adapter.open()
		if err != nil {
			return nil, err
		}

		runtime.SetFinalizer(adapter, finalizer)
	}
	return adapter, nil
}

func (a *Adapter) open() error {
	var err error
	var db gdb.DB

	gdb.SetConfig(gdb.Config{
		"casbin": gdb.ConfigGroup{
			gdb.ConfigNode{
				Type:     a.driverName,
				LinkInfo: a.dataSourceName,
				Role:     "master",
				Weight:   100,
			},
		},
	})
	db, err = gdb.New("casbin")

	if err != nil {
		return err
	}

	a.db = db

	return a.createTable()
}

// close
func (a *Adapter) close() error {
	a.db = nil // 注意不用的时候不需要使用Close方法关闭数据库连接(并且gdb也没有提供Close方法)，数据库引擎底层采用了链接池设计，当链接不再使用时会自动关闭
	return nil
}

// createTable Create a data table
// createTable 创建数据表
func (a *Adapter) createTable() error {
	tables, err := a.db.Tables(a.dataSourceName)
	if err != nil {
		panic(err)
	}
	for _, v := range tables {
		if v == a.tableName {
			return nil
		}
	}
	_, err = a.db.Exec(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (ptype VARCHAR(10), v0 VARCHAR(256), v1 VARCHAR(256), v2 VARCHAR(256), v3 VARCHAR(256), v4 VARCHAR(256), v5 VARCHAR(256))", a.tableName))
	return err
}

// dropTable Delete table
// dropTable 删除表
func (a *Adapter) dropTable() error {
	_, err := a.db.Exec(fmt.Sprintf("DROP TABLE %s", a.tableName))
	return err
}

// loadPolicyLine
func loadPolicyLine(rule CasbinRule, model model.Model) {
	ruleText := rule.PType
	if rule.V0 != "" {
		ruleText += ", " + rule.V0
	}
	if rule.V1 != "" {
		ruleText += ", " + rule.V1
	}
	if rule.V2 != "" {
		ruleText += ", " + rule.V2
	}
	if rule.V3 != "" {
		ruleText += ", " + rule.V3
	}
	if rule.V4 != "" {
		ruleText += ", " + rule.V4
	}
	if rule.V5 != "" {
		ruleText += ", " + rule.V5
	}
	persist.LoadPolicyLine(ruleText, model)
}

// LoadPolicy loads policy from database.
// LoadPolicy 从数据库加载策略。
func (a *Adapter) LoadPolicy(model model.Model) error {
	var lines []CasbinRule

	if err := a.db.Table(a.tableName).Scan(&lines); err != nil {
		return err
	}

	for _, line := range lines {
		loadPolicyLine(line, model)
	}

	return nil
}

func savePolicyLine(ptype string, rule []string) CasbinRule {
	line := CasbinRule{}

	line.PType = ptype
	if len(rule) > 0 {
		line.V0 = rule[0]
	}
	if len(rule) > 1 {
		line.V1 = rule[1]
	}
	if len(rule) > 2 {
		line.V2 = rule[2]
	}
	if len(rule) > 3 {
		line.V3 = rule[3]
	}
	if len(rule) > 4 {
		line.V4 = rule[4]
	}
	if len(rule) > 5 {
		line.V5 = rule[5]
	}

	return line
}

// SavePolicy saves policy to database.
// SavePolicy 将策略保存到数据库。
func (a *Adapter) SavePolicy(model model.Model) error {
	if err := a.dropTable(); err != nil {
		return err
	}
	if err := a.createTable(); err != nil {
		return err
	}

	for ptype, ast := range model["p"] {
		for _, rule := range ast.Policy {
			line := savePolicyLine(ptype, rule)
			_, err := a.db.Table(a.tableName).Data(&line).Insert()
			if err != nil {
				return err
			}
		}
	}

	for ptype, ast := range model["g"] {
		for _, rule := range ast.Policy {
			line := savePolicyLine(ptype, rule)
			_, err := a.db.Table(a.tableName).Data(&line).Insert()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// AddPolicy adds a policy rule to the storage.
// AddPolicy 向存储添加策略规则。
func (a *Adapter) AddPolicy(sec string, ptype string, rule []string) error {
	line := savePolicyLine(ptype, rule)
	_, err := a.db.Table(a.tableName).Data(&line).Insert()
	return err
}

// RemovePolicy removes a policy rule from the storage.
// RemovePolicy 从存储中删除策略规则。
func (a *Adapter) RemovePolicy(sec string, ptype string, rule []string) error {
	line := savePolicyLine(ptype, rule)
	err := rawDelete(a, line)
	return err
}

// RemoveFilteredPolicy removes policy rules that match the filter from the storage.
// RemoveFilteredPolicy 从存储中删除与筛选器匹配的策略规则。
func (a *Adapter) RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error {
	line := CasbinRule{}

	line.PType = ptype
	if fieldIndex <= 0 && 0 < fieldIndex+len(fieldValues) {
		line.V0 = fieldValues[0-fieldIndex]
	}
	if fieldIndex <= 1 && 1 < fieldIndex+len(fieldValues) {
		line.V1 = fieldValues[1-fieldIndex]
	}
	if fieldIndex <= 2 && 2 < fieldIndex+len(fieldValues) {
		line.V2 = fieldValues[2-fieldIndex]
	}
	if fieldIndex <= 3 && 3 < fieldIndex+len(fieldValues) {
		line.V3 = fieldValues[3-fieldIndex]
	}
	if fieldIndex <= 4 && 4 < fieldIndex+len(fieldValues) {
		line.V4 = fieldValues[4-fieldIndex]
	}
	if fieldIndex <= 5 && 5 < fieldIndex+len(fieldValues) {
		line.V5 = fieldValues[5-fieldIndex]
	}
	err := rawDelete(a, line)
	return err
}

func rawDelete(a *Adapter, line CasbinRule) error {
	db := a.db.Table(a.tableName).Safe()

	db.Where("ptype = ?", line.PType)
	if line.V0 != "" {
		db.Where("v0 = ?", line.V0)
	}
	if line.V1 != "" {
		db.Where("v1 = ?", line.V1)
	}
	if line.V2 != "" {
		db.Where("v2 = ?", line.V2)
	}
	if line.V3 != "" {
		db.Where("v3 = ?", line.V3)
	}
	if line.V4 != "" {
		db.Where("v4 = ?", line.V4)
	}
	if line.V5 != "" {
		db.Where("v5 = ?", line.V5)
	}

	_, err := db.Delete()
	return err
}
