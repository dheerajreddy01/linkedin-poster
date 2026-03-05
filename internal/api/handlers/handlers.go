package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"linkedin-poster/internal/ai"
	"linkedin-poster/internal/db"
	"linkedin-poster/internal/linkedin"
	"linkedin-poster/internal/models"
)

type Handler struct {
	db *db.Database
	ai *ai.Service
	li *linkedin.Client
}

func New(database *db.Database, aiSvc *ai.Service) *Handler {
	return &Handler{
		db: database,
		ai: aiSvc,
		li: linkedin.New(),
	}
}

// GET /api/posts
func (h *Handler) GetPosts(c *gin.Context) {
	status := c.Query("status")
	q := h.db.DB.Model(&models.Post{})
	if status != "" {
		q = q.Where("status = ?", status)
	}
	var posts []models.Post
	q.Order("created_at DESC").Limit(50).Find(&posts)
	c.JSON(http.StatusOK, posts)
}

// GET /api/posts/stats
func (h *Handler) GetStats(c *gin.Context) {
	var stats models.PostStats
	h.db.DB.Model(&models.Post{}).Count(&stats.Total)
	h.db.DB.Model(&models.Post{}).Where("status = ?", "draft").Count(&stats.Drafts)
	h.db.DB.Model(&models.Post{}).Where("status = ?", "approved").Count(&stats.Approved)
	h.db.DB.Model(&models.Post{}).Where("status = ?", "posted").Count(&stats.Posted)
	h.db.DB.Model(&models.Post{}).Where("status = ?", "rejected").Count(&stats.Rejected)
	c.JSON(http.StatusOK, stats)
}

// PUT /api/posts/:id/approve
func (h *Handler) ApprovePost(c *gin.Context) {
	h.db.DB.Model(&models.Post{}).Where("id = ?", c.Param("id")).Update("status", "approved")
	c.JSON(http.StatusOK, gin.H{"message": "approved"})
}

// PUT /api/posts/:id/reject
func (h *Handler) RejectPost(c *gin.Context) {
	h.db.DB.Model(&models.Post{}).Where("id = ?", c.Param("id")).Update("status", "rejected")
	c.JSON(http.StatusOK, gin.H{"message": "rejected"})
}

// PUT /api/posts/:id/edit
func (h *Handler) EditPost(c *gin.Context) {
	var body struct {
		Content string `json:"content"`
	}
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.db.DB.Model(&models.Post{}).Where("id = ?", c.Param("id")).Update("content", body.Content)
	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}

// POST /api/posts/:id/regenerate
func (h *Handler) RegeneratePost(c *gin.Context) {
	var post models.Post
	if err := h.db.DB.First(&post, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	var body struct {
		Instruction string `json:"instruction"`
	}
	c.BindJSON(&body)

	instruction := body.Instruction
	if instruction == "" {
		instruction = "rewrite with a fresh angle and different opening hook"
	}

	newContent, err := h.ai.RegeneratePost(post.Content, instruction)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.db.DB.Model(&post).Updates(map[string]interface{}{
		"content": newContent,
		"status":  "draft",
	})

	post.Content = newContent
	c.JSON(http.StatusOK, post)
}

// POST /api/posts/:id/post
func (h *Handler) PublishPost(c *gin.Context) {
	var post models.Post
	if err := h.db.DB.First(&post, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	if err := h.li.PostToLinkedIn(post.Content); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	now := time.Now()
	h.db.DB.Model(&post).Updates(map[string]interface{}{
		"status":    "posted",
		"posted_at": &now,
	})

	c.JSON(http.StatusOK, gin.H{"message": "posted to LinkedIn!"})
}

// GET /api/settings
func (h *Handler) GetSettings(c *gin.Context) {
	var settings []models.Settings
	h.db.DB.Find(&settings)
	result := map[string]string{}
	for _, s := range settings {
		result[s.Key] = s.Value
	}
	c.JSON(http.StatusOK, result)
}

// PUT /api/settings
func (h *Handler) UpdateSettings(c *gin.Context) {
	var body map[string]string
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	for k, v := range body {
		h.db.Set(k, v)
	}
	c.JSON(http.StatusOK, gin.H{"message": "saved"})
}
