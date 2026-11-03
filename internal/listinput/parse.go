package listinput

import (
	"errors"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

var listPath = regexp.MustCompile(`^/(?:[1-9][0-9]*/list|lists)/([1-9][0-9]*)(?:/[^/]+)?/?$`)

func Parse(input string) (int64, error) {
	s := strings.TrimSpace(input)
	if id, err := strconv.ParseInt(s, 10, 64); err == nil && id > 0 {
		return id, nil
	}
	u, err := url.Parse(s)
	if err != nil || u.Scheme != "https" || u.Host != "simkl.com" || u.User != nil || u.Fragment != "" {
		return 0, errors.New("Bitte eine Listen-ID oder einen https://simkl.com-Listenlink eingeben")
	}
	match := listPath.FindStringSubmatch(u.Path)
	if len(match) != 2 {
		return 0, errors.New("Simkl-Listenlink nicht erkannt")
	}
	id, err := strconv.ParseInt(match[1], 10, 64)
	if err != nil || id <= 0 {
		return 0, errors.New("Ungültige Listen-ID")
	}
	return id, nil
}
