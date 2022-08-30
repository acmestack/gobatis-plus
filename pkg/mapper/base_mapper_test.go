package mapper

import (
	"encoding/json"
	"fmt"
	"github.com/acmestack/gobatis"
	"github.com/acmestack/gobatis/datasource"
	"github.com/acmestack/gobatis/factory"
	_ "github.com/go-sql-driver/mysql"
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
			Password: "123456",
			Charset:  "utf8",
		}))
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
	queryWrapper.In("id", 1, 2, 3)
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
	queryWrapper.Eq("username", "zouchangfu").Select("username", "password")
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
	table := TestTable{Username: "hello", Password: "123456"}
	ret, id, err := userMapper.Save(table)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(ret, id)

}

func TestUserMapperImpl_SaveBatch(t *testing.T) {
	mgr := gobatis.NewSessionManager(connect())
	userMapper := BaseMapper[TestTable]{SessMgr: mgr}
	var entities []TestTable
	table1 := TestTable{Username: "zouchangfu1", Password: "123456"}
	table2 := TestTable{Username: "zouchangfu2", Password: "123456"}
	table3 := TestTable{Username: "zouchangfu3", Password: "123456"}
	entities = append(entities, table1)
	entities = append(entities, table2)
	entities = append(entities, table3)
	ret, id, err := userMapper.SaveBatch(entities...)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(ret, id)
}

func TestUserMapperImpl_Delete(t *testing.T) {
	mgr := gobatis.NewSessionManager(connect())
	userMapper := BaseMapper[TestTable]{SessMgr: mgr}
	ret, err := userMapper.DeleteById(138)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(ret)
}

func TestUserMapperImpl_DeleteBatch(t *testing.T) {
	mgr := gobatis.NewSessionManager(connect())
	userMapper := BaseMapper[TestTable]{SessMgr: mgr}
	var ids []any
	ids = append(ids, 135)
	ids = append(ids, 136)
	ids = append(ids, 137)
	ret, err := userMapper.DeleteBatchIds(ids)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(ret)
}

func TestUserMapperImpl_UpdateById(t *testing.T) {
	mgr := gobatis.NewSessionManager(connect())
	userMapper := BaseMapper[TestTable]{SessMgr: mgr}
	var entity = TestTable{Id: 1, Username: "zouchangfu000", Password: "123456"}
	ret, err := userMapper.UpdateById(entity)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(ret)
}
