package engine

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"api-workbench/internal/model"
	"api-workbench/internal/repository"
	"api-workbench/internal/variable"

	"github.com/gorilla/websocket"
)

var wsClients = make(map[uint][]*websocket.Conn)
var wsMu sync.RWMutex

type WSMessage struct {
	Type    string      `json:"type"`
	RunID   uint        `json:"run_id"`
	Data    interface{} `json:"data"`
}

func RegisterWSClient(runID uint, conn *websocket.Conn) {
	wsMu.Lock()
	defer wsMu.Unlock()
	wsClients[runID] = append(wsClients[runID], conn)
}

func UnregisterWSClient(runID uint, conn *websocket.Conn) {
	wsMu.Lock()
	defer wsMu.Unlock()
	clients := wsClients[runID]
	for i, c := range clients {
		if c == conn {
			wsClients[runID] = append(clients[:i], clients[i+1:]...)
			break
		}
	}
}

func broadcastToRun(runID uint, msg WSMessage) {
	wsMu.RLock()
	clients := make([]*websocket.Conn, len(wsClients[runID]))
	copy(clients, wsClients[runID])
	wsMu.RUnlock()

	data, _ := json.Marshal(msg)
	for _, conn := range clients {
		if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
			log.Printf("WebSocket write error: %v", err)
		}
	}
}

func RunTestCaseWithWS(testCaseID uint, opts *RunOptions) (*model.TestRun, error) {
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

	go func() {
		executeTestRunWithWS(run, tc.ID, nil)
	}()

	return run, nil
}

func RunTestSuiteWithWS(suiteID uint, opts *RunOptions) (*model.TestRun, error) {
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

	go func() {
		executeTestSuiteRunWithWS(run, ts)
	}()

	return run, nil
}

func executeTestRunWithWS(run *model.TestRun, testCaseID uint, envVars map[string]string) {
	var tcaList []model.TestCaseAPI
	repository.GetTestCaseAPIs(testCaseID, &tcaList)

	run.Total = len(tcaList)
	repository.UpdateTestRun(run)

	broadcastToRun(run.ID, WSMessage{
		Type:  "progress",
		RunID: run.ID,
		Data: map[string]interface{}{
			"total":  run.Total,
			"passed": 0,
			"failed": 0,
			"status": "running",
		},
	})

	passed := 0
	failed := 0

	for i, tca := range tcaList {
		var apiDef model.API
		if err := repository.GetAPIByID(tca.APIID, &apiDef); err != nil {
			failed++
			continue
		}

		detail := executeAPI(&apiDef, run.ID, testCaseID, 0, envVars)
		if detail.Status == "passed" {
			passed++
		} else {
			failed++
		}

		broadcastToRun(run.ID, WSMessage{
			Type:  "detail",
			RunID: run.ID,
			Data: map[string]interface{}{
				"index":    i + 1,
				"total":    run.Total,
				"passed":   passed,
				"failed":   failed,
				"api_name": apiDef.Name,
				"status":   detail.Status,
			},
		})
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

	broadcastToRun(run.ID, WSMessage{
		Type:  "complete",
		RunID: run.ID,
		Data: map[string]interface{}{
			"status":     run.Status,
			"passed":     run.Passed,
			"failed":     run.Failed,
			"duration_ms": run.DurationMs,
		},
	})
}

func executeTestSuiteRunWithWS(run *model.TestRun, suite *model.TestSuite) {
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

	broadcastToRun(run.ID, WSMessage{
		Type:  "progress",
		RunID: run.ID,
		Data: map[string]interface{}{
			"total":  run.Total,
			"passed": 0,
			"failed": 0,
			"status": "running",
		},
	})

	passed := 0
	failed := 0
	currentIndex := 0

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
						currentIndex++
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
					currentIndex++
					idx := currentIndex
					p := passed
					f := failed
					mu.Unlock()

					broadcastToRun(run.ID, WSMessage{
						Type:  "detail",
						RunID: run.ID,
						Data: map[string]interface{}{
							"index":    idx,
							"total":    run.Total,
							"passed":   p,
							"failed":   f,
							"api_name": apiDef.Name,
							"status":   detail.Status,
						},
					})
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
					currentIndex++
					continue
				}
				detail := executeAPI(&apiDef, run.ID, sc.TestCaseID, 0, nil)
				if detail.Status == "passed" {
					passed++
				} else {
					failed++
				}
				currentIndex++

				broadcastToRun(run.ID, WSMessage{
					Type:  "detail",
					RunID: run.ID,
					Data: map[string]interface{}{
						"index":    currentIndex,
						"total":    run.Total,
						"passed":   passed,
						"failed":   failed,
						"api_name": apiDef.Name,
						"status":   detail.Status,
					},
				})
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

	broadcastToRun(run.ID, WSMessage{
		Type:  "complete",
		RunID: run.ID,
		Data: map[string]interface{}{
			"status":     run.Status,
			"passed":     run.Passed,
			"failed":     run.Failed,
			"duration_ms": run.DurationMs,
		},
	})
}

func ReplaceVars(text string, envVars map[string]string) string {
	return variable.Replace(text, envVars, nil)
}
