package gox

import (
	`encoding/json`
)

type (
	pageData struct {
		// CurrentPage 当前页码
		CurrentPage int `json:"currentPage"`
		// HasNext 是否还有下一页数据
		HasNext bool `json:"hasNext"`
		// HasPrev 是否有上一页数据
		HasPrev bool `json:"hasPrev"`
		// TotalNum 总共数量
		TotalNum int64 `json:"totalNum"`
		// TotalPage 总共页数
		TotalPage int64 `json:"totalPage"`
		// Items 数据列表
		Items interface{} `json:"items"`
		// Extras 额外数据
		Extras []extraPageData `json:"extras"`
	}
	extraPageData struct {
		// Key 键
		Key string
		// Value 值
		Value interface{}
	}
)

// NewPage 生成新的分页数据对象
func NewPage(items interface{}, totalNum int64, perPage int, page int, extras ...extraPageData) *pageData {
	totalPage := totalNum / int64(perPage)
	if (totalNum % int64(perPage)) > 0 {
		totalPage += 1
	}

	hasPrev := false
	if page > 1 {
		hasPrev = true
	}

	hasNext := false
	if int64(page) < totalPage {
		hasNext = true
	}

	return &pageData{
		CurrentPage: page,
		HasNext:     hasNext,
		HasPrev:     hasPrev,
		TotalNum:    totalNum,
		TotalPage:   totalPage,
		Items:       items,
		Extras:      extras,
	}
}

func NewPageExtra(key string, value interface{}) *extraPageData {
	return &extraPageData{
		Key:   key,
		Value: value,
	}
}

func (pd pageData) MarshalJSON() ([]byte, error) {
	data := make(map[string]interface{}, 7)
	data["currentPage"] = pd.CurrentPage
	data["hasNext"] = pd.HasNext
	data["hasPrev"] = pd.HasPrev
	data["totalNum"] = pd.TotalNum
	data["items"] = pd.Items
	data["currentPage"] = pd.CurrentPage
	data["totalPage"] = pd.TotalPage
	for _, extra := range pd.Extras {
		data[extra.Key] = extra.Value
	}

	return json.Marshal(data)
}
