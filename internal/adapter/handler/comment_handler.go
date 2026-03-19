package handler

import (
	"net/http"

	usecase "comment-service/internal/usecase/comment"

	"github.com/gin-gonic/gin"
)

type CommentHandler struct {
	uc *usecase.UseCase
}

func NewCommentHandler(uc *usecase.UseCase) *CommentHandler {
	return &CommentHandler{uc: uc}
}

type createCommentRequest struct {
	AuthorName string `json:"author_name"`
	Content    string `json:"content"`
}

func (h *CommentHandler) GetAllComments(c *gin.Context) {
	comments, err := h.uc.GetAllComments(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": comments})
}

func (h *CommentHandler) CreateComment(c *gin.Context) {
	var req createCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "author_name and content are required"})
		return
	}

	comment, err := h.uc.CreateComment(c.Request.Context(), req.AuthorName, req.Content)
	if err != nil {
		if err.Error() == "author_name and content are required" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": comment})
}
