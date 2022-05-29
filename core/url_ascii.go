package core

// import (
// 	"errors"
// 	"fmt"
// 	"strconv"
// 	"strings"
// )

// type Error struct {
// 	Op  string
// 	URL string
// 	Err error
// }

// func (e *Error) Unwrap() error { return e.Err }
// func (e *Error) Error() string { return fmt.Sprintf("%s %q: %s", e.Op, e.URL, e.Err) }

// func (e *Error) Timeout() bool {
// 	t, ok := e.Err.(interface {
// 		Timeout() bool
// 	})
// 	return ok && t.Timeout()
// }

// func (e *Error) Temporary() bool {
// 	t, ok := e.Err.(interface {
// 		Temporary() bool
// 	})
// 	return ok && t.Temporary()
// }

// type EscapeError string

// func (e EscapeError) Error() string {
// 	return "invalid URL escape " + strconv.Quote(string(e))
// }

// type InvalidHostError string

// func (e InvalidHostError) Error() string {
// 	return "invalid character " + strconv.Quote(string(e)) + " in host name"
// }

// const upperhex = "0123456789ABCDEF"

// func ishex(c byte) bool {
// 	switch {
// 	case '0' <= c && c <= '9':
// 		return true
// 	case 'a' <= c && c <= 'f':
// 		return true
// 	case 'A' <= c && c <= 'F':
// 		return true
// 	}
// 	return false
// }

// func unhex(c byte) byte {
// 	switch {
// 	case '0' <= c && c <= '9':
// 		return c - '0'
// 	case 'a' <= c && c <= 'f':
// 		return c - 'a' + 10
// 	case 'A' <= c && c <= 'F':
// 		return c - 'A' + 10
// 	}
// 	return 0
// }

// type encoding int

// const (
// 	encodePath encoding = 1 + iota
// 	encodeQueryComponent
// 	encodeFragment
// )

// func shouldEscape(c byte, mode encoding) bool {
// 	if 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || '0' <= c && c <= '9' {
// 		return false
// 	}

// 	switch c {
// 	case '-', '_', '.':
// 		return false

// 	case '$', '&', '+', ',', '/', ':', ';', '=', '?', '@', '~', '*',' ':
// 		switch mode {
// 		case encodePath:
// 			return c == '?'

// 		case encodeQueryComponent:
// 			return true

// 		case encodeFragment:
// 			return true
// 		}
// 	}
// 	return true
// }

// func QueryUnescape(s string) (string, error) {
// 	return unescape(s, encodeQueryComponent)
// }

// func unescape(s string, mode encoding) (string, error) {
// 	n := 0
// 	for i := 0; i < len(s); {
// 		switch s[i] {
// 		case '%':
// 			n++
// 			if i+2 >= len(s) || !ishex(s[i+1]) || !ishex(s[i+2]) {
// 				s = s[i:]
// 				if len(s) > 3 {
// 					s = s[:3]
// 				}
// 				return "", EscapeError(s)
// 			}
// 			i += 3
// 		default:
// 			if s[i] < 0x80 && shouldEscape(s[i], mode) {
// 				return "", InvalidHostError(s[i : i+1])
// 			}
// 			i++
// 		}
// 	}

// 	var t strings.Builder
// 	t.Grow(len(s) - 2*n)
// 	for i := 0; i < len(s); i++ {
// 		switch s[i] {
// 		case '%':
// 			t.WriteByte(unhex(s[i+1])<<4 | unhex(s[i+2]))
// 			i += 2
// 		default:
// 			t.WriteByte(s[i])
// 		}
// 	}
// 	return t.String(), nil
// }

// func escape(s string, mode encoding) string {
// 	hexCount := 0
// 	for i := 0; i < len(s); i++ {
// 		c := s[i]
// 		if shouldEscape(c, mode) {
// 			hexCount++
// 		}
// 	}

// 	if hexCount == 0 {
// 		return s
// 	}

// 	var buf [64]byte
// 	var t []byte

// 	required := len(s) + 2*hexCount
// 	if required <= len(buf) {
// 		t = buf[:required]
// 	} else {
// 		t = make([]byte, required)
// 	}

// 	j := 0
// 	for i := 0; i < len(s); i++ {
// 		switch c := s[i]; {
// 		case shouldEscape(c, mode):
// 			t[j] = '%'
// 			t[j+1] = upperhex[c>>4]
// 			t[j+2] = upperhex[c&15]
// 			j += 3
// 		default:
// 			t[j] = s[i]
// 			j++
// 		}
// 	}
// 	return string(t)
// }

// type URL struct {
// 	Scheme      string
// 	Method      string
// 	Path        string // path (relative paths may omit leading slash)
// 	RawPath     string // encoded path hint (see EscapedPath method)
// 	ForceQuery  bool   // append a query ('?') even if RawQuery is empty
// 	RawQuery    string // encoded query values, without '?'
// 	Fragment    string // fragment for references, without '#'
// 	RawFragment string // encoded fragment hint (see EscapedFragment method)
// }

