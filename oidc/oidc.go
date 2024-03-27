package oidc

import (
	"context"
	"fmt"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/debyten/apierr"
	"github.com/debyten/httplayer"
	"github.com/kelseyhightower/envconfig"
	"log/slog"
	"net/http"
	"strings"
)

// New creates a httplayer.Middleware which verifies an incoming authorization bearer token by using the remote
// jwks uri built with the environment variables injected by the Entando deployer.
// The bearer tokens are issued from the Entando Keycloak instance.
// When an error occurs from the verification phase, 401 unauthorized is returned.
// After the verification, the oidc.IDToken will be parsed and the resulting principal is injected in the request context; to retrieve
// this token, use PrincipalFrom. If the parser return a nil principal the middleware panics.
//
// Below the default environment variables injected in an MS pod by Entando:
//   - KEYCLOAK_AUTH_URL: the keycloak authorization url
//   - KEYCLOAK_REALM: the target keycloak realm
//   - KEYCLOAK_CLIENT_ID: the keycloak client id for this microservice
//
// Use PrincipalParser for the default behaviour.
func New[T any](parser IDTokenParser[T], logger *slog.Logger) httplayer.Middleware {
	l := logger.With("scope", "oidc")
	verifier := newVerifier()
	return func(h http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			auth = strings.TrimPrefix(auth, "Bearer ")
			idToken, err := verifier.Verify(r.Context(), auth)
			if err != nil {
				l.Error("could not verify the token", "err", err)
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			principal, err := parser(idToken)
			if err != nil {
				l.Error("could not parse the token", "err", err)
				apierr.HandleISE(err, w, r)
				return
			}
			if principal == nil {
				l.Error("nil principal")
				panic("the parser has returned a nil principal")
				return
			}
			ctx := context.WithValue(r.Context(), CtxVal, principal)
			h(w, r.WithContext(ctx))
		}
	}
}

func newVerifier() *oidc.IDTokenVerifier {
	var c configuration
	envconfig.MustProcess("KEYCLOAK", &c)
	rks := oidc.NewRemoteKeySet(context.Background(), c.jwks())
	return oidc.NewVerifier(c.issuer(), rks, &oidc.Config{
		ClientID: c.ClientID,
	})
}

type configuration struct {
	AuthURL string `envconfig:"AUTH_URL" default:"http://localhost:9080/auth"`
	Realm   string `default:"entando-go"`
	/* fixme
	   keycloak creates access token with `aud`: [account].
	   but the oidc verifier expects that the token contains also the client id in
	   the audience field.
	*/
	ClientID string `default:"account"`
}

func (c configuration) jwks() string {
	return fmt.Sprintf("%s/realms/%s/protocol/openid-connect/certs", c.AuthURL, c.Realm)
}

func (c configuration) issuer() string {
	return fmt.Sprintf("%s/realms/%s", c.AuthURL, c.Realm)
}
