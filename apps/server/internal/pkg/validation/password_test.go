package validation

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsPwdComplex(t *testing.T) {
	tests := []struct {
		name     string
		password string
		expect   bool
	}{
		{
			name:     "Empty password",
			password: "",
			expect:   false,
		},
		{
			name:     "Lowercase only, too short",
			password: "abcde",
			expect:   false,
		},
		{
			name:     "Lowercase only, sufficient length",
			password: "abcdefghijabcdefghij",
			expect:   true,
		},
		{
			name:     "Lowercase and digits, too short",
			password: "abc123",
			expect:   false,
		},
		{
			name:     "Lowercase and digits, sufficient length",
			password: "abc123abc123",
			expect:   true,
		},
		{
			name:     "Lowercase, uppercase, and digits, too short",
			password: "Abc123",
			expect:   false,
		},
		{
			name:     "Lowercase, uppercase, and digits, sufficient length",
			password: "Abc123Abc123",
			expect:   true,
		},
		{
			name:     "Lowercase, uppercase, digits, and special chars, sufficient length",
			password: "Abc123!@#Abc123!@#",
			expect:   true,
		},
		{
			name:     "Only special characters, too short",
			password: "!@#$%^",
			expect:   false,
		},
		{
			name:     "Only special characters, sufficient length",
			password: "!@#$%^&*()_+!@#$%^&*()_+",
			expect:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := IsPwdComplex(tt.password)
			assert.Equalf(t, tt.expect, actual, "IsPwdComplex() = %v, expect %v", actual, tt.expect)
		})
	}
}
