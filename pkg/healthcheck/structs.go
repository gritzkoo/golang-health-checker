package healthcheck

/*
	ApplicationHealth used to check all application integrations
	and return status of each of then
*/
type ApplicationHealthDetailed struct {
	Name         string        `json:"name,omitempty"`
	Status       bool          `json:"status"`
	Version      string        `json:"version,omitempty"`
	Date         string        `json:"date"`
	Duration     float64       `json:duration"`
	Integrations []Integration `json:"integrations"`
}

// ApplicationHealthSimple used to simple return a string 'OK'
type ApplicationHealthSimple struct {
	Status string `json:"status"`
}

// Integration is the type result for requests
type Integration struct {
	Name         string  `json:"name"`
	Kind         string  `json:"kind"`
	Status       bool    `json:"status"`
	ResponseTime float64 `json:"response_time"`
	URL          string  `json:"url"`
	Error        error   `json:"errors,omitempty"`
}

// Auth is a default struct to map user/pass protocol
type Auth struct {
	User     string
	Password string
}

// ApplicationConfig is a config contract to init health caller
type ApplicationConfig struct {
	Name         string              `json:"name"`
	Version      string              `json:"version"`
	Integrations []IntegrationConfig `json:"integrations"`
}

// IntegrationConfig used to inform each integration config
type IntegrationConfig struct {
	Type    string       `json:"Type"` // must be web | redis | memcache
	Name    string       `json:"name"`
	Host    string       `json:"host"` // yes you can concat host:port here
	Port    string       `json:"port,omitempty"`
	Headers []HTTPHeader `json:"headers,omitempty"`
	DB      int          `json:"db,omitempty"`      // default value is 0
	TimeOut int          `json:"timeout,omitempty"` // default value: 10
	Auth    Auth         `json:"auth,omitempty"`
}

// HTTPHeader used to setup webservices integrations
type HTTPHeader struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"Value,omitempty"`
}

// Mapped types for IntegrationConfig
const (
	Redis     = "redis"
	Memcached = "memcached"
	Web       = "web"
)

// Mapped typs for kinds of integrations
const (
	RedisIntegration      = "Redis DB"
	MemcachedIntegration  = "Memcached DB"
	WebServiceIntegration = "Web service API"
)
