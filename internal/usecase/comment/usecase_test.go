package comment_test

import (
	"context"
	"errors"
	"testing"

	domain "comment-service/internal/domain/comment"
	usecase "comment-service/internal/usecase/comment"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository is a mock implementation of domain.Repository
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(ctx context.Context, comment *domain.Comment) (*domain.Comment, error) {
	args := m.Called(ctx, comment)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Comment), args.Error(1)
}

func (m *MockRepository) FindAll(ctx context.Context) ([]*domain.Comment, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Comment), args.Error(1)
}

func TestCreateComment_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	uc := usecase.NewUseCase(mockRepo)

	input := &domain.Comment{AuthorName: "Blend 285", Content: "have a good day"}
	expected := &domain.Comment{ID: 1, AuthorName: "Blend 285", Content: "have a good day"}

	mockRepo.On("Create", mock.Anything, input).Return(expected, nil)

	result, err := uc.CreateComment(context.Background(), "Blend 285", "have a good day")

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}

func TestCreateComment_EmptyContent(t *testing.T) {
	mockRepo := new(MockRepository)
	uc := usecase.NewUseCase(mockRepo)

	result, err := uc.CreateComment(context.Background(), "Blend 285", "")

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertNotCalled(t, "Create")
}

func TestCreateComment_EmptyAuthorName(t *testing.T) {
	mockRepo := new(MockRepository)
	uc := usecase.NewUseCase(mockRepo)

	result, err := uc.CreateComment(context.Background(), "", "have a good day")

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertNotCalled(t, "Create")
}

func TestGetAllComments_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	uc := usecase.NewUseCase(mockRepo)

	expected := []*domain.Comment{
		{ID: 1, AuthorName: "Blend 285", Content: "have a good day"},
		{ID: 2, AuthorName: "Alice", Content: "hello world"},
	}

	mockRepo.On("FindAll", mock.Anything).Return(expected, nil)

	result, err := uc.GetAllComments(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}

func TestGetAllComments_Empty(t *testing.T) {
	mockRepo := new(MockRepository)
	uc := usecase.NewUseCase(mockRepo)

	mockRepo.On("FindAll", mock.Anything).Return([]*domain.Comment{}, nil)

	result, err := uc.GetAllComments(context.Background())

	assert.NoError(t, err)
	assert.Empty(t, result)
	assert.NotNil(t, result)
	mockRepo.AssertExpectations(t)

	// verify no DB error
	assert.False(t, errors.Is(err, errors.New("db error")))
}
