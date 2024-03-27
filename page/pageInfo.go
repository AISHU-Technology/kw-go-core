package page

import (
	"github.com/acmestack/gorm-plus/gplus"
)

type Page[T any] struct {
	PageNo   int //当期页码
	PageSize int
	Count    int64
	Pages    int
	Data     []*T
	//RecordsMap []T
	IsLastPage  bool //是否最后一页
	NextPage    int  //下一页页码
	PrePage     int  //上一页
	IsFirstPage bool //是否第一页
}

// gplus分页
func PageGpInfo[T any](page gplus.Page[T]) *Page[T] {
	pageInfo := &Page[T]{PageNo: page.Current, PageSize: page.Size, Count: page.Total, Data: page.Records}
	return PageInfos[T](pageInfo)
}

func PageInfo[T any](PageNo int, PageSize int, count int64, data []*T) *Page[T] {
	pageInfo := &Page[T]{PageNo: PageNo, PageSize: PageSize, Count: count, Data: data}
	return PageInfos[T](pageInfo)
}

func PageInfos[T any](page *Page[T]) *Page[T] {
	if page.Count == 0 {
		page.IsLastPage = true
		page.NextPage = page.PageNo
		page.PrePage = page.PageNo
		page.IsFirstPage = true
		return page
	}
	var count = int(page.Count)
	pages := count / page.PageSize
	page.Pages = pages
	if count%page.PageSize != 0 {
		pages++
		page.Pages = pages
	}
	if page.PageNo > 1 {
		page.PrePage = page.PrePage - 1
		page.IsFirstPage = false
	} else {
		page.PrePage = 1
		page.IsFirstPage = true
	}
	if page.Pages > page.PageNo {
		page.NextPage = page.NextPage + 1
		page.IsLastPage = false
	} else {
		page.NextPage = page.Pages
		page.IsLastPage = true
	}
	return page
}
