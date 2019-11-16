package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// NewJobHandler is the handler to accept new screenshot jobs from the user
func NewJobHandler(c dispatcher) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// receive urls from request
		r.ParseForm()

		urls := r.FormValue("urls")

		if urls == "" {
			log.Debug("supplied request has empty urls")
			w.WriteHeader(400)
			return
		}

		urlsSlice := strings.Split(urls, ";")
		urlHashes := make([]string, len(urlsSlice))

		for i, url := range urlsSlice {
			if err := c.Enqueue(url); err != nil {
				log.Debugf("failed to queue %s - %v", url, err)
				w.WriteHeader(500)
				w.Write([]byte("failed queueing some urls"))
				return
			}
			urlHashes[i] = getBaseURL(r, "api/result/") + base64.StdEncoding.EncodeToString([]byte(url))
		}

		response := map[string]interface{}{
			"results": urlHashes,
		}
		responseBody, _ := json.Marshal(response)

		w.WriteHeader(200)
		w.Write(responseBody)
	}
}

// GetResultHandler retreives the result of the request and returns proper
// response to the user
func GetResultHandler(d dispatcher) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		urlHash := vars["hash"]

		url, err := base64.StdEncoding.DecodeString(urlHash)

		if err != nil {
			w.WriteHeader(422)
			return
		}

		exists, isPending, isFinished, err := d.FetchStatus(string(url))

		if err != nil {
			log.Errorf("failed checking the status of url: %v", err)
			w.WriteHeader(500)
			return
		}

		if !exists {
			w.WriteHeader(404)
			return
		}

		if isPending {
			w.WriteHeader(204)
			return
		}

		if !isFinished {
			w.WriteHeader(503)
			return
		}

		path, err := d.FetchResult(string(url))

		if err != nil {
			log.Errorf("failed fetching the result of url: %v", err)
			w.WriteHeader(500)
			return
		}

		w.WriteHeader(200)
		w.Write([]byte(getBaseURL(r, path)))
	}
}

func getBaseURL(r *http.Request, tail string) string {
	return fmt.Sprintf("%s://%s/%s", "http", r.Host, tail)
}
