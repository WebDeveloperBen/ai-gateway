package middleware

import (
	"context"

	"github.com/WebDeveloperBen/ai-gateway/internal/logger"
	"github.com/WebDeveloperBen/ai-gateway/internal/model"
)

func GetOrgQueries(ctx context.Context) *OrgQueries {
	return ctx.Value(orgQueriesKey).(*OrgQueries)
}

func GetScopedToken(ctx context.Context) (model.ScopedToken, bool) {
	claims, ok := ctx.Value(ScopedTokenKey).(model.ScopedToken)
	logger.Logger.Info().Msgf("[Get Scoped Token Claims]: %+v", claims)
	return claims, ok
}

func GetOrgIDFromSession(ctx context.Context) (string, bool) {
	claims, ok := ctx.Value(ScopedTokenKey).(model.ScopedToken)
	if !ok || claims.OrgID == "" {
		return "", false
	}
	return claims.OrgID, true
}
