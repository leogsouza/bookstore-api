package repository

import (
	"bookstore-api/internal/model"
	"database/sql"
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-test/deep"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type BookSuite struct {
	suite.Suite
	db         *gorm.DB
	mock       sqlmock.Sqlmock
	repository Reader[model.Book]
}

func (s *BookSuite) SetupTest() {
	var (
		db  *sql.DB
		err error
	)

	db, s.mock, err = sqlmock.New()
	require.NoError(s.T(), err)

	dialector := postgres.New(postgres.Config{
		Conn:       db,
		DriverName: "postgres",
	})

	s.db, err = gorm.Open(dialector)
	require.NoError(s.T(), err)

	s.repository = NewBookRepo[model.Book](s.db)
}

func (s *BookSuite) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func TestInitBook(t *testing.T) {
	suite.Run(t, new(BookSuite))
}

func (s *BookSuite) TestFind() {

	queryBook := `
		SELECT * FROM "books" WHERE "books"."id" = $1 ORDER BY "books"."id" LIMIT $2
	`
	id := 1
	title := "Book 1"
	author := "Author Name"
	price := 10.75

	s.mock.ExpectQuery(regexp.QuoteMeta(queryBook)).WithArgs(id, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "author", "price"}).
			AddRow(id, title, author, price))

	res, err := s.repository.Find(id)
	require.NoError(s.T(), err)
	require.Nil(s.T(), deep.Equal(model.Book{ID: id, Title: title, Author: author, Price: price}, res))

	s.mock.ExpectQuery(regexp.QuoteMeta(queryBook)).
		WithArgs(id, 1).
		WillReturnError(fmt.Errorf("book not found"))

	_, err = s.repository.Find(id)

	require.NotNil(s.T(), err)
}

func (s *BookSuite) TestFindAll() {
	queryBook := `
		SELECT * FROM "books"
	`
	rows := []model.Book{
		{ID: 1, Title: "Book 1", Author: "Author 1", Price: 10.25},
		{ID: 2, Title: "Book 2", Author: "Author 2", Price: 13.25},
		{ID: 3, Title: "Book 3", Author: "Author 3", Price: 12.25},
	}

	s.mock.ExpectQuery(regexp.QuoteMeta(queryBook)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "author", "price"}).
			AddRow(rows[0].ID, rows[0].Title, rows[0].Author, rows[0].Price).
			AddRow(rows[1].ID, rows[1].Title, rows[1].Author, rows[0].Price).
			AddRow(rows[2].ID, rows[2].Title, rows[2].Author, rows[0].Price))

	res, err := s.repository.FindAll()

	require.Nil(s.T(), err)
	require.Len(s.T(), res, len(rows))

}
