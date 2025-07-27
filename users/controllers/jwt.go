package controllers

import (
	"users/config"
	"users/models"
	"users/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type JwtController struct {
	DB  *gorm.DB
	Cfg *config.Config
}

func NewJwtController(db *gorm.DB, cfg *config.Config) *JwtController {
	return &JwtController{
		DB:  db,
		Cfg: cfg,
	}
}

func (jc *JwtController) setCookie(c *gin.Context, accessToken string, refreshToken string) {
	c.SetCookie(
		"accessToken",
		accessToken,
		jc.Cfg.JWT.AccessExpiration*24,
		"/",
		"",
		true,
		true,
	)
	c.SetCookie(
		"refreshToken",
		refreshToken,
		jc.Cfg.JWT.RefreshExpiration*24*7,
		"/",
		"",
		true,
		true,
	)
}

func (jc *JwtController) Login(c *gin.Context) {
	var body models.LoginRequest
	if err := c.ShouldBindBodyWithJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Get the user by username or email
	var user models.User
	if body.Username != "" {
		if err := jc.DB.Where("username = ?", body.Username).First(&user).Error; err != nil {
			c.JSON(404, gin.H{"error": "User not found"})
			return
		}
	} else {
		c.JSON(400, gin.H{"error": "Either username or email must be provided"})
		return
	}

	// Check the password
	if err := utils.CheckPasswordHash(body.Password, user.Password); err != nil {
		c.JSON(401, gin.H{"error": "Invalid password"})
		return
	}
	// Generate JWT accessToken
	accessToken, refreshToken, err := jc.generateTokens(&user)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to generate tokens"})
		return
	}
	jc.setCookie(c, accessToken, refreshToken)

	c.JSON(200, gin.H{"access": accessToken, "refresh": refreshToken})

}

func (jc *JwtController) generateTokens(user *models.User) (string, string, error) {
	accessToken, err := utils.GenerateJWT(user, jc.Cfg.JWT.Secret, jc.Cfg.JWT.AccessExpiration)
	if err != nil {
		return "", "", err
	}
	refreshToken, err := utils.GenerateJWT(user, jc.Cfg.JWT.Secret, jc.Cfg.JWT.RefreshExpiration)
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

func (jc *JwtController) RefreshToken(c *gin.Context) {
	claims, code, err := utils.ValidateJwtHelper(c, jc.Cfg.JWT)
	if err != nil {
		c.JSON(code, gin.H{"error": err.Error()})
		return
	}

	user := &models.User{
		Id:       claims.UserID,
		Username: claims.Username,
	}
	// Generate new access token
	newAccess, newRefresh, err := jc.generateTokens(user)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to generate new tokens"})
		return
	}
	jc.setCookie(c, newAccess, newRefresh)

	c.JSON(200, gin.H{"access": newAccess, "refresh": newRefresh})
}

func (jc *JwtController) ValidateJwt(c *gin.Context) {
	user, code, err := utils.ValidateJwtHelper(c, jc.Cfg.JWT)
	if err != nil {
		c.JSON(code, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, user)
}

func (jc *JwtController) createJwtRouter(server *gin.Engine) *gin.RouterGroup {
	api := server.Group("/api")
	jwt := api.Group("/jwt")
	return jwt
}

func (jc *JwtController) Bind(server *gin.Engine) {
	router := jc.createJwtRouter(server)
	router.POST("/login", jc.Login)
	router.GET("/refresh", jc.RefreshToken)
	router.GET("/validate", jc.ValidateJwt)
}
