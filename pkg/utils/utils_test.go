package utils

import (
	"testing"
)

func TestParseCmd(t *testing.T) {
	var cmd1 = []string{"filename=wd-es.yaml", "msg=Gift", "from", "me"}
	var cmd2 = []string{"msg=Gift", "from", "me", "filename=wd-es.yaml"}
	var cmd3 = []string{"filename=wd-es.yaml"}
	var cmdArray = [][]string{cmd1, cmd2, cmd3}

	for _, c := range cmdArray {
		results := preParseUserOption(c)
		t.Log(results)
	}
}
