package ginmiddleware

import (
	"encoding/json"
	"net/http"

	"github.com/Paxman23l/golang-api-tools/utils/ginutils"
	"github.com/nats-io/nats.go"

	"github.com/Paxman23l/golang-api-tools/models"
	"github.com/gin-gonic/gin"
)

var _nc *nats.Conn

// InitGinMiddleware instantiates the nats connection for checking authorization
func InitGinMiddleware(nc *nats.Conn) {
	_nc = nc
}

// IsInRequiredRoles checks to see if the jwt has the correct roles
func IsInRequiredRoles(roles []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		res := models.NatsResponse{}
		reqModel := models.RolesRequest{}
		claims, err := ginutils.GetClaims(c)
		reqModel.RequestorID = claims.Subject
		reqModel.Roles = roles
		byteReqModel, err := json.Marshal(reqModel)
		req, err := _nc.RequestWithContext(c, "identity.authorization.isinmultipleroles", byteReqModel)
		if err != nil || json.Unmarshal(req.Data, &res) != nil || !res.Metadata.Success {
			var metadata models.Metadata
			metadata.Message = "User is not in required role"
			ginutils.GenerateResponse(
				http.StatusUnauthorized,
				c,
				nil,
				&metadata,
			)
			c.Abort()
		}
		c.Next()
	}
}

// IsInOneRole checks to see if the jwt is is one of the roles listed
func IsInOneRole(roles []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		res := models.NatsResponse{}
		reqModel := models.RolesRequest{}
		claims, err := ginutils.GetClaims(c)
		reqModel.RequestorID = claims.Subject
		reqModel.Roles = roles
		byteReqModel, err := json.Marshal(reqModel)
		req, err := _nc.RequestWithContext(c, "identity.authorization.isinrole", byteReqModel)
		if err != nil || json.Unmarshal(req.Data, &res) != nil || !res.Metadata.Success {
			var metadata models.Metadata
			metadata.Message = "User is not in required role"
			ginutils.GenerateResponse(
				http.StatusUnauthorized,
				c,
				nil,
				&metadata,
			)
			c.Abort()
		}
		c.Next()
	}
}