// func (u *URL) setFragment(f string) error {
// 	frag, err := unescape(f, encodeFragment)
// 	if err != nil {
// 		return err
// 	}
// 	u.Fragment = frag
// 	if escf := escape(frag, encodeFragment); f == escf {
// 		u.RawFragment = ""
// 	} else {
// 		u.RawFragment = f
// 	}
// 	return nil
// }

// func (u *URL) setPath(p string) error {
// 	path, err := unescape(p, encodePath)
// 	if err != nil {
// 		return err
// 	}
// 	u.Path = path
// 	if escp := escape(path, encodePath); p == escp {
// 		u.RawPath = ""
// 	} else {
// 		u.RawPath = p
// 	}
// 	return nil
// }


// func Parse(rawURL string) (*URL, error) {
// 	u, frag, _ := strings.Cut(rawURL, "#")
// 	url, err := parse(u, false)
// 	if err != nil {
// 		return nil, &Error{"parse", u, err}
// 	}
// 	if frag == "" {
// 		return url, nil
// 	}
// 	if err = url.setFragment(frag); err != nil {
// 		return nil, &Error{"parse", rawURL, err}
// 	}
// 	return url, nil
// }

// func parse(rawURL string, viaRequest bool) (*URL, error) {
// 	var rest string
// 	var err error

// 	if stringContainsCTLByte(rawURL) {
// 		return nil, errors.New("core/url: invalid control character in URL")
// 	}

// 	if rawURL == "" && viaRequest {
// 		return nil, errors.New("empty url")
// 	}
// 	url := new(URL)

// 	if url.Scheme, rest, err = getScheme(rawURL); err != nil {
// 		return nil, err
// 	}

// 	url.Scheme = strings.ToLower(url.Scheme)

// 	if url.Scheme != "did" {
// 		return nil, errors.New("unkonw scheme")
// 	}

// 	if url.Method, rest, err = getMethod(rest); err != nil {
// 		return nil, err
// 	}
// 	// rest = example:123?service=agent&relativeRef=/credentials
// 	if strings.HasSuffix(rest, "?") && strings.Count(rest, "?") == 1 {
// 		url.ForceQuery = true
// 		rest = rest[:len(rest)-1]
// 	} else {
// 		// rest = example:123/test url.RawQuery= service=agent&relativeRef=/credentials
// 		rest, url.RawQuery, _ = strings.Cut(rest, "?")
// 	}

// 	if err := url.setPath(rest); err != nil {
// 		return nil, err
// 	}
// 	return url, nil
// }

// func getMethod(path string) (method, restpath string, err error) {
// 	for i := 0; i < len(path); i++ {
// 		c := path[i]
// 		switch {
// 		case 'a' <= c && c <= 'z' || '0' <= c && c <= '9':
// 		case c == ':':
// 			if i == 0 {
// 				return "", "", errors.New("missing protocol method")
// 			}
// 			return path[:i], path[i+1:], nil
// 		default:
// 			return "", path, nil
// 		}
// 	}
// 	return "", path, nil
// }

// func getScheme(rawURL string) (scheme, path string, err error) {
// 	for i := 0; i < len(rawURL); i++ {
// 		c := rawURL[i]
// 		switch {
// 		case 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z':
// 		case '0' <= c && c <= '9' || c == '+' || c == '-' || c == '.':
// 			if i == 0 {
// 				return "", rawURL, nil
// 			}
// 		case c == ':':
// 			if i == 0 {
// 				return "", "", errors.New("missing protocol scheme")
// 			}
// 			return rawURL[:i], rawURL[i+1:], nil
// 		default:
// 			return "", rawURL, nil
// 		}
// 	}
// 	return "", rawURL, nil
// }

// type Values map[string][]string

// func ParseQuery(query string) (Values, error) {
// 	m := make(Values)
// 	err := parseQuery(m, query)
// 	return m, err
// }

// func parseQuery(m Values, query string) (err error) {
// 	for query != "" {
// 		var key string
// 		key, query, _ = strings.Cut(query, "&")
// 		if strings.Contains(key, ";") {
// 			err = fmt.Errorf("invalid semicolon separator in query")
// 			continue
// 		}
// 		if key == "" {
// 			continue
// 		}
// 		key, value, _ := strings.Cut(key, "=")
// 		key, err1 := QueryUnescape(key)
// 		if err1 != nil {
// 			if err == nil {
// 				err = err1
// 			}
// 			continue
// 		}
// 		value, err1 = QueryUnescape(value)
// 		if err1 != nil {
// 			if err == nil {
// 				err = err1
// 			}
// 			continue
// 		}
// 		m[key] = append(m[key], value)
// 	}
// 	return err
// }

// func stringContainsCTLByte(s string) bool {
// 	for i := 0; i < len(s); i++ {
// 		b := s[i]
// 		if b < ' ' || b == 0x7f {
// 			return true
// 		}
// 	}
// 	return false
// }
