// Copyright (C) 2022, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description:

package mapper

type Wrapper[T any] interface {
	Eq(column string, val any) Wrapper[T]

	Ne(column string, val any) Wrapper[T]

	Gt(column string, val any) Wrapper[T]

	Ge(column string, val any) Wrapper[T]

	Lt(column string, val any) Wrapper[T]

	Le(column string, val any) Wrapper[T]

	Like(column string, val any) Wrapper[T]

	NotLike(column string, val any) Wrapper[T]

	LikeLeft(column string, val any) Wrapper[T]

	LikeRight(column string, val1 any) Wrapper[T]

	Between(column string, val1 any, val2 any) Wrapper[T]

	NotBetween(column string, val1 any, val2 any) Wrapper[T]

	Select(columns ...string) Wrapper[T]
}
