package v1

import "net/http"

func GetAdmin(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
