package authentication

import (
	"context"
	"github.com/go-pkgz/auth"
	"github.com/go-pkgz/auth/avatar"
	"github.com/go-pkgz/auth/token"
	log "github.com/go-pkgz/lgr"
	"github.com/spf13/viper"
	"githug.com/gauravgola96/shellshare/pkg/utils"
	"strings"
	"time"
)

type AuthHandler struct {
	Service *auth.Service
}

var Auth AuthHandler

func Initialize(ctx context.Context) error {

	addr := utils.GetHostAddress()

	log.Setup(log.Debug, log.Msec, log.LevelBraces, log.CallerFile, log.CallerFunc) // setup default logger with go-pkgz/lgr

	options := auth.Opts{
		SecretReader: token.SecretFunc(func(_ string) (string, error) { // secret key for JWT, ignores aud
			return "secret", nil
		}),
		TokenDuration:     time.Minute,                                 // short token, refreshed automatically
		CookieDuration:    time.Hour * 24,                              // cookie fine to keep for long time
		DisableXSRF:       false,                                       // don't disable XSRF in real-life applications!
		Issuer:            "shellshare-service",                        // part of token, just informational
		URL:               addr,                                        // base url of the protected Service
		AvatarStore:       avatar.NewLocalFS("/tmp/demo-auth-service"), // stores avatars locally
		AvatarResizeLimit: 200,
		ClaimsUpd: token.ClaimsUpdFunc(func(claims token.Claims) token.Claims { // modify issued token
			if claims.User != nil {
				adminUsers := viper.GetStringSlice("admins")
				if utils.ItemExists(adminUsers, claims.User.ID) {
					// set attributes for admin
					claims.User.SetAdmin(true)
				}
			}
			return claims
		}), // resizes avatars to 200x200
		Validator: token.ValidatorFunc(func(_ string, claims token.Claims) bool { // rejects some tokens
			if claims.User != nil {
				if strings.HasPrefix(claims.User.ID, "github_") { // allow all users with github authentication
					return true
				}
				if strings.HasPrefix(claims.User.ID, "google_") { // allow all users with github authentication
					return true
				}
				return false
			}
			return false
		}),
		Logger:      log.Default(), // optional logger for auth library
		UseGravatar: true,          // for verified provider use gravatar service
	}

	service := auth.NewService(options)
	service.AddProvider("github", viper.GetString("auth.github.client_id"), viper.GetString("auth.github.client_secret"))
	service.AddProvider("google", viper.GetString("auth.google.client_id"), viper.GetString("auth.google.client_secret"))

	Auth = AuthHandler{
		Service: service,
	}
	return nil
}
