package worker_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"rpa-backend/internal/database"
	"rpa-backend/internal/models"
	"rpa-backend/internal/worker"
)

func TestRunJobSuccess(t *testing.T) {
	db := database.ConnectTest()
	job := models.Job{Name: "echo-job", Script: "echo hello-world", Enabled: true}
	require.NoError(t, db.Create(&job).Error)

	worker.RunJob(db, job)

	var run models.Run
	require.NoError(t, db.First(&run, "job_id = ?", job.ID).Error)
	assert.Equal(t, "success", run.Status)
	assert.Contains(t, run.Logs, "hello-world")
	assert.NotNil(t, run.StartedAt)
	assert.NotNil(t, run.FinishedAt)
	assert.True(t, run.FinishedAt.After(*run.StartedAt) || run.FinishedAt.Equal(*run.StartedAt))
}

func TestRunJobFailure(t *testing.T) {
	db := database.ConnectTest()
	job := models.Job{Name: "fail-job", Script: "exit 1", Enabled: true}
	require.NoError(t, db.Create(&job).Error)

	worker.RunJob(db, job)

	var run models.Run
	require.NoError(t, db.First(&run, "job_id = ?", job.ID).Error)
	assert.Equal(t, "failed", run.Status)
	assert.Contains(t, run.Logs, "[error]")
}

func TestRunJobRecordsTimestamps(t *testing.T) {
	db := database.ConnectTest()
	job := models.Job{Name: "ts-job", Script: "sleep 0", Enabled: true}
	require.NoError(t, db.Create(&job).Error)

	before := time.Now()
	worker.RunJob(db, job)

	var run models.Run
	require.NoError(t, db.First(&run, "job_id = ?", job.ID).Error)
	assert.True(t, run.StartedAt.After(before.Add(-time.Second)))
}
