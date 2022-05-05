package concurrency

import (
	`fmt`
	`reflect`
	`sync`
)

type FillFunc func(i int)

// Fills 并发将slice切片（需要元素为相同类型）里的每一个元素，作为fun的参数并调用fun
func Fills(slice interface{}, fun FillFunc) {
	value := reflect.ValueOf(slice)
	if value.Kind() != reflect.Slice {
		return
	}

	var wg sync.WaitGroup

	wg.Add(value.Len())
	for i := 0; i < value.Len(); i++ {
		go func(i int) {
			if r := recover(); r != nil {
				fmt.Println("Fills方法发生panic")
			}

			fun(i)

			wg.Done()
		}(i)
	}
	wg.Wait()
}
