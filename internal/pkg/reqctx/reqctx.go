package reqctx

import (
	"at-backend-announcement/internal/pkg/roles"
	"context"

	"github.com/google/uuid"
)

type keyType int

const ReqCtxKey = keyType(0)

type ReqCtx struct {
	CorrelationID *string
	UserID        *uuid.UUID
	ResponseCode  *int
	Role          *roles.Role
}

func WithCorrelationID(ctx context.Context, correlationID string) context.Context {
	if c, ok := ctx.Value(ReqCtxKey).(ReqCtx); ok {
		c.CorrelationID = &correlationID
		return context.WithValue(ctx, ReqCtxKey, c)
	}

	return context.WithValue(ctx, ReqCtxKey, ReqCtx{CorrelationID: &correlationID})
}

func GetCorrelationID(ctx context.Context) (string, bool) {
	if c, ok := ctx.Value(ReqCtxKey).(ReqCtx); ok && c.CorrelationID != nil {
		return *c.CorrelationID, true
	}

	return "", false
}

func WithUserID(ctx context.Context, userID uuid.UUID) context.Context {
	if c, ok := ctx.Value(ReqCtxKey).(ReqCtx); ok {
		c.UserID = &userID
		return context.WithValue(ctx, ReqCtxKey, c)
	}

	return context.WithValue(ctx, ReqCtxKey, ReqCtx{UserID: &userID})
}

func GetUserID(ctx context.Context) *uuid.UUID {
	if c, ok := ctx.Value(ReqCtxKey).(ReqCtx); ok {
		return c.UserID
	}

	return nil
}

func WithResponseCode(ctx context.Context, responseCode int) context.Context {
	if c, ok := ctx.Value(ReqCtxKey).(ReqCtx); ok {
		c.ResponseCode = &responseCode
		return context.WithValue(ctx, ReqCtxKey, c)
	}

	return context.WithValue(ctx, ReqCtxKey, ReqCtx{ResponseCode: &responseCode})
}

func GetResponseCode(ctx context.Context) *int {
	if c, ok := ctx.Value(ReqCtxKey).(ReqCtx); ok {
		return c.ResponseCode
	}

	return nil
}

func WithRole(ctx context.Context, r roles.Role) context.Context {
	if c, ok := ctx.Value(ReqCtxKey).(ReqCtx); ok {
		c.Role = &r
		return context.WithValue(ctx, ReqCtxKey, c)
	}

	return context.WithValue(ctx, ReqCtxKey, ReqCtx{Role: &r})
}

func GetRole(ctx context.Context) (roles.Role, bool) {
	if c, ok := ctx.Value(ReqCtxKey).(ReqCtx); ok && c.Role != nil {
		return *c.Role, true
	}

	return roles.Unknown, false
}
