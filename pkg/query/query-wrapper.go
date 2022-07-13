package query

import (
	"context"
	"github.com/acmestack/gobatis-plus/pkg/mapper"
)

type QueryWrapper[T any] struct {
}

func (queryWrapper *QueryWrapper[T]) Eq(ctx context.Context, column any, val any) mapper.Wrapper[T] {
	return nil
}

func (queryWrapper *QueryWrapper[T]) Ne(ctx context.Context, column any, val any) mapper.Wrapper[T] {
	return nil
}

func (queryWrapper *QueryWrapper[T]) Gt(ctx context.Context, column any, val any) mapper.Wrapper[T] {
	return nil
}

func (queryWrapper *QueryWrapper[T]) Ge(ctx context.Context, column any, val any) mapper.Wrapper[T] {
	return nil
}

func (queryWrapper *QueryWrapper[T]) Lt(ctx context.Context, column any, val any) mapper.Wrapper[T] {
	return nil
}

func (queryWrapper *QueryWrapper[T]) Le(ctx context.Context, column any, val any) mapper.Wrapper[T] {
	return nil
}

func (queryWrapper *QueryWrapper[T]) Like(ctx context.Context, column any, val any) mapper.Wrapper[T] {
	return nil
}

func (queryWrapper *QueryWrapper[T]) NotLike(ctx context.Context, column any, val any) mapper.Wrapper[T] {
	return nil
}

func (queryWrapper *QueryWrapper[T]) LikeLeft(ctx context.Context, column any, val any) mapper.Wrapper[T] {
	return nil
}

func (queryWrapper *QueryWrapper[T]) LikeRight(ctx context.Context, column any, val1 any) mapper.Wrapper[T] {
	return nil
}

func (queryWrapper *QueryWrapper[T]) Between(ctx context.Context, column any, val1 any, val2 any) mapper.Wrapper[T] {
	return nil
}

func (queryWrapper *QueryWrapper[T]) NotBetween(ctx context.Context, column any, val1 any, val2 any) mapper.Wrapper[T] {
	return nil
}
