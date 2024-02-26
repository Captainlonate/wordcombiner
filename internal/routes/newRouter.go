package routes

import (
	"net/http"
)

func NewRouter() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /combine", validateQueryParams(
		http.HandlerFunc(combineRouteHandler), []string{"one", "two"}),
	)

	return mux
}
