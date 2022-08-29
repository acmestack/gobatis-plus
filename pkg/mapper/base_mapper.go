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

package mapper

import (
	"context"
	"fmt"
	"github.com/acmestack/gobatis"
	"github.com/acmestack/gobatis-plus/pkg/constants"
	"github.com/acmestack/godkits/gox/stringsx"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type BaseMapper[T any] struct {
	SessMgr      *gobatis.SessionManager
	Ctx          context.Context
	Columns      []string
	ParamNameSeq int
}

type BuildSqlFunc func(columns string, tableName string) string

func (userMapper *BaseMapper[T]) SelectList(queryWrapper *QueryWrapper[T]) ([]T, error) {
	// 初始化queryWrapper，如果queryWrapper是空的，需要初始化一个新的
	queryWrapper = userMapper.init(queryWrapper)

	// 构建Select查询语句
	paramMap, sql, sqlId := userMapper.buildSelectSql(queryWrapper, "")

	err := gobatis.RegisterSql(sqlId, sql)
	defer gobatis.UnregisterSql(sqlId)
	if err != nil {
		return nil, err
	}

	// 创建会话查询数据
	sess := userMapper.SessMgr.NewSession()
	var results []T
	err = sess.Select(sqlId).Param(paramMap).Result(&results)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (userMapper *BaseMapper[T]) SelectById(id any) (T, error) {
	queryWrapper := userMapper.init(nil)
	switch v := id.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		queryWrapper.Eq(constants.ID, fmt.Sprintf("%d", v))
	case string:
		queryWrapper.Eq(constants.ID, v)
	}

	// 构建Select查询语句
	paramMap, sql, sqlId := userMapper.buildSelectSql(queryWrapper, "")

	// 注册sql

	err := gobatis.RegisterSql(sqlId, sql)
	defer gobatis.UnregisterSql(sqlId)
	var entity T
	if err != nil {
		return entity, err
	}

	// 创建会话查询数据
	sess := userMapper.SessMgr.NewSession()
	err = sess.Select(sqlId).Param(paramMap).Result(&entity)
	if err != nil {
		return entity, err
	}
	return entity, nil
}

func (userMapper *BaseMapper[T]) SelectBatchIds(ids []any) ([]T, error) {
	queryWrapper := userMapper.init(nil)
	queryWrapper.In(constants.ID, ids...)

	// 构建Select查询语句
	paramMap, sql, sqlId := userMapper.buildSelectSql(queryWrapper, "")

	err := gobatis.RegisterSql(sqlId, sql)
	defer gobatis.UnregisterSql(sqlId)
	if err != nil {
		return nil, err
	}

	// 创建会话查询数据
	sess := userMapper.SessMgr.NewSession()
	var arr []T
	err = sess.Select(sqlId).Param(paramMap).Result(&arr)
	if err != nil {
		return nil, err
	}

	return arr, nil
}

func (userMapper *BaseMapper[T]) SelectOne(queryWrapper *QueryWrapper[T]) (T, error) {
	// 初始化queryWrapper，如果queryWrapper是空的，需要初始化一个新的
	queryWrapper = userMapper.init(queryWrapper)

	// 构建Select查询语句
	paramMap, sql, sqlId := userMapper.buildSelectSql(queryWrapper, "")

	err := gobatis.RegisterSql(sqlId, sql)
	defer gobatis.UnregisterSql(sqlId)
	var entity T
	if err != nil {
		return entity, err
	}

	// 创建会话查询数据
	sess := userMapper.SessMgr.NewSession()
	err = sess.Select(sqlId).Param(paramMap).Result(&entity)
	if err != nil {
		return entity, err
	}

	return entity, nil
}

func (userMapper *BaseMapper[T]) SelectCount(queryWrapper *QueryWrapper[T]) (int64, error) {
	// 初始化queryWrapper，如果queryWrapper是空的，需要初始化一个新的
	queryWrapper = userMapper.init(queryWrapper)

	// 构建Select查询语句
	paramMap, sql, sqlId := userMapper.buildSelectSql(queryWrapper, constants.COUNT)

	err := gobatis.RegisterSql(sqlId, sql)
	defer gobatis.UnregisterSql(sqlId)
	if err != nil {
		return 0, err
	}

	// 创建会话查询数据
	sess := userMapper.SessMgr.NewSession()
	var count int64
	err = sess.Select(sqlId).Param(paramMap).Result(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (userMapper *BaseMapper[T]) Save(entity T) (int, int64, error) {
	sess := userMapper.SessMgr.NewSession()

	// 获取表名
	tableName := userMapper.getTableName()

	// 获取插入字段
	// eg：columnName1,columnName2,columnName3
	columns := userMapper.buildInsertColumns()

	// 构建插入语句后半部分
	// eg：(#{mapping1},#{mapping2},#{mapping3})

	endPartBuilder := userMapper.buildInsertColumnValue(entity, paramMap)

	// 构建插件语句
	insertSql := columns.String() + endPartBuilder.String()
	fmt.Println(insertSql)

	// 构建sqlId
	sqlId := buildSqlId(constants.INTO)

	// 注册sql
	err := gobatis.RegisterSql(sqlId, insertSql)
	if err != nil {
		return 0, 0, err
	}
	var ret int
	selectRunner := sess.Insert(sqlId).Param(paramMap)
	err = selectRunner.Result(&ret)
	if err != nil {
		return 0, 0, err
	}
	insertId := selectRunner.LastInsertId()
	return ret, insertId, nil
}

func (userMapper *BaseMapper[T]) SaveBatch(entities ...T) (int64, int64, error) {
	sess := userMapper.SessMgr.NewSession()
	// 构建insert语句前部分
	// eg：insert field1,field2 into tableName values
	builder := userMapper.buildInsertColumns()

	var paramMap = map[string]any{}
	for i, entity := range entities {
		// 构建插入语句后半部分
		// eg：(#{mapping1},#{mapping2},#{mapping3})
		endPartBuilder := userMapper.buildInsertColumnValue(entity, paramMap)
		if i != len(entities)-1 {
			builder.WriteString(endPartBuilder.String() + constants.COMMA)
		} else {
			builder.WriteString(endPartBuilder.String())
		}
	}
	fmt.Println(builder.String())

	// 构建sqlId
	sqlId := buildSqlId(constants.INTO)

	// 注册sql
	err := gobatis.RegisterSql(sqlId, builder.String())
	if err != nil {
		return 0, 0, err
	}
	var ret int64
	selectRunner := sess.Insert(sqlId).Param(paramMap)
	err = selectRunner.Result(&ret)
	if err != nil {
		return 0, 0, err
	}
	insertId := selectRunner.LastInsertId()
	return ret, insertId, nil
}

func (userMapper *BaseMapper[T]) DeleteById(id any) (int64, error) {
	tableName := userMapper.getTableName()
	builder := strings.Builder{}
	builder.WriteString(constants.DELETE + constants.SPACE + constants.FROM + constants.SPACE + tableName + constants.SPACE +
		constants.WHERE + constants.SPACE + constants.ID + constants.Eq)
	var paramMap = map[string]any{}
	// 构建idSql
	// eg： #{mapping1}
	idSqlBuilder := userMapper.buildIdSql(id, paramMap)

	sql := builder.String() + idSqlBuilder.String()
	fmt.Println(sql)

	sqlId := buildSqlId(constants.DELETE)

	err := gobatis.RegisterSql(sqlId, sql)
	if err != nil {
		return 0, err
	}

	sess := userMapper.SessMgr.NewSession()
	var ret int64
	err = sess.Delete(sqlId).Param(paramMap).Result(&ret)
	if err != nil {
		return 0, err
	}

	// delete sqlId
	gobatis.UnregisterSql(sqlId)
	return ret, nil
}

func (userMapper *BaseMapper[T]) DeleteBatchIds(ids []any) (int64, error) {
	tableName := userMapper.getTableName()
	builder := strings.Builder{}
	builder.WriteString(constants.DELETE + constants.SPACE + constants.FROM + constants.SPACE + tableName + constants.SPACE +
		constants.WHERE + constants.SPACE + constants.ID + constants.SPACE + constants.In + constants.SPACE + constants.LEFT_BRACKET)
	var paramMap = map[string]any{}
	// 构建idSql
	// eg： #{mapping1},#{mapping1},#{mapping1}
	for i, id := range ids {
		idSqlBuilder := userMapper.buildIdSql(id, paramMap)
		if i != len(ids)-1 {
			builder.WriteString(idSqlBuilder.String())
			builder.WriteString(",")
		} else {
			builder.WriteString(idSqlBuilder.String())
		}
	}
	builder.WriteString(constants.RIGHT_BRACKET)

	fmt.Println(builder.String())

	sqlId := buildSqlId(constants.DELETE)

	err := gobatis.RegisterSql(sqlId, builder.String())
	if err != nil {
		return 0, err
	}

	sess := userMapper.SessMgr.NewSession()
	var ret int64
	err = sess.Delete(sqlId).Param(paramMap).Result(&ret)
	if err != nil {
		return 0, err
	}

	// delete sqlId
	gobatis.UnregisterSql(sqlId)
	return ret, nil
}

func (userMapper *BaseMapper[T]) UpdateById(entity T) (int64, error) {
	tableName := userMapper.getTableName()
	sess := userMapper.SessMgr.NewSession()
	builder := strings.Builder{}
	builder.WriteString(constants.UPDATE + constants.SPACE + tableName + constants.SPACE + constants.SET + constants.SPACE)
	entityType := reflect.TypeOf(entity)
	entityValue := reflect.ValueOf(entity)
	numField := entityType.NumField()
	paramMap := map[string]any{}

	// 拼接需要更新的值
	// 字段名 = #{mapping}，字段名 = #{mapping}，字段名 = #{mapping}
	var idValue any
	for i := 0; i < numField; i++ {
		tag := entityType.Field(i).Tag
		column := tag.Get(constants.COLUMN)
		// 如果是等于Id的话也需要跳过
		if column == "" || constants.ID == column {
			idValue = entityValue.Field(i).Interface()
			continue
		}
		builder.WriteString(column + constants.Eq + constants.HASH_LEFT_BRACE)
		fieldValue := entityValue.Field(i).Interface()
		switch v := fieldValue.(type) {
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			idStr := fmt.Sprintf("%d", v)
			mapping := userMapper.getMappingSeq()
			builder.WriteString(mapping)
			paramMap[mapping] = idStr
		case string:
			mapping := userMapper.getMappingSeq()
			builder.WriteString(mapping)
			paramMap[mapping] = v
		}
		builder.WriteString(constants.RIGHT_BRACE)
		if i != numField-1 {
			builder.WriteString(constants.COMMA)
		}
	}

	// 拼接后半部分sql
	// eg： where id = #{mapping}
	switch v := idValue.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		idStr := fmt.Sprintf("%d", v)
		mapping := userMapper.getMappingSeq()

		builder.WriteString(constants.SPACE + constants.WHERE + constants.SPACE + constants.ID + constants.Eq + constants.HASH_LEFT_BRACE)
		builder.WriteString(mapping)
		builder.WriteString(constants.RIGHT_BRACE)
		paramMap[mapping] = idStr
	case string:
		mapping := userMapper.getMappingSeq()
		builder.WriteString(constants.WHERE + constants.SPACE + constants.ID + constants.Eq + constants.HASH_LEFT_BRACE)
		builder.WriteString(mapping)
		builder.WriteString(constants.RIGHT_BRACE)
		paramMap[mapping] = v
	}

	// 构建sqlId
	sqlId := buildSqlId(constants.UPDATE)

	fmt.Println(builder.String())
	// 注册sql
	err := gobatis.RegisterSql(sqlId, builder.String())
	if err != nil {
		return 0, err
	}
	var ret int64
	selectRunner := sess.Update(sqlId).Param(paramMap)
	err = selectRunner.Result(&ret)
	if err != nil {
		return 0, err
	}
	return ret, nil
}

func (userMapper *BaseMapper[T]) buildIdSql(id any, paramMap map[string]any) strings.Builder {
	builder := strings.Builder{}
	builder.WriteString(constants.HASH_LEFT_BRACE)
	switch v := id.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		idStr := fmt.Sprintf("%d", v)
		mapping := userMapper.getMappingSeq()
		builder.WriteString(mapping)
		paramMap[mapping] = idStr
	case string:
		mapping := userMapper.getMappingSeq()
		builder.WriteString(mapping)
		paramMap[mapping] = v
	}
	builder.WriteString(constants.RIGHT_BRACE)
	return builder
}

func (userMapper *BaseMapper[T]) buildSelectSql(queryWrapper *QueryWrapper[T], columns string) (map[string]any, string, string) {
	// 构建需要查询的字段，如果查询条件没有传入的话，默认查询所有
	if stringsx.Empty(columns) {
		// eg: columnName1,columnName2,columnName3
		columns = userMapper.buildSelectColumns(queryWrapper)
	}

	// 获取表名称
	tableName := userMapper.getTableName()

	// 构建查询条件
	// eg: columnName1 = #{mapping1} and columnName2 = #{mapping1}
	sqlCondition, paramMap := userMapper.buildCondition(queryWrapper)

	// 构建sql
	// eg: SELECT * FROM WHERE columnName = #{mapping1} and columnName = #{mapping1}
	sql := userMapper.buildSql(columns, tableName, sqlCondition)

	// 构建sqlId
	sqlId := buildSqlId(constants.SELECT)
	return paramMap, sql, sqlId
}

func (userMapper *BaseMapper[T]) buildSql(columns string, tableName string, sqlCondition string) string {
	sql := strings.Replace(constants.SELECT_SQL, constants.COLUMN_HASH, columns, -1)
	sql = strings.Replace(sql, constants.TABLE_NAME_HASH, tableName, -1)
	sql = strings.Replace(sql, constants.CONDITIONS_HASH, sqlCondition, -1)
	if sqlCondition == "" {
		sql = strings.Replace(sql, constants.WHERE, "", -1)
	}
	return sql
}

func (userMapper *BaseMapper[T]) buildInsertColumnValue(entity T) (map[string]any, string) {
	var paramMap = map[string]any{}
	entityType := reflect.TypeOf(entity)
	entityValue := reflect.ValueOf(entity)
	entityValueNum := entityValue.NumField()
	for i := 0; i < entityValueNum; i++ {
		tag := entityType.Field(i).Tag
		column := tag.Get(constants.COLUMN)
		if column == "" {
			continue
		}
		// 构建values值
		v := entityValue.Field(i).Interface()
		mapping := userMapper.getMappingSeq()
		switch value := v.(type) {
		case string:
			paramMap[mapping] = value
		case int64:
			paramMap[mapping] = value
		}
		if i != entityValueNum-1 {
			builder.WriteString(constants.HASH_LEFT_BRACE + mapping + constants.RIGHT_BRACE + constants.COMMA)
		} else {
			builder.WriteString(constants.HASH_LEFT_BRACE + mapping + constants.RIGHT_BRACE + constants.RIGHT_BRACKET)
		}
	}
	return builder
}

func (userMapper *BaseMapper[T]) buildInsertColumns() string {
	entityType := reflect.TypeOf(new(T)).Elem()
	entityTypeNum := entityType.NumField()
	var columns []string
	for i := 0; i < entityTypeNum; i++ {
		tag := entityType.Field(i).Tag
		column := tag.Get(constants.COLUMN)
		if stringsx.Empty(column) {
			continue
		}
		columns = append(columns, column)
	}
	return strings.Join(columns, ",")
}

func (userMapper *BaseMapper[T]) buildSelectColumns(queryWrapper *QueryWrapper[T]) string {
	var columns string
	if len(queryWrapper.Columns) > 0 {
		columns = strings.Join(queryWrapper.Columns, ",")
	} else {
		columns = constants.ASTERISK
	}
	return columns
}

func (userMapper *BaseMapper[T]) init(queryWrapper *QueryWrapper[T]) *QueryWrapper[T] {
	if queryWrapper == nil {
		queryWrapper = &QueryWrapper[T]{}
	}
	return queryWrapper
}

// 构建查询条件
func (userMapper *BaseMapper[T]) buildCondition(queryWrapper *QueryWrapper[T]) (string, map[string]any) {
	var paramMap = map[string]any{}
	conditions := queryWrapper.Conditions
	build := strings.Builder{}
	// 遍历所有的条件参数
	for _, v := range conditions {
		// 如果是ParamValue的话，通过#{} 拼接数据，并且把value值存储到paramMap中，方便后面使用
		// ParamValue 存储的是查询条件具体的值
		if paramValue, ok := v.(ParamValue); ok {
			rt := reflect.TypeOf(paramValue.value)
			rv := reflect.ValueOf(paramValue.value)

			if rt.Kind() == reflect.Slice {
				l := rv.Len()
				build.WriteString(constants.LEFT_BRACKET)
				for i := 0; i < l; i++ {
					elemV := rv.Index(i)
					if !elemV.CanInterface() {
						elemV = reflect.Indirect(elemV)
					}
					mapping := userMapper.getMappingSeq()
					switch iv := elemV.Interface().(type) {
					case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
						paramMap[mapping] = fmt.Sprintf("%d", iv)
					case string:
						paramMap[mapping] = iv
					}
					if i != l-1 {
						build.WriteString(constants.HASH_LEFT_BRACE + mapping + constants.RIGHT_BRACE + constants.COMMA)
					} else {
						build.WriteString(constants.HASH_LEFT_BRACE + mapping + constants.RIGHT_BRACE)
					}
				}
				build.WriteString(constants.RIGHT_BRACKET)
			} else {
				mapping := userMapper.getMappingSeq()
				paramMap[mapping] = paramValue.value
				build.WriteString(constants.HASH_LEFT_BRACE + mapping + constants.RIGHT_BRACE + constants.SPACE)
			}
		} else {
			build.WriteString(v.(string) + constants.SPACE)
		}
	}
	return build.String(), paramMap
}

func (userMapper *BaseMapper[T]) getTableName() string {
	entityRef := reflect.TypeOf(new(T)).Elem()
	tableNameTag := entityRef.Field(0).Tag
	tableName := string(tableNameTag)
	return tableName
}

// build sql id ,may need to select a better implementation
func buildSqlId(sqlType string) string {
	sqlId := sqlType + constants.CONNECTION + strconv.Itoa(time.Now().Nanosecond())
	return sqlId
}

func (userMapper *BaseMapper[T]) getMappingSeq() string {
	userMapper.ParamNameSeq = userMapper.ParamNameSeq + 1
	mapping := constants.MAPPING + strconv.Itoa(userMapper.ParamNameSeq)
	return mapping
}
