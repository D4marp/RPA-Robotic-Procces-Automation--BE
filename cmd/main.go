package main

import (
	_ "rpa-backend/internal/initenv"

	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"rpa-backend/internal/config"
	"rpa-backend/internal/database"
	"rpa-backend/internal/handlers"
	"rpa-backend/internal/scheduler"
)

func main() {
	_ = godotenv.Load()

	db := database.Connect()
	database.Seed(db)
	scheduler.Start(db)

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
	}))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := r.Group("/api")
	{
		api.GET("/stats", handlers.GetStats(db))

		jobs := api.Group("/jobs")
		jobs.GET("", handlers.ListJobs(db))
		jobs.POST("", handlers.CreateJob(db))
		jobs.GET("/:id", handlers.GetJob(db))
		jobs.PUT("/:id", handlers.UpdateJob(db))
		jobs.DELETE("/:id", handlers.DeleteJob(db))
		jobs.POST("/:id/run", handlers.TriggerRun(db))

		runs := api.Group("/runs")
		runs.GET("", handlers.ListRuns(db))
		runs.GET("/:id", handlers.GetRun(db))
	}

	port := config.GetEnv("PORT", "8080")
	log.Fatal(r.Run(":" + port))
}
