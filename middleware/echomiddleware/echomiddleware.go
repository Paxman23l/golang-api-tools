package echomiddleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Paxman23l/golang-api-tools/models"
	"github.com/Paxman23l/golang-api-tools/utils/echoutils"
	"github.com/Paxman23l/golang-api-tools/utils/genericutils"

	"github.com/labstack/echo"
)

// IsInOneRole checks to see if the jwt is is one of the roles listed
func IsInOneRole(roles []string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(e echo.Context) error {
			userRoles := echoutils.GetRoles(e)
			isInRoles := false
			for _, role := range roles {
				if genericutils.IsInArray(userRoles, strings.ToLower(role)) == true {
					isInRoles = true
					break
				}
			}

			if isInRoles == false {

				var metadata models.Metadata
				metadata.Message = "User is not in required role"
				echoutils.GenerateResponse(
					http.StatusUnauthorized,
					e,
					nil,
					&metadata,
				)

				return nil
			}
			return next(e)
		}
	}
}

// IsInRequiredRoles checks to see if the jwt has the correct roles
func IsInRequiredRoles(roles []string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userRoles := echoutils.GetRoles(c)
			isInRoles := true
			var missingRoles []string
			for _, role := range roles {
				if genericutils.IsInArray(userRoles, strings.ToLower(role)) == false {
					isInRoles = false
					missingRoles = append(missingRoles, role)
				}
			}

			if isInRoles == false {

				var metadata models.Metadata
				metadata.Message = "User is not in required roles"
				for _, role := range missingRoles {
					metadata.Errors = append(metadata.Errors, fmt.Sprintf("User must be in %s", role))
				}
				echoutils.GenerateResponse(
					http.StatusUnauthorized,
					c,
					nil,
					&metadata,
				)
				return nil
			}
			next(c)
			return nil
		}
	}
}
