package page

import (
	"github.com/acmestack/gorm-plus/gplus"
	"gorm.io/gorm"
)

func NewPage[T any](PageNo, PageSize int) *Page[T] {
	pageInfo := gplus.NewPage[T](PageNo, PageSize)
	return PageGpInfo[T](*pageInfo)
}

func SelectPageOpts[T any](PageNo int, PageSize int, q *gplus.QueryCond[T], opts ...gplus.OptionFunc) (*Page[T], *gorm.DB) {
	pageInfo := gplus.NewPage[T](PageNo, PageSize)
	result, db := gplus.SelectPage(pageInfo, q, opts...)
	return PageGpInfo[T](*result), db
}

func SelectPage[T any](PageNo int, PageSize int, q *gplus.QueryCond[T]) (*Page[T], *gorm.DB) {
	pageInfo := gplus.NewPage[T](PageNo, PageSize)
	result, db := gplus.SelectPage(pageInfo, q, gplus.IgnoreTotal())
	return PageGpInfo[T](*result), db
}
