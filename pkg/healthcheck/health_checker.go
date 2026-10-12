package healthcheck

import (
	"context"
	"sync"
	"time"
)

const internalStatus = "fully functional"

// Liveness performs a simple check of the application
func Liveness() LivenessResponse {
	return LivenessResponse{
		Status: internalStatus,
	}
}

// Readiness perform a check for every integration informed
func Readiness(ctx context.Context, config RedinessConfig) ReadinessResponse {
	var (
		start     = time.Now()
		wg        sync.WaitGroup
		total     = len(config.Integrations)
		checklist = make(chan Integration, total)
		result    = ReadinessResponse{
			Name:    config.Name,
			Version: config.Version,
			Status:  true,
			Date:    start.Format(time.RFC3339),
		}
	)
	wg.Add(total)
	for idx, integration := range config.Integrations {
		go func(idx int) {
			internalValidation(ctx, integration, &wg, checklist)
		}(idx)
	}
	wg.Wait()
	close(checklist)
	result.Duration = time.Since(start).Seconds()
	for chk := range checklist {
		result.Integrations = append(result.Integrations, chk)
		if !chk.Status && !chk.IsOptional {
			result.Status = false
		}
	}
	return result
}

func internalValidation(ctx context.Context, item IntegrationItem, wg *sync.WaitGroup, checklist chan Integration) {
	defer wg.Done()
	var start = time.Now()
	response, err := item.Check(ctx)
	response.ResponseTime = time.Since(start).Seconds()
	if err != nil {
		response.Status = false
		response.Error = err.Error()
	}
	checklist <- response
}
