package core

import `strings`

// RegionCodeCountry 全国
const RegionCodeCountry RegionCode = "100000"

const (
	// RegionTypeCountry 全国
	RegionTypeCountry RegionLevel = 0
	// RegionTypeProvince 省、直辖市、自治区、特别行政区
	RegionTypeProvince RegionLevel = 1
	// RegionTypeCity 市、直辖市、自治州、省直辖县
	RegionTypeCity RegionLevel = 2
	// RegionTypeCounty 区、旗、县
	RegionTypeCounty RegionLevel = 3
	// RegionTypeCustom 自定义类型
	RegionTypeCustom RegionLevel = 4
)

// RegionCode 行政区划码
type RegionCode string

// RegionLevel 行政区划等级
type RegionLevel int8

func (rc RegionCode) Prefix() string {
	code := string(rc)
	for {
		if strings.HasSuffix(code, "0") {
			code = code[:len(code)-1]

			continue
		}

		break
	}

	return code
}

func (rl RegionLevel) Name() string {
	switch rl {
	case RegionTypeCountry:
		return "全国"
	case RegionTypeProvince:
		return "省"
	case RegionTypeCity:
		return "市"
	case RegionTypeCounty:
		return "区"
	default:
		return "自定义"
	}
}
