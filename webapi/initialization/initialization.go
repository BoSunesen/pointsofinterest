package initialization

import (
	"github.com/BoSunesen/pointsofinterest/webapi/handlers"
	"html"
	"log"
	"net/http"
)

func Initialize() {
	handle("/poi/", handlers.NewPoiHandler())

	const pingRoute string = "/ping/"
	handle(pingRoute, handlers.PingHandler{})

	//TODO Better routing
	//TODO favicon.ico
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := html.EscapeString(r.URL.Path)
		log.Printf("Redirecting %v", path)
		defer log.Printf("Redirected %v", path)
		http.Redirect(w, r, pingRoute, http.StatusFound)
	})
}

func handle(route string, handler handlers.HttpHandler) {
	http.Handle(route, &handlers.LoggingDecorator{handler, route})
}
