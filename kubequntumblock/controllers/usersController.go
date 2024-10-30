package controllers

import (
	"kubequntumblock/internal/initializer"
	"kubequntumblock/models"
	"kubequntumblock/schemas"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c *gin.Context) {
	var body schemas.SignupSchemas
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read pass",
		})

		return
	}

	user := models.User{Email: body.Email, Password: string(hash)}
	result := initializer.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create the user",
		})

		return
	}

	c.JSON(http.StatusAccepted, gin.H{"message": "success"})

}

func Login(c *gin.Context) {
	var body schemas.SignupSchemas
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": "Failed to read body",
		})
		return
	}
	var user models.User
	initializer.DB.First(&user, "email = ?", body.Email)

	if user.ID == 0 {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": "Failed to read email",
		})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": "Failed to read email",
		})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,jwt.MapClaims{
		"sub" : user.ID,
		"exp" : time.Now().Add(time.Hour *24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECERET")))

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error" : "Failed to create the token",
		})
		return 
	}
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString , 3600 * 24, "","http://localhost:3000",true,true)
	c.JSON(http.StatusAccepted,gin.H{})
}

func Validate(c *gin.Context) {
	c.JSON(http.StatusAccepted,gin.H{
		"message":" validated",
	})
}
