package main

import (
	"fmt"
	"github.com/acmestack/gobatis-plus/example"
	"github.com/acmestack/gobatis-plus/pkg/query"
)

func main() {
	queryWrapper := query.QueryWrapper[example.UserDo]{}
	fmt.Println(queryWrapper)
}
