package echomiddleware

import (
	"encoding/json"
	"net/http"

	"github.com/Paxman23l/golang-api-tools/models"
	"github.com/Paxman23l/golang-api-tools/utils/echoutils"
	"github.com/nats-io/nats.go"

	"github.com/labstack/echo"
)

var _nc *nats.Conn

// InitEchoMiddleware instantiates the nats connection for checking authorization
func InitEchoMiddleware(nc *nats.Conn) {
	_nc = nc
}

// IsInOneRole checks to see if the jwt is is one of the roles listed
func IsInOneRole(roles []string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(e echo.Context) error {
			res := models.NatsResponse{}
			reqModel := models.RolesRequest{}
			claims, err := echoutils.GetClaims(e)
			reqModel.RequestorID = claims.Subject
			reqModel.Roles = roles
			byteReqModel, err := json.Marshal(reqModel)
			req, err := _nc.RequestWithContext(e.Request().Context(), "identity.authorization.isinrole", byteReqModel)
			if err != nil || json.Unmarshal(req.Data, &res) != nil || !res.Metadata.Success {
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
			res := models.NatsResponse{}
			reqModel := models.RolesRequest{}
			claims, err := echoutils.GetClaims(c)
			reqModel.RequestorID = claims.Subject
			reqModel.Roles = roles
			byteReqModel, err := json.Marshal(reqModel)
			req, err := _nc.RequestWithContext(c.Request().Context(), "identity.authorization.isinmultipleroles", byteReqModel)
			if err != nil || json.Unmarshal(req.Data, &res) != nil || !res.Metadata.Success {
				var metadata models.Metadata
				metadata.Message = "User is not in required roles"
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
