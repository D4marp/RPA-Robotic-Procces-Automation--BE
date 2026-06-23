package worker_test

import (
	"context"
	"fmt"
	"os/exec"
	"sync"
	"testing"
	"time"

	"gorm.io/gorm/logger"

	"rpa-backend/internal/database"
	"rpa-backend/internal/models"
	"rpa-backend/internal/worker"
)

func TestRPABenchmark(t *testing.T) {
	db := database.ConnectTest() // isolates tests in memory
	
	// Limit connection pool to 1 to avoid separate in-memory DB connections in SQLite
	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.SetMaxOpenConns(1)
	}

	// Silence GORM debug logs to keep test output clean
	db.Logger = db.Logger.LogMode(logger.Silent)

	fmt.Println("\n=== RPA SAAS ACADEMIC BENCHMARK RUN ===")

	// 1. LATENCY BY COMPLEXITY (Section 4.4.2)
	fmt.Println("\n[1] Pengukuran Latensi Eksekusi Skrip Bot berdasarkan Kompleksitas")
	scenarios := []struct {
		name       string
		complexity string
		script     string
	}{
		{"Simple Echo Bot", "Rendah", "echo 'Simple output'"},
		{"Medium Calculation Bot", "Sedang", "sleep 0.1 && echo 'Processed 100 items'"},
		{"High Data Processing Bot", "Tinggi", "sleep 0.5 && echo 'Completed heavy batch export'"},
	}

	for _, sc := range scenarios {
		job := models.Job{Name: sc.name, Script: sc.script, Enabled: true}
		if err := db.Create(&job).Error; err != nil {
			t.Fatalf("Failed to create job: %v", err)
		}

		start := time.Now()
		worker.RunJob(db, job)
		duration := time.Since(start)

		var run models.Run
		if err := db.First(&run, "job_id = ?", job.ID).Error; err != nil {
			t.Fatalf("Failed to find run: %v", err)
		}

		fmt.Printf("- Job: %-25s | Kompleksitas: %-6s | Latensi: %8.2f ms | Status: %s\n",
			sc.name, sc.complexity, float64(duration.Microseconds())/1000.0, run.Status)
	}

	// 2. DATABASE THROUGHPUT / BATCH PROCESSING (Section 4.4.2)
	fmt.Println("\n[2] Analisis Throughput Pangkalan Data (Concurrent Batch Processing)")
	concurrencies := []int{10, 20, 50}
	for _, c := range concurrencies {
		var wg sync.WaitGroup
		start := time.Now()

		for i := 0; i < c; i++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				job := models.Job{Name: fmt.Sprintf("Batch Job %d-%d", c, idx), Script: "echo 'batch'", Enabled: true}
				db.Create(&job)
				worker.RunJob(db, job)
			}(i)
		}
		wg.Wait()
		duration := time.Since(start)
		runsCreated := c
		throughput := float64(runsCreated) / duration.Seconds()

		fmt.Printf("- Concurrency: %3d | Waktu Total: %8.2f ms | Throughput: %8.2f jobs/sec\n",
			c, float64(duration.Milliseconds()), throughput)
	}

	// 3. EVENT-DRIVEN PIPELINE STABILITY (Section 4.4.3)
	fmt.Println("\n[3] Pengukuran Jeda Event Broker (Simulasi Apache Kafka & Consumer)")
	// We simulate a message broker delay (producer writes to broker channel, consumer processes from channel)
	eventBrokerChan := make(chan models.Run, 100)
	var eventDelays []time.Duration
	var delayWg sync.WaitGroup
	var delayMu sync.Mutex

	// Consumer Service Simulator
	go func() {
		for run := range eventBrokerChan {
			receivedTime := time.Now()
			// Simulate updating database asynchronously
			run.Status = "success"
			db.Save(&run)
			delay := receivedTime.Sub(*run.StartedAt)
			
			delayMu.Lock()
			eventDelays = append(eventDelays, delay)
			delayMu.Unlock()
			
			delayWg.Done()
		}
	}()

	// Producer triggers runs
	numEvents := 30
	for i := 0; i < numEvents; i++ {
		delayWg.Add(1)
		now := time.Now()
		run := models.Run{
			JobID:     1,
			Status:    "pending",
			StartedAt: &now,
		}
		db.Create(&run)
		
		// Simulate network propagation delay (100us - 900us)
		time.Sleep(time.Duration(100+ (i%5)*200) * time.Microsecond)
		eventBrokerChan <- run
	}

	delayWg.Wait()
	close(eventBrokerChan)

	var totalDelay time.Duration
	var minDelay = eventDelays[0]
	var maxDelay = eventDelays[0]
	for _, d := range eventDelays {
		totalDelay += d
		if d < minDelay {
			minDelay = d
		}
		if d > maxDelay {
			maxDelay = d
		}
	}
	avgDelay := float64(totalDelay.Nanoseconds()) / float64(numEvents) / 1e6 // in ms

	fmt.Printf("- Total Event Terkirim: %d\n", numEvents)
	fmt.Printf("- Rata-rata Jeda Transmisi: %.4f ms\n", avgDelay)
	fmt.Printf("- Jeda Minimum: %.4f ms\n", float64(minDelay.Nanoseconds())/1e6)
	fmt.Printf("- Jeda Maksimum: %.4f ms\n", float64(maxDelay.Nanoseconds())/1e6)

	// 4. FAULT TOLERANCE & RETRY STRATEGY (Section 4.4.5)
	fmt.Println("\n[4] Pengujian Toleransi Kesalahan (Third-party Failure & Retry Strategy)")

	// Simulate system timeout handling
	timeoutJob := models.Job{Name: "Timeout Bot", Script: "sleep 2", Enabled: true}
	db.Create(&timeoutJob)

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	start := time.Now()
	errChan := make(chan error, 1)

	go func() {
		cmd := exec.CommandContext(ctx, "sh", "-c", timeoutJob.Script)
		_, err := cmd.CombinedOutput()
		errChan <- err
	}()

	var executionErr error
	select {
	case err := <-errChan:
		executionErr = err
	case <-ctx.Done():
		executionErr = ctx.Err()
	}
	elapsed := time.Since(start)

	fmt.Printf("- Skenario Timeout   | Batas Waktu: 200.00 ms | Terjadi Setelah: %6.2f ms | Status Error: %v\n",
		float64(elapsed.Milliseconds()), executionErr)

	// Simulate Retry Strategy with Rollback
	fmt.Println("- Simulasi Retry Strategy (Maksimum 3 Percobaan) dengan Data Rollback:")
	maxRetries := 3
	attempt := 0
	successOnAttempt := -1 // will never succeed in this test to demonstrate rollback
	
	tx := db.Begin() // start transaction for rollback demonstration
	jobError := models.Job{Name: "Unstable Payment Gateway Bot", Script: "exit 5", Enabled: true}
	tx.Create(&jobError)

	for attempt = 1; attempt <= maxRetries; attempt++ {
		fmt.Printf("  * Percobaan %d: Memicu eksekusi bot...", attempt)
		
		// Run job mock
		cmd := exec.Command("sh", "-c", jobError.Script)
		_, err := cmd.CombinedOutput()
		
		if err == nil {
			fmt.Println(" [SUKSES]")
			successOnAttempt = attempt
			tx.Commit()
			break
		} else {
			fmt.Printf(" [GAGAL] - Status Error: %s. Melakukan retry dalam %d ms...\n", err.Error(), attempt*10)
			time.Sleep(time.Duration(attempt*10) * time.Millisecond)
		}
	}

	if successOnAttempt == -1 {
		fmt.Println("  * Semua percobaan gagal. Melakukan ROLLBACK transaksi basis data...")
		tx.Rollback()
		
		// Verify if the job indeed got rolled back (should not be in DB)
		var checkJob models.Job
		err := db.First(&checkJob, "name = ?", jobError.Name).Error
		fmt.Printf("  * Status Rollback Basis Data: %t (Job tidak tersimpan karena transaksi di-rollback)\n", err != nil)
	}

	fmt.Println("\n=== BENCHMARK RUN COMPLETED ===")
}
