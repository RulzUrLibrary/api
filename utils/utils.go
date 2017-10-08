package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/url"
	"strings"
)

func Dump(i interface{}) string {
	bytes, err := json.MarshalIndent(i, "", "  ")
	if err != nil {
		return fmt.Sprintf("error marshaling entity: %s", err)
	}
	return string(bytes)
}

func Ceil(i, j int) int {
	if j == 0 {
		return 0
	}
	return int(math.Ceil(float64(i) / float64(j)))
}

func UrlQueryS(path string, values ...interface{}) *url.URL {
	u, _ := url.Parse(path)
	q := url.Values{}
	for i := 0; i < len(values)-1; i += 2 {
		q.Add(fmt.Sprint(values[i]), fmt.Sprint(values[i+1]))
	}
	u.RawQuery = q.Encode()
	return u
}

func UrlQuery(path *url.URL, values ...interface{}) *url.URL {
	return UrlQueryS(path.String(), values...)
}

func SanitizeIsbn(isbn string) string {
	return strings.Map(func(r rune) rune {
		if (r >= '0' && r <= '9') || r == 'X' {
			return r
		}
		return -1
	}, isbn)
}

func DecodeJSON(rc io.Reader, i interface{}) error {
	return json.NewDecoder(rc).Decode(i)
}

func DefaultS(s ...string) string {
	for _, s := range s {
		if s != "" {
			return s
		}
	}
	return ""
}
