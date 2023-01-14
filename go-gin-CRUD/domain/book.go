package domain

type Book struct {
	ID     string `json:"id"`
	Title  string `json:"title" binding:"required"`
	Year   int    `json:"year" binding:"required"`
	Author string `json:"author" binding:"required"`
}
