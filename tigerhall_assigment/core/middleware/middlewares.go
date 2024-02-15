package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/nitin/tigerhall/core/utils"
)

func IsAuthorized() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("token")

		if err != nil {
			c.JSON(401, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		err = utils.ParseToken(cookie)
		if err != nil {
			c.JSON(401, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}
		//Proceed to next controller chained
		c.Next()
	}
}
