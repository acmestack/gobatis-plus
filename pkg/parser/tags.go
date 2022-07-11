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

package parser

type TableName string

const (
	TagTableName = "tableName"

	TagPrimaryKey = "primaryKey"

	// Example: UserId int64 `tableId:"user_id,idType=auto"`
	TagTableId = "tableId"

	// For TAG [TableId], options: [ auto | none | input | assign_id | assign_uuid ]
	TagTableIdType = "idType"

	// Example: UserName string `tableField:"user_name,fill=insert"`
	TagTableField = "tableField"

	// For TAG [TagTableField], options: [ default | insert | update | insert_update ]
	TagTableFieldFill = "fill"

	// Example: Status int8 `tableField:"status" tableLogic:"0,delval=1"`
	TagTableLogic = "tableLogic"

	// Define the value to mark record deleted
	TagTableLogicDelVal = "delval"

	// Record version
	TagRecordVersion = "tableRecVer"
)
