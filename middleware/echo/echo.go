package echo

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Paxman23l/golang-api-tools/models"
	"github.com/Paxman23l/golang-api-tools/utils"
	"github.com/labstack/echo"
)

// IsInOneRole checks to see if the jwt is is one of the roles listed
func IsInOneRole(roles []string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(e echo.Context) error {
			userRoles := utils.GetRoles(e)
			isInRoles := false
			for _, role := range roles {
				if utils.IsInArray(userRoles, strings.ToLower(role)) == true {
					isInRoles = true
					break
				}
			}

			if isInRoles == false {

				var metadata models.Metadata
				metadata.Message = "User is not in required role"
				// for _, role := range missingRoles {
				// 	// metadata.Errors = append(metadata.Errors, fmt.Sprintf("User must be in %s", role))
				// }
				utils.EchoGenerateResponse(
					http.StatusUnauthorized,
					e,
					nil,
					&metadata,
				)
				// e.Response().
				// 	e.Abort()
				return next(e)
			}
			// e.Next()
			return next(e)
		}
	}
}

// IsInRequiredRoles checks to see if the jwt has the correct roles
func IsInRequiredRoles(roles []string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) {
			userRoles := utils.GetRoles(c)
			isInRoles := true
			var missingRoles []string
			for _, role := range roles {
				if utils.IsInArray(userRoles, strings.ToLower(role)) == false {
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
				utils.GenerateResponse(
					http.StatusUnauthorized,
					c,
					nil,
					&metadata,
				)
				c.Abort()
				return
			}
			c.Next()
		}
	}
}

// // IsInOneRole checks to see if the jwt is is one of the roles listed
// func IsInOneRole(roles []string) gin.HandlerFunc {
// 	return func(c *gin.Context) {

// 		userRoles := utils.GetRoles(c)
// 		isInRoles := false
// 		for _, role := range roles {
// 			if utils.IsInArray(userRoles, strings.ToLower(role)) == true {
// 				isInRoles = true
// 				break
// 			}
// 		}

// 		if isInRoles == false {

// 			var metadata models.Metadata
// 			metadata.Message = "User is not in required role"
// 			// for _, role := range missingRoles {
// 			// 	// metadata.Errors = append(metadata.Errors, fmt.Sprintf("User must be in %s", role))
// 			// }
// 			utils.GenerateResponse(
// 				http.StatusUnauthorized,
// 				c,
// 				nil,
// 				&metadata,
// 			)
// 			c.Abort()
// 			return
// 		}
// 		c.Next()
// 	}
// }
