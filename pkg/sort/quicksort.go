package sort

import (
	"slices"
)

type QuickSort struct {
	list []uint
}

func NewQuickSort(list []uint) *QuickSort {
	return &QuickSort{list: list}
}

func (q *QuickSort) Sort() []uint {
	q.sort(0, len(q.list)-1)
	return q.list
}

func (q *QuickSort) sort(left, right int) {
	if right <= left {
		return
	}
	pivotIndex := q.partition(left, right)
	q.sort(left, pivotIndex-1)
	q.sort(pivotIndex+1, right)
}

func (q *QuickSort) partition(left, right int) int {
	pivotIndex := selectPivot(q.list[left:right+1]) + left
	q.swap(pivotIndex, right)
	pivot := q.list[right]
	leftP := left
	rightP := right - 1

	for leftP <= rightP {
		for leftP <= rightP && q.list[leftP] <= pivot {
			leftP++
		}
		for leftP <= rightP && q.list[rightP] > pivot {
			rightP--
		}
		if leftP < rightP {
			q.swap(leftP, rightP)
		}
	}
	q.swap(leftP, right)
	return leftP
}

func (q *QuickSort) swap(index1, index2 int) {
	q.list[index1], q.list[index2] = q.list[index2], q.list[index1]
}

func selectPivot(array []uint) int {
	if len(array) <= 5 {
		sorted := mergeSort(array)
		return slices.Index(array, sorted[len(sorted)/2])
	}

	tmp := make([]uint, 0, (len(array)+4)/5)

	for i := 0; i < len(array); i += 5 {
		end := i + 5
		if end > len(array) {
			end = len(array)
		}
		chunk := array[i:end]
		sortedChunk := mergeSort(chunk)
		mid := len(sortedChunk) / 2
		tmp = append(tmp, sortedChunk[mid])
	}

	sortedTmp := mergeSort(tmp)
	pivotValue := sortedTmp[len(sortedTmp)/2]
	return slices.Index(array, pivotValue)
}
