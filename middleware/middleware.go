package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"

	"github.com/Paxman23l/golang-api-tools/models"
	"github.com/Paxman23l/golang-api-tools/utils"
)

// AuthenticationRequired checks to see if the jwt has the correct roles
func AuthenticationRequired(h httprouter.Handle, roles []string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		userRoles := utils.GetRoles(r)
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
