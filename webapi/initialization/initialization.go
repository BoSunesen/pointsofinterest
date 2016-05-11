package initialization

import (
	"github.com/BoSunesen/pointsofinterest/webapi/handlers"
	"net/http"
)

func Initialize() {
	const pingUrl string = "/ping/"
	http.Handle(pingUrl, handlers.LoggingDecorator{handlers.PingHandler{}, "Ping"})
	http.Handle("/poi/", handlers.LoggingDecorator{handlers.PoiHandler{}, "POI"})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, pingUrl, http.StatusFound)
	})
}
