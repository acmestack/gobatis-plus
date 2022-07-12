package query

type query[Children any, T any, R any] interface {
	Select(columns ...R) Children

	GetSqlSelect() string
}
