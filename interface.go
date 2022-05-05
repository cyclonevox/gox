package gox

import `reflect`

func ToInterfaceSlice(slice interface{}) []interface{} {
	if reflect.Slice == reflect.TypeOf(slice).Kind() {
		value := reflect.ValueOf(slice)
		rsp := make([]interface{}, 0, value.Len())
		for i := 0; i < value.Len(); i++ {
			rsp = append(rsp, value.Index(i).Interface())
		}

		return rsp
	}

	return []interface{}{slice}
}
