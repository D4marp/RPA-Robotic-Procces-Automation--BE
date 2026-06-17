package worker

import (
	"os/exec"
	"time"

	"gorm.io/gorm"

	"rpa-backend/internal/models"
)

func RunJob(db *gorm.DB, job models.Job) {
	run := models.Run{JobID: job.ID, Status: "running"}
	now := time.Now()
	run.StartedAt = &now
	db.Create(&run)

	cmd := exec.Command("sh", "-c", job.Script)
	out, err := cmd.CombinedOutput()

	fin := time.Now()
	run.FinishedAt = &fin
	run.Logs = string(out)
	if err != nil {
		run.Status = "failed"
		run.Logs += "\n[error] " + err.Error()
	} else {
		run.Status = "success"
	}
	db.Save(&run)
}
