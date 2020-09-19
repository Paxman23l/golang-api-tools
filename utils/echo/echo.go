package echo

import (
	"fmt"
	"strings"

	"github.com/labstack/echo"

	// "gopkg.in/square/go-jose.v2/jwt"
	"github.com/Paxman23l/golang-api-tools/models"
	"github.com/dgrijalva/jwt-go/v4"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// Functions of type `txnFunc` are passed as arguments to our
// `runTransaction` wrapper that handles transaction retries for us
type txnFunc func(*gorm.DB) error

// IsInRole returns if a user is in specified Role
func IsInRole(c echo.Context, role string) bool {
	roles := GetRoles(c)
	lowerCaseRole := strings.ToLower(role)
	for _, n := range roles {
		if lowerCaseRole == strings.ToLower(n) {
			return true
		}
	}

	return false
}

// GetRoles returns the roles for the user
func GetRoles(c echo.Context) []string {
	claims, err := GetClaims(c)
	if err != nil || claims.Role == nil {
		return []string{}
	}
	return claims.Role
}

// GetClaims returns info out of jwt token
func GetClaims(c echo.Context) (models.ISClaims, error) {
	var claims models.ISClaims
	reqToken := c.Request().Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer")
	if len(splitToken) != 2 {
		// Error: Bearer token not in proper format
		return claims, fmt.Errorf("Bearer header not in correct format or does not exist")
	}

	reqToken = strings.TrimSpace(splitToken[1])

	// Ignore these errors since it's just returning a map
	jwt.ParseWithClaims(reqToken, &claims, nil)

	return claims, nil
}

// GenerateResponse generates a response for the api
func GenerateResponse(status int, c echo.Context, data interface{}, meta *models.Metadata) {

	// Set for successful status codes
	if status >= 200 && status <= 299 {
		meta.Success = true
	}

	c.JSON(
		status,
		gin.H{
			"data":     data,
			"metadata": meta,
		},
	)
}
