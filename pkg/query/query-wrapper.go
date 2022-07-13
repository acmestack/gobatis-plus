package query

import (
	"github.com/acmestack/gobatis-plus/pkg/mapper"
)

type QueryWrapper[T any] struct {
	MapCondition map[string]any
}

func (queryWrapper *QueryWrapper[T]) Eq(column string, val any) mapper.Wrapper[T] {
	queryWrapper.initMap()
	queryWrapper.MapCondition[column] = val
	return queryWrapper
}

func (queryWrapper *QueryWrapper[T]) Ne(column string, val any) mapper.Wrapper[T] {
	return nil
}

func (queryWrapper *QueryWrapper[T]) Gt(column string, val any) mapper.Wrapper[T] {
	return nil
}

func (queryWrapper *QueryWrapper[T]) Ge(column string, val any) mapper.Wrapper[T] {
	return nil
}

func (queryWrapper *QueryWrapper[T]) Lt(column string, val any) mapper.Wrapper[T] {
	return nil
}

func (queryWrapper *QueryWrapper[T]) Le(column string, val any) mapper.Wrapper[T] {
	return nil
}

func (queryWrapper *QueryWrapper[T]) Like(column string, val any) mapper.Wrapper[T] {
	return nil
}

func (queryWrapper *QueryWrapper[T]) NotLike(column string, val any) mapper.Wrapper[T] {
	return nil
}

func (queryWrapper *QueryWrapper[T]) LikeLeft(column string, val any) mapper.Wrapper[T] {
	return nil
}

func (queryWrapper *QueryWrapper[T]) LikeRight(column string, val1 any) mapper.Wrapper[T] {
	return nil
}

func (queryWrapper *QueryWrapper[T]) Between(column string, val1 any, val2 any) mapper.Wrapper[T] {
	return nil
}

func (queryWrapper *QueryWrapper[T]) NotBetween(column string, val1 any, val2 any) mapper.Wrapper[T] {
	return nil
}

func (queryWrapper *QueryWrapper[T]) initMap() {
	if len(queryWrapper.MapCondition) == 0 {
		queryWrapper.MapCondition = make(map[string]any, 16)
	}
}
