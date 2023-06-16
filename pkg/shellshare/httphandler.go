package shellshare

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"net/http"
)

func HttpRoutes() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/health", Health)
	r.Get("/{id}", DownloadFile)
	return r
}

func Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":     "Healthy upstream !!!",
		"status_code": http.StatusOK,
	})
}

func DownloadFile(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Healthy upstream !!!"))
}
