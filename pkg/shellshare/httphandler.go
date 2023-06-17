package shellshare

import (
	"fmt"
	"github.com/go-chi/chi"
	t "githug.com/gauravgola96/shellshare/pkg/tunnel"
	"githug.com/gauravgola96/shellshare/pkg/utils"
	"net/http"
	"os"
)

func HttpRoutes() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/health", HandleHealthCheck)
	r.Get("/download/{id}", HandleDownloadFileFromLink)
	return r
}

func HandleHealthCheck(w http.ResponseWriter, r *http.Request) {
	utils.WriteJson(w, http.StatusOK, "Healthy Upstream !!!", nil, utils.ResponseVar{Key: "Version", Val: os.Getenv("Version")})
}

func HandleDownloadFileFromLink(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	tunnel, ok := t.Tunnel.Get(id)
	if !ok {
		utils.WriteJson(w, http.StatusNotFound, "Id not found", nil, utils.ResponseVar{
			Key: "Id",
			Val: id,
		})
	}
	defer t.Tunnel.Delete(id)
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.zip\"", "shellshare"))

	doneChan := make(chan struct{}, 1)

	tunnel <- t.SSHTunnel{W: w, Done: doneChan}
	<-doneChan
}
