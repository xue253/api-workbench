package scheduler

import (
	"log"
	"sync"

	"api-workbench/internal/db"
	"api-workbench/internal/engine"
	"api-workbench/internal/model"

	"github.com/robfig/cron/v3"
)

var (
	cronScheduler *cron.Cron
	jobs          map[uint]cron.EntryID
	mu            sync.RWMutex
)

func Init() {
	cronScheduler = cron.New()
	jobs = make(map[uint]cron.EntryID)
	cronScheduler.Start()
	loadExistingTasks()
}

func loadExistingTasks() {
	var tasks []model.ScheduledTask
	db.DB.Where("enabled = ?", true).Find(&tasks)
	for _, task := range tasks {
		addJob(task)
	}
	log.Printf("Loaded %d scheduled tasks", len(tasks))
}

func addJob(task model.ScheduledTask) {
	mu.Lock()
	defer mu.Unlock()

	if _, exists := jobs[task.ID]; exists {
		return
	}

	entryID, err := cronScheduler.AddFunc(task.CronExpr, func() {
		runScheduledTask(task)
	})
	if err != nil {
		log.Printf("Failed to schedule task %d: %v", task.ID, err)
		return
	}

	jobs[task.ID] = entryID
	log.Printf("Scheduled task %d with cron: %s", task.ID, task.CronExpr)
}

func runScheduledTask(task model.ScheduledTask) {
	log.Printf("Executing scheduled task %d", task.ID)

	opts := &engine.RunOptions{
		EnvID: task.EnvironmentID,
	}

	switch task.TargetType {
	case "test_case":
		_, err := engine.RunTestCaseWithWS(task.TargetID, opts)
		if err != nil {
			log.Printf("Failed to run test case %d: %v", task.TargetID, err)
		}
	case "test_suite":
		_, err := engine.RunTestSuiteWithWS(task.TargetID, opts)
		if err != nil {
			log.Printf("Failed to run test suite %d: %v", task.TargetID, err)
		}
	}
}

func AddTask(task model.ScheduledTask) {
	if task.Enabled {
		addJob(task)
	}
}

func RemoveTask(taskID uint) {
	mu.Lock()
	defer mu.Unlock()

	if entryID, exists := jobs[taskID]; exists {
		cronScheduler.Remove(entryID)
		delete(jobs, taskID)
		log.Printf("Removed scheduled task %d", taskID)
	}
}

func UpdateTask(task model.ScheduledTask) {
	RemoveTask(task.ID)
	if task.Enabled {
		addJob(task)
	}
}

func Stop() {
	if cronScheduler != nil {
		cronScheduler.Stop()
	}
}
