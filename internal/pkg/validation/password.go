package validation

import (
	"math"
	"unicode"
)

func IsPwdComplex(password string) bool {
	isDigit := false
	isUpper := false
	isLower := false
	isSpecial := false

	for _, c := range password {
		if unicode.IsDigit(c) {
			isDigit = true
		}
		if unicode.IsUpper(c) {
			isUpper = true
		}
		if unicode.IsLower(c) {
			isLower = true
		}
		if unicode.IsPunct(c) || unicode.IsSymbol(c) {
			isSpecial = true
		}
	}

	var symbolPool int
	if isLower && isUpper && isDigit && isSpecial {
		symbolPool = 95 // contains (a-z, A-Z, ASCII, space)
	} else if isLower && isUpper && isDigit {
		symbolPool = 62 // contains (a-z, A-Z, 0-9)
	} else if isLower && isDigit {
		symbolPool = 36 // contains (a-z, 0-9)
	} else {
		symbolPool = 26 // contains (a-z)
	}

	pwdComplexity := math.Log2(float64(symbolPool)) * float64(len(password))

	const minComplexity = 40.0
	if pwdComplexity < minComplexity {
		return false
	}

	return true
}
