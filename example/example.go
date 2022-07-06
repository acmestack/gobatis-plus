// Copyright (C) 2019-2022, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package example

import "time"

// +gobatis:data
type UserDo struct {
	//
	UserId     int64     `tableId:"user_id,idType=auto"`
	UserName   string    `tableField:"user_name"`
	Status     int8      `tableField:"status" tableLogic:"0,delval=1"`
	CreateTime time.Time `tableField:"create_time,fill=insert"`
	RecVersion uint64    `tableRecVer:"rec_var"`
}

// +gobatis:mapper
type UserMapper interface {
	// +gobatis:select="select * from "

	Insert(vo ...UserDo) (err error)
}
