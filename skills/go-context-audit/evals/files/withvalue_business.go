package auditfixture

import "context"

type contextKey string

const (
	userIDKey  contextKey = "userID"
	orderIDKey contextKey = "orderID"
	traceIDKey contextKey = "traceID"
)

func BuildContext(ctx context.Context, userID, orderID, traceID string) context.Context {
	ctx = context.WithValue(ctx, userIDKey, userID)
	ctx = context.WithValue(ctx, orderIDKey, orderID)
	ctx = context.WithValue(ctx, traceIDKey, traceID)
	return ctx
}

func PlaceOrder(ctx context.Context) error {
	userID := ctx.Value(userIDKey).(string)
	orderID := ctx.Value(orderIDKey).(string)
	traceID := ctx.Value(traceIDKey).(string)

	_, _, _ = userID, orderID, traceID
	return nil
}
