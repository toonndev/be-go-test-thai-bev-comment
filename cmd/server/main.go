package main

import (
	"log"
	"os"

	commentHandler "comment-service/internal/adapter/handler"
	commentRepo "comment-service/internal/adapter/repository"
	"comment-service/internal/infrastructure/database"
	commentUseCase "comment-service/internal/usecase/comment"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	db, err := database.NewPostgresDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto migrate
	if err := db.AutoMigrate(&database.CommentTable{}); err != nil {
		log.Fatalf("Failed to auto migrate: %v", err)
	}

	// Wire dependencies
	repo := commentRepo.NewPostgresCommentRepo(db)
	uc := commentUseCase.NewUseCase(repo)
	h := commentHandler.NewCommentHandler(uc)

	r := gin.Default()

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	v1 := r.Group("/api/v1")
	{
		v1.GET("/comments", h.GetAllComments)
		v1.POST("/comments", h.CreateComment)
	}

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
