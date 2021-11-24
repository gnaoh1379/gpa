package gpa

import "gorm.io/gorm"

type Option func(q *QueryContext)

func WithStatement(statement func(db *gorm.DB) *gorm.DB) Option {
	return func(q *QueryContext) {
		q.Tx = statement(q.Tx)
	}
}

type FieldNameGetter interface {
	FieldName(field string) string
}

func WithEqual(field string, value interface{}) Option {
	return func(q *QueryContext) {
		if getter, ok := q.Dest.(FieldNameGetter); ok {
			field = getter.FieldName(field)
		}
		q.AddStatement(func(db *gorm.DB) *gorm.DB {
			return db.Where(field+" = ?", value)
		})
	}
}
