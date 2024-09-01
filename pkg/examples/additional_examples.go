package examples

// AdditionalData Example is a struct embedding multiple other structs and showcasing a variety of types.
type AdditionalData struct {
	Users           []User        `json:"users"`
	Articles        []Article     `json:"articles"`
	ArticleComments map[int][]int `json:"article_comments"`
}

type extraArgs int64
