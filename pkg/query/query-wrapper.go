package query

import "github.com/acmestack/gobatis-plus/pkg/mapper"

type QueryWrapper[T any] struct {
	mapper.Wrapper[T]
}
