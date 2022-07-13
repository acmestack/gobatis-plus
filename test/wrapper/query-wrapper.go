package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/acmestack/gobatis-plus/pkg/mapper"
	"github.com/acmestack/gobatis-plus/test/wrapper/impl"
)

func main() {
	userMapper := impl.UserMapperImpl[impl.TestTable]{}
	queryWrapper := mapper.QueryWrapper[impl.TestTable]{}
	queryWrapper.Eq("id", 4)
	list := userMapper.SelectList(context.Background(), queryWrapper)
	marshal, _ := json.Marshal(list)
	fmt.Println(string(marshal))
}
