package api

import (
	"embed"
	"encoding/json"
	"math/rand"
	"net/http"

	"github.com/cloudflare/cloudflare-go"
	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
	"golang.org/x/exp/slog"
)

type (
	Application struct {
		CloudFlare *cloudflare.API
	}

	Image struct {
		Id       string   `json:"id"`
		FileName string   `json:"filename"`
		Variants []string `json:"variants"`
	}
)

func handlerTest(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Jojo is Awesome!"))
}

func (app *Application) handlerPostImage(w http.ResponseWriter, r *http.Request) {
	slog.Info("handlePostImage")
	err := r.ParseMultipartForm(32 << 20) // 32MB is the maximum size of a file we can upload
	if err != nil {
		slog.Error("failed to parse multipart form", "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	mForm := r.MultipartForm
	for k := range mForm.File {
		file, fileHeader, err := r.FormFile(k)
		if err != nil {
			slog.Error("failed to get image from form", "error", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer file.Close()

		img, err := app.CloudFlare.UploadImage(r.Context(), &cloudflare.ResourceContainer{
			Identifier: "203752570d3d905ee071d7857cc2989d", // JOJO: Remove hardcode
			Level:      cloudflare.AccountRouteLevel,
		}, cloudflare.UploadImageParams{
			File: file,
			Name: fileHeader.Filename,
			Metadata: map[string]interface{}{
				"upload": "api",
			},
		})
		if err != nil {
			slog.Error("failed to upload image", "error", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		slog.Info("uploaded image", "id", img.ID, "filename", img.Filename)

		imgByte, err := json.Marshal(Image{
			Id:       img.ID,
			FileName: img.Filename,
			Variants: img.Variants,
		})
		if err != nil {
			slog.Error("failed to marshal image", "error", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(imgByte)
	}
}

func (app *Application) handlerGetRandomImage(w http.ResponseWriter, r *http.Request) {
	slog.Info("handleGetImage")

	imgs, err := app.CloudFlare.ListImages(r.Context(), &cloudflare.ResourceContainer{
		Identifier: "203752570d3d905ee071d7857cc2989d", // TODO(JOJO): Remove hardcode
		Level:      cloudflare.AccountRouteLevel,
	}, cloudflare.ListImagesParams{})

	if err != nil {
		slog.Error("failed to list images", "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	randomImage := imgs[rand.Intn(len(imgs))]

	imgByte, err := json.Marshal(Image{
		Id:       randomImage.ID,
		FileName: randomImage.Filename,
		Variants: randomImage.Variants,
	})
	if err != nil {
		slog.Error("failed to marshal image", "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(imgByte)
}

//go:embed openapi.yaml
var swaggerFs embed.FS

func (app *Application) Routes() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/v1/images", app.handlerPostImage).Methods("POST")
	router.HandleFunc("/v1/images", app.handlerGetRandomImage).Methods("GET")
	router.HandleFunc("/test", handlerTest).Methods("GET")
	opts := middleware.SwaggerUIOpts{SpecURL: "openapi.yaml"}
	sh := middleware.SwaggerUI(opts, nil)
	router.Handle("/docs", sh)
	router.Handle("/openapi.yaml", http.FileServer(http.FS(swaggerFs)))
	return router

}
