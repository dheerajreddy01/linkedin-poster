package db

import (
	"os"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"linkedin-poster/internal/models"
)

type Database struct {
	DB *gorm.DB
}

func Init(dbPath string) (*Database, error) {
	if dbPath == "" {
		dbPath = "./data/poster.db"
	}
	os.MkdirAll(filepath.Dir(dbPath), 0755)

	gormDB, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}

	gormDB.AutoMigrate(&models.Post{}, &models.NewsItem{}, &models.Settings{})

	// Default settings
	defaults := map[string]string{
		"author_name":      "Dheeraj Reddy",
		"post_frequency":   "daily",
		"auto_fetch":       "true",
		"newsapi_key":      "",
	}
	for k, v := range defaults {
		gormDB.FirstOrCreate(&models.Settings{Key: k, Value: v}, models.Settings{Key: k})
	}

	return &Database{DB: gormDB}, nil
}

func (d *Database) Get(key, def string) string {
	var s models.Settings
	if err := d.DB.Where("key = ?", key).First(&s).Error; err != nil {
		return def
	}
	return s.Value
}

func (d *Database) Set(key, value string) {
	var s models.Settings
	d.DB.Where(models.Settings{Key: key}).FirstOrCreate(&s)
	d.DB.Model(&s).Update("value", value)
}
