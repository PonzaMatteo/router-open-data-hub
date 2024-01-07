package app

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/PonzaMatteo/router-open-data-hub/pkg/router"
)

func Start() error {
	log.Println("--- Starting Open Data Hub Router Application ---")

	defaultRouter, err := router.NewDefaultRouter()
	if err != nil {
		return fmt.Errorf("failed to create router: %w", err)
	}
	http.HandleFunc("/", handleRequest(defaultRouter))

	log.Println("Starting application at: http://localhost:8080")
	http.ListenAndServe(":8080", nil)
	if err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func handleRequest(defaultRouter router.Router) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		response, err := defaultRouter.RouteRequest(r.URL.Path+ "?"+r.URL.RawQuery, r.Method)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			var errorResponse struct {
				ErrorMessage string
			}
			errorResponse.ErrorMessage = err.Error()
			json.NewEncoder(w).Encode(errorResponse)
			return
		}
		w.WriteHeader(response.StatusCode)
		fmt.Fprintln(w, response.Body)
	}
}
