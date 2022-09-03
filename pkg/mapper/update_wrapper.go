package mapper

import "github.com/acmestack/gobatis-plus/pkg/constants"

type UpdateWrapper[T any] struct {
	Columns           []string
	ValuesMap         map[string]any
	Conditions        []any
	LastConditionType string
}

func (updateWrapper *UpdateWrapper[T]) Set(column string, val any) Wrapper[T] {
	updateWrapper.ValuesMap[column] = val
	return updateWrapper
}

func (updateWrapper *UpdateWrapper[T]) Eq(column string, val any) Wrapper[T] {
	updateWrapper.addCondition(column, val, constants.Eq)
	return updateWrapper
}

func (updateWrapper *UpdateWrapper[T]) Ne(column string, val any) Wrapper[T] {
	updateWrapper.addCondition(column, val, constants.Ne)
	return updateWrapper
}

func (updateWrapper *UpdateWrapper[T]) Gt(column string, val any) Wrapper[T] {
	updateWrapper.addCondition(column, val, constants.Gt)
	return updateWrapper
}

func (updateWrapper *UpdateWrapper[T]) Ge(column string, val any) Wrapper[T] {
	updateWrapper.addCondition(column, val, constants.Ge)
	return updateWrapper
}

func (updateWrapper *UpdateWrapper[T]) Lt(column string, val any) Wrapper[T] {
	updateWrapper.addCondition(column, val, constants.Lt)
	return updateWrapper
}

func (updateWrapper *UpdateWrapper[T]) Le(column string, val any) Wrapper[T] {
	updateWrapper.addCondition(column, val, constants.Le)
	return updateWrapper
}

func (updateWrapper *UpdateWrapper[T]) Like(column string, val any) Wrapper[T] {
	s := val.(string)
	updateWrapper.addCondition(column, "%"+s+"%", constants.Like)
	return updateWrapper
}

func (updateWrapper *UpdateWrapper[T]) NotLike(column string, val any) Wrapper[T] {
	s := val.(string)
	updateWrapper.addCondition(column, "%"+s+"%", constants.Not+constants.Like)
	return updateWrapper
}

func (updateWrapper *UpdateWrapper[T]) LikeLeft(column string, val any) Wrapper[T] {
	s := val.(string)
	updateWrapper.addCondition(column, "%"+s, constants.Like)
	return updateWrapper
}

func (updateWrapper *UpdateWrapper[T]) LikeRight(column string, val any) Wrapper[T] {
	s := val.(string)
	updateWrapper.addCondition(column, s+"%", constants.Like)
	return updateWrapper
}

func (updateWrapper *UpdateWrapper[T]) In(column string, val ...any) Wrapper[T] {
	updateWrapper.addCondition(column, val, constants.In)
	return updateWrapper
}

func (updateWrapper *UpdateWrapper[T]) And() Wrapper[T] {
	updateWrapper.Conditions = append(updateWrapper.Conditions, constants.Eq)
	updateWrapper.LastConditionType = constants.Eq
	return updateWrapper
}

func (updateWrapper *UpdateWrapper[T]) Or() Wrapper[T] {
	updateWrapper.Conditions = append(updateWrapper.Conditions, constants.Or)
	updateWrapper.LastConditionType = constants.Or
	return updateWrapper
}

func (updateWrapper *UpdateWrapper[T]) Select(columns ...string) Wrapper[T] {
	updateWrapper.Columns = append(updateWrapper.Columns, columns...)
	return updateWrapper
}

func (updateWrapper *UpdateWrapper[T]) addCondition(column string, val any, conditionType string) {

	if updateWrapper.LastConditionType != constants.And && updateWrapper.LastConditionType != constants.Or && len(updateWrapper.Conditions) > 0 {
		updateWrapper.Conditions = append(updateWrapper.Conditions, constants.And)
	}

	updateWrapper.Conditions = append(updateWrapper.Conditions, column)

	updateWrapper.Conditions = append(updateWrapper.Conditions, conditionType)

	updateWrapper.Conditions = append(updateWrapper.Conditions, ParamValue{val})
}
