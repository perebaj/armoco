package api

import (
	"embed"
	"net/http"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
	"golang.org/x/exp/slog"
)

func handlePostImage(w http.ResponseWriter, r *http.Request) {
	slog.Info("handlePostImage")

}

func handleGetImage(w http.ResponseWriter, r *http.Request) {
	slog.Info("handleGetImage")

}

//go:embed openapi.yaml
var swaggerFs embed.FS

func Routes() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/v1/images", handlePostImage).Methods("POST")
	router.HandleFunc("/v1/images", handleGetImage).Methods("GET")

	opts := middleware.SwaggerUIOpts{SpecURL: "openapi.yaml"}
	sh := middleware.SwaggerUI(opts, nil)
	router.Handle("/docs", sh)
	router.Handle("/openapi.yaml", http.FileServer(http.FS(swaggerFs)))
	return router

}
