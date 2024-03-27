# entando-go
This library provides a set of utilities to build an Entando MS in go.

## Details
> From Entando MS specifications

To build a minimal microservice to be deployed on the Entando platform it is mandatory to:
- Expose the health endpoint (usually `/api/health`)
- Handle the special variable `SERVER_SERVLET_CONTEXT_PATH`
- Expose the http protocol to the port `8081`

The first two requirements are handled at code level:
- use `entando.NewMux` to build the `http.Handler` which includes two http endpoints (metrics and health)
- build the http middleware with `entando.ContextPath` to strip the context path from each request.
This is particularly useful when running the MS in production.

### Features
To build a production grade MS, you can use these additional packages:
- `datasource`: useful to normalize the `SPRING_DATASOURCE_URL` env injected in production.
- `oidc`: useful to parse the authorization headers (keycloak tokens)
- `tracing`: useful to implement tracing capabilities for the MS.

#### datasource

```go
import "github.com/almartino/entando-go/datasource"

dbConfig := datasource.NewConfiguration()
// use dbConfig properties, e.g. dbConfig.Host, dbConfig.Port
```

#### tracing
```go
import "github.com/almartino/entando-go/tracing"

mw := tracing.NewOtel()
defer mw.Close()
handler := entando.NewMux()
// use mw with http.Handler
handler = mw(handler)
```

#### oidc

```go
import "github.com/almartino/entando-go/oidc"
import "github.com/debyten/httplayer"

logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
oidcMW := oidc.New(oidc.PrincipalParser, logger)
// or build your custom principal parser:

type MyPrincipal struct {
	ID string
	Role string
	Email string
}

func MyPrincipalParser(tok *oidc.IDToken) (*MyPrincipal, error) {
  // implement me
}

oidcMW := oidc.New(MyPrincipalParser, logger)
handler := entando.NewMux()
myService := ...
httplayer.NewServiceBuilder(oidcMW).
    Add(myService). // oidcMW is executed for each myService route
    MountTo(mux)
```

# Examples

Complete Entando bundles are available under examples directory