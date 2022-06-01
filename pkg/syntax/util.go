package syntax

// containsRune - if one character (rune) inside a list
func containsRune(ch rune, list []rune) bool {
	for _, item := range list {
		if item == ch {
			return true
		}
	}

	return false
}

// ContainsRune -
func ContainsRune(ch rune, list []rune) bool {
	return containsRune(ch, list)
}

// ContainsInt - if one character (int) inside a list
func ContainsInt(input int, list []int) bool {
	for _, item := range list {
		if item == input {
			return true
		}
	}

	return false
}
