package main

import (
	"bytes"
	"net/http"
)

func relay(path string, body []byte) error {

	resp, err := http.Post("http://127.0.0.1:8080/"+path, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
