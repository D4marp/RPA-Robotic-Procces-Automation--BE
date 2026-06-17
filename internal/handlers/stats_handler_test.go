package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"rpa-backend/internal/database"
	"rpa-backend/internal/handlers"
	"rpa-backend/internal/models"
)

func TestGetStats(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := database.ConnectTest()

	db.Create(&models.Job{Name: "j1", Script: "echo 1", Enabled: true})
	j2 := models.Job{Name: "j2", Script: "echo 2"}
	db.Create(&j2)
	db.Model(&j2).Update("enabled", false)
	db.Create(&models.Run{JobID: 1, Status: "success", CreatedAt: time.Now()})
	db.Create(&models.Run{JobID: 1, Status: "failed", CreatedAt: time.Now()})

	r := gin.New()
	r.GET("/api/stats", handlers.GetStats(db))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/stats", nil)
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var stats map[string]int64
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &stats))
	assert.Equal(t, int64(2), stats["total_jobs"])
	assert.Equal(t, int64(1), stats["enabled_jobs"])
	assert.Equal(t, int64(2), stats["runs_today"])
	assert.Equal(t, int64(1), stats["success_today"])
	assert.Equal(t, int64(1), stats["failed_today"])
}
