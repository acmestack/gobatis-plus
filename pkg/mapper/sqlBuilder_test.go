package mapper

import (
	"fmt"
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
