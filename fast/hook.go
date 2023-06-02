package fast

import (
	"encoding/json"
	"fmt"
	"io"
)

type validator interface {
	Validate() (interface{}, bool)
}

func Validate[T validator](reader io.Reader, data *T) (bool, error) {
	if err := json.NewDecoder(reader).Decode(data); err != nil {
		// Handle the error
		return false, err
	}
	fmt.Println(data)
	if errs, ok := (*data).Validate(); !ok {
		return false, fmt.Errorf("validation error: %v", errs)
	}

	return true, nil
}

/*
func ValidateRequestParams(ctx *Ctx, v any) bool {

	if v, ok := v.(Validator); ok {
		if errs, ok := v.Validate(); !ok {
			ctx.Status(http.StatusBadRequest).JSON(errs)
			return false
		}
	}
	return true
}

*/
