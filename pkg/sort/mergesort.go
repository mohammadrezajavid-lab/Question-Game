package sort

/*func mergeSort(inputArray []uint) []uint {
	if len(inputArray) <= 2 {
		if len(inputArray) == 2 {
			if inputArray[1] < inputArray[0] {
				inputArray[0], inputArray[1] = inputArray[1], inputArray[0]
			}
		}
		return inputArray
	}

	// Divide
	mid := len(inputArray) / 2
	divideL := mergeSort(inputArray[:mid]) // T(n/2) -> each half has n/2 elements,
	// we have two recursive calls with input size as (n/2).
	divideR := mergeSort(inputArray[mid:]) // T(n/2)

	var result []uint
	var i, j = 0, 0

	// Conquer and Merge // O(n) -> for merge the two sorted halves
	for i < len(divideL) && j < len(divideR) {
		if divideL[i] < divideR[j] {
			result = append(result, divideL[i])
			i += 1
			continue
		}
		if divideR[j] < divideL[i] {
			result = append(result, divideR[j])
			j += 1
			continue
		}
		if divideL[i] == divideR[j] {
			result = append(result, divideL[i])
			i += 1
			continue
		}
	}
	if i < len(divideL) {
		result = append(result, divideL[i:]...)
	}
	if j < len(divideR) {
		result = append(result, divideR[j:]...)
	}

	return result
}*/

func mergeSort(inputArray []uint) []uint {
	if len(inputArray) <= 1 {
		return inputArray
	}

	mid := len(inputArray) / 2
	left := mergeSort(inputArray[:mid])
	right := mergeSort(inputArray[mid:])

	result := make([]uint, 0, len(inputArray))
	i, j := 0, 0
	for i < len(left) && j < len(right) {
		if left[i] <= right[j] {
			result = append(result, left[i])
			i++
		} else {
			result = append(result, right[j])
			j++
		}
	}
	result = append(result, left[i:]...)
	result = append(result, right[j:]...)
	return result
}
