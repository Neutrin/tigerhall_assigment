package utils

import (
	"testing"

	"github.com/nitin/tigerhall/core/internal/config"
)

func TestGenImgURL(t *testing.T) {
	testCases := []struct {
		name        string
		escapedPath string
		expected    string
	}{
		{name: "Simple case",
			escapedPath: "/image.jpg",
			expected:    config.ImageHost + "/image.jpg",
		},
		{name: "Empty path",
			escapedPath: "",
			expected:    config.ImageHost,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := GenImgURL(tc.escapedPath)
			if result != tc.expected {
				t.Errorf("GenImgURL(%s) returned %s, expected %s", tc.escapedPath, result, tc.expected)
			}
		})
	}
}
