package comment

import (
	"context"
	"errors"

	domain "comment-service/internal/domain/comment"
)

type UseCase struct {
	repo domain.Repository
}

func NewUseCase(repo domain.Repository) *UseCase {
	return &UseCase{repo: repo}
}

func (uc *UseCase) CreateComment(ctx context.Context, authorName, content string) (*domain.Comment, error) {
	if authorName == "" {
		return nil, errors.New("author_name and content are required")
	}
	if content == "" {
		return nil, errors.New("author_name and content are required")
	}

	c := &domain.Comment{
		AuthorName: authorName,
		Content:    content,
	}
	return uc.repo.Create(ctx, c)
}

func (uc *UseCase) GetAllComments(ctx context.Context) ([]*domain.Comment, error) {
	return uc.repo.FindAll(ctx)
}
