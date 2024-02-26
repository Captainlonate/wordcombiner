package routes

import (
	ce "captainlonate/wordcombiner/internal/customError"
	"net/http"
)

// Route handler middleware for net/http that will ensure that all
// required query parameters are present in the request.
//
// Example:
//
//	mux.Handle("GET /combine", validateQueryParams(
//		http.HandlerFunc(combineRouteHandler), []string{"one", "two"}),
//	)
func validateQueryParams(next http.Handler, queryParams []string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if all query parameters in queryParms were provided
		for _, param := range queryParams {
			if r.URL.Query().Get(param) == "" {
				sendJSON(w, apiResponseFailure(ce.BadRequestQueryParam, "'"+param+"' query param is required"))
				return
			}
		}

		// Call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}
