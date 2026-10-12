package healthcheck

import "context"

// ReadinessResponse used to check all application integrations and return status of each of then
type ReadinessResponse struct {
	Name         string        `json:"name,omitempty"`
	Status       bool          `json:"status"`
	Version      string        `json:"version,omitempty"`
	Date         string        `json:"date"`
	Duration     float64       `json:"duration"`
	Integrations []Integration `json:"integrations"`
}

// LivenessResponse used to simple return a string 'OK'
type LivenessResponse struct {
	Status string `json:"status"`
}

// Integration is the type result for requests
type Integration struct {
	Name         string  `json:"name"`
	Kind         string  `json:"kind"`
	Status       bool    `json:"status"`
	IsOptional   bool    `json:"is_optional"`
	ResponseTime float64 `json:"response_time"` //in seconds
	Url          string  `json:"url,omitempty"`
	Error        string  `json:"error,omitempty"` // error.Error()
}

// RedinessConfig is a config contract to init health caller
type RedinessConfig struct {
	// The name of your application
	Name string
	// The version to trace changes
	Version string
	// The list of integrations to be tested
	Integrations []IntegrationItem
}

// IntegrationItem interface defines the contract of each Integration that will be tested
// Now you can create your own integration and use as you wish
type IntegrationItem interface {
	// Check function execute the integraion validation
	Check(ctx context.Context) (Integration, error)
}
