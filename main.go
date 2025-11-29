package main

import (
	"fmt"
	"time"
)

type Book struct {
	Title       string
	Author      string
	ISBN        string
	IsAvailable bool
}

func (b *Book) MarkAsLoaned() {
	b.IsAvailable = false
}

func (b *Book) MarkAsAvailable() {
	b.IsAvailable = true
}

type Reader struct {
	Id    int
	Name  string
	Email string
}

func (r *Reader) BorrowBook(book *Book, librarian *Librarian) error {
	if !book.IsAvailable {
		return fmt.Errorf("книга '%s' недоступна", book.Title)
	}
	loan := &Loan{
		Book:      book,
		Reader:    r,
		LoanDate:  time.Now(),
		ReturnDate: time.Time{}, 
	}
	loan.IssueLoan()
	librarian.loans = append(librarian.loans, loan)
	return nil
}

func (r *Reader) ReturnBook(book *Book, librarian *Librarian) error {
	for _, loan := range librarian.loans {
		if loan.Book.ISBN == book.ISBN && loan.Reader.Id == r.Id && loan.ReturnDate.IsZero() {
			loan.CompleteLoan()
			return nil
		}
	}
	return fmt.Errorf("выдача не найдена")
}

type Loan struct {
	Book       *Book
	Reader     *Reader
	LoanDate   time.Time
	ReturnDate time.Time
}

func (l *Loan) IssueLoan() {
	l.Book.MarkAsLoaned()
}

func (l *Loan) CompleteLoan() {
	l.ReturnDate = time.Now()
	l.Book.MarkAsAvailable()
}

type Librarian struct {
	Id       int
	Name     string
	Position string

	books    []*Book
	readers  []*Reader
	loans    []*Loan
}

func NewLibrarian(id int, name, position string) *Librarian {
	return &Librarian{
		Id:       id,
		Name:     name,
		Position: position,
		books:    make([]*Book, 0),
		readers:  make([]*Reader, 0),
		loans:    make([]*Loan, 0),
	}
}

func (l *Librarian) AddBook(book *Book) {
	l.books = append(l.books, book)
}

func (l *Librarian) RemoveBook(isbn string) {
	for i, b := range l.books {
		if b.ISBN == isbn {
			l.books = append(l.books[:i], l.books[i+1:]...)
			return
		}
	}
}

func (l *Librarian) AddReader(reader *Reader) {
	l.readers = append(l.readers, reader)
}

func (l *Librarian) RemoveReader(id int) {
	for i, r := range l.readers {
		if r.Id == id {
			l.readers = append(l.readers[:i], l.readers[i+1:]...)
			return
		}
	}
}

func (l *Librarian) SearchBooks(query string) []*Book {
	var result []*Book
	query = fmt.Sprintf("%s", query)
	for _, b := range l.books {
		if b.Title == query || b.Author == query {
			result = append(result, b)
		}
	}
	return result
}

func (l *Librarian) GetAvailableBooks() []*Book {
	var result []*Book
	for _, b := range l.books {
		if b.IsAvailable {
			result = append(result, b)
		}
	}
	return result
}

func (l *Librarian) GetLoanedBooks() []*Book {
	var result []*Book
	for _, loan := range l.loans {
		if loan.ReturnDate.IsZero() {
			result = append(result, loan.Book)
		}
	}
	return result
}

func main() {
	librarian := NewLibrarian(1, "Анна Иванова", "Главный библиотекарь")

	book1 := &Book{Title: "1984", Author: "Джордж Оруэлл", ISBN: "978-0-452-28423-4", IsAvailable: true}
	book2 := &Book{Title: "Гордость и предубеждение", Author: "Джейн Остин", ISBN: "978-0-14-143951-8", IsAvailable: true}
	librarian.AddBook(book1)
	librarian.AddBook(book2)

	reader := &Reader{Id: 101, Name: "Иван Петров", Email: "ivan@example.com"}
	librarian.AddReader(reader)

	err := reader.BorrowBook(book1, librarian)
	if err != nil {
		fmt.Println("Ошибка:", err)
	} else {
		fmt.Println("Книга выдана:", book1.Title)
	}

	found := librarian.SearchBooks("1984")
	fmt.Println("Найдено книг:", len(found))

	fmt.Println("Доступные книги:", len(librarian.GetAvailableBooks()))
	fmt.Println("Выданные книги:", len(librarian.GetLoanedBooks()))
}
