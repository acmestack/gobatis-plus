// Copyright (C) 2022, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description:

package mapper

import "context"

type Wrapper[T any] interface {
	Eq(ctx context.Context, column any, val any) T

	Ne(ctx context.Context, column any, val any) T

	Gt(ctx context.Context, column any, val any) T

	Ge(ctx context.Context, column any, val any) T

	Lt(ctx context.Context, column any, val any) T

	Le(ctx context.Context, column any, val any) T

	Like(ctx context.Context, column any, val any) T

	NotLike(ctx context.Context, column any, val any) T

	LikeLeft(ctx context.Context, column any, val any) T

	LikeRight(ctx context.Context, column any, val1 any) T

	Between(ctx context.Context, column any, val1 any, val2 any) T

	NotBetween(ctx context.Context, column any, val1 any, val2 any) T
}
