package concurrency

import (
	`fmt`
	`reflect`
	`sync`
)

type FillFunc func(i int)

// Fills 并发将slice切片（需要元素为相同类型）里的每一个元素，作为fun的参数并调用fun
func Fills(slice interface{}, fun FillFunc, maxGoroutines ...int) {
	max := 15
	if len(maxGoroutines) != 0 && maxGoroutines[0] > 0 {
		max = maxGoroutines[0]
	}

	value := reflect.ValueOf(slice)
	if value.Kind() != reflect.Slice {
		return
	}

	if value.Len() == 0 {
		slice = reflect.MakeSlice(value.Type(), 0, 0).Interface()

		return
	}

	var (
		offset       int
		length       = value.Len()
		groups       = length/max + 1
		lastGroupNum = length % max
	)

	for i := 0; i < groups; i++ {
		groupNum := max
		if i == groups-1 {
			groupNum = lastGroupNum
		}

		wg := sync.WaitGroup{}
		wg.Add(groupNum)
		for j := 0; j < groupNum; j++ {
			go func(j int) {
				defer func() {
					if r := recover(); r != nil {
						fmt.Println("调用fills方法发生panic")
					}

					wg.Done()
				}()

				fun(j + offset)
			}(j)
		}
		wg.Wait()

		offset += groupNum
	}
}
