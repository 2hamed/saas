package api

import "net/http"

func Handle(qm QManager) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// receive urls from request

		if err := qm.Enqueue(""); err != nil {
			// return error to user
			return
		}

		// return success to user
	}
}
