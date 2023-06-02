package launch

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type Validator interface {
	Validate() (any, bool)
}

type ErrorHandler func(*Ctx, error)

type PostHandler[T any] func(*PostCtx[T]) error

type Handler func(*Ctx) error

const (
	ApplicationJSON   = "application/json"
	MultipartFormData = "multipart/form-data"
)

func validateRequestParams(ctx *Ctx, v any) bool {
	if v, ok := v.(Validator); ok {
		if errs, ok := v.Validate(); !ok {
			ctx.Status(http.StatusBadRequest).JSON(errs)
			return false
		}
	}
	return true
}

type Ctx struct {
	r          *http.Request
	w          http.ResponseWriter
	params     httprouter.Params
	statusCode int
}

func newCtx(w http.ResponseWriter, r *http.Request, params httprouter.Params) *Ctx {
	return &Ctx{
		r:          r,
		w:          w,
		statusCode: http.StatusOK,
		params:     params,
	}
}

type Param string

func (p Param) AsInt(ctx *Ctx) int {
	val, err := strconv.Atoi(string(p))
	if err != nil {
		m := map[string]string{"error": fmt.Sprintf("param <%s> not of type <int>", p)}
		ctx.Status(http.StatusBadRequest).JSON(m)
		return 0
	}
	return val
}

func (c *Ctx) Param(name string) Param {
	return Param(c.params.ByName(name))
}
func (c *PostCtx[T]) RequestParams() T {
	return c.params
}

func (c *Ctx) Status(s int) *Ctx {
	c.statusCode = s
	return c
}

func (c *PostCtx[T]) FormValue(name string) string {
	return c.r.FormValue(name)
}

func (c *PostCtx[T]) ManyFormValue(name ...string) []string {
	multiData := make([]string, 5, 10)
	for _, v := range name {
		multiData = append(multiData, c.r.FormValue(v))
	}
	return multiData
}

func (c *Ctx) JSON(v any) error {
	c.w.WriteHeader(c.statusCode)
	c.w.Header().Add("Content-Type", ApplicationJSON)
	return json.NewEncoder(c.w).Encode(v)
}

type launch struct {
	errorHandler ErrorHandler
	router       *httprouter.Router
}

var App = &launch{
	errorHandler: func(ctx *Ctx, err error) {},
	router:       httprouter.New(),
}

type PostCtx[T any] struct {
	Ctx
	params T
}

func Start() {
	http.ListenAndServe(":5000", App.router)
}

func Get(path string, h Handler) {
	fn := func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		ctx := newCtx(w, r, p)
		if r.Method != "GET" {
			ctx.Status(http.StatusMethodNotAllowed).JSON(map[string]string{
				"error": fmt.Sprintf("method <%s> not allowed", r.Method),
			})
			return
		}
		h(ctx)
	}
	App.router.GET(path, fn)
}

func Delete(path string, h Handler) {
	fn := func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		ctx := newCtx(w, r, p)
		if r.Method != "Delete" {
			ctx.Status(http.StatusMethodNotAllowed).JSON(map[string]string{
				"error": fmt.Sprintf("method <%s> not allowed", r.Method),
			})
			return
		}
		h(ctx)
	}
	App.router.DELETE(path, fn)
}

func Post[T any](path string, h PostHandler[T]) {
	fn := func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		var (
			params T
			ctx    = newCtx(w, r, p)
		)
		if r.Method != "POST" {
			ctx.Status(http.StatusMethodNotAllowed).JSON(map[string]string{
				"error": fmt.Sprintf("method <%s> not allowed", r.Method),
			})
			return
		}
		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			App.errorHandler(ctx, err)
			return
		}
		if !validateRequestParams(ctx, params) {
			ctx.Status(http.StatusBadRequest).JSON(map[string]string{
				"error": "some field in validation is missing",
			})
			return
		}

		postCtx := &PostCtx[T]{
			Ctx:    *ctx,
			params: params,
		}

		err := h(postCtx)
		if err != nil {
			App.errorHandler(ctx, err)
			return
		}

	}
	App.router.POST(path, fn)
}

func Put[T any](path string, h PostHandler[T]) {
	fn := func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		var (
			params T
			ctx    = newCtx(w, r, p)
		)
		if r.Method != "Put" {
			ctx.Status(http.StatusMethodNotAllowed).JSON(map[string]string{
				"error": fmt.Sprintf("method <%s> not allowed", r.Method),
			})
			return
		}
		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			App.errorHandler(ctx, err)
			return
		}
		if !validateRequestParams(ctx, params) {
			ctx.Status(http.StatusBadRequest).JSON(map[string]string{
				"error": "some field in validation is missing",
			})
			return
		}

		postCtx := &PostCtx[T]{
			Ctx:    *ctx,
			params: params,
		}

		err := h(postCtx)
		if err != nil {
			App.errorHandler(ctx, err)
			return
		}

	}
	App.router.PUT(path, fn)
}

func Patch[T any](path string, h PostHandler[T]) {
	fn := func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		var (
			params T
			ctx    = newCtx(w, r, p)
		)
		if r.Method != "Patch" {
			ctx.Status(http.StatusMethodNotAllowed).JSON(map[string]string{
				"error": fmt.Sprintf("method <%s> not allowed", r.Method),
			})
			return
		}
		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			App.errorHandler(ctx, err)
			return
		}
		if !validateRequestParams(ctx, params) {
			ctx.Status(http.StatusBadRequest).JSON(map[string]string{
				"error": "some field in validation is missing",
			})
			return
		}

		postCtx := &PostCtx[T]{
			Ctx:    *ctx,
			params: params,
		}

		err := h(postCtx)
		if err != nil {
			App.errorHandler(ctx, err)
			return
		}

	}
	App.router.PATCH(path, fn)
}
