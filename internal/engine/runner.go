package engine

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"api-workbench/internal/model"
	"api-workbench/internal/repository"
	"api-workbench/internal/variable"
)

type RunOptions struct {
	EnvID uint
}

func RunTestCase(testCaseID uint, opts *RunOptions) (*model.TestRun, error) {
	tc := &model.TestCase{}
	err := repository.GetTestCaseByID(testCaseID, tc)
	if err != nil {
		return nil, err
	}

	run := &model.TestRun{
		TargetType:    "test_case",
		TargetID:      testCaseID,
		EnvironmentID: opts.EnvID,
		Status:        "running",
		TriggerType:   "manual",
		StartedAt:     time.Now(),
	}
	if err := repository.CreateTestRun(run); err != nil {
		return nil, err
	}

	go executeTestRun(run, tc.ID, nil)

	return run, nil
}

func RunTestSuite(suiteID uint, opts *RunOptions) (*model.TestRun, error) {
	ts := &model.TestSuite{}
	err := repository.GetTestSuiteByID(suiteID, ts)
	if err != nil {
		return nil, err
	}

	run := &model.TestRun{
		TargetType:    "test_suite",
		TargetID:      suiteID,
		EnvironmentID: opts.EnvID,
		Status:        "running",
		TriggerType:   "manual",
		StartedAt:     time.Now(),
	}
	if err := repository.CreateTestRun(run); err != nil {
		return nil, err
	}

	go executeTestSuiteRun(run, ts)

	return run, nil
}

func executeTestRun(run *model.TestRun, testCaseID uint, envVars map[string]string) {
	var tcaList []model.TestCaseAPI
	repository.GetTestCaseAPIs(testCaseID, &tcaList)

	run.Total = len(tcaList)
	repository.UpdateTestRun(run)

	passed := 0
	failed := 0

	for _, tca := range tcaList {
		var apiDef model.API
		if err := repository.GetAPIByID(tca.APIID, &apiDef); err != nil {
			failed++
			continue
		}

		var datasets []model.TestDataSet
		repository.GetTestDataSets(tca.ID, &datasets)

		if len(datasets) > 0 {
			for i, ds := range datasets {
				var data map[string]string
				_ = json.Unmarshal([]byte(ds.Data), &data)
				detail := executeAPI(&apiDef, run.ID, testCaseID, i, data)
				if detail.Status == "passed" {
					passed++
				} else {
					failed++
				}
			}
		} else {
			detail := executeAPI(&apiDef, run.ID, testCaseID, 0, envVars)
			if detail.Status == "passed" {
				passed++
			} else {
				failed++
			}
		}
	}

	now := time.Now()
	run.FinishedAt = &now
	run.Passed = passed
	run.Failed = failed
	run.DurationMs = now.Sub(run.StartedAt).Milliseconds()
	if failed > 0 {
		run.Status = "failed"
	} else {
		run.Status = "done"
	}
	repository.UpdateTestRun(run)
}

func executeTestSuiteRun(run *model.TestRun, suite *model.TestSuite) {
	var suiteCases []model.TestSuiteCase
	repository.GetTestSuiteCases(suite.ID, &suiteCases)

	totalTests := 0
	for _, sc := range suiteCases {
		var tcaList []model.TestCaseAPI
		repository.GetTestCaseAPIs(sc.TestCaseID, &tcaList)
		totalTests += len(tcaList)
	}

	run.Total = totalTests
	repository.UpdateTestRun(run)

	passed := 0
	failed := 0

	if suite.RunMode == "parallel" {
		var mu sync.Mutex
		var wg sync.WaitGroup
		sem := make(chan struct{}, suite.MaxConcurrency)

		for _, sc := range suiteCases {
			wg.Add(1)
			sem <- struct{}{}
			go func(sc model.TestSuiteCase) {
				defer wg.Done()
				defer func() { <-sem }()
				var tcaList []model.TestCaseAPI
				repository.GetTestCaseAPIs(sc.TestCaseID, &tcaList)
				for _, tca := range tcaList {
					var apiDef model.API
					if err := repository.GetAPIByID(tca.APIID, &apiDef); err != nil {
						mu.Lock()
						failed++
						mu.Unlock()
						continue
					}
					detail := executeAPI(&apiDef, run.ID, sc.TestCaseID, 0, nil)
					mu.Lock()
					if detail.Status == "passed" {
						passed++
					} else {
						failed++
					}
					mu.Unlock()
				}
			}(sc)
		}
		wg.Wait()
	} else {
		for _, sc := range suiteCases {
			var tcaList []model.TestCaseAPI
			repository.GetTestCaseAPIs(sc.TestCaseID, &tcaList)
			for _, tca := range tcaList {
				var apiDef model.API
				if err := repository.GetAPIByID(tca.APIID, &apiDef); err != nil {
					failed++
					continue
				}
				detail := executeAPI(&apiDef, run.ID, sc.TestCaseID, 0, nil)
				if detail.Status == "passed" {
					passed++
				} else {
					failed++
				}
			}
		}
	}

	now := time.Now()
	run.FinishedAt = &now
	run.Passed = passed
	run.Failed = failed
	run.DurationMs = now.Sub(run.StartedAt).Milliseconds()
	if failed > 0 {
		run.Status = "failed"
	} else {
		run.Status = "done"
	}
	repository.UpdateTestRun(run)
}

