package shellshare

import (
	"fmt"
	"github.com/gliderlabs/ssh"
	"github.com/gofrs/uuid"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	t "githug.com/gauravgola96/shellshare/pkg/tunnel"
	"io"
)

func HandleSSHSession(s ssh.Session) {
	subLogger := log.With().Str("module", "ssh_handler.HandleSShRequest").Logger()
	uid, err := uuid.NewV7()
	if err != nil || len(uid.Bytes()) == 0 {
		subLogger.Error().Err(err).Msg("Error in new uuid")
		return
	}
	subLogger.Debug().Msgf("Tunnel Id : %s", uid.String())
	t.Tunnel.Store(uid.String(), make(chan t.SSHTunnel))

	address := fmt.Sprintf("%s:%d", viper.GetString("http.hostname"), viper.GetInt("http.port"))

	s.Write([]byte(fmt.Sprintf("\n Your download link: %s", fmt.Sprintf("http://%s/download/%s", address, uid.String()))))

	tunnel := <-t.Tunnel.GetWaitTunnel(uid.String())
	subLogger.Debug().Msgf("Tunnel ready : %s", uid.String())

	_, err = io.Copy(tunnel.W, s)
	if err != nil {
		subLogger.Error().Err(err).Msg("Error in session writer")
		return
	}
	close(tunnel.Done)
	s.Write([]byte("\n File transfer completed !!!"))
}
