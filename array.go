package gox

import "reflect"

// IsInArray 通用的判断数据是否在数组中
func IsInArray(val interface{}, array interface{}) (exists bool, index int) {
	exists = false
	index = -1

	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
				index = i
				exists = true
				return
			}
		}
	}

	return
}

func UnionArray(array1 interface{}, array2 interface{}) interface{} {
	s1 := reflect.ValueOf(array1)
	s2 := reflect.ValueOf(array2)
	if s1.Len() < s2.Len() {
		s1, s2 = s2, s1
	}

	cache := make(map[interface{}]struct{})

	typ := s1.Type()

	for i := 0; i < s1.Len(); i++ {
		cache[s1.Index(i).Interface()] = struct{}{}
	}

	rsp := reflect.MakeSlice(typ, 0, s1.Len())
	for i := 0; i < s2.Len(); i++ {
		if _, ok := cache[s2.Index(i).Interface()]; ok {
			rsp = reflect.Append(rsp, s2.Index(i))
		}
	}

	return rsp.Interface()
}
