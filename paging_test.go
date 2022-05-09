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

	sortedData := NewManualPaging(originData, Paging{
		Page:    2,
		PerPage: 3,
	}).Sort("sum DESC,count DESC,proportion DESC,id ASC")

	results := sortedData.GetInterface().([]test)
	fmt.Println("Slice By Sort:")
	for _, result := range results {
		fmt.Println(result)
	}

	for index := range expected {
		if results[index] != expected[index] {
			t.Fatalf("sort(),第%v元素,期望：%v，实际：%v", index, expected[index], results[index])
		}
	}

	pagingResults := sortedData.Paging().([]test)
	expectedPaging := []test{
		{id: 4, sum: 76, count: 3, proportion: 0.4},
		{id: 2, sum: 70, count: 2, proportion: 0.6},
		{id: 5, sum: 70, count: 1, proportion: 0.8},
	}

	fmt.Println("Slice By Sort And Paging:")
	for _, result := range pagingResults {
		fmt.Println(result)
	}

	for index := range expectedPaging {
		if pagingResults[index] != expectedPaging[index] {
			t.Fatalf("sort(),第%v元素,期望：%v，实际：%v", index, pagingResults[index], expectedPaging[index])
		}
	}

}
