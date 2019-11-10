package api

import "net/http"

func Handle(c coordinator) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// receive urls from request

		if err := c.CaptureAsync(""); err != nil {
			// return error to user
			return
		}

		// return success to user
	}
}
