package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"rpa-backend/internal/models"
)

func GetStats(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var totalJobs, enabledJobs int64
		db.Model(&models.Job{}).Count(&totalJobs)
		db.Model(&models.Job{}).Where("enabled = ?", true).Count(&enabledJobs)

		today := time.Now().Truncate(24 * time.Hour)
		var runsToday, success, failed int64
		db.Model(&models.Run{}).Where("created_at >= ?", today).Count(&runsToday)
		db.Model(&models.Run{}).Where("created_at >= ? AND status = ?", today, "success").Count(&success)
		db.Model(&models.Run{}).Where("created_at >= ? AND status = ?", today, "failed").Count(&failed)

		c.JSON(http.StatusOK, gin.H{
			"total_jobs":    totalJobs,
			"enabled_jobs":  enabledJobs,
			"runs_today":    runsToday,
			"success_today": success,
			"failed_today":  failed,
		})
	}
}
