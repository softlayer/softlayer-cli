package utils

import (
	"strings"

	"github.com/softlayer/softlayer-go/filter"
)

func QueryFilter(query string, path string) filter.Filter {
	switch {
	case strings.HasPrefix(query, "^="):
		queryString := strings.TrimLeft(query, "^=")
		return filter.Path(path).StartsWith(queryString)
	case strings.HasPrefix(query, "$="):
		queryString := strings.TrimLeft(query, "$=")
		return filter.Path(path).EndsWith(queryString)
	case strings.HasPrefix(query, "*="):
		queryString := strings.TrimLeft(query, "*=")
		return filter.Path(path).Contains(queryString)
	case strings.HasPrefix(query, "!*="):
		queryString := strings.TrimLeft(query, "!*=")
		return filter.Path(path).NotContains(queryString)

	case strings.HasPrefix(query, ">="):
		queryString := strings.TrimLeft(query, ">=")
		return filter.Path(path).GreaterThanOrEqual(queryString)
	case strings.HasPrefix(query, ">"):
		queryString := strings.TrimLeft(query, ">")
		return filter.Path(path).GreaterThan(queryString)
	case strings.HasPrefix(query, "<"):
		queryString := strings.TrimLeft(query, "<")
		return filter.Path(path).LessThan(queryString)
	case strings.HasPrefix(query, "<="):
		queryString := strings.TrimLeft(query, "<=")
		return filter.Path(path).LessThanOrEqual(queryString)
	case strings.HasPrefix(query, "~"):
		queryString := strings.TrimLeft(query, "~")
		return filter.Path(path).Like(queryString)
	case strings.HasPrefix(query, "!~"):
		queryString := strings.TrimLeft(query, "!~")
		return filter.Path(path).NotLike(queryString)

	case strings.HasPrefix(query, "*") && strings.HasSuffix(query, "*"):
		queryString := strings.Trim(query, "*")
		return filter.Path(path).Contains(queryString)
	case strings.HasSuffix(query, "*"):
		queryString := strings.TrimRight(query, "*")
		return filter.Path(path).StartsWith(queryString)
	case strings.HasPrefix(query, "*"):
		queryString := strings.TrimLeft(query, "*")
		return filter.Path(path).EndsWith(queryString)

	default:
		return filter.Path(path).Eq(query)
	}
}
