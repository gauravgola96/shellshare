package utils

import (
	"encoding/json"
	"fmt"
	"github.com/TwiN/go-color"
	"github.com/enescakir/emoji"
	"github.com/spf13/viper"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

const (
	MaxTimoutMinutes = 15
	MaxCacheTTL      = 20
	MaxBytesSize     = 2147483648 //2 GB
)

func WriteJson(w http.ResponseWriter, status int, message any, err error, rvars ...ResponseVar) {

	resp := map[string]interface{}{}
	if err != nil {
		resp["error"] = err.Error()
	}
	resp["message"] = message
	resp["status"] = status
	for _, v := range rvars {
		resp[v.Key] = v.Val
	}

	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Connection", "keep-alive")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(resp)
}

type ResponseVar struct {
	Key string
	Val any
}

func RandomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length+2)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[2 : length+2]
}

func GetHostAddress() string {
	address := fmt.Sprintf("https://%s", viper.GetString("http.remote_host"))
	if !viper.GetBool("production") {
		address = fmt.Sprintf("http://%s:%d", viper.GetString("http.hostname"), viper.GetInt("http.port"))
	}
	return address
}

func BuildDownloadLinkStr(address string, id string, timeout int) string {
	var msg strings.Builder
	msg.WriteString("\n")
	msg.WriteString(fmt.Sprintf(color.Ize(color.Red, fmt.Sprintf("WARNING : Only file size upto %dGB is allowed !!! ", MaxBytesSize/1024/1024/1024))))
	msg.WriteString(fmt.Sprintf("%s ", emoji.Parse(":warning: ")))
	msg.WriteString("\n")
	msg.WriteString("Your download link ")
	msg.WriteString(fmt.Sprintf("%s ", emoji.Parse(":eyes: :")))
	msg.WriteString(fmt.Sprintf(color.Ize(color.Green, fmt.Sprintf("%s/v1/redirect/download/%s", address, id))))
	msg.WriteString("\n \n \n")
	msg.WriteString(fmt.Sprintf(color.Ize(color.Yellow, "Please don't kill this session \n")))
	msg.WriteString(fmt.Sprintf("Your link will expire in %d minutes, If not used %s \n", timeout, emoji.Parse(":hugging_face:")))
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
	msg.WriteString(fmt.Sprintf(color.Ize(color.Red, "Sorry something went wrong!")))
	if err != nil {
		msg.WriteString("\n")
		msg.WriteString(fmt.Sprintf("%s %s ", err.Error(), emoji.Parse(":cold_sweat:")))
		return msg.String()
	}
	msg.WriteString(fmt.Sprintf("%s ", emoji.Parse(":face_with_head_bandage:")))
	return msg.String()
}

func BuildCloseSessionTimeoutStr() string {
	var msg strings.Builder
	msg.WriteString("\n \n")
	msg.WriteString("Closing session due to timeout ")
	msg.WriteString(fmt.Sprintf("%s ", emoji.Parse(":disappointed: ")))
	return msg.String()
}
