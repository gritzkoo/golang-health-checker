package healthcheck

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type httpChecker struct {
	params HttpCheckerParams
}

func (h *httpChecker) Check(ctx context.Context) (Integration, error) {
	var response = Integration{
		Name:       h.params.Name,
		Kind:       "Http",
		Url:        h.params.Url,
		Status:     true,
		IsOptional: h.params.IsOptional,
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, h.params.Url, nil)
	if err != nil {
		return response, err
	}
	if h.params.Timeout == 0 {
		h.params.Timeout = time.Second * 10
	}
	if len(h.params.Headers) > 0 {
		for key, val := range h.params.Headers {
			req.Header.Add(key, val)
		}
	}
	cli := &http.Client{
		Timeout: h.params.Timeout,
	}
	r, err := cli.Do(req)
	if err != nil {
		return response, err
	}
	defer r.Body.Close()
	if r.StatusCode != 200 {
		return response, fmt.Errorf("unexpected status code, expect: 200 got: %d", r.StatusCode)
	}
	return response, nil
}

// HttpCheckerParams
type HttpCheckerParams struct {
	// Name to be displaied on readiness response integration name
	Name string
	// The Url to be used on http.get function
	Url string
	// Headers is a optional parameter
	Headers map[string]string
	// Timeout is a optional parameter, if not informed a default 10 seconds will be used
	Timeout time.Duration
	// IsOptional flag is used to mark dependenci optional so when it fails, do not changes the root status to false
	IsOptional bool
}

// NewHttpChecker func return a new IntegrationItem with http configurations
func NewHttpChecker(params HttpCheckerParams) IntegrationItem {
	return &httpChecker{params: params}
}
