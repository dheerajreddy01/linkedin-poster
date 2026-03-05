package scheduler

import (
	"log"
	"time"

	"github.com/robfig/cron/v3"
	"linkedin-poster/internal/ai"
	"linkedin-poster/internal/db"
	"linkedin-poster/internal/models"
	"linkedin-poster/internal/news"
)

type Scheduler struct {
	cron    *cron.Cron
	db      *db.Database
	ai      *ai.Service
	fetcher *news.Fetcher
}

func New(database *db.Database, aiSvc *ai.Service) *Scheduler {
	return &Scheduler{
		cron:    cron.New(),
		db:      database,
		ai:      aiSvc,
		fetcher: news.New(),
	}
}

func (s *Scheduler) Start() {
	// Fetch news & generate drafts every 6 hours
	s.cron.AddFunc("0 */6 * * *", s.fetchAndGenerate)
	s.cron.Start()

	// Run once on startup after 3 seconds
	go func() {
		time.Sleep(3 * time.Second)
		s.fetchAndGenerate()
	}()

	log.Println("✅ Scheduler started — fetching every 6 hours")
}

func (s *Scheduler) Stop() { s.cron.Stop() }

func (s *Scheduler) fetchAndGenerate() {
	log.Println("⏰ Fetching news and generating post drafts...")

	apiKey := s.db.Get("newsapi_key", "")
	authorName := s.db.Get("author_name", "Dheeraj Reddy")

	// Fetch fresh news
	items := s.fetcher.FetchAll(apiKey)
	if len(items) == 0 {
		log.Println("No news items fetched")
		return
	}

	newItems := 0
	generatedPosts := 0

	for _, item := range items {
		// Skip if already processed
		var existing models.NewsItem
		if err := s.db.DB.Where("url = ?", item.URL).First(&existing).Error; err == nil {
			continue
		}

		// Save news item
		s.db.DB.Create(&item)
		newItems++

		// Generate LinkedIn post draft
		content, err := s.ai.GeneratePost(item, authorName)
		if err != nil {
			log.Printf("❌ GPT error for '%s': %v", item.Title, err)
			continue
		}

		post := models.Post{
			Title:      item.Title,
			Content:    content,
			Topic:      item.Topic,
			SourceURL:  item.URL,
			SourceName: item.Source,
			Status:     "draft",
		}
		s.db.DB.Create(&post)

		// Mark news item as processed
		s.db.DB.Model(&item).Update("processed", true)
		generatedPosts++
	}

	log.Printf("✅ %d new articles → %d post drafts generated", newItems, generatedPosts)
}
