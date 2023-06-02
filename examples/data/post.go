package data

import "fmt"

type Post struct {
	ID          int
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (p Post) Validate() (interface{}, bool) {
	fmt.Println("valudate", p)

	if len(p.Title) < 5 {
		return map[string]string{
			"title": "too short",
		}, false
	}

	if len(p.Description) < 5 {
		return map[string]string{
			"description": "too short",
		}, false
	}

	return nil, true
}
