package mapper

import (
	"context"
	"github.com/acmestack/gobatis-plus/pkg/constants"
	"reflect"
)

type QueryWrapper[T any] struct {
	MapCondition map[string]any
	ctx          context.Context
	Columns      []string
	Entity       *T
}

func (queryWrapper *QueryWrapper[T]) Eq(column string, val any) Wrapper[T] {
	queryWrapper.init()
	entityValueRef := reflect.ValueOf(queryWrapper.Entity).Elem()
	entityRef := reflect.TypeOf(queryWrapper.Entity).Elem()
	name := entityRef.Name()
	numField := entityRef.NumField()
	for i := 0; i < numField; i++ {
		field := entityRef.Field(i)
		filedName := field.Tag.Get("xfield")
		if filedName == column {
			entityValueRef.FieldByName(field.Name).SetString("123")
		}
	}
	key := column + constants.Eq + "#{" + name + "." + column + "}"
	queryWrapper.MapCondition[key] = val
	return queryWrapper
}

func (queryWrapper *QueryWrapper[T]) Ne(column string, val any) Wrapper[T] {
	return nil
}

func (queryWrapper *QueryWrapper[T]) Gt(column string, val any) Wrapper[T] {
	return nil
}

func (queryWrapper *QueryWrapper[T]) Ge(column string, val any) Wrapper[T] {
	return nil
}

func (queryWrapper *QueryWrapper[T]) Lt(column string, val any) Wrapper[T] {
	return nil
}

func (queryWrapper *QueryWrapper[T]) Le(column string, val any) Wrapper[T] {
	return nil
}

func (queryWrapper *QueryWrapper[T]) Like(column string, val any) Wrapper[T] {
	return nil
}

func (queryWrapper *QueryWrapper[T]) NotLike(column string, val any) Wrapper[T] {
	return nil
}

func (queryWrapper *QueryWrapper[T]) LikeLeft(column string, val any) Wrapper[T] {
	return nil
}

func (queryWrapper *QueryWrapper[T]) LikeRight(column string, val1 any) Wrapper[T] {
	return nil
}

func (queryWrapper *QueryWrapper[T]) Between(column string, val1 any, val2 any) Wrapper[T] {
	return nil
}

func (queryWrapper *QueryWrapper[T]) NotBetween(column string, val1 any, val2 any) Wrapper[T] {
	return nil
}

func (queryWrapper *QueryWrapper[T]) Select(columns ...string) Wrapper[T] {
	for _, v := range columns {
		queryWrapper.Columns = append(queryWrapper.Columns, v)
	}
	return queryWrapper
}

func (queryWrapper *QueryWrapper[T]) init() {
	if len(queryWrapper.MapCondition) == 0 {
		queryWrapper.MapCondition = make(map[string]any, 16)
	}
	if queryWrapper.Entity == nil {
		queryWrapper.Entity = new(T)
	}
	if queryWrapper.ctx == nil {
		queryWrapper.ctx = context.Background()
	}
}
