package utils

func IsAlpha(s string) bool {
	for _, c := range s {
		if (c < 'A' || c > 'Z') && (c < 'a' || c > 'z') && c != ' ' {
			return false
		}
	}
	return true
}

func IsNumeric(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}
func IsAlphanumeric(s string) bool {
	for _, c := range s {
		if (c < '0' || c > '9') && (c < 'A' || c > 'Z') && (c < 'a' || c > 'z') && c != ' ' {
			return false
		}
	}
	return true
}
