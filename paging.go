package gox

import (
	`fmt`
	`strings`

	`github.com/olivere/elastic/v7`
)

type SortBy string

type sorters []sorter

type sorter struct {
	Field     string
	Ascending bool
}

// Paging 分页对象
type Paging struct {
	// 当前页
	Page int `default:"1" json:"page" param:"page" query:"page" form:"page" xml:"page" validate:"min=1"`
	// 每页个数
	PerPage int `default:"20" json:"perPage" param:"perPage" query:"perPage" form:"perPage" xml:"perPage" validate:"min=1"`
	// 查询关键字
	Keyword string `json:"keyword" param:"keyword" query:"keyword" form:"keyword" xml:"keyword" `
	// 排序顺序
	SortOrder string `default:"DESC" json:"sortOrder" param:"sortOrder" query:"sortOrder" form:"sortOrder" xml:"sortOrder"  validate:"oneof=asc ASC ascending ASCENDING desc DESC descending DESCENDING"`
}

func (s SortBy) OrderBy() string {
	return s.sorters().orderBy()
}

func (s SortBy) Sorters() []elastic.Sorter {
	return s.sorters().Sorters()
}

func (s SortBy) sorters() sorters {
	fields := strings.Split(string(s), ",")
	if len(fields) == 0 {
		return nil
	}

	sorters := make(sorters, 0, len(fields))
	for _, field := range fields {
		sortBy := strings.Split(field, " ")
		if len(sortBy) > 2 {
			continue
		}

		sorters = append(sorters, sorter{
			Field:     sortBy[0],
			Ascending: len(sortBy) == 2 && strings.ToUpper(sortBy[1]) == "ASC",
		})
	}

	return sorters
}

func (s sorters) orderBy() string {
	if len(s) == 0 {
		return ""
	}

	b := strings.Builder{}
	for _, sorter := range s {
		order := "DESC"
		if sorter.Ascending {
			order = "ASC"
		}

		b.WriteString(fmt.Sprintf("`%s` %s,", sorter.Field, order))
	}

	orderBy := b.String()

	return orderBy[:len(orderBy)-1]
}

func (s sorters) Sorters() []elastic.Sorter {
	sorters := make([]elastic.Sorter, 0, len(s))
	for _, sorter := range s {
		sorters = append(sorters, elastic.SortInfo{
			Field:     sorter.Field,
			Ascending: sorter.Ascending,
		})
	}

	return sorters
}

// OrderBy 排序字符串
func (p *Paging) OrderBy(sorter Sorter) string {
	order := "ASC"
	if strings.HasPrefix(strings.ToLower(p.SortOrder), "desc") {
		order = "DESC"
	}

	return fmt.Sprintf("`%s` %s", sorter.SortFieldName(), order)
}

// MySQL 获得MySQL需要的分页参数
func (p *Paging) MySQL() (start int, offset int) {
	return p.PerPage, (p.Page - 1) * p.PerPage
}

// Start 获得开始下标
func (p *Paging) Start() int {
	return (p.Page - 1) * p.PerPage
}

// Limit 获得限制个数
func (p *Paging) Limit() int {
	return p.PerPage
}
