package rabbit

import (
	"net/http"
	"strconv"

	"fmt"

	"github.com/jmshal/go-applicationinsights/appinsights"
)

func (a *rabbit) logRequestApplicationInsights(r *http.Request) {
	info := GetRequestInfo(r)

	telemetry := appinsights.NewRequestTelemetry(
		info.ID,
		fmt.Sprintf("%v %v", r.Method, info.URL.String()),
		info.StartTime,
		info.Duration,
		r.Method,
		info.URL.String(),
		strconv.Itoa(info.ResponseStatus),
		info.ResponseStatus != 0 && info.ResponseStatus >= 200 && info.ResponseStatus < 400,
		nil,
		func(c *appinsights.TelemetryContext) {})

	a.applicationInsights.Track(telemetry)
}
