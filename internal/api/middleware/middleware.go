package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/lestrrat-go/jwx/v3/jwa"
	"github.com/lestrrat-go/jwx/v3/jwt"
)

func NewAuthMiddleware(api huma.API, secret []byte) func(ctx huma.Context, next func(huma.Context)) {

	return func(ctx huma.Context, next func(huma.Context)) {
		isAuthorizationRequired := false
		for _, opScheme := range ctx.Operation().Security {
			var anyOfNeededScopes []string
			var ok bool
			if anyOfNeededScopes, ok = opScheme["myAuth"]; ok {
				isAuthorizationRequired = true
				break
			}
			fmt.Printf("Security scheme: %v, anyOfNeededScopes: %v\n", opScheme, anyOfNeededScopes)
		}

		if !isAuthorizationRequired {
			next(ctx)
			return
		}

		token := strings.TrimPrefix(ctx.Header("Authorization"), "Bearer ")
		if len(token) == 0 {
			huma.WriteErr(api, ctx, http.StatusUnauthorized, "Unauthorized 1")
			return
		}

		jwtSignedParsedOptions := jwt.WithKey(jwa.HS256(), secret)

		// Parse and validate the JWT.
		parsed, err := jwt.ParseString(token,
			jwt.WithValidate(true),
			jwtSignedParsedOptions,
		)
		if err != nil {
			huma.WriteErr(api, ctx, http.StatusUnauthorized, "Unauthorized 2")
			return
		}
		fmt.Printf("Parsed JWT: %v\n", parsed)

		next(ctx)
	}
}
