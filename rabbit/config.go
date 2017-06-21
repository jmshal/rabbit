package rabbit

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"

	bugsnag "github.com/bugsnag/bugsnag-go"
)

type config struct {
	Routes  []route  `json:"routes"`
	Ports   ports    `json:"ports"`
	Bugsnag bugsnag_ `json:"bugsnag"`
	Certs   []certs  `json:"certs"`
}

type certs struct {
	Cert string `json:"cert"`
	Key  string `json:"key"`
}

type bugsnag_ struct {
	APIKey       string `json:"apiKey"`
	ReleaseStage string `json:"releaseStage"`
	Endpoint     string `json:"endpoint"`
}

type route struct {
	Entrypoints []entrypoint `json:"entrypoints"`
	Endpoints   []endpoint   `json:"endpoints"`
}

type entrypoint struct {
	Host string `json:"host"`
	Port uint16 `json:"port"`
	Path string `json:"path"`
}

type endpoint struct {
	Protocol string `json:"protocol"`
	Host     string `json:"host"`
	Port     uint16 `json:"port"`
	Path     string `json:"path"`
	Origin   string `json:"origin"`
}

type ports struct {
	HTTP  uint16 `json:"http"`
	HTTPS uint16 `json:"https"`
}

func (c *config) applyDefaults() {
	if c.Ports.HTTP == 0 {
		port := 80
		for _, name := range []string{
			"PORT",
			"HTTP_PLATFORM_PORT",
		} {
			if env := os.Getenv(name); env != "" {
				port, _ = strconv.Atoi(env)
				break
			}
		}
		c.Ports.HTTP = uint16(port)
	}
	if c.Ports.HTTPS == 0 && len(c.Certs) > 0 {
		c.Ports.HTTPS = 443
	}
}

func LoadConfigString(text string) (*config, error) {
	var config config
	config.applyDefaults()

	err := json.Unmarshal([]byte(text), &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func LoadConfigFile(path string) (*config, error) {
	file, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, err
	}

	return LoadConfigString(string(file))
}

func (c *config) configureBugsnag() {
	bugsnag.Configure(bugsnag.Configuration{
		APIKey:       c.Bugsnag.APIKey,
		ReleaseStage: c.Bugsnag.ReleaseStage,
		Endpoint:     c.Bugsnag.Endpoint,
		AppVersion:   Version,
	})
}
