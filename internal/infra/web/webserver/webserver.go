package webserver

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Handler struct {
	method  string
	path    string
	handler http.HandlerFunc
}

type WebServer struct {
	router        chi.Router
	handlers      []Handler
	webServerPort string
}

func NewWebServer(serverPort string) *WebServer {
	return &WebServer{
		router:        chi.NewRouter(),
		handlers:      make([]Handler, 0),
		webServerPort: serverPort,
	}
}

func (s *WebServer) AddHandler(method string, path string, handler http.HandlerFunc) {
	s.handlers = append(s.handlers, Handler{method: method, path: path, handler: handler})
}

// loop through the handlers and add them to the router
// register middeleware logger
// start the server
func (s *WebServer) Start() {
	s.router.Use(middleware.Logger)
	for _, handler := range s.handlers {
		s.router.Method(handler.method, handler.path, handler.handler)
	}

	http.ListenAndServe(s.webServerPort, s.router)
}
