package server

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func Handle(w http.ResponseWriter, r *http.Request) {
	t := time.Now()
	log.Printf("%v  %v use time %v content-length %v",
		r.Method,
		r.URL.String(),
		time.Now().Sub(t).String(),
		r.ContentLength)
	fmt.Fprintf(w, "%v\n", "200 Ok")
}
