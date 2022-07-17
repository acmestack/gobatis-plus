package mapper

import (
	"context"
	"github.com/xfali/gobatis"
	"github.com/xfali/gobatis/builder"
)

type BaseMapper[T any] struct {
	SessMgr *gobatis.SessionManager
}

func (userMapper *BaseMapper[T]) Insert(ctx context.Context, entity T) int64 {
	return 0
}
func (userMapper *BaseMapper[T]) InsertBatch(ctx context.Context, entities ...T) (int64, int64) {
	return 0, 0
}
func (userMapper *BaseMapper[T]) DeleteById(ctx context.Context, id any) int64 {
	return 0
}
func (userMapper *BaseMapper[T]) DeleteBatchIds(ctx context.Context, ids []any) int64 {
	return 0
}
func (userMapper *BaseMapper[T]) UpdateById(ctx context.Context, entity T) int64 {
	return 0
}
func (userMapper *BaseMapper[T]) SelectById(ctx context.Context, id any) T {

	return *new(T)
}
func (userMapper *BaseMapper[T]) SelectBatchIds(ctx context.Context, ids []any) []T {
	var arr []T
	return arr
}
func (userMapper *BaseMapper[T]) SelectOne(ctx context.Context, entity T) T {
	return *new(T)
}
func (userMapper *BaseMapper[T]) SelectCount(ctx context.Context, entity T) int64 {
	return 0
}
func (userMapper *BaseMapper[T]) SelectList(ctx context.Context, queryWrapper QueryWrapper[T]) []T {
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
