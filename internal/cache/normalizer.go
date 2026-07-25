package cache

import (
	"net/url"
	"strings"
)

func NormalizeURL(urlStr string) (string, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}

	scheme := strings.ToLower(u.Scheme)
	host := strings.ToLower(u.Host)

	// Remove default ports
	if scheme == "http" && strings.HasSuffix(host, ":80") {
		host = host[:len(host)-3]
	} else if scheme == "https" && strings.HasSuffix(host, ":443") {
		host = host[:len(host)-4]
	}

	path := u.Path
	if path == "" {
		path = "/"
	}
	// Strip trailing slash unless it's just "/"
	if len(path) > 1 && strings.HasSuffix(path, "/") {
		path = path[:len(path)-1]
	}

	// Sort query parameters to make keys deterministic
	q := u.Query()
	rawQuery := q.Encode()

	// Reconstruct
	normalized := &url.URL{
		Scheme:   scheme,
		Host:     host,
		Path:     path,
		RawQuery: rawQuery,
	}

	return normalized.String(), nil
}
