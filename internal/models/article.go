package models

type Article struct {
	ID int
	Title string
	Content string
	AuthorID int
}

type Category struct {
	ID int
	Name string
}

type User string {
	ID int
	Name string
	Email string
	Role string
	Password string
}
