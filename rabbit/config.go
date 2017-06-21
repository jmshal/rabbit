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
	Server  server   `json:"server"`
	Bugsnag bugsnag_ `json:"bugsnag"`
	TLS     []tls_   `json:"tls"`
}

type tls_ struct {
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

type server struct {
	Port uint16 `json:"port"`
	TLS  bool   `json:"tls"`
}

func (c *config) applyDefaults() {
	if c.Server.Port == 0 {
		var port int
		for _, name := range []string{
			"PORT",
			"HTTP_PLATFORM_PORT",
		} {
			if env := os.Getenv(name); env != "" {
				port, _ = strconv.Atoi(env)
				break
			}
		}
		c.Server.Port = uint16(port)
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
