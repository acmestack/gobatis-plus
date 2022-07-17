package mapper

import (
	"context"
	"github.com/xfali/gobatis"
)

type BaseMapper[T any] struct {
	SessMgr *gobatis.SessionManager
	Ctx     context.Context
	Columns []string
}

func (userMapper *BaseMapper[T]) Insert(entity T) int64 {
	return 0
}

func (userMapper *BaseMapper[T]) InsertBatch(entities ...T) (int64, int64) {
	return 0, 0
}
func (userMapper *BaseMapper[T]) DeleteById(id any) int64 {
	return 0
}
func (userMapper *BaseMapper[T]) DeleteBatchIds(ids []any) int64 {
	return 0
}
func (userMapper *BaseMapper[T]) UpdateById(entity T) int64 {
	return 0
}
func (userMapper *BaseMapper[T]) SelectById(id any) T {
	return *new(T)
}
func (userMapper *BaseMapper[T]) SelectBatchIds(ids []any) []T {
	var arr []T
	return arr
}
func (userMapper *BaseMapper[T]) SelectOne(entity T) T {
	return *new(T)
}
func (userMapper *BaseMapper[T]) SelectCount(entity T) int64 {
	return 0
}

func (userMapper *BaseMapper[T]) SelectList(queryWrapper *QueryWrapper[T]) ([]T, error) {
	if queryWrapper == nil {
		queryWrapper = &QueryWrapper[T]{}
		queryWrapper.init()
	}
	sess := userMapper.SessMgr.NewSession()
	var arr []T
	err := sess.Select(queryWrapper.SqlBuild.String()).Param(queryWrapper.Entity).Result(&arr)
	if err != nil {
		return nil, err
	}
	return arr, nil
}
