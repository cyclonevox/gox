package gox

import (
	`fmt`
	`reflect`
	`sort`
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

func (p Paging) manual() ManualPagingInfo {
	return ManualPagingInfo{
		Page:    p.Page,
		PerPage: p.PerPage,
	}
}

func (s SortBy) OrderBy() string {
	return s.sorters().orderBy()
}

func (s SortBy) Sorters() []elastic.Sorter {
	return s.sorters().Sorters()
}

func (s SortBy) sorters() sorters {
	if strings.TrimSpace(string(s)) == "" {
		return nil
	}

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

// 排序字段转换 下划线转大驼峰
func (s sorters) fieldMarshal() sorters {
	for index, singleSorter := range s {
		field := singleSorter.Field
		if field == "" {
			continue
		}

		temp := strings.Split(field, "_")
		var str string
		for _, value := range temp {
			runeValue := []rune(value)
			if len(runeValue) > 0 {
				if runeValue[0] >= 'a' && runeValue[0] <= 'z' {
					// 首字母大写
					runeValue[0] -= 32
				}
				str += string(runeValue)
			}
		}

		s[index] = sorter{
			Field:     str,
			Ascending: singleSorter.Ascending,
		}
	}

	return s
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

type manualPaging struct {
	slice      interface{}
	typ        reflect.Type
	sliceValue reflect.Value
	sorters    sorters
	info       ManualPagingInfo
}

type ManualPagingInfo struct {
	Page    int
	PerPage int
}

func NewManualPaging(slice interface{}, paging Paging) *manualPaging {
	typ := reflect.TypeOf(slice)
	if typ.Kind() != reflect.Slice {
		panic(typ.Name() + "不是切片")
	}

	return &manualPaging{
		typ:        typ,
		slice:      slice,
		sliceValue: reflect.ValueOf(slice),
		info:       paging.manual(),
	}
}

func (mp *manualPaging) Sort(sortBy SortBy) *manualPaging {
	sorters := sortBy.sorters().fieldMarshal()
	if len(sorters) == 0 {
		panic(sortBy + ":缺少合法的排序字段")
	}
	mp.sorters = sorters

	sort.Slice(mp.slice, func(i, j int) bool {
		return mp.compare(i, j, 0)
	})

	return mp
}

func structValue(value reflect.Value) reflect.Value {
	if reflect.Ptr == value.Kind() {
		return value.Elem()
	}

	return value
}

func (mp *manualPaging) compare(i, j, sorterIndex int) bool {
	if sorterIndex > len(mp.sorters) {
		// It means two same values when function runs in this judge true
		return false
	}
	firstValue := structValue(mp.sliceValue.Index(i)).FieldByName(mp.sorters[sorterIndex].Field)
	secondValue := structValue(mp.sliceValue.Index(j)).FieldByName(mp.sorters[sorterIndex].Field)

	if firstValue.Kind() != secondValue.Kind() {
		panic("cannot sort the values with different types")
	}

	switch firstValue.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64:
		if firstValue.Int() == secondValue.Int() {
			sorterIndex++
			return mp.compare(i, j, sorterIndex)
		}

		if mp.sorters[sorterIndex].Ascending {
			b := firstValue.Int() < secondValue.Int()
			return b
		} else {
			b := firstValue.Int() > secondValue.Int()
			return b
		}

	case reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		if firstValue.Uint() == secondValue.Uint() {
			sorterIndex++
			return mp.compare(i, j, sorterIndex)
		}

		if mp.sorters[sorterIndex].Ascending {
			return firstValue.Uint() < secondValue.Uint()
		} else {
			return firstValue.Uint() > secondValue.Uint()
		}

	case reflect.Float32, reflect.Float64:
		if firstValue.Float() == secondValue.Float() {
			sorterIndex++
			return mp.compare(i, j, sorterIndex)
		}

		if mp.sorters[sorterIndex].Ascending {
			return firstValue.Float() < secondValue.Float()
		} else {
			return firstValue.Float() > secondValue.Float()
		}

	case reflect.Bool:
		if firstValue.Bool() == secondValue.Bool() {
			sorterIndex++
			return mp.compare(i, j, sorterIndex)
		}

		if mp.sorters[sorterIndex].Ascending {
			return firstValue.Bool()
		} else {
			return secondValue.Bool()
		}
	default:
		panic(mp.sorters[sorterIndex].Field + ":sort function just supports int float bool")
	}
}

// GetInterface 获取slice结果
func (mp *manualPaging) GetInterface() interface{} {
	return mp.sliceValue.Interface()
}

// Paging 手动分页操作
func (mp *manualPaging) Paging() interface{} {
	var (
		length int
		rsp    = reflect.MakeSlice(mp.typ, 0, mp.info.PerPage)
	)

	if length = mp.sliceValue.Len(); length == 0 {
		return mp.sliceValue.Interface()
	}

	start := (mp.info.Page - 1) * mp.info.PerPage
	if start < length {
		for i := 0; start+i < length && i < mp.info.PerPage; i++ {
			rsp = reflect.Append(rsp, mp.sliceValue.Index(start+i))
		}
	}

	return rsp.Interface()
}
