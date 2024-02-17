package utils

import (
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/nitin/tigerhall/core/internal/config"
)

// Mocking the model and config for testing
type Claim struct {
	jwt.StandardClaims
}

func TestGenerateSignedTokens(t *testing.T) {
	// Test cases
	testCases := []struct {
		name            string
		uniqueKey       string
		expireTimeMin   int
		expectedError   bool
		expectedToken   string
		expirationDelta time.Duration
	}{
		{
			"Valid case",
			"unique_key_123",
			30,
			false,
			"",
			1 * time.Minute,
		},
		{
			"Expired token",
			"expired_key",
			-1,
			true,
			"",
			1 * time.Minute,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			tokenString, err := GenerateSignedTokens(tc.uniqueKey, tc.expireTimeMin)
			// Check if error matches expectation
			if (err != nil) != tc.expectedError {
				t.Errorf("Expected error: %v, but got: %v", tc.expectedError, err)
			}
			// If no error, validate the token
			if !tc.expectedError {
				token, parseErr := jwt.ParseWithClaims(tokenString, &Claim{}, func(token *jwt.Token) (interface{}, error) {
					return []byte(config.JwtKey), nil
				})
				if parseErr != nil || !token.Valid {
					t.Errorf("Generated token is not valid: %v", parseErr)
				}

				claims, ok := token.Claims.(*Claim)
				if !ok {
					t.Error("Failed to parse claims")
				}

				expirationTime := time.Unix(claims.ExpiresAt, 0)
				expectedExpirationTime := time.Now().Add(time.Duration(tc.expireTimeMin) * time.Minute)
				if expirationTime.Before(expectedExpirationTime.Add(-tc.expirationDelta)) ||
					expirationTime.After(expectedExpirationTime.Add(tc.expirationDelta)) {
					t.Errorf("Expiration time not within expected range. Expected: %v, Actual: %v", expectedExpirationTime, expirationTime)
				}
			}
		})
	}
}
