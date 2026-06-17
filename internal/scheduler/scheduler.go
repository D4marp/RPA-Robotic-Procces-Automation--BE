package scheduler

import (
	"log"

	"github.com/robfig/cron/v3"
	"gorm.io/gorm"

	"rpa-backend/internal/models"
	"rpa-backend/internal/worker"
)

func Start(db *gorm.DB) {
	c := cron.New()

	var jobs []models.Job
	db.Where("enabled = ? AND schedule != ?", true, "").Find(&jobs)

	for _, job := range jobs {
		j := job
		_, err := c.AddFunc(j.Schedule, func() {
			log.Printf("[scheduler] running job: %s", j.Name)
			go worker.RunJob(db, j)
		})
		if err != nil {
			log.Printf("[scheduler] invalid cron for job %s: %v", j.Name, err)
		}
	}

	c.Start()
	log.Printf("[scheduler] started, %d job(s) scheduled", len(jobs))
}
