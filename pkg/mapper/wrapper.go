// Copyright (C) 2022, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description:

package mapper

import "context"

type Wrapper[T any] interface {
	Eq(ctx context.Context, column any, val any) Wrapper[T]

	Ne(ctx context.Context, column any, val any) Wrapper[T]

	Gt(ctx context.Context, column any, val any) Wrapper[T]

	Ge(ctx context.Context, column any, val any) Wrapper[T]

	Lt(ctx context.Context, column any, val any) Wrapper[T]

	Le(ctx context.Context, column any, val any) Wrapper[T]

	Like(ctx context.Context, column any, val any) Wrapper[T]

	NotLike(ctx context.Context, column any, val any) Wrapper[T]

	LikeLeft(ctx context.Context, column any, val any) Wrapper[T]

	LikeRight(ctx context.Context, column any, val1 any) Wrapper[T]

	Between(ctx context.Context, column any, val1 any, val2 any) Wrapper[T]

	NotBetween(ctx context.Context, column any, val1 any, val2 any) Wrapper[T]
}
