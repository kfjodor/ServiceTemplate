package handlers

import (
	"encoding/json"
	"net/http"
)

// error structure
type ForeignError struct {
	Code    string `json:"code"`
	Payload string `json:"payload"`
}

//
func buildForeignError(w http.ResponseWriter, s int, code, payload string) {
	w.WriteHeader(s)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ForeignError{
		Code:    code,
		Payload: payload,
	})
}

// Currency not supported
// ERROR_CURRENCY_NOT_SUPPORTED
func ERROR_CURRENCY_NOT_SUPPORTED(w http.ResponseWriter, payload string) {
	buildForeignError(w, http.StatusNotImplemented, "ERROR_CURRENCY_NOT_SUPPORTED", payload)
}

// Asset not found
// ERROR_ASSET_NOT_FOUND
func ERROR_ASSET_NOT_FOUND(w http.ResponseWriter, payload string) {
	buildForeignError(w, http.StatusNotFound, "ERROR_ASSET_NOT_FOUND", payload)
}

// Asset state not found
// ERROR_ASSET_STATE_NOT_FOUND
func ERROR_ASSET_STATE_NOT_FOUND(w http.ResponseWriter, payload string) {
	buildForeignError(w, http.StatusNotFound, "ERROR_ASSET_STATE_NOT_FOUND", payload)
}

// Asset state invalid
// ERROR_ASSET_STATE_INVALID
func ERROR_ASSET_STATE_INVALID(w http.ResponseWriter, payload string) {
	buildForeignError(w, http.StatusNotFound, "ERROR_ASSET_STATE_INVALID", payload)
}

// Not found
// ERROR_NOT_FOUND
func ERROR_NOT_FOUND(w http.ResponseWriter, payload string) {
	buildForeignError(w, http.StatusNotFound, "ERROR_NOT_FOUND", payload)
}

//  Address not found
// ERROR_ADDRESS_NOT_FOUND
func ERROR_ADDRESS_NOT_FOUND(w http.ResponseWriter, payload string) {
	buildForeignError(w, http.StatusNotFound, "ERROR_ADDRESS_NOT_FOUND", payload)
}

// Bad request
// ERROR_BAD_REQUEST
func ERROR_BAD_REQUEST(w http.ResponseWriter, payload string) {
	buildForeignError(w, http.StatusBadRequest, "ERROR_BAD_REQUEST", payload)
}

// Currency not set
// ERROR_CURRENCY_NOT_SET
func ERROR_CURRENCY_NOT_SET(w http.ResponseWriter) {
	buildForeignError(w, http.StatusBadRequest, "ERROR_CURRENCY_NOT_SET", "")
}

// Currency not found
// ERROR_CURRENCY_NOT_FOUND
func ERROR_CURRENCY_NOT_FOUND(w http.ResponseWriter) {
	buildForeignError(w, http.StatusNotFound, "ERROR_CURRENCY_NOT_FOUND", "")
}

// Address not set
// ERROR_ADDRESS_NOT_SET
func ERROR_ADDRESS_NOT_SET(w http.ResponseWriter) {
	buildForeignError(w, http.StatusBadRequest, "ERROR_ADDRESS_NOT_SET", "")
}

// Transaction not set
// ERROR_TRX_NOT_SET
func ERROR_TRX_NOT_SET(w http.ResponseWriter) {
	buildForeignError(w, http.StatusBadRequest, "ERROR_TRX_NOT_SET", "")
}

// Transaction not found
// ERROR_TRX_NOT_FOUND
func ERROR_TRX_NOT_FOUND(w http.ResponseWriter) {
	buildForeignError(w, http.StatusNotFound, "ERROR_TRX_NOT_FOUND", "")
}

// Asset id is not set
// ERROR_ASSET_NOT_SET
func ERROR_ASSET_NOT_SET(w http.ResponseWriter) {
	buildForeignError(w, http.StatusBadRequest, "ERROR_ASSET_NOT_SET", "")
}

// Asset invalid
// ERROR_ASSET_INVALID
func ERROR_ASSET_INVALID(w http.ResponseWriter, p string) {
	buildForeignError(w, http.StatusBadRequest, "ERROR_ASSET_INVALID", p)
}

// Authorization missing
// ERROR_AUTH_MISSING
func ERROR_AUTH_MISSING(w http.ResponseWriter) {
	buildForeignError(w, http.StatusForbidden, "ERROR_AUTH_MISSING", "")
}

// Invalid authorization
// ERROR_AUTH_INVALID
func ERROR_AUTH_INVALID(w http.ResponseWriter, payload string) {
	buildForeignError(w, http.StatusForbidden, "ERROR_AUTH_INVALID", payload)
}

// Cannot parse JWT token
// ERROR_AUTH_CANNOT_PARSE_TOKEN
func ERROR_AUTH_CANNOT_PARSE_TOKEN(w http.ResponseWriter) {
	buildForeignError(w, http.StatusForbidden, "ERROR_AUTH_CANNOT_PARSE_TOKEN", "")
}

// Invalid authorization token
// ERROR_AUTH_TOKEN_INVALID
func ERROR_AUTH_TOKEN_INVALID(w http.ResponseWriter) {
	buildForeignError(w, http.StatusForbidden, "ERROR_AUTH_TOKEN_INVALID", "")
}

// User not found
// ERROR_AUTH_USER_NOT_FOUND
func ERROR_AUTH_USER_NOT_FOUND(w http.ResponseWriter, pl string) {
	buildForeignError(w, http.StatusForbidden, "ERROR_AUTH_USER_NOT_FOUND", pl)
}

// No permission
// ERROR_AUTH_NO_PERMISSION
func ERROR_AUTH_NO_PERMISSION(w http.ResponseWriter, pl string) {
	buildForeignError(w, http.StatusForbidden, "ERROR_AUTH_NO_PERMISSION", pl)
}

// Bad password
// ERROR_AUTH_BAD_PASSWORD
func ERROR_AUTH_BAD_PASSWORD(w http.ResponseWriter, pl string) {
	buildForeignError(w, http.StatusForbidden, "ERROR_AUTH_BAD_PASSWORD", pl)
}

// Authorization forbidden
// ERROR_AUTH_FORBIDDEN
func ERROR_AUTH_FORBIDDEN(w http.ResponseWriter, pl string) {
	buildForeignError(w, http.StatusForbidden, "ERROR_AUTH_FORBIDDEN", pl)
}

// Seed is not set
// ERROR_AUTH_SEED_NOT_SET
func ERROR_AUTH_SEED_NOT_SET(w http.ResponseWriter, pl string) {
	buildForeignError(w, http.StatusForbidden, "ERROR_AUTH_SEED_NOT_SET", pl)
}

// Seed not found
// ERROR_AUTH_SEED_NOT_FOUND
func ERROR_AUTH_SEED_NOT_FOUND(w http.ResponseWriter, pl string) {
	buildForeignError(w, http.StatusForbidden, "ERROR_AUTH_SEED_NOT_FOUND", pl)
}

// Internal server error
// ERROR_INTERNAL_SERVER
func ERROR_INTERNAL_SERVER(w http.ResponseWriter, pl string) {
	buildForeignError(w, http.StatusInternalServerError, "ERROR_INTERNAL_SERVER", pl)
}

//
