package api

import (
	"fmt"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
)

func NewJobHandler(c dispatcher) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// receive urls from request
		r.ParseForm()

		urls := r.FormValue("urls")
		urlsSlice := strings.Split(urls, ";")

		for _, u := range urlsSlice {
			if err := c.Enqueue(u); err != nil {
				log.Debug("failed to queue", u, err)
				w.WriteHeader(500)
				w.Write([]byte("failed queueing some urls"))
				return
			}
		}

		w.WriteHeader(200)
		w.Write([]byte("success"))
	}
}

func GetResultHandler(d dispatcher) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		path, err := d.GetResult("")

		if err != nil {
			// show error message
		}

		fmt.Println(path)
	}
}
