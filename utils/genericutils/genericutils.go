package genericutils

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"

	"github.com/nats-io/nats.go"

	// "gopkg.in/square/go-jose.v2/jwt"
	"github.com/Paxman23l/golang-api-tools/models"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)

// Functions of type `txnFunc` are passed as arguments to our
// `runTransaction` wrapper that handles transaction retries for us
type txnFunc func(*gorm.DB) error

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

// GinGenerateResponse generates a response for the api
func GinGenerateResponse(status int, c *gin.Context, data interface{}, meta *models.Metadata) {

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

// GenerateNatsResponse builds responses for nats requests
func GenerateNatsResponse(status int, msg *nats.Msg, data interface{}, meta *models.Metadata) {

	// Set for successful status codes
	if status >= 200 && status <= 299 {
		meta.Success = true
	}

	response := &models.NatsResponse{}
	response.NatsData = &models.NatsData{}
	response.Status = status
	response.Data = data
	response.Metadata = meta

	jsonData, err := json.Marshal(response)
	if err != nil {
		msg.Respond([]byte(err.Error()))
		return
	}
	msg.Respond(jsonData)
}

// ArrayFind takes a slice and looks for an element in it. If found it will
// return it's key, otherwise it will return -1 and a bool of false.
func ArrayFind(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

// IsInArray takes a slice and looks for an element in it. If found it will
// return it's key, otherwise it will return -1 and a bool of false.
func IsInArray(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}
