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

func (userMapper *BaseMapper[T]) Save(entity T) (int64, error) {
	tableName := userMapper.getTableName()
	sess := userMapper.SessMgr.NewSession()
	builder := strings.Builder{}
	builder.WriteString(constants.INSERT + constants.SPACE + constants.INTO + constants.SPACE + tableName + constants.SPACE + constants.LEFT_BRACKET)
	entityType := reflect.TypeOf(entity)
	entityTypeNum := entityType.NumField()
	for i := 0; i < entityTypeNum; i++ {
		tag := entityType.Field(i).Tag
		column := tag.Get("column")
		if column == "" {
			continue
		}
		if i != entityTypeNum-1 {
			// 如果不是最后一个字段，需要在字段后面拼接逗号
			builder.WriteString(column + constants.COMMA)
		} else {
			// 如果是最后一个字段，需要在字段后面拼接右括号
			builder.WriteString(column + constants.RIGHT_BRACKET)
		}
	}

	builder.WriteString(constants.SPACE + constants.VALUES + constants.SPACE + constants.LEFT_BRACKET)
	var paramMap = map[string]any{}
	entityValue := reflect.ValueOf(entity)
	entityValueNum := entityValue.NumField()

	for i := 0; i < entityValueNum; i++ {
		tag := entityType.Field(i).Tag
		column := tag.Get("column")
		if column == "" {
			continue
		}
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

	fmt.Println(builder.String())

	sqlId := buildSqlId(constants.INTO)
	err := gobatis.RegisterSql(sqlId, builder.String())
	if err != nil {
		return 0, err
	}
	var ret int64
	selectRunner := sess.Insert(sqlId).Param(paramMap)
	err = selectRunner.Result(&ret)
	if err != nil {
		return 0, nil
	}
	insertId := selectRunner.LastInsertId()
	return insertId, nil
}

func (userMapper *BaseMapper[T]) SaveBatch(entities ...T) (int64, int64) {
	return 0, 0
}
func (userMapper *BaseMapper[T]) DeleteById(id any) int64 {
	return 0
}
func (userMapper *BaseMapper[T]) DeleteBatchIds(ids []any) int64 {
	return 0
}
func (userMapper *BaseMapper[T]) UpdateById(entity T) int64 {
	return 0
}
func (userMapper *BaseMapper[T]) SelectById(id any) (T, error) {
	queryWrapper := userMapper.init(nil)
	queryWrapper.Eq(constants.ID, strconv.Itoa(id.(int)))
	columns := userMapper.buildSelectColumns(queryWrapper)

	sqlId, sql, paramMap := userMapper.buildSelectSql(queryWrapper, columns, buildSelectSqlFirstPart)

	var entity T
	err := gobatis.RegisterSql(sqlId, sql)
	if err != nil {
		return entity, err
	}

	sess := userMapper.SessMgr.NewSession()

	err = sess.Select(sqlId).Param(paramMap).Result(&entity)
	if err != nil {
		return entity, err
	}

	// delete sqlId
	gobatis.UnregisterSql(sqlId)

	return entity, nil
}
func (userMapper *BaseMapper[T]) SelectBatchIds(ids []any) ([]T, error) {
	tableName := userMapper.getTableName()
	// 构建第一部分sql
	sqlFirstPart := buildSelectSqlFirstPart(constants.ASTERISK, tableName)
	var paramMap = map[string]any{}
	build := strings.Builder{}
	// 拼接sql
	build.WriteString(constants.SPACE + constants.WHERE + constants.SPACE + constants.ID +
		constants.SPACE + constants.In + constants.LEFT_BRACKET + constants.SPACE)

	for index, id := range ids {
		mapping := userMapper.getMappingSeq()
		paramMap[mapping] = strconv.Itoa(id.(int))
		if index == len(ids)-1 {
			build.WriteString(constants.HASH_LEFT_BRACE + mapping + constants.RIGHT_BRACE)
		} else {
			build.WriteString(constants.HASH_LEFT_BRACE + mapping + constants.RIGHT_BRACE + constants.COMMA)
		}
	}
	build.WriteString(constants.SPACE + constants.RIGHT_BRACKET)
	sqlId := buildSqlId(constants.SELECT)
	sql := sqlFirstPart + build.String()

	err := gobatis.RegisterSql(sqlId, sql)
	if err != nil {
		return nil, err
	}

	sess := userMapper.SessMgr.NewSession()
	var arr []T
	err = sess.Select(sqlId).Param(paramMap).Result(&arr)
	if err != nil {
		return nil, err
	}
	return arr, nil
}

func (userMapper *BaseMapper[T]) getMappingSeq() string {
	userMapper.ParamNameSeq = userMapper.ParamNameSeq + 1
	mapping := constants.MAPPING + strconv.Itoa(userMapper.ParamNameSeq)
	return mapping
}

func (userMapper *BaseMapper[T]) SelectOne(queryWrapper *QueryWrapper[T]) (T, error) {
	queryWrapper = userMapper.init(queryWrapper)

	columns := userMapper.buildSelectColumns(queryWrapper)

	sqlId, sql, paramMap := userMapper.buildSelectSql(queryWrapper, columns, buildSelectSqlFirstPart)

	var entity T
	err := gobatis.RegisterSql(sqlId, sql)
	if err != nil {
		return entity, err
	}

	sess := userMapper.SessMgr.NewSession()

	err = sess.Select(sqlId).Param(paramMap).Result(&entity)
	if err != nil {
		return entity, err
	}

	// delete sqlId
	gobatis.UnregisterSql(sqlId)
	return entity, nil
}

func (userMapper *BaseMapper[T]) SelectCount(queryWrapper *QueryWrapper[T]) (int64, error) {
	queryWrapper = userMapper.init(queryWrapper)

	sqlId, sql, paramMap := userMapper.buildSelectSql(queryWrapper, constants.COUNT, buildSelectSqlFirstPart)

	err := gobatis.RegisterSql(sqlId, sql)
	if err != nil {
		return 0, err
	}

	sess := userMapper.SessMgr.NewSession()
	var count int64
	err = sess.Select(sqlId).Param(paramMap).Result(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (userMapper *BaseMapper[T]) SelectList(queryWrapper *QueryWrapper[T]) ([]T, error) {
	// 初始化queryWrapper，如果queryWrapper是空的，需要初始化一个新的
	queryWrapper = userMapper.init(queryWrapper)

	// 构建需要查询的字段，如果查询条件没有传入的话，默认查询所有
	columns := userMapper.buildSelectColumns(queryWrapper)

	// 构建sql语句
	// sqlId：sql的key值
	// sql：需要执行的sql
	// paramMap：所有的参数封装在这个map中
	sqlId, sql, paramMap := userMapper.buildSelectSql(queryWrapper, columns, buildSelectSqlFirstPart)

	err := gobatis.RegisterSql(sqlId, sql)
	if err != nil {
		return nil, err
	}

	sess := userMapper.SessMgr.NewSession()
	var arr []T
	err = sess.Select(sqlId).Param(paramMap).Result(&arr)
	if err != nil {
		return nil, err
	}

	// delete sqlId
	gobatis.UnregisterSql(sqlId)
	return arr, nil
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
	expression := queryWrapper.Expression
	build := strings.Builder{}
	// 遍历所有的条件参数
	for _, v := range expression {
		// 如果是ParamValue的话，通过#{} 拼接数据，并且把value值存储到paramMap中，方便后面使用
		// ParamValue 存储的是查询条件具体的值
		if paramValue, ok := v.(ParamValue); ok {
			mapping := userMapper.getMappingSeq()
			paramMap[mapping] = paramValue.value
			build.WriteString(constants.HASH_LEFT_BRACE + mapping + constants.RIGHT_BRACE + constants.SPACE)
		} else {
			// 拼接字段和条件
			build.WriteString(v.(string) + constants.SPACE)
		}
	}
	return build.String(), paramMap
}

func (userMapper *BaseMapper[T]) buildSelectSql(queryWrapper *QueryWrapper[T], columns string, buildSqlFunc BuildSqlFunc) (string, string, map[string]any) {

	// 构建查询条件
	sqlCondition, paramMap := userMapper.buildCondition(queryWrapper)

	// 获取表名称
	tableName := userMapper.getTableName()

	// 构建sql语句
	sqlId := buildSqlId(constants.SELECT)

	// 执行函数构建sql的第一部分
	sqlFirstPart := buildSqlFunc(columns, tableName)

	var sql string
	if len(queryWrapper.Expression) > 0 {
		// 拼接sql语句
		sql = sqlFirstPart + constants.SPACE + constants.WHERE + constants.SPACE + sqlCondition
	} else {
		// 如果没有条件，查询所有数据
		sql = sqlFirstPart
	}

	return sqlId, sql, paramMap
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

func buildSelectSqlFirstPart(columns string, tableName string) string {
	return constants.SELECT + constants.SPACE + columns + constants.SPACE + constants.FROM + constants.SPACE + tableName
}
