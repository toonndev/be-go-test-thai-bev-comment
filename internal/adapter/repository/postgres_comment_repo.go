package repository

import (
	"context"
	"time"

	domain "comment-service/internal/domain/comment"

	"gorm.io/gorm"
)

type postgresCommentRepo struct {
	db *gorm.DB
}

type commentModel struct {
	ID         uint      `gorm:"primaryKey;autoIncrement"`
	AuthorName string    `gorm:"column:author_name;not null"`
	Content    string    `gorm:"column:content;not null"`
	CreatedAt  time.Time `gorm:"column:created_at;not null;autoCreateTime"`
	UpdatedAt  time.Time `gorm:"column:updated_at;not null;autoUpdateTime"`
}

func (commentModel) TableName() string {
	return "comments"
}

func NewPostgresCommentRepo(db *gorm.DB) domain.Repository {
	return &postgresCommentRepo{db: db}
}

func (r *postgresCommentRepo) Create(ctx context.Context, comment *domain.Comment) (*domain.Comment, error) {
	model := &commentModel{
		AuthorName: comment.AuthorName,
		Content:    comment.Content,
	}
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return nil, err
	}
	comment.ID = uint(model.ID)
	comment.CreatedAt = model.CreatedAt
	comment.UpdatedAt = model.UpdatedAt
	return comment, nil
}

func (r *postgresCommentRepo) FindAll(ctx context.Context) ([]*domain.Comment, error) {
	var models []commentModel
	if err := r.db.WithContext(ctx).Order("created_at DESC").Find(&models).Error; err != nil {
		return nil, err
	}

	comments := make([]*domain.Comment, len(models))
	for i, m := range models {
		comments[i] = &domain.Comment{
			ID:         uint(m.ID),
			AuthorName: m.AuthorName,
			Content:    m.Content,
			CreatedAt:  m.CreatedAt,
			UpdatedAt:  m.UpdatedAt,
		}
	}
	return comments, nil
}
