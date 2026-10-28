package http1

import (
	"net/http"
	"strings"
	"sync"
)

var DefaultMux = NewServeMux()

type MyMux struct {
	mu sync.RWMutex
	m  map[string]muxEntry //holds a map of URL patterns
}

type muxEntry struct {
	h       Handler //holds a handler
	pattern Pattern // matching pattern
}

type Pattern struct {
	uri     string
	smethod string
}

func (m *MyMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h := m.match(r.URL.Path)
	if h == nil {
		http.NotFound(w, r)
		return
	}
	h.ServeHTTP(w, r)
}

func (m *MyMux) Handle(pattern string, handler Handler) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if pattern == "" {
		panic("http: empty pattern")
	}
	if handler == nil {
		panic("http: nil handler")
	}
	if m.m == nil {
		m.m = make(map[string]muxEntry)
	}

	// check if Method Prefix exists
	method := ""
	uri := pattern
	if !strings.HasPrefix(pattern, "/") {
		splitPatter := strings.Split(pattern, " ")
		method = splitPatter[0]
		if _, ok := supporedHTTPMethods[strings.ToUpper(method)]; !ok {
			panic("http: method not supported ")
		}
		uri = splitPatter[1]
	}
	// move '/' back to end
	uri = rotateLeftEfficient(uri)

	if _, ok := m.m[pattern]; ok {
		if method != "" && m.m[pattern].pattern.smethod == method {
			m.m[pattern] = muxEntry{h: handler, pattern: Pattern{uri, method}}
		} else {
			panic("method is already registered")
		}
		panic("same pattern is already registered: " + pattern)
	}

	m.m[pattern] = muxEntry{h: handler, pattern: Pattern{uri, method}}
}

func (m *MyMux) match(path string) Handler {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if entry, ok := m.m[path]; ok {
		return entry.h
	}
	return nil
}

type HandlerFunc func(http.ResponseWriter, *http.Request)

type Handler interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}

func (m *MyMux) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	m.Handle(pattern, HandlerFunc(handler))
}

func (f HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f(w, r)
}

func NewServeMux() *MyMux {
	return &MyMux{m: make(map[string]muxEntry)}
}

func HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	DefaultMux.HandleFunc(pattern, handler)
}
