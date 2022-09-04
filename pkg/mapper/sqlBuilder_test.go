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
	"github.com/acmestack/gobatis-plus/pkg/constants"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSqlBuilder_BuildSelectSql(t *testing.T) {
	type args struct {
		queryWrapper *QueryWrapper[TestTable]
		columns      string
	}
	queryWrapper := &QueryWrapper[TestTable]{}
	queryWrapper.Eq("username", "acmestack").In("password", "123456", "pw5")
	var wantParamMap = make(map[string]any)
	wantParamMap["mapping1"] = "acmestack"
	wantParamMap["mapping2"] = "123456"
	wantParamMap["mapping3"] = "pw5"
	tests := []struct {
		name         string
		args         args
		wantParamMap map[string]any
		wantSql      string
		wantSqlId    string
	}{
		{
			name:         "buildSelectSql",
			args:         args{queryWrapper: queryWrapper, columns: ""},
			wantParamMap: wantParamMap,
			wantSql:      "SELECT * FROM test_table WHERE username = #{mapping1} and password in (#{mapping2},#{mapping3})",
			wantSqlId:    "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sqlBuilder := &SqlBuilder[TestTable]{}
			paramMap, sql, _ := sqlBuilder.BuildSelectSql(tt.args.queryWrapper, tt.args.columns)
			assert.Equal(t, tt.wantParamMap, paramMap, "they should be equal")
			assert.Equal(t, tt.wantSql, sql, "they should be equal")
		})
	}
}

func TestSqlBuilder_BuildInsertSql(t *testing.T) {
	table1 := TestTable{Username: "gobatis", Password: "123456"}
	table2 := TestTable{Username: "acmestack", Password: "654321"}
	var testTables []TestTable
	testTables = append(testTables, table1, table2)
	type args struct {
		entity []TestTable
	}
	var wantParamMap = make(map[string]any)
	wantParamMap["mapping1"] = "0"
	wantParamMap["mapping2"] = "gobatis"
	wantParamMap["mapping3"] = "123456"
	wantParamMap["mapping4"] = "0"
	wantParamMap["mapping5"] = "acmestack"
	wantParamMap["mapping6"] = "654321"

	tests := []struct {
		name         string
		args         args
		wantParamMap map[string]any
		wantSql      string
		wantSqlId    string
	}{
		{
			name:         "BuildInsertSql",
			args:         args{entity: testTables},
			wantParamMap: wantParamMap,
			wantSql:      "INSERT INTO test_table (id,username,password) VALUES (#{mapping1},#{mapping2},#{mapping3}),(#{mapping4},#{mapping5},#{mapping6})",
			wantSqlId:    "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sqlBuilder := &SqlBuilder[TestTable]{}
			paramMap, sql, _ := sqlBuilder.BuildInsertSql(tt.args.entity...)
			fmt.Println(paramMap)
			fmt.Println(sql)
			assert.Equal(t, tt.wantParamMap, paramMap, "they should be equal")
			assert.Equal(t, tt.wantSql, sql, "they should be equal")
		})
	}
}

func TestSqlBuilder_BuildUpdateSql(t *testing.T) {
	type args struct {
		entity        TestTable
		updateWrapper *UpdateWrapper[TestTable]
	}

	var entity = TestTable{Id: 1, Username: "gobatis", Password: "123456"}
	updateWrapper := &UpdateWrapper[TestTable]{}
	updateWrapper.Eq(constants.ID, 1)

	var wantParamMap = make(map[string]any)
	wantParamMap["mapping1"] = "gobatis"
	wantParamMap["mapping2"] = "123456"
	wantParamMap["mapping3"] = "1"
	tests := []struct {
		name         string
		args         args
		wantParamMap map[string]any
		wantSql      string
		wantSqlId    string
	}{
		{
			name:         "BuildUpdateSql",
			args:         args{entity: entity, updateWrapper: updateWrapper},
			wantParamMap: wantParamMap,
			wantSql:      "UPDATE test_table SET username=#{mapping1},password=#{mapping2} WHERE id = #{mapping3} ",
			wantSqlId:    "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sqlBuilder := &SqlBuilder[TestTable]{}
			paramMap, sql, _ := sqlBuilder.BuildUpdateSql(tt.args.entity, tt.args.updateWrapper)
			fmt.Println(paramMap)
			fmt.Println(sql)
			assert.Equal(t, tt.wantParamMap, paramMap, "they should be equal")
			assert.Equal(t, tt.wantSql, sql, "they should be equal")
		})
	}
}

func TestSqlBuilder_BuildDeleteSql(t *testing.T) {
	type args struct {
		conditions []any
	}

	var conditions []any
	conditions = append(conditions, constants.ID)
	conditions = append(conditions, constants.Eq)
	conditions = append(conditions, ParamValue{1})

	var wantParamMap = make(map[string]any)
	wantParamMap["mapping1"] = "1"

	tests := []struct {
		name         string
		args         args
		wantParamMap map[string]any
		wantSql      string
		wantSqlId    string
	}{
		{
			name:         "BuildDeleteSql",
			args:         args{conditions: conditions},
			wantParamMap: wantParamMap,
			wantSql:      "delete from test_table where id = #{mapping1} ",
			wantSqlId:    "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sqlBuilder := &SqlBuilder[TestTable]{}
			paramMap, sql, _ := sqlBuilder.BuildDeleteSql(tt.args.conditions)
			fmt.Println(paramMap)
			fmt.Println(sql)
			assert.Equal(t, tt.wantParamMap, paramMap, "they should be equal")
			assert.Equal(t, tt.wantSql, sql, "they should be equal")
		})
	}
}
