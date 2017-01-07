// Package token enables setting and reading
// authentication token contexts
package token

import (
	"net/http"
	"strings"

	"golang.org/x/net/context"
)

type key int

const tokenKey key = 0

// Token is the concrete type which represents kolide session tokens
type Token string

// String return a string representation of a token
func (t Token) String() string {
	return string(t)
}

// FromHTTPRequest extracts an Authorization
// from an HTTP header if present.
func FromHTTPRequest(r *http.Request) Token {
	headers := r.Header.Get("Authorization")
	headerParts := strings.Split(headers, " ")
	if len(headerParts) != 2 || strings.ToUpper(headerParts[0]) != "BEARER" {
		return ""
	}
	return Token(headerParts[1])
}

// NewContext returns a new context carrying the Authorization Bearer token.
func NewContext(ctx context.Context, token Token) context.Context {
	if token == "" {
		return ctx
	}
	return context.WithValue(ctx, tokenKey, token)
}

// FromContext extracts the Authorization Bearer token if present.
func FromContext(ctx context.Context) (Token, bool) {
	token, ok := ctx.Value(tokenKey).(Token)
	return token, ok
}
