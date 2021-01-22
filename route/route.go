package route

import (
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/arganaphangquestian/gin-jwt/model"
	"github.com/arganaphangquestian/gin-jwt/repository"
	"github.com/dgrijalva/jwt-go"
)

type userRepository struct {
	repo repository.UserRepository
}

func (r *userRepository) register(c *gin.Context) {
	p := new(model.InputUser)
	if err := c.ShouldBind(p); err != nil {
		c.JSON(500, gin.H{
			"success": false,
			"message": err,
		})
	}
	response, err := r.repo.Register(*p)
	if err != nil {
		c.JSON(500, gin.H{
			"success": false,
			"message": err,
		})
	}
	c.JSON(201, gin.H{
		"success": true,
		"message": "Register endpoint reached",
		"data":    response,
	})
}

func (r *userRepository) users(c *gin.Context) {
	response, err := r.repo.Users()
	if err != nil {
		c.JSON(500, gin.H{
			"success": false,
			"message": err,
		})
	}
	c.JSON(200, gin.H{
		"success": true,
		"message": "Get All Users endpoint reached",
		"data":    response,
	})
}

func (r *userRepository) login(c *gin.Context) {
	p := new(model.Login)
	if err := c.ShouldBind(p); err != nil {
		c.JSON(500, gin.H{
			"success": false,
			"message": err,
		})
	}
	user, err := r.repo.Login(*p)
	token, err := createToken(*user)
	if err != nil {
		c.JSON(500, gin.H{
			"success": false,
			"message": err,
		})
	}
	c.JSON(200, gin.H{
		"success": true,
		"message": "Login Successfully",
		"data":    token,
	})
}

func (r *userRepository) dashboard(c *gin.Context) {
	authorizationHeader := c.Request.Header.Get("Authorization")
	if !strings.Contains(authorizationHeader, "Bearer") {
		c.JSON(403, gin.H{
			"success": false,
			"message": "Authorization must be valid",
		})
	}
	user, err := extractToken(authorizationHeader)
	if err != nil {
		c.JSON(500, gin.H{
			"success": false,
			"message": err,
		})
	}
	// After check extractToken is Successfully,
	// just pass `user` variable into repository or another function that you like
	c.JSON(200, gin.H{
		"success": true,
		"message": "DASHBOARD",
		"data":    user,
	})
}

func createToken(user model.User) (string, error) {
	claims := jwt.MapClaims{}
	claims["user"] = user
	claims["exp"] = time.Now().Add(time.Minute * 60 * 24 * 30).Unix() // 1 month
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := at.SignedString([]byte("MY_SUPER_SECRET_KEY"))
	if err != nil {
		return "", err
	}
	return token, nil
}

func extractToken(authorizationHeader string) (interface{}, error) {
	claims := jwt.MapClaims{}
	tokenString := strings.Replace(authorizationHeader, "Bearer ", "", -1)
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if method, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Signing method invalid")
		} else if method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("Signing method invalid")
		}
		return []byte("MY_SUPER_SECRET_KEY"), nil
	})
	if err != nil {
		return nil, fmt.Errorf("Error : %s", err)
	}
	if !token.Valid {
		return nil, fmt.Errorf("Token not valid")
	}
	return claims["user"], nil
}

// New Route
func New(repository repository.UserRepository) *gin.Engine {
	app := gin.Default()
	repo := &userRepository{repository}
	app.GET("/user", repo.users)
	app.POST("/register", repo.register)
	app.POST("/login", repo.login)
	app.GET("/dashboard", repo.dashboard)
	return app
}
