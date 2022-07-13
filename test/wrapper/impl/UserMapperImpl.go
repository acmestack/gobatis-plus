package impl

import (
	"context"
	"github.com/acmestack/gobatis-plus/pkg/mapper"
)

type UserMapperImpl[T any, This any] struct {
}

func (userMapper *UserMapperImpl[T, This]) Insert(ctx context.Context, entity T) int64 {
	return 0
}
func (userMapper *UserMapperImpl[T, This]) InsertBatch(ctx context.Context, entities ...T) (int64, int64) {
	return 0, 0
}
func (userMapper *UserMapperImpl[T, This]) DeleteById(ctx context.Context, id any) int64 {
	return 0
}
func (userMapper *UserMapperImpl[T, This]) DeleteBatchIds(ctx context.Context, ids []any) int64 {
	return 0
}
func (userMapper *UserMapperImpl[T, This]) UpdateById(ctx context.Context, entity T) int64 {
	return 0
}
func (userMapper *UserMapperImpl[T, This]) SelectById(ctx context.Context, id any) T {

	return nil
}
func (userMapper *UserMapperImpl[T, This]) SelectBatchIds(ctx context.Context, ids []any) []T {
	return nil
}
func (userMapper *UserMapperImpl[T, This]) SelectOne(ctx context.Context, entity T) T {
	return nil
}
func (userMapper *UserMapperImpl[T, This]) SelectCount(ctx context.Context, entity T) int64 {
	return 0
}
func (userMapper *UserMapperImpl[T, This]) SelectList(ctx context.Context, queryWrapper mapper.Wrapper[T]) []T {
	return nil
}
