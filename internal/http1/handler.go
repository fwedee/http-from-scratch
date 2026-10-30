package http1

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

var DefaultMux = NewServeMux()

type MyMux struct {
	mu sync.RWMutex
	m  map[string]map[string]muxEntry //holds a map of URL patterns
	dm map[string]map[string]muxEntry // holds the dynamics maps
}

type muxEntry struct {
	h       Handler //holds a handler
	pattern string
}

func (m *MyMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h, err := m.match(r.URL.Path, r.Method)

	if err != nil {
		http.NotFound(w, r)
		return
	}
	h.ServeHTTP(w, r)
}

func (m *MyMux) Handle(pattern string, handler Handler) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if pattern == "" {
		log.Fatal("HTTP1Handler: empty route detected")
		os.Exit(3)
	}
	if handler == nil {
		log.Fatal("HTTP1Handler: handler is nil detected")
		os.Exit(3)
	}
	if m.m == nil {
		m.m = make(map[string]map[string]muxEntry)
	}

	// check if Method Prefix exists
	method := ""
	uri := pattern
	if !strings.HasPrefix(pattern, "/") {
		splitPatter := strings.Split(pattern, " ")
		method = splitPatter[0]
		if _, ok := supporedHTTPMethods[strings.ToUpper(method)]; !ok {
			log.Fatal("HTTP1Handler: http method not supported")
			os.Exit(3)
		}
		uri = splitPatter[1]
	}
	// move '/' back to end
	// uri = rotateLeftEfficient(uri)

	if _, ok := m.m[uri]; ok {
		if method != "" {
			if _, exists := m.m[uri][method]; !exists {
				m.m[uri][method] = muxEntry{h: handler, pattern: pattern}
			} else {
				log.Fatal("Pattern is already registered")
				os.Exit(3)
			}
		} else {
			log.Fatal("Method is already registered")
			os.Exit(3)
		}
	}
	if _, exists := m.m[uri][method]; !exists {
		new_map := map[string]muxEntry{
			method: {h: handler, pattern: pattern}}
		m.m[uri] = new_map
	}
	m.m[uri][method] = muxEntry{h: handler, pattern: pattern}
}

func (m *MyMux) match(path string, method string) (Handler, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if entry, ok := m.m[path]; ok {

		if _, exists := entry[method]; exists {
			return entry[method].h, nil
		}

		if _, exists := entry[""]; exists {
			return entry[""].h, nil
		}
		return nil, errors.New("This method is not permitted")
	}

	// Todo: Check if the entry is a dynamic route
	// test/{1}/{2}
	//

	return nil, errors.New("This route is not here")
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
	return &MyMux{m: make(map[string]map[string]muxEntry)}
}

func HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	DefaultMux.HandleFunc(pattern, handler)
}
