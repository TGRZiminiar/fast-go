package fast

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"reflect"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/patrickmn/go-cache"
)

const (
	ApplicationJSON   = "application/json"
	MultipartFormData = "multipart/form-data"
)

var (
	cacheData *cache.Cache
)

func init() {
	cacheData = cache.New(5*time.Minute, 10*time.Minute) // Initialize the cache
}

type CustomRequest struct {
	*http.Request
	Temp interface{} // Assuming Temp Variable
}

type Ctx struct {
	W          http.ResponseWriter
	R          *CustomRequest
	params     httprouter.Params
	statusCode int
}

func newCtx(w http.ResponseWriter, r *http.Request, params httprouter.Params) (*Ctx, bool) {
	path := r.URL.Path
	if cachedData, found := cacheData.Get(path); found {
		w.Header().Set("Content-Type", ApplicationJSON)
		w.Write(cachedData.([]byte))
		return nil, true // Return an empty Ctx, as the response is served from the cache
	}
	return &Ctx{
		W:          w,
		R:          &CustomRequest{Request: r},
		params:     params,
		statusCode: http.StatusOK,
	}, false
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

func (c *Ctx) JSON(v interface{}) error {
	c.W.Header().Set("Content-Type", ApplicationJSON)

	// Check if the response data is already cached
	cacheKey := c.R.URL.Path

	jsonData, err := json.Marshal(v)
	if err != nil {
		return err
	}

	// Cache the response data
	cacheData.Set(cacheKey, jsonData, cache.DefaultExpiration)

	c.W.WriteHeader(c.statusCode)
	c.W.Write(jsonData)
	return nil
}

func (c *Ctx) FormValue(name string) string {
	return c.R.FormValue(name)
}

func (c *Ctx) ManyFormKeyValue(name ...string) map[string]string {
	multiData := make(map[string]string)
	// c.R.ParseForm()

	// test := make(map[string](chan string))

	// go func() {
	for _, v := range name {
		multiData[v] = c.R.FormValue(v)
		// test[v] <- c.R.FormValue(v)
	}
	// }()

	// temp := ChanToSlice(test).(map[string]string)

	return multiData
}

func (c *Ctx) ManyFormValue(name ...string) []string {

	multiData := make(chan string)
	// multiData := []string{}
	go func() {
		for _, v := range name {
			// multiData = append(multiData, c.R.FormValue(v))
			multiData <- c.R.FormValue(v)
		}
		close(multiData)
	}()
	temp := ChanToSlice(multiData).([]string)
	return temp
}

func (c *Ctx) FormFile(name string) (multipart.File, *multipart.FileHeader, error) {

	// c.R.ParseForm()

	file, header, err := c.R.FormFile(name)

	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	return file, header, nil
}

func ChanToSlice(ch interface{}) interface{} {
	chv := reflect.ValueOf(ch)
	slv := reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(ch).Elem()), 0, 0)
	for {
		v, ok := chv.Recv()
		if !ok {
			return slv.Interface()
		}
		slv = reflect.Append(slv, v)
	}
}

func (c *Ctx) FormManyFiles(key ...string) ([]*multipart.FileHeader, error) {

	c.R.ParseMultipartForm(32 << 20)

	fhs := c.R.MultipartForm.File["files"]
	if len(fhs) == 0 {
		return nil, fmt.Errorf("no value in this key")
	}
	multiFile := make(chan *multipart.FileHeader)

	go func() {
		for _, fh := range fhs {
			multiFile <- fh
		}
		close(multiFile)
	}()

	temp := ChanToSlice(multiFile).([]*multipart.FileHeader)

	return temp, nil

}

func (c *Ctx) SetCookie(data http.Cookie) {

	cook := &http.Cookie{
		Name:       data.Name,
		Value:      data.Value,
		MaxAge:     data.MaxAge,
		Path:       data.Path,
		Domain:     data.Domain,
		Secure:     data.Secure,
		HttpOnly:   data.HttpOnly,
		Expires:    data.Expires,
		RawExpires: data.RawExpires,
		SameSite:   data.SameSite,
		Raw:        data.Raw,
		Unparsed:   data.Unparsed,
	}
	http.SetCookie(c.W, cook)

}

func (c *Ctx) DeleteCookie(cookieName string) {
	http.SetCookie(c.W, &http.Cookie{
		Name:   cookieName,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
}
