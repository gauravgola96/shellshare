package utils

import (
	"encoding/json"
	"net/http"
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
