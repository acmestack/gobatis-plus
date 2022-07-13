package impl

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/acmestack/gobatis-plus/pkg/mapper"
	_ "github.com/go-sql-driver/mysql"
	"github.com/xfali/gobatis"
	"github.com/xfali/gobatis/builder"
	"github.com/xfali/gobatis/datasource"
	"github.com/xfali/gobatis/factory"
	"strconv"
)

type UserMapperImpl[T any] struct {
}

func (userMapper *UserMapperImpl[T]) Insert(ctx context.Context, entity T) int64 {
	return 0
}
func (userMapper *UserMapperImpl[T]) InsertBatch(ctx context.Context, entities ...T) (int64, int64) {
	return 0, 0
}
func (userMapper *UserMapperImpl[T]) DeleteById(ctx context.Context, id any) int64 {
	return 0
}
func (userMapper *UserMapperImpl[T]) DeleteBatchIds(ctx context.Context, ids []any) int64 {
	return 0
}
func (userMapper *UserMapperImpl[T]) UpdateById(ctx context.Context, entity T) int64 {
	return 0
}
func (userMapper *UserMapperImpl[T]) SelectById(ctx context.Context, id any) T {

	return *new(T)
}
func (userMapper *UserMapperImpl[T]) SelectBatchIds(ctx context.Context, ids []any) []T {
	var arr []T
	return arr
}
func (userMapper *UserMapperImpl[T]) SelectOne(ctx context.Context, entity T) T {
	return *new(T)
}
func (userMapper *UserMapperImpl[T]) SelectCount(ctx context.Context, entity T) int64 {
	return 0
}
func (userMapper *UserMapperImpl[T]) SelectList(ctx context.Context, queryWrapper mapper.QueryWrapper[T]) []*T {
	mgr := gobatis.NewSessionManager(connect())
	sess := mgr.NewSession()
	t := new(T)
	for k, v := range queryWrapper.MapCondition {
		i := v.(int)
		condition := k + strconv.Itoa(i)
		str := builder.Select("id", "username", "password").From("test_table").Where(condition).String()
		sess.Select(str).Param().Result(t)
	}
	marshal, _ := json.Marshal(t)
	fmt.Println(string(marshal))
	var arr []*T
	arr = append(arr, t)
	return arr
}

type TestTable struct {
	TestTable gobatis.ModelName "test_table"
	Id        int64             `xfield:"id"`
	Username  string            `xfield:"username"`
	Password  string            `xfield:"password"`
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
