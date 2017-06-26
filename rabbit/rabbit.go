package rabbit

import (
	"log"
	"net/http"
	"net/http/httputil"

	bugsnag "github.com/bugsnag/bugsnag-go"
	"github.com/jmshal/go-applicationinsights/appinsights"
	"github.com/jmshal/wsutil"
)

type rabbit struct {
	config              *config
	httpProxy           *httputil.ReverseProxy
	wsProxy             *wsutil.ReverseProxy
	applicationInsights appinsights.TelemetryClient
}

func (a *rabbit) Log(format string, v ...interface{}) {
	if a.config.Debug {
		log.Printf(format, v...)
	}
}

func (a *rabbit) setup() {
	if a.config.Logging.Bugsnag != nil {
		c := a.config.Logging.Bugsnag
		bugsnag.Configure(bugsnag.Configuration{
			APIKey:       c.Key,
			ReleaseStage: c.ReleaseStage,
			Endpoint:     c.Endpoint,
			AppVersion:   Version,
		})
	}
	if c := a.config.Logging.ApplicationInsights; c != nil {
		a.applicationInsights = appinsights.NewTelemetryClient(c.InstrumentationKey)
	}
}

func NewRabbit(c *config) *rabbit {
	a := &rabbit{
		config: c,
		wsProxy: &wsutil.ReverseProxy{
			Director: func(r *http.Request) {},
		},
	}
	a.httpProxy = &httputil.ReverseProxy{
		Director:       func(r *http.Request) {},
		ModifyResponse: a.modifyResponse, // strip server identifying headers
	}
	a.setup()
	return a
}
