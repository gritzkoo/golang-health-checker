package healthcheck

import (
	"context"
)

// CustomCheckerHandler type to standardize the NewCustomChecker func parameter callback
// The callback is the raw function that will be executed using the Check function of IntegrationItem interface
type CustomCheckerHandler func(ctx context.Context) error

type customChecker struct {
	Name       string
	Handler    CustomCheckerHandler
	IsOptional bool
}

func (c *customChecker) Check(ctx context.Context) (Integration, error) {
	var response = Integration{
		Name:       c.Name,
		Kind:       "Custom Checker",
		Status:     true,
		IsOptional: c.IsOptional,
	}
	if err := c.Handler(ctx); err != nil {
		return response, err
	}

	return response, nil
}

// NewCustomChecker return a customized configuration of IntegrationItem interface with default Check interface
func NewCustomChecker(name string, isOptional bool, callback CustomCheckerHandler) IntegrationItem {
	return &customChecker{
		Name:       name,
		Handler:    callback,
		IsOptional: isOptional,
	}
}
