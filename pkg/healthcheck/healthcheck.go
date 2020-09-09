package healthcheck

// ApplicationHealth is a type for map the result
type ApplicationHealth struct {
	Name    string  `json:"name"`
	Status  bool    `json:"status"`
	Version string  `json:"version"`
	Date    string  `json:"date"`
	Checks  []Check `json:"checks"`
}

// Check is the type result for requests
type Check struct {
	Name         string `json:"name"`
	Kind         string `json:"kind"`
	Status       bool   `json:"status"`
	ResponseTime int64  `json:"response_time"`
	Optional     bool   `json:"optional"`
	URL          string `json:"url"`
}
