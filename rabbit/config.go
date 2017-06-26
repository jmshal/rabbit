package rabbit

import (
	"encoding/json"
	"io/ioutil"
)

type config struct {
	Debug    bool      `json:"debug"`
	Ports    *Ports    `json:"ports"`
	Certs    []*Cert   `json:"certs"`
	Routes   []*Route  `json:"routes"`
	Logging  *Logging  `json:"logging"`
	Database *Database `json:"database"`
}

type Ports struct {
	HTTP  int `json:"http"`
	HTTPS int `json:"https"`
}

type Cert struct {
	Cert string `json:"cert"`
	Key  string `json:"key"`
}

type Route struct {
	Entrypoints []*Entrypoint  `json:"entrypoints"`
	Endpoints   []*Endpoint    `json:"endpoints"`
	Caches      []*CachePolicy `json:"caches"`
}

type Entrypoint struct {
	Secure     bool     `json:"secure"`     // upgrades non-https requests to https
	Websockets bool     `json:"websockets"` // default: disable websockets
	Methods    []string `json:"methods"`
	Host       string   `json:"host"`
	Path       string   `json:"path"`
}

type Endpoint struct {
	Scheme string `json:"scheme"` // default http
	Host   string `json:"host"`
	Port   int    `json:"port"`
	Path   string `json:"path"`
	Origin string `json:"origin"`
}

type Logging struct {
	ApplicationInsights *ApplicationInsights `json:"applicationInsights"`
	Bugsnag             *Bugsnag             `json:"bugsnag"`
}

type ApplicationInsights struct {
	InstrumentationKey string `json:"instrumentationKey"`
}

type Bugsnag struct {
	Key          string `json:"apiKey"`
	ReleaseStage string `json:"releaseStage"`
	Endpoint     string `json:"endpoint"`
}

type CachePolicy struct {
	Match  string `json:"match"` // glob matcher string (full url)
	TTL    int    `json:"ttl"`
	Search bool   `json:"search"` // include search string in cache id
}

type Database struct {
	Redis *Redis `json:"redis"`
}

type Redis struct {
	URL string `json:"url"`
}

func (c *config) ApplyDefaults() {
	c.Debug = false
	c.Ports = &Ports{
		HTTP: 80,
	}
	c.Certs = make([]*Cert, 0)
	c.Routes = make([]*Route, 0)
	c.Logging = &Logging{}
	c.Database = &Database{}
}

func NewConfig() *config {
	c := &config{}
	c.ApplyDefaults()
	return c
}

func NewConfigFromFile(path string) (*config, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	c := NewConfig()
	err = json.Unmarshal(b, c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

/*

{
	"database": {
		"redis": {
			"url": "tcp://redis:6379"
		}
	},
	"logging": {
		"appInsights": {
			"instrumentationKey": "..."
		},
		"bugsnag": {
			"apiKey": "...",
			"releaseStage": "production"
		}
	},
	"ports": {
		"http": 80,
		"https": 443
	},
	"certs": [{
		"cert": "/run/secrets/example.com.crt",
		"key": "/run/secrets/example.com.key"
	}],
	"routes": [
		{
			"entrypoints": [{
				"secure": true,
				"host": "example.com",
				"path": "/blog"
			}],
			"endpoints": [{
				"scheme": "https",
				"host": "wordpress", // assume docker overlay network, and "wordpress" is the hostname
				"path": "/"
			}],
			"caches": [{
				"match": "*\/wp-content/*", // whole url matching
				"ttl": 3600, // 1 hour ttl
				"search": false
			}]
		},
		{
			"entrypoints": [{
				"secure": true,
				"host": "example.com",
				"path": "/"
			}],
			"endpoints": [{
				"scheme": "https",
				"host": "homepage", // same as above
				"path": "/"
			}]
		}
	]
}

*/
