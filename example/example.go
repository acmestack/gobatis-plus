/*
 * Copyright (c) 2022, AcmeStack
 * All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package example

import (
	"time"

	"github.com/acmestack/gobatis-plus/pkg/mapper"
)

// +gobatis:data:tablename=tbl_user
type UserDo struct {
	// +gobatis:tableid:value=user_id,idType=auto
	UserId int64

	// +gobatis:tablefield:value=user_name
	UserName string

	// +gobatis:tablefield:value=status
	// +gobatis:tablelogic:value=0,delval=1
	Status int8

	// +gobatis:tablefield:value=create_time,fill=insert
	CreateTime time.Time

	// +gobatis:tablefield:value=rec_var
	// +gobatis:version
	RecVersion uint64
}

// +gobatis:mapper
type UserMapper interface {
	mapper.BaseMapper[UserDo]

	// +gobatis:select="select * from tbl_user where id = #{UserDo.UserId}"
	Select(user UserDo) (users []UserDo, err error)
}
