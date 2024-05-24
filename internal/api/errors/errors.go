package errors

import (
	"log/slog"
	"net/http"
)

func InternalServerError(w http.ResponseWriter, _ *http.Request, ierr error) {
	w.WriteHeader(http.StatusInternalServerError)
	_, err := w.Write([]byte("internal server error"))
	if err != nil {
		slog.Error("Failed to write response body")
	}
	if ierr != nil {
		slog.Error("", ierr)
	}
}

func Unauthorized(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusUnauthorized)
	_, err := w.Write([]byte("unauthorized"))
	if err != nil {
		InternalServerError(w, nil, nil)
		return
	}
}

func Forbidden(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusForbidden)
	_, err := w.Write([]byte("forbidden"))
	if err != nil {
		InternalServerError(w, nil, nil)
		return
	}
}

func NotFound(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	_, err := w.Write([]byte("not found"))
	if err != nil {
		InternalServerError(w, nil, err)
		return
	}
}

func MethodNotAllowed(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	_, err := w.Write([]byte("method is not allowed"))
	if err != nil {
		slog.Error("Failed to write response body")
	}
}
