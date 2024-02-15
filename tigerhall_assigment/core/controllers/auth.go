package controllers

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nitin/tigerhall/core/internal/config"
	"github.com/nitin/tigerhall/core/internal/model"
	repositiories "github.com/nitin/tigerhall/core/internal/repositiories"
	"github.com/nitin/tigerhall/core/utils"
)

type AuthController struct {
	repo repositiories.UserRepo
}

func NewAuthController(repo repositiories.UserRepo) AuthController {
	return AuthController{repo: repo}
}

func (controller AuthController) Signup(c *gin.Context) {
	var (
		user model.User
		err  error
	)

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if controller.repo.UserExists(user.Email) {
		c.JSON(400, gin.H{"error": "user already exists"})
		return
	}

	if _, err = controller.repo.Create(user); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"success": fmt.Sprintf(" user email = %s created !!!", user.Email)})
}

func (controller AuthController) Login(c *gin.Context) {
	var (
		user model.User
	)

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	existingUser := controller.repo.User(user.Email)
	if existingUser.ID == 0 {
		c.JSON(400, gin.H{"error": "user does not exist"})
		return
	}
	if !utils.CompareHashPassword(user.Password, existingUser.Password) {
		c.JSON(400, gin.H{"error": "invalid password"})
		return
	}
	jwtToken, err := utils.GenerateSignedTokens(user.Email, config.ExpiryTimeInMinutes)
	if err != nil {
		c.JSON(500, gin.H{"error": "could not generate token"})
		return
	}
	c.SetCookie("token", jwtToken, int(time.Now().Add(config.ExpiryTimeInMinutes*time.Minute).Unix()), "/", "localhost", false, true)
	c.JSON(200, gin.H{"success": "user logged in"})
}
