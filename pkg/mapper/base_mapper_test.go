/*
 * Licensed to the AcmeStack under one or more contributor license
 * agreements. See the NOTICE file distributed with this work for
 * additional information regarding copyright ownership.
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package mapper

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/acmestack/gobatis"
	"github.com/acmestack/gobatis/datasource"
	"github.com/acmestack/gobatis/factory"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func connect() factory.Factory {
	return gobatis.NewFactory(
		gobatis.SetMaxConn(100),
		gobatis.SetMaxIdleConn(50),
		gobatis.SetDataSource(&datasource.MysqlDataSource{
			Host:     "localhost",
			Port:     3306,
			DBName:   "test",
			Username: "root",
			Password: "test",
			Charset:  "utf8",
		}))
}

func TestInitTable(t *testing.T) {
	sql_table := "CREATE TABLE IF NOT EXISTS `test_table` (" +
		"`id` int(11) NOT NULL AUTO_INCREMENT," +
		"`username` varchar(255) DEFAULT NULL," +
		"`password` varchar(255) DEFAULT NULL," +
		"`createTime` datetime DEFAULT NULL," +
		"PRIMARY KEY (`id`)" +
		") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;"

	db, err := sql.Open("mysql", "test:test@tcp(localhost:3306)/test?charset=utf8")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(sql_table)
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec("DELETE FROM test_table")
	if err != nil {
		t.Fatal(err)
	}
}

type TestTable struct {
	TableName gobatis.TableName `test_table`
	Id        int64             `column:"id"`
	Username  string            `column:"username"`
	Password  string            `column:"password"`
}

func TestUserMapperImpl_SelectList(t *testing.T) {
	mgr := gobatis.NewSessionManager(connect())
	userMapper := BaseMapper[TestTable]{SessMgr: mgr}
	queryWrapper := &QueryWrapper[TestTable]{}
	queryWrapper.Eq("username", "acmestack").In("password", "123456", "pw5")
	list, err := userMapper.SelectList(queryWrapper)
	if err != nil {
		fmt.Println(err.Error())
	}
	marshal, _ := json.Marshal(list)
	fmt.Println(string(marshal))
}

func TestUserMapperImpl_SelectOne(t *testing.T) {
	mgr := gobatis.NewSessionManager(connect())
	userMapper := BaseMapper[TestTable]{SessMgr: mgr}
	queryWrapper := &QueryWrapper[TestTable]{}
	queryWrapper.Eq("username", "gobatis").Select("username", "password")
	entity, err := userMapper.SelectOne(queryWrapper)
	if err != nil {
		fmt.Println(err.Error())
	}
	marshal, _ := json.Marshal(entity)
	fmt.Println(string(marshal))
}

func TestUserMapperImpl_SelectCount(t *testing.T) {
	mgr := gobatis.NewSessionManager(connect())
	userMapper := BaseMapper[TestTable]{SessMgr: mgr}
	queryWrapper := &QueryWrapper[TestTable]{}
	queryWrapper.Eq("username", "user123")
	count, err := userMapper.SelectCount(queryWrapper)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(count)
}

func TestUserMapperImpl_SelectById(t *testing.T) {
	mgr := gobatis.NewSessionManager(connect())
	userMapper := BaseMapper[TestTable]{SessMgr: mgr}
	entity, err := userMapper.SelectById(103)
	if err != nil {
		fmt.Println(err.Error())
	}
	marshal, _ := json.Marshal(entity)
	fmt.Println(string(marshal))
}

func TestUserMapperImpl_SelectBatchIds(t *testing.T) {
	mgr := gobatis.NewSessionManager(connect())
	userMapper := BaseMapper[TestTable]{SessMgr: mgr}
	var arr []any
	arr = append(arr, 1)
	arr = append(arr, 103)
	list, err := userMapper.SelectBatchIds(arr)
	if err != nil {
		fmt.Println(err.Error())
	}
	marshal, _ := json.Marshal(list)
	fmt.Println(string(marshal))
}

func TestUserMapperImpl_Save(t *testing.T) {
	mgr := gobatis.NewSessionManager(connect())
	userMapper := BaseMapper[TestTable]{SessMgr: mgr}
	uuid := fmt.Sprintf("%d", random())
	table := TestTable{Username: "gobatis" + uuid, Password: "123456"}
	ret, id, err := userMapper.Save(table)
	if err != nil {
		t.Fail()
	}
	table.Id = id
	queryWrapper := &QueryWrapper[TestTable]{}
	queryWrapper.Eq("username", "gobatis"+uuid).Eq("password", "123456")
	one, err := userMapper.SelectOne(queryWrapper)
	if err != nil {
		t.Fail()
	}
	fmt.Println("ret:", ret)
	assert.Equal(t, table, one, "they should be equal")
}

func TestUserMapperImpl_SaveBatch(t *testing.T) {
	mgr := gobatis.NewSessionManager(connect())
	userMapper := BaseMapper[TestTable]{SessMgr: mgr}
	var entities []TestTable
	username1 := "gobatis" + fmt.Sprintf("%d", random())
	username2 := "gobatis" + fmt.Sprintf("%d", random())
	username3 := "gobatis" + fmt.Sprintf("%d", random())
	id1 := random()
	id2 := random()
	id3 := random()
	table1 := TestTable{Id: id1, Username: username1, Password: "123456"}
	table2 := TestTable{Id: id2, Username: username2, Password: "123456"}
	table3 := TestTable{Id: id3, Username: username3, Password: "123456"}
	entities = append(entities, table1)
	entities = append(entities, table2)
	entities = append(entities, table3)
	ret, id, err := userMapper.SaveBatch(entities...)
	if err != nil {
		t.Fail()
	}
	fmt.Println(ret, id)
	queryWrapper := &QueryWrapper[TestTable]{}
	queryWrapper.In("username", username1, username2, username3).Eq("password", "123456")
	list, err := userMapper.SelectList(queryWrapper)
	fmt.Println(entities)
	fmt.Println(list)
}

func TestUserMapperImpl_Delete(t *testing.T) {
	mgr := gobatis.NewSessionManager(connect())
	userMapper := BaseMapper[TestTable]{SessMgr: mgr}
	username := "gobatis" + fmt.Sprintf("%d", random())
	table := TestTable{Username: username, Password: "123456"}
	ret, id, err := userMapper.Save(table)
	if err != nil {
		t.Fail()
	}

	ret2, err := userMapper.DeleteById(id)
	if err != nil {
		t.Fail()
	}

	if ret2 != 1 {
		t.Fail()
	}
	fmt.Println(ret)
}

func TestUserMapperImpl_DeleteBatch(t *testing.T) {
	mgr := gobatis.NewSessionManager(connect())
	userMapper := BaseMapper[TestTable]{SessMgr: mgr}

	var entities []TestTable
	username1 := "gobatis" + fmt.Sprintf("%d", random())
	username2 := "gobatis" + fmt.Sprintf("%d", random())
	username3 := "gobatis" + fmt.Sprintf("%d", random())
	id1 := random()
	id2 := random()
	id3 := random()
	table1 := TestTable{Id: id1, Username: username1, Password: "123456"}
	table2 := TestTable{Id: id2, Username: username2, Password: "123456"}
	table3 := TestTable{Id: id3, Username: username3, Password: "123456"}
	entities = append(entities, table1)
	entities = append(entities, table2)
	entities = append(entities, table3)

	ret, id, err := userMapper.SaveBatch(entities...)
	if err != nil {
		t.Fail()
	}
	fmt.Println(ret, id)

	var ids []any
	ids = append(ids, id1)
	ids = append(ids, id2)
	ids = append(ids, id3)
	ret, err = userMapper.DeleteBatchIds(ids)
	if err != nil {
		t.Fail()
	}
	if ret != 3 {
		t.Fail()
	}
	fmt.Println("ret", ret)
}

func TestUserMapperImpl_UpdateById(t *testing.T) {
	mgr := gobatis.NewSessionManager(connect())
	userMapper := BaseMapper[TestTable]{SessMgr: mgr}

	uuid := fmt.Sprintf("%d", random())
	table := TestTable{Username: "gobatis" + uuid, Password: "123456"}
	ret, id, err := userMapper.Save(table)
	if err != nil {
		t.Fail()
	}
	fmt.Println(ret, id)

	var entity = TestTable{Id: id, Username: "gobatis", Password: "123456"}
	id, err = userMapper.UpdateById(entity)
	if err != nil {
		t.Fail()
	}

	if id != 1 {
		t.Fail()
	}
	fmt.Println(ret)
}

func random() int64 {
	intn := rand.Intn(100000000)
	return int64(intn)
}
