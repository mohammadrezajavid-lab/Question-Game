package search

import "golang.project/go-fundamentals/gameapp/param/presenceparam"

func BinarySearch(list []presenceparam.PresenceItem, userId uint) (presenceparam.PresenceItem, bool) {
	if index := search(list, userId); index != -1 {
		return list[index], true
	}

	return presenceparam.PresenceItem{}, false
}
func search(list []presenceparam.PresenceItem, userId uint) int {
	low := 0
	high := len(list) - 1

	for low <= high {
		mid := (low + high) / 2

		if list[mid].UserId == userId {
			return mid
		}

		if list[mid].UserId < userId {
			low = mid + 1
			continue
		}

		if list[mid].UserId > userId {
			high = mid - 1
		}
	}

	return -1
}
