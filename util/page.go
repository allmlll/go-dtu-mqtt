package util

import (
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Page struct {
	Page int `form:"page" binding:"required"`
	Size int `form:"size" binding:"required"`
}

func (p Page) Check() bool {
	return p.Page >= 1 && p.Size >= 1
}

func (p Page) GetOpts() *options.FindOptions {
	return options.Find().SetSkip(int64((p.Page - 1) * p.Size)).SetLimit(int64(p.Size))
}
