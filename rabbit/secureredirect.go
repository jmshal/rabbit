package rabbit

import (
	"fmt"
	"net/http"
)

func (a *rabbit) SecureRedirect(w http.ResponseWriter, r *http.Request) {
	info := GetRequestInfo(r)
	url := info.URL

	if info.Websocket {
		url.Scheme = "wss" // TODO check if websockets allow http redirects
	} else {
		url.Scheme = "https"
	}

	if a.config.Ports.HTTPS != 443 {
		// If the port is non-standard, update the host (usually not production)
		url.Host = fmt.Sprintf("%v:%v", url.Hostname(), a.config.Ports.HTTPS)
	}

	http.Redirect(w, r, url.String(), http.StatusTemporaryRedirect)
}
