package shellshare

import (
	"archive/zip"
	"context"
	"errors"
	"fmt"
	"github.com/gliderlabs/ssh"
	"github.com/gofrs/uuid"
	"github.com/rs/zerolog/log"
	st "githug.com/gauravgola96/shellshare/pkg/storage"
	t "githug.com/gauravgola96/shellshare/pkg/tunnel"
	"githug.com/gauravgola96/shellshare/pkg/utils"
	gossh "golang.org/x/crypto/ssh"
	"io"
	"os"
	"time"
)

func HandleSSHSession(s ssh.Session) {
	subLogger := log.With().Str("module", "ssh_handler.HandleSShRequest").Logger()

	authorizedKey := gossh.MarshalAuthorizedKey(s.PublicKey())
	subLogger.Info().Msgf("SSH request from %s : %s", s.User(), authorizedKey)
	uid, err := uuid.NewV7()
	if err != nil || len(uid.Bytes()) == 0 {
		s.Write([]byte(utils.BuildDownloadErrorStr(nil)))
		subLogger.Error().Err(err).Msg("Error in new uuid")
		return
	}
	subLogger.Debug().Msgf("Tunnel Id : %s", uid.String())
	t.Tunnel.Store(uid.String(), make(chan t.ConnectionTunnel))

	address := utils.GetHostAddress()
	option, err := utils.ParseUserOption(s.Command())
	if err != nil {
		s.Write([]byte(utils.BuildDownloadErrorStr(err)))
		subLogger.Error().Err(err).Msg("Error in user options")
		return
	}
	//store in cache
	st.S.Cache.Put(uid.String(), st.ValueItem{FileName: option.FileName, Message: option.Message}, utils.MaxCacheTTLMinutes*time.Minute)
	defer st.S.Cache.Delete(uid.String())

	s.Write([]byte(utils.BuildDownloadLinkStr(address, uid.String(), utils.MaxTimoutSSHMinutes)))

	var (
		written int64
		status  st.Status = st.Success
	)
	defer func() {
		err := st.UpdateDownloadDetail(context.TODO(), st.Download{
			SSHKeys:      string(authorizedKey),
			Source:       st.SourceSSH,
			BytesWritten: written,
			Status:       status,
		})
		if err != nil {
			subLogger.Error().Err(err).Msg("Error in mongo update")
		}
	}()

	ticker := time.NewTicker(utils.MaxTimoutSSHMinutes * time.Minute)
	for {
		select {
		case <-s.Context().Done():
			status = st.SessionClose
			subLogger.Info().Msg("Session closed from client")
			return

		case <-ticker.C:
			subLogger.Info().Msg("Session timeout")
			status = st.Timeout
			s.Write([]byte(utils.BuildCloseSessionTimeoutStr()))
			t.Tunnel.Delete(uid.String())
			s.Close()
			return

		case tunnel := <-t.Tunnel.GetWaitTunnel(uid.String()):
			defer func() {
				close(tunnel.Done)
				s.Close()
			}()

			subLogger.Debug().Msgf("Tunnel ready : %s", uid.String())

			written, err = ZipAndWriteFile(option.FileName, tunnel.W, s)
			if err != nil {
				status = st.Failed
				s.Write([]byte(utils.BuildDownloadErrorStr(err)))
				subLogger.Error().Err(err).Msg("Error in session writer")
				return
			}

			s.Write([]byte(utils.BuildDownloadFinishedStr()))
			return
		default:
			//pass
		}
	}
}

// ZipAndWriteFile Zip file and update it to io.writer
func ZipAndWriteFile(filename string, w io.Writer, r io.Reader) (int64, error) {
	subLogger := log.With().Str("module", "ssh_handler.ZipAndWriteFile").Logger()

	if filename == "" {
		filename = utils.RandomString(8)
	}

	f, err := os.Create(filename)
	if err != nil {
		subLogger.Error().Err(err).Msg("Error in file creation")
		return 0, err
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
		return 0, err
	}

	//_, err = io.Copy(cf, r)
	var written int64
	written, err = CopyBuffer(cf, r, nil)
	if err != nil {
		subLogger.Error().Err(err).Msg("Error in copy reader")
		return 0, err
	}
	err = zw.Close()
	if err != nil {
		subLogger.Error().Err(err).Msg("Error in zip closure")
		return 0, err
	}
	return written, nil
}

// CopyBuffer implementation of io.Copy with max byte limit check
func CopyBuffer(dst io.Writer, src io.Reader, buf []byte) (written int64, err error) {
	// If the reader has a WriteTo method, use it to do the copy.
	// Avoids an allocation and a copy.
	if wt, ok := src.(io.WriterTo); ok {
		return wt.WriteTo(dst)
	}
	// Similarly, if the writer has a ReadFrom method, use it to do the copy.
	if rt, ok := dst.(io.ReaderFrom); ok {
		return rt.ReadFrom(src)
	}
	if buf == nil {
		size := 32 * 1024
		if l, ok := src.(*io.LimitedReader); ok && int64(size) > l.N {
			if l.N < 1 {
				size = 1
			} else {
				size = int(l.N)
			}
		}
		buf = make([]byte, size)
	}
	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			if nw < 0 || nr < nw {
				nw = 0
				if ew == nil {
					ew = errors.New("invalid write result")
				}
			}
			written += int64(nw)
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}

			if written > utils.MaxBytesSize {
				return written, errors.New(fmt.Sprintf("File size cannot be more than %dGB", utils.MaxBytesSize/1024/1024/1024))
			}

		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
	}
	return written, err
}
