package impl

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/acmestack/gobatis-plus/pkg/mapper"
	"github.com/xfali/gobatis"
	"github.com/xfali/gobatis/datasource"
	"github.com/xfali/gobatis/factory"
	"testing"
)

func TestUserMapperImpl_SelectList(t *testing.T) {
	mgr := gobatis.NewSessionManager(connect())
	userMapper := UserMapperImpl[TestTable]{mapper.BaseMapper[TestTable]{SessMgr: mgr}}
	queryWrapper := mapper.QueryWrapper[TestTable]{}
	queryWrapper.Eq("username", 4).Select("id", "username", "password")
	list := userMapper.SelectList(context.Background(), queryWrapper)
	marshal, _ := json.Marshal(list)
	fmt.Println(string(marshal))
}

func connect() factory.Factory {
	return gobatis.NewFactory(
		gobatis.SetMaxConn(100),
		gobatis.SetMaxIdleConn(50),
		gobatis.SetDataSource(&datasource.MysqlDataSource{
			Host:     "123.57.13.246",
			Port:     3306,
			DBName:   "http_info",
			Username: "root",
			Password: "root-abcd-1234",
			Charset:  "utf8",
		}))
}

type TestTable struct {
	TestTable gobatis.ModelName "test_table"
	Id        int64             `xfield:"id"`
	Username  string            `xfield:"username"`
	Password  string            `xfield:"password"`
}
