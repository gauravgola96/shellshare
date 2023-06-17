package utils

import (
	"errors"
	"fmt"
	"strings"
)

const (
	Filename     string = "filename"
	Message      string = "message"
	maxMsgLength int    = 15
)

type UserOptions struct {
	FileName string
	Message  string
}

func preParseUserOption(cmd []string) []string {
	var options []string
	var previousCmd string
	for _, c := range cmd {
		if ok := strings.Contains(c, "="); ok {
			if previousCmd != "" {
				options = append(options, previousCmd)
			}
			previousCmd = c
		} else {
			previousCmd = previousCmd + " " + c
		}
	}
	if previousCmd != "" {
		options = append(options, previousCmd)
	}

	return options
}

func ParseUserOption(cmd []string) (*UserOptions, error) {
	userOpts := &UserOptions{}
	option := preParseUserOption(cmd)
	for _, o := range option {
		parts := strings.Split(o, "=")
		if len(parts) < 2 {
			return nil, errors.New("invalid options")
		}

		key := parts[0]
		value := parts[1]
		switch key {
		case Filename:
			userOpts.FileName = value
		case Message:
			if len(value) > maxMsgLength {
				return nil, errors.New(fmt.Sprintf("message size is more than %d characters", maxMsgLength))
			}
			userOpts.Message = value
		}
	}
	return userOpts, nil
}
