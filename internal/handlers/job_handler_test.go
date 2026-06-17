package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gorm.io/gorm"

	"rpa-backend/internal/database"
	"rpa-backend/internal/handlers"
	"rpa-backend/internal/models"
)

func setupRouter() (*gin.Engine, *gorm.DB) {
	gin.SetMode(gin.TestMode)
	db := database.ConnectTest()
	r := gin.New()
	api := r.Group("/api")
	jobs := api.Group("/jobs")
	jobs.GET("", handlers.ListJobs(db))
	jobs.POST("", handlers.CreateJob(db))
	jobs.GET("/:id", handlers.GetJob(db))
	jobs.PUT("/:id", handlers.UpdateJob(db))
	jobs.DELETE("/:id", handlers.DeleteJob(db))
	jobs.POST("/:id/run", handlers.TriggerRun(db))
	return r, db
}

func TestJobCRUD(t *testing.T) {
	r, _ := setupRouter()

	body, _ := json.Marshal(models.Job{
		Name:    "test-job",
		Script:  "echo hello",
		Enabled: true,
	})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/jobs", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	var created models.Job
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &created))
	assert.Equal(t, "test-job", created.Name)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/api/jobs", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	var jobs []models.Job
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &jobs))
	assert.Len(t, jobs, 1)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/api/jobs/1", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	updateBody, _ := json.Marshal(models.Job{
		Name:    "updated-job",
		Script:  "echo updated",
		Enabled: false,
	})
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPut, "/api/jobs/1", bytes.NewReader(updateBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	var updated models.Job
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &updated))
	assert.Equal(t, "updated-job", updated.Name)
	assert.False(t, updated.Enabled)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodDelete, "/api/jobs/1", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/api/jobs/1", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetJobNotFound(t *testing.T) {
	r, _ := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/jobs/999", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestTriggerRun(t *testing.T) {
	r, db := setupRouter()

	job := models.Job{Name: "run-test", Script: "echo run", Enabled: true}
	require.NoError(t, db.Create(&job).Error)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/jobs/1/run", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusAccepted, w.Code)
}
