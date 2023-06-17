package shellshare

import (
	"archive/zip"
	"fmt"
	"github.com/TwiN/go-color"
	"github.com/enescakir/emoji"
	"github.com/gliderlabs/ssh"
	"github.com/gofrs/uuid"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	t "githug.com/gauravgola96/shellshare/pkg/tunnel"
	"githug.com/gauravgola96/shellshare/pkg/utils"
	"io"
	"os"
	"strings"
)

func HandleSSHSession(s ssh.Session) {
	subLogger := log.With().Str("module", "ssh_handler.HandleSShRequest").Logger()
	uid, err := uuid.NewV7()
	if err != nil || len(uid.Bytes()) == 0 {
		s.Write([]byte(BuildDownloadErrorStr(nil)))
		subLogger.Error().Err(err).Msg("Error in new uuid")
		return
	}
	subLogger.Debug().Msgf("Tunnel Id : %s", uid.String())
	t.Tunnel.Store(uid.String(), make(chan t.SSHTunnel))

	address := fmt.Sprintf("%s:%d", viper.GetString("http.hostname"), viper.GetInt("http.port"))

	option, err := utils.ParseUserOption(s.Command())
	if err != nil {
		s.Write([]byte(BuildDownloadErrorStr(err)))
		subLogger.Error().Err(err).Msg("Error in user options")
		return
	}

	s.Write([]byte(BuildDownloadLinkStr(address, uid.String())))

	tunnel := <-t.Tunnel.GetWaitTunnel(uid.String())
	defer func() {
		close(tunnel.Done)
		s.Close()
	}()

	subLogger.Debug().Msgf("Tunnel ready : %s", uid.String())

	err = ZipAndWriteFile(option.FileName, tunnel.W, s)
	if err != nil {
		s.Write([]byte(BuildDownloadErrorStr(nil)))
		subLogger.Error().Err(err).Msg("Error in session writer")
		return
	}

	s.Write([]byte(BuildDownloadFinishedStr()))
}

func BuildDownloadLinkStr(address string, id string) string {
	var msg strings.Builder
	msg.WriteString("\n \n")
	msg.WriteString("Your download link ")
	msg.WriteString(fmt.Sprintf("%s ", emoji.Parse(":eyes:")))
	msg.WriteString(fmt.Sprintf(color.Ize(color.Green, fmt.Sprintf("http://%s/download/%s", address, id))))
	return msg.String()
}

func BuildDownloadFinishedStr() string {
	var msg strings.Builder
	msg.WriteString("\n \n")
	msg.WriteString(fmt.Sprintf("%s ", emoji.Parse(":sunglasses:")))
	msg.WriteString("We are done !!! ")
	msg.WriteString(fmt.Sprintf("%s ", emoji.Parse(":tada:")))
	return msg.String()
}

func BuildDownloadErrorStr(err error) string {
	var msg strings.Builder
	msg.WriteString("\n \n")
	msg.WriteString(fmt.Sprintf(color.Ize(color.Red, "Sorry something wend wrong!")))
	if err != nil {
		msg.WriteString("\n")
		msg.WriteString(fmt.Sprintf("%s %s ", err.Error(), emoji.Parse(":cold_sweat:")))
		return msg.String()
	}
	msg.WriteString(fmt.Sprintf("%s ", emoji.Parse(":face_with_head_bandage:")))
	return msg.String()
}

// ZipAndWriteFile Zip file and update it to io.writer
func ZipAndWriteFile(filename string, w io.Writer, r io.Reader) error {
	subLogger := log.With().Str("module", "ssh_handler.ZipAndWriteFile").Logger()

	if filename == "" {
		filename = utils.RandomString(8)
	}

	f, err := os.Create(filename)
	if err != nil {
		subLogger.Error().Err(err).Msg("Error in file creation")
		return err
	}

	defer func() {
		if err := f.Close(); err != nil {
			subLogger.Error().Err(err).Msg("Error in file closure")
		}
		os.Remove(filename)
	}()

	// write straight to the http.ResponseWriter
	zw := zip.NewWriter(w)
	cf, err := zw.Create(f.Name())
	if err != nil {
		subLogger.Error().Err(err).Msg("Error in zip create")
		return err
	}

	_, err = io.Copy(cf, r)
	if err != nil {
		subLogger.Error().Err(err).Msg("Error in copy reader")
		return err
	}
	err = zw.Close()
	if err != nil {
		subLogger.Error().Err(err).Msg("Error in zip closure")
		return err
	}
	return nil
}
