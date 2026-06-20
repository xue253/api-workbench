package model

import "time"

type Project struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"size:100;not null"`
	Description string    `json:"description" gorm:"size:500"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Environment struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	ProjectID   uint   `json:"project_id" gorm:"index"`
	Name        string `json:"name" gorm:"size:100;not null"`
	Description string `json:"description" gorm:"size:500"`
	SortOrder   int    `json:"sort_order" gorm:"default:0"`
}

type EnvironmentVariable struct {
	ID            uint   `json:"id" gorm:"primaryKey"`
	EnvironmentID uint   `json:"environment_id" gorm:"index"`
	Key           string `json:"key" gorm:"size:200;not null"`
	Value         string `json:"value" gorm:"type:text"`
}

type Collection struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	ProjectID   uint   `json:"project_id" gorm:"index"`
	ParentID    *uint  `json:"parent_id" gorm:"index"`
	Name        string `json:"name" gorm:"size:100;not null"`
	Description string `json:"description" gorm:"size:500"`
	SortOrder   int    `json:"sort_order" gorm:"default:0"`
}

type API struct {
	ID             uint      `json:"id" gorm:"primaryKey"`
	CollectionID   uint      `json:"collection_id" gorm:"index"`
	Name           string    `json:"name" gorm:"size:200;not null"`
	Description    string    `json:"description" gorm:"size:500"`
	Protocol       string    `json:"protocol" gorm:"size:20;default:http"`
	Method         string    `json:"method" gorm:"size:10"`
	URL            string    `json:"url" gorm:"size:2000"`
	Headers        string    `json:"headers" gorm:"type:text"`
	PathParams     string    `json:"path_params" gorm:"type:text"`
	QueryParams    string    `json:"query_params" gorm:"type:text"`
	BodyType       string    `json:"body_type" gorm:"size:20"`
	Body           string    `json:"body" gorm:"type:text"`
	ProtoService   string    `json:"proto_service" gorm:"size:200"`
	ProtoMethod    string    `json:"proto_method" gorm:"size:200"`
	ExpectedStatus int       `json:"expected_status"`
	TimeoutMs      int       `json:"timeout_ms" gorm:"default:30000"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type Assertion struct {
	ID         uint   `json:"id" gorm:"primaryKey"`
	APIID      uint   `json:"api_id" gorm:"index"`
	TargetType string `json:"target_type" gorm:"size:30;not null"`
	Operator   string `json:"operator" gorm:"size:20;not null"`
	Path       string `json:"path" gorm:"size:500"`
	Expected   string `json:"expected" gorm:"size:1000"`
	Enabled    bool   `json:"enabled" gorm:"default:true"`
}

type TestCase struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	ProjectID   uint      `json:"project_id" gorm:"index"`
	Name        string    `json:"name" gorm:"size:200;not null"`
	Description string    `json:"description" gorm:"size:500"`
	CreatedAt   time.Time `json:"created_at"`
}

type TestCaseAPI struct {
	ID         uint `json:"id" gorm:"primaryKey"`
	TestCaseID uint `json:"test_case_id" gorm:"index"`
	APIID      uint `json:"api_id" gorm:"index"`
	SortOrder  int  `json:"sort_order" gorm:"default:0"`
}

type TestDataSet struct {
	ID           uint   `json:"id" gorm:"primaryKey"`
	TestCaseAPIID uint  `json:"test_case_api_id" gorm:"index"`
	Data         string `json:"data" gorm:"type:text"`
	SortOrder    int    `json:"sort_order" gorm:"default:0"`
}

type TestSuite struct {
	ID             uint      `json:"id" gorm:"primaryKey"`
	ProjectID      uint      `json:"project_id" gorm:"index"`
	Name           string    `json:"name" gorm:"size:200;not null"`
	Description    string    `json:"description" gorm:"size:500"`
	RunMode        string    `json:"run_mode" gorm:"size:20;default:sequential"`
	MaxConcurrency int       `json:"max_concurrency" gorm:"default:5"`
	CreatedAt      time.Time `json:"created_at"`
}

type TestSuiteCase struct {
	ID          uint `json:"id" gorm:"primaryKey"`
	TestSuiteID uint `json:"test_suite_id" gorm:"index"`
	TestCaseID  uint `json:"test_case_id" gorm:"index"`
	SortOrder   int  `json:"sort_order" gorm:"default:0"`
}

type ScheduledTask struct {
	ID            uint   `json:"id" gorm:"primaryKey"`
	ProjectID     uint   `json:"project_id" gorm:"index"`
	TargetType    string `json:"target_type" gorm:"size:20"`
	TargetID      uint   `json:"target_id"`
	CronExpr      string `json:"cron_expr" gorm:"size:100"`
	Enabled       bool   `json:"enabled" gorm:"default:true"`
	EnvironmentID uint   `json:"environment_id"`
}

type TestRun struct {
	ID            uint       `json:"id" gorm:"primaryKey"`
	TargetType    string     `json:"target_type" gorm:"size:20"`
	TargetID      uint       `json:"target_id"`
	EnvironmentID uint       `json:"environment_id"`
	Status        string     `json:"status" gorm:"size:20;default:running"`
	TriggerType   string     `json:"trigger_type" gorm:"size:20"`
	Total         int        `json:"total"`
	Passed        int        `json:"passed"`
	Failed        int        `json:"failed"`
	Skipped       int        `json:"skipped"`
	DurationMs    int64      `json:"duration_ms"`
	StartedAt     time.Time  `json:"started_at"`
	FinishedAt    *time.Time `json:"finished_at"`
}

type TestRunDetail struct {
	ID              uint       `json:"id" gorm:"primaryKey"`
	TestRunID       uint       `json:"test_run_id" gorm:"index"`
	APIID           uint       `json:"api_id"`
	TestCaseID      uint       `json:"test_case_id"`
	DataIndex       int        `json:"data_index"`
	Status          string     `json:"status" gorm:"size:20"`
	StatusCode      int        `json:"status_code"`
	ResponseHeaders string     `json:"response_headers" gorm:"type:text"`
	ResponseBody    string     `json:"response_body" gorm:"type:text"`
	DurationMs      int64      `json:"duration_ms"`
	ErrorMessage    string     `json:"error_message" gorm:"type:text"`
	RetryCount      int        `json:"retry_count"`
	ExecutedAt      time.Time  `json:"executed_at"`
}
