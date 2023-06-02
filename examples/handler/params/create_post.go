package params

type CreatePost struct {
	Name string `formData:"name"`
}

func (p CreatePost) Validate() (any, bool) {

	if len(p.Name) < 3 {
		return map[string]string{
			"title": "too short",
		}, false
	}
	return nil, true
}
