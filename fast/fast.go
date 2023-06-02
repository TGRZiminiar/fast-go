package fast

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
)

const (
	ApplicationJSON   = "application/json"
	MultipartFormData = "multipart/form-data"
)

type Ctx struct {
	w          http.ResponseWriter
	r          *http.Request
	params     httprouter.Params
	statusCode int
}

func newCtx(w http.ResponseWriter, r *http.Request, params httprouter.Params) *Ctx {
	return &Ctx{
		w:          w,
		r:          r,
		params:     params,
		statusCode: http.StatusOK,
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

func (c *Ctx) Status(s int) *Ctx {
	c.statusCode = s
	return c
}

func (c *Ctx) JSON(v any) error {
	c.w.WriteHeader(c.statusCode)
	c.w.Header().Add("Content-Type", ApplicationJSON)
	return json.NewEncoder(c.w).Encode(v)
}

func (c *Ctx) Request() *http.Request {
	return c.r
}

func (c *Ctx) Writer() http.ResponseWriter {
	return c.w
}

func (c *Ctx) FormValue(name string) string {
	return c.r.FormValue(name)
}

func (c *Ctx) ManyFormKeyValue(name ...string) map[string]string {
	multiData := make(map[string]string)

	for _, v := range name {
		multiData[v] = c.r.FormValue(v)
	}

	return multiData
}

func (c *Ctx) ManyFormValue(name ...string) []string {
	multiData := []string{}
	for _, v := range name {
		multiData = append(multiData, c.r.FormValue(v))
	}
	return multiData
}

func (c *Ctx) FormFile(name string) (multipart.File, *multipart.FileHeader, error) {
	file, header, err := c.r.FormFile(name)
	if err != nil {
		return nil, nil, err
	}
	return file, header, nil
}

type Cookie struct {
	Name        string    `json:"name"`
	Value       string    `json:"value"`
	Path        string    `json:"path"`
	Domain      string    `json:"domain"`
	MaxAge      int       `json:"max_age"`
	Expires     time.Time `json:"expires"`
	Secure      bool      `json:"secure"`
	HTTPOnly    bool      `json:"http_only"`
	SameSite    string    `json:"same_site"`
	SessionOnly bool      `json:"session_only"`
}

func (c *Ctx) Cookie(name, value string, maxAge int, path, domain string, secure, httpOnly bool) {

	cook := &http.Cookie{
		Name:     name,
		Value:    value,
		MaxAge:   maxAge,
		Path:     path,
		Domain:   domain,
		Secure:   secure,
		HttpOnly: httpOnly,
	}

	http.SetCookie(c.w, cook)

}
