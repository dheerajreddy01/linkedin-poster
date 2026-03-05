package models

import (
	"time"
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Title      string    `json:"title"`
	Content    string    `json:"content"`     // AI-generated LinkedIn post
	Topic      string    `json:"topic"`       // Go, DS, AWS, AI, etc.
	SourceURL  string    `json:"source_url"`
	SourceName string    `json:"source_name"`
	Status     string    `json:"status" gorm:"default:'draft'"` // draft, approved, posted, rejected
	PostedAt   *time.Time `json:"posted_at"`
	Likes      int       `json:"likes"`
	Comments   int       `json:"comments"`
}

type NewsItem struct {
	gorm.Model
	Title       string    `json:"title"`
	URL         string    `json:"url"`
	Source      string    `json:"source"`
	Topic       string    `json:"topic"`
	Summary     string    `json:"summary"`
	PublishedAt time.Time `json:"published_at"`
	Processed   bool      `json:"processed" gorm:"default:false"`
}

type Settings struct {
	gorm.Model
	Key   string `json:"key" gorm:"uniqueIndex"`
	Value string `json:"value"`
}

type PostStats struct {
	Total    int64 `json:"total"`
	Drafts   int64 `json:"drafts"`
	Approved int64 `json:"approved"`
	Posted   int64 `json:"posted"`
	Rejected int64 `json:"rejected"`
}
