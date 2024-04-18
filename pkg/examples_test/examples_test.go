package examples_test

import (
	"fmt"
	"time"
)

// Another Example Is just an example Struct
type Another struct{}

// DoSomething Examples prints a message indicating that 'Another' is doing something.
// This method belongs to the 'Another' struct.
func (a *Another) DoSomething() {
	fmt.Println("Another doing something")
}

// Metadata Example provides common metadata fields for various entities.
type Metadata struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CachedAt  time.Time `json:"-"`
}

// User Example defines a user in the system.
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Verified bool   `json:"verified"`
	Metadata        // Embedding Metadata struct
}

// Article Example represents a piece of written content.
type Article struct {
	ID                 int                                   `json:"id" validate:"required"`
	Title              string                                `json:"title"`
	Content            string                                `json:"content"`
	Author             User                                  `json:"author"`
	Tags               []string                              `json:"tags"`
	Comments           []Comment                             `json:"comments"`
	PublishedAt        time.Time                             `json:"published_at"`
	Status             PublicationStatus                     `json:"status"`
	Errors             error                                 `json:"error"`
	Anything           interface{}                           `json:"anything"`
	DoSomething        map[string]interface{ DoSomething() } `json:"do_something"`
	RandomInt          uint64                                `json:"random_int"`
	AnotherRandomInt64 uint64                                `json:"another_random_int64"`
	Metadata                                                 // Embedding Metadata struct
}

// Comment Example represents a user's comment on an article.
type Comment struct {
	ID        int    `json:"id"`
	ArticleID int    `json:"article_id"`
	Author    User   `json:"author"`
	Content   string `json:"content"`
	Metadata         // Embedding Metadata struct
}

// PublicationStatus Example is an enum representing the publication state of an article.
type PublicationStatus string

// CMSData Example is a struct embedding multiple other structs and showcasing a variety of types.
type CMSData struct {
	Users           []User        `json:"users"`
	Articles        []Article     `json:"articles"`
	ArticleComments map[int][]int `json:"article_comments"`
}
