package middleware

import (
	"context"
	"errors"
)

func GetUserID(ctx context.Context) (int, error) {
	id, ok := ctx.Value(UserIdKey).(int)
	if !ok {
		return 0, errors.New("You not allowed to see this page")
	}
	return id, nil
}