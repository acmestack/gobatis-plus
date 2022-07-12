package wrapper

import (
	"context"
	"github.com/acmestack/gobatis-plus/pkg/mapper"
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

	return nil
}
func (userMapper *UserMapperImpl[T]) SelectBatchIds(ctx context.Context, ids []any) []T {
	return nil
}
func (userMapper *UserMapperImpl[T]) SelectOne(ctx context.Context, entity T) T {
	return nil
}
func (userMapper *UserMapperImpl[T]) SelectCount(ctx context.Context, entity T) int64 {
	return 0
}
func (userMapper *UserMapperImpl[T]) SelectList(ctx context.Context, queryWrapper mapper.Wrapper[T]) []T {
	return nil
}
