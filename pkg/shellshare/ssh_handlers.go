package shellshare

import (
	"fmt"
	"github.com/gliderlabs/ssh"
	"io"
)

func HandleSShRequest(s ssh.Session) {
	io.WriteString(s, fmt.Sprintf("Hello World %s", s.User()))
}
