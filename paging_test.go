package gox

import (
	`fmt`
	`testing`
)

func TestManualPaging_Sort(t *testing.T) {
	type test struct {
		id         int64
		sum        int
		count      int
		proportion float32
	}

	originData := []test{
		{id: 3, sum: 76, count: 3, proportion: 0.4},
		{id: 6, sum: 86, count: 2, proportion: 0.3},
		{id: 2, sum: 70, count: 2, proportion: 0.6},
		{id: 4, sum: 76, count: 3, proportion: 0.4},
		{id: 5, sum: 70, count: 1, proportion: 0.8},
		{id: 1, sum: 86, count: 2, proportion: 0.5},
	}

	expected := []test{
		{id: 1, sum: 86, count: 2, proportion: 0.5},
		{id: 6, sum: 86, count: 2, proportion: 0.3},
		{id: 3, sum: 76, count: 3, proportion: 0.4},
		{id: 4, sum: 76, count: 3, proportion: 0.4},
		{id: 2, sum: 70, count: 2, proportion: 0.6},
		{id: 5, sum: 70, count: 1, proportion: 0.8},
	}

	sortedData := NewManualPaging(originData, ManualPagingInfo{
		page:    2,
		perPage: 5,
		SortBy:  "sum DESC,count DESC,proportion DESC,id ASC",
	}).Sort()

	results := sortedData.GetSliceInterface().([]test)
	for _, result := range results {
		fmt.Println(result)
	}

	for index := range expected {
		if results[index] != expected[index] {
			t.Fatalf("第%v元素,期望：%v，实际：%v", index, expected[index], results[index])
		}
	}

}
