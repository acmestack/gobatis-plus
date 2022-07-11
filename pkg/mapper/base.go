// Copyright (C) 2022, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description:

package mapper

import "context"

type BaseMapper[T any] interface {
	Insert(ctx context.Context, entity T) int64

	InsertBatch(ctx context.Context, entities ...T) (int64, int64)

	DeleteById(ctx context.Context, id any) int64

	DeleteBatchIds(ctx context.Context, ids []any) int64

	UpdateById(ctx context.Context, entity T) int64

	SelectById(ctx context.Context, id any) T

	SelectBatchIds(ctx context.Context, ids []any) []T

	SelectOne(ctx context.Context, entity T) T

	SelectCount(ctx context.Context, entity T) int64

	SelectList(ctx context.Context, queryWrapper Wrapper[T]) []T
}
