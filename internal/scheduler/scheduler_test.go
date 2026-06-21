package scheduler

import (
	"testing"
	"time"

	"github.com/robfig/cron/v3"
)

func TestJobManagement(t *testing.T) {
	jobs = make(map[uint]cron.EntryID)

	jobs[1] = cron.EntryID(1)
	jobs[2] = cron.EntryID(2)

	if len(jobs) != 2 {
		t.Errorf("jobs count = %v, want 2", len(jobs))
	}

	delete(jobs, 1)
	if len(jobs) != 1 {
		t.Errorf("jobs count after delete = %v, want 1", len(jobs))
	}

	if _, exists := jobs[1]; exists {
		t.Error("job 1 should be deleted")
	}

	if _, exists := jobs[2]; !exists {
		t.Error("job 2 should exist")
	}
}

func TestTimeFormat(t *testing.T) {
	now := time.Now()
	formatted := now.Format("2006-01-02 15:04:05")

	if len(formatted) != 19 {
		t.Errorf("formatted time length = %v, want 19", len(formatted))
	}
}
