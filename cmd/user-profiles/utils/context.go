package utils

import (
	"context"
	"errors"
)

type contextKey string

const (
	UserIdKey   contextKey = "user_id"
	UserRoleKey contextKey = "user_role"
	UserBanKey  contextKey = "user_ban"
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

func GetUserBanStatus(ctx context.Context) (*bool, error) {
	ban, ok := ctx.Value(UserBanKey).(bool)
	if !ok {
		return nil, errors.New("You not allowed to see this page")
	}
	return &ban, nil
}
