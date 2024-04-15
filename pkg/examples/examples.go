package examples

import (
	"fmt"
	"time"
)

type Another struct{}

func (a *Another) DoSomething() {
	fmt.Println("Another doing something")
}

// Metadata provides common metadata fields for various entities.
type Metadata struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// User defines a user in the system.
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Verified bool   `json:"verified"`
	Metadata        // Embedding Metadata struct
}

// Article represents a piece of written content.
type Article struct {
	ID          int                                   `json:"id"`
	Title       string                                `json:"title"`
	Content     string                                `json:"content"`
	Author      User                                  `json:"author"`
	Tags        []string                              `json:"tags"`
	Comments    []Comment                             `json:"comments"`
	PublishedAt time.Time                             `json:"published_at"`
	Status      PublicationStatus                     `json:"status"`
	Errors      error                                 `json:"error"`
	Anything    interface{}                           `json:"anything"`
	DoSomething map[string]interface{ DoSomething() } `json:"do_something"`
	Metadata                                          // Embedding Metadata struct
}

// Comment represents a user's comment on an article.
type Comment struct {
	ID        int    `json:"id"`
	ArticleID int    `json:"article_id"`
	Author    User   `json:"author"`
	Content   string `json:"content"`
	Metadata         // Embedding Metadata struct
}

// PublicationStatus is an enum representing the publication state of an article.
type PublicationStatus string

// CMSData is a struct embedding multiple other structs and showcasing a variety of types.
type CMSData struct {
	Users           []User        `json:"users"`
	Articles        []Article     `json:"articles"`
	ArticleComments map[int][]int `json:"article_comments"`
}
