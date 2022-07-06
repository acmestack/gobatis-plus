// Copyright (C) 2022, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description:

package mapper

type BaseMapper[T any] interface {
	Insert(entity T) int64

	InsertBatch(entities ...T) (int64, int64)

	DeleteById(id any) int64

	DeleteBatchIds(ids []any) int64

	UpdateById(entity T) int64

	SelectById(id any) T

	SelectBatchIds(ids []any) []T

	SelectOne(entity T) T

	SelectCount(entity T) int64
}
