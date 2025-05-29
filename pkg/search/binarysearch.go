package search

func BinarySearch(list []uint, item uint) bool {
	if index := search(list, item); index != -1 {
		return true
	}
	return false
}
func search(list []uint, target uint) int {
	low := 0
	high := len(list) - 1

	for low <= high {
		mid := (low + high) / 2

		if list[mid] == target {
			return mid
		}

		if list[mid] < target {
			low = mid + 1
			continue
		}

		if list[mid] > target {
			high = mid - 1
		}
	}

	return -1
}
