/*
 * Licensed to the AcmeStack under one or more contributor license
 * agreements. See the NOTICE file distributed with this work for
 * additional information regarding copyright ownership.
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
	"fmt"
	"github.com/acmestack/gobatis"
	"github.com/acmestack/gobatis-plus/pkg/constants"
	"reflect"
)

type BaseMapper[T any] struct {
	SessMgr *gobatis.SessionManager
}

func (userMapper *BaseMapper[T]) SelectList(queryWrapper *QueryWrapper[T]) ([]T, error) {
	// if queryWrapper is nil ,need to build a new queryWrapper
	queryWrapper = userMapper.initQueryWrapper(queryWrapper)

	builder := SqlBuilder[T]{}
	paramMap, sql, sqlId := builder.BuildSelectSql(queryWrapper, "")

	err := gobatis.RegisterSql(sqlId, sql)
	defer gobatis.UnregisterSql(sqlId)
	if err != nil {
		return nil, err
	}

	sess := userMapper.SessMgr.NewSession()
	var results []T
	err = sess.Select(sqlId).Param(paramMap).Result(&results)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (userMapper *BaseMapper[T]) SelectById(id any) (T, error) {
	queryWrapper := userMapper.initQueryWrapper(nil)
	switch v := id.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		queryWrapper.Eq(constants.ID, fmt.Sprintf("%d", v))
	case string:
		queryWrapper.Eq(constants.ID, v)
	}

	builder := SqlBuilder[T]{}
	paramMap, sql, sqlId := builder.BuildSelectSql(queryWrapper, "")

	err := gobatis.RegisterSql(sqlId, sql)
	defer gobatis.UnregisterSql(sqlId)
	var entity T
	if err != nil {
		return entity, err
	}

	sess := userMapper.SessMgr.NewSession()
	err = sess.Select(sqlId).Param(paramMap).Result(&entity)
	if err != nil {
		return entity, err
	}
	return entity, nil
}

func (userMapper *BaseMapper[T]) SelectBatchIds(ids []any) ([]T, error) {
	queryWrapper := userMapper.initQueryWrapper(nil)
	queryWrapper.In(constants.ID, ids...)

	builder := SqlBuilder[T]{}
	paramMap, sql, sqlId := builder.BuildSelectSql(queryWrapper, "")

	err := gobatis.RegisterSql(sqlId, sql)
	defer gobatis.UnregisterSql(sqlId)
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

func (userMapper *BaseMapper[T]) SelectOne(queryWrapper *QueryWrapper[T]) (T, error) {
	queryWrapper = userMapper.initQueryWrapper(queryWrapper)

	builder := SqlBuilder[T]{}
	paramMap, sql, sqlId := builder.BuildSelectSql(queryWrapper, "")

	err := gobatis.RegisterSql(sqlId, sql)
	defer gobatis.UnregisterSql(sqlId)
	var entity T
	if err != nil {
		return entity, err
	}

	sess := userMapper.SessMgr.NewSession()
	err = sess.Select(sqlId).Param(paramMap).Result(&entity)
	if err != nil {
		return entity, err
	}

	return entity, nil
}

func (userMapper *BaseMapper[T]) SelectCount(queryWrapper *QueryWrapper[T]) (int64, error) {

	queryWrapper = userMapper.initQueryWrapper(queryWrapper)

	builder := SqlBuilder[T]{}
	paramMap, sql, sqlId := builder.BuildSelectSql(queryWrapper, constants.COUNT)

	err := gobatis.RegisterSql(sqlId, sql)
	defer gobatis.UnregisterSql(sqlId)
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

func (userMapper *BaseMapper[T]) Save(entity T) (int, int64, error) {
	builder := SqlBuilder[T]{}
	paramMap, sql, sqlId := builder.BuildInsertSql(entity)

	err := gobatis.RegisterSql(sqlId, sql)
	defer gobatis.UnregisterSql(sqlId)
	if err != nil {
		return 0, 0, err
	}

	sess := userMapper.SessMgr.NewSession()
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
	builder := SqlBuilder[T]{}
	paramMap, sql, sqlId := builder.BuildInsertSql(entities...)

	err := gobatis.RegisterSql(sqlId, sql)
	if err != nil {
		return 0, 0, err
	}

	sess := userMapper.SessMgr.NewSession()
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
	var conditions []any
	conditions = append(conditions, constants.ID)
	conditions = append(conditions, constants.Eq)
	conditions = append(conditions, ParamValue{id})

	builder := SqlBuilder[T]{}
	paramMap, sql, sqlId := builder.BuildDeleteSql(conditions)

	err := gobatis.RegisterSql(sqlId, sql)
	defer gobatis.UnregisterSql(sqlId)
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
	var conditions []any
	conditions = append(conditions, constants.ID)
	conditions = append(conditions, constants.In)
	conditions = append(conditions, ParamValue{ids})

	builder := SqlBuilder[T]{}
	paramMap, sql, sqlId := builder.BuildDeleteSql(conditions)

	err := gobatis.RegisterSql(sqlId, sql)
	defer gobatis.UnregisterSql(sql)
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
	updateWrapper := userMapper.initUpdateWrapper(nil)
	value := userMapper.getIdValue(entity)
	updateWrapper.Eq(constants.ID, value)

	builder := SqlBuilder[T]{}
	paramMap, sql, sqlId := builder.BuildUpdateSql(entity, updateWrapper)

	sess := userMapper.SessMgr.NewSession()

	err := gobatis.RegisterSql(sqlId, sql)
	defer gobatis.UnregisterSql(sqlId)
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

func (userMapper *BaseMapper[T]) initQueryWrapper(queryWrapper *QueryWrapper[T]) *QueryWrapper[T] {
	if queryWrapper == nil {
		queryWrapper = &QueryWrapper[T]{}
	}
	return queryWrapper
}

func (userMapper *BaseMapper[T]) initUpdateWrapper(updateWrapper *UpdateWrapper[T]) *UpdateWrapper[T] {
	if updateWrapper == nil {
		updateWrapper = &UpdateWrapper[T]{}
	}
	return updateWrapper
}

func (userMapper *BaseMapper[T]) getIdValue(entity T) any {
	entityType := reflect.TypeOf(entity)
	entityValue := reflect.ValueOf(entity)
	numField := entityType.NumField()
	for i := 0; i < numField; i++ {
		tag := entityType.Field(i).Tag
		column := tag.Get(constants.COLUMN)
		if constants.ID == column {
			return entityValue.Field(i).Interface()
		}
	}
	return nil
}
