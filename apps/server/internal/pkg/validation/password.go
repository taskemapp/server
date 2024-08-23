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

	symbolPool := calcSymbolPool(isDigit, isUpper, isLower, isSpecial)

	pwdComplexity := math.Log2(float64(symbolPool)) * float64(len(password))

	const minComplexity = 40.0

	return pwdComplexity > minComplexity
}

func calcSymbolPool(
	isDigit,
	isUpper,
	isLower,
	isSpecial bool,
) int {
	switch {
	case isLower && isUpper && isDigit && isSpecial:
		return 95 // contains (a-z, A-Z, 0-9, special)
	case isLower && isUpper && isDigit:
		return 62 // contains (a-z, A-Z, 0-9)
	case (isLower || isUpper) && isDigit:
		return 36 // contains (a-z or A-Z, 0-9)
	case isSpecial:
		return 32 // contains (special)
	case isLower || isUpper:
		return 26 // contains (a-z or A-Z)
	case isDigit:
		return 10 // contains (0-9)
	}

	return 0
}
