package gpa

import (
	"context"
	"gorm.io/gorm"
)

const TransactionContextKey = "$__TX_CTX_KEY__$"

func GetTxFromContext(ctx context.Context) (tx *gorm.DB, found bool) {
	value := ctx.Value(TransactionContextKey)
	if tx, found = value.(*gorm.DB); found && tx != nil {
		return tx, true
	}
	return nil, false
}

type Hook func(c *QueryContext) error
type QueryContext struct {
	Dest   interface{}
	Tx     *gorm.DB
	hooks  []Hook
	action func(tx *gorm.DB) error
	curr   int
}

func (c *QueryContext) Next() error {
	c.curr++
	if c.curr == len(c.hooks) {
		return c.action(c.Tx)
	}
	return c.hooks[c.curr](c)
}

func (c *QueryContext) Execute(action func(tx *gorm.DB) error) error {
	c.action = action
	return c.Next()
}

func (c *QueryContext) Use(hooks ...Hook) {
	c.hooks = append(c.hooks, hooks...)
}

func (c *QueryContext) AddStatement(f func(db *gorm.DB) *gorm.DB) {
	c.Tx = f(c.Tx)
}
