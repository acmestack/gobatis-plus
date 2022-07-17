package impl

import (
	"github.com/acmestack/gobatis-plus/pkg/mapper"
	_ "github.com/go-sql-driver/mysql"
)

type UserMapperImpl[T any] struct {
	mapper.BaseMapper[T]
}
