package guards

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/insurgence-ai/llm-gateway/internal/admin/models"
	"github.com/insurgence-ai/llm-gateway/internal/config"
	"github.com/insurgence-ai/llm-gateway/internal/lib/exceptions"
)

// wrapper func
func Use(api huma.API, mw func(huma.API) func(huma.Context, func(huma.Context))) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		handler := mw(api) // create new instance on each request
		handler(ctx, next)
	}
}

// TODO: update this to be relevant to this app
func AuthenticationMiddleware(api huma.API) func(ctx huma.Context, next func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		authHeader := ctx.Header("Authorization")
		if authHeader == "" {
			huma.WriteErr(api, ctx, http.StatusUnauthorized, "Missing Authorization header")
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims := &models.ScopedTokenClaims{}

		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(config.Envs.AuthSecret), nil
		})
		if err != nil || !token.Valid {
			huma.WriteErr(api, ctx, http.StatusUnauthorized, "Invalid token")
			return
		}

		userID, err := uuid.Parse(claims.Subject)
		if err != nil {
			huma.WriteErr(api, ctx, http.StatusUnauthorized, "Invalid user ID in subject")
			return
		}

		session := models.UserClaims{
			UserID: userID,
			Email:  claims.Email,
		}

		if claims.Organisation != nil && claims.Organisation.OrganisationID != uuid.Nil {
			session.Organisations = []models.OrganisationMembership{*claims.Organisation}
		} else if len(claims.Organisations) > 0 {
			session.Organisations = claims.Organisations
		}

		ctx = huma.WithValue(ctx, models.UserClaimsKey, session)
		next(ctx)
	}
}

// Accessors for handlers
func GetUserSession(ctx context.Context) (models.UserClaims, bool) {
	claims, ok := ctx.Value(models.UserClaimsKey).(models.UserClaims)
	return claims, ok
}

func GetScopedOrganisationID(ctx context.Context) (uuid.UUID, error) {
	session, ok := GetUserSession(ctx)
	if !ok {
		return uuid.Nil, exceptions.Unauthorized("unauthenticated or session missing")
	}

	if len(session.Organisations) == 0 {
		return uuid.Nil, exceptions.Forbidden("no organisation context in token")
	}

	if len(session.Organisations) > 1 {
		return uuid.Nil, exceptions.Forbidden("token not scoped to a single organisation")
	}

	return session.Organisations[0].OrganisationID, nil
}
