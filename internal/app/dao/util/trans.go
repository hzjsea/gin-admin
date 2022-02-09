package util

import (
	"context"

	"github.com/LyricTian/gin-admin/v8/internal/app/contextx"
	"github.com/google/wire"
	"gorm.io/gorm"
)

var TransSet = wire.NewSet(wire.Struct(new(Trans), "*"))

type Trans struct {
	DB *gorm.DB
}

func (a *Trans) Exec(ctx context.Context, fn func(context.Context) error) error {
	// 如果是不需要事务的直接返回执行的函数
	if _, ok := contextx.FromTrans(ctx); ok {
		return fn(ctx)
	}

	// 如果是需要事务的返回一个锁表的db连接
	return a.DB.Transaction(func(db *gorm.DB) error {
		return fn(contextx.NewTrans(ctx, db))
	})
}
