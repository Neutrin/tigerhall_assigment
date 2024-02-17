package utils

import (
	"fmt"

	"github.com/nitin/tigerhall/core/internal/config"
)

func GenImgURL(escapedPath string) string {
	return fmt.Sprintf("%s%s", config.ImageHost, escapedPath)
}
