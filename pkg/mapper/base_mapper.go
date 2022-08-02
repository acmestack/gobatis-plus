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
	"context"

	"github.com/acmestack/gobatis"
)

type BaseMapper[T any] struct {
	SessMgr *gobatis.SessionManager
	Ctx     context.Context
	Columns []string
}

func (userMapper *BaseMapper[T]) Insert(entity T) int64 {
	return 0
}

func (userMapper *BaseMapper[T]) InsertBatch(entities ...T) (int64, int64) {
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
func (userMapper *BaseMapper[T]) SelectById(id any) T {
	return *new(T)
}
func (userMapper *BaseMapper[T]) SelectBatchIds(ids []any) []T {
	var arr []T
	return arr
}
func (userMapper *BaseMapper[T]) SelectOne(entity T) T {
	return *new(T)
}
func (userMapper *BaseMapper[T]) SelectCount(entity T) int64 {
	return 0
}

func (userMapper *BaseMapper[T]) SelectList(queryWrapper *QueryWrapper[T]) ([]T, error) {
	if queryWrapper == nil {
		queryWrapper = &QueryWrapper[T]{}
		queryWrapper.init()
	}
	sess := userMapper.SessMgr.NewSession()
	var arr []T
	err := sess.Select(queryWrapper.SqlBuild.String()).Param(queryWrapper.Entity).Result(&arr)
	if err != nil {
		return nil, err
	}
	return arr, nil
}
