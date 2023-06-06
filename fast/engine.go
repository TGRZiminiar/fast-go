package fast

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
)

type Config struct {
	//Default Port : 8080
	Port int
	//Default Logger : true
	Logger bool
}
type ErrorHandler func(*Ctx, error)

type Handler func(*Ctx) error

type Engine struct {
	config       Config
	errorHandler ErrorHandler
	router       *httprouter.Router
	corsConfig   CorsConfig
}
type CorsConfig struct {
	// Next defines a function to skip this middleware when returned true.
	//
	// Optional. Default: nil

	// AllowOriginsFunc defines a function that will set the 'access-control-allow-origin'
	// response header to the 'origin' request header when returned true.
	//
	// Optional. Default: nil
	AllowOriginsFunc func(origin string) bool

	// AllowOrigin defines a list of origins that may access the resource.
	//
	// Optional. Default value "*"
	AllowOrigins string

	// AllowMethods defines a list methods allowed when accessing the resource.
	// This is used in response to a preflight request.
	//
	// Optional. Default value "GET,POST,HEAD,PUT,DELETE,PATCH"
	AllowMethods string

	// AllowHeaders defines a list of request headers that can be used when
	// making the actual request. This is in response to a preflight request.
	//
	// Optional. Default value "".
	AllowHeaders string

	// AllowCredentials indicates whether or not the response to the request
	// can be exposed when the credentials flag is true. When used as part of
	// a response to a preflight request, this indicates whether or not the
	// actual request can be made using credentials.
	//
	// Optional. Default value false.
	AllowCredentials bool

	// ExposeHeaders defines a whitelist headers that clients are allowed to
	// access.
	//
	// Optional. Default value "".
	ExposeHeaders string

	// MaxAge indicates how long (in seconds) the results of a preflight request
	// can be cached.
	//
	// Optional. Default value 0.
	MaxAge int

	AllowCookies bool
}

const (
	MethodGet    string = "GET"
	MethodPost   string = "POST"
	MethodPut    string = "PUT"
	MethodPatch  string = "PATCH"
	MethodDelete string = "DELETE"
)

var ConfigDefault = CorsConfig{
	AllowOriginsFunc: nil,
	AllowOrigins:     "*",
	AllowMethods: strings.Join([]string{
		MethodGet,
		MethodPost,
		MethodPut,
		MethodPatch,
		MethodDelete,
	}, ","),
	AllowHeaders:     "Content-Type, withCredentials",
	AllowCredentials: true,
	ExposeHeaders:    "",
	MaxAge:           3600,
	AllowCookies:     true,
}

func Init(c Config) *Engine {

	if c.Port == 0 {
		c.Port = 8080
	}

	return &Engine{
		router:       httprouter.New(),
		errorHandler: func(ctx *Ctx, err error) {},
		config:       c,
		corsConfig:   ConfigDefault,
	}
}

func (e *Engine) Start() {
	fmt.Printf("server start on %d\n", e.config.Port)
	http.ListenAndServe(":"+strconv.Itoa(e.config.Port), e.router)
}

func (e *Engine) UpdateCorsConfig(config CorsConfig) {
	e.corsConfig = config
}

func (e *Engine) CORS(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		// Set the 'Access-Control-Allow-Origin' header to the origin of the client
		// making the request instead of a wildcard '*' when withCredentials is true
		// origin := r.Header.Get("Origin")
		// if origin != "" {
		// 	w.Header().Set("Access-Control-Allow-Origin", origin)
		// } else {
		// 	w.Header().Set("Access-Control-Allow-Origin", e.corsConfig.AllowOrigins)
		// }
		// fmt.Println(origin)

		w.Header().Set("Access-Control-Allow-Origin", e.corsConfig.AllowOrigins)
		w.Header().Set("Access-Control-Allow-Credentials", strconv.FormatBool(e.corsConfig.AllowCredentials))
		w.Header().Set("Access-Control-Allow-Headers", e.corsConfig.AllowHeaders)
		w.Header().Set("Access-Control-Allow-Methods", e.corsConfig.AllowMethods)
		w.Header().Set("Access-Control-Max-Age", strconv.Itoa(e.corsConfig.MaxAge))
		w.Header().Set("Access-Control-Expose-Headers", e.corsConfig.ExposeHeaders)

		if r.Method == "OPTIONS" {
			http.Error(w, "No Content", http.StatusNoContent)
			return
		}

		next(w, r, p)
	}
}

func (e *Engine) Logger(method string, path string, startTime time.Time) {
	elapsedTime := time.Since(startTime)
	switch method {
	case "GET":
		log.Printf("\x1b[42m%s\x1b[0m request for path %s completed in %s", method, path, elapsedTime)
	case "POST":
		log.Printf("\x1b[43m%s\x1b[0m request for path %s completed in %s", method, path, elapsedTime)
	case "PATCH":
		log.Printf("\x1b[44m%s\x1b[0m request for path %s completed in %s", method, path, elapsedTime)
	case "PUT":
		log.Printf("\x1b[45m%s\x1b[0m request for path %s completed in %s", method, path, elapsedTime)
	case "DELETE":
		log.Printf("\x1b[101m%s\x1b[0m request for path %s completed in %s", method, path, elapsedTime)
	}
}

func handleCtxMethod(path string, w http.ResponseWriter, r *http.Request, p httprouter.Params, e *Engine, method string) (*Ctx, bool, error) {
	ctx, servedFromCache := newCtx(w, r, p)
	if r.Method != method {
		return nil, servedFromCache, fmt.Errorf("method <%s> not allowed", r.Method)
	}

	return ctx, servedFromCache, nil

}
func createHandler(e *Engine, path string, h []Handler, method string) httprouter.Handle {

	fn := func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		if e.config.Logger {
			startTime := time.Now()
			defer e.Logger(r.Method, path, startTime)
		}

		ctx, cacheServe, err := handleCtxMethod(path, w, r, p, e, method)
		if err != nil {
			ctx.Status(http.StatusMethodNotAllowed).JSON(map[string]string{
				"error": err.Error(),
			})
			return
		}

		if cacheServe {
			return
		}

		for _, handler := range h {
			handler(ctx)
		}
	}

	return e.CORS(fn)
}

func (e *Engine) Get(path string, h ...Handler) {
	e.router.GET(path, createHandler(e, path, h, "GET"))
}

func (e *Engine) Post(path string, h ...Handler) {
	e.router.GET(path, createHandler(e, path, h, "POST"))
}

func (e *Engine) Patch(path string, h ...Handler) {
	e.router.GET(path, createHandler(e, path, h, "PATCH"))
}

func (e *Engine) Put(path string, h ...Handler) {
	e.router.GET(path, createHandler(e, path, h, "PUT"))
}

func (e *Engine) Delete(path string, h ...Handler) {
	e.router.GET(path, createHandler(e, path, h, "DELETE"))
}
