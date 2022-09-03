package mapper

import (
	"fmt"
	"github.com/acmestack/gobatis-plus/pkg/constants"
	"github.com/acmestack/godkits/gox/stringsx"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type SqlBuilder[T any] struct {
	ParamNameSeq int
}

func (sqlBuilder *SqlBuilder[T]) BuildSelectSql(queryWrapper *QueryWrapper[T], columns string) (map[string]any, string, string) {
	// 构建需要查询的字段，如果查询条件没有传入的话，默认查询所有
	if stringsx.Empty(columns) {
		// eg: columnName1,columnName2,columnName3
		columns = sqlBuilder.buildSelectColumns(queryWrapper)
	}

	// 获取表名称
	tableName := sqlBuilder.getTableName()

	// 构建查询条件
	// eg: columnName1 = #{mapping1} and columnName2 = #{mapping1}
	sqlCondition, paramMap := sqlBuilder.buildCondition(queryWrapper.Conditions)

	// 构建sql
	// eg: SELECT * FROM WHERE columnName = #{mapping1} and columnName = #{mapping1}
	sql := sqlBuilder.onBuildSelectSql(columns, tableName, sqlCondition)

	// 构建sqlId
	sqlId := sqlBuilder.buildSqlId(constants.SELECT)
	return paramMap, sql, sqlId
}

func (sqlBuilder *SqlBuilder[T]) BuildInsertSql(entity ...T) (map[string]any, string, string) {
	// 获取表名
	tableName := sqlBuilder.getTableName()

	// 获取插入字段
	// eg：columnName1,columnName2,columnName3
	columns := sqlBuilder.buildInsertColumns()

	// 构建插入语句后半部分
	// eg：(#{mapping1},#{mapping2},#{mapping3})
	paramMap, columnMappings := sqlBuilder.buildInsertColumnMapping(entity...)

	// 构建sql
	sql := sqlBuilder.onBuildInsertSql(tableName, columns, columnMappings)

	// 构建sqlId
	sqlId := sqlBuilder.buildSqlId(constants.INSERT)
	return paramMap, sql, sqlId
}

func (sqlBuilder *SqlBuilder[T]) BuildUpdateSql(entity T, updateWrapper *UpdateWrapper[T]) (map[string]any, string, string) {
	tableName := sqlBuilder.getTableName()
	paramMap, columnMapping := sqlBuilder.buildUpdateColumnMapping(entity)

	// 构建查询条件
	// eg: columnName1 = #{mapping1} and columnName2 = #{mapping1}
	sqlCondition, paramConditionMap := sqlBuilder.buildCondition(updateWrapper.Conditions)
	for k, v := range paramConditionMap {
		paramMap[k] = v
	}
	// 构建更新语句
	sql := sqlBuilder.onBuildUpdateSql(tableName, columnMapping, sqlCondition)

	// 构建sqlId
	sqlId := sqlBuilder.buildSqlId(constants.UPDATE)
	return paramMap, sql, sqlId
}

func (sqlBuilder *SqlBuilder[T]) BuildDeleteSql(conditions []any) (map[string]any, string, string) {
	tableName := sqlBuilder.getTableName()

	conditionMapping, paramMap := sqlBuilder.buildCondition(conditions)

	sql := sqlBuilder.onBuildDeleteSql(tableName, conditionMapping)

	sqlId := sqlBuilder.buildSqlId(constants.DELETE)
	return paramMap, sql, sqlId
}

func (sqlBuilder *SqlBuilder[T]) onBuildSelectSql(columns string, tableName string, sqlCondition string) string {
	sql := strings.Replace(constants.SELECT_SQL, constants.COLUMN_HASH, columns, -1)
	sql = strings.Replace(sql, constants.TABLE_NAME_HASH, tableName, -1)
	sql = strings.Replace(sql, constants.CONDITIONS_HASH, sqlCondition, -1)
	if sqlCondition == "" {
		sql = strings.Replace(sql, constants.WHERE, "", -1)
	}
	return sql
}

func (sqlBuilder *SqlBuilder[T]) onBuildInsertSql(tableName string, columns string, columnMappings []string) string {
	sql := strings.Replace(constants.INSERT_SQL, constants.TABLE_NAME_HASH, tableName, -1)
	sql = strings.Replace(sql, constants.COLUMN_HASH, columns, -1)
	sql = strings.Replace(sql, constants.COLUMN_MAPPING_HASH, columnMappings[0], -1)

	builder := stringsx.Builder{}
	builder.JoinString(sql)
	for i, columnMapping := range columnMappings {
		// 跳过第一次，因为上面已经使用了
		if i == 0 {
			continue
		}
		builder.JoinString(constants.COMMA + constants.LEFT_BRACKET + columnMapping + constants.RIGHT_BRACKET)
	}
	return builder.String()
}

func (sqlBuilder *SqlBuilder[T]) onBuildUpdateSql(tableName string, columnMapping string, sqlCondition string) string {
	sql := strings.Replace(constants.UPDATEBYID_SQL, constants.TABLE_NAME_HASH, tableName, -1)
	sql = strings.Replace(sql, constants.COLUMN_MAPPING_HASH, columnMapping, -1)
	sql = strings.Replace(sql, constants.CONDITIONS_HASH, sqlCondition, -1)
	return sql
}

func (sqlBuilder *SqlBuilder[T]) onBuildDeleteSql(tableName string, conditionMapping string) string {
	sql := strings.Replace(constants.DELETEBYID_SQL, constants.TABLE_NAME_HASH, tableName, -1)
	sql = strings.Replace(sql, constants.CONDITIONS_HASH, conditionMapping, -1)
	return sql
}

func (sqlBuilder *SqlBuilder[T]) buildSelectColumns(queryWrapper *QueryWrapper[T]) string {
	var columns string
	if len(queryWrapper.Columns) > 0 {
		columns = strings.Join(queryWrapper.Columns, ",")
	} else {
		columns = constants.ASTERISK
	}
	return columns
}

func (sqlBuilder *SqlBuilder[T]) buildInsertColumns() string {
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

func (sqlBuilder *SqlBuilder[T]) buildInsertColumnMapping(entities ...T) (map[string]any, []string) {
	var paramMap = map[string]any{}
	var allColumnMappings []string
	for _, entity := range entities {
		entityType := reflect.TypeOf(entity)
		entityValue := reflect.ValueOf(entity)
		entityValueNum := entityValue.NumField()
		var columnMappings []string
		for i := 0; i < entityValueNum; i++ {
			tag := entityType.Field(i).Tag
			column := tag.Get(constants.COLUMN)
			if column == "" {
				continue
			}
			// 构建columnMapping值
			v := entityValue.Field(i)
			mapping := sqlBuilder.getMappingSeq()
			switch iv := v.Interface().(type) {
			case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
				paramMap[mapping] = fmt.Sprintf("%d", iv)
			case string:
				paramMap[mapping] = iv
			}
			mapping = constants.HASH_LEFT_BRACE + mapping + constants.RIGHT_BRACE
			columnMappings = append(columnMappings, mapping)
		}
		allColumnMappings = append(allColumnMappings, strings.Join(columnMappings, ","))
	}

	return paramMap, allColumnMappings
}

func (sqlBuilder *SqlBuilder[T]) buildUpdateColumnMapping(entity T) (map[string]any, string) {
	entityType := reflect.TypeOf(entity)
	entityValue := reflect.ValueOf(entity)
	numField := entityType.NumField()
	paramMap := map[string]any{}
	var columnMappings []string
	for i := 0; i < numField; i++ {
		tag := entityType.Field(i).Tag
		column := tag.Get(constants.COLUMN)
		// 如果是等于Id的话也需要跳过
		if column == "" || constants.ID == column {
			continue
		}
		fieldValue := entityValue.Field(i).Interface()
		var mapping string
		switch v := fieldValue.(type) {
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
			idStr := fmt.Sprintf("%d", v)
			mapping = sqlBuilder.getMappingSeq()
			paramMap[mapping] = idStr
		case string:
			mapping = sqlBuilder.getMappingSeq()
			paramMap[mapping] = v
		}
		var columnMapping = column + constants.Eq + constants.HASH_LEFT_BRACE + mapping + constants.RIGHT_BRACE
		columnMappings = append(columnMappings, columnMapping)
	}
	str := strings.Join(columnMappings, ",")
	return paramMap, str
}

func (sqlBuilder *SqlBuilder[T]) getTableName() string {
	entityRef := reflect.TypeOf(new(T)).Elem()
	tableNameTag := entityRef.Field(0).Tag
	tableName := string(tableNameTag)
	return tableName
}

func (sqlBuilder *SqlBuilder[T]) buildCondition(conditions []any) (string, map[string]any) {
	var paramMap = map[string]any{}
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
					mapping := sqlBuilder.getMappingSeq()
					switch iv := elemV.Interface().(type) {
					case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
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
				mapping := sqlBuilder.getMappingSeq()
				switch iv := paramValue.value.(type) {
				case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
					paramMap[mapping] = fmt.Sprintf("%d", iv)
				case string:
					paramMap[mapping] = iv
				}
				build.WriteString(constants.HASH_LEFT_BRACE + mapping + constants.RIGHT_BRACE + constants.SPACE)
			}
		} else {
			build.WriteString(v.(string) + constants.SPACE)
		}
	}
	return build.String(), paramMap
}

func (sqlBuilder *SqlBuilder[T]) getMappingSeq() string {
	sqlBuilder.ParamNameSeq = sqlBuilder.ParamNameSeq + 1
	mapping := constants.MAPPING + strconv.Itoa(sqlBuilder.ParamNameSeq)
	return mapping
}

func (sqlBuilder *SqlBuilder[T]) buildSqlId(sqlType string) string {
	sqlId := sqlType + constants.CONNECTION + strconv.Itoa(time.Now().Nanosecond())
	return sqlId
}
