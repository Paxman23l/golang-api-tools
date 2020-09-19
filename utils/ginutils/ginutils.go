package ginutils

import (
	"encoding/base64"
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"

	// "gopkg.in/square/go-jose.v2/jwt"
	"github.com/Paxman23l/golang-api-tools/models"
	"github.com/dgrijalva/jwt-go/v4"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)

// Functions of type `txnFunc` are passed as arguments to our
// `runTransaction` wrapper that handles transaction retries for us
type txnFunc func(*gorm.DB) error

// IsInRole returns if a user is in specified Role
func IsInRole(c *gin.Context, role string) bool {
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
func GetRoles(c *gin.Context) []string {
	claims, err := GetClaims(c)
	if err != nil || claims.Role == nil {
		return []string{}
	}
	return claims.Role
}

// GetClaims returns info out of jwt token
func GetClaims(c *gin.Context) (models.ISClaims, error) {
	var claims models.ISClaims
	reqToken := c.Request.Header.Get("Authorization")
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

func base64Decode(src string) (string, error) {
	if l := len(src) % 4; l > 0 {
		src += strings.Repeat("=", 4-l)
	}
	decoded, err := base64.URLEncoding.DecodeString(src)
	if err != nil {
		errMsg := fmt.Errorf("Decoding Error %s", err)
		return "", errMsg
	}
	return string(decoded), nil
}

// RunTransaction is a wrapper for a transaction.  This automatically re-calls `fn` with
// the open transaction as an argument as long as the database server
// asks for the transaction to be retried.
func RunTransaction(db *gorm.DB, fn txnFunc) error {
	var maxRetries = 3
	for retries := 0; retries <= maxRetries; retries++ {
		if retries == maxRetries {
			return fmt.Errorf("hit max of %d retries, aborting", retries)
		}
		txn := db.Begin()
		if err := fn(txn); err != nil {
			// We need to cast GORM's db.Error to *pq.Error so we can
			// detect the Postgres transaction retry error code and
			// handle retries appropriately.
			pqErr := err.(*pq.Error)
			if pqErr.Code == "40001" {
				// Since this is a transaction retry error, we
				// ROLLBACK the transaction and sleep a little before
				// trying again.  Each time through the loop we sleep
				// for a little longer than the last time
				// (A.K.A. exponential backoff).
				txn.Rollback()
				var sleepMs = math.Pow(2, float64(retries)) * 100 * (rand.Float64() + 0.5)
				fmt.Println("Hit 40001 transaction retry error, sleeping", sleepMs, "milliseconds")
				time.Sleep(time.Millisecond * time.Duration(sleepMs))
			} else {
				// If it's not a retry error, it's some other sort of
				// DB interaction error that needs to be handled by
				// the caller.
				return err
			}
		} else {
			// All went well, so we try to commit and break out of the
			// retry loop if possible.
			if err := txn.Commit().Error; err != nil {
				pqErr := err.(*pq.Error)
				if pqErr.Code == "40001" {
					// However, our attempt to COMMIT could also
					// result in a retry error, in which case we
					// continue back through the loop and try again.
					continue
				} else {
					// If it's not a retry error, it's some other sort
					// of DB interaction error that needs to be
					// handled by the caller.
					return err
				}
			}
			break
		}
	}
	return nil
}

// IndexOfString finds the index of an item in a string array.
// This is case insensitive and returns -1 if not found
func IndexOfString(element string, data []string) int {
	element = strings.ToLower(element)
	for k, v := range data {
		if element == strings.ToLower(v) {
			return k
		}
	}
	return -1 //not found.
}

// GenerateResponse generates a response for the api
func GenerateResponse(status int, c *gin.Context, data interface{}, meta *models.Metadata) {

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
