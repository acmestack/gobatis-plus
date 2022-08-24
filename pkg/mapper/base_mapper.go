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

func (userMapper *BaseMapper[T]) Save(entity T) (int, int64, error) {
	sess := userMapper.SessMgr.NewSession()

	firstPartBuilder := userMapper.buildInsertFirstPart()

	var paramMap = map[string]any{}
	endPartBuilder := userMapper.buildInsertEndPart(entity, paramMap)

	insertSql := firstPartBuilder.String() + endPartBuilder.String()
	fmt.Println(insertSql)

	sqlId := buildSqlId(constants.INTO)

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

	// eg：insert field1,field2 into tableName values
	builder := userMapper.buildInsertFirstPart()

	var paramMap = map[string]any{}
	for i, entity := range entities {

		// eg：(#{mapping1},#{mapping2},#{mapping3})
		endPartBuilder := userMapper.buildInsertEndPart(entity, paramMap)
		if i != len(entities)-1 {
			builder.WriteString(endPartBuilder.String() + constants.COMMA)
		} else {
			builder.WriteString(endPartBuilder.String())
		}
	}
	fmt.Println(builder.String())

	sqlId := buildSqlId(constants.INTO)

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

func (userMapper *BaseMapper[T]) DeleteBatchIds(ids []any) (int64, error) {
	tableName := userMapper.getTableName()
	builder := strings.Builder{}
	builder.WriteString(constants.DELETE + constants.SPACE + constants.FROM + constants.SPACE + tableName + constants.SPACE +
		constants.WHERE + constants.SPACE + constants.ID + constants.SPACE + constants.In + constants.SPACE + constants.LEFT_BRACKET)
	var paramMap = map[string]any{}

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

	// field = #{mapping}，field = #{mapping}，field = #{mapping}
	var idValue any
	for i := 0; i < numField; i++ {
		tag := entityType.Field(i).Tag
		column := tag.Get(constants.COLUMN)
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

	sqlId := buildSqlId(constants.UPDATE)

	fmt.Println(builder.String())

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

	sqlFirstPart := buildSelectSqlFirstPart(constants.ASTERISK, tableName)
	var paramMap = map[string]any{}
	build := strings.Builder{}

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

	queryWrapper = userMapper.init(queryWrapper)

	columns := userMapper.buildSelectColumns(queryWrapper)

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

func (userMapper *BaseMapper[T]) buildInsertEndPart(entity T, paramMap map[string]any) strings.Builder {
	builder := strings.Builder{}
	builder.WriteString(constants.LEFT_BRACKET)
	entityType := reflect.TypeOf(entity)
	entityValue := reflect.ValueOf(entity)
	entityValueNum := entityValue.NumField()
	for i := 0; i < entityValueNum; i++ {
		tag := entityType.Field(i).Tag
		column := tag.Get(constants.COLUMN)
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
	return builder
}

func (userMapper *BaseMapper[T]) buildInsertFirstPart() *strings.Builder {
	builder := strings.Builder{}
	tableName := userMapper.getTableName()
	builder.WriteString(constants.INSERT + constants.SPACE + constants.INTO + constants.SPACE + tableName + constants.SPACE + constants.LEFT_BRACKET)
	entityType := reflect.TypeOf(new(T)).Elem()
	entityTypeNum := entityType.NumField()
	for i := 0; i < entityTypeNum; i++ {
		tag := entityType.Field(i).Tag
		column := tag.Get(constants.COLUMN)
		if column == "" {
			continue
		}
		if i != entityTypeNum-1 {
			builder.WriteString(column + constants.COMMA)
		} else {
			builder.WriteString(column + constants.RIGHT_BRACKET)
		}
	}
	builder.WriteString(constants.SPACE + constants.VALUES + constants.SPACE)
	return &builder
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

func (userMapper *BaseMapper[T]) buildCondition(queryWrapper *QueryWrapper[T]) (string, map[string]any) {
	var paramMap = map[string]any{}
	expression := queryWrapper.Expression
	build := strings.Builder{}

	for _, v := range expression {
		if paramValue, ok := v.(ParamValue); ok {
			mapping := userMapper.getMappingSeq()
			paramMap[mapping] = paramValue.value
			build.WriteString(constants.HASH_LEFT_BRACE + mapping + constants.RIGHT_BRACE + constants.SPACE)
		} else {
			build.WriteString(v.(string) + constants.SPACE)
		}
	}
	return build.String(), paramMap
}

func (userMapper *BaseMapper[T]) buildSelectSql(queryWrapper *QueryWrapper[T], columns string, buildSqlFunc BuildSqlFunc) (string, string, map[string]any) {

	sqlCondition, paramMap := userMapper.buildCondition(queryWrapper)

	tableName := userMapper.getTableName()

	sqlId := buildSqlId(constants.SELECT)

	sqlFirstPart := buildSqlFunc(columns, tableName)

	var sql string
	if len(queryWrapper.Expression) > 0 {

		sql = sqlFirstPart + constants.SPACE + constants.WHERE + constants.SPACE + sqlCondition
	} else {
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
