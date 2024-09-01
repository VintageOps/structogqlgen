package examples

import "time"

// AdditionalData Example is a struct embedding multiple other structs and showcasing a variety of types.
type AdditionalData struct {
	Users           []User        `json:"users"`
	Articles        []Article     `json:"articles"`
	ArticleComments map[int][]int `json:"article_comments"`
}

type extraArgs int64

// Non exported Structs that has an anonymous unexported Struct
type structWithAnonymous struct {
	Name        string `json:"name"`
	Age         int    `json:"age"`
	UnNamedType struct {
		Name      string    `json:"name"`
		CreatedAt time.Time `json:"created_at"`
		DeletedAt time.Time `json:"deleted_at"`
	} `json:"unnamed_type"`
}
