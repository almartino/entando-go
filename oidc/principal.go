package oidc

import (
	"context"
	"github.com/coreos/go-oidc/v3/oidc"
)

// Key used for context propagation
type Key int

// CtxVal is the value used to inject the principal.
const CtxVal Key = iota

// PrincipalFrom retrieves the principal from the context if any.
func PrincipalFrom[T any](ctx context.Context) (*T, bool) {
	tok, ok := ctx.Value(CtxVal).(*T)
	return tok, ok
}

// Principal is the default implementation for IDTokenParser
// to be used when the application doesn't have high security requirements.
type Principal struct {
	ID string
}

// PrincipalParser implement IDTokenParser
func PrincipalParser(tok *oidc.IDToken) (*Principal, error) {
	return &Principal{
		ID: tok.Subject,
	}, nil
}

// IDTokenParser must convert the id token to the MS principal.
type IDTokenParser[T any] func(tok *oidc.IDToken) (*T, error)
