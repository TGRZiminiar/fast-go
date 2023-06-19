package params

import (
	"encoding/json"
	"fmt"

	"github.com/TGRZiminiar/based/fast"
)

type CreatePost struct {
	Name string `json:"name"`
}

func (p CreatePost) Validate(ctx *fast.Ctx) (CreatePost, error) {

	var params CreatePost
	if err := json.NewDecoder(ctx.R.Body).Decode(&params); err != nil {
		return CreatePost{}, err
	}

	if len(params.Name) < 3 {
		return CreatePost{}, fmt.Errorf("name is too short %d", len(params.Name))
	}
	return params, nil
}
