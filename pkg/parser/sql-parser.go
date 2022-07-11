package parser

import (
	"context"
	"github.com/acmestack/gobatis-plus/pkg/mapper"
)

type SqlParser[T any] struct {
}

func (sqlParser *SqlParser[T]) parser(ctx context.Context, queryWrapper mapper.Wrapper[T]) {

}
