// Copyright (C) 2019-2022, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package parser

type TableName string

const (
	TagTableName   = "tableName"

	TagPrimaryKey  = "primaryKey"

	// Example: UserId int64 `tableId:"user_id,idType=auto"`
	TagTableId     = "tableId"

	// For TAG [TableId], options: [ auto | none | input | assign_id | assign_uuid ]
	TagTableIdType = "idType"

	// Example: UserName string `tableField:"user_name,fill=insert"`
	TagTableField  = "tableField"

	// For TAG [TagTableField], options: [ default | insert | update | insert_update ]
	TagTableFieldFill = "fill"

	// Example: Status int8 `tableField:"status" tableLogic:"0,delval=1"`
	TagTableLogic = "tableLogic"

	// Define the value to mark record deleted
	TagTableLogicDelVal = "delval"

	// Record version
	TagRecordVersion = "tableRecVer"
)
