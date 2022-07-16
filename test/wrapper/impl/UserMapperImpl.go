package impl

import (
	"context"
	"github.com/acmestack/gobatis-plus/pkg/mapper"
	_ "github.com/go-sql-driver/mysql"
	"github.com/xfali/gobatis"
	"github.com/xfali/gobatis/builder"
)

type UserMapperImpl[T any] struct {
	SessMgr *gobatis.SessionManager
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
func (userMapper *UserMapperImpl[T]) SelectList(ctx context.Context, queryWrapper mapper.QueryWrapper[T]) []T {
	sess := userMapper.SessMgr.NewSession()
	var arr []T
	for k, _ := range queryWrapper.MapCondition {
		condition := k
		sql := builder.Select(queryWrapper.Columns...).From("test_table").Where(condition).String()
		err := sess.Select(sql).Param(queryWrapper.Entity).Result(&arr)
		if err != nil {
			return nil
		}
	}
	return arr
}
