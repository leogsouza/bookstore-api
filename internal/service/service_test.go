package service

import (
	"bookstore-api/internal/model"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository[T model.Entity] struct {
	mock.Mock
}

func (m *MockRepository[T]) Find(id int) (model.User, error) {
	args := m.Called()
	result := args.Get(0)
	return result.(model.User), args.Error(1)
}

func (m *MockRepository[T]) FindAll() ([]model.User, error) {
	args := m.Called()
	result := args.Get(0)
	return result.([]model.User), args.Error(1)
}

func (m *MockRepository[T]) FindByCondition(condition string, args interface{}) (*model.User, error) {
	res := m.Called(condition, args)
	result := res.Get(0)
	return result.(*model.User), res.Error(1)
}

func (m *MockRepository[T]) Create(user model.User) (*model.User, error) {
	args := m.Called()
	result := args.Get(0)
	return result.(*model.User), args.Error(1)
}

func TestFind(t *testing.T) {
	// was not able to mock generic
	mockRepo := &MockRepository[model.User]{}
	id := 1

	entity := model.User{ID: 1, Name: "Customer 1", Email: "customer@gmail.com"}
	mockRepo.On("Find").Return(entity, nil)

	testService := NewService[model.User](mockRepo)

	result, _ := testService.Find(id)

	mockRepo.AssertExpectations(t)

	assert.Equal(t, 1, result.ID)
	assert.Equal(t, entity.Name, result.Name)
	assert.Equal(t, entity.Email, result.Email)

}

func TestFindAll(t *testing.T) {
	mockRepo := &MockRepository[model.User]{}

	entity := model.User{ID: 1, Name: "Customer 1", Email: "customer@gmail.com"}
	mockRepo.On("FindAll").Return([]model.User{entity}, nil)

	testService := NewService(mockRepo)

	result, _ := testService.FindAll()

	mockRepo.AssertExpectations(t)

	assert.Equal(t, 1, result[0].ID)
	assert.Equal(t, entity.Name, result[0].Name)
	assert.Equal(t, entity.Email, result[0].Email)

}

func TestCreate(t *testing.T) {
	mockRepo := &MockRepository[model.User]{}

	entity := model.User{ID: 1, Name: "Customer 1", Email: "customer@gmail.com"}
	mockRepo.On("Create").Return(entity.ID, nil)

	testService := NewService(mockRepo)

	result, _ := testService.Create(entity)

	mockRepo.AssertExpectations(t)

	assert.Equal(t, entity.ID, result)

}
