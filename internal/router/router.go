package router

import (
	"dev-challenge/internal/cell"
	"dev-challenge/internal/sheet"
	"dev-challenge/internal/utils"
	"encoding/json"
	"log"
	"net/http"
	"regexp"
)

type Ctx struct {
	Response http.ResponseWriter
	Request  *http.Request
	Params   map[string]string
}

type Executor func(*Ctx)

type handler struct {
	pattern  *regexp.Regexp
	executor Executor
}

type Router struct {
	sheetService *sheet.Service
	cellService  *cell.Service

	// http method to slice of handlers map
	handlers map[string][]handler
}

func New(sheetService *sheet.Service, cellService *cell.Service) *Router {
	rt := &Router{
		sheetService: sheetService,
		cellService:  cellService,

		handlers: make(map[string][]handler),
	}

	rt.establishRoutes()

	return rt
}

func (rt *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handlers, ok := rt.handlers[r.Method]
	if !ok {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
		return
	}

	for _, handler := range handlers {
		if handler.pattern.MatchString(r.URL.Path) {
			params := utils.ParseUrlPath(handler.pattern, r.URL.Path)

			ctx := Ctx{
				Response: w,
				Request:  r,
				Params:   params,
			}

			log.Printf("%s %s", r.Method, r.URL.Path)
			handler.executor(&ctx)
			return
		}
	}

	log.Printf("miss. 404 %s %s", r.Method, r.URL.Path)
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(http.StatusText(http.StatusNotFound)))
}

func (rt *Router) regisetHandler(method, pattern string, executor Executor) {
	h := handler{
		pattern:  regexp.MustCompile(pattern),
		executor: executor,
	}

	if rt.handlers == nil {
		rt.handlers = make(map[string][]handler)
	}

	if _, ok := rt.handlers[method]; !ok {
		rt.handlers[method] = make([]handler, 0)
	}

	rt.handlers[method] = append(rt.handlers[method], h)
}

func (rt *Router) Get(pattern string, executor Executor) {
	rt.regisetHandler(http.MethodGet, pattern, executor)
}

func (rt *Router) Post(pattern string, executor Executor) {
	rt.regisetHandler(http.MethodPost, pattern, executor)
}

func respondJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Println(err)
	}
}
