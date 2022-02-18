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