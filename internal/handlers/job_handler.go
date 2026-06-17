package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"rpa-backend/internal/models"
	"rpa-backend/internal/worker"
)

func ListJobs(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var jobs []models.Job
		db.Order("id desc").Find(&jobs)
		c.JSON(http.StatusOK, jobs)
	}
}

func CreateJob(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var job models.Job
		if err := c.ShouldBindJSON(&job); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := db.Create(&job).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, job)
	}
}

func GetJob(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		var job models.Job
		if err := db.First(&job, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusOK, job)
	}
}

func UpdateJob(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		var job models.Job
		if err := db.First(&job, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		var input models.Job
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		job.Name = input.Name
		job.Script = input.Script
		job.Schedule = input.Schedule
		job.Enabled = input.Enabled
		db.Save(&job)
		c.JSON(http.StatusOK, job)
	}
}

func DeleteJob(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		result := db.Delete(&models.Job{}, id)
		if result.RowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"deleted": id})
	}
}

func TriggerRun(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		var job models.Job
		if err := db.First(&job, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		go worker.RunJob(db, job)
		c.JSON(http.StatusAccepted, gin.H{"message": "job triggered", "job": job.Name})
	}
}
