package config

import (
	"bytes"
	"errors"
	"io"
	"net/http"
)

func DoPostRequest(url string, jsonData string) ([]byte, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(jsonData)))
	if err != nil {
		return []byte{}, err
	}
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	}
	if res.StatusCode < 200 || res.StatusCode > 200 {
		if res.StatusCode == 401 {
			return []byte{}, errors.New("PAYMENT_NOT_COMPLETE")
		}
		return []byte{}, errors.New("status code not in range 200")
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return []byte{}, err
	}
	return body, nil
}
