package handler

import (
	"bookstore-api/internal/model"
	"bytes"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type userMockService struct {
	mock.Mock
}

func (m *userMockService) Find(id int) (model.User, error) {
	args := m.Called(id)
	var user model.User

	if rf, ok := args.Get(0).(func(int) model.User); ok {
		user = rf(id)
	} else if args.Get(0) != nil {
		user = args.Get(0).(model.User)
	}

	var err error
	if rf, ok := args.Get(1).(func(int) error); ok {
		err = rf(id)
	} else {
		err = args.Error(1)
	}

	return user, err
}

func (m *userMockService) FindAll() ([]model.User, error) {
	args := m.Called()
	result := args.Get(0)
	return result.([]model.User), args.Error(1)
}

func (m *userMockService) Create(user model.User) (int, error) {
	args := m.Called(user)

	var result int
	if rf, ok := args.Get(0).(func(model.User) int); ok {
		result = rf(user)
	} else if args.Get(0) != nil {
		result = args.Get(0).(int)
	}

	var err error
	if rf, ok := args.Get(1).(func(model.User) error); ok {
		err = rf(user)
	} else {
		err = args.Error(1)
	}
	return result, err
}

func TestGetAllUsers(t *testing.T) {
	t.Run("GetAllUsers - Success", func(t *testing.T) {
		mockService := new(userMockService)
		mockService.On("FindAll").Return(getUsers(), nil)

		userHandler := NewUserHandler(mockService)

		app := fiber.New()

		app.Get("/users", userHandler.GetAll)

		req := httptest.NewRequest("GET", "/users", nil)

		resp, err := app.Test(req)
		assert.Nil(t, err)
		assert.Equal(t, resp.StatusCode, fiber.StatusOK)

		buf := &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		enc.Encode(getUsers())

		usersResponse := []model.User{}

		if err := json.NewDecoder(resp.Body).Decode(&usersResponse); err != nil {
			assert.Nil(t, err)
		}

		assert.NotEmpty(t, usersResponse)
		assert.Len(t, usersResponse, len(getUsers()))

		mockService.AssertExpectations(t)

	})

}

func TestGetUser(t *testing.T) {
	t.Run("GetUser - Success", func(t *testing.T) {
		mockService := new(userMockService)
		mockService.On("Find", 1).Return(getUser(), nil)

		userHandler := NewUserHandler(mockService)

		app := fiber.New()

		app.Get("/users/:id", userHandler.Get)

		req := httptest.NewRequest("GET", "/users/1", nil)

		resp, err := app.Test(req)
		assert.Nil(t, err)
		assert.Equal(t, resp.StatusCode, fiber.StatusOK)

		buf := &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		enc.Encode(getUser())

		userResponse := model.User{}

		if err := json.NewDecoder(resp.Body).Decode(&userResponse); err != nil {
			assert.Nil(t, err)
		}

		user := getUser()

		assert.NotEmpty(t, userResponse)
		assert.Equal(t, userResponse.Name, user.Name)
		assert.Equal(t, userResponse.Email, user.Email)

		mockService.AssertExpectations(t)

	})

	t.Run("GetUser - Error ", func(t *testing.T) {
		mockService := new(userMockService)
		mockService.On("Find", 1).Return(nil, errors.New("user not found"), nil)

		userHandler := NewUserHandler(mockService)

		app := fiber.New()

		app.Get("/users/:id", userHandler.Get)

		req := httptest.NewRequest("GET", "/users/1", nil)

		resp, err := app.Test(req)
		assert.Nil(t, err)
		assert.Equal(t, resp.StatusCode, fiber.StatusNotFound)

		mockService.AssertExpectations(t)

	})

}

func TestCreateUser(t *testing.T) {
	t.Run("CreateUser - Success", func(t *testing.T) {

		userReq := userRequest{
			Name:     "Customer 1",
			Email:    "customer@test.com",
			Password: "custpass123",
		}
		buf := &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		enc.Encode(userReq)

		b := bytes.NewReader(buf.Bytes())

		userModel := model.User{
			Name:     userReq.Name,
			Email:    userReq.Email,
			Password: userReq.Password,
		}

		mockService := new(userMockService)
		mockService.On("Create", userModel).Return(1, nil)

		userHandler := NewUserHandler(mockService)

		app := fiber.New()

		app.Post("/users", userHandler.Post)

		req := httptest.NewRequest("POST", "/users", b)
		req.Header.Add(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		resp, err := app.Test(req)
		assert.Nil(t, err)
		assert.Equal(t, resp.StatusCode, fiber.StatusOK)

		userResponse := model.User{}

		if err := json.NewDecoder(resp.Body).Decode(&userResponse); err != nil {
			assert.Nil(t, err)
		}

		user := getUser()

		assert.NotEmpty(t, userResponse)
		assert.Equal(t, userResponse.Name, user.Name)
		assert.Equal(t, userResponse.Email, user.Email)

		mockService.AssertExpectations(t)

	})

	t.Run("CreateUser - Error ", func(t *testing.T) {

		userReq := userRequest{
			Name:     "Customer 1",
			Email:    "customer@test.com",
			Password: "",
		}
		buf := &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		enc.Encode(userReq)

		b := bytes.NewReader(buf.Bytes())

		userModel := model.User{
			Name:     userReq.Name,
			Email:    userReq.Email,
			Password: userReq.Password,
		}

		mockService := new(userMockService)
		mockService.On("Create", userModel).Return(nil, errors.New("user not found"))

		userHandler := NewUserHandler(mockService)

		app := fiber.New()

		app.Post("/users", userHandler.Post)

		req := httptest.NewRequest("POST", "/users", b)
		req.Header.Add(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		resp, err := app.Test(req)
		assert.Nil(t, err)
		assert.Equal(t, resp.StatusCode, fiber.StatusInternalServerError)

		mockService.AssertExpectations(t)

	})

}

func getUser() model.User {
	return model.User{ID: 1, Name: "Customer 1", Email: "customer@test.com", Password: "custpass123"}
}

func getUsers() []model.User {
	return []model.User{getUser()}
}
