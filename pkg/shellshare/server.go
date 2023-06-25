package shellshare

import (
	"context"
	"fmt"
	"github.com/gliderlabs/ssh"
	"github.com/go-chi/chi"
	_ "github.com/go-oauth2/oauth2/v4/server"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	auth2 "githug.com/gauravgola96/shellshare/pkg/authentication"
	"githug.com/gauravgola96/shellshare/pkg/middleware"
	"githug.com/gauravgola96/shellshare/pkg/storage"
	gossh "golang.org/x/crypto/ssh"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

//var (
//	DeadlineTimeout = 30 * time.Second
//	IdleTimeout     = 10 * time.Second
//)

func ServerAll() error {
	subLogger := log.With().Str("module", "shellshare.ServerAll").Logger()

	//storage initialize
	err := storage.Initialize()
	if err != nil {
		subLogger.Error().Err(err).Msgf("Error in storage initialization")
		return err
	}

	//SSH Server
	ssh.Handle(HandleSSHSession)
	sshAddr := fmt.Sprintf("%s:%d", viper.GetString("ssh.hostname"), viper.GetInt("ssh.port"))

	//Adding private keys so SSH don't recreate new private keys on every restart
	b, err := ioutil.ReadFile("private.pem")
	if err != nil {
		return err
	}
	key, err := gossh.ParsePrivateKey(b)
	if err != nil {
		return err
	}
	sshServer := &ssh.Server{
		Addr: sshAddr,
		//MaxTimeout:  DeadlineTimeout,
		//IdleTimeout: IdleTimeout,
		HostSigners: []ssh.Signer{key},
		PublicKeyHandler: func(ctx ssh.Context, key ssh.PublicKey) bool {
			return true
		},
	}

	go func() {
		subLogger.Info().Msgf("Listening SSH Server on %s", sshAddr)
		err := sshServer.ListenAndServe()
		if err != nil {
			subLogger.Error().Err(err).Msgf("Error in sshServer start")
			return
		}
	}()

	//HTTP
	err = auth2.Initialize(context.TODO())
	if err != nil {
		subLogger.Error().Err(err).Msgf("Error in authentication initialization")
		return err
	}

	addr := fmt.Sprintf("%s:%d", viper.GetString("http.hostname"), viper.GetInt("http.port"))
	router := chi.NewRouter()
	middleware.DefaultMiddleware(router)

	// setup ui
	workDir, _ := os.Getwd()
	filesDir := filepath.Join(workDir, "frontend-working")
	fileServer(router, "/", http.Dir(filesDir))

	// setup routes
	authRoutes, _ := auth2.Auth.Service.Handlers()
	router.Mount("/auth", authRoutes)
	router.Mount("/v1", HttpRoutes())

	httpServer := http.Server{
		Addr:              addr,
		Handler:           router,
		ReadHeaderTimeout: 20 * time.Second,
		ReadTimeout:       20 * time.Minute,
		WriteTimeout:      20 * time.Minute,
	}

	go func() {
		subLogger.Info().Msgf("Listening http Server on %s", addr)
		err := httpServer.ListenAndServe()
		if err != nil {
			subLogger.Error().Err(err).Msgf("Error in sshServer start")
			return
		}
	}()

	ShutDown(subLogger, "SSH & HTTP")
	return nil
}

func SSHServer() error {
	subLogger := log.With().Str("module", "shellshare.sshserver").Logger()

	//SSH Server
	ssh.Handle(HandleSSHSession)
	sshAddr := fmt.Sprintf("%s:%d", viper.GetString("ssh.hostname"), viper.GetInt("ssh.port"))

	//Adding private keys so SSH don't recreate new private keys on every restart
	b, err := ioutil.ReadFile("private.pem")
	if err != nil {
		return err
	}
	key, err := gossh.ParsePrivateKey(b)
	if err != nil {
		return err
	}
	sshServer := &ssh.Server{
		Addr: sshAddr,
		//MaxTimeout:  DeadlineTimeout,
		//IdleTimeout: IdleTimeout,
		HostSigners: []ssh.Signer{key},
	}

	go func() {
		subLogger.Info().Msgf("Listening SSH Server on %s", sshAddr)
		err := sshServer.ListenAndServe()
		if err != nil {
			subLogger.Error().Err(err).Msgf("Error in sshServer start")
			return
		}
	}()

	ShutDown(subLogger, "HTTP")
	return nil
}

func HttpServer() error {
	subLogger := log.With().Str("module", "shellshare.httpserver").Logger()

	addr := fmt.Sprintf("%s:%d", viper.GetString("http.hostname"), viper.GetInt("http.port"))
	router := chi.NewRouter()
	middleware.DefaultMiddleware(router)
	err := auth2.Initialize(context.TODO())
	if err != nil {
		subLogger.Error().Err(err).Msgf("Error in authentication initialization")
		return err
	}
	middleware.DefaultMiddleware(router)

	// setup ui
	workDir, _ := os.Getwd()
	filesDir := filepath.Join(workDir, "frontend")
	fileServer(router, "/", http.Dir(filesDir))

	// setup routes
	authRoutes, _ := auth2.Auth.Service.Handlers()
	router.Mount("/auth", authRoutes)
	router.Mount("/v1", HttpRoutes())

	httpServer := http.Server{
		Addr:              addr,
		Handler:           router,
		ReadHeaderTimeout: 20 * time.Second,
		ReadTimeout:       1 * time.Minute,
		WriteTimeout:      2 * time.Minute,
	}
	go func() {
		subLogger.Info().Msgf("Listening http Server on %s", addr)
		err := httpServer.ListenAndServe()
		if err != nil {
			subLogger.Error().Err(err).Msgf("Error in sshServer start")
			return
		}
	}()
	ShutDown(subLogger, "SSH")
	return nil

}

func ShutDown(l zerolog.Logger, service string) {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	sig := <-quit
	l.Info().Msgf("%s Server shutdown due to %s", service, sig)
}

// FileServer conveniently sets up a http.FileServer handler to serve static files from a http.FileSystem.
// Borrowed from https://github.com/go-chi/chi/blob/master/_examples/fileserver/main.go
func fileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	log.Printf("[INFO] serving static files from %v", root)
	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", http.StatusMovedPermanently).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	})
}
