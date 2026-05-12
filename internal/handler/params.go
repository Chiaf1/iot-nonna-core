package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// Parse time value from url
func parseTimeParam(r *http.Request, key string, defaultVal time.Time) (time.Time, error) {
	val := r.URL.Query().Get(key)
	if val == "" {
		return defaultVal, nil
	}
	t, err := time.Parse(time.RFC3339, val)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid format for '%s': use RFC3339 (e.g. 2006-01-02T15:04:05Z)", key)
	}
	return t, nil
}

// Parse int value from url
func parseIntParam(r *http.Request, key string, defaultVal int) (int, error) {
	val := r.URL.Query().Get(key)
	if val == "" {
		return defaultVal, nil
	}
	n, err := strconv.Atoi(val)
	if err != nil {
		return 0, fmt.Errorf("invalid value format for '%s': must be integer", key)
	}
	return n, nil
}
