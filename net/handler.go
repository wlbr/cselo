package net

import (
	"fmt"
	"net/http"
	"time"

	log "github.com/wlbr/commons/log"
)

func Handler(w http.ResponseWriter, r *http.Request) {

	log.Info("Receiving request.")
	log.Info("Request %+v", r)

	keys, ok := r.URL.Query()["id"]

	if ok && len(keys) >= 1 {
		id := keys[0]
		log.Info("Got ip %s from request url.", id)
	} else {

	}
	var err error
	if err != nil {

	}
	fmt.Fprintf(w, "<html>start<br>\n")
	for i := 0; i < 5; i++ {
		fmt.Fprintf(w, "%d<br>\n", i)
		time.Sleep(1 * time.Second)
	}
	fmt.Fprintf(w, "end\n</html")

}
