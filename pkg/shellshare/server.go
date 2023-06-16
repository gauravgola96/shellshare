package shellshare

import (
	"fmt"
	"github.com/gliderlabs/ssh"
	"github.com/go-chi/chi"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"githug.com/gauravgola96/shellshare/pkg/middleware"
	gossh "golang.org/x/crypto/ssh"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	DeadlineTimeout = 30 * time.Second
	IdleTimeout     = 10 * time.Second
)

func SSHServer() error {
	subLogger := log.With().Str("module", "shellshare.sshserver").Logger()

	//SSH Server
	ssh.Handle(HandleSShRequest)
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
		Addr:        sshAddr,
		MaxTimeout:  DeadlineTimeout,
		IdleTimeout: IdleTimeout,
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

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	sig := <-quit
	subLogger.Info().Msgf("SSH Server shutdown due to %s", sig)

	return nil
}

func HttpServer() error {
	subLogger := log.With().Str("module", "shellshare.httpserver").Logger()

	addr := fmt.Sprintf("%s:%d", viper.GetString("http.hostname"), viper.GetInt("http.port"))
	router := chi.NewRouter()
	middleware.DefaultMiddleware(router)

	httpServer := http.Server{
		Addr:              addr,
		Handler:           router,
		ReadHeaderTimeout: 20 * time.Second,
		ReadTimeout:       1 * time.Minute,
		WriteTimeout:      2 * time.Minute,
	}
	router.Mount("/", HttpRoutes())

	go func() {
		subLogger.Info().Msgf("Listening http Server on %s", addr)
		err := httpServer.ListenAndServe()
		if err != nil {
			subLogger.Error().Err(err).Msgf("Error in sshServer start")
			return
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	sig := <-quit
	subLogger.Info().Msgf("HTTP Server shutdown due to %s", sig)

	return nil

}
