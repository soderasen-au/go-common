package util

import (
	"net/url"
	"path"
	"strings"
)

func MaybeAssignStr(s **string, defStr string) {
	if *s == nil || **s == "" {
		*s = &defStr
	}
}

func StrMaybeNil(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func StrMaybeDefault(s *string, defStr string) string {
	if s == nil {
		return defStr
	}
	return *s
}

func IntMaybeDefault(s *int, def int) int {
	if s == nil {
		return def
	}
	return *s
}

func BoolMaybeDeNil(s *bool) bool {
	if s == nil {
		return false
	}
	return *s
}

func MaybeNil[T any](s *T) (res T) {
	if s == nil {
		return res
	}
	return *s
}

func MaybeDefault[T any](s *T, v T) T {
	if s == nil {
		return v
	}
	return *s
}

// Max returns the larger of x or y.
func Max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

// Min returns the smaller of x or y.
func Min(x, y int) int {
	if x > y {
		return y
	}
	return x
}

func GetAPIUrl(baseUrl, endpoint string) (string, *Result) {
	u, err := url.Parse(baseUrl)
	if err != nil {
		return "", Error("ParseURL", err)
	}
	if strings.HasPrefix(endpoint, "/api/v1") {
		endpoint = strings.Replace(endpoint, "/api/v1", "", 1)
	}
	u.Path = path.Join("/api/v1", endpoint)
	return u.String(), nil
}

func GetUrl(baseUrl, endpoint string) (string, *Result) {
	u, err := url.Parse(baseUrl)
	if err != nil {
		return "", Error("ParseURL", err)
	}
	u.Path = path.Join(u.Path, endpoint)
	return u.String(), nil
}

func NewString(v string) *string {
	return &v
}

func Ptr[T any](v T) *T {
	return &v
}

func Int32to64(i *int32) *int64 {
	var o int64 = int64(*i)
	return &o
}
