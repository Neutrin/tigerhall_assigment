package utils

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestGenerateHashPassword(t *testing.T) {
	// Test cases
	testCases := []struct {
		name          string
		password      string
		expectedError bool
	}{
		{"Simple password", "password123", false},
		{"Empty password", "", true},
		{"Complex password", "P@ssw0rd!123", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			hashedPassword, err := GenerateHashPassword(tc.password)
			if (err != nil) != tc.expectedError {
				t.Errorf("Expected error: %v, but got: %v", tc.expectedError, err)
			}
			if !tc.expectedError {
				err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(tc.password))
				if err != nil {
					t.Errorf("Generated hash is not valid: %v", err)
				}
			}
		})
	}
}
