package slice

func DoesExist(list []uint, value uint) bool {

	for _, val := range list {
		if val == value {

			return true
		}
	}

	return false
}
