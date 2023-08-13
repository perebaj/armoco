package api

import (
	"embed"
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
	"os"

	"github.com/cloudflare/cloudflare-go"
	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
	"golang.org/x/exp/slog"
)

type Application struct {
	Cloudflare *cloudflare.API
}

func (app *Application) handlerPostImage(w http.ResponseWriter, r *http.Request) {
	slog.Info("handlePostImage")
	err := r.ParseMultipartForm(32 << 20) // 32MB is the maximum size of a file we can upload
	if err != nil {
		slog.Error("failed to parse multipart form", "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	slog.Info("form", "form", r.Form)
	slog.Info("form", "form name", r.Form["name"])

	file, handler, err := r.FormFile("image")
	if err != nil {
		slog.Error("failed to get image from form", "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer file.Close()

	f, err := os.OpenFile("./uploads/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		slog.Error("failed to open file", "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()

	_, err = io.Copy(f, file)
	if err != nil {
		slog.Error("failed to copy file", "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write([]byte("File uploaded successfully :)"))
}

func (app *Application) handlerGetImage(w http.ResponseWriter, r *http.Request) {
	images, err := app.Cloudflare.ListImages(r.Context(), nil, cloudflare.ListImagesParams{})

	if err != nil {
		slog.Error("failed to list images", "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	slog.Info("images", "images", images)

	selectedImage := images[rand.Intn(len(images))]

	image, err := app.Cloudflare.GetImage(r.Context(), nil, selectedImage.ID)

	if err != nil {
		slog.Error("failed to get image", "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	imageByte, err := json.Marshal(image)
	if err != nil {
		slog.Error("failed to marshal image", "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(imageByte)
}

//go:embed openapi.yaml
var swaggerFs embed.FS

func (app *Application) Routes() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/v1/images", app.handlerPostImage).Methods("POST")
	router.HandleFunc("/v1/images", app.handlerGetImage).Methods("GET")

	opts := middleware.SwaggerUIOpts{SpecURL: "openapi.yaml"}
	sh := middleware.SwaggerUI(opts, nil)
	router.Handle("/docs", sh)
	router.Handle("/openapi.yaml", http.FileServer(http.FS(swaggerFs)))
	return router

}
