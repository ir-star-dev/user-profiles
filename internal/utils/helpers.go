package utils

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

func DetectDevice(ua string) string {
	ua = strings.ToLower(ua)
	switch {
	case strings.Contains(ua, "iphone"):
		return "iPhone"

	case strings.Contains(ua, "android"):
		return "Android"

	case strings.Contains(ua, "ipad"):
		return "iPad"

	case strings.Contains(ua, "windows"):
		return "Windows PC"

	case strings.Contains(ua, "mac"):
		return "Mac"

	default:
		return "Unknown"
	}
}

func WithQuery(base string, params map[string]string, key string, value any) string {
	q := url.Values{}
	for k, v := range params {
		if v != "" {
			q.Set(k, v)
		}
	}
	if value == nil {
		return base + "?" + q.Encode()
	}
	q.Set(key, fmt.Sprint(value))
	return base + "?" + q.Encode()
}

func GetPostIdFromSlug(slug string) (*int, error) {
	parts := strings.Split(slug, "-")
	idStr := parts[len(parts)-1]
	pId, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, err
	}
	return &pId, nil
}
