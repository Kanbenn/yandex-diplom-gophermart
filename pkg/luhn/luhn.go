package luhn

// IsValidLuhnNumber returns true if the given number is compliant with the Luhn formula.
// it allows empty spaces in the input and skips them.
func IsValidLuhnNumber(in []byte) bool {
	var digit, result int
	idx := 1
	for i := len(in) - 1; i >= 0; i-- {
		if in[i] == ' ' {
			continue
		}
		if in[i] < '0' || in[i] > '9' {
			return false
		}
		digit = int(in[i] - '0')
		if isSecondDigit(idx) {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		result += digit
		idx++
	}
	//  check the input string's length after all the empty spaces where skipped
	if idx <= 2 {
		return false
	}
	return result%10 == 0
}

// isSecondDigit returns true on every second digit of the number,
// counting from the right side using another idx, not just in[i].
func isSecondDigit(index int) bool {
	return index%2 == 0
}
