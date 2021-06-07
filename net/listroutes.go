package net

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/gorilla/mux"
	log "github.com/wlbr/commons/log"
)

var router *mux.Router

func SetRouter(r *mux.Router) {
	router = r
}

func GetRouter() *mux.Router {
	return router
}

func ListRoutes(router *mux.Router) []string {
	var ollist []string
	router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		path, err := route.GetPathTemplate()
		log.Info("Route: %s \t error: %v", path, err)
		if err == nil {
			ollist = append(ollist, path)
		}
		return err
	})
	sort.Strings(ollist)
	return ollist
}

func ListRoutesToResponse(w http.ResponseWriter, r *http.Request) {
	log.Info("Receiving request: %s %s", r.Method, r.URL)

	fmt.Fprintf(w, "<html>\n<body>\n\n<h1>All Routes</h1>\n\n <ol>")

	for _, route := range ListRoutes(GetRouter()) {
		fmt.Fprintf(w, "\n   <li><a href=\"%s\">%s</li>", route, route)
	}
	fmt.Fprintf(w, "\n </ol>\n\n</body>\n</html>")
}
