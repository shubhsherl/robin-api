package router

import (
	"errors"
	"net/http"
	"strings"

	"github.com/RobinHoodArmyHQ/robin-api/internal/util"
	"github.com/RobinHoodArmyHQ/robin-api/pkg/nanoid"
	"github.com/gin-gonic/gin"
)

func extractBearerToken(header string) (string, error) {
	if header == "" {
		return "", errors.New("missing authorization header")
	}

	jwtToken := strings.Split(header, " ")
	if len(jwtToken) != 2 {
		return "", errors.New("incorrectly formatted authorization header")
	}

	return jwtToken[1], nil
}

func isUserLoggedIn(c *gin.Context) {
	token, err := extractBearerToken(c.GetHeader("Authorization"))

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	claims, err := util.VerifyJwt(token)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		c.Abort()
		return
	}

	c.Set("user_id", nanoid.NanoID(claims.UserId))
	c.Set("user_roles", claims.UserRoles)
	c.Next()
}

func isAdminUser(c *gin.Context) {
	c.Next()
}