func executeAPI(apiDef *model.API, runID uint, testCaseID uint, dataIndex int, envVars map[string]string) *model.TestRunDetail {
	headers := make(map[string]string)
	if apiDef.Headers != "" {
		_ = json.Unmarshal([]byte(apiDef.Headers), &headers)
	}
	queryParams := make(map[string]string)
	if apiDef.QueryParams != "" {
		_ = json.Unmarshal([]byte(apiDef.QueryParams), &queryParams)
	}

	url := variable.Replace(apiDef.URL, envVars, nil)
	body := variable.Replace(apiDef.Body, envVars, nil)
	headers = variable.ReplaceMap(headers, envVars, nil)
	queryParams = variable.ReplaceMap(queryParams, envVars, nil)

	req := &DebugRequest{
		Method:      apiDef.Method,
		URL:         url,
		Headers:     headers,
		QueryParams: queryParams,
		BodyType:    apiDef.BodyType,
		Body:        body,
		TimeoutMs:   apiDef.TimeoutMs,
	}

	resp := ExecuteHTTP(req)

	detail := &model.TestRunDetail{
		TestRunID:       runID,
		APIID:           apiDef.ID,
		TestCaseID:      testCaseID,
		DataIndex:       dataIndex,
		StatusCode:      resp.StatusCode,
		DurationMs:      resp.DurationMs,
		ExecutedAt:      time.Now(),
	}

	respHeadersJSON, _ := json.Marshal(resp.Headers)
	detail.ResponseHeaders = string(respHeadersJSON)
	detail.ResponseBody = resp.Body

	if resp.Error != "" {
		detail.Status = "failed"
		detail.ErrorMessage = resp.Error
	} else {
		detail.Status = "passed"
		assertions := getAssertionsForAPI(apiDef.ID)
		for _, assertion := range assertions {
			if !checkAssertion(assertion, resp) {
				detail.Status = "failed"
				detail.ErrorMessage = "断言失败: " + assertion.TargetType + " " + assertion.Expected
				break
			}
		}
	}

	repository.CreateTestRunDetail(detail)
	return detail
}

func getAssertionsForAPI(apiID uint) []model.Assertion {
	var list []model.Assertion
	repository.GetAssertionsByAPIID(apiID, &list)
	return list
}

func checkAssertion(a model.Assertion, resp *DebugResponse) bool {
	if !a.Enabled {
		return true
	}
	switch a.TargetType {
	case "status_code":
		return checkStatusAssertion(a, resp.StatusCode)
	case "response_time":
		return checkTimeAssertion(a, resp.DurationMs)
	case "response_body":
		return checkBodyAssertion(a, resp.Body)
	default:
		return true
	}
}

func checkStatusAssertion(a model.Assertion, actual int) bool {
	expected := 0
	if a.Expected != "" {
		if err := json.Unmarshal([]byte(a.Expected), &expected); err != nil {
			if n, err := strconv.Atoi(a.Expected); err == nil {
				expected = n
			}
		}
	}
	switch a.Operator {
	case "equals":
		return actual == expected
	case "not_equals":
		return actual != expected
	default:
		return actual == expected
	}
}

func checkTimeAssertion(a model.Assertion, actual int64) bool {
	expected := int64(0)
	if a.Expected != "" {
		if err := json.Unmarshal([]byte(a.Expected), &expected); err != nil {
			if n, err := strconv.ParseInt(a.Expected, 10, 64); err == nil {
				expected = n
			}
		}
	}
	switch a.Operator {
	case "less_than":
		return actual < expected
	case "greater_than":
		return actual > expected
	case "equals":
		return actual == expected
	default:
		return true
	}
}

func checkBodyAssertion(a model.Assertion, body string) bool {
	switch a.Operator {
	case "contains":
		return len(body) > 0 && len(a.Expected) > 0 && strings.Contains(body, a.Expected)
	case "not_contains":
		return !(len(body) > 0 && len(a.Expected) > 0 && strings.Contains(body, a.Expected))
	case "equals":
		return body == a.Expected
	default:
		return true
	}
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
