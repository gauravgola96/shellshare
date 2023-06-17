package shellshare

import (
	"fmt"
	"github.com/TwiN/go-color"
	"github.com/enescakir/emoji"
	"github.com/gliderlabs/ssh"
	"github.com/gofrs/uuid"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	t "githug.com/gauravgola96/shellshare/pkg/tunnel"
	"io"
	"strings"
)

func HandleSSHSession(s ssh.Session) {
	subLogger := log.With().Str("module", "ssh_handler.HandleSShRequest").Logger()
	uid, err := uuid.NewV7()
	if err != nil || len(uid.Bytes()) == 0 {
		s.Write([]byte(BuildDownloadErrorStr()))
		subLogger.Error().Err(err).Msg("Error in new uuid")
		return
	}
	subLogger.Debug().Msgf("Tunnel Id : %s", uid.String())
	t.Tunnel.Store(uid.String(), make(chan t.SSHTunnel))

	address := fmt.Sprintf("%s:%d", viper.GetString("http.hostname"), viper.GetInt("http.port"))

	s.Write([]byte(BuildDownloadLinkStr(address, uid.String())))

	tunnel := <-t.Tunnel.GetWaitTunnel(uid.String())
	subLogger.Debug().Msgf("Tunnel ready : %s", uid.String())

	_, err = io.Copy(tunnel.W, s)
	if err != nil {
		s.Write([]byte(BuildDownloadErrorStr()))
		subLogger.Error().Err(err).Msg("Error in session writer")
		return
	}
	close(tunnel.Done)
	s.Write([]byte(BuildDownloadFinisedStr()))
}

func BuildDownloadLinkStr(address string, id string) string {
	var msg strings.Builder
	msg.WriteString("Your download link ")
	msg.WriteString(fmt.Sprintf("%s ", emoji.Parse(":eyes:")))
	msg.WriteString(fmt.Sprintf(color.Ize(color.Green, fmt.Sprintf("http://%s/download/%s", address, id))))
	return msg.String()
}

func BuildDownloadFinisedStr() string {
	var msg strings.Builder
	msg.WriteString("\n \n")
	msg.WriteString(fmt.Sprintf("%s ", emoji.Parse(":sunglasses:")))
	msg.WriteString("We are done !!! ")
	msg.WriteString(fmt.Sprintf("%s ", emoji.Parse(":tada:")))
	return msg.String()
}

func BuildDownloadErrorStr() string {
	var msg strings.Builder
	msg.WriteString("\n \n")
	msg.WriteString("Sorry something wend wrong! ")
	msg.WriteString(fmt.Sprintf("%s ", emoji.Parse(":face_with_head_bandage:")))
	return msg.String()
}
