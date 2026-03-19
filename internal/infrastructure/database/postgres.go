package database

import (
	"fmt"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// CommentTable is the GORM model used for auto migration.
type CommentTable struct {
	ID         uint      `gorm:"primaryKey;autoIncrement"`
	AuthorName string    `gorm:"column:author_name;not null"`
	Content    string    `gorm:"column:content;not null"`
	CreatedAt  time.Time `gorm:"column:created_at;not null;autoCreateTime"`
	UpdatedAt  time.Time `gorm:"column:updated_at;not null;autoUpdateTime"`
}

func (CommentTable) TableName() string { return "comments" }

func NewPostgresDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=UTC",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}
