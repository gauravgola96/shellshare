package shellshare

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-pkgz/auth/token"
	"github.com/rs/zerolog/log"
	auth2 "githug.com/gauravgola96/shellshare/pkg/authentication"
	st "githug.com/gauravgola96/shellshare/pkg/storage"
	t "githug.com/gauravgola96/shellshare/pkg/tunnel"
	"githug.com/gauravgola96/shellshare/pkg/utils"
	"net/http"
	"os"
)

func HttpRoutes() *chi.Mux {
	r := chi.NewRouter()
	m := auth2.Auth.Service.Middleware()

	r.Get("/health", HandleHealthCheck)
	r.Get("/download/{id}", HandleDirectDownload)
	r.Get("/redirect/download/{id}", HandleRedirectDownload)
	r.Route("/user/", func(r chi.Router) {
		r.Use(m.Auth)

		r.Get("/info", HandleUserInfo)
		r.Get("/register", HandleRegisterUser)
	})
	return r
}

func InternalRoutes() *chi.Mux {
	r := chi.NewRouter()
	m := auth2.Auth.Service.Middleware()
	r.Use(m.AdminOnly)
	r.Route("/users", func(r chi.Router) {
		r.Get("/", HandleUserList)
	})
	return r
}

func HandleHealthCheck(w http.ResponseWriter, r *http.Request) {
	utils.WriteJson(w, http.StatusOK, "Healthy Upstream !!!", nil, utils.ResponseVar{Key: "Version", Val: os.Getenv("Version")})
}

func HandleDirectDownload(w http.ResponseWriter, r *http.Request) {
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

func HandleUserInfo(w http.ResponseWriter, r *http.Request) {
	subLogger := log.With().Str("module", "ssh_handler.HandleSShRequest").Logger()
	userInfo, err := token.GetUserInfo(r)
	if err != nil {
		subLogger.Error().Err(err).Msg("Error in get user info")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	utils.WriteJson(w, http.StatusOK, "successfully fetched user info", nil, utils.ResponseVar{"user_info", userInfo})
}

func HandleRedirectDownload(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	_, startTime, err := st.S.Cache.Get(id)
	if err == st.ErrNilCache {
		utils.WriteJson(w, http.StatusNotFound, fmt.Sprintf("Download is either completed or timed out"), nil)
		return
	} else if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, "something went wrong", err, utils.ResponseVar{
			Key: "Id",
			Val: id,
		})
		return
	}

	address := utils.GetHostAddress()
	downloadLink := fmt.Sprintf("%s/v1/download/%s", address, id)

	utils.WriteJson(w, http.StatusOK, "successfully fetched download details", nil, utils.ResponseVar{"download_link", downloadLink},
		utils.ResponseVar{"start_time", startTime})
}

func HandleRegisterUser(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
}

func HandleUserList(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	utils.WriteJson(w, http.StatusOK, "successfully fetched user list", nil, utils.ResponseVar{})

}
