package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"linkedin-poster/internal/ai"
	"linkedin-poster/internal/api/handlers"
	"linkedin-poster/internal/db"
	"linkedin-poster/internal/scheduler"
)

func main() {
	godotenv.Load()

	database, err := db.Init(os.Getenv("DB_PATH"))
	if err != nil {
		log.Fatalf("DB error: %v", err)
	}

	aiSvc := ai.New()
	sched := scheduler.New(database, aiSvc)
	sched.Start()
	defer sched.Stop()

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3000", "http://127.0.0.1:5500"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders: []string{"Content-Type"},
	}))

	h := handlers.New(database, aiSvc)

	api := r.Group("/api")
	{
		api.GET("/posts", h.GetPosts)
		api.GET("/posts/stats", h.GetStats)
		api.PUT("/posts/:id/approve", h.ApprovePost)
		api.PUT("/posts/:id/reject", h.RejectPost)
		api.PUT("/posts/:id/edit", h.EditPost)
		api.POST("/posts/:id/regenerate", h.RegeneratePost)
		api.POST("/posts/:id/post", h.PublishPost)
		api.GET("/settings", h.GetSettings)
		api.PUT("/settings", h.UpdateSettings)
	}

	// Serve frontend
	r.Static("/app", "./frontend")
	r.GET("/", func(c *gin.Context) {
		c.Redirect(302, "/app/index.html")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("🚀 LinkedIn Poster running at http://localhost:%s", port)
	r.Run(":" + port)
}
