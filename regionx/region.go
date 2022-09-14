package regionx

import (
	_ `embed`
	`encoding/json`
	`strings`

	`github.com/storezhang/gox/core`
)

var (
	//go:embed region.json
	regionData []byte

	// 行政区划缓存
	regionMap = make(map[core.RegionCode]Region, 0)
)

type Region struct {
	Code      core.RegionCode  `json:"code"`
	Name      string           `json:"name"`
	Type      core.RegionLevel `json:"type"`
	Longitude float64          `json:"longitude"`
	Latitude  float64          `json:"latitude"`
	Parents   []Region         `json:"parents,omitempty"`
	Children  []Region         `json:"children,omitempty"`
}

func Init(data []byte) {
	var (
		err     error
		regions []Region
	)

	if len(data) == 0 {
		data = regionData
	}

	if err = json.Unmarshal(data, &regions); err != nil {
		panic(err)
	}

	if len(regions) != 0 {
		set(regions[0], regionMap)
	}
}

func GetByCode(code interface{}) Region {
	switch code := code.(type) {
	case string:
		if code == "" {
			code = string(core.RegionCodeCountry)
		}

		return regionMap[core.RegionCode(code)]
	case core.RegionCode:
		if code == "" {
			code = core.RegionCodeCountry
		}

		return regionMap[code]
	}

	return Region{}
}

func ListChildren(code core.RegionCode) []Region {
	region := GetByCode(code)
	if len(region.Children) == 0 {
		return make([]Region, 0)
	}

	return region.Children
}

func IsValidCode(code string) bool {
	_, ok := regionMap[core.RegionCode(code)]

	return ok
}

func (r Region) FullName(sep ...string) string {
	names := make([]string, 0)
	if len(r.Parents) != 0 {
		for i := len(r.Parents) - 1; i >= 0; i-- {
			names = append(names, r.Parents[i].Name)
		}
	}

	names = append(names, r.Name)

	separator := "-"
	if len(sep) != 0 && sep[0] != "" {
		separator = sep[0]
	}

	return strings.Join(names, separator)
}

func (r Region) IsLeaf() bool {
	return len(r.Children) == 0
}

// Parent 返回上一级
func (r Region) Parent() Region {
	if len(r.Parents) != 0 {
		return r.Parents[0]
	}

	return Region{}
}

// Top 返回顶层（非全国）
func (r Region) Top() Region {
	if len(r.Parents) != 0 {
		return r.Parents[len(r.Parents)-1]
	}

	return r
}

func set(region Region, m map[core.RegionCode]Region) {
	d := Region{
		Code:      region.Code,
		Name:      region.Name,
		Type:      region.Type,
		Longitude: region.Longitude,
		Latitude:  region.Latitude,
		Parents:   region.Parents,
	}

	for _, child := range region.Children {
		d.Children = append(d.Children, Region{
			Code:      child.Code,
			Name:      child.Name,
			Type:      child.Type,
			Longitude: child.Longitude,
			Latitude:  child.Latitude,
		})

		set(child, m)
	}

	m[region.Code] = d
}
