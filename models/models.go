package models

import (
	"github.com/dgrijalva/jwt-go/v4"
)

// ISClaims is the basic structure of a jwt from identityserver4
type ISClaims struct {
	*jwt.StandardClaims
	// jwt.StandardClaims
	ClientID string `json:"client_id,omitempty"`
	// Subject is also the user Id
	AuthorizedTime       int      `json:"auth_time,omitempty"`
	IdentityProvider     string   `json:"idp,omitempty"`
	SecurityStamp        string   `json:"AspNet.Identity.SecurityStamp,omitempty"`
	Role                 []string `json:"role"`
	PreferredUsername    string   `json:"preferred_username,omitempty"`
	Name                 string   `json:"name,omitempty"`
	Email                string   `json:"email,omitempty"`
	EmailVerified        bool     `json:"email_verified,omitempty"`
	Phone                string   `json:"phone_number,omitempty"`
	PhoneVerified        bool     `json:"phone_number_verified,omitempty"`
	Scope                []string `json:"scope,omitempty"`
	AuthenticationMethod []string `json:"amr,omitempty"`
	SchoolID             string   `json:"school_id,omitempty"`
}

// Error is a basic structure for all errors coming out of the api
type Error struct {
	Message string `json:"message"`
	Field   string `json:"field,omitempty"`
}

// Cursor is used for pagination
type Cursor struct {
	Prev string `json:"prev"`
	Next string `json:"next"`
}

// Metadata is basic data about the data being returned from the api
type Metadata struct {
	Message string  `json:"message"`
	Errors  []Error `json:"errors"`
	Success bool    `json:"success"`
	Cursor  Cursor  `json:"cursor"`
	// fields: string `json:message,omitempty`
}

// NatsResponse is the parent response
type NatsResponse struct {
	Status int `json:"status"`
	*NatsData
}

// NatsData is a child to NatsResponse that holds the data and metadata
type NatsData struct {
	Data     interface{} `json:"data"`
	Metadata *Metadata   `json:"metadata"`
}

// RolesRequest is used in the request to the identity server for authorization
type RolesRequest struct {
	Roles       []string `json:"roles"`
	RequestorID string   `json:"requestorId"`
}
