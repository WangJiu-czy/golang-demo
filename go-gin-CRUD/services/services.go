package services

import (
	"github.com/google/uuid"
	"go-gin-CRUD/domain"
	"go-gin-CRUD/repostiory"
)

type Services struct {
	Repostiory *repostiory.Repository
}

func NewServices() *Services {
	return &Services{
		Repostiory: repostiory.NewRepository(),
	}
}
func (s *Services) GetAllBooks() ([]domain.Book, error) {
	return s.Repostiory.GetAllBooks()
}
func (s *Services) GetBook(id string) (domain.Book, error) {
	return s.Repostiory.GetBook(id)
}
func (s *Services) CreateBook(book domain.Book) (domain.Book, int, error) {
	book.ID = uuid.New().String()
	book, err := s.Repostiory.CreateBook(book)
	if err != nil {
		return domain.Book{}, 400, err
	}
	return book, 201, nil
}
func (s *Services) UpdateBook(book domain.Book) (domain.Book, error) {
	if err := s.Repostiory.UpdateBook(book); err != nil {
		return domain.Book{}, err
	}
	return book, nil
}
func (s Services) DeleteBook(id string) error {
	if err := s.Repostiory.DeleteBook(id); err != nil {
		return err
	}
	return nil

}
