package main

import (
	"encoding/json"
	"fmt"
	"github.com/acmestack/gobatis-plus/example"
	"github.com/acmestack/gobatis-plus/pkg/query"
)

func main() {
	queryWrapper := query.QueryWrapper[example.UserDo]{}
	queryWrapper.Eq("age", 1).Eq("aa", 2)
	marshal, _ := json.Marshal(queryWrapper)
	fmt.Println(string(marshal))
}
