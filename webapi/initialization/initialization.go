package initialization

import (
	"github.com/BoSunesen/pointsofinterest/webapi/handlers"
	"html"
	"log"
	"net/http"
)

func Initialize() {
	//TODO Extract method
	const poiUrl string = "/poi/"
	poiHandler := handlers.NewPoiHandler()
	http.Handle(poiUrl, &handlers.LoggingDecorator{poiHandler, poiUrl})

	const pingUrl string = "/ping/"
	http.Handle(pingUrl, &handlers.LoggingDecorator{handlers.PingHandler{}, pingUrl})
	//TODO Better routing
	//TODO favicon.ico
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := html.EscapeString(r.URL.Path)
		log.Printf("Redirecting %q", path)
		defer log.Printf("Redirected %q", path)
		http.Redirect(w, r, pingUrl, http.StatusFound)
	})
}
