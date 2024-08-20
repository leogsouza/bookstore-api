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

type UserSuite struct {
	suite.Suite
	db         *gorm.DB
	mock       sqlmock.Sqlmock
	repository Repository[model.User]
}

func (s *UserSuite) SetupTest() {
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

	s.repository = NewUserRepo[model.User](s.db)
}

func (s *UserSuite) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func TestInit(t *testing.T) {
	suite.Run(t, new(UserSuite))
}

func (s *UserSuite) TestFind() {

	queryUser := `
		SELECT * FROM "users" WHERE "users"."id" = $1 ORDER BY "users"."id" LIMIT $2
	`
	id := 1
	name := "User 1"
	email := "test@user.com"

	s.mock.ExpectQuery(regexp.QuoteMeta(queryUser)).WithArgs(id, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email"}).
			AddRow(id, name, email))

	res, err := s.repository.Find(id)
	require.NoError(s.T(), err)
	require.Nil(s.T(), deep.Equal(model.User{ID: id, Name: name, Email: email}, res))

	s.mock.ExpectQuery(regexp.QuoteMeta(queryUser)).
		WithArgs(id, 1).
		WillReturnError(fmt.Errorf("user not found"))

	_, err = s.repository.Find(id)

	require.NotNil(s.T(), err)
}

func (s *UserSuite) TestFindAll() {
	queryUser := `
		SELECT * FROM "users"
	`
	rows := []model.User{
		{ID: 1, Name: "Customer 1", Email: "customer1@test.com"},
		{ID: 2, Name: "Customer 2", Email: "customer2@test.com"},
		{ID: 3, Name: "Customer 3", Email: "customer3@test.com"},
	}

	s.mock.ExpectQuery(regexp.QuoteMeta(queryUser)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email"}).
			AddRow(rows[0].ID, rows[0].Name, rows[0].Email).
			AddRow(rows[1].ID, rows[1].Name, rows[1].Email).
			AddRow(rows[2].ID, rows[2].Name, rows[2].Email))

	res, err := s.repository.FindAll()

	require.Nil(s.T(), err)
	require.Len(s.T(), res, len(rows))

}

func (s *UserSuite) TestCreate() {
	user := model.User{
		ID:    1,
		Name:  "Customer xyz",
		Email: "customer@xyz.com",
	}

	s.mock.ExpectBegin()
	s.mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "users" ("name","email","password","id") VALUES ($1,$2,$3,$4) RETURNING "id"`)).
		WithArgs(user.Name, user.Email, user.Password, user.ID).
		WillReturnRows(
			sqlmock.NewRows([]string{"id"}).AddRow(user.ID))

	s.mock.ExpectCommit()
	id, err := s.repository.Create(user)
	require.NoError(s.T(), err)
	require.Equal(s.T(), user.ID, id)

	s.mock.ExpectBegin()
	s.mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "users" ("name","email","password","id") VALUES ($1,$2,$3,$4) RETURNING "id"`)).
		WithArgs(user.Name, user.Email, user.Password, user.ID).
		WillReturnError(fmt.Errorf("user already exists"))

	s.mock.ExpectRollback()
	_, err = s.repository.Create(user)
	require.NotNil(s.T(), err)

}
