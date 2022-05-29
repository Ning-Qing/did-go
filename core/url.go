package core

import (
	"errors"
	"fmt"
	"strings"
)

type URL struct {
	Scheme   string
	Method   string
	Path     string
	Query    string
	Fragment string
}

func (u *URL) DID() string {
	return fmt.Sprintf("%s:%s:%s", u.Scheme, u.Method, u.Path)
}

func (u *URL) DIDServer() string {
	return fmt.Sprintf("%s:%s:%s#%s", u.Scheme, u.Method, u.Path, u.Fragment)
}

func (u *URL) setFragment(f string) {
	u.Fragment = f
}

func (u *URL) setPath(p string) {
	u.Path = p
}

func Parse(rawURL string) (*URL, error) {
	u, frag, _ := strings.Cut(rawURL, "#")
	url, err := parse(u, false)
	if err != nil {
		return nil, err
	}
	if frag == "" {
		return url, nil
	}
	url.setFragment(frag)
	return url, nil
}

func parse(rawURL string, viaRequest bool) (*URL, error) {
	var rest string
	var err error

	if stringContainsCTLByte(rawURL) {
		return nil, errors.New("core/url: invalid control character in URL")
	}

	if rawURL == "" && viaRequest {
		return nil, errors.New("empty url")
	}
	url := new(URL)

	if url.Scheme, rest, err = getScheme(rawURL); err != nil {
		return nil, err
	}

	url.Scheme = strings.ToLower(url.Scheme)

	if url.Scheme != "did" {
		return nil, errors.New("unkonw scheme")
	}

	if url.Method, rest, err = getMethod(rest); err != nil {
		return nil, err
	}

	rest, url.Query, _ = strings.Cut(rest, "?")

	url.setPath(rest)
	return url, nil
}

func getMethod(path string) (method, restpath string, err error) {
	for i := 0; i < len(path); i++ {
		c := path[i]
		switch {
		case 'a' <= c && c <= 'z' || '0' <= c && c <= '9':
		case c == ':':
			if i == 0 {
				return "", "", errors.New("missing protocol method")
			}
			return path[:i], path[i+1:], nil
		default:
			return "", path, nil
		}
	}
	return "", path, nil
}

func getScheme(rawURL string) (scheme, path string, err error) {
	for i := 0; i < len(rawURL); i++ {
		c := rawURL[i]
		switch {
		case 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z':
		case '0' <= c && c <= '9' || c == '+' || c == '-' || c == '.':
			if i == 0 {
				return "", rawURL, nil
			}
		case c == ':':
			if i == 0 {
				return "", "", errors.New("missing protocol scheme")
			}
			return rawURL[:i], rawURL[i+1:], nil
		default:
			return "", rawURL, nil
		}
	}
	return "", rawURL, nil
}

type Values map[string][]string

func ParseQuery(query string) (Values, error) {
	m := make(Values)
	err := parseQuery(m, query)
	return m, err
}

func parseQuery(m Values, query string) (err error) {
	for query != "" {
		var key string
		key, query, _ = strings.Cut(query, "&")
		if strings.Contains(key, ";") {
			err = fmt.Errorf("invalid semicolon separator in query")
			continue
		}
		if key == "" {
			continue
		}
		key, value, _ := strings.Cut(key, "=")
		m[key] = append(m[key], value)
	}
	return err
}

func stringContainsCTLByte(s string) bool {
	for i := 0; i < len(s); i++ {
		b := s[i]
		if b < ' ' || b == 0x7f {
			return true
		}
	}
	return false
}
