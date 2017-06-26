package rabbit

import (
	"net/http"
	"time"
)

func (a *rabbit) AfterRequest(r *http.Request) {
	info := GetRequestInfo(r)
	info.Duration = time.Since(info.StartTime)

	if a.config.Logging.ApplicationInsights != nil {
		a.logRequestApplicationInsights(r)
	}
}
