package wrapper

import (
	"context"
	"fmt"
	"github.com/acmestack/gobatis-plus/example"
	"github.com/acmestack/gobatis-plus/pkg/mapper"
	"github.com/acmestack/gobatis-plus/pkg/query"
)

func Test_(ctx context.Context) {
	wrapper := query.QueryWrapper[example.UserDo]{}
	wrapper.Eq(ctx, "age", "1")
	var baseMapper mapper.BaseMapper[example.UserDo]
	baseMapper = &UserMapperImpl[example.UserDo]{}
	fmt.Println(baseMapper)
}
