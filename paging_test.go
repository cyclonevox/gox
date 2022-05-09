package gox

import (
	`fmt`
	`testing`
)

func TestManualPaging_Sort(t *testing.T) {
	type test struct {
		Id           int64
		CourseSum    int
		TeacherCount int
		Proportion   float32
	}

	originData := []test{
		{Id: 3, CourseSum: 76, TeacherCount: 3, Proportion: 0.4},
		{Id: 6, CourseSum: 86, TeacherCount: 2, Proportion: 0.3},
		{Id: 2, CourseSum: 70, TeacherCount: 2, Proportion: 0.6},
		{Id: 4, CourseSum: 76, TeacherCount: 3, Proportion: 0.4},
		{Id: 5, CourseSum: 70, TeacherCount: 1, Proportion: 0.8},
		{Id: 1, CourseSum: 86, TeacherCount: 2, Proportion: 0.5},
	}

	expected := []test{
		{Id: 1, CourseSum: 86, TeacherCount: 2, Proportion: 0.5},
		{Id: 6, CourseSum: 86, TeacherCount: 2, Proportion: 0.3},
		{Id: 3, CourseSum: 76, TeacherCount: 3, Proportion: 0.4},
		{Id: 4, CourseSum: 76, TeacherCount: 3, Proportion: 0.4},
		{Id: 2, CourseSum: 70, TeacherCount: 2, Proportion: 0.6},
		{Id: 5, CourseSum: 70, TeacherCount: 1, Proportion: 0.8},
	}

	sortedData := NewManualPaging(originData, Paging{
		Page:    2,
		PerPage: 3,
	}).Sort("course_sum DESC,teacher_count DESC,proportion DESC,id ASC")

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
		{Id: 4, CourseSum: 76, TeacherCount: 3, Proportion: 0.4},
		{Id: 2, CourseSum: 70, TeacherCount: 2, Proportion: 0.6},
		{Id: 5, CourseSum: 70, TeacherCount: 1, Proportion: 0.8},
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
