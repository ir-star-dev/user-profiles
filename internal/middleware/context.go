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

func GetUserRole(ctx context.Context) (string, error) {
	role, ok := ctx.Value(UserRoleKey).(string)
	if !ok {
		return "", errors.New("You not allowed to see this page")
	}
	return role, nil
}