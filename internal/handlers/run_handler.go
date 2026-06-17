package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"rpa-backend/internal/models"
)

func ListRuns(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var runs []models.Run
		db.Preload("Job").Order("id desc").Limit(100).Find(&runs)
		c.JSON(http.StatusOK, runs)
	}
}

func GetRun(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		var run models.Run
		if err := db.Preload("Job").First(&run, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusOK, run)
	}
}
